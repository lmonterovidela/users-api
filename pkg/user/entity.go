package user

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type UserList struct {
	Data   []User `json:"data"`
	Size   int    `json:"size"`
	Offset int    `json:"offset"`
}

type User struct {
	Id        int       `db:"id" json:"id,omitempty"`
	Address   string    `db:"address" json:"address,omitempty"`
	Dob       time.Time `db:"dob" json:"dob,omitempty"`
	Name      string    `db:"name" json:"name,omitempty" validate:"required,min=3"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (u *User) Validate() error {
	return validator.New().Struct(u)
}
