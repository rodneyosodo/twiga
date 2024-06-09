// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import (
	context "context"

	notifications "github.com/rodneyosodo/twiga/notifications"
	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// CreateNotification provides a mock function with given fields: ctx, notification
func (_m *Service) CreateNotification(ctx context.Context, notification notifications.Notification) (notifications.Notification, error) {
	ret := _m.Called(ctx, notification)

	if len(ret) == 0 {
		panic("no return value specified for CreateNotification")
	}

	var r0 notifications.Notification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, notifications.Notification) (notifications.Notification, error)); ok {
		return rf(ctx, notification)
	}
	if rf, ok := ret.Get(0).(func(context.Context, notifications.Notification) notifications.Notification); ok {
		r0 = rf(ctx, notification)
	} else {
		r0 = ret.Get(0).(notifications.Notification)
	}

	if rf, ok := ret.Get(1).(func(context.Context, notifications.Notification) error); ok {
		r1 = rf(ctx, notification)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteNotification provides a mock function with given fields: ctx, token, id
func (_m *Service) DeleteNotification(ctx context.Context, token string, id string) error {
	ret := _m.Called(ctx, token, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteNotification")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, token, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadAllNotifications provides a mock function with given fields: ctx, token, page
func (_m *Service) ReadAllNotifications(ctx context.Context, token string, page notifications.Page) error {
	ret := _m.Called(ctx, token, page)

	if len(ret) == 0 {
		panic("no return value specified for ReadAllNotifications")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, notifications.Page) error); ok {
		r0 = rf(ctx, token, page)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadNotification provides a mock function with given fields: ctx, token, id
func (_m *Service) ReadNotification(ctx context.Context, token string, id string) error {
	ret := _m.Called(ctx, token, id)

	if len(ret) == 0 {
		panic("no return value specified for ReadNotification")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, token, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RetrieveAllNotifications provides a mock function with given fields: ctx, token, page
func (_m *Service) RetrieveAllNotifications(ctx context.Context, token string, page notifications.Page) (notifications.NotificationsPage, error) {
	ret := _m.Called(ctx, token, page)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveAllNotifications")
	}

	var r0 notifications.NotificationsPage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, notifications.Page) (notifications.NotificationsPage, error)); ok {
		return rf(ctx, token, page)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, notifications.Page) notifications.NotificationsPage); ok {
		r0 = rf(ctx, token, page)
	} else {
		r0 = ret.Get(0).(notifications.NotificationsPage)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, notifications.Page) error); ok {
		r1 = rf(ctx, token, page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveNotification provides a mock function with given fields: ctx, token, id
func (_m *Service) RetrieveNotification(ctx context.Context, token string, id string) (notifications.Notification, error) {
	ret := _m.Called(ctx, token, id)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveNotification")
	}

	var r0 notifications.Notification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (notifications.Notification, error)); ok {
		return rf(ctx, token, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) notifications.Notification); ok {
		r0 = rf(ctx, token, id)
	} else {
		r0 = ret.Get(0).(notifications.Notification)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, token, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
