// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	store "github.com/mr-linch/go-tg-bot/internal/store"
	mock "github.com/stretchr/testify/mock"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// Migrator provides a mock function with given fields:
func (_m *Store) Migrator() store.Migrator {
	ret := _m.Called()

	var r0 store.Migrator
	if rf, ok := ret.Get(0).(func() store.Migrator); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.Migrator)
		}
	}

	return r0
}

// Tx provides a mock function with given fields: ctx, txFunc
func (_m *Store) Tx(ctx context.Context, txFunc store.TxFunc) error {
	ret := _m.Called(ctx, txFunc)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, store.TxFunc) error); ok {
		r0 = rf(ctx, txFunc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// User provides a mock function with given fields:
func (_m *Store) User() store.User {
	ret := _m.Called()

	var r0 store.User
	if rf, ok := ret.Get(0).(func() store.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.User)
		}
	}

	return r0
}

type mockConstructorTestingTNewStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStore(t mockConstructorTestingTNewStore) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
