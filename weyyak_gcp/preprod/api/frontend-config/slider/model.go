package slider

import (
	"time"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type SliderAvailability struct {
	Id         string `json:"id" binding:"required"`
	IsDisabled *bool  `json:"isDisabled"  binding:"required"`
}

type SliderUpdate struct {
	IsDisabled *bool `json:"isDisabled" binding:"required"`
}

//slider_type
type sliderTypes struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type SliderList struct {
	ID                  string `json:"id" gorm:"primary_key"`
	Name                string `json:"name"`
	IsDisabled          bool   `json:"isDisabled"`
	SchedulingEndDate   string `json:"schedulingEndDate"`
	IsDefaultForAnyPage bool   `json:"isDefaultForAnyPage"`
	AvailableOn         string `json:"availableOn"`
	Region              string `json:"region"`
	HasMoreRegions      bool   `json:"hasMoreRegions"`
	FoundIn             string `json:"foundIn"`
}
type SliderDetails struct {
	Id                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Pages               []SliderPageDetails    `json:"pages"`
	BlackAreaPlaylistId string                 `json:"-"`
	GreenAreaPlaylistId string                 `json:"-"`
	RedAreaPlaylistId   string                 `json:"-"`
	BlackAreaPlaylist   *SliderPlaylistDetails `json:"blackAreaPlaylist"`
	GreenAreaPlaylist   *SliderPlaylistDetails `json:"greenAreaPlaylist"`
	RedAreaPlaylist     *SliderPlaylistDetails `json:"redAreaPlaylist"`
	PublishingPlatforms []int                  `json:"publishingPlatforms"`
	Regions             []int                  `json:"regions"`
	SchedulingEndDate   time.Time              `json:"schedulingEndDate"`
	SchedulingStartDate time.Time              `json:"schedulingStartDate"`
	Type                int                    `json:"type"`
}
type SliderPlaylistDetails struct {
	Id                  string    `json:"id"`
	ArabicTitle         string    `json:"arabicTitle"`
	EnglishTitle        string    `json:"englishTitle"`
	IsDisabled          bool      `json:"isDisabled"`
	PublishingPlatforms *int      `json:"publishingPlatforms"`
	SchedulingStartDate time.Time `json:"-"`
	SchedulingEndDate   time.Time `json:"-"`
}
type SliderPageDetails struct {
	IsDefault           bool   `json:"isDefault"`
	EnglishTitle        string `json:"englishTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	IsDisabled          bool   `json:"isDisabled"`
	IsHome              bool   `json:"isHome"`
	PublishingPlatforms []int  `json:"publishingPlatforms"`
	Id                  string `json:"id"`
	Platforms           string `json:"platforms,omitempty"`
}
type SliderTargetPlatform struct {
	SliderId       string `json:"slider_id"`
	TargetPlatform int    `json:"target_platform"`
}
type SliderCountry struct {
	SliderId  string `json:"slider_id"`
	CountryId int    `json:"country_id"`
}

// for Update slider regions type is []int
type CreateUpdateSliderRequest struct {
	Name                string   `json:"name"`
	Type                int      `json:"type"`
	BlackAreaPlaylistId string   `json:"blackAreaPlaylistId"`
	RedAreaPlaylistId   string   `json:"redAreaPlaylistId"`
	GreenAreaPlaylistId string   `json:"greenAreaPlaylistId"`
	SchedulingStartDate string   `json:"schedulingStartDate"`
	SchedulingEndDate   string   `json:"schedulingEndDate"`
	PublishingPlatforms *[]int   `json:"publishingPlatforms"`
	Regions             *[]int   `json:"regions"`
	PagesIds            []string `json:"pagesIds"`
}

// for create new slider regions type is []string
type CreateUpdateSliderRequestCreate struct {
	SliderId            string   `json:"id"` // sliderid for creating old sliders with .net
	SliderKey           int      `json:"sliderKey"`
	Name                string   `json:"name"`
	Type                int      `json:"type"`
	BlackAreaPlaylistId string   `json:"blackAreaPlaylistId"`
	RedAreaPlaylistId   string   `json:"redAreaPlaylistId"`
	GreenAreaPlaylistId string   `json:"greenAreaPlaylistId"`
	SchedulingStartDate string   `json:"schedulingStartDate"`
	SchedulingEndDate   string   `json:"schedulingEndDate"`
	PublishingPlatforms *[]int   `json:"publishingPlatforms"`
	Regions             *[]int   `json:"regions"`
	PagesIds            []string `json:"pagesIds"`
}
type Slider struct {
	Id                  string    `json:"id" gorm:"primary_key"`
	Name                string    `json:"name"`
	Type                int       `json:"type"`
	CreatedAt           time.Time `json:"created_at"`
	DeletedByUserId     *string   `json:"deleted_by_user_id"`
	IsDisabled          bool      `json:"is_disabled"`
	SchedulingStartDate string    `json:"scheduling_start_date"`
	SchedulingEndDate   string    `json:"scheduling_end_date"`
	BlackAreaPlaylistId string    `json:"black_area_playlist_id"`
	RedAreaPlaylistId   string    `json:"red_area_playlist_id"`
	GreenAreaPlaylistId string    `json:"green_area_playlist_id"`
	SliderKey           int       `json:"slider_key"`
	ModifiedAt          time.Time `json:"modified_at"`
}
type PageSlider struct {
	PageId   string `json:"page_id"`
	SliderId string `json:"slider_id"`
	Order    int    `json:"order"`
}
type PreviewLayouts struct {
	PreviewImageKey string `json:"previewImageUrl"`
	SliderType      int    `json:"sliderType"`
	Platform        int    `json:"platform"`
	Id              string `json:"id"`
}

type PageSummary struct {
	Name       string `json:"name"`
	IsDisabled bool   `json:"isDisabled"`
	Id         string `json:"id"`
	Details    string `json:"details"`
}

type UpdateDetails struct {
	DeletedByUserId string    `json:"deleted_by_user_id"`
	ModifiedAt      time.Time `json:"modified_at"`
}
type PlaylistContents struct {
	ContentId    string `json:"content_id"`
	PlaylistName string `json:"playlist_name"`
}
type SliderNotificationError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Type    int    `json:"type"`
}
type PlaylistContentsCount struct {
	ContentKey          int    `jsom:"id"`
	ContentType         string `json:"content_type"`
	TransliteratedTitle string `json:"transliterated_title"`
}

/*sliderRegions*/
type Country struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}
