FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o wallet-api ./cmd/server

CMD ["./wallet-api"]
