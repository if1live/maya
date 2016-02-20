package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var mode string
var filePath string

func init() {
	flag.StringVar(&mode, "mode", "pelican-md", "document mode")
	flag.StringVar(&filePath, "file", "demo.md", "file path")
}

func main() {
	flag.Parse()

	loader := NewTemplateLoader()

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	text := string(data[:])

	article := NewArticle(text)
	metadata := article.Metadata()
	header := loader.Execute(metadata, mode)

	content := article.Content()
	body := content.String()

	output := strings.Join([]string{header, "", body}, "\n")
	fmt.Println(output)

}
