package maya

import (
	"reflect"
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

func autoFillCmd(c cmd, args *cmdArgs) cmd {
	t := reflect.TypeOf(c).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("maya")
		tokens := strings.Split(tag, ",")

		key := tokens[0]
		switch field.Type {
		case reflect.TypeOf(""):
			defaultValue := ""
			if len(tokens) >= 2 {
				defaultValue = tokens[1]
			}

			v := args.stringVal(key, defaultValue)
			reflect.ValueOf(c).Elem().Field(i).SetString(v)
			break

		case reflect.TypeOf(1):
			defaultValue := 0
			if len(tokens) >= 2 {
				defaultValue, _ = strconv.Atoi(tokens[1])
			}
			v := args.intVal(key, defaultValue)
			reflect.ValueOf(c).Elem().Field(i).SetInt(int64(v))
			break
		}
	}
	return c
}
