package viewactivity

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)


type ViewActivityDetails struct {
	ContentTypeName string `json:content_type_name`
}