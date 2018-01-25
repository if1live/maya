package maya

import "github.com/op/go-logging"

type cmdUnknown struct {
	Action string
}

func newCmdUnknown(action string, args *cmdArgs) cmd {
	return &cmdUnknown{
		Action: action,
	}
}

func (c *cmdUnknown) output() []string {
	log := logging.MustGetLogger("maya")
	log.Warningf("Command Unknown: %v", c)
	tokens := []string{
		"Action=" + c.Action,
	}
	return tokens
}

func (c *cmdUnknown) execute() string {
	f := newFormatter(formatBlockquote)
	return f.format(c.output())
}
