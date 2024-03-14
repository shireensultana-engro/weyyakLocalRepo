package content

import (
	"frontend_service/menu"
	"time"

	_ "github.com/google/uuid"
	_ "github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// SearchDetails
type SearchDetails struct {
	ContentKey          int       `json:"content_key"`
	Id                  string    `json:"id"`
	ContentType         string    `json:"content_type"`
	DigitalRightsType   int       `json:"digital_rights_type"`
	ContentTier         int       `json:"content_tier"`
	CreatedAt           time.Time `json:"created_at"`
	VideoContentId      string    `json:"video_content_id"`
	TransliteratedTitle string    `json:"transliterated_title"`
	ArabicTitle         string    `json:"arabic_title"`
	HasPosterImage      bool      `json:"has_poster_image"`
	VarianceId          string    `json:"variance_id"`
}
type Search struct {
	ID               int                   `json:"id"`
	VideoID          string                `json:"video_id"`
	FriendlyURL      string                `json:"friendly_url"`
	ContentType      string                `json:"content_type"`
	Title            string                `json:"title"`
	Imagery          ContentImageryDetails `json:"imagery"`
	DigitalRighttype int                   `json:"digitalRighttype"`
	Geoblock         bool                  `json:"geoblock"`
}
type MultiTierCheck struct {
	EpisodeCount int    `json:"episode_count"`
	TrailerCount int    `json:"trailer_count"`
	Id           string `json:"id"`
}

// ContentIdDetails -- struct for DB binding
type ContentIdDetails struct {
	Id             string    `json:"id"`
	ContentKey     int       `json:"content_key"`
	ContentTier    int       `json:"content_tier"`
	CreatedAt      time.Time `json:"created_at"`
	PlaybackItemId string    `json:"playback_item_id,omitempty"`
}

// ContentDetails -- struct for DB binding
type ContentDetails struct {
	ContentId          string                `json:"-"`
	ContentTier        int                   `json:"-"`
	SeasonOrVarienceId string                `json:"-"`
	ContentKey         int                   `json:"id"`
	VideoId            string                `json:"video_id"`
	Title              string                `json:"title"`
	TranslatedTitle    string                `json:"translated_title,omitempty"`
	CreatedAt          *time.Time            `json:"createdAt,omitempty"`
	InsertedAt         *time.Time            `json:"insertedAt,omitempty"`
	ModifiedAt         *time.Time            `json:"modifiedAt,omitempty"`
	ContentType        string                `json:"content_type"`
	DigitalRightsType  int                   `json:"digitalRighttype"`
	FriendlyUrl        string                `json:"friendly_url"`
	Imagery            ContentImageryDetails `json:"imagery"`
	Geoblock           bool                  `json:"geoblock"`
}

// MenuPlaylists struct for DB binding
type MenuPlaylists struct {
	ID           int               `json:"id"`
	Title        string            `json:"title"`
	Content      []PlaylistContent `json:"content"`
	PlaylistType string            `json:"playlisttype"`
	PageContent  []PageContent     `json:"pagecontent"`
}

// PageContent struct for DB binding
type PageContent struct {
	Key            string         `json:"key"`
	ID             int            `json:"id"`
	FriendlyUrl    string         `json:"friendly_url"`
	SeoDescription string         `jsfieldson:"seo_description"`
	Title          string         `json:"title"`
	Type           string         `json:"type"`
	Imagery        ImageryDetails `json:"imagery"`
}

// ImageryDetails struct for DB binding
type ImageryDetails struct {
	MobileMenu            string `json:"mobile-menu"`
	MobilePosterImage     string `json:"menu-poster-image"`
	MobileMenuPosterImage string `json:"mobile-menu-poster-image"`
}

// PlaylistContent struct for DB binding
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
type PlaylistContentForWithoutEpisode struct {
	ContentId              string                                    `json:"content_id"`
	ID                     int                                       `json:"id"`
	AgeRating              string                                    `json:"age_rating"`
	VideoId                string                                    `json:"video_id"`
	FriendlyUrl            string                                    `json:"friendly_url"`
	ContentType            string                                    `json:"content_type"`
	Synopsis               string                                    `json:"synopsis"`
	ProductionYear         *int32                                    `json:"production_year"`
	Length                 *int32                                    `json:"length"`
	Title                  string                                    `json:"title"`
	Cast                   []string                                  `json:"cast"`
	SeoDescription         string                                    `json:"seo_description"`
	TranslatedTitle        string                                    `json:"translated_title"`
	Genres                 []string                                  `json:"genres"`
	Tags                   []string                                  `json:"tags"`
	SeoTitle               string                                    `json:"seo_title"`
	Imagery                interface{}                               `json:"imagery"`
	Seasons                []PlaylistContentSeasonsForWithoutEpisode `json:"seasons,omitempty"`
	Movies                 []PlaylistMovie                           `json:"movies,omitempty"`
	InsertedAt             time.Time                                 `json:"insertedAt"`
	ModifiedAt             time.Time                                 `json:"modifiedAt"`
	Geoblock               bool                                      `json:"geoblock"`
	MainActor              string                                    `json:"main_actor"`
	MainActress            string                                    `json:"main_actress"`
	SchedulingDateTime     *time.Time                                `json:"_"`
	DigitalRightsStartDate *time.Time                                `json:"_"`
	DigitalRightsEndDate   *time.Time                                `json:"_"`
}

// Content Releated
// ContentImageryDetails for DB binding
type ContentImageryDetails struct {
	Thumbnail   string `json:"thumbnail"`
	Backdrop    string `json:"backdrop"`
	MobileImg   string `json:"mobile_img"`
	FeaturedImg string `json:"featured_img,omitempty"`
	Banner      string `json:"banner"`
}

// PlaylistContentSeasons struct for DB binding
type PlaylistContentSeasons struct {
	ID                    int                    `json:"id"`
	SeasonNumber          int                    `json:"season_number"`
	Dubbed                bool                   `json:"dubbed"`
	SeoDescription        string                 `json:"seo_description"`
	SeoTitle              string                 `json:"seo_title"`
	Title                 string                 `json:"title"`
	Geoblock              bool                   `json:"geoblock"`
	DigitalRightType      int                    `json:"digitalRighttype"`
	DigitalRightsRegions  []int                  `json:"digitalRightsRegions"`
	SubscriptiontPlans    []int                  `json:"subscriptiontPlans"`
	SubscriptionPlansName []string               `json:"subscriptionPlansName"`
	Imagery               *ContentImageryDetails `json:"imagery,omitempty"`
	Episodes              []SeasonEpisodes       `json:"episodes,omitempty"`
	IntroDuration         string                 `json:"introDuration"`
	IntroStart            string                 `json:"introStart"`
	OutroDuration         string                 `json:"outroDuration"`
	OutroStart            string                 `json:"outroStart"`
}
type PlaylistContentSeasonsForWithoutEpisode struct {
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
	Episodes             []SeasonEpisodes       `json:"episodes"`
	IntroDuration        string                 `json:"introDuration"`
	IntroStart           string                 `json:"introStart"`
	OutroDuration        string                 `json:"outroDuration"`
	OutroStart           string                 `json:"outroStart"`
}

// PlaylistMovie struct for DB binding
type PlaylistMovie struct {
	ID                    int       `json:"id"`
	Title                 string    `json:"title"`
	Geoblock              bool      `json:"geoblock"`
	DigitalRightType      int       `json:"digitalRighttype"`
	DigitalRightsRegions  []int     `json:"digitalRightsRegions"`
	SubscriptiontPlans    []int     `json:"subscriptiontPlans"`
	SubscriptionPlansName []string  `json:"subscriptionPlansName"`
	InsertedAt            time.Time `json:"insertedAt"`
	IntroDuration         string    `json:"introDuration"`
	IntroStart            string    `json:"introStart"`
}

// ContentSeasonDetails struct for DB binding
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
	IntroDuration              string `json:"intro_duration"`
	IntroStart                 string `json:"intro_start"`
	OutroDuration              string `json:"outro_duration"`
	OutroStart                 string `json:"outro_start"`
}

