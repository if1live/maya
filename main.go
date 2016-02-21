package main

import (
	"flag"
	"log"
	"os"
)

var mode string
var filePath string

func init() {
	flag.StringVar(&mode, "mode", "", "document mode: pelican/hugo")
	flag.StringVar(&filePath, "file", "", "file path: xxx.md")
}

func main() {
	flag.Parse()

	if filePath == "" {
		log.Fatal("file path required. use -h")
		os.Exit(1)
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	article := NewArticleFromReader(f, mode)
	out := os.Stdout
	article.Output(out)
}
