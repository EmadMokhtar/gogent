package builtin

import (
	"context"
	"testing"
)

func TestEcho_Execute(t *testing.T) {
	echo := NewEcho()
	ctx := context.Background()
	
	tests := []struct {
		name    string
		input   map[string]interface{}
		want    string
		wantErr bool
	}{
		{
			name: "echo string message",
			input: map[string]interface{}{
				"message": "Hello, World!",
			},
			want:    "Hello, World!",
			wantErr: false,
		},
		{
			name: "missing message parameter",
			input: map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "non-string message",
			input: map[string]interface{}{
				"message": 123,
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := echo.Execute(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Echo.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				result, ok := got.(string)
				if !ok {
					t.Errorf("Echo.Execute() result is not string")
					return
				}
				if result != tt.want {
					t.Errorf("Echo.Execute() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestEcho_Name(t *testing.T) {
	echo := NewEcho()
	if echo.Name() != "echo" {
		t.Errorf("Echo.Name() = %v, want %v", echo.Name(), "echo")
	}
}

func TestEcho_Description(t *testing.T) {
	echo := NewEcho()
	desc := echo.Description()
	if desc == "" {
		t.Error("Echo.Description() should not be empty")
	}
}

func TestEcho_Parameters(t *testing.T) {
	echo := NewEcho()
	params := echo.Parameters()
	if len(params) != 1 {
		t.Errorf("Echo.Parameters() returned %d parameters, want 1", len(params))
	}
	if params[0].Name != "message" {
		t.Errorf("Echo.Parameters()[0].Name = %s, want message", params[0].Name)
	}
}
