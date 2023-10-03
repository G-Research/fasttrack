// Code generated by mockery v2.27.1. DO NOT EDIT.

package repositories

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	models "github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
)

// MockRunRepositoryProvider is an autogenerated mock type for the RunRepositoryProvider type
type MockRunRepositoryProvider struct {
	mock.Mock
}

// Archive provides a mock function with given fields: ctx, run
func (_m *MockRunRepositoryProvider) Archive(ctx context.Context, run *models.Run) error {
	ret := _m.Called(ctx, run)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Run) error); ok {
		r0 = rf(ctx, run)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ArchiveBatch provides a mock function with given fields: ctx, ids
func (_m *MockRunRepositoryProvider) ArchiveBatch(ctx context.Context, ids []string) error {
	ret := _m.Called(ctx, ids)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) error); ok {
		r0 = rf(ctx, ids)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: ctx, run
func (_m *MockRunRepositoryProvider) Create(ctx context.Context, run *models.Run) error {
	ret := _m.Called(ctx, run)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Run) error); ok {
		r0 = rf(ctx, run)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, run
func (_m *MockRunRepositoryProvider) Delete(ctx context.Context, run *models.Run) error {
	ret := _m.Called(ctx, run)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Run) error); ok {
		r0 = rf(ctx, run)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteBatch provides a mock function with given fields: ctx, ids
func (_m *MockRunRepositoryProvider) DeleteBatch(ctx context.Context, ids []string) error {
	ret := _m.Called(ctx, ids)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) error); ok {
		r0 = rf(ctx, ids)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *MockRunRepositoryProvider) GetByID(ctx context.Context, id string) (*models.Run, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.Run
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.Run, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Run); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Run)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByNamespaceIDAndRunID provides a mock function with given fields: ctx, namespaceID, runID
func (_m *MockRunRepositoryProvider) GetByNamespaceIDAndRunID(ctx context.Context, namespaceID uint, runID string) (*models.Run, error) {
	ret := _m.Called(ctx, namespaceID, runID)

	var r0 *models.Run
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) (*models.Run, error)); ok {
		return rf(ctx, namespaceID, runID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) *models.Run); ok {
		r0 = rf(ctx, namespaceID, runID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Run)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, string) error); ok {
		r1 = rf(ctx, namespaceID, runID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByNamespaceIDRunIDAndLifecycleStage provides a mock function with given fields: ctx, namespaceID, runID, lifecycleStage
func (_m *MockRunRepositoryProvider) GetByNamespaceIDRunIDAndLifecycleStage(ctx context.Context, namespaceID uint, runID string, lifecycleStage models.LifecycleStage) (*models.Run, error) {
	ret := _m.Called(ctx, namespaceID, runID, lifecycleStage)

	var r0 *models.Run
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, string, models.LifecycleStage) (*models.Run, error)); ok {
		return rf(ctx, namespaceID, runID, lifecycleStage)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, string, models.LifecycleStage) *models.Run); ok {
		r0 = rf(ctx, namespaceID, runID, lifecycleStage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Run)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, string, models.LifecycleStage) error); ok {
		r1 = rf(ctx, namespaceID, runID, lifecycleStage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDB provides a mock function with given fields:
func (_m *MockRunRepositoryProvider) GetDB() *gorm.DB {
	ret := _m.Called()

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func() *gorm.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// Restore provides a mock function with given fields: ctx, run
func (_m *MockRunRepositoryProvider) Restore(ctx context.Context, run *models.Run) error {
	ret := _m.Called(ctx, run)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Run) error); ok {
		r0 = rf(ctx, run)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RestoreBatch provides a mock function with given fields: ctx, ids
func (_m *MockRunRepositoryProvider) RestoreBatch(ctx context.Context, ids []string) error {
	ret := _m.Called(ctx, ids)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) error); ok {
		r0 = rf(ctx, ids)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetRunTagsBatch provides a mock function with given fields: ctx, run, batchSize, tags
func (_m *MockRunRepositoryProvider) SetRunTagsBatch(ctx context.Context, run *models.Run, batchSize int, tags []models.Tag) error {
	ret := _m.Called(ctx, run, batchSize, tags)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Run, int, []models.Tag) error); ok {
		r0 = rf(ctx, run, batchSize, tags)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, run
func (_m *MockRunRepositoryProvider) Update(ctx context.Context, run *models.Run) error {
	ret := _m.Called(ctx, run)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Run) error); ok {
		r0 = rf(ctx, run)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateWithTransaction provides a mock function with given fields: ctx, tx, run
func (_m *MockRunRepositoryProvider) UpdateWithTransaction(ctx context.Context, tx *gorm.DB, run *models.Run) error {
	ret := _m.Called(ctx, tx, run)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *models.Run) error); ok {
		r0 = rf(ctx, tx, run)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockRunRepositoryProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockRunRepositoryProvider creates a new instance of MockRunRepositoryProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockRunRepositoryProvider(t mockConstructorTestingTNewMockRunRepositoryProvider) *MockRunRepositoryProvider {
	mock := &MockRunRepositoryProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
