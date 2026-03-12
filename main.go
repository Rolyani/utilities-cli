package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type opItem struct {
	title string
	op    operation
}

func (i opItem) Title() string       { return i.title }
func (i opItem) Description() string { return "" }
func (i opItem) FilterValue() string { return i.title }

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
