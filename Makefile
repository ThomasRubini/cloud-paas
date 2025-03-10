.PHONY: backend_docs cli backend
swag = go run github.com/swaggo/swag/cmd/swag@v1.16.4

backend_docs:
	($(swag) init -d internal/paas_backend -d internal/comm -g mod.go -o internal/paas_backend/docs)
backend: backend_docs
	go build -o backend ./cmd/backend
cli:
	go build -o cli ./cmd/cli
