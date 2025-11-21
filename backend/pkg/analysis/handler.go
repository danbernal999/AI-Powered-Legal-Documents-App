package analysis

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/danbernal999/kiradoc-backend/pkg/repository"
	"github.com/google/uuid"
)

type UploadDocumentRequest struct {
	DocumentID string `json:"document_id"`
}

type UploadDocumentResponse struct {
	DocumentID string `json:"document_id"`
	ChunksCount int    `json:"chunks_count"`
	Message    string `json:"message"`
}

func NewUploadDocumentHandler(repo *repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to parse form"})
			return
		}

		docID := r.FormValue("document_id")
		if docID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "document_id is required"})
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "file is required"})
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to read file"})
			return
		}

		ctx := r.Context()
		text, err := ExtractTextFromDocument(ctx, fileBytes, "text/plain")
		if err != nil {
			text = string(fileBytes)
		}

		chunks := ChunkText(text, 1000, 100)
		if len(chunks) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "no text extracted from document"})
			return
		}

		savedCount := 0
		for i, chunkText := range chunks {
			embedding, err := GenerateEmbedding(ctx, chunkText)
			if err != nil {
				continue
			}

			chunk := &repository.Chunk{
				ID:         uuid.New().String(),
				DocumentID: docID,
				Content:    chunkText,
				Metadata: map[string]interface{}{
					"chunk_index": i,
					"total_chunks": len(chunks),
				},
			}

			err = repo.SaveChunk(chunk, embedding)
			if err != nil {
				continue
			}
			savedCount++
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UploadDocumentResponse{
			DocumentID: docID,
			ChunksCount: savedCount,
			Message:    fmt.Sprintf("Successfully uploaded and processed %d chunks", savedCount),
		})
	}
}
