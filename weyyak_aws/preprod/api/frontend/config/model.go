package config

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

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
	Web     interface{} `json:"web"`
}

type RedisCacheRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RedisCacheResponse struct {
	Value string `json:"value"`
	Error string `json:"error"`
}
