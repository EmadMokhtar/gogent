package memory

import "time"

// Message represents a single message in the conversation.
type Message struct {
	Role      string                 `json:"role"`      // "user", "assistant", "system"
	Content   string                 `json:"content"`   // The message content
	Timestamp time.Time              `json:"timestamp"` // When the message was created
	Metadata  map[string]interface{} `json:"metadata"`  // Additional metadata

	// Memory-specific fields
	Importance  float64       `json:"importance"`   // 0.0 - 1.0 (for long-term memory)
	TTL         time.Duration `json:"ttl"`          // For short-term memory
	AccessCount int           `json:"access_count"` // How many times accessed
}

// Stats represents statistics about the memory.
type Stats struct {
	TotalMessages     int        `json:"total_messages"`     // Total number of messages
	UserMessages      int        `json:"user_messages"`      // Number of user messages
	AssistantMessages int        `json:"assistant_messages"` // Number of assistant messages
	SystemMessages    int        `json:"system_messages"`    // Number of system messages
	OldestMessage     *time.Time `json:"oldest_message"`     // Timestamp of oldest message
	NewestMessage     *time.Time `json:"newest_message"`     // Timestamp of newest message
	TotalSize         int64      `json:"total_size"`         // Approximate bytes
}
