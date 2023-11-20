package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/config"
	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterUser - Register User as a Member
func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			FullName string `json:"full_name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error while hashing password", http.StatusInternalServerError)
			return
		}

		// User role is hardcoded as 'member'
		user := models.User{
			FullName: requestBody.FullName,
			Email:    requestBody.Email,
			Password: string(hashedPassword),
			Role:     "member",
		}

		result := db.Create(&user)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		config.SendJSONResponse(w, map[string]interface{}{
			"id":         user.ID,
			"full_name":  user.FullName,
			"email":      user.Email,
			"created_at": user.CreatedAt.Format(time.RFC3339),
		})
	}
}

func InitializeAdminUser(db *gorm.DB) {
	adminPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	adminUser := models.User{
		FullName: "Admin User",
		Email:    "admin@gmail.com",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	// Check if admin exists
	var count int64
	db.Model(&models.User{}).Where("email = ?", adminUser.Email).Count(&count)
	if count == 0 {
		// Create admin if not exists
		db.Create(&adminUser)
	}
}

// Login User
func LoginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user models.User
		result := db.Where("email = ?", requestBody.Email).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "User not found", http.StatusUnauthorized)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(1 * time.Hour)
		claims := &config.Claims{
			Email: user.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(config.JwtKey)
		if err != nil {
			http.Error(w, "Error while signing the token", http.StatusInternalServerError)
			return
		}

		config.SendJSONResponse(w, map[string]string{
			"token": tokenString,
		})
	}
}

func UpdateUserAccount(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mengurai token dari header Authorization
		authHeader := r.Header.Get("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			http.Error(w, "Authorization header must be in the format 'Bearer {token}'", http.StatusUnauthorized)
			return
		}
		tokenString := authHeaderParts[1]

		// Parsing token JWT
		claims := &config.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Mendekode request body
		var requestBody struct {
			FullName string `json:"full_name"`
			Email    string `json:"email"`
		}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Memperbarui user
		user := models.User{}
		db.Model(&user).Where("email = ?", claims.Email).Updates(models.User{FullName: requestBody.FullName, Email: requestBody.Email})
		db.Where("email = ?", requestBody.Email).First(&user)

		// Mengirimkan respon
		config.SendJSONResponse(w, map[string]interface{}{
			"id":         user.ID,
			"full_name":  user.FullName,
			"email":      user.Email,
			"updated_at": user.UpdatedAt,
		})

	}
}

func DeleteUserAccount(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate the user
		authHeader := r.Header.Get("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			http.Error(w, "Authorization header must be in the format 'Bearer {token}'", http.StatusUnauthorized)
			return
		}
		tokenString := authHeaderParts[1]

		claims := &config.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check if there are any tasks associated with the user
		var taskCount int64
		if err := db.Model(&models.Task{}).Where("user_id = ?", claims.UserID).Count(&taskCount).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if taskCount > 0 {
			http.Error(w, "Cannot delete user with associated tasks", http.StatusConflict)
			return
		}

		// Delete the user
		result := db.Where("email = ?", claims.Email).Delete(&models.User{})
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Send a JSON response indicating successful deletion
		config.SendJSONResponse(w, map[string]string{
			"message": "Your account has been successfully deleted",
		})
	}
}
