package ui

import (
	"github.com/rs/zerolog"
)

type UI struct {
	lg *zerolog.Logger
}

func New(logger *zerolog.Logger) UI {
	return UI{
		lg: logger,
	}
}
