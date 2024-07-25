package slogmodule

import (
	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var FxOptions = fx.Options(
	fx.Provide(func() *slog.Logger {
		return slog.Default()
	}),
	fx.WithLogger(func(slog *slog.Logger) fxevent.Logger {
		return &fxevent.SlogLogger{Logger: slog}
	}),
)
