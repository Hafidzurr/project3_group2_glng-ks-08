package routes

import (
	"fmt"
	"net/http"

	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/controllers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(router *mux.Router, db *gorm.DB) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API Project 3 Kelompok 2")
	})
	// User routes
	router.HandleFunc("/users/register", controllers.RegisterUser(db)).Methods("POST")
	router.HandleFunc("/users/login", controllers.LoginUser(db)).Methods("POST")
	router.HandleFunc("/users/update-account", controllers.UpdateUserAccount(db)).Methods("PUT")
	router.HandleFunc("/users/delete-account", controllers.DeleteUserAccount(db)).Methods("DELETE")

	// Category routes
	router.HandleFunc("/categories", controllers.CreateCategory(db)).Methods("POST")
	router.HandleFunc("/categories", controllers.GetCategories(db)).Methods("GET")
	router.HandleFunc("/categories/{categoryId}", controllers.UpdateCategory(db)).Methods("PATCH")
	router.HandleFunc("/categories/{categoryId}", controllers.DeleteCategory(db)).Methods("DELETE")

	// Task routes
	router.HandleFunc("/tasks", controllers.CreateTask(db)).Methods("POST")
	router.HandleFunc("/tasks/{taskId}", controllers.UpdateTask(db)).Methods("PUT")
	router.HandleFunc("/tasks/update-status/{taskId}", controllers.UpdateTaskStatus(db)).Methods("PATCH")
	router.HandleFunc("/tasks/update-category/{taskId}", controllers.UpdateTaskCategory(db)).Methods("PATCH")
	router.HandleFunc("/tasks/{taskId}", controllers.DeleteTask(db)).Methods("DELETE")
	router.HandleFunc("/tasks", controllers.GetTasks(db)).Methods("GET")

}
