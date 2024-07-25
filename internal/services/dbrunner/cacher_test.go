package dbrunnerservice_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/database-playground/backend/internal/dbrunner"
	dbrunnerservice "github.com/database-playground/backend/internal/services/dbrunner"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheModule_GetOutputHash(t *testing.T) {
	t.Parallel()

	t.Run("if there is no such output hash, returns not found", func(t *testing.T) {
		t.Parallel()

		client, mock := redismock.NewClientMock()
		cm := dbrunnerservice.NewCacheModule(client)
		mock.ExpectGetEx("dbrunner:sql-input:input-hash", 1*time.Hour).SetErr(redis.Nil)

		_, err := cm.GetOutputHash(context.TODO(), "input-hash")

		assert.ErrorIs(t, err, dbrunnerservice.ErrNotFound)
	})

	t.Run("if there is an output hash, returns it", func(t *testing.T) {
		t.Parallel()

		client, mock := redismock.NewClientMock()
		cm := dbrunnerservice.NewCacheModule(client)
		mock.ExpectGetEx("dbrunner:sql-input:input-hash", 1*time.Hour).SetVal("output-hash")

		outputHash, err := cm.GetOutputHash(context.TODO(), "input-hash")

		assert.NoError(t, err)
		assert.Equal(t, "output-hash", outputHash)
	})
}

func TestCacheModule_HasOutput(t *testing.T) {
	t.Parallel()

	t.Run("if there is no such output hash, returns false", func(t *testing.T) {
		t.Parallel()

		client, mock := redismock.NewClientMock()
		cm := dbrunnerservice.NewCacheModule(client)
		mock.ExpectGetEx("dbrunner:sql-output:output-hash", 1*time.Hour).SetErr(redis.Nil)

		hasOutput := cm.HasOutput(context.TODO(), "output-hash")

		assert.False(t, hasOutput)
	})

	t.Run("if there is an output hash, returns true", func(t *testing.T) {
		t.Parallel()

		client, mock := redismock.NewClientMock()
		cm := dbrunnerservice.NewCacheModule(client)
		mock.ExpectGetEx("dbrunner:sql-output:output-hash", 1*time.Hour).SetVal("output")

		hasOutput := cm.HasOutput(context.TODO(), "output-hash")

		assert.True(t, hasOutput)
	})
}

func TestCacheModule_GetOutput(t *testing.T) {
	t.Parallel()

	t.Run("if there is output hash, returns the unmarshaled value", func(t *testing.T) {
		t.Parallel()

		expected := dbrunner.Output{
			Header: []string{"column", "column2"},
			Data: [][]*string{
				{nil, lo.ToPtr("Hello!")},
			},
		}

		client, mock := redismock.NewClientMock()
		cm := dbrunnerservice.NewCacheModule(client)
		mock.ExpectGetEx("dbrunner:sql-output:output-hash", 1*time.Hour).SetVal(string(lo.Must(json.Marshal(expected))))

		actual, err := cm.GetOutput(context.TODO(), "output-hash")
		require.NoError(t, err)

		// convert to json to workaround the nil issue
		expectedJSON, _ := json.Marshal(expected)
		actualJSON, _ := json.Marshal(actual)
		assert.JSONEq(t, string(expectedJSON), string(actualJSON))
	})

	t.Run("if there is no such output hash, returns not found", func(t *testing.T) {
		t.Parallel()

		client, mock := redismock.NewClientMock()
		cm := dbrunnerservice.NewCacheModule(client)
		mock.ExpectGetEx("dbrunner:sql-output:output-hash", 1*time.Hour).SetErr(redis.Nil)

		_, err := cm.GetOutput(context.TODO(), "output-hash")

		assert.ErrorIs(t, err, dbrunnerservice.ErrNotFound)
	})
}

