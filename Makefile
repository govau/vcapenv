GO ?= go

DIST       ?= dist
TARGETS	    = vcapenv vcapenvwrapper
TARGETS_BIN = $(TARGETS:%=$(DIST)/%)

SRC = $(wildcard *.go)

all: build cmds

build:
	$(GO) build ./...

test:
	$(GO) test ./...

cmds: $(TARGETS)

linux-cmds:
	$(MAKE) GOOS=linux GOARCH=amd64 DIST=$(DIST)-linux $(TARGETS)

$(TARGETS): %: $(DIST)/%

.SECONDEXPANSION:
$(TARGETS_BIN): $(DIST)/%: cmd/%/main.go $(SRC) $$(wildcard cmd/%/*.go)
	@mkdir -p $(@D)
	$(GO) build -o $@ ./cmd/$*

clean: $(DIST)
	rm -r $^

.PHONY: clean build test
