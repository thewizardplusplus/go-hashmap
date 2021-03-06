// Code generated by mockery v1.0.0. DO NOT EDIT.

package hashmap

import mock "github.com/stretchr/testify/mock"

// MockHandlerInterface is an autogenerated mock type for the HandlerInterface type
type MockHandlerInterface struct {
	mock.Mock
}

// Handle provides a mock function with given fields: key, value
func (_m *MockHandlerInterface) Handle(key Key, value interface{}) bool {
	ret := _m.Called(key, value)

	var r0 bool
	if rf, ok := ret.Get(0).(func(Key, interface{}) bool); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
