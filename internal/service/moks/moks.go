package moks

import "github.com/stretchr/testify/mock"

// MockStorage представляет заглушку для хранилища данных.
type MockStorage struct {
	mock.Mock
}

// Get имитирует метод Get интерфейса app.Storage.
func (m *MockStorage) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

// Put имитирует метод Put интерфейса app.Storage.
func (m *MockStorage) Put(key, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

// Load имитирует метод Load интерфейса app.Storage.
func (m *MockStorage) Load(data string) error {
	args := m.Called(data)
	return args.Error(0)
}
