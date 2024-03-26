package config

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config - структура конфигурации.
// Содержит все конфигурационные данные о сервисе
type Config struct {
	BotConfig        BotConfig        `yaml:"bot" mapstructure:"bot"`
	ClientsConfig    ClientsConfig    `yaml:"clients" mapstructure:"clients"`
	MonitoringConfig MonitoringConfig `yaml:"listener" mapstructure:"listener"`
}

// BotConfig -  конфигурация бота в VK
type BotConfig struct {
	Token string `yaml:"token" mapstructure:"token"`
}

// ClientsConfig - конфигурация клиентов
type ClientsConfig struct {
	// PostgresConfig - клиент СУБД
	PostgresConfig PostgresConfig `yaml:"postgres" mapstructure:"postgres"`
}

// MonitoringConfig - конфигурация для мониторинговых систем
type MonitoringConfig struct {
	// Address - адрес, по которому будем принимать запросы касательно метрик
	Address string `yaml:"address" mapstructure:"address"`
}

// PostgresConfig - конфигурация для клиента PostgreSQL
type PostgresConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"pass" mapstructure:"pass"`
	Name     string `yaml:"name" mapstructure:"name"`
}

func Read(ctx context.Context, path string) (Config, error) {
	v := viper.New()

	if path == "" {
		return Config{}, fmt.Errorf("read: empty configuration path")
	}

	v.SetConfigType("yaml")
	v.SetConfigFile(path)
	v.WatchConfig()

	err := v.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("read: ReadInConfig error: %w", err)
	}

	cfg := &Config{}
	err = v.Unmarshal(cfg)
	if err != nil {
		return Config{}, fmt.Errorf("read: Unmarshal error: %w", err)
	}

	log.Ctx(ctx).Debug().Interface("cfg", cfg).Msg("config parsed")

	return *cfg, nil
}
