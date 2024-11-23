package tui

import "github.com/charmbracelet/bubbles/key"

type Keymap struct {
	Start      key.Binding
	Stop       key.Binding
	Reset      key.Binding
	Quit       key.Binding
	Focus      key.Binding
	Next       key.Binding
	Prev       key.Binding
	ToggleDone key.Binding
	Insert     key.Binding
	Delete     key.Binding
}
