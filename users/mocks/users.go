// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import (
	context "context"

	users "github.com/rodneyosodo/twiga/users"
	mock "github.com/stretchr/testify/mock"
)

// UsersRepository is an autogenerated mock type for the UsersRepository type
type UsersRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, user
func (_m *UsersRepository) Create(ctx context.Context, user users.User) (users.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 users.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, users.User) (users.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, users.User) users.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, users.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *UsersRepository) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RetrieveAll provides a mock function with given fields: ctx, page
func (_m *UsersRepository) RetrieveAll(ctx context.Context, page users.Page) (users.UsersPage, error) {
	ret := _m.Called(ctx, page)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveAll")
	}

	var r0 users.UsersPage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, users.Page) (users.UsersPage, error)); ok {
		return rf(ctx, page)
	}
	if rf, ok := ret.Get(0).(func(context.Context, users.Page) users.UsersPage); ok {
		r0 = rf(ctx, page)
	} else {
		r0 = ret.Get(0).(users.UsersPage)
	}

	if rf, ok := ret.Get(1).(func(context.Context, users.Page) error); ok {
		r1 = rf(ctx, page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveByEmail provides a mock function with given fields: ctx, email
func (_m *UsersRepository) RetrieveByEmail(ctx context.Context, email string) (users.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveByEmail")
	}

	var r0 users.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (users.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) users.User); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveByID provides a mock function with given fields: ctx, id
func (_m *UsersRepository) RetrieveByID(ctx context.Context, id string) (users.User, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveByID")
	}

	var r0 users.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (users.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) users.User); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, user
func (_m *UsersRepository) Update(ctx context.Context, user users.User) (users.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 users.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, users.User) (users.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, users.User) users.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, users.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateBio provides a mock function with given fields: ctx, user
func (_m *UsersRepository) UpdateBio(ctx context.Context, user users.User) (users.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdateBio")
	}

	var r0 users.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, users.User) (users.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, users.User) users.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, users.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEmail provides a mock function with given fields: ctx, user
func (_m *UsersRepository) UpdateEmail(ctx context.Context, user users.User) (users.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdateEmail")
	}

	var r0 users.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, users.User) (users.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, users.User) users.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, users.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePassword provides a mock function with given fields: ctx, user
func (_m *UsersRepository) UpdatePassword(ctx context.Context, user users.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePassword")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, users.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePictureURL provides a mock function with given fields: ctx, user
func (_m *UsersRepository) UpdatePictureURL(ctx context.Context, user users.User) (users.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePictureURL")
	}

	var r0 users.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, users.User) (users.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, users.User) users.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, users.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePreferences provides a mock function with given fields: ctx, user
func (_m *UsersRepository) UpdatePreferences(ctx context.Context, user users.User) (users.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePreferences")
	}

	var r0 users.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, users.User) (users.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, users.User) users.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, users.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUsername provides a mock function with given fields: ctx, user
func (_m *UsersRepository) UpdateUsername(ctx context.Context, user users.User) (users.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUsername")
	}

	var r0 users.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, users.User) (users.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, users.User) users.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, users.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUsersRepository creates a new instance of UsersRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUsersRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UsersRepository {
	mock := &UsersRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
