package rag

import "context"

// Embedder defines the interface for generating text embeddings
type Embedder interface {
	// Embed generates an embedding for a single text
	Embed(ctx context.Context, text string) ([]float32, error)

	// EmbedBatch generates embeddings for multiple texts
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)

	// Dimension returns the dimensionality of the embeddings
	Dimension() int
}
