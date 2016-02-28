package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/op/go-logging"
)

type CommandExecute struct {
	Cmd       string
	AttachCmd bool
	Format    string
}

func NewCommandExecute(action string, args *CommandArguments) Command {
	return &CommandExecute{
		Cmd:       args.StringVal("cmd", "echo empty"),
		AttachCmd: args.BoolVal("attach_cmd", false),
		Format:    args.StringVal("format", OutputFormatCode),
	}
}

func (c *CommandExecute) RawOutput() []string {
	log := logging.MustGetLogger("maya")
	log.Infof("Command execute: %v", c)

	elems := []string{}
	if c.AttachCmd {
		elems = append(elems, "$ "+c.Cmd)
	}

	tmpfile, err := ioutil.TempFile("", "maya")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(c.Cmd)); err != nil {
		panic(err)
	}
	if err := tmpfile.Close(); err != nil {
		panic(err)
	}
	out, err := exec.Command("bash", tmpfile.Name()).CombinedOutput()

	if err != nil {
		if _, ok := err.(*exec.Error); ok {
			elems = append(elems, err.Error())
			return elems
		}
	}

	elems = append(elems, strings.Split(string(out[:]), "\n")...)
	return elems
}

func (c *CommandExecute) Formatter() *OutputFormatter {
	return &OutputFormatter{c.Format}
}

func (c *CommandExecute) Execute() string {
	formatter := c.Formatter()
	return formatter.Format(c.RawOutput())
}
