.PHONY: all clean

all: build

go: build run

build: build_plugins build_oracle

build_plugins:
	@echo "Building plugins..."
	@cd plugins/bestapi && go build -buildmode=plugin -o ../../build/plugins/bestapi.so
	@cd plugins/goerli && go build -buildmode=plugin -o ../../build/plugins/goerli.so

build_oracle:
	@echo "Building main project..."
	@go build -o build/oracle

clean:
	@echo "Cleaning up..."
	@rm -rf build/plugins/*.so
	@rm -rf build/oracle

run:
	@echo "Starting Oracle"
	@./build/oracle
