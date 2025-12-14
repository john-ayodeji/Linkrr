FROM golang:1.25.3
WORKDIR /app

COPY . .
#RUN go install github.com/pressly/goose/v3/cmd/goose@v3.25.0

ENV PORT=8080
ENV IPSTACK_URL="https://api.ipstack.com/"
ENV PLATFORM="docker"

CMD ["./linkrr"]
