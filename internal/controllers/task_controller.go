package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/config"
	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type CreateTaskResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Status      bool      `json:"status"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func CreateTask(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Autentikasi pengguna
		claims, err := config.Authenticate(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if category exists
		var category models.Category
		if err := db.First(&category, task.CategoryID).Error; err != nil {
			http.Error(w, "Category not found", http.StatusBadRequest)
			return
		}

		task.UserID = claims.UserID // Set user ID from JWT claims
		task.Status = false         // Set status to false by default

		// Save task to database
		if err := db.Create(&task).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create the response struct with only the required fields
		response := CreateTaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Status:      task.Status,
			Description: task.Description,
			UserID:      task.UserID,
			CategoryID:  task.CategoryID,
			CreatedAt:   task.CreatedAt,
		}

		w.WriteHeader(http.StatusCreated)
		config.SendJSONResponse(w, response)
	}
}

type GetTasksResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Status      bool      `json:"status"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	User        struct {
		ID       uint   `json:"id"`
		Email    string `json:"email"`
		FullName string `json:"full_name"`
	} `json:"User"`
}

func GetTasks(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := config.Authenticate(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var tasks []models.Task
		if err := db.Where("user_id = ?", claims.UserID).Preload("User").Find(&tasks).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a slice for the custom response
		var response []GetTasksResponse

		// Map tasks to custom response struct
		for _, task := range tasks {
			response = append(response, GetTasksResponse{
				ID:          task.ID,
				Title:       task.Title,
				Status:      task.Status,
				Description: task.Description,
				UserID:      task.UserID,
				CategoryID:  task.CategoryID,
				CreatedAt:   task.CreatedAt,
				User: struct {
					ID       uint   `json:"id"`
					Email    string `json:"email"`
					FullName string `json:"full_name"`
				}{
					ID:       task.User.ID,
					Email:    task.User.Email,
					FullName: task.User.FullName,
				},
			})
		}

		config.SendJSONResponse(w, response)
	}
}

func UpdateTask(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := config.Authenticate(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		taskID, err := strconv.Atoi(vars["taskId"])
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		var updateData struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var task models.Task
		if err := db.First(&task, taskID).Error; err != nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		if task.UserID != claims.UserID {
			http.Error(w, "Unauthorized to update this task", http.StatusUnauthorized)
			return
		}

		// Update task with new data
		db.Model(&task).Updates(updateData)

		// Retrieve updated task data
		var updatedTask models.Task
		db.First(&updatedTask, taskID)

		// Create the response struct with only the required fields
		response := struct {
			ID          uint      `json:"id"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Status      bool      `json:"status"`
			UserID      uint      `json:"user_id"`
			CategoryID  uint      `json:"category_id"`
			UpdatedAt   time.Time `json:"updated_at"`
		}{
			ID:          updatedTask.ID,
			Title:       updatedTask.Title,
			Description: updatedTask.Description,
			Status:      updatedTask.Status,
			UserID:      updatedTask.UserID,
			CategoryID:  updatedTask.CategoryID,
			UpdatedAt:   updatedTask.UpdatedAt,
		}

		config.SendJSONResponse(w, response)
	}
}

type UpdateTaskStatusResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func UpdateTaskStatus(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := config.Authenticate(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		taskId, err := strconv.Atoi(vars["taskId"])
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		var updateData struct {
			Status bool `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var task models.Task
		if err := db.First(&task, taskId).Error; err != nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		if task.UserID != claims.UserID {
			http.Error(w, "Unauthorized to update status of this task", http.StatusUnauthorized)
			return
		}

		db.Model(&task).Update("status", updateData.Status)

		// Create the response struct with only the required fields
		response := UpdateTaskStatusResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			UserID:      task.UserID,
			CategoryID:  task.CategoryID,
			UpdatedAt:   task.UpdatedAt,
		}

		config.SendJSONResponse(w, response)
	}
}

func UpdateTaskCategory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := config.Authenticate(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		taskId, err := strconv.Atoi(vars["taskId"])
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		var updateData struct {
			CategoryID uint `json:"category_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var category models.Category
		if err := db.First(&category, updateData.CategoryID).Error; err != nil {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		var task models.Task
		if err := db.First(&task, taskId).Error; err != nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		// Check if the task belongs to the authenticated user
		if task.UserID != claims.UserID {
			http.Error(w, "Unauthorized to change the category of this task", http.StatusUnauthorized)
			return
		}

		// Update the category ID in the task
		db.Model(&task).Update("category_id", updateData.CategoryID)

		// Fetch the updated task with preloaded user data
		var updatedTask models.Task
		if err := db.Where("id = ?", taskId).Preload("User").First(&updatedTask).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create the response struct with the required fields
		response := struct {
			ID          uint      `json:"id"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Status      bool      `json:"status"`
			UserID      uint      `json:"user_id"`
			CategoryID  uint      `json:"category_id"`
			UpdatedAt   time.Time `json:"updated_at"`
		}{
			ID:          updatedTask.ID,
			Title:       updatedTask.Title,
			Description: updatedTask.Description,
			Status:      updatedTask.Status,
			UserID:      updatedTask.UserID,
			CategoryID:  updatedTask.CategoryID,
			UpdatedAt:   updatedTask.UpdatedAt,
		}

		// Send the JSON response
		config.SendJSONResponse(w, response)
	}
}
func DeleteTask(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate the user
		claims, err := config.Authenticate(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Extract task ID from the request parameters
		vars := mux.Vars(r)
		taskId, err := strconv.Atoi(vars["taskId"])
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		// Retrieve the task from the database
		var task models.Task
		if err := db.First(&task, taskId).Error; err != nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		// Check if the task belongs to the authenticated user
		if task.UserID != claims.UserID {
			http.Error(w, "Unauthorized to delete this task", http.StatusUnauthorized)
			return
		}

		// Delete the task
		if err := db.Delete(&task).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send a JSON response indicating successful deletion
		config.SendJSONResponse(w, map[string]string{"message": "Task has been successfully deleted"})
	}
}
