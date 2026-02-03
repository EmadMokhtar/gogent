package llm

import "context"

// Provider is the interface that all LLM providers must implement
type Provider interface {
	// Generate sends messages to the LLM and returns a response
	Generate(ctx context.Context, messages []Message, opts ...Option) (*Response, error)
	
	// Stream sends messages to the LLM and returns a channel of tokens
	Stream(ctx context.Context, messages []Message, opts ...Option) (<-chan Token, error)
}
