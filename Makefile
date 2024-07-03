.SILENT:

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o proxy cmd/proxy/main.go

run: build
	./proxy

docker:
	docker build . -t snmp-proxy && docker run -p 161:161 snmp-proxy