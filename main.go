package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"Rolyani/utilities-cli/transform"
)

type opItem struct {
	title string
	op    operation
}

func (i opItem) Title() string       { return i.title }
func (i opItem) Description() string { return "" }
func (i opItem) FilterValue() string { return i.title }

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
	// file input
	fi	:= textinput.New()
	fi.Placeholder	= "/path/to/file.txt"
	fi.Prompt		= "> "
	fi.CharLimit = 500
	fi.Focus()

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
		opItem{title: "Space", op: ""},
		opItem{title: "Comma", op: ""},
		opItem{title: "Space and comma", op: ""},
	}
	splitL := list.New(splitItems, list.NewDefaultDelegate(), 0, 0)
	splitL.Title = "Split after which delimiter?"
	splitL.SetShowHelp(false)

	// save 
	saveItems := []list.Item {
		opItem{title: "Save as a new file", op: ""},
		opItem{title: "Overwrite original file", op: ""},
	}
	saveL := list.New(saveItems, list.NewDefaultDelegate(), 0, 0)
	saveL.Title = "How do you want to save the output?"
	saveL.SetShowHelp(false)

	oi := textinput.New()
	oi.Placeholder = "output file path"
	oi.Prompt = "> "
	oi.CharLimit = 500

	return model{
		stage:       stagePickFile,
		fileInput: 	 fi,
		opList:      opL,
		prefixInput: p,
		suffixInput: s,
		splitList:   splitL,
		saveList:		 saveL,
		outputInput:	oi,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) View() string {
	if m.quitting {
		return ""
	}

	header := titleStyle.Render("utilities-cli") + "\n" +
		hintStyle.Render("ctrl + c: quit • esc: back • enter: next") + "\n\n"

	switch m.stage {
	case stagePickFile:
		body := "Enter the path of the file to edit.\n\n" + m.fileInput.View()
		if m.errMsg != "" {
			body += "\n\n" + m.errMsg
		}
		return header + boxStyle.Render(body)

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
			return header + boxStyle.Render("No further configuration needed.\n\nPress Enter to continue.")
		default:
			return header + boxStyle.Render("Unknown operation.\n\nPress esc to go back.")
		}

	case stagePreview:
		out := m.outputBytes

		left := "INPUT (first 10 lines)\n\n" + transform.FirstNLines(m.fileBytes, 10)
		right := "OUTPUT (first 10 lines)\n\n" + transform.FirstNLines(out, 10)

		colW := (m.width -6) /2
		leftBox := lipgloss.NewStyle().
			Width(colW).
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Render(left)

		rightBox := lipgloss.NewStyle().
			Width(colW).
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Render(right)

		content := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

		footer := "\n\nPress Enter to continue. Press esc to go back. Press ctrl + c to quit."
		return header + content + footer

	case stageSave:
		if m.saveChoice == "" {
			return header + boxStyle.Render(m.saveList.View())
		}

		if m.saveChoice == "new" {
			body := "Enter output file path:\n\n" + m.outputInput.View()
			if m.statusMsg != "" {
				body += "\n\n" + m.statusMsg
			}
			body += "\n\nPress enter to save."
			return header + boxStyle.Render(body)
		}

		if m.saveChoice == "overwrite" {
			body := "overwrite original file:\n\n" + m.filePath + "\n\nA .bak backup will be created first. \n\nPress enter to confirm."
			if m.statusMsg != "" {
				body += "\n\n" + m.statusMsg
			}
			return header + boxStyle.Render(body)

		}


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

func readFile(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("The path is a directory")
	}
	return os.ReadFile(path)
}

func writeNewFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func backupFile(path string) (string, error) {
	original, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	backupPath := path + ".bak"
	err = os.WriteFile(backupPath, original, 0644)
	if err != nil {
		return "", err
	}

	return backupPath, nil
}

func overwriteFile(path string, data []byte) error {
	dir := "."
	if idx := strings.LastIndex(path, string(os.PathSeparator)); idx != -1 {
		dir = path[:idx]
		if dir == "" {
			dir = string(os.PathSeparator)
		}
	}

	tmp, err := os.CreateTemp(dir, "utilities-cli-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	defer func() {
		tmp.Close()
		os.Remove(tmpName)
	}()

	if _, err := tmp.Write(data); err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, path)

}
