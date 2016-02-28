{{range $_,$elem := .Table}}
{{if $elem.IsListValue}}
{{title $elem.Key}}: {{join $elem.ListValue ", "}}
{{else}}
{{title $elem.Key}}: {{$elem.Value}}
{{end}}
{{end}}
