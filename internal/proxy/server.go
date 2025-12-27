package proxy

import (
	"log"
	"net/http"

	"github.com/A-NGJ/claude-code-proxy/internal/config"
)

type Server struct {
	config *config.Config
	client *http.Client
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		client: &http.Client{},
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/messages", s.handleMessages)
	mux.HandleFunc("/v1/messages/count_tokens", s.handleCountTokens)

	log.Printf("Starting proxy on %s -> %s (model: %s)",
		s.config.ListenAddr, s.config.OllamaURL, s.config.Model)
	return http.ListenAndServe(s.config.ListenAddr, mux)
}
