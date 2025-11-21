-- 2_add_chunks.down.sql: revert chunks table and indexes

DROP INDEX IF EXISTS idx_chunks_embedding;
DROP INDEX IF EXISTS idx_chunks_document_id;
DROP TABLE IF EXISTS chunks;

-- Note: we do not DROP EXTENSION vector globally here to avoid affecting other migrations.
