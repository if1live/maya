package main

import (
	"regexp"
	"strings"
)

type ArticleContent struct {
	raw    string
	blocks []ContentBlock
}

type ContentBlock struct {
	command string
	lines   []string
}

func (cb *ContentBlock) Lines() []string {
	if cb.command == "" {
		return cb.lines
	}

	re := regexp.MustCompile(`^(\w+)\s*=(.*)$`)
	params := map[string]string{}
	for _, line := range cb.lines {
		m := re.FindStringSubmatch(line)
		if len(m) == 0 {
			continue
		}
		key, value := m[1], m[2]
		params[key] = value
	}
	cmd := NewCommand(cb.command, params)
	formatter := cmd.Formatter()
	return []string{formatter.Format(cmd.RawOutput())}
}

func NewContent(text string) *ArticleContent {
	rawlines := strings.Split(text, "\n")

	cmdStartRe := regexp.MustCompile(`^~~~maya:(\w+)\s*$`)
	cmdEndRe := regexp.MustCompile(`^~~~\s*$`)
	buffer := []string{}
	blocks := []ContentBlock{}

	state := ""

	for _, line := range rawlines {
		switch state {
		case "":
			m := cmdStartRe.FindStringSubmatch(line)
			if len(m) > 0 {
				blocks = append(blocks, ContentBlock{state, buffer})

				state = m[1]
				buffer = []string{line}
			} else {
				buffer = append(buffer, line)
			}
		default:
			m := cmdEndRe.FindString(line)
			if m != "" {
				buffer = append(buffer, line)
				blocks = append(blocks, ContentBlock{state, buffer})

				state = ""
				buffer = []string{}
			} else {
				buffer = append(buffer, line)
			}
		}
	}
	blocks = append(blocks, ContentBlock{state, buffer})

	return &ArticleContent{
		raw:    text,
		blocks: blocks,
	}
}

func (c *ArticleContent) String() string {
	lines := []string{}
	for _, block := range c.blocks {
		lines = append(lines, block.Lines()...)
	}
	return strings.Join(lines, "\n")
}
