package memory

import (
	"context"
	"time"
)

// Message represents a message in memory
type Message struct {
	Role      string                 `json:"role"`
	Content   string                 `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Memory is the interface for memory operations
type Memory interface {
	// Add adds a message to memory
	Add(ctx context.Context, msg Message) error
	
	// Get retrieves the most recent messages up to the limit
	Get(ctx context.Context, limit int) ([]Message, error)
	
	// Search searches for messages matching a query
	Search(ctx context.Context, query string) ([]Message, error)
	
	// Clear clears all messages from memory
	Clear(ctx context.Context) error
	
	// Count returns the total number of messages in memory
	Count(ctx context.Context) (int, error)
}

// NewMessage creates a new message with the current timestamp
func NewMessage(role, content string) Message {
	return Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// NewUserMessage creates a new user message
func NewUserMessage(content string) Message {
	return NewMessage("user", content)
}

// NewAssistantMessage creates a new assistant message
func NewAssistantMessage(content string) Message {
	return NewMessage("assistant", content)
}

// NewSystemMessage creates a new system message
func NewSystemMessage(content string) Message {
	return NewMessage("system", content)
}
