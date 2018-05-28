GO ?= go

TARGETS 	= vcapenv vcapenvwrapper
TARGETS_BIN = $(TARGETS:%=build/%)

SRC = $(wildcard *.go)

all: $(TARGETS_BIN)

build/%: cmd/%/main.go $(SRC) $(wildcard % *.go)
	@mkdir -p $(@D)
	$(GO) build -o $@ $<

clean: build
	rm -r $^

.PHONY: clean
