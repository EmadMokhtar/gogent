# Gogent

Gogent is a high-performance, multi-agent framework for Go, inspired by Pydantic AI. It provides a powerful and flexible foundation for building AI agents with tool execution capabilities, LLM integration, and conversation management.

## Features

- 🤖 **Agent System**: Core agent abstraction with context management and conversation history
- 🔧 **Tool System**: Flexible tool registration and execution with input validation
- 🌐 **LLM Integration**: Abstract provider interface with OpenAI implementation
- 🔄 **Multi-turn Conversations**: Automatic handling of tool calls and iterative execution
- 🎯 **Type-Safe**: Leverage Go's type system for safe and reliable code
- 📦 **Composable**: Modular components that can be easily composed
- ⚙️ **Configurable**: Functional options pattern for flexible configuration

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
    "os"

    "github.com/EmadMokhtar/gogent/agent"
    "github.com/EmadMokhtar/gogent/llm/providers/openai"
)

func main() {
    ctx := context.Background()

    // Create LLM provider
    llmProvider := openai.New(
        openai.WithModel("gpt-4"),
        openai.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
    )

    // Create agent
    ag := agent.New(
        agent.WithLLM(llmProvider),
        agent.WithSystemPrompt("You are a helpful AI assistant."),
    )

    // Run agent
    result, err := ag.Run(ctx, "Hello! Can you help me?")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response:", result.Content)
}
```

### Agent with Tools

```go
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

    // Create LLM provider
    llmProvider := openai.New(
        openai.WithModel("gpt-4"),
        openai.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
    )

    // Create agent with tools
    ag := agent.New(
        agent.WithLLM(llmProvider),
        agent.WithTools(
            builtin.NewCalculator(),
            builtin.NewEcho(),
        ),
        agent.WithSystemPrompt("You are a helpful AI assistant with access to tools."),
    )

    // Run agent with a calculation
    result, err := ag.Run(ctx, "What is 42 multiplied by 17?")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response:", result.Content)
    
    // Continue conversation
    result2, err := ag.Run(ctx, "Now add 100 to that result")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Follow-up:", result2.Content)
}
```

## Architecture

### Core Components

#### Agent
The `Agent` is the main component that:
- Manages conversation history
- Orchestrates LLM calls
- Executes tools when needed
- Handles multi-turn interactions

```go
agent := agent.New(
    agent.WithLLM(llmProvider),
    agent.WithTools(tool1, tool2),
    agent.WithSystemPrompt("You are a helpful assistant"),
    agent.WithMaxIterations(10),
    agent.WithTemperature(0.7),
)
```

#### Tools
Tools extend agent capabilities by allowing them to perform specific actions:

```go
type Tool interface {
    Name() string
    Description() string
    Parameters() []Parameter
    Execute(ctx context.Context, input map[string]interface{}) (interface{}, error)
}
```

Built-in tools:
- **Calculator**: Performs arithmetic operations (add, subtract, multiply, divide)
- **Echo**: Echoes back messages (useful for testing)

#### LLM Provider
Abstract interface for LLM integration:

```go
type Provider interface {
    Generate(ctx context.Context, messages []Message, opts ...Option) (*Response, error)
    Stream(ctx context.Context, messages []Message, opts ...Option) (<-chan Token, error)
}
```

Currently supported providers:
- **OpenAI**: Full support for chat completions and tool calling

### Directory Structure

```
gogent/
├── agent/              # Core agent implementation
│   ├── agent.go
│   └── agent_test.go
├── tools/              # Tool system
│   ├── tool.go
│   ├── registry.go
│   ├── tool_test.go
│   └── builtin/        # Built-in tools
│       ├── calculator.go
│       ├── calculator_test.go
│       └── echo.go
├── llm/                # LLM provider abstraction
│   ├── provider.go
│   ├── types.go
│   ├── options.go
│   └── providers/      # Provider implementations
│       └── openai/
│           ├── openai.go
│           └── types.go
└── examples/           # Usage examples
    ├── basic_agent/
    └── tool_agent/
```

## Core Concepts

### Messages
Messages represent the conversation between the user, assistant, and tools:

```go
type Message struct {
    Role      string                 // "user", "assistant", "system", "tool"
    Content   string
    ToolCalls []ToolCall
    Metadata  map[string]interface{}
}
```

### Tool Execution
The agent automatically handles tool execution in a loop:

1. User sends a query
2. LLM decides if tools are needed
3. Agent executes requested tools
4. Results are sent back to LLM
5. LLM generates final response
6. Process repeats if more tool calls are needed (up to max iterations)

### Configuration
All components use the functional options pattern for configuration:

```go
// Agent options
agent.WithLLM(provider)
agent.WithTools(tool1, tool2)
agent.WithSystemPrompt("prompt")
agent.WithMaxIterations(10)
agent.WithTemperature(0.7)
agent.WithMaxTokens(1000)

// OpenAI provider options
openai.WithModel("gpt-4")
openai.WithAPIKey("key")
openai.WithBaseURL("https://api.openai.com/v1")

// LLM call options
llm.WithTemperature(0.7)
llm.WithMaxTokens(1000)
llm.WithTopP(0.9)
llm.WithStop([]string{"STOP"})
```

## Creating Custom Tools

Implement the `Tool` interface:

```go
package mytools

import (
    "context"
    "fmt"
    
    "github.com/EmadMokhtar/gogent/tools"
)

type MyTool struct{}

func NewMyTool() *MyTool {
    return &MyTool{}
}

func (t *MyTool) Name() string {
    return "my_tool"
}

func (t *MyTool) Description() string {
    return "Description of what my tool does"
}

func (t *MyTool) Parameters() []tools.Parameter {
    return []tools.Parameter{
        {
            Name:        "param1",
            Description: "Description of param1",
            Type:        "string",
            Required:    true,
        },
    }
}

func (t *MyTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
    param1, ok := input["param1"].(string)
    if !ok {
        return nil, fmt.Errorf("param1 must be a string")
    }
    
    // Do something with param1
    result := "Processed: " + param1
    
    return result, nil
}
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

## Examples

Check out the [examples](./examples) directory for more usage examples:

- [Basic Agent](./examples/basic_agent) - Simple agent without tools
- [Tool Agent](./examples/tool_agent) - Agent with calculator and echo tools

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Guidelines

1. Follow Go idioms and best practices
2. Add tests for new functionality
3. Update documentation as needed
4. Ensure all tests pass before submitting

## Roadmap

- [x] Core agent implementation
- [x] Tool system with registry
- [x] OpenAI provider
- [x] Built-in tools (Calculator, Echo)
- [ ] Memory system (Part 2)
- [ ] RAG system (Part 3)
- [ ] Additional LLM providers (Anthropic, etc.)
- [ ] Advanced streaming support
- [ ] Observability and tracing

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Inspired by [Pydantic AI](https://github.com/pydantic/pydantic-ai) - a Python framework for building production-grade applications with LLMs.
