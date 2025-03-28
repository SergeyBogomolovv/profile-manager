// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// SetupTokenRepo is an autogenerated mock type for the SetupTokenRepo type
type SetupTokenRepo struct {
	mock.Mock
}

type SetupTokenRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *SetupTokenRepo) EXPECT() *SetupTokenRepo_Expecter {
	return &SetupTokenRepo_Expecter{mock: &_m.Mock}
}

// CheckUserID provides a mock function with given fields: ctx, token
func (_m *SetupTokenRepo) CheckUserID(ctx context.Context, token string) (string, error) {
	ret := _m.Called(ctx, token)

	if len(ret) == 0 {
		panic("no return value specified for CheckUserID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetupTokenRepo_CheckUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckUserID'
type SetupTokenRepo_CheckUserID_Call struct {
	*mock.Call
}

// CheckUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - token string
func (_e *SetupTokenRepo_Expecter) CheckUserID(ctx interface{}, token interface{}) *SetupTokenRepo_CheckUserID_Call {
	return &SetupTokenRepo_CheckUserID_Call{Call: _e.mock.On("CheckUserID", ctx, token)}
}

func (_c *SetupTokenRepo_CheckUserID_Call) Run(run func(ctx context.Context, token string)) *SetupTokenRepo_CheckUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *SetupTokenRepo_CheckUserID_Call) Return(_a0 string, _a1 error) *SetupTokenRepo_CheckUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SetupTokenRepo_CheckUserID_Call) RunAndReturn(run func(context.Context, string) (string, error)) *SetupTokenRepo_CheckUserID_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, userID
func (_m *SetupTokenRepo) Create(ctx context.Context, userID string) (string, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetupTokenRepo_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type SetupTokenRepo_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
func (_e *SetupTokenRepo_Expecter) Create(ctx interface{}, userID interface{}) *SetupTokenRepo_Create_Call {
	return &SetupTokenRepo_Create_Call{Call: _e.mock.On("Create", ctx, userID)}
}

func (_c *SetupTokenRepo_Create_Call) Run(run func(ctx context.Context, userID string)) *SetupTokenRepo_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *SetupTokenRepo_Create_Call) Return(_a0 string, _a1 error) *SetupTokenRepo_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SetupTokenRepo_Create_Call) RunAndReturn(run func(context.Context, string) (string, error)) *SetupTokenRepo_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Revoke provides a mock function with given fields: ctx, token
func (_m *SetupTokenRepo) Revoke(ctx context.Context, token string) error {
	ret := _m.Called(ctx, token)

	if len(ret) == 0 {
		panic("no return value specified for Revoke")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetupTokenRepo_Revoke_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Revoke'
type SetupTokenRepo_Revoke_Call struct {
	*mock.Call
}

// Revoke is a helper method to define mock.On call
//   - ctx context.Context
//   - token string
func (_e *SetupTokenRepo_Expecter) Revoke(ctx interface{}, token interface{}) *SetupTokenRepo_Revoke_Call {
	return &SetupTokenRepo_Revoke_Call{Call: _e.mock.On("Revoke", ctx, token)}
}

func (_c *SetupTokenRepo_Revoke_Call) Run(run func(ctx context.Context, token string)) *SetupTokenRepo_Revoke_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *SetupTokenRepo_Revoke_Call) Return(_a0 error) *SetupTokenRepo_Revoke_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SetupTokenRepo_Revoke_Call) RunAndReturn(run func(context.Context, string) error) *SetupTokenRepo_Revoke_Call {
	_c.Call.Return(run)
	return _c
}

// NewSetupTokenRepo creates a new instance of SetupTokenRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSetupTokenRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *SetupTokenRepo {
	mock := &SetupTokenRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
