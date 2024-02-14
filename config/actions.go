package config

import "strings"

type ReplaceLine struct {
	Old string
	New string
}

type ReplaceFullLine struct {
	*ReplaceLine
}

type ReplaceSubLine struct {
	*ReplaceLine
}

type ReplaceWord struct {
	*ReplaceLine
}

type ReplaceFirstWord struct {
	*ReplaceLine
}

func (l *ReplaceFullLine) Action(line string) string {
	if line == l.ReplaceLine.Old {
		return l.ReplaceLine.New
	}
	return line
}

func (l *ReplaceSubLine) Action(line string) string {
	if strings.Contains(line, l.ReplaceLine.Old) {
		return l.ReplaceLine.New
	}
	return line
}

func (l *ReplaceWord) Action(line string) string {
	if strings.Contains(line, l.ReplaceLine.Old) {
		return strings.Replace(line, l.ReplaceLine.Old, l.ReplaceLine.New, -1)
	}
	return line
}
