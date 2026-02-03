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
		log.Fatal("OPENAI_API_KEY environment variable is not set")
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
		agent.WithTemperature(0.7),
	)

	// Run agent with a simple query
	result, err := ag.Run(ctx, "Hello! Can you tell me what you can help me with?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Agent Response:", result.Content)
	fmt.Printf("Tokens used: %d\n", result.Usage.TotalTokens)
	
	// Continue the conversation
	result2, err := ag.Run(ctx, "What's the weather like in San Francisco?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nFollow-up Response:", result2.Content)
	fmt.Printf("Total messages in conversation: %d\n", result2.MessageCount)
}
