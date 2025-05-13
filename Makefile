# Makefile

BINARY_NAME=docky
BUILD_DIR=build

build:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

clean:
	rm -rf $(BUILD_DIR)
