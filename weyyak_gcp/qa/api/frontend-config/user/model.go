package user

import (
	"time"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type NewMenu struct {
	Url string `json:"url"`
}

type Menu struct {
	ID              string `json:"id" gorm:"primary_key"`
	Device          string `json:"device" binding:"required"`
	MenuType        string `json:"menu_type"`
	MenuEnglishName string `json:"menu_english_name"`
	MenuArabicName  string `json:"menu_arabic_name"`
	SliderKey       int    `json:"slider_key"`
	Url             string `json:"url "`
	Order           int    `json:"order "`
	IsPublished     bool   `json:"is_published"`
}

/*Get User Ratings details with search text*/
type RatingRecords struct {
	Id          string    `json:"id"`
	RatedAt     time.Time `json:"ratedAt"`
	Title       string    `json:"title"`
	ContentType string    `json:"contentType"`
	DeviceId    string    `json:"deviceId"`
	Rating      float64   `json:"rating"`
	IsHidden    bool      `json:"isHidden"`
}

type RatingByUser struct {
	RatedAt             time.Time `json:"ratedAt"`
	Title               string    `json:"title"`
	ContentType         string    `json:"contentType"`
	Genres              []string  `json:"genres"`
	RatedOnPlatformName string    `json:"ratedOnPlatformName"`
	Rating              float64   `json:"rating"`
	IsHidden            bool      `json:"isHidden"`
}
type Genname struct {
	EnglishName string `json:"englishName"`
}
type PlatformName struct {
	RatedOnPlatformName string `json:"ratedOnPlatformName"`
}
// watching issue
type WatchingIssue struct {
	ReportedAt      time.Time `json:"reportedAt"`
	IsVideo         bool      `json:"isVideo"`
	IsSound         bool      `json:"isSound"`
	IsTranslation   bool      `json:"isTranslation"`
	IsCommunication bool      `json:"isCommunication"`
	Description     string    `json:"description"`
}
