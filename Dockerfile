# Build Stage
FROM golang:1.21.5 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g ./main.go -o ./api
    
RUN go build -o /build/app .

# Final Stage
FROM golang:1.21.5

# Crie o diretório /app se ele não existir
RUN mkdir -p /app

COPY --from=builder /build/app /app

CMD ["/app/app"]
