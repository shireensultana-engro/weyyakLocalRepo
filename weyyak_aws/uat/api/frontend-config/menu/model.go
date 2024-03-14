package menu

import (
	"time"

	_ "github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type NewMenu struct {
	Url string `json:"url"`
}

type Cont struct {
	Cnt postgres.Jsonb
}

type Menu struct {
	ID              string `json:"id" gorm:"primary_key"`
	Device          string `json:"device" binding:"required"`
	MenuType        string `json:"menu_type"`
	MenuEnglishName string `json:"menu_english_name"`
	MenuArabicName  string `json:"menu_arabic_name"`
	SliderKey       int    `json:"slider_key"`
	Url             string `json:"url "`
	Order           int    `json:"order "`
	IsPublished     bool   `json:"is_published"`
}

type Publishingplatforms struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

/*page-fragment */
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
	Imagery        *ImageryDetails  `json:"imagery"`
	Featured       *FeaturedDetails `json:"featured"`
	Playlists      []MenuPlaylists  `json:"playlists"`
}

//MenuPageDetails struct for DB binding
type MenuPageDetails struct {
	ID                       string         `json:"id"`
	FriendlyUrl              string         `json:"friendly_url"`
	SeoDescription           string         `json:"seo_description"`
	Title                    string         `json:"title"`
	Type                     int            `json:"type"`
	PageType                 string         `json:"pageType,omitempty"`
	PageKey                  int            `json:"page_key"`
	Featured                 postgres.Jsonb `json:"Featured,omitempty"`
	PageOrderNumber          int            `json:"page_order_number,omitempty"`
	HasMobileMenu            bool           `json:"has_mobile_menu"`
	HasMenuPosterImage       bool           `json:"has_menu_poster_image"`
	HasMobileMenuPosterImage bool           `json:"has_mobile_menu_poster_image"`
	Playlists                string         `json:"playlists"`
}

//FeaturedResponse
type FeaturedResponse struct {
	ID        int      `json:"id"`
	Type      string   `json:"type"`
	Playlists []string `json:"playlists"`
}
type FeaturedPlaylistsResponse struct {
	ID           int            `json:"id"`
	PlaylistType string         `json:"playlist_type"`
	Title        string         `json:"title"`
	Content      postgres.Jsonb `json:"content"`
}

type PlaylistsResponse struct {
	ID      int            `json:"id"`
	Title   string         `json:"title"`
	Content postgres.Jsonb `json:"content"`
}

//ImageryDetails struct for DB binding
type ImageryDetails struct {
	MobileMenu            string `json:"mobile-menu"`
	MobilePosterImage     string `json:"menu-poster-image"`
	MobileMenuPosterImage string `json:"mobile-menu-poster-image"`
}
type Featured struct {
	Featured postgres.Jsonb `json:"Featured"`
}

type Plays struct {
	Playlists string `json:"playlists"`
}

//PlaylistContent struct for DB binding
type Content struct {
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

// type PlaylistX struct {
// 	Content PlaylistContent
// }

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

//FeaturedDetails struct for DB binding
type FeaturedDetails struct {
	ID        int                 `json:"id"`
	Type      string              `json:"type"`
	Playlists []FeaturedPlaylists `json:"playlists"`
}

//FeaturedPlaylists struct for DB binding
type FeaturedPlaylists struct {
	ID           int            `json:"id"`
	PlaylistType string         `json:"playlist_type"`
	Title        string         `json:"title"`
	Content      postgres.Jsonb `json:"content"`
}

type MenuPlaylists struct {
	ID      int            `json:"id"`
	Title   string         `json:"title"`
	Content postgres.Jsonb `json:"content"`
	// PlaylistType string        `json:"playlisttype"`
	// PageContent  []PageContent `json:"pagecontent"`
}
type PageContent struct {
	Key            string         `json:"key"`
	ID             int            `json:"id"`
	FriendlyUrl    string         `json:"friendly_url"`
	SeoDescription string         `jsfieldson:"seo_description"`
	Title          string         `json:"title"`
	Type           string         `json:"type"`
	Imagery        ImageryDetails `json:"imagery"`
}

type ContentImageryDetails struct {
	Thumbnail   string `json:"thumbnail"`
	Backdrop    string `json:"backdrop"`
	MobileImg   string `json:"mobile_img"`
	FeaturedImg string `json:"featured_img,omitempty"`
	Banner      string `json:"banner"`
}
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
}

/*side Menu */
type SideMenuDetails struct {
	Total       int64  `json:"total"`
	PerPage     int64  `json:"per_page"`
	CurrentPage int64  `json:"current_page"`
	LastPage    int64  `json:"last_page"`
	NextPageUrl string `json:"next_page_url"`
	PrevPageUrl string `json:"prev_page_url"`
	From        int64  `json:"from"`
	To          int64  `json:"to"`
	Data        []Data `json:"data"`
}

/* Redis-struct */
type RedisCacheResponse struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

/* HomePage- TV */
type Result struct {
	Data Data `json:"data"`
}

