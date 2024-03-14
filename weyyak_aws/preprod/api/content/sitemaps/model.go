package sitemaps

import (
	"encoding/xml"
)

/* Get All Sitemap xml files */
//github.com/aws/aws-sdk-go/aws"
// "github.com/aws/aws-sdk-go/aws/credentials"
// "github.com/aws/aws-sdk-go/aws/session"
// "github.com/aws/aws-sdk-go/service/s3")
type Sitemapindex struct {
	XMLName string `xml:"sitemapindex"`
	Xmlns   string `xml:"xmlns,attr"`
	Sitemap []struct {
		Location     string `xml:"loc"`
		Lastmod string `xml:"lastmod"`
	} `xml:"sitemap"`
}

/* Get All available Pages IN Arabic-XML file */
type ActivePages struct {
	EnglishTitle string `json:"english_title"`
}
type ActiveArabicPages struct {
	ArabicTitle string `json:"arabic_title"`
}

type Urlset struct {
	XMLName        string `xml:"urlset"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Xmlns          string   `xml:"xmlns,attr"`
	URL            []Result `xml:"url"`
}

type Result struct {
	SchemaLocation        string `xml:"loc"`
	Changefrequency string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

/* Get all movie titles in arabic-xml file */
type activeMovieTitlesArabic struct {
	ContentKey  string `json:"content_key"`
	ArabicTitle string `json:"arabic_title"`
}
type Url struct {
	XMLName xml.Name `xml:"urlset"`
	Xsi            string `xml:"xsi,attr"`
	SchemaLocation string `xml:"schemaLocation,attr"`
	Xmlns          string `xml:"xmlns,attr"`
	URL            []struct {
		Location        string `xml:"loc"`
		Changefreqency string `xml:"changefreq"`
		Priority   string `xml:"priority"`
	} `xml:"url"`
}
type activeMovieTitlesEnglish struct {
	ContentKey          string `json:"content_key"`
	TransliteratedTitle string `json:"transliterated_title"`
}
type MovieEN struct {
	XMLName string `xml:"urlset"`
	Xsi            string `xml:"xsi,attr"`
	SchemaLocation string `xml:"schemaLocation,attr"`
	Xmlns          string `xml:"xmlns,attr"`
	URL            []struct {
		Location        string `xml:"loc"`
		Changefrequency string `xml:"changefreq"`
		Priority   string `xml:"priority"`
	} `xml:"url"`
}
type activeSeriesTitlesArabic struct {
	ContentKey  string `json:"content_key"`
	ArabicTitle string `json:"arabic_title"`
}
type SeriesAR struct {
	XMLName string `xml:"urlset"`
	Xsi            string `xml:"xsi,attr"`
	SchemaLocation string `xml:"schemaLocation,attr"`
	Xmlns          string `xml:"xmlns,attr"`
	URL            []struct {
		Location        string `xml:"loc"`
		Changefreqency string `xml:"changefreq"`
		Priority   string `xml:"priority"`
	} `xml:"url"`
}
type activeSeriesTitlesEnglish struct {
	ContentKey          string `json:"content_key"`
	TransliteratedTitle string `json:"transliterated_title"`
}
type SeriesEN struct {
	XMLName xml.Name `xml:"urlset"`
	Xsi            string `xml:"xsi,attr"`
	SchemaLocation string `xml:"schemaLocation,attr"`
	Xmlns          string `xml:"xmlns,attr"`
	URL            []struct {
		Location        string `xml:"loc"`
		Changefreqency string `xml:"changefreq"`
		Priority   string `xml:"priority"`
	} `xml:"url"`
}
type activeEpisodeDetailsEnglish struct {
	EpisodeKey          string `json:"episode_key"`
	TransliteratedTitle string `json:"transliterated_title,omitempty"`
	ArabicTitle         string `json:"arabic_title,omitempty"`
}
type Episode struct {
	XMLName        string   `xml:"urlset"`
	Text           string       `xml:",chardata"`
	Xsi            string       `xml:"xsi,attr"`
	SchemaLocation string       `xml:"schemaLocation,attr"`
	Xmlns          string       `xml:"xmlns,attr"`
	URL            []EpisodeUrl `xml:"url"`
}

type EpisodeUrl struct {
	Text       string `xml:",chardata"`
	Location        string `xml:"loc"`
	Changefreqency string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

