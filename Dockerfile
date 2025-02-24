FROM golang:1.24.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app

FROM debian:trixie-slim AS release
RUN apt-get -y update && apt-get -y upgrade
RUN apt-get install -y ffmpeg
COPY --from=builder /app /app
COPY .env ./

EXPOSE 8080
CMD [ "/app/convert-thing" ]
