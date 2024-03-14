package musiccomposer

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//MusicComposer - struct for DB binding
type MusicComposer struct {
	Id          string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
}
