package main

import (
	"context"
	"fmt"
	"time"

	"github.com/EmadMokhtar/gogent/memory"
)

func main() {
	ctx := context.Background()

	fmt.Println("=== Conversation Memory Demo ===")
	demoConversationMemory(ctx)

	fmt.Println("\n=== Short-term Memory Demo ===")
	demoShortTermMemory(ctx)

	fmt.Println("\n=== Long-term Memory Demo ===")
	demoLongTermMemory(ctx)
}

func demoConversationMemory(ctx context.Context) {
	// Create conversation memory with optional size limit
	mem := memory.NewConversation(memory.WithMaxSize(100))

	// Add messages
	mem.Add(ctx, memory.Message{
		Role:    "user",
		Content: "Hello, how are you?",
	})

	mem.Add(ctx, memory.Message{
		Role:    "assistant",
		Content: "I'm doing well, thank you!",
	})

	mem.Add(ctx, memory.Message{
		Role:    "user",
		Content: "What's the weather like?",
	})

	// Get all messages
	messages, _ := mem.Get(ctx, 0)
	fmt.Printf("Total messages: %d\n", len(messages))
	for i, msg := range messages {
		fmt.Printf("%d. [%s]: %s\n", i+1, msg.Role, msg.Content)
	}

	// Search messages
	results, _ := mem.Search(ctx, "weather", 5)
	fmt.Printf("\nSearch results for 'weather': %d matches\n", len(results))

	// Get stats
	stats, _ := mem.Stats(ctx)
	fmt.Printf("Stats: %d total, %d user, %d assistant\n",
		stats.TotalMessages, stats.UserMessages, stats.AssistantMessages)
}

func demoShortTermMemory(ctx context.Context) {
	// Create short-term memory with 5 second TTL
	mem := memory.NewShortTerm(
		memory.WithTTL(5*time.Second),
		memory.WithCleanupInterval(2*time.Second),
	)
	defer mem.Stop()

	// Add a temporary message
	mem.Add(ctx, memory.Message{
		Role:    "user",
		Content: "This will expire in 5 seconds",
	})

	// Check immediately
	messages, _ := mem.Get(ctx, 0)
	fmt.Printf("Messages immediately after add: %d\n", len(messages))

	// Wait 3 seconds
	fmt.Println("Waiting 3 seconds...")
	time.Sleep(3 * time.Second)
	messages, _ = mem.Get(ctx, 0)
	fmt.Printf("Messages after 3 seconds: %d\n", len(messages))

	// Wait another 3 seconds (total 6 seconds)
	fmt.Println("Waiting another 3 seconds...")
	time.Sleep(3 * time.Second)
	messages, _ = mem.Get(ctx, 0)
	fmt.Printf("Messages after 6 seconds (expired): %d\n", len(messages))

	// Add message with custom TTL
	mem.Add(ctx, memory.Message{
		Role:    "user",
		Content: "Custom TTL message",
		TTL:     2 * time.Second,
	})

	messages, _ = mem.Get(ctx, 0)
	fmt.Printf("Messages with custom TTL: %d\n", len(messages))

	time.Sleep(3 * time.Second)
	messages, _ = mem.Get(ctx, 0)
	fmt.Printf("After custom TTL expired: %d\n", len(messages))
}

func demoLongTermMemory(ctx context.Context) {
	// Create long-term memory with capacity of 3
	mem := memory.NewLongTerm(
		memory.WithCapacity(3),
		memory.WithImportanceThreshold(0.3),
	)

	// Add messages with different importance
	messages := []struct {
		content    string
		importance float64
	}{
		{"User prefers dark mode", 0.9},
		{"User likes coffee", 0.7},
		{"Random comment", 0.2}, // Below threshold, won't be stored
		{"User's birthday is May 15", 0.8},
		{"Another preference", 0.6},
	}

	for _, msg := range messages {
		mem.Add(ctx, memory.Message{
			Role:       "user",
			Content:    msg.content,
			Importance: msg.importance,
		})
	}

	// Get all messages (should only have 3 most important)
	stored, _ := mem.Get(ctx, 0)
	fmt.Printf("Stored messages (capacity=3): %d\n", len(stored))
	for i, msg := range stored {
		fmt.Printf("%d. [Importance: %.1f] %s\n", i+1, msg.Importance, msg.Content)
	}

	// Search by importance
	results, _ := mem.Search(ctx, "user", 5)
	fmt.Printf("\nSearch results for 'user': %d matches\n", len(results))
	for i, msg := range results {
		fmt.Printf("%d. [Importance: %.1f] %s\n", i+1, msg.Importance, msg.Content)
	}

	// Get stats
	stats, _ := mem.Stats(ctx)
	fmt.Printf("\nStats: %d messages, %d bytes\n", stats.TotalMessages, stats.TotalSize)
}
