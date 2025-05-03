package repository

import (
	"database/sql"
	"fmt"

	// Import math for ceiling calculation
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/models"
)

// GetAllGrupos retrieves a paginated list of all groups.
func GetAllGrupos(db *sql.DB, limit, offset int) ([]models.Grupo, int, error) {
	// Query for the data page
	query := `SELECT idGrupo, nombre, numeroResolucion, lineaInvestigacion, tipoInvestigacion, fechaRegistro, archivo, createdAt, updatedAt FROM grupo ORDER BY nombre LIMIT $1 OFFSET $2`
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying groups page: %w", err)
	}
	defer rows.Close()

	grupos := []models.Grupo{}
	for rows.Next() {
		var g models.Grupo
		if err := rows.Scan(&g.ID, &g.Nombre, &g.NumeroResolucion, &g.LineaInvestigacion, &g.TipoInvestigacion, &g.FechaRegistro, &g.Archivo, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("error scanning group row: %w", err)
		}
		grupos = append(grupos, g)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error after iterating through group rows: %w", err)
	}

	// Query for the total count
	var total int
	countQuery := `SELECT COUNT(*) FROM grupo`
	if err := db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("error querying total group count: %w", err)
	}

	return grupos, total, nil
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

// SearchGrupos searches for groups with pagination and returns them with investigators and roles.
func SearchGrupos(db *sql.DB, groupName, investigatorName, year, lineaInvestigacion, tipoInvestigacion string, limit, offset int) ([]models.GrupoWithInvestigadores, int, error) {
	args := []interface{}{}
	placeholderCount := 1

	// --- Build WHERE clause dynamically (for the CTE) ---
	whereConditions := ""

	if groupName != "" {
		// Apply unaccent to both column and search term
		whereConditions += fmt.Sprintf(` AND unaccent(g.nombre) ILIKE unaccent($%d)`, placeholderCount)
		args = append(args, "%"+groupName+"%")
		placeholderCount++
	}

	if investigatorName != "" {
		// Apply unaccent to both column and search term
		whereConditions += fmt.Sprintf(` AND (unaccent(i.nombre) ILIKE unaccent($%d) OR unaccent(i.apellido) ILIKE unaccent($%d))`, placeholderCount, placeholderCount+1)
		args = append(args, "%"+investigatorName+"%", "%"+investigatorName+"%")
		placeholderCount += 2
	}

	if year != "" {
		whereConditions += fmt.Sprintf(` AND EXTRACT(YEAR FROM g.fechaRegistro) = $%d`, placeholderCount)
		args = append(args, year)
		placeholderCount++
	}

	if lineaInvestigacion != "" {
		// Apply unaccent to both column and search term
		whereConditions += fmt.Sprintf(` AND unaccent(g.lineaInvestigacion) ILIKE unaccent($%d)`, placeholderCount)
		args = append(args, "%"+lineaInvestigacion+"%")
		placeholderCount++
	}

	if tipoInvestigacion != "" {
		// Apply unaccent to both column and search term
		whereConditions += fmt.Sprintf(` AND unaccent(g.tipoInvestigacion) ILIKE unaccent($%d)`, placeholderCount)
		args = append(args, "%"+tipoInvestigacion+"%")
		placeholderCount++
	}
	// --- End WHERE clause build ---

	// CTE to find matching group IDs based on filters
	cteQuery := `WITH FilteredGroups AS (
		SELECT DISTINCT g.idGrupo
		FROM grupo g
		LEFT JOIN Grupo_Investigador dgi ON g.idGrupo = dgi.idGrupo
		LEFT JOIN investigador i ON dgi.idInvestigador = i.idInvestigador
		WHERE 1=1` + whereConditions + `
	)`

	// Query for the data page using the CTE
	dataQuery := cteQuery + fmt.Sprintf(`
	SELECT
		g.idGrupo, g.nombre, g.numeroResolucion, g.lineaInvestigacion, g.tipoInvestigacion, g.fechaRegistro, g.archivo, g.createdAt, g.updatedAt,
		i.idInvestigador, i.nombre, i.apellido, i.createdAt as invCreatedAt, i.updatedAt as invUpdatedAt,
		dgi.rol
	FROM grupo g
	JOIN Grupo_Investigador dgi ON g.idGrupo = dgi.idGrupo
	JOIN investigador i ON dgi.idInvestigador = i.idInvestigador
	WHERE g.idGrupo IN (SELECT idGrupo FROM FilteredGroups)
	ORDER BY g.idGrupo, i.idInvestigador
	LIMIT $%d OFFSET $%d`, placeholderCount, placeholderCount+1)

	finalArgs := append(args, limit, offset)
	rows, err := db.Query(dataQuery, finalArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("error searching groups page with details: %w", err)
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
			return nil, 0, fmt.Errorf("error scanning group/investigator row during search: %w", err)
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
		return nil, 0, fmt.Errorf("error after iterating through group search rows: %w", err)
	}

	// Query for the total count using the CTE
	var total int
	countQuery := cteQuery + ` SELECT COUNT(*) FROM FilteredGroups`
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil { // Use original args for count
		return nil, 0, fmt.Errorf("error searching total group count: %w", err)
	}

	// Convert []*models.GrupoWithInvestigadores to []models.GrupoWithInvestigadores
	result := make([]models.GrupoWithInvestigadores, len(orderedGrupos))
	for i, ptr := range orderedGrupos {
		result[i] = *ptr
	}

	return result, total, nil
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
