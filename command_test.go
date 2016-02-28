package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		// 임시폴더를 쓰기때문에 임의의 경로가 나온다
		//{
		//	CommandExecute{"invalid-cmd", false, OutputFormatCode},
		//	[]string{"exec: \"invalid-cmd\": executable file not found in $PATH"},
		//},
		// local path
		{
			CommandExecute{"./demo.sh", true, OutputFormatCode},
			[]string{"$ ./demo.sh", "hello-world!", ""},
		},
		// complex
		{
			CommandExecute{"ls | sort | head -n 1", false, OutputFormatCode},
			[]string{"Godeps", ""},
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
