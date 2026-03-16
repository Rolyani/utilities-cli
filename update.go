package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"Rolyani/utilities-cli/transform"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// handle terminal size
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.fileSourceList.SetSize(m.width-8, m.height-8)
		m.fileList.SetSize(m.width-8, m.height-8)
		m.opList.SetSize(m.width-8, m.height-8)
		m.splitList.SetSize(m.width-8, m.height-8)
		m.saveList.SetSize(m.width-8, m.height-8)

		return m, nil
	}

	// global keys
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.stage == stagePickOp {
				m.stage = stagePickFileSource
				return m, nil
			}
			if m.stage == stagePickFile {
				m.stage = stagePickFileSource
				m.fileInput.Blur()
				return m, nil
			}
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
			if m.stage == stageSave {
				if m.saveChoice != "" {
					m.saveChoice = ""
					m.statusMsg = ""
					m.outputInput.Blur()
					return m, nil
				}
				m.stage = stagePreview
				return m, nil
			}
		}
	}

	switch m.stage {
	case stagePickFileSource:
		var cmd tea.Cmd
		m.fileSourceList, cmd = m.fileSourceList.Update(msg)

		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
			switch m.fileSourceList.Index() {
			case 0:
				items, err := listFiles("files")
				if err != nil {
					m.errMsg = fmt.Sprintf("Cannot read files in directory: %v", err)
					return m, nil
				}

				m.fileList.SetItems(items)
				m.fileList.SetSize(m.width-8, m.height-8)

				m.errMsg = ""
				m.stage = stagePickFile
				return m, nil
			case 1:
				m.stage = stagePickFile
				m.fileInput.Focus()
				return m, nil
			}
		}

		return m, cmd

	case stagePickFile:
		// file list found
		if len(m.fileList.Items()) > 0 {
			var cmd tea.Cmd
			m.fileList, cmd = m.fileList.Update(msg)

			if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
				selected := m.fileList.SelectedItem()
				it, ok := selected.(fileItem)
				if !ok {
					m.errMsg = "Could not read selected file item."
					return m, nil
				}

				b, err := readFile(it.path)
				if err != nil {
					m.errMsg = fmt.Sprintf("Cannot read file: %v", err)
					return m, nil
				}

				m.filePath = it.path
				m.fileBytes = b
				m.errMsg = ""

				m.stage = stagePickOp
				return m, nil
			}

			return m, cmd

		}

		// no file list
		var cmd tea.Cmd
		m.fileInput, cmd = m.fileInput.Update(msg)

		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
			path := m.fileInput.Value()
			b, err := readFile(path)
			if err != nil {
				m.errMsg = fmt.Sprintf("Cannot read file: %v", err)
				return m, nil
			}
			m.filePath = path
			m.fileBytes = b
			m.errMsg = ""
			m.fileInput.Blur()

			m.stage = stagePickOp
			return m, nil
		}
		return m, cmd

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
				m.outputBytes = transform.Apply(
					m.fileBytes,
					transform.Operation(m.op),
					m.prefixInput.Value(),
					m.suffixInput.Value(),
					m.splitChoice,
				)
				m.stage = stagePreview
				return m, nil
			}
			return m, cmd

		case opSuffix:
			var cmd tea.Cmd
			m.suffixInput, cmd = m.suffixInput.Update(msg)
			if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
				m.outputBytes = transform.Apply(
					m.fileBytes,
					transform.Operation(m.op),
					m.prefixInput.Value(),
					m.suffixInput.Value(),
					m.splitChoice,
				)
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
				m.outputBytes = transform.Apply(
					m.fileBytes,
					transform.Operation(m.op),
					m.prefixInput.Value(),
					m.suffixInput.Value(),
					m.splitChoice,
				)
				m.stage = stagePreview
				return m, nil
			}
			return m, cmd

		case opAddComma:
			if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
				m.outputBytes = transform.Apply(
					m.fileBytes,
					transform.Operation(m.op),
					m.prefixInput.Value(),
					m.suffixInput.Value(),
					m.splitChoice,
				)
				m.stage = stagePreview
				return m, nil
			}
			return m, nil
		}

	case stagePreview:
		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
			m.stage = stageSave

			if m.filePath != "" {
				m.outputInput.SetValue(m.filePath + ".out.txt")
			}
			m.outputInput.Focus()

			return m, nil
		}

	case stageSave:
		if m.saveChoice == "" {
			var cmd tea.Cmd
			m.saveList, cmd = m.saveList.Update(msg)

			if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
				m.statusMsg = ""
				switch m.saveList.Index() {
				case 0:
					m.saveChoice = "new"
					m.outputInput.Focus()
				case 1:
					m.saveChoice = "overwrite"
				}
				return m, nil
			}

			return m, cmd

		}
		// ask for path if new file chosen
		if m.saveChoice == "new" {
			var cmd tea.Cmd
			m.outputInput, cmd = m.outputInput.Update(msg)

			if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
				outPath := m.outputInput.Value()
				if outPath == "" {
					m.statusMsg = "Output path cannot be empty."
					return m, nil
				}

				err := writeNewFile(outPath, m.outputBytes)
				if err != nil {
					m.statusMsg = fmt.Sprintf("Save Failed: %v", err)
					return m, nil
				}

				m.statusMsg = "Saved successfully to: " + outPath
				return m, nil
			}

			return m, cmd
		}

		// Overwrite file chosen
		if m.saveChoice == "overwrite" {
			if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
				backupPath, err := backupFile(m.filePath)
				if err != nil {
					m.statusMsg = fmt.Sprintf("Backup failed: %v", err)
					return m, nil
				}

				err = overwriteFile(m.filePath, m.outputBytes)
				if err != nil {
					m.statusMsg = fmt.Sprintf("Overwrite failed: %v", err)
					return m, nil
				}

				m.statusMsg = "Overwritten successfully. Backup created: " + backupPath
				return m, nil

			}
		}

	}

	return m, nil
}


