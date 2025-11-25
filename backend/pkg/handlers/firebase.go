package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/danbernal999/kiradoc-backend/pkg/models"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

var firebaseAuth *auth.Client

// InitializeFirebase initializes Firebase Admin SDK
// Tries multiple approaches:
// 1. Uses service account file from GOOGLE_APPLICATION_CREDENTIALS env var
// 2. Falls back to Application Default Credentials
func InitializeFirebase() error {
	ctx := context.Background()

	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	
	var app *firebase.App
	var err error
	
	if credentialsPath != "" {
		// Try with the service account file from env var
		opt := option.WithCredentialsFile(credentialsPath)
		app, err = firebase.NewApp(ctx, nil, opt)
		if err != nil {
			log.Printf("Error initializing Firebase with credentials file (%s): %v", credentialsPath, err)
			// Try without credentials file as fallback
			app, err = firebase.NewApp(ctx, nil)
			if err != nil {
				log.Printf("Error initializing Firebase app (fallback): %v", err)
				return err
			}
			log.Println("Firebase initialized with default credentials (fallback)")
		} else {
			log.Println("Firebase initialized successfully with service account file")
		}
	} else {
		// Initialize Firebase App with default credentials
		app, err = firebase.NewApp(ctx, nil)
		if err != nil {
			log.Printf("Error initializing Firebase app: %v", err)
			return err
		}
	}

	// Get Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Printf("Error getting Auth client: %v", err)
		return err
	}

	firebaseAuth = authClient
	log.Println("Firebase Auth client initialized successfully")
	return nil
}

// InitializeFirebaseWithServiceAccount initializes Firebase with a service account key file
func InitializeFirebaseWithServiceAccount(serviceAccountPath string) error {
	ctx := context.Background()

	opt := option.WithCredentialsFile(serviceAccountPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("Error initializing Firebase app with service account: %v", err)
		return err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Printf("Error getting Auth client: %v", err)
		return err
	}

	firebaseAuth = authClient
	log.Println("Firebase initialized successfully with service account")
	return nil
}

// GoogleLoginHandler handles Google OAuth login via Firebase
func (h *Handlers) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GoogleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.IDToken == "" {
		http.Error(w, "ID token is required", http.StatusBadRequest)
		return
	}

	// Verify the ID token with Firebase
	if firebaseAuth == nil {
		http.Error(w, "Google authentication is not configured on the server", http.StatusServiceUnavailable)
		return
	}

	ctx := context.Background()
	token, err := firebaseAuth.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		log.Printf("Failed to verify ID token: %v", err)
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		return
	}

	// Extract user information from the token
	email, ok := token.Claims["email"].(string)
	if !ok || email == "" {
		log.Printf("Email not found in token claims")
		http.Error(w, "Email not found in token", http.StatusBadRequest)
		return
	}

	name, _ := token.Claims["name"].(string)
	if name == "" {
		name = email // Fallback to email if name is not provided
	}

	// Check if user exists in our database
	user, err := h.services.GetUserByEmail(email)
	if err != nil {
		// User doesn't exist, create a new one
		userID := uuid.New().String()

		// For Google auth users, we don't have a password
		// We'll use a random UUID as a placeholder password (they won't use it)
		placeholderPassword := uuid.New().String()

		if err := h.services.CreateUser(userID, email, placeholderPassword, name); err != nil {
			log.Printf("Failed to create user: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Fetch the newly created user
		user, err = h.services.GetUserByEmail(email)
		if err != nil {
			log.Printf("Failed to fetch created user: %v", err)
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
			return
		}
		log.Printf("Created new user from Google auth: %s (%s)", name, email)
	} else {
		log.Printf("Existing user logged in via Google: %s (%s)", user.Name, user.Email)
	}

	// Generate JWT token for our application
	appToken, err := generateToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: appToken,
		User:  models.User{ID: user.ID, Email: user.Email, Name: user.Name},
	})
}
