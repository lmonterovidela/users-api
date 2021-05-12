package user

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRepository_Insert(t *testing.T) {
	setTestEnvironment()

	type args struct {
		user *User
	}

	tests := []struct {
		name        string
		args        args
		assertError func(*testing.T, error)
		assertFunc  func(*testing.T, int64)
	}{
		{
			name: "Success - insert first user",
			args: args{
				user: &User{
					Name:    "test",
					Dob:     time.Now(),
					Address: "address",
				},
			},
			assertError: func(t *testing.T, e error) {
				assert.Nil(t, e)
			},
			assertFunc: func(t *testing.T, id int64) {
				assert.NotZero(t, id)
			},
		},
		{
			name: "Error - name too long",
			args: args{
				user: &User{
					Name:    "test1test1test1test1test1test1test1test1test1test1",
					Dob:     time.Now(),
					Address: "address",
				},
			},
			assertError: func(t *testing.T, e error) {
				assert.NotNil(t, e)
			},
			assertFunc: func(t *testing.T, id int64) {
				assert.Zero(t, id)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := NewRepository()
			userId, err := repository.Insert(tt.args.user)
			tt.assertFunc(t, userId)
			tt.assertError(t, err)
		})
	}
}

func TestRepository_Update(t *testing.T) {
	setTestEnvironment()

	type args struct {
		userId int
		user   *User
	}

	tests := []struct {
		name        string
		args        args
		assertError func(*testing.T, error)
		assertFunc  func(*testing.T, User, User)
	}{
		{
			name: "Success - update user",
			args: args{
				userId: 1,
				user: &User{
					Name:    "test",
					Dob:     time.Now(),
					Address: "address",
				},
			},
			assertError: func(t *testing.T, e error) {
				assert.Nil(t, e)
			},
			assertFunc: func(t *testing.T, user User, userAfterUpdated User) {
				assert.Equal(t, user.Name, user.Name)
				assert.Equal(t, user.Address, user.Address)
				assert.Equal(t, user.Dob, user.Dob)
			},
		},
		{
			name: "Error - name too long",
			args: args{
				userId: 1,
				user: &User{
					Name:    "test1test1test1test1test1test1test1test1test1test1",
					Dob:     time.Now(),
					Address: "address",
				},
			},
			assertError: func(t *testing.T, e error) {
				assert.NotNil(t, e)
			},
			assertFunc: func(t *testing.T, user User, userAfterUpdated User) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var userUpdated User
			repository := NewRepository()
			err := repository.Update(tt.args.userId, tt.args.user)
			if err == nil {
				userUpdated, _ = repository.Get(tt.args.userId)
			}

			tt.assertFunc(t, *tt.args.user, userUpdated)
			tt.assertError(t, err)
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	setTestEnvironment()

	type args struct {
		userId int
	}

	tests := []struct {
		name        string
		initMocks   func()
		args        args
		assertError func(*testing.T, error)
	}{
		{
			name: "Success - delete user",
			args: args{
				userId: 1,
			},
			assertError: func(t *testing.T, e error) {
				assert.Nil(t, e)
			},
		},
		{
			name: "Error - not found",
			args: args{
				userId: 1001,
			},
			assertError: func(t *testing.T, e error) {
				assert.NotNil(t, e)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := NewRepository()
			err := repository.Delete(tt.args.userId)
			tt.assertError(t, err)
		})
	}
}

func TestRepository_Get(t *testing.T) {
	setTestEnvironment()

	type args struct {
		userId int
	}

	tests := []struct {
		name        string
		initMocks   func()
		args        args
		assertError func(*testing.T, error)
		assertFunc  func(*testing.T, User)
	}{
		{
			name: "Success - get user",
			args: args{
				userId: 1,
			},
			assertError: func(t *testing.T, e error) {
				assert.Nil(t, e)
			},
			assertFunc: func(t *testing.T, user User) {
				assert.Equal(t, user.Name, "Jhon")
			},
		},
		{
			name: "Error - not found",
			args: args{
				userId: 1001,
			},
			assertError: func(t *testing.T, e error) {
				assert.NotNil(t, e)
			},
			assertFunc: func(t *testing.T, id User) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := NewRepository()
			user, err := repository.Get(tt.args.userId)
			tt.assertError(t, err)
			tt.assertFunc(t, user)
		})
	}
}

func TestRepository_Find(t *testing.T) {
	setTestEnvironment()

	type args struct {
		name   string
		offset int
		size   int
	}

	tests := []struct {
		name       string
		initMocks  func()
		args       args
		assertFunc func(*testing.T, UserList)
	}{
		{
			name: "Success - retrieve 5 elements",
			args: args{
				name:   "Jhon",
				offset: 0,
				size:   5,
			},
			assertFunc: func(t *testing.T, user UserList) {
				assert.Equal(t, len(user.Data), 5)
				assert.Equal(t, user.Data[0].Id, 1)
				assert.Equal(t, user.Data[4].Id, 5)
			},
		},
		{
			name: "Success - retrieve 3 elements, from 3 to 5",
			args: args{
				name:   "Jhon",
				offset: 2,
				size:   3,
			},
			assertFunc: func(t *testing.T, user UserList) {
				assert.Equal(t, len(user.Data), 3)
				assert.Equal(t, user.Data[0].Id, 3)
				assert.Equal(t, user.Data[2].Id, 5)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := NewRepository()
			user, _ := repository.Find(tt.args.name, tt.args.size, tt.args.offset)
			tt.assertFunc(t, user)
		})
	}
}

func setTestEnvironment() {
	viper.Set("env", "test")
	viper.Set("database.host", "localhost:3305")
	viper.Set("database.pass", "root")
	viper.Set("database.user", "root")
	viper.Set("database.name", "challenge")
}
