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
	params map[string]string
}

func (ca *CommandArguments) IntVal(key string, defaultVal int) int {
	val, err := strconv.Atoi(ca.params[key])
	if err != nil {
		return defaultVal
	}
	return val
}

func (ca *CommandArguments) StringVal(key string, defaultVal string) string {
	if val, ok := ca.params[key]; ok {
		return val
	} else {
		return defaultVal
	}
}

func (ca *CommandArguments) BoolVal(key string, defaultVal bool) bool {
	if val, ok := ca.params[key]; ok {
		if val == "false" || val == "f" {
			return false
		} else {
			return true
		}
	} else {
		return defaultVal
	}
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

func NewCommand(action string, args *CommandArguments) Command {
	switch action {
	case "view":
		filePath := args.StringVal("file", "")
		defaultLang := strings.Replace(filepath.Ext(filePath), ".", "", -1)
		language := args.StringVal("lang", defaultLang)

		return &CommandView{
			FilePath:  filePath,
			StartLine: args.IntVal("start_line", 0),
			EndLine:   args.IntVal("end_line", 0),
			Language:  language,
			Format:    args.StringVal("format", OutputFormatCode),
		}
	case "execute":
		return &CommandExecute{
			Cmd:       args.StringVal("cmd", "echo empty"),
			AttachCmd: args.BoolVal("attach_cmd", false),
			Format:    args.StringVal("format", OutputFormatCode),
		}
	default:
		return &CommandUnknown{
			Action: action,
		}
	}
}
