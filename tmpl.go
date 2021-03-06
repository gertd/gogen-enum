package main

const tmpl = `
{{- define "Header" -}}
//
// Code generated by {{.Generator}} DO NOT EDIT.
// 
{{- end -}}

{{- define "PackageBlock" -}}
package {{ .PackageName }}
{{- end -}}

{{ define "ImportsBlock" }}
import (
	{{ with .Imports -}}
    {{ range . -}} 
	{{ . | printf "%q" }}
	{{ end -}}
	{{ end -}}
)
{{- end -}}

{{ define "EnumBlock" }} {{- /* BOF EnumBlock */ -}}
{{ range $name, $props := .EnumMap -}}
{{ template "EnumType" args "name" $name "props" $props }}
{{ template "EnumNumericConstants" args "name" $name "props" $props }}
{{ template "EnumStringConstants" args "name" $name "props" $props }}
{{ template "EnumMapID2String" args "name" $name "props" $props }}
{{ if $props.CaseInsensitive }}
{{ template "EnumMapString2ID-CI" args "name" $name "props" $props }}
{{ template "EnumNewFunc-CI" args "name" $name "props" $props }}
{{ else }}
{{ template "EnumMapString2ID" args "name" $name "props" $props }}
{{ template "EnumNewFunc" args "name" $name "props" $props }}
{{ end }}
{{ template "EnumStringFunc" args "name" $name "props" $props }}
{{ if $props.Marshal }}
{{ template "EnumMarshalFunc" args "name" $name "props" $props }}
{{ template "EnumUnmarshalFunc" args "name" $name "props" $props }}
{{ end }} {{- /* EOF if $props.Marshal */ -}}
{{ if $props.Bitmask }}
{{ template "EnumBitmaskSetFunc" args "name" $name "props" $props }}
{{ template "EnumBitmaskClearFunc" args "name" $name "props" $props }}
{{ template "EnumBitmaskToggleFunc" args "name" $name "props" $props }}
{{ template "EnumBitmaskHasFunc" args "name" $name "props" $props }}
{{ end }} {{- /* EOF if $props.Bitmask */ -}}
{{ end }} {{- /* EOF range $name, $props := .EnumMap */ -}}
{{ end }} {{- /* EOF EnumBlock */ -}}

{{ define "EnumType" }}
// {{ $.name | pubIdent }} -- enum type
type {{ $.name | pubIdent }} {{ $.props.Base }}
{{- end -}}

{{ define "EnumNumericConstants" }}
// {{ $.name | pubIdent }} -- enum constants
const (
	{{ $items := $.props.Items -}}
	{{ range $id, $item := $items -}} 
	{{ $.name | pubIdent }}{{ $item | pubIdent }}{{ if not $id }} {{ $.name | pubIdent }} = {{ if $.props.Bitmask }} 1 << iota {{ else }} 0 + iota {{ end }} {{ end }}
	{{ end -}}
	{{ if $.props.Bitmask -}}
	{{ $.name | pubIdent }}All = {{ bitmaskComposite $.name $items }}
	{{- end }} 
)
{{- end -}}

{{ define "EnumStringConstants" }}
// {{ $.name | privIdent }} -- enum string representation constants
const (
	{{ $items := $.props.Items -}}
	{{ range $id, $item := $items -}}
	{{ $.name | privIdent }}{{ $item | pubIdent }} = "{{ $item }}"
	{{ end -}}
	{{ if $.props.Bitmask }}{{ $.name | privIdent }}All = "All"{{- end }} 
)
{{- end -}}

{{ define "EnumMapID2String" }}
// {{ $.name | pubIdent }}ID -- map enum constant to string
var {{ $.name | privIdent }}ID = map[{{ $.name | pubIdent }}]string{
	{{ $items := $.props.Items -}}
	{{ range $id, $item := $items -}}
	{{ $.name | pubIdent }}{{ $item | pubIdent }}: {{ $.name | privIdent }}{{ $item | pubIdent }},
	{{ end -}}
	{{ if $.props.Bitmask }}{{ $.name | pubIdent }}All: {{ $.name | privIdent }}All,{{- end }} 
}
{{- end -}}

{{ define "EnumMapString2ID" }}
// {{ $.name | pubIdent }}Name -- map string to enum constant
var {{ $.name | privIdent }}Name = map[string]{{ $.name | pubIdent }}{
	{{ $items := $.props.Items -}}
	{{ range $id, $item := $items -}}
	{{ $.name | privIdent }}{{ $item | pubIdent }}: {{ $.name | pubIdent }}{{ $item | pubIdent }},
	{{ end -}}
	{{ if $.props.Bitmask }}{{ $.name | privIdent }}All: {{ $.name | pubIdent }}All,{{- end }} 
}
{{- end -}}

{{ define "EnumMapString2ID-CI" }}
// {{ $.name | pubIdent }}Name -- map string to enum constant
var {{ $.name | privIdent }}Name = map[string]{{ $.name | pubIdent }}{
	{{ $items := $.props.Items -}}
	{{ range $id, $item := $items -}}
	strings.ToLower({{ $.name | privIdent }}{{ $item | pubIdent }}): {{ $.name | pubIdent }}{{ $item | pubIdent }},
	{{ end -}}
	{{ if $.props.Bitmask }}strings.ToLower({{ $.name | privIdent }}All): {{ $.name | pubIdent }}All,{{- end }} 

}
{{- end -}}

{{ define "EnumNewFunc" }}
// New{{ $.name | pubIdent }} -- Create {{ $.name | pubIdent }} instance from string representation
func New{{ $.name | pubIdent }}(k string) {{ $.name | pubIdent }} {
	if kind, ok := {{ $.name | privIdent }}Name[k]; ok {
		return kind
	}
	return {{ $.name | pubIdent }}Unknown
}
{{- end -}}

{{ define "EnumNewFunc-CI" }}
// New{{ $.name | pubIdent }} -- Create {{ $.name | pubIdent }} instance from string representation
func New{{ $.name | pubIdent }}(k string) {{ $.name | pubIdent }} {
	if kind, ok := {{ $.name | privIdent }}Name[strings.ToLower(k)]; ok {
		return kind
	}
	return {{ $.name | pubIdent }}Unknown
}
{{- end -}}

{{ define "EnumStringFunc" }}
// String -- {{ $.name | pubIdent }}
func (t {{ $.name | pubIdent }}) String() string { return {{ $.name | privIdent }}ID[t] }
{{- end -}}

{{ define "EnumMarshalFunc" }}
// MarshalJSON -- {{ $.name | pubIdent }}
func (t {{ $.name | pubIdent }}) MarshalJSON() ([]byte, error) {

	buffer := bytes.NewBufferString("\"")
	buffer.WriteString({{ $.name | privIdent }}ID[t])
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}
{{- end -}}

{{ define "EnumUnmarshalFunc" }}
// UnmarshalJSON -- {{ $.name | pubIdent }}
func (t *{{ $.name | pubIdent }}) UnmarshalJSON(b []byte) (err error) {

	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*t = {{ $.name | privIdent }}Name[s]
	return nil
}
{{- end -}}

{{ define "EnumBitmaskSetFunc" }}
// Set -- set flag
func (t *{{ $.name | pubIdent }}) Set(flag {{ $.name | pubIdent }}) {
	*t |= flag
}
{{- end -}}

{{ define "EnumBitmaskClearFunc" }}
// Clear -- clear flag
func (t *{{ $.name | pubIdent }}) Clear(flag {{ $.name | pubIdent }}) {
	*t &^= flag
}
{{- end -}}

{{ define "EnumBitmaskToggleFunc" }}
// Toggle -- toggle flag state
func (t *{{ $.name | pubIdent }}) Toggle(flag {{ $.name | pubIdent }}) {
	*t ^= flag
}
{{- end -}}

{{ define "EnumBitmaskHasFunc" }}
// Has -- is flag set?
func (t {{ $.name | pubIdent }}) Has(flag {{ $.name | pubIdent }}) bool {
	return t&flag != 0
}
{{- end -}}

{{ define "Document" }}
{{ template "Header" pipeline }}
{{ template "PackageBlock" pipeline }}
{{ template "ImportsBlock" pipeline }}
{{ template "EnumBlock" pipeline }}
{{ end }}

{{ template "Document" }}
`
