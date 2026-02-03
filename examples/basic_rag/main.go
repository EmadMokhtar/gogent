package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// Create components
	embedder := openai.NewEmbedder(openai.WithAPIKey(apiKey))
	store := rag.NewInMemoryVectorStore()
	chunker := rag.NewFixedSizeChunker(
		rag.WithChunkSize(500),
		rag.WithOverlap(50),
	)

	// Create RAG system
	ragSys := rag.New(
		rag.WithVectorStore(store),
		rag.WithEmbedder(embedder),
		rag.WithChunker(chunker),
	)

	// Add documents
	docs := []rag.Document{
		{
			ID:      "doc1",
			Content: "Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency.",
			Metadata: map[string]interface{}{
				"source":   "wikipedia",
				"category": "programming",
			},
		},
		{
			ID:      "doc2",
			Content: "Agents are autonomous systems that perceive their environment and take actions to achieve goals. In AI, intelligent agents can learn from experience and adapt their behavior. They can be simple reactive agents or complex reasoning agents with planning capabilities.",
			Metadata: map[string]interface{}{
				"source":   "textbook",
				"category": "ai",
			},
		},
		{
			ID:      "doc3",
			Content: "Vector embeddings are numerical representations of text that capture semantic meaning. They allow us to perform similarity search and retrieve relevant documents based on meaning rather than just keyword matching.",
			Metadata: map[string]interface{}{
				"source":   "documentation",
				"category": "ai",
			},
		},
	}

	fmt.Println("Adding documents to RAG system...")
	if err := ragSys.AddDocuments(ctx, docs); err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}

	// Get stats
	stats, err := ragSys.Stats(ctx)
	if err != nil {
		log.Fatalf("Failed to get stats: %v", err)
	}
	fmt.Printf("Vector store stats: %d documents, dimension: %d\n\n", stats.TotalDocuments, stats.Dimension)

	// Retrieve relevant documents
	queries := []string{
		"What is Go programming language?",
		"Tell me about AI agents",
		"How do embeddings work?",
	}

	for _, query := range queries {
		fmt.Printf("Query: %s\n", query)
		results, err := ragSys.Retrieve(ctx, query, rag.WithRetrievalK(2))
		if err != nil {
			log.Printf("Failed to retrieve: %v", err)
			continue
		}

		for _, result := range results {
			fmt.Printf("  [Score: %.3f] %s\n", result.Score, truncate(result.Document.Content, 80))
		}
		fmt.Println()
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
