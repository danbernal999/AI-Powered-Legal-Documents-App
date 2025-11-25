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
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Type        string                   `json:"type"`
	Content     string                   `json:"content"`
	Variables   []map[string]interface{} `json:"variables"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type Document struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`
	TemplateID string                 `json:"template_id"`
	Title      string                 `json:"title"`
	Type       string                 `json:"type"`
	Content    string                 `json:"content"`
	Variables  map[string]interface{} `json:"variables"`
	Status     string                 `json:"status"`
	Version    int                    `json:"version"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
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

type GoogleLoginRequest struct {
	IDToken string `json:"id_token"`
}


type Signature struct {
	ID            string    `json:"id"`
	DocumentID    string    `json:"document_id"`
	SignerName    string    `json:"signer_name"`
	SignerEmail   string    `json:"signer_email"`
	SignatureData string    `json:"signature_data"`
	Timestamp     string    `json:"timestamp"`
	SignedAt      time.Time `json:"signed_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SignDocumentRequest struct {
	SignerName    string `json:"signer_name"`
	SignerEmail   string `json:"signer_email"`
	SignatureData string `json:"signature_data"`
}

type DocumentShare struct {
	ID               string    `json:"id"`
	DocumentID       string    `json:"document_id"`
	SharedByUserID   string    `json:"shared_by_user_id"`
	SharedWithUserID string    `json:"shared_with_user_id"`
	Permission       string    `json:"permission"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ShareLink struct {
	ID            string     `json:"id"`
	DocumentID    string     `json:"document_id"`
	CreatedByUser string     `json:"created_by_user_id"`
	Token         string     `json:"token"`
	Permission    string     `json:"permission"`
	ExpiresAt     *time.Time `json:"expires_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type DocumentComment struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DocumentVersion struct {
	ID              string    `json:"id"`
	DocumentID      string    `json:"document_id"`
	UserID          string    `json:"user_id"`
	UserName        string    `json:"user_name"`
	PreviousContent string    `json:"previous_content"`
	CurrentContent  string    `json:"current_content"`
	ChangeSummary   string    `json:"change_summary"`
	CreatedAt       time.Time `json:"created_at"`
}

type ShareDocumentRequest struct {
	SharedWithEmail string `json:"shared_with_email"`
	Permission      string `json:"permission"`
}

type CreateShareLinkRequest struct {
	Permission string     `json:"permission"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

type CommentRequest struct {
	Content string `json:"content"`
}
