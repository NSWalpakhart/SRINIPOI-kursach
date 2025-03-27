FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum* ./
COPY *.go ./
COPY templates/ ./templates/

RUN go mod tidy && \
    go build -o main .

EXPOSE 8888

CMD ["./main"]
