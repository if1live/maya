package main

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
func (c *CommandUnknown) Formatter() *OutputFormatter {
	return &OutputFormatter{OutputFormatBlockquote}
}

func (c *CommandUnknown) Execute() string {
	formatter := c.Formatter()
	return formatter.Format(c.RawOutput())
}
