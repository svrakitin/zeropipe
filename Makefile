.PHONY: test
test:
	CGO_ENABLED=0 go test -v ./...