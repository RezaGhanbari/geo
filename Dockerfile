FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get github.com/gomodule/redigo/redis
RUN go get github.com/cnjack/throttle
RUN go get github.com/gin-gonic/gin
RUN go get github.com/gin-contrib/cache
RUN go get github.com/gin-contrib/cache/persistence
RUN go get github.com/satori/go.uuid

ENV API_TOKEN=WP5D&s3ftd^NU3TG@JH2n?!!@!MLmquD5t?V7vCPdANyY4Vrq5F \
    MAP_NAME=CEDAR \
    REDIS_URL=localhost:6379 \
#    SERVER=localhost \
    PORT=3001 \
    LIMIT=100
EXPOSE 3001:3001
RUN go build -o main .
CMD ["/app/main"]