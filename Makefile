.PHONY:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t tg-bot-pocketer:0.1 .

start-container:
	docker run --env-file .env -p 8088:8088 tg-bot-pocketer:0.1