// Package utils provides internal utility functions for the Gogent framework.
package utils

import (
	"testing"
)

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		required []string
		wantErr  bool
	}{
		{
			name: "all required fields present",
			input: map[string]interface{}{
				"field1": "value1",
				"field2": "value2",
			},
			required: []string{"field1", "field2"},
			wantErr:  false,
		},
		{
			name: "missing required field",
			input: map[string]interface{}{
				"field1": "value1",
			},
			required: []string{"field1", "field2"},
			wantErr:  true,
		},
		{
			name: "nil value",
			input: map[string]interface{}{
				"field1": nil,
			},
			required: []string{"field1"},
			wantErr:  true,
		},
		{
			name: "empty string",
			input: map[string]interface{}{
				"field1": "",
			},
			required: []string{"field1"},
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequired(tt.input, tt.required)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateType(t *testing.T) {
	tests := []struct {
		name         string
		input        map[string]interface{}
		field        string
		expectedType string
		wantErr      bool
	}{
		{
			name:         "valid string",
			input:        map[string]interface{}{"field": "value"},
			field:        "field",
			expectedType: "string",
			wantErr:      false,
		},
		{
			name:         "invalid string",
			input:        map[string]interface{}{"field": 123},
			field:        "field",
			expectedType: "string",
			wantErr:      true,
		},
		{
			name:         "valid number (int)",
			input:        map[string]interface{}{"field": 123},
			field:        "field",
			expectedType: "number",
			wantErr:      false,
		},
		{
			name:         "valid number (float)",
			input:        map[string]interface{}{"field": 123.45},
			field:        "field",
			expectedType: "number",
			wantErr:      false,
		},
		{
			name:         "invalid number",
			input:        map[string]interface{}{"field": "not a number"},
			field:        "field",
			expectedType: "number",
			wantErr:      true,
		},
		{
			name:         "valid boolean",
			input:        map[string]interface{}{"field": true},
			field:        "field",
			expectedType: "boolean",
			wantErr:      false,
		},
		{
			name:         "invalid boolean",
			input:        map[string]interface{}{"field": "true"},
			field:        "field",
			expectedType: "boolean",
			wantErr:      true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateType(tt.input, tt.field, tt.expectedType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCoerceToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    float64
		wantErr bool
	}{
		{"float64", float64(123.45), 123.45, false},
		{"float32", float32(123.45), float64(float32(123.45)), false},
		{"int", int(123), 123.0, false},
		{"int64", int64(123), 123.0, false},
		{"uint", uint(123), 123.0, false},
		{"string", "123", 0, true},
		{"bool", true, 0, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CoerceToFloat64(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoerceToFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("CoerceToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}
