# text reset
NO_COLOR=\033[0m
# green
OK_COLOR=\033[32;01m
# red
ERROR_COLOR=\033[31;01m
# cyan
WARN_COLOR=\033[36;01m
# yellow
ATTN_COLOR=\033[33;01m

ROOT_DIR := $(git rev-parse --show-toplevel)
BIN_DIR  := ./bin

LINTER := $(BIN_DIR)/golangci-lint

GOOS :=
ifeq ($(OS),Windows_NT)
	GOOS = windows
else 
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS = linux
	endif
	ifeq ($(UNAME_S),Darwin)
		GOOS = darwin
	endif
endif
GOARCH ?= amd64

BINARY := gogen-enum
VERSION:=`git describe --tags --dirty`
COMMIT :=`git rev-parse --short HEAD 2>/dev/null`
DATE   :=`date "+%FT%T%z"`
PLATFORMS := windows linux darwin
OS = $(word 1, $@)

LDFLAGS := -ldflags "-w -s -X github.com/gertd/$(BINARY)/cmd.version=${VERSION} -X github.com/gertd/$(BINARY)/cmd.date=${DATE} -X github.com/gertd/$(BINARY)/cmd.commit=${COMMIT}"

.PHONY: all
all: build test lint

.PHONY: build
build:
	@echo "$(ATTN_COLOR)==> build GOOS=$(GOOS) GOARCH=$(GOARCH) VERSION=$(VERSION)@$(COMMIT) $(NO_COLOR)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=on go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY) ./

.PHONY: test
test:
	@echo "$(ATTN_COLOR)==> test $(NO_COLOR)"
	@go test -v ./...

$(LINTER):
	@echo "$(ATTN_COLOR)==> get  $(NO_COLOR)"
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.15.0
 
.PHONY: lint
lint: $(LINTER)
	@echo "$(ATTN_COLOR)==> lint $(NO_COLOR)"
	@$(LINTER) run
	@echo "$(NO_COLOR)\c"

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@echo "$(ATTN_COLOR)==> release GOOS=$(GOOS) GOARCH=$(GOARCH) release/$(BINARY)-$(OS)-$(GOARCH) $(NO_COLOR)"
	@mkdir -p release
	@GOOS=$(OS) GOARCH=$(GOARCH) GO111MODULE=on go build $(LDFLAGS) -o release/$(BINARY)-$(OS)-$(GOARCH)$(if $(findstring $(OS),windows),".exe","")

.PHONY: release
release: windows linux darwin

.PHONY: install
install:
	@echo "$(ATTN_COLOR)==> install $(NO_COLOR)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=on go install $(LDFLAGS) ./
