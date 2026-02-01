package rag

import (
	"context"
)

// Embedder is the interface for generating embeddings
type Embedder interface {
	// Embed generates an embedding for the given text
	Embed(ctx context.Context, text string) ([]float32, error)
	
	// EmbedBatch generates embeddings for multiple texts
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
	
	// Dimensions returns the dimensionality of the embeddings
	Dimensions() int
}

// MockEmbedder is a mock embedder for testing
type MockEmbedder struct {
	dimensions int
}

// NewMockEmbedder creates a new mock embedder
func NewMockEmbedder(dimensions int) *MockEmbedder {
	return &MockEmbedder{
		dimensions: dimensions,
	}
}

// Embed generates a mock embedding
func (m *MockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	// Generate a simple mock embedding based on text length
	embedding := make([]float32, m.dimensions)
	for i := range embedding {
		embedding[i] = float32(len(text)%100) / 100.0
	}
	return embedding, nil
}

// EmbedBatch generates mock embeddings for multiple texts
func (m *MockEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))
	for i, text := range texts {
		embedding, err := m.Embed(ctx, text)
		if err != nil {
			return nil, err
		}
		embeddings[i] = embedding
	}
	return embeddings, nil
}

// Dimensions returns the dimensionality of the embeddings
func (m *MockEmbedder) Dimensions() int {
	return m.dimensions
}
