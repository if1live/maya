package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

func (c *CommandExecute) CacheFileName() string {
	// 실행 경로를 캐시 생성 경로로 이용
	pwd, _ := os.Getwd()

	cachePath := filepath.Join(pwd, "cache")
	os.MkdirAll(cachePath, 0755)

	data := []byte(c.Cmd)
	filename := fmt.Sprintf("%x.txt", md5.Sum(data))

	cacheFile := filepath.Join(cachePath, filename)
	return cacheFile
}

func (c *CommandExecute) RawOutput() []string {
	cacheFilePath := c.CacheFileName()
	_, err := os.Stat(cacheFilePath)

	outputLines := []string{}
	if os.IsNotExist(err) {
		outputLines = c.ExecuteImmediately()
		data := []byte(strings.Join(outputLines, "\n"))
		ioutil.WriteFile(cacheFilePath, data, 0644)

	} else {
		data, _ := ioutil.ReadFile(cacheFilePath)
		text := string(data[:])
		outputLines = strings.Split(text, "\n")
	}

	elems := []string{}
	if c.AttachCmd {
		elems = append(elems, "$ "+c.Cmd)
	}
	elems = append(elems, outputLines...)
	return elems
}

func (c *CommandExecute) ExecuteImmediately() []string {
	log := logging.MustGetLogger("maya")
	log.Infof("Command execute: %v", c)

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

	elems := []string{}
	if err != nil {
		if _, ok := err.(*exec.Error); ok {
			elems = append(elems, err.Error())
			return elems
		}
	}

	elems = strings.Split(string(out[:]), "\n")
	return elems
}

func (c *CommandExecute) Formatter() *OutputFormatter {
	return &OutputFormatter{c.Format}
}

func (c *CommandExecute) Execute() string {
	formatter := c.Formatter()
	return formatter.Format(c.RawOutput())
}
