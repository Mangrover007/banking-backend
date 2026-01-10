FROM golang:tip-trixie

WORKDIR /app

COPY ./app/go.mod ./
RUN go mod download

COPY ./app .
RUN go build -o ./bin/ ./cmd/main.go

ENV DB_USER=postgres
ENV DB_PASSWORD=pass
ENV DB_HOST=bank-db

EXPOSE 8080

CMD [ "./bin/main" ]

