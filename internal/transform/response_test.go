package transform

import (
	"testing"

	"github.com/A-NGJ/claude-code-proxy/internal/types"
)

func TestResponse_SimpleText(t *testing.T) {
	resp := types.OpenAIResponse{
		ID:    "chatcmpl-123",
		Model: "qwen2.5-coder",
		Choices: []types.OpenAIChoice{
			{
				Index: 0,
				Message: types.OpenAIMessage{
					Role:    "assistant",
					Content: "Hello, world!",
				},
				FinishReason: "stop",
			},
		},
		Usage: &types.OpenAIUsage{
			PromptTokens:     5,
			CompletionTokens: 3,
		},
	}

	result := Response(resp, "qwen2.5-coder")

	if result.Type != "message" {
		t.Errorf("Expected Type 'message', got '%s'", result.Type)
	}

	if len(result.Content) != 1 {
		t.Fatalf("Expected 1 content block, got %d", len(result.Content))
	}

	if result.Content[0].Text != "Hello, world!" {
		t.Errorf("Expected content text 'Hello, world!', got '%s'", result.Content[0].Text)
	}

	if result.StopReason != "end_turn" {
		t.Errorf("Expected StopReason 'end_turn', got '%s'", result.StopReason)
	}
}

func TestResponse_ToolUse(t *testing.T) {
	resp := types.OpenAIResponse{
		ID:    "chatcmpl-123",
		Model: "qwen2.5-coder",
		Choices: []types.OpenAIChoice{
			{
				Index: 0,
				Message: types.OpenAIMessage{
					Role:    "assistant",
					Content: "",
					ToolCalls: []types.ToolCall{
						{
							ID:   "call_1",
							Type: "function",
							Function: types.FunctionCall{
								Name:      "read_file",
								Arguments: `{"path": "/tmp/test.txt"}`,
							},
						},
					},
				},
				FinishReason: "tool_calls",
			},
		},
	}

	result := Response(resp, "qwen2.5-coder")

	if result.StopReason != "tool_use" {
		t.Errorf("Expected 'tool_use', got %s", result.StopReason)
	}

	if len(result.Content) != 1 {
		t.Fatalf("Expected 1 content block, got %d", len(result.Content))
	}

	if result.Content[0].Type != "tool_use" {
		t.Errorf("Expected 'tool_use' type, got %s", result.Content[0].Type)
	}

	if result.Content[0].Name != "read_file" {
		t.Errorf("Expected 'read_file', got %s", result.Content[0].Name)
	}
}

func TestMapStopReason(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"stop", "end_turn"},
		{"length", "max_tokens"},
		{"tool_calls", "tool_use"},
		{"unknown", "end_turn"},
	}

	for _, tt := range tests {
		result := mapStopReason(tt.input)
		if result != tt.expected {
			t.Errorf("mapStopReason(%s) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}
