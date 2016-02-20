package main

import (
	"regexp"
	"strings"
)

type Article struct {
	MetadataText string
	ContentText  string
}

func NewArticle(text string) *Article {
	metadataLines := []string{}
	contentLines := []string{}
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
	}
}

func (a *Article) Metadata() *ArticleMetadata {
	return NewMetadata(a.MetadataText)
}

func (a *Article) Content() *ArticleContent {
	return NewContent(a.ContentText)
}
