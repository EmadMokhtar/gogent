// Package types provides common types used across the Gogent framework.
package types

// Message represents a conversation message in the agent system.
type Message struct {
	Role       string                 `json:"role"`        // "user", "assistant", "system", "tool"
	Content    string                 `json:"content"`     // Message content
	ToolCalls  []ToolCall             `json:"tool_calls,omitempty"`  // For assistant messages requesting tool calls
	ToolCallID string                 `json:"tool_call_id,omitempty"` // For tool response messages
	Name       string                 `json:"name,omitempty"`         // Optional name (for tool messages)
	Metadata   map[string]interface{} `json:"metadata,omitempty"`     // Additional metadata
}

// ToolCall represents a request to execute a tool.
type ToolCall struct {
	ID       string       `json:"id"`       // Unique identifier for this tool call
	Type     string       `json:"type"`     // Type of tool call (usually "function")
	Function FunctionCall `json:"function"` // Function call details
}

// FunctionCall represents the details of a function/tool to call.
type FunctionCall struct {
	Name      string `json:"name"`      // Name of the function/tool
	Arguments string `json:"arguments"` // JSON-encoded arguments
}

// Response represents the agent's response to a query.
type Response struct {
	Content   string                 `json:"content"`   // The response content
	Messages  []Message              `json:"messages"`  // Complete message history
	ToolCalls []ToolCall             `json:"tool_calls,omitempty"` // Any tool calls made
	Metadata  map[string]interface{} `json:"metadata,omitempty"`   // Additional metadata
}
