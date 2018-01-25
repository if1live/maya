package maya

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_commandArgs_intVal(t *testing.T) {
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
		ca := cmdArgs{c.params}
		assert.Equal(t, c.expected, ca.intVal(c.key, c.defaultVal))
	}
}

func Test_commandArgs_stringVal(t *testing.T) {
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
		ca := cmdArgs{c.params}
		assert.Equal(t, c.expected, ca.stringVal(c.key, c.defaultVal))
	}
}

func Test_CommandArgs_boolVal(t *testing.T) {
	cases := []struct {
		params     map[string]string
		key        string
		defaultVal bool
		expected   bool
	}{
		{map[string]string{"key": "123"}, "key", true, true},
		{map[string]string{"key": "123"}, "not-exist", true, true},

		{map[string]string{"key": "false"}, "key", true, false},
		{map[string]string{"key": "FaLsE"}, "key", true, false},

		{map[string]string{"key": "TrUe"}, "key", false, true},
		{map[string]string{"key": "t"}, "key", false, true},
	}
	for _, c := range cases {
		ca := cmdArgs{c.params}
		assert.Equal(t, c.expected, ca.boolVal(c.key, c.defaultVal))
	}
}

func TestRawOutpoutCommandExecute(t *testing.T) {
	cases := []struct {
		supportWindows bool
		cmd            cmdExecute
		output         []string
	}{
		{
			true,
			cmdExecute{"echo hello", false, formatCode},
			[]string{"hello", ""},
		},
		// stderr
		{
			false,
			cmdExecute{"./demo_stderr.py", false, formatCode},
			[]string{"this is stderr", ""},
		},
		{
			false,
			cmdExecute{"./demo_stderr.py", true, formatCode},
			[]string{"$ ./demo_stderr.py", "this is stderr", ""},
		},
		// command not exist
		// 임시폴더를 쓰기때문에 임의의 경로가 나온다
		//{
		//	CommandExecute{"invalid-cmd", false, OutputFormatCode},
		//	[]string{"exec: \"invalid-cmd\": executable file not found in $PATH"},
		//},
		// local path
		{
			false,
			cmdExecute{"./demo.sh", true, formatCode},
			[]string{"$ ./demo.sh", "hello-world!", ""},
		},
		// complex
		{
			false,
			cmdExecute{"ls | sort | grep \".go\" | head -n 1", false, formatCode},
			[]string{"article.go", ""},
		},
	}
	for _, c := range cases {
		switch runtime.GOOS {
		case "windows":
			if c.supportWindows {
				assert.Equal(t, c.output, c.cmd.RawOutput())
			}
		default:
			assert.Equal(t, c.output, c.cmd.RawOutput())
		}
	}
}

func TestRawOutputCommandView(t *testing.T) {
	cases := []struct {
		cmd    cmdView
		output []string
	}{
		{
			cmdView{
				filePath:  "cmd_test.go",
				startLine: 1,
				endLine:   3,
				format:    formatCode,
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
		cmd    cmdUnknown
		output []string
	}{
		{
			cmdUnknown{"foo"},
			[]string{"Action=foo"},
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.output, c.cmd.RawOutput())
	}
}

func TestNewCommand(t *testing.T) {
	cases := []struct {
		actual   cmd
		expected cmd
	}{
		{
			newCmd("view", &cmdArgs{map[string]string{"file": "hello.txt"}}),
			&cmdView{"hello.txt", 0, 0, "txt", formatCode},
		},
		{
			newCmd("view", &cmdArgs{map[string]string{
				"file":       "foo.txt",
				"start_line": "1",
				"end_line":   "10",
				"format":     "blockquote",
			}}),
			&cmdView{"foo.txt", 1, 10, "txt", formatBlockquote},
		},
		{
			newCmd("view", &cmdArgs{map[string]string{
				"file": "hello.txt",
				"lang": "lisp",
			}}),
			&cmdView{"hello.txt", 0, 0, "lisp", formatCode},
		},
		{
			newCmd("execute", &cmdArgs{map[string]string{
				"cmd": "echo hello",
			}}),
			&cmdExecute{"echo hello", false, formatCode},
		},
		{
			newCmd("execute", &cmdArgs{map[string]string{
				"cmd":    "echo hello",
				"format": "blockquote",
			}}),
			&cmdExecute{"echo hello", false, formatBlockquote},
		},
		{
			newCmd("execute", &cmdArgs{map[string]string{
				"cmd":        "echo hello",
				"format":     "blockquote",
				"attach_cmd": "t",
			}}),
			&cmdExecute{"echo hello", true, formatBlockquote},
		},
		{
			newCmd("youtube", &cmdArgs{map[string]string{
				"video_id": "id",
				"width":    "480",
				"height":   "320",
			}}),
			&cmdYoutube{"id", 480, 320},
		},
		{
			newCmd("gist", &cmdArgs{map[string]string{
				"id":   "3254906",
				"file": "brew-update-notifier.sh",
			}}),
			&cmdGist{"3254906", "brew-update-notifier.sh"},
		},
		{
			newCmd("hello", &cmdArgs{map[string]string{"key": "value"}}),
			&cmdUnknown{"hello"},
		},
	}
	for _, c := range cases {
		if !reflect.DeepEqual(c.actual, c.expected) {
			t.Errorf("CreateCommand - expected %Q, got %Q", c.expected, c.actual)
		}
	}
}
