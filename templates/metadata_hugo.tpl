+++
{{range $_,$elem := .Table}}
{{if $elem.IsListValue}}
{{$elem.Key}} = ["{{join $elem.ListValue "\", \""}}"]
{{else}}
{{$elem.Key}} = "{{escape $elem.Value}}"
{{end}}
{{end}}
+++
