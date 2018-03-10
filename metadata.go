package maya

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/op/go-logging"
	yaml "gopkg.in/yaml.v2"
)

const (
	ModePelican = "pelican"
	ModeHugo    = "hugo"
	ModeEmpty   = "empty"
)

const templateHugo = `
+++
{{range $_,$elem := .Table}}
{{if $elem.IsListValue}}
{{$elem.Key}} = ["{{join $elem.ListValue "\", \""}}"]
{{else}}
{{$elem.Key}} = "{{escape $elem.Value}}"
{{end}}
{{end}}
+++
`

const templatePelican = `
{{range $_,$elem := .Table}}
{{if $elem.IsListValue}}
{{title $elem.Key}}: {{join $elem.ListValue ", "}}
{{else}}
{{title $elem.Key}}: {{$elem.Value}}
{{end}}
{{end}}
`

func getTemplateHugo() string {
	return strings.TrimLeft(templateHugo, " ")
}
func getTemplatePelican() string {
	return strings.TrimLeft(templatePelican, " ")
}
func getTemplateEmpty() string {
	return ""
}

type MetadataKeyValue struct {
	Key string

	isList    bool
	singleVal string
	multiVal  []string
}

func (kv *MetadataKeyValue) IsListValue() bool {
	return kv.isList
}

func (kv *MetadataKeyValue) Value() string {
	return kv.singleVal
}

func (kv *MetadataKeyValue) ListValue() []string {
	return kv.multiVal
}

type ArticleMetadata struct {
	Table []MetadataKeyValue
}

type MetadataTemplateLoader struct {
	texts     map[string]string
	templates map[string]*template.Template
}

func NewMetadata(text string) *ArticleMetadata {
	m := yaml.MapSlice{}
	err := yaml.Unmarshal([]byte(text), &m)
	if err != nil {
		panic(err)
	}
	dict := NewDict(m)
	keys := dict.GetStrKeys()

	table := []MetadataKeyValue{}
	for _, key := range keys {
		pair := MetadataKeyValue{
			Key: key,
		}
		switch dict.GetValueType(key) {
		case valueTypeStr:
			pair.isList = false
			val, _ := dict.GetStr(key)
			pair.singleVal = val
			break

		case valueTypeInt:
			pair.isList = false
			tmp, _ := dict.GetInt(key)
			str := strconv.Itoa(tmp)
			pair.singleVal = str
			break

		case valueTypeStrList:
			pair.isList = true
			val, _ := dict.GetStrList(key)
			pair.multiVal = val
			break

		case valueTypeIntList:
			pair.isList = true
			tmp, _ := dict.GetIntList(key)
			val := make([]string, len(tmp))
			for i, num := range tmp {
				val[i] = strconv.Itoa(num)
			}
			pair.multiVal = val
		}
		table = append(table, pair)
	}

	return &ArticleMetadata{
		Table: table,
	}
}

func (m *ArticleMetadata) Preprocess(mode string) {
	type Func func(*ArticleMetadata)
	funcs := map[string]Func{
		"hugo": preprocessHugo,
	}
	if fn, ok := funcs[mode]; ok {
		fn(m)
	}
}

func preprocessHugo(m *ArticleMetadata) {
	type Func func(string) string
	funcs := map[string]Func{
		"date": preprocessHugo_date,
	}

	for i, t := range m.Table {
		if fn, ok := funcs[t.Key]; ok {
			t.singleVal = fn(t.singleVal)
			m.Table[i] = t
		}
	}
}

func preprocessHugo_date(val string) string {
	// to support date-only format
	// YYYY-MM-DD => YYYY-MM-DDT00:00:00+00:00
	dateRe := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if dateRe.MatchString(val) {
		return val + "T00:00:00+00:00"
	}
	return val
}

func NewTemplateLoader() MetadataTemplateLoader {
	loader := MetadataTemplateLoader{
		texts:     map[string]string{},
		templates: map[string]*template.Template{},
	}

	targets := []struct {
		mode string
		text string
	}{
		{ModePelican, getTemplatePelican()},
		{ModeHugo, getTemplateHugo()},
		{ModeEmpty, getTemplateEmpty()},
	}

	for _, target := range targets {
		loader.RegisterTemplate(target.mode, target.text)
	}

	return loader
}

func (l *MetadataTemplateLoader) RegisterFile(mode, filepath string) bool {
	text, err := l.readFile(filepath)
	if err != nil {
		return false
	}
	ok := l.RegisterTemplate(mode, text)

	log := logging.MustGetLogger("maya")
	log.Infof("Metadata Template Load Success [%s] %s", mode, filepath)
	return ok
}

func (l *MetadataTemplateLoader) RegisterTemplate(mode, text string) bool {
	funcMap := l.createFuncMap()
	l.texts[mode] = text
	l.templates[mode] = template.Must(
		template.New(mode).Funcs(funcMap).Parse(l.texts[mode]),
	)
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
	metadataClone := *metadata
	metadataClone.Preprocess(mode)

	t := l.templates[mode]
	if t == nil {
		msg := "Unknown document mode : " + mode
		panic(msg)
	}
	var b bytes.Buffer
	t.Execute(&b, &metadataClone)
	text := string(b.Bytes())
	lines := strings.Split(text, "\n")

	result := []string{}
	for _, line := range lines {
		line = strings.Replace(line, "\r", "", -1)
		line = strings.Trim(line, " ")
		if len(line) > 0 {
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
