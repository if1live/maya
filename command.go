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

const (
	CommandTypeView    = "view"
	CommandTypeExecute = "execute"
	CommandTypeYoutube = "youtube"
	CommandTypeUnknown = "unknown"
)

type Command interface {
	Run() string
	RawOutput() []string
}

type CommandArguments struct {
	Args map[string]string
}

type OutputFormatter struct {
	format string
}

func (f *OutputFormatter) Format(lines []string, args ...string) string {
	if len(lines) == 0 {
		return ""
	}

	switch f.format {
	case OutputFormatCode:
		lang := ""
		if len(args) > 0 {
			lang = args[0]
		}
		return f.formatCode(lines, lang)
	case OutputFormatBlockquote:
		return f.formatBlockquote(lines)
	case OutputFormatBold:
		return f.formatBold(lines)
	default:
		msg := "unknown format : " + f.format
		panic(msg)
	}
	return ""
}

func (f *OutputFormatter) formatCode(lines []string, lang string) string {
	headLine := "```" + lang
	tailLine := "```"

	contents := make([]string, len(lines)+2)
	contents[0] = headLine
	for i, line := range lines {
		contents[i+1] = line
	}
	contents[len(contents)-1] = tailLine
	return strings.Join(contents, "\n")
}

func (f *OutputFormatter) formatBlockquote(lines []string) string {
	contents := make([]string, len(lines)*2-1)
	for i, line := range lines {
		contents[i*2+0] = "> " + line
		if i != len(lines)-1 {
			contents[i*2+1] = "> "
		}
	}
	for i, line := range contents {
		contents[i] = strings.Trim(line, " ")
	}
	return strings.Join(contents, "\n")
}

func (f *OutputFormatter) formatBold(lines []string) string {
	contents := make([]string, len(lines))
	for i, line := range lines {
		contents[i] = "**" + line + "**"
	}
	return strings.Join(contents, "\n")
}

func (f *OutputFormatter) Run(text string, args ...string) string {
	lines := strings.Split(text, "\n")
	return f.Format(lines, args...)
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

func (c *CommandView) Run() string {
	text := strings.Join(c.RawOutput(), "\n")
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

func (c *CommandExecute) RawOutput() []string {
	log := logging.MustGetLogger("maya")
	log.Infof("Command execute: %v", c)

	elems := []string{}
	if c.AttachCmd {
		elems = append(elems, "$ "+c.Cmd)
	}

	name, args := c.SplitCommand()
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.Error); ok {
			elems = append(elems, err.Error())
			return elems
		}
	}

	elems = append(elems, strings.Split(string(out[:]), "\n")...)
	return elems
}

func (c *CommandExecute) Run() string {
	text := strings.Join(c.RawOutput(), "\n")
	formatter := OutputFormatter{c.Format}
	return formatter.Run(text)
}

type CommandUnknown struct {
	Action string
	Params string
}

func (c *CommandUnknown) RawOutput() []string {
	log := logging.MustGetLogger("maya")
	log.Warningf("Command Unknown: %v", c)
	tokens := []string{
		"Action=" + c.Action,
		"Params=" + c.Params,
	}
	return tokens
}
func (c *CommandUnknown) Run() string {
	text := strings.Join(c.RawOutput(), "\n")
	formatter := OutputFormatter{OutputFormatBlockquote}
	return formatter.Run(text)
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

		return &CommandView{
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
