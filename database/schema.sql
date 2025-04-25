-- Tabla Investigador
-- Nombres en minúscula y snake_case
CREATE TABLE investigador (
    id_investigador SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    rol VARCHAR(100) NOT NULL
);

-- Tabla Grupos de Investigación
CREATE TABLE grupo (
    id_grupo SERIAL PRIMARY KEY,
    nombre VARCHAR(150) NOT NULL,
    numero_resolucion VARCHAR(100) NOT NULL,
    linea_investigacion VARCHAR(200) NOT NULL,
    tipo_investigacion VARCHAR(100) NOT NULL,
    fecha_registro DATE NOT NULL,
    archivo VARCHAR(255) -- Puede ser NULL si no siempre hay archivo
);

-- Tabla Detalle_Grupo de Investigadores
CREATE TABLE detalle_grupo_investigador (
    id_detalle_gi SERIAL PRIMARY KEY,
    id_grupo INT NOT NULL,
    id_investigador INT NOT NULL,
    tipo_relacion VARCHAR(100) NOT NULL,
    -- Las referencias FOREIGN KEY también usan nombres en minúsculas
    FOREIGN KEY (id_grupo) REFERENCES grupo(id_grupo) ON DELETE CASCADE,
    FOREIGN KEY (id_investigador) REFERENCES investigador(id_investigador) ON DELETE CASCADE
);

-- Example INSERT statement (nombres de columna en minúscula)
-- Nota: No insertamos los IDs seriales, PostgreSQL los genera.
INSERT INTO investigador (nombre, apellido, rol) VALUES ('John', 'Doe', 'Professor');

-- Ejemplo asumiendo que los primeros IDs serán 1:
INSERT INTO grupo (nombre, numero_resolucion, linea_investigacion, tipo_investigacion, fecha_registro, archivo) VALUES ('Research Group 1', 'Res-2023-001', 'AI in Education', 'Applied', '2023-01-15', '/files/research_group_1_resolution.pdf');
INSERT INTO detalle_grupo_investigador (id_grupo, id_investigador, tipo_relacion) VALUES (1, 1, 'coordinator');
