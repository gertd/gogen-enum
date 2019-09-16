// !build

package tests

//go:generate go get -u github.com/gertd/gogen-enum

//go:generate gogen-enum -input ./architecture.yaml -package tests -output ./architecture-gen.go
//go:generate gofmt -w architecture-gen.go

//go:generate gogen-enum -input ./addresskind.yaml -package tests -output ./addresskind-gen.go
//go:generate gofmt -w addresskind-gen.go

//go:generate gogen-enum -input ./packagemanager.yaml -package tests -output ./packagemanager-gen.go
//go:generate gofmt -w packagemanager-gen.go

//go:generate gogen-enum -input ./pluralize.yaml -package tests -output ./pluralize-gen.go
//go:generate gofmt -w pluralize-gen.go

//go:generate gogen-enum -input ./sourcetype.yaml -package tests -output ./sourcetype-gen.go
//go:generate gofmt -w sourcetype-gen.go
