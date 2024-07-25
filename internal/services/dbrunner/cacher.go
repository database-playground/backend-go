package dbrunnerservice

import (
	"context"
	"encoding/json"
	"time"

	"github.com/database-playground/backend/internal/dbrunner"
	"github.com/redis/go-redis/v9"
)

// dbrunner:sql-input:<input-hash> -> <output-hash>
// dbrunner:sql-output:<output-hash> -> <output-marshaled-json>
// <output-hash#1> == <output-hash#2> means the output is the same.

const (
	inputHashPrefix  = "dbrunner:sql-input:"
	outputHashPrefix = "dbrunner:sql-output:"
)

type CacheModule struct {
	redis *redis.Client
}

func NewCacheModule(redis *redis.Client) *CacheModule {
	return &CacheModule{
		redis: redis,
	}
}

// GetOutputHash returns the output hash of the input hash.
//
// The input hash is better to be [dbrunner.Input.Normalize]d.
func (c *CacheModule) GetOutputHash(ctx context.Context, inputHash string) (outputHash string, err error) {
	key := inputHashPrefix + inputHash

	result := c.redis.GetEx(ctx, key, time.Hour*1)
	if result.Err() == redis.Nil {
		return "", ErrNotFound
	}
	if result.Err() != nil {
		return "", result.Err()
	}

	outputHash, err = result.Result()
	if err != nil {
		return "", err
	}

	return outputHash, nil
}

// HasOutput returns whether the output hash is existed.
//
// It does not unmarshal the output, therefore it is much faster than
// [cacheModule.getOutput] if you don't need the output.
func (c *CacheModule) HasOutput(ctx context.Context, outputHash string) bool {
	key := outputHashPrefix + outputHash

	result := c.redis.GetEx(ctx, key, time.Hour*1)
	return result.Err() == nil
}

// GetOutput returns the output of the output hash.
//
// It is better to use [cacheModule.HasOutput] if you don't need the output.
func (c *CacheModule) GetOutput(ctx context.Context, outputHash string) (output *dbrunner.Output, err error) {
	key := outputHashPrefix + outputHash

	result := c.redis.GetEx(ctx, key, time.Hour*1)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, ErrNotFound
		}

		return output, result.Err()
	}

	err = json.Unmarshal([]byte(result.Val()), &output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// writeOutput writes (overrides) the output to the cache.
func (c *CacheModule) writeOutput(ctx context.Context, output dbrunner.Output) (hash string, err error) {
	hash, err = output.Hash()
	if err != nil {
		return "", err
	}

	key := outputHashPrefix + hash
	outputMarshaled, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	err = c.redis.SetEx(ctx, key, string(outputMarshaled), time.Hour*1).Err()
	if err != nil {
		return "", err
	}

	return hash, nil
}

// reuseOrWriteOutput reuses the output if it is existed, otherwise writes it.
func (c *CacheModule) reuseOrWriteOutput(ctx context.Context, output dbrunner.Output) (hash string, err error) {
	hash, err = output.Hash()
	if err != nil {
		return "", err
	}

	if c.HasOutput(ctx, hash) {
		return hash, nil
	}

	return c.writeOutput(ctx, output)
}

// writeInput writes the input with its output hash map to the cache.
func (c *CacheModule) writeInput(ctx context.Context, input dbrunner.Input, outputHash string) (hash string, err error) {
	normalizedInput, err := input.Normalize()
	if err != nil {
		return "", err
	}

	hash = normalizedInput.Hash()
	key := inputHashPrefix + hash

	err = c.redis.SetEx(ctx, key, outputHash, time.Hour*1).Err()
	if err != nil {
		return "", err
	}

	return hash, nil
}

// WriteToCache writes the input and output to the cache.
func (c *CacheModule) WriteToCache(ctx context.Context, input dbrunner.Input, output dbrunner.Output) (hash string, err error) {
	outputHash, err := c.reuseOrWriteOutput(ctx, output)
	if err != nil {
		return "", err
	}

	return c.writeInput(ctx, input, outputHash)
}
