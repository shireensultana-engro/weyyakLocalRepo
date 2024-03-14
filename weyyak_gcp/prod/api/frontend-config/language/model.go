package language

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//LanguageOriginType - struct for DB binding
type LanguageOriginType struct {
	Id   string `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

//LanguageDialect - struct for DB binding
type LanguageDialect struct {
	Id          int    `json:"id" gorm:"primary_key"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
}

//LanguageDubbing - struct for DB binding
type LanguageDubbing struct {
	Id          int    `json:"id" gorm:"primary_key"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
	Code        string `json:"code"`
}

//LanguageSubtitles - struct for DB binding
type LanguageSubtitles struct {
	Id          int    `json:"id" gorm:"primary_key"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
	Code        string `json:"code"`
}

//Language - struct for DB binding
type Language struct {
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
	Code        string `json:"code"`
}

//LanguageAvailable - struct for DB binding
type LanguageAvailable struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
