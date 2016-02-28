package main

import "strings"

const (
	OutputFormatCode       = "code"
	OutputFormatBlockquote = "blockquote"
	OutputFormatBold       = "bold"
	OutputFormatText       = "text"
)

type OutputFormatter struct {
	format string
}

func (f *OutputFormatter) Format(lines []string, args ...string) string {
	if len(lines) == 0 {
		return ""
	}

	switch f.format {
	case OutputFormatCode:
		lang := ""
		if len(args) > 0 {
			lang = args[0]
		}
		return f.formatCode(lines, lang)
	case OutputFormatBlockquote:
		return f.formatBlockquote(lines)
	case OutputFormatBold:
		return f.formatBold(lines)
	case OutputFormatText:
		return f.formatText(lines)
	default:
		msg := "unknown format : " + f.format
		panic(msg)
	}
	return ""
}

func (f *OutputFormatter) formatCode(lines []string, lang string) string {
	headLine := "```" + f.convertLanguage(lang)
	tailLine := "```"

	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	contents := make([]string, len(lines)+2)
	contents[0] = headLine
	for i, line := range lines {
		contents[i+1] = line
	}
	contents[len(contents)-1] = tailLine
	return strings.Join(contents, "\n")
}

func (f *OutputFormatter) convertLanguage(lang string) string {
	table := map[string]string{
		"cs": "csharp",
		"py": "python",
		"rb": "ruby",
	}

	found := table[lang]
	if found != "" {
		return found
	} else {
		return lang
	}
}

func (f *OutputFormatter) formatBlockquote(lines []string) string {
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

func (f *OutputFormatter) formatBold(lines []string) string {
	contents := make([]string, len(lines))
	for i, line := range lines {
		contents[i] = "**" + line + "**"
	}
	return strings.Join(contents, "\n")
}

func (f *OutputFormatter) formatText(lines []string) string {
	return strings.Join(lines, "\n")
}

func (f *OutputFormatter) Run(text string, args ...string) string {
	lines := strings.Split(text, "\n")
	return f.Format(lines, args...)
}
