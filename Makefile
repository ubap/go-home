# ====================================================================================
#  Makefile for the Go Project
# ====================================================================================

# --- Variables ---

# The name of the final binary
BINARY_NAME=myapp

# The directory to place the final binary into
BUILD_DIR=target


# --- Targets ---

# The 'all' target is the default one executed when you just run 'make'
.PHONY: all
all: build

# Builds the Go application
.PHONY: build
build:
	@echo "==> Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "==> Done! Binary is at $(BUILD_DIR)/$(BINARY_NAME)"

# Cleans the build artifacts
.PHONY: clean
clean:
	@echo "==> Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "==> Done!"

# Runs the application
.PHONY: run
run: build
	@echo "==> Running..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Runs the tests
.PHONY: test
test:
	@echo "==> Testing..."
	go test ./... -v

# Installs the binary to the Go bin path
.PHONY: install
install:
	@echo "==> Installing..."
	go install .