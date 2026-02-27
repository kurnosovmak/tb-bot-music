package tgevents

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kurnosovmak/tb-bot-music/internal/config"
)

type HandleFunc func(ctx context.Context, newMessage tgbotapi.Message) error

type Listener struct {
	botApi         *tgbotapi.BotAPI
	newMessageFunc HandleFunc
	cfg            config.Bot
	logger         *slog.Logger
}

func NewListener(cfg config.Bot, logger *slog.Logger, botApi *tgbotapi.BotAPI, newMessageFunc HandleFunc) (*Listener, error) {
	return &Listener{
		botApi:         botApi,
		newMessageFunc: newMessageFunc,
		cfg:            cfg,
		logger:         logger,
	}, nil
}

func (l *Listener) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = l.cfg.EventLongpollTimeout

	updates := l.botApi.GetUpdatesChan(u)

	wg := &sync.WaitGroup{}

OuterLoop:
	for {
		select {
		case <-ctx.Done():
			break OuterLoop
		case update, isOpen := <-updates:
			if !isOpen {
				return errors.New("updated chan closed")
			}
			if update.Message == nil || update.Message.Text == "" {
				continue
			}

			wg.Go(func() {
				defer func() {
					if r := recover(); r != nil {
						l.logger.Error(
							"panic handler",
							slog.Any("err", r),
							slog.String("trace", string(debug.Stack())),
						)
					}
				}()
				requestCtx, cancel := context.WithTimeout(ctx, time.Duration(l.cfg.Timeout)*time.Second)
				defer cancel()
				if err := l.newMessageFunc(requestCtx, *update.Message); err != nil {
					l.logger.Error("error handle", slog.Any("err", err), slog.Any("message", *update.Message))
				} else {
					l.logger.Info("ok handle", slog.Any("message", *update.Message))
				}
			})

		}
	}
	wg.Wait()
	return nil
}
