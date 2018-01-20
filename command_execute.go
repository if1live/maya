package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

func (c *CommandExecute) cacheFileName() string {
	data := []byte(c.Cmd)
	return fmt.Sprintf("%x.txt", md5.Sum(data))
}

func (c *CommandExecute) cacheDir() string {
	// 실행 경로를 캐시 생성 경로로 이용
	pwd, _ := os.Getwd()
	cachePath := filepath.Join(pwd, "cache")
	return cachePath
}

func (c *CommandExecute) cacheFilePath() string {
	dir := c.cacheDir()
	filename := c.cacheFileName()
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, filename)
}

func (c *CommandExecute) cacheExists() bool {
	_, err := os.Stat(c.cacheFilePath())
	return !os.IsNotExist(err)
}

func (c *CommandExecute) readCache() []string {
	data, _ := ioutil.ReadFile(c.cacheFilePath())
	text := string(data[:])
	lines := strings.Split(text, "\n")

	retval := []string{}
	for _, line := range lines {
		if !strings.HasPrefix(line, "# ") {
			retval = append(retval, line)
		}
	}
	return retval
}

func (c *CommandExecute) writeCache(lines []string) bool {
	cacheLines := []string{
		"# " + c.Cmd,
	}
	cacheLines = append(cacheLines, lines...)
	data := []byte(strings.Join(cacheLines, "\n"))
	ioutil.WriteFile(c.cacheFilePath(), data, 0644)
	return true
}

func (c *CommandExecute) RawOutput() []string {
	outputLines := []string{}
	if c.cacheExists() {
		outputLines = c.readCache()
	} else {
		outputLines = c.ExecuteImmediately()
		c.writeCache(outputLines)
	}

	elems := []string{}
	if c.AttachCmd {
		elems = append(elems, "$ "+c.Cmd)
	}
	elems = append(elems, outputLines...)
	elems = sanitizeLineFeedMultiLine(elems)
	return elems
}

func (c *CommandExecute) executeImmediatelyUnix() []string {
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

func (c *CommandExecute) executeImmediatelyWindows() []string {
	// https://groups.google.com/forum/#!topic/golang-nuts/Qtaw8r3Sx68
	out, err := exec.Command("cmd", "/c", c.Cmd).CombinedOutput()
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

func (c *CommandExecute) ExecuteImmediately() []string {
	log := logging.MustGetLogger("maya")
	log.Infof("Command execute: %v", c)

	switch runtime.GOOS {
	case "windows":
		return c.executeImmediatelyWindows()
	default:
		return c.executeImmediatelyUnix()
	}
}

func (c *CommandExecute) Formatter() *OutputFormatter {
	return &OutputFormatter{c.Format}
}

func (c *CommandExecute) Execute() string {
	formatter := c.Formatter()
	return formatter.Format(c.RawOutput(), "bash")
}
