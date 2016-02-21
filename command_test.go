package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommandExecute(t *testing.T) {
	cases := []struct {
		cmd    CommandExecute
		output string
	}{
		{
			CommandExecute{"echo hello", false, OutputFormatCode},
			"```\nhello\n```",
		},
		{
			CommandExecute{"echo hello", false, OutputFormatBlockquote},
			"> hello\n",
		},
		{
			CommandExecute{"echo hello", true, OutputFormatBlockquote},
			"> $ echo hello\n> hello\n",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.output, c.cmd.Run())
	}
}

func TestRunCommandViewFile(t *testing.T) {
	cases := []struct {
		cmd    CommandViewFile
		output string
	}{
		{
			CommandViewFile{
				FilePath:  "command_test.go",
				StartLine: 1,
				EndLine:   3,
				Format:    OutputFormatCode,
			},
			"```\n\nimport (\n```",
		},
		{
			CommandViewFile{
				FilePath:  "command_test.go",
				StartLine: 1,
				EndLine:   3,
				Format:    OutputFormatBlockquote,
			},
			"> \n> import (",
		},
		{
			CommandViewFile{
				FilePath:  "command_test.go",
				StartLine: 1,
				EndLine:   3,
				Language:  "go",
				Format:    OutputFormatCode,
			},
			"```go\n\nimport (\n```",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.output, c.cmd.Run())
	}
}

func TestRunCommandUnknown(t *testing.T) {
	cases := []struct {
		cmd    CommandUnknown
		output string
	}{
		{
			CommandUnknown{"foo", "bar"},
			"**Action=foo, Params=bar**",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.output, c.cmd.Run())
	}
}

func TestNewCommand(t *testing.T) {
	cases := []struct {
		actual   Command
		expected Command
	}{
		{
			NewCommand("view", "file=hello.txt"),
			&CommandViewFile{"hello.txt", 0, 0, "txt", OutputFormatCode},
		},
		{
			NewCommand("view", "file=foo.txt,start=1,end=10,fmt=blockquote"),
			&CommandViewFile{"foo.txt", 1, 10, "txt", OutputFormatBlockquote},
		},
		{
			NewCommand("view", "file=hello.txt,lang=lisp"),
			&CommandViewFile{"hello.txt", 0, 0, "lisp", OutputFormatCode},
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
