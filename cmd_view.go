package maya

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
)

type cmdView struct {
	FilePath  string `maya:"file"`
	StartLine int    `maya:"start_line,0"`
	EndLine   int    `maya:"end_line,0"`
	Language  string
	Format    string `maya:"format,code"`
}

func newCmdView(args *cmdArgs) cmd {
	c := &cmdView{}
	fillCmd(c, args)
	defaultLang := strings.Replace(filepath.Ext(c.FilePath), ".", "", -1)
	c.Language = args.stringVal("lang", defaultLang)
	return c
}

func (c *cmdView) output() []string {
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

func (c *cmdView) execute() string {
	f := newFormatter(c.Format)
	return f.format(c.output(), c.Language)
}
