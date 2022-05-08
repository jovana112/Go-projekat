FROM golang:latest AS builder

LABEL maintainer="Luka Matovic <lukamatovicfinal@gmail.com>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./main"]
