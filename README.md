# Gogent - Agentic Framework for Go

Gogent is a comprehensive agentic framework in Golang inspired by Pydantic AI, featuring memory management and in-memory RAG (Retrieval-Augmented Generation) capabilities.

## Features

- 🤖 **Agent Core**: Flexible agent abstraction with tool execution and memory integration
- 🧠 **Memory System**: Conversation, short-term, and long-term memory with different retention strategies
- 📚 **In-Memory RAG**: Vector store with similarity search, document chunking, and retrieval
- 🛠️ **Tool System**: Extensible tool registry with built-in tools (calculator, web search)
- 🔌 **LLM Integration**: Provider abstraction supporting OpenAI, Anthropic, and custom providers
- ⚡ **Concurrent Processing**: Parallel embedding generation and batch operations
- 🎯 **Type-Safe**: Leverages Go's type system for compile-time safety
- 🧪 **Well-Tested**: Comprehensive test coverage for core components

## Installation

```bash
go get github.com/EmadMokhtar/gogent
```

## Quick Start

### Basic Agent

```go
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
    ctx := context.Background()

    // Create LLM provider
    llm := openai.New("gpt-4", openai.WithAPIKey("your-api-key"))

    // Create agent with tools
    ag := agent.New(
        agent.WithLLM(llm),
        agent.WithTools(
            builtin.NewCalculator(),
            builtin.NewWebSearch(),
        ),
        agent.WithSystemPrompt("You are a helpful AI assistant."),
    )

    // Run agent
    result, err := ag.Run(ctx, "What is 25 + 17?")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response:", result.Response)
}
```

### RAG-Enhanced Agent

```go
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
    ctx := context.Background()

    // Create RAG system
    ragSystem := rag.New(
        rag.WithVectorStore(rag.NewInMemoryVectorStore()),
        rag.WithEmbedder(openai.NewEmbedder("text-embedding-ada-002")),
        rag.WithChunkSize(500),
        rag.WithChunkOverlap(50),
    )

    // Add documents
    docs := []rag.Document{
        {
            ID:      "doc1",
            Content: "Go is a statically typed, compiled programming language.",
            Metadata: map[string]interface{}{"source": "wikipedia"},
        },
    }
    
    if err := ragSystem.AddDocuments(ctx, docs); err != nil {
        log.Fatal(err)
    }

    // Create agent with RAG
    llm := openai.New("gpt-4", openai.WithAPIKey("your-api-key"))
    ag := agent.New(
        agent.WithLLM(llm),
        agent.WithRAG(ragSystem),
    )

    // Query with RAG enhancement
    result, err := ag.Run(ctx, "What is Go?")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response:", result.Response)
    fmt.Println("Sources:", len(result.Sources))
}
```

### Memory-Enabled Agent

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/EmadMokhtar/gogent/agent"
    "github.com/EmadMokhtar/gogent/llm/providers/openai"
    "github.com/EmadMokhtar/gogent/memory"
)

func main() {
    ctx := context.Background()

    // Create conversation memory
    mem := memory.NewConversation()

    // Create agent with memory
    llm := openai.New("gpt-4", openai.WithAPIKey("your-api-key"))
    ag := agent.New(
        agent.WithLLM(llm),
        agent.WithMemory(mem),
    )

    // Conversation with memory
    result1, _ := ag.Run(ctx, "My name is Alice.")
    fmt.Println(result1.Response)

    result2, _ := ag.Run(ctx, "What is my name?")
    fmt.Println(result2.Response) // Will remember "Alice"
}
```

## Architecture

### Core Components

#### 1. Agent Core (`agent/`)

The agent is the main orchestrator that coordinates LLM calls, tool execution, memory, and RAG.

```go
type Agent struct {
    opts     Options
    runner   *Runner
    registry *tools.Registry
    executor *tools.Executor
}
```

**Key Features:**
- Synchronous and asynchronous execution
- Tool calling with validation
- Memory integration
- RAG-enhanced responses
- Error handling and retry logic
- Streaming support

#### 2. Memory System (`memory/`)

Three types of memory with different retention strategies:

**Conversation Memory**: Simple FIFO with max size
```go
mem := memory.NewConversation(memory.WithMaxSize(100))
```

**Short-term Memory**: Time-based decay with TTL
```go
mem := memory.NewShortTerm(memory.WithTTL(30 * time.Minute))
```

**Long-term Memory**: Importance-based retention
```go
mem := memory.NewLongTerm()
mem.AddWithImportance(ctx, message, 0.9) // High importance
```

#### 3. RAG System (`rag/`)

Complete RAG pipeline with document chunking, embedding, and retrieval:

**Components:**
- **Vector Store**: In-memory storage with similarity search
- **Embedder**: Interface for embedding generation
- **Chunker**: Document splitting with overlap
- **Retriever**: Query processing and top-k search

**Similarity Metrics:**
- Cosine similarity (default)
- Dot product
- Euclidean distance

```go
ragSystem := rag.New(
    rag.WithVectorStore(rag.NewInMemoryVectorStore(
        rag.WithSimilarityMetric(rag.CosineSimilarity),
    )),
    rag.WithChunkSize(500),
    rag.WithChunkOverlap(50),
)
```

#### 4. Tool System (`tools/`)

Extensible tool system with registry and executor:

**Built-in Tools:**
- Calculator: Basic arithmetic operations
- Web Search: Web search (mock implementation)

**Creating Custom Tools:**
```go
type MyTool struct{}

