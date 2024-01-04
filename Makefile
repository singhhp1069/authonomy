# Variables
BINARY_FOLDER = build
BINARY_NAME = authonomy
DOCKER_IMAGE_NAME = authonomy-image

# Run the server
run:
	@echo "Running the server..."
	@go run main.go start

# Build the Go binary
build:
	@echo "Building the Go binary..."
	@go build -o $(BINARY_FOLDER)/$(BINARY_NAME)

# Lint the Go code
lint:
	@echo "Linting the Go code..."
	@golangci-lint run

## TESTS
TEST_PACKAGES=$(shell go list ./...)
TEST_TARGETS := test-unit test-race
BASE_FLAGS=-mod=readonly -timeout=5m
test-unit: ARGS=-tags=norace
test-race: ARGS=-race
$(TEST_TARGETS): run-tests

run-tests:
	@echo "--> Running tests $(BASE_FLAGS) $(ARGS)"
ifneq (,$(shell which tparse 2>/dev/null))
	@go test $(BASE_FLAGS) -json $(ARGS) $(TEST_PACKAGES) | tparse
else
	@go test $(BASE_FLAGS) $(ARGS) $(TEST_PACKAGES)
endif


# Build a Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE_NAME) .

# Run the application inside a Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 7322:7322 $(DOCKER_IMAGE_NAME)

# Clean up
clean:
	@echo "Cleaning up..."
	@rm $(BINARY_NAME)

.PHONY: run build lint docker-build docker-run clean
