.DEFAULT_GOAL = build
BUILD_FILE = ./bytecode/build.go

build:
	echo "package bytecode" > $(BUILD_FILE)
	echo "const (" >> $(BUILD_FILE)
	echo "\tAGORA_BUILD = \"$(shell git rev-parse --short HEAD)\"" >> $(BUILD_FILE)
	echo ")" >> $(BUILD_FILE)
	go install ./...

.PHONY: build

