package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/danbernal999/kiradoc-backend/pkg/handlers"
	"github.com/danbernal999/kiradoc-backend/pkg/middleware"
	"github.com/danbernal999/kiradoc-backend/pkg/repository"
	"github.com/danbernal999/kiradoc-backend/pkg/services"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")

	dbURL := os.Getenv("DATABASE_URL")
	// Log DB host and database only (mask credentials) for local debug
	masked := maskDSN(dbURL)
	log.Printf("DB connection: %q", masked)
	if dbURL == "" {
		dbURL = "postgres://user:password@postgres:5432/kiradoc?sslmode=disable"
	}

	// Run migrations using golang-migrate (migrations copied into /migrations in the image)
	// we run migrations before opening DB connection
	if dbURL != "" {
		// migrate uses the database URL directly
		m, err := migrate.New("file:///migrations", dbURL)
		if err != nil {
			log.Printf("migrate.New error: %v", err)
		} else {
			err = m.Up()
			if err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Migration failed: %v", err)
			}
			log.Println("DB migrations applied (or no change)")
		}
	}

	db, err := repository.NewDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	repo := repository.NewRepository(db)

	if err := repo.SeedTemplatesIfEmpty(); err != nil {
		log.Printf("Warning: Failed to seed templates: %v", err)
	}
	svc := services.NewServices(repo)
	h := handlers.NewHandlers(svc)

	router := mux.NewRouter()

	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET", "OPTIONS")

	authRoutes := router.PathPrefix("/api/v1/auth").Subrouter()
	authRoutes.HandleFunc("/register", h.RegisterHandler).Methods("POST", "OPTIONS")
	authRoutes.HandleFunc("/login", h.LoginHandler).Methods("POST", "OPTIONS")

	apiRoutes := router.PathPrefix("/api/v1").Subrouter()
	apiRoutes.Use(middleware.JWTMiddleware)

	apiRoutes.HandleFunc("/templates", h.GetTemplatesHandler).Methods("GET", "OPTIONS")
	apiRoutes.HandleFunc("/templates", h.CreateTemplateHandler).Methods("POST", "OPTIONS")
	apiRoutes.HandleFunc("/templates/{id}", h.GetTemplateHandler).Methods("GET", "OPTIONS")

	apiRoutes.HandleFunc("/generate", h.GenerateDocumentHandler).Methods("POST", "OPTIONS")

	// Document routes - more specific paths first
	apiRoutes.HandleFunc("/documents/{id}/sign", h.SignDocumentHandler).Methods("POST", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}/signatures", h.GetDocumentSignaturesHandler).Methods("GET", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}/share", h.ShareDocumentHandler).Methods("POST", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}/shares", h.GetDocumentSharesHandler).Methods("GET", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}/comments", h.GetDocumentCommentsHandler).Methods("GET", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}/comments", h.CreateDocumentCommentHandler).Methods("POST", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}/versions", h.GetDocumentVersionsHandler).Methods("GET", "OPTIONS")

	// Generic document routes - less specific paths last
	apiRoutes.HandleFunc("/documents", h.ListDocumentsHandler).Methods("GET", "OPTIONS")
	apiRoutes.HandleFunc("/documents", h.SaveDocumentHandler).Methods("POST", "OPTIONS")
	apiRoutes.HandleFunc("/documents/shared-with-me", h.GetSharedWithMeHandler).Methods("GET", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}", h.GetDocumentHandler).Methods("GET", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}", h.UpdateDocumentHandler).Methods("PUT", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}", h.DeleteDocumentHandler).Methods("DELETE", "OPTIONS")

	// Standalone resource routes
	apiRoutes.HandleFunc("/signatures/{id}", h.DeleteSignatureHandler).Methods("DELETE", "OPTIONS")
	apiRoutes.HandleFunc("/shares/{id}", h.DeleteDocumentShareHandler).Methods("DELETE", "OPTIONS")
	apiRoutes.HandleFunc("/comments/{id}", h.DeleteDocumentCommentHandler).Methods("DELETE", "OPTIONS")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	corsRouter := middleware.CORSMiddleware(router)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsRouter))
}

// maskDSN returns a redacted DSN showing only host and database (safe for logs)
func maskDSN(dsn string) string {
	if dsn == "" {
		return "(empty)"
	}
	// crude parse: look for @ then /dbname
	at := strings.Index(dsn, "@")
	if at == -1 {
		return "(invalid dsn)"
	}
	afterAt := dsn[at+1:]
	slash := strings.Index(afterAt, "/")
	if slash == -1 {
		return "(invalid dsn)"
	}
	hostPart := afterAt[:slash]
	rest := afterAt[slash+1:]
	q := strings.Index(rest, "?")
	dbName := rest
	if q != -1 {
		dbName = rest[:q]
	}
	return hostPart + "/" + dbName
}
