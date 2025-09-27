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

.PHONY: itest
itest:
	@echo "==> Starting mock server in the background..."
	# 1. Uruchom serwer w tle za pomocą '&'.
	# 2. Natychmiast zapisz jego PID do zmiennej MOCK_PID.
	#    - Używamy `$$!` zamiast `$!`, ponieważ Makefile interpretuje pojedynczy `$`
	#      jako swoją zmienną. `$$` przekazuje dosłowny `$` do powłoki (shell).
	# 3. Ustaw 'trap', aby komenda 'kill' wykonała się na końcu,
	#    niezależnie od wyniku testów. To gwarantuje posprzątanie po sobie.
	# 4. Uruchom testy z odpowiednią zmienną środowiskową.
	#
	# Wszystko jest połączone w jedną linię za pomocą `\` i wykonywane w jednej powłoce.
	@go run ./mocks/ & \
	MOCK_PID=$$! ; \
	\
	@echo "Mock server started with PID: $$MOCK_PID" ; \
	@echo "Waiting for server to become available..." ; \
	# Opcjonalnie, ale zalecane: pętla czekająca na gotowość serwera
	for i in $$(seq 1 10); do \
		if curl -s -o /dev/null http://localhost:8080; then \
			echo "Mock server is ready!"; \
			break; \
		fi; \
		sleep 0.5; \
	done; \
	\
	@echo "==> Running tests..." ; \
	# `trap` to mechanizm powłoki, który wykonuje polecenie, gdy skrypt się kończy.
	# `EXIT` oznacza, że wykona się zawsze - po sukcesie, błędzie lub przerwaniu.
	trap 'echo "==> Stopping mock server (PID $$MOCK_PID)..."; kill $$MOCK_PID' EXIT; \
	\
	# Uruchamiamy testy. Zmienna środowiskowa jest ustawiona tylko dla tego polecenia.
	# Wynik (kod wyjścia) tego polecenia zadecyduje o sukcesie lub porażce całego kroku.
	RECUPERATOR_ADDRESS="http://localhost:8080" go test ./... -v

# Możesz dodać też prosty target do samego uruchomienia serwera
.PHONY: mock-server
mock-server:
	@echo "==> Starting mock server (press Ctrl+C to stop)..."
	@go run ./mocks/

# Installs the binary to the Go bin path
.PHONY: install
install:
	@echo "==> Installing..."
	go install .