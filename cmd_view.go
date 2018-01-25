package maya

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
)

type cmdView struct {
	filePath  string
	startLine int
	endLine   int
	language  string
	format    string
}

func newCmdView(action string, args *cmdArgs) cmd {
	filePath := args.stringVal("file", "")
	defaultLang := strings.Replace(filepath.Ext(filePath), ".", "", -1)
	language := args.stringVal("lang", defaultLang)

	return &cmdView{
		filePath:  filePath,
		startLine: args.intVal("start_line", 0),
		endLine:   args.intVal("end_line", 0),
		language:  language,
		format:    args.stringVal("format", formatCode),
	}
}

func (c *cmdView) RawOutput() []string {
	log := logging.MustGetLogger("maya")
	log.Infof("Command ViewFile: %v", c)
	data, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data[:]), "\n")

	if c.startLine == 0 && c.endLine == 0 {
		c.endLine = len(lines)
	}

	elems := lines[c.startLine:c.endLine]
	elems = sanitizeLineFeedMultiLine(elems)
	return elems
}

func (c *cmdView) execute() string {
	f := newFormatter(c.format)
	return f.format(c.RawOutput(), c.language)
}
