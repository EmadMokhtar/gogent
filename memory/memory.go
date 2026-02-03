package memory

import "context"

// Memory defines the interface for memory implementations.
// It provides methods for storing, retrieving, and searching messages.
type Memory interface {
	// Add stores a message in memory.
	Add(ctx context.Context, msg Message) error

	// Get retrieves recent messages from memory.
	// If limit is 0, all messages are returned.
	// Messages are returned in chronological order (oldest first).
	Get(ctx context.Context, limit int) ([]Message, error)

	// Search finds messages matching the query string.
	// Returns up to 'limit' matching messages.
	Search(ctx context.Context, query string, limit int) ([]Message, error)

	// Clear removes all messages from memory.
	Clear(ctx context.Context) error

	// Stats returns statistics about the memory.
	Stats(ctx context.Context) (*Stats, error)
}
