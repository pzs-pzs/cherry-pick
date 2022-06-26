
.PHONY: build
build:
	@echo "Building..."
	$Q CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./ ./...

build-local:
	@echo "Building..."
	$Q CGO_ENABLED=0 go build -ldflags="-s -w" -o ./ ./...