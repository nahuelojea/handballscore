build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go

copy:
	mkdir -p bin
	rm -f bin/main.zip
	zip -r bin/main.zip main
