.PHONY: build clean run docs lint format

run:
	go run main.go create test
build:
	go build -o bin/reporadio-cli

clean:
	rm -rf bin/
	rm -rf .reporadio/test

install: build
	go install

docs:
	@echo "Installing pkgsite if not present..."
	@go install golang.org/x/pkgsite/cmd/pkgsite@latest
	@echo "Starting documentation server on http://localhost:8080"
	@echo "Press Ctrl+C to stop the server"
	pkgsite -http=localhost:8080

lint:
	trunk check

format:
	trunk fmt
