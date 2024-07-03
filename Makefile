.SILENT:

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o proxy cmd/proxy/main.go

run: build
	./proxy