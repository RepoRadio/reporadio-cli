.PHONY: build clean run

run:
	go run main.go create test
build:
	go build -o bin/reporadio-cli 

clean:
	rm -rf bin/
	rm -rf .reporadio/test

install: build
	go install