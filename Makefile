# core-banking-go development commands.
# Run "make" or "make help" to list them.

# Load variables from .env, when the file exists, and export them to the commands.
ifneq (,$(wildcard .env))
include .env
export
endif

BINARY_DIR := bin
COVER_FILE := coverage.out

# Extra flags for go test. The CI enables the race detector with: make test RACE=-race
RACE ?=

# Windows uses cmd.exe so the commands behave the same in PowerShell and Git Bash.
ifeq ($(OS),Windows_NT)
SHELL     := cmd.exe
BINARY    := $(BINARY_DIR)/api.exe
CLEAN_CMD := if exist $(BINARY_DIR) rmdir /s /q $(BINARY_DIR) & if exist $(COVER_FILE) del /q $(COVER_FILE)
else
BINARY    := $(BINARY_DIR)/api
CLEAN_CMD := rm -rf $(BINARY_DIR) $(COVER_FILE)
endif

.DEFAULT_GOAL := help
.PHONY: help run build test cover fmt vet tidy clean up down logs lint

help:
	$(info Usage: make <target>)
	$(info )
	$(info help     Show this help)
	$(info run      Run the API (loads .env when it exists))
	$(info build    Build the API binary into $(BINARY_DIR)/)
	$(info test     Run the tests with coverage (use RACE=-race to enable the race detector))
	$(info cover    Run the tests and show the coverage per function)
	$(info fmt      Format the code)
	$(info vet      Run go vet)
	$(info tidy     Sync go.mod and go.sum with the code)
	$(info clean    Remove the binary and the coverage file)
	$(info up       Start the local dependencies (PostgreSQL and Keycloak))
	$(info down     Stop the local dependencies (data is kept))
	$(info logs     Follow the logs of the local dependencies)
	$(info lint     Run the linters (golangci-lint))

run:
	go run ./cmd/api

build:
	go build -o $(BINARY) ./cmd/api

test:
	go test $(RACE) -coverprofile=$(COVER_FILE) ./...

cover: test
	go tool cover -func=$(COVER_FILE)

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	@$(CLEAN_CMD)

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

lint:
	golangci-lint run ./...
