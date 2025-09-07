FROM golang:1.24-alpine

# Cria um diretório de trabalho dentro do container
WORKDIR /app

# Copia go.mod e go.sum primeiro (boa prática p/ cache das dependências)
COPY go.mod go.sum ./

RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário a partir do main.go
RUN go build -o watcher ./cmd/api

# Comando padrão ao rodar o container
CMD ["./watcher"]