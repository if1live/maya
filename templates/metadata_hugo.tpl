+++
{{range $_,$elem := .Table}}
{{if $elem.IsListValue}}
{{$elem.Key}} = ["{{join $elem.ListValue "\", \""}}"]
{{else}}
{{$elem.Key}} = "{{$elem.Value}}"
{{end}}
{{end}}
+++
