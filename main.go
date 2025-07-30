package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listStyle  = lipgloss.NewStyle().Margin(1, 2)
	timerStyle = lipgloss.NewStyle().Margin(1, 18)

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
	timerProgress     progress.Model
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
		var timerView string

		if m.timer.Timedout() {
			doneText := fmt.Sprintf("Your %s tea is done brewing.\n", m.chosenTeaName)
			return timerStyle.Render(doneText)
		}

		timerView = timerStyle.Render(m.timer.View())
		progressView := m.timerProgress.View()
		return timerView + "\n " + progressView

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

	list := list.New(
		teaItems,
		list.NewDefaultDelegate(),
		9, 0,
	)

	list.Title = "Tea Timer Options"
	list.SetShowStatusBar(false)
	list.Styles.Title = titleStyle
	list.Styles.PaginationStyle = paginationStyle
	list.Styles.HelpStyle = helpStyle

	progress := progress.New(
		progress.WithSolidFill("10"),
		progress.WithoutPercentage(),
	)

	model := model{
		list:          list,
		timerProgress: progress,
	}

	// debug log
	if len(os.Getenv("DEBUG")) > 0 {
		file, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Println("fatal:", err)
			os.Exit(1)
		}
		defer file.Close()
	}

	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatalln("fatal:", err)
		os.Exit(1)
	}
}
