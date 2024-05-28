package logging

import (
	"github.com/monorepo/common/configloader"
)

// Config hold the logger configuration.
type Config struct {
	Level      string `mapstructure:"level"`
	ShowCaller bool   `mapstructure:"show-caller"`
}

// Defaults sets default configuration values.
func (*Config) Defaults(l *configloader.Loader) {
	l.SetDefault("level", "info")
	l.SetDefault("show-caller", true)
}
