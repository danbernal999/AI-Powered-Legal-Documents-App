# KiraDoc — Setup & Migrations

Este documento describe cómo preparar el proyecto localmente, ejecutar la base de datos y aplicar migraciones.

**Requisitos previos**
- Docker & Docker Compose (o Docker Desktop)
- Go 1.21 (para construir el backend localmente)
- Node.js (para el frontend)
- Opcional: `migrate` CLI (https://github.com/golang-migrate/migrate) si prefieres ejecutarlo localmente

**Resumen rápido (Makefile disponible)**
- `make up` — levanta la pila (Postgres + backend)
- `make db-migrate` — aplica las migraciones (usa binario local `migrate` si está instalado, si no usa la imagen Docker de `migrate`)
- `make db-down` — revierte las migraciones
- `make down` — baja la pila Docker

1) Copia el archivo de ejemplo de variables de entorno

```powershell
cp .env.example .env
# Edita .env si necesitas cambiar valores locales
```

2) Levantar la base y servicios (desarrollo)

```powershell
# Levanta servicios (el backend aplicará migraciones al arrancar)
make up
```

3) Aplicar migraciones manualmente (opcional)

```powershell
# Usando Makefile (prefiere binario local o Docker)
make db-migrate

# O con el binario migrate si lo tienes instalado:
cd backend
migrate -path ./migrations -database "$env:DATABASE_URL" up
```

4) Parar la pila

```powershell
make down
```

**Detalles importantes**
- Las migraciones están en `backend/migrations/` (archivos SQL `1_init.up.sql` y `1_init.down.sql`).
- El contenedor del backend copia `/migrations` dentro de la imagen y ejecuta `migrate` al iniciar (útil en dev). En producción, ejecuta migraciones como paso explícito en tu CI/CD — no confíes en migraciones automáticas sin pruebas.
- El antiguo método que hacía migraciones desde código fue removido: ahora usamos `golang-migrate` y archivos SQL versionados.

**Uso en Windows/WSL**
- Ejecuta los comandos desde WSL para mejor compatibilidad con Make y mount de volúmenes Docker.

**Recomendaciones de seguridad**
- No expongas secretos en logs. El backend ya enmascara las credenciales del DSN en los logs.
- Mantén `.env` fuera del control de versiones (usa `.env.example` en repo).

**Comandos útiles**
- Ver la URL que usa la app (dentro del contenedor backend):

```powershell
docker compose exec backend env | Select-String DATABASE_URL
```

- Listar bases en Postgres (conectando a la base `kiradoc`):

```powershell
docker compose exec postgres psql -U user -d kiradoc -c "\l"
```

- Ejecutar migraciones con Docker (si no tienes `migrate` instalado):

```powershell
docker run --rm -v ${PWD}/backend/migrations:/migrations -e DATABASE_URL="$env:DATABASE_URL" migrate/migrate:v4.15.2 -path=/migrations -database "$env:DATABASE_URL" up
```

¿Quieres que añada un job de CI (GitHub Actions) que verifique y aplique/valide migraciones en un entorno de pruebas? Si es así, lo puedo crear y probar con un contenedor de Postgres en el workflow.
