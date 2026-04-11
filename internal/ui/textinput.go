package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	input textinput.Model
	done  bool
	value string
}

func NewInput(prompt string) inputModel {
	ti := textinput.New()
	ti.Placeholder = prompt
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return inputModel{input: ti}
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.value = m.input.Value()
			m.done = true
			return m, tea.Quit
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	return fmt.Sprintf("Enter remote directory:\n\n%s\n\n(enter to confirm)", m.input.View())
}

func RunTextInput(prompt string) (string, error) {
	m := NewInput(prompt)
	p := tea.NewProgram(m)

	final, err := p.Run()
	if err != nil {
		return "", err
	}

	return final.(inputModel).value, nil
}
