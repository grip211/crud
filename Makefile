PROJECT_NAME=$(shell basename "$(PWD)")
SCRIPT_AUTHOR=Grigorev Pavel <grip211@gmail.com>
SCRIPT_VERSION=0.0.1.dev

all: tests

tests:
	source .env
	go test ./... -v