FROM golang:1.18-alpine AS builder

COPY . /github.com/snsvistunov/tg-bot-pocketer/
WORKDIR /github.com/snsvistunov/tg-bot-pocketer/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM golang:1.18-alpine

WORKDIR /root/

COPY --from=builder /github.com/snsvistunov/tg-bot-pocketer/bin/bot .
COPY --from=builder /github.com/snsvistunov/tg-bot-pocketer/configs configs/

EXPOSE 8088

CMD ["./bot"]