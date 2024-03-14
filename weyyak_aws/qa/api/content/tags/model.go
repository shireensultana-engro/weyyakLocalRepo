package tags

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//TextualDataTag - struct for DB binding
type TextualDataTag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type Name struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
