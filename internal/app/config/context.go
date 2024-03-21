package config

import "context"

type configContextKeyType struct {
}

var configContextKey = configContextKeyType{}

// FromContext Возвращает конфиг из контекста
func FromContext(ctx context.Context) *Config {
	cfgRaw := ctx.Value(configContextKey)
	cfg, ok := cfgRaw.(*Config)
	if ok {
		return cfg
	}
	return nil
}

// WrapContext Обогащает контекст конфигом
func WrapContext(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, configContextKey, cfg)
}
