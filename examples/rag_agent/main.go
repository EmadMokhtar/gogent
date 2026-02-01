package main

import (
	"context"
	"fmt"
	"log"

	"github.com/EmadMokhtar/gogent/agent"
	"github.com/EmadMokhtar/gogent/llm/providers/openai"
	"github.com/EmadMokhtar/gogent/rag"
)

func main() {
	fmt.Println("=== RAG Agent Example ===")
	fmt.Println()
	
	ctx := context.Background()
	
	// Create LLM provider
	llm := openai.New("gpt-4", openai.WithAPIKey("your-api-key-here"))
	
	// Create embedder
	embedder := openai.NewEmbedder("text-embedding-ada-002", 
		openai.WithEmbedderAPIKey("your-api-key-here"))
	
	// Create RAG system with in-memory vector store
	ragSystem := rag.New(
		rag.WithVectorStore(rag.NewInMemoryVectorStore()),
		rag.WithEmbedder(embedder),
		rag.WithChunkSize(500),
		rag.WithChunkOverlap(50),
	)
	
	// Add documents to RAG
	docs := []rag.Document{
		{
			ID:      "doc1",
			Content: "Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency.",
			Metadata: map[string]interface{}{
				"source": "wikipedia",
				"topic":  "programming",
			},
		},
		{
			ID:      "doc2",
			Content: "Agents are autonomous systems that can perceive their environment and take actions to achieve specific goals. In AI, agents can be simple or complex, ranging from basic reactive agents to sophisticated goal-based agents with learning capabilities.",
			Metadata: map[string]interface{}{
				"source": "ai-textbook",
				"topic":  "artificial-intelligence",
			},
		},
		{
			ID:      "doc3",
			Content: "The Go programming language is particularly well-suited for building concurrent systems and microservices. Its goroutines and channels make it easy to write programs that do many things at once. Go is also known for its fast compilation and excellent standard library.",
			Metadata: map[string]interface{}{
				"source": "tech-blog",
				"topic":  "programming",
			},
		},
	}
	
	fmt.Println("Adding documents to RAG system...")
	if err := ragSystem.AddDocuments(ctx, docs); err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}
	
	count, _ := ragSystem.Count(ctx)
	fmt.Printf("Added %d document chunks to vector store\n\n", count)
	
	// Create agent with RAG
	ag := agent.New(
		agent.WithLLM(llm),
		agent.WithRAG(ragSystem),
		agent.WithSystemPrompt("You are a helpful AI assistant. Use the provided context to answer questions accurately."),
		agent.WithTemperature(0.7),
		agent.WithMaxTokens(500),
	)
	
	// Query with RAG enhancement
	query := "What is Go and how can it be used for building agents?"
	fmt.Printf("User: %s\n", query)
	
	result, err := ag.Run(ctx, query)
	if err != nil {
		log.Fatalf("Error running agent: %v", err)
	}
	
	fmt.Printf("\nAgent: %s\n", result.Response)
	fmt.Printf("\nSources retrieved: %d\n", len(result.Sources))
	fmt.Printf("Duration: %v\n", result.Duration)
	
	// Test RAG retrieval directly
	fmt.Println("\n--- Direct RAG Query ---")
	ragQuery := "concurrent programming"
	fmt.Printf("Query: %s\n", ragQuery)
	
	ragResults, err := ragSystem.Query(ctx, ragQuery, 3)
	if err != nil {
		log.Fatalf("Error querying RAG: %v", err)
	}
	
	fmt.Printf("\nTop %d relevant chunks:\n", len(ragResults))
	for i, result := range ragResults {
		fmt.Printf("\n%d. Score: %.4f\n", i+1, result.Score)
		fmt.Printf("   Content: %s...\n", truncate(result.Chunk.Content, 100))
		if source, ok := result.Chunk.Metadata["source"]; ok {
			fmt.Printf("   Source: %v\n", source)
		}
	}
}

// truncate truncates a string to the specified length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
