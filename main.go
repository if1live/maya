package main

import (
	"flag"
	"os"

	"github.com/op/go-logging"
)

var _mode string
var _filePath string
var _logLevel string

func init() {
	flag.StringVar(&_mode, "mode", "", "document mode: pelican/hugo")
	flag.StringVar(&_filePath, "file", "", "file path: xxx.md")
	flag.StringVar(&_logLevel, "log", "ERROR", "log level: critical, error, warning, notice, info, debug")
}

var _formatter = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	flag.Parse()

	logLevel, _ := logging.LogLevel(_logLevel)
	logging.SetLevel(logLevel, "maya")
	logging.SetFormatter(_formatter)

	log := logging.MustGetLogger("maya")
	if _filePath == "" {
		log.Fatal("file path required. use -h")
	}

	f, err := os.Open(_filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	article := NewArticleFromReader(f, _mode)
	out := os.Stdout
	article.Output(out)
}
