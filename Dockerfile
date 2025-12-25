FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o lofam ./cmd/server

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/lofam .

ENV PORT=8080
ENV DB_PATH=/data/lofam.db

EXPOSE 8080

VOLUME ["/data"]

CMD ["./lofam"]
