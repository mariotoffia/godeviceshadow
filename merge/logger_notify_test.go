package merge_test

import (
	"context"
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

func (m *MockMergeLoggerWithPreparePost) Managed(ctx context.Context, path string, operation model.MergeOperation, oldValue, newValue model.ValueAndTimestamp, oldTimeStamp, newTimeStamp time.Time) {
	m.Called(ctx, path, operation, oldValue, newValue, oldTimeStamp, newTimeStamp)
}

func (m *MockMergeLoggerWithPreparePost) Plain(ctx context.Context, path string, operation model.MergeOperation, oldValue, newValue any) {
	m.Called(ctx, path, operation, oldValue, newValue)
}

func (m *MockMergeLoggerWithPreparePost) Prepare(ctx context.Context) error {
	m.Called(ctx)
	return m.PrepareError
}

func (m *MockMergeLoggerWithPreparePost) Post(ctx context.Context, err error) error {
	m.Called(ctx, err)
	return m.PostError
}

func TestLoggerNotifyPrepare(t *testing.T) {
	// Test with no error
	mockLogger := &MockMergeLoggerWithPreparePost{}
	mockLogger.On("Prepare").Return(nil).Once()

	loggers := merge.MergeLoggers{mockLogger}
	err := loggers.NotifyPrepare(context.Background())
	assert.NoError(t, err)
	mockLogger.AssertExpectations(t)

	// Test with error
	mockLogger = &MockMergeLoggerWithPreparePost{
		PrepareError: errors.New("prepare error"),
	}
	mockLogger.On("Prepare").Return(mockLogger.PrepareError).Once()

	loggers = merge.MergeLoggers{mockLogger}
	err = loggers.NotifyPrepare(context.Background())
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
	err = loggers.NotifyPrepare(context.Background())
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
	err := loggers.NotifyPost(context.Background(), inputErr)
	assert.NoError(t, err)
	mockLogger.AssertExpectations(t)

	// Test with error
	mockLogger = &MockMergeLoggerWithPreparePost{
		PostError: errors.New("post error"),
	}
	mockLogger.On("Post", inputErr).Return(mockLogger.PostError).Once()

	loggers = merge.MergeLoggers{mockLogger}
	err = loggers.NotifyPost(context.Background(), inputErr)
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
	err = loggers.NotifyPost(context.Background(), inputErr)
	assert.Error(t, err)
	assert.Equal(t, "post error from logger 2", err.Error())
	mockLogger1.AssertExpectations(t)
	mockLogger2.AssertExpectations(t)
}
