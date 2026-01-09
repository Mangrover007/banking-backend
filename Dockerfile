FROM golang:tip-trixie
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o bin cmd/main.go
CMD ["./bin/main"]

