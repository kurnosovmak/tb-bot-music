package start

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kurnosovmak/tb-bot-music/internal/text"
)

var (
	startMessage string
)

func init() {
	startMessage = text.MustLoad("start.html")
}

type Handler struct {
	botApi *tgbotapi.BotAPI
}

func NewHandler(botApi *tgbotapi.BotAPI) *Handler {
	return &Handler{
		botApi: botApi,
	}
}

func (h *Handler) Handle(ctx context.Context, newMessage tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(newMessage.Chat.ID, startMessage)
	msg.ParseMode = "html"

	_, err := h.botApi.Send(msg)
	return err
}
