package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (d *DB) Migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		name VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS templates (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		type VARCHAR(50) NOT NULL,
		content TEXT NOT NULL,
		variables TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS documents (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL REFERENCES users(id),
		template_id VARCHAR(36) REFERENCES templates(id),
		title VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		content TEXT NOT NULL,
		variables TEXT,
		status VARCHAR(50),
		version INTEGER DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_documents_user_id ON documents(user_id);
	CREATE INDEX IF NOT EXISTS idx_documents_template_id ON documents(template_id);
	`

	_, err := d.Exec(schema)
	return err
}

type Repository struct {
	db *DB
}

func NewRepository(db *DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(id, email, password, name string) error {
	query := `INSERT INTO users (id, email, password, name) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, id, email, password, name)
	return err
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	var user User
	query := `SELECT id, email, password, name, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id string) (*User, error) {
	var user User
	query := `SELECT id, email, password, name, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type User struct {
	ID        string
	Email     string
	Password  string
	Name      string
	CreatedAt interface{}
	UpdatedAt interface{}
}

func (r *Repository) GetTemplates() ([]*Template, error) {
	query := `SELECT id, name, description, type, content, variables, created_at, updated_at FROM templates`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*Template
	for rows.Next() {
		var t Template
		err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Type, &t.Content, &t.Variables, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &t)
	}
	return templates, nil
}

func (r *Repository) GetTemplateByID(id string) (*Template, error) {
	var t Template
	query := `SELECT id, name, description, type, content, variables, created_at, updated_at FROM templates WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&t.ID, &t.Name, &t.Description, &t.Type, &t.Content, &t.Variables, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

type Template struct {
	ID          string
	Name        string
	Description string
	Type        string
	Content     string
	Variables   string
	CreatedAt   interface{}
	UpdatedAt   interface{}
}

func (r *Repository) CreateDocument(doc *Document) error {
	query := `
	INSERT INTO documents (id, user_id, template_id, title, type, content, variables, status, version)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query, doc.ID, doc.UserID, doc.TemplateID, doc.Title, doc.Type, doc.Content, doc.Variables, doc.Status, doc.Version)
	return err
}

func (r *Repository) GetDocumentByID(id string) (*Document, error) {
	var doc Document
	query := `SELECT id, user_id, template_id, title, type, content, variables, status, version, created_at, updated_at FROM documents WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&doc.ID, &doc.UserID, &doc.TemplateID, &doc.Title, &doc.Type, &doc.Content, &doc.Variables, &doc.Status, &doc.Version, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *Repository) GetDocumentsByUserID(userID string) ([]*Document, error) {
	query := `SELECT id, user_id, template_id, title, type, content, variables, status, version, created_at, updated_at FROM documents WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []*Document
	for rows.Next() {
		var doc Document
		err := rows.Scan(&doc.ID, &doc.UserID, &doc.TemplateID, &doc.Title, &doc.Type, &doc.Content, &doc.Variables, &doc.Status, &doc.Version, &doc.CreatedAt, &doc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		documents = append(documents, &doc)
	}
	return documents, nil
}

func (r *Repository) UpdateDocument(doc *Document) error {
	query := `
	UPDATE documents
	SET title = $1, content = $2, variables = $3, status = $4, version = $5, updated_at = CURRENT_TIMESTAMP
	WHERE id = $6
	`
	_, err := r.db.Exec(query, doc.Title, doc.Content, doc.Variables, doc.Status, doc.Version, doc.ID)
	return err
}

func (r *Repository) DeleteDocument(id string) error {
	query := `DELETE FROM documents WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

type Document struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	TemplateID string      `json:"template_id"`
	Title      string      `json:"title"`
	Type       string      `json:"type"`
	Content    string      `json:"content"`
	Variables  string      `json:"variables"`
	Status     string      `json:"status"`
	Version    int         `json:"version"`
	CreatedAt  interface{} `json:"created_at"`
	UpdatedAt  interface{} `json:"updated_at"`
}
