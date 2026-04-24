# SeriesTracker Backend

API REST construida en Go con SQLite para gestionar series.

## Correr localmente

Requisitos:

- Go 1.21+
- GCC o MinGW-w64 por `go-sqlite3`

Pasos:

```bash
cd backend
go run main.go
```

El backend queda disponible en `http://localhost:8080`.

Endpoints principales:

- `GET /health`
- `GET /series`
- `POST /series`
- `GET /series/:id`
- `PUT /series/:id`
- `DELETE /series/:id`

## Variables de entorno

- `PORT`: puerto HTTP. En local usa `8080` por defecto.
- `DATABASE_PATH`: ruta absoluta o relativa del archivo SQLite.
- `ALLOWED_ORIGINS`: lista separada por comas para CORS.

Ejemplo:

```env
PORT=8080
DATABASE_PATH=../series.db
ALLOWED_ORIGINS=http://localhost:3000
```

## Despliegue en Railway

Este backend ya quedó preparado para Railway con:

- [railway.toml](/c:/Users/chaar/OneDrive/Desktop/Archivos%20UVG/Semestre%205/Web/Proyecto1/backend/railway.toml)
- [Procfile](/c:/Users/chaar/OneDrive/Desktop/Archivos%20UVG/Semestre%205/Web/Proyecto1/backend/Procfile)

Pasos:

1. Sube el proyecto a GitHub.
2. Crea un proyecto en Railway desde ese repositorio.
3. En el servicio del backend configura `Root Directory = backend`.
4. Agrega un Volume, por ejemplo montado en `/data`.
5. Define `DATABASE_PATH=/data/series.db`.
6. Define `ALLOWED_ORIGINS=https://tu-frontend.vercel.app`.
7. Despliega y prueba `https://tu-backend.up.railway.app/health`.

Nota:

- Si no montas un Volume, la base SQLite puede perderse al redeployar.
- El backend también detecta `RAILWAY_VOLUME_MOUNT` y `RAILWAY_VOLUME_MOUNT_PATH`.
