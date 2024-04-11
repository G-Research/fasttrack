// Code generated by mockery v2.34.0. DO NOT EDIT.

package repositories

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	models "github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
)

// MockNamespaceRepositoryProvider is an autogenerated mock type for the NamespaceRepositoryProvider type
type MockNamespaceRepositoryProvider struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, namespace
func (_m *MockNamespaceRepositoryProvider) Create(ctx context.Context, namespace *models.Namespace) error {
	ret := _m.Called(ctx, namespace)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Namespace) error); ok {
		r0 = rf(ctx, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, namespace
func (_m *MockNamespaceRepositoryProvider) Delete(ctx context.Context, namespace *models.Namespace) error {
	ret := _m.Called(ctx, namespace)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Namespace) error); ok {
		r0 = rf(ctx, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByCode provides a mock function with given fields: ctx, code
func (_m *MockNamespaceRepositoryProvider) GetByCode(ctx context.Context, code string) (*models.Namespace, error) {
	ret := _m.Called(ctx, code)

	var r0 *models.Namespace
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.Namespace, error)); ok {
		return rf(ctx, code)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Namespace); ok {
		r0 = rf(ctx, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Namespace)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *MockNamespaceRepositoryProvider) GetByID(ctx context.Context, id uint) (*models.Namespace, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.Namespace
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) (*models.Namespace, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) *models.Namespace); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Namespace)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDB provides a mock function with given fields:
func (_m *MockNamespaceRepositoryProvider) GetDB() *gorm.DB {
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

// List provides a mock function with given fields: ctx
func (_m *MockNamespaceRepositoryProvider) List(ctx context.Context) ([]models.Namespace, error) {
	ret := _m.Called(ctx)

	var r0 []models.Namespace
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.Namespace, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.Namespace); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Namespace)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, namespace
func (_m *MockNamespaceRepositoryProvider) Update(ctx context.Context, namespace *models.Namespace) error {
	ret := _m.Called(ctx, namespace)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Namespace) error); ok {
		r0 = rf(ctx, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockNamespaceRepositoryProvider creates a new instance of MockNamespaceRepositoryProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockNamespaceRepositoryProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockNamespaceRepositoryProvider {
	mock := &MockNamespaceRepositoryProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
