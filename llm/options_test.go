package llm

import "testing"

func TestWithTemperature(t *testing.T) {
	temp := 0.7
	opts := ApplyOptions(WithTemperature(temp))
	
	if opts.Temperature == nil {
		t.Error("WithTemperature() did not set temperature")
		return
	}
	if *opts.Temperature != temp {
		t.Errorf("WithTemperature() = %v, want %v", *opts.Temperature, temp)
	}
}

func TestWithMaxTokens(t *testing.T) {
	maxTokens := 1000
	opts := ApplyOptions(WithMaxTokens(maxTokens))
	
	if opts.MaxTokens == nil {
		t.Error("WithMaxTokens() did not set max tokens")
		return
	}
	if *opts.MaxTokens != maxTokens {
		t.Errorf("WithMaxTokens() = %v, want %v", *opts.MaxTokens, maxTokens)
	}
}

func TestWithTopP(t *testing.T) {
	topP := 0.9
	opts := ApplyOptions(WithTopP(topP))
	
	if opts.TopP == nil {
		t.Error("WithTopP() did not set top-p")
		return
	}
	if *opts.TopP != topP {
		t.Errorf("WithTopP() = %v, want %v", *opts.TopP, topP)
	}
}

func TestWithStop(t *testing.T) {
	stop := []string{"STOP", "END"}
	opts := ApplyOptions(WithStop(stop))
	
	if len(opts.Stop) != 2 {
		t.Errorf("WithStop() stop length = %d, want 2", len(opts.Stop))
	}
}

func TestWithTools(t *testing.T) {
	tools := []ToolDefinition{
		{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "test",
				Description: "test function",
			},
		},
	}
	opts := ApplyOptions(WithTools(tools))
	
	if len(opts.Tools) != 1 {
		t.Errorf("WithTools() tools length = %d, want 1", len(opts.Tools))
	}
	if opts.Tools[0].Function.Name != "test" {
		t.Errorf("WithTools() tool name = %s, want test", opts.Tools[0].Function.Name)
	}
}

func TestApplyOptions_Multiple(t *testing.T) {
	temp := 0.7
	maxTokens := 1000
	
	opts := ApplyOptions(
		WithTemperature(temp),
		WithMaxTokens(maxTokens),
	)
	
	if opts.Temperature == nil || *opts.Temperature != temp {
		t.Error("ApplyOptions() did not apply temperature")
	}
	if opts.MaxTokens == nil || *opts.MaxTokens != maxTokens {
		t.Error("ApplyOptions() did not apply max tokens")
	}
}
