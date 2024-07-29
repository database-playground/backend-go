package main

import (
	"github.com/database-playground/backend/internal/clients"
	slogmodule "github.com/database-playground/backend/internal/modules/slog"
	gatewayservice "github.com/database-playground/backend/internal/services/gateway"
	"go.uber.org/fx"
)

func main() {
	fx.New(slogmodule.FxOptions, clients.QuestionManagerClientFxModule, gatewayservice.FxModule).Run()
}
