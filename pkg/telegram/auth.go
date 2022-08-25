package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/snsvistunov/tg-bot-pocketer/pkg/repository"
)

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	authlink, err := b.generateAutorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(b.messages.Responses.Start, authlink))
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.TokenRepository.Get(chatID, repository.AccessTokens)
}

func (b *Bot) generateAutorizationLink(chatID int64) (string, error) {
	redirectURL := b.generateRedirectUrl(chatID)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), b.redirectURL)
	if err != nil {

		return "", err
	}
	if err := b.TokenRepository.Save(chatID, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

func (b *Bot) generateRedirectUrl(chatID int64) string {

	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
}
