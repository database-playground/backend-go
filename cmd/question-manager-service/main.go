package main

import (
	"github.com/database-playground/backend/gen/questionmanager/v1/questionmanagerv1connect"
	"github.com/database-playground/backend/internal/database"
	httpservermodule "github.com/database-playground/backend/internal/modules/httpserver"
	slogmodule "github.com/database-playground/backend/internal/modules/slog"
	questionmanagerservice "github.com/database-playground/backend/internal/services/question_manager"
	"go.uber.org/fx"
)

func main() {
	fx.New(slogmodule.FxOptions, database.FxModule, questionmanagerservice.FxModule, fx.Provide(func(s *questionmanagerservice.Service) httpservermodule.HTTPHandler {
		return httpservermodule.WrapHTTPHandler[questionmanagerv1connect.QuestionManagerServiceHandler](questionmanagerv1connect.NewQuestionManagerServiceHandler, s)
	}), httpservermodule.FxModule).Run()
}
