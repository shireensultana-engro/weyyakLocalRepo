package ageRatings

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//AgeRatings - struct for DB binding
type AgeRatings struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Code string `json:"code"`
	Name string `json:"name"`
}
