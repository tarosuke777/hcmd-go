# Doc

https://go.dev/doc/

# Build Command

go build -o home.exe cmd/main.go
GOOS=windows GOARCH=amd64 go build -o home.exe cmd/main.go

go build -o home cmd/main.go
GOOS=linux GOARCH=amd64 go build -o home cmd/main.go

# Web Upload Command

curl -X POST -F "myImage=@touka_2.jpg" http://localhost:8080/images/

# Jenkins Build Command With Parameters
home jenkins build hms-build
