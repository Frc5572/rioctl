package ui

import (
	"fmt"
	"rioctl/internal/utils"
	"strings"

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
	progress         progress.Model
	done             bool
	total            int
	finished         int
	currentFile      string
	currentFileBytes int64
	totalFileBytes   int64
	completedFiles   []string
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
	case tea.KeyPressMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.SetWidth(msg.Width - padding*2 - 4)
		if m.progress.Width() > maxWidth {
			m.progress.SetWidth(maxWidth)
		}
		return m, nil
	case int:
		m.finished = msg
		if m.finished >= m.total {
			m.done = true
			return m, tea.Quit
		}
	case utils.FileUpload:
		m.currentFileBytes = msg.Current
		m.totalFileBytes = msg.Size
		m.currentFile = msg.Name
		if m.currentFileBytes >= m.totalFileBytes {
			m.completedFiles = append(m.completedFiles, m.currentFile)
		}
	}

	percent := float64(m.finished) / float64(m.total)
	cmd := m.progress.SetPercent(percent)

	return m, cmd
}

func (m progressModel) View() tea.View {
	pad := strings.Repeat(" ", padding)
	// val := float64(m.finished) / float64(m.total)
	current := float64(m.currentFileBytes) / float64(m.totalFileBytes)
	a := ""
	for i, file := range m.completedFiles {
		a += fmt.Sprintf("%s (%d/%d)\n", file, i+1, m.total)
	}
	if m.done {
		return tea.NewView("\n\n" + a)
	}
	return tea.NewView("\n\n" + a +
		fmt.Sprintf("Downloading %s (%d/%d)", m.currentFile, m.finished+1, m.total) + "\n\n" +
		m.progress.ViewAs(current) + "\n\n" +
		pad + helpStyle(fmt.Sprintf("%s / %s", utils.Humanize(m.currentFileBytes), utils.Humanize((m.totalFileBytes)))))
}

func RunProgress(total int, current <-chan utils.FileUpload, finished <-chan int) {
	m := NewProgress(total)
	p := tea.NewProgram(m)

	go func() {
		for u := range current {
			p.Send(u)
		}
	}()
	go func() {
		for u := range finished {
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
