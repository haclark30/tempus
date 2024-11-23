package tui

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var selectedTask = lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true)

type task struct {
	text string
	done bool
}
type tasklistModel struct {
	tasks    []task
	selected int
}

type NextTaskMsg struct{}
type PrevTaskMsg struct{}
type ToggleDoneMsg struct{}
type DeleteTaskMsg struct{}

func (m model) updateTasklist(msg tea.Msg) (model, tea.Cmd) {
	switch msg.(type) {
	case NextTaskMsg:
		m.tasklist.selected++
		if int(m.tasklist.selected) >= len(m.tasklist.tasks) {
			m.tasklist.selected = 0
		}
	case PrevTaskMsg:
		m.tasklist.selected--
		if int(m.tasklist.selected) < 0 {
			m.tasklist.selected = len(m.tasklist.tasks) - 1
		}
	case ToggleDoneMsg:
		if len(m.tasklist.tasks) > 0 {
			newDone := !m.tasklist.tasks[m.tasklist.selected].done
			m.tasklist.tasks[m.tasklist.selected].done = newDone
			slog.Info("toggle task", "task", m.tasklist.tasks[m.tasklist.selected])
		}
	case DeleteTaskMsg:
		if len(m.tasklist.tasks) > 0 {
			idx := m.tasklist.selected
			m.tasklist.tasks = append(m.tasklist.tasks[:idx], m.tasklist.tasks[idx+1:]...)
			m.tasklist.selected--
			if int(m.tasklist.selected) < 0 {
				m.tasklist.selected = len(m.tasklist.tasks) - 1
			}
		}
	}
	return m, nil
}

func (t tasklistModel) view() string {
	var s string
	for i, task := range t.tasks {
		style := lipgloss.NewStyle()
		if task.done {
			style = style.Strikethrough(true)
		}
		if t.selected == int(i) {
			style = style.Inherit(selectedTask)
		}
		s += style.Render(task.text) + "\n"
	}
	return s
}
