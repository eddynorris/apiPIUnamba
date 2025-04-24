package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/database"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/handlers"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	log.Print("starting server...")

	// Initialize database connection
	var err error
	db, err = database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Set up routes using Gorilla Mux for better routing capabilities
	r := mux.NewRouter()

	// Investigador routes (GET will now handle search)
	r.HandleFunc("/investigadores", handlers.GetInvestigadoresHandler(db)).Methods("GET")
	r.HandleFunc("/investigadores", handlers.CreateInvestigadorHandler(db)).Methods("POST")
	r.HandleFunc("/investigadores/{id}", handlers.GetInvestigadorHandler(db)).Methods("GET")
	r.HandleFunc("/investigadores/{id}", handlers.UpdateInvestigadorHandler(db)).Methods("PUT")
	r.HandleFunc("/investigadores/{id}", handlers.DeleteInvestigadorHandler(db)).Methods("DELETE")

	// Grupo routes (GET will now handle search)
	r.HandleFunc("/grupos", handlers.GetGruposHandler(db)).Methods("GET")	r.HandleFunc("/grupos", handlers.CreateGrupoHandler(db)).Methods("POST")
	r.HandleFunc("/grupos/with-details", handlers.CreateGrupoWithDetailsHandler(db)).Methods("POST") // New endpoint for combined creation
	r.HandleFunc("/grupos/{id}", handlers.GetGrupoHandler(db)).Methods("GET")
	r.HandleFunc("/grupos/{id}", handlers.UpdateGrupoHandler(db)).Methods("PUT")
	r.HandleFunc("/grupos/{id}", handlers.DeleteGrupoHandler(db)).Methods("DELETE")
	r.HandleFunc("/grupos/{id}/details", handlers.GetGrupoDetailsHandler(db)).Methods("GET") // Get group details with members

	// DetalleGrupoInvestigador routes
	r.HandleFunc("/detalles", handlers.CreateDetalleGrupoInvestigadorHandler(db)).Methods("POST") // Still useful for adding members to existing groups
	r.HandleFunc("/detalles/{id}", handlers.GetDetalleGrupoInvestigadorHandler(db)).Methods("GET")
	r.HandleFunc("/detalles/{id}", handlers.UpdateDetalleGrupoInvestigadorHandler(db)).Methods("PUT")
	r.HandleFunc("/detalles/{id}", handlers.DeleteDetalleGrupoInvestigadorHandler(db)).Methods("DELETE")
	r.HandleFunc("/grupos/{grupoID}/detalles", handlers.GetDetallesByGrupoHandler(db)).Methods("GET") // Get details by group ID

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
