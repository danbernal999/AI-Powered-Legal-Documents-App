package analysis

import "time"

// Chunk represents a document fragment stored in the vector DB
type Chunk struct {
    ID         string                 `json:"id"`
    DocumentID string                 `json:"document_id"`
    Content    string                 `json:"content"`
    Metadata   map[string]interface{} `json:"metadata"`
    // Embedding is intentionally omitted from JSON responses by default
    CreatedAt  time.Time              `json:"created_at"`
}

// AnalysisResult is a structured result returned by the RAG pipeline
type AnalysisResult struct {
    Summary       string                 `json:"summary"`
    Findings      []map[string]interface{} `json:"findings"`
    Recommendations []map[string]interface{} `json:"recommendations"`
    Confidence    string                 `json:"confidence"`
    RawOutput     string                 `json:"raw_output"`
}
