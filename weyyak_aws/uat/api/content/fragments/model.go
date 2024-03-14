package fragments

import (
	"time"

	_ "github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

/*Content Fragment */
type ContentType struct {
	ContentTier int
}
type ContentImageryDetails struct {
	Thumbnail   string `json:"thumbnail"`
	Backdrop    string `json:"backdrop"`
	MobileImg   string `json:"mobile_img"`
	FeaturedImg string `json:"featured_img"`
	Banner      string `json:"banner"`
}
type Movie struct {
	ID                   int       `json:"id"`
	Title                string    `json:"title"`
	Geoblock             bool      `json:"geoblock"`
	DigitalRightType     int       `json:"digitalRighttype"`
	DigitalRightsRegions []int     `json:"digitalRightsRegions"`
	SubscriptiontPlans   []int     `json:"subscriptiontPlans"`
	InsertedAt           time.Time `json:"insertedAt"`
}
type ContentFragmentDetails struct {
	Id              int         `json:"id"`
	Cast            []string    `json:"cast"`
	Tags            []string    `json:"tags"`
	Title           string      `json:"title"`
	Genres          []string    `json:"genres"`
	Length          int         `json:"length"`
	Movies          []Movie     `json:"movies,omitempty"`
	Imagery         interface{} `json:"imagery"`
	Geoblock        bool        `json:"geoblock"`
	Synopsis        string      `json:"synopsis"`
	VideoId         string      `json:"video_id"`
	SeoTitle        string      `json:"seo_title"`
	ContentId       string      `json:"content_id"`
	AgeRating       string      `json:"age_rating"`
	InsertedAt      string      `json:"insertedAt"`
	MainActor       string      `json:"main_actor"`
	ModifiedAt      string      `json:"modifiedAt"`
	ContentType     string      `json:"content_type"`
	FriendlyUrl     string      `json:"friendly_url"`
	MainActress     string      `json:"main_actress"`
	Seasons         []Seasons   `json:"seasons,omitempty"`
	ProductionYear  int         `json:"production_year"`
	SeoDescription  string      `json:"seo_description"`
	TranslatedTitle string      `json:"translated_title"`
}
type ContentFragment struct {
	Id                string         `json:"id"`
	ContentId         string         `json:"contentId"`
	ContentVarianceId string         `json:"contentVarianceId"`
	RightsStartDate   time.Time      `json:"rightsStartDate"`
	RightsEndDate     time.Time      `json:"rightsEndDate"`
	Details           postgres.Jsonb `json:"details"`
	Country           string         `json:"country"`
	Platform          string         `json:"platform"`
	Language          string         `json:"language"`
	ContentType       string         `json:"contentType"`
	ContentKey        int            `json:"contentKey"`
}
type ContentDetails struct {
	ContentId              string
	Id                     int
	EnCast                 string
	ArCast                 string
	Title                  string
	ArabicTitlen           string //for fetching arabic title
	EnGenres               string
	ArGenres               string
	Length                 int
	EnSeoTitle             string
	ArSeoTitle             string
	FriendlyUrl            string
	ProductionYear         int
	EnglishSynopsis        string
	ArabicSynopsis         string
	ContentType            string
	HasAllRights           bool
	Country                string
	ModifiedAt             time.Time
	InsertedAt             time.Time
	EnglishSeoDescription  string
	ArabicSeoDescription   string
	EnAgeRating            string
	ArAgeRating            string
	EnMainActor            string
	ArMainActor            string
	EnMainActress          string
	ArMainActress          string
	VideoId                string
	Tags                   string
	Platforms              string
	Geoblock               bool
	DigitalRightType       int
	TranslatedTitle        string
	ContentVarianceId      string
	SubscriptionPlans      string
	HasPosterImage         bool
	Imagery                postgres.Jsonb
	Dubbed                 bool
	DigitalRightsStartDate time.Time
	DigitalRightsEndDate   time.Time
	SeasonId               string
}
type MissingRights struct {
	DigitalRightType       int
	DigitalRightsStartDate time.Time
	DigitalRightsEndDate   time.Time
}
type Seasons struct {
	DigitalRightsRegions *string `json:"digitalRightsRegions"`
	DigitalRighttype     int     `json:"digitalRighttype"`
	Dubbed               bool    `json:"dubbed"`
	Geoblock             bool    `json:"geoblock"`
	Id                   int     `json:"id"`
	SeasonNumber         int     `json:"season_number"`
	SeoDescription       string  `json:"seo_description"`
	SeoTitle             string  `json:"seo_title"`
	SubscriptionPlans    []int   `json:"subscriptionPlans"`
	Title                string  `josn:"title"`
}
type RegionsAndPlatforms struct {
	Country   string
	Platforms string
}

type RequestIds struct {
	Ids []string `json:"ids"`
}

type Ids struct {
	Id string `json:"id"`
}
