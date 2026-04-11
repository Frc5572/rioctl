package ui

import (
	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	yellow    = lipgloss.Color("#FDFF8C")
	pink      = lipgloss.Color("#FF7CCB")
)

type progressModel struct {
	progress progress.Model
	total    int
	current  int
	done     bool
}

func NewProgress(total int) progressModel {
	return progressModel{
		progress: progress.New(progress.WithScaled(true)),
		total:    total,
	}
}

type tickMsg struct{}

func (m progressModel) Init() tea.Cmd {
	return nil
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case int:
		m.current = msg
		if m.current >= m.total {
			m.done = true
			return m, tea.Quit
		}
	}

	percent := float64(m.current) / float64(m.total)
	cmd := m.progress.SetPercent(percent)

	return m, cmd
}

func (m progressModel) View() tea.View {
	// pad := strings.Repeat(" ", padding)
	val := float64(m.current) / float64(m.total)
	return tea.NewView(m.progress.ViewAs(val))
}

func RunProgress(total int, updates <-chan int) {
	m := NewProgress(total)
	p := tea.NewProgram(m)

	go func() {
		for u := range updates {
			p.Send(u)
		}
	}()

	p.Run()
}
func RunProgress1(total int, stuff []string) {
	m := NewProgress(total)
	p := tea.NewProgram(m)

	go func() {
		for u := range stuff {
			p.Send(u)
		}
	}()

	p.Run()
}
