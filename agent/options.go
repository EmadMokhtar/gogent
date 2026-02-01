package agent

import (
	"github.com/EmadMokhtar/gogent/llm"
	"github.com/EmadMokhtar/gogent/memory"
	"github.com/EmadMokhtar/gogent/rag"
	"github.com/EmadMokhtar/gogent/tools"
)

// Options contains configuration options for an agent
type Options struct {
	LLM          llm.Provider
	Memory       memory.Memory
	RAG          *rag.RAG
	Tools        []tools.Tool
	SystemPrompt string
	MaxTokens    int
	Temperature  float32
	Streaming    bool
	MaxRetries   int
}

// Option is a functional option for configuring an agent
type Option func(*Options)

// DefaultOptions returns default options
func DefaultOptions() Options {
	return Options{
		SystemPrompt: "You are a helpful AI assistant.",
		MaxTokens:    1000,
		Temperature:  0.7,
		Streaming:    false,
		MaxRetries:   3,
	}
}

// WithLLM sets the LLM provider
func WithLLM(llmProvider llm.Provider) Option {
	return func(o *Options) {
		o.LLM = llmProvider
	}
}

// WithMemory sets the memory system
func WithMemory(mem memory.Memory) Option {
	return func(o *Options) {
		o.Memory = mem
	}
}

// WithRAG sets the RAG system
func WithRAG(ragSystem *rag.RAG) Option {
	return func(o *Options) {
		o.RAG = ragSystem
	}
}

// WithTools sets the tools
func WithTools(toolList ...tools.Tool) Option {
	return func(o *Options) {
		o.Tools = toolList
	}
}

// WithSystemPrompt sets the system prompt
func WithSystemPrompt(prompt string) Option {
	return func(o *Options) {
		o.SystemPrompt = prompt
	}
}

// WithMaxTokens sets the maximum tokens
func WithMaxTokens(max int) Option {
	return func(o *Options) {
		o.MaxTokens = max
	}
}

// WithTemperature sets the temperature
func WithTemperature(temp float32) Option {
	return func(o *Options) {
		o.Temperature = temp
	}
}

// WithStreaming enables streaming
func WithStreaming() Option {
	return func(o *Options) {
		o.Streaming = true
	}
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(max int) Option {
	return func(o *Options) {
		o.MaxRetries = max
	}
}
