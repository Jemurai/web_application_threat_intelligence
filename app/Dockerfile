FROM golang:1.10

WORKDIR /go/src/app
COPY . .
RUN go get github.com/gorilla/context
RUN go get github.com/gorilla/handlers
RUN go get github.com/gorilla/mux

EXPOSE 8080

CMD ["go", "run", "app.go"]
