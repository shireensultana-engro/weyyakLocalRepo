package contentType

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//ContentOnetierTypes - struct for DB binding
type ContentOnetierTypes struct {
	Id          int    `json:"id" gorm:"primary_key"`
	ContentType string `json:"contentType"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
}

//ContentMultitierTypes - struct for DB binding
type ContentMultitierTypes struct {
	Id          int    `json:"id" gorm:"primary_key"`
	ContentType string `json:"contentType"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
}

// TextualDataTags
type TextualDataTag struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

