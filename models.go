package main

import "time"

// Investigador represents an investigator in the database.
type Investigador struct {
	ID        int    `json:"idInvestigador" db:"idInvestigador"`
	Nombre    string `json:"nombre" db:"nombre"`
	Apellido  string `json:"apellido" db:"apellido"`
	Rol       string `json:"rol" db:"rol"`
}

// Grupo represents a research group in the database.
type Grupo struct {
	ID                 int       `json:"idGrupo" db:"idGrupo"`
	Nombre             string    `json:"nombre" db:"nombre"`
	NumeroResolucion   string    `json:"numeroResolucion" db:"numeroResolucion"`
	LineaInvestigacion string    `json:"lineaInvestigacion" db:"lineaInvestigacion"`
	TipoInvestigacion  string    `json:"tipoInvestigacion" db:"tipoInvestigacion"`
	FechaRegistro      time.Time `json:"fechaRegistro" db:"fechaRegistro"`
	Archivo            string    `json:"archivo" db:"archivo"`
}

// DetalleGrupoInvestigador represents the relationship between a group and an investigator.
type DetalleGrupoInvestigador struct {
	ID           int    `json:"idDetalleGI" db:"idDetalleGI"`
	IDGrupo      int    `json:"idGrupo" db:"idGrupo"`
	IDInvestigador int    `json:"idInvestigador" db:"idInvestigador"`
	TipoRelacion string `json:"tipoRelacion" db:"tipoRelacion"`
}
