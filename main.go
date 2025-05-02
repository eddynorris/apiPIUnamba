package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/database"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/routes" // Import the new routes package
	"github.com/joho/godotenv"                                            // Para cargar variables de entorno desde .env
	// "github.com/GoogleCloudPlatform/golang-samples/run/helloworld/handlers" // Removed old handlers import
	// "github.com/gorilla/mux" // Mux is now used within routes package
)

var db *sql.DB

func main() {
	log.Print("starting server...")

	// Cargar variables de entorno desde .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	db, err = database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Setup routes using the new routes package
	r := routes.SetupRoutes(db)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
