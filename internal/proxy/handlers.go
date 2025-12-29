package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/A-NGJ/claude-code-proxy/internal/transform"
	"github.com/A-NGJ/claude-code-proxy/internal/types"
)

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.AntrhopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Transform request
	openaiReq := transform.Request(req, s.config.Model)

	if req.Stream {
		s.handleStreamingRequest(w, openaiReq)
	} else {
		s.handleNonStreamingRequest(w, openaiReq)
	}
}

func (s *Server) handleCountTokens(w http.ResponseWriter, r *http.Request) {
	// Stub implementation - return plausible fake numbers
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"input_tokens":  100,
	})
}

func (s *Server) handleNonStreamingRequest(w http.ResponseWriter, req types.OpenAIRequest) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	ollamaURL := s.config.OllamaURL + "/v1/chat/completions"
	resp, err := s.client.Post(ollamaURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		http.Error(w, "Failed to reach Ollama: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Ollama error: %s", body), resp.StatusCode)
		return
	}

	var openaiResp types.OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		http.Error(w, "Failed to decode Ollama response", http.StatusInternalServerError)
		return
	}

	anthropicResp := transform.Response(openaiResp, req.Model)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anthropicResp)
}

func (s *Server) handleStreamingRequest(w http.ResponseWriter, req types.OpenAIRequest) {
	// TODO : Implement in pahse 5
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
