.PHONY: backend_docs cli backend
swag = go run github.com/swaggo/swag/cmd/swag@v1.16.4

backend_docs:
	git clean -fdx internal/paas_backend/docs
	($(swag) init -d internal -o internal/paas_backend/docs)
backend: backend_docs
	go build -o backend ./cmd/backend
cli:
	go build -o cli ./cmd/cli
