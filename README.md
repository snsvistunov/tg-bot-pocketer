# Telegram Bot Client for Pocket  

## General

Simple Telegram Bot Client written in Golang for [Pocket](https://getpocket.com/)

Designed to quickly add links to your Pocket account via [Telegram](https://web.telegram.org/).

## Stack

 - Go 1.18
 - Docker
 - BoltDB
 - Makefile

## Makefile

`Makefile` - contains commands to build and run application.

```
build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t tg-bot-pocketer:0.1 .

start-container:
	docker run --env-file .env -p 8088:8088 tg-bot-pocketer:0.1

```

`make build` - build application from binaries.

`make run` - build and run application from binaries.

`make build-image` - build Docker Image.

`start-container` - start Docker container.

##Config

Your should use `.env` variables.


`TOKEN` - Telegram bot token from BotFather. 

`CONSUMER_KEY` Consumer key for your application recieved from Pocket. 

`AUTH_SERVER_URL` - your application server address, e.g `http://localhost:8088/`.

###Additional

[Pocket API: Documentation](https://getpocket.com/developer/docs/overview).