// ContentMovieDetails struct for DB binding
type ContentMovieDetails struct {
	Title             string    `json:"title"`
	DigitalRightsType int       `json:"digital_rights_type"`
	InsertedAt        time.Time `json:"insertedAt"`
	Id                string    `json:"id"`
	IntroDuration     string    `json:"intro_duration"`
	IntroStart        string    `json:"intro_start"`
}

// OnetierContentResult -- struct for fetching details
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
type ActorIds struct {
	MainActorId   string `json:"main_actor_id"`
	MainActressId string `json:"main_actress_id"`
	Actors        string `json:"actors"`
}

// Content - struct for DB binding
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
type PaginationResult struct {
	Size   int64 `json:"size"`
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}
type ContentFragment struct {
	ContentId              string         `json:"content_id"`
	ContentTier            int            `json:"content_tier"`
	ContentKey             int            `json:"content_key"`
	Language               string         `json:"language"`
	SchedulingDateTime     *time.Time     `json:"scheduling_date_time"`
	DigitalRightsStartDate *time.Time     `json:"digital_rights_start_date"`
	DigitalRightsEndDate   *time.Time     `json:"digital_rights_end_date"`
	Details                postgres.Jsonb `json:"details"`
	Country                string         `json:"country"`
}
type ContentFragmentDetails struct {
	ContentId              string    `json:"content_id"`
	ContentTier            int       `json:"content_tier"`
	ContentKey             int       `json:"content_key"`
	Language               string    `json:"language"`
	SchedulingDateTime     time.Time `json:"scheduling_date_time"`
	DigitalRightsStartDate time.Time `json:"digital_rights_start_date"`
	DigitalRightsEndDate   time.Time `json:"digital_rights_end_date"`
	Details                string    `json:"details"`
}
type ContentSubsPlans struct {
	SubscriptionPlanId int `json:"subscription_plan_id"`
}
type ContentRatingDetails struct {
	Content  ContentRating `json:"content"`
	UserData *UserRating   `json:"userData,omitempty"`
}
type ResumbleContentRatingDetails struct {
	Content  ContentRating       `json:"content"`
	UserData *ResumbleUserRating `json:"userData,omitempty"`
}
type ContentRatingQuery struct {
	AverageRating       float64 `json:"average_rating"`
	TransliteratedTitle string  `json:"transliterated_title"`
	ContentType         string  `json:"content_type"`
	Length              int     `json:"length"`
	Id                  string  `json:"id"`
	ContentKey          int     `json:"content_key"`
	ContentTier         int     `json:"content_tier"`
	DigitalRightsType   int     `json:"digital_rights_type"`
}
type ContentRatingQueryEpisode struct {
	AverageRating       float64 `json:"average_rating"`
	TransliteratedTitle string  `json:"transliterated_title"`
	ContentType         string  `json:"content_type"`
	Length              int     `json:"length"`
	Id                  string  `json:"id"`
	ContentKey          int     `json:"content_key"`
	ContentTier         int     `json:"content_tier"`
	DigitalRightsType   int     `json:"digital_rights_type"`
	PlaybackItemId      string  `json:"playback_item_id"`
}
type ContentRating struct {
	AverageRating     float64   `json:"averageRating"`
	Title             string    `json:"title"`
	ContentType       string    `json:"contentType"`
	DigitalRightsType int       `json:"digitalRightsType"`
	Duration          int       `json:"duration"`
	Genres            []string  `json:"genres"`
	Id                int       `json:"id"`
	ContentId         string    `json:"-"`
	ViewedAt          time.Time `json:"-"`
	LastWatchPosition int       `json:"-"`
}
type UserRating struct {
	Rating            *float64            `json:"rating"`
	RatedAt           *time.Time          `json:"ratedAt"`
	IsTailored        bool                `json:"isTailored"`
	AddedToPlaylistAt *time.Time          `json:"addedToPlaylistAt"`
	ViewActivity      *RatingViewActivity `json:"viewActivity"`
}
type ResumbleUserRating struct {
	Rating            *float64                   `json:"rating"`
	RatedAt           *time.Time                 `json:"ratedAt"`
	IsTailored        bool                       `json:"isTailored"`
	AddedToPlaylistAt *time.Time                 `json:"addedToPlaylistAt"`
	ViewActivity      ResumbleRatingViewActivity `json:"viewActivity"`
}

