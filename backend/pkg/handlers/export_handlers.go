package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/danbernal999/kiradoc-backend/pkg/models"
	"github.com/danbernal999/kiradoc-backend/pkg/repository"
	"github.com/danbernal999/kiradoc-backend/pkg/services"
	"github.com/gorilla/mux"
)

func (h *Handlers) ExportDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	documentID := vars["id"]

	userID := r.Context().Value("userID").(string)

	doc, err := h.services.GetDocumentByID(documentID)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	if doc.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var req models.ExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Format == "" {
		req.Format = "pdf"
	}

	exportService := services.NewExportService(h.services.GetRepository())
	opts := exportService.BuildExportOptions(&req)

	var result *services.ExportResult
	var exportErr error

	switch req.Format {
	case "pdf":
		result, exportErr = exportService.ExportPDF(doc, opts)
	case "docx":
		result, exportErr = exportService.ExportDOCX(doc, opts)
	case "html":
		result, exportErr = exportService.ExportHTML(doc, opts)
	default:
		http.Error(w, "Unsupported format: "+req.Format, http.StatusBadRequest)
		return
	}

	if exportErr != nil {
		log.Printf("Export error: %v", exportErr)
		http.Error(w, fmt.Sprintf("Failed to export document: %v", exportErr), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", result.FileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(result.Data)))

	if _, err := w.Write(result.Data); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (h *Handlers) BatchExportHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req models.BatchExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.DocumentIDs) == 0 {
		http.Error(w, "No documents specified", http.StatusBadRequest)
		return
	}

	if req.Format == "" {
		req.Format = "pdf"
	}

	var documents []*repository.Document
	for _, docID := range req.DocumentIDs {
		doc, err := h.services.GetDocumentByID(docID)
		if err != nil {
			continue
		}

		if doc.UserID != userID {
			continue
		}

		documents = append(documents, doc)
	}

	if len(documents) == 0 {
		http.Error(w, "No valid documents found", http.StatusNotFound)
		return
	}

	exportService := services.NewExportService(h.services.GetRepository())
	opts := exportService.BuildExportOptions(&req.Options)

	result, err := exportService.ExportBatch(documents, req.Format, opts)
	if err != nil {
		log.Printf("Batch export error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to export documents: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", result.FileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(result.Data)))

	if _, err := w.Write(result.Data); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (h *Handlers) ExportTemplateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateID := vars["id"]

	_, err := h.services.GetTemplateByID(templateID)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	var req models.ExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Format == "" {
		req.Format = "pdf"
	}

	http.Error(w, "Template export not yet implemented", http.StatusNotImplemented)
}
