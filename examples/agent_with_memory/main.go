package main

import (
	"context"
	"fmt"
	"log"

	"github.com/EmadMokhtar/gogent/agent"
	"github.com/EmadMokhtar/gogent/llm/providers/openai"
	"github.com/EmadMokhtar/gogent/memory"
)

func main() {
	ctx := context.Background()

	// Create conversation memory
	mem := memory.NewConversation()

	// Create LLM (using stub implementation)
	llm := openai.New(openai.WithAPIKey("your-api-key"))

	// Create agent with memory
	ag := agent.New(
		agent.WithLLM(llm),
		agent.WithMemory(mem),
		agent.WithSystemPrompt("You are a helpful assistant with memory."),
	)

	// First interaction - introduce yourself
	fmt.Println("=== First Interaction ===")
	resp1, err := ag.Run(ctx, "My name is Emad and I'm a software engineer")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("User: My name is Emad and I'm a software engineer\n")
	fmt.Printf("Assistant: %s\n", resp1.Content)

	// Second interaction - ask about name
	fmt.Println("\n=== Second Interaction ===")
	resp2, err := ag.Run(ctx, "What's my name?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("User: What's my name?\n")
	fmt.Printf("Assistant: %s\n", resp2.Content)

	// Third interaction - ask about profession
	fmt.Println("\n=== Third Interaction ===")
	resp3, err := ag.Run(ctx, "What do I do for work?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("User: What do I do for work?\n")
	fmt.Printf("Assistant: %s\n", resp3.Content)

	// Check memory statistics
	fmt.Println("\n=== Memory Statistics ===")
	stats, _ := mem.Stats(ctx)
	fmt.Printf("Total messages in memory: %d\n", stats.TotalMessages)
	fmt.Printf("User messages: %d\n", stats.UserMessages)
	fmt.Printf("Assistant messages: %d\n", stats.AssistantMessages)
	fmt.Printf("System messages: %d\n", stats.SystemMessages)

	// Retrieve and display all messages
	fmt.Println("\n=== Conversation History ===")
	messages, _ := mem.Get(ctx, 0)
	for i, msg := range messages {
		if msg.Role == "system" {
			continue // Skip system prompt in display
		}
		fmt.Printf("%d. [%s]: %s\n", i+1, msg.Role, msg.Content)
	}

	// Demonstrate memory search
	fmt.Println("\n=== Memory Search ===")
	searchResults, _ := mem.Search(ctx, "name", 5)
	fmt.Printf("Found %d messages containing 'name':\n", len(searchResults))
	for i, msg := range searchResults {
		fmt.Printf("%d. [%s]: %s\n", i+1, msg.Role, msg.Content)
	}
}
