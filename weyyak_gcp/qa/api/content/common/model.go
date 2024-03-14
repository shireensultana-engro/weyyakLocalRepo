package common

import (
	"time"

	_ "github.com/google/uuid"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// ServerError -- binding struct for error response
type ServerError struct {
	Error       string `json:"error"`
	Description string `json:"description"`
	Code        string `json:"code"`
	RequestId   string `json:"requestId"`
}

/* Input Error Codes*/
type EnglishName struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
type ArabicName struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Invalid struct {
	EnglishName interface{} `json:"englishName,omitempty"`
	ArabicName  interface{} `json:"arabicName,omitempty"`
}
type InpurError struct {
	Error       string  `json:"error"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	RequestId   string  `json:"requestId"`
	Invalid     Invalid `json:"invalid"`
}

// Error codes for one tier
type EnglishTitleError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ArabicTitleError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type PrimaryInfoError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ContentGenresError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ContentVarianceError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type CastError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type MusicError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type TaginfoError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type AbouttheContentError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ContentTypeError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type NonTextualDataError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

type Invalids struct {
	Id                   interface{} `json:"id,omitempty"`
	PrimaryInfoError     interface{} `json:"textualData.primaryInfo,omitempty"`
	ContentTypeError     interface{} `json:"textualData.primaryInfo.contentType,omitempty"`
	ContentGenresError   interface{} `json:"textualData.contentGenres,omitempty"`
	ContentVarianceError interface{} `json:"textualData.contentVariances,omitempty"`
	CastError            interface{} `json:"textualData.cast,omitempty"`
	MusicError           interface{} `json:"textualData.music,omitempty"`
	TaginfoError         interface{} `json:"textualData.tagInfo,omitempty"`
	AbouttheContentError interface{} `json:"textualData.aboutTheContent,omitempty"`
	NonTextualDataError  interface{} `json:"nonTextualData,omitempty"`
	ArabicTitleError     interface{} `json:"arabicTitle,omitempty"`
	EnglishTitleError    interface{} `json:"transliteratedTitle,omitempty"`
}
type FinalErrorResponse struct {
	Error       string   `json:"error"`
	Description string   `json:"description"`
	Code        string   `json:"code"`
	RequestId   string   `json:"requestId"`
	Invalid     Invalids `json:"invalid,omitempty"`
}

type InvalidError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

// Error codes for episode and season
type Contentiderror struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type RigthsError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type SeasonGenresError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type TranslationError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ProductsError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type Invalidsepisode struct {
	PrimaryInfoError    interface{} `json:"primaryInfo,omitempty"`
	CastError           interface{} `json:"cast,omitempty"`
	MusicError          interface{} `json:"music,omitempty"`
	TaginfoError        interface{} `json:"tagInfo,omitempty"`
	NonTextualDataError interface{} `json:"nonTextualData,omitempty"`
	RightsError         interface{} `json:"rights,omitempty"`
	SeasonGenresError   interface{} `json:"seasonGenres,omitempty"`
	ProductsError       interface{} `json:"products,omitempty"`
	TranslationError    interface{} `json:"translation,omitempty"`
	Contentiderror      interface{} `json:"contentId,omitempty"`
	GenresError         interface{} `json:"Genres,omitempty"`
}
type FinalErrorResponseepisode struct {
	Error       string          `json:"error"`
	Description string          `json:"description"`
	Code        string          `json:"code"`
	RequestId   string          `json:"requestId"`
	Invalid     Invalidsepisode `json:"invalid,omitempty"`
}

type GenresError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

/* Redis cache  contenttype */
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
type ContentImageryDetails struct {
	Thumbnail   string `json:"thumbnail"`
	Backdrop    string `json:"backdrop"`
	MobileImg   string `json:"mobile_img"`
	FeaturedImg string `json:"featured_img,omitempty"`
	Banner      string `json:"banner"`
}
type Regions struct {
	Region string `json:"region"`
}
type RedisCacheRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PlaylistContent struct for Redis
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

// PlaylistContentSeasons struct for DB binding
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
	IntroDuration        string                 `json:"introDuration"`
	IntroStart           string                 `json:"introStart"`
	OutroDuration        string                 `json:"outroDuration"`
	OutroStart           string                 `json:"outroStart"`
}

// PlaylistMovie struct for DB binding
type PlaylistMovie struct {
	ID                   int       `json:"id"`
	Title                string    `json:"title"`
	Geoblock             bool      `json:"geoblock"`
	DigitalRightType     int       `json:"digitalRighttype"`
	DigitalRightsRegions []int     `json:"digitalRightsRegions"`
	SubscriptiontPlans   []int     `json:"subscriptiontPlans"`
	InsertedAt           time.Time `json:"insertedAt"`
	IntroDuration        string    `json:"introDuration"`
	IntroStart           string    `json:"introStart"`
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
	EpisodeId            string                `json:"-"`
	HasPosterImage       bool                  `json:"has_poster_image"`
	IntroStart           string                `json:"intro_start"`
	OutroStart           string                `json:"outro_start"`
}
type ContentType struct {
	ContentType string `json:"content_type"`
}
type ContentID struct {
	ContentId string `json:"content_id"`
}
type ContentKey struct {
	ContentKey int `json:"content_key"`
}
type EpisodeDetails struct {
	ContentKey  int    `json:"content_key"`
	ContentType string `json:"content_type"`
	ContentId   string `json:"content_id"`
}
type PlaylistSync struct {
	PlaylistId string `json:"playlistId"`
	DirtyCount int    `json:"dirtyCount"`
}
type SliderSync struct {
	SliderId   string `json:"sliderId"`
	DirtyCount int    `json:"dirtyCount"`
}
