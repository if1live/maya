package maya

import (
	"fmt"
)

type cmdGist struct {
	ID   string `maya:"id"`
	File string `maya:"file"`
}

func newCmdGist(action string, args *cmdArgs) cmd {
	return autoFillCmd(&cmdGist{}, args)
}

func (c *cmdGist) RawOutput() []string {
	url := fmt.Sprintf("https://gist.github.com/%s.js", c.ID)
	if c.File != "" {
		url += fmt.Sprintf("?file=%s", c.File)
	}
	return []string{
		`<div class="maya-gist">`,
		fmt.Sprintf(`<script src="%s"></script>`, url),
		`</div>`,
	}
}

func (c *cmdGist) execute() string {
	f := newFormatter(formatText)
	return f.format(c.RawOutput())
}
