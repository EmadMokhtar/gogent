// Package tools provides the tool system for agent tool execution.
package tools

import (
	"context"
	"testing"
)

// mockTool is a mock tool for testing.
type mockTool struct {
	name string
}

func (m *mockTool) Name() string {
	return m.name
}

func (m *mockTool) Description() string {
	return "A mock tool for testing"
}

func (m *mockTool) Parameters() map[string]interface{} {
	return map[string]interface{}{}
}

func (m *mockTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	return "mock result", nil
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	tool := &mockTool{name: "test"}
	
	err := registry.Register(tool)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !registry.Has("test") {
		t.Error("tool should be registered")
	}
}

func TestRegistry_Register_Duplicate(t *testing.T) {
	registry := NewRegistry()
	tool := &mockTool{name: "test"}
	
	registry.Register(tool)
	err := registry.Register(tool)
	
	if err == nil {
		t.Error("expected error for duplicate registration")
	}
}

func TestRegistry_Register_Nil(t *testing.T) {
	registry := NewRegistry()
	
	err := registry.Register(nil)
	if err == nil {
		t.Error("expected error for nil tool")
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()
	tool := &mockTool{name: "test"}
	registry.Register(tool)
	
	retrieved, err := registry.Get("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if retrieved.Name() != "test" {
		t.Errorf("expected tool name 'test', got %s", retrieved.Name())
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	registry := NewRegistry()
	
	_, err := registry.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent tool")
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
		t.Errorf("expected 2 tools, got %d", len(tools))
	}
}

func TestRegistry_Unregister(t *testing.T) {
	registry := NewRegistry()
	tool := &mockTool{name: "test"}
	registry.Register(tool)
	
	err := registry.Unregister("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if registry.Has("test") {
		t.Error("tool should be unregistered")
	}
}

func TestRegistry_Unregister_NotFound(t *testing.T) {
	registry := NewRegistry()
	
	err := registry.Unregister("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent tool")
	}
}
