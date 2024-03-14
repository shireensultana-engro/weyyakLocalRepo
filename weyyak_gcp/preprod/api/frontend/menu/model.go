package menu

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type MenuFragmentResponse struct {
	PageId           string `json:"page_id"`
	Country          string `json:"country"`
	Platform         string `json:"platform"`
	FragmentResponse string `json:"fragment_response"`
	PageOrder        int    `json:"page_order"`
	Language         string `json:"language"`
	FragmentType     string `json:"fragment_type"`
	PageKey          int    `json:"page_key"`
}
type MenuFragment struct {
	PageId           string         `json:"page_id"`
	Country          string         `json:"country"`
	Platform         string         `json:"platform"`
	FragmentResponse postgres.Jsonb `json:"fragment_response"`
	PageOrder        int            `json:"page_order"`
	Language         string         `json:"language"`
	FragmentType     string         `json:"fragment_type"`
	PageKey          int            `json:"page_key"`
}
type PageFragment struct {
	PageId    string         `json:"page_id"`
	PageOrder int            `json:"page_order"`
	PageKey   int            `json:"page_key"`
	Country   string         `json:"country"`
	Platform  string         `json:"platform"`
	Details   postgres.Jsonb `json:"details"`
	Language  string         `json:"language"`
}
type GetPageFragment struct {
	PageId    string `json:"page_id"`
	PageOrder int    `json:"page_order"`
	PageKey   int    `json:"page_key"`
	Country   string `json:"country"`
	Platform  string `json:"platform"`
	Details   string `json:"details"`
	Language  string `json:"language"`
}
type SliderFragment struct {
	PageId   string         `json:"page_id"`
	PageKey  int            `json:"page_key"`
	SliderId string         `json:"slider_id"`
	Country  string         `json:"country"`
	Platform string         `json:"platform"`
	Details  postgres.Jsonb `json:"details"`
	Language string         `json:"language"`
}
type PlaylistFragment struct {
	PageId     string         `json:"page_id"`
	PageKey    int            `json:"page_key"`
	PlaylistId string         `json:"playlist_id"`
	Country    string         `json:"country"`
	Platform   string         `json:"platform"`
	Details    postgres.Jsonb `json:"details"`
	Language   string         `json:"language"`
}

// SideMenuDetails struct for DB binding
type SideMenuDetails struct {
	Total       int64      `json:"total"`
	PerPage     int64      `json:"per_page"`
	CurrentPage int64      `json:"current_page"`
	LastPage    int64      `json:"last_page"`
	NextPageUrl string     `json:"next_page_url"`
	PrevPageUrl string     `json:"prev_page_url"`
	From        int64      `json:"from"`
	To          int64      `json:"to"`
	Data        []MenuPage `json:"data"`
}

//GetMenu struct for DB binding
type GetMenu struct {
	Data MenuPageDetails `json:"data"`
}

//MenuPageDetails struct for DB binding
type MenuPageDetails struct {
	ID             int              `json:"id"`
	FriendlyUrl    string           `json:"friendly_url"`
	SeoDescription string           `json:"seo_description"`
	Title          string           `json:"title"`
	Type           string           `json:"type"`
	Featured       *FeaturedDetails `json:"featured"`
	Playlists      []MenuPlaylists  `json:"playlists"`
	Imagery        ImageryDetails   `json:"imagery"`
}

//ImageryDetails struct for DB binding
type ImageryDetails struct {
	MobileMenu            string `json:"mobile-menu"`
	MobilePosterImage     string `json:"menu-poster-image"`
	MobileMenuPosterImage string `json:"mobile-menu-poster-image"`
}

//ContentImageryDetails for DB binding
type ContentImageryDetails struct {
	Thumbnail   string `json:"thumbnail"`
	Backdrop    string `json:"backdrop"`
	MobileImg   string `json:"mobile_img"`
	FeaturedImg string `json:"featured_img"`
	Banner      string `json:"banner"`
}
type PageContent struct {
	Key            string         `json:"key"`
	ID             int            `json:"id"`
	FriendlyUrl    string         `json:"friendly_url"`
	SeoDescription string         `json:"seo_description"`
	Title          string         `json:"title"`
	Type           string         `json:"type"`
	Imagery        ImageryDetails `json:"imagery"`
}

//MenuPlaylists struct for DB binding
type MenuPlaylists struct {
	ID           int32             `json:"id"`
	Title        string            `json:"title"`
	Content      []PlaylistContent `json:"content"`
	PlaylistType string            `json:"playlisttype"`
	PageContent  []PageContent     `json:"pagecontent"`
}

