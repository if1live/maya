package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntVal(t *testing.T) {
	cases := []struct {
		params     map[string]string
		key        string
		defaultVal int
		expected   int
	}{
		{map[string]string{"key": "123"}, "key", 1, 123},
		{map[string]string{"key": "123"}, "not-exist", 1, 1},
		{map[string]string{"key": "invalid"}, "key", 1, 1},
	}
	for _, c := range cases {
		ca := CommandArguments{c.params}
		assert.Equal(t, c.expected, ca.IntVal(c.key, c.defaultVal))
	}
}

func TestStringVal(t *testing.T) {
	cases := []struct {
		params     map[string]string
		key        string
		defaultVal string
		expected   string
	}{
		{map[string]string{"key": "123"}, "key", "1", "123"},
		{map[string]string{"key": "123"}, "not-exist", "default", "default"},
	}
	for _, c := range cases {
		ca := CommandArguments{c.params}
		assert.Equal(t, c.expected, ca.StringVal(c.key, c.defaultVal))
	}
}

func TestBoolVal(t *testing.T) {
	cases := []struct {
		params     map[string]string
		key        string
		defaultVal bool
		expected   bool
	}{
		{map[string]string{"key": "123"}, "key", true, true},
		{map[string]string{"key": "123"}, "not-exist", true, true},
		{map[string]string{"key": "false"}, "key", true, false},
	}
	for _, c := range cases {
		ca := CommandArguments{c.params}
		assert.Equal(t, c.expected, ca.BoolVal(c.key, c.defaultVal))
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
			CommandExecute{"ls | sort | grep \".go\" | head -n 1", false, OutputFormatCode},
			[]string{"article.go", ""},
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
			CommandUnknown{"foo"},
			[]string{"Action=foo"},
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
			NewCommand("view", &CommandArguments{map[string]string{"file": "hello.txt"}}),
			&CommandView{"hello.txt", 0, 0, "txt", OutputFormatCode},
		},
		{
			NewCommand("view", &CommandArguments{map[string]string{
				"file":       "foo.txt",
				"start_line": "1",
				"end_line":   "10",
				"format":     "blockquote",
			}}),
			&CommandView{"foo.txt", 1, 10, "txt", OutputFormatBlockquote},
		},
		{
			NewCommand("view", &CommandArguments{map[string]string{
				"file": "hello.txt",
				"lang": "lisp",
			}}),
			&CommandView{"hello.txt", 0, 0, "lisp", OutputFormatCode},
		},
		{
			NewCommand("execute", &CommandArguments{map[string]string{
				"cmd": "echo hello",
			}}),
			&CommandExecute{"echo hello", false, OutputFormatCode},
		},
		{
			NewCommand("execute", &CommandArguments{map[string]string{
				"cmd":    "echo hello",
				"format": "blockquote",
			}}),
			&CommandExecute{"echo hello", false, OutputFormatBlockquote},
		},
		{
			NewCommand("execute", &CommandArguments{map[string]string{
				"cmd":        "echo hello",
				"format":     "blockquote",
				"attach_cmd": "t",
			}}),
			&CommandExecute{"echo hello", true, OutputFormatBlockquote},
		},
		{
			NewCommand("youtube", &CommandArguments{map[string]string{
				"video_id": "id",
				"width":    "480",
				"height":   "320",
			}}),
			&CommandYoutube{"id", 480, 320},
		},
		{
			NewCommand("hello", &CommandArguments{map[string]string{"key": "value"}}),
			&CommandUnknown{"hello"},
		},
	}
	for _, c := range cases {
		if !reflect.DeepEqual(c.actual, c.expected) {
			t.Errorf("CreateCommand - expected %Q, got %Q", c.expected, c.actual)
		}
	}
}
