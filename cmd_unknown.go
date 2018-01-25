package maya

import "github.com/op/go-logging"

type cmdUnknown struct {
	Action string
	Args   *cmdArgs
}

func newCmdUnknown(action string, args *cmdArgs) cmd {
	return &cmdUnknown{
		Action: action,
		Args:   args,
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
