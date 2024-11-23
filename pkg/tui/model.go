package tui

import (
	"fmt"
	"log/slog"
	"tempus/pkg/webhook"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

const width = 45

var (
	modelStyle = lipgloss.NewStyle().
			Width(width).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder()).Render
	focusedModelStyle = lipgloss.NewStyle().
				Width(width).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69")).Render
)

type WebhookHandler interface {
	SendEvent(webhook.WebhookRequest)
}

type focusedState uint

const (
	timerFocused focusedState = iota
	tasklistFocused
)

type model struct {
	timer          timer.Model
	keymap         Keymap
	help           help.Model
	quitting       bool
	timeout        time.Duration
	webhookHandler WebhookHandler
	focusedState   focusedState
	muted          bool
	tasklist       tasklistModel
	insertTask     textinput.Model
}

func NewModel(timeout time.Duration, keymap Keymap, webhookHandler WebhookHandler, muted bool) model {
	keymap.Start.SetEnabled(false)
	return model{
		timer:          timer.NewWithInterval(timeout, time.Second),
		keymap:         keymap,
		help:           help.New(),
		timeout:        timeout,
		webhookHandler: webhookHandler,
		focusedState:   tasklistFocused,
		muted:          muted,
		tasklist:       tasklistModel{tasks: []task{}},
		insertTask:     textinput.New(),
	}
}

func (m model) sendStartStopEvent() {
	if !m.muted {
		var startStop string
		if m.timer.Running() {
			startStop = "Start"
		} else {
			startStop = "Pause"
		}
		webhookReq := webhook.WebhookRequest{Round: "Work Session", Type: startStop}
		m.webhookHandler.SendEvent(webhookReq)
	}
}

func (m model) Init() tea.Cmd {
	m.sendStartStopEvent()
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Info("update", "timer", m.timer.Running(), "remaining", m.timer.Timeout)
	if m.insertTask.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				newTask := task{m.insertTask.Value(), false}
				m.tasklist.tasks = append(m.tasklist.tasks, newTask)
				m.insertTask.Reset()
				m.insertTask.Blur()
				return m, nil
			}
		}

		return m.updateInsertTaskList(msg)
	}
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.Stop.SetEnabled(m.timer.Running())
		m.keymap.Start.SetEnabled(!m.timer.Running())

		m.sendStartStopEvent()
		return m, cmd

	case timer.TimeoutMsg:
		m.quitting = true
		if !m.muted {
			webhookReq := webhook.WebhookRequest{Round: "Work Session", Type: "Complete"}
			m.webhookHandler.SendEvent(webhookReq)
		}
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			if !m.muted {
				webhookReq := webhook.WebhookRequest{Round: "Work Session", Type: "Quit"}
				m.webhookHandler.SendEvent(webhookReq)
			}
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Reset):
			m.timer.Timeout = m.timeout
		case key.Matches(msg, m.keymap.Start, m.keymap.Stop):
			return m, m.timer.Toggle()
		case key.Matches(msg, m.keymap.Focus):
			return m.updateFocusState(), nil
		case key.Matches(msg, m.keymap.Next):
			return m.updateTasklist(NextTaskMsg{})
		case key.Matches(msg, m.keymap.Prev):
			return m.updateTasklist(PrevTaskMsg{})
		case key.Matches(msg, m.keymap.ToggleDone):
			return m.updateTasklist(ToggleDoneMsg{})
		case key.Matches(msg, m.keymap.Insert):
			return m.focusInsertTaskList()
		case key.Matches(msg, m.keymap.Delete):
			return m.updateTasklist(DeleteTaskMsg{})
		}
	}

	return m, nil
}

func (m model) updateInsertTaskList(msg tea.Msg) (model, tea.Cmd) {
	var cmds []tea.Cmd
	insertTask, cmd := m.insertTask.Update(msg)
	cmds = append(cmds, cmd)
	m.insertTask = insertTask
	slog.Info("update task list", "timer", m.timer.Running())
	var timerCmd tea.Cmd
	m.timer, timerCmd = m.timer.Update(msg)
	cmds = append(cmds, timerCmd)
	return m, tea.Batch(cmds...)

}
func (m model) focusInsertTaskList() (model, tea.Cmd) {
	focusCmd := m.insertTask.Focus()
	return m, focusCmd
}

func (m model) helpView() string {
	helpStyle := lipgloss.NewStyle().Width(width).Render
	slog.Info("styles", "style", m.help.Styles.ShortDesc.String())
	return helpStyle("\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.Start,
		m.keymap.Stop,
		m.keymap.Reset,
		m.keymap.Quit,
		m.keymap.Focus,
		m.keymap.Next,
		m.keymap.Prev,
		m.keymap.ToggleDone,
		m.keymap.Insert,
	}))
}

func (m model) updateFocusState() model {
	switch m.focusedState {
	case timerFocused:
		m.focusedState = tasklistFocused
		m.keymap.Focus.SetHelp("tab", "timer")
		m.keymap.Next.SetEnabled(true)
		m.keymap.Prev.SetEnabled(true)
		m.keymap.ToggleDone.SetEnabled(true)
		m.keymap.Insert.SetEnabled(true)
		m.keymap.Start.SetEnabled(false)
		m.keymap.Stop.SetEnabled(false)
		m.keymap.Reset.SetEnabled(false)
	case tasklistFocused:
		m.focusedState = timerFocused
		m.keymap.Focus.SetHelp("tab", "tasklist")
		m.keymap.Next.SetEnabled(false)
		m.keymap.Prev.SetEnabled(false)
		m.keymap.ToggleDone.SetEnabled(false)
		m.keymap.Insert.SetEnabled(false)
		m.keymap.Start.SetEnabled(true)
		m.keymap.Stop.SetEnabled(m.timer.Running())
		m.keymap.Start.SetEnabled(!m.timer.Running())
		m.keymap.Reset.SetEnabled(true)
	}
	return m
}

func (m model) tasklistView() string {
	s := m.tasklist.view()
	if m.focusedState == tasklistFocused {
		return focusedModelStyle(s)
	}
	return modelStyle(s)
}

func (m model) timerView() string {
	s := fmt.Sprintf("Exiting in %v", m.timer.View())
	if m.focusedState == timerFocused {
		return focusedModelStyle(s)
	}
	return modelStyle(s)
}

func (m model) View() string {
	s := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			m.timerView(),
			m.tasklistView(),
		),
		m.helpView(),
	)

	if m.insertTask.Focused() {
		s += "\n" + m.insertTask.View()
	}

	if m.timer.Timedout() {
		s = "All done!"
	}
	// if !m.quitting {
	// 	s += m.helpView()
	// }

	return s
}
