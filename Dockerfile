FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add make

COPY go.mod go.sum ./

RUN go mod download

COPY internal internal
COPY cmd cmd
COPY Makefile .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    make backend

RUN adduser -D -g '' user

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
USER user

WORKDIR /app
COPY --from=builder /app/backend /app/backend

EXPOSE 8080
CMD ["/app/backend"]
