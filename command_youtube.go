package maya

import "fmt"

type CommandYoutube struct {
	VideoId string
	Width   int
	Height  int
}

func NewCommandYoutube(action string, args *CommandArguments) Command {
	return &CommandYoutube{
		VideoId: args.StringVal("video_id", ""),
		Width:   args.IntVal("width", 640),
		Height:  args.IntVal("height", 480),
	}
}

func (c *CommandYoutube) RawOutput() []string {
	return []string{
		`<div class="maya-youtube">`,
		fmt.Sprintf(`<iframe width="%d" height="%d" src="//www.youtube.com/embed/%s" frameborder="0" allowfullscreen></iframe>`, c.Width, c.Height, c.VideoId),
		`</div>`,
	}
}

func (c *CommandYoutube) execute() string {
	f := newFormatter(formatText)
	return f.format(c.RawOutput())
}
