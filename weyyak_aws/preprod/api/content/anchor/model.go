package anchor

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Anchor struct {
	Id            string        `json:"id" gorm:"primary_key"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	PhotoUrl      string        `json:"photo_url"`
	Email         string        `json:"email"`
	Shows         []AnchorShows `json:"shows"`
	AbouTheAnchor string        `json:"aboutTheAnchor"`
	Status        bool          `json:"status"`
	Photo         bool          `json:"photo"`
	HasDeleted    bool          `json:"hasDeleted"`
}

type AnchorShows struct {
	AnchorId string `json:"anchor_id"`
	ShowName string `json:"show_name"`
	Timing   string `json:"timing"`
}

type CreateAnchor struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	AbouTheAnchor string `json:"aboutTheAnchor"`
	Status        bool   `json:"status"`
	Photo         string `json:"photo"`
}
