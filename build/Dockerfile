FROM golang:1.22-alpine

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod download

RUN go install github.com/cosmtrek/air@latest

CMD ["air", "-c", "build/air.toml"]
