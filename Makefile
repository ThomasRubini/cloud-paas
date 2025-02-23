.PHONY: backend_docs cli backend
swag = go run github.com/swaggo/swag/cmd/swag@latest

backend_docs:
	(cd internal/paas_backend && $(swag) init)
backend: backend_docs
	go build -o backend ./cmd/backend
cli:
	go build -o cli ./cmd/cli
