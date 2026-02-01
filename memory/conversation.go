package memory

import (
	"context"
	"strings"
	"sync"
)

// Conversation implements conversation memory storage
type Conversation struct {
	mu       sync.RWMutex
	messages []Message
	maxSize  int
}

// ConversationOption is a functional option for configuring conversation memory
type ConversationOption func(*Conversation)

// NewConversation creates a new conversation memory
func NewConversation(opts ...ConversationOption) *Conversation {
	c := &Conversation{
		messages: make([]Message, 0),
		maxSize:  100, // Default max size
	}
	
	for _, opt := range opts {
		opt(c)
	}
	
	return c
}

// WithMaxSize sets the maximum number of messages to store
func WithMaxSize(size int) ConversationOption {
	return func(c *Conversation) {
		c.maxSize = size
	}
}

// Add adds a message to the conversation
func (c *Conversation) Add(ctx context.Context, msg Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.messages = append(c.messages, msg)
	
	// Trim if exceeds max size
	if len(c.messages) > c.maxSize {
		c.messages = c.messages[len(c.messages)-c.maxSize:]
	}
	
	return nil
}

// Get retrieves the most recent messages up to the limit
func (c *Conversation) Get(ctx context.Context, limit int) ([]Message, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if limit <= 0 || limit > len(c.messages) {
		limit = len(c.messages)
	}
	
	// Return the most recent messages
	start := len(c.messages) - limit
	if start < 0 {
		start = 0
	}
	
	result := make([]Message, limit)
	copy(result, c.messages[start:])
	
	return result, nil
}

// Search searches for messages containing the query string
func (c *Conversation) Search(ctx context.Context, query string) ([]Message, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	query = strings.ToLower(query)
	var results []Message
	
	for _, msg := range c.messages {
		if strings.Contains(strings.ToLower(msg.Content), query) {
			results = append(results, msg)
		}
	}
	
	return results, nil
}

// Clear clears all messages from the conversation
func (c *Conversation) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.messages = make([]Message, 0)
	return nil
}

// Count returns the total number of messages
func (c *Conversation) Count(ctx context.Context) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return len(c.messages), nil
}

// GetAll returns all messages (for internal use)
func (c *Conversation) GetAll(ctx context.Context) ([]Message, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	result := make([]Message, len(c.messages))
	copy(result, c.messages)
	
	return result, nil
}
