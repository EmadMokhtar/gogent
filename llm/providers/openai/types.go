package openai

// chatCompletionRequest represents a request to the OpenAI chat completion API
type chatCompletionRequest struct {
	Model       string                   `json:"model"`
	Messages    []chatMessage            `json:"messages"`
	Temperature *float64                 `json:"temperature,omitempty"`
	MaxTokens   *int                     `json:"max_tokens,omitempty"`
	TopP        *float64                 `json:"top_p,omitempty"`
	Stop        []string                 `json:"stop,omitempty"`
	Tools       []map[string]interface{} `json:"tools,omitempty"`
	Stream      bool                     `json:"stream,omitempty"`
}

// chatMessage represents a message in the chat
type chatMessage struct {
	Role      string                   `json:"role"`
	Content   string                   `json:"content"`
	ToolCalls []chatToolCall           `json:"tool_calls,omitempty"`
	ToolCallID string                  `json:"tool_call_id,omitempty"`
}

// chatToolCall represents a tool call in a message
type chatToolCall struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"`
	Function chatFunctionCall    `json:"function"`
}

// chatFunctionCall represents a function call
type chatFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// chatCompletionResponse represents a response from the OpenAI chat completion API
type chatCompletionResponse struct {
	ID      string               `json:"id"`
	Object  string               `json:"object"`
	Created int64                `json:"created"`
	Model   string               `json:"model"`
	Choices []chatChoice         `json:"choices"`
	Usage   usageInfo            `json:"usage"`
}

// chatChoice represents a choice in the response
type chatChoice struct {
	Index        int         `json:"index"`
	Message      chatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// usageInfo represents token usage information
type usageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// errorResponse represents an error response from OpenAI
type errorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}
