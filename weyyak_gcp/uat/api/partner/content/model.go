package content

import (
	"time"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserInfo struct {
	Email string `json:"email"`
}

type Pagination struct {
	Size   int `json:"size"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

/*Fetch one-tier content*/
type OnetierContent struct {
	ContentKey       int             `json:"contentKey"`
	PrimaryInfo      PrimaryInfo     `json:"primaryInfo"`
	ContentGenres    []interface{}   `json:"contentGenres"`
	ContentVariances []interface{}   `json:"contentVariances"`
	Cast             Cast            `json:"cast"`
	Music            Music           `json:"music"`
	TagInfo          TagInfo         `json:"tagInfo"`
	AboutTheContent  AboutTheContent `json:"aboutTheContent"`
	SeoDetails       SeoDetails      `json:"seoDetails"`
	NonTextualData   NonTextualData  `json:"nonTextualData"`
	CreatedAt        time.Time       `json:"createdAt"`
	// InsertedAt       time.Time       `json:"insertedAt,omitempty"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Id         string    `json:"id"`
}
type AllOnetierContent struct {
	ContentKey       int                      `json:"contentKey"`
	ContentVariances []map[string]interface{} `json:"contentVariances"`
	PrimaryInfo      PrimaryInfo              `json:"primaryInfo"`
	ContentGenres    []interface{}            `json:"contentGenres"`
	Cast             Cast                     `json:"cast"`
	Music            Music                    `json:"music"`
	TagInfo          TagInfo                  `json:"tagInfo"`
	AboutTheContent  AboutTheContent          `json:"aboutTheContent"`
	SeoDetails       SeoDetails               `json:"seoDetails"`
	NonTextualData   NonTextualData           `json:"nonTextualData"`
	CreatedAt        time.Time                `json:"createdAt"`
	ModifiedAt       time.Time                `json:"modifiedAt"`
	Id               string                   `json:"id"`
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

type FinalSeasonResultOneTire struct {
	Id                         string    `json:"id"`
	ContentKey                 int       `json:"content_key"`
	PrimaryInfoId              string    `json:"primary_info_id"`
	ContentType                string    `json:"content_type"`
	OriginalTitle              string    `json:"original_title"`
	AlternativeTitle           string    `json:"alternative_title"`
	ArabicTitle                string    `json:"arabic_title"`
	TransliteratedTitle        string    `json:"transliterated_title"`
	Notes                      string    `json:"notes"`
	CastId                     string    `json:"cast_id"`
	MusicId                    string    `json:"music_id"`
	TagInfoId                  string    `json:"tag_info_id"`
	OriginalLanguage           string    `json:"original_language"`
	Supplier                   string    `json:"supplier"`
	AcquisitionDepartment      string    `json:"acquisition_department"`
	EnglishSynopsis            string    `json:"english_synopsis"`
	ArabicSynopsis             string    `json:"arabic_synopsis"`
	ProductionYear             int       `json:"production_year"`
	ProductionHouse            string    `json:"production_house"`
	AgeGroup                   int       `json:"age_group"`
	AboutTheContentInfoId      string    `json:"about_the_content_info_id"`
	EnglishMetaTitle           string    `json:"english_meta_title"`
	ArabicMetaTitle            string    `json:"arabic_meta_title"`
	EnglishMetaDescription     string    `json:"english_meta_description"`
	ArabicMetaDescription      string    `json:"arabic_meta_description"`
	HasPosterImage             bool      `json:"has_poster_image"`
	HasDetailsBackground       bool      `json:"has_details_background"`
	HasMobileDetailsBackground bool      `json:"has_mobile_details_background"`
	ModifiedAt                 time.Time `json:"modified_at"`
	InsertedAt                 time.Time `json:"insertedAt,omitempty"`
	CreatedAt                  time.Time `json:"created_at"`
}

type ContentVariancesSource struct {
	Id                    string `json:"id"`
	Length                int    `json:"length"`
	VideoContentId        string `json:"video_content_id"`
	LanguageType          int    `json:"languageType"`
	HasDubbingScript      bool   `json:"has_dubbing_script"`
	HasSubtitlingScript   bool   `json:"has_subtitling_script"`
	VarianceId            int    `json:"varianceId"`
	DubbingLanguage       string `json:"dubbing_language"`
	DubbingDialectId      int    `json:"dubbing_dialect_id"`
	RightsId              string `json:"rights_id"`
	TagInfoId             string `json:"tagInfoId"`
	HasOverlayPosterImage bool   `json:"has_overlay_poster_image"`
}

/*fetch multi-tier content*/
type MultiTierContent struct {
	ContentKey     int                `json:"contentKey"`
	PrimaryInfo    ContentPrimaryInfo `json:"primaryInfo"`
	ContentGenres  []interface{}      `json:"contentGenres"`
	ContentSeasons [1]ContentSeasons  `json:"contentSeasons"`
	SeoDetails     SeoDetails         `json:"seoDetails"`
	CreatedAt      time.Time          `json:"createdAt"`
	ModifiedAt     time.Time          `json:"modifiedAt"`
	ContentId      string             `json:"id"` //content ID
}

type MultiTierContent2 struct {
	ContentKey     int                      `json:"contentKey"`
	PrimaryInfo    ContentPrimaryInfo       `json:"primaryInfo"`
	ContentGenres  []MultiTierContentGenres `json:"contentGenres"`
	ContentSeasons []ContentSeasons2         `json:"contentSeasons"`
	SeoDetails     SeoDetails               `json:"seoDetails"`
	CreatedAt      time.Time                `json:"createdAt"`
	ModifiedAt     time.Time                `json:"modifiedAt"`
	Id             string                   `json:"id"` //content ID
}
type AllMultiTierContent struct {
	ContentKey     int                 `json:"contentKey"`
	ContentSeasons []AllContentSeasons `json:"contentSeasons"`
	PrimaryInfo    ContentPrimaryInfo  `json:"primaryInfo"`
	ContentGenres  []interface{}       `json:"contentGenres"`
	SeoDetails     SeoDetails          `json:"seoDetails"`
	CreatedAt      time.Time           `json:"createdAt"`
	ModifiedAt     time.Time           `json:"modifiedAt"`
	ContentId      string              `json:"id"` //content ID
}
type ContentPrimaryInfo struct {
	ContentType         string `json:"contentType"`
	OriginalTitle       string `json:"originalTitle"`
	AlternativeTitle    string `json:"alternativeTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	TransliteratedTitle string `json:"transliteratedTitle"`
	Notes               string `json:"notes"`
}
type SeasonPrimaryInfo struct {
	SeasonNumber        int    `json:"seasonNumber"`
	OriginalTitle       string `json:"originalTitle"`
	AlternativeTitle    string `json:"alternativeTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	TransliteratedTitle string `json:"transliteratedTitle"`
	Notes               string `json:"notes"`
}
type EpisodeDetailsMultiTier struct {
	EpisodeNumber   int            `json:"episdoeNumber"`
	EpisodeKey      int            `json:"episodeKey"`
	Length          int            `json:"length"`
	VideoContentId  string         `json:"videoContentId,omitempty"`
	VideoContentUrl string         `json:"videoContentUrl"`
	SynopsisEnglish string         `json:"synopsisEnglish"`
	SynopsisArabic  string         `json:"synopsisArabic"`
	NonTextualData  NonTextualData `json:"nonTextualData"`
	EpisodeId       string         `json:"id"`
}
type FetchEpisodeDetailsMultiTier struct {
	EpisodeNumber       int            `json:"episdoeNumber"`
	EpisodeKey          int            `json:"episodeKey"`
	Length              int            `json:"length"`
	VideoContentId      string         `json:"videoContentId,omitempty"`
	VideoContentUrl     string         `json:"videoContentUrl"`
	SynopsisEnglish     string         `json:"synopsisEnglish"`
	SynopsisArabic      string         `json:"synopsisArabic"`
	NonTextualData      NonTextualData `json:"nonTextualData"`
	EpisodeId           string         `json:"id"`
	HasPosterImage      bool           `json:"hasPosterImage"`
	HasDubbingScript    bool           `json:"hasDubbingScript"`
	HasSubtitlingScript bool           `json:"hasSubtitlingScript"`
}

type ContentSeasons struct {
	ContentId             string                `json:"contentId"`
	SeasonKey             int                   `json:"seasonKey"`
	SeasonNumber          int                   `json:"seasonNumber"`
	CreatedAt             time.Time             `json:"createdAt"`
	ModifiedAt            time.Time             `json:"modifiedAt"`
	PrimaryInfo           SeasonPrimaryInfo     `json:"primaryInfo"`
	Cast                  Cast                  `json:"cast"`
	Music                 Music                 `json:"music"`
	TagInfo               TagInfo               `json:"tagInfo"`
	SeasonGenres          []interface{}         `json:"seasonGeneres"`
	TrailerInfo           []interface{}         `json:"trailersInfo"`
	AboutTheContent       AboutTheContent       `json:"aboutTheContent"`
	Translation           Translation           `json:"translation"`
	EpisodeResult         []interface{}         `json:"episodes"`
	ContentNonTextualData ContentNonTextualData `json:"nonTextualData"`
	DigitalRightsRegions  []int                 `json:"digitalRightsRegions"`
	SeasonId              string                `json:"id"` //season ID
}
type ContentSeasons2 struct {
	ContentId             string                   `json:"contentId"`
	SeasonKey             int                      `json:"seasonKey"`
	SeasonNumber          int                      `json:"seasonNumber"`
	CreatedAt             time.Time                `json:"createdAt"`
	ModifiedAt            time.Time                `json:"modifiedAt"`
	PrimaryInfo           SeasonPrimaryInfo        `json:"primaryInfo"`
	Cast                  Cast                     `json:"cast"`
	Music                 Music                    `json:"music"`
	TagInfo               TagInfo                  `json:"tagInfo"`
	SeasonGenres          []MultiTierContentGenres `json:"seasonGenres"`
	TrailerInfo           []interface{}            `json:"trailersInfo"`
	AboutTheContent       AboutTheContent          `json:"aboutTheContent"`
	Translation           Translation2              `json:"translation"`
	EpisodeResult         []ContentEpisode         `json:"episodes"`
	ContentNonTextualData ContentNonTextualData    `json:"nonTextualData"`
	DigitalRightsRegions  []int                    `json:"digitalRightsRegions"`
	SeasonId              string                   `json:"id"` //season ID
}
type Rights struct {
	DigitalRightsRegions []int `json:"digitalRightsRegion"`
}

/*for get all multi tire content*/
type AllContentSeasons struct {
	ContentId             string                `json:"contentId"`
	SeasonKey             int                   `json:"seasonKey"`
	SeasonNumber          int                   `json:"seasonNumber"`
	CreatedAt             time.Time             `json:"createdAt"`
	ModifiedAt            time.Time             `json:"modifiedAt"`
	PrimaryInfo           SeasonPrimaryInfo     `json:"primaryInfo"`
	TrailerInfo           []interface{}         `json:"trailersInfo"`
	Cast                  Cast                  `json:"cast"`
	Music                 Music                 `json:"music"`
	TagInfo               TagInfo               `json:"tagInfo"`
	AboutTheContent       AboutTheContent       `json:"aboutTheContent"`
	Translation           Translation           `json:"translation"`
	ContentNonTextualData ContentNonTextualData `json:"nonTextualData"`
	Rights                Rights                `json:"rights"`
	SeasonId              string                `json:"id"` //season ID
}

/*fetch multi tire all*/
type AllMultiTier struct {
	ContentKey     string             `json:"contentKey"`
	ContentSeasons AllContentSeasons  `json:"contentSeasons"`
	PrimaryInfo    ContentPrimaryInfo `json:"primaryInfo"`
	ContentGenres  []interface{}      `json:"contentGenres"`
	SeoDetails     SeoDetails         `json:"seoDetails"`
	CreatedAt      time.Time          `json:"createdAt"`
	ModifiedAt     time.Time          `json:"modifiedAt"`
	ContentId      string             `json:"id"`
}

type ContentVariances struct {
	VideoContentUrl      string  `json:"videoContentUrl"`
	Length               int     `json:"length"`
	LanguageType         string  `json:"languageType"`
	DubbingScript        string  `json:"dubbingScript"`
	SubtitlingScript     string  `json:"subtitlingScript"`
	DubbingLanguage      *string `json:"dubbingLanguage"`
	DubbingDialectId     *int    `json:"dubbingDialectId ,omitempty"`
	DubbingDialectName   string  `json:"dubbingDialectName"`
	SubtitlingLanguage   string  `json:"subtitlingLanguage,omitempty"`
	DigitalRightsRegions []int   `json:"digitalRightsRegions"`
	VarianceId           int     `json:"varianceId,omitempty"`
	Id                   string  `json:"id"`
}

type PrimaryInfo struct {
	SeasonNumber        int    `json:"seasonNumber,omitempty"`
	ContentTier         int    `json:"contentTier,omitempty"`
	Number              int    `json:"number,omitempty"`
	VideoContentId      string `json:"videoContentId,omitempty"`
	SynopsisEnglish     string `json:"synopsisEnglish,omitempty"`
	SynopsisArabic      string `json:"synopsisArabic,omitempty"`
	ContentType         string `json:"contentType,omitempty"`
	OriginalTitle       string `json:"originalTitle"`
	AlternativeTitle    string `json:"alternativeTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	TransliteratedTitle string `json:"transliteratedTitle"`
	Notes               string `json:"notes"`
}
type Cast struct {
	CastId             string   `json:"castId"`
	MainActorId        string   `json:"mainActorId"`
	MainActressId      string   `json:"mainActressId"`
	MainActorEnglish   string   `json:"mainActorEnglish"`
	MainActorArabic    string   `json:"mainActorArabic"`
	MainActressEnglish string   `json:"mainActressEnglish"`
	MainActressArabic  string   `json:"mainActressArabic"`
	ActorIds           []string `json:"actorIds"`
	ActorEnglish       []string `json:"actorEnglish"`
	ActorArabic        []string `json:"actorArabic"`
	WriterId           []string `json:"writerIds"`
	WriterEnglish      []string `json:"writerEnglish"`
	WriterArabic       []string `json:"writerArabic"`
	DirectorIds        []string `json:"directorIds"`
	DirectorEnglish    []string `json:"directorEnglish"`
	DirectorArabic     []string `json:"directorArabic"`
}
type ContentActor struct {
	ActorId         string `json:"actorId"`
	ActorEnglish    string `json:"actorEnglish"`
	ActorArabic     string `json:"actorArabic"`
	WriterId        string `json:"writerId"`
	WriterEnglish   string `json:"writerEnglish"`
	WriterArabic    string `json:"writerArabic"`
	DirectorId      string `json:"directorId"`
	DirectorEnglish string `json:"directorEnglish"`
	DirectorArabic  string `json:"directorArabic"`
}
type Music struct {
	MusicId               string   `json:"musicId"`
	SingerIds             []string `json:"singerIds"`
	SingersEnglish        []string `json:"singersEnglish"`
	SingersArabic         []string `json:"singersArabic"`
	MusicComposerIds      []string `json:"musicComposerIds"`
	MusicComposersEnglish []string `json:"musicComposersEnglish"`
	MusicComposersArabic  []string `json:"musicComposersArabic"`
	SongWriterIds         []string `json:"songWriterIds"`
	SongWritersEnglish    []string `json:"songWritersEnglish"`
	SongWritersArabic     []string `json:"songWritersArabic"`
}
type ContentMusic struct {
	SingerIds             string `json:"actorEnglish"`
	SingersEnglish        string `json:"singersEnglish"`
	SingersArabic         string `json:"singersArabic"`
	MusicComposerIds      string `json:"musicComposerIds"`
	MusicComposersEnglish string `json:"musicComposersEnglish"`
	MusicComposersArabic  string `json:"musicComposersArabic"`
	SongWriterIds         string `json:"songWriterIds"`
	SongWritersEnglish    string `json:"songWritersEnglish"`
	SongWritersArabic     string `json:"songWritersArabic"`
}
type TagInfo struct {
	Tags []string `json:"tag"`
}
type AboutTheContent struct {
	OriginalLanguage      string `json:"originalLanguage"`
	Supplier              string `json:"supplier"`
	AcquisitionDepartment string `json:"acquisitionDepartment"`
	EnglishSynopsis       string `json:"englishSynopsis"`
	ArabicSynopsis        string `json:"arabicSynopsis"`
	ProductionYear        int    `json:"productionYear"`
	ProductionHouse       string `json:"productionHouse"`
	AgeGroup              string `json:"ageGroups"`
	ProductionCountries   []int  `json:"productionCountries"`
}
type SeoDetails struct {
	EnglishMetaTitle       string `json:"englishMetaTitle"`
	ArabicMetaTitle        string `json:"arabicMetaTitle"`
	EnglishMetaDescription string `json:"englishMetaDescription"`
	ArabicMetaDescription  string `json:"arabicMetaDescription"`
}
type DigitalRightsRegions struct {
	CountryId int `json:"country_Id"`
}
type FinalSeasonResult struct {
	SeasonNumber               int       `json:"seasonNumber"`
	SeasonKey                  int       `json:"seasonKey"`
	MultiTierContentKey        int       `json:"multiTiercontentKey"`
	ModifiedAt                 time.Time `json:"modifiedAt"`
	InsertedAt                 time.Time `json:"insertedAt,omitempty"`
	OriginalTitle              string    `json:"originalTitle"`
	AlternativeTitle           string    `json:"alternativeTitle"`
	ArabicTitle                string    `json:"arabicTitle"`
	TransliteratedTitle        string    `json:"transliteratedTitle"`
	Notes                      string    `json:"notes"`
	LanguageType               int       `json:"languageType"`
	MultiTierLanguageType      int       `json:"multiTierLanguageType,omitempty"`
	DubbingLanguage            *string   `json:"dubbingLanguage"`
	DubbingDialectId           int       `json:"dubbingDialectId"`
	DubbingDialectName         *string   `json:"dubbingDialectName"`
	SubtitlingLanguage         *string   `json:"subtitlingLanguage"`
	Id                         string    `json:"id"`
	RightsId                   string    `json:"rightsId,omitempty"`
	CastId                     string    `json:"castId,omitempty"`
	MusicId                    string    `json:"musicId,omitempty"`
	TagInfoId                  string    `json:"tagInfoId,omitempty"`
	OriginalLanguage           string    `json:"originalLanguage,omitempty"`
	Supplier                   string    `json:"supplier,omitempty"`
	AcquisitionDepartment      string    `json:"acquisitionDepartment"`
	EnglishSynopsis            string    `json:"englishSynopsis,omitempty"`
	ArabicSynopsis             string    `json:"arabicSynopsis,omitempty"`
	ProductionYear             int       `json:"productionYear,omitempty"`
	ProductionHouse            string    `json:"productionHouse,omitempty"`
	AgeGroup                   int       `json:"ageGroups,omitempty"`
	AboutTheContentInfoId      string    `json:"aboutTheContentInfoId,omitempty"`
	ContentTier                int       `json:"contentTier"`
	ContentKey                 int       `json:"contentKey,omitempty"`
	Duration                   int       `json:"duration,omitempty"`
	ContentType                string    `json:"contentType,omitempty"`
	VideoContentId             string    `json:"VideoContentId,omitempty"`
	HasOverlayPosterImage      bool      `json:"hasOverlayPosterImage,omitempty"`
	HasDubbingScript           bool      `json:"hasDubbingScript,omitempty"`
	HasSubtitlingScript        bool      `json:"hasSubtitlingScript,omitempty"`
	VarianceId                 string    `json:"varianceId"`
	CreatedByUserId            *string   `json:"createdByUserId"`
	PlaybackItemId             string    `json:"playbackItemId"`
	HasPosterImage             bool      `json:"hasPosterImage"`
	HasDetailsBackground       bool      `json:"hasDetailsBackground"`
	HasMobileDetailsBackground bool      `json:"hasMobileDetailsBackground"`
	CreatedAt                  time.Time `json:"createdAt"`
	EpisodeNumber              int       `json:"episodeNumber,omitempty"`
	ContentId                  string    `json:"contentId,omitempty"`
	EpisodeKey                 int       `json:"episodeKey,omitempty"`
	EpisodeLength              int       `json:"episodeLength"`
	SeasonId                   string    `json:"seasonId,omitempty"`
	SynopsisEnglish            string    `json:"synopsisEnglish"`
	SynopsisArabic             string    `json:"synopsisArabic"`
	EnglishMetaTitle           string    `json:"englishMetaTitle"`
	ArabicMetaTitle            string    `json:"arabicMetaTitle"`
	EnglishMetaDescription     string    `json:"englishMetaDescription"`
	ArabicMetaDescription      string    `json:"arabicMetaDescription"`
	SeasonOriginalTitle        string    `json:"seasonOriginalTitle"`
	SeasonAlternativeTitle     string    `json:"seasonAlternativeTitle"`
	SeasonArabicTitle          string    `json:"seasonArabicName"`
	SeasonTransliteratedTitle  string    `json:"seasonTransliteratedTitle"`
	SeasonNotes                string    `json:"seasonNotes"`
	SeasonMainActorId          *string   `json:"seasonMainActorId"`
	SeasonMainActressId        *string   `json:"seasonMainActressId"`
	SeasonMainActorEnglish     *string   `json:"seasonMainActorEnglish"`
	SeasonMainActorArabic      *string   `json:"seasonMainActorArabic"`
	SeasonMainActressEnglish   *string   `json:"seasonMainActressEnglish"`
	SeasonMainActressArabic    *string   `json:"seasonMainActressArabic"`
	SeasonGenreId              string    `json:"seasonGenreId"`
}
type NonTextualData struct {
	PosterImage             string  `json:"posterImage"`
	DubbingScript           *string `json:"dubbingScript,omitempty"`
	SubtitlingScript        *string `json:"subtitlingScript,omitempty"`
	OverlayPosterImage      string  `json:"overlayPosterImage,omitempty"`
	DetailsBackground       string  `json:"detailsBackground,omitempty"`
	MobileDetailsBackground string  `json:"mobileDetailsBackground,omitempty"`
}
type NowTextualDataEpisode struct {
	PosterImage             string `json:"posterImage"`
	DubbingScript           string `json:"dubbingScript"`
	SubtitlingScript        string `json:"subtitlingScript"`
	OverlayPosterImage      string `json:"overlayPosterImage,omitempty"`
	DetailsBackground       string `json:"detailsBackground,omitempty"`
	MobileDetailsBackground string `json:"mobileDetailsBackground,omitempty"`
}
type ContentNonTextualData struct {
	PosterImage             string `json:"posterImage"`
	OverlayPosterImage      string `json:"overlayPosterImage,omitempty"`
	DetailsBackground       string `json:"detailsBackground,omitempty"`
	MobileDetailsBackground string `json:"mobileDetailsBackground,omitempty"`
}
type ContentRightsCountry struct {
	Id              string `json:"id"`
	ContentRightsId string `json:"contentRightsId"`
	CountryId       int    `json:"countryId"`
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
	Name string `json:"name"`
}
type SeasonGenres struct {
	GenerEnglishName string `json:"generEnglishName"`
	GenerArabicName  string `json:"generArabicName"`
	Id               string `json:"id"`
}
type SeasonSubgenre struct {
	SubGenerEnglish string `json:"subGenerEnglish"`
	SubGenerArabic  string `json:"subGenerArabic"`
}
type NewSeasonGenres struct {
	GenerEnglishName string   `json:"generEnglishName"`
	GenerArabicName  string   `json:"generArabicName"`
	SubGenerEnglish  []string `json:"subGenerEnglish"`
	SubGenerArabic   []string `json:"subGenerArabic"`
	Id               string   `json:"id"`
}
type ProductionCountry struct {
	CountryId int `json:"countryId"`
}

/*Fetch episode Details By episode key*/
type EpisodeResult struct {
	EpisodeNumber         int                   `json:"episodeNumber"`
	ContentKey            string                `json:"contentId"`
	EpisodeKey            int                   `json:"episodeKey"`
	Length                int                   `json:"length"`
	VideoContentUrl       string                `json:"videoContentUrl"`
	SynopsisEnglish       string                `json:"synopsisEnglish"`
	SynopsisArabic        string                `json:"synopsisArabic"`
	SeasonId              string                `json:"seasonId"`
	DigitalRightsRegions  []int                 `json:"digitalRightsRegions"`
	PrimaryInfo           PrimaryInfo           `json:"primaryInfo"`
	Cast                  Cast                  `json:"cast"`
	Music                 Music                 `json:"music"`
	TagInfo               TagInfo               `json:"tagInfo"`
	NonTextualDataEpisode NowTextualDataEpisode `json:"nonTextualData"`
	Translation           Translation           `json:"translation"`
	SeoDetails            SeoDetails            `json:"seoDetails"`
	CreatedAt             time.Time             `json:"createdAt"`
	ModifiedAt            time.Time             `json:"modifiedAt"`
	Id                    string                `json:"id"`
}
type Translation struct {
	LanguageType       string  `json:"languageType"`
	DubbingLanguage    *string `json:"dubbingLanguage"`
	DubbingDialectId   *int    `json:"dubbingDialectId"`
	SubtitlingLanguage *string `json:"subtitlingLanguage"`
}

type Translation2 struct {
	LanguageType       string  `json:"languageType"`
	DubbingLanguage    *string `json:"dubbingLanguage"`
	DubbingDialectId   int     `json:"dubbingDialectId"`
	DubbingDialectName string  `json:"dubbingDialectName"`
	SubtitlingLanguage *string `json:"subtitlingLanguage"`
}
type MenuDetails struct {
	Total       int        `json:"total"`
	PerPage     int        `json:"per_page"`
	CurrentPage int        `json:"current_page"`
	LastPage    int        `json:"last_page"`
	NextPageUrl *string    `json:"next_page_url"`
	PrevPageUrl *string    `json:"prev_page_url"`
	From        int        `json:"from"`
	To          int        `json:"to"`
	Data        []MenuData `json:"data"`
}
type MenuData struct {
	Id                 int              `json:"id"`
	FriendlyUrlEnglish string           `json:"friendly_url_english"`
	FriendlyUrlArabic  string           `json:"friendly_url_arabic"`
	SeoDescription     string           `json:"seo_description"`
	TitleEnglish       string           `json:"title_english"`
	TitleArabic        string           `json:"title_arabic"`
	Type               string           `json:"type"`
	Featured           *FeaturedDetails `json:"featured,omitempty"`
	Playlists          []MenuPlaylists  `json:"playlists,omitempty"`
}
type MenuDatas struct {
	Id                 int              `json:"id"`
	FriendlyUrlEnglish string           `json:"friendly_url_english"`
	FriendlyUrlArabic  string           `json:"friendly_url_arabic"`
	SeoDescription     string           `json:"seo_description"`
	TitleEnglish       string           `json:"title_english"`
	TitleArabic        string           `json:"title_arabic"`
	Type               string           `json:"type"`
	Featured           *FeaturedDetails `json:"featured"`
	Playlists          []MenuPlaylists  `json:"playlists"`
}
type PageDetails struct {
	Id                     string `json:"id"`
	ThirdPartyPageKey      int    `json:"third_party_page_key"`
	PageKey                int    `json:"page_key"`
	EnglishPageFriendlyUrl string `json:"english_page_friendly_url"`
	ArabicPageFriendlyUrl  string `json:"arabic_page_friendly_url"`
	EnglishMetaDescription string `json:"english_meta_description"`
	ArabicMetaDescription  string `json:"arabic_meta_description"`
	EnglishTitle           string `json:"english_title"`
	ArabicTitle            string `json:"arabic_title"`
	PageType               int    `json:"page_type"`
}
type FeaturedDetails struct {
	ID        int64               `json:"id"`
	Type      string              `json:"type"`
	Playlists []FeaturedPlaylists `json:"playlists"`
}

// FeaturedPlaylists struct for DB binding
type FeaturedPlaylists struct {
	ID           int32             `json:"id"`
	PlaylistType string            `json:"playlist_type"`
	Content      []PlaylistContent `json:"content"`
}

// PlaylistContent struct for DB binding
type PlaylistContent struct {
	ContentId             string                `json:"content_id"`
	ContentKey            int                   `json:"content_key"`
	AgeRating             string                `json:"age_rating"`
	VideoId               string                `json:"video_id"`
	FriendlyUrl           string                `json:"friendly_url"`
	ContentType           string                `json:"content_type"`
	SynopsisEnglish       string                `json:"synopsis_english"`
	SynopsisArabic        string                `json:"synopsis_arabic"`
	SeoTitleEnglish       string                `json:"seo_title_english"`
	SeoTitleArabic        string                `json:"seo_title_arabic"`
	SeoDescriptionEnglish string                `json:"seo_description_english"`
	SeoDescriptionArabic  string                `json:"seo_description_arabic"`
	Length                *int32                `json:"length"`
	TitleEnglish          string                `json:"title_english"`
	TitleArabic           string                `json:"title_arabic"`
	SeoTitle              string                `json:"seo_title"`
	Imagery               ContentImageryDetails `json:"imagery"`
	InsertedAt            time.Time             `json:"inserted_at"`
	ModifiedAt            time.Time             `json:"modified_at"`
}

// ContentImageryDetails for DB binding
type ContentImageryDetails struct {
	Thumbnail     string `json:"thumbnail"`
	Backdrop      string `json:"backdrop"`
	MobileImg     string `json:"mobile_img"`
	FeaturedImg   string `json:"featured_img"`
	OverlayPoster string `json:"overlayPoster"`
}

// MenuPlaylists struct for DB binding
type MenuPlaylists struct {
	ID           int32             `json:"id"`
	TitleEnglish string            `json:"title_english"`
	TitleArabic  string            `json:"title_arabic"`
	Content      []PlaylistContent `json:"content"`
}
type ContentDetails struct {
	Id                    string                `json:"id"`
	ContentKey            int                   `json:"content_key"`
	AgeRating             string                `json:"age_rating"`
	VideoId               string                `json:"video_id"`
	FriendlyUrl           string                `json:"friendly_url"`
	ContentType           string                `json:"content_type"`
	ContentTier           int                   `json:"content_tier"`
	SynopsisEnglish       string                `json:"synopsis_english"`
	SynopsisArabic        string                `json:"synopsis_arabic"`
	SeoTitleEnglish       string                `json:"seo_title_english"`
	SeoTitleArabic        string                `json:"seo_title_arabic"`
	SeoDescriptionEnglish string                `json:"seo_description_english"`
	SeoDescriptionArabic  string                `json:"seo_description_arabic"`
	Length                *int32                `json:"length"`
	TitleEnglish          string                `json:"title_english"`
	TitleArabic           string                `json:"title_arabic"`
	SeoTitle              string                `json:"seo_title"`
	Imagery               ContentImageryDetails `json:"imagery"`
	InsertedAt            time.Time             `json:"inserted_at,omitempty"`
	ModifiedAt            time.Time             `json:"modified_at"`
	VarienceId            string                `json:"varience_id"`
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

// Playlist - binding for db
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
type FunctionResponse struct {
	ContentDetails []PlaylistContent
	Err            error
}

type Data struct {
	Pagination Pagination          `json:"pagination"`
	data       []AllOnetierContent `json:"data"`
}
type Response struct {
	Pagination Pagination          `json:"pagination"`
	Data       []AllOnetierContent `json:"data"`
}
type ErrorResponse struct {
	Description string `json:"description"`
	Code        string `json:"code"`
	Error       string `json:"error"`
	Invalid     struct {
		AdditionalProp1 struct {
			Description string `json:"description"`
			Code        string `json:"code"`
		} `json:"additionalProp1"`
		AdditionalProp2 struct {
			Description string `json:"description"`
			Code        string `json:"code"`
		} `json:"additionalProp2"`
		AdditionalProp3 struct {
			Description string `json:"description"`
			Code        string `json:"code"`
		} `json:"additionalProp3"`
	} `json:"invalid"`
}

type FinalSeasonResultContentOneTire struct {
	MultiTierContentKey    int       `json:"multi_tier_content_key"`
	ContentType            string    `json:"content_type"`
	OriginalTitle          string    `json:"original_title"`
	AlternativeTitle       string    `json:"alternative_title"`
	ArabicTitle            string    `json:"arabic_title"`
	TransliteratedTitle    string    `json:"transliterated_title"`
	EnglishMetaTitle       string    `json:"english_meta_title"`
	ArabicMetaTitle        string    `json:"arabic_meta_title"`
	EnglishMetaDescription string    `json:"english_meta_description"`
	ArabicMetaDescription  string    `json:"arabic_meta_description"`
	Notes                  string    `json:"notes"`
	ContentId              string    `json:"content_id"`
	ModifiedAt             time.Time `json:"modified_at"`
	CreatedAt              time.Time `json:"created_at"`
}

type MultiTierContentGenres struct {
	GenerEnglishName string   `json:"generEnglishName"`
	GenerArabicName  string   `json:"generArabicName"`
	SubGenerEnglish  []string `json:"subGenerEnglish"`
	SubGenerArabic   []string `json:"subGenerArabic"`
	Id               string   `json:"id"`
}

type Seasons struct {
	Id                         string    `json:"id"`
	ContentId                  string    `json:"content_id"`
	Status                     int       `json:"status"`
	ModifiedAt                 time.Time `json:"modified_at"`
	CreatedByUserId            string    `json:"created_by_user_id"`
	PrimaryInfoId              string    `json:"primary_info_id"`
	Number                     int       `json:"number"`
	TranslationId              string    `json:"translation_id"`
	AboutTheContentInfoId      string    `json:"about_the_content_info_id"`
	HasPosterImage             bool      `json:"has_poster_image"`
	HasOverlayPosterImage      bool      `json:"has_overlay_poster_image"`
	HasDetailsBackground       bool      `json:"has_details_background"`
	HasMobileDetailsBackground bool      `json:"has_mobile_details_background"`
	CastId                     string    `json:"cast_id"`
	MusicId                    string    `json:"music_id"`
	TagInfoId                  string    `json:"tag_info_id"`
	RightsId                   string    `json:"rights_id"`
	DeletedByUserId            string    `json:"deleted_by_user_id"`
	SeasonKey                  int       `json:"season_key"`
	CreatedAt                  time.Time `json:"created_at"`
	EnglishMetaTitle           string    `json:"english_meta_title"`
	ArabicMetaTitle            string    `json:"arabic_meta_title"`
	EnglishMetaDescription     string    `json:"english_meta_description"`
	ArabicMetaDescription      string    `json:"arabic_meta_description"`
	HasAllRights               bool      `json:"has_all_rights"`
	ThirdPartySeasonKey        int       `json:"third_party_season_key"`
}

type AgeRatingsCode struct {
	Code string `json:"code"`
}

type ContentEpisode struct {
	EpisodeNumber   int            `json:"episodeNumber"`
	EpisodeKey      int            `json:"episodeKey"`
	Length          int            `json:"length"`
	VideoContentUrl string         `json:"videoContentUrl"`
	SynopsisEnglish string         `json:"synopsisEnglish"`
	SynopsisArabic  string         `json:"synopsisArabic"`
	NonTextualData  NonTextualData `json:"nonTextualData"`
	EpisodeId       string         `json:"id"`
}
