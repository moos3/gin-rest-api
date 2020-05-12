FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN go mod download -x

RUN go get github.com/githubnemo/CompileDaemon

ENV port 5000

ENTRYPOINT CompileDaemon -exclude-dir=.git -exclude-dir=docs --build="go build main.go" --command=./main
