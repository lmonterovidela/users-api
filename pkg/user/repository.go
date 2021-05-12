package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/users-api/infrastructure"
	"net/http"
	"time"
)

type IRepository interface {
	Insert(user *User) (int64, error)
	Update(userId int, user *User) error
	Get(userId int) (User, error)
	Find(name string, size int, offset int) (UserList, error)
	Delete(userId int) error

	GetLocation(userId int) (Location, error)
}

type Repository struct {
	db           *sqlx.DB
	mapApiClient infrastructure.RestClient
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

	locationUrl       string = "/geocoding/v5/mapbox.places/%s.json?access_token=%s"
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

func (r *Repository) GetLocation(userId int) (Location, error) {
	user := User{}
	err := r.db.Get(&user, getUserDataSQL, userId)
	if err != nil {
		return Location{}, err
	}

	if err == sql.ErrNoRows {
		return Location{}, &NotFoundError{Message: fmt.Errorf(userNotFound, userId)}
	}

	path:= fmt.Sprintf(locationUrl, user.Address, viper.Get("clients.map.token"))
	response, err := r.mapApiClient.Get(path, nil, nil)
	defer response.Body.Close()

	if !(response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices){
		return Location{}, fmt.Errorf("unexpected error")
	}

	var location Location
	err = json.NewDecoder(response.Body).Decode(&location)

	if err != nil {
		logrus.Errorf("error while unmarshalling response path: %s, body: %v - response: %v", path, err, response.Body)
		return Location{}, err
	}

	return location, err
}

func NewRepository() IRepository {

	return &Repository{
		db: infrastructure.ConnectDatabase(),
		mapApiClient: infrastructure.NewRestClient(viper.GetString("clients.map.base_url"), time.Duration(viper.GetInt("clients.map.timeout"))),
	}
}
