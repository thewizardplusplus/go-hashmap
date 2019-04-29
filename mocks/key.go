// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Key is an autogenerated mock type for the Key type
type Key struct {
	mock.Mock
}

// Equals provides a mock function with given fields: key
func (_m *Key) Equals(key interface{}) bool {
	ret := _m.Called(key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Hash provides a mock function with given fields:
func (_m *Key) Hash() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}
