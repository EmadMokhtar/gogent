package openai

import (
	"context"
	"fmt"

	"github.com/EmadMokhtar/gogent/llm"
)

// Provider implements the LLM Provider interface for OpenAI
type Provider struct {
	model  string
	apiKey string
}

// ProviderOption is a functional option for configuring the OpenAI provider
type ProviderOption func(*Provider)

// New creates a new OpenAI provider
func New(model string, opts ...ProviderOption) *Provider {
	p := &Provider{
		model: model,
	}
	
	for _, opt := range opts {
		opt(p)
	}
	
	return p
}

// WithAPIKey sets the API key
func WithAPIKey(key string) ProviderOption {
	return func(p *Provider) {
		p.apiKey = key
	}
}

// Generate creates a completion for the given messages
func (p *Provider) Generate(ctx context.Context, messages []llm.Message, opts ...llm.Option) (*llm.Response, error) {
	// This is a mock implementation for demonstration purposes
	// In a real implementation, this would call the OpenAI API
	
	if p.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key not set")
	}
	
	// Apply options
	options := llm.DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	
	// Mock response
	return &llm.Response{
		Content: "This is a mock response from OpenAI. In a real implementation, this would call the actual OpenAI API.",
		Role:    llm.RoleAssistant,
		Usage: &llm.Usage{
			PromptTokens:     100,
			CompletionTokens: 50,
			TotalTokens:      150,
		},
	}, nil
}

// Stream creates a streaming completion for the given messages
func (p *Provider) Stream(ctx context.Context, messages []llm.Message, opts ...llm.Option) (<-chan llm.Token, error) {
	// This is a mock implementation for demonstration purposes
	// In a real implementation, this would call the OpenAI streaming API
	
	if p.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key not set")
	}
	
	ch := make(chan llm.Token, 10)
	
	go func() {
		defer close(ch)
		
		// Mock streaming response
		mockResponse := "This is a mock streaming response from OpenAI."
		for _, char := range mockResponse {
			select {
			case <-ctx.Done():
				ch <- llm.Token{Error: ctx.Err()}
				return
			case ch <- llm.Token{Content: string(char), Done: false}:
			}
		}
		
		ch <- llm.Token{Done: true}
	}()
	
	return ch, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "openai"
}

// Model returns the model name
func (p *Provider) Model() string {
	return p.model
}

// Embedder implements the Embedder interface for OpenAI
type Embedder struct {
	model      string
	apiKey     string
	dimensions int
}

// EmbedderOption is a functional option for configuring the OpenAI embedder
type EmbedderOption func(*Embedder)

// NewEmbedder creates a new OpenAI embedder
func NewEmbedder(model string, opts ...EmbedderOption) *Embedder {
	e := &Embedder{
		model:      model,
		dimensions: 1536, // Default for text-embedding-ada-002
	}
	
	for _, opt := range opts {
		opt(e)
	}
	
	return e
}

// WithEmbedderAPIKey sets the API key for the embedder
func WithEmbedderAPIKey(key string) EmbedderOption {
	return func(e *Embedder) {
		e.apiKey = key
	}
}

// WithDimensions sets the embedding dimensions
func WithDimensions(dim int) EmbedderOption {
	return func(e *Embedder) {
		e.dimensions = dim
	}
}

// Embed generates an embedding for the given text
func (e *Embedder) Embed(ctx context.Context, text string) ([]float32, error) {
	// This is a mock implementation for demonstration purposes
	// In a real implementation, this would call the OpenAI embeddings API
	
	if e.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key not set")
	}
	
	// Return a mock embedding vector
	embedding := make([]float32, e.dimensions)
	for i := range embedding {
		embedding[i] = 0.01 // Mock values
	}
	
	return embedding, nil
}

// EmbedBatch generates embeddings for multiple texts
func (e *Embedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	// This is a mock implementation for demonstration purposes
	// In a real implementation, this would call the OpenAI embeddings API
	
	if e.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key not set")
	}
	
	embeddings := make([][]float32, len(texts))
	for i := range texts {
		embedding, err := e.Embed(ctx, texts[i])
		if err != nil {
			return nil, err
		}
		embeddings[i] = embedding
	}
	
	return embeddings, nil
}

// Dimensions returns the dimensionality of the embeddings
func (e *Embedder) Dimensions() int {
	return e.dimensions
}
