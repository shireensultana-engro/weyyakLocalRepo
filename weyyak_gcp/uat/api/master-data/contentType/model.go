package contentType

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//ContentOneTier - struct for DB binding
type ContentOneTier struct {
	Id          int    `json:"id" gorm:"primary_key"`
	ContentType string `json:"contentType"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
}

//ContentMultiTier - struct for DB binding
type ContentMultiTier struct {
	Id          int    `json:"id" gorm:"primary_key"`
	ContentType string `json:"contentType"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
}
