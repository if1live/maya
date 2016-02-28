package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		format string
		lines  []string
		args   []string
		output string
	}{
		{OutputFormatCode, []string{}, []string{}, ""},
		{
			OutputFormatCode,
			[]string{"hello", "world"},
			[]string{},
			"```\nhello\nworld\n```",
		},
		{
			OutputFormatCode,
			[]string{"hello", "world"},
			[]string{"python"},
			"```python\nhello\nworld\n```",
		},
		{
			OutputFormatBlockquote,
			[]string{"hello", "world"},
			[]string{},
			"> hello\n>\n> world",
		},
		{
			OutputFormatBlockquote,
			[]string{"hello", "", "world"},
			[]string{},
			"> hello\n>\n>\n>\n> world",
		},
		{
			OutputFormatBold,
			[]string{"hello", "world"},
			[]string{},
			"**hello**\n**world**",
		},
	}
	for _, c := range cases {
		f := OutputFormatter{c.format}
		assert.Equal(t, c.output, f.Format(c.lines, c.args...))
	}
}
