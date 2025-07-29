package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var listStyle = lipgloss.NewStyle().Margin(1, 2)

type teaItem struct {
	title, description string
	timerDuration      int
}

func (ti teaItem) Title() string       { return ti.title }
func (ti teaItem) Description() string { return ti.description }
func (ti teaItem) FilterValue() string { return ti.title }

type model struct {
	list list.Model
}

// Init
func (m model) Init() tea.Cmd {
	return nil
}

// Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View
func (m model) View() string {
	return listStyle.Render(m.list.View())
}

// main
func main() {
	teaItems := []list.Item{
		teaItem{title: "Black Tea", description: "highest caffeine", timerDuration: 5},
		teaItem{title: "Green Tea", description: "very delicate", timerDuration: 2},
		teaItem{title: "Fruit Tea", description: "many flavors", timerDuration: 8},
	}

	m := model{list: list.New(teaItems, list.NewDefaultDelegate(), 0, 0)}

	// debug log
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln("fatal:", err)
		os.Exit(1)
	}
}
