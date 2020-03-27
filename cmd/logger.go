package cmd

import (
	"context"
	"io"
	"os"

	appctx "github.com/brave-intl/bat-go/utils/context"
	"github.com/rs/zerolog"
)

func setupLogger(ctx context.Context) (context.Context, *zerolog.Logger) {
	var output io.Writer
	if os.Getenv("ENV") != "local" {
		output = os.Stdout
	} else {
		output = zerolog.ConsoleWriter{Out: os.Stdout}
	}

	// always print out timestamp
	log := zerolog.New(output).With().Timestamp().Logger()

	debug := os.Getenv("DEBUG")
	if debug == "" || debug == "f" || debug == "n" || debug == "0" {
		log = log.Level(zerolog.InfoLevel)
	}

	ctx = context.WithValue(ctx, appctx.LoggerCTXKey, log)

	return log.WithContext(ctx), &log
}
