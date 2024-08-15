all: clean cleanall fmt sourcemaps

.PHONY: clean cleanall fmt generate sourcemaps

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

generate:
	@templ generate

sourcemaps:
	@templ generate --source-map-visualisations

