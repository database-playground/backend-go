package database

import "github.com/samber/lo"

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
