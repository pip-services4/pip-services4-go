.PHONY: all build clean install uninstall fmt simplify check run test

install:
	@go install main.go

run: install
	@go run main.go

test:
	@go clean -testcache && go test  -v ./test/...
