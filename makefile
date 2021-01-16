GOCMD=go
GOBUILD=$(GOCMD) build

BINARY_NAME=bp
BINARY_UNIX=$(BINARY_NAME)_unix
INIT_FILE=src/main.go


build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./dist/$(BINARY_NAME) $(INIT_FILE)