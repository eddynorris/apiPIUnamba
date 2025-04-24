-- Tabla Investigador
CREATE TABLE Investigador (
    idInvestigador INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    rol VARCHAR(100) NOT NULL
);

-- Tabla Grupos de Investigación
CREATE TABLE Grupo (
    idGrupo INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(150) NOT NULL,
    numeroResolucion VARCHAR(100) NOT NULL,
    lineaInvestigacion VARCHAR(200) NOT NULL,
    tipoInvestigacion VARCHAR(100) NOT NULL,
    fechaRegistro DATE NOT NULL,
    archivo VARCHAR(255)
);

-- Tabla Detalle_Grupo de Investigadores
CREATE TABLE Detalle_GrupoInvestigador (
    idDetalleGI INT AUTO_INCREMENT PRIMARY KEY,
    idGrupo INT NOT NULL,
    idInvestigador INT NOT NULL,
    tipoRelacion VARCHAR(100) NOT NULL,
    FOREIGN KEY (idGrupo) REFERENCES Grupo(idGrupo) ON DELETE CASCADE,
    FOREIGN KEY (idInvestigador) REFERENCES Investigador(idInvestigador) ON DELETE CASCADE
);

-- Example INSERT statement
INSERT INTO Investigador (nombre, apellido, rol) VALUES ('John', 'Doe', 'Professor');
INSERT INTO Grupo (nombre, numeroResolucion, lineaInvestigacion, tipoInvestigacion, fechaRegistro, archivo) VALUES ('Research Group 1', 'Res-2023-001', 'AI in Education', 'Applied', '2023-01-15', '/files/research_group_1_resolution.pdf');
INSERT INTO Detalle_GrupoInvestigador (idGrupo, idInvestigador, tipoRelacion) VALUES (1, 1, 'coordinator');
