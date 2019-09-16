package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
)

const (
	name = "gogen-enum"
)

var (
	input       = flag.String("input", "", "input file name;")
	output      = flag.String("output", "", "output file name;")
	packageName = flag.String("package", "", "package name;")
)

var funcMap = template.FuncMap{
	"pubIdent":            publicIdentifier,
	"privIdent":           privateIdentifier,
	"first":               first,
	"toLower":             strings.ToLower,
	"toUpper":             strings.ToUpper,
	"titleCase":           strings.ToTitle,
	"camelCase":           strcase.ToCamel,
	"lowerCamelCase":      strcase.ToLowerCamel,
	"snakeCase":           snakeCase,
	"lintName":            lintName,
	"jsonEncode":          jsonEncode,
	"jsonEncodeOmitEmpty": jsonEncodeOmitEmpty,
	"backTick":            backTick,
	"base64Encode":        base64Encode,
	"args":                args,
	"bitmaskComposite":    bitmask,
}

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(name + ": ")
	flag.Usage = Usage
	flag.Parse()

	// imports := []string{"bytes", "encoding/json"}

	enums, err := loadYAMLFile(input)
	if err != nil {
		log.Fatalln(err)
	}

	importsMap := make(map[string]interface{})

	for _, v := range enums {
		if v.Bitmask {
			// nothing
		}
		if v.Marshal {
			importsMap["bytes"] = true
			importsMap["encoding/json"] = true
		}
		if v.CaseInsensitive {
			importsMap["strings"] = true
		}
	}

	imports := []string{}
	for k := range importsMap {
		imports = append(imports, k)
	}

	gen := generator{
		Generator:   name,
		PackageName: *packageName,
		InputFile:   *input,
		OutputFile:  *output,
		Imports:     imports,
		EnumMap:     enums,
	}

	pipeline := func() generator {
		return gen
	}
	funcMap["pipeline"] = pipeline

	t := template.Must(template.New("").Funcs(funcMap).Parse(tmpl))

	var f *os.File
	if len(*output) > 0 {
		f, err = os.Create(*output)
		defer f.Close()
	} else {
		f = os.Stdout
	}

	err = t.ExecuteTemplate(f, "Document", gen)
	if err != nil {
		log.Fatal("Execute: ", err)
		return
	}

}

func loadYAMLFile(filePath *string) (enumsMap, error) {
	if filePath == nil || *filePath == "" {
		return nil, fmt.Errorf("filePath not set")
	}

	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", *filePath)
	}

	buf, err := ioutil.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}

	var m enumsMap
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

type generator struct {
	Generator   string
	PackageName string
	InputFile   string
	OutputFile  string
	Imports     []string
	EnumMap     map[string]enum
}

type enumsMap map[string]enum
type enum struct {
	Base            string   `yaml:"base"`
	Marshal         bool     `yaml:"marshal"`
	Bitmask         bool     `yaml:"bitmask"`
	CaseInsensitive bool     `yaml:"caseinsensitive"`
	Items           []string `yaml:"items"`
}

func pp(p interface{}) string {

	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		log.Printf("err %v", err)
		return ""
	}
	return string(b)
}

func publicIdentifier(s string) string {

	return lintName(strcase.ToCamel(s))
}

func privateIdentifier(s string) string {
	return lintName(strcase.ToLowerCamel(s))
}

func first(s string) string {
	return strings.ToLower(string(s[0]))
}

func snakeCase(s string) string {

	offset := strings.Index(s, "IPv")
	if offset > 0 {
		return toSnakeCase(s[:offset]) + "_" + strings.ToLower(s[offset:])
	}
	return toSnakeCase(s)
}

func toSnakeCase(in string) string {
	runes := []rune(in)

	var out []rune
	for i := 0; i < len(runes); i++ {
		if i > 0 && (unicode.IsUpper(runes[i]) || unicode.IsNumber(runes[i])) && ((i+1 < len(runes) && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

func jsonEncode(s string) string {
	return fmt.Sprintf("json:\"%s\"", s)
}

func jsonEncodeOmitEmpty(s string) string {
	return fmt.Sprintf("json:\"%s,omitempty\"", s)
}

func backTick(s string) string {
	return fmt.Sprintf("`%s`", s)
}

func base64Encode() string {
	return "encode:\"base64\""
}

func args(kvs ...interface{}) (map[string]interface{}, error) {

	if len(kvs)%2 != 0 {
		return nil, fmt.Errorf("fnArgs requires even number of arguments")
	}

	m := make(map[string]interface{})
	for i := 0; i < len(kvs); i += 2 {
		s, ok := kvs[i].(string)
		if !ok {
			return nil, errors.New("even args to args must be strings")
		}
		m[s] = kvs[i+1]
	}
	return m, nil
}

func bitmask(s string, v interface{}) string {

	var (
		items  []string
		result []string
		ok     bool
	)
	items, ok = v.([]string)
	if !ok {
		return ""
	}
	// [1:] to skip over the Unknown element
	for _, v := range items[1:] {
		result = append(result, strcase.ToCamel(s)+lintName(v))
	}
	ss := strings.Join(result, " + ")
	return ss
}