type Data struct {
	ID             int    `json:"id"`
	FriendlyURL    string `json:"friendly_url"`
	SeoDescription string `json:"seo_description"`
	Title          string `json:"title"`
	Type           string `json:"type"`
	Imagery        *struct {
		MobileMenu            string `json:"mobile-menu"`
		MenuPosterImage       string `json:"menu-poster-image"`
		MobileMenuPosterImage string `json:"mobile-menu-poster-image"`
	} `json:"imagery"`
	Featured *struct {
		ID        int    `json:"id"`
		Type      string `json:"type"`
		Playlists []struct {
			ID           int    `json:"id"`
			PlaylistType string `json:"playlist_type"`
			Title        string `json:"title"`
			Content      []struct {
				ID     int      `json:"id"`
				Cast   []string `json:"cast"`
				Tags   []string `json:"tags"`
				Title  string   `json:"title"`
				Genres []string `json:"genres"`
				Length int      `json:"length"`
				Movies []struct {
					ID                   int           `json:"id"`
					Title                string        `json:"title"`
					Geoblock             bool          `json:"geoblock"`
					InsertedAt           time.Time     `json:"insertedAt"`
					DigitalRighttype     int           `json:"digitalRighttype"`
					SubscriptiontPlans   []interface{} `json:"subscriptiontPlans"`
					DigitalRightsRegions interface{}   `json:"digitalRightsRegions"`
				} `json:"movies,omitempty"`
				Imagery struct {
					Banner      string `json:"banner"`
					Backdrop    string `json:"backdrop"`
					Thumbnail   string `json:"thumbnail"`
					MobileImg   string `json:"mobile_img"`
					FeaturedImg string `json:"featured_img"`
				} `json:"imagery"`
				Geoblock        bool   `json:"geoblock"`
				Synopsis        string `json:"synopsis"`
				VideoID         string `json:"video_id"`
				SeoTitle        string `json:"seo_title"`
				AgeRating       string `json:"age_rating"`
				ContentID       string `json:"content_id"`
				InsertedAt      string `json:"insertedAt"`
				MainActor       string `json:"main_actor"`
				ModifiedAt      string `json:"modifiedAt"`
				ContentType     string `json:"content_type"`
				FriendlyURL     string `json:"friendly_url"`
				MainActress     string `json:"main_actress"`
				ProductionYear  int    `json:"production_year"`
				SeoDescription  string `json:"seo_description"`
				TranslatedTitle string `json:"translated_title"`
				Seasons         []struct {
					ID                   int         `json:"id"`
					Title                string      `json:"Title"`
					Dubbed               bool        `json:"dubbed"`
					Geoblock             bool        `json:"geoblock"`
					SeoTitle             string      `json:"seo_title"`
					SeasonNumber         int         `json:"season_number"`
					SeoDescription       string      `json:"seo_description"`
					DigitalRighttype     int         `json:"digitalRighttype"`
					SubscriptionPlans    interface{} `json:"subscriptionPlans"`
					DigitalRightsRegions interface{} `json:"digitalRightsRegions"`
				} `json:"seasons,omitempty"`
			} `json:"content"`
		} `json:"playlists"`
	} `json:"featured"`
	Playlists []struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content []struct {
			ID      int      `json:"id"`
			Cast    []string `json:"cast"`
			Tags    []string `json:"tags"`
			Title   string   `json:"title"`
			Genres  []string `json:"genres"`
			Length  int      `json:"length"`
			Imagery struct {
				Banner      string `json:"banner"`
				Backdrop    string `json:"backdrop"`
				Thumbnail   string `json:"thumbnail"`
				MobileImg   string `json:"mobile_img"`
				FeaturedImg string `json:"featured_img"`
			} `json:"imagery"`
			Seasons []struct {
				ID                   int         `json:"id"`
				Title                string      `json:"Title"`
				Dubbed               bool        `json:"dubbed"`
				Geoblock             bool        `json:"geoblock"`
				SeoTitle             string      `json:"seo_title"`
				SeasonNumber         int         `json:"season_number"`
				SeoDescription       string      `json:"seo_description"`
				DigitalRighttype     int         `json:"digitalRighttype"`
				SubscriptionPlans    interface{} `json:"subscriptionPlans"`
				DigitalRightsRegions interface{} `json:"digitalRightsRegions"`
			} `json:"seasons,omitempty"`
			Movies []struct {
				ID                   int           `json:"id"`
				Title                string        `json:"title"`
				Geoblock             bool          `json:"geoblock"`
				InsertedAt           time.Time     `json:"insertedAt"`
				DigitalRighttype     int           `json:"digitalRighttype"`
				SubscriptiontPlans   []interface{} `json:"subscriptiontPlans"`
				DigitalRightsRegions interface{}   `json:"digitalRightsRegions"`
			} `json:"movies,omitempty"`
			Geoblock        bool   `json:"geoblock"`
			Synopsis        string `json:"synopsis"`
			VideoID         string `json:"video_id"`
			SeoTitle        string `json:"seo_title"`
			AgeRating       string `json:"age_rating"`
			ContentID       string `json:"content_id"`
			InsertedAt      string `json:"insertedAt"`
			MainActor       string `json:"main_actor"`
			ModifiedAt      string `json:"modifiedAt"`
			ContentType     string `json:"content_type"`
			FriendlyURL     string `json:"friendly_url"`
			MainActress     string `json:"main_actress"`
			ProductionYear  int    `json:"production_year"`
			SeoDescription  string `json:"seo_description"`
			TranslatedTitle string `json:"translated_title"`
		} `json:"content"`
	} `json:"playlists"`
}

