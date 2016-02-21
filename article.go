package main

import (
	"bytes"
	"io"
	"regexp"
	"strings"

	"github.com/op/go-logging"
)

type Article struct {
	MetadataText string
	ContentText  string
	MetadataMode string
	loader       MetadataTemplateLoader
}

func NewArticleFromReader(r io.Reader, mode string) *Article {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return NewArticle(buf.String(), mode)
}

func NewArticle(text string, mode string) *Article {
	metadataLines := []string{}
	contentLines := []string{}

	if mode == "" {
		log := logging.MustGetLogger("maya")
		log.Fatal("mode required. use -h")
	}

	lines := strings.Split(text, "\n")

	const (
		LineParseStateMetadata = 1
		LineParseStateContent  = 2
	)

	state := LineParseStateMetadata
	re := regexp.MustCompile(`^(.+):(.*)$`)
	for _, line := range lines {
		switch state {
		case LineParseStateMetadata:
			if len(strings.Trim(line, " \t")) == 0 {
				continue
			}
			m := re.FindString(line)
			if m == "" {
				state = LineParseStateContent
				contentLines = append(contentLines, line)
			} else {
				metadataLines = append(metadataLines, line)
			}
		case LineParseStateContent:
			contentLines = append(contentLines, line)
		}
	}

	return &Article{
		MetadataText: strings.Join(metadataLines, "\n"),
		ContentText:  strings.Join(contentLines, "\n"),
		MetadataMode: mode,
		loader:       NewTemplateLoader(),
	}
}

func (a *Article) Metadata() *ArticleMetadata {
	return NewMetadata(a.MetadataText)
}

func (a *Article) Content() *ArticleContent {
	return NewContent(a.ContentText)
}

func (a *Article) Output(w io.Writer) {
	metadata := a.Metadata()
	header := a.loader.Execute(metadata, a.MetadataMode)

	content := a.Content()
	body := content.String()

	output := strings.Join([]string{header, "", body}, "\n")
	output = strings.TrimLeft(output, "\n")
	w.Write([]byte(output))
}
