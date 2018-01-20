package maya

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArticle(t *testing.T) {
	cases := []struct {
		text         string
		metadataText string
		contentText  string
	}{
		// basic
		{
			strings.Join([]string{
				"+++",
				"title: hello",
				"+++",
				"this is content",
			}, "\n"),
			"title: hello",
			"this is content",
		},
		// metadata
		{
			strings.Join([]string{
				"+++",
				"title: hello",
				"+++",
				"",
			}, "\n"),
			"title: hello",
			"",
		},
		// content only
		{
			strings.Join([]string{
				"this is content",
			}, "\n"),
			"",
			"this is content",
		},
		{
			strings.Join([]string{
				"+++",
				"+++",
				"this is content",
			}, "\n"),
			"",
			"this is content",
		},
	}

	for _, c := range cases {
		article := NewArticle(c.text, ModePelican)
		assert.Equal(t, c.metadataText, article.MetadataText)
		assert.Equal(t, c.contentText, article.ContentText)
	}
}
