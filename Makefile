run:
	go run cmd/web/main.go

update:
	go mod tidy
	go mod vendor

.PHONY: run, update
