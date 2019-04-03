package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"
)

var (
	input       = flag.String("input", "", "input file name;")
	output      = flag.String("output", "", "output file name;")
	packageName = flag.String("package", "", "package name;")
)

var funcMap = template.FuncMap{
	"first": func(s string) string {
		return strings.ToLower(string(s[0]))
	},
	"toLower": func(s string) string {
		return strings.ToLower(s)
	},
	"camelCase": func(s string) string {
		return strcase.ToCamel(s)
	},
	"lowerCamelCase": func(s string) string {
		return strcase.ToLowerCamel(s)
	},
	"snakeCase": func(s string) string {
		return snakeCase(s)
	},
	"lintName": func(s string) string {
		return lintName(s)
	},
	"jsonEncode": func(s string) string {
		return fmt.Sprintf("json:\"%s\"", s)
	},
	"jsonEncodeOmitEmpty": func(s string) string {
		return fmt.Sprintf("json:\"%s,omitempty\"", s)
	},
	"backTick": func(s string) string {
		return fmt.Sprintf("`%s`", s)
	},
	"base64Encode": func() string {
		return "encode:\"base64\""
	},
}

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("gogen-enum: ")
	flag.Usage = Usage
	flag.Parse()

	enums, err := loadYAMLFile(input)
	if err != nil {
		log.Fatalln(err)
	}

	t := template.Must(template.New("").Funcs(funcMap).Parse(tmpl))

	generator := generator{
		PackageName: *packageName,
		InputFile:   *input,
		OutputFile:  *output,
		Imports:     []string{"bytes", "encoding/json"},
		EnumMap:     enums,
	}

	var f *os.File
	if len(*output) > 0 {
		f, err = os.Create(*output)
		defer f.Close()
	} else {
		f = os.Stdout
	}

	err = t.Execute(f, generator)
	if err != nil {
		log.Fatal("Execute: ", err)
		return
	}

}

func loadYAMLFile(filePath *string) (enumsMap, error) {
	if filePath == nil {
		return nil, fmt.Errorf("filePath is NIL pointer")
	}

	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("YAML file [%s] does not exist", *filePath)
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
	PackageName string
	InputFile   string
	OutputFile  string
	Imports     []string
	EnumMap     map[string]enum
}

type enumsMap map[string]enum
type enum struct {
	Base  string   `yaml:"base"`
	Items []string `yaml:"items"`
}

func pp(p interface{}) string {

	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		log.Printf("err %v", err)
		return ""
	}
	return string(b)
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

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
// copied from golint (https://github.com/golang/lint/blob/master/lint.go)
// NOTE: this list has been augemented with: [CA, IPv4, IPv6, OS, URN]
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ARM":   true,
	"ASCII": true,
	"CA":    true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"IPv4":  true,
	"IPv6":  true,
	"JSON":  true,
	"LHS":   true,
	"OS":    true,
	"PPC":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"URN":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

// lintName returns a different name if it should be different.
func lintName(name string) (should string) {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}
	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		return name
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word
		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// Keep consistent case, which is lowercase only at the start.
			if w == 0 && unicode.IsLower(runes[w]) {
				u = strings.ToLower(u)
			}
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))
		} else if w > 0 && strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}
