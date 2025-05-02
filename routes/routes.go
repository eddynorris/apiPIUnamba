package routes

import (
	"database/sql"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/controllers"
	"github.com/gorilla/mux"
)

// SetupRoutes configures the application routes.
func SetupRoutes(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Investigador routes
	r.HandleFunc("/investigadores", controllers.GetInvestigadoresHandler(db)).Methods("GET")
	r.HandleFunc("/investigadores", controllers.CreateInvestigadorHandler(db)).Methods("POST")
	r.HandleFunc("/investigadores/{id}", controllers.GetInvestigadorHandler(db)).Methods("GET")
	r.HandleFunc("/investigadores/{id}", controllers.UpdateInvestigadorHandler(db)).Methods("PUT")
	r.HandleFunc("/investigadores/{id}", controllers.DeleteInvestigadorHandler(db)).Methods("DELETE")

	// Grupo routes
	r.HandleFunc("/grupos", controllers.GetGruposHandler(db)).Methods("GET")
	r.HandleFunc("/grupos", controllers.CreateGrupoHandler(db)).Methods("POST")
	r.HandleFunc("/grupos/with-details", controllers.CreateGrupoWithDetailsHandler(db)).Methods("POST")
	r.HandleFunc("/grupos/{id}", controllers.GetGrupoHandler(db)).Methods("GET")
	r.HandleFunc("/grupos/{id}", controllers.UpdateGrupoHandler(db)).Methods("PUT")
	r.HandleFunc("/grupos/{id}", controllers.DeleteGrupoHandler(db)).Methods("DELETE")
	r.HandleFunc("/grupos/{id}/details", controllers.GetGrupoDetailsHandler(db)).Methods("GET")

	// DetalleGrupoInvestigador routes
	r.HandleFunc("/detalles", controllers.CreateDetalleGrupoInvestigadorHandler(db)).Methods("POST")
	r.HandleFunc("/detalles/{id}", controllers.GetDetalleGrupoInvestigadorHandler(db)).Methods("GET")
	r.HandleFunc("/detalles/{id}", controllers.UpdateDetalleGrupoInvestigadorHandler(db)).Methods("PUT")
	r.HandleFunc("/detalles/{id}", controllers.DeleteDetalleGrupoInvestigadorHandler(db)).Methods("DELETE")
	r.HandleFunc("/grupos/{grupoID}/detalles", controllers.GetDetallesByGrupoHandler(db)).Methods("GET")

	return r
}
