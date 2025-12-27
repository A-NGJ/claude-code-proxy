package proxy

import (
	"encoding/json"
	"net/http"

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
	// openaiReq := transform.Request(req, s.config.Model)
	//
	// if req.Stream {
	// 	s.handleStreamingRequest(w, openaiReq)
	// } else {
	// 	s.handleNonStreamingRequest(w, openaiReq)
	// }
}

func (s *Server) handleCountTokens(w http.ResponseWriter, r *http.Request) {
	// Stub implementation - return plausible fake numbers
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"input_tokens":  100,
	})
}

func (s *Server) handleNonStreamingRequest(w http.ResponseWriter, req types.OpenAIRequest) {
	// TODO : Implement in pahse 4
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) handleStreamingRequest(w http.ResponseWriter, req types.OpenAIRequest) {
	// TODO : Implement in pahse 5
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
