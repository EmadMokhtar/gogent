package tools

import (
	"context"
	"testing"

	"github.com/EmadMokhtar/gogent/llm"
)

// mockTool is a simple mock tool for testing
type mockTool struct {
	name string
}

func (m *mockTool) Name() string {
	return m.name
}

func (m *mockTool) Description() string {
	return "A mock tool for testing"
}

func (m *mockTool) Parameters() []Parameter {
	return []Parameter{
		{Name: "input", Description: "Test input", Type: "string", Required: true},
	}
}

func (m *mockTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	return "mock result", nil
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	tool := &mockTool{name: "test_tool"}
	
	err := registry.Register(tool)
	if err != nil {
		t.Errorf("Registry.Register() error = %v", err)
	}
	
	// Try to register the same tool again
	err = registry.Register(tool)
	if err == nil {
		t.Error("Registry.Register() should return error when registering duplicate tool")
	}
	
	// Try to register nil tool
	err = registry.Register(nil)
	if err == nil {
		t.Error("Registry.Register() should return error when registering nil tool")
	}
	
	// Try to register tool with empty name
	emptyNameTool := &mockTool{name: ""}
	err = registry.Register(emptyNameTool)
	if err == nil {
		t.Error("Registry.Register() should return error when registering tool with empty name")
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()
	tool := &mockTool{name: "test_tool"}
	
	registry.Register(tool)
	
	retrieved, ok := registry.Get("test_tool")
	if !ok {
		t.Error("Registry.Get() should return true for registered tool")
	}
	if retrieved.Name() != "test_tool" {
		t.Errorf("Registry.Get() returned wrong tool, got %v", retrieved.Name())
	}
	
	_, ok = registry.Get("nonexistent")
	if ok {
		t.Error("Registry.Get() should return false for nonexistent tool")
	}
}

func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()
	tool1 := &mockTool{name: "tool1"}
	tool2 := &mockTool{name: "tool2"}
	
	registry.Register(tool1)
	registry.Register(tool2)
	
	tools := registry.List()
	if len(tools) != 2 {
		t.Errorf("Registry.List() returned %d tools, want 2", len(tools))
	}
}

func TestRegistry_Execute(t *testing.T) {
	registry := NewRegistry()
	tool := &mockTool{name: "test_tool"}
	registry.Register(tool)
	
	ctx := context.Background()
	
	// Test successful execution
	result, err := registry.Execute(ctx, llm.ToolCall{
		ID:    "1",
		Name:  "test_tool",
		Input: map[string]interface{}{"input": "test"},
	})
	
	if err != nil {
		t.Errorf("Registry.Execute() error = %v", err)
	}
	if result != "mock result" {
		t.Errorf("Registry.Execute() = %v, want %v", result, "mock result")
	}
	
	// Test execution with nonexistent tool
	_, err = registry.Execute(ctx, llm.ToolCall{
		ID:    "2",
		Name:  "nonexistent",
		Input: map[string]interface{}{},
	})
	
	if err == nil {
		t.Error("Registry.Execute() should return error for nonexistent tool")
	}
}

func TestRegistry_ToLLMTools(t *testing.T) {
	registry := NewRegistry()
	tool := &mockTool{name: "test_tool"}
	registry.Register(tool)
	
	tools := registry.ToLLMTools()
	if len(tools) != 1 {
		t.Errorf("Registry.ToLLMTools() returned %d tools, want 1", len(tools))
	}
	
	if tools[0].Type != "function" {
		t.Errorf("Registry.ToLLMTools() tool type = %s, want function", tools[0].Type)
	}
	
	if tools[0].Function.Name != "test_tool" {
		t.Errorf("Registry.ToLLMTools() tool name = %s, want test_tool", tools[0].Function.Name)
	}
}

