.SILENT:

GOIMPORTS_VERSION ?= v0.24.0
GOIMPORTS := $(shell go env GOPATH)/bin/goimports
GO_FILES := $(shell find backend -type f -name '*.go' ! -path 'backend/prisma/db/*')

COMPOSE=docker-compose
DEV_COMPOSE=$(COMPOSE) -f docker-compose.yml -f docker-compose.dev.yml

.PHONY: dev up down setup generate-prisma generate-swagger test format docker-clean

dev:
	$(DEV_COMPOSE) up --abort-on-container-exit

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

setup:
	go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
	cd backend && go mod download
	cd bruno && npm ci
	git config core.hooksPath .githooks

generate-prisma:
	cd backend && go run github.com/steebchen/prisma-client-go db push
	cd backend && go generate ./...
	cd backend && go mod tidy

generate-swagger:
	cd backend && swag init

test:
	cd bruno && npm run test

format:
	$(GOIMPORTS) -w $(GO_FILES)
	cd backend && go run github.com/steebchen/prisma-client-go format --schema=prisma/schema.prisma

docker-clean:
	$(DEV_COMPOSE) down --volumns --remove-orphans
	$(COMPOSE) down --remove-orphans --rmi local
