.PHONY: build-ConnectraApiFunction build-lambda clean test deploy help

# SAM expects this target name format: build-{FunctionName}
# This target is used when BuildMethod: makefile is specified in template.yaml
# The binary must be copied to $(ARTIFACTS_DIR) for SAM to package it
build-ConnectraApiFunction:
	@echo "Building Lambda function for ConnectraApiFunction..."
	@GOOS=linux GOARCH=amd64 go build -o bootstrap \
		-ldflags="-s -w" \
		./cmd/lambda/main.go
	@if [ -n "$$ARTIFACTS_DIR" ]; then \
		cp bootstrap $$ARTIFACTS_DIR; \
		echo "Copied bootstrap to $$ARTIFACTS_DIR"; \
	else \
		echo "Build complete! Binary: bootstrap (not in SAM build context)"; \
	fi

# Build Lambda function binary (alias for convenience)
# Use this for manual builds outside of SAM
build-lambda:
	@echo "Building Lambda function (manual build)..."
	@GOOS=linux GOARCH=amd64 go build -o bootstrap \
		-ldflags="-s -w" \
		./cmd/lambda/main.go
	@echo "Build complete! Binary: bootstrap"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f bootstrap
	@rm -rf .aws-sam
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Build and deploy using SAM
deploy: build-lambda
	@echo "Deploying with SAM..."
	@sam build
	@sam deploy

# Help target
help:
	@echo "Available targets:"
	@echo "  build-lambda  - Build Lambda function binary"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  deploy        - Build and deploy with SAM"
