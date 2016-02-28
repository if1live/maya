package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/op/go-logging"
)

const (
	CommandTypeView    = "view"
	CommandTypeExecute = "execute"
	CommandTypeYoutube = "youtube"
	CommandTypeUnknown = "unknown"
)

type Command interface {
	Formatter() *OutputFormatter
	RawOutput() []string
}

type CommandArguments struct {
	Args map[string]string
}

type CommandView struct {
	FilePath  string
	StartLine int
	EndLine   int
	Language  string
	Format    string
}

func (c *CommandView) RawOutput() []string {
	log := logging.MustGetLogger("maya")
	log.Infof("Command ViewFile: %v", c)
	data, err := ioutil.ReadFile(c.FilePath)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data[:]), "\n")

	if c.StartLine == 0 && c.EndLine == 0 {
		c.EndLine = len(lines)
	}

	return lines[c.StartLine:c.EndLine]
}

func (c *CommandView) Formatter() *OutputFormatter {
	return &OutputFormatter{c.Format}
}

type CommandExecute struct {
	Cmd       string
	AttachCmd bool
	Format    string
}

func (c *CommandExecute) RawOutput() []string {
	log := logging.MustGetLogger("maya")
	log.Infof("Command execute: %v", c)

	elems := []string{}
	if c.AttachCmd {
		elems = append(elems, "$ "+c.Cmd)
	}

	tmpfile, err := ioutil.TempFile("", "maya")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(c.Cmd)); err != nil {
		panic(err)
	}
	if err := tmpfile.Close(); err != nil {
		panic(err)
	}
	out, err := exec.Command("bash", tmpfile.Name()).CombinedOutput()

	if err != nil {
		if _, ok := err.(*exec.Error); ok {
			elems = append(elems, err.Error())
			return elems
		}
	}

	elems = append(elems, strings.Split(string(out[:]), "\n")...)
	return elems
}

func (c *CommandExecute) Formatter() *OutputFormatter {
	return &OutputFormatter{c.Format}
}

type CommandUnknown struct {
	Action string
}

func (c *CommandUnknown) RawOutput() []string {
	log := logging.MustGetLogger("maya")
	log.Warningf("Command Unknown: %v", c)
	tokens := []string{
		"Action=" + c.Action,
	}
	return tokens
}
func (c *CommandUnknown) Formatter() *OutputFormatter {
	return &OutputFormatter{OutputFormatBlockquote}
}

func NewCommand(action string, params map[string]string) Command {
	// 파싱하기 귀찮은 관계로 url query string에 묻어가자
	// 나중에 개선하기
	format := params["fmt"]
	if format == "" {
		format = OutputFormatCode
	}

	switch action {
	case "view":
		startLine, _ := strconv.Atoi(params["start"])
		endLine, _ := strconv.Atoi(params["end"])
		language := params["lang"]
		filePath := params["file"]
		if language == "" {
			language = strings.Replace(filepath.Ext(filePath), ".", "", -1)
		}

		return &CommandView{
			FilePath:  filePath,
			StartLine: startLine,
			EndLine:   endLine,
			Language:  language,
			Format:    format,
		}
	case "execute":
		attachCmd := false
		if len(params["attach_cmd"]) > 0 {
			attachCmd = true
		}
		return &CommandExecute{
			Cmd:       params["cmd"],
			AttachCmd: attachCmd,
			Format:    format,
		}
	default:
		return &CommandUnknown{
			Action: action,
		}
	}
}
