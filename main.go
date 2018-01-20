package main

import (
	"flag"
	"os"

	"github.com/op/go-logging"
)

var _mode string
var _filePath string
var _logLevel string
var _outputPath string

func init() {
	flag.StringVar(&_mode, "mode", "", "document mode: pelican/hugo")
	flag.StringVar(&_filePath, "file", "", "file path: xxx.md")
	flag.StringVar(&_logLevel, "log", "ERROR", "log level: critical, error, warning, notice, info, debug")
	flag.StringVar(&_outputPath, "output", "stdout", "output path: xxx.md")
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

	infile, err := os.Open(_filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	outfile := os.Stdout
	if _outputPath != "stdout" {
		outfile, err = os.Create(_outputPath)
		if err != nil {
			panic(err)
		}
		defer outfile.Close()
	}

	article := NewArticleFromReader(infile, _mode)
	article.Output(outfile)
}
