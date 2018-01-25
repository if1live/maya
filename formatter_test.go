package maya

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_codeFormatter_format(t *testing.T) {
	cases := []struct {
		lines  []string
		args   []string
		output string
	}{
		{[]string{}, []string{}, "```\n```"},
		{
			[]string{"hello", "world"},
			[]string{},
			"```\nhello\nworld\n```",
		},
		{
			[]string{"", "", "hello", "world", "", ""},
			[]string{},
			"```\nhello\nworld\n```",
		},
		{
			[]string{"hello", "", "world"},
			[]string{},
			"```\nhello\n\nworld\n```",
		},
		{
			[]string{"hello", "world"},
			[]string{"python"},
			"```python\nhello\nworld\n```",
		},
		{
			[]string{"hello", "world"},
			[]string{"py"},
			"```python\nhello\nworld\n```",
		},
	}
	for _, c := range cases {
		f := codeFormatter{}
		assert.Equal(t, c.output, f.format(c.lines, c.args...))
	}
}

func Test_blockquoteFormatter_format(t *testing.T) {
	cases := []struct {
		lines  []string
		args   []string
		output string
	}{
		{
			[]string{"hello", "world"},
			[]string{},
			"> hello\n>\n> world",
		},
		{
			[]string{"hello", "", "world"},
			[]string{},
			"> hello\n>\n>\n>\n> world",
		},
	}
	for _, c := range cases {
		f := blockquoteFormatter{}
		assert.Equal(t, c.output, f.format(c.lines, c.args...))
	}
}

func Test_textFormatter_format(t *testing.T) {
	cases := []struct {
		lines  []string
		args   []string
		output string
	}{
		{
			[]string{"hello", "world"},
			[]string{},
			"hello\nworld",
		},
	}
	for _, c := range cases {
		f := textFormatter{}
		assert.Equal(t, c.output, f.format(c.lines, c.args...))
	}
}

func Test_boldFormatter_format(t *testing.T) {
	cases := []struct {
		lines  []string
		args   []string
		output string
	}{
		{
			[]string{"hello", "world"},
			[]string{},
			"**hello**\n**world**",
		},
	}
	for _, c := range cases {
		f := boldFormatter{}
		assert.Equal(t, c.output, f.format(c.lines, c.args...))
	}
}
