package maya

import "fmt"

type CommandGist struct {
	Id   string
	File string
}

func NewCommandGist(action string, args *CommandArguments) Command {
	return &CommandGist{
		Id:   args.StringVal("id", ""),
		File: args.StringVal("file", ""),
	}
}

func (c *CommandGist) RawOutput() []string {
	url := fmt.Sprintf("https://gist.github.com/%s.js", c.Id)
	if c.File != "" {
		url += fmt.Sprintf("?file=%s", c.File)
	}
	return []string{
		`<div class="maya-gist">`,
		fmt.Sprintf(`<script src="%s"></script>`, url),
		`</div>`,
	}
}

func (c *CommandGist) execute() string {
	f := newFormatter(formatText)
	return f.format(c.RawOutput())
}
