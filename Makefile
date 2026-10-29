.PHONY: run check test
run:
	go run ./cmd/app/main.go
check:
	go build ./...
test:
	go test -race ./...