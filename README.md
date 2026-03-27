# Doc

https://go.dev/doc/

# Build Command

go build -o home.exe cmd/main.go

go build -ldflags="-H windowsgui" -o home.exe ./cmd/main.go

# Web Upload Command

curl -X POST -F "myImage=@touka_2.jpg" http://localhost:8080/images/