package config

import "github.com/urfave/cli/v2"

type Config struct {
	ListenAddr string
	OllamaURL  string
	Model      string
}

func NewConfigFromCLI(c *cli.Context) *Config {
	return &Config{
		ListenAddr: c.String("listen"),
		OllamaURL:  c.String("ollama-url"),
		Model:      c.String("model"),
	}
}

func CLIFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "listen",
			Aliases: []string{"l"},
			Usage:   "Address to listen on",
			Value:   ":3456",
			EnvVars: []string{"LISTEN_ADDR"},
		},
		&cli.StringFlag{
			Name:    "ollama-url",
			Aliases: []string{"o"},
			Usage:   "Ollama API base URL",
			Value:   "http://localhost:11434",
			EnvVars: []string{"OLLAMA_URL"},
		},
		&cli.StringFlag{
			Name:    "model",
			Aliases: []string{"m"},
			Usage:   "Model to use for all requests",
			Value:   "qwen2.5-coder:latest",
			EnvVars: []string{"MODEL"},
		},
	}
}
