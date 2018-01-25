package maya

import (
	"strconv"
	"strings"
)

type cmd interface {
	RawOutput() []string
	execute() string
}

type cmdArgs struct {
	params map[string]string
}

func (args *cmdArgs) intVal(key string, defaultVal int) int {
	val, err := strconv.Atoi(args.params[key])
	if err != nil {
		return defaultVal
	}
	return val
}

func (args *cmdArgs) stringVal(key string, defaultVal string) string {
	if val, ok := args.params[key]; ok {
		return val
	}
	return defaultVal
}

func (args *cmdArgs) boolVal(key string, defaultVal bool) bool {
	trueStr := []string{
		"true",
		"t",
	}

	falseStr := []string{
		"false",
		"f",
	}

	if val, ok := args.params[key]; ok {
		v := strings.ToLower(val)
		for _, s := range trueStr {
			if v == strings.ToLower(s) {
				return true
			}
		}
		for _, s := range falseStr {
			if v == strings.ToLower(s) {
				return false
			}
		}
	}
	return defaultVal
}

func newCmd(action string, args *cmdArgs) cmd {
	type CreateFunc func(string, *cmdArgs) cmd

	table := map[string]CreateFunc{
		"view":    newCmdView,
		"execute": newCmdExecute,
		"youtube": newCommandYoutube,
		"gist":    newCmdGist,
	}
	if fn, ok := table[action]; ok {
		return fn(action, args)
	}
	return newCmdUnknown(action, args)
}
