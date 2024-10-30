FROM golang:1.22

WORKDIR /consultant-microservices

COPY go.mod go.sum ./

RUN go mod download

COPY . .