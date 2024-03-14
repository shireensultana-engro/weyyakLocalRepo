package country

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Country - struct for DB binding
type Country struct {
	Id          int    `json:"id" gorm:"primary_key"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
	RegionId    string `json:"regionId"`
	CallingCode string `json:"callingCode"`
	Alpha2Code  string `json:"alpha2Code"`
}

type CountryDetails struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
