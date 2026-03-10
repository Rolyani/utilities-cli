package transform

import "strings"

type Operation string

const (
	AddComma Operation = "Add comma to end of every line"
	Prefix   Operation = "Add text to beginning of every line"
	Suffix   Operation = "Add text to end of every line"
	Split    Operation = "Split after space or comma"
)

func Apply(in []byte, op Operation, prefix, suffix, splitChoice string) []byte {
	s := string(in)

	switch op {

	case AddComma:
		lines := strings.Split(s, "\n")
		for i := 0; i < len(lines); i++ {
			if i == len(lines)-1 && lines[i] == "" {
				continue
			}
			lines[i] += ","
		}
		return []byte(strings.Join(lines, "\n"))

	case Prefix:
		lines := strings.Split(s, "\n")
		for i := range lines {
			if i == len(lines)-1 && lines[i] == "" {
				continue
			}
			lines[i] = prefix + lines[i]
		}
		return []byte(strings.Join(lines, "\n"))

	case Suffix:
		lines := strings.Split(s, "\n")
		for i := range lines {
			if i == len(lines)-1 && lines[i] == "" {
				continue
			}
			lines[i] = lines[i] + suffix
		}
		return []byte(strings.Join(lines, "\n"))

	case Split:
		switch splitChoice {
		case "space":
			s = strings.ReplaceAll(s, " ", " \n")
		case "comma":
			s = strings.ReplaceAll(s, ",", ",\n")
		case "both":
			s = strings.ReplaceAll(s, " ", " \n")
			s = strings.ReplaceAll(s, ",", ",\n")
		}
		return []byte(s)

	default:
		return in
	}
}

func FirstNLines(b []byte, n int) string {
	lines := strings.Split(string(b), "\n")
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, "\n")
}
