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
			key.WithDisabled(),
		),
		Stop: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "stop"),
			key.WithDisabled(),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset"),
			key.WithDisabled(),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Focus: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch to timer"),
		),
		Next: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j", "next"),
		),
		Prev: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k", "prev"),
		),
		ToggleDone: key.NewBinding(
			key.WithKeys(" ", "t"),
			key.WithHelp("space", "toggle task"),
		),
		Insert: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "insert task"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
		),
	}

	webhookHandler := webhook.HttpWebhookHandler{
		Client:     http.DefaultClient,
		WebhookUrl: cfg.WebhookUrl,
	}
	m := tui.NewModel(timeout, keymap, webhookHandler, cfg.Muted)
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal(fmt.Sprintf("error running program: %v", err))
	}
}
