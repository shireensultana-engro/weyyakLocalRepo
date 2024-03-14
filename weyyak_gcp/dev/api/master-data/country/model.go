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
	RegionId	string `json:"regionId"`
	CallingCode string `json:"callingCode"`
	Alpha2Code  string `json:"alpha2Code"`
}

//Country - EN
type CountryEN struct {
	Id          int8   `json:"id" gorm:"primary_key"`
	EnglishName string `json:"name"`
	CallingCode string `json:"code"`
}

//Country - AR
type CountryAR struct {
	Id          int8   `json:"id" gorm:"primary_key"`
	ArabicName  string `json:"name"`
	CallingCode string `json:"code"`
}
// 
type Countries struct {
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`	
	Id          int    `json:"id" gorm:"primary_key"`
}
