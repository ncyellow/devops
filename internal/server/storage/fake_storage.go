package storage

import (
	"context"
	"errors"
)

// FakeStorage Пустая реализация хранилища - если нет ни файла ни базы,
// то используем эту реализацию которая ничего не делает
type FakeStorage struct {
}

func (m *FakeStorage) Ping() error {
	return errors.New("not supported operation")
}

func (m *FakeStorage) Close() {
}

func (m *FakeStorage) Load() error {
	return nil
}

func (m *FakeStorage) Save(context.Context) error {
	return nil
}

func NewFakeStorage() (PersistentStorage, error) {
	return &FakeStorage{}, nil
}
