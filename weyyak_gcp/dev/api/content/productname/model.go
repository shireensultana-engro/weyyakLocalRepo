package productname

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ProductName struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}


