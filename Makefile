OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)
FILENAME=aoc2020-$(OS)-$(ARCH)
FILE_LOCATION=./bin/$(FILENAME)

build:
	go build -o ./bin/ddns-job ./

buildWithArch:
	go build -o $(FILE_LOCATION) ./

buildLocation:
	@echo $(FILE_LOCATION)

buildName:
	@echo $(FILENAME)

test:
	go test ./...