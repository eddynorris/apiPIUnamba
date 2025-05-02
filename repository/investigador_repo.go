package repository

import (
	"database/sql"
	"fmt"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/models"
)

// GetAllInvestigadores retrieves all investigators from the database.
func GetAllInvestigadores(db *sql.DB) ([]models.Investigador, error) {
	rows, err := db.Query(`SELECT idInvestigador, nombre, apellido, createdAt, updatedAt FROM investigador`)
	if err != nil {
		return nil, fmt.Errorf("error querying investigators: %w", err)
	}
	defer rows.Close()

	investigadores := []models.Investigador{}
	for rows.Next() {
		var inv models.Investigador
		if err := rows.Scan(&inv.ID, &inv.Nombre, &inv.Apellido, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning investigator row: %w", err)
		}
		investigadores = append(investigadores, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through investigator rows: %w", err)
	}

	return investigadores, nil
}

// GetInvestigadorByID retrieves a single investigator by their ID.
func GetInvestigadorByID(db *sql.DB, id int) (*models.Investigador, error) {
	var inv models.Investigador
	err := db.QueryRow(`SELECT idInvestigador, nombre, apellido, createdAt, updatedAt FROM investigador WHERE idInvestigador = $1`, id).Scan(&inv.ID, &inv.Nombre, &inv.Apellido, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil for both when not found
		}
		return nil, fmt.Errorf("error getting investigator by ID: %w", err)
	}
	return &inv, nil
}

// CreateInvestigador inserts a new investigator into the database.
func CreateInvestigador(db *sql.DB, inv *models.Investigador) error {
	query := `INSERT INTO investigador (nombre, apellido) VALUES ($1, $2) RETURNING idInvestigador, createdAt, updatedAt`
	err := db.QueryRow(query, inv.Nombre, inv.Apellido).Scan(&inv.ID, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error inserting investigator: %w", err)
	}
	return nil
}

// UpdateInvestigador updates an existing investigator in the database.
func UpdateInvestigador(db *sql.DB, inv *models.Investigador) error {
	_, err := db.Exec(`UPDATE investigador SET nombre = $1, apellido = $2, updatedAt = CURRENT_TIMESTAMP WHERE idInvestigador = $3`, inv.Nombre, inv.Apellido, inv.ID)
	if err != nil {
		return fmt.Errorf("error updating investigator: %w", err)
	}
	return nil
}

// DeleteInvestigador deletes an investigator from the database.
func DeleteInvestigador(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM investigador WHERE idInvestigador = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting investigator: %w", err)
	}
	return nil
}

// SearchInvestigadores searches for investigators based on optional criteria.
func SearchInvestigadores(db *sql.DB, name string) ([]models.Investigador, error) {
	query := `SELECT idInvestigador, nombre, apellido, createdAt, updatedAt FROM investigador WHERE 1=1`
	args := []interface{}{}
	placeholderCount := 1

	if name != "" {
		query += fmt.Sprintf(` AND (nombre ILIKE $%d OR apellido ILIKE $%d)`, placeholderCount, placeholderCount+1)
		searchPattern := "%" + name + "%"
		args = append(args, searchPattern, searchPattern)
		placeholderCount += 2 // Increment by 2 because we added two placeholders
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error searching investigators: %w", err)
	}
	defer rows.Close()

	investigadores := []models.Investigador{}
	for rows.Next() {
		var inv models.Investigador
		if err := rows.Scan(&inv.ID, &inv.Nombre, &inv.Apellido, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning investigator row during search: %w", err)
		}
		investigadores = append(investigadores, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through investigator search rows: %w", err)
	}

	return investigadores, nil
}
