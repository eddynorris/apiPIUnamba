package models

import "time"

// Investigador represents an investigator in the database.
type Investigador struct {
	ID       int    `json:"idInvestigador" db:"id_investigador"`
	Nombre   string `json:"nombre" db:"nombre"`
	Apellido string `json:"apellido" db:"apellido"`
	Rol      string `json:"rol" db:"rol"`
}

// Grupo represents a research group in the database.
type Grupo struct {
	ID                 int       `json:"idGrupo" db:"id_grupo"`
	Nombre             string    `json:"nombre" db:"nombre"`
	NumeroResolucion   string    `json:"numeroResolucion" db:"numero_resolucion"`
	LineaInvestigacion string    `json:"lineaInvestigacion" db:"linea_investigacion"`
	TipoInvestigacion  string    `json:"tipoInvestigacion" db:"tipo_investigacion"`
	FechaRegistro      time.Time `json:"fechaRegistro" db:"fecha_registro"`
	Archivo            string    `json:"archivo" db:"archivo"`
}

// DetalleGrupoInvestigador represents the relationship between a group and an investigator.
type DetalleGrupoInvestigador struct {
	ID             int    `json:"idDetalleGI" db:"id_detalle_gi"`
	IDGrupo        int    `json:"idGrupo" db:"id_grupo"`
	IDInvestigador int    `json:"idInvestigador" db:"id_investigador"`
	TipoRelacion   string `json:"tipoRelacion" db:"tipo_relacion"`
}

// GrupoWithInvestigadores represents a group with its associated investigators.
type GrupoWithInvestigadores struct {
	Grupo          Grupo          `json:"grupo"`
	Investigadores []Investigador `json:"investigadores"`
}
