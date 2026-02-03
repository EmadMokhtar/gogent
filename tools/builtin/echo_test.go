// Package builtin provides built-in tools for the Gogent framework.
package builtin

import (
	"context"
	"testing"
)

func TestEcho_Name(t *testing.T) {
	echo := NewEcho()
	if echo.Name() != "echo" {
		t.Errorf("expected name 'echo', got %s", echo.Name())
	}
}

func TestEcho_Description(t *testing.T) {
	echo := NewEcho()
	desc := echo.Description()
	if desc == "" {
		t.Error("description should not be empty")
	}
}

func TestEcho_Parameters(t *testing.T) {
	echo := NewEcho()
	params := echo.Parameters()
	if params == nil {
		t.Error("parameters should not be nil")
	}
}

func TestEcho_Execute(t *testing.T) {
	echo := NewEcho()
	ctx := context.Background()
	
	result, err := echo.Execute(ctx, map[string]interface{}{
		"message": "Hello, World!",
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	resultMap, ok := result.(map[string]string)
	if !ok {
		t.Fatalf("expected map[string]string, got %T", result)
	}
	
	if resultMap["echoed"] != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got %s", resultMap["echoed"])
	}
}

func TestEcho_Execute_InvalidInput(t *testing.T) {
	echo := NewEcho()
	ctx := context.Background()
	
	_, err := echo.Execute(ctx, map[string]interface{}{
		"message": 123, // Not a string
	})
	
	if err == nil {
		t.Error("expected error for invalid input")
	}
}
