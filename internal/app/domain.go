package app

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kurnosovmak/tb-bot-music/internal/config"
	tgevent "github.com/kurnosovmak/tb-bot-music/internal/handler/tg-event"
	"github.com/kurnosovmak/tb-bot-music/internal/handler/tg-event/start"
	vkcomaudio "github.com/kurnosovmak/tb-bot-music/internal/handler/tg-event/vkcom-audio"
	"github.com/kurnosovmak/tb-bot-music/internal/listener/tg-events/router"
)

func DI(
	ctx context.Context,
	cfg *config.Config,
	logger *slog.Logger,
	botApi *tgbotapi.BotAPI,
	router *router.Router,
) error {

	startHandler := start.NewHandler(botApi)
	vkAudioHandler := vkcomaudio.NewHandler(cfg.VK, botApi, logger)

	tgevent.InitRoutes(router, startHandler, vkAudioHandler)

	return nil
}
