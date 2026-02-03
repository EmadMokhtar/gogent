// Package main provides a basic example of using the Gogent agent.
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

	// Run agent
	response, err := ag.Run(ctx, "What is the capital of France?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Response:", response.Content)
}
