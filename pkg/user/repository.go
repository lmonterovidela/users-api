package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/users-api/infrastructure"
)

type IRepository interface {
	Insert(user *User) (int64, error)
	Update(userId int, user *User) error
	Get(userId int) (User, error)
	Find(name string, size int, offset int) (UserList, error)
	Delete(userId int) error
}

type Repository struct {
	db *sqlx.DB
}

type NotFoundError struct {
	Message error
}

func (e *NotFoundError) Error() string {
	return e.Message.Error()
}

const (
	insertUserSQL     string = "INSERT INTO user (name, address, dob) VALUES (:name,:address,:dob)"
	getUserDataSQL    string = "SELECT id, name, address, dob, created_at, updated_at FROM user WHERE id = ?"
	updateUserDataSQL string = "UPDATE user SET name=:name, address=:address, dob=:dob WHERE id =:id"
	findUserDataSQL   string = "SELECT id, name, address, dob, created_at, updated_at FROM user WHERE name = ? Limit ?  Offset ?"
	deleteUserSQL     string = "DELETE FROM user where id = ?"

	userNotFound string = "user with id=%d not found"
	defaultSize  int    = 20
)

func (r *Repository) Insert(user *User) (int64, error) {

	result, err := r.db.NamedExec(insertUserSQL, user)
	if err != nil {
		return 0, errors.New("error while creating service")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("error while creating service")
	}

	return id, nil
}

func (r *Repository) Update(userId int, user *User) error {
	user.Id = userId
	_, err := r.db.NamedExec(updateUserDataSQL, user)

	return err
}

func (r *Repository) Get(userId int) (User, error) {
	user := User{}
	err := r.db.Get(&user, getUserDataSQL, userId)

	if err == sql.ErrNoRows {
		return user, &NotFoundError{Message: fmt.Errorf(userNotFound, userId)}
	}

	return user, err
}

func (r *Repository) Find(name string, size int, offset int) (UserList, error) {
	if size == 0 {
		size = defaultSize
	}

	users := make([]User, 0)

	err := r.db.Select(&users, findUserDataSQL, name, size, offset)
	if err != nil {
		return UserList{}, err
	}

	return UserList{Data: users, Size: size, Offset: offset}, nil
}

func (r *Repository) Delete(userId int) error {
	//physical deletion is developed instead of logical deletion due to lack of context information
	result, err := r.db.Exec(deleteUserSQL, userId)
	if err == nil {
		if rowsAffected, _ := result.RowsAffected(); rowsAffected < 1 {
			return &NotFoundError{Message: fmt.Errorf(userNotFound, userId)}
		}
	}

	return err
}

func NewRepository() IRepository {

	return &Repository{
		db: infrastructure.ConnectDatabase(),
	}
}
