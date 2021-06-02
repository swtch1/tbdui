package conf

import "github.com/gizak/termui/v3"

type Config struct {
	DefaultPrimaryColor   termui.Color
	DefaultSecondaryColor termui.Color
}

// NewDefault initializes a new default configuration.
func NewDefault() Config {
	return Config{
		DefaultPrimaryColor:   termui.ColorGreen,
		DefaultSecondaryColor: termui.ColorCyan,
	}
}
