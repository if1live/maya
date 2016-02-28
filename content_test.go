package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContent(t *testing.T) {
	cases := []struct {
		text   string
		blocks []ContentBlock
	}{
		{
			strings.Trim(`
hello
world
`, "\n"),
			[]ContentBlock{
				{"", []string{"hello", "world"}},
			},
		},
		{
			strings.Trim(`
hello
~~~maya:view
file=x.py
~~~
world
`, "\n"),
			[]ContentBlock{
				{"", []string{"hello"}},
				{"view", []string{"~~~maya:view", "file=x.py", "~~~"}},
				{"", []string{"world"}},
			},
		},
		{
			strings.Trim(`
hello
~~~maya:view
file=x.py
`, "\n"),
			[]ContentBlock{
				{"", []string{"hello"}},
				{"view", []string{"~~~maya:view", "file=x.py"}},
			},
		},
	}
	for _, c := range cases {
		content := NewContent(c.text)
		assert.Equal(t, c.blocks, content.blocks)
	}
}
