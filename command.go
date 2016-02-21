package main

import (
	"io/ioutil"
	"net/url"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/op/go-logging"
)

const (
	OutputFormatCode       = "code"
	OutputFormatBlockquote = "blockquote"
	OutputFormatBold       = "bold"
)

type Command interface {
	Run() string
}

type OutputFormatter struct {
	Format string
}

func (f *OutputFormatter) Run(text string, args ...string) string {
	switch f.Format {
	case OutputFormatCode:
		content := "```"
		if len(args) > 0 {
			content += args[0]
		}

		content = content + "\n" + text

		if text[len(text)-1] == byte('\n') {
			content += "```"
		} else {
			content += "\n```"
		}
		return content

	case OutputFormatBold:
		return "**" + text + "**"

	case OutputFormatBlockquote:
		lines := strings.Split(text, "\n")
		for i, line := range lines {
			isLastLine := (i == len(lines)-1)
			if !isLastLine || (isLastLine && len(line) != 0) {
				lines[i] = strings.TrimRight("> "+line, " ")
			}
		}
		return strings.Join(lines, "\n")

	default:
		msg := "unknown format : " + f.Format
		panic(msg)
	}
}

type CommandViewFile struct {
	FilePath  string
	StartLine int
	EndLine   int
	Language  string
	Format    string
}

func (c *CommandViewFile) Run() string {
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

	text := strings.Join(lines[c.StartLine:c.EndLine], "\n")
	formatter := OutputFormatter{c.Format}
	return formatter.Run(text, c.Language)
}

type CommandExecute struct {
	Cmd       string
	AttachCmd bool
	Format    string
}

func (c *CommandExecute) SplitCommand() (string, []string) {
	tokens := strings.Split(c.Cmd, " ")
	return tokens[0], tokens[1:]
}

func (c *CommandExecute) Run() string {
	log := logging.MustGetLogger("maya")
	log.Infof("Command execute: %v", c)
	name, args := c.SplitCommand()
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		panic(err)
	}

	elems := []string{}
	if c.AttachCmd {
		elems = append(elems, "$ "+c.Cmd)
	}
	elems = append(elems, string(out[:]))
	text := strings.Join(elems, "\n")
	formatter := OutputFormatter{c.Format}
	return formatter.Run(text)
}

type CommandUnknown struct {
	Action string
	Params string
}

func (c *CommandUnknown) Run() string {
	log := logging.MustGetLogger("maya")
	log.Warningf("Command Unknown: %v", c)
	tokens := []string{
		"**",
		"Action=" + c.Action,
		", ",
		"Params=" + c.Params,
		"**",
	}
	return strings.Join(tokens, "")
}

func NewCommand(action string, params string) Command {
	// 파싱하기 귀찮은 관계로 url query string에 묻어가자
	// 나중에 개선하기
	values, err := url.ParseQuery(strings.Replace(params, ",", "&", -1))
	if err != nil {
		panic(err)
	}

	format := values.Get("fmt")
	if format == "" {
		format = OutputFormatCode
	}

	switch action {
	case "view":
		startLine, _ := strconv.Atoi(values.Get("start"))
		endLine, _ := strconv.Atoi(values.Get("end"))
		language := values.Get("lang")
		filePath := values.Get("file")
		if language == "" {
			language = strings.Replace(filepath.Ext(filePath), ".", "", -1)
		}

		return &CommandViewFile{
			FilePath:  filePath,
			StartLine: startLine,
			EndLine:   endLine,
			Language:  language,
			Format:    format,
		}
	case "execute":
		attachCmd := false
		if len(values.Get("attach_cmd")) > 0 {
			attachCmd = true
		}
		return &CommandExecute{
			Cmd:       values.Get("cmd"),
			AttachCmd: attachCmd,
			Format:    format,
		}
	default:
		return &CommandUnknown{
			Action: action,
			Params: params,
		}
	}
}
