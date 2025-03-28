// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// SetupSubsRepo is an autogenerated mock type for the SetupSubsRepo type
type SetupSubsRepo struct {
	mock.Mock
}

type SetupSubsRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *SetupSubsRepo) EXPECT() *SetupSubsRepo_Expecter {
	return &SetupSubsRepo_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, userID, subType
func (_m *SetupSubsRepo) Delete(ctx context.Context, userID string, subType domain.SubscriptionType) error {
	ret := _m.Called(ctx, userID, subType)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.SubscriptionType) error); ok {
		r0 = rf(ctx, userID, subType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetupSubsRepo_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type SetupSubsRepo_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - subType domain.SubscriptionType
func (_e *SetupSubsRepo_Expecter) Delete(ctx interface{}, userID interface{}, subType interface{}) *SetupSubsRepo_Delete_Call {
	return &SetupSubsRepo_Delete_Call{Call: _e.mock.On("Delete", ctx, userID, subType)}
}

func (_c *SetupSubsRepo_Delete_Call) Run(run func(ctx context.Context, userID string, subType domain.SubscriptionType)) *SetupSubsRepo_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(domain.SubscriptionType))
	})
	return _c
}

func (_c *SetupSubsRepo_Delete_Call) Return(_a0 error) *SetupSubsRepo_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SetupSubsRepo_Delete_Call) RunAndReturn(run func(context.Context, string, domain.SubscriptionType) error) *SetupSubsRepo_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// IsExists provides a mock function with given fields: ctx, userID, subType
func (_m *SetupSubsRepo) IsExists(ctx context.Context, userID string, subType domain.SubscriptionType) (bool, error) {
	ret := _m.Called(ctx, userID, subType)

	if len(ret) == 0 {
		panic("no return value specified for IsExists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.SubscriptionType) (bool, error)); ok {
		return rf(ctx, userID, subType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.SubscriptionType) bool); ok {
		r0 = rf(ctx, userID, subType)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, domain.SubscriptionType) error); ok {
		r1 = rf(ctx, userID, subType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetupSubsRepo_IsExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsExists'
type SetupSubsRepo_IsExists_Call struct {
	*mock.Call
}

// IsExists is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - subType domain.SubscriptionType
func (_e *SetupSubsRepo_Expecter) IsExists(ctx interface{}, userID interface{}, subType interface{}) *SetupSubsRepo_IsExists_Call {
	return &SetupSubsRepo_IsExists_Call{Call: _e.mock.On("IsExists", ctx, userID, subType)}
}

func (_c *SetupSubsRepo_IsExists_Call) Run(run func(ctx context.Context, userID string, subType domain.SubscriptionType)) *SetupSubsRepo_IsExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(domain.SubscriptionType))
	})
	return _c
}

func (_c *SetupSubsRepo_IsExists_Call) Return(_a0 bool, _a1 error) *SetupSubsRepo_IsExists_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SetupSubsRepo_IsExists_Call) RunAndReturn(run func(context.Context, string, domain.SubscriptionType) (bool, error)) *SetupSubsRepo_IsExists_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: ctx, userID, subType
func (_m *SetupSubsRepo) Save(ctx context.Context, userID string, subType domain.SubscriptionType) error {
	ret := _m.Called(ctx, userID, subType)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.SubscriptionType) error); ok {
		r0 = rf(ctx, userID, subType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetupSubsRepo_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type SetupSubsRepo_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - subType domain.SubscriptionType
func (_e *SetupSubsRepo_Expecter) Save(ctx interface{}, userID interface{}, subType interface{}) *SetupSubsRepo_Save_Call {
	return &SetupSubsRepo_Save_Call{Call: _e.mock.On("Save", ctx, userID, subType)}
}

func (_c *SetupSubsRepo_Save_Call) Run(run func(ctx context.Context, userID string, subType domain.SubscriptionType)) *SetupSubsRepo_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(domain.SubscriptionType))
	})
	return _c
}

func (_c *SetupSubsRepo_Save_Call) Return(_a0 error) *SetupSubsRepo_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SetupSubsRepo_Save_Call) RunAndReturn(run func(context.Context, string, domain.SubscriptionType) error) *SetupSubsRepo_Save_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, userID, subType, enabled
func (_m *SetupSubsRepo) Update(ctx context.Context, userID string, subType domain.SubscriptionType, enabled bool) error {
	ret := _m.Called(ctx, userID, subType, enabled)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.SubscriptionType, bool) error); ok {
		r0 = rf(ctx, userID, subType, enabled)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetupSubsRepo_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type SetupSubsRepo_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - subType domain.SubscriptionType
//   - enabled bool
func (_e *SetupSubsRepo_Expecter) Update(ctx interface{}, userID interface{}, subType interface{}, enabled interface{}) *SetupSubsRepo_Update_Call {
	return &SetupSubsRepo_Update_Call{Call: _e.mock.On("Update", ctx, userID, subType, enabled)}
}

func (_c *SetupSubsRepo_Update_Call) Run(run func(ctx context.Context, userID string, subType domain.SubscriptionType, enabled bool)) *SetupSubsRepo_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(domain.SubscriptionType), args[3].(bool))
	})
	return _c
}

func (_c *SetupSubsRepo_Update_Call) Return(_a0 error) *SetupSubsRepo_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SetupSubsRepo_Update_Call) RunAndReturn(run func(context.Context, string, domain.SubscriptionType, bool) error) *SetupSubsRepo_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewSetupSubsRepo creates a new instance of SetupSubsRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSetupSubsRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *SetupSubsRepo {
	mock := &SetupSubsRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
