package maya

import "fmt"

type cmdYoutube struct {
	VideoId string `maya:"video_id"`
	Width   int    `maya:"width,640"`
	Height  int    `maya:"height,480"`
}

func newCommandYoutube(action string, args *cmdArgs) cmd {
	return autoFillCmd(&cmdYoutube{}, args)
}

func (c *cmdYoutube) RawOutput() []string {
	return []string{
		`<div class="maya-youtube">`,
		fmt.Sprintf(`<iframe width="%d" height="%d" src="//www.youtube.com/embed/%s" frameborder="0" allowfullscreen></iframe>`, c.Width, c.Height, c.VideoId),
		`</div>`,
	}
}

func (c *cmdYoutube) execute() string {
	f := newFormatter(formatText)
	return f.format(c.RawOutput())
}
