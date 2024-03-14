package actor

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Actor - struct for DB binding
type Actor struct {
	Id          string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
}
