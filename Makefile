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
	go test test/integration_test.go

test: test-unit test-integration