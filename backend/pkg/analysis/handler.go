package analysis

import (
    "encoding/json"
    "net/http"
)

// UploadDocumentHandler is a minimal HTTP handler that will accept a document
// upload and return a placeholder response while the ingestion pipeline is implemented.
func UploadDocumentHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotImplemented)
    json.NewEncoder(w).Encode(map[string]string{"error": "upload/ingestion not implemented yet"})
}
