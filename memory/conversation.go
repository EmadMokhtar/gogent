package memory

import (
	"context"
	"strings"
	"sync"
	"time"
)

// ConversationMemory stores complete conversation history.
// It provides thread-safe access to messages and supports optional size limits.
type ConversationMemory struct {
	mu       sync.RWMutex
	messages []Message
	maxSize  int
}

// NewConversation creates a new conversation memory with the given options.
func NewConversation(opts ...Option) *ConversationMemory {
	options := &ConversationOptions{
		MaxSize: 0, // unlimited by default
	}
	for _, opt := range opts {
		opt(options)
	}

	return &ConversationMemory{
		messages: make([]Message, 0),
		maxSize:  options.MaxSize,
	}
}

// Add stores a message in the conversation memory.
func (cm *ConversationMemory) Add(ctx context.Context, msg Message) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Set timestamp if not provided
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	// Add the message
	cm.messages = append(cm.messages, msg)

	// Apply size limit if configured
	if cm.maxSize > 0 && len(cm.messages) > cm.maxSize {
		// Remove oldest messages to stay within limit
		cm.messages = cm.messages[len(cm.messages)-cm.maxSize:]
	}

	return nil
}

// Get retrieves recent messages from memory.
// If limit is 0, all messages are returned.
func (cm *ConversationMemory) Get(ctx context.Context, limit int) ([]Message, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if limit <= 0 || limit > len(cm.messages) {
		// Return all messages
		result := make([]Message, len(cm.messages))
		copy(result, cm.messages)
		return result, nil
	}

	// Return the last 'limit' messages
	start := len(cm.messages) - limit
	result := make([]Message, limit)
	copy(result, cm.messages[start:])
	return result, nil
}

// Search finds messages containing the query string in their content.
func (cm *ConversationMemory) Search(ctx context.Context, query string, limit int) ([]Message, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var results []Message
	query = strings.ToLower(query)

	for _, msg := range cm.messages {
		if strings.Contains(strings.ToLower(msg.Content), query) {
			results = append(results, msg)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}

	return results, nil
}

// Clear removes all messages from memory.
func (cm *ConversationMemory) Clear(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.messages = make([]Message, 0)
	return nil
}

// Stats returns statistics about the conversation memory.
func (cm *ConversationMemory) Stats(ctx context.Context) (*Stats, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	stats := &Stats{
		TotalMessages: len(cm.messages),
	}

	for i, msg := range cm.messages {
		// Count by role
		switch msg.Role {
		case "user":
			stats.UserMessages++
		case "assistant":
			stats.AssistantMessages++
		case "system":
			stats.SystemMessages++
		}

		// Track oldest and newest
		if i == 0 {
			oldest := msg.Timestamp
			stats.OldestMessage = &oldest
		}
		if i == len(cm.messages)-1 {
			newest := msg.Timestamp
			stats.NewestMessage = &newest
		}

		// Approximate size
		stats.TotalSize += int64(len(msg.Content))
	}

	return stats, nil
}
