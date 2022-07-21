GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOTOOL=$(GOCMD) tool
BINARY_NAME=simple_service
LINTER=golangci-lint

all: test build

mod:
	$(GOMOD) tidy

test:
	$(GOTEST) ./... -v

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

lint:
	$(LINTER) run

coverage:
	$(GOTEST) -coverprofile=coverage.out ./... && $(GOTOOL) cover -html=coverage.out