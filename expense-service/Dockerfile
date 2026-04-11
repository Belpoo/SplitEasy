# Archivo: Dockerfile para Expense Service

FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o expense-svc main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/expense-svc .

EXPOSE 8081
# USE_BROKER=false por default si no se provee. ENV inyección Render
CMD ["./expense-svc"]
