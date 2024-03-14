package digitalRights

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//DigitalRights - struct for DB binding
type DigitalRights struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

// DisplayStatus - struct for DB binding
type DisplayStatus struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
