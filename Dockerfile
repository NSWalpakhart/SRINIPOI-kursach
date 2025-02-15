FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod init guess-game && \
    go mod tidy && \
    go build -o main .

EXPOSE 8888

CMD ["./main"]
