package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/models"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/repository"
	"github.com/gorilla/mux"
)

// Struct to represent group details with associated investigators for response
type GrupoDetail struct {
	models.Grupo
	Investigadores []models.Investigador `json:"investigadores"`
}

// GetGruposHandler handles fetching all groups or searching based on criteria.
func GetGruposHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupName := r.URL.Query().Get("grupo")
		investigatorName := r.URL.Query().Get("investigador")
		year := r.URL.Query().Get("año") // Assuming 'año' is the query parameter for year

		var grupos []models.Grupo
		var err error

		if groupName != "" || investigatorName != "" || year != "" {
			// Perform search if any query parameter is provided
			grupos, err = repository.SearchGrupos(db, groupName, investigatorName, year)
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

		grupoDetail, err := repository.GetGrupoDetails(db, id)
		if err != nil {
			log.Printf("Error getting group details: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if grupoDetail == nil {
			http.Error(w, "Grupo not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(grupoDetail)
	}
}

// Struct to represent the investigator relationship in the combined creation request
type InvestigatorRelationshipRequest struct {
	IDInvestigador int    `json:"idInvestigador"`
	TipoRelacion   string `json:"tipoRelacion"`
}

// Struct to represent the combined group and details creation request body
type CreateGrupoWithDetailsRequest struct {
	models.Grupo             `json:"grupo"`
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
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(r)
			} else if err != nil {
				tx.Rollback()
			} else {
				err = tx.Commit()
				if err != nil {
					log.Printf("Error committing transaction: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}
		}()

		// Create the group within the transaction
		grupoToCreate := requestBody.Grupo
		result, err := tx.Exec("INSERT INTO Grupo (nombre, numeroResolucion, lineaInvestigacion, tipoInvestigacion, fechaRegistro, archivo) VALUES (?, ?, ?, ?, ?, ?)", grupoToCreate.Nombre, grupoToCreate.NumeroResolucion, grupoToCreate.LineaInvestigacion, grupoToCreate.TipoInvestigacion, grupoToCreate.FechaRegistro, grupoToCreate.Archivo)
		if err != nil {
			log.Printf("Error inserting group in transaction: %v", err)
			return // Rollback handled by defer
		}

		grupoID, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error getting last insert ID for group in transaction: %v", err)
			return // Rollback handled by defer
		}

		// Create the detailed relationships within the transaction
		for _, invRel := range requestBody.Investigadores {
			_, err := tx.Exec("INSERT INTO Detalle_GrupoInvestigador (idGrupo, idInvestigador, tipoRelacion) VALUES (?, ?, ?)", grupoID, invRel.IDInvestigador, invRel.TipoRelacion)
			if err != nil {
				log.Printf("Error inserting group-investigator detail in transaction: %v", err)
				return // Rollback handled by defer
			}
		}

		// If we reach here, everything was successful (before explicit commit)
		// The defer function will handle the commit or rollback based on the 'err' variable.

		// Prepare the response (you might want to return the created group with its ID)
		grupoToCreate.ID = int(grupoID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(grupoToCreate)
	}
}