// GenreNameResponse
type GenreNameResponse struct {
	Genre string `json:"genre"`
}

// EpisodeResponse
type EpisodeResponse struct {
	Id                int                   `json:"id"`
	SeriesId          int                   `json:"series_id"`
	SeasonNumber      int                   `json:"season_number"`
	SeasonId          int                   `json:"season_id"`
	EpisodeNumber     int                   `json:"episode_number"`
	VideoId           string                `json:"video_id"`
	Length            int                   `json:"length"`
	Genres            []string              `json:"genres"`
	Title             string                `json:"title"`
	TranslatedTitle   string                `json:"translated_title"`
	Imagery           ContentImageryDetails `json:"imagery"`
	CreatedAt         time.Time             `json:"insertedAt"`
	Geoblock          bool                  `json:"geoblock"`
	Tags              []string              `json:"tags"`
	DigitalRightsType int                   `json:"digitalRighttype"`
	SubscriptionPlans []string              `json:"subscriptiontPlans"`
	IntroStart        string                `json:"introStart"`
	OutroStart        string                `json:"outroStart"`
	Contentid         string                `json:"-"`
	Seasonid          string                `json:"-"`
	Episodeid         string                `json:"-"`
	Tagid             string                `json:"-"`
	Rightsid          string                `json:"-"`
}

