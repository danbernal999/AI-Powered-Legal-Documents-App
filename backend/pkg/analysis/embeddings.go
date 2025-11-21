package analysis

import (
    "context"
    "errors"
)

// GenerateEmbedding calls the configured embeddings provider and returns a vector.
// This is a placeholder; integrate your ai_client or OpenAI SDK here.
func GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
    // TODO: call OpenAI embeddings or internal provider
    return nil, errors.New("GenerateEmbedding: not implemented")
}
