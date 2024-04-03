FROM golang:1.22.1-bullseye

RUN mkdir /app
WORKDIR /app
COPY ./ .

RUN go install github.com/githubnemo/CompileDaemon@latest
RUN go mod download

ENTRYPOINT CompileDaemon --build="go build -o main ./" --command=./main