// SubscriptionPlan
type SubscriptionPlan struct {
	Name string `json:"name"`
}

// TagResponse
type TagResponse struct {
	TagInfoId string `json:"tag_info_id"`
	Tags      string `json:"tags"`
}

// ImageKey
type ImageKey struct {
	Thumbnail   string `json:"thumbnail"`
	Backdrop    string `json:"backdrop"`
	MobileImg   string `json:"mobile_img"`
	FeaturedImg string `json:"featured_img"`
	Banner      string `json:"banner"`
}

// EpisodeResponse
type SeasonEpisodes struct {
	Id                   int                   `json:"id"`
	SeriesId             int                   `json:"series_id"`
	EpisodeNumber        int                   `json:"episode_number"`
	Synopsis             string                `json:"synopsis"`
	VideoId              string                `json:"video_id"`
	Length               int                   `json:"length"`
	Title                string                `json:"title"`
	Imagery              ContentImageryDetails `json:"imagery"`
	InsertedAt           time.Time             `json:"insertedAt"`
	Geoblock             bool                  `json:"geoblock"`
	Tags                 []string              `json:"tags"`
	DigitalRightsType    int                   `json:"digitalRighttype"`
	SubscriptionPlans    []int                 `json:"subscriptionPlans"`
	DigitalRightsRegions string                `json:"digitalRightsRegions"`
	LastWatchPosition    int                   `json:"lastWatchPosition"`
	PlaybackItemId       string                `json:"playbackItemId"`
	EpisodeId            string                `json:"-"`
	HasPosterImage       bool                  `json:"has_poster_image"`
	IntroStart           string                `json:"introStart"`
	OutroStart           string                `json:"outroStart"`
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
type TrailerImagery struct {
	TrailerPosterImage string `json:"trailerPosterImage"`
}
type ContentTrailers struct {
	Id          string         `json:"-"`
	VarianceId  string         `json:"-"`
	ContentId   string         `json:"-"`
	ContentTier int            `json:"-"`
	Title       string         `json:"title"`
	VideoId     string         `json:"video_id"`
	Length      int            `json:"length"`
	Imagery     TrailerImagery `json:"imagery"`
	Geoblock    bool           `json:"geoblock"`
}
type MediaObjectDetails struct {
	ContentTier     int                   `json:"-"`
	ContentId       string                `json:"-"`
	VarianceId      string                `json:"-"`
	EpisodeId       string                `json:"-"`
	Id              int                   `json:"id"`
	SeriesId        int                   `json:"series_id,omitempty"`
	SeasonNumber    int                   `json:"season_number,omitempty"`
	AgeRating       string                `json:"age_rating"`
	SeasonId        int                   `json:"season_id,omitempty"`
	EpisodeNumber   int                   `json:"episode_number,omitempty"`
	FriendlyUrl     string                `json:"friendly_url,omitempty"`
	ContentType     string                `json:"content_type"`
	BaseContentType string                `json:"base_content_type"`
	VideoId         string                `json:"video_id"`
	Length          int                   `json:"length"`
	SeriesTitle     string                `json:"series_title,omitempty"`
	Title           string                `json:"title"`
	TranslatedTitle string                `json:"translated_title,omitempty"`
	Imagery         ContentImageryDetails `json:"imagery"`
	InsertedAt      time.Time             `json:"insertedAt"`
	Geoblock        bool                  `json:"geoblock"`
}
type RelatedContentGenres struct {
	OriginalLanguage string `json:"original_language"`
	GenreId          string `json:"genre_id"`
	SubgenreId       string `json:"subgenre_id"`
}
type ViewActivityDetails struct {
	LastWatchPosition int        `json:"last_watch_position"`
	ViewedAt          *time.Time `json:"viewed_at"`
	EpisodeKey        int        `json:"episode_key"`
	ContentType       string     `json:"content_type"`
	Duration          int        `json:"duration"`
}

type RatedContent struct {
	Id        string    `json:"id" grom:"primary_key"`
	RatedAt   time.Time `json:"rated_at"`
	Rating    float64   `json:"rating"`
	ContentId string    `json:"content_id"`
	UserId    string    `json:"user_id"`
	DeviceId  string    `json:"device_id"`
	IsHidden  bool      `json:"is_hidden"`
}
type ContentKeys struct {
	ContentKey string `json:"content_key"`
	Id         string `json:"id"`
}
type PlaylistedContent struct {
	UserId    string    `json:"user_id"`
	ContentId string    `json:"content_id"`
	AddedAt   time.Time `json:"added_at"`
}

// PlaylistedContentRequest -- Request for playlisted content
type PlaylistedContentRequest struct {
	Id          int           `json:"id"`
	Title       string        `json:"title"`
	ContentType string        `json:"contentType"`
	Genres      []interface{} `json:"genres"`
}
type RatingViewActivity struct {
	ViewedAt            *time.Time `json:"viewedAt"`
	ResumeWatchPosition *int       `json:"resumeWatchPosition"`
}
type ResumbleRatingViewActivity struct {
	ViewedAt            time.Time `json:"viewedAt"`
	ResumeWatchPosition int       `json:"resumeWatchPosition"`
}

type AddViewActivityRequest struct {
	Content           ViewActivityRequestContent `json:"content"`
	LastWatchPosition int                        `json:"lastWatchPosition"`
	WatchSessionId    string                     `json:"watchSessionId"`
}
type ViewActivityRequestContent struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	ContentType string   `json:"contentType"`
	Duration    int      `json:"duration"`
	Genres      []string `json:"genres"`
}
type ViewActivity struct {
	Id                string    `json:"id" gorm:"primary_id"`
	ContentId         string    `json:"content_id"`
	UserId            string    `json:"user_id"`
	DeviceId          string    `json:"device_id"`
	LastWatchPosition int       `json:"last_watch_position"`
	ViewedAt          time.Time `json:"viewed_at"`
	WatchSessionId    string    `json:"watch_session_id"`
	IsHidden          string    `json:"is_hidden"`
	PlaybackItemId    string    `json:"playback_item_id"`
}
type ViewActivityHistory struct {
	Id                string    `json:"id,omitempty" gorm:"primary_id"`
	ContentTypeName   string    `json:"content_type_name"`
	ContentKey        int       `json:"content_key"`
	UserId            string    `json:"user_id"`
	DeviceId          string    `json:"device_id"`
	LastWatchPosition int       `json:"last_watch_position"`
	ViewedAt          time.Time `json:"viewed_at"`
	WatchSessionId    string    `json:"watch_session_id"`
}

