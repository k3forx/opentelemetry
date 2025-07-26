FROM golang:1.24.1-bullseye

RUN mkdir /app
WORKDIR /app

# Copy the entire project
COPY . .

RUN go mod download && \
    go clean -modcache && \
    rm -rf /tmp/* /var/tmp/* /root/.cache

WORKDIR /app/services/author-service

# Set Go environment variables to prevent tmp directory creation
RUN mkdir -p /tmp/go-cache /tmp/go-tmp
ENV GOCACHE=/tmp/go-cache
ENV GOTMPDIR=/tmp/go-tmp

CMD ["go", "run", "main.go"] 
