// Package main provides an example of using streaming with the Gogent agent.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/EmadMokhtar/gogent/agent"
	"github.com/EmadMokhtar/gogent/llm/providers/openai"
)

func main() {
	ctx := context.Background()

	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create LLM provider
	llmProvider := openai.New(
		openai.WithModel("gpt-4"),
		openai.WithAPIKey(apiKey),
	)

	// Create agent
	ag := agent.New(
		agent.WithLLM(llmProvider),
		agent.WithSystemPrompt("You are a helpful AI assistant."),
	)

	// Stream response
	fmt.Print("Response: ")
	tokenChan, err := ag.Stream(ctx, "Write a short poem about Go programming language.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Print tokens as they arrive
	for token := range tokenChan {
		if token.Error != nil {
			log.Fatalf("Streaming error: %v", token.Error)
		}
		
		fmt.Print(token.Delta)
		
		if token.Done {
			fmt.Println("\n\nStreaming complete.")
			break
		}
	}
}
