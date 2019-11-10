GOCMD=go
GOFMT=goreturns
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GODOC=$(GOCMD) doc
GOGET=$(GOCMD) get

all: fmt build test cover bench doc
dev: testdev benchdev doc
ci: build test cover bench
doc:
	$(GODOC) -all .
fmt:
	$(GOFMT) -l -w .
build:
	$(GOCMD) build -v
test:
	$(GOTEST) -race -v
testdev:
	$(GOTEST) -race -short -cover -covermode=atomic -count 1
bench:
	$(GOTEST) -parallel=4 -run="none" -benchtime="2s" -benchmem -bench=.
benchdev:
	$(GOTEST) -parallel=8 -run="none" -benchtime="1s" -benchmem -bench=.
cover:
	$(GOTEST) -race -cover -covermode=atomic -coverprofile=coverage.txt
