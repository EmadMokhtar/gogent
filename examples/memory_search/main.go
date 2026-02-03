package main

import (
	"context"
	"fmt"
	"time"

	"github.com/EmadMokhtar/gogent/memory"
)

func main() {
	ctx := context.Background()

	fmt.Println("=== Memory Search Examples ===")

	// Example 1: Search in Conversation Memory
	fmt.Println("\n1. Conversation Memory Search")
	demoConversationSearch(ctx)

	// Example 2: Search in Short-term Memory
	fmt.Println("\n2. Short-term Memory Search")
	demoShortTermSearch(ctx)

	// Example 3: Search in Long-term Memory (sorted by importance)
	fmt.Println("\n3. Long-term Memory Search (by importance)")
	demoLongTermSearch(ctx)
}

func demoConversationSearch(ctx context.Context) {
	mem := memory.NewConversation()

	// Add sample messages
	messages := []string{
		"I love programming in Go",
		"Python is great for data science",
		"JavaScript is essential for web development",
		"I prefer Go for backend services",
		"Rust is interesting for systems programming",
	}

	for _, content := range messages {
		mem.Add(ctx, memory.Message{
			Role:    "user",
			Content: content,
		})
	}

	// Search for "Go"
	results, err := mem.Search(ctx, "Go", 0)
	if err != nil {
		fmt.Printf("Search failed: %v\n", err)
		return
	}

	fmt.Printf("Search results for 'Go': %d matches\n", len(results))
	for i, msg := range results {
		fmt.Printf("  %d. %s\n", i+1, msg.Content)
	}

	// Search with limit
	results, _ = mem.Search(ctx, "programming", 2)
	fmt.Printf("\nSearch results for 'programming' (limit 2): %d matches\n", len(results))
	for i, msg := range results {
		fmt.Printf("  %d. %s\n", i+1, msg.Content)
	}

	// Case-insensitive search
	results, _ = mem.Search(ctx, "PYTHON", 0)
	fmt.Printf("\nCase-insensitive search for 'PYTHON': %d matches\n", len(results))
	for i, msg := range results {
		fmt.Printf("  %d. %s\n", i+1, msg.Content)
	}
}

func demoShortTermSearch(ctx context.Context) {
	mem := memory.NewShortTerm(memory.WithTTL(10 * time.Second))
	defer mem.Stop()

	// Add messages
	messages := []string{
		"Meeting scheduled for tomorrow",
		"Don't forget to buy milk",
		"Call the dentist about appointment",
		"Meeting notes from yesterday",
	}

	for _, content := range messages {
		mem.Add(ctx, memory.Message{
			Role:    "user",
			Content: content,
		})
	}

	// Search for "meeting"
	results, _ := mem.Search(ctx, "meeting", 0)
	fmt.Printf("Search results for 'meeting': %d matches\n", len(results))
	for i, msg := range results {
		fmt.Printf("  %d. %s (AccessCount: %d)\n", i+1, msg.Content, msg.AccessCount)
	}

	// Search again to see access count increase
	results, _ = mem.Search(ctx, "meeting", 0)
	fmt.Printf("\nSecond search (access count should increase):\n")
	for i, msg := range results {
		fmt.Printf("  %d. %s (AccessCount: %d)\n", i+1, msg.Content, msg.AccessCount)
	}
}

func demoLongTermSearch(ctx context.Context) {
	mem := memory.NewLongTerm(memory.WithCapacity(10))

	// Add messages with importance scores
	type importantMsg struct {
		content    string
		importance float64
	}

	messages := []importantMsg{
		{"User prefers dark mode in all applications", 0.9},
		{"User mentioned liking dark chocolate", 0.3},
		{"User's email is user@example.com", 0.8},
		{"User likes to work in dark, quiet environments", 0.7},
		{"Random note about the dark knight movie", 0.2},
	}

	for _, msg := range messages {
		mem.Add(ctx, memory.Message{
			Role:       "user",
			Content:    msg.content,
			Importance: msg.importance,
		})
	}

	// Search for "dark" - results sorted by importance
	results, _ := mem.Search(ctx, "dark", 0)
	fmt.Printf("Search results for 'dark' (sorted by importance): %d matches\n", len(results))
	for i, msg := range results {
		fmt.Printf("  %d. [Importance: %.1f] %s\n", i+1, msg.Importance, msg.Content)
	}

	// Search with limit - get only most important matches
	results, _ = mem.Search(ctx, "dark", 2)
	fmt.Printf("\nTop 2 most important results for 'dark':\n")
	for i, msg := range results {
		fmt.Printf("  %d. [Importance: %.1f] %s\n", i+1, msg.Importance, msg.Content)
	}

	// Search for user preferences
	results, _ = mem.Search(ctx, "user", 0)
	fmt.Printf("\nSearch results for 'user' (by importance):\n")
	for i, msg := range results {
		fmt.Printf("  %d. [Importance: %.1f, Access: %d] %s\n",
			i+1, msg.Importance, msg.AccessCount, msg.Content)
	}
}
