.PHONY: test

test:
	@echo "Testing all submodules..."
	@find . -type f -name 'go.mod' -not -path "./vendor/*" | while read modfile; do \
		dir=$$(dirname $$modfile); \
		echo "Testing in $$dir"; \
		(cd $$dir && go test ./...); \
	done