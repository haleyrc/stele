all: clean cleanall fmt generate test

.PHONY: clean cleanall fmt install generate test

SOURCEMAPS := $(shell fd --type file --color never --extension html --no-ignore _templ_sourcemap)
GENERATED := $(shell fd --type file --color never --extension go _templ)

echo:
	@echo $(SOURCEMAPS)
	@echo $(GENERATED)

clean:
	@rm -rf $(SOURCEMAPS)

cleanall: clean
	@rm -rf $(GENERATED)

fmt:
	@templ fmt .

install:
	go install ./cmd/stele

generate:
	@templ generate --source-map-visualisations -include-version=false

test:
	go test -v -count=1 ./...
