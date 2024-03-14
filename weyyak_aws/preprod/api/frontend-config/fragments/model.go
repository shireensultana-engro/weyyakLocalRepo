package fragments

import (
	"time"

	_ "github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type PageFragment struct {
	PageId    string         `json:"page_id"`
	PageOrder int            `json:"page_order"`
	PageKey   int            `json:"page_key"`
	Country   string         `json:"country"`
	Platform  string         `json:"platform"`
	Details   postgres.Jsonb `json:"details"`
	Language  string         `json:"language"`
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

//MenuPageDetails struct for DB binding
type MenuPageDetails struct {
	ID                       string `json:"id"`
	FriendlyUrl              string `json:"friendly_url"`
	SeoDescription           string `json:"seo_description"`
	Title                    string `json:"title"`
	Type                     int    `json:"type"`
	PageKey                  int    `json:"page_key"`
	PageOrderNumber          int    `json:"page_order_number,omitempty"`
	HasMobileMenu            bool   `json:"has_mobile_menu"`
	HasMenuPosterImage       bool   `json:"has_menu_poster_image"`
	HasMobileMenuPosterImage bool   `json:"has_mobile_menu_poster_image"`
}

//play_list_platform - db binding struct
type PlayListPlatform struct {
	PlayListId     string `json:"playListId"`
	TargetPlatform int    `json:"targetPlatform"`
	PlayListKey    string `json:"playListKey"`
}

type PagePlaylist struct {
	PageId     string `json:"pageId"`
	PlaylistId string `json:"playlistId"`
	Order      int    `json:"order"`
}

type PageOrder struct {
	Order  int    `json:"order"`
	PageId string `json:"pageId"`
}

type PlayListCountry struct {
	PlayListId string `json:"play_list_id"`
	CountryId  int    `json:"country_id"`
}
type PlaylistContentRequest struct {
	Ids      []string `json:"ids"`
	Language string   `json:"language"`
	Country  int      `json:"country"`
}
type PlaylistContentIds struct {
	ContentId string `json:"content_id"`
}

type PlaylistFragment struct {
	PageId     string         `json:"page_id"`
	PageKey    int            `json:"page_key"`
	PlaylistId string         `json:"playlist_id"`
	Details    postgres.Jsonb `json:"details"`
	Language   string         `json:"language"`
	Country    string         `json:"country"`
	Platform   string         `json:"platform"`
}

/* new are */
//MenuPlaylists struct for DB binding
type MenuPlaylists struct {
	ID           int               `json:"id"`
	Title        *string           `json:"title"`
	Content      []PlaylistContent `json:"content"`
	PlaylistType string            `json:"playlisttype"`
	PageContent  []PageContent     `json:"pagecontent"`
}

//PageContent struct for DB binding
type PageContent struct {
	Key            string         `json:"key"`
	ID             int            `json:"id"`
	FriendlyUrl    string         `json:"friendly_url"`
	SeoDescription string         `jsfieldson:"seo_description"`
	Title          string         `json:"title"`
	Type           string         `json:"type"`
	Imagery        ImageryDetails `json:"imagery"`
}

//ImageryDetails struct for DB binding
type ImageryDetails struct {
	MobileMenu            string `json:"mobile-menu"`
	MobilePosterImage     string `json:"menu-poster-image"`
	MobileMenuPosterImage string `json:"mobile-menu-poster-image"`
}

//PlaylistContent struct for DB binding
type PlaylistContent struct {
	ContentId              string                   `json:"content_id"`
	ID                     int                      `json:"id"`
	AgeRating              string                   `json:"age_rating"`
	VideoId                string                   `json:"video_id"`
	FriendlyUrl            string                   `json:"friendly_url"`
	ContentType            string                   `json:"content_type"`
	Synopsis               string                   `json:"synopsis"`
	ProductionYear         *int32                   `json:"production_year"`
	Length                 *int32                   `json:"length"`
	Title                  string                   `json:"title"`
	Cast                   []string                 `json:"cast"`
	SeoDescription         string                   `json:"seo_description"`
	TranslatedTitle        string                   `json:"translated_title"`
	Genres                 []string                 `json:"genres"`
	Tags                   []string                 `json:"tags"`
	SeoTitle               string                   `json:"seo_title"`
	Imagery                interface{}              `json:"imagery"`
	Seasons                []PlaylistContentSeasons `json:"seasons,omitempty"`
	Movies                 []PlaylistMovie          `json:"movies,omitempty"`
	InsertedAt             time.Time                `json:"insertedAt"`
	ModifiedAt             time.Time                `json:"modifiedAt"`
	Geoblock               bool                     `json:"geoblock"`
	MainActor              string                   `json:"main_actor"`
	MainActress            string                   `json:"main_actress"`
	SchedulingDateTime     *time.Time               `json:"_"`
	DigitalRightsStartDate *time.Time               `json:"_"`
	DigitalRightsEndDate   *time.Time               `json:"_"`
}

//PlaylistContentSeasons struct for DB binding
type PlaylistContentSeasons struct {
	ID                   int                    `json:"id"`
	SeasonNumber         int                    `json:"season_number"`
	Dubbed               bool                   `json:"dubbed"`
	SeoDescription       string                 `json:"seo_description"`
	SeoTitle             string                 `json:"seo_title"`
	Title                string                 `json:"title"`
	Geoblock             bool                   `json:"geoblock"`
	DigitalRightType     int                    `json:"digitalRighttype"`
	DigitalRightsRegions []int                  `json:"digitalRightsRegions"`
	SubscriptiontPlans   []int                  `json:"subscriptiontPlans"`
	Imagery              *ContentImageryDetails `json:"imagery,omitempty"`
	Episodes             []SeasonEpisodes       `json:"episodes,omitempty"`
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

//Content Releated
//ContentImageryDetails for DB binding
type ContentImageryDetails struct {
	Thumbnail   string `json:"thumbnail"`
	Backdrop    string `json:"backdrop"`
	MobileImg   string `json:"mobile_img"`
	FeaturedImg string `json:"featured_img,omitempty"`
	Banner      string `json:"banner"`
}

//EpisodeResponse
type SeasonEpisodes struct {
	Id                   int                   `json:"id"`
	SeriesId             int                   `json:"series_id"`
	EpisodeNumber        int                   `json:"episode_number"`
	Synopsis             string                `json:"synopsis"`
	VideoId              string                `json:"video_id"`
	Length               int                   `json:"length"`
	Title                string                `json:"title"`
	Imagery              ContentImageryDetails `json:"imagery"`
	insertedAt           time.Time             `json:"insertedAt"`
	Geoblock             bool                  `json:"geoblock"`
	Tags                 []string              `json:"tags"`
	digitalRightsType    int                   `json:"digitalRighttype"`
	subscriptionPlans    []int                 `json:"subscriptionPlans"`
	digitalRightsRegions string                `json:"digitalRightsRegions"`
	EpisodeId            string                `json:"-"`
	HasPosterImage       bool                  `json:"has_poster_image"`
}

//Content - struct for DB binding
type Content struct {
	//uuid.UUID
	Id                         string    `json:"id" gorm:"primary_key" swaggerignore:"true"`
	AverageRating              float64   `json:"average_rating"`
	AverageRatingUpdatedAt     time.Time `json:"average_rating_updated_at"`
	ContentKey                 int       `json:"content_key"`
	ContentType                string    `json:"content_type"`
	Status                     int       `json:"status"`
	ModifiedAt                 time.Time `json:"modified_at"`
	HasPosterImage             bool      `json:"has_poster_image"`
	HasDetailsBackground       bool      `json:"has_details_background"`
	HasMobileDetailsBackground bool      `json:"has_mobile_details_background"`
	CreatedByUserId            string    `json:"created_by_user_id" gorm:"default:00000000-0000-0000-0000-000000000000;"` //TODO:dependency with token
	ContentTier                int       `json:"content_tier"`
	PrimaryInfoId              string    `json:"primary_info_id"`
	AboutTheContentInfoId      string    `json:"about_the_content_info_id"`
	CastId                     string    `json:"cast_id"`
	MusicId                    string    `json:"music_id"`
	TagInfoId                  string    `json:"tag_info_id"`
	DeletedByUserId            string    `json:"deleted_by_user_id" gorm:"default:00000000-0000-0000-0000-000000000000;"` //TODO:dependency with token
	CreatedAt                  time.Time `json:"created_at"`
	EnglishMetaTitle           string    `json:"english_meta_title"`
	ArabicMetaTitle            string    `json:"arabic_meta_title"`
	EnglishMetaDescription     string    `json:"english_meta_description"`
	ArabicMetaDescription      string    `json:"arabic_meta_description"`
}

//OnetierContentResult -- struct for fetching details
type OnetierContentResult struct {
	ID                         int        `json:"id"`
	AgeRating                  int        `json:"age_rating"`
	VideoId                    string     `json:"video_id"`
	FriendlyUrl                string     `json:"friendly_url"`
	ContentType                string     `json:"content_type"`
	Synopsis                   string     `json:"synopsis"`
	ProductionYear             *int32     `json:"production_year"`
	Length                     *int32     `json:"length"`
	Title                      string     `json:"title"`
	SeoDescription             string     `json:"seo_description"`
	TranslatedTitle            string     `json:"translated_title"`
	SeoTitle                   string     `json:"seo_title"`
	InsertedAt                 time.Time  `json:"insertedAt"`
	ModifiedAt                 time.Time  `json:"modifiedAt"`
	Geoblock                   bool       `json:"geoblock"`
	CastId                     string     `json:"cast_id"`
	DigitalRightsType          int32      `json:"digital_rights_type"`
	RightsCountrys             string     `json:"rights_countrys"`
	RightsPlans                string     `json:"rights_plans"`
	ContentVersionId           string     `json:"content_version_id"`
	SchedulingDateTime         *time.Time `json:"scheduling_date_time"`
	DigitalRightsStartDate     *time.Time `json:"digital_rights_start_date"`
	DigitalRightsEndDate       *time.Time `json:"digital_rights_end_date"`
	HasPosterImage             bool       `json:"has_poster_image"`
	HasDetailsBackground       bool       `json:"has_details_background"`
	HasMobileDetailsBackground bool       `json:"has_mobile_details_background"`
}

//ContentMovieDetails struct for DB binding
type ContentMovieDetails struct {
	Title             string    `json:"title"`
	DigitalRightsType int       `json:"digital_rights_type"`
	InsertedAt        time.Time `json:"insertedAt"`
	Id                string    `json:"id"`
}

type ContentSubsPlans struct {
	SubscriptionPlanId int `json:"subscription_plan_id"`
}

type ActorIds struct {
	MainActorId   string `json:"main_actor_id"`
	MainActressId string `json:"main_actress_id"`
	Actors        string `json:"actors"`
}

//ContentSeasonDetails struct for DB binding
type ContentSeasonDetails struct {
	ID                         string `json:"id"`
	SeasonKey                  int    `json:"season_key"`
	SeasonNumber               int    `json:"season_number,omitempty"`
	LanguageType               int    `json:"language_type,omitempty"`
	SeoDescription             string `json:"seo_description,omitempty"`
	SeoTitle                   string `json:"seo_title,omitempty"`
	Title                      string `json:"title"`
	DigitalRightsType          int    `json:"digital_rights_type"`
	Synopsis                   string `json:"synopsis"`
	HasPosterImage             bool   `json:"has_poster_image"`
	HasDetailsBackground       bool   `json:"has_details_background"`
	HasMobileDetailsBackground bool   `json:"has_mobile_details_background"`
}

//Fragment details models
type SliderFragment struct {
	PageId   string         `json:"page_id"`
	PageKey  int            `json:"page_key"`
	SliderId string         `json:"slider_id"`
	Details  postgres.Jsonb `json:"details"`
	Language string         `json:"language"`
	Country  string         `json:"country"`
	Platform string         `json:"platform"`
}

//FeaturedDetails struct for DB binding
type FeaturedDetails struct {
	ID        int                 `json:"id"`
	Type      string              `json:"type"`
	Playlists []FeaturedPlaylists `json:"playlists"`
}

//FeaturedPlaylists struct for DB binding
type FeaturedPlaylists struct {
	ID           int               `json:"id"`
	PlaylistType string            `json:"playlist_type"`
	Content      []PlaylistContent `json:"content"`
}

type Country struct {
	Id          int    `json:"id"`
	EnglishName string `json:"english_name"`
	ArabicName  string `json:"arabic_name"`
	RegionId    string `json:"region_id"`
	CallingCode string `json:"calling_code"`
	Alpha2code  string `json:"alpha2code"`
}
type PublishPlatform struct {
	Id       int    `json:"id"`
	Platform string `json:"platform"`
}

type FragmentUpdate struct {
	Response string
	Err      error
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

// create  or update page details

type Page struct {
	Id                     string `json:"id,omitempty"`
	EnglishTitle           string `json:"english_title,omitempty"`
	ArabicTitle            string `json:"arabic_title,omitempty"`
	PageOrderNumber        int    `json:"page_order_number,omitempty"`
	PageType               int    `json:"page_type,omitempty"`
	EnglishPageFriendlyUrl string `json:"english_page_friendly_url,omitempty"`
	ArabicPageFriendlyUrl  string `json:"arabic_page_friendly_url,omitempty"`
	EnglishMetaTitle       string `json:"english_meta_title,omitempty"`
	ArabicMetaTitle        string `json:"arabic_meta_title,omitempty"`
	EnglishMetaDescription string `json:"english_meta_description,omitempty"`
	ArabicMetaDescription  string `json:"arabic_meta_description,omitempty"`
	//	DeletedByUserId          string    `json:"deleted_by_user_id,"`
	CreatedAt                time.Time `json:"created_at,omitempty"`
	IsDisabled               bool      `json:"is_disabled,omitempty"`
	PageKey                  int       `json:"page_key,omitempty"`
	ModifiedAt               time.Time `json:"modified_at,omitempty"`
	HasMobileMenu            bool      `json:"has_mobile_menu,omitempty"`
	HasMenuPosterImage       bool      `json:"has_menu_poster_image,omitempty"`
	HasMobileMenuPosterImage bool      `json:"has_mobile_menu_poster_image,omitempty"`
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
}
type PageSlider struct {
	PageId   string `json:"page_id,omitempty"`
	SliderId string `json:"slider_id,omitempty"`
	Order    int    `json:"order,omitempty"`
}
