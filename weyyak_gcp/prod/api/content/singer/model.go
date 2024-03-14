package singer

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Singer - struct for DB binding
type Singer struct {
	Id          string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
}
