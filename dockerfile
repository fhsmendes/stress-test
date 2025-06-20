# Etapa de build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app .

# Etapa final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
ENTRYPOINT ["./app"]