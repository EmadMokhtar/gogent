// Package main provides an example of using the Gogent agent with tools.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/EmadMokhtar/gogent/agent"
	"github.com/EmadMokhtar/gogent/llm/providers/openai"
	"github.com/EmadMokhtar/gogent/tools/builtin"
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

	// Create agent with calculator tool
	ag := agent.New(
		agent.WithLLM(llmProvider),
		agent.WithTools(
			builtin.NewCalculator(),
		),
		agent.WithSystemPrompt("You are a helpful assistant with access to a calculator. Use it when you need to perform calculations."),
	)

	// Run agent with a calculation query
	response, err := ag.Run(ctx, "What is 157 * 23 + 891?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Response:", response.Content)
	fmt.Printf("\nMetadata: %+v\n", response.Metadata)
}
