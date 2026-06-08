FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o forum .


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/forum ./forum
COPY templates/ ./templates/
COPY static/ ./static/

EXPOSE 8080

CMD ["./forum"]
