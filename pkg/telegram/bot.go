package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/snsvistunov/tg-bot-pocketer/pkg/config"
	"github.com/snsvistunov/tg-bot-pocketer/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	TokenRepository repository.TokenRepository
	redirectURL     string

	messages config.Messages
	commands config.Commands
}

func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocket.Client, tr repository.TokenRepository, redirectURL string, messages config.Messages, commands config.Commands) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, redirectURL: redirectURL, TokenRepository: tr, messages: messages, commands: commands}
}

func (b *Bot) Start() {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates := b.initUpdatesChannel()
	b.handleUpdates(updates)
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // If we didn't got a message
			continue
		}
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}

		}
		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
}
