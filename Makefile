.PHONY: all clean

BUILD_DIR := ./build
PLUGIN_DIR := ./build/plugins
PLUGINS_SRC := $(wildcard plugins/*/main.go)
PLUGINS := $(patsubst plugins/%/main.go,%,$(PLUGINS_SRC))

all: build

go: build run

build: build_plugins build_oracle copy_env

run:
	@echo "Starting Oracle..."
	@$(BUILD_DIR)/oracle

build_plugins: $(addprefix build_plugin_, $(PLUGINS))

build_plugin_%:
	@echo "Building $*..."
	@go build -buildmode=plugin -o $(PLUGIN_DIR)/$*.so plugins/$*/main.go

build_oracle:
	@echo "Building oracle..."
	@go build -o $(BUILD_DIR)/oracle

copy_env:
	@echo "Copying .env..."
	@cp .env $(BUILD_DIR)

clean:
	@echo "Cleaning up..."
	@rm -rf $(PLUGIN_DIR)/*.so
	@rm -rf $(BUILD_DIR)/oracle
	@rm -rf $(BUILD_DIR)/.env
