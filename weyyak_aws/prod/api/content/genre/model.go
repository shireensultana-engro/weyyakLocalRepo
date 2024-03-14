package genre

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Genre - struct for DB binding
type Genre struct {
	Id          string `json:"id" gorm:"primary_key"`
	EnglishName string `json:"english_name"`
	ArabicName  string `json:"arabic_name"`
}

type GenreList struct {	
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
	Id          string `json:"id"`
}

type SubGenre struct {
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
	Id          string `json:"id"`
}
