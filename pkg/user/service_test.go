package user

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestService_Create(t *testing.T) {

	repositoryMock := &RepositoryMock{}
	user := &User{
		Address: "address",
		Dob:     time.Now(),
		Name:    "Jhon",
	}

	type args struct {
		user *User
	}

	tests := []struct {
		name        string
		initMocks   func()
		args        args
		assertMocks func(*testing.T)
		assertError func(*testing.T, error)
		assertFunc  func(*testing.T, int64)
	}{
		{
			name: "Success - repository response ok",
			initMocks: func() {
				repositoryMock.On("Insert", user).
					Return(int64(1), nil).Once()
			},
			args: args{
				user: user,
			},
			assertMocks: func(t *testing.T) {
				repositoryMock.AssertExpectations(t)
			},
			assertError: func(t *testing.T, e error) {
				assert.Nil(t, e)
			},
			assertFunc: func(t *testing.T, userId int64) {
				assert.Equal(t, userId, int64(1))
			},
		},
		{
			name: "Error - repository response err",
			initMocks: func() {
				repositoryMock.On("Insert", user).
					Return(int64(0), errors.New("some error")).Once()
			},
			args: args{
				user: user,
			},
			assertMocks: func(t *testing.T) {
				repositoryMock.AssertExpectations(t)
			},
			assertError: func(t *testing.T, e error) {
				assert.NotNil(t, e)
			},
			assertFunc: func(t *testing.T, userId int64) {
				assert.Equal(t, userId, int64(0))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initMocks()
			service := Service{
				repository: repositoryMock,
			}

			userId, err := service.Create(tt.args.user)
			tt.assertMocks(t)
			tt.assertError(t, err)
			tt.assertFunc(t, userId)
		})
	}
}

func TestService_Get(t *testing.T) {

	repositoryMock := &RepositoryMock{}

	type args struct {
		userId int
	}

	tests := []struct {
		name        string
		initMocks   func()
		args        args
		assertMocks func(*testing.T)
		assertError func(*testing.T, error)
		assertFunc  func(*testing.T, User)
	}{
		{
			name: "Success - repository response ok",
			initMocks: func() {
				repositoryMock.On("Get", 1).
					Return(User{
						Id:      1,
						Address: "address",
						Name:    "Jhon",
					}, nil).Once()
			},
			args: args{
				userId: 1,
			},
			assertMocks: func(t *testing.T) {
				repositoryMock.AssertExpectations(t)
			},
			assertError: func(t *testing.T, e error) {
				assert.Nil(t, e)
			},
			assertFunc: func(t *testing.T, user User) {
				assert.Equal(t, user, User{
					Id:      1,
					Address: "address",
					Name:    "Jhon",
				})
			},
		},
		{
			name: "Error - repository response err",
			initMocks: func() {
				repositoryMock.On("Get", 1).
					Return(User{}, errors.New("some error")).Once()
			},
			args: args{
				userId: 1,
			},
			assertMocks: func(t *testing.T) {
				repositoryMock.AssertExpectations(t)
			},
			assertError: func(t *testing.T, e error) {
				assert.NotNil(t, e)
			},
			assertFunc: func(t *testing.T, user User) {
				assert.Equal(t, user.Id, 0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initMocks()
			service := Service{
				repository: repositoryMock,
			}

			userId, err := service.Get(tt.args.userId)
			tt.assertMocks(t)
			tt.assertError(t, err)
			tt.assertFunc(t, userId)
		})
	}
}

func TestService_Find(t *testing.T) {

	repositoryMock := &RepositoryMock{}

	type args struct {
		name   string
		offset int
		size   int
	}

	tests := []struct {
		name        string
		initMocks   func()
		args        args
		assertMocks func(*testing.T)
		assertError func(*testing.T, error)
		assertFunc  func(*testing.T, UserList)
	}{
		{
			name: "Success - repository response ok",
			initMocks: func() {
				repositoryMock.On("Find", "Jhon", 1, 2).
					Return(UserList{Data: make([]User, 5)}, nil).Once()
			},
			args: args{
				name:   "Jhon",
				offset: 2,
				size:   1,
			},
			assertMocks: func(t *testing.T) {
				repositoryMock.AssertExpectations(t)
			},
			assertError: func(t *testing.T, e error) {
				assert.Nil(t, e)
			},
			assertFunc: func(t *testing.T, users UserList) {
				assert.Equal(t, len(users.Data), 5)
			},
		},
		{
			name: "Error - repository response err",
			initMocks: func() {
				repositoryMock.On("Find", "Jhon", 1, 2).
					Return(UserList{}, errors.New("some error")).Once()
			},
			args: args{
				name:   "Jhon",
				offset: 2,
				size:   1,
			},
			assertMocks: func(t *testing.T) {
				repositoryMock.AssertExpectations(t)
			},
			assertError: func(t *testing.T, e error) {
				assert.NotNil(t, e)
			},
			assertFunc: func(t *testing.T, users UserList) {
				assert.Equal(t, len(users.Data), 0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initMocks()
			service := Service{
				repository: repositoryMock,
			}

			userId, err := service.Find(tt.args.name, tt.args.size, tt.args.offset)
			tt.assertMocks(t)
			tt.assertError(t, err)
			tt.assertFunc(t, userId)
		})
	}
}
