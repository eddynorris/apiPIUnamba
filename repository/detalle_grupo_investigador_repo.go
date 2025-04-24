package repository

import (
	"database/sql"
	"fmt"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/models"
)

// CreateDetalleGrupoInvestigador inserts a new relationship between a group and an investigator.
func CreateDetalleGrupoInvestigador(db *sql.DB, detalle *models.DetalleGrupoInvestigador) error {
	result, err := db.Exec("INSERT INTO Detalle_GrupoInvestigador (idGrupo, idInvestigador, tipoRelacion) VALUES (?, ?, ?)", detalle.IDGrupo, detalle.IDInvestigador, detalle.TipoRelacion)
	if err != nil {
		return fmt.Errorf("error inserting group-investigator detail: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID for group-investigator detail: %w", err)
	}

	detalle.ID = int(id)
	return nil
}

// GetDetallesByGrupoID retrieves all relationship details for a given group ID.
func GetDetallesByGrupoID(db *sql.DB, grupoID int) ([]models.DetalleGrupoInvestigador, error) {
	rows, err := db.Query("SELECT idDetalleGI, idGrupo, idInvestigador, tipoRelacion FROM Detalle_GrupoInvestigador WHERE idGrupo = ?", grupoID)
	if err != nil {
		return nil, fmt.Errorf("error querying group-investigator details by group ID: %w", err)
	}
	defer rows.Close()

	detalles := []models.DetalleGrupoInvestigador{}
	for rows.Next() {
		var d models.DetalleGrupoInvestigador
		if err := rows.Scan(&d.ID, &d.IDGrupo, &d.IDInvestigador, &d.TipoRelacion); err != nil {
			return nil, fmt.Errorf("error scanning group-investigator detail row: %w", err)
		}
		detalles = append(detalles, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through group-investigator detail rows: %w", err)
	}

	return detalles, nil
}

// DeleteDetalleGrupoInvestigador deletes a specific relationship detail by its ID.
func DeleteDetalleGrupoInvestigador(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM Detalle_GrupoInvestigador WHERE idDetalleGI = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting group-investigator detail: %w", err)
	}
	return nil
}

// GetDetalleGrupoInvestigadorByID retrieves a single relationship detail by its ID.
// This might be useful for updating a specific relationship (e.g., changing a role).
func GetDetalleGrupoInvestigadorByID(db *sql.DB, id int) (*models.DetalleGrupoInvestigador, error) {
	var d models.DetalleGrupoInvestigador
	err := db.QueryRow("SELECT idDetalleGI, idGrupo, idInvestigador, tipoRelacion FROM Detalle_GrupoInvestigador WHERE idDetalleGI = ?", id).Scan(&d.ID, &d.IDGrupo, &d.IDInvestigador, &d.TipoRelacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil for both when not found
		}
		return nil, fmt.Errorf("error getting group-investigator detail by ID: %w", err)
	}
	return &d, nil
}

// UpdateDetalleGrupoInvestigador updates an existing relationship detail.
func UpdateDetalleGrupoInvestigador(db *sql.DB, detalle *models.DetalleGrupoInvestigador) error {
	_, err := db.Exec("UPDATE Detalle_GrupoInvestigador SET idGrupo = ?, idInvestigador = ?, tipoRelacion = ? WHERE idDetalleGI = ?", detalle.IDGrupo, detalle.IDInvestigador, detalle.TipoRelacion, detalle.ID)
	if err != nil {
		return fmt.Errorf("error updating group-investigator detail: %w", err)
	}
	return nil
}