func (t *MyTool) Name() string {
    return "my_tool"
}

func (t *MyTool) Description() string {
    return "Description of what the tool does"
}

func (t *MyTool) Parameters() []tools.Parameter {
    return []tools.Parameter{
        {
            Name:        "input",
            Type:        tools.TypeString,
            Description: "Input parameter",
            Required:    true,
        },
    }
}

func (t *MyTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
    // Tool implementation
    return "result", nil
}
```

#### 5. LLM Integration (`llm/`)

Provider abstraction supporting multiple LLM providers:

```go
type Provider interface {
    Generate(ctx context.Context, messages []Message, opts ...Option) (*Response, error)
    Stream(ctx context.Context, messages []Message, opts ...Option) (<-chan Token, error)
    Name() string
    Model() string
}
```

**Current Implementations:**
- OpenAI (mock implementation for demonstration)
- Extensible for Anthropic, local models, etc.

### Utilities (`internal/utils/`)

Vector math operations:
- Dot product
- Magnitude (L2 norm)
- Cosine similarity
- Euclidean distance
- Vector normalization

## Examples

The `examples/` directory contains three working examples:

1. **basic_agent**: Simple agent with tools
2. **rag_agent**: Agent with RAG capabilities
3. **memory_agent**: Agent with memory management

Run an example:
```bash
cd examples/basic_agent
go run main.go
```

## Testing

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Project Structure

```
gogent/
├── agent/              # Agent core implementation
│   ├── agent.go        # Main agent
│   ├── context.go      # Context management
│   ├── options.go      # Configuration options
│   └── runner.go       # Agent runner
├── memory/             # Memory system
│   ├── memory.go       # Memory interface
│   ├── conversation.go # Conversation memory
│   ├── shortterm.go    # Short-term memory
│   └── longterm.go     # Long-term memory
├── rag/                # RAG system
│   ├── rag.go          # RAG orchestrator
│   ├── vectorstore.go  # Vector store
│   ├── embeddings.go   # Embedding interface
│   ├── retriever.go    # Retrieval logic
│   ├── chunker.go      # Document chunking
│   ├── document.go     # Document types
│   └── similarity.go   # Similarity functions
├── tools/              # Tool system
│   ├── tool.go         # Tool interface
│   ├── registry.go     # Tool registry
│   ├── executor.go     # Tool executor
│   └── builtin/        # Built-in tools
├── llm/                # LLM integration
│   ├── provider.go     # Provider interface
│   ├── message.go      # Message types
│   ├── options.go      # LLM options
│   └── providers/      # Provider implementations
├── internal/utils/     # Internal utilities
├── examples/           # Example applications
└── README.md
```

## Design Principles

1. **Type Safety**: Leverage Go's type system for compile-time safety
2. **Composability**: Modular components that can be composed together
3. **Flexibility**: Support for custom implementations via interfaces
4. **Concurrency**: Goroutines and channels for parallel operations
5. **Error Handling**: Explicit error returns with context
6. **Observability**: Clear interfaces for logging and monitoring

## Dependencies

```go
require (
    golang.org/x/sync v0.6.0  // For errgroup
)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See LICENSE file for details.

## Acknowledgments

Inspired by [Pydantic AI](https://ai.pydantic.dev/), this framework brings similar agentic capabilities to the Go ecosystem.
