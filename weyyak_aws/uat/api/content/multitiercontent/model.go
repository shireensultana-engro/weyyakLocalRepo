package multitiercontent

import (
	"time"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Multitier struct {
	TransliteratedTitle string `json:"transliteratedTitle"`
	Id                  string `json:"id"`
}
type Details struct {
	Id string `json:"id"`
}
type UpdateStatus struct {
	Status     string    `json:"status"`
	ModifiedAt time.Time `json:"modified_at"`
}
type UpdateDetails struct {
	DeletedByUserId string    `json:"deleted_by_user_id"`
	ModifiedAt      time.Time `json:"modified_at"`
}
type ContentType struct {
	ContentType string `json:"content_type"`
}

// create or update multitier content details
type StatusDetails struct {
	Id        string `json:"id"`
	Status    int    `json:"status"`
	CastId    string `json:"cast_id"`
	MusicId   string `json:"music_id"`
	TagInfoId string `json:"tag_info"`
}

type MainResponse struct {
	ContentId  string `json:"contentId"`  // taking content id from request body for creating old contents
	ContentKey int    `json:"contentKey"` // taking content key from request body for creating old contents
	//	CreatedByUserId string          `json:"userId"` // taking created by user id from request body for creating old contents
	PrimaryInfo   PrimaryInfo     `json:"primaryInfo"`
	ContentGenres []ContentGenres `json:"contentGenres"`
	SeoDetails    SeoDetails      `json:"seoDetails"`
}
type PrimaryInfo struct {
	Id                  string `json:"id"`
	ContentType         string `json:"contentType"`
	OriginalTitle       string `json:"originalTitle"`
	AlternativeTitle    string `json:"alternativeTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	TransliteratedTitle string `json:"transliteratedTitle"`
	Notes               string `json:"notes"`
	IntroStart          string `json:"introStart"`
	OutroStart          string `json:"outroStart"`
}

type PrimaryInfoRequest struct {
	Id                  string `json:"id"`
	OriginalTitle       string `json:"originalTitle"`
	AlternativeTitle    string `json:"alternativeTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	TransliteratedTitle string `json:"transliteratedTitle"`
	Notes               string `json:"notes"`
	IntroStart          string `json:"introStart"`
	OutroStart          string `json:"outroStart"`
}

type SeoDetails struct {
	EnglishMetaTitle       string `json:"englishMetaTitle"`
	ArabicMetaTitle        string `json:"arabicMetaTitle"`
	EnglishMetaDescription string `json:"englishMetaDescription"`
	ArabicMetaDescription  string `json:"arabicMetaDescription"`
}

type ContentGenres struct {
	GenreId    string   `json:"genreId"`
	SubgenreId []string `json:"subgenresId"`
	Id         string   `json:"id"`
}
type ContentGenresUpdate struct {
	GenreId string `json:"genreId"`
	Id      string `json:"id"`
}
type ContentCasts struct {
	Id            string  `json:"id"`
	MainActorId   *string `json:"main_actor_id,omitempty"`
	MainActressId *string `json:"main_actress_id,omitempty"`
}
type SeoDetailsResponse struct {
	Id                     string    `json:"id"`
	ContentKey             int       `json:"contentkey"`
	Status                 int       `json:"status"`
	ContentType            string    `json:"contenttype"`
	ContentTier            int       `json:"contenttier"`
	EnglishMetaTitle       string    `json:"englishMetaTitle"`
	ArabicMetaTitle        string    `json:"arabicMetaTitle"`
	EnglishMetaDescription string    `json:"englishMetaDescription"`
	ArabicMetaDescription  string    `json:"arabicMetaDescription"`
	CreatedByUserId        string    `json:"created_by_user_id"`
	CreatedAt              time.Time `json:"created_at,omitempty"`
	ModifiedAt             time.Time `json:"modified_at,omitempty"`
	CastId                 string    `json:"cast_id"`
	MusicId                string    `json:"music_id"`
	TagInfoId              string    `json:"tag_info_id"`
}

type PrimaryInfoIdDetails struct {
	PrimaryInfoId string `json:"primary_info_id"`
}
type ContentGenreResponse struct {
	Id        string `json:"id"`
	ContentId string `json:"content_id"`
	Order     int    `json:"order"`
	GenreId   string `json:"genreId"`
}

type SubGenreResponse struct {
	ContentGenreId string `json:"content_genre_id"`
	Order          int    `json:"order"`
	SubgenreId     string `json:"subgenreid"`
}

type ContentKeyResponse struct {
	ContentKey int `json:"contentkey"`
}

// create or update one tier content

//Onetier content request
type OnetierContentRequest struct {
	ContentId  string `json:"contentId"`  // taking content id from request body for creating old contents
	ContentKey int    `json:"contentKey"` // taking content key from request body for creating old contents
	//	CreatedByUserId string          `json:"userId"`    // taking created by user id from request body for creating old contents
	TextualData    TextualData     `json:"textualData" binding:"required"`
	NonTextualData *NonTextualData `json:"nonTextualData" binding:"required"`
	VarianceIds    []string        `json:"varianceIds"` // for sync
}

//TextualData content request
type TextualData struct {
	PrimaryInfo      *PrimaryInforequest     `json:"primaryInfo"  binding:"required"`
	SeoDetails       *SeoDetailsrequest      `json:"seoDetails" binding:"required"`
	ContentGenres    *[]ContentGenresrequest `json:"contentGenres" binding:"required"`
	ContentVariances *[]ContentVariances     `json:"contentVariances" binding:"required"`
	AboutTheContent  *AboutTheContent        `json:"aboutTheContent" binding:"required"`
	Cast             *Cast                   `json:"cast" binding:"required"`
	Music            *Music                  `json:"music" binding:"required"`
	TagInfo          *TagInfo                `json:"tagInfo" binding:"required"`
}
type PrimaryInforequest struct {
	ContentType         string `json:"contentType"`
	OriginalTitle       string `json:"originalTitle"`
	AlternativeTitle    string `json:"alternativeTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	TransliteratedTitle string `json:"transliteratedTitle"`
	Notes               string `json:"notes"`
	IntroStart          string `json:"introStart"`
	OutroStart          string `json:"outroStart"`
}
type SeoDetailsrequest struct {
	EnglishMetaTitle       string `json:"englishMetaTitle"`
	ArabicMetaTitle        string `json:"arabicMetaTitle"`
	EnglishMetaDescription string `json:"englishMetaDescription"`
	ArabicMetaDescription  string `json:"arabicMetaDescription"`
}
type ContentGenresrequest struct {
	GenreId     string   `json:"genreId"`
	SubgenresId []string `json:"subgenresId"`
}
type ContentVariances struct {
	//	ID                     string             `json:"id"`
	Status                 int                `json:"status"`
	StatusCanBeChanged     bool               `json:"statusCanBeChanged"`
	SubStatusName          string             `json:"subStatusName"`
	VideoContentId         string             `json:"videoContentId"`
	LanguageType           string             `json:"languageType"`
	OverlayPosterImage     string             `json:"overlayPosterImage"`
	DubbingScript          string             `json:"dubbingScript"`
	SubtitlingScript       string             `json:"subtitlingScript"`
	DubbingLanguage        string             `json:"dubbingLanguage"`
	DubbingDialectId       int                `json:"dubbingDialectId"`
	SubtitlingLanguage     string             `json:"subtitlingLanguage"`
	DigitalRightsType      int                `json:"digitalRightsType"`
	DigitalRightsStartDate string             `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   string             `json:"digitalRightsEndDate"`
	DigitalRightsRegions   []int              `json:"digitalRightsRegions"`
	SchedulingDateTime     *time.Time         `json:"schedulingDateTime"`
	CreatedBy              string             `json:"createdBy"`
	PublishingPlatforms    []int              `json:"publishingPlatforms"`
	Products               []int              `json:"products"`
	SubscriptionPlans      []int              `json:"subscriptionPlans"`
	CountryCheck           bool               `json:"countryCheck"`
	IntroDuration          string             `json:"introDuration"`
	IntroStart             string             `json:"introStart"`
	VarianceTrailers       []VarianceTrailers `json:"varianceTrailers"`
	Id                     string             `json:"id"`
	VarianceTrailerIds     []string           `json:"varianceTrailerIds"` // for sync
}

type VarianceTrailers struct {
	Order                 int    `json:"order,omitempty"`
	VideoTrailerId        string `json:"videoTrailerId,omitempty"`
	EnglishTitle          string `json:"englishTitle,omitempty"`
	ArabicTitle           string `json:"arabicTitle,omitempty"`
	Duration              int    `json:"duration,omitempty"`
	HasTrailerPosterImage bool   `json:"hasTrailerPosterImage,omitempty"`
	TrailerPosterImage    string `json:"trailerposterImage,omitempty"`
	Id                    string `json:"Id,omitempty"`
}
type AboutTheContent struct {
	OriginalLanguage      string `json:"originalLanguage"`
	Supplier              string `json:"supplier"`
	AcquisitionDepartment string `json:"acquisitionDepartment"`
	EnglishSynopsis       string `json:"englishSynopsis"`
	ArabicSynopsis        string `json:"arabicSynopsis"`
	ProductionYear        int    `json:"productionYear"`
	ProductionHouse       string `json:"productionHouse"`
	AgeGroup              int    `json:"ageGroup"`
	IntroDuration         string `json:"introDuration"`
	IntroStart            string `json:"introStart"`
	OutroDuration         string `json:"outroDuration"`
	OutroStart            string `json:"outroStart"`
	ProductionCountries   []int  `json:"productionCountries"`
}
type Cast struct {
	MainActorId   string   `json:"mainActorId"`
	MainActressId string   `json:"mainActressId"`
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
	LanguageType       int    `json:"languageType"`
	DubbingLanguage    string `json:"dubbingLanguage"`
	DubbingDialectId   int    `json:"dubbingDialectId"`
	SubtitlingLanguage string `json:"subtitlingLanguage"`
}
type NonTextualData struct {
	PosterImage             string `json:"posterImage"`
	DetailsBackground       string `json:"detailsBackground"`
	MobileDetailsBackground string `json:"mobileDetailsBackground"`
}

//Content - struct for DB binding
type Content struct {
	//uuid.UUID
	Id                         string    `json:"id"`
	AverageRating              float64   `json:"average_rating"`
	AverageRatingUpdatedAt     time.Time `json:"average_rating_updated_at"`
	ContentKey                 int       `json:"content_key"`
	ContentType                string    `json:"content_type"`
	Status                     int       `json:"status"`
	ModifiedAt                 time.Time `json:"modified_at"`
	HasPosterImage             bool      `json:"has_poster_image"`
	HasDetailsBackground       bool      `json:"has_details_background"`
	HasMobileDetailsBackground bool      `json:"has_mobile_details_background"`
	CreatedByUserId            string    `json:"created_by_user_id"` //TODO:dependency with token
	ContentTier                int       `json:"content_tier"`
	PrimaryInfoId              string    `json:"primary_info_id" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	AboutTheContentInfoId      string    `json:"about_the_content_info_id" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	CastId                     string    `json:"cast_id" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	MusicId                    string    `json:"music_id" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	TagInfoId                  string    `json:"tag_info_id" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	//	DeletedByUserId            *string    `json:"deleted_by_user_id" gorm:"default:00000000-0000-0000-0000-000000000000;"` //TODO:dependency with token
	CreatedAt              time.Time `json:"created_at"`
	EnglishMetaTitle       string    `json:"english_meta_title"`
	ArabicMetaTitle        string    `json:"arabic_meta_title"`
	EnglishMetaDescription string    `json:"english_meta_description"`
	ArabicMetaDescription  string    `json:"arabic_meta_description"`
}

type ContentPrimaryInfo struct {
	Id                  string `json:"id"`
	OriginalTitle       string `json:"originalTitle"`
	AlternativeTitle    string `json:"alternativeTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	TransliteratedTitle string `json:"transliteratedTitle"`
	Notes               string `json:"notes"`
	IntroStart          string `json:"intro_start"`
	OutroStart          string `json:"outro_start"`
}
type ContentRights struct {
	Id                     string `json:"id,omitempty"`
	DigitalRightsType      int    `json:"digitalRightsType,omitempty"`
	DigitalRightsStartDate string `json:"digitalRightsStartDate,omitempty"`
	DigitalRightsEndDate   string `json:"digitalRightsEndDate,omitempty"`
}
type ContentTranslation struct {
	Id                 string `json:"id,omitempty"`
	LanguageType       int    `json:"language_type,omitempty"`
	DubbingLanguage    string `json:"dubbingLanguage,omitempty"`
	DubbingDialectId   int    `json:"dubbingDialectId,omitempty"`
	SubtitlingLanguage string `json:"subtitlingLanguage"`
}

type AboutTheContentInfo struct {
	Id                    string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	OriginalLanguage      string `json:"originalLanguage" binding:"required"`
	Supplier              string `json:"supplier" binding:"required"`
	AcquisitionDepartment string `json:"acquisitionDepartment"`
	EnglishSynopsis       string `json:"englishSynopsis" binding:"required"`
	ArabicSynopsis        string `json:"arabicSynopsis" binding:"required"`
	ProductionYear        int    `json:"productionYear"`
	ProductionHouse       string `json:"productionHouse"`
	AgeGroup              int    `json:"ageGroup" binding:"required"`
	IntroDuration         string `json:"introDuration"`
	IntroStart            string `json:"introStart"`
	OutroDuration         string `json:"outroDuration"`
	OutroStart            string `json:"outroStart"`
}
type ContentVariance struct {
	ID                    string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	ContentId             string `json:"content_id"`
	PlaybackItemId        string `json:"playback_item_id"`
	Order                 int    `json:"order"`
	HasOverlayPosterImage bool   `json:"has_overlay_poster_image"`
	HasDubbingScript      bool   `json:"has_dubbing_script"`
	HasSubtitlingScript   bool   `json:"has_subtitling_script"`
	//	DeletedByUserId       string    `json:"deleted_by_user_id" gorm:"default:00000000-0000-0000-0000-000000000000;"` //TODO:dependency with token
	Status        int       `json:"status"`
	HasAllRights  bool      `json:"has_all_rights"`
	IntroStart    string    `json:"intro_start"`
	IntroDuration string    `json:"intro_duration"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
}

type VarianceTrailer struct {
	Id                    string `json:"Id"`
	Order                 int    `json:"order,omitempty"`
	VideoTrailerId        string `json:"videoTrailerId,omitempty"`
	EnglishTitle          string `json:"englishTitle,omitempty"`
	ArabicTitle           string `json:"arabicTitle,omitempty"`
	Duration              int    `json:"duration,omitempty"`
	HasTrailerPosterImage bool   `json:"hasTrailerPosterImage,omitempty"`
	ContentVarianceId     string `json:"id,omitempty"`
}
type ContentCast struct {
	Id            string `json:"id"`
	MainActorId   string `json:"main_actor_id,omitempty"`
	MainActressId string `json:"main_actress_id,omitempty"`
}
type ContentRightsCountry struct {
	Id              string `json:"id"`
	ContentRightsId string `json:"content_rights_id"`
	CountryId       int    `json:"country_id"`
}

type PlaybackItem struct {
	Id                 string     `json:"Id"`
	VideoContentId     string     `json:"video_content_id"`
	SchedulingDateTime *time.Time `json:"SchedulingDateTime"`
	CreatedByUserId    string     `json:"CreatedByUserId"`
	TranslationId      string     `json:"TranslationId"`
	RightsId           string     `json:"rightsid"`
	Duration           int        `json:"duration"`
}

type ContentGenre struct {
	ContentId string `json:"content_id"`
	Id        string `json:"id"`
	Order     int    `json:"order"`
	GenreId   string `json:"genreid"`
}
type PlaybackItemTargetPlatform struct {
	PlaybackItemId string `json:"playbackitemid"`
	TargetPlatform int    `json:"targetplataform"`
	RightsId       string `json:"rights_id"`
}
type ContentMusic struct {
	Id string `json:"id"`
}
type Actors struct {
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}
type Writers struct {
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}
type Directors struct {
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}
type Singers struct {
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}
type MusicComposers struct {
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}
type SongWriters struct {
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}

type Actor struct {
	Id          string `json:"id"`
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}
type Writer struct {
	Id          string `json:"id"`
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}
type Director struct {
	Id          string `json:"id"`
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}
type Singer struct {
	Id          string `json:"id"`
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}
type MusicComposer struct {
	Id          string `json:"id"`
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}
type SongWriter struct {
	Id          string `json:"id"`
	EnglishName string `json:"englishname"`
	ArabicName  string `json:"arabicname"`
}

type ContentRightsPlan struct {
	RightsId           string `json:"rightsId"`
	SubscriptionPlanId int    `json:"subscriptionplanid"`
}
type RightsProduct struct {
	RightsId    string `json:"rights_id"`
	ProductName int    `json:"product_name"`
}
type ProductionCountry struct {
	AboutTheContentInfoId string `json:"aboutthecontentinfoid"`
	CountryId             int    `json:"country_id"`
}
type ContentActor struct {
	CastId  string `json:"cast_id"`
	ActorId string `json:"actor_id"`
}
type ContentWriter struct {
	CastId   string `json:"cast_id"`
	WriterId string `json:"actor_id"`
}
type ContentDirector struct {
	CastId     string `json:"cast_id"`
	DirectorId string `json:"actor_id"`
}
type ContentSinger struct {
	MusicId  string `json:"cast_id"`
	SingerId string `json:"actor_id"`
}
type ContentMusicComposer struct {
	MusicId         string `json:"cast_id"`
	MusicComposerId string `json:"actor_id"`
}
type ContentSongWriter struct {
	MusicId      string `json:"cast_id"`
	SongWriterId string `json:"actor_id"`
}

type ContentTagInfo struct {
	Id string `json:"id"`
}
type ContentTag struct {
	TagInfoId        string `json:"tagInfoId"`
	TextualDataTagId string `json:"textauldatatagid"`
}
type ContentgenreId struct {
	Id string `json:"id"`
}

// structs for validations
type OnetierContentRequestValidtion struct {
	TextualData    TextualDataValidation     `json:"textualData" binding:"required"`
	NonTextualData *NonTextualDataValidation `json:"nonTextualData" binding:"required"`
}
type TextualDataValidation struct {
	PrimaryInfo      *PrimaryInforequest    `json:"primaryInfo"`
	SeoDetails       *SeoDetailsrequest     `json:"seoDetails" binding:"required"`
	ContentGenres    []ContentGenresrequest `json:"contentGenres" binding:"required"`
	ContentVariances []ContentVariances     `json:"contentVariances" binding:"required"`
	AboutTheContent  *AboutTheContent       `json:"aboutTheContent" binding:"required"`
	Cast             *Cast                  `json:"cast" binding:"required"`
	Music            *Music                 `json:"music" binding:"required"`
	TagInfo          *TagInfo               `json:"tagInfo" binding:"required"`
}

type NonTextualDataValidation struct {
	PosterImage             string `json:"posterImage"`
	DetailsBackground       string `json:"detailsBackground"`
	MobileDetailsBackground string `json:"mobileDetailsBackground"`
}

/* This struct For uploading Images to ContentVariance */
type Variance struct {
	Id                 string
	OverlayPosterImage string
	DubbingScript      string
	SubtitlingScript   string
}

type ContentSubgenre struct {
	ContentGenreId string
	Order          int
	SubgenreId     string
}
