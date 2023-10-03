// Code generated by mockery v2.27.1. DO NOT EDIT.

package storage

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// MockArtifactStorageProvider is an autogenerated mock type for the ArtifactStorageProvider type
type MockArtifactStorageProvider struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, artifactURI, path
func (_m *MockArtifactStorageProvider) Get(ctx context.Context, artifactURI string, path string) (io.ReadCloser, error) {
	ret := _m.Called(ctx, artifactURI, path)

	var r0 io.ReadCloser
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (io.ReadCloser, error)); ok {
		return rf(ctx, artifactURI, path)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) io.ReadCloser); ok {
		r0 = rf(ctx, artifactURI, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, artifactURI, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, artifactURI, path
func (_m *MockArtifactStorageProvider) List(ctx context.Context, artifactURI string, path string) ([]ArtifactObject, error) {
	ret := _m.Called(ctx, artifactURI, path)

	var r0 []ArtifactObject
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]ArtifactObject, error)); ok {
		return rf(ctx, artifactURI, path)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []ArtifactObject); ok {
		r0 = rf(ctx, artifactURI, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ArtifactObject)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, artifactURI, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockArtifactStorageProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockArtifactStorageProvider creates a new instance of MockArtifactStorageProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockArtifactStorageProvider(t mockConstructorTestingTNewMockArtifactStorageProvider) *MockArtifactStorageProvider {
	mock := &MockArtifactStorageProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
