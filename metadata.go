package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/kardianos/osext"
)

const (
	ModePelican = "pelican"
	ModeHugo    = "hugo"
	ModeEmpty   = "empty"
)

type ArticleMetadata struct {
	table map[string]string
}

type MetadataTemplateLoader struct {
	texts     map[string]string
	templates map[string]*template.Template
}

func NewMetadata(text string) *ArticleMetadata {
	table := make(map[string]string)

	lines := strings.Split(text, "\n")
	re := regexp.MustCompile(`(.+):(.*)`)
	for _, line := range lines {
		m := re.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		key, value := m[1], m[2]
		key = strings.Trim(key, " ")
		value = strings.Trim(value, " ")
		key = strings.ToLower(key)

		table[key] = value
	}

	return &ArticleMetadata{
		table: table,
	}
}

func (m *ArticleMetadata) Get(key string) string {
	return m.table[key]
}

func (m *ArticleMetadata) GetList(key string) []string {
	vals := strings.Split(m.Get(key), ",")
	for i, val := range vals {
		vals[i] = strings.Trim(val, " ")
	}

	keys := map[string]bool{}
	result := vals[:0]
	for _, val := range vals {
		if len(val) > 0 && keys[val] == false {
			keys[val] = true
			result = append(result, val)
		}
	}
	return result
}

func NewTemplateLoader() MetadataTemplateLoader {
	loader := MetadataTemplateLoader{
		texts:     map[string]string{},
		templates: map[string]*template.Template{},
	}

	targets := []struct {
		mode     string
		filepath string
	}{
		{ModePelican, "templates/metadata_pelican.tpl"},
		{ModeHugo, "templates/metadata_hugo.tpl"},
		{ModeEmpty, "templates/metadata_empty.tpl"},
	}
	for _, target := range targets {
		dir, _ := osext.ExecutableFolder()
		candidates := []string{
			target.filepath,
			filepath.Join(dir, target.filepath),
		}
		for _, candidate := range candidates {
			success := loader.Register(target.mode, candidate)
			if success {
				break
			}
		}
	}

	return loader
}

func (l *MetadataTemplateLoader) Register(mode, filepath string) bool {
	funcMap := l.createFuncMap()
	text, err := l.readFile(filepath)
	if err != nil {
		return false
	}
	l.texts[mode] = text
	l.templates[mode] = template.Must(
		template.New(mode).Funcs(funcMap).Parse(l.texts[mode]),
	)
	log.Printf("Metadata Template Load Success [%s] %s", mode, filepath)
	return true

}

func (l *MetadataTemplateLoader) createFuncMap() template.FuncMap {
	return template.FuncMap{
		"title": strings.Title,
		"join":  strings.Join,
		"seperator": func(text string, sep string) string {
			count := 0

			for i, w := 0, 0; i < len(text); i += w {
				_, width := utf8.DecodeRuneInString(text[i:])
				if width == 1 {
					count += 1
				} else {
					count += 2
				}
				w = width
			}

			tokens := make([]string, count)
			for i := 0; i < count; i++ {
				tokens[i] = sep
			}
			return strings.Join(tokens, "")
		},
		"isString": func(x interface{}) bool {
			_, ok := x.(string)
			return ok
		},
	}
}

func (l *MetadataTemplateLoader) Execute(metadata *ArticleMetadata, mode string) string {
	t := l.templates[mode]
	if t == nil {
		msg := "Unknown document mode : " + mode
		panic(msg)
	}
	var b bytes.Buffer
	t.Execute(&b, metadata)
	text := string(b.Bytes())
	lines := strings.Split(text, "\n")

	result := []string{}
	for _, line := range lines {
		if len(strings.Trim(line, " ")) > 0 {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

func (l *MetadataTemplateLoader) readFile(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(data[:]), nil
}
