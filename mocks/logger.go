package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/cryptoPickle/go-ddd-example/internal/infrastructure/logger"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(...any)                   {}
func (m *MockLogger) Debugf(format string, v ...any) {}
func (m *MockLogger) Info(...any)                    {}
func (m *MockLogger) Infof(format string, v ...any)  {}
func (m *MockLogger) Warn(...any)                    {}
func (m *MockLogger) Warnf(format string, v ...any)  {}
func (m *MockLogger) Error(...any)                   {}
func (m *MockLogger) Errorf(format string, v ...any) {}
func (m *MockLogger) Fatal(...any)                   {}
func (m *MockLogger) Fatalf(format string, v ...any) {}
func (m *MockLogger) SetLevel(lvl logger.LogLevel) error {
	return nil
}
