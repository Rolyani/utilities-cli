package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
		cursor int
		choices []string
		quitting bool
}

var (
		titleStyle = lipgloss.NewStyle().Bold(true)
		itemStyle = lipgloss.NewStyle().PaddingLeft(2)
		cursorStyle = lipgloss.NewStyle().Bold(true)
)

func initialModel() model {
		return model {
			cursor: 0,
			choices: []string {
				"Add comma to end of every line",
				"Add text to beginning of every line",
				"Add text to end of every line",
				"Create a new line after a space or comma",
				"Quit",
			},
		}
}

func (m model) Init() tea.Cmd {return nil}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
		switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {

					case "ctrl+c", "q":
						m.quitting = true
						return m, tea.Quit

					case "up", "k":
						if m.cursor > 0 {
							m.cursor--
						}

					case "down", "j":
						if m.cursor < len(m.choices) -1 {
							m.cursor++
						}

					case "enter":
						// testing
						m.quitting = true
						return m, tea.Quit


				}
		}
		return m, nil

}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	s := titleStyle.Render("Utilities-cli") + "\n\n"
	s += "Use ↑/↓ (or j/k). Press Enter.\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = cursorStyle.Render(">")
		}
		s += fmt.Sprintf("%s %s\n", cursor, itemStyle.Render(choice))
		
	}

	s += "\nPress q to quit.\n"
	return s

}

func main() {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	//testing
	finalModel, ok := m.(model)
	if ok && finalModel.cursor >= 0 && finalModel.cursor < len(finalModel.choices) {
		fmt.Println("Selected:", finalModel.choices[finalModel.cursor])
	}
}
