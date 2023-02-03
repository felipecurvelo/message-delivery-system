build-client:
	go build -o ./build/client ./cmd/client

build-server:
	go build -o ./build/server ./cmd/server

build: build-server build-client

test-unit:
	go clean -testcache
	go test internal/...

test-integration:
	go clean -testcache
	go test --race test/integration_test.go

test: test-unit test-integration

start-server: build-server
	go run ./cmd/server 1234

start-client: build-client
	go run ./cmd/client 127.0.0.1 1234