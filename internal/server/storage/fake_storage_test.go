package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewFakeStorage проверяем работу с FakeStorage
func TestNewFakeStorage(t *testing.T) {
	storage, err := NewFakeStorage()
	assert.NoError(t, err)
	defer storage.Close()
	assert.Nil(t, storage.Ping())
	assert.Nil(t, storage.Load())
	assert.Nil(t, storage.Save(context.Background()))
}
