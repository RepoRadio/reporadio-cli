.PHONY: build clean

build:
	go build -o bin/reporadio-cli 

clean:
	rm -rf bin/

install: build
	go install