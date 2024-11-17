package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		slog.Error("not enough args")
		os.Exit(1)
	}

	cfg, err := readConfig()
	if err != nil && errors.Is(err, os.ErrNotExist) {
		log.Fatal(fmt.Sprintf("no config file found: %v", err))
	}

	sleepTime, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(fmt.Sprintf("not a number: %v", err))
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
		help:    help.New(),
		timeout: timeout,
		webhookHandler: HttpWebhookHandler{
			client:     http.DefaultClient,
			webhookUrl: cfg.WebhookUrl,
		},
	}

	m.keymap.start.SetEnabled(false)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}
