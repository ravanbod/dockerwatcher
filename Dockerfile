FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o build/dockerwatcher cmd/dockerwatcher/main.go

CMD ["./build/dockerwatcher"]