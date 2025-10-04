BINARY_NAME=myapp

# The directory to place the final binary into
BUILD_DIR=target

# The 'all' target is the default one executed when you just run 'make'
.PHONY: all
all: build

.PHONY: build
build:
	@echo "==> Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "==> Done! Binary is at $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: clean
clean:
	@echo "==> Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "==> Done!"

.PHONY: run
run: build
	@echo "==> Running..."
	./$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: test
test:
	@echo "==> Testing..."
	go test ./... -v

.PHONY: deploy
deploy:
	@echo "==> Deploying..."
	./cmd/deploy.sh

.PHONY: logs
logs:
	@echo "==> Downloading logs..."
	./cmd/logs.sh > logs.txt