// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import (
	context "context"

	users "github.com/rodneyosodo/twiga/users"
	mock "github.com/stretchr/testify/mock"
)

// FollowingRepository is an autogenerated mock type for the FollowingRepository type
type FollowingRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, following
func (_m *FollowingRepository) Create(ctx context.Context, following users.Following) (users.Following, error) {
	ret := _m.Called(ctx, following)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 users.Following
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, users.Following) (users.Following, error)); ok {
		return rf(ctx, following)
	}
	if rf, ok := ret.Get(0).(func(context.Context, users.Following) users.Following); ok {
		r0 = rf(ctx, following)
	} else {
		r0 = ret.Get(0).(users.Following)
	}

	if rf, ok := ret.Get(1).(func(context.Context, users.Following) error); ok {
		r1 = rf(ctx, following)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, following
func (_m *FollowingRepository) Delete(ctx context.Context, following users.Following) error {
	ret := _m.Called(ctx, following)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, users.Following) error); ok {
		r0 = rf(ctx, following)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RetrieveAll provides a mock function with given fields: ctx, page
func (_m *FollowingRepository) RetrieveAll(ctx context.Context, page users.Page) (users.FollowingsPage, error) {
	ret := _m.Called(ctx, page)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveAll")
	}

	var r0 users.FollowingsPage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, users.Page) (users.FollowingsPage, error)); ok {
		return rf(ctx, page)
	}
	if rf, ok := ret.Get(0).(func(context.Context, users.Page) users.FollowingsPage); ok {
		r0 = rf(ctx, page)
	} else {
		r0 = ret.Get(0).(users.FollowingsPage)
	}

	if rf, ok := ret.Get(1).(func(context.Context, users.Page) error); ok {
		r1 = rf(ctx, page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFollowingRepository creates a new instance of FollowingRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFollowingRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *FollowingRepository {
	mock := &FollowingRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}