package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// import tea "github.com/charmbracelet/bubbletea"
// bubbletea is always imported as tea.
// tea was named after 'The Elm Architecture'

/*
 * Model is the type that gets used for the Init, Update, and View functions.
 * It should store the apps state and can be or contain any types.
 */
type model struct {
	// status int
	// err error
}

/*
 * initialModel should return the model struct when the app starts.
 * This can also be replaced with model{} in the call to tea.NewProgram()
 */
func initialModel() model {
	return model{}
}

/*
 * The Init method can return a tea.Cmd to perform initial I/O.
 * Or, it can just return nil, which means no initial I/O.
 */
func (m model) Init() tea.Cmd {
	return nil
}

/*
 * The Update method is a very important one.
 * Update will decide, based on messages passed in, what to do.
 *
 * Any update to the application is a tea.Msg,
 * and is typically interpreted through a switch on its type.
 *
 * model should always be returned, after its been updated.
 */
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Because tea.Msg is of type Any, it needs type assertions.
	// These assertions let you decide how to perform actions,
	// based on its type or value.
	//
	// msg.(type) is special syntax only within switch statements.
	switch msg := msg.(type) {

	// This case checks if the parameter is 'tea.KeyMsg'
	case tea.KeyMsg:
		return m, nil

	// This case checks if the parameter is 'MyCustomMsg'
	case MyCustomMsg:
		return m, nil

	// This case checks if the parameter is an error
	case error:
		return m, nil

	} // You can have a default case, for any remaining types

	// And for any cases that did not return, they can reach the return at the bottom
	return m, nil
}

/*
 * View decides how to render the UI of the application.
 * Basic apps can just return a formatted string.
 *
 * The model is checked for its state,
 * and based on it, will decide what to show to the user.
 *
 * This is typically where instructions are provided to the user.
 */
func (m model) View() string {
	return ""
}

// Main is the last necessary function.
// It will actually start the program and can even intialize state.
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("there was an error: %v\n", err)
		os.Exit(1)
	}
}
