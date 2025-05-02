package models

import "time"

// Grupo represents a research group in the database.
type Grupo struct {
	ID                 int       `json:"idGrupo" db:"idGrupo"`
	Nombre             string    `json:"nombre" db:"nombre"`
	NumeroResolucion   string    `json:"numeroResolucion" db:"numero_resolucion"`
	LineaInvestigacion string    `json:"lineaInvestigacion" db:"linea_investigacion"`
	TipoInvestigacion  string    `json:"tipoInvestigacion" db:"tipo_investigacion"`
	FechaRegistro      time.Time `json:"fechaRegistro" db:"fecha_registro"`
	Archivo            string    `json:"archivo" db:"archivo"`
	CreatedAt          time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" db:"updatedAt"`
}

// GrupoWithInvestigadores represents a group with its associated investigators.
type GrupoWithInvestigadores struct {
	Grupo          Grupo          `json:"grupo"`
	Investigadores []Investigador `json:"investigadores"`
}