func TestWriteToCache(t *testing.T) {
	mockInput, _ := dbrunner.Input{
		Init:  "CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT); INSERT INTO test (name) VALUES ('Hello!');",
		Query: "SELECT * FROM test;",
	}.Normalize()
	mockOutput := dbrunner.Output{
		Header: []string{"id", "name"},
		Data:   [][]*string{{lo.ToPtr("1"), lo.ToPtr("Hello!")}},
	}

	mockInputHash := mockInput.Hash()
	mockOutputHash, err := mockOutput.Hash()
	require.NoError(t, err)

	t.Run("if there is no input and output hash, we write ours", func(t *testing.T) {
		t.Parallel()

		client, mock := redismock.NewClientMock()
		cm := dbrunnerservice.NewCacheModule(client)

		mock.MatchExpectationsInOrder(false)
		mock.ExpectGetEx("dbrunner:sql-input:"+mockInputHash, 1*time.Hour).SetErr(redis.Nil)
		mock.ExpectGetEx("dbrunner:sql-output:"+mockOutputHash, 1*time.Hour).SetErr(redis.Nil)
		mock.ExpectSetEx("dbrunner:sql-input:"+mockInputHash, mockOutputHash, 1*time.Hour).SetVal("OK")
		mock.ExpectSetEx("dbrunner:sql-output:"+mockOutputHash, string(lo.Must(json.Marshal(mockOutput))), 1*time.Hour).SetVal("OK")

		hash, err := cm.WriteToCache(context.TODO(), mockInput, mockOutput)
		require.NoError(t, err)
		assert.Equal(t, mockInputHash, hash)
	})

	t.Run("if there is input hash but no output hash, we override the current values", func(t *testing.T) {
		t.Parallel()

		client, mock := redismock.NewClientMock()
		cm := dbrunnerservice.NewCacheModule(client)

		mock.MatchExpectationsInOrder(false)
		mock.ExpectGetEx("dbrunner:sql-input:"+mockInputHash, 1*time.Hour).SetVal(mockOutputHash)
		mock.ExpectGetEx("dbrunner:sql-output:"+mockOutputHash, 1*time.Hour).SetErr(redis.Nil)
		mock.ExpectSetEx("dbrunner:sql-input:"+mockInputHash, mockOutputHash, 1*time.Hour).SetVal("OK")
		mock.ExpectSetEx("dbrunner:sql-output:"+mockOutputHash, string(lo.Must(json.Marshal(mockOutput))), 1*time.Hour).SetVal("OK")

		hash, err := cm.WriteToCache(context.TODO(), mockInput, mockOutput)
		require.NoError(t, err)
		assert.Equal(t, mockInputHash, hash)
	})

	// FIXME: Usually we don't need to update the output hash
	t.Run("if there is input and output hash, we only update the input hash", func(t *testing.T) {
		t.Parallel()

		client, mock := redismock.NewClientMock()
		cm := dbrunnerservice.NewCacheModule(client)

		mock.MatchExpectationsInOrder(false)
		mock.ExpectGetEx("dbrunner:sql-input:"+mockInputHash, 1*time.Hour).SetVal(mockOutputHash)
		mock.ExpectGetEx("dbrunner:sql-output:"+mockOutputHash, 1*time.Hour).SetVal(string(lo.Must(json.Marshal(mockOutput))))
		mock.ExpectSetEx("dbrunner:sql-input:"+mockInputHash, mockOutputHash, 1*time.Hour).SetVal("OK")

		hash, err := cm.WriteToCache(context.TODO(), mockInput, mockOutput)
		require.NoError(t, err)
		assert.Equal(t, mockInputHash, hash)
	})

	t.Run("if there is no input hash but there is output hash, we write only input hash", func(t *testing.T) {
		t.Parallel()

		client, mock := redismock.NewClientMock()
		cm := dbrunnerservice.NewCacheModule(client)

		mock.MatchExpectationsInOrder(false)
		mock.ExpectGetEx("dbrunner:sql-input:"+mockInputHash, 1*time.Hour).SetErr(redis.Nil)
		mock.ExpectGetEx("dbrunner:sql-output:"+mockOutputHash, 1*time.Hour).SetVal(string(lo.Must(json.Marshal(mockOutput))))
		mock.ExpectSetEx("dbrunner:sql-input:"+mockInputHash, mockOutputHash, 1*time.Hour).SetVal("OK")

		hash, err := cm.WriteToCache(context.TODO(), mockInput, mockOutput)
		require.NoError(t, err)
		assert.Equal(t, mockInputHash, hash)
	})

	t.Run("the written cache should be retrievable", func(t *testing.T) {
		t.Parallel()

		client, mock := redismock.NewClientMock()
		cm := dbrunnerservice.NewCacheModule(client)

		mock.MatchExpectationsInOrder(true)
		mock.ExpectGetEx("dbrunner:sql-output:"+mockOutputHash, 1*time.Hour).SetErr(redis.Nil)
		mock.ExpectSetEx("dbrunner:sql-output:"+mockOutputHash, string(lo.Must(json.Marshal(mockOutput))), 1*time.Hour).SetVal("OK")
		mock.ExpectSetEx("dbrunner:sql-input:"+mockInputHash, mockOutputHash, 1*time.Hour).SetVal("OK")
		mock.ExpectGetEx("dbrunner:sql-input:"+mockInputHash, 1*time.Hour).SetVal(mockOutputHash)
		mock.ExpectGetEx("dbrunner:sql-output:"+mockOutputHash, 1*time.Hour).SetVal(string(lo.Must(json.Marshal(mockOutput))))

		hash, err := cm.WriteToCache(context.TODO(), mockInput, mockOutput)
		require.NoError(t, err)
		assert.Equal(t, mockInputHash, hash)

		outputHash, err := cm.GetOutputHash(context.TODO(), mockInputHash)
		require.NoError(t, err)
		assert.Equal(t, mockOutputHash, outputHash)

		output, err := cm.GetOutput(context.TODO(), outputHash)
		require.NoError(t, err)

		// convert to json to workaround the nil issue
		expectedJSON, _ := json.Marshal(mockOutput)
		actualJSON, _ := json.Marshal(output)
		assert.JSONEq(t, string(expectedJSON), string(actualJSON))
	})
}
