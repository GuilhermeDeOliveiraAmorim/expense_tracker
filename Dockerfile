# Build Stage
FROM golang:1.23.2 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /build/app .

# Final Stage
FROM golang:1.23.2

# Crie o diretório /app se ele não existir
RUN mkdir -p /app

COPY --from=builder /build/app /app

CMD ["/app/app"]
