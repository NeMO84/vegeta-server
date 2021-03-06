// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import io "io"
import mock "github.com/stretchr/testify/mock"
import models "vegeta-server/models"
import time "time"
import vegeta "vegeta-server/pkg/vegeta"

// ITask is an autogenerated mock type for the ITask type
type ITask struct {
	mock.Mock
}

// Cancel provides a mock function with given fields:
func (_m *ITask) Cancel() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Complete provides a mock function with given fields: _a0
func (_m *ITask) Complete(_a0 io.Reader) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Reader) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatedAt provides a mock function with given fields:
func (_m *ITask) CreatedAt() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// Fail provides a mock function with given fields:
func (_m *ITask) Fail() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ID provides a mock function with given fields:
func (_m *ITask) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Params provides a mock function with given fields:
func (_m *ITask) Params() models.AttackParams {
	ret := _m.Called()

	var r0 models.AttackParams
	if rf, ok := ret.Get(0).(func() models.AttackParams); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(models.AttackParams)
	}

	return r0
}

// Run provides a mock function with given fields: _a0
func (_m *ITask) Run(_a0 vegeta.AttackFunc) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(vegeta.AttackFunc) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Status provides a mock function with given fields:
func (_m *ITask) Status() models.AttackStatus {
	ret := _m.Called()

	var r0 models.AttackStatus
	if rf, ok := ret.Get(0).(func() models.AttackStatus); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(models.AttackStatus)
	}

	return r0
}

// UpdatedAt provides a mock function with given fields:
func (_m *ITask) UpdatedAt() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}
