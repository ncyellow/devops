package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewFakeStorage проверяем работу с FakeStorage
func TestNewFakeStorage(t *testing.T) {
	storage, err := NewFakeStorage()
	defer storage.Close()
	assert.NoError(t, err)
	assert.Nil(t, storage.Ping())
	assert.Nil(t, storage.Load())
	assert.Nil(t, storage.Save(context.Background()))
}
