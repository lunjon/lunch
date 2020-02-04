
test:
	go test ./internal/...


install: test
	go install ./cmd/lunch/

