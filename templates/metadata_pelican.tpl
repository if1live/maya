{{with .Get "title"}}Title: {{.}}{{end}}
{{with .Get "subtitle"}}Subtitle: {{.}}{{end}}
{{with .Get "slug"}}Slug: {{.}}{{end}}
{{with .GetList "tags"}}Tags: {{join . ", "}}{{end}}
{{with .Get "date"}}Date: {{.}}{{end}}
{{with .Get "author"}}Author: {{.}}{{end}}
