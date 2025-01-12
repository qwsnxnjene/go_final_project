FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go_final_project .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/go_final_project .
COPY --from=builder /app/.env .

EXPOSE 7540

CMD ["./go_final_project"]