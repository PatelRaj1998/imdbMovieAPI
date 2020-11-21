package models

type Movie struct {
	Title        string   `json:"title,omitempty" bson:"title,omitempty"`
	ReleasedYear int      `json:"releasedYear,omitempty" bson:"releasedYear,omitempty"`
	Rating       float64  `json:"rating,omitempty" bson:"rating,omitempty"`
	Id           string   `json:"id,omitempty" bson:"id,omitempty"`
	Genres       []string `json:"genres,omitempty" bson:"genres,omitempty"`
}
