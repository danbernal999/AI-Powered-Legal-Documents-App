# Configuración Rápida de Firebase

## ✅ Paso 1: Frontend (Ya completado)

Tu `.env.local` ya tiene las credenciales correctas.

---

## 🔧 Paso 2: Backend - Configurar Service Account

El backend necesita credenciales de Firebase Admin SDK.

### Opción A: Para Desarrollo Local (Más Simple)

1. **Descargar Service Account Key:**
   - Ve a https://console.firebase.google.com/
   - Selecciona tu proyecto "kiradoc-1bb3f"
   - Ve a **Project Settings** (⚙️) → **Service accounts**
   - Click en **"Generate new private key"**
   - Click **"Generate key"** - se descargará un archivo JSON

2. **Guardar el archivo:**
   - Renombra el archivo descargado a: `firebase-service-account.json`
   - Muévelo a la carpeta `backend/` de tu proyecto
   - **IMPORTANTE**: Este archivo ya está en `.gitignore`, nunca lo subas a Git

3. **Configurar variable de entorno:**
   
   En Windows PowerShell:
   ```powershell
   $env:GOOGLE_APPLICATION_CREDENTIALS="C:\Users\danbe\OneDrive\Documentos\Projects\AI-Powered-Legal-Documents-App\backend\firebase-service-account.json"
   ```

4. **Ejecutar el backend:**
   ```bash
   cd backend
   go run cmd/main.go
   ```

   Deberías ver:
   ```
   Firebase initialized successfully
   ```

### Opción B: Sin Service Account (Solo para pruebas)

Si quieres probar sin configurar Firebase en el backend, puedes comentar temporalmente la inicialización:

En `backend/cmd/main.go`, líneas 63-69, comenta el código:

```go
// if err := handlers.InitializeFirebase(); err != nil {
//     log.Printf("Warning: Failed to initialize Firebase: %v", err)
//     log.Println("Google authentication will not be available")
// }
```

**NOTA**: Con esta opción, el login con Google NO funcionará en el backend, pero podrás probar el resto de la aplicación.

---

## 🚀 Paso 3: Reiniciar Todo

1. **Detén el frontend** (Ctrl+C)
2. **Inicia el frontend de nuevo:**
   ```bash
   cd frontend
   npm run dev
   ```

3. **En otra terminal, inicia el backend:**
   ```bash
   cd backend
   go run cmd/main.go
   ```

---

## ✅ Verificación

1. Frontend en http://localhost:3000 - NO debe mostrar error `auth/invalid-api-key`
2. Backend debe mostrar: `Firebase initialized successfully`
3. Click en "Continuar con Google" debe abrir el popup de Google

---

## 🐛 Si aún tienes problemas

**Error en frontend**: Verifica que `.env.local` tenga EXACTAMENTE estos valores:
```
NEXT_PUBLIC_FIREBASE_API_KEY=AIzaSyCEYk9zth9oItH-bzOGMyuq6YL59m8XQrI
```

**Error en backend**: Verifica la ruta del archivo JSON en la variable de entorno.
