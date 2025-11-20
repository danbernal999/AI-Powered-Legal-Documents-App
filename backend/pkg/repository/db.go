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

func (r *Repository) CreateTemplate(t *Template) error {
	query := `INSERT INTO templates (id, name, description, type, content, variables) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, t.ID, t.Name, t.Description, t.Type, t.Content, t.Variables)
	return err
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

type Signature struct {
	ID            string
	DocumentID    string
	SignerName    string
	SignerEmail   string
	SignatureData []byte
	Timestamp     string
	SignedAt      interface{}
	CreatedAt     interface{}
	UpdatedAt     interface{}
}

func (r *Repository) CreateSignature(sig *Signature) error {
	query := `
	INSERT INTO signatures (id, document_id, signer_name, signer_email, signature_data, timestamp, signed_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, sig.ID, sig.DocumentID, sig.SignerName, sig.SignerEmail, sig.SignatureData, sig.Timestamp, sig.SignedAt)
	return err
}

func (r *Repository) GetSignaturesByDocumentID(documentID string) ([]*Signature, error) {
	query := `SELECT id, document_id, signer_name, signer_email, signature_data, timestamp, signed_at, created_at, updated_at FROM signatures WHERE document_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signatures []*Signature
	for rows.Next() {
		var sig Signature
		err := rows.Scan(&sig.ID, &sig.DocumentID, &sig.SignerName, &sig.SignerEmail, &sig.SignatureData, &sig.Timestamp, &sig.SignedAt, &sig.CreatedAt, &sig.UpdatedAt)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, &sig)
	}
	return signatures, nil
}

func (r *Repository) GetSignatureByID(id string) (*Signature, error) {
	var sig Signature
	query := `SELECT id, document_id, signer_name, signer_email, signature_data, timestamp, signed_at, created_at, updated_at FROM signatures WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&sig.ID, &sig.DocumentID, &sig.SignerName, &sig.SignerEmail, &sig.SignatureData, &sig.Timestamp, &sig.SignedAt, &sig.CreatedAt, &sig.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &sig, nil
}

func (r *Repository) DeleteSignature(id string) error {
	query := `DELETE FROM signatures WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

type DocumentShare struct {
	ID               string      `json:"id"`
	DocumentID       string      `json:"document_id"`
	SharedByUserID   string      `json:"shared_by_user_id"`
	SharedWithUserID string      `json:"shared_with_user_id"`
	Permission       string      `json:"permission"`
	CreatedAt        interface{} `json:"created_at"`
	UpdatedAt        interface{} `json:"updated_at"`
}

func (r *Repository) CreateDocumentShare(share *DocumentShare) error {
	query := `
	INSERT INTO document_shares (id, document_id, shared_by_user_id, shared_with_user_id, permission)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, share.ID, share.DocumentID, share.SharedByUserID, share.SharedWithUserID, share.Permission)
	return err
}

func (r *Repository) GetDocumentShares(documentID string) ([]*DocumentShare, error) {
	query := `SELECT id, document_id, shared_by_user_id, shared_with_user_id, permission, created_at, updated_at FROM document_shares WHERE document_id = $1`
	rows, err := r.db.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []*DocumentShare
	for rows.Next() {
		var share DocumentShare
		err := rows.Scan(&share.ID, &share.DocumentID, &share.SharedByUserID, &share.SharedWithUserID, &share.Permission, &share.CreatedAt, &share.UpdatedAt)
		if err != nil {
			return nil, err
		}
		shares = append(shares, &share)
	}
	return shares, nil
}

func (r *Repository) GetDocumentsSharedWithUser(userID string) ([]*DocumentShare, error) {
	query := `SELECT id, document_id, shared_by_user_id, shared_with_user_id, permission, created_at, updated_at FROM document_shares WHERE shared_with_user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []*DocumentShare
	for rows.Next() {
		var share DocumentShare
		err := rows.Scan(&share.ID, &share.DocumentID, &share.SharedByUserID, &share.SharedWithUserID, &share.Permission, &share.CreatedAt, &share.UpdatedAt)
		if err != nil {
			return nil, err
		}
		shares = append(shares, &share)
	}
	return shares, nil
}

func (r *Repository) UpdateDocumentShare(share *DocumentShare) error {
	query := `UPDATE document_shares SET permission = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, share.Permission, share.ID)
	return err
}

