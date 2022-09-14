package storage

import (
	"context"
)

// FakeStorage Пустая реализация хранилища - если нет ни файла ни базы, реализует интерфейс PersistentStorage
// то используем эту реализацию которая ничего не делает
type FakeStorage struct {
}

// NewFakeStorage конструктор пустого хранилища, явно не используется, только через фабрику
func NewFakeStorage() (*FakeStorage, error) {
	return &FakeStorage{}, nil
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
