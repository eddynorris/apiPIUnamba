-- Table: Investigador (Researchers)
CREATE TABLE Investigador (
    idInvestigador SERIAL PRIMARY KEY, -- SERIAL is PostgreSQL's auto-incrementing integer
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Sets timestamp on creation only
);

-- Table: Grupo (Research Groups)
CREATE TABLE Grupo (
    idGrupo SERIAL PRIMARY KEY,
    nombre VARCHAR(150) NOT NULL,
    numeroResolucion VARCHAR(100) NOT NULL,
    lineaInvestigacion VARCHAR(200) NOT NULL,
    tipoInvestigacion VARCHAR(100) NOT NULL,
    fechaRegistro DATE NOT NULL,
    archivo VARCHAR(255), -- Assuming this stores a file path or name
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Sets timestamp on creation only
);

-- Table: Grupo_Investigador (Associative table for Groups and Researchers)
CREATE TABLE Grupo_Investigador (
    idGrupo_Investigador SERIAL PRIMARY KEY,
    idGrupo INT NOT NULL,
    idInvestigador INT NOT NULL,
    rol VARCHAR(50) NOT NULL, -- e.g., 'Coordinador' or 'Integrante'
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Sets timestamp on creation only
    FOREIGN KEY (idGrupo) REFERENCES Grupo(idGrupo) ON DELETE CASCADE,
    FOREIGN KEY (idInvestigador) REFERENCES Investigador(idInvestigador) ON DELETE CASCADE
);