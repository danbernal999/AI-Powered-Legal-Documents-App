package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Template struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Content     string                 `json:"content"`
	Variables   []map[string]interface{} `json:"variables"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type Document struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	TemplateID  string                 `json:"template_id"`
	Title       string                 `json:"title"`
	Type        string                 `json:"type"`
	Content     string                 `json:"content"`
	Variables   map[string]interface{} `json:"variables"`
	Status      string                 `json:"status"`
	Version     int                    `json:"version"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type GenerateDocumentRequest struct {
	TemplateID string                 `json:"template_id"`
	Variables  map[string]interface{} `json:"variables"`
	Title      string                 `json:"title"`
}

type GenerateDocumentResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
