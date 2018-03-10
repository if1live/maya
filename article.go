package maya

import (
	"bytes"
	"io"
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
	text := buf.String()
	text = strings.Replace(text, "\r", "", -1)
	return NewArticle(text, mode)
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
		LineParseStateInit     = 0
		LineParseStateMetadata = 1
		LineParseStateContent  = 2
	)

	// github는 --- ~ --- 구역을 yaml로 파싱한다
	// ---
	// metadatas
	// ---
	// content

	firstLine := 0
	for i, line := range lines {
		if len(strings.Trim(line, " ")) > 0 {
			firstLine = i
			break
		}
	}

	if lines[firstLine] != "---" {
		// no metadata
		return &Article{
			MetadataText: "",
			ContentText:  text,
			MetadataMode: mode,
			loader:       NewTemplateLoader(),
		}
	}

	state := LineParseStateInit
	for i, line := range lines {
		if i < firstLine {
			continue
		}
		switch state {
		case LineParseStateInit:
			if strings.Trim(line, " ") == "---" {
				state = LineParseStateMetadata
				metadataLines = []string{}
			} else {
				contentLines = append(contentLines, line)
			}
		case LineParseStateMetadata:
			if strings.Trim(line, " ") == "---" {
				state = LineParseStateContent
				contentLines = []string{}
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
	output := a.OutputString()
	w.Write([]byte(output))
}

func (a *Article) OutputString() string {
	metadata := a.Metadata()
	header := a.loader.Execute(metadata, a.MetadataMode)

	content := a.Content()
	body := content.String()

	output := strings.Join([]string{header, "", body}, "\n")
	output = strings.TrimLeft(output, "\n")

	return output
}
