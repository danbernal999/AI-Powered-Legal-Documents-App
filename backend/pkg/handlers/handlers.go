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
