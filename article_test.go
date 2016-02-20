package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArticle(t *testing.T) {
	cases := []struct {
		text         string
		metadataText string
		contentText  string
	}{
		// metadata + empty line + content
		{
			"title: hello\n\nthis is content",
			"title: hello",
			"this is content",
		},
		// metadata + content
		{
			"title: hello\nthis is content",
			"title: hello",
			"this is content",
		},
		// metadata
		{
			"title: hello\n",
			"title: hello",
			"",
		},
		// content
		{
			"this is content\n",
			"",
			"this is content\n",
		},
	}

	for _, c := range cases {
		article := NewArticle(c.text)
		assert.Equal(t, c.metadataText, article.MetadataText)
		assert.Equal(t, c.contentText, article.ContentText)
	}
}
