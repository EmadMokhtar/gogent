package builtin

import (
	"context"
	"testing"
)

func TestCalculator_Execute(t *testing.T) {
	calc := NewCalculator()
	ctx := context.Background()
	
	tests := []struct {
		name    string
		input   map[string]interface{}
		want    float64
		wantErr bool
	}{
		{
			name: "add two numbers",
			input: map[string]interface{}{
				"operation": "add",
				"a":         float64(5),
				"b":         float64(3),
			},
			want:    8,
			wantErr: false,
		},
		{
			name: "subtract two numbers",
			input: map[string]interface{}{
				"operation": "subtract",
				"a":         float64(10),
				"b":         float64(3),
			},
			want:    7,
			wantErr: false,
		},
		{
			name: "multiply two numbers",
			input: map[string]interface{}{
				"operation": "multiply",
				"a":         float64(4),
				"b":         float64(5),
			},
			want:    20,
			wantErr: false,
		},
		{
			name: "divide two numbers",
			input: map[string]interface{}{
				"operation": "divide",
				"a":         float64(20),
				"b":         float64(4),
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "divide by zero",
			input: map[string]interface{}{
				"operation": "divide",
				"a":         float64(10),
				"b":         float64(0),
			},
			wantErr: true,
		},
		{
			name: "unknown operation",
			input: map[string]interface{}{
				"operation": "modulo",
				"a":         float64(10),
				"b":         float64(3),
			},
			wantErr: true,
		},
		{
			name: "missing parameter a",
			input: map[string]interface{}{
				"operation": "add",
				"b":         float64(3),
			},
			wantErr: true,
		},
		{
			name: "missing parameter b",
			input: map[string]interface{}{
				"operation": "add",
				"a":         float64(3),
			},
			wantErr: true,
		},
		{
			name: "missing operation",
			input: map[string]interface{}{
				"a": float64(3),
				"b": float64(5),
			},
			wantErr: true,
		},
		{
			name: "int parameters",
			input: map[string]interface{}{
				"operation": "multiply",
				"a":         42,
				"b":         17,
			},
			want:    714,
			wantErr: false,
		},
		{
			name: "int64 parameters",
			input: map[string]interface{}{
				"operation": "add",
				"a":         int64(100),
				"b":         int64(200),
			},
			want:    300,
			wantErr: false,
		},
		{
			name: "invalid type for a",
			input: map[string]interface{}{
				"operation": "add",
				"a":         "invalid",
				"b":         float64(3),
			},
			wantErr: true,
		},
		{
			name: "invalid type for b",
			input: map[string]interface{}{
				"operation": "add",
				"a":         float64(3),
				"b":         "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid operation type",
			input: map[string]interface{}{
				"operation": 123,
				"a":         float64(3),
				"b":         float64(5),
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calc.Execute(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculator.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				result, ok := got.(float64)
				if !ok {
					t.Errorf("Calculator.Execute() result is not float64")
					return
				}
				if result != tt.want {
					t.Errorf("Calculator.Execute() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestCalculator_Name(t *testing.T) {
	calc := NewCalculator()
	if calc.Name() != "calculator" {
		t.Errorf("Calculator.Name() = %v, want %v", calc.Name(), "calculator")
	}
}

func TestCalculator_Description(t *testing.T) {
	calc := NewCalculator()
	desc := calc.Description()
	if desc == "" {
		t.Error("Calculator.Description() should not be empty")
	}
}

func TestCalculator_Parameters(t *testing.T) {
	calc := NewCalculator()
	params := calc.Parameters()
	if len(params) != 3 {
		t.Errorf("Calculator.Parameters() returned %d parameters, want 3", len(params))
	}
}
