package page

import (
	"time"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PublishPlatformDetails struct {
	Id int `json:"platform"`
}

// type PublishTargetPlatformDetails struct{
// 	Targetplatform  int `json:"target_platform"`
// }

type PlatformDetails struct {
	Id              string `json:"pageId"`
	Englishtitle    string `json:"englishTitle"`
	Pageordernumber int    `json:"pageOrderNumber"`
	Ishome          bool   `json:"isHome"`
}

type PageDetails struct {
	EnglishTitle           string `json:"englishTitle"`
	ArabicTitle            string `json:"arabicTitle" `
	PageOrderNumber        int    `json:"pageOrderNumber"`
	EnglishPageFriendlyUrl string `json:"englishPageFriendlyUrl"`
	ArabicPageFriendlyUrl  string `json:"arabicPageFriendlyUrl"`
	EnglishMetaTitle       string `json:"englishMetaTitle"`
	ArabicMetaTitle        string `json:"arabicMetaTitle"`
	EnglishMetaDescription string `json:"englishMetaDescription"`
	ArabicMetaDescription  string `json:"arabicMetaDescription"`
	Ishome                 bool   `json:"isHome"`
	IsDisabled             *bool  `json:"isDisabled"`
	Mobilemenu             string `json:"mobileMenu"`
	Menuposterimage        string `json:"menuPosterImage"`
	Mobilemenuposterimage  string `json:"mobileMenuPosterImage"`
	TargetPlatform         int    `json:"targetplatform"`
	CountryId              int    `json:"countryid"`
	Platforms              string `json:"platforms,omitempty"`
	Regions                string `json:"regions,omitempty"`
	PlatformOrder          string `json:"platformOrder,omitempty"`
}
type PlaylistIds struct {
	PlaylistId string `json:"playlistid"`
}
type Playlistdetails struct {
	EnglishTitle        string `json:"englishTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	IsDisabled          *bool  `json:"isDisabled"`
	PublishingPlatforms *[]int `json:"publishingPlatforms"`
	Id                  string `json:"id"`
	Platforms           string `json:"platforms,omitempty"`
}
type Targetplatform struct {
	TargetPlatform int `json:"targetplatform"`
}
type Regions struct {
	CountryId int `json:"countryid"`
}
type Sliderid struct {
	SliderId string `json:"sliderid"`
	Order    int    `json:"order"`
}
type Sliderdetails struct {
	Name       string `json:"name"`
	IsDisabled *bool  `json:"isDisabled"`
	Id         string `json:"id"`
}
type PageOrder struct {
	PageOrderNumber int `json:"pageordernumber"`
}
type PublishPlatformdata struct {
	PlayListId     string `json:"playlistid"`
	TargetPlatform int    `json:"targetplatform"`
}
type Sample struct {
	PlayListId     string `json:"playlistid"`
	TargetPlatform []int  `json:"targetplatform"`
}

type DeletePageDetails struct {
	Id              string `json:"id"`
	DeletedByUserId string `json:"deleted_by_user_id"`
}

//Disable or enable based on page id
type PageAvailability struct {
	Id         string `json:"id" binding:"required"`
	IsDisabled *bool  `json:"isDisabled"  binding:"required"`
}
type PageId struct {
	Id string
}
type CheckPageDetails struct {
	DeletedByUserId string
}

type PagelistSummary struct {
	ID                  string `json:"id"`
	EnglishTitle        string `json:"englishTitle"`
	ArabicTitle         string `json:"arabicTitle"`
	IsDisabled          bool   `json:"isDisabled"`
	Ishome              bool   `json:"isHome"`
	PublishingPlatforms []int  `json:"publishingPlatforms"`
}

type PlanNames struct {
	TargetPlatform int `json:"target_platform"`
}

type PageDetailsSummary struct {
	PageOrderNumber        int    `json:"pageOrderNumber"`
	Region                 string `json:"region"`
	HasMoreRegions         bool   `json:"hasMoreRegions"`
	HasPublishingPlatforms bool   `json:"hasPublishingPlatforms"`
	HasPlaylists           bool   `json:"hasPlaylists"`
	HasSliders             bool   `json:"hasSliders"`
	EnglishTitle           string `json:"englishTitle"`
	ArabicTitle            string `json:"arabicTitle"`
	IsDisabled             bool   `json:"isDisabled"`
	IsHome                 bool   `json:"isHome"`
	PublishingPlatforms    []int  `json:"publishingPlatforms"`
	ID                     string `json:"id"`
}
type PageRegionDetails struct {
	Details string `json:"details"`
}
type HasPlaylist struct {
	HasPlaylistId bool
	HasSliderId   bool
}

// update page order details
type PageResponse struct {
	PageId          string `json:"pageId"`
	EnglishTitle    string `json:"englishTitle"`
	PageOrderNumber int    `json:"pageOrderNumber"`
	IsHome          bool   `json:"isHome"`
}

type Response struct {
	TargetPlatform0  []PageResponse `json:"0"`
	TargetPlatform1  []PageResponse `json:"1"`
	TargetPlatform2  []PageResponse `json:"2"`
	TargetPlatform3  []PageResponse `json:"3"`
	TargetPlatform4  []PageResponse `json:"4"`
	TargetPlatform5  []PageResponse `json:"5"`
	TargetPlatform6  []PageResponse `json:"6"`
	TargetPlatform7  []PageResponse `json:"7"`
	TargetPlatform9  []PageResponse `json:"9"`
	TargetPlatform10 []PageResponse `json:"10"`
}
type TargetPlatformResponse struct {
	TargetPlatform int `json:"target_platform"`
}

type FinalResponse struct {
	PublishingPlatformOrderedDetails Response `json:"publishingPlatformOrderedDetails"`
	Id                               string   `json:"id"`
}

type PageTargetPlatform struct {
	PageId          string `json:"pageId"`
	TargetPlatform  int    `json:"targetPlatform"`
	PageOrderNumber int    `json:"page_order_number"`
}

//Error Codes
type Id struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type EnglishTitleError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ArabicTitleError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type EnglishPageFriendlyError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ArabicPageFriendlyError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type EnglishMetaTitleError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ArabicMetaTitleError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

type EnglishMetaDescriptionError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ArabicMetaDescriptionError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type PageOrderNUmberError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type PublishingPlatformsError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type RegionssError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type SlidersIds struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type Invalid struct {
	Id                          *Id         `json:"id,omitempty"`
	EnglishTitleError           interface{} `json:"englishTitle,omitempty"`
	ArabicTitleError            interface{} `json:"arabicTitle,omitempty"`
	EnglishPageFriendlyError    interface{} `json:"englishPageFriendlyUrl,omitempty"`
	ArabicPageFriendlyError     interface{} `json:"arabicPageFriendlyUrl,omitempty"`
	EnglishMetaTitleError       interface{} `json:"englishMetaTitle,omitempty"`
	ArabicMetaTitleError        interface{} `json:"arabicMetaTitle,omitempty"`
	EnglishMetaDescriptionError interface{} `json:"englishMetaDescription,omitempty"`
	ArabicMetaDescriptionError  interface{} `json:"arbicMetaDescription,omitempty"`
	PageOrderNUmberError        interface{} `json:"pageOrderNumber,omitempty"`
	PublishingPlatformsError    interface{} `json:"publishingPlatforms,omitempty"`
	RegionssError               interface{} `json:"regions,omitempty"`
	SlidersIds                  interface{} `json:"slidersIds,omitempty"`
}

type FinalErrorResponse struct {
	Error       string  `json:"error"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	RequestId   string  `json:"requestId"`
	Invalid     Invalid `json:"invalid,omitempty"`
}
type FinalError struct {
	Error       string `json:"error"`
	Description string `json:"description"`
	Code        string `json:"code"`
	RequestId   string `json:"requestId"`
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
	ThirdPartyPageKey                  int       `json:"third_party_page_key,omitempty"`

}

type PageRequest struct {
	// for sync remove
	//	PageId                 string         `json:"id"` // for creating old pages with .net
	EnglishTitle           string         `json:"englishTitle"`
	ArabicTitle            string         `json:"arabicTitle"`
	PageOrderNumber        int            `json:"pageOrderNumber"`
	IsHome                 bool           `json:"isHome"`
	PlaylistsIds           []string       `json:"playlistsIds"`
	EnglishPageFriendlyUrl string         `json:"englishPageFriendlyUrl"`
	ArabicPageFriendlyUrl  string         `json:"arabicPageFriendlyUrl"`
	EnglishMetaTitle       string         `json:"englishMetaTitle"`
	ArabicMetaTitle        string         `json:"arabicMetaTitle"`
	EnglishMetaDescription string         `json:"englishMetaDescription"`
	ArabicMetaDescription  string         `json:"arabicMetaDescription"`
	PublishingPlatforms    []int          `json:"publishingPlatforms"`
	Regions                []int          `json:"regions"`
	DefaultSliderId        string         `json:"defaultSliderId"`
	NonTextualData         NonTextualData `json:"nonTextualData"`
	SlidersIds             []string       `json:"slidersIds"`
	PageKey                int            `json:"pageKey"`
}

type NonTextualData struct {
	MenuPosterImage       string `json:"menuPosterImage,omitempty"`
	MobileMenuPosterImage string `json:"mobileMenuPosterImage,omitempty"`
	MobileMenu            string `json:"mobileMenu,omitempty"`
}

type PagePlaylist struct {
	PageId     string `json:"page_id,omitempty"`
	PlaylistId string `json:"playlist_id,omitempty"`
	Order      int    `json:"order,omitempty"`
}

type PageCountry struct {
	PageId    string `json:"page_id,omitempty"`
	CountryId int    `json:"country_id,omitempty"`
}
type PageSlider struct {
	PageId   string `json:"page_id,omitempty"`
	SliderId string `json:"slider_id,omitempty"`
	Order    int    `json:"order,omitempty"`
}

/*pageRegions*/
type Country struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}
