package openai

import (
	"context"

	"github.com/EmadMokhtar/gogent/memory"
)

// Client represents an OpenAI API client.
type Client struct {
	apiKey string
	model  string
}

// Option configures the OpenAI client.
type Option func(*Client)

// WithAPIKey sets the API key.
func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

// WithModel sets the model name.
func WithModel(model string) Option {
	return func(c *Client) {
		c.model = model
	}
}

// New creates a new OpenAI client.
func New(opts ...Option) *Client {
	c := &Client{
		model: "gpt-3.5-turbo",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Generate produces a response using OpenAI API.
// This is a minimal stub implementation for the memory system demo.
func (c *Client) Generate(ctx context.Context, messages []memory.Message) (string, error) {
	// In a real implementation, this would call the OpenAI API
	// For now, return a simple response
	return "This is a simulated response from OpenAI API", nil
}
