package transform

import (
	"testing"

	"github.com/A-NGJ/claude-code-proxy/internal/types"
)

func TestRequest_SimpleText(t *testing.T) {
	req := types.AntrhopicRequest{
		Model: "claude-3-opus",
		MaxTokens: 1024,
		System: "You are a helpful assistant.",
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
