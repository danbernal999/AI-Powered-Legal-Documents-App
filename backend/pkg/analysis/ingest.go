package analysis

import (
    "context"
    "errors"
)

// ExtractTextFromDocument should extract text from uploaded documents (PDF/DOCX/TXT).
// This is a placeholder for the real implementation (use pdfcpu, unidoc, or external OCR).
func ExtractTextFromDocument(ctx context.Context, fileBytes []byte, contentType string) (string, error) {
    // TODO: implement PDF/DOCX parsing and OCR fallback
    return "", errors.New("ExtractTextFromDocument: not implemented")
}

// ChunkText splits a large text into chunks suitable for embeddings and retrieval.
func ChunkText(text string, maxTokens int, overlap int) []string {
    // Very simple chunking by fixed size characters as placeholder.
    // Replace with token-based chunking (using tiktoken or similar) in production.
    if len(text) == 0 {
        return []string{}
    }
    var chunks []string
    step := maxTokens - overlap
    if step <= 0 {
        step = maxTokens
    }
    for i := 0; i < len(text); i += step {
        end := i + maxTokens
        if end > len(text) {
            end = len(text)
        }
        chunks = append(chunks, text[i:end])
        if end == len(text) {
            break
        }
    }
    return chunks
}
