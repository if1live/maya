+++
{{with .Get "title"}}title = "{{.}}"{{end}}
{{with .Get "subtitle"}}subtitle = "{{.}}"{{end}}
{{with .Get "slug"}}slug = "{{.}}"{{end}}
{{with .GetList "tags"}}tags = ["{{join . "\", \""}}"]{{end}}
{{with .Get "date"}}date = "{{.}}"{{end}}
{{with .Get "author"}}author = "{{.}}"{{end}}
+++
