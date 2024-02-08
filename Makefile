APP_NAME=handballscore-app

LAMBDA_FUNCTION_NAME=handball-score

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go
	mkdir -p bin
	rm -f bin/$(APP_NAME).zip

deploy:
	zip -r bin/$(APP_NAME).zip main
	aws lambda update-function-code --function-name $(LAMBDA_FUNCTION_NAME) --zip-file fileb://bin/$(APP_NAME).zip
