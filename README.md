# Gogent

Gogent is a high-performance, multi-agent framework for Go with comprehensive memory capabilities.

## Features

- **Multiple Memory Types**: Conversation, Short-term (TTL-based), and Long-term (importance-based)
- **Agent System**: AI agents with built-in memory integration
- **Thread-Safe**: All memory operations are safe for concurrent use
- **Flexible Configuration**: Customizable options for each memory type
- **Search Capabilities**: Search through memories with different strategies

## Installation

```bash
go get github.com/EmadMokhtar/gogent
```

## Quick Start

### Basic Agent with Conversation Memory

```go
package main

import (
    "context"
    "fmt"

    "github.com/EmadMokhtar/gogent/agent"
    "github.com/EmadMokhtar/gogent/memory"
    "github.com/EmadMokhtar/gogent/llm/providers/openai"
)

func main() {
    ctx := context.Background()

    // Create memory
    mem := memory.NewConversation()

    // Create LLM
    llm := openai.New(openai.WithAPIKey("your-api-key"))

    // Create agent
    ag := agent.New(
        agent.WithLLM(llm),
        agent.WithMemory(mem),
        agent.WithSystemPrompt("You are a helpful assistant."),
    )

    // Use the agent
    resp, _ := ag.Run(ctx, "My name is Alice")
    fmt.Println(resp.Content)

    resp, _ = ag.Run(ctx, "What's my name?")
    fmt.Println(resp.Content) // Agent remembers: "Alice"
}
```

## Memory Types

### 1. Conversation Memory

Stores complete conversation history with optional size limits.

```go
mem := memory.NewConversation(
    memory.WithMaxSize(1000), // Optional: limit to 1000 messages
)

// Add messages
mem.Add(ctx, memory.Message{
    Role:    "user",
    Content: "Hello!",
})

// Get all messages
messages, _ := mem.Get(ctx, 0)

// Search messages
results, _ := mem.Search(ctx, "hello", 10)
```

**Use case**: Full context retention for sessions, chat history.

### 2. Short-term Memory

Session-based storage with automatic expiration (TTL).

```go
mem := memory.NewShortTerm(
    memory.WithTTL(30 * time.Minute),        // Default TTL
    memory.WithCleanupInterval(5 * time.Minute), // Cleanup frequency
)
defer mem.Stop() // Important: stop background cleanup

// Add message with default TTL
mem.Add(ctx, memory.Message{
    Role:    "user",
    Content: "Temporary note",
})

// Add message with custom TTL
mem.Add(ctx, memory.Message{
    Role:    "user",
    Content: "Expires in 5 minutes",
    TTL:     5 * time.Minute,
})

// Messages automatically expire after TTL
```

**Use case**: Working memory for current tasks, temporary context.

### 3. Long-term Memory

Persistent storage with importance-based retention.

```go
mem := memory.NewLongTerm(
    memory.WithCapacity(100),                  // Keep 100 most important
    memory.WithImportanceThreshold(0.5),      // Min importance: 0.5
)

// Add important information
mem.Add(ctx, memory.Message{
    Role:       "user",
    Content:    "User prefers dark mode",
    Importance: 0.9, // High importance (0.0 - 1.0)
})

// Less important messages are pruned when capacity is reached
// Search returns results sorted by importance
results, _ := mem.Search(ctx, "dark mode", 5)
```

**Use case**: Key facts, user preferences, important events.

## Memory Interface

All memory types implement the same interface:

```go
type Memory interface {
    // Add a message to memory
    Add(ctx context.Context, msg Message) error
    
    // Get recent messages (limit = 0 means all)
    Get(ctx context.Context, limit int) ([]Message, error)
    
    // Search memory by query string
    Search(ctx context.Context, query string, limit int) ([]Message, error)
    
    // Clear all memory
    Clear(ctx context.Context) error
    
    // Get memory statistics
    Stats(ctx context.Context) (*Stats, error)
}
```

## Memory Statistics

Get insights about your memory:

```go
stats, _ := mem.Stats(ctx)

fmt.Printf("Total messages: %d\n", stats.TotalMessages)
fmt.Printf("User messages: %d\n", stats.UserMessages)
fmt.Printf("Assistant messages: %d\n", stats.AssistantMessages)
fmt.Printf("Total size: %d bytes\n", stats.TotalSize)
```

## When to Use Which Memory Type

| Memory Type | Best For | Retention | Capacity |
|-------------|----------|-----------|----------|
| **Conversation** | Chat history, full context | Unlimited (or max size) | Manual limit |
| **Short-term** | Temporary context, current tasks | Time-based (TTL) | Unlimited |
| **Long-term** | Important facts, preferences | Importance-based | Fixed capacity |

## Examples

See the `examples/` directory for complete working examples:

- `examples/memory_types/` - Demonstration of all memory types
- `examples/agent_with_memory/` - Agent with conversation memory
- `examples/memory_search/` - Advanced search examples

Run an example:

```bash
cd examples/memory_types
go run main.go
```

## Thread Safety

All memory implementations are thread-safe and can be used concurrently from multiple goroutines.

```go
// Safe to use from multiple goroutines
go mem.Add(ctx, msg1)
go mem.Add(ctx, msg2)
go mem.Get(ctx, 10)
```

## Testing

Run all tests:

```bash
go test ./...
```

Run with race detector:

```bash
go test -race ./...
```

## License

See LICENSE file for details.
