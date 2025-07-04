package merge_test

import (
	"errors"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMergeLoggerWithPreparePost implements both prepare and post interfaces
type MockMergeLoggerWithPreparePost struct {
	mock.Mock
	PrepareError error
	PostError    error
}

func (m *MockMergeLoggerWithPreparePost) Managed(path string, operation model.MergeOperation, oldValue, newValue model.ValueAndTimestamp, oldTimeStamp, newTimeStamp time.Time) {
	m.Called(path, operation, oldValue, newValue, oldTimeStamp, newTimeStamp)
}

func (m *MockMergeLoggerWithPreparePost) Plain(path string, operation model.MergeOperation, oldValue, newValue any) {
	m.Called(path, operation, oldValue, newValue)
}

func (m *MockMergeLoggerWithPreparePost) Prepare() error {
	m.Called()
	return m.PrepareError
}

func (m *MockMergeLoggerWithPreparePost) Post(err error) error {
	m.Called(err)
	return m.PostError
}

func TestLoggerNotifyPrepare(t *testing.T) {
	// Test with no error
	mockLogger := &MockMergeLoggerWithPreparePost{}
	mockLogger.On("Prepare").Return(nil).Once()

	loggers := merge.MergeLoggers{mockLogger}
	err := loggers.NotifyPrepare()
	assert.NoError(t, err)
	mockLogger.AssertExpectations(t)

	// Test with error
	mockLogger = &MockMergeLoggerWithPreparePost{
		PrepareError: errors.New("prepare error"),
	}
	mockLogger.On("Prepare").Return(mockLogger.PrepareError).Once()

	loggers = merge.MergeLoggers{mockLogger}
	err = loggers.NotifyPrepare()
	assert.Error(t, err)
	assert.Equal(t, "prepare error", err.Error())
	mockLogger.AssertExpectations(t)

	// Test with multiple loggers, one with error
	mockLogger1 := &MockMergeLoggerWithPreparePost{}
	mockLogger1.On("Prepare").Return(nil).Once()

	mockLogger2 := &MockMergeLoggerWithPreparePost{
		PrepareError: errors.New("prepare error from logger 2"),
	}
	mockLogger2.On("Prepare").Return(mockLogger2.PrepareError).Once()

	loggers = merge.MergeLoggers{mockLogger1, mockLogger2}
	err = loggers.NotifyPrepare()
	assert.Error(t, err)
	assert.Equal(t, "prepare error from logger 2", err.Error())
	mockLogger1.AssertExpectations(t)
	mockLogger2.AssertExpectations(t)
}

func TestLoggerNotifyPost(t *testing.T) {
	inputErr := errors.New("input error")

	// Test with no error
	mockLogger := &MockMergeLoggerWithPreparePost{}
	mockLogger.On("Post", inputErr).Return(nil).Once()

	loggers := merge.MergeLoggers{mockLogger}
	err := loggers.NotifyPost(inputErr)
	assert.NoError(t, err)
	mockLogger.AssertExpectations(t)

	// Test with error
	mockLogger = &MockMergeLoggerWithPreparePost{
		PostError: errors.New("post error"),
	}
	mockLogger.On("Post", inputErr).Return(mockLogger.PostError).Once()

	loggers = merge.MergeLoggers{mockLogger}
	err = loggers.NotifyPost(inputErr)
	assert.Error(t, err)
	assert.Equal(t, "post error", err.Error())
	mockLogger.AssertExpectations(t)

	// Test with multiple loggers, one with error
	mockLogger1 := &MockMergeLoggerWithPreparePost{}
	mockLogger1.On("Post", inputErr).Return(nil).Once()

	mockLogger2 := &MockMergeLoggerWithPreparePost{
		PostError: errors.New("post error from logger 2"),
	}
	mockLogger2.On("Post", inputErr).Return(mockLogger2.PostError).Once()

	loggers = merge.MergeLoggers{mockLogger1, mockLogger2}
	err = loggers.NotifyPost(inputErr)
	assert.Error(t, err)
	assert.Equal(t, "post error from logger 2", err.Error())
	mockLogger1.AssertExpectations(t)
	mockLogger2.AssertExpectations(t)
}
