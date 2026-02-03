package llm

// Options holds configuration for LLM calls
type Options struct {
	Temperature *float64
	MaxTokens   *int
	TopP        *float64
	Stop        []string
	Tools       []ToolDefinition
}

// ToolDefinition defines a tool that can be used by the LLM
type ToolDefinition struct {
	Type     string                 `json:"type"`
	Function FunctionDefinition     `json:"function"`
}

// FunctionDefinition defines a function tool
type FunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Option is a functional option for configuring LLM calls
type Option func(*Options)

// WithTemperature sets the temperature for the LLM call
func WithTemperature(temperature float64) Option {
	return func(o *Options) {
		o.Temperature = &temperature
	}
}

// WithMaxTokens sets the maximum number of tokens for the LLM call
func WithMaxTokens(maxTokens int) Option {
	return func(o *Options) {
		o.MaxTokens = &maxTokens
	}
}

// WithTopP sets the top-p sampling parameter
func WithTopP(topP float64) Option {
	return func(o *Options) {
		o.TopP = &topP
	}
}

// WithStop sets stop sequences
func WithStop(stop []string) Option {
	return func(o *Options) {
		o.Stop = stop
	}
}

// WithTools sets the tools available to the LLM
func WithTools(tools []ToolDefinition) Option {
	return func(o *Options) {
		o.Tools = tools
	}
}

// ApplyOptions applies the given options to an Options struct
func ApplyOptions(opts ...Option) *Options {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
