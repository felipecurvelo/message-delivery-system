test-integration:
	go clean -testcache
	go test test/integration_test.go

test: test-integration