func (r *Repository) DeleteDocumentShare(id string) error {
	query := `DELETE FROM document_shares WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

type ShareLink struct {
	ID             string
	DocumentID     string
	CreatedByUserID string
	Token          string
	Permission     string
	ExpiresAt      interface{}
	CreatedAt      interface{}
	UpdatedAt      interface{}
}

func (r *Repository) CreateShareLink(link *ShareLink) error {
	query := `
	INSERT INTO share_links (id, document_id, created_by_user_id, token, permission, expires_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, link.ID, link.DocumentID, link.CreatedByUserID, link.Token, link.Permission, link.ExpiresAt)
	return err
}

func (r *Repository) GetShareLinkByToken(token string) (*ShareLink, error) {
	var link ShareLink
	query := `SELECT id, document_id, created_by_user_id, token, permission, expires_at, created_at, updated_at FROM share_links WHERE token = $1`
	err := r.db.QueryRow(query, token).Scan(&link.ID, &link.DocumentID, &link.CreatedByUserID, &link.Token, &link.Permission, &link.ExpiresAt, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *Repository) GetShareLinks(documentID string) ([]*ShareLink, error) {
	query := `SELECT id, document_id, created_by_user_id, token, permission, expires_at, created_at, updated_at FROM share_links WHERE document_id = $1`
	rows, err := r.db.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*ShareLink
	for rows.Next() {
		var link ShareLink
		err := rows.Scan(&link.ID, &link.DocumentID, &link.CreatedByUserID, &link.Token, &link.Permission, &link.ExpiresAt, &link.CreatedAt, &link.UpdatedAt)
		if err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	return links, nil
}

func (r *Repository) DeleteShareLink(id string) error {
	query := `DELETE FROM share_links WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

type DocumentComment struct {
	ID        string
	DocumentID string
	UserID    string
	Content   string
	CreatedAt interface{}
	UpdatedAt interface{}
}

func (r *Repository) CreateDocumentComment(comment *DocumentComment) error {
	query := `
	INSERT INTO document_comments (id, document_id, user_id, content)
	VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, comment.ID, comment.DocumentID, comment.UserID, comment.Content)
	return err
}

func (r *Repository) GetDocumentComments(documentID string) ([]*DocumentComment, error) {
	query := `
	SELECT dc.id, dc.document_id, dc.user_id, dc.content, dc.created_at, dc.updated_at
	FROM document_comments dc
	WHERE dc.document_id = $1
	ORDER BY dc.created_at DESC
	`
	rows, err := r.db.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*DocumentComment
	for rows.Next() {
		var comment DocumentComment
		err := rows.Scan(&comment.ID, &comment.DocumentID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (r *Repository) DeleteDocumentComment(id string) error {
	query := `DELETE FROM document_comments WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

type DocumentVersion struct {
	ID              string
	DocumentID      string
	UserID          string
	PreviousContent string
	CurrentContent  string
	ChangeSummary   string
	CreatedAt       interface{}
}

func (r *Repository) CreateDocumentVersion(version *DocumentVersion) error {
	query := `
	INSERT INTO document_versions (id, document_id, user_id, previous_content, current_content, change_summary)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, version.ID, version.DocumentID, version.UserID, version.PreviousContent, version.CurrentContent, version.ChangeSummary)
	return err
}

func (r *Repository) GetDocumentVersions(documentID string) ([]*DocumentVersion, error) {
	query := `SELECT id, document_id, user_id, previous_content, current_content, change_summary, created_at FROM document_versions WHERE document_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*DocumentVersion
	for rows.Next() {
		var version DocumentVersion
		err := rows.Scan(&version.ID, &version.DocumentID, &version.UserID, &version.PreviousContent, &version.CurrentContent, &version.ChangeSummary, &version.CreatedAt)
		if err != nil {
			return nil, err
		}
		versions = append(versions, &version)
	}
	return versions, nil
}

