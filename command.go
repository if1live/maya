package maya

import (
	"strconv"
)

type Command interface {
	RawOutput() []string
	execute() string
}

type CommandArguments struct {
	params map[string]string
}

func (ca *CommandArguments) IntVal(key string, defaultVal int) int {
	val, err := strconv.Atoi(ca.params[key])
	if err != nil {
		return defaultVal
	}
	return val
}

func (ca *CommandArguments) StringVal(key string, defaultVal string) string {
	if val, ok := ca.params[key]; ok {
		return val
	}
	return defaultVal
}

func (ca *CommandArguments) BoolVal(key string, defaultVal bool) bool {
	if val, ok := ca.params[key]; ok {
		return !(val == "false" || val == "f")
	}
	return defaultVal
}

func NewCommand(action string, args *CommandArguments) Command {
	type CommandCreateFunc func(string, *CommandArguments) Command

	table := map[string]CommandCreateFunc{
		"view":    NewCommandView,
		"execute": NewCommandExecute,
		"youtube": NewCommandYoutube,
		"gist":    NewCommandGist,
	}
	if fn, ok := table[action]; ok {
		return fn(action, args)
	}
	return NewCommandUnknown(action, args)
}
