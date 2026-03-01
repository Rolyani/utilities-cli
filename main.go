package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type stage int

const (
	stagePickOp stage = iota
	stageConfig
	stagePreview
)

type operation string

const (
	opAddComma operation = "Add comma to end of every line"
	opPrefix   operation = "Add text to beginning of every line"
	opSuffix   operation = "Add text to end of every line"
	opSplit    operation = "Split after space or comma"
)

type opItem struct {
	title string
	op    operation
}

func (i opItem) Title() string       { return i.title }
func (i opItem) Description() string { return "" }
func (i opItem) FilterValue() string { return i.title }

type model struct {
	width		int
	height	int


	stage stage

	// pick operation
	opList list.Model
	op     operation

	// config
	prefixInput textinput.Model
	suffixInput textinput.Model
	splitList   list.Model
	splitChoice string // "space", "comma", "both"

	quitting bool
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	hintStyle  = lipgloss.NewStyle().Faint(true)
	boxStyle   = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder())
)

func initialModel() model {
	// operation list
	opItems := []list.Item{
		opItem{title: string(opAddComma), op: opAddComma},
		opItem{title: string(opPrefix), op: opPrefix},
		opItem{title: string(opSuffix), op: opSuffix},
		opItem{title: string(opSplit), op: opSplit},
	}
	opL := list.New(opItems, list.NewDefaultDelegate(), 0, 0)
	opL.Title = "Choose an operation"
	opL.SetShowHelp(false)

	// prefix input
	p := textinput.New()
	p.Placeholder = "Enter prefix text…"
	p.Prompt = "> "
	p.CharLimit = 200

	// suffix input
	s := textinput.New()
	s.Placeholder = "Enter suffix text…"
	s.Prompt = "> "
	s.CharLimit = 200

	// split list
	splitItems := []list.Item{
		opItem{title: "Space", op: ""},          // reuse opItem purely as list item
		opItem{title: "Comma", op: ""},          // (op not used here)
		opItem{title: "Space and comma", op: ""}, // (op not used here)
	}
	splitL := list.New(splitItems, list.NewDefaultDelegate(), 0, 0)
	splitL.Title = "Split after which delimiter?"
	splitL.SetShowHelp(false)

	return model{
		stage:       stagePickOp,
		opList:      opL,
		prefixInput: p,
		suffixInput: s,
		splitList:   splitL,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// handle terminal size
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.opList.SetSize(m.width-4, m.height-6)
		m.splitList.SetSize(m.width-4, m.height-6)

		return m, nil
	}

	// global keys
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.stage == stageConfig {
				m.stage = stagePickOp
				m.prefixInput.Blur()
				m.suffixInput.Blur()
				return m, nil
			}
			if m.stage == stagePreview {
				m.stage = stageConfig
				return m, nil
			}
		}
	}

	switch m.stage {
	case stagePickOp:
		var cmd tea.Cmd
		m.opList, cmd = m.opList.Update(msg)

		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
			if it, ok := m.opList.SelectedItem().(opItem); ok {
				m.op = it.op
				m.stage = stageConfig

				// focus correct input
				m.prefixInput.Blur()
				m.suffixInput.Blur()
				if m.op == opPrefix {
					m.prefixInput.Focus()
				} else if m.op == opSuffix {
					m.suffixInput.Focus()
				}
				return m, nil
			}
		}
		return m, cmd

	case stageConfig:
		switch m.op {
		case opPrefix:
			var cmd tea.Cmd
			m.prefixInput, cmd = m.prefixInput.Update(msg)
			if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
				m.stage = stagePreview
				return m, nil
			}
			return m, cmd

		case opSuffix:
			var cmd tea.Cmd
			m.suffixInput, cmd = m.suffixInput.Update(msg)
			if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
				m.stage = stagePreview
				return m, nil
			}
			return m, cmd

		case opSplit:
			var cmd tea.Cmd
			m.splitList, cmd = m.splitList.Update(msg)
			if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
				switch m.splitList.Index() {
				case 0:
					m.splitChoice = "space"
				case 1:
					m.splitChoice = "comma"
				case 2:
					m.splitChoice = "both"
				}
				m.stage = stagePreview
				return m, nil
			}
			return m, cmd

		case opAddComma:
			if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
				m.stage = stagePreview
				return m, nil
			}
			return m, nil
		}

	case stagePreview:
		// placeholder: Enter quits
		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	header := titleStyle.Render("utilities-cli") + "\n" +
		hintStyle.Render("q: quit • esc: back • enter: next") + "\n\n"

	switch m.stage {
	case stagePickOp:
		return header + boxStyle.Render(m.opList.View())

	case stageConfig:
		switch m.op {
		case opPrefix:
			return header + boxStyle.Render("Prefix text:\n\n"+m.prefixInput.View())
		case opSuffix:
			return header + boxStyle.Render("Suffix text:\n\n"+m.suffixInput.View())
		case opSplit:
			return header + boxStyle.Render(m.splitList.View())
		case opAddComma:
			return header + boxStyle.Render("No configuration needed.\n\nPress Enter to continue.")
		default:
			return header + boxStyle.Render("Unknown operation.\n\nPress esc to go back.")
		}

	case stagePreview:
		// placeholder preview screen
		config := "No extra config"
		switch m.op {
		case opPrefix:
			config = fmt.Sprintf("Prefix: %q", m.prefixInput.Value())
		case opSuffix:
			config = fmt.Sprintf("Suffix: %q", m.suffixInput.Value())
		case opSplit:
			config = fmt.Sprintf("Split: %s", m.splitChoice)
		}

		body := "Preview (placeholder)\n\n" +
			fmt.Sprintf("Operation: %s\n%s\n\n", m.op, config) +
			"Next stage will show input/output sample.\n\n" +
			"Press Enter to quit (for now)."

		return header + boxStyle.Render(body)
	}

	return header + "Unknown state."
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
