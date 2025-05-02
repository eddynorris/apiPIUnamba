package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/models"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/repository"
	"github.com/gorilla/mux"
)

// GetGruposHandler handles fetching all groups or searching based on criteria.
func GetGruposHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupName := r.URL.Query().Get("grupo")
		investigatorName := r.URL.Query().Get("investigador")
		year := r.URL.Query().Get("año") // Assuming 'año' is the query parameter for year
		lineaInvestigacion := r.URL.Query().Get("lineaInvestigacion")

		var grupos []models.Grupo
		var err error

		if groupName != "" || investigatorName != "" || year != "" || lineaInvestigacion != "" {
			// Perform search if any query parameter is provided
			grupos, err = repository.SearchGrupos(db, groupName, investigatorName, year, lineaInvestigacion)
		} else {
			// Otherwise, get all groups
			grupos, err = repository.GetAllGrupos(db)
		}

		if err != nil {
			log.Printf("Error getting/searching groups: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(grupos)
	}
}

// GetGrupoHandler handles fetching a single group by ID.
func GetGrupoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}

		grupo, err := repository.GetGrupoByID(db, id)
		if err != nil {
			log.Printf("Error getting group by ID: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if grupo == nil {
			http.Error(w, "Grupo not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(grupo)
	}
}

// CreateGrupoHandler handles creating a new group.
func CreateGrupoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var g models.Grupo
		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := repository.CreateGrupo(db, &g); err != nil {
			log.Printf("Error creating group: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(g)
	}
}

// UpdateGrupoHandler handles updating an existing group.
func UpdateGrupoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}

		var g models.Grupo
		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Ensure the ID in the body matches the ID in the URL
		g.ID = id

		if err := repository.UpdateGrupo(db, &g); err != nil {
			log.Printf("Error updating group: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(g)
	}
}

// DeleteGrupoHandler handles deleting a group by ID.
func DeleteGrupoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}

		if err := repository.DeleteGrupo(db, id); err != nil {
			log.Printf("Error deleting group: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetGrupoDetailsHandler retrieves a group's details along with its associated investigators.
func GetGrupoDetailsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}

		// Use the repository function that returns the combined struct
		grupoWithInvestigadores, err := repository.GetGrupoDetails(db, id)
		if err != nil {
			// Log the specific error from the repository
			log.Printf("Error getting group details from repository: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Check if the group was found (repository returns nil, nil if not found)
		if grupoWithInvestigadores == nil {
			http.Error(w, "Grupo not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// Encode the combined struct directly
		json.NewEncoder(w).Encode(grupoWithInvestigadores)
	}
}

// Struct to represent the investigator relationship in the combined creation request
type InvestigatorRelationshipRequest struct {
	IDInvestigador int    `json:"idInvestigador"`
	TipoRelacion   string `json:"tipoRelacion"`
}

// Struct to represent the combined group and details creation request body
type CreateGrupoWithDetailsRequest struct {
	models.Grupo   `json:"grupo"`
	Investigadores []InvestigatorRelationshipRequest `json:"investigadores"`
}

// Handler for creating a group with associated investigator details
func CreateGrupoWithDetailsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody CreateGrupoWithDetailsRequest
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		// Use a deferred function for commit/rollback based on error
		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
				panic(p) // Re-panic after rollback
			} else if err != nil {
				// Log the error that caused the rollback
				log.Printf("Rolling back transaction due to error: %v", err)
				tx.Rollback() // Rollback on any error
			} else {
				err = tx.Commit() // Commit otherwise
				if err != nil {
					log.Printf("Error committing transaction: %v", err)
					// Don't send HTTP error here as response might have already been written
				}
			}
		}()

		// Create the group within the transaction using QueryRow with RETURNING
		grupoToCreate := requestBody.Grupo
		// Use lowercase snake_case names and $n placeholders
		groupInsertQuery := `INSERT INTO grupo (nombre, numeroResolucion, lineaInvestigacion, tipoInvestigacion, fechaRegistro, archivo) VALUES ($1, $2, $3, $4, $5, $6) RETURNING idGrupo`
		var grupoID int64 // Use int64 for Scan with RETURNING

		err = tx.QueryRow(groupInsertQuery, grupoToCreate.Nombre, grupoToCreate.NumeroResolucion, grupoToCreate.LineaInvestigacion, grupoToCreate.TipoInvestigacion, grupoToCreate.FechaRegistro, grupoToCreate.Archivo).Scan(&grupoID)
		if err != nil {
			// Error is logged and transaction rolled back by defer
			log.Printf("Error inserting group in transaction: %v", err)
			http.Error(w, "Internal server error during group creation", http.StatusInternalServerError)
			return
		}

		// Create the detailed relationships within the transaction using Exec
		// Use lowercase snake_case names and $n placeholders
		detailInsertQuery := `INSERT INTO Grupo_Investigador (idGrupo, idInvestigador, tipo_relacion) VALUES ($1, $2, $3)`
		for _, invRel := range requestBody.Investigadores {
			_, err = tx.Exec(detailInsertQuery, grupoID, invRel.IDInvestigador, invRel.TipoRelacion)
			if err != nil {
				// Error is logged and transaction rolled back by defer
				log.Printf("Error inserting group-investigator detail in transaction: %v", err)
				http.Error(w, "Internal server error during detail creation", http.StatusInternalServerError)
				return
			}
		}

		// If we reach here without error, the defer func will handle the commit.

		// Prepare the response
		grupoToCreate.ID = int(grupoID) // Convert int64 back to int for the response model
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(grupoToCreate)
	}
}

// GetGruposByInvestigadorHandler maneja la obtención de todos los grupos a los que pertenece un investigador.
func GetGruposByInvestigadorHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["idInvestigador"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de investigador inválido", http.StatusBadRequest)
			return
		}

		gruposConIntegrantes, err := repository.GetGruposByInvestigadorID(db, id)
		if err != nil {
			log.Printf("Error obteniendo grupos por investigador: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		// Enriquecer la respuesta para incluir los integrantes con su rol
		var respuesta []map[string]interface{}
		for _, grupoConInt := range gruposConIntegrantes {
			respuesta = append(respuesta, map[string]interface{}{
				"grupo":       grupoConInt["grupo"],
				"integrantes": grupoConInt["integrantes"],
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respuesta)
	}
}
