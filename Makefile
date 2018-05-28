GO ?= go

DIST       ?= dist
TARGETS 	= vcapenv vcapenvwrapper
TARGETS_BIN = $(TARGETS:%=$(DIST)/%)

SRC = $(wildcard *.go)

all: build $(TARGETS_BIN)

build:
	$(GO) build ./...

build-linux:
	$(MAKE) GOOS=linux GOARCH=amd64 DIST=$(DIST)-linux $(TARGETS)

test:
	$(GO) test ./...

.SECONDEXPANSION:
$(TARGETS_BIN): $(DIST)/%: $(SRC) $$(wildcard cmd/%/*.go)
	@mkdir -p $(@D)
	env GOBIN=$(abspath $(DIST)) $(GO) install ./cmd/$*

$(TARGETS): %: $(DIST)/%

clean: $(DIST)
	rm -r $^

.PHONY: clean build test
