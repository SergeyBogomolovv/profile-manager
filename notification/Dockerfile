FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o notifications cmd/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /app/notifications /notifications

USER nonroot:nonroot

CMD ["/notifications"]