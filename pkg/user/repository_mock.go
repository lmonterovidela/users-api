package user

import (
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) Insert(user *User) (int64, error) {
	args := m.Called(user)
	err := args.Error(1)
	if args.Get(0) == nil {
		return 0, err
	}
	return args.Get(0).(int64), err
}

func (m *RepositoryMock) Update(userId int, user *User) error {
	args := m.Called(userId, user)
	return args.Error(0)
}

func (m *RepositoryMock) Get(userId int) (User, error) {
	args := m.Called(userId)
	err := args.Error(1)
	if args.Get(0) == nil {
		return User{}, err
	}
	return args.Get(0).(User), err
}

func (m *RepositoryMock) Find(name string, size int, offset int) (UserList, error) {
	args := m.Called(name, size, offset)
	err := args.Error(1)
	if args.Get(0) == nil {
		return UserList{}, err
	}
	return args.Get(0).(UserList), err
}

func (m *RepositoryMock) Delete(userId int) error {
	args := m.Called(userId)
	return args.Error(0)
}

func (m *RepositoryMock) GetLocation(userId int) (Location, error) {
	args := m.Called(userId)
	err := args.Error(1)
	if args.Get(0) == nil {
		return Location{}, err
	}
	return args.Get(0).(Location), err
}
