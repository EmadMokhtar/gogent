package tools

import (
	"context"
	"testing"
)

func TestRegistry(t *testing.T) {
	t.Run("register tool", func(t *testing.T) {
		registry := NewRegistry()
		tool := &mockTool{name: "test_tool"}
		
		err := registry.Register(tool)
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
		
		if registry.Count() != 1 {
			t.Errorf("Count() = %d, want 1", registry.Count())
		}
	})
	
	t.Run("register duplicate tool", func(t *testing.T) {
		registry := NewRegistry()
		tool := &mockTool{name: "test_tool"}
		
		_ = registry.Register(tool)
		err := registry.Register(tool)
		
		if err == nil {
			t.Error("Register() expected error for duplicate tool")
		}
	})
	
	t.Run("get tool", func(t *testing.T) {
		registry := NewRegistry()
		tool := &mockTool{name: "test_tool"}
		
		_ = registry.Register(tool)
		
		got, err := registry.Get("test_tool")
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}
		
		if got.Name() != "test_tool" {
			t.Errorf("Get() name = %s, want test_tool", got.Name())
		}
	})
	
	t.Run("get non-existent tool", func(t *testing.T) {
		registry := NewRegistry()
		
		_, err := registry.Get("non_existent")
		if err == nil {
			t.Error("Get() expected error for non-existent tool")
		}
	})
	
	t.Run("unregister tool", func(t *testing.T) {
		registry := NewRegistry()
		tool := &mockTool{name: "test_tool"}
		
		_ = registry.Register(tool)
		err := registry.Unregister("test_tool")
		
		if err != nil {
			t.Errorf("Unregister() error = %v", err)
		}
		
		if registry.Count() != 0 {
			t.Errorf("Count() = %d, want 0", registry.Count())
		}
	})
}

// mockTool is a mock implementation of the Tool interface for testing
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
	return []Parameter{}
}

func (m *mockTool) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	return "mock result", nil
}
