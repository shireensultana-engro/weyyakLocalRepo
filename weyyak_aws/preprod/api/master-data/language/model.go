package language

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//LanguageOriginTypes - struct for DB binding
type LanguageOriginTypes struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

//LanguageDialects - struct for DB binding
type LanguageDialects struct {
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

//LanguageSubtitle - struct for DB binding
type LanguageSubtitle struct {
	Id          int    `json:"id" gorm:"primary_key"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
	Code        string `json:"code"`
}

//LanguageOriginal - struct for DB binding
type LanguageOriginal struct {
	Id          int    `json:"id" gorm:"primary_key"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
	Code        string `json:"code"`
}
