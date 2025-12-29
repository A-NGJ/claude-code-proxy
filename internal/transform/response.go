package transform

import (
	"encoding/json"

	"github.com/A-NGJ/claude-code-proxy/internal/types"
	"github.com/google/uuid"
)

func Response(resp types.OpenAIResponse, requestModel string) types.AnthropicResponse {
	var content []types.ContentBlock

	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]

		// Add text content if present
		if choice.Message.Content != "" {
			content = append(content, types.ContentBlock{
				Type: "text",
				Text: choice.Message.Content,
			})
		}

		// Add tool use blocks
		for _, tc := range choice.Message.ToolCalls {
			content = append(content, types.ContentBlock{
				Type:  "tool_use",
				ID:    tc.ID,
				Name:  tc.Function.Name,
				Input: json.RawMessage(tc.Function.Arguments),
			})
		}

	}

	// Ensure at least empty text block
	if len(content) == 0 {
		content = append(content, types.ContentBlock{Type: "text", Text: ""})
	}

	stopReason := "end_turn"
	if len(resp.Choices) > 0 {
		stopReason = mapStopReason(resp.Choices[0].FinishReason)
	}

	usage := types.Usage{InputTokens: 0, OutputTokens: 0}
	if resp.Usage != nil {
		usage.InputTokens = resp.Usage.PromptTokens
		usage.OutputTokens = resp.Usage.CompletionTokens
	}

	return types.AnthropicResponse{
		ID:         "msg_" + uuid.New().String(),
		Type:       "message",
		Role:       "assistant",
		Content:    content,
		Model:      requestModel,
		StopReason: stopReason,
		Usage:      usage,
	}
}

func mapStopReason(reason string) string {
	switch reason {
	case "stop":
		return "end_turn"
	case "length":
		return "max_tokens"
	case "tool_calls":
		return "tool_use"
	default:
		return "end_turn"
	}
}
