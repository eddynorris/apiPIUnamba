package repository

import (
	"database/sql"
	"fmt"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/models"
)

// CreateDetalleGrupoInvestigador inserts a new relationship between a group and an investigator.
func CreateDetalleGrupoInvestigador(db *sql.DB, detalle *models.DetalleGrupoInvestigador) error {
	// Use lowercase snake_case, $n placeholders, and RETURNING
	query := `INSERT INTO detalle_grupo_investigador (id_grupo, id_investigador, tipo_relacion) VALUES ($1, $2, $3) RETURNING id_detalle_gi`
	err := db.QueryRow(query, detalle.IDGrupo, detalle.IDInvestigador, detalle.TipoRelacion).Scan(&detalle.ID)
	if err != nil {
		return fmt.Errorf("error inserting group-investigator detail: %w", err)
	}
	return nil
}

// GetDetallesByGrupoID retrieves all relationship details for a given group ID.
func GetDetallesByGrupoID(db *sql.DB, grupoID int) ([]models.DetalleGrupoInvestigador, error) {
	// Use lowercase snake_case and $1 placeholder
	rows, err := db.Query(`SELECT id_detalle_gi, id_grupo, id_investigador, tipo_relacion FROM detalle_grupo_investigador WHERE id_grupo = $1`, grupoID)
	if err != nil {
		return nil, fmt.Errorf("error querying group-investigator details by group ID: %w", err)
	}
	defer rows.Close()

	detalles := []models.DetalleGrupoInvestigador{}
	for rows.Next() {
		var d models.DetalleGrupoInvestigador
		// Ensure SELECT order matches struct fields
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
	// Use lowercase snake_case and $1 placeholder
	_, err := db.Exec(`DELETE FROM detalle_grupo_investigador WHERE id_detalle_gi = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting group-investigator detail: %w", err)
	}
	return nil
}

// GetDetalleGrupoInvestigadorByID retrieves a single relationship detail by its ID.
// This might be useful for updating a specific relationship (e.g., changing a role).
func GetDetalleGrupoInvestigadorByID(db *sql.DB, id int) (*models.DetalleGrupoInvestigador, error) {
	var d models.DetalleGrupoInvestigador
	// Use lowercase snake_case and $1 placeholder
	err := db.QueryRow(`SELECT id_detalle_gi, id_grupo, id_investigador, tipo_relacion FROM detalle_grupo_investigador WHERE id_detalle_gi = $1`, id).Scan(&d.ID, &d.IDGrupo, &d.IDInvestigador, &d.TipoRelacion)
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
	// Use lowercase snake_case and $n placeholders
	_, err := db.Exec(`UPDATE detalle_grupo_investigador SET id_grupo = $1, id_investigador = $2, tipo_relacion = $3 WHERE id_detalle_gi = $4`, detalle.IDGrupo, detalle.IDInvestigador, detalle.TipoRelacion, detalle.ID)
	if err != nil {
		return fmt.Errorf("error updating group-investigator detail: %w", err)
	}
	return nil
}