//PlaylistContent struct for DB binding
type PlaylistContent struct {
	ContentId       string                   `json:"content_id"`
	ID              int                      `json:"id"`
	AgeRating       string                   `json:"age_rating"`
	VideoId         string                   `json:"video_id"`
	FriendlyUrl     string                   `json:"friendly_url"`
	ContentType     string                   `json:"content_type"`
	Synopsis        string                   `json:"synopsis"`
	ProductionYear  *int32                   `json:"production_year"`
	Length          *int32                   `json:"length"`
	Title           string                   `json:"title"`
	Cast            []string                 `json:"cast"`
	SeoDescription  string                   `json:"seo_description"`
	TranslatedTitle string                   `json:"translated_title"`
	Genres          []string                 `json:"genres"`
	Tags            []string                 `json:"tags"`
	SeoTitle        string                   `json:"seo_title"`
	Imagery         interface{}              `json:"imagery"`
	Seasons         []PlaylistContentSeasons `json:"seasons,omitempty"`
	Movies          []PlaylistMovie          `json:"movies,omitempty"`
	// InsertedAt      time.Time                `json:"insertedAt"`
	// ModifiedAt      time.Time                `json:"modifiedAt"`
	Geoblock    bool   `json:"geoblock"`
	MainActor   string `json:"main_actor"`
	MainActress string `json:"main_actress"`
}

//PlaylistContentSeasons struct for DB binding
type PlaylistContentSeasons struct {
	ID                   int32  `json:"id"`
	SeasonNumber         int32  `json:"season_number,omitempty"`
	Dubbed               bool   `json:"dubbed"`
	SeoDescription       string `json:"seo_description,omitempty"`
	SeoTitle             string `json:"seo_title,omitempty"`
	Title                string `json:"title"`
	Geoblock             bool   `json:"geoblock"`
	DigitalRightType     int    `json:"digitalRighttype"`
	DigitalRightsRegions []int  `json:"digitalRightsRegions"`
	SubscriptiontPlans   []int  `json:"subscriptiontPlans"`
}

//FeaturedDetails struct for DB binding
type FeaturedDetails struct {
	ID               int64               `json:"id"`
	Type             string              `json:"type"`
	PlayWatchTrailer bool                `json:"play_watch_trailer"`
	MoreButton       bool                `json:"more_button"`
	Playlists        []FeaturedPlaylists `json:"playlists"`
}

//FeaturedPlaylists struct for DB binding
type FeaturedPlaylists struct {
	ID           int32             `json:"id"`
	PlaylistType string            `json:"playlist_type"`
	Content      []PlaylistContent `json:"content"`
}
type FragmentResponse struct {
	EnResponseData postgres.Jsonb `json:"en_response"`
	ArResponseData postgres.Jsonb `json:"ar_response"`
}
type PageDataSyncRequest struct {
	PageId          string                `json:"page_id"`
	PageKey         int                   `json:"page_key"`
	Country         string                `json:"country"`
	PageOrder       []PageOrderDetails    `json:"page_order"`
	PageDetails     PageDetails           `json:"page_details"`
	SliderDetails   []PageSliderDetails   `json:"slider_details"`
	PlaylistDetails []PagePlaylistDetails `json:"playlist_details"`
}
type PageOrderDetails struct {
	TargetPlarform  int `json:"target_plarform"`
	PageOrderNumber int `json:"page_order_number"`
}
type PageDetails struct {
	En PageLanguageDetails `json:"en"`
	Ar PageLanguageDetails `json:"ar"`
}
type PageLanguageDetails struct {
	ID             int            `json:"id"`
	FriendlyUrl    string         `json:"friendly_url"`
	SeoDescription string         `json:"seo_description"`
	Title          string         `json:"title"`
	Type           string         `json:"type"`
	Imagery        ImageryDetails `json:"imagery"`
}
type PageSliderDetails struct {
	SliderId string        `json:"slider_id"`
	Details  SliderDetails `json:"details"`
}
type SliderDetails struct {
	En FeaturedDetails `json:"en"`
	Ar FeaturedDetails `json:"ar"`
}
type PagePlaylistDetails struct {
	PlaylistId string          `json:"playlist_id"`
	Details    PlaylistDetails `json:"details"`
}
type PlaylistDetails struct {
	En MenuPlaylists `json:"en"`
	Ar MenuPlaylists `json:"ar"`
}
type SliderDataSyncRequest struct {
	SliderId             string                 `json:"slider_id"`
	Country              string                 `json:"country"`
	PublishingPlatforms  []int                  `json:"publishing_platforms"`
	SliderDetails        SliderDetails          `json:"slider_details"`
	SliderAvailablePages []SliderAvailablePages `json:"slider_available_pages"`
}
type SliderAvailablePages struct {
	PageId  string `json:"page_id"`
	PageKey int    `json:"page_key"`
}
type PlaylistDataSyncRequest struct {
	PlaylistId             string                 `json:"playlist_id"`
	Country                string                 `json:"country"`
	PublishingPlatforms    []int                  `json:"publishing_platforms"`
	PlaylistDetails        PlaylistDetails        `json:"playlist_details"`
	PlaylistAvailablePages []SliderAvailablePages `json:"playlist_available_pages"`
}

