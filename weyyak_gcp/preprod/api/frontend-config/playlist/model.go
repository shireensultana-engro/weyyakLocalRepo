package playlist

import (
	"time"

	_ "github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Playlists struct {
	ID                      string `json:"id"`
	EnglishTitle            string `json:"englishTitle"`
	ArabicTitle             string `json:"arabicTitle"`
	SchedulingEndDate       string `json:"schedulingEndDate"`
	IsDisabled              bool   `json:"isDisabled"`
	Region                  string `json:"region"`
	IsAssignedToAnyHomePage bool   `json:"isAssignedToAnyHomePage"`
	TheOnlyPlaylistFor      string `json:"theOnlyPlaylistFor"`
	FoundIn                 string `json:"foundIn"`
	PublishingPlatforms     []int  `json:"publishingPlatforms"`
	Platforms               string `json:"platforms,omitempty"`
}

type PlanNames struct {
	TargetPlatform int `json:"target_platform"`
}

type SourceItemTypes struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Playlistupdate struct {
	ID                string    `json:"id"`
	EnglishTitle      string    `json:"englishTitle"`
	ArabicTitle       string    `json:"arabicTitle"`
	SchedulingEndDate string    `json:"schedulingEndDate"`
	IsDisabled        *bool     `json:"isDisabled"`
	DeletedByUserId   string    `json:"deletedByUserId"`
	ModifiedAt        time.Time `json:"modifiedAt"`
}

type PlaylistDisable struct {
	ID         string `json:"id"`
	IsDisabled *bool  `json:"isDisabled"`
}

type PlaylistSummary struct {
	ID                  string `json:"id"`
	EnglishTitle        string `json:"englishTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	IsDisabled          bool   `json:"isDisabled"`
	PublishingPlatforms []int  `json:"publishingPlatforms"`
}

type PlaylistDetails struct {
	ID                     string          `json:"id"`
	EnglishTitle           *string         `json:"englishTitle"`
	ArabicTitle            *string         `json:"arabicTitle"`
	EnglishMetaTitle       *string         `json:"englishMetaTitle"`
	ArabicMetaTitle        *string         `json:"arabicMetaTitle"`
	EnglishMetaDescription *string         `json:"englishMetaDescription"`
	ArabicMetaDescription  *string         `json:"arabicMetaDescription"`
	IsDisabled             bool            `json:"isDisabled"`
	PlaylistType           string          `json:"playlisttype"`
	SchedulingStartDate    *string         `json:"schedulingStartDate"`
	SchedulingEndDate      *string         `json:"schedulingEndDate"`
	PublishingPlatforms    []int           `json:"publishingPlatforms"`
	Regions                []int           `json:"regions"`
	PlaylistItems          []PlaylistItems `json:"playlistItems"`
	Pages                  []Pages         `json:"pages"`
	Platforms              string          `json:"platforms,omitempty"`
	Country                string          `json:"country,omitempty"`
}

type Platforms struct {
	TargetPlatform int `json:"target_platform"`
}

type PlaylistRegion struct {
	CountryId int `json:"countryId"`
}

type Pages struct {
	EnglishTitle        string `json:"englishTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	IsDisabled          bool   `json:"isDisabled"`
	IsHome              bool   `json:"isHome"`
	PublishingPlatforms []int  `json:"publishingPlatforms"`
	ID                  string `json:"id"`
}
type PlaylistItems struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	ID   string `json:"id"`
}

type ContentIds struct {
	ContentId                   string `json:"content_id"`
	SeasonId                    string `json:"season_id"`
	GroupByGenreId              string `json:"group_by_genre_id"`
	GroupBySubgenreId           string `json:"group_by_subgenre_id"`
	GroupByActorId              string `json:"group_by_actor_id"`
	GroupByWriterId             string `json:"group_by_writer_id"`
	GroupByDirectorId           string `json:"group_by_director_id"`
	GroupBySingerId             string `json:"group_by_singer_id"`
	GroupByMusicComposerId      string `json:"group_by_music_composer_id"`
	GroupBySongWriterId         string `json:"group_by_song_writer_id"`
	GroupByProductionYear       string `json:"group_by_production_year"`
	GroupByOriginalLanguageCode string `json:"group_by_original_language_code"`
	GroupByPageId               string `json:"group_by_page_id"`
}

