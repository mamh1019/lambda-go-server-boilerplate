APP_NAME := lambda-go-server-boilerplate
LAMBDA_DIR := lambda
BUILD_DIR := build
BINARY_NAME := bootstrap

# Local run (always using Air via go run)
.PHONY: run
run:
	go run github.com/air-verse/air@latest -c .air.toml

# Build Lambda artifact (Linux) into build/function.zip
.PHONY: deploy
deploy: clean
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(LAMBDA_DIR)
	cd $(BUILD_DIR) && zip function.zip $(BINARY_NAME)

# Build & upload to AWS Lambda (requires aws CLI and FUNCTION_NAME env)
.PHONY: deploy-aws
deploy-aws: deploy
	aws lambda update-function-code --function-name $$FUNCTION_NAME --zip-file fileb://$(BUILD_DIR)/function.zip

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

