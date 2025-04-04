// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	events "github.com/SergeyBogomolovv/profile-manager/common/api/events"
	mock "github.com/stretchr/testify/mock"
)

// Broker is an autogenerated mock type for the Broker type
type Broker struct {
	mock.Mock
}

type Broker_Expecter struct {
	mock *mock.Mock
}

func (_m *Broker) EXPECT() *Broker_Expecter {
	return &Broker_Expecter{mock: &_m.Mock}
}

// PublishUserLogin provides a mock function with given fields: user
func (_m *Broker) PublishUserLogin(user events.UserLogin) error {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for PublishUserLogin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(events.UserLogin) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Broker_PublishUserLogin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishUserLogin'
type Broker_PublishUserLogin_Call struct {
	*mock.Call
}

// PublishUserLogin is a helper method to define mock.On call
//   - user events.UserLogin
func (_e *Broker_Expecter) PublishUserLogin(user interface{}) *Broker_PublishUserLogin_Call {
	return &Broker_PublishUserLogin_Call{Call: _e.mock.On("PublishUserLogin", user)}
}

func (_c *Broker_PublishUserLogin_Call) Run(run func(user events.UserLogin)) *Broker_PublishUserLogin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(events.UserLogin))
	})
	return _c
}

func (_c *Broker_PublishUserLogin_Call) Return(_a0 error) *Broker_PublishUserLogin_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Broker_PublishUserLogin_Call) RunAndReturn(run func(events.UserLogin) error) *Broker_PublishUserLogin_Call {
	_c.Call.Return(run)
	return _c
}

// PublishUserRegister provides a mock function with given fields: user
func (_m *Broker) PublishUserRegister(user events.UserRegister) error {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for PublishUserRegister")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(events.UserRegister) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Broker_PublishUserRegister_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishUserRegister'
type Broker_PublishUserRegister_Call struct {
	*mock.Call
}

// PublishUserRegister is a helper method to define mock.On call
//   - user events.UserRegister
func (_e *Broker_Expecter) PublishUserRegister(user interface{}) *Broker_PublishUserRegister_Call {
	return &Broker_PublishUserRegister_Call{Call: _e.mock.On("PublishUserRegister", user)}
}

func (_c *Broker_PublishUserRegister_Call) Run(run func(user events.UserRegister)) *Broker_PublishUserRegister_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(events.UserRegister))
	})
	return _c
}

func (_c *Broker_PublishUserRegister_Call) Return(_a0 error) *Broker_PublishUserRegister_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Broker_PublishUserRegister_Call) RunAndReturn(run func(events.UserRegister) error) *Broker_PublishUserRegister_Call {
	_c.Call.Return(run)
	return _c
}

// NewBroker creates a new instance of Broker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBroker(t interface {
	mock.TestingT
	Cleanup(func())
}) *Broker {
	mock := &Broker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
