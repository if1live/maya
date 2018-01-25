package maya

import (
	"reflect"
	"strconv"
	"strings"
)

type cmd interface {
	output() []string
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
	type CreateFunc func(*cmdArgs) cmd

	table := map[string]CreateFunc{
		"view":    newCmdView,
		"execute": newCmdExecute,
		"youtube": newCmdYoutube,
		"gist":    newCmdGist,
	}
	if fn, ok := table[action]; ok {
		return fn(args)
	}
	return newCmdUnknown(action, args)
}

func fillCmd(c cmd, args *cmdArgs) cmd {
	t := reflect.TypeOf(c).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("maya")
		tokens := strings.Split(tag, ",")

		keyIdx := 0
		defaultValIdx := 1

		key := tokens[keyIdx]

		switch field.Type {
		case reflect.TypeOf(""):
			defaultVal := ""
			if len(tokens) >= defaultValIdx+1 {
				defaultVal = tokens[defaultValIdx]
			}

			v := args.stringVal(key, defaultVal)
			reflect.ValueOf(c).Elem().Field(i).SetString(v)
			break

		case reflect.TypeOf(1):
			defaultVal := 0
			if len(tokens) >= defaultValIdx+1 {
				defaultVal, _ = strconv.Atoi(tokens[defaultValIdx])
			}
			v := args.intVal(key, defaultVal)
			reflect.ValueOf(c).Elem().Field(i).SetInt(int64(v))
			break

		case reflect.TypeOf(true):
			defaultVal := false
			if len(tokens) >= defaultValIdx+1 {
				v := tokens[defaultValIdx]
				if v == "true" {
					defaultVal = true
				} else if v == "false" {
					defaultVal = false
				}
			}
			v := args.boolVal(key, defaultVal)
			reflect.ValueOf(c).Elem().Field(i).SetBool(v)
			break
		}
	}
	return c
}
