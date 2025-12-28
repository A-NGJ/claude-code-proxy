package transform

import (
	"encoding/json"

	"github.com/A-NGJ/claude-code-proxy/internal/types"
)

func Request(req types.AntrhopicRequest, modelOverride string) types.OpenAIRequest {
	var messages []types.OpenAIMessage

	// Handle system prompt
	if req.System != nil {
		systemText := extractSystemText(req.System)
		if systemText != "" {
			messages = append(messages, types.OpenAIMessage{
				Role:    "system",
				Content: systemText,
			})
		}
	}

	// Convert messages
	for _, msg := range req.Messages {
		converted := convertMessage(msg)
		messages = append(messages, converted...)
	}

	// Determine model
	model := modelOverride
	if model == "" {
		model = req.Model
	}

	return types.OpenAIRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Stream:      req.Stream,
		Tools:       convertTools(req.Tools),
		ToolChoice:  convertToolChoice(req.ToolChoice),
		Stop:        req.StopSequences,
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}
}

func extractSystemText(system any) string {
	switch v := system.(type) {
	case string:
		return v
	case []interface{}:
		var text string
		for _, block := range v {
			if m, ok := block.(map[string]interface{}); ok {
				if t, ok := m["text"].(string); ok {
					text += t + "\n"
				}
			}
		}
		return text
	}
	return ""
}

func convertMessage(msg types.AnthropicMessage) []types.OpenAIMessage {
	var result []types.OpenAIMessage

	switch content := msg.Content.(type) {
	case string:
		result = append(result, types.OpenAIMessage{
			Role:    msg.Role,
			Content: content,
		})
	case []interface{}:
		// Handle array of content blocks
		var textContent string
		var toolCalls []types.ToolCall
		var toolResults []types.OpenAIMessage

		for _, block := range content {
			if m, ok := block.(map[string]interface{}); ok {
				blockType, _ := m["type"].(string)

				switch blockType {
				case "text":
					if t, ok := m["text"].(string); ok {
						textContent += t
					}
				case "tool_use":
					id, _ := m["id"].(string)
					name, _ := m["name"].(string)
					input, _ := json.Marshal(m["input"])
					toolCalls = append(toolCalls, types.ToolCall{
						ID:   id,
						Type: "function",
						Function: types.FunctionCall{
							Name:      name,
							Arguments: string(input),
						},
					})
				case "tool_result":
					toolUseID, _ := m["tool_use_id"].(string)
					resultContent := extractToolResultContent(m["content"])
					toolResults = append(toolResults, types.OpenAIMessage{
						Role:       "tool",
						Content:    resultContent,
						ToolCallID: toolUseID,
					})
				}
			}
		}

		// Add assistant message with text and/or tool calls
		if textContent != "" || len(toolCalls) > 0 {
			result = append(result, types.OpenAIMessage{
				Role:      msg.Role,
				Content:   textContent,
				ToolCalls: toolCalls,
			})
		}

		result = append(result, toolResults...)
	}

	return result
}

func extractToolResultContent(content any) string {
	switch v := content.(type) {
	case string:
		return v
	case []interface{}:
		var text string
		for _, block := range v {
			if m, ok := block.(map[string]interface{}); ok {
				if t, ok := m["text"].(string); ok {
					text += t
				}
			}
		}
		return text
	}
	return ""
}

func convertTools(tools []types.AnthropicTool) []types.OpenAITool {
	if len(tools) == 0 {
		return nil
	}

	result := make([]types.OpenAITool, len(tools))
	for i, tool := range tools {
		result[i] = types.OpenAITool{
			Type: "function",
			Function: types.OpenAIFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.InputSchema,
			},
		}
	}
	return result
}

func convertToolChoice(choice any) any {
	if choice == nil {
		return nil
	}

	switch v := choice.(type) {
	case string:
		// "auto", "any", "none" -> OpenAI equivalents
		switch v {
		case "any":
			return "required"
		case "none":
			return "none"
		default:
			return "auto"
		}
	case map[string]interface{}:
		// Specific tool choice: {"type": "tool", "name": "..."}
		if name, ok := v["name"].(string); ok {
			return map[string]interface{}{
				"type": "function",
				"function": map[string]string{
					"name": name,
				},
			}
		}
	}

	return "auto"
}
