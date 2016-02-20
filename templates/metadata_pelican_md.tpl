Title: {{.Get "title"}}
{{with .Get "subtitle"}}Subtitle: {{.}}{{end}}
Slug: {{.Get "slug"}}
{{with .GetList "tags"}}Tags: {{join . ", "}}{{end}}
{{with .Get "date"}}Date: {{.}}{{end}}
{{with .Get "author"}}Author: {{.}}{{end}}
