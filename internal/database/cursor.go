package database

import "github.com/samber/lo"

type Cursor struct {
	Offset int
	Limit  int
}

func (c Cursor) GetOffset() int {
	return c.Offset
}

func (c Cursor) GetLimit() int {
	return min(lo.CoalesceOrEmpty(c.Limit, 10), 100)
}
