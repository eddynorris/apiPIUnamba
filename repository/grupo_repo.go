package repository

import (
	"database/sql"
	"fmt"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/models"
)

// GetAllGrupos retrieves all groups from the database.
func GetAllGrupos(db *sql.DB) ([]models.Grupo, error) {
	rows, err := db.Query(`SELECT id_grupo, nombre, numero_resolucion, linea_investigacion, tipo_investigacion, fecha_registro, archivo FROM grupo`)
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
	err := db.QueryRow(`SELECT id_grupo, nombre, numero_resolucion, linea_investigacion, tipo_investigacion, fecha_registro, archivo FROM grupo WHERE id_grupo = $1`, id).Scan(&g.ID, &g.Nombre, &g.NumeroResolucion, &g.LineaInvestigacion, &g.TipoInvestigacion, &g.FechaRegistro, &g.Archivo)
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
	query := `INSERT INTO grupo (nombre, numero_resolucion, linea_investigacion, tipo_investigacion, fecha_registro, archivo) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_grupo`
	err := db.QueryRow(query, g.Nombre, g.NumeroResolucion, g.LineaInvestigacion, g.TipoInvestigacion, g.FechaRegistro, g.Archivo).Scan(&g.ID)
	if err != nil {
		return fmt.Errorf("error inserting group: %w", err)
	}
	return nil
}

// UpdateGrupo updates an existing group in the database.
func UpdateGrupo(db *sql.DB, g *models.Grupo) error {
	_, err := db.Exec(`UPDATE grupo SET nombre = $1, numero_resolucion = $2, linea_investigacion = $3, tipo_investigacion = $4, fecha_registro = $5, archivo = $6 WHERE id_grupo = $7`, g.Nombre, g.NumeroResolucion, g.LineaInvestigacion, g.TipoInvestigacion, g.FechaRegistro, g.Archivo, g.ID)
	if err != nil {
		return fmt.Errorf("error updating group: %w", err)
	}
	return nil
}

// DeleteGrupo deletes a group from the database.
func DeleteGrupo(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM grupo WHERE id_grupo = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting group: %w", err)
	}
	return nil
}

// SearchGrupos searches for groups based on optional criteria.
func SearchGrupos(db *sql.DB, groupName, investigatorName, year string) ([]models.Grupo, error) {
	query := `SELECT DISTINCT g.id_grupo, g.nombre, g.numero_resolucion, g.linea_investigacion, g.tipo_investigacion, g.fecha_registro, g.archivo
			 FROM grupo g
			 JOIN detalle_grupo_investigador dgi ON g.id_grupo = dgi.id_grupo
			 JOIN investigador i ON dgi.id_investigador = i.id_investigador
			 WHERE 1=1`
	args := []interface{}{}
	placeholderCount := 1

	if groupName != "" {
		query += fmt.Sprintf(` AND g.nombre ILIKE $%d`, placeholderCount)
		args = append(args, "%"+groupName+"%")
		placeholderCount++
	}

	if investigatorName != "" {
		query += fmt.Sprintf(` AND (i.nombre ILIKE $%d OR i.apellido ILIKE $%d)`, placeholderCount, placeholderCount+1)
		args = append(args, "%"+investigatorName+"%", "%"+investigatorName+"%")
		placeholderCount += 2
	}

	if year != "" {
		query += fmt.Sprintf(` AND EXTRACT(YEAR FROM g.fecha_registro) = $%d`, placeholderCount)
		args = append(args, year)
		placeholderCount++
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

// GetGrupoDetails retrieves a group and its associated investigators.
func GetGrupoDetails(db *sql.DB, id int) (*models.GrupoWithInvestigadores, error) {
	// 1. Get the group details
	grupo, err := GetGrupoByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found, return nil for the result
		}
		return nil, fmt.Errorf("error in GetGrupoByID called from GetGrupoDetails: %w", err)
	}
	if grupo == nil { // Should not happen if GetGrupoByID returns nil, nil for not found
		return nil, nil
	}

	// 2. Get associated investigators
	query := `
		SELECT i.id_investigador, i.nombre, i.apellido, i.rol
		FROM investigador i
		JOIN detalle_grupo_investigador dgi ON i.id_investigador = dgi.id_investigador
		WHERE dgi.id_grupo = $1
	`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying investigators for group details: %w", err)
	}
	defer rows.Close()

	investigadores := []models.Investigador{}
	for rows.Next() {
		var inv models.Investigador
		if err := rows.Scan(&inv.ID, &inv.Nombre, &inv.Apellido, &inv.Rol); err != nil {
			return nil, fmt.Errorf("error scanning investigator row for group details: %w", err)
		}
		investigadores = append(investigadores, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating investigator rows for group details: %w", err)
	}

	// 3. Combine results
	grupoDetail := &models.GrupoWithInvestigadores{
		Grupo:          *grupo,
		Investigadores: investigadores,
	}

	return grupoDetail, nil
}
