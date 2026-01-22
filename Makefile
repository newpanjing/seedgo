build:
	CGO_ENABLED=0 go build -p 10 -ldflags="-s -w" -o main main.go
