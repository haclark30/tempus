package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"tempus/pkg/config"
	"tempus/pkg/tui"
	"tempus/pkg/webhook"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("not enough args")
	}

	cfg, err := config.ReadConfig()
	if err != nil && errors.Is(err, os.ErrNotExist) {
		log.Fatal(fmt.Sprintf("no config file found: %v", err))
	}

	sleepTime, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(fmt.Sprintf("not a number: %v", err))
	}

	timeout := time.Minute * time.Duration(sleepTime)

	keymap := tui.Keymap{
		Start: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "start"),
		),
		Stop: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "stop"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}

	webhookHandler := webhook.HttpWebhookHandler{
		Client:     http.DefaultClient,
		WebhookUrl: cfg.WebhookUrl,
	}
	m := tui.NewModel(timeout, keymap, webhookHandler)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal(fmt.Sprintf("error running program: %v", err))
	}
}
