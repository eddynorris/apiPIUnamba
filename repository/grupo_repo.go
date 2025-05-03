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

// SearchGrupos searches for groups based on optional criteria and returns them with their investigators and roles.
// Returns []models.GrupoWithInvestigadores
func SearchGrupos(db *sql.DB, groupName, investigatorName, year, lineaInvestigacion string) ([]models.GrupoWithInvestigadores, error) {
	// Query needs to select fields from grupo, investigador, and the linking table (Grupo_Investigador)
	// No DISTINCT needed here, we group in Go. ORDER BY g.idGrupo is helpful.
	args := []interface{}{}
	placeholderCount := 1

	// --- Build WHERE clause dynamically ---
	whereConditions := ""

	if groupName != "" {
		whereConditions += fmt.Sprintf(` AND g.nombre ILIKE $%d`, placeholderCount)
		args = append(args, "%"+groupName+"%")
		placeholderCount++
	}

	if investigatorName != "" {
		// Important: This condition needs to be applied correctly. We might get multiple rows for the same group if different investigators match.
		// This WHERE clause will filter the *rows* returned, not necessarily limit the groups to *only* those containing the matching investigator.
		// To strictly filter groups *containing* the investigator, a subquery or EXISTS might be needed, making the query more complex.
		// For now, we filter the rows and reconstruct.
		whereConditions += fmt.Sprintf(` AND (i.nombre ILIKE $%d OR i.apellido ILIKE $%d)`, placeholderCount, placeholderCount+1)
		args = append(args, "%"+investigatorName+"%", "%"+investigatorName+"%")
		placeholderCount += 2
	}

	if year != "" {
		whereConditions += fmt.Sprintf(` AND EXTRACT(YEAR FROM g.fechaRegistro) = $%d`, placeholderCount)
		args = append(args, year)
		placeholderCount++
	}

	if lineaInvestigacion != "" {
		whereConditions += fmt.Sprintf(` AND g.lineaInvestigacion ILIKE $%d`, placeholderCount)
		args = append(args, "%"+lineaInvestigacion+"%")
		placeholderCount++
	}
	// --- End WHERE clause build ---

	// We need to find the groups that match the criteria first, then get their details.
	// A more robust approach uses a subquery to filter groups first.
	finalQuery := `WITH FilteredGroups AS (
		SELECT DISTINCT g.idGrupo
		FROM grupo g
		LEFT JOIN Grupo_Investigador dgi ON g.idGrupo = dgi.idGrupo
		LEFT JOIN investigador i ON dgi.idInvestigador = i.idInvestigador
		WHERE 1=1` + whereConditions + `
	)
	SELECT
		g.idGrupo, g.nombre, g.numeroResolucion, g.lineaInvestigacion, g.tipoInvestigacion, g.fechaRegistro, g.archivo, g.createdAt, g.updatedAt,
		i.idInvestigador, i.nombre, i.apellido, i.createdAt as invCreatedAt, i.updatedAt as invUpdatedAt,
		dgi.rol
	FROM grupo g
	JOIN Grupo_Investigador dgi ON g.idGrupo = dgi.idGrupo
	JOIN investigador i ON dgi.idInvestigador = i.idInvestigador
	WHERE g.idGrupo IN (SELECT idGrupo FROM FilteredGroups)
	ORDER BY g.idGrupo, i.idInvestigador -- Order to help grouping in Go
	`

	rows, err := db.Query(finalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error searching groups with details: %w", err)
	}
	defer rows.Close()

	// Map to group investigators by group ID
	grupoMap := make(map[int]*models.GrupoWithInvestigadores)
	// Slice to maintain order of groups found
	orderedGrupos := []*models.GrupoWithInvestigadores{}

	for rows.Next() {
		var g models.Grupo
		var inv models.InvestigadorConRol           // Use the struct with role
		var invCreatedAt, invUpdatedAt sql.NullTime // Use sql.NullTime for potentially null timestamps from joins

		if err := rows.Scan(
			&g.ID, &g.Nombre, &g.NumeroResolucion, &g.LineaInvestigacion, &g.TipoInvestigacion, &g.FechaRegistro, &g.Archivo, &g.CreatedAt, &g.UpdatedAt,
			&inv.ID, &inv.Nombre, &inv.Apellido, &invCreatedAt, &invUpdatedAt,
			&inv.Rol,
		); err != nil {
			return nil, fmt.Errorf("error scanning group/investigator row during search: %w", err)
		}

		// Handle potentially null timestamps for investigator
		if invCreatedAt.Valid {
			inv.CreatedAt = invCreatedAt.Time
		}
		if invUpdatedAt.Valid {
			inv.UpdatedAt = invUpdatedAt.Time
		}

		// Check if we've already seen this group
		if _, exists := grupoMap[g.ID]; !exists {
			// First time seeing this group
			grupoWithDetails := &models.GrupoWithInvestigadores{
				Grupo:          g,
				Investigadores: []models.InvestigadorConRol{}, // Initialize empty slice
			}
			grupoMap[g.ID] = grupoWithDetails
			orderedGrupos = append(orderedGrupos, grupoWithDetails) // Add to ordered list
		}

		// Add the current investigator (with role) to this group's list
		grupoMap[g.ID].Investigadores = append(grupoMap[g.ID].Investigadores, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through group search rows: %w", err)
	}

	// Convert []*models.GrupoWithInvestigadores to []models.GrupoWithInvestigadores
	result := make([]models.GrupoWithInvestigadores, len(orderedGrupos))
	for i, ptr := range orderedGrupos {
		result[i] = *ptr
	}

	return result, nil
}

// GetGrupoDetails retrieves a group and its associated investigators including their roles.
func GetGrupoDetails(db *sql.DB, id int) (*models.GrupoWithInvestigadores, error) {
	// 1. Get the group details
	grupo, err := GetGrupoByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("error in GetGrupoByID called from GetGrupoDetails: %w", err)
	}
	if grupo == nil { // Should not happen
		return nil, nil
	}

	// 2. Get associated investigators with their roles in this specific group
	query := `
		SELECT i.idInvestigador, i.nombre, i.apellido, dgi.rol, i.createdAt, i.updatedAt
		FROM investigador i
		JOIN Grupo_Investigador dgi ON i.idInvestigador = dgi.idInvestigador
		WHERE dgi.idGrupo = $1
	`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying investigators for group details: %w", err)
	}
	defer rows.Close()

	investigadores := []models.InvestigadorConRol{}
	for rows.Next() {
		var inv models.InvestigadorConRol
		// Scan id, nombre, apellido, rol, createdAt, updatedAt
		if err := rows.Scan(&inv.ID, &inv.Nombre, &inv.Apellido, &inv.Rol, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning investigator row with role for group details: %w", err)
		}
		investigadores = append(investigadores, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating investigator rows for group details: %w", err)
	}

	// 3. Combine results
	grupoDetail := &models.GrupoWithInvestigadores{
		Grupo:          *grupo,
		Investigadores: investigadores, // Now contains investigators with roles
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
