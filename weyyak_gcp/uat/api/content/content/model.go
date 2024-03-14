package content

import (
	"time"

	_ "github.com/google/uuid"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

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

// ContentPrimaryInfo - db binding struct
type ContentPrimaryInfo struct {
	ID                  string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	OriginalTitle       string `json:"original_title"`
	AlternativeTitle    string `json:"alternative_title"`
	ArabicTitle         string `json:"arabic_title"`
	TransliteratedTitle string `json:"transliterated_title"`
	Notes               string `json:"notes"`
	IntroStart          string `json:"intro_start,omitempty"`
	OutroStart          string `json:"outro_start,omitempty"`
}

// ContentVariance - db binding struct
type ContentVariance struct {
	ID                    string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	ContentId             string `json:"content_id"`
	PlaybackItemId        string `json:"playback_item_id"`
	Order                 int    `json:"order"`
	HasOverlayPosterImage bool   `json:"has_overlay_poster_image"`
	HasDubbingScript      bool   `json:"has_dubbing_script"`
	HasSubtitlingScript   bool   `json:"has_subtitling_script"`
	DeletedByUserId       string `json:"deleted_by_user_id" gorm:"default:00000000-0000-0000-0000-000000000000;"` //TODO:dependency with token
	Status                int    `json:"status"`
	HasAllRights          bool   `json:"has_all_rights"`
	IntroStart            int    `json:"intro_start,omitempty"`
	IntroDuration         int    `json:"intro_duration,omitempty"`
}

// PlaybackItem - db binding struct
type PlaybackItem struct {
	Id             string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	VideoContentId string `json:"video_content_id"`
	// YoutubeVideoId     string    `json:"youtube_video_id"`
	SchedulingDateTime time.Time `json:"scheduling_date_time"`
	CreatedByUserId    string    `json:"created_by_user_id" gorm:"default:00000000-0000-0000-0000-000000000000;"` //TODO:dependency with token
	TranslationId      string    `json:"translation_id"`
	RightsId           string    `json:"rights_id"`
	Duration           int       `json:"duration"`
}

// ContentRights - db binding struct
type ContentRights struct {
	Id                     string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	DigitalRightsType      int    `json:"digital_rights_type"`
	DigitalRightsStartDate string `json:"digital_rights_start_date"`
	DigitalRightsEndDate   string `json:"digital_rights_end_date"`
}

// ContentRightsCountry - db binding struct
type ContentRightsCountry struct {
	Id              string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	ContentRightsId string `json:"content_rights_id"`
	CountryId       int    `json:"country_id"`
}

// ContentRightsPlan - db binding struct
type ContentRightsPlan struct {
	RightsId           string `json:"rights_id" gorm:"primary_key" swaggerignore:"true"`
	SubscriptionPlanId int    `json:"subscription_plan_id"`
}

// RightsProduct - db binding struct
type RightsProduct struct {
	RightsId    string `json:"rights_id"`
	ProductName int    `json:"product_name"`
}

// PlaybackItemTargetPlatform - db binding struct
type PlaybackItemTargetPlatform struct {
	PlaybackItemId string `json:"playback_item_id"`
	TargetPlatform int    `json:"target_platform"`
	RightsId       string `json:"rights_id"`
}

// AboutTheContentInfo - db binding struct
type AboutTheContentInfo struct {
	Id                    string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	OriginalLanguage      string `json:"originalLanguage" binding:"required"`
	Supplier              string `json:"supplier" binding:"required"`
	AcquisitionDepartment string `json:"acquisitionDepartment"`
	EnglishSynopsis       string `json:"englishSynopsis" binding:"required"`
	ArabicSynopsis        string `json:"arabicSynopsis" binding:"required"`
	ProductionYear        string `json:"productionYear"`
	ProductionHouse       string `json:"productionHouse"`
	AgeGroup              int    `json:"ageGroup" binding:"required"`
	IntroDuration         string `json:"introDuration,omitempty"`
	IntroStart            string `json:"introStart,omitempty"`
	OutroDuration         string `json:"outroDuration,omitempty"`
	OutroStart            string `json:"outroStart,omitempty"`
}

type AboutTheContentInfoUpdate struct {
	Id                    string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	OriginalLanguage      string `json:"originalLanguage" binding:"required"`
	Supplier              string `json:"supplier" binding:"required"`
	AcquisitionDepartment string `json:"acquisitionDepartment"`
	EnglishSynopsis       string `json:"englishSynopsis" binding:"required"`
	ArabicSynopsis        string `json:"arabicSynopsis" binding:"required"`
	ProductionYear        *int   `json:"productionYear"`
	ProductionHouse       string `json:"productionHouse"`
	AgeGroup              int    `json:"ageGroup" binding:"required"`
	IntroDuration         string `json:"introDuration,omitempty"`
	IntroStart            string `json:"introStart,omitempty"`
	OutroDuration         string `json:"outroDuration,omitempty"`
	OutroStart            string `json:"outroStart,omitempty"`
}

// ProductionCountry - db binding struct
type ProductionCountry struct {
	AboutTheContentInfoId string `json:"about_the_content_info"`
	CountryId             int    `json:"country_id"`
}

// Onetier content request
type OnetierContentRequest struct {
	TextualData    TextualData    `json:"textualData" binding:"required"`
	NonTextualData NonTextualData `json:"nonTextualData" binding:"required"`
}

// TextualData content request
type TextualData struct {
	PrimaryInfo      PrimaryInfo        `json:"primaryInfo"  binding:"required"`
	SeoDetails       SeoDetails         `json:"seoDetails" binding:"required"`
	ContentGenres    []ContentGenres    `json:"contentGenres" binding:"required"`
	ContentVariances []ContentVariances `json:"contentVariances" binding:"required"`
	AboutTheContent  AboutTheContent    `json:"aboutTheContent" binding:"required"`
	Cast             Cast               `json:"cast" binding:"required"`
	Music            Music              `json:"music" binding:"required"`
	TagInfo          TagInfo            `json:"tagInfo" binding:"required"`
}
type PrimaryInfo struct {
	SeasonNumber        int     `json:"seasonNumber,omitempty"`
	ContentType         string  `json:"contentType,omitempty"`
	Number              int     `json:"number,omitempty"`
	VideoContentId      string  `json:"videoContentId,omitempty"`
	SynopsisEnglish     string  `json:"synopsisEnglish,omitempty"`
	SynopsisArabic      string  `json:"synopsisArabic,omitempty"`
	OriginalTitle       string  `json:"originalTitle"`
	AlternativeTitle    string  `json:"alternativeTitle"`
	ArabicTitle         string  `json:"arabicTitle"`
	TransliteratedTitle string  `json:"transliteratedTitle"`
	Notes               string  `json:"notes"`
	IntroStart          *string `json:"introStart,omitempty"`
	OutroStart          *string `json:"outroStart,omitempty"`
}
type SeoDetails struct {
	EnglishMetaTitle       string `json:"englishMetaTitle"`
	ArabicMetaTitle        string `json:"arabicMetaTitle"`
	EnglishMetaDescription string `json:"englishMetaDescription"`
	ArabicMetaDescription  string `json:"arabicMetaDescription"`
}
type ContentGenres struct {
	GenreId     string   `json:"genreId"`
	SubgenresId []string `json:"subgenresId"`
	Id          string   `json:"id,omitempty"`
}
type ContentVariances struct {
	ID                 string `json:"id"`
	Status             int    `json:"status"`
	StatusCanBeChanged bool   `json:"statusCanBeChanged"`
	SubStatusName      string `json:"subStatusName"`
	VideoContentId     string `json:"videoContentId"`
	// YoutubeVideoId         string   `json:"youtubeVideoId"`
	LanguageType           string   `json:"languageType"`
	OverlayPosterImage     string   `json:"overlayPosterImage"`
	DubbingScript          string   `json:"dubbingScript"`
	SubtitlingScript       string   `json:"subtitlingScript"`
	DubbingLanguage        string   `json:"dubbingLanguage"`
	DubbingDialectId       int      `json:"dubbingDialectId"`
	SubtitlingLanguage     string   `json:"subtitlingLanguage"`
	DigitalRightsType      int      `json:"digitalRightsType"`
	DigitalRightsStartDate string   `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   string   `json:"digitalRightsEndDate"`
	DigitalRightsRegions   []int    `json:"digitalRightsRegions"`
	SchedulingDateTime     string   `json:"schedulingDateTime"`
	CreatedBy              string   `json:"createdBy"`
	PublishingPlatforms    []int    `json:"publishingPlatforms"`
	Products               []int    `json:"products"`
	SubscriptionPlans      []int    `json:"subscriptionPlans"`
	CountryCheck           bool     `json:"countryCheck"`
	VarianceTrailers       []string `json:"varianceTrailers"`
}
type AboutTheContent struct {
	OriginalLanguage      string  `json:"originalLanguage"`
	Supplier              string  `json:"supplier"`
	AcquisitionDepartment string  `json:"acquisitionDepartment"`
	EnglishSynopsis       string  `json:"englishSynopsis"`
	ArabicSynopsis        string  `json:"arabicSynopsis"`
	ProductionYear        string  `json:"productionYear"`
	ProductionHouse       string  `json:"productionHouse"`
	AgeGroup              int     `json:"ageGroup"`
	IntroDuration         *string `json:"introDuration,omitempty"`
	IntroStart            *string `json:"introStart,omitempty"`
	OutroDuration         *string `json:"outroDuration,omitempty"`
	OutroStart            *string `json:"outroStart,omitempty"`
	ProductionCountries   []int   `json:"productionCountries"`
}
type Cast struct {
	MainActorId   *string  `json:"mainActorId"`
	MainActressId *string  `json:"mainActressId"`
	Actors        []string `json:"actors"`
	Writers       []string `json:"writers"`
	Directors     []string `json:"directors"`
}
type Music struct {
	Singers        []string `json:"singers"`
	MusicComposers []string `json:"musicComposers"`
	SongWriters    []string `json:"songWriters"`
}
type TagInfo struct {
	Tags []string `json:"tags"`
}
type ContentTranslationRequest struct {
	Id                 string  `json:"id" gorm:"primary_key"`
	LanguageType       int     `json:"language_type"`
	DubbingLanguage    *string `json:"dubbing_language"`
	DubbingDialectId   *int    `json:"dubbing_dialect_id"`
	SubtitlingLanguage *string `json:"subtitling_language"`
}
type NonTextualData struct {
	PosterImage             string `json:"posterImage,omitempty"`
	OverlayPosterImage      string `json:"overlayPosterImage,omitempty"`
	DetailsBackground       string `json:"detailsBackground,omitempty"`
	MobileDetailsBackground string `json:"mobileDetailsBackground,omitempty"`
	SeasonLogo              string `json:"seasonLogo,omitempty"`
	DubbingScript           string `json:"dubbingScript,omitempty"`
	SubtitlingScript        string `json:"subtitlingScript,omitempty"`
}

// CurlCall response
type CurlResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"Status"`
	InsertedID string `json:"insert_id"`
}
type ContentTranslation struct {
	LanguageType       string  `json:"languageType"`
	DubbingLanguage    *string `json:"dubbingLanguage"`
	DubbingDialectId   *int    `json:"dubbingDialectId"`
	SubtitlingLanguage *string `json:"subtitlingLanguage"`
}
type Rights struct {
	DigitalRightsType      int       `json:"digitalRightsType"`
	DigitalRightsStartDate time.Time `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   time.Time `json:"digitalRightsEndDate"`
	DigitalRightsRegions   []int     `json:"digitalRightsRegions"`
	SubscriptionPlans      []int     `json:"subscriptionPlans"`
}
type SeasonEpisodes struct {
	IsPrimary              bool               `json:"isPrimary"`
	UserId                 string             `json:"userId"`
	SecondarySeasonId      string             `json:"secondarySeasonId"`
	VarianceIds            *string            `json:"varianceIds"`
	EpisodeIds             *string            `json:"episodeIds"`
	SecondaryEpisodeId     string             `json:"secondaryEpisodeId"`
	ContentId              string             `json:"contentId"`
	EpisodeKey             int                `json:"episodeKey"`
	SeasonId               string             `json:"seasonId"`
	Status                 int                `json:"status"`
	StatusCanBeChanged     bool               `json:"statusCanBeChanged"`
	SubStatus              int                `json:"subStatus"`
	SubStatusName          *string            `json:"subStatusName"`
	DigitalRightsType      *int               `json:"digitalRightsType"`
	DigitalRightsStartDate *time.Time         `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   *time.Time         `json:"digitalRightsEndDate"`
	CreatedBy              *time.Time         `json:"createdBy"`
	PrimaryInfo            PrimaryInfo        `json:"primaryInfo"`
	Cast                   Cast               `json:"cast"`
	Music                  Music              `json:"music"`
	TagInfo                TagInfo            `json:"tagInfo"`
	NonTextualData         *NonTextualData    `json:"nonTextualData"`
	Translation            ContentTranslation `json:"translation"`
	SchedulingDateTime     *time.Time         `json:"schedulingDateTime"`
	PublishingPlatforms    []int              `json:"publishingPlatforms"`
	SeoDetails             *SeoDetails        `json:"seoDetails"`
	Id                     string             `json:"id"`
}
type SeasonEpisode struct {
	IsPrimary              bool               `json:"isPrimary"`
	UserId                 string             `json:"userId"`
	SecondarySeasonId      string             `json:"secondarySeasonId"`
	VarianceIds            *string            `json:"varianceIds"`
	EpisodeIds             *string            `json:"episodeIds"`
	SecondaryEpisodeId     string             `json:"secondaryEpisodeId"`
	ContentId              string             `json:"contentId"`
	EpisodeKey             int                `json:"episodeKey"`
	SeasonId               string             `json:"seasonId"`
	Status                 int                `json:"status"`
	StatusCanBeChanged     bool               `json:"statusCanBeChanged"`
	SubStatus              int                `json:"subStatus"`
	SubStatusName          *string            `json:"subStatusName"`
	DigitalRightsType      *int               `json:"digitalRightsType"`
	DigitalRightsStartDate *time.Time         `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   *time.Time         `json:"digitalRightsEndDate"`
	CreatedBy              *time.Time         `json:"createdBy"`
	PrimaryInfo            PrimaryInfo        `json:"primaryInfo"`
	Cast                   Cast               `json:"cast"`
	Music                  Music              `json:"music"`
	TagInfo                TagInfo            `json:"tagInfo"`
	NonTextualData         NonTextualData     `json:"nonTextualData"`
	Translation            ContentTranslation `json:"translation"`
	SchedulingDateTime     *time.Time         `json:"schedulingDateTime"`
	PublishingPlatforms    []int              `json:"publishingPlatforms"`
	SeoDetails             SeoDetails         `json:"seoDetails"`
	Id                     string             `json:"id"`
	SubscriptionPlanId       *int  `json:"subscriptionPlans"`
}
type ContentSeasons struct {
	ContentId          string             `json:"contentId"`
	SeasonKey          int                `json:"seasonKey"`
	Status             int                `json:"status"`
	StatusCanBeChanged bool               `json:"statusCanBeChanged"`
	SubStatusName      *string            `json:"subStatusName"`
	ModifiedAt         time.Time          `json:"modifiedAt"`
	PrimaryInfo        PrimaryInfo        `json:"primaryInfo"`
	Cast               Cast               `json:"cast"`
	Music              Music              `json:"music"`
	TagInfo            TagInfo            `json:"tagInfo"`
	SeasonGenres       []ContentGenres    `json:"seasonGenres"`
	AboutTheContent    AboutTheContent    `json:"aboutTheContent"`
	Translation        ContentTranslation `json:"translation"`
	Episodes           []SeasonEpisodes   `json:"episodes"`
	NonTextualData     *NonTextualData    `json:"nonTextualData"`
	Rights             Rights             `json:"rights"`
	CreatedBy          *string            `json:"createdBy"`
	IntroDuration      string             `json:"introDuration,omitempty"`
	IntroStart         string             `json:"introStart,omitempty"`
	OutroDuration      string             `json:"outroDuration,omitempty"`
	OutroStart         string             `json:"outroStart,omitempty"`
	Products           []int              `json:"products"`
	SeoDetails         *SeoDetails        `json:"seoDetails"`
	VarianceTrailers   []interface{}      `json:"varianceTrailers"`
	Id                 string             `json:"id"`
}
type MultitierContentDetails struct {
	ContentKey     int              `json:"contentKey"`
	Duration       *int             `json:"duration"`
	Status         int              `json:"status"`
	PrimaryInfo    PrimaryInfo      `json:"primaryInfo"`
	ContentGenres  []ContentGenres  `json:"contentGenres"`
	ContentSeasons []ContentSeasons `json:"contentSeasons"`
	SeoDetails     *SeoDetails      `json:"seoDetails"`
	Id             string           `json:"id"`
}
type MultitierContentQueryDetails struct {
	ContentKey             int     `json:"content_key"`
	Duration               *int    `json:"duration"`
	Status                 int     `json:"status"`
	ContentType            string  `json:"content_type"`
	OriginalTitle          string  `json:"original_title"`
	AlternativeTitle       string  `json:"alternative_title"`
	ArabicTitle            string  `json:"arabic_title"`
	TransliteratedTitle    string  `json:"transliterated_title"`
	Notes                  string  `json:"notes"`
	IntroStart             *string `json:"intro_start,omitempty"`
	OutroStart             *string `json:"outro_start,omitempty"`
	EnglishMetaTitle       string  `json:"english_meta_title"`
	ArabicMetaTitle        string  `json:"arabic_meta_title"`
	EnglishMetaDescription string  `json:"english_meta_descriptipon"`
	ArabicMetaDescription  string  `json:"arabic_meta_descriptipon"`
	Id                     string  `json:"id"`
}
type ContentGeneresQueryDetails struct {
	GenreId       string `json:"genre_id"`
	SubgenresId   string `json:"subgenres_id"`
	SubGenreOrder string `json:"sub_genre_order"`
	Id            string `json:"id"`
}
type ContentSeasonsQueryDetails struct {
	ContentId              string    `json:"contentId"`
	SeasonKey              int       `json:"seasonKey"`
	Status                 int       `json:"status"`
	ModifiedAt             time.Time `json:"modifiedAt"`
	Number                 int       `json:"number"`
	OriginalTitle          string    `json:"original_title"`
	AlternativeTitle       string    `json:"alternative_title"`
	ArabicTitle            string    `json:"arabic_title"`
	TransliteratedTitle    string    `json:"transliterated_title"`
	Notes                  string    `json:"notes"`
	IntroStart             *string   `json:"intro_start,omitempty"`
	OutroStart             *string   `json:"outro_start"`
	OriginalLanguage       string    `json:"original_language"`
	Supplier               string    `json:"supplier"`
	AcquisitionDepartment  string    `json:"acquisition_department"`
	EnglishSynopsis        string    `json:"english_synopsis"`
	ArabicSynopsis         string    `json:"arabic_synopsis"`
	ProductionYear         string    `json:"production_year"`
	ProductionHouse        string    `json:"production_house"`
	AgeGroup               int       `json:"age_group"`
	IntroDuration          *string   `json:"intro_duration,omitempty"`
	AtciIntroStart         *string   `json:"atci_intro_start,omitempty"`
	OutroDuration          *string   `json:"outro_duration,omitempty"`
	AtciOutroStart         *string   `json:"atci_outro_start,omitempty"`
	ProductionCountries    string    `json:"production_countries"`
	LanguageType           int       `json:"language_type"`
	DubbingLanguage        *string   `json:"dubbing_language"`
	DubbingDialectId       *int      `json:"dubbing_dialect_id"`
	SubtitlingLanguage     *string   `json:"subtitling_language"`
	DigitalRightsType      int       `json:"digital_rights_type"`
	DigitalRightsStartDate time.Time `json:"digital_rights_start_date"`
	DigitalRightsEndDate   time.Time `json:"digital_rights_end_date"`
	DigitalRightsRegions   string    `json:"digital_rights_regions"`
	SubscriptionPlans      string    `json:"subscription_plans"`
	CreatedByUserId        *string   `json:"created_by_user_id"`
	EnglishMetaTitle       string    `json:"english_meta_title"`
	ArabicMetaTitle        string    `json:"arabic_meta_title"`
	EnglishMetaDescription string    `json:"english_meta_descriptipon"`
	ArabicMetaDescription  string    `json:"arabic_meta_descriptipon"`
	Id                     string    `json:"id"`
	RightsId               string    `json:"rights_id"`
	DeletedByUserId        string    `json:"deleted_by_user_id"`
	AboutTheContentInfoId  string    `json:"about_the_content_info_id"`
}
type CastQueryDetails struct {
	MainActorId   *string `json:"main_actor_id"`
	MainActressId *string `json:"main_actress_id"`
	Actors        string  `json:"actors"`
	Writers       string  `json:"writers"`
	Directors     string  `json:"directors"`
}
type MusicQueryDetails struct {
	Singers        string `json:"singers"`
	MusicComposers string `json:"music_composers"`
	SongWriters    string `json:"song_writers"`
}
type TagInfoQueryDetails struct {
	Tags string `json:"tags"`
}
type SeasonEpisodesQueryDetails struct {
	EpisodeKey             int        `json:"episode_key"`
	SeasonId               string     `json:"season_id"`
	Status                 int        `json:"status"`
	DigitalRightsType      *int       `json:"digital_rights_type"`
	DigitalRightsStartDate *time.Time `json:"digital_rights_start_date"`
	DigitalRightsEndDate   *time.Time `json:"digital_rights_end_date"`
	CreatedAt              *time.Time `json:"created_at"`
	Number                 int        `json:"number"`
	VideoContentId         string     `json:"video_content_id"`
	EnglishSynopsis        string     `json:"english_synopsis"`
	ArabicSynopsis         string     `json:"arabic_synopsis"`
	OriginalTitle          string     `json:"original_title"`
	AlternativeTitle       string     `json:"alternative_title"`
	ArabicTitle            string     `json:"arabic_title"`
	TransliteratedTitle    string     `json:"transliterated_title"`
	Notes                  string     `json:"notes"`
	IntroStart             *string    `json:"intro_start,omitempty"`
	OutroStart             *string    `json:"outro_start,omitempty"`
	LanguageType           int        `json:"language_type"`
	DubbingLanguage        *string    `json:"dubbing_language"`
	DubbingDialectId       *int       `json:"dubbing_dialect_id"`
	SubtitlingLanguage     *string    `json:"subtitling_language"`
	SchedulingDateTime     *time.Time `json:"scheduling_date_time"`
	PublishingPlatforms    string     `json:"publishing_platforms"`
	EnglishMetaTitle       string     `json:"english_meta_title"`
	ArabicMetaTitle        string     `json:"arabic_meta_title"`
	EnglishMetaDescription string     `json:"english_meta_descriptipon"`
	ArabicMetaDescription  string     `json:"arabic_meta_descriptipon"`
	Id                     string     `json:"id"`
	MainActorId            *string    `json:"main_actor_id"`
	MainActressId          *string    `json:"main_actress_id"`
	Actors                 string     `json:"actors"`
	Writers                string     `json:"writers"`
	Directors              string     `json:"directors"`
	Singers                string     `json:"singers"`
	MusicComposers         string     `json:"music_composers"`
	SongWriters            string     `json:"song_writers"`
	Tags                   string     `json:"tags"`
	ContentId              string     `json:"contentId"`
	HasPosterImage         bool       `json:"hasPosterImage"`
	HasDubbingScript       bool       `json:"hasDubbingScript"`
	HasSubtitlingScript    bool       `json:"HasSubtitlingScript"`
	SubscriptionPlanId      *int        `json:"subscriptionPlans"`
}
type PrimaryInfoRequest struct {
	OriginalTitle       string `json:"originalTitle"`
	AlternativeTitle    string `json:"alternativeTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	TransliteratedTitle string `json:"transliteratedTitle"`
	Notes               string `json:"notes"`
	SeasonNumber        int    `json:"seasonNumber"`
}

// AboutTheContentInfoRequest - db binding struct
type AboutTheContentInfoRequest struct {
	OriginalLanguage      string `json:"originalLanguage" binding:"required"`
	Supplier              string `json:"supplier" binding:"required"`
	AcquisitionDepartment string `json:"acquisitionDepartment"`
	EnglishSynopsis       string `json:"englishSynopsis" binding:"required"`
	ArabicSynopsis        string `json:"arabicSynopsis" binding:"required"`
	ProductionYear        string `json:"productionYear"`
	ProductionHouse       string `json:"productionHouse"`
	AgeGroup              int    `json:"ageGroup" binding:"required"`
	IntroDuration         string `json:"introDuration,omitempty"`
	IntroStart            string `json:"introStart,omitempty"`
	OutroDuration         string `json:"outroDuration" binding:"required"`
	OutroStart            string `json:"outroStart" binding:"required"`
	ProductionCountries   []int  `json:"productionCountries"`
}
type RightsRequest struct {
	DigitalRightsType       int    `json:"digitalRightsType"`
	DigitalRightsStartDate  string `json:"digitalRightsStartDate"`
	DigitalRightsEndDate    string `json:"digitalRightsendDate"`
	DigitalRightsRegionsint []int  `json:"digitalRightsRegions"`
	SubscriptionPlans       []int  `json:"subscriptionPlans"`
}
type RightsRequestvariance struct {
	DigitalRightsType          int      `json:"digitalRightsType"`
	DigitalRightsStartDate     string   `json:"digitalRightsStartDate"`
	DigitalRightsEndDate       string   `json:"digitalRightsendDate"`
	DigitalRightsRegionsstring []string `json:"digitalRightsRegions"`
	SubscriptionPlans          []int    `json:"subscriptionPlans"`
}
type CreateSeasonRequest struct {
	SeasonId         string                      `json:"Seasonid"`
	ContentId        *string                     `json:"contentId"`
	PrimaryInfo      *PrimaryInfoRequest         `json:"primaryInfo"`
	SeasonGenres     []ContentGenres             `json:"contentGenres"`
	VarianceTrailers []interface{}               `json:"varianceTrailers"`
	AboutTheContent  *AboutTheContentInfoRequest `json:"aboutTheContent"`
	IntroStart       string                      `json:"introStart,omitempty"`
	Translation      *ContentTranslation         `json:"translation"`
	NonTextualData   *NonTextualData             `json:"nonTextualData"`
	Cast             *Cast                       `json:"cast"`
	Music            *Music                      `json:"music"`
	TagInfo          *TagInfo                    `json:"tagInfo"`
	Rights           *RightsRequest              `json:"rights"`
	CountryCheck     bool                        `json:"countryCheck"`
	Products         *[]int                      `json:"products"`
	SeoDetails       SeoDetails                  `json:"seoDetails"`
}
type CreateSeasonRequestvariance struct {
	SeasonId         string                      `json:"Seasonid"`
	ContentId        *string                     `json:"contentId"`
	PrimaryInfo      *PrimaryInfoRequest         `json:"primaryInfo"`
	SeasonGenres     []ContentGenres             `json:"contentGenres"`
	VarianceTrailers []interface{}               `json:"varianceTrailers"`
	AboutTheContent  *AboutTheContentInfoRequest `json:"aboutTheContent"`
	IntroStart       string                      `json:"introStart,omitempty"`
	Translation      *ContentTranslation         `json:"translation"`
	NonTextualData   *NonTextualData             `json:"nonTextualData"`
	Cast             *Cast                       `json:"cast"`
	Music            *Music                      `json:"music"`
	TagInfo          *TagInfo                    `json:"tagInfo"`
	Rights           *RightsRequestvariance      `json:"rights"`
	CountryCheck     bool                        `json:"countryCheck"`
	Products         *[]int                      `json:"products"`
	SeoDetails       SeoDetails                  `json:"seoDetails"`
}

// Season - struct for DB binding
type Season struct {
	Id                         string    `json:"id" gorm:"primary_key" swaggerignore:"true"`
	ContentId                  string    `json:"content_id"`
	SeasonKey                  int       `json:"content_key"`
	Status                     int       `json:"status"`
	ModifiedAt                 time.Time `json:"modified_at"`
	HasPosterImage             string    `json:"has_poster_image"`
	HasOverlayPosterImage      string    `json:"has_overlay_poster_image"`
	HasDetailsBackground       string    `json:"has_details_background"`
	HasMobileDetailsBackground string    `json:"has_mobile_details_background"`
	HasSeasonLogo              string    `json:"has_season_logo"`
	CreatedByUserId            string    `json:"created_by_user_id" gorm:"default:00000000-0000-0000-0000-000000000000;"` //TODO:dependency with token
	PrimaryInfoId              string    `json:"primary_info_id"`
	AboutTheContentInfoId      string    `json:"about_the_content_info_id"`
	Number                     int       `json:"number"`
	TranslationId              string    `json:"translation_id"`
	CastId                     string    `json:"cast_id"`
	MusicId                    string    `json:"music_id"`
	TagInfoId                  string    `json:"tag_info_id"`
	RightsId                   string    `json:"rights_id"`
	//	DeletedByUserId            string    `json:"deleted_by_user_id" gorm:"default:00000000-0000-0000-0000-000000000000;"` //TODO:dependency with token
	CreatedAt              time.Time `json:"created_at"`
	EnglishMetaTitle       string    `json:"english_meta_title"`
	ArabicMetaTitle        string    `json:"arabic_meta_title"`
	EnglishMetaDescription string    `json:"english_meta_description"`
	ArabicMetaDescription  string    `json:"arabic_meta_description"`
	HasAllRights           bool      `json:"has_all_rights"`
	ThirdPartySeasonKey    int       `json:"third_party_season_key"`
}
type ContentCast struct {
	Id            string  `json:"id" gorm:"primary_key" swaggerignore:"true"`
	MainActorId   *string `json:"main_actor_id"`
	MainActressId *string `json:"main_actress_id"`
}
type ContentActor struct {
	CastId  string `json:"cast_id"`
	ActorId string `json:"actor_id"`
}
type ContentWriter struct {
	CastId   string `json:"cast_id"`
	WriterId string `json:"writer_id"`
}
type ContentDirector struct {
	CastId     string `json:"cast_id"`
	DirectorId string `json:"director_id"`
}
type ContentSinger struct {
	MusicId  string `json:"music_id"`
	SingerId string `json:"singer_id"`
}
type ContentMusic struct {
	Id string `json:"id"`
}
type ContentMusicComposer struct {
	MusicId         string `json:"music_id"`
	MusicComposerId string `json:"music_composer_id"`
}
type ContentSongWriter struct {
	MusicId      string `json:"music_id"`
	SongWriterId string `json:"song_writer_id"`
}
type ContentTagInfo struct {
	Id string `json:"id"`
}
type ContentTag struct {
	TagInfoId        string `json:"tag_info_id"`
	TextualDataTagId string `json:"textual_data_tag_id"`
}
type TextualDataTag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type SeasonGenre struct {
	Id       string `json:"id" gorm:"primary_key"`
	GenreId  string `json:"genre_id"`
	SeasonId string `json:"season_id"`
	Order    int    `json:"order"`
}
type SeasonSubgenre struct {
	SeasonGenreId string `json:"season_genre_id"`
	SubgenreId    string `json:"subgenre_id"`
	Order         int    `json:"order"`
}
type ContentVariancesDetails struct {
	Status                 int        `json:"status"`
	StatusCanBeChanged     bool       `json:"statusCanBeChanged"`
	SubStatusName          string     `json:"subStatusName"`
	VideoContentId         *string    `json:"videoContentId"`
	LanguageType           *int       `json:"languageType"`
	OverlayPosterImage     *string    `json:"overlayPosterImage"`
	DubbingScript          *string    `json:"dubbingScript"`
	SubtitlingScript       *string    `json:"subtitlingScript"`
	DubbingLanguage        *string    `json:"dubbingLanguage"`
	DubbingDialectId       *int       `json:"dubbingDialectId"`
	SubtitlingLanguage     *string    `json:"subtitlingLanguage"`
	DigitalRightsType      int        `json:"digitalRightsType"`
	DigitalRightsStartDate *time.Time `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   *time.Time `json:"digitalRightsEndDate"`
	DigitalRights          string     `json:"-"`
	DigitalRightsRegions   []int      `json:"digitalRightsRegions"`
	SchedulingDateTime     *time.Time `json:"schedulingDateTime"`
	CreatedBy              string     `json:"createdBy"`
	PublishingPlatforms    *int       `json:"publishingPlatforms"`
	Products               *int       `json:"products"`
	SubscriptionPlans      *int       `json:"subscriptionPlans"`
	CountryCheck           bool       `json:"countryCheck"`
	IntroDuration          *string    `json:"introDuration,omitempty"`
	IntroStart             *string    `json:"introStart,omitempty"`
	VarianceTrailers       *string    `json:"varianceTrailers"`
	Id                     string     `json:"id"`
}
type AboutTheContentDetails struct {
	OriginalLanguage      string  `json:"originalLanguage"`
	Supplier              *string `json:"supplier"`
	AcquisitionDepartment *string `json:"acquisitionDepartment"`
	EnglishSynopsis       *string `json:"englishSynopsis"`
	ArabicSynopsis        *string `json:"arabicSynopsis"`
	ProductionYear        *int    `json:"productionYear"`
	ProductionHouse       *string `json:"productionHouse"`
	AgeGroup              *int    `json:"ageGroup"`
	IntroDuration         *string `json:"introDuration,omitempty"`
	IntroStart            *string `json:"introStart,omitempty"`
	OutroDuration         *string `json:"outroDuration"`
	OutroStart            *string `json:"outroStart"`
	ProductionCountries   *int    `json:"productionCountries"`
}
type PrimaryInfoDetails struct {
	SeasonNumber        int     `json:"seasonNumber,omitempty"`
	Number              int     `json:"number,omitempty"`
	VideoContentId      *string `json:"videoContentId"`
	SynopsisEnglish     *string `json:"synopsisEnglish"`
	SynopsisArabic      *string `json:"synopsisArabic"`
	OriginalTitle       *string `json:"originalTitle"`
	AlternativeTitle    *string `json:"alternativeTitle"`
	ArabicTitle         *string `json:"arabicTitle"`
	TransliteratedTitle string  `json:"transliteratedTitle"`
	Notes               *string `json:"notes"`
	IntroStart          *string `json:"introStart,omitempty"`
	OutroStart          *string `json:"outroStart,omitempty"`
}
type ContentTranslationDetails struct {
	LanguageType       string  `json:"languageType"`
	DubbingLanguage    *string `json:"dubbingLanguage"`
	DubbingDialectId   *int    `json:"dubbingDialectId"`
	SubtitlingLanguage *string `json:"subtitlingLanguage"`
}
type SeasonEpisodesDetails struct {
	IsPrimary              bool               `json:"isPrimary"`
	UserId                 string             `json:"userId"`
	SecondarySeasonId      string             `json:"secondarySeasonId"`
	VarianceIds            *string            `json:"varianceIds"`
	EpisodeIds             *string            `json:"episodeIds"`
	SecondaryEpisodeId     string             `json:"secondaryEpisodeId"`
	ContentId              string             `json:"contentId"`
	EpisodeKey             int                `json:"episodeKey"`
	Number                 int                `json:"-"`
	TransliteratedTitle    string             `json:"-"`
	SeasonId               string             `json:"seasonId"`
	Status                 int                `json:"status"`
	StatusCanBeChanged     bool               `json:"statusCanBeChanged"`
	SubStatus              int                `json:"subStatus"`
	SubStatusName          *string            `json:"subStatusName"`
	DigitalRightsType      *int               `json:"digitalRightsType"`
	DigitalRightsStartDate *time.Time         `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   *time.Time         `json:"digitalRightsEndDate"`
	CreatedBy              *string            `json:"createdBy"`
	PrimaryInfo            PrimaryInfoDetails `json:"primaryInfo"`
	Cast                   *int               `json:"cast"`
	Music                  *int               `json:"music"`
	TagInfo                *int               `json:"tagInfo"`
	NonTextualData         *string            `json:"nonTextualData"`
	Translation            *int               `json:"translation"`
	SchedulingDateTime     *time.Time         `json:"schedulingDateTime"`
	PublishingPlatforms    *int               `json:"publishingPlatforms"`
	SeoDetails             *SeoDetails        `json:"seoDetails"`
	Id                     string             `json:"id"`
}
type ContentSeasonsDetails struct {
	ContentId              string                    `json:"contentId"`
	SeasonKey              int                       `json:"seasonKey"`
	Status                 int                       `json:"status"`
	StatusCanBeChanged     bool                      `json:"statusCanBeChanged"`
	SubStatusName          *string                   `json:"subStatusName"`
	ModifiedAt             time.Time                 `json:"modifiedAt"`
	SeasonNumber           int                       `json:"-"`
	TransliteratedTitle    string                    `json:"-"`
	OriginalLanguage       string                    `json:"-"`
	LanguageType           int                       `json:"-"`
	DubbingLanguage        string                    `json:"-"`
	DubbingDialectId       int                       `json:"-"`
	DigitalRightsType      int                       `json:"-"`
	DigitalRightsStartDate time.Time                 `json:"-"`
	DigitalRightsEndDate   time.Time                 `json:"-"`
	DigitalRights          string                    `json:"-"`
	PrimaryInfo            PrimaryInfoDetails        `json:"primaryInfo"`
	Cast                   *string                   `json:"cast"`
	Music                  *string                   `json:"music"`
	TagInfo                *string                   `json:"tagInfo"`
	SeasonGenres           *string                   `json:"seasonGenres"`
	AboutTheContent        AboutTheContentDetails    `json:"aboutTheContent"`
	Translation            ContentTranslationDetails `json:"translation"`
	Episodes               []SeasonEpisodesDetails   `json:"episodes"`
	NonTextualData         *string                   `json:"nonTextualData"`
	Rights                 Rights                    `json:"rights"`
	CreatedBy              *string                   `json:"createdBy"`
	IntroDuration          string                    `json:"introDuration,omitempty"`
	IntroStart             string                    `json:"introStart,omitempty"`
	OutroDuration          string                    `json:"outroDuration,omitempty"`
	OutroStart             string                    `json:"outroStart,omitempty"`
	Products               *int                      `json:"products"`
	SeoDetails             *string                   `json:"seoDetails"`
	VarianceTrailers       *string                   `json:"varianceTrailers"`
	Id                     string                    `json:"id"`
}
type GetAllContentDetails struct {
	Type                int                       `json:"type"`
	ContentKey          int                       `json:"contentKey"`
	Status              int                       `json:"status"`
	StatusCanBeChanged  bool                      `json:"statusCanBeChanged"`
	SubStatusName       string                    `json:"subStatusName"`
	TransliteratedTitle string                    `json:"transliteratedTitle"`
	CreatedBy           string                    `json:"createdBy"`
	ContentVariances    []ContentVariancesDetails `json:"contentVariances"`
	ContentSeasons      []ContentSeasonsDetails   `json:"contentSeasons"`
	Id                  string                    `json:"id"`
	CreatedByUserId     string                    `json:"-"`
}
type UserDetails struct {
	UserName string `json:"user_name"`
}
type Pagination struct {
	Size   int `json:"size"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// error code validation for episode
type CreateSeasonRequestValidation struct {
	// for removing sync below fields are commented -- SecondarySeasonId,SeasonKey,VarianceTrailerIds
	//	SecondarySeasonId  string                      `json:"secondarySeasonId"` // taking content id from request body for creating old contents
	//	SeasonKey          int                         `json:"seasonKey"`
	//	VarianceTrailerIds []string                    `json:"varianceTrailerIds"` // taking season trailer id from request body for creating old contents
	SeasonId         string                      `json:"Seasonid"`
	ContentId        *string                     `json:"contentId"`
	PrimaryInfo      *PrimaryInfoRequest         `json:"primaryInfo"`
	SeasonGenres     []ContentGenres             `json:"seasonGenres"`
	VarianceTrailers *[]VarianceTrailers         `json:"varianceTrailers"`
	AboutTheContent  *AboutTheContentInfoRequest `json:"aboutTheContent"`
	IntroStart       string                      `json:"introStart,omitempty"`
	Translation      *ContentTranslation         `json:"translation"`
	NonTextualData   *NonTextualData             `json:"nonTextualData"`
	Cast             *Cast                       `json:"cast"`
	Music            *Music                      `json:"music"`
	TagInfo          *TagInfo                    `json:"tagInfo"`
	Rights           *RightsRequest              `json:"rights"`
	CountryCheck     bool                        `json:"countryCheck"`
	Products         *[]int                      `json:"products"`
	SeoDetails       SeoDetails                  `json:"seoDetails"`
}

// season variance
type CreateSeasonVarainceRequest struct {
	SeasonId         string                      `json:"Seasonid"`
	ContentId        string                      `json:"contentId"`
	PrimaryInfo      *PrimaryInfoRequest         `json:"primaryInfo"`
	SeasonGenres     []ContentGenres             `json:"contentGenres"`
	VarianceTrailers *[]interface{}              `json:"varianceTrailers"`
	AboutTheContent  *AboutTheContentInfoRequest `json:"aboutTheContent"`
	IntroStart       string                      `json:"introStart,omitempty"`
	Translation      *ContentTranslation         `json:"translation"`
	NonTextualData   *NonTextualData             `json:"nonTextualData"`
	Cast             *Cast                       `json:"cast"`
	Music            *Music                      `json:"music"`
	TagInfo          *TagInfo                    `json:"tagInfo"`
	Rights           *RightsRequests             `json:"rights"`
	Products         *[]int                      `json:"products"`
	SeoDetails       SeoDetails                  `json:"seoDetails"`
}
type Episode struct {
	Id                     string    `json:"id"`
	SeasonId               string    `json:"seasonId"`
	Number                 int       `json:"number"`
	PrimaryInfoId          string    `json:"primaryInfoId"`
	PlaybackItemId         string    `json:"playbackItemId"`
	Status                 int       `json:"status"`
	SynopsisEnglish        string    `json:"synopsisEnglish"`
	SynopsisArabic         string    `json:"synopsisArabic"`
	CastId                 string    `json:"csatId"`
	MusicId                string    `json:"musicId"`
	TagInfoId              string    `json:"tagInfoId"`
	HasPosterImage         bool      `json:"hasPosterImage"`
	HasDubbingScript       bool      `json:"hasDubbingScript"`
	HasSubtitlingScript    bool      `json:"hasSubtitlingScript"`
	EpisodeKey             int       `json:"episodekey"`
	CreatedAt              time.Time `json:"createdAt"`
	ModifiedAt             time.Time `json:"modifiedAt"`
	EnglishMetaTitle       string    `json:"englishMetaTitle"`
	ArabicMetaTitle        string    `json:"arabicMetaTitle"`
	EnglishMetaDescription string    `json:"englishMetaDescription"`
	ArabicMetaDescription  string    `json:"arabicMetaDescription"`
}
type ContentTranslationData struct {
	Id                 string  `json:"id"`
	LanguageType       string  `json:"languageType"`
	DubbingLanguage    *string `json:"dubbingLanguage"`
	DubbingDialectId   *int    `json:"dubbingDialectId"`
	SubtitlingLanguage *string `json:"subtitlingLanguage"`
}

// season variance rights request
type RightsRequests struct {
	DigitalRightsType      int      `json:"digitalRightsType"`
	DigitalRightsStartDate string   `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   string   `json:"digitalRightsendDate"`
	DigitalRightsRegions   []string `json:"digitalRightsRegions"`
	SubscriptionPlans      []int    `json:"subscriptionPlans"`
}
type ContentRightsCountries struct {
	Id              string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	ContentRightsId string `json:"content_rights_id"`
	CountryId       string `json:"country_id"`
}

type VarianceTrailers struct {
	Order                 int    `json:"order,omitempty"`
	VideoTrailerId        string `json:"videoTrailerId,omitempty"`
	EnglishTitle          string `json:"englishTitle,omitempty"`
	ArabicTitle           string `json:"arabicTitle,omitempty"`
	Duration              int    `json:"duration,omitempty"`
	HasTrailerPosterImage bool   `json:"hasTrailerPosterImage,omitempty"`
	TrailerPosterImage    string `json:"trailerposterImage"`
	Id                    string `json:"Id,omitempty"`
	SeasonId              string `json:"seasonid,omitempty"`
}
type VarianceTrailer struct {
	Order                 int    `json:"order,omitempty"`
	VideoTrailerId        string `json:"videoTrailerId,omitempty"`
	EnglishTitle          string `json:"englishTitle,omitempty"`
	ArabicTitle           string `json:"arabicTitle,omitempty"`
	Duration              int    `json:"duration,omitempty"`
	HasTrailerPosterImage bool   `json:"hasTrailerPosterImage,omitempty"`
	Id                    string `json:"Id,omitempty"`
	SeasonId              string `json:"seasonid,omitempty"`
}

// season variance episode image upload
type Images struct {
	Imagename string `json:"imagename"`
	HasImage  bool   `json:"hasimage"`
}

type RedisErrorResponse struct {
	Error string `json:"error`
}
