# Nombre de tu aplicación
APP_NAME=handballscore-app

# Nombre del archivo binario generado
BINARY_NAME=$(APP_NAME)

# Nombre de la función Lambda en AWS
LAMBDA_FUNCTION_NAME=handball-score

# Regla para compilar la aplicación de Go
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_NAME) main.go

# Regla para desplegar la función Lambda en AWS
deploy:
	zip $(BINARY_NAME).zip $(BINARY_NAME)
	aws lambda update-function-code --function-name $(LAMBDA_FUNCTION_NAME) --zip-file fileb://$(BINARY_NAME).zip

# Regla para limpiar los archivos binarios generados
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).zip
