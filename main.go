package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/database"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/routes" // Import the new routes package
	"github.com/joho/godotenv"                                            // Para cargar variables de entorno desde .env
	"github.com/rs/cors"                                                  // Importar CORS
	// "github.com/GoogleCloudPlatform/golang-samples/run/helloworld/handlers" // Removed old handlers import
	// "github.com/gorilla/mux" // Mux is now used within routes package
)

var db *sql.DB

func main() {
	log.Print("starting server...")

	// Cargar variables de entorno desde .env (Principalmente para desarrollo)
	err := godotenv.Load() // Cargar solo si existe, no fallar si no está
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize database connection
	db, err = database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Setup routes using the new routes package
	r := routes.SetupRoutes(db)

	// --- Configuración de CORS ---
	// Define los orígenes permitidos (¡sé específico en producción!)
	allowedOrigins := []string{"http://localhost:5173", "http://localhost:3000"} // Añade la URL de tu frontend de desarrollo
	// Si tienes un dominio de producción, añádelo: "https://tu-dominio.com"
	if os.Getenv("FRONTEND_URL") != "" {
		allowedOrigins = append(allowedOrigins, os.Getenv("FRONTEND_URL"))
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Permitir el header de Auth
		AllowCredentials: true,
		// Enable Debugging for testing, disable in production
		// Debug:            true,
	})

	// Aplicar el middleware CORS al router principal
	httpHandler := c.Handler(r)
	// --- Fin Configuración de CORS ---

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server with CORS handler
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, httpHandler); err != nil {
		log.Fatal(err)
	}
}
