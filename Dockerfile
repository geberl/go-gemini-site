FROM golang:1.18-alpine3.15 as builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o hn.bin ./cmd

FROM alpine:3.15
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/hn.bin hn
RUN mkdir /app/certs
EXPOSE 1965
ENTRYPOINT ["/app/hn"]
