package transform

import (
	"encoding/json"
	"testing"

	"github.com/A-NGJ/claude-code-proxy/internal/types"
)

func TestRequest_SimpleText(t *testing.T) {
	req := types.AntrhopicRequest{
		Model:     "claude-3-opus",
		MaxTokens: 1024,
		System:    "You are a helpful assistant.",
		Messages: []types.AnthropicMessage{
			{Role: "user", Content: "Hello"},
		},
	}

	result := Request(req, "qwen2.5-coder")

	if result.Model != "qwen2.5-coder" {
		t.Errorf("Expected model 'qwen2.5-coder', got '%s'", result.Model)
	}

	if len(result.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(result.Messages))
	}

	if result.Messages[0].Role != "system" {
		t.Errorf("Expected system role, got %s", result.Messages[0].Role)
	}

	if result.Messages[1].Content != "Hello" {
		t.Errorf("Expected user content 'Hello', got '%s'", result.Messages[1].Content)
	}
}

func TestRequest_ToolUse(t *testing.T) {
	toolInput := json.RawMessage(`{"path": "/tmp""}`)
	req := types.AntrhopicRequest{
		Model:     "claude-3-opus",
		MaxTokens: 1024,
		Messages: []types.AnthropicMessage{
			{Role: "user", Content: "List files"},
			{
				Role: "assistant",
				Content: []any{
					map[string]any{
						"type":  "tool_use",
						"id":    "tool_1",
						"name":  "ls",
						"input": map[string]any{"path": "/tmp"},
					},
				},
			},
			{
				Role: "user",
				Content: []any{
					map[string]any{
						"type":        "tool_result",
						"tool_use_id": "tool_id",
						"content":     "file1.txt\nfile2.txt",
					},
				},
			},
		},
		Tools: []types.AnthropicTool{
			{Name: "ls", Description: "List files in a directory", InputSchema: toolInput},
		},
	}

	result := Request(req, "")

	// Should have: user, assistant with tool_calls, tool result
	if len(result.Messages) != 3 {
		t.Fatalf("Expected 3 messages, got %d", len(result.Messages))
	}

	// Check tool calls on assistant message
	if len(result.Messages[1].ToolCalls) != 1 {
		t.Errorf("Expected 1 tool call, got %d", len(result.Messages[1].ToolCalls))
	}

	// Check tool result
	if result.Messages[2].Role != "tool" {
		t.Errorf("Expected tool role, got %s", result.Messages[2].Role)
	}
}

func TestConvertToolChoice(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{"auto", "auto"},
		{"any", "required"},
		{"none", "none"},
		{nil, nil},
	}

	for _, tt := range tests {
		result := convertToolChoice(tt.input)
		if result != tt.expected {
			t.Errorf("convertToolChoice(%v) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
