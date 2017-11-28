.PHONY: bootstrap build-parser watch

GRAMMARFILE=grammar.peg
PARSERFILE=lang/parser.go
SAMPLEFILE=samples/sample1.clj

bootstrap:
	mkdir -p parser
	go get -u github.com/mna/pigeon
	go get -u github.com/cespare/reflex

watch:
	reflex --inverse-glob=$(PARSERFILE) -- $(MAKE) test

$(PARSERFILE): $(GRAMMARFILE)
	@echo "== Rebuilding parser"
	echo "package lang" > $@
	pigeon <$< | gofmt >>$@

build-parser: $(PARSERFILE)

test: $(PARSERFILE)
	@echo "== Running unit tests"
	go test ./...
