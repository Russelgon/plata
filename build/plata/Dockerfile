FROM golang:1.24rc1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o myapp ./cmd/plata

EXPOSE 8080

CMD ["./myapp"]
