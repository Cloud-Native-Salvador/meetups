
FROM golang:1.24

WORKDIR /app
COPY cmd/main.go .

RUN go mod init app && go mod tidy && go build -o app main.go

CMD ["./app"]
