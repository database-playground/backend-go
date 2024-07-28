package questionmanagerservice

import (
	"github.com/database-playground/backend/gen/questionmanager/v1/questionmanagerv1connect"
	"github.com/database-playground/backend/internal/database"
	"github.com/database-playground/backend/internal/models"
	"github.com/database-playground/backend/internal/models/generated"
	"go.uber.org/fx"
)

var FxModule = fx.Module("question-manager-service", fx.Provide(New))

type Service struct {
	questionmanagerv1connect.UnimplementedQuestionManagerServiceHandler

	db        *database.Database
	converter models.Converter
}

func New(database *database.Database) *Service {
	return &Service{
		db:        database,
		converter: &generated.ConverterImpl{},
	}
}
