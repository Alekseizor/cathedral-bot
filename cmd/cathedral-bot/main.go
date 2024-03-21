package main

import (
	"context"
	"flag"
	"os"

	"github.com/Alekseizor/cathedral-bot/internal/app/config"
	"github.com/Alekseizor/cathedral-bot/internal/pkg/app"
	"github.com/rs/zerolog/log"
)

var (
	cfgPath = flag.String("config", "", "path to config file")
)

func main() {
	ctx := context.Background()

	cfg, err := config.Read(ctx, *cfgPath)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("")

		os.Exit(2)
	}

	// Создание приложения
	application := app.New(ctx, cfg)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("")

		os.Exit(2)
	}
	// Запуск приложения
	err = application.Run(ctx)
	if err != nil {
		log.Ctx(ctx).Error(err).Msg("app startup error")

		os.Exit(2)
	}
}
