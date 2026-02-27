package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/kurnosovmak/tb-bot-music/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // очищаем ресурсы после выхода

	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}
