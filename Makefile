SUBMODULES := array-utils async-utils map-utils math-utils number opt queue ref-utils str-utils
GOBIN ?= $$(go env GOPATH)/bin

.PHONY: test
test:
	@echo "Testing all submodules..."
	@find . -type f -name 'go.mod' -not -path "./vendor/*" | while read modfile; do \
		dir=$$(dirname $$modfile); \
		echo "Testing in $$dir"; \
		(cd $$dir && go test ./...); \
	done

.PHONY: tag
tag:
ifndef VERSION
	$(error VERSION is undefined. Usage: make tag VERSION=v1.2.3)
endif
	@for submodule in $(SUBMODULES); do \
		echo "Tagging $$submodule with $(VERSION)"; \
		git tag $$submodule/$(VERSION); \
		git push origin $$submodule/$(VERSION); \
	done

.PHONY: tidy
tidy:
	@echo "Go mod tidy in all submodules..."
	@find . -type f -name 'go.mod' -not -path "./vendor/*" | while read modfile; do \
		dir=$$(dirname $$modfile); \
		echo "go mod tidy in $$dir"; \
		(cd $$dir && go mod tidy); \
	done

vet:
	@echo "Go vet + staticcheck in all submodules..."
	@find . -type f -name 'go.mod' -not -path "./vendor/*" | while read modfile; do \
		dir=$$(dirname $$modfile); \
		echo "go vet in $$dir"; \
		(cd $$dir && go vet ./...); \
		echo "staticcheck in $$dir"; \
		(cd $$dir && staticcheck .); \
	done

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml