package storage

type FakeStorage struct {
}

func (m *FakeStorage) Close() {
}

func (m *FakeStorage) Load() error {
	return nil
}

func (m *FakeStorage) Save() error {
	return nil
}

func NewFakeStorage() (PersistentStorage, error) {
	return &FakeStorage{}, nil
}
