package memory

import (
	"context"
	"strings"
	"sync"
)

// LongTerm implements long-term memory with importance-based retention
type LongTerm struct {
	mu       sync.RWMutex
	messages []MessageWithImportance
	maxSize  int
}

// MessageWithImportance extends Message with importance score
type MessageWithImportance struct {
	Message
	Importance float64 // 0.0 to 1.0
}

// LongTermOption is a functional option for configuring long-term memory
type LongTermOption func(*LongTerm)

// NewLongTerm creates a new long-term memory
func NewLongTerm(opts ...LongTermOption) *LongTerm {
	l := &LongTerm{
		messages: make([]MessageWithImportance, 0),
		maxSize:  1000, // Default max size
	}
	
	for _, opt := range opts {
		opt(l)
	}
	
	return l
}

// WithLongTermMaxSize sets the maximum number of messages to store
func WithLongTermMaxSize(size int) LongTermOption {
	return func(l *LongTerm) {
		l.maxSize = size
	}
}

// Add adds a message to long-term memory with default importance
func (l *LongTerm) Add(ctx context.Context, msg Message) error {
	return l.AddWithImportance(ctx, msg, 0.5)
}

// AddWithImportance adds a message with a specific importance score
func (l *LongTerm) AddWithImportance(ctx context.Context, msg Message, importance float64) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	msgWithImportance := MessageWithImportance{
		Message:    msg,
		Importance: importance,
	}
	
	l.messages = append(l.messages, msgWithImportance)
	
	// If exceeds max size, remove least important messages
	if len(l.messages) > l.maxSize {
		l.trimByImportance()
	}
	
	return nil
}

// Get retrieves the most recent messages up to the limit
func (l *LongTerm) Get(ctx context.Context, limit int) ([]Message, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	if limit <= 0 || limit > len(l.messages) {
		limit = len(l.messages)
	}
	
	// Return the most recent messages
	start := len(l.messages) - limit
	if start < 0 {
		start = 0
	}
	
	result := make([]Message, limit)
	for i := 0; i < limit; i++ {
		result[i] = l.messages[start+i].Message
	}
	
	return result, nil
}

// Search searches for messages containing the query string
func (l *LongTerm) Search(ctx context.Context, query string) ([]Message, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	query = strings.ToLower(query)
	var results []Message
	
	for _, msg := range l.messages {
		if strings.Contains(strings.ToLower(msg.Content), query) {
			results = append(results, msg.Message)
		}
	}
	
	return results, nil
}

// Clear clears all messages from long-term memory
func (l *LongTerm) Clear(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	l.messages = make([]MessageWithImportance, 0)
	return nil
}

// Count returns the total number of messages
func (l *LongTerm) Count(ctx context.Context) (int, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	return len(l.messages), nil
}

// trimByImportance removes the least important messages
func (l *LongTerm) trimByImportance() {
	// Sort by importance (descending)
	// Simple bubble sort for demonstration
	for i := 0; i < len(l.messages)-1; i++ {
		for j := 0; j < len(l.messages)-i-1; j++ {
			if l.messages[j].Importance < l.messages[j+1].Importance {
				l.messages[j], l.messages[j+1] = l.messages[j+1], l.messages[j]
			}
		}
	}
	
	// Keep only the most important messages
	l.messages = l.messages[:l.maxSize]
}

// GetByImportance retrieves messages with importance above threshold
func (l *LongTerm) GetByImportance(ctx context.Context, threshold float64, limit int) ([]Message, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	var results []Message
	
	for _, msg := range l.messages {
		if msg.Importance >= threshold {
			results = append(results, msg.Message)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}
	
	return results, nil
}
