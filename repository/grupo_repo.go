package repository

import (
	"database/sql"
	"fmt"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/models"
)

// GetAllGrupos retrieves all groups from the database.
func GetAllGrupos(db *sql.DB) ([]models.Grupo, error) {
	rows, err := db.Query("SELECT idGrupo, nombre, numeroResolucion, lineaInvestigacion, tipoInvestigacion, fechaRegistro, archivo FROM Grupo")
	if err != nil {
		return nil, fmt.Errorf("error querying groups: %w", err)
	}
	defer rows.Close()

	grupos := []models.Grupo{}
	for rows.Next() {
		var g models.Grupo
		if err := rows.Scan(&g.ID, &g.Nombre, &g.NumeroResolucion, &g.LineaInvestigacion, &g.TipoInvestigacion, &g.FechaRegistro, &g.Archivo); err != nil {
			return nil, fmt.Errorf("error scanning group row: %w", err)
		}
		grupos = append(grupos, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through group rows: %w", err)
	}

	return grupos, nil
}

// GetGrupoByID retrieves a single group by its ID.
func GetGrupoByID(db *sql.DB, id int) (*models.Grupo, error) {
	var g models.Grupo
	err := db.QueryRow("SELECT idGrupo, nombre, numeroResolucion, lineaInvestigacion, tipoInvestigacion, fechaRegistro, archivo FROM Grupo WHERE idGrupo = ?", id).Scan(&g.ID, &g.Nombre, &g.NumeroResolucion, &g.LineaInvestigacion, &g.TipoInvestigacion, &g.FechaRegistro, &g.Archivo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil for both when not found
		}
		return nil, fmt.Errorf("error getting group by ID: %w", err)
	}
	return &g, nil
}

// CreateGrupo inserts a new group into the database.
func CreateGrupo(db *sql.DB, g *models.Grupo) error {
	result, err := db.Exec("INSERT INTO Grupo (nombre, numeroResolucion, lineaInvestigacion, tipoInvestigacion, fechaRegistro, archivo) VALUES (?, ?, ?, ?, ?, ?)", g.Nombre, g.NumeroResolucion, g.LineaInvestigacion, g.TipoInvestigacion, g.FechaRegistro, g.Archivo)
	if err != nil {
		return fmt.Errorf("error inserting group: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}

	g.ID = int(id)
	return nil
}

// UpdateGrupo updates an existing group in the database.
func UpdateGrupo(db *sql.DB, g *models.Grupo) error {
	_, err := db.Exec("UPDATE Grupo SET nombre = ?, numeroResolucion = ?, lineaInvestigacion = ?, tipoInvestigacion = ?, fechaRegistro = ?, archivo = ? WHERE idGrupo = ?", g.Nombre, g.NumeroResolucion, g.LineaInvestigacion, g.TipoInvestigacion, g.FechaRegistro, g.Archivo, g.ID)
	if err != nil {
		return fmt.Errorf("error updating group: %w", err)
	}
	return nil
}

// DeleteGrupo deletes a group from the database.
func DeleteGrupo(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM Grupo WHERE idGrupo = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting group: %w", err)
	}
	return nil
}

// SearchGrupos searches for groups based on optional criteria.
func SearchGrupos(db *sql.DB, groupName, investigatorName, year string) ([]models.Grupo, error) {
	query := "SELECT DISTINCT g.idGrupo, g.nombre, g.numeroResolucion, g.lineaInvestigacion, g.tipoInvestigacion, g.fechaRegistro, g.archivo FROM Grupo g JOIN Detalle_GrupoInvestigador dgi ON g.idGrupo = dgi.idGrupo JOIN Investigador i ON dgi.idInvestigador = i.idInvestigador WHERE 1=1"
	args := []interface{}{}

	if groupName != "" {
		query += " AND g.nombre LIKE ?"
		args = append(args, "%"+groupName+"%")
	}

	if investigatorName != "" {
		query += " AND (i.nombre LIKE ? OR i.apellido LIKE ?)"
		args = append(args, "%"+investigatorName+"%", "%"+investigatorName+"%")
	}

	if year != "" {
		query += " AND YEAR(g.fechaRegistro) = ?"
		args = append(args, year)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error searching groups: %w", err)
	}
	defer rows.Close()

	grupos := []models.Grupo{}
	for rows.Next() {
		var g models.Grupo
		if err := rows.Scan(&g.ID, &g.Nombre, &g.NumeroResolucion, &g.LineaInvestigacion, &g.TipoInvestigacion, &g.FechaRegistro, &g.Archivo); err != nil {
			return nil, fmt.Errorf("error scanning group row during search: %w", err)
		}
		grupos = append(grupos, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through group search rows: %w", err)
	}
	return grupos, nil
}
