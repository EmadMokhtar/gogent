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
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}

	// Create LLM provider
	llmProvider := openai.New(
		openai.WithModel("gpt-4"),
		openai.WithAPIKey(apiKey),
	)

	// Create agent with tools
	ag := agent.New(
		agent.WithLLM(llmProvider),
		agent.WithTools(
			builtin.NewCalculator(),
			builtin.NewEcho(),
		),
		agent.WithSystemPrompt("You are a helpful AI assistant with access to tools. Use the calculator tool for math operations."),
		agent.WithTemperature(0.7),
	)

	// Run agent with a calculation query
	fmt.Println("Query: What is 42 multiplied by 17?")
	result, err := ag.Run(ctx, "What is 42 multiplied by 17?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", result.Content)
	fmt.Printf("Tokens used: %d\n", result.Usage.TotalTokens)
	
	// Continue conversation with another calculation
	fmt.Println("\nQuery: Now add 100 to that result")
	result2, err := ag.Run(ctx, "Now add 100 to that result")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Follow-up:", result2.Content)
	fmt.Printf("Total messages: %d\n", result2.MessageCount)
	
	// Test the echo tool
	fmt.Println("\nQuery: Can you echo the message 'Hello from Gogent!'")
	result3, err := ag.Run(ctx, "Can you echo the message 'Hello from Gogent!'")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Echo Response:", result3.Content)
}
