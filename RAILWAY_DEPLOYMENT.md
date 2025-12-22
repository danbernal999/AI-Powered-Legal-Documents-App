# 🚂 Guía de Despliegue en Railway - KiraDoc

Esta guía te llevará paso a paso para desplegar tu aplicación **KiraDoc** en Railway.

## 📋 Prerrequisitos

Antes de comenzar, asegúrate de tener:

- ✅ Cuenta en [Railway](https://railway.app/) (puedes usar GitHub para registrarte)
- ✅ [Railway CLI](https://docs.railway.app/develop/cli) instalado (opcional pero recomendado)
- ✅ Repositorio Git con tu código
- ✅ API Keys necesarias:
  - OpenAI API Key
  - Groq API Key (opcional)
  - Firebase credentials JSON

## 🏗️ Arquitectura de Despliegue

Railway desplegará **3 servicios** separados:

1. **PostgreSQL** - Base de datos
2. **Backend** - API Go (puerto 8080)
3. **Frontend** - Next.js (puerto 3000)

## 🚀 Paso 1: Preparar el Repositorio

### 1.1 Asegurar que tienes todos los archivos de configuración

Verifica que existan estos archivos en tu proyecto:

```
AI-Powered-Legal-Documents-App/
├── railway.json                    # ✅ Configuración de Railway
├── backend/
│   ├── nixpacks.toml              # ✅ Configuración de build para Go
│   ├── Dockerfile                 # ✅ Ya existe
│   └── migrations/                # ✅ Migraciones SQL
└── frontend/
    └── next.config.js             # ✅ Configuración de Next.js
```

### 1.2 Subir cambios a Git

```bash
git add .
git commit -m "Add Railway deployment configuration"
git push origin main
```

## 🎯 Paso 2: Crear Proyecto en Railway

### Opción A: Desde la Web UI

1. Ve a [railway.app](https://railway.app/)
2. Haz clic en **"New Project"**
3. Selecciona **"Deploy from GitHub repo"**
4. Autoriza Railway para acceder a tu repositorio
5. Selecciona el repositorio `AI-Powered-Legal-Documents-App`

### Opción B: Desde Railway CLI

```bash
# Instalar Railway CLI (si no lo tienes)
npm i -g @railway/cli

# Login
railway login

# Crear proyecto
railway init

# Vincular con tu repositorio
railway link
```

## 🗄️ Paso 3: Configurar PostgreSQL

### 3.1 Agregar PostgreSQL

1. En tu proyecto de Railway, haz clic en **"+ New"**
2. Selecciona **"Database"** → **"Add PostgreSQL"**
3. Railway creará automáticamente la base de datos

### 3.2 Obtener la URL de conexión

Railway generará automáticamente estas variables:
- `DATABASE_URL` - URL completa de conexión
- `PGHOST`, `PGPORT`, `PGUSER`, `PGPASSWORD`, `PGDATABASE` - Componentes individuales

**No necesitas configurar nada más**, Railway las hace disponibles automáticamente.

## ⚙️ Paso 4: Configurar el Backend (Go)

### 4.1 Crear servicio de Backend

1. En tu proyecto, haz clic en **"+ New"** → **"GitHub Repo"**
2. Selecciona tu repositorio
3. Railway detectará automáticamente que es un proyecto Go

### 4.2 Configurar el Root Directory

1. Ve a **Settings** del servicio Backend
2. En **"Root Directory"**, establece: `backend`
3. En **"Build Command"**, Railway usará automáticamente `nixpacks.toml`

### 4.3 Configurar Variables de Entorno

En la pestaña **"Variables"** del servicio Backend, agrega:

```bash
# Base de datos (Railway la vincula automáticamente)
DATABASE_URL=${{Postgres.DATABASE_URL}}

# Puerto
PORT=8080

# JWT Secret (genera uno seguro)
JWT_SECRET=tu-secret-key-super-segura-aqui-cambiar

# OpenAI
OPENAI_API_KEY=sk-tu-api-key-de-openai

# Groq (opcional)
GROQ_API_KEY=tu-groq-api-key

# Firebase Credentials (ver siguiente sección)
GOOGLE_APPLICATION_CREDENTIALS=/app/firebase-credentials.json
FIREBASE_CREDENTIALS_JSON={"type":"service_account",...}
```

### 4.4 Configurar Firebase Credentials

Railway no soporta archivos directamente, así que usaremos una variable de entorno:

1. Abre tu archivo `firebase-credentials.json`
2. Copia **TODO** el contenido (debe ser un JSON válido)
3. En Railway, crea la variable `FIREBASE_CREDENTIALS_JSON` con el contenido completo
4. Agrega esta variable: `GOOGLE_APPLICATION_CREDENTIALS=/app/firebase-credentials.json`

**Nota**: Necesitarás modificar ligeramente el código del backend para escribir el JSON a un archivo temporal. Ver sección de "Ajustes de Código" más abajo.

### 4.5 Configurar Health Check

1. Ve a **Settings** → **"Health Check"**
2. Establece:
   - **Path**: `/health`
   - **Timeout**: 10 segundos

### 4.6 Exponer el servicio

1. Ve a **Settings** → **"Networking"**
2. Haz clic en **"Generate Domain"**
3. Copia la URL generada (algo como `https://backend-production-xxxx.up.railway.app`)

## 🎨 Paso 5: Configurar el Frontend (Next.js)

### 5.1 Crear servicio de Frontend

1. En tu proyecto, haz clic en **"+ New"** → **"GitHub Repo"**
2. Selecciona el mismo repositorio
3. Railway detectará Next.js automáticamente

### 5.2 Configurar el Root Directory

1. Ve a **Settings** del servicio Frontend
2. En **"Root Directory"**, establece: `frontend`

### 5.3 Configurar Variables de Entorno

En la pestaña **"Variables"** del servicio Frontend:

```bash
# URL del Backend (usa la URL que generaste en el paso 4.6)
NEXT_PUBLIC_API_URL=https://backend-production-xxxx.up.railway.app/api/v1

# Firebase Config (del frontend)
NEXT_PUBLIC_FIREBASE_API_KEY=tu-firebase-api-key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=tu-proyecto.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=tu-proyecto-id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=tu-proyecto.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=123456789
NEXT_PUBLIC_FIREBASE_APP_ID=1:123456789:web:abcdef
```

### 5.4 Configurar Build Settings

Railway usará automáticamente:
- **Build Command**: `npm run build`
- **Start Command**: `npm start`

### 5.5 Exponer el servicio

1. Ve a **Settings** → **"Networking"**
2. Haz clic en **"Generate Domain"**
3. Esta será tu URL pública (ej: `https://frontend-production-xxxx.up.railway.app`)

## 🔧 Paso 6: Ajustes de Código Necesarios

### 6.1 Modificar Backend para Firebase Credentials

Crea o modifica `backend/pkg/config/firebase.go`:

```go
package config

import (
    "context"
    "encoding/json"
    "os"
    "path/filepath"
    
    firebase "firebase.google.com/go/v4"
    "google.golang.org/api/option"
)

func InitFirebase() (*firebase.App, error) {
    ctx := context.Background()
    
    // Check if we're in Railway (using env var)
    credJSON := os.Getenv("FIREBASE_CREDENTIALS_JSON")
    
    if credJSON != "" {
        // Railway: write JSON to temp file
        tmpDir := os.TempDir()
        credPath := filepath.Join(tmpDir, "firebase-credentials.json")
        
        if err := os.WriteFile(credPath, []byte(credJSON), 0600); err != nil {
            return nil, err
        }
        
        opt := option.WithCredentialsFile(credPath)
        return firebase.NewApp(ctx, nil, opt)
    }
    
    // Local: use file directly
    credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
    if credPath == "" {
        credPath = "../firebase-credentials.json"
    }
    
    opt := option.WithCredentialsFile(credPath)
    return firebase.NewApp(ctx, nil, opt)
}
```

### 6.2 Actualizar CORS en el Backend

Asegúrate de que el backend permita requests desde tu dominio de Railway.

En `backend/cmd/main.go` o donde configures CORS:

```go
c := cors.New(cors.Options{
    AllowedOrigins: []string{
        "http://localhost:3000",
        "https://frontend-production-xxxx.up.railway.app", // Tu dominio de Railway
    },
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"*"},
    AllowCredentials: true,
})
```

## 🔄 Paso 7: Ejecutar Migraciones

### Opción A: Automático (Recomendado)

Modifica `backend/cmd/main.go` para ejecutar migraciones al inicio:

```go
import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(databaseURL string) error {
    m, err := migrate.New(
        "file://migrations",
        databaseURL,
    )
    if err != nil {
        return err
    }
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    
    return nil
}

func main() {
    // ... load env vars ...
    
    // Run migrations
    if err := runMigrations(os.Getenv("DATABASE_URL")); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }
    
    // ... rest of your code ...
}
```

### Opción B: Manual desde Railway CLI

```bash
# Conectarse a la base de datos
railway connect postgres

# Ejecutar migraciones manualmente
psql $DATABASE_URL -f backend/migrations/001_initial_schema.up.sql
```

## ✅ Paso 8: Verificar el Despliegue

### 8.1 Verificar Backend

1. Ve a la URL del backend: `https://backend-production-xxxx.up.railway.app/health`
2. Deberías ver una respuesta exitosa

### 8.2 Verificar Frontend

1. Ve a la URL del frontend: `https://frontend-production-xxxx.up.railway.app`
2. Deberías ver la página de inicio de KiraDoc

### 8.3 Verificar Logs

En Railway, ve a cada servicio y revisa los **"Logs"** para asegurarte de que no hay errores.

## 🔍 Troubleshooting

### Problema: Backend no se conecta a la base de datos

**Solución**: Verifica que la variable `DATABASE_URL` esté correctamente vinculada:
```bash
DATABASE_URL=${{Postgres.DATABASE_URL}}
```

### Problema: Frontend no puede hacer requests al backend

**Solución**: 
1. Verifica que `NEXT_PUBLIC_API_URL` apunte a la URL correcta del backend
2. Verifica CORS en el backend

### Problema: Firebase credentials no funcionan

**Solución**:
1. Verifica que `FIREBASE_CREDENTIALS_JSON` contenga un JSON válido
2. Asegúrate de que el código escriba el archivo temporal correctamente

### Problema: Migraciones no se ejecutan

**Solución**:
1. Verifica que la carpeta `migrations` esté incluida en el build
2. Ejecuta migraciones manualmente usando Railway CLI

### Ver logs en tiempo real

```bash
# Backend
railway logs --service backend

# Frontend
railway logs --service frontend

# Postgres
railway logs --service postgres
```

## 🎉 ¡Listo!

Tu aplicación KiraDoc ahora está desplegada en Railway. Puedes:

- 🌐 Compartir la URL del frontend con otros
- 📊 Monitorear el uso en el dashboard de Railway
- 🔄 Hacer push a tu repositorio para desplegar automáticamente
- 📈 Escalar tus servicios según sea necesario

## 📚 Recursos Adicionales

- [Railway Docs](https://docs.railway.app/)
- [Railway CLI](https://docs.railway.app/develop/cli)
- [Railway Templates](https://railway.app/templates)
- [Nixpacks](https://nixpacks.com/)

## 💰 Costos

Railway ofrece:
- **$5 de crédito gratis** al mes para el plan Hobby
- **$500 horas de ejecución** gratis
- Después de eso, pagas por uso

Para una app pequeña como KiraDoc, el plan gratuito debería ser suficiente para desarrollo y pruebas.

---

**¿Necesitas ayuda?** Revisa los logs en Railway o contacta al soporte de Railway.
