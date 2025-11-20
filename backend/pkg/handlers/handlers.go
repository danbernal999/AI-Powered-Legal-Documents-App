package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/danbernal999/kiradoc-backend/pkg/middleware"
	"github.com/danbernal999/kiradoc-backend/pkg/models"
	"github.com/danbernal999/kiradoc-backend/pkg/repository"
	"github.com/danbernal999/kiradoc-backend/pkg/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handlers struct {
	services *services.Services
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{services: services}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID := uuid.New().String()
	if err := h.services.CreateUser(userID, req.Email, req.Password, req.Name); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	user := models.User{ID: userID, Email: req.Email, Name: req.Name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.services.GetUserByEmail(req.Email)
	if err != nil || !h.services.ComparePasswords(user.Password, req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		User:  models.User{ID: user.ID, Email: user.Email, Name: user.Name},
	})
}

func (h *Handlers) GetTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	repoTemplates, err := h.services.GetTemplates()
	if err != nil {
		http.Error(w, "Failed to get templates", http.StatusInternalServerError)
		return
	}

	templates := make([]*models.Template, 0)
	for _, rt := range repoTemplates {
		mt := convertRepoTemplateToModel(rt)
		templates = append(templates, mt)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func (h *Handlers) GetTemplateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	repoTemplate, err := h.services.GetTemplateByID(id)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	template := convertRepoTemplateToModel(repoTemplate)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

func (h *Handlers) GenerateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GenerateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	template, err := h.services.GetTemplateByID(req.TemplateID)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	generatedContent := template.Content
	for key, value := range req.Variables {
		placeholder := "{{" + key + "}}"
		generatedContent = replaceAll(generatedContent, placeholder, toString(value))
	}

	docID := uuid.New().String()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.GenerateDocumentResponse{
		ID:      docID,
		Content: generatedContent,
	})
}

func (h *Handlers) ListDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		log.Printf("Invalid user_id type: %T", userIDVal)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	documents, err := h.services.GetDocumentsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

