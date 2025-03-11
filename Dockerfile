FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add make

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make backend

FROM scratch
WORKDIR /app
COPY --from=builder /app/backend /app/backend

EXPOSE 8080
CMD ["/app/backend"]
