TAG=$(shell git describe --tags --always)
VERSION=$(TAG:v%=%)
NAME=convector
DOCKER_REGISTRY=ghcr.io
REPO=shawntoffel/$(NAME)
GO=go
BUILD=GOARCH=amd64 $(GO) build -ldflags="-s -w -X 'github.com/$(REPO)/internal.Version=$(VERSION)'" 

.PHONY: all deps test build build-linux docker-build docker-save docker-deploy clean

all: deps test build 
deps:
	$(GO) mod download

test:
	$(GO) vet ./...
	$(GO) test -v -race ./...

build:
	$(BUILD) -o bin/$(NAME)-$(VERSION) ./cmd/...

build-linux:
	CGO_ENABLED=0 GOOS=linux $(BUILD) -a -installsuffix cgo -o bin/$(NAME) ./cmd/...

docker-build:
	docker build -t $(DOCKER_REGISTRY)/$(REPO):$(VERSION) .

docker-push:
	docker push $(DOCKER_REGISTRY)/$(REPO):$(VERSION)

clean:
	@find bin -type f -delete -print
