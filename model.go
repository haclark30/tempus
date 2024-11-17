package main

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"

	tea "github.com/charmbracelet/bubbletea"
)

type WebhookHandler interface {
	SendEvent(WebhookRequest)
}

type model struct {
	timer          timer.Model
	keymap         keymap
	help           help.Model
	quitting       bool
	timeout        time.Duration
	webhookHandler WebhookHandler
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

func (m model) sendStartStopEvent() {
	var startStop string
	if m.timer.Running() {
		startStop = "Start"
	} else {
		startStop = "Pause"
	}
	webhookReq := WebhookRequest{"Work Session", startStop}
	m.webhookHandler.SendEvent(webhookReq)
}

func (m model) Init() tea.Cmd {
	m.sendStartStopEvent()
	return m.timer.Init()
	// return m.timer.Stop()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())

		m.sendStartStopEvent()
		return m, cmd

	case timer.TimeoutMsg:
		m.quitting = true
		webhookReq := WebhookRequest{"Work Session", "Complete"}
		m.webhookHandler.SendEvent(webhookReq)
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			webhookReq := WebhookRequest{"Work Session", "Quit"}
			m.webhookHandler.SendEvent(webhookReq)
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.timer.Timeout = m.timeout
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			return m, m.timer.Toggle()
		}
	}

	return m, nil
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func (m model) View() string {
	// For a more detailed timer view you could read m.timer.Timeout to get
	// the remaining time as a time.Duration and skip calling m.timer.View()
	// entirely.
	s := m.timer.View()

	if m.timer.Timedout() {
		s = "All done!"
	}
	s += "\n"
	if !m.quitting {
		s = "Exiting in " + s
		s += m.helpView()
	}
	return s
}
