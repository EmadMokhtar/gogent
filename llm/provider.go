// Package llm provides the LLM provider abstraction and implementations.
package llm

import (
	"context"

	"github.com/EmadMokhtar/gogent/types"
)

// Provider defines the interface for LLM providers.
type Provider interface {
	// Generate produces a response from the LLM for the given messages.
	Generate(ctx context.Context, messages []types.Message, opts ...Option) (*types.Response, error)
	
	// Stream produces a streaming response from the LLM.
	Stream(ctx context.Context, messages []types.Message, opts ...Option) (<-chan Token, error)
}

// Token represents a streaming token from the LLM.
type Token struct {
	Content string      // Token content
	Delta   string      // Delta content for this token
	Done    bool        // Whether streaming is complete
	Error   error       // Any error that occurred
}

// Option represents a configuration option for LLM calls.
type Option func(*Options)

// Options holds configuration for LLM calls.
type Options struct {
	Temperature float64   // Sampling temperature (0.0 to 2.0)
	MaxTokens   int       // Maximum tokens to generate
	TopP        float64   // Nucleus sampling parameter
	Stop        []string  // Stop sequences
	Model       string    // Model name override
	Tools       []ToolDef // Available tools for function calling
}

// ToolDef defines a tool/function that can be called by the LLM.
type ToolDef struct {
	Type     string       `json:"type"`     // Usually "function"
	Function FunctionDef  `json:"function"` // Function definition
}

// FunctionDef defines a function/tool for the LLM.
type FunctionDef struct {
	Name        string                 `json:"name"`        // Function name
	Description string                 `json:"description"` // Function description
	Parameters  map[string]interface{} `json:"parameters"`  // JSON schema for parameters
}

// DefaultOptions returns default LLM options.
func DefaultOptions() *Options {
	return &Options{
		Temperature: 0.7,
		MaxTokens:   1000,
		TopP:        1.0,
		Stop:        []string{},
		Tools:       []ToolDef{},
	}
}

// WithTemperature sets the sampling temperature.
func WithTemperature(temp float64) Option {
	return func(o *Options) {
		o.Temperature = temp
	}
}

// WithMaxTokens sets the maximum tokens to generate.
func WithMaxTokens(tokens int) Option {
	return func(o *Options) {
		o.MaxTokens = tokens
	}
}

// WithTopP sets the nucleus sampling parameter.
func WithTopP(topP float64) Option {
	return func(o *Options) {
		o.TopP = topP
	}
}

// WithStop sets stop sequences.
func WithStop(stop []string) Option {
	return func(o *Options) {
		o.Stop = stop
	}
}

// WithModel sets the model name override.
func WithModel(model string) Option {
	return func(o *Options) {
		o.Model = model
	}
}

// WithTools sets available tools for function calling.
func WithTools(tools []ToolDef) Option {
	return func(o *Options) {
		o.Tools = tools
	}
}
