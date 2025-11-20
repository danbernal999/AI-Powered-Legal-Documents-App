package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danbernal999/kiradoc-backend/pkg/handlers"
	"github.com/danbernal999/kiradoc-backend/pkg/middleware"
	"github.com/danbernal999/kiradoc-backend/pkg/repository"
	"github.com/danbernal999/kiradoc-backend/pkg/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost/kiradoc"
	}

	db, err := repository.NewDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

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
	apiRoutes.HandleFunc("/templates/{id}", h.GetTemplateHandler).Methods("GET", "OPTIONS")

	apiRoutes.HandleFunc("/generate", h.GenerateDocumentHandler).Methods("POST", "OPTIONS")

	apiRoutes.HandleFunc("/documents", h.ListDocumentsHandler).Methods("GET", "OPTIONS")
	apiRoutes.HandleFunc("/documents", h.SaveDocumentHandler).Methods("POST", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}", h.GetDocumentHandler).Methods("GET", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}", h.UpdateDocumentHandler).Methods("PUT", "OPTIONS")
	apiRoutes.HandleFunc("/documents/{id}", h.DeleteDocumentHandler).Methods("DELETE", "OPTIONS")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	corsRouter := middleware.CORSMiddleware(router)

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsRouter))
}
