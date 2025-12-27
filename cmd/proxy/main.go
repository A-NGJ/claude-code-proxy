package main

import (
	"log"
	"os"

	"github.com/A-NGJ/claude-code-proxy/internal/config"
	"github.com/A-NGJ/claude-code-proxy/internal/proxy"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "claude-proxy",
		Usage: "Proxy Claude Code requests to Ollama",
		Flags: config.CLIFlags(),
		Action: func(c *cli.Context) error {
			cfg := config.NewConfigFromCLI(c)
			server := proxy.NewServer(cfg)
			return server.Run()
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
