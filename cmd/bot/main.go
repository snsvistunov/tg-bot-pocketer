package main

import (
	"log"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/snsvistunov/tg-bot-pocketer/pkg/config"
	"github.com/snsvistunov/tg-bot-pocketer/pkg/repository"
	"github.com/snsvistunov/tg-bot-pocketer/pkg/repository/boltdb"
	"github.com/snsvistunov/tg-bot-pocketer/pkg/server"
	"github.com/snsvistunov/tg-bot-pocketer/pkg/telegram"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.ConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerURL, cfg.Messages, cfg.Commands)

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg.TelegramBotURL)

	go func() {
		if err := authorizationServer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	telegramBot.Start()
}

func initDB(cfg *config.Config) (*bolt.DB, error) {

	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
