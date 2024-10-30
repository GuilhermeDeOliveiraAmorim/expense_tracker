# Build Stage
FROM golang:1.23.2 AS builder

RUN apt-get update && apt-get install -y \
    tesseract-ocr \
    libleptonica-dev \
    libtesseract-dev \
    && apt-get clean

# Crie o diretório do aplicativo
WORKDIR /app

# Copia e instala dependências Go
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário
RUN go build -o app .

# Defina o ponto de entrada
CMD ["./app"]
