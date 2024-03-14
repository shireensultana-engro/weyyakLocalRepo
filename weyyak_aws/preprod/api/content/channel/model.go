package channel

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//channel program list - struct for DB binding
type ChannelProgramList struct {
	Name       string `json:"name"`
	Url        string `json:"url"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Duration   string `json:"duration"`
	Channel    string `json:"channel"`
	Site       string `json:"site"`
	Lang       string `json:"lang"`
	Logo       bool   `json:"logo"`
	HasDeleted bool   `json:"hasDeleted"`
}

//channel program list - struct for DB binding
type CreateChannelProgramList struct {
	Name       string `json:"name"`
	Url        string `json:"url"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Duration   string `json:"duration"`
	Channel    string `json:"channel"`
	Site       string `json:"site"`
	Lang       string `json:"lang"`
	Logo       string `json:"logo"`
	HasDeleted bool   `json:"hasDeleted"`
}

