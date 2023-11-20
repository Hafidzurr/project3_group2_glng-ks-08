package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/config"
	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/controllers"
	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/migrations"
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

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v\n", err)
	}

	router := mux.NewRouter()
	routes.RegisterRoutes(router, db)

	// Mendapatkan port dari variabel lingkungan PORT
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Printf("Server is running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