type SearchContent struct {
	Id              string `json:"id"`
	contentType     string `json:"contentType"`
	ContentKey      string `json:"content_key"`
	CreatedByUserId string `json:"created_by_user_id"`
	ContentTier     int    `json:"content_tier"`
}

// add rating for content
type Addcontent struct {
	ContentRequest ContentRequest `json:"content"`
	Rating         int            `json:"rating"`
}
type ContentRequest struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	ContentType string   `json:"contentType"`
	Duration    int      `json:"duration"`
	Genres      []string `json:"genres"`
}
type ContentId struct {
	Id string `json:"id"`
}
type CreateRatedContent struct {
	Rating    int       `json:"rating"`
	RatedAt   time.Time `json:"rated_at"`
	ContentId string    `json:"content_id"`
	UserID    string    `json:"user_id"`
	IsHidden  bool      `json:"is_hidden"`
	DeviceId  string    `json:"device_id"`
}
type AddRatingToHistory struct {
	Rating    int       `json:"rating"`
	RatedAt   time.Time `json:"rated_at"`
	ContentId string    `json:"content_id"`
	UserID    string    `json:"user_id"`
	DeviceId  string    `json:"device_id"`
}
type ResponseContent struct {
	AverageRating          float64 `json:"average_rating"`
	AverageRatingUpdatedAt time.Time
}
type RatedContentByUser struct {
	Id        string
	ContentId string
}

// End of Rated_content
type ViewActivityDetailRequest struct {
	Description     string `json:"description" gorm:"not null;"`
	IsCommunication *bool  `json:"isWithCommunication" gorm:"not null;"`
	IsSound         *bool  `json:"isWithSound" gorm:"not null;`
	IsTranslation   *bool  `json:"isWithTranslation" gorm:"not null;`
	IsVideo         *bool  `json:"isWithVideo" gorm:"not null;`
}
type ViewActivityDetailss struct {
	Id              string `json:"id"`
	Description     string `json:"description" gorm:"not null;"`
	IsCommunication *bool  `json:"isWithCommunication" gorm:"not null;"`
	IsSound         *bool  `json:"isWithSound" gorm:"not null;`
	IsTranslation   *bool  `json:"isWithTranslation" gorm:"not null;`
	IsVideo         *bool  `json:"isWithVideo" gorm:"not null;`
	ReportedAt      time.Time
	ViewActivityId  string `json:"viewactivityid`
}
type ViewDetails struct {
	Id string `json:"id"`
}
type WatchDetails struct {
	ViewActivityId string `json:"id"`
}

