package maya

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	metadataText := `
title: "제목"
subtitle: subtitle-1
date: 2016-02-20
tags: [foo, bar]
slug: slug-1
status: draft
`
	metadata := NewMetadata(metadataText)
	loader := NewTemplateLoader()

	cases := []struct {
		actual   string
		expected string
	}{
		{
			loader.Execute(metadata, ModePelican),
			strings.Trim(`
Title: 제목
Subtitle: subtitle-1
Date: 2016-02-20
Tags: foo, bar
Slug: slug-1
Status: draft
`, "\n"),
		},
		{
			loader.Execute(metadata, ModeHugo),
			strings.Trim(`
+++
title = "제목"
subtitle = "subtitle-1"
date = "2016-02-20T00:00:00+00:00"
tags = ["foo", "bar"]
slug = "slug-1"
status = "draft"
+++
`, "\n"),
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual)
	}
}

func Test_makeSeperator(t *testing.T) {
	cases := []struct {
		input    string
		sep      string
		expected string
	}{
		{"abc", "=", "==="},
		{"abc", "-", "---"},
		{"", "-", ""},
		{"한글12", "=", "======"},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, makeSeperator(c.input, c.sep))
	}
}

func Test_isString(t *testing.T) {
	cases := []struct {
		input    interface{}
		expected bool
	}{
		{"hello", true},
		{123, false},
		{nil, false},
		{[]bool{}, false},
		{map[string]string{}, false},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, isString(c.input))
	}
}

func Test_escape(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{`"hello"`, `\"hello\"`},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, escape(c.input))
	}
}
