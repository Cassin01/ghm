package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/Cassin01/ghm/pkg/config"
	"github.com/urfave/cli/v2"
)

func TestRootCommand(t *testing.T) {
	cfg := &config.Config{
		Root:            "/test/root",
		DefaultProtocol: "https",
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "root",
				Action: func(c *cli.Context) error {
					return rootCommand(c, cfg)
				},
			},
		},
	}

	err := app.Run([]string{"ghm", "root"})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	if err != nil {
		t.Errorf("rootCommand() error = %v", err)
	}

	output := buf.String()
	expected := "/test/root\n"
	if output != expected {
		t.Errorf("rootCommand() output = %q, want %q", output, expected)
	}
}
