package llm

import (
	"time"
)

// Role represents the role of a message sender
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message represents a single message in a conversation
type Message struct {
	Role      Role                   `json:"role"`
	Content   string                 `json:"content"`
	Name      string                 `json:"name,omitempty"`       // Optional name for tool messages
	ToolCalls []ToolCall             `json:"tool_calls,omitempty"` // Tool calls made by assistant
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ToolCall represents a tool call made by the LLM
type ToolCall struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"` // "function"
	Function FunctionCall           `json:"function"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// FunctionCall represents a function call within a tool call
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string of arguments
}

// Token represents a streaming token response
type Token struct {
	Content string
	Done    bool
	Error   error
}

// Response represents the complete response from an LLM
type Response struct {
	Content   string                 `json:"content"`
	Role      Role                   `json:"role"`
	ToolCalls []ToolCall             `json:"tool_calls,omitempty"`
	Usage     *Usage                 `json:"usage,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewMessage creates a new message with the current timestamp
func NewMessage(role Role, content string) Message {
	return Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// NewSystemMessage creates a new system message
func NewSystemMessage(content string) Message {
	return NewMessage(RoleSystem, content)
}

// NewUserMessage creates a new user message
func NewUserMessage(content string) Message {
	return NewMessage(RoleUser, content)
}

// NewAssistantMessage creates a new assistant message
func NewAssistantMessage(content string) Message {
	return NewMessage(RoleAssistant, content)
}
