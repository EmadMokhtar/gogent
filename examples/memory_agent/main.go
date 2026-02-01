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
	fmt.Println("=== Memory Agent Example ===")
	fmt.Println()
	
	ctx := context.Background()
	
	// Create LLM provider
	llm := openai.New("gpt-4", openai.WithAPIKey("your-api-key-here"))
	
	// Create conversation memory
	mem := memory.NewConversation(memory.WithMaxSize(50))
	
	// Create agent with memory
	ag := agent.New(
		agent.WithLLM(llm),
		agent.WithMemory(mem),
		agent.WithSystemPrompt("You are a helpful AI assistant. Remember the conversation context."),
		agent.WithTemperature(0.7),
		agent.WithMaxTokens(500),
	)
	
	// Simulate a conversation
	queries := []string{
		"Hi! My name is Alice and I'm learning about Go programming.",
		"What did I just tell you my name was?",
		"Can you help me understand goroutines?",
		"What topic did I ask about in my previous message?",
	}
	
	fmt.Println("--- Conversation with Memory ---")
	fmt.Println()
	for i, query := range queries {
		fmt.Printf("[Turn %d]\n", i+1)
		fmt.Printf("User: %s\n", query)
		
		result, err := ag.Run(ctx, query)
		if err != nil {
			log.Fatalf("Error running agent: %v", err)
		}
		
		fmt.Printf("Agent: %s\n", result.Response)
		fmt.Printf("Duration: %v\n\n", result.Duration)
	}
	
	// Check memory contents
	fmt.Println("--- Memory Contents ---")
	count, _ := mem.Count(ctx)
	fmt.Printf("Total messages in memory: %d\n\n", count)
	
	allMessages, _ := mem.Get(ctx, 10)
	for i, msg := range allMessages {
		fmt.Printf("%d. [%s] %s\n", i+1, msg.Role, truncate(msg.Content, 60))
	}
	
	// Test memory search
	fmt.Println("\n--- Memory Search ---")
	searchQuery := "name"
	fmt.Printf("Searching for: %s\n", searchQuery)
	
	searchResults, err := mem.Search(ctx, searchQuery)
	if err != nil {
		log.Fatalf("Error searching memory: %v", err)
	}
	
	fmt.Printf("Found %d matching messages:\n", len(searchResults))
	for i, msg := range searchResults {
		fmt.Printf("%d. [%s] %s\n", i+1, msg.Role, msg.Content)
	}
	
	// Demonstrate different memory types
	fmt.Println("\n--- Different Memory Types ---")
	
	// Short-term memory example
	fmt.Println("\n1. Short-term Memory (with TTL):")
	shortMem := memory.NewShortTerm()
	shortMem.Add(ctx, memory.NewUserMessage("This is a short-term message"))
	count, _ = shortMem.Count(ctx)
	fmt.Printf("   Messages in short-term memory: %d\n", count)
	
	// Long-term memory example
	fmt.Println("\n2. Long-term Memory (with importance):")
	longMem := memory.NewLongTerm()
	longMem.AddWithImportance(ctx, memory.NewUserMessage("This is an important message"), 0.9)
	longMem.AddWithImportance(ctx, memory.NewUserMessage("This is a less important message"), 0.3)
	count, _ = longMem.Count(ctx)
	fmt.Printf("   Messages in long-term memory: %d\n", count)
	
	// Retrieve high-importance messages
	importantMsgs, _ := longMem.GetByImportance(ctx, 0.7, 10)
	fmt.Printf("   High-importance messages (>0.7): %d\n", len(importantMsgs))
}

// truncate truncates a string to the specified length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
