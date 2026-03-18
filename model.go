
package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type stage int

const (
	stagePickFileSource stage = iota
	stagePickFileList
	stagePickFilePath
	stagePickOp
	stageConfig
	stagePreview
	stageSave
)

type operation string

const (
	opAddComma operation = "Add comma to end of every line"
	opPrefix   operation = "Add text to beginning of every line"
	opSuffix   operation = "Add text to end of every line"
	opSplit    operation = "Split after space or comma"
)

type model struct {
	// window size
	width		int
	height	int

	stage stage

	// File Source
	fileSourceList list.Model
	fileList list.Model
	selectedFile string

	// File
	fileInput	textinput.Model
	filePath	string
	fileBytes	[]byte
	errMsg		string
	defaultFilesDir string
	fileListEmpty bool

	// pick operation
	opList list.Model
	op     operation

	// config
	prefixInput textinput.Model
	suffixInput textinput.Model
	splitList   list.Model
	splitChoice string

	// saving
	saveList list.Model
	saveChoice string
	outputInput textinput.Model
	statusMsg string
	outputBytes []byte

	quitting bool
}


