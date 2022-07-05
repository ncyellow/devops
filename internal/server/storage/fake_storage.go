package storage

import (
	"context"
)

// FakeStorage Пустая реализация хранилища - если нет ни файла ни базы,
// то используем эту реализацию которая ничего не делает
type FakeStorage struct {
}

func (m *FakeStorage) Ping() error {
	return nil
}

func (m *FakeStorage) Close() {
}

func (m *FakeStorage) Load() error {
	return nil
}

func (m *FakeStorage) Save(context.Context) error {
	return nil
}

func NewFakeStorage() (*FakeStorage, error) {
	return &FakeStorage{}, nil
}
