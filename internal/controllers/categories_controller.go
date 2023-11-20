package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/config"
	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type CreateCategoryResponse struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateCategory creates a new category
func CreateCategory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Autentikasi dan autorisasi
		authorized, err := config.AuthenticateAndAuthorize(r, db)
		if err != nil || !authorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var category models.Category
		if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := db.Create(&category)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Create the response struct with only the required fields
		response := CreateCategoryResponse{
			ID:        category.ID,
			Type:      category.Type,
			CreatedAt: category.CreatedAt,
		}

		w.WriteHeader(http.StatusCreated)
		config.SendJSONResponse(w, response)
	}
}

type CategoryResponse struct {
	ID        uint       `json:"id"`
	Type      string     `json:"type"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
	Tasks     []TaskInfo `json:"Tasks"`
}

type TaskInfo struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetCategories(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate the user
		claims, err := config.Authenticate(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Check if the user is an admin
		var user models.User
		if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Get categories with tasks based on user role
		var categories []models.Category
		if user.Role == "admin" {
			// Admin can view all tasks
			result := db.Preload("Tasks").Find(&categories)
			if result.Error != nil {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// Regular user can view only their tasks
			result := db.Preload("Tasks", "user_id = ?", claims.UserID).Find(&categories)
			if result.Error != nil {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Create a slice of CategoryResponse without user information
		var responseCategories []CategoryResponse
		for _, category := range categories {
			var taskInfos []TaskInfo
			for _, task := range category.Tasks {
				taskInfo := TaskInfo{
					ID:          task.ID,
					Title:       task.Title,
					Description: task.Description,
					UserID:      task.UserID,
					CategoryID:  task.CategoryID,
					CreatedAt:   task.CreatedAt,
					UpdatedAt:   task.UpdatedAt,
				}

				taskInfos = append(taskInfos, taskInfo)
			}

			responseCategory := CategoryResponse{
				ID:        category.ID,
				Type:      category.Type,
				UpdatedAt: category.UpdatedAt,
				CreatedAt: category.CreatedAt,
				Tasks:     taskInfos,
			}
			responseCategories = append(responseCategories, responseCategory)
		}

		config.SendJSONResponse(w, responseCategories)
	}
}

// UpdateCategory updates an existing category
type UpdateCategoryResponse struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UpdateCategory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Autentikasi dan autorisasi
		authorized, err := config.AuthenticateAndAuthorize(r, db)
		if err != nil || !authorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		categoryID := vars["categoryId"]

		var updateData struct {
			Type string `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := db.Model(&models.Category{}).Where("id = ?", categoryID).Updates(updateData)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch the updated category
		var updatedCategory models.Category
		if err := db.First(&updatedCategory, categoryID).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create the response struct with only the required fields
		response := UpdateCategoryResponse{
			ID:        updatedCategory.ID,
			Type:      updatedCategory.Type,
			UpdatedAt: updatedCategory.UpdatedAt,
		}

		config.SendJSONResponse(w, response)
	}
}

// DeleteCategory deletes a category
func DeleteCategory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Autentikasi dan autorisasi
		authorized, err := config.AuthenticateAndAuthorize(r, db)
		if err != nil || !authorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		categoryID := vars["categoryId"]

		// Check if the category exists
		var category models.Category
		if err := db.First(&category, categoryID).Error; err != nil {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		// Check if there are associated tasks
		var tasks []models.Task
		if err := db.Where("category_id = ?", categoryID).Find(&tasks).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Delete associated tasks
		for _, task := range tasks {
			if err := db.Delete(&task).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Delete the category
		if err := db.Delete(&category).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		config.SendJSONResponse(w, map[string]string{"message": "Category and associated tasks have been successfully deleted"})
	}
}
