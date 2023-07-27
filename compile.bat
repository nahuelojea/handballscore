set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build main.go
del bin/main.zip
tar.exe -a -cf bin/main.zip main