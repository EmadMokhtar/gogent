package llm

import (
	"context"
)

// Provider is the interface that LLM providers must implement
type Provider interface {
	// Generate creates a completion for the given messages
	Generate(ctx context.Context, messages []Message, opts ...Option) (*Response, error)
	
	// Stream creates a streaming completion for the given messages
	Stream(ctx context.Context, messages []Message, opts ...Option) (<-chan Token, error)
	
	// Name returns the provider name
	Name() string
	
	// Model returns the model name
	Model() string
}

// Embedder is the interface for generating embeddings
type Embedder interface {
	// Embed generates an embedding for the given text
	Embed(ctx context.Context, text string) ([]float32, error)
	
	// EmbedBatch generates embeddings for multiple texts
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
	
	// Dimensions returns the dimensionality of the embeddings
	Dimensions() int
}
