package maya

import "fmt"

type cmdYoutube struct {
	VideoId string `maya:"video_id"`
	Width   int    `maya:"width,640"`
	Height  int    `maya:"height,480"`
}

func newCmdYoutube(args *cmdArgs) cmd {
	return fillCmd(&cmdYoutube{}, args)
}

func (c *cmdYoutube) output() []string {
	return []string{
		`<div class="maya-youtube">`,
		fmt.Sprintf(`<iframe width="%d" height="%d" src="//www.youtube.com/embed/%s" frameborder="0" allowfullscreen></iframe>`, c.Width, c.Height, c.VideoId),
		`</div>`,
	}
}

func (c *cmdYoutube) execute() string {
	f := newFormatter(formatText)
	return f.format(c.output())
}
