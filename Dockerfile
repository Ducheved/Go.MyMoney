FROM golang:1.22.3-alpine3.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gomymoney ./main.go

FROM alpine:3.20
WORKDIR /app

ARG {BOT_TOKEN}
ARG {DATABASE_URL}

ENV BOT_TOKEN=${BOT_TOKEN}
ENV DATABASE_URL=${DATABASE_URL}

COPY --from=builder /app/gomymoney .
CMD ["./gomymoney"]