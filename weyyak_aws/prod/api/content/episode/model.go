package episode

import (
	"encoding/json"
	"time"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Get Episode Details Based on SeasonId
type UserInfo struct {
	UserName string `json:"userName"`
}
type EpisodeDetailsSummary struct {
	IsPrimary              bool        `json:"isPrimary"`
	UserId                 string      `json:"userId"`
	SecondarySeasonId      string      `json:"secondarySeasonId" `
	VarianceIds            []int       `json:"varianceIds"`
	EpisodeIds             []int       `json:"episodeIds"`
	SecondaryEpisodeId     string      `json:"secondaryEpisodeId"`
	ContentId              string      `json:"contentId"`
	EpisodeKey             int         `json:"episodeKey"`
	SeasonId               string      `json:"seasonId"`
	Status                 int         `json:"status"`
	StatusCanBeChanged     bool        `json:"statusCanBeChanged"`
	SubStatus              int         `json:"subStatus"`
	SubStatusName          string      `json:"subStatusName"`
	DigitalRightsType      int         `json:"digitalRightsType"`
	DigitalRightsStartDate time.Time   `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   time.Time   `json:"digitalRightsEndDate"`
	CreatedBy              string      `json:"createdBy"`
	PrimaryInfo            PrimaryInfo `json:"primaryInfo"`
	Cast                   []int       `json:"cast"`
	Music                  []int       `json:"music"`
	TagInfo                []int       `json:"tagInfo"`
	NonTextualData         []int       `json:"nonTextualData"`
	Translation            Translation `json:"translation"`
	SchedulingDateTime     *time.Time  `json:"schedulingDateTime"`
	PublishingPlatforms    []int       `json:"publishingPlatforms"`
	SeoDetails             []int       `json:"seoDetails"`
	Id                     string      `json:"id"`
}
type PrimaryInfo struct {
	Id                     string                   `json:"id,omitempty"`
	Number                 int                      `json:"number,omitempty"`
	VideoContentId         string                   `json:"videoContentId,omitempty"`
	SynopsisEnglish        string                   `json:"synopsisEnglish,omitempty"`
	SynopsisArabic         string                   `json:"synopsisArabic,omitempty"`
	SeasonNumber           int                      `json:"seasonNumber,omitempty"`
	OriginalTitle          string                   `json:"originalTitle"`
	AlternativeTitle       string                   `json:"alternativeTitle"`
	ArabicTitle            string                   `json:"arabicTitle"`
	TransliteratedTitle    string                   `json:"transliteratedTitle"`
	Notes                  string                   `json:"notes"`
	IntroStart             *string                  `json:"introStart"`
	OutroStart             *string                  `json:"outroStart"`
	SavedEpisodes          []int                    `json:"savedEpisodes,omitempty"`
	SavedEpisodesAndTitles []SavedEpisodesAndTitles `json:"savedEpisodesAndTitles,omitempty"`
	ContentType            string                   `json:"contentType,omitempty"`
}
type SavedEpisodesAndTitles struct {
	Number              int    `json:"number"`
	TransliteratedTitle string `json:"transliteratedTitle"`
}
type Translation struct {
	LanguageType       string  `json:"languageType"`
	DubbingLanguage    *string `json:"dubbingLanguage"`
	DubbingDialectId   *int    `json:"dubbingDialectId"`
	SubtitlingLanguage *string `json:"subtitlingLanguage"`
}

type FinalEpisodesResult struct {
	IsPrimary              bool       `json:"isPrimary"`
	UserId                 string     `json:"userId" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	SecondarySeasonId      string     `json:"secondarySeasonId" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	VarianceIds            []int      `json:"varianceIds"`
	EpisodeIds             []int      `json:"episodeIds"`
	SecondaryEpisodeId     string     `json:"secondaryEpisodeId" gorm:"default:00000000-0000-0000-0000-000000000000;"`
	ContentId              string     `json:"contentId"`
	EpisodeKey             int        `json:"episodeKey"`
	SeasonId               string     `json:"seasonId"`
	SeasonStatus           int        `json:"seasonstatus"`
	Status                 int        `json:"status"`
	StatusCanBeChanged     bool       `json:"statusCanBeChanged"`
	SubStatus              int        `json:"subStatus"`
	SubStatusName          string     `json:"subStatusName"`
	DigitalRightsType      int        `json:"digitalRightsType"`
	DigitalRightsStartDate time.Time  `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   time.Time  `json:"digitalRightsEndDate"`
	CreatedByUserId        string     `json:"createdByUserId"`
	Number                 int        `json:"number"`
	VideoContentId         string     `json:"videoContentId"`
	SynopsisEnglish        string     `json:"synopsisEnglish"`
	SynopsisArabic         string     `json:"synopsisArabic"`
	OriginalTitle          string     `json:"originalTitle"`
	AlternativeTitle       string     `json:"alternativeTitle"`
	ArabicTitle            string     `json:"arabicTitle"`
	TransliteratedTitle    string     `json:"transliteratedTitle"`
	Notes                  string     `json:"notes"`
	IntroStart             *string    `json:"introStart"`
	OutroStart             *string    `json:"outroStart"`
	Cast                   []int      `json:"cast"`
	Music                  []int      `json:"music"`
	TagInfo                []int      `json:"tagInfo"`
	NonTextualData         []int      `json:"nonTextualData"`
	LanguageType           int        `json:"languageType"`
	DubbingLanguage        *string    `json:"dubbingLanguage"`
	DubbingDialectId       *int       `json:"dubbingDialectId"`
	SubtitlingLanguage     *string    `json:"subtitlingLanguage"`
	SchedulingDateTime     *time.Time `json:"schedulingDateTime"`
	PublishingPlatforms    []int      `json:"publishingPlatforms"`
	SeoDetails             []int      `json:"seoDetails"`
	Id                     string     `json:"id"`
	RightsId               string     `json:"rightsId,omitempty"`
	CastId                 string     `json:"castId,omitempty"`
	MusicId                string     `json:"musicId,omitempty"`
	TagInfoId              string     `json:"tagInfoId,omitempty"`
	PlaybackItemId         string     `json:"playbackItemId,omitempty"`
	MainActorId            *string    `json:"mainActorId,omitempty"`
	MainActressId          *string    `json:"mainActressId,omitempty"`
}

//Get Episode Details Based in contentId
type SeasonDetailsSummary struct {
	ContentId          string      `json:"contentId"`
	SeasonKey          int         `json:"seasonKey"`
	Status             int         `json:"status"`
	StatusCanBeChanged bool        `json:"statusCanBeChanged"`
	SubStatusName      string      `json:"subStatusName"`
	ModifiedAt         time.Time   `json:"modifiedAt"`
	PrimaryInfo        PrimaryInfo `json:"primaryInfo"`
	Cast               *string     `json:"cast"`
	Music              *string     `json:"music"`
	TagInfo            *string     `json:"tagInfo"`
	SeasonGenres       *string     `json:"seasonGenres"`
	AboutTheContent    *string     `json:"aboutTheContent"`
	Translation        Translation `json:"translation"`
	Episodes           *string     `json:"episodes"`
	NonTextualData     *string     `json:"nonTextualData"`
	Rights             Rights      `json:"rights"`
	CreatedBy          string      `json:"createdBy"`
	IntroDuration      string      `json:"introDuration"`
	IntroStart         string      `json:"introStart"`
	OutroDuration      string      `json:"outroDuration"`
	OutroStart         string      `json:"outroStart"`
	Products           *string     `json:"products"`
	SeoDetails         *string     `json:"seoDetails"`
	VarianceTrailers   *string     `json:"varianceTrailers"`
	Id                 string      `json:"id"`
}
type Rights struct {
	Id                     string    `json:"id ,omitempty"`
	DigitalRightsType      int       `json:"digitalRightsType"`
	DigitalRightsStartDate time.Time `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   time.Time `json:"digitalRightsEndDate"`
	DigitalRightsRegions   []int     `json:"digitalRightsRegions"`
	SubscriptionPlans      []int     `json:"subscriptionPlans"`
}
type DigitalRightsRegions struct {
	CountryId int `json:"country_Id"`
}
type FinalSeasonResult struct {
	ContentId                  string     `json:"contentId"`
	SeasonKey                  int        `json:"seasonKey"`
	ContentStatus              int        `json:"contentStatus"`
	Status                     int        `json:"status"`
	StatusCanBeChanged         bool       `json:"statusCanBeChanged"`
	SubStatusName              string     `json:"subStatusName"`
	ModifiedAt                 time.Time  `json:"modifiedAt"`
	SeasonNumber               int        `json:"seasonNumber"`
	OriginalTitle              string     `json:"originalTitle"`
	AlternativeTitle           string     `json:"alternativeTitle"`
	ArabicTitle                string     `json:"arabicTitle"`
	TransliteratedTitle        string     `json:"transliteratedTitle"`
	Notes                      string     `json:"notes"`
	IntroStart                 *string    `json:"introStart"`
	OutroStart                 *string    `json:"outroStart"`
	Cast                       *string    `json:"cast"`
	MainActorId                *string    `json:"mainActorId,omitempty"`
	MainActressId              *string    `json:"mainActressId,omitempty"`
	Music                      *string    `json:"music"`
	TagInfo                    *string    `json:"tagInfo"`
	SeasonGenres               *string    `json:"seasonGenres"`
	AboutTheContent            *string    `json:"aboutTheContent"`
	LanguageType               int        `json:"languageType"`
	DubbingLanguage            *string    `json:"dubbingLanguage"`
	DubbingDialectId           *int       `json:"dubbingDialectId"`
	SubtitlingLanguage         *string    `json:"subtitlingLanguage"`
	Episodes                   *string    `json:"episodes"`
	NonTextualData             *string    `json:"nonTextualData"`
	DigitalRightsType          int        `json:"digitalRightsType"`
	DigitalRightsStartDate     time.Time  `json:"digitalRightsStartDate"`
	DigitalRightsEndDate       time.Time  `json:"digitalRightsEndDate"`
	SubscriptionPlans          int        `json:"subscriptionPlans"`
	CreatedBy                  string     `json:"createdBy"`
	IntroDuration              string     `json:"introDuration"`
	AboutIntroStart            *string    `json:"aboutintroStart"`
	OutroDuration              string     `json:"outroDuration"`
	AboutOutroStart            *string    `json:"aboutoutroStart"`
	Products                   *string    `json:"products"`
	SeoDetails                 *string    `json:"seoDetails"`
	VarianceTrailers           *string    `json:"varianceTrailers"`
	Id                         string     `json:"id"`
	RightsId                   string     `json:"rightsId,omitempty"`
	CastId                     string     `json:"castId,omitempty"`
	MusicId                    string     `json:"musicId,omitempty"`
	TagInfoId                  string     `json:"tagInfoId,omitempty"`
	OriginalLanguage           string     `json:"originalLanguage,omitempty"`
	Supplier                   string     `json:"supplier,omitempty"`
	AcquisitionDepartment      string     `json:"acquisitionDepartment,omitempty"`
	EnglishSynopsis            string     `json:"englishSynopsis,omitempty"`
	ArabicSynopsis             string     `json:"arabicSynopsis,omitempty"`
	ProductionYear             int        `json:"productionYear,omitempty"`
	ProductionHouse            string     `json:"productionHouse,omitempty"`
	AgeGroup                   int        `json:"ageGroup,omitempty"`
	AboutIntroDuration         *string    `json:"aboutIntroDuration,omitempty"`
	AboutOutroDuration         *string    `json:"aboutOutroDuration,omitempty"`
	EnglishMetaTitle           string     `json:"englishMetaTitle,omitempty"`
	ArabicMetaTitle            string     `json:"arabicMetaTitle,omitempty"`
	EnglishMetaDescription     string     `json:"englishMetaDescription,omitempty"`
	ArabicMetaDescription      string     `json:"arabicMetaDescription,omitempty"`
	AboutTheContentInfoId      string     `json:"aboutTheContentInfoId,omitempty"`
	ContentKey                 int        `json:"contentKey,omitempty"`
	Duration                   *int       `json:"duration,omitempty"`
	ContentType                string     `json:"contentType,omitempty"`
	VarianceStatus             int        `json:"varianceStatus,omitempty"`
	VideoContentId             string     `json:"VideoContentId,omitempty"`
	HasOverlayPosterImage      bool       `json:"hasOverlayPosterImage,omitempty"`
	HasDubbingScript           bool       `json:"hasDubbingScript,omitempty"`
	HasSubtitlingScript        bool       `json:"hasSubtitlingScript,omitempty"`
	VarianceId                 string     `json:"varianceId"`
	VarianceIntroDuration      string     `json:"varianceIntroDuration"`
	VarianceIntroStart         string     `json:"arianceIntroStart"`
	SchedulingDateTime         *time.Time `json:"schedulingDateTime"`
	CreatedByUserId            *string    `json:"createdByUserId"`
	PlaybackItemId             string     `json:"playbackItemId"`
	HasPosterImage             bool       `json:"hasPosterImage"`
	HasDetailsBackground       bool       `json:"hasDetailsBackground"`
	HasMobileDetailsBackground bool       `json:"hasMobileDetailsBackground"`
}

type FinalSeasonResultNew struct {
	ContentId                  string     `json:"contentId"`
	SeasonKey                  int        `json:"seasonKey"`
	ContentStatus              int        `json:"contentStatus"`
	Status                     int        `json:"status"`
	StatusCanBeChanged         bool       `json:"statusCanBeChanged"`
	SubStatusName              string     `json:"subStatusName"`
	ModifiedAt                 time.Time  `json:"modifiedAt"`
	SeasonNumber               int        `json:"seasonNumber"`
	OriginalTitle              string     `json:"originalTitle"`
	AlternativeTitle           string     `json:"alternativeTitle"`
	ArabicTitle                string     `json:"arabicTitle"`
	TransliteratedTitle        string     `json:"transliteratedTitle"`
	Notes                      string     `json:"notes"`
	IntroStart                 *string    `json:"introStart"`
	OutroStart                 *string    `json:"outroStart"`
	Cast                       *string    `json:"cast"`
	MainActorId                *string    `json:"mainActorId,omitempty"`
	MainActressId              *string    `json:"mainActressId,omitempty"`
	Music                      *string    `json:"music"`
	TagInfo                    *string    `json:"tagInfo"`
	SeasonGenres               *string    `json:"seasonGenres"`
	AboutTheContent            *string    `json:"aboutTheContent"`
	LanguageType               int        `json:"languageType"`
	DubbingLanguage            *string    `json:"dubbingLanguage"`
	DubbingDialectId           *int       `json:"dubbingDialectId"`
	SubtitlingLanguage         *string    `json:"subtitlingLanguage"`
	Episodes                   *string    `json:"episodes"`
	NonTextualData             *string    `json:"nonTextualData"`
	DigitalRightsType          int        `json:"digitalRightsType"`
	DigitalRightsStartDate     *time.Time `json:"digitalRightsStartDate"`
	DigitalRightsEndDate       *time.Time `json:"digitalRightsEndDate"`
	SubscriptionPlans          int        `json:"subscriptionPlans"`
	CreatedBy                  string     `json:"createdBy"`
	IntroDuration              string     `json:"introDuration"`
	AboutIntroStart            *string    `json:"aboutintroStart"`
	OutroDuration              string     `json:"outroDuration"`
	AboutOutroStart            *string    `json:"aboutoutroStart"`
	Products                   *string    `json:"products"`
	SeoDetails                 *string    `json:"seoDetails"`
	VarianceTrailers           *string    `json:"varianceTrailers"`
	Id                         string     `json:"id"`
	RightsId                   string     `json:"rightsId,omitempty"`
	CastId                     string     `json:"castId,omitempty"`
	MusicId                    string     `json:"musicId,omitempty"`
	TagInfoId                  string     `json:"tagInfoId,omitempty"`
	OriginalLanguage           string     `json:"originalLanguage,omitempty"`
	Supplier                   string     `json:"supplier,omitempty"`
	AcquisitionDepartment      string     `json:"acquisitionDepartment,omitempty"`
	EnglishSynopsis            string     `json:"englishSynopsis,omitempty"`
	ArabicSynopsis             string     `json:"arabicSynopsis,omitempty"`
	ProductionYear             int        `json:"productionYear,omitempty"`
	ProductionHouse            string     `json:"productionHouse,omitempty"`
	AgeGroup                   int        `json:"ageGroup,omitempty"`
	AboutIntroDuration         *string    `json:"aboutIntroDuration,omitempty"`
	AboutOutroDuration         *string    `json:"aboutOutroDuration,omitempty"`
	EnglishMetaTitle           string     `json:"englishMetaTitle,omitempty"`
	ArabicMetaTitle            string     `json:"arabicMetaTitle,omitempty"`
	EnglishMetaDescription     string     `json:"englishMetaDescription,omitempty"`
	ArabicMetaDescription      string     `json:"arabicMetaDescription,omitempty"`
	AboutTheContentInfoId      string     `json:"aboutTheContentInfoId,omitempty"`
	ContentKey                 int        `json:"contentKey,omitempty"`
	Duration                   *int       `json:"duration,omitempty"`
	ContentType                string     `json:"contentType,omitempty"`
	VarianceStatus             int        `json:"varianceStatus,omitempty"`
	VideoContentId             string     `json:"VideoContentId,omitempty"`
	HasOverlayPosterImage      bool       `json:"hasOverlayPosterImage,omitempty"`
	HasDubbingScript           bool       `json:"hasDubbingScript,omitempty"`
	HasSubtitlingScript        bool       `json:"hasSubtitlingScript,omitempty"`
	VarianceId                 string     `json:"varianceId"`
	VarianceIntroDuration      string     `json:"varianceIntroDuration"`
	VarianceIntroStart         string     `json:"arianceIntroStart"`
	SchedulingDateTime         *time.Time `json:"schedulingDateTime"`
	CreatedByUserId            *string    `json:"createdByUserId"`
	PlaybackItemId             string     `json:"playbackItemId"`
	HasPosterImage             bool       `json:"hasPosterImage"`
	HasDetailsBackground       bool       `json:"hasDetailsBackground"`
	HasMobileDetailsBackground bool       `json:"hasMobileDetailsBackground"`
}

//Create or Update episodes
type Episodes struct {
	SecondaryEpisodeId string `json:"secondaryEpisodeId"` // episodeId for creating old contents
	EpisodeKey         int    `json:"episodeKey"`         // episodeId for creating old contents
	//	CreatedByUserId     string          `json:"userId"`             // userid for creating old contents
	SeasonId            string          `json:"seasonId"`
	PrimaryInfo         *PrimaryInfo    `json:"primaryInfo"`
	NonTextualData      *NonTextualData `json:"nonTextualData"`
	Rights              *Rights         `json:"rights"`
	Cast                *Cast           `json:"cast"`
	Music               *Music          `json:"music"`
	TagInfo             *TagInfo        `json:"tagInfo"`
	SchedulingDateTime  *time.Time      `json:"schedulingDateTime"`
	PublishingPlatforms []int           `json:"publishingPlatforms"`
	SeoDetails          SeoDetails      `json:"seoDetails"`
}
type NonTextualData struct {
	PosterImage             string `json:"posterImage"`
	DubbingScript           string `json:"dubbingScript,omitempty"`
	SubtitlingScript        string `json:"subtitlingScript,omitempty"`
	OverlayPosterImage      string `json:"overlayPosterImage,omitempty"`
	DetailsBackground       string `json:"detailsBackground,omitempty"`
	MobileDetailsBackground string `json:"mobileDetailsBackground,omitempty"`
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
type PublishingPlatforms struct {
	PublishingPlatforms int `json:"publishingPlatforms"`
}
type SeoDetails struct {
	EnglishMetaTitle       string `json:"englishMetaTitle"`
	ArabicMetaTitle        string `json:"arabicMetaTitle"`
	EnglishMetaDescription string `json:"englishMetaDescription"`
	ArabicMetaDescription  string `json:"arabicMetaDescription"`
}

/*create or update record*/
type CreateEpisode struct {
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
type CreatePrimaryInfo struct {
	Id                  string  `json:"id ,omitempty"`
	OriginalTitle       string  `json:"originalTitle"`
	AlternativeTitle    string  `json:"alternativeTitle"`
	ArabicTitle         string  `json:"arabicTitle"`
	TransliteratedTitle string  `json:"transliteratedTitle"`
	Notes               string  `json:"notes"`
	IntroStart          *string `json:"introStart"`
	OutroStart          *string `json:"outroStart"`
}
type FetchEpisodeDetails struct {
	Id             string `json:"id"`
	EpisodeKey     int    `json:"episodeKey"`
	PlaybackItemId string `json:"playbackItemId"`
	PrimaryInfoId  string `json:"primaryInfoId"`
	CastId         string `json:"castId"`
	MusicId        string `json:"musicId"`
	TagInfoId      string `json:"tagInfoId"`
	RightsId       string `json:"rightsId"`
}
type InsertCast struct {
	MainActorId   string `json:"mainActorId"`
	MainActressId string `json:"mainActressId"`
	Id            string `json:"Id"`
}
type ContentRights struct {
	Id                     string    `json:"id ,omitempty"`
	DigitalRightsType      int       `json:"digitalRightsType"`
	DigitalRightsStartDate time.Time `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   time.Time `json:"digitalRightsEndDate"`
}
type ContentRightsCountry struct {
	Id              string `json:"id"`
	ContentRightsId string `json:"contentRightsId"`
	CountryId       int    `json:"countryId"`
}
type ContentActor struct {
	CastId  string `json:"castId,omitempty"`
	ActorId string `json:"actorId"`
}
type ContentWriter struct {
	CastId   string `json:"castId"`
	WriterId string `json:"WriterId"`
}
type ContentDirector struct {
	CastId     string `json:"castId"`
	DirectorId string `json:"directorId"`
}

type ContentSinger struct {
	SingerId string `json:"singerId"`
	MusicId  string `json:"musicId"`
}
type ContentMusicComposer struct {
	MusicComposerId string `json:"musicComposerId"`
	MusicId         string `json:"musicId"`
}
type ContentSongWriter struct {
	SongWriterId string `json:"songWriterId"`
	MusicId      string `json:"musicId"`
}
type ContentTag struct {
	TagInfoId        string `json:"tagInfoId"`
	TextualDataTagId string `json:"textualDataTagId"`
}
type PlaybackItem struct {
	Id                 string     `json:"id ,omitempty"`
	VideoContentId     string     `json:"videoContentId"`
	SchedulingDateTime *time.Time `json:"schedulingDateTime"`
	CreatedByUserId    string     `json:"createdByUserId"`
	TranslationId      string     `json:"translationId"`
	RightsId           string     `json:"rightsId"`
	Duration           int        `json:"duration"`
}
type PlaybackItemTargetPlatform struct {
	PlaybackItemId string `json:"playbackItemId"`
	TargetPlatform int    `json:"targetPlatform"`
	RightsId       string `json:"rightsId"`
}

/*Get Season details by season id*/
type SeasonResult struct {
	ContentId          string          `json:"contentId"`
	SeasonKey          int             `json:"seasonKey"`
	Status             int             `json:"status"`
	StatusCanBeChanged bool            `json:"statusCanBeChanged"`
	ModifiedAt         time.Time       `json:"modifiedAt"`
	PrimaryInfo        PrimaryInfo     `json:"primaryInfo"`
	Cast               Cast            `json:"cast"`
	Music              Music           `json:"music"`
	TagInfo            TagInfo         `json:"tagInfo"`
	SeasonGenres       []interface{}   `json:"seasonGenres"`
	AboutTheContent    AboutTheContent `json:"aboutTheContent"`
	Translation        Translation     `json:"translation"`
	Episodes           interface{}     `json:"episodes"`
	NonTextualData     NonTextualData  `json:"nonTextualData"`
	Rights             Rights          `json:"rights"`
	CreatedBy          string          `json:"createdBy"`
	IntroDuration      string          `json:"introDuration"`
	IntroStart         string          `json:"introStart"`
	OutroDuration      string          `json:"outroDuration"`
	OutroStart         string          `json:"outroStart"`
	Products           []int           `json:"products"`
	SeoDetails         SeoDetails      `json:"seoDetails"`
	VarianceTrailers   []interface{}   `json:"varianceTrailers"`
	Id                 string          `json:"id"`
}
type AboutTheContent struct {
	OriginalLanguage      string  `json:"originalLanguage"`
	Supplier              string  `json:"supplier"`
	AcquisitionDepartment string  `json:"acquisitionDepartment"`
	EnglishSynopsis       string  `json:"englishSynopsis"`
	ArabicSynopsis        string  `json:"arabicSynopsis"`
	ProductionYear        int     `json:"productionYear"`
	ProductionHouse       string  `json:"productionHouse"`
	AgeGroup              int     `json:"ageGroup"`
	IntroDuration         *string `json:"introDuration"`
	IntroStart            *string `json:"introStart"`
	OutroDuration         *string `json:"outroDuration"`
	OutroStart            *string `json:"outroStart"`
	ProductionCountries   []int   `json:"productionCountries"`
}
type RightProduct struct {
	ProductName int `json:"productName"`
}
type SeasonGenres struct {
	GenreId string `json:"genreId"`
	Id      string `json:"id"`
}
type SeasonSubgenre struct {
	SubgenreId string `jso:"SubgenreId"`
}
type NewSeasonGenres struct {
	GenreId     string   `json:"genreId"`
	SubgenresId []string `json:"subgenresId"`
	Id          string   `json:"id"`
}
type SubscriptionPlans struct {
	SubscriptionPlanId int `json:"subscriptionPlanId"`
}
type ProductionCountry struct {
	CountryId int `json:"countryId"`
}
type VarianceTrailers struct {
	Order                 int    `json:"order"`
	VideoTrailerId        string `json:"videoTrailerId"`
	EnglishTitle          string `json:"englishTitle"`
	ArabicTitle           string `json:"arabicTitle"`
	Duration              int    `json:"duration"`
	HasTrailerPosterImage bool   `json:"hasTrailerPosterImage"`
	TrailerposterImage    string `json:"trailerposterImage,omitempty"`
	Id                    string `json:"id"`
}
type EpisodeDetailsByseasonId struct {
	IsPrimary              bool        `json:"isPrimary"`
	UserId                 string      `json:"userId"`
	SecondarySeasonId      string      `json:"secondarySeasonId" `
	VarianceIds            []int       `json:"varianceIds"`
	EpisodeIds             []int       `json:"episodeIds"`
	SecondaryEpisodeId     string      `json:"secondaryEpisodeId"`
	ContentId              string      `json:"contentId"`
	EpisodeKey             int         `json:"episodeKey"`
	SeasonId               string      `json:"seasonId"`
	Status                 int         `json:"status"`
	StatusCanBeChanged     bool        `json:"statusCanBeChanged"`
	SubStatus              int         `json:"subStatus"`
	SubStatusName          string      `json:"subStatusName"`
	DigitalRightsType      int         `json:"digitalRightsType"`
	DigitalRightsStartDate time.Time   `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   time.Time   `json:"digitalRightsEndDate"`
	CreatedBy              string      `json:"createdBy"`
	PrimaryInfo            PrimaryInfo `json:"primaryInfo"`
	Cast                   Cast        `json:"cast"`
	Music                  Music       `json:"music"`
	TagInfo                TagInfo     `json:"tagInfo"`
	NonTextualData         []int       `json:"nonTextualData"`
	Translation            Translation `json:"translation"`
	SchedulingDateTime     *time.Time  `json:"schedulingDateTime"`
	PublishingPlatforms    []int       `json:"publishingPlatforms"`
	SeoDetails             []int       `json:"seoDetails"`
	Id                     string      `json:"id"`
}

/*Fetch one-tier content*/
type OnetireContent struct {
	ContentKey     int            `json:"contentKey"`
	Duration       *int           `json:"duration"`
	Status         int            `json:"status"`
	TextualData    TextualData    `json:"textualData"`
	NonTextualData NonTextualData `json:"nonTextualData"`
	Id             string         `json:"id"`
}
type TextualData struct {
	PrimaryInfo      PrimaryInfo        `json:"primaryInfo"`
	ContentGenres    []interface{}      `json:"contentGenres"`
	ContentVariances []ContentVariances `json:"contentVariances"`
	Cast             Cast               `json:"cast"`
	Music            Music              `json:"music"`
	TagInfo          TagInfo            `json:"tagInfo"`
	AboutTheContent  AboutTheContent    `json:"aboutTheContent"`
	SeoDetails       SeoDetails         `json:"seoDetails"`
}
type ContentVariances struct {
	Status                 int                `json:"status"`
	StatusCanBeChanged     bool               `json:"statusCanBeChanged"`
	SubStatusName          *string            `json:"subStatusName"`
	VideoContentId         string             `json:"videoContentId"`
	LanguageType           string             `json:"languageType"`
	OverlayPosterImage     string             `json:"overlayPosterImage"`
	DubbingScript          string             `json:"dubbingScript"`
	SubtitlingScript       string             `json:"subtitlingScript"`
	DubbingLanguage        *string            `json:"dubbingLanguage"`
	DubbingDialectId       *int               `json:"dubbingDialectId"`
	SubtitlingLanguage     *string            `json:"subtitlingLanguage"`
	DigitalRightsType      int                `json:"digitalRightsType"`
	DigitalRightsStartDate *time.Time         `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   *time.Time         `json:"digitalRightsEndDate"`
	DigitalRightsRegions   []int              `json:"digitalRightsRegions"`
	SchedulingDateTime     *time.Time         `json:"schedulingDateTime"`
	CreatedBy              *string            `json:"createdBy"`
	PublishingPlatforms    []int              `json:"publishingPlatforms"`
	Products               []int              `json:"products"`
	SubscriptionPlans      []int              `json:"subscriptionPlans"`
	CountryCheck           bool               `json:"countryCheck"` /*pending*/
	IntroDuration          string             `json:"introDuration"`
	IntroStart             string             `json:"introStart"`
	VarianceTrailers       []VarianceTrailers `json:"varianceTrailers"`
	Id                     string             `json:"id"`
}

type ContentGenre struct {
	GenreId string `json:"genreId"`
	Id      string `json:"id"`
	// Order     int
}

type ContentSubgenre struct {
	// ContentGenreId string
	SubgenreId string
	// Order          int
}
type ContentGeneresQueryDetails struct {
	GenreId       string `json:"genre_id"`
	SubgenresId   string `json:"subgenres_id"`
	SubGenreOrder string `json:"sub_genre_order"`
	Id            string `json:"id"`
}
type ContentGenres struct {
	GenreId     string   `json:"genreId"`
	SubgenresId []string `json:"subgenresId"`
	Id          string   `json:"id,omitempty"`
}

func JsonStringToStringSliceOrMap(data string) ([]string, error) {
	output := make([]string, 1000)
	err := json.Unmarshal([]byte(data), &output)
	if err != nil {
		return nil, err
	}
	//sort.Strings(output)
	return output, nil
}
