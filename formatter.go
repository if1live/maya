package maya

import (
	"strings"
)

const (
	formatCode       = "code"
	formatBlockquote = "blockquote"
	formatBold       = "bold"
	formatText       = "text"
)

type Formatter interface {
	format(lines []string, args ...string) string
}

func format(format string, text string, args ...string) string {
	f := newFormatter(format)
	lines := strings.Split(text, "\n")
	return f.format(lines, args...)
}

func newFormatter(format string) Formatter {
	switch format {
	case formatCode:
		return &codeFormatter{}
	case formatBlockquote:
		return &blockquoteFormatter{}
	case formatBold:
		return &boldFormatter{}
	case formatText:
		return &textFormatter{}
	default:
		msg := "unknown format : " + format
		panic(msg)
	}
}

type codeFormatter struct{}

func (f *codeFormatter) getLanguage(args ...string) string {
	lang := ""
	if len(args) > 0 {
		lang = args[0]
	}
	return f.convertLanguage(lang)
}

func (f *codeFormatter) convertLanguage(lang string) string {
	table := map[string]string{
		"cs": "csharp",
		"py": "python",
		"rb": "ruby",
	}

	found, ok := table[lang]
	if ok {
		return found
	}
	return lang
}

func (f *codeFormatter) format(lines []string, args ...string) string {
	lang := f.getLanguage(args...)
	headLine := "```" + lang
	tailLine := "```"

	if len(lines) == 0 {
		return strings.Join([]string{headLine, tailLine}, "\n")
	}

	isBlankLine := func(line string) bool {
		return strings.Trim(line, "\t\r ") == ""
	}

	startIdx := 0
	for i := 0; i < len(lines); i++ {
		if !isBlankLine(lines[i]) {
			startIdx = i
			break
		}
	}

	endIdx := len(lines) - 1
	for i := len(lines) - 1; i >= 0; i-- {
		if !isBlankLine(lines[i]) {
			endIdx = i
			break
		}
	}

	newLines := []string{}
	newLines = append(newLines, headLine)
	newLines = append(newLines, lines[startIdx:endIdx+1]...)
	newLines = append(newLines, tailLine)
	return strings.Join(newLines, "\n")
}

type blockquoteFormatter struct{}

func (f *blockquoteFormatter) format(lines []string, args ...string) string {
	contents := make([]string, len(lines)*2-1)
	for i, line := range lines {
		contents[i*2+0] = "> " + line
		if i != len(lines)-1 {
			contents[i*2+1] = "> "
		}
	}
	for i, line := range contents {
		contents[i] = strings.Trim(line, " ")
	}
	return strings.Join(contents, "\n")
}

type boldFormatter struct{}

func (f *boldFormatter) format(lines []string, argss ...string) string {
	contents := make([]string, len(lines))
	for i, line := range lines {
		contents[i] = "**" + line + "**"
	}
	return strings.Join(contents, "\n")
}

type textFormatter struct{}

func (f *textFormatter) format(lines []string, args ...string) string {
	return strings.Join(lines, "\n")
}
