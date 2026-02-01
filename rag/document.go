package rag

import (
	"time"
)

// Document represents a document to be stored and retrieved
type Document struct {
	ID        string                 `json:"id"`
	Content   string                 `json:"content"`
	Embedding []float32              `json:"embedding,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// Chunk represents a chunk of a document
type Chunk struct {
	ID         string                 `json:"id"`
	DocumentID string                 `json:"document_id"`
	Content    string                 `json:"content"`
	Embedding  []float32              `json:"embedding,omitempty"`
	Index      int                    `json:"index"`      // Position in the original document
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Score      float32                `json:"score,omitempty"` // Similarity score (used in retrieval)
}

// RetrievalResult represents a retrieved document chunk with its score
type RetrievalResult struct {
	Chunk     Chunk                  `json:"chunk"`
	Score     float32                `json:"score"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NewDocument creates a new document with timestamps
func NewDocument(id, content string) Document {
	now := time.Now()
	return Document{
		ID:        id,
		Content:   content,
		Metadata:  make(map[string]interface{}),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewChunk creates a new chunk
func NewChunk(id, documentID, content string, index int) Chunk {
	return Chunk{
		ID:         id,
		DocumentID: documentID,
		Content:    content,
		Index:      index,
		Metadata:   make(map[string]interface{}),
	}
}
