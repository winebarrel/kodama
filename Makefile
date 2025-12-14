.PHONY: all
all: vet build
# all: vet test build

.PHONY: build
build:
	go build ./cmd/kodama

.PHONY: vet
vet:
	go vet ./...

# .PHONY: test
# test:
# 	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run
