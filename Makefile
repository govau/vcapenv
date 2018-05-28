GO ?= go

TARGETS 	= vcapenv vcapenvwrapper
TARGETS_BIN = $(TARGETS:%=dist/%)

SRC = $(wildcard *.go)

all: build $(TARGETS_BIN)

build:
	$(GO) build ./...

test:
	$(GO) test ./...

dist/%: cmd/%/main.go $(SRC) $(wildcard % *.go)
	@mkdir -p $(@D)
	$(GO) build -o $@ $<

clean: dist
	rm -r $^

.PHONY: clean build test
