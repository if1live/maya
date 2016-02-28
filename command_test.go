package main

import (
	"reflect"
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

func TestRawOutpoutCommandExecute(t *testing.T) {
	cases := []struct {
		cmd    CommandExecute
		output []string
	}{
		{
			CommandExecute{"echo hello", false, OutputFormatCode},
			[]string{"hello", ""},
		},
		// stderr
		{
			CommandExecute{"clang", false, OutputFormatCode},
			[]string{"clang: error: no input files", ""},
		},
		{
			CommandExecute{"clang", true, OutputFormatCode},
			[]string{"$ clang", "clang: error: no input files", ""},
		},
		// command not exist
		{
			CommandExecute{"invalid-cmd", false, OutputFormatCode},
			[]string{"exec: \"invalid-cmd\": executable file not found in $PATH"},
		},
		{
			CommandExecute{"invalid", true, OutputFormatCode},
			[]string{"$ invalid", "exec: \"invalid\": executable file not found in $PATH"},
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.output, c.cmd.RawOutput())
	}
}

func TestRawOutputCommandView(t *testing.T) {
	cases := []struct {
		cmd    CommandView
		output []string
	}{
		{
			CommandView{
				FilePath:  "command_test.go",
				StartLine: 1,
				EndLine:   3,
				Format:    OutputFormatCode,
			},
			[]string{"", "import ("},
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.output, c.cmd.RawOutput())
	}
}

func TestRawOutputCommandUnknown(t *testing.T) {
	cases := []struct {
		cmd    CommandUnknown
		output []string
	}{
		{
			CommandUnknown{"foo", "bar"},
			[]string{"Action=foo", "Params=bar"},
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.output, c.cmd.RawOutput())
	}
}

func TestNewCommand(t *testing.T) {
	cases := []struct {
		actual   Command
		expected Command
	}{
		{
			NewCommand("view", "file=hello.txt"),
			&CommandView{"hello.txt", 0, 0, "txt", OutputFormatCode},
		},
		{
			NewCommand("view", "file=foo.txt,start=1,end=10,fmt=blockquote"),
			&CommandView{"foo.txt", 1, 10, "txt", OutputFormatBlockquote},
		},
		{
			NewCommand("view", "file=hello.txt,lang=lisp"),
			&CommandView{"hello.txt", 0, 0, "lisp", OutputFormatCode},
		},
		{
			NewCommand("execute", "cmd=echo hello"),
			&CommandExecute{"echo hello", false, OutputFormatCode},
		},
		{
			NewCommand("execute", "cmd=echo hello,fmt=blockquote"),
			&CommandExecute{"echo hello", false, OutputFormatBlockquote},
		},
		{
			NewCommand("execute", "cmd=echo hello,fmt=blockquote,attach_cmd=t"),
			&CommandExecute{"echo hello", true, OutputFormatBlockquote},
		},
		{
			NewCommand("hello", "world"),
			&CommandUnknown{"hello", "world"},
		},
	}
	for _, c := range cases {
		if !reflect.DeepEqual(c.actual, c.expected) {
			t.Errorf("CreateCommand - expected %Q, got %Q", c.expected, c.actual)
		}
	}
}
