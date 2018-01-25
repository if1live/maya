package maya

import "github.com/op/go-logging"

type CommandUnknown struct {
	Action string
}

func NewCommandUnknown(action string, args *CommandArguments) Command {
	return &CommandUnknown{
		Action: action,
	}
}

func (c *CommandUnknown) RawOutput() []string {
	log := logging.MustGetLogger("maya")
	log.Warningf("Command Unknown: %v", c)
	tokens := []string{
		"Action=" + c.Action,
	}
	return tokens
}

func (c *CommandUnknown) execute() string {
	f := newFormatter(formatBlockquote)
	return f.format(c.RawOutput())
}
