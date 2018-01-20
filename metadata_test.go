package maya

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsListValue(t *testing.T) {
	cases := []struct {
		value    string
		expected bool
	}{
		{"[a, b, c]", true},
		{"a,b,c", false},
		{"[]", true},
	}
	for _, c := range cases {
		kv := MetadataKeyValue{"key", c.value}
		assert.Equal(t, c.expected, kv.IsListValue(), c.value)
	}
}

func TestListValue(t *testing.T) {
	cases := []struct {
		value    string
		expected []string
	}{
		{"[this is title]", []string{"this is title"}},
		{"[1, 2, 3]", []string{"1", "2", "3"}},
		{"[1,2,1]", []string{"1", "2"}},
	}
	for _, c := range cases {
		kv := MetadataKeyValue{"key", c.value}
		assert.Equal(t, c.expected, kv.ListValue())
	}
}

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
		{"title : foo : bar : spam", "title", "foo : bar : spam"},
		{"title :  strip   ", "title", "strip"},
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
		{"tags: [1,2,3]", "tags", []string{"1", "2", "3"}},
		{"tags: [1, 2, 3]", "not-exist", []string{}},
	}
	for _, c := range cases {
		m := NewMetadata(c.line)
		assert.Equal(t, c.val, m.GetList(c.key))
	}
}

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
Title: "제목"
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
title = "\"제목\""
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
