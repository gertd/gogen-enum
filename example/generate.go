// !build

package test

//go:generate gogen-enum -input ./test.yaml -package main -output ./test.go
//go:generate gofmt -w test.go
//go:generate golangci-lint run test.go
