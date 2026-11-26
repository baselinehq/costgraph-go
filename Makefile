OAPI_CODEGEN_VERSION := v2.9.0
OAPI_CODEGEN := go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)

.PHONY: fmt build lint generate check-generated

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

build:
	go build ./...

lint:
	go vet ./...

generate:
	$(OAPI_CODEGEN) -config api/oapi-codegen.yaml -o internal/api/client.gen.go api/management.openapi.json

check-generated:
	@tmp=$$(mktemp); trap 'rm -f "$$tmp"' EXIT; \
	$(OAPI_CODEGEN) -config api/oapi-codegen.yaml -o "$$tmp" api/management.openapi.json && \
	diff -u internal/api/client.gen.go "$$tmp"
