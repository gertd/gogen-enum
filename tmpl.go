package main

const tmpl = `
// Package {{ .PackageName }} -- generated by gogen-enum from source: [{{ .InputFile }}]
// 
// !! DO NOT EDIT !! 
// 
package {{ .PackageName }}

import (
	{{ with .Imports -}}
    {{ range . -}} 
	{{ . | printf "%q" }}
	{{ end -}}
	{{ end -}}
)

{{ $enums := .EnumMap }}
{{ range $key, $value := $enums -}}
// {{ $key | camelCase | lintName }} -- enum
type {{ $key | camelCase | lintName }} {{ $value.Base }}

// {{ $key | camelCase | lintName }} -- enum constants
const (
	{{ $items := $value.Items -}}
	{{ range $id, $item := $items -}}
	{{ $key | camelCase | lintName }}{{ $item | camelCase | lintName }}{{ if not $id }} {{ $key | camelCase | lintName }} = 0 + iota{{ end }}
	{{ end -}}
)

// {{ $key | camelCase | lintName }}ID -- map enum constant to string
var {{ $key | lowerCamelCase | lintName }}ID = map[{{ $key | camelCase | lintName }}]string{
	{{ range $id, $item := $items -}}
	{{ $key | camelCase | lintName }}{{ $item | camelCase | lintName }}: "{{ $item }}",
	{{ end -}}
}

// {{ $key | camelCase | lintName }}Name -- map string to enum constant
var {{ $key | lowerCamelCase | lintName }}Name = map[string]{{ $key | camelCase | lintName }}{
	{{ range $id, $item := $items -}}
	"{{ $item }}":{{ $key | camelCase | lintName }}{{ $item | camelCase | lintName }},
	{{ end -}}
}

// String -- {{ $key | camelCase | lintName }}
func (t {{ $key | camelCase | lintName }}) String() string { return {{ $key | lowerCamelCase | lintName }}ID[t] }

// MarshalJSON -- {{ $key | camelCase | lintName }}
func (t {{ $key | camelCase | lintName }}) MarshalJSON() ([]byte, error) {

	buffer := bytes.NewBufferString("\"")
	buffer.WriteString({{ $key | lowerCamelCase | lintName }}ID[t])
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

// UnmarshalJSON -- {{ $key | camelCase | lintName }}
func (t *{{ $key | camelCase | lintName }}) UnmarshalJSON(b []byte) (err error) {

	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*t = {{ $key | lowerCamelCase | lintName }}Name[s]
	return nil
}

{{ end -}}
`