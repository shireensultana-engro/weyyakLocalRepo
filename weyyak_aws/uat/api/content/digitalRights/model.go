package digitalRights

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//ContentRights - struct for DB binding
type ContentRights struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

// DisplayStatus - struct for DB binding
type DisplayStatus struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

/*  Get digital rights regions*/
type Continent struct {
	Id   string `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
type Region struct {
	Id          string `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	ContinentId string `json:"continentId"`
}
type Country struct {
	Id          int    `json:"id" gorm:"primary_key"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
	RegionId    string `json:"regionId"`
	CallingCode string `json:"callingCode"`
}

type DigitalrightsResponse struct {
	Name      string            `json:"name"`
	Regions   []RegionsResponse `json:"regions"`
	Countries *string           `json:"countries"`
}
type RegionsResponse struct {
	Name      string             `json:"name"`
	Countries []CountrysResponse `json:"countries"`
}
type CountrysResponse struct {
	Name        string `json:"name"`
	CallingCode *string `json:"callingCode"`
	Id          int    `json:"id"`
}
