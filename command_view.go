package maya

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
)

type CommandView struct {
	FilePath  string
	StartLine int
	EndLine   int
	Language  string
	Format    string
}

func NewCommandView(action string, args *CommandArguments) Command {
	filePath := args.StringVal("file", "")
	defaultLang := strings.Replace(filepath.Ext(filePath), ".", "", -1)
	language := args.StringVal("lang", defaultLang)

	return &CommandView{
		FilePath:  filePath,
		StartLine: args.IntVal("start_line", 0),
		EndLine:   args.IntVal("end_line", 0),
		Language:  language,
		Format:    args.StringVal("format", formatCode),
	}
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

	elems := lines[c.StartLine:c.EndLine]
	elems = sanitizeLineFeedMultiLine(elems)
	return elems
}

func (c *CommandView) execute() string {
	f := newFormatter(c.Format)
	return f.format(c.RawOutput(), c.Language)
}
