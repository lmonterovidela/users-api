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


type Location struct {
	Type     string   `json:"type"`
	Query    []string `json:"query"`
	Features []struct {
		ID         string   `json:"id"`
		Type       string   `json:"type"`
		PlaceType  []string `json:"place_type"`
		Text      string    `json:"text"`
		PlaceName string    `json:"place_name"`
		Bbox      []float64 `json:"bbox,omitempty"`
		Center    []float64 `json:"center"`
		Geometry  struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Context []struct {
			ID        string `json:"id"`
			Wikidata  string `json:"wikidata"`
			ShortCode string `json:"short_code"`
			Text      string `json:"text"`
		} `json:"context"`
	} `json:"features"`
}
