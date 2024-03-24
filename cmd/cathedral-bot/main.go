package main

import (
	"context"
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Alekseizor/cathedral-bot/internal/app/config"
	"github.com/Alekseizor/cathedral-bot/internal/pkg/app"
)

var (
	cfgPath = flag.String("config", "", "path to config file")
)

func main() {
	ctx := context.Background()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	ctx = log.Logger.WithContext(ctx)

	// считали конфиг
	flag.Parse()
	cfg, err := config.Read(ctx, *cfgPath)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("[config.Read]")

		os.Exit(2)
	}

	// Создание приложения
	application := app.New(ctx, cfg)

	err = application.Init()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("[application.Init]")

		os.Exit(2)
	}

	// Запуск приложения
	err = application.Run()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("[application.Run]")

		os.Exit(2)
	}
}
