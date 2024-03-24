package app

import (
	"context"
	"fmt"

	"github.com/Alekseizor/cathedral-bot/internal/app/config"
	"github.com/Alekseizor/cathedral-bot/internal/app/endpoint"
	"github.com/rs/zerolog/log"
)

type App struct {
	ctx      context.Context
	cfg      config.Config
	endpoint *endpoint.Endpoint
}

func New(ctx context.Context, cfg config.Config) *App {
	return &App{
		ctx:      ctx,
		cfg:      cfg,
		endpoint: endpoint.New(cfg),
	}
}

func (a *App) Init() error {
	err := a.endpoint.Init(a.ctx)
	if err != nil {
		return fmt.Errorf("[Endpoint.Init]: %w", err)
	}

	return nil
}

func (a *App) Run() error {
	log.Ctx(a.ctx).Info().Msg("[app.Run]: the application is running")

	err := a.endpoint.Run()
	if err != nil {
		return fmt.Errorf("[endpoint.Run]: %w", err)
	}

	return nil
}
