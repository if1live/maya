package main

import (
	"regexp"
	"strings"
)

type ArticleContent struct {
	raw   string
	lines []string
}

func NewContent(text string) *ArticleContent {
	rawlines := strings.Split(text, "\n")
	lines := []string{}

	cmdRe := regexp.MustCompile(`^{{(\w+):(.+)}}$`)

	for _, line := range rawlines {
		cmdMatch := cmdRe.FindStringSubmatch(line)
		if len(cmdMatch) > 0 {
			action, params := cmdMatch[1], cmdMatch[2]
			cmd := NewCommand(action, params)
			lines = append(lines, cmd.Run())
		} else {
			lines = append(lines, line)
		}
	}

	return &ArticleContent{
		raw:   text,
		lines: lines,
	}
}

func (c *ArticleContent) String() string {
	return strings.Join(c.lines, "\n")
}
