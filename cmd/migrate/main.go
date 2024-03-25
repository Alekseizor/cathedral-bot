package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Alekseizor/cathedral-bot/internal/app/config"
)

const (
	migrationsPath = "migrations"
	driver         = "postgres"
	sslMode        = "disable"
)

var (
	cfgPath = flag.String("config", "", "path to config file")
)

func main() {
	ctx := context.Background()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	ctx = log.Logger.WithContext(ctx)

	log.Ctx(ctx).Info().Msg("Starting migrations")

	// считали конфиг
	flag.Parse()
	cfg, err := config.Read(ctx, *cfgPath)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("[config.Read]")

		os.Exit(2)
	}

	db, err := connect(cfg.ClientsConfig.PostgresConfig)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("[connect]")

		os.Exit(2)
	}

	// устанавливаем свой логер
	goose.SetLogger(&gooseLogger{ctx: ctx})

	// запускаем миграции
	log.Ctx(ctx).Info().Msg("Upping migrations")
	err = goose.SetDialect(driver)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("[goose.SetDialect]")

		os.Exit(2)
	}

	err = goose.Up(db.DB, migrationsPath)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("[goose.Up]")
	}

	log.Ctx(ctx).Info().Msg("DB migration completed")
}

// Выполняет подключение к БД
func connect(repoCfg config.PostgresConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", repoCfg.Host, repoCfg.Port, repoCfg.User, repoCfg.Password, repoCfg.Name, sslMode)
	log.Info().Msg(dsn)

	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("[sqlx.Connect]: %w", err)
	}
	return db, nil
}

// Реализация интерфйса goose.Logger
type gooseLogger struct {
	ctx context.Context
}

func (gl *gooseLogger) Fatal(v ...interface{}) {
	log.Fatal().Msgf("%v", v)
}
func (gl *gooseLogger) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v)
}
func (gl *gooseLogger) Print(v ...interface{}) {
	log.Info().Msgf("%v", v)
}
func (gl *gooseLogger) Println(v ...interface{}) {
	log.Info().Msgf("%v", v)
}
func (gl *gooseLogger) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v)
}
