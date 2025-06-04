# Makefile

BINARY_NAME=docky
BUILD_DIR=bin
APP=cmd/docky/main.go

build:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) $(APP)

run: 
	go run $(APP) $(filter-out $@,$(MAKECMDGOALS))

clean:
	rm -rf $(BUILD_DIR)

%:
	@: