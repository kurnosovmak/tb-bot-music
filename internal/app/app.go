package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/kurnosovmak/tb-bot-music/internal/config"
	tgevents "github.com/kurnosovmak/tb-bot-music/internal/listener/tg-events"
	"github.com/kurnosovmak/tb-bot-music/internal/listener/tg-events/router"
	"github.com/kurnosovmak/tb-bot-music/internal/logger"
	tgapi "github.com/kurnosovmak/tb-bot-music/internal/tg-api"
)

func Run(ctx context.Context) error {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		return fmt.Errorf("error init config: %w", err)
	}

	logger, err := logger.NewLogger(cfg.Logger.Level)
	if err != nil {
		return fmt.Errorf("error init logger: %w", err)
	}

	logger.Debug("App config", slog.Any("cfg", cfg))

	botApi, err := tgapi.NewTgApiClient(cfg.Bot)
	if err != nil {
		return fmt.Errorf("error init tg bot api: %w", err)
	}

	logger.Info("Authorized on account", slog.String("username", botApi.Self.UserName))

	router := router.NewRouter()

	DI(ctx, cfg, logger, botApi, router)

	tgEventListener, err := tgevents.NewListener(cfg.Bot, logger, botApi, router.Handle)
	if err != nil {
		return fmt.Errorf("error init tg event listener: %w", err)
	}

	wg := &sync.WaitGroup{}

	wg.Go(func() {
		logger.Info("tg event listener start")
		err := tgEventListener.Run(ctx)
		if err != nil {
			logger.Error("error run tg event listener", slog.Any("err", err))
		}
		logger.Info("tg event listener stop")
	})

	wg.Wait()

	botApi.StopReceivingUpdates()
	return nil
}
