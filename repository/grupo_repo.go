package repository

import (
	"database/sql"
	"fmt"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/models"
)

// GetAllGrupos retrieves all groups from the database.
func GetAllGrupos(db *sql.DB) ([]models.Grupo, error) {
	rows, err := db.Query(`SELECT idGrupo, nombre, numeroResolucion, lineaInvestigacion, tipoInvestigacion, fechaRegistro, archivo, createdAt, updatedAt FROM grupo`)
	if err != nil {
		return nil, fmt.Errorf("error querying groups: %w", err)
	}
	defer rows.Close()

	grupos := []models.Grupo{}
	for rows.Next() {
		var g models.Grupo
		if err := rows.Scan(&g.ID, &g.Nombre, &g.NumeroResolucion, &g.LineaInvestigacion, &g.TipoInvestigacion, &g.FechaRegistro, &g.Archivo, &g.CreatedAt, &g.UpdatedAt); err != nil {
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
	err := db.QueryRow(`SELECT idGrupo, nombre, numeroResolucion, lineaInvestigacion, tipoInvestigacion, fechaRegistro, archivo, createdAt, updatedAt FROM grupo WHERE idGrupo = $1`, id).Scan(&g.ID, &g.Nombre, &g.NumeroResolucion, &g.LineaInvestigacion, &g.TipoInvestigacion, &g.FechaRegistro, &g.Archivo, &g.CreatedAt, &g.UpdatedAt)
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
	query := `INSERT INTO grupo (nombre, numeroResolucion, lineaInvestigacion, tipoInvestigacion, fechaRegistro, archivo) VALUES ($1, $2, $3, $4, $5, $6) RETURNING idGrupo, createdAt, updatedAt`
	err := db.QueryRow(query, g.Nombre, g.NumeroResolucion, g.LineaInvestigacion, g.TipoInvestigacion, g.FechaRegistro, g.Archivo).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error inserting group: %w", err)
	}
	return nil
}

// UpdateGrupo updates an existing group in the database.
func UpdateGrupo(db *sql.DB, g *models.Grupo) error {
	_, err := db.Exec(`UPDATE grupo SET nombre = $1, numeroResolucion = $2, lineaInvestigacion = $3, tipoInvestigacion = $4, fechaRegistro = $5, archivo = $6, updatedAt = CURRENT_TIMESTAMP WHERE idGrupo = $7`, g.Nombre, g.NumeroResolucion, g.LineaInvestigacion, g.TipoInvestigacion, g.FechaRegistro, g.Archivo, g.ID)
	if err != nil {
		return fmt.Errorf("error updating group: %w", err)
	}
	return nil
}

// DeleteGrupo deletes a group from the database.
func DeleteGrupo(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM grupo WHERE idGrupo = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting group: %w", err)
	}
	return nil
}

// SearchGrupos searches for groups based on optional criteria.
func SearchGrupos(db *sql.DB, groupName, investigatorName, year, lineaInvestigacion string) ([]models.Grupo, error) {
	query := `SELECT DISTINCT g.idGrupo, g.nombre, g.numeroResolucion, g.lineaInvestigacion, g.tipoInvestigacion, g.fechaRegistro, g.archivo, g.createdAt, g.updatedAt
			 FROM grupo g
			 JOIN Grupo_Investigador dgi ON g.idGrupo = dgi.idGrupo
			 JOIN investigador i ON dgi.idInvestigador = i.idInvestigador
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
		query += fmt.Sprintf(` AND EXTRACT(YEAR FROM g.fechaRegistro) = $%d`, placeholderCount)
		args = append(args, year)
		placeholderCount++
	}

	if lineaInvestigacion != "" {
		query += fmt.Sprintf(` AND g.lineaInvestigacion ILIKE $%d`, placeholderCount)
		args = append(args, "%"+lineaInvestigacion+"%")
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
		if err := rows.Scan(&g.ID, &g.Nombre, &g.NumeroResolucion, &g.LineaInvestigacion, &g.TipoInvestigacion, &g.FechaRegistro, &g.Archivo, &g.CreatedAt, &g.UpdatedAt); err != nil {
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
		SELECT i.idInvestigador, i.nombre, i.apellido
		FROM investigador i
		JOIN Grupo_Investigador dgi ON i.idInvestigador = dgi.idInvestigador
		WHERE dgi.idGrupo = $1
	`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying investigators for group details: %w", err)
	}
	defer rows.Close()

	investigadores := []models.Investigador{}
	for rows.Next() {
		var inv models.Investigador
		if err := rows.Scan(&inv.ID, &inv.Nombre, &inv.Apellido); err != nil {
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

// GetGruposByInvestigadorID obtiene todos los grupos a los que pertenece un investigador dado su id.
func GetGruposByInvestigadorID(db *sql.DB, idInvestigador int) ([]map[string]interface{}, error) {
	query := `SELECT g.idGrupo, g.nombre, g.numeroResolucion, g.lineaInvestigacion, g.tipoInvestigacion, g.fechaRegistro, g.archivo, g.createdAt, g.updatedAt
				 , dgi.rol
			 FROM grupo g
			 JOIN Grupo_Investigador dgi ON g.idGrupo = dgi.idGrupo
			 WHERE dgi.idInvestigador = $1`
	rows, err := db.Query(query, idInvestigador)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo grupos por idInvestigador: %w", err)
	}
	defer rows.Close()

	var gruposConIntegrantes []map[string]interface{}
	for rows.Next() {
		var g models.Grupo
		var rol string
		if err := rows.Scan(&g.ID, &g.Nombre, &g.NumeroResolucion, &g.LineaInvestigacion, &g.TipoInvestigacion, &g.FechaRegistro, &g.Archivo, &g.CreatedAt, &g.UpdatedAt, &rol); err != nil {
			return nil, fmt.Errorf("error escaneando grupo: %w", err)
		}

		// Obtener los integrantes y sus roles para este grupo
		queryIntegrantes := `SELECT i.idInvestigador, i.nombre, i.apellido, dgi.rol
			FROM investigador i
			JOIN Grupo_Investigador dgi ON i.idInvestigador = dgi.idInvestigador
			WHERE dgi.idGrupo = $1`
		rowsIntegrantes, err := db.Query(queryIntegrantes, g.ID)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo integrantes del grupo: %w", err)
		}
		var integrantesConRol []map[string]interface{}
		for rowsIntegrantes.Next() {
			var idInvestigador int
			var nombre, apellido, rolIntegrante string
			if err := rowsIntegrantes.Scan(&idInvestigador, &nombre, &apellido, &rolIntegrante); err != nil {
				rowsIntegrantes.Close()
				return nil, fmt.Errorf("error escaneando integrante: %w", err)
			}
			integrantesConRol = append(integrantesConRol, map[string]interface{}{
				"idInvestigador": idInvestigador,
				"nombre":         nombre,
				"apellido":       apellido,
				"rol":            rolIntegrante,
			})
		}
		rowsIntegrantes.Close()

		grupoMap := map[string]interface{}{
			"grupo":       g,
			"integrantes": integrantesConRol,
		}
		gruposConIntegrantes = append(gruposConIntegrantes, grupoMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error después de iterar los grupos: %w", err)
	}
	return gruposConIntegrantes, nil
}
