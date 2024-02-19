.PHONY: test tag
SUBMODULES := array-utils async-utils map-utils math-utils number opt queue ref-utils str-utils

test:
	@echo "Testing all submodules..."
	@find . -type f -name 'go.mod' -not -path "./vendor/*" | while read modfile; do \
		dir=$$(dirname $$modfile); \
		echo "Testing in $$dir"; \
		(cd $$dir && go test ./...); \
	done

tag:
ifndef VERSION
	$(error VERSION is undefined. Usage: make tag VERSION=v1.2.3)
endif
	@for submodule in $(SUBMODULES); do \
		echo "Tagging $$submodule with $(VERSION)"; \
		git tag $$submodule/$(VERSION); \
		git push origin $$submodule/$(VERSION); \
	done