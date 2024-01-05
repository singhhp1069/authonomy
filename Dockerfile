FROM golang:1.21-alpine3.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/authonomy .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/build/authonomy .

COPY --from=builder /app/config.yaml .

COPY --from=builder /app/web ./web

# sample schemas
COPY --from=builder /app/ssi/schemas ./ssi/schemas

EXPOSE 8081

CMD ["./authonomy"]