func (d Data) Equals(other Data) bool {
	return d.Title == other.Title && d.FriendlyURL == other.FriendlyURL
}

type Page struct {
	PageKey                int
	EnglishTitle           string
	ArabicTitle            string
	EnglishPageFriendlyUrl string
	ArabicPageFriendlyUrl  string
	EnglishMetaDescription string
	ArabicMetaDescription  string
}
type Playlists []struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content []struct {
		ID      int      `json:"id"`
		Cast    []string `json:"cast"`
		Tags    []string `json:"tags"`
		Title   string   `json:"title"`
		Genres  []string `json:"genres"`
		Length  int      `json:"length"`
		Imagery struct {
			Banner      string `json:"banner"`
			Backdrop    string `json:"backdrop"`
			Thumbnail   string `json:"thumbnail"`
			MobileImg   string `json:"mobile_img"`
			FeaturedImg string `json:"featured_img"`
		} `json:"imagery"`
		Seasons []struct {
			ID                   int         `json:"id"`
			Title                string      `json:"Title"`
			Dubbed               bool        `json:"dubbed"`
			Geoblock             bool        `json:"geoblock"`
			SeoTitle             string      `json:"seo_title"`
			SeasonNumber         int         `json:"season_number"`
			SeoDescription       string      `json:"seo_description"`
			DigitalRighttype     int         `json:"digitalRighttype"`
			SubscriptionPlans    interface{} `json:"subscriptionPlans"`
			DigitalRightsRegions interface{} `json:"digitalRightsRegions"`
		} `json:"seasons,omitempty"`
		Movies []struct {
			ID                   int           `json:"id"`
			Title                string        `json:"title"`
			Geoblock             bool          `json:"geoblock"`
			InsertedAt           time.Time     `json:"insertedAt"`
			DigitalRighttype     int           `json:"digitalRighttype"`
			SubscriptiontPlans   []interface{} `json:"subscriptiontPlans"`
			DigitalRightsRegions interface{}   `json:"digitalRightsRegions"`
		} `json:"movies,omitempty"`
		Geoblock        bool   `json:"geoblock"`
		Synopsis        string `json:"synopsis"`
		VideoID         string `json:"video_id"`
		SeoTitle        string `json:"seo_title"`
		AgeRating       string `json:"age_rating"`
		ContentID       string `json:"content_id"`
		InsertedAt      string `json:"insertedAt"`
		MainActor       string `json:"main_actor"`
		ModifiedAt      string `json:"modifiedAt"`
		ContentType     string `json:"content_type"`
		FriendlyURL     string `json:"friendly_url"`
		MainActress     string `json:"main_actress"`
		ProductionYear  int    `json:"production_year"`
		SeoDescription  string `json:"seo_description"`
		TranslatedTitle string `json:"translated_title"`
	} `json:"content"`
}
type seasondetails struct {
	ID string `json:"id"`
}
type seasondetail struct {
	ContentID string `json:"content_id"`
}
type ContentIds struct {
	Ids postgres.Jsonb
}
type FilteredIds struct {
	Filteredids string
	Id          string
}
type VerifyPlaylist []struct {
	ID     int      `json:"id"`
	Cast   []string `json:"cast"`
	Tags   []string `json:"tags"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
	Length int      `json:"length"`
	Movies []struct {
		ID                   int           `json:"id"`
		Title                string        `json:"title"`
		Geoblock             bool          `json:"geoblock"`
		InsertedAt           time.Time     `json:"insertedAt"`
		DigitalRighttype     int           `json:"digitalRighttype"`
		SubscriptiontPlans   []interface{} `json:"subscriptiontPlans"`
		DigitalRightsRegions interface{}   `json:"digitalRightsRegions"`
	} `json:"movies"`
	Imagery struct {
		Banner      string `json:"banner"`
		Backdrop    string `json:"backdrop"`
		Thumbnail   string `json:"thumbnail"`
		MobileImg   string `json:"mobile_img"`
		FeaturedImg string `json:"featured_img"`
	} `json:"imagery"`
	Geoblock        bool   `json:"geoblock"`
	Synopsis        string `json:"synopsis"`
	VideoID         string `json:"video_id"`
	SeoTitle        string `json:"seo_title"`
	AgeRating       string `json:"age_rating"`
	ContentID       string `json:"content_id"`
	InsertedAt      string `json:"insertedAt"`
	MainActor       string `json:"main_actor"`
	ModifiedAt      string `json:"modifiedAt"`
	ContentType     string `json:"content_type"`
	FriendlyURL     string `json:"friendly_url"`
	MainActress     string `json:"main_actress"`
	ProductionYear  int    `json:"production_year"`
	SeoDescription  string `json:"seo_description"`
	TranslatedTitle string `json:"translated_title"`
}
