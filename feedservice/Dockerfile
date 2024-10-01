FROM golang:1.22.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o feedservice

EXPOSE 8081

CMD ["./feedservice"]