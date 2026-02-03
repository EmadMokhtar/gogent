// Package utils provides internal utility functions for the Gogent framework.
package utils

import (
	"fmt"
	"reflect"
)

// ValidateRequired checks if required fields in a map are present and non-empty.
func ValidateRequired(input map[string]interface{}, required []string) error {
	for _, field := range required {
		value, exists := input[field]
		if !exists {
			return fmt.Errorf("required field '%s' is missing", field)
		}
		
		// Check if value is empty/nil
		if value == nil {
			return fmt.Errorf("required field '%s' cannot be nil", field)
		}
		
		// Check for empty strings
		if str, ok := value.(string); ok && str == "" {
			return fmt.Errorf("required field '%s' cannot be empty", field)
		}
	}
	
	return nil
}

// ValidateType checks if a field has the expected type.
func ValidateType(input map[string]interface{}, field string, expectedType string) error {
	value, exists := input[field]
	if !exists {
		return nil // Field doesn't exist, skip validation
	}
	
	actualType := reflect.TypeOf(value).Kind().String()
	
	// Handle numeric types
	if expectedType == "number" {
		switch reflect.TypeOf(value).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return nil
		default:
			return fmt.Errorf("field '%s' must be a number, got %s", field, actualType)
		}
	}
	
	// Handle boolean
	if expectedType == "boolean" {
		if reflect.TypeOf(value).Kind() != reflect.Bool {
			return fmt.Errorf("field '%s' must be a boolean, got %s", field, actualType)
		}
		return nil
	}
	
	// Handle string
	if expectedType == "string" {
		if reflect.TypeOf(value).Kind() != reflect.String {
			return fmt.Errorf("field '%s' must be a string, got %s", field, actualType)
		}
		return nil
	}
	
	// Handle array
	if expectedType == "array" {
		kind := reflect.TypeOf(value).Kind()
		if kind != reflect.Slice && kind != reflect.Array {
			return fmt.Errorf("field '%s' must be an array, got %s", field, actualType)
		}
		return nil
	}
	
	// Handle object
	if expectedType == "object" {
		if reflect.TypeOf(value).Kind() != reflect.Map {
			return fmt.Errorf("field '%s' must be an object, got %s", field, actualType)
		}
		return nil
	}
	
	return nil
}

// CoerceToFloat64 attempts to convert a value to float64.
func CoerceToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("cannot coerce %T to float64", value)
	}
}