//Menu - struct for DB binding
type Menu struct {
	Id              string `json:"id" gorm:"primary_key"`
	Device          string `json:"device"`
	MenuType        string `json:"menu_type"`
	MenuEnglishName string `json:"menu_english_name"`
	MenuArabicName  string `json:"menu_arabic_name"`
	SliderKey       int    `json:"slider_key"`
	Url             string `json:"url"`
	Order           int    `json:"order"`
}
type MenuDetails struct {
	Device    string `json:"device"`
	Menutype  string `json:"menuType"`
	Title     string `json:"title"`
	Sliderkey int    `json:"sliderKey"`
	Url       string `json:"url"`
	Order     int    `json:"order"`
}

//PlaylistMovie struct for DB binding
type PlaylistMovie struct {
	ID                   int       `json:"id"`
	Title                string    `json:"title"`
	Geoblock             bool      `json:"geoblock"`
	DigitalRightType     int       `json:"digitalRighttype"`
	DigitalRightsRegions []int     `json:"digitalRightsRegions"`
	SubscriptiontPlans   []int     `json:"subscriptiontPlans"`
	InsertedAt           time.Time `json:"insertedAt"`
}
type Slider struct {
	Id                  string    `json:"id" gorm:"primary_key"`
	Name                string    `json:"name"`
	Type                int       `json:"type"`
	CreatedAt           time.Time `json:"created_at"`
	DeletedByUserId     *string   `json:"deleted_by_user_id"`
	IsDisabled          bool      `json:"is_disabled"`
	SchedulingStartDate time.Time `json:"scheduling_start_date"`
	SchedulingEndDate   time.Time `json:"scheduling_end_date"`
	BlackAreaPlaylistId string    `json:"black_area_playlist_id"`
	RedAreaPlaylistId   string    `json:"red_area_playlist_id"`
	GreenAreaPlaylistId string    `json:"green_area_playlist_id"`
	SliderKey           int       `json:"slider_key"`
	ModifiedAt          time.Time `json:"modified_at"`
	PlayWatchTrailer    bool      `json:"play_watch_trailer"`
	MoreButton          bool      `json:"more_button"`
}

//Playlist - binding for db
type Playlist struct {
	ID                  string    `json:"id" gorm:"primary_key"`
	EnglishTitle        string    `json:"english_title"`
	ArabicTitle         string    `json:"arabic_title"`
	SchedulingStartDate time.Time `json:"scheduling_start_date"`
	SchedulingEndDate   time.Time `json:"scheduling_end_date"`
	DeletedByUserId     string    `json:"deleted_by_user_id"`
	IsDisabled          bool      `json:"is_disabled"`
	CreatedAt           time.Time `json:"created_at"`
	PlaylistKey         int       `json:"playlist_key"`
	PodifiedAt          time.Time `json:"modified_at"`
	PlaylistType        string    `json:"playlist_type"`
}
type PlaylistContentIds struct {
	ContentId string `json:"content_id"`
}

//MenuPage struct for DB binding
type MenuPage struct {
	ID             int              `json:"id"`
	FriendlyUrl    string           `json:"friendly_url"`
	SeoDescription string           `json:"seo_description"`
	Title          string           `json:"title"`
	Type           string           `json:"type"`
	Imagery        ImageryDetails   `json:"imagery"`
	Featured       *FeaturedDetails `json:"featured,omitempty"`
	Playlists      []MenuPlaylists  `json:"playlists,omitempty"`
}
type AllAvailableSeasons struct {
	ContentId   string                `json:"-"`
	SeasonId    string                `json:"-"`
	ContentTier int                   `json:"-"`
	Id          int                   `json:"id"`
	VideoId     string                `json:"video_id"`
	FriendlyUrl string                `json:"friendly_url"`
	Title       string                `json:"title"`
	Imagery     ContentImageryDetails `json:"imagery"`
	Geoblock    bool                  `json:"geoblock"`
}
