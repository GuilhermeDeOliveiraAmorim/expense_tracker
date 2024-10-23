# Build Stage
FROM golang:1.21.5 AS build

WORKDIR /app

# Copia e baixa as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário de forma estática
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun

# Final Stage
FROM scratch

WORKDIR /app

# Copia o binário da etapa de build
COPY --from=build /app/cloudrun .

# Expõe a porta 8080
EXPOSE 8080

# Define o ENTRYPOINT
ENTRYPOINT ["./cloudrun"]
