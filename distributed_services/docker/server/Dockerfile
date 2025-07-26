FROM golang:1.24.1-bullseye

RUN mkdir /app
WORKDIR /app
COPY ./gin .

RUN go install github.com/air-verse/air@latest
RUN go mod download

CMD ["air", "-c", ".air.toml"]
