package handlers

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

// GetInvestigadoresHandler handles fetching all investigators or searching by name.
func GetInvestigadoresHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name") // Assuming 'name' is the query parameter for investigator name

		var investigadores []models.Investigador
		var err error

		if name != "" {
			// Perform search if name parameter is provided
			investigadores, err = repository.SearchInvestigadores(db, name)
		} else {
			// Otherwise, get all investigators
			investigadores, err = repository.GetAllInvestigadores(db)
		}

		if err != nil {
			log.Printf("Error getting/searching investigators: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(investigadores)
	}
}

// GetInvestigadorHandler handles fetching a single investigator by ID.
func GetInvestigadorHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid investigator ID", http.StatusBadRequest)
			return
		}

		investigador, err := repository.GetInvestigadorByID(db, id)
		if err != nil {
			log.Printf("Error getting investigator by ID: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if investigador == nil {
			http.Error(w, "Investigador not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(investigador)
	}
}

// CreateInvestigadorHandler handles creating a new investigator.
func CreateInvestigadorHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var inv models.Investigador
		if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := repository.CreateInvestigador(db, &inv); err != nil {
			log.Printf("Error creating investigator: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(inv)
	}
}

// UpdateInvestigadorHandler handles updating an existing investigator.
func UpdateInvestigadorHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid investigator ID", http.StatusBadRequest)
			return
		}

		var inv models.Investigador
		if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Ensure the ID in the body matches the ID in the URL
		inv.ID = id

		if err := repository.UpdateInvestigador(db, &inv); err != nil {
			log.Printf("Error updating investigator: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(inv)
	}
}

// DeleteInvestigadorHandler handles deleting an investigator by ID.
func DeleteInvestigadorHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid investigator ID", http.StatusBadRequest)
			return
		}

		if err := repository.DeleteInvestigador(db, id); err != nil {
			log.Printf("Error deleting investigator: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
