package llm

// Options represents configuration options for LLM calls
type Options struct {
	Temperature      float32
	MaxTokens        int
	TopP             float32
	FrequencyPenalty float32
	PresencePenalty  float32
	Stop             []string
	Stream           bool
	Tools            []ToolDefinition
}

// ToolDefinition represents a tool that can be called by the LLM
type ToolDefinition struct {
	Type     string           `json:"type"` // "function"
	Function FunctionDef      `json:"function"`
}

// FunctionDef represents a function definition for tool calling
type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Option is a functional option for configuring LLM calls
type Option func(*Options)

// DefaultOptions returns default options
func DefaultOptions() Options {
	return Options{
		Temperature: 0.7,
		MaxTokens:   1000,
		TopP:        1.0,
	}
}

// WithTemperature sets the temperature
func WithTemperature(t float32) Option {
	return func(o *Options) {
		o.Temperature = t
	}
}

// WithMaxTokens sets the maximum number of tokens
func WithMaxTokens(max int) Option {
	return func(o *Options) {
		o.MaxTokens = max
	}
}

// WithTopP sets the top-p sampling parameter
func WithTopP(topP float32) Option {
	return func(o *Options) {
		o.TopP = topP
	}
}

// WithFrequencyPenalty sets the frequency penalty
func WithFrequencyPenalty(penalty float32) Option {
	return func(o *Options) {
		o.FrequencyPenalty = penalty
	}
}

// WithPresencePenalty sets the presence penalty
func WithPresencePenalty(penalty float32) Option {
	return func(o *Options) {
		o.PresencePenalty = penalty
	}
}

// WithStop sets stop sequences
func WithStop(stop []string) Option {
	return func(o *Options) {
		o.Stop = stop
	}
}

// WithStream enables streaming
func WithStream() Option {
	return func(o *Options) {
		o.Stream = true
	}
}

// WithTools sets available tools for function calling
func WithTools(tools []ToolDefinition) Option {
	return func(o *Options) {
		o.Tools = tools
	}
}
