package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/kardianos/osext"
	"github.com/op/go-logging"
)

const (
	ModePelican = "pelican"
	ModeHugo    = "hugo"
	ModeEmpty   = "empty"
)

type MetadataKeyValue struct {
	Key   string
	Value string
}

func (kv *MetadataKeyValue) IsListValue() bool {
	re := regexp.MustCompile(`^\[.*\]$`)
	return re.FindString(kv.Value) != ""
}

func (kv *MetadataKeyValue) ListValue() []string {
	re := regexp.MustCompile(`^\[(.*)\]$`)
	m := re.FindStringSubmatch(kv.Value)
	if len(m) == 0 {
		return []string{}
	}

	vals := strings.Split(m[1], ",")
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

type ArticleMetadata struct {
	Table []MetadataKeyValue
}

type MetadataTemplateLoader struct {
	texts     map[string]string
	templates map[string]*template.Template
}

func NewMetadata(text string) *ArticleMetadata {
	table := []MetadataKeyValue{}

	lines := strings.Split(text, "\n")
	re := regexp.MustCompile(`(\w+)\s*:(.*)`)
	for _, line := range lines {
		m := re.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		key, value := m[1], m[2]
		value = strings.Trim(value, " ")
		key = strings.ToLower(key)

		t := MetadataKeyValue{key, value}
		table = append(table, t)
	}

	return &ArticleMetadata{
		Table: table,
	}
}

func (m *ArticleMetadata) Get(key string) string {
	for _, t := range m.Table {
		if t.Key == key {
			return t.Value
		}
	}
	return ""
}

func (m *ArticleMetadata) GetList(key string) []string {
	for _, t := range m.Table {
		if t.Key == key {
			return t.ListValue()
		}
	}
	return []string{}
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

	const packageName = "github.com/if1live/maya"
	executableDir, _ := osext.ExecutableFolder()

	for _, target := range targets {
		candidatePaths := []string{
			".",
			executableDir,
			filepath.Join(os.Getenv("GOPATH"), "src", packageName),
		}

		found := false
		for _, path := range candidatePaths {
			candidate := filepath.Join(path, target.filepath)
			if loader.Register(target.mode, candidate) {
				found = true
				break
			}
		}
		if found == false {
			log := logging.MustGetLogger("maya")
			log.Fatalf("Metadata Template Load Fail [%s] cannot find any candidate", target.mode)
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

	log := logging.MustGetLogger("maya")
	log.Infof("Metadata Template Load Success [%s] %s", mode, filepath)
	return true
}

func makeSeperator(text string, sep string) string {
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
}

func isString(x interface{}) bool {
	_, ok := x.(string)
	return ok
}

func escape(val string) string {
	table := []struct {
		in  string
		out string
	}{
		{`"`, `\"`},
	}
	for _, t := range table {
		val = strings.Replace(val, t.in, t.out, -1)
	}
	return val
}

func (l *MetadataTemplateLoader) createFuncMap() template.FuncMap {
	return template.FuncMap{
		"title":     strings.Title,
		"join":      strings.Join,
		"seperator": makeSeperator,
		"isString":  isString,
		"escape":    escape,
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
