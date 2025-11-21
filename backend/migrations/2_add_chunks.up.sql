-- 2_add_chunks.up.sql: add chunks table and pgvector support

-- Enable pgvector extension (if not already enabled)
CREATE EXTENSION IF NOT EXISTS vector;

-- Table to store document chunks and their embeddings
CREATE TABLE IF NOT EXISTS chunks (
    id VARCHAR(36) PRIMARY KEY,
    document_id VARCHAR(36),
    content TEXT NOT NULL,
    metadata JSONB,
    embedding VECTOR(1536),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_chunks_document_id ON chunks(document_id);

-- ivfflat index for fast approximate nearest neighbors on embeddings
-- Adjust 'lists' parameter based on dataset size; 100 is a reasonable default for small corpora
CREATE INDEX IF NOT EXISTS idx_chunks_embedding ON chunks USING ivfflat (embedding) WITH (lists = 100);
