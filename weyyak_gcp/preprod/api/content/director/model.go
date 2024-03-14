package director

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Director - struct for DB binding
type Director struct {
	Id          string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
}
