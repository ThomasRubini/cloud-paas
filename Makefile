.PHONY: backend_docs backend_run

backend_docs:
	(cd cmd/backend && swag init)
backend_run: backend_docs
	go run ./cmd/backend
