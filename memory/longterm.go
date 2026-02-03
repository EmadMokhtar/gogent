package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"
)

// LongTermMemory stores messages based on importance scores.
// It keeps a fixed capacity of the most important messages.
type LongTermMemory struct {
	mu                  sync.RWMutex
	messages            []Message
	capacity            int
	importanceThreshold float64
}

// NewLongTerm creates a new long-term memory with the given options.
func NewLongTerm(opts ...Option) *LongTermMemory {
	options := &LongTermOptions{
		Capacity:            100,
		ImportanceThreshold: 0.0,
	}
	for _, opt := range opts {
		opt(options)
	}

	return &LongTermMemory{
		messages:            make([]Message, 0),
		capacity:            options.Capacity,
		importanceThreshold: options.ImportanceThreshold,
	}
}

// Add stores a message in long-term memory based on importance.
func (ltm *LongTermMemory) Add(ctx context.Context, msg Message) error {
	ltm.mu.Lock()
	defer ltm.mu.Unlock()

	// Set timestamp if not provided
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	// Skip if below importance threshold
	if msg.Importance < ltm.importanceThreshold {
		return nil
	}

	// Add the message
	ltm.messages = append(ltm.messages, msg)

	// Prune if over capacity
	if len(ltm.messages) > ltm.capacity {
		ltm.pruneMessages()
	}

	return nil
}

// pruneMessages keeps only the most important messages up to capacity.
// Must be called while holding the write lock.
func (ltm *LongTermMemory) pruneMessages() {
	// Sort by importance (descending) and access count (descending)
	sort.Slice(ltm.messages, func(i, j int) bool {
		if ltm.messages[i].Importance != ltm.messages[j].Importance {
			return ltm.messages[i].Importance > ltm.messages[j].Importance
		}
		return ltm.messages[i].AccessCount > ltm.messages[j].AccessCount
	})

	// Keep only top capacity messages
	ltm.messages = ltm.messages[:ltm.capacity]
}

// Get retrieves messages from long-term memory.
// Messages are returned in chronological order.
func (ltm *LongTermMemory) Get(ctx context.Context, limit int) ([]Message, error) {
	ltm.mu.Lock()
	defer ltm.mu.Unlock()

	// Increment access count first (need write lock for this)
	for i := range ltm.messages {
		ltm.messages[i].AccessCount++
	}

	// Create a copy and sort by timestamp
	messages := make([]Message, len(ltm.messages))
	copy(messages, ltm.messages)

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})

	// Apply limit
	if limit > 0 && len(messages) > limit {
		messages = messages[len(messages)-limit:]
	}

	return messages, nil
}

// Search finds messages containing the query string, sorted by importance.
func (ltm *LongTermMemory) Search(ctx context.Context, query string, limit int) ([]Message, error) {
	ltm.mu.Lock()
	defer ltm.mu.Unlock()

	var results []Message
	query = strings.ToLower(query)

	// Collect matching messages and increment access count
	for i := range ltm.messages {
		if strings.Contains(strings.ToLower(ltm.messages[i].Content), query) {
			ltm.messages[i].AccessCount++
			results = append(results, ltm.messages[i])
		}
	}

	// Sort by importance (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Importance > results[j].Importance
	})

	// Apply limit
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// Clear removes all messages from long-term memory.
func (ltm *LongTermMemory) Clear(ctx context.Context) error {
	ltm.mu.Lock()
	defer ltm.mu.Unlock()

	ltm.messages = make([]Message, 0)
	return nil
}

// Stats returns statistics about the long-term memory.
func (ltm *LongTermMemory) Stats(ctx context.Context) (*Stats, error) {
	ltm.mu.RLock()
	defer ltm.mu.RUnlock()

	stats := &Stats{
		TotalMessages: len(ltm.messages),
	}

	for i, msg := range ltm.messages {
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
			newest := msg.Timestamp
			stats.NewestMessage = &newest
		} else {
			if msg.Timestamp.Before(*stats.OldestMessage) {
				oldest := msg.Timestamp
				stats.OldestMessage = &oldest
			}
			if msg.Timestamp.After(*stats.NewestMessage) {
				newest := msg.Timestamp
				stats.NewestMessage = &newest
			}
		}

		stats.TotalSize += int64(len(msg.Content))
	}

	return stats, nil
}
