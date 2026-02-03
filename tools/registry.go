package tools

import (
	"context"
	"fmt"
	"sync"

	"github.com/EmadMokhtar/gogent/llm"
)

// Registry manages available tools
type Registry struct {
	tools map[string]Tool
	mu    sync.RWMutex
}

// NewRegistry creates a new tool registry
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry
func (r *Registry) Register(tool Tool) error {
	if tool == nil {
		return fmt.Errorf("cannot register nil tool")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	name := tool.Name()
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %s already registered", name)
	}
	
	r.tools[name] = tool
	return nil
}

// Get retrieves a tool by name
func (r *Registry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	tool, ok := r.tools[name]
	return tool, ok
}

// List returns all registered tools
func (r *Registry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// Execute executes a tool call
func (r *Registry) Execute(ctx context.Context, call llm.ToolCall) (interface{}, error) {
	tool, ok := r.Get(call.Name)
	if !ok {
		return nil, &ToolError{
			Tool:    call.Name,
			Message: "tool not found",
		}
	}
	
	result, err := tool.Execute(ctx, call.Input)
	if err != nil {
		return nil, &ToolError{
			Tool:    call.Name,
			Message: "execution failed",
			Err:     err,
		}
	}
	
	return result, nil
}

// ToLLMTools converts registered tools to LLM tool definitions
func (r *Registry) ToLLMTools() []llm.ToolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	definitions := make([]llm.ToolDefinition, 0, len(r.tools))
	for _, tool := range r.tools {
		params := tool.Parameters()
		properties := make(map[string]interface{})
		required := []string{}
		
		for _, param := range params {
			properties[param.Name] = map[string]interface{}{
				"type":        param.Type,
				"description": param.Description,
			}
			if param.Required {
				required = append(required, param.Name)
			}
		}
		
		definition := llm.ToolDefinition{
			Type: "function",
			Function: llm.FunctionDefinition{
				Name:        tool.Name(),
				Description: tool.Description(),
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": properties,
					"required":   required,
				},
			},
		}
		definitions = append(definitions, definition)
	}
	
	return definitions
}
