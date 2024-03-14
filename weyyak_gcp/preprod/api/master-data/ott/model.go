package ott

import ()

//TopMenuDetails -- for db response
type TopMenuDetails struct {
	Device    string `json:"device"`
	MenuType  string `json:"menuType"`
	Title     string `json:"title"`
	SliderKey int    `json:"sliderKey"`
	Url       string `json:"url"`
	Order     int    `json:"order"`
}

//ApplicationSetting -- for db binding
type ApplicationSetting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//ConfigurationDetails
type ConfigurationDetails struct {
	Default interface{} `json:"default"`
	Ios     interface{} `json:"ios"`
	Android interface{} `json:"android"`
}
