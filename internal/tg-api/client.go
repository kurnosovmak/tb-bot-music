package tgapi

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kurnosovmak/tb-bot-music/internal/config"
)

func NewTgApiClient(cfg config.Bot) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("error init bot: %w", err)
	}
	bot.Debug = cfg.Debug
	return bot, nil
}
