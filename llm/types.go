package llm

// Message represents a message in a conversation
type Message struct {
	Role      string                 `json:"role"`      // "user", "assistant", "system"
	Content   string                 `json:"content"`
	ToolCalls []ToolCall             `json:"tool_calls,omitempty"` // For tool execution
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Response represents an LLM response
type Response struct {
	Content   string                 `json:"content"`
	ToolCalls []ToolCall             `json:"tool_calls,omitempty"`
	Usage     Usage                  `json:"usage"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ToolCall represents a tool call request from the LLM
type ToolCall struct {
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Token represents a streaming token
type Token struct {
	Content   string
	ToolCalls []ToolCall
	Done      bool
	Error     error
}
