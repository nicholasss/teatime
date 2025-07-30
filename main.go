package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listStyle  = lipgloss.NewStyle().Margin(1, 2)
	timerStyle = lipgloss.NewStyle().Margin(1, 2)

	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Bold(true)
	paginationStyle = list.DefaultStyles().PaginationStyle.
			PaddingLeft(4)
	helpStyle = list.DefaultStyles().HelpStyle.
			PaddingLeft(4).
			PaddingBottom(1)
)

type teaItem struct {
	title, description string
	timerDuration      int
}

func (ti teaItem) Title() string       { return ti.title }
func (ti teaItem) Description() string { return ti.description }
func (ti teaItem) FilterValue() string { return ti.title }

type model struct {
	list              list.Model
	chosenTeaName     string
	chosenTeaDesc     string
	chosenTeaDuration int
	timer             timer.Model
	quitting          bool
}

// Init
func (m model) Init() tea.Cmd {
	return nil
}

// Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	// case timer.TimeoutMsg:
	// 	m.quitting = true
	// 	return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter", " ":
			// do not start new timer if a tea was already selected
			if m.chosenTeaName != "" {
				return m, nil
			}

			selectedTea, ok := m.list.SelectedItem().(teaItem)
			if !ok {
				log.Printf("not able to find in list -> '%v'", m.list.SelectedItem())
				return m, nil
			}

			m.chosenTeaName = selectedTea.title
			m.chosenTeaDesc = selectedTea.description
			m.chosenTeaDuration = selectedTea.timerDuration
			log.Printf("%q tea selected, timer duration of %d minutes\n", m.chosenTeaName, m.chosenTeaDuration)

			// returning model and starting timer
			teaTimerDuration := time.Minute * time.Duration(m.chosenTeaDuration)
			m.timer = timer.NewWithInterval(teaTimerDuration, time.Second)
			return m, m.timer.Init()
		}

	// used by list on directly from Init
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
	// timer view
	if m.chosenTeaName != "" {
		timerView := timerStyle.Render(m.timer.View())

		if m.timer.Timedout() {
			s := fmt.Sprintf("Your %s tea is done brewing.\n", m.chosenTeaName)
			timerView = timerStyle.Render(s)
		}
		return timerView

	}
	// list view
	return listStyle.Render(m.list.View())
}

// main
func main() {
	teaItems := []list.Item{
		teaItem{title: "Black Tea", description: "100C for 5 minutes", timerDuration: 5},
		teaItem{title: "Green Tea", description: "80C for 3 minutes", timerDuration: 3},
		teaItem{title: "Herbal Tea", description: "85C for 8 minutes", timerDuration: 8},
		teaItem{title: "Oolong Tea", description: "85C for 4 mintues", timerDuration: 4},
	}

	l := list.New(
		teaItems,
		list.NewDefaultDelegate(),
		9, 0,
	)

	l.Title = "Tea Timer Options"
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

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
