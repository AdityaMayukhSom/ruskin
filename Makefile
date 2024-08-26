BUILD_FLAGS = -race
OUTPUT_FILE = ./bin/ruskin

build:
	@echo "Building the source files"
	CGO_ENABLED=1 go build $(BUILD_FLAGS) -o $(OUTPUT_FILE) .

run: build
	@echo "Starting Ruskin"
	$(OUTPUT_FILE)

clean:
	@echo "Removing Ruskin executable"
	rm -fr $(OUTPUT_FILE)

rebuild: clean build

cleanrun: clean run