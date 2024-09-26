APP_NAME=handballscore-app
LAMBDA_FUNCTION_NAME=handball-score

# Build process
build:
	@echo "Starting build process for $(APP_NAME)..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go
	@echo "Build completed. Preparing binary..."
	mkdir -p bin
	@echo "Removing old zip file if it exists..."
	rm -f bin/$(APP_NAME).zip
	@echo "Creating new zip archive for $(APP_NAME)..."
	zip -r bin/$(APP_NAME).zip main
	@echo "Build process completed. The zip file is ready at bin/$(APP_NAME).zip"

# Deployment process
deploy:
	@echo "Starting deployment to AWS Lambda for function $(LAMBDA_FUNCTION_NAME)..."
	aws lambda update-function-code --function-name $(LAMBDA_FUNCTION_NAME) --zip-file fileb://bin/$(APP_NAME).zip
	@echo "Deployment completed successfully for $(LAMBDA_FUNCTION_NAME)."
