FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o sso cmd/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /app/sso /sso

USER nonroot:nonroot

CMD ["/sso"]