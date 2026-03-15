package main

import (
	"Rolyani/utilities-cli/transform"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	hintStyle  = lipgloss.NewStyle().Faint(true)
	boxStyle   = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder())
)

func (m model) View() string {
	if m.quitting {
		return ""
	}

	header := titleStyle.Render("utilities-cli") + "\n" +
		hintStyle.Render("ctrl + c: quit • esc: back • enter: next") + "\n\n"

	switch m.stage {
	case stagePickFileSource:
		return header + boxStyle.Render(m.fileSourceList.View())
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


