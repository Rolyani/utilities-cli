package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type fileItem struct {
	title string
	path string
}

type opItem struct {
	title string
	op    operation
}

func (i fileItem) Title() string { return i.title }
func (i fileItem) Description() string { return i.path }
func (i fileItem) FilterValue() string { return i.title }

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

	// file source
	sourceItems := []list.Item {
		opItem {title: "Choose from files folder", op: ""},
		opItem {title: "Enter a custom file path", op: ""},
	}
	sourceList := list.New(sourceItems, list.NewDefaultDelegate(), 0, 0)
	sourceList.Title = "How would you like to choose a file?"
	sourceList.SetShowHelp(false)

	fileL := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	fileL.Title = "Choose a file"
	fileL.SetShowHelp(false)

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
		stage:       stagePickFileSource,
		fileSourceList: sourceList,
		fileList: fileL,
		fileInput: 	 fi,
		opList:      opL,
		prefixInput: p,
		suffixInput: s,
		splitList:   splitL,
		saveList:		 saveL,
		outputInput:	oi,
		defaultFilesDir: "files",
	}
}

func (m model) Init() tea.Cmd { return nil }


