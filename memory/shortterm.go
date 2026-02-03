package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// messageTTL wraps a message with its expiration time.
type messageTTL struct {
	Message   Message
	ExpiresAt time.Time
}

// ShortTermMemory stores messages with time-to-live (TTL).
// Messages automatically expire and are removed after their TTL.
type ShortTermMemory struct {
	mu              sync.RWMutex
	messages        map[string]*messageTTL
	ttl             time.Duration
	cleanupInterval time.Duration
	cleanup         *time.Ticker
	stopCleanup     chan struct{}
}

// NewShortTerm creates a new short-term memory with the given options.
func NewShortTerm(opts ...Option) *ShortTermMemory {
	options := &ShortTermOptions{
		TTL:             30 * time.Minute,
		CleanupInterval: 5 * time.Minute,
	}
	for _, opt := range opts {
		opt(options)
	}

	stm := &ShortTermMemory{
		messages:        make(map[string]*messageTTL),
		ttl:             options.TTL,
		cleanupInterval: options.CleanupInterval,
		cleanup:         time.NewTicker(options.CleanupInterval),
		stopCleanup:     make(chan struct{}),
	}

	// Start background cleanup goroutine
	go stm.runCleanup()

	return stm
}

// runCleanup periodically removes expired messages.
func (stm *ShortTermMemory) runCleanup() {
	for {
		select {
		case <-stm.cleanup.C:
			stm.removeExpired()
		case <-stm.stopCleanup:
			return
		}
	}
}

// removeExpired removes all expired messages from memory.
func (stm *ShortTermMemory) removeExpired() {
	stm.mu.Lock()
	defer stm.mu.Unlock()

	now := time.Now()
	for id, mt := range stm.messages {
		if now.After(mt.ExpiresAt) {
			delete(stm.messages, id)
		}
	}
}

// Stop stops the background cleanup goroutine.
func (stm *ShortTermMemory) Stop() {
	stm.cleanup.Stop()
	close(stm.stopCleanup)
}

// Add stores a message with TTL in short-term memory.
func (stm *ShortTermMemory) Add(ctx context.Context, msg Message) error {
	stm.mu.Lock()
	defer stm.mu.Unlock()

	// Set timestamp if not provided
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	// Use message TTL if provided, otherwise use default
	ttl := stm.ttl
	if msg.TTL > 0 {
		ttl = msg.TTL
	}

	// Generate unique ID for the message
	id := uuid.New().String()

	stm.messages[id] = &messageTTL{
		Message:   msg,
		ExpiresAt: time.Now().Add(ttl),
	}

	return nil
}

// Get retrieves recent non-expired messages from memory.
func (stm *ShortTermMemory) Get(ctx context.Context, limit int) ([]Message, error) {
	stm.mu.Lock()
	defer stm.mu.Unlock()

	now := time.Now()
	var messages []Message

	// Collect non-expired messages and increment access count
	for _, mt := range stm.messages {
		if now.Before(mt.ExpiresAt) {
			mt.Message.AccessCount++
			messages = append(messages, mt.Message)
		}
	}

	// Sort by timestamp (oldest first) - O(n log n)
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})

	// Apply limit
	if limit > 0 && len(messages) > limit {
		messages = messages[len(messages)-limit:]
	}

	return messages, nil
}

// Search finds non-expired messages containing the query string.
func (stm *ShortTermMemory) Search(ctx context.Context, query string, limit int) ([]Message, error) {
	stm.mu.Lock()
	defer stm.mu.Unlock()

	now := time.Now()
	var results []Message
	query = strings.ToLower(query)

	for _, mt := range stm.messages {
		if now.Before(mt.ExpiresAt) && strings.Contains(strings.ToLower(mt.Message.Content), query) {
			mt.Message.AccessCount++
			results = append(results, mt.Message)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}

	return results, nil
}

// Clear removes all messages from memory.
func (stm *ShortTermMemory) Clear(ctx context.Context) error {
	stm.mu.Lock()
	defer stm.mu.Unlock()

	stm.messages = make(map[string]*messageTTL)
	return nil
}

// Stats returns statistics about the short-term memory.
func (stm *ShortTermMemory) Stats(ctx context.Context) (*Stats, error) {
	stm.mu.RLock()
	defer stm.mu.RUnlock()

	now := time.Now()
	stats := &Stats{}

	for _, mt := range stm.messages {
		// Only count non-expired messages
		if now.Before(mt.ExpiresAt) {
			stats.TotalMessages++

			switch mt.Message.Role {
			case "user":
				stats.UserMessages++
			case "assistant":
				stats.AssistantMessages++
			case "system":
				stats.SystemMessages++
			}

			// Track oldest and newest
			if stats.OldestMessage == nil || mt.Message.Timestamp.Before(*stats.OldestMessage) {
				oldest := mt.Message.Timestamp
				stats.OldestMessage = &oldest
			}
			if stats.NewestMessage == nil || mt.Message.Timestamp.After(*stats.NewestMessage) {
				newest := mt.Message.Timestamp
				stats.NewestMessage = &newest
			}

			stats.TotalSize += int64(len(mt.Message.Content))
		}
	}

	return stats, nil
}