func (h *Handlers) SaveDocumentHandler(w http.ResponseWriter, r *http.Request) {
	var doc models.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		log.Printf("Invalid user_id type: %T", userIDVal)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	doc.ID = uuid.New().String()
	doc.UserID = userID
	doc.Status = "draft"
	doc.Version = 1

	if err := h.services.CreateDocument(convertModelToRepo(&doc)); err != nil {
		log.Printf("Failed to create document: %v", err)
		http.Error(w, "Failed to save document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(doc)
}

func (h *Handlers) GetDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	document, err := h.services.GetDocumentByID(id)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *Handlers) UpdateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var doc models.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	doc.ID = id
	if err := h.services.UpdateDocument(convertModelToRepo(&doc)); err != nil {
		http.Error(w, "Failed to update document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (h *Handlers) DeleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.services.DeleteDocument(id); err != nil {
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) CreateTemplateHandler(w http.ResponseWriter, r *http.Request) {
	var t models.Template
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// convert to repo template
	varsBytes, _ := json.Marshal(t.Variables)
	repoT := &repository.Template{
		ID:          uuid.New().String(),
		Name:        t.Name,
		Description: t.Description,
		Type:        t.Type,
		Content:     t.Content,
		Variables:   string(varsBytes),
	}

	if err := h.services.CreateTemplate(repoT); err != nil {
		http.Error(w, "Failed to create template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *Handlers) SignDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	var req models.SignDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	sig := &repository.Signature{
		ID:            uuid.New().String(),
		DocumentID:    docID,
		SignerName:    req.SignerName,
		SignerEmail:   req.SignerEmail,
		SignatureData: []byte(req.SignatureData),
		Timestamp:     time.Now().Format(time.RFC3339),
	}

	if err := h.services.CreateSignature(sig); err != nil {
		http.Error(w, "Failed to sign document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": sig.ID})
}

func (h *Handlers) GetDocumentSignaturesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	sigs, err := h.services.GetSignaturesByDocumentID(docID)
	if err != nil {
		http.Error(w, "Failed to get signatures", http.StatusInternalServerError)
		return
	}

	// convert to models
	out := make([]*models.Signature, 0)
	for _, s := range sigs {
		out = append(out, &models.Signature{
			ID:            s.ID,
			DocumentID:    s.DocumentID,
			SignerName:    s.SignerName,
			SignerEmail:   s.SignerEmail,
			SignatureData: string(s.SignatureData),
			Timestamp:     s.Timestamp,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *Handlers) DeleteSignatureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.services.DeleteSignature(id); err != nil {
		http.Error(w, "Failed to delete signature", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) ShareDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]
	log.Printf("ShareDocumentHandler called with docID: %s", docID)

	var req models.ShareDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("Request decoded: shared_with_email=%s, permission=%s", req.SharedWithEmail, req.Permission)

	// find user to share with
	user, err := h.services.GetUserByEmail(req.SharedWithEmail)
	if err != nil {
		log.Printf("User not found for email %s: %v", req.SharedWithEmail, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	log.Printf("User found: %s (%s)", user.Name, user.Email)

	share := &repository.DocumentShare{
		ID:               uuid.New().String(),
		DocumentID:       docID,
		SharedByUserID:   "", // will be set from context if needed
		SharedWithUserID: user.ID,
		Permission:       req.Permission,
	}

	// try to get shared-by from context
	if v := r.Context().Value("user_id"); v != nil {
		if uid, ok := v.(string); ok {
			share.SharedByUserID = uid
			log.Printf("SharedByUserID from context: %s", uid)
		}
	}

	log.Printf("Creating document share: id=%s, docID=%s, sharedByUserID=%s, sharedWithUserID=%s, permission=%s",
		share.ID, share.DocumentID, share.SharedByUserID, share.SharedWithUserID, share.Permission)

	if err := h.services.CreateDocumentShare(share); err != nil {
		log.Printf("Failed to create document share: %v", err)
		http.Error(w, "Failed to share document", http.StatusInternalServerError)
		return
	}

	log.Printf("Document share created successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(share)
}

func (h *Handlers) GetDocumentSharesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	shares, err := h.services.GetDocumentShares(docID)
	if err != nil {
		http.Error(w, "Failed to get shares", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shares)
}

func (h *Handlers) DeleteDocumentShareHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.services.DeleteDocumentShare(id); err != nil {
		http.Error(w, "Failed to delete share", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) GetDocumentCommentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	comments, err := h.services.GetDocumentComments(docID)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *Handlers) CreateDocumentCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	var req models.CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID := ""
	if v := r.Context().Value("user_id"); v != nil {
		if uid, ok := v.(string); ok {
			userID = uid
		}
	}

	comment := &repository.DocumentComment{
		ID:         uuid.New().String(),
		DocumentID: docID,
		UserID:     userID,
		Content:    req.Content,
	}

	if err := h.services.CreateDocumentComment(comment); err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *Handlers) DeleteDocumentCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.services.DeleteDocumentComment(id); err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) GetDocumentVersionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	versions, err := h.services.GetDocumentVersions(docID)
	if err != nil {
		http.Error(w, "Failed to get versions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versions)
}

func generateToken(userID string) (string, error) {
	return middleware.GenerateToken(userID)
}

func convertModelToRepo(doc *models.Document) *repository.Document {
	vars, _ := json.Marshal(doc.Variables)
	return &repository.Document{
		ID:         doc.ID,
		UserID:     doc.UserID,
		TemplateID: doc.TemplateID,
		Title:      doc.Title,
		Type:       doc.Type,
		Content:    doc.Content,
		Variables:  string(vars),
		Status:     doc.Status,
		Version:    doc.Version,
	}
}

func convertRepoTemplateToModel(rt *repository.Template) *models.Template {
	var vars []map[string]interface{}
	err := json.Unmarshal([]byte(rt.Variables), &vars)
	if err != nil {
		log.Printf("Failed to parse template variables: %v", err)
		vars = []map[string]interface{}{}
	}

	createdAt := rt.CreatedAt
	updatedAt := rt.UpdatedAt

	return &models.Template{
		ID:          rt.ID,
		Name:        rt.Name,
		Description: rt.Description,
		Type:        rt.Type,
		Content:     rt.Content,
		Variables:   vars,
		CreatedAt:   toTime(createdAt),
		UpdatedAt:   toTime(updatedAt),
	}
}

func toTime(v interface{}) time.Time {
	switch val := v.(type) {
	case time.Time:
		return val
	default:
		return time.Now()
	}
}

func (h *Handlers) GetSharedWithMeHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by JWT middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("GetSharedWithMeHandler called for userID: %s", userID)

	// Get all document shares where this user is the recipient
	shares, err := h.services.GetDocumentsSharedWithUser(userID)
	if err != nil {
		log.Printf("Failed to get shared documents: %v", err)
		http.Error(w, "Failed to get shared documents", http.StatusInternalServerError)
		return
	}

	// Fetch full document details for each share
	type DocumentWithShare struct {
		Share    *repository.DocumentShare `json:"share"`
		Document *repository.Document      `json:"document"`
		SharedBy *models.User              `json:"shared_by"`
	}

	var result []DocumentWithShare
	for _, share := range shares {
		// Get document details
		doc, err := h.services.GetDocumentByID(share.DocumentID)
		if err != nil {
			log.Printf("Failed to get document %s: %v", share.DocumentID, err)
			continue
		}

		// Get shared by user details
		sharedByUser, err := h.services.GetUserByID(share.SharedByUserID)
		if err != nil {
			log.Printf("Failed to get user %s: %v", share.SharedByUserID, err)
			continue
		}

		result = append(result, DocumentWithShare{
			Share:    share,
			Document: doc,
			SharedBy: &models.User{ID: sharedByUser.ID, Email: sharedByUser.Email, Name: sharedByUser.Name},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func replaceAll(content, old, new string) string {
	for {
		idx := strings.Index(content, old)
		if idx == -1 {
			break
		}
		content = content[:idx] + new + content[idx+len(old):]
	}
	return content
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%v", int(val))
	default:
		return fmt.Sprintf("%v", val)
	}
}
