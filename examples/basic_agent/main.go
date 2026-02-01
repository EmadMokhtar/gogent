package main

import (
	"context"
	"fmt"
	"log"

	"github.com/EmadMokhtar/gogent/agent"
	"github.com/EmadMokhtar/gogent/llm/providers/openai"
	"github.com/EmadMokhtar/gogent/tools/builtin"
)

func main() {
	fmt.Println("=== Basic Agent Example ===")
	fmt.Println()
	
	ctx := context.Background()
	
	// Create LLM provider
	// Note: In a real application, you would set a valid API key
	llm := openai.New("gpt-4", openai.WithAPIKey("your-api-key-here"))
	
	// Create agent with tools
	ag := agent.New(
		agent.WithLLM(llm),
		agent.WithTools(
			builtin.NewCalculator(),
			builtin.NewWebSearch(),
		),
		agent.WithSystemPrompt("You are a helpful AI assistant with access to tools."),
		agent.WithTemperature(0.7),
		agent.WithMaxTokens(500),
	)
	
	// Run agent with a simple query
	query := "Hello! Can you help me with some calculations?"
	fmt.Printf("User: %s\n", query)
	
	result, err := ag.Run(ctx, query)
	if err != nil {
		log.Fatalf("Error running agent: %v", err)
	}
	
	fmt.Printf("\nAgent: %s\n", result.Response)
	fmt.Printf("\nDuration: %v\n", result.Duration)
	
	// Example of using a tool through the agent
	fmt.Println("\n--- Testing Calculator Tool ---")
	calcQuery := "What is 25 + 17?"
	fmt.Printf("User: %s\n", calcQuery)
	
	result2, err := ag.Run(ctx, calcQuery)
	if err != nil {
		log.Fatalf("Error running agent: %v", err)
	}
	
	fmt.Printf("\nAgent: %s\n", result2.Response)
	fmt.Printf("Duration: %v\n", result2.Duration)
	
	// Test the calculator tool directly
	fmt.Println("\n--- Direct Tool Usage ---")
	executor := ag.GetExecutor()
	calcResult, err := executor.Execute(ctx, "calculator", map[string]interface{}{
		"operation": "multiply",
		"a":         12,
		"b":         8,
	})
	
	if err != nil {
		log.Fatalf("Error executing calculator: %v", err)
	}
	
	fmt.Printf("Direct calculator result (12 * 8): %v\n", calcResult.Output)
}
