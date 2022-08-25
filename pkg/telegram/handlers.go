package telegram

import (
	"context"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
)

// const (
// 	commandStart           = "start"
// 	replyStartTemplate     = "Hi! To save your links in a Pocket account you should give me an access for this. Follow the link:\n%s"
// 	replyAlreadyAuthorized = "You are already authorized. Send me a link I will save!"
// 	replyLinkSaved         = "Link saved successfully"
// )

func (b *Bot) handleMessage(message *tgbotapi.Message) error {

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return errunAuthorized
	}

	_, err = url.ParseRequestURI(message.Text)
	if err != nil {
		return errInvalidURL
	}

	if err = b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		return errUnableToSave
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.SavedSuccessfully)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {

	switch message.Command() {
	case b.commands.Start:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.AlreadyAuthorized)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.UnknownCommand)
	_, err := b.bot.Send(msg)
	return err
}
