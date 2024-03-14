package mrssfeed

import (
	"encoding/xml"
	"time"
)

// import (
// 	"time"

// 	_ "github.com/jinzhu/gorm/dialects/postgres"
// )

// type Rss struct {
// 	//Pagination Pagination `json:"pagination"`
// 	Xsd        string    `xml:"xsd,attr"`
// 	Xsi        string    `xml:"xsi,attr"`
// 	Media      string    `xml:"media,attr"`
// 	Atom       string    `xml:"atom,attr"`
// 	OpenSearch string    `xml:"openSearch,attr"`
// 	Dfpvideo   string    `xml:"dfpvideo,attr"`
// 	Version    string    `xml:"version,attr"`
// 	Channel    []channel `xml:"channel"`
// }
// type channel struct {
// 	Title   string `xml:"title"`
// 	Version string `xml:"version"`
// 	Rel     Rel    `xml:"rel,omitempty"`
// 	Href    Href   `xml:"hrel,omitempty"`
// 	Item    []item `xml:"item"`
// }
// type Rel struct {
// 	Rel  string
// 	Href string
// }
// type Href struct {
// 	Rel  string
// 	Href string
// }
// type item struct {
// 	PubDate          time.Time `xml:"pubDate"`
// 	Title            string    `xml:"title"`
// 	Thumbnail        thumbnail `xml:"media:thumbnail"`
// 	Content          content   `xml:"media:Content"`
// 	Keyvalues        interface{}
// 	ContentId        string `xml:"dfpvideo:contentId"`
// 	LastModifiedDate string `xml:"dfpvideo:lastModifiedDate"`
// 	Cuepoints        string `xml:"dfpvideo:cuepoints"`
// }

// type thumbnail struct {
// 	URL string `xml:"thumbnail,attr"`
// }
// type content struct {
// 	Eng eng `xml:"media:Content"`
// 	Arb arb `xml:"media:content"`
// }
// type eng struct {
// 	Duration int    `xml:"duration,attr"`
// 	Url      string `xml:"url,attr"`
// }
// type arb struct {
// 	Duration int    `xml:"duration,attr"`
// 	Url      string `xml:"url,attr"`
// }
// type dfpvideo struct {
// 	Keyvalues []keyvaluee `xml:"dfpvideo"`
// }
// type keyvaluee struct {
// 	Key   string `xml:"key,attr"`
// 	Value string `xml:"value,attr"`
// 	//val   int    `xml:"val,attr"`
// 	Type string `xml:"type,attr"`
// }
// type Activecontent struct {
// 	ContentId   string `xml:"content_Id"`
// 	ContentType string `xml:"content_type"`
// 	ContentTier int    `xml:"content_tier"`
// }
// type AllContentDetails struct {
// 	CreatedAt       time.Time `xml:"createdAt"`
// 	Title           string    `xml:"title"`
// 	Thumbnail       string    `xml:"url,attr"`
// 	Duration        int       `xml:"duration"`
// 	Urlen           string    `xml:"linke,attr"`
// 	Urlar           string    `xml:"linka,attr"`
// 	SeasonNumber    string    `xml:"seasonNumber,omitempty"`
// 	EpisodeNumber   string    `xml:"episodeNumber,omitempty"`
// 	EnglishGenre    string    `xml:"englishGenre"`
// 	ArabicGenre     string    `xml:"arabicGenre"`
// 	EnglishSubgenre string    `xml:"englishSubGenre"`
// 	ArabicSubgenre  string    `xml:"arabicSubGenre"`
// 	//Subgenre        string    `xml:"subgenre"`
// 	EnglishTitle string `xml:"englishTitle"`
// 	ArabicTitle  string `xml:"arabicTitle"`
// 	ContentType  string `xml:"contenType"`
// 	Language     string `xml:"language"`
// 	ContentId    string `xml:"contentId"`
// 	SeasonId     string `xml:"seasonId"`
// 	EpisodeId    string `xml:"episodeId"`
// 	EpisodeKey   string `xml:"episodeKey"`
// 	ContentKey   string `xml:"contentKey"`
// 	ModifiedAt   string `xml:"modifiedAt"`
// 	CuePoint     string `xml:"cuePoint`
// }

// type Pagination struct {
// 	Size   int `json:"size"`
// 	Offset int `json:"offset"`
// 	Limit  int `json:"limit"`
// }
// type data struct {
// 	Dfpvideo []struct {
// 		Key   string `xml:"key,attr"`
// 		Value string `xml:"value,omitempty,attr"`
// 		Val   string `xml:"val,omitempty,attr"`
// 		Type  string `xml:"type,attr"`
// 	} `xml:"dfpvideo"`
// }
// type ContentDetails struct {
// 	Type                int    `json:"type"`
// 	ContentKey          int    `json:"contentKey"`
// 	Status              int    `json:"status"`
// 	StatusCanBeChanged  bool   `json:"statusCanBeChanged"`
// 	SubStatusName       string `json:"subStatusName"`
// 	TransliteratedTitle string `json:"transliteratedTitle"`
// 	CreatedBy           string `json:"createdBy"`
// 	Id                  string `json:"id"`
// 	CreatedByUserId     string `json:"-"`
// }
// type count struct {
// }

type Rss struct {
	XMLName    xml.Name `xml:"rss"`
	Text       string   `xml:",chardata"`
	Xsd        string   `xml:"xmlns:xsd,attr"`
	Xsi        string   `xml:"xmlns:xsi,attr"`
	Media      string   `xml:"xmlns:media,attr"`
	Atom       string   `xml:"xmlns:atom,attr"`
	OpenSearch string   `xml:"xmlns:openSearch,attr"`
	Dfpvideo   string   `xml:"xmlns:dfpvideo,attr"`
	Version    string   `xml:"version,attr"`
	Channel    Channel  `xml:"channel"`
}

type Channel struct {
	Text    string `xml:",chardata"`
	Title   string `xml:"title"`
	Version string `xml:"dfpvideo:version"`
	Link    []Link `xml:"atom:link"`
	Item    []Item `xml:"item"`
}

type Link struct {
	Text string `xml:",chardata"`
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type Item struct {
	Text      string `xml:",chardata"`
	PubDate   string `xml:"pubDate"`
	Title     string `xml:"title"`
	Thumbnail struct {
		Text string `xml:",chardata"`
		URL  string `xml:"url,attr"`
	} `xml:"media:thumbnail"`
	Content          []Content   `xml:"media:content"`
	Keyvalues        []Keyvalues `xml:"dfpvideo:keyvalues"`
	ContentId        string      `xml:"dfpvideo:contentId"`
	LastModifiedDate string      `xml:"dfpvideo:lastModifiedDate"`
	Cuepoints        string      `xml:"dfpvideo:cuepoints"`
}

type Keyvalues struct {
	Text  string `xml:",chardata"`
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
	Type  string `xml:"type,attr"`
}

type Content struct {
	Text     string `xml:",chardata"`
	Duration string `xml:"duration,attr"`
	URL      string `xml:"url,attr"`
}

type ContentDetails struct {
	ModifiedAt          time.Time
	CreatedAt           time.Time
	ContentKey          int
	ContentType         string
	TransliteratedTitle string
	Id                  string
	SeasonId            string
	EpisodeId           string
	Duration            int
	SeasonNumber        int
	EpisodeNumber       int
	ArabicTitle         string
	OriginalLanguage    string
}

type GenreDetails struct {
	ContentId     string
	GenresEnglish string
	GenresArabic  string
}

type SubGenreDetails struct {
	ContentId        string
	SubgenresEnglish string
	SubgenresArabic  string
}
