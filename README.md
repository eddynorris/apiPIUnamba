# API PI UNAMBA

API REST en Go para la gestión de investigadores y grupos de investigación. Permite administrar investigadores, grupos de investigación y la relación entre ambos (rol dentro del grupo: coordinador o integrante), orientada al sistema de investigación de la UNAMBA.

## Características principales

- CRUD de investigadores (`/investigadores`).
- CRUD de grupos de investigación (`/grupos`) con número de resolución, línea de investigación y tipo.
- Asociación de investigadores a grupos con rol y detalle (`/detalles`), con CRUD completo.
- Consulta de grupos por investigador y de detalles por grupo.
- Paginación de resultados.
- CORS configurable mediante variable de entorno y carga de configuración vía `.env`.

## Tecnologías usadas

- Go 1.21+
- gorilla/mux (enrutador HTTP)
- PostgreSQL (driver lib/pq)
- rs/cors
- godotenv

## Requisitos previos

- Go 1.21+
- PostgreSQL con la base de datos creada.

## Cómo ejecutar

1. Configurar las variables de entorno (`.env`):

   ```
   DB_USER=postgres
   DB_PASSWORD=123456
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=db_PIUnamba
   DB_SSLMODE=disable
   PORT=3000
   FRONTEND_URL=https://tu-frontend.com
   ```

2. Crear las tablas ejecutando `database/schema.sql` en tu PostgreSQL.

3. Ejecutar el servidor:

   ```bash
   go run main.go
   ```

El servidor escucha en el puerto 3000 por defecto (configurable con `PORT`).

## Estructura del proyecto

- `main.go` — punto de entrada, configuración de CORS y arranque del servidor.
- `routes/` — definición de rutas.
- `controllers/` — handlers HTTP.
- `models/` — estructuras de datos.
- `repository/` — acceso a datos.
- `database/` — conexión a PostgreSQL y esquema SQL.
- `utils/` — utilidades (paginación).