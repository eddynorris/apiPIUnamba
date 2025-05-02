package models

import "time"

// Investigador represents an investigator in the database.
type Investigador struct {
	ID        int       `json:"idInvestigador" db:"idInvestigador"`
	Nombre    string    `json:"nombre" db:"nombre"`
	Apellido  string    `json:"apellido" db:"apellido"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}