type SourceItemDatas struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	Id   string `json:"id"`
}
type sourceItemTypes struct {
	Name string `json:"name"`
	Id   int    `json:"id" gorm:"primary_key"`
}

/*Get Playlist item based on playlist item type and item type ID*/
type GroupByPageId struct {
	ContentId string `json:"contentId"`
}
type ContentID struct {
	ContentId string `json:"contentId"`
}

/*create update Playlist api */
type PlaylistInputs struct {
	PlaylistId          string            `json:"id"` // id for creating old playlist with .net
	EnglishTitle        string            `json:"englishTitle"`
	ArabicTitle         string            `json:"arabicTitle"`
	PlaylistItems       []PlaylistItemarr `json:"PlaylistItems"`
	SchedulingEndDate   *time.Time        `json:"schedulingEndDate"`
	SchedulingStartDate *time.Time        `json:"schedulingStartDate"`
	PublishingPlatforms []int             `json:"publishingPlatforms"`
	Regions             []int             `json:"regions"`
	PagesIds            []string          `json:"pagesIds"`
	Playlisttype        string            `json:"playlisttype"`
	PlaylistKey         int               `json:"playlistKey"`
}

type Playlist struct {
	Id                  string     `json:"id"`
	EnglishTitle        string     `json:"english_title"`
	ArabicTitle         string     `json:"arabic_title"`
	SchedulingEndDate   *time.Time `json:"scheduling_end_Date"`
	SchedulingStartDate *time.Time `json:"scheduling_start_date"`
	ModifiedAt          time.Time  `json:"modified_at"`
	CreatedAt           time.Time  `json:"created_at"`
	PlaylistKey         int        `json:"playlist_key"`
	PlaylistType        string     `json:"playlist_type"`
}

/*type Pages struct {
	EnglishTitle        string `json:"englishTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	IsDisabled          bool   `json:"isDisabled"`
	IsHome              bool   `json:"isHome"`
	PublishingPlatforms []int  `json:"publishingPlatforms"`
	Id                  string `json:"id"`
}*/

type PlaylistItemarr struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	Id   string `json:"id"`
}

type PlaylistItem struct {
	Id                          string  `json:"id" gorm:"primary_key"`
	PlaylistId                  string  `json:"playlist_id"`
	OneTierContentId            *string `json:"one_tier_content_id"`
	MultiTierContentId          *string `json:"multi_tier_content_id"`
	SeasonId                    *string `json:"season_id" `
	Order                       int     `json:"order"`
	GroupByGenreId              *string `json:"group_by_genre_id" `
	GroupBySubgenreId           *string `json:"group_by_subgenre_id"`
	GroupByActorId              *string `json:"group_by_actor_id" `
	GroupByWriterId             *string `json:"group_by_writer_id" `
	GroupByDirectorId           *string `json:"group_by_director_id" `
	GroupBySingerId             *string `json:"group_by_singer_id" `
	GroupByMusicComposerId      *string `json:"group_by_music_composer_id" `
	GroupBySongWriterId         *string `json:"group_by_song_writer_id"`
	GroupByOriginalLanguageCode *string `json:"group_by_original_language_code" `
	GroupByProductionYear       int     `json:"group_by_production_year" `
	GroupByPageId               *string `json:"group_by_page_id" `
}

type PublishPlatformdata struct {
	PlayListId     string `json:"playlistid"`
	TargetPlatform int    `json:"targetplatform"`
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

// PlaylistItemContent struct
type PlaylistItemContent struct {
	PlaylistItemId string `json:"playlist_item_id"`
	ContentId      *string `json:"content_id"`
	SeasonId       *string `json:"season_id"`
}
