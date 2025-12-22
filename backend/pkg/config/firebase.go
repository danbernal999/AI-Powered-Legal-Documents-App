package config

import (
	"context"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// InitFirebase initializes Firebase app with support for both local and Railway environments
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
