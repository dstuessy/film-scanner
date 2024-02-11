FROM golang:1.22-alpine

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]
