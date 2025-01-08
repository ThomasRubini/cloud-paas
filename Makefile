.PHONY: backend_docs backend_run

backend_docs:
	(cd internal/backend && swag init)
backend_run: backend_docs
	go run ./cmd/backend
