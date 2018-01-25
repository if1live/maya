package maya

import "fmt"

type cmdYoutube struct {
	videoId string
	width   int
	height  int
}

func newCommandYoutube(action string, args *cmdArgs) cmd {
	return &cmdYoutube{
		videoId: args.stringVal("video_id", ""),
		width:   args.intVal("width", 640),
		height:  args.intVal("height", 480),
	}
}

func (c *cmdYoutube) RawOutput() []string {
	return []string{
		`<div class="maya-youtube">`,
		fmt.Sprintf(`<iframe width="%d" height="%d" src="//www.youtube.com/embed/%s" frameborder="0" allowfullscreen></iframe>`, c.width, c.height, c.videoId),
		`</div>`,
	}
}

func (c *cmdYoutube) execute() string {
	f := newFormatter(formatText)
	return f.format(c.RawOutput())
}
