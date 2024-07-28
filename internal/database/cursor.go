package database

import (
	commonv1 "github.com/database-playground/backend/gen/common/v1"
	"github.com/samber/lo"
)

type Cursor struct {
	Offset int64
	Limit  int64
}

func (c Cursor) GetOffset() int64 {
	return c.Offset
}

func (c Cursor) GetLimit() int64 {
	return min(lo.CoalesceOrEmpty(c.Limit, 10), 100)
}

func CursorFromProto(cursor *commonv1.Cursor) Cursor {
	if cursor == nil {
		return Cursor{}
	}

	return Cursor{
		Limit:  cursor.GetLimit(),
		Offset: cursor.GetOffset(),
	}
}
