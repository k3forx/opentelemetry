FROM golang:1.22.1-bullseye

RUN mkdir /app
WORKDIR /app
COPY ./ .

RUN go install github.com/cosmtrek/air@latest
RUN go mod download

CMD ["air", "-c", ".air.toml"]
