package slogmodule

import (
	"log/slog"
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var FxOptions = fx.Options(
	fx.Provide(func() *slog.Logger {
		if os.Getenv("DEBUG") == "1" {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}

		return slog.Default()
	}),
	fx.WithLogger(func(slog *slog.Logger) fxevent.Logger {
		return &fxevent.SlogLogger{Logger: slog}
	}),
)
