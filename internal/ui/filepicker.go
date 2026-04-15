package ui

import (
	"fmt"
	"os"
	"rioctl/internal/utils"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var defaultDelegate = list.NewDefaultDelegate()

type item struct {
	title    string
	size     string
	selected bool
}

func (i item) Title() string {
	if i.selected {
		return fmt.Sprintf("[x] %s", i.title)
	}
	return fmt.Sprintf("[ ] %s", i.title)
	// return i.title
}
func (i item) Description() string {
	return "    " + i.size
}

func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	choices  []*item
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) View() tea.View {
	v := tea.NewView(docStyle.Render(m.list.View()))
	v.AltScreen = true
	return v
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			m.quitting = true
			m.choices = []*item{}
			return m, tea.Quit

		case "enter":
			if m.list.FilterState() == list.Filtering {
				// Let the list handle it
				break
			}
			return m, tea.Quit

		case "space":
			f, ok := m.list.SelectedItem().(*item)
			if ok {
				f.selected = !f.selected
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// func (m model) View() string {
// 	if m.quitting {
// 		return ""
// 	}
// 	return m.list.View()
// }

func (m model) Selected() []string {
	var out []string
	for _, c := range m.choices {
		if c.selected {
			out = append(out, c.title)
		}
	}
	return out
}

func NewFilePicker(files []utils.File) model {
	items := make([]list.Item, len(files))
	choices := make([]*item, len(files))

	for i, f := range files {
		it := &item{title: f.Name, size: f.HumanizedSize()}
		items[i] = it
		choices[i] = it
	}

	// defaultDelegate.ShowDescription = false
	l := list.New(items, defaultDelegate, 0, 0)
	l.Title = "Select log files (space to toggle, enter to confirm)"

	return model{
		list:    l,
		choices: choices,
	}
}

func RunFilePicker(files []utils.File) ([]string, error) {
	m := NewFilePicker(files)

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	resultModel := finalModel.(model)
	selected := resultModel.Selected()

	if len(selected) == 0 {
		fmt.Println("No files selected.")
		os.Exit(0)
	}

	return selected, nil
}
