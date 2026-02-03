package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/EmadMokhtar/gogent/agent"
	"github.com/EmadMokhtar/gogent/llm/providers/openai"
	"github.com/EmadMokhtar/gogent/rag"
)

func main() {
	ctx := context.Background()

	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Setup RAG
	embedder := openai.NewEmbedder(openai.WithAPIKey(apiKey))
	store := rag.NewInMemoryVectorStore()
	ragSys := rag.New(
		rag.WithVectorStore(store),
		rag.WithEmbedder(embedder),
	)

	// Add knowledge base
	fmt.Println("Building knowledge base...")
	docs := []rag.Document{
		{
			ID:      "1",
			Content: "The company was founded in 2020 by Jane Smith and John Doe.",
			Metadata: map[string]interface{}{
				"category": "company_info",
			},
		},
		{
			ID:      "2",
			Content: "Our main product is an AI-powered assistant that helps businesses automate customer support.",
			Metadata: map[string]interface{}{
				"category": "products",
			},
		},
		{
			ID:      "3",
			Content: "We have 50 employees across 10 countries, with offices in San Francisco, London, and Tokyo.",
			Metadata: map[string]interface{}{
				"category": "company_info",
			},
		},
		{
			ID:      "4",
			Content: "Our pricing starts at $99/month for the basic plan and goes up to $999/month for the enterprise plan.",
			Metadata: map[string]interface{}{
				"category": "pricing",
			},
		},
	}

	if err := ragSys.AddDocuments(ctx, docs); err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}

	// Create agent with RAG
	llm := openai.New(openai.WithLLMAPIKey(apiKey))
	ag := agent.New(
		agent.WithLLM(llm),
		agent.WithRAG(ragSys),
		agent.WithSystemPrompt("You are a helpful assistant. Use the provided context to answer questions accurately and concisely."),
		agent.WithRAGTopK(3),
	)

	// Query with RAG context
	queries := []string{
		"When was the company founded?",
		"How many employees do we have?",
		"What is the pricing?",
	}

	for _, query := range queries {
		fmt.Printf("\nQ: %s\n", query)
		response, err := ag.Run(ctx, query)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		fmt.Printf("A: %s\n", response.Content)

		// Show sources
		if len(response.Sources) > 0 {
			fmt.Println("\nSources:")
			for _, src := range response.Sources {
				fmt.Printf("  [Score: %.3f] %s\n", src.Score, truncate(src.Document.Content, 60))
			}
		}
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
