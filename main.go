package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/joho/godotenv/autoload"
)

type WebhookRequest struct {
	Round string `json:"round"`
	Type  string `json:"type"`
}

func main() {

	webhookUrl := os.Getenv("WEBHOOK_URL")
	args := os.Args[1:]
	if len(args) < 1 {
		slog.Error("not enough args")
		os.Exit(1)
	}

	sleepTime, err := strconv.Atoi(args[0])
	if err != nil {
		slog.Error("not a number")
		os.Exit(1)
	}

	timeout := time.Minute * time.Duration(sleepTime)

	m := model{
		timer: timer.NewWithInterval(timeout, time.Second),
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		help:       help.New(),
		timeout:    timeout,
		webhookUrl: webhookUrl,
	}

	m.keymap.start.SetEnabled(false)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}
