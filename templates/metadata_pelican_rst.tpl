{{.Get "title"}}
{{seperator (.Get "title") "="}}
{{with .Get "subtitle"}}:subtitle: {{.}}{{end}}
:slug: {{.Get "slug"}}
{{with .GetList "tags"}}:tags: {{join . ", "}}{{end}}
{{with .Get "date"}}:date: {{.}}{{end}}
{{with .Get "author"}}:author: {{.}}{{end}}
