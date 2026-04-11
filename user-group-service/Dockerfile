# Dockerfile del User & Group Service

# Imagen base oficial constructora (multistage base compilada Go)
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copia de todo y construcción binaria final (peso minúsculo e hyper rápido)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-group-svc main.go

# Imagen liviana ejecutoria
FROM alpine:latest
WORKDIR /app
# Trae el bin de la etapa 'builder'
COPY --from=builder /app/user-group-svc .

# Expose su puerto (predet. 8080)
EXPOSE 8080
CMD ["./user-group-svc"]
