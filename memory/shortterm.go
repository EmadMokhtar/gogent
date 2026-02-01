package memory

import (
	"context"
	"strings"
	"sync"
	"time"
)

// ShortTerm implements short-term memory with time-based decay
type ShortTerm struct {
	mu       sync.RWMutex
	messages []MessageWithExpiry
	ttl      time.Duration
}

// MessageWithExpiry extends Message with expiration time
type MessageWithExpiry struct {
	Message
	ExpiresAt time.Time
}

// ShortTermOption is a functional option for configuring short-term memory
type ShortTermOption func(*ShortTerm)

// NewShortTerm creates a new short-term memory
func NewShortTerm(opts ...ShortTermOption) *ShortTerm {
	s := &ShortTerm{
		messages: make([]MessageWithExpiry, 0),
		ttl:      30 * time.Minute, // Default TTL
	}
	
	for _, opt := range opts {
		opt(s)
	}
	
	// Start cleanup goroutine
	go s.cleanup()
	
	return s
}

// WithTTL sets the time-to-live for messages
func WithTTL(ttl time.Duration) ShortTermOption {
	return func(s *ShortTerm) {
		s.ttl = ttl
	}
}

// Add adds a message to short-term memory
func (s *ShortTerm) Add(ctx context.Context, msg Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	msgWithExpiry := MessageWithExpiry{
		Message:   msg,
		ExpiresAt: time.Now().Add(s.ttl),
	}
	
	s.messages = append(s.messages, msgWithExpiry)
	return nil
}

// Get retrieves the most recent non-expired messages up to the limit
func (s *ShortTerm) Get(ctx context.Context, limit int) ([]Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	now := time.Now()
	var validMessages []Message
	
	// Filter out expired messages
	for _, msg := range s.messages {
		if msg.ExpiresAt.After(now) {
			validMessages = append(validMessages, msg.Message)
		}
	}
	
	if limit <= 0 || limit > len(validMessages) {
		limit = len(validMessages)
	}
	
	// Return the most recent messages
	start := len(validMessages) - limit
	if start < 0 {
		start = 0
	}
	
	result := make([]Message, limit)
	copy(result, validMessages[start:])
	
	return result, nil
}

// Search searches for non-expired messages containing the query string
func (s *ShortTerm) Search(ctx context.Context, query string) ([]Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	now := time.Now()
	query = strings.ToLower(query)
	var results []Message
	
	for _, msg := range s.messages {
		if msg.ExpiresAt.After(now) && strings.Contains(strings.ToLower(msg.Content), query) {
			results = append(results, msg.Message)
		}
	}
	
	return results, nil
}

// Clear clears all messages from short-term memory
func (s *ShortTerm) Clear(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.messages = make([]MessageWithExpiry, 0)
	return nil
}

// Count returns the number of non-expired messages
func (s *ShortTerm) Count(ctx context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	now := time.Now()
	count := 0
	
	for _, msg := range s.messages {
		if msg.ExpiresAt.After(now) {
			count++
		}
	}
	
	return count, nil
}

// cleanup periodically removes expired messages
func (s *ShortTerm) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		validMessages := make([]MessageWithExpiry, 0)
		
		for _, msg := range s.messages {
			if msg.ExpiresAt.After(now) {
				validMessages = append(validMessages, msg)
			}
		}
		
		s.messages = validMessages
		s.mu.Unlock()
	}
}
