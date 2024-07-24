package dbrunner

import (
	lru "github.com/hashicorp/golang-lru/v2"
)

type ExecutionResultCache struct {
	*lru.Cache[Input, Output]
}

func NewExecutionResultCache(size int) (*ExecutionResultCache, error) {
	cache, err := lru.New[Input, Output](size)
	if err != nil {
		return nil, err
	}

	return &ExecutionResultCache{cache}, nil
}
