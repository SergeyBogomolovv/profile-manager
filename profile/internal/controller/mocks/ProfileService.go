// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/SergeyBogomolovv/profile-manager/profile/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// ProfileService is an autogenerated mock type for the ProfileService type
type ProfileService struct {
	mock.Mock
}

type ProfileService_Expecter struct {
	mock *mock.Mock
}

func (_m *ProfileService) EXPECT() *ProfileService_Expecter {
	return &ProfileService_Expecter{mock: &_m.Mock}
}

// GetProfile provides a mock function with given fields: ctx, userID
func (_m *ProfileService) GetProfile(ctx context.Context, userID string) (domain.Profile, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetProfile")
	}

	var r0 domain.Profile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.Profile, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Profile); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(domain.Profile)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileService_GetProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProfile'
type ProfileService_GetProfile_Call struct {
	*mock.Call
}

// GetProfile is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
func (_e *ProfileService_Expecter) GetProfile(ctx interface{}, userID interface{}) *ProfileService_GetProfile_Call {
	return &ProfileService_GetProfile_Call{Call: _e.mock.On("GetProfile", ctx, userID)}
}

func (_c *ProfileService_GetProfile_Call) Run(run func(ctx context.Context, userID string)) *ProfileService_GetProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ProfileService_GetProfile_Call) Return(_a0 domain.Profile, _a1 error) *ProfileService_GetProfile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileService_GetProfile_Call) RunAndReturn(run func(context.Context, string) (domain.Profile, error)) *ProfileService_GetProfile_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, userID, dto
func (_m *ProfileService) Update(ctx context.Context, userID string, dto domain.UpdateProfileDTO) (domain.Profile, error) {
	ret := _m.Called(ctx, userID, dto)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 domain.Profile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.UpdateProfileDTO) (domain.Profile, error)); ok {
		return rf(ctx, userID, dto)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.UpdateProfileDTO) domain.Profile); ok {
		r0 = rf(ctx, userID, dto)
	} else {
		r0 = ret.Get(0).(domain.Profile)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, domain.UpdateProfileDTO) error); ok {
		r1 = rf(ctx, userID, dto)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileService_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type ProfileService_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - dto domain.UpdateProfileDTO
func (_e *ProfileService_Expecter) Update(ctx interface{}, userID interface{}, dto interface{}) *ProfileService_Update_Call {
	return &ProfileService_Update_Call{Call: _e.mock.On("Update", ctx, userID, dto)}
}

func (_c *ProfileService_Update_Call) Run(run func(ctx context.Context, userID string, dto domain.UpdateProfileDTO)) *ProfileService_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(domain.UpdateProfileDTO))
	})
	return _c
}

func (_c *ProfileService_Update_Call) Return(_a0 domain.Profile, _a1 error) *ProfileService_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileService_Update_Call) RunAndReturn(run func(context.Context, string, domain.UpdateProfileDTO) (domain.Profile, error)) *ProfileService_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewProfileService creates a new instance of ProfileService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProfileService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProfileService {
	mock := &ProfileService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
