package maya

import "fmt"

type cmdGist struct {
	id   string
	file string
}

func newCmdGist(action string, args *cmdArgs) cmd {
	return &cmdGist{
		id:   args.stringVal("id", ""),
		file: args.stringVal("file", ""),
	}
}

func (c *cmdGist) RawOutput() []string {
	url := fmt.Sprintf("https://gist.github.com/%s.js", c.id)
	if c.file != "" {
		url += fmt.Sprintf("?file=%s", c.file)
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
