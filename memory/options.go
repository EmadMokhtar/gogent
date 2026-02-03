package memory

import "time"

// Option is a function that configures a memory implementation.
type Option func(interface{})

// ConversationOptions holds configuration for conversation memory.
type ConversationOptions struct {
	MaxSize int // Maximum number of messages to keep (0 = unlimited)
}

// ShortTermOptions holds configuration for short-term memory.
type ShortTermOptions struct {
	TTL             time.Duration // Default time-to-live for messages
	CleanupInterval time.Duration // How often to clean up expired messages
}

// LongTermOptions holds configuration for long-term memory.
type LongTermOptions struct {
	Capacity            int     // Maximum number of messages to keep
	ImportanceThreshold float64 // Minimum importance score to keep (0.0 - 1.0)
}

// WithMaxSize sets the maximum number of messages for conversation memory.
func WithMaxSize(size int) Option {
	return func(opts interface{}) {
		if co, ok := opts.(*ConversationOptions); ok {
			co.MaxSize = size
		}
	}
}

// WithTTL sets the default time-to-live for short-term memory messages.
func WithTTL(ttl time.Duration) Option {
	return func(opts interface{}) {
		if sto, ok := opts.(*ShortTermOptions); ok {
			sto.TTL = ttl
		}
	}
}

// WithCleanupInterval sets the cleanup interval for short-term memory.
func WithCleanupInterval(interval time.Duration) Option {
	return func(opts interface{}) {
		if sto, ok := opts.(*ShortTermOptions); ok {
			sto.CleanupInterval = interval
		}
	}
}

// WithCapacity sets the maximum capacity for long-term memory.
func WithCapacity(capacity int) Option {
	return func(opts interface{}) {
		if lto, ok := opts.(*LongTermOptions); ok {
			lto.Capacity = capacity
		}
	}
}

// WithImportanceThreshold sets the minimum importance threshold for long-term memory.
func WithImportanceThreshold(threshold float64) Option {
	return func(opts interface{}) {
		if lto, ok := opts.(*LongTermOptions); ok {
			lto.ImportanceThreshold = threshold
		}
	}
}