type GenreContent struct {
	GenreId    string `json:"genre_id"`
	ContentId  string `json:"content_id"`
	ContentKey string `json:"content_key"`
}

type CastContent struct {
	GenreId    string `json:"genre_id"`
	ContentId  string `json:"content_id"`
	ContentKey string `json:"content_key"`
}
type UserEvent struct {
	Timestamp time.Time `json:"@timestamp"`
	UserID    string    `json:"userid"`
	EventType string    `json:"event"`
	Details   string    `json:"details"`
}

// Get Episode Details Based on SeasonId
type UserInfo struct {
	Email string `json:"email"`
}
type EpisodeDetailsSummary struct {
	IsPrimary              bool   `json:"isPrimary"`
	UserId                 string `json:"userId"`
	SecondarySeasonId      string `json:"secondarySeasonId" `
	VarianceIds            []int  `json:"varianceIds"`
	EpisodeIds             []int  `json:"episodeIds"`
	SecondaryEpisodeId     string `json:"secondaryEpisodeId"`
	ContentId              string `json:"contentId"`
	EpisodeKey             int    `json:"episodeKey"`
	SeasonId               string `json:"seasonId"`
	Status                 int    `json:"status"`
	StatusCanBeChanged     bool   `json:"statusCanBeChanged"`
	SubStatus              int    `json:"subStatus"`
	SubStatusName          string `json:"subStatusName"`
	DigitalRightsType      int    `json:"digitalRightsType"`
	DigitalRightsStartDate string `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   string `json:"digitalRightsEndDate"`
	CreatedBy              string `json:"createdBy"`
	PrimaryInfo            PrimaryInfo
	Cast                   []int `json:"cast"`
	Music                  []int `json:"music"`
	TagInfo                []int `json:"tagInfo"`
	NonTextualData         []int `json:"nonTextualData"`
	Translation            Translation
	SchedulingDateTime     []int  `json:"schedulingDateTime"`
	PublishingPlatforms    []int  `json:"publishingPlatforms"`
	SeoDetails             []int  `json:"seoDetails"`
	Id                     string `json:"id"`
}
type PrimaryInfo struct {
	Number              int     `json:"number ,omitempty"`
	VideoContentId      string  `json:"videoContentId ,omitempty"`
	SynopsisEnglish     string  `json:"synopsisEnglish ,omitempty"`
	SynopsisArabic      string  `json:"synopsisArabic ,omitempty"`
	SeasonNumber        int     `json:"seasonNumber ,omitempty"`
	OriginalTitle       string  `json:"originalTitle"`
	AlternativeTitle    string  `json:"alternativeTitle"`
	ArabicTitle         string  `json:"arabicTitle"`
	TransliteratedTitle string  `json:"transliteratedTitle"`
	Notes               string  `json:"notes"`
	IntroStart          *string `json:"introStart"`
	OutroStart          *string `json:"outroStart"`
}
type Translation struct {
	LanguageType       string  `json:"languageType"`
	DubbingLanguage    *string `json:"dubbingLanguage"`
	DubbingDialectId   *int    `json:"dubbingDialectId"`
	SubtitlingLanguage *string `json:"subtitlingLanguage"`
}

type FinalEpisodesResult struct {
	IsPrimary              bool    `json:"isPrimary"`
	UserId                 string  `json:"userId" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	SecondarySeasonId      string  `json:"secondarySeasonId" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	VarianceIds            []int   `json:"varianceIds"`
	EpisodeIds             []int   `json:"episodeIds"`
	SecondaryEpisodeId     string  `json:"secondaryEpisodeId" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	ContentId              string  `json:"contentId"`
	EpisodeKey             int     `json:"episodeKey"`
	SeasonId               string  `json:"seasonId"`
	Status                 int     `json:"status"`
	StatusCanBeChanged     bool    `json:"statusCanBeChanged"`
	SubStatus              int     `json:"subStatus"`
	SubStatusName          string  `json:"subStatusName"`
	DigitalRightsType      int     `json:"digitalRightsType"`
	DigitalRightsStartDate string  `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   string  `json:"digitalRightsEndDate"`
	CreatedByUserId        string  `json:"createdByUserId"`
	Number                 int     `json:"number"`
	VideoContentId         string  `json:"videoContentId"`
	SynopsisEnglish        string  `json:"synopsisEnglish"`
	SynopsisArabic         string  `json:"synopsisArabic"`
	OriginalTitle          string  `json:"originalTitle"`
	AlternativeTitle       string  `json:"alternativeTitle"`
	ArabicTitle            string  `json:"arabicTitle"`
	TransliteratedTitle    string  `json:"transliteratedTitle"`
	Notes                  string  `json:"notes"`
	IntroStart             *string `json:"introStart"`
	OutroStart             *string `json:"outroStart"`
	Cast                   []int   `json:"cast"`
	Music                  []int   `json:"music"`
	TagInfo                []int   `json:"tagInfo"`
	NonTextualData         []int   `json:"nonTextualData"`
	LanguageType           int     `json:"languageType"`
	DubbingLanguage        *string `json:"dubbingLanguage"`
	DubbingDialectId       *int    `json:"dubbingDialectId"`
	SubtitlingLanguage     *string `json:"subtitlingLanguage"`
	SchedulingDateTime     []int   `json:"schedulingDateTime"`
	PublishingPlatforms    []int   `json:"publishingPlatforms"`
	SeoDetails             []int   `json:"seoDetails"`
	Id                     string  `json:"id"`
}

// Get Episode Details Based in contentId
type SeasonDetailsSummary struct {
	ContentId          string    `json:"contentId"`
	SeasonKey          int       `json:"seasonKey"`
	Status             int       `json:"status"`
	StatusCanBeChanged bool      `json:"statusCanBeChanged"`
	SubStatusName      string    `json:"subStatusName"`
	ModifiedAt         time.Time `json:"modifiedAt"`
	PrimaryInfo        PrimaryInfo
	Cast               *string `json:"cast"`
	Music              *string `json:"music"`
	TagInfo            *string `json:"tagInfo"`
	SeasonGenres       *string `json:"seasonGenres"`
	AboutTheContent    *string `json:"aboutTheContent"`
	Translation        Translation
	Episodes           *string `json:"episodes"`
	NonTextualData     *string `json:"nonTextualData"`
	Rights             Rights  `json:"rights"`
	CreatedBy          string  `json:"createdBy"`
	IntroDuration      string  `json:"introDuration"`
	IntroStart         *string `json:"introStart"`
	OutroDuration      string  `json:"outroDuration"`
	OutroStart         *string `json:"outroStart"`
	Products           *string `json:"products"`
	SeoDetails         *string `json:"seoDetails"`
	VarianceTrailers   *string `json:"varianceTrailers"`
	Id                 string  `json:"id"`
}
type Rights struct {
	DigitalRightsType      int    `json:"digitalRightsType"`
	DigitalRightsStartDate string `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   string `json:"digitalRightsEndDate"`
	DigitalRightsRegions   []int  `json:"digitalRightsRegions"`
	SubscriptionPlans      []int  `json:"subscriptionPlans"`
}
type DigitalRightsRegions struct {
	CountryId int `json:"country_Id"`
}
type FinalSeasonResult struct {
	ContentId              string    `json:"contentId"`
	SeasonKey              int       `json:"seasonKey"`
	Status                 int       `json:"status"`
	StatusCanBeChanged     bool      `json:"statusCanBeChanged"`
	SubStatusName          string    `json:"subStatusName"`
	ModifiedAt             time.Time `json:"modifiedAt"`
	SeasonNumber           int       `json:"seasonNumber"`
	OriginalTitle          string    `json:"originalTitle"`
	AlternativeTitle       string    `json:"alternativeTitle"`
	ArabicTitle            string    `json:"arabicTitle"`
	TransliteratedTitle    string    `json:"transliteratedTitle"`
	Notes                  string    `json:"notes"`
	IntroStart             *string   `json:"introStart"`
	OutroStart             *string   `json:"outroStart"`
	Cast                   *string   `json:"cast"`
	Music                  *string   `json:"music"`
	TagInfo                *string   `json:"tagInfo"`
	SeasonGenres           *string   `json:"seasonGenres"`
	AboutTheContent        *string   `json:"aboutTheContent"`
	LanguageType           int       `json:"languageType"`
	DubbingLanguage        *string   `json:"dubbingLanguage"`
	DubbingDialectId       *int      `json:"dubbingDialectId"`
	SubtitlingLanguage     *string   `json:"subtitlingLanguage"`
	Episodes               *string   `json:"episodes"`
	NonTextualData         *string   `json:"nonTextualData"`
	DigitalRightsType      int       `json:"digitalRightsType"`
	DigitalRightsStartDate string    `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   string    `json:"digitalRightsEndDate"`
	// DigitalRightsRegions   string    `json:"DigitalRightsRegions"`
	SubscriptionPlans int     `json:"subscriptionPlans"`
	CreatedBy         string  `json:"createdBy"`
	IntroDuration     string  `json:"introDuration"`
	AboutIntroStart   *string `json:"aboutintroStart"`
	OutroDuration     string  `json:"outroDuration"`
	AboutOutroStart   *string `json:"aboutoutroStart"`
	Products          *string `json:"products"`
	SeoDetails        *string `json:"seoDetails"`
	VarianceTrailers  *string `json:"varianceTrailers"`
	Id                string  `json:"id"`
	RightsId          string  `json:"rightsId"`
}
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
type PlaylistDetailsResponse struct {
	Total       int                  `json:"total"`
	PerPage     int                  `json:"per_page"`
	CurrentPage int                  `json:"current_page"`
	LastPage    int                  `json:"last_page"`
	NextPageUrl *string              `json:"next_page_url"`
	PrevPageUrl *string              `json:"prev_page_url"`
	From        int                  `json:"from"`
	To          int                  `json:"to"`
	Data        MenuPlaylistsDetails `json:"data"`
}

// MenuPlaylists struct for DB binding
type MenuPlaylistsDetails struct {
	ID           int                `json:"id"`
	Title        string             `json:"title"`
	Content      postgres.Jsonb     `json:"content"`
	PlaylistType string             `json:"playlisttype"`
	PageContent  []menu.PageContent `json:"pagecontent"`
	// []PlaylistContent
}

type Cont struct {
	Cnt postgres.Jsonb
}

type PlaylistPage struct {
	ID   string `json:"id"`
	Key  int    `json:"key"`
	Name string `json:"name"`
}

type ContentRights struct {
	Countries string `json:"countries"`
	Platform  string `json:"platform"`
}

type ContentImagesFlutter struct {
	ImageCategory string   `json:"imageCategory"`
	ImageUrl      []string `json:"imageUrl"`
}

type ContentRatingFlutter struct {
	BaseContentType    string                 `json:"baseContentType"`
	ContentType        string                 `json:"contentType"`
	ID                 string                 `json:"id"`
	Images             []ContentImagesFlutter `json:"images"`
	Key                int                    `json:"key"`
	MultiTierContentId string                 `json:"multi_tier_content_id"`
	Name               string                 `json:"name"`
	OneTierContentId   string                 `json:"one_tier_content_id"`
	Isaddtoplaylist    bool                   `json:"isaddtoplaylist"`
}

type ResumbleContentRatingDetailsForFlutter struct {
	ID              string                 `json:"id"`
	Content         []ContentRatingFlutter `json:"content"`
	CreatedAt       time.Time              `json:"createdAt"`
	DeletedByUserId string                 `json:"deletedByUserId" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	EndDate         *time.Time             `json:"endDate"`
	Key             int                    `json:"key"`
	Language        string                 `json:"language"`
	ModifiedAt      *time.Time             `json:"modifiedAt"`
	Name            string                 `json:"name"`
	Pages           []PlaylistPage         `json:"pages"`
	Rights          []ContentRights        `json:"rights"`
	StartDate       *time.Time             `json:"startDate"`
	Status          string                 `json:"status"`
	Title           string                 `json:"title"`
	Type            string                 `json:"type"`
}

type ContentRatingQueryFlutter struct {
	AverageRating      float64   `json:"averageRating"`
	Title              string    `json:"title"`
	ContentType        string    `json:"contentType"`
	DigitalRightsType  int       `json:"digitalRightsType"`
	Duration           int       `json:"duration"`
	Genres             []string  `json:"genres"`
	Id                 int       `json:"id"`
	OneTierContentId   string    `json:"one_tier_content_id"`
	MultiTierContentId string    `json:"multi_tier_content_id"`
	ContentId          string    `json:"-"`
	ViewedAt           time.Time `json:"-"`
	LastWatchPosition  int       `json:"-"`
}
