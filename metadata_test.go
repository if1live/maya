package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	cases := []struct {
		line string
		key  string
		val  string
	}{
		{"title: foo", "title", "foo"},
		{"Subtitle : bar", "subtitle", "bar"},
		{"SLUG   :   spam", "slug", "spam"},
		{"title : this is title", "title", "this is title"},
		{"title : xxx", "", ""},
		{"title : xxx", "not-exist", ""},
	}
	for _, c := range cases {
		m := NewMetadata(c.line)
		assert.Equal(t, c.val, m.Get(c.key))
	}
}

func TestGetList(t *testing.T) {
	cases := []struct {
		line string
		key  string
		val  []string
	}{
		{"tags: this is title", "tags", []string{"this is title"}},
		{"tags: 1,2,3", "tags", []string{"1", "2", "3"}},
		{"tags: 1, 2, 3", "not-exist", []string{}},
		{"tags: 1,2,1", "tags", []string{"1", "2"}},
	}
	for _, c := range cases {
		m := NewMetadata(c.line)
		assert.Equal(t, c.val, m.GetList(c.key))
	}
}

func TestExecute(t *testing.T) {
	metadataText := `
title: 제목
subtitle: subtitle-1
date: 2016-02-20
tags: foo, bar
slug: slug-1
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
Slug: slug-1
Tags: foo, bar
Date: 2016-02-20
`, "\n"),
		},
		{
			loader.Execute(metadata, ModeHugo),
			strings.Trim(`
+++
title = "제목"
subtitle = "subtitle-1"
slug = "slug-1"
tags = ["foo", "bar"]
date = "2016-02-20"
+++
`, "\n"),
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual)
	}
}
