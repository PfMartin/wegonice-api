package logging

import (
	"github.com/rs/zerolog"
)

func NewLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
