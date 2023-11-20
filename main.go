package main

import (
	"log"
	"net/http"

	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/config"
	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/controllers"
	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/routes"
	"github.com/gorilla/mux"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	// Initialize Admin User
	controllers.InitializeAdminUser(db)

	router := mux.NewRouter()
	routes.RegisterRoutes(router, db)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
