.PHONY: all clean

BUILD_DIR := ./build
PLUGIN_DIR := ./build/plugins

all: build

go: build run

build: build_plugins build_oracle copy_env

run:
	@echo "Starting Oracle"
	@$(BUILD_DIR)/oracle


build_plugins:
	@echo "Building plugins..."
	@go build -buildmode=plugin -o $(PLUGIN_DIR)/bestapi.so plugins/bestapi/main.go
	@go build -buildmode=plugin -o $(PLUGIN_DIR)/goerli.so plugins/goerli/main.go
	@go build -buildmode=plugin -o $(PLUGIN_DIR)/mathapi.so plugins/mathapi/main.go
	@go build -buildmode=plugin -o $(PLUGIN_DIR)/exchangerateapi.so plugins/exchangerateapi/main.go

build_oracle:
	@echo "Building main project..."
	@go build -o $(BUILD_DIR)/oracle

copy_env:
	@cp .env $(BUILD_DIR)

clean:
	@echo "Cleaning up..."
	@rm -rf $(PLUGIN_DIR)/*.so
	@rm -rf $(BUILD_DIR)/oracle
	@rm -rf $(BUILD_DIR)/.env
