# Gogent

Gogent is a high-performance, multi-agent framework for Go, inspired by Pydantic AI. It provides a clean, extensible architecture for building AI agents with tool execution capabilities.

## Features

- 🤖 **Simple Agent API**: Easy-to-use agent interface for LLM interactions
- 🔧 **Tool System**: Extensible tool registry with built-in tools (calculator, echo)
- 🔄 **Automatic Tool Calling**: Agents automatically detect and execute tool calls
- 🌊 **Streaming Support**: Real-time streaming responses from LLMs
- 🔌 **Provider Abstraction**: Clean LLM provider interface (OpenAI included)
- ⚡ **Parallel Tool Execution**: Execute multiple tools concurrently
- 🎯 **Type-Safe**: Strongly typed Go implementation with proper error handling

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
)

func main() {
    ctx := context.Background()

    // Create LLM provider
    llm := openai.New(
        openai.WithModel("gpt-4"),
        openai.WithAPIKey("your-api-key"),
    )

    // Create agent
    ag := agent.New(
        agent.WithLLM(llm),
        agent.WithSystemPrompt("You are a helpful AI assistant."),
    )

    // Run agent
    response, err := ag.Run(ctx, "What is the capital of France?")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(response.Content)
}
```

### Agent with Tools

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

    llm := openai.New(
        openai.WithModel("gpt-4"),
        openai.WithAPIKey("your-api-key"),
    )

    ag := agent.New(
        agent.WithLLM(llm),
        agent.WithTools(
            builtin.NewCalculator(),
        ),
        agent.WithSystemPrompt("You are a helpful assistant with access to a calculator."),
    )

    response, err := ag.Run(ctx, "What is 157 * 23 + 891?")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(response.Content)
}
```

### Streaming Responses

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/EmadMokhtar/gogent/agent"
    "github.com/EmadMokhtar/gogent/llm/providers/openai"
)

func main() {
    ctx := context.Background()

    llm := openai.New(
        openai.WithModel("gpt-4"),
        openai.WithAPIKey("your-api-key"),
    )

    ag := agent.New(
        agent.WithLLM(llm),
        agent.WithSystemPrompt("You are a helpful AI assistant."),
    )

    tokenChan, err := ag.Stream(ctx, "Write a short poem.")
    if err != nil {
        log.Fatal(err)
    }

    for token := range tokenChan {
        if token.Error != nil {
            log.Fatal(token.Error)
        }
        fmt.Print(token.Delta)
    }
}
```

## Architecture

### Core Components

#### Agent
The `Agent` is the main orchestrator that:
- Manages conversation context
- Calls the LLM provider
- Detects and executes tool calls
- Handles iteration loops for multi-step reasoning

#### Tools
Tools are functions that agents can call. The tool system includes:
- **Tool Interface**: Standard interface for all tools
- **Tool Registry**: Manages registered tools
- **Tool Executor**: Executes tools with validation
- **Built-in Tools**: Calculator, Echo, and more

#### LLM Providers
Provider abstraction allows integration with different LLM services:
- **OpenAI**: Full support for GPT models with function calling
- Extensible for other providers (Anthropic, local models, etc.)

### Agent Workflow

1. User provides a query
2. Agent prepends system prompt and builds message history
3. Agent calls LLM provider with available tools
4. If LLM requests tool calls:
   - Tools are executed in parallel
   - Results are added to message history
   - Agent calls LLM again with results
   - Process repeats (up to max iterations)
5. Final response is returned

## Configuration

### Agent Options

```go
agent.New(
    agent.WithLLM(provider),              // Set LLM provider
    agent.WithTools(tool1, tool2),        // Register tools
    agent.WithSystemPrompt("..."),        // Set system prompt
    agent.WithMaxIterations(10),          // Max tool call iterations
    agent.WithTimeout(2 * time.Minute),   // Execution timeout
)
```

### LLM Options

```go
// Per-call options
response, err := agent.Run(ctx, query, 
    llm.WithTemperature(0.8),
    llm.WithMaxTokens(2000),
    llm.WithTopP(0.95),
)
```

## Creating Custom Tools

Implement the `Tool` interface:

```go
type MyTool struct{}

func (t *MyTool) Name() string {
    return "my_tool"
}

func (t *MyTool) Description() string {
    return "Does something useful"
}

func (t *MyTool) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "param1": map[string]interface{}{
                "type": "string",
                "description": "First parameter",
            },
        },
        "required": []string{"param1"},
    }
}

func (t *MyTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
    // Tool implementation
    return result, nil
}
```

## Examples

See the [examples](./examples) directory for complete examples:
- [Basic Agent](./examples/basic_agent/main.go)
- [Tool Agent](./examples/tool_agent/main.go)
- [Streaming Agent](./examples/streaming_agent/main.go)

## Testing

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [x] Core agent framework
- [x] Tool system with built-in tools
- [x] OpenAI provider with streaming
- [ ] Memory system (Part 2)
- [ ] RAG support (Part 3)
- [ ] Additional LLM providers
- [ ] More built-in tools
- [ ] Agent orchestration
- [ ] Multi-agent collaboration
