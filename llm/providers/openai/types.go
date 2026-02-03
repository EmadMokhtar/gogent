// Package openai provides an OpenAI LLM provider implementation.
package openai

import (
	"github.com/EmadMokhtar/gogent/llm"
	"github.com/EmadMokhtar/gogent/types"
)

// chatMessage represents an OpenAI chat message.
type chatMessage struct {
	Role       string         `json:"role"`
	Content    string         `json:"content,omitempty"`
	Name       string         `json:"name,omitempty"`
	ToolCalls  []toolCall     `json:"tool_calls,omitempty"`
	ToolCallID string         `json:"tool_call_id,omitempty"`
}

// toolCall represents an OpenAI tool call.
type toolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function functionCall `json:"function"`
}

// functionCall represents an OpenAI function call.
type functionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// chatCompletionRequest represents the request to OpenAI chat completions API.
type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
	Stop        []string      `json:"stop,omitempty"`
	Tools       []tool        `json:"tools,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// tool represents an OpenAI tool definition.
type tool struct {
	Type     string       `json:"type"`
	Function functionDef  `json:"function"`
}

// functionDef represents an OpenAI function definition.
type functionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// chatCompletionResponse represents the response from OpenAI chat completions API.
type chatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []choice `json:"choices"`
	Usage   usage    `json:"usage"`
}

// choice represents a completion choice.
type choice struct {
	Index        int         `json:"index"`
	Message      chatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// usage represents token usage information.
type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// chatCompletionChunk represents a streaming chunk from OpenAI.
type chatCompletionChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []deltaChoice `json:"choices"`
}

// deltaChoice represents a streaming choice with delta.
type deltaChoice struct {
	Index        int          `json:"index"`
	Delta        deltaMessage `json:"delta"`
	FinishReason string       `json:"finish_reason,omitempty"`
}

// deltaMessage represents a streaming message delta.
type deltaMessage struct {
	Role      string     `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []toolCall `json:"tool_calls,omitempty"`
}

// errorResponse represents an OpenAI API error response.
type errorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// buildChatCompletionRequest builds a chat completion request.
func buildChatCompletionRequest(messages []types.Message, model string, opts *llm.Options) map[string]interface{} {
	req := map[string]interface{}{
		"model":    model,
		"messages": convertToOpenAIMessages(messages),
	}
	
	if opts.Temperature > 0 {
		req["temperature"] = opts.Temperature
	}
	if opts.MaxTokens > 0 {
		req["max_tokens"] = opts.MaxTokens
	}
	if opts.TopP > 0 && opts.TopP != 1.0 {
		req["top_p"] = opts.TopP
	}
	if len(opts.Stop) > 0 {
		req["stop"] = opts.Stop
	}
	if len(opts.Tools) > 0 {
		req["tools"] = convertToOpenAITools(opts.Tools)
	}
	
	return req
}

// convertToOpenAIMessages converts agent messages to OpenAI format.
func convertToOpenAIMessages(messages []types.Message) []chatMessage {
	result := make([]chatMessage, len(messages))
	for i, msg := range messages {
		result[i] = chatMessage{
			Role:       msg.Role,
			Content:    msg.Content,
			Name:       msg.Name,
			ToolCallID: msg.ToolCallID,
		}
		
		if len(msg.ToolCalls) > 0 {
			result[i].ToolCalls = convertToOpenAIToolCalls(msg.ToolCalls)
		}
	}
	return result
}

// convertToOpenAIToolCalls converts agent tool calls to OpenAI format.
func convertToOpenAIToolCalls(toolCalls []types.ToolCall) []toolCall {
	result := make([]toolCall, len(toolCalls))
	for i, tc := range toolCalls {
		result[i] = toolCall{
			ID:   tc.ID,
			Type: tc.Type,
			Function: functionCall{
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			},
		}
	}
	return result
}

// convertToAgentMessage converts an OpenAI message to agent format.
func convertToAgentMessage(msg chatMessage) types.Message {
	agentMsg := types.Message{
		Role:       msg.Role,
		Content:    msg.Content,
		Name:       msg.Name,
		ToolCallID: msg.ToolCallID,
	}
	
	if len(msg.ToolCalls) > 0 {
		agentMsg.ToolCalls = convertToAgentToolCalls(msg.ToolCalls)
	}
	
	return agentMsg
}

// convertToAgentToolCalls converts OpenAI tool calls to agent format.
func convertToAgentToolCalls(toolCalls []toolCall) []types.ToolCall {
	result := make([]types.ToolCall, len(toolCalls))
	for i, tc := range toolCalls {
		result[i] = types.ToolCall{
			ID:   tc.ID,
			Type: tc.Type,
			Function: types.FunctionCall{
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			},
		}
	}
	return result
}

// convertToOpenAITools converts LLM tools to OpenAI format.
func convertToOpenAITools(tools []llm.ToolDef) []tool {
	result := make([]tool, len(tools))
	for i, t := range tools {
		result[i] = tool{
			Type: t.Type,
			Function: functionDef{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
			},
		}
	}
	return result
}
