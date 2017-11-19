.PHONY: bootstrap build-parser watch

PARSERFILE=lang/parser.go
SAMPLEFILE=samples/sample1.clj

bootstrap:
	mkdir -p parser
	go get -u github.com/mna/pigeon
	go get -u github.com/cespare/reflex

watch:
	reflex --inverse-glob=$(PARSERFILE) -- $(MAKE) test-parser

build-parser:
	echo "package lang" > $(PARSERFILE)
	pigeon <grammar.peg | gofmt >>$(PARSERFILE)

test-parser: build-parser
	go run main.go $(SAMPLEFILE)

