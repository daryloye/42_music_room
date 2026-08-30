.SILENT:

GOIMPORTS_VERSION ?= v0.24.0
GOIMPORTS := $(shell go env GOPATH)/bin/goimports
GO_FILES := $(shell find backend -type f -name '*.go' ! -path 'backend/prisma/db/*')

.PHONY: run generate-prisma create-db setup format format-check

run:
	cd backend && go run main.go

generate-prisma:
	cd backend && go run github.com/steebchen/prisma-client-go db push
	cd backend && go generate ./...
	cd backend && go mod tidy

create-db:
	docker run --name music-room-postgres \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_DB=music_room \
    -p 5432:5432 \
    -d postgres:16-alpine

setup:
	go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
	cd backend && go mod download
	cd backend && npm ci
	git config core.hooksPath .githooks

format:
	$(GOIMPORTS) -w $(GO_FILES)
	cd backend && npm run format

format-check:
	test -x "$(GOIMPORTS)" || { echo "goimports is missing; run 'make setup'"; exit 1; }
	unformatted="$$($(GOIMPORTS) -l $(GO_FILES))"; \
	if test -n "$$unformatted"; then \
		echo "Go files need formatting:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	cd backend && npm run format:check
