# Cloud-PaaS
This is the PaaS project for the DO3 Cloud Technologies subject

Test coverage: [![codecov](https://codecov.io/gh/ThomasRubini/cloud-paas/graph/badge.svg?token=40TRSMIVVE)](https://codecov.io/gh/ThomasRubini/cloud-paas)

## [Project Deployment](./deployment/README.md)

## [Assets Documentation](./Assets/README.md)

## [Enduser documentation](./cmd/README.md)


# Seting up the the developpement environment

## Set up pre-commit
- `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- `pipx install pre-commit`
- `pre-commit install`

## How to Run

### Backend
- Start the database: ``docker compose up -d``
- Configure `.env`
- Start the backend: `./run_backend.sh`

### CLI
- Configure `paas_cli_config.yml`
- Run `./run_cli.sh`  
