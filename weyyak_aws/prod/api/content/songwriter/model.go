package songwriter

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//SongWriter - struct for DB binding
type SongWriter struct {
	Id          string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
}
