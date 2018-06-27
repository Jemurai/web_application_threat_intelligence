FROM golang:1.10

WORKDIR /src
COPY . .
RUN go get github.com/fzzy/radix/redis

CMD ["go", "run", "main.go"]