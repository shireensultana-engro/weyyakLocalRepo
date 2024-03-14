package sitemaps

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	//"github.com/sabloger/sitemap-generator/smg"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	qrg := r.Group("/api")
	qrg.GET("/sitemap/index.xml", hs.Indexsitemap)
	qrg.GET("/sitemap/generic-en.xml", hs.GetAllavailablePagesINEnglishXMLfile)
	qrg.GET("/sitemap/generic-ar.xml", hs.GetAllavailablePagesINArabicXMLfile)
	qrg.GET("/sitemap/sitemap_Movie-ar.xml", hs.GetAllmoviedetailsINArabicXMLfile)
	qrg.GET("/sitemap/sitemap_Movie-en.xml", hs.GetAllmoviedetailsINEnglishXMLfile)
	qrg.GET("/sitemap/sitemap_series-ar.xml", hs.GetAllseriesdetailsINArabicXMLfile)
	qrg.GET("/sitemap/sitemap_series-en.xml", hs.GetAllseriesdetailsINEnglishXMLfile)
	qrg.GET("/sitemap/sitemap_shows-en.xml", hs.GetAllEpisodedetailsINEnglishXMLfile)
	qrg.GET("/sitemap/sitemap_shows-ar.xml", hs.GetAllEpisodedetailsINArabicXMLfile)
}

//  Indexsitemap
// GET /api/sitemap/index.xml
// @Summary Get All Sitemap xml files
// @Description Get All Sitemap xml files
// @Tags sitemap
// @Accept xml
// @Produce  xml
// @Success 200 {array} Sitemapindex{}
// @Router /api/sitemap/index.xml [GET]
/* Get All Sitemap xml files */
func (hs *HandlerService) Indexsitemap(c *gin.Context) {
	url := os.Getenv("SITE_URL") + "/sitemap/index.xml"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Status error: %v", res.StatusCode)
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Read body: %v", err)
		return
	}
	var result Sitemapindex
	xml.Unmarshal(data, &result)
	xmlout, err := xml.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(xmlout))
	filename := "Index-wyk2.xml"
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Println(err)
	}
	UploadFileToS3(data, filename)
	c.XML(http.StatusOK, result)
}

//  availablePagesINEnglish
// GET /api/sitemap/generic-en.xml
// @Summary Get All available Pages IN English-XML file
// @Description Get All available Pages IN English-XML file
// @Tags sitemap
// @Accept xml
// @Produce  xml
// @Success 200 {array} object c.xml
// @Router /api/sitemap/generic-en.xml [GET]
/* Get All available Pages IN English-XML file */
func (hs *HandlerService) GetAllavailablePagesINEnglishXMLfile(c *gin.Context) {
	db := c.MustGet("FCDB").(*gorm.DB)
	var activePages []ActivePages
	var result Result
	var finalResult []Result
	var urlset Urlset
	url := os.Getenv("SITE_URL") + "/en/"
	if err := db.Debug().Table("page").Select("english_title").Where("is_disabled=false and deleted_by_user_id is null").Order("page_order_number asc").Find(&activePages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	urlset.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	urlset.SchemaLocation = "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	urlset.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	for _, value := range activePages {
		result.SchemaLocation = url + value.EnglishTitle
		result.Changefrequency = "weekly"
		result.Priority = "0.8"
		finalResult = append(finalResult, result)
	}
	urlset.URL = finalResult

	xmlout, _ := xml.Marshal(urlset)
	_ = xml.Unmarshal([]byte(xmlout), &urlset)
	buf, _ := xml.MarshalIndent(urlset, "", "\t")
	filename := "generic_en_wyk2.xml"
	err1 := ioutil.WriteFile(filename, buf, 0644)
	if err1 != nil {
		log.Println(err1)
	}
	UploadFileToS3(buf, filename)
	c.XML(http.StatusOK, urlset)
}

//  availablePagesINArabic
// GET /api/sitemap/generic-ar.xml
// @Summary Get All available Pages IN Arabic-XML file
// @Description Get All available Pages IN Arabic-XML file
// @Tags sitemap
// @Accept xml
// @Produce  xml
// @Success 200 {array} object c.xml
// @Router /api/sitemap/generic-ar.xml [GET]
/* Get All available Pages IN Arabic-XML file */
func (hs *HandlerService) GetAllavailablePagesINArabicXMLfile(c *gin.Context) {
	db := c.MustGet("FCDB").(*gorm.DB)
	var activePages []ActiveArabicPages
	var result Result
	var finalResult []Result
	var urlset Urlset
	url := os.Getenv("SITE_URL") + "/en/"
	if err := db.Debug().Table("page").Select("arabic_title").Where("is_disabled=false and deleted_by_user_id is null").Order("page_order_number asc").Find(&activePages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	urlset.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	urlset.SchemaLocation = "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	urlset.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	for _, value := range activePages {
		result.SchemaLocation = url + value.ArabicTitle
		result.Changefrequency = "weekly"
		result.Priority = "0.8"
		finalResult = append(finalResult, result)
	}
	urlset.URL = finalResult
	xmlout, _ := xml.Marshal(urlset)
	_ = xml.Unmarshal([]byte(xmlout), &urlset)
	buf, _ := xml.MarshalIndent(urlset, "", "\t")
	filename := "generic_ar_wyk2.xml"
	err1 := ioutil.WriteFile(filename, buf, 0644)
	if err1 != nil {
		log.Println(err1)
	}
	UploadFileToS3(buf, filename)
	c.XML(http.StatusOK, urlset)
}

//  availableMoviesINArabic
// GET /api/sitemap/sitemap_Movie-ar.xml
// @Summary Get All movie details IN Arabic-XML file
// @Description Get All movie details IN Arabic-XML file
// @Tags sitemap
// @Accept xml
// @Produce  xml
// @Success 200 {array} object c.xml
// @Router /api/sitemap/sitemap_Movie-ar.xml [GET]
/* Get All movie details IN Arabic-XML file */
func (hs *HandlerService) GetAllmoviedetailsINArabicXMLfile(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var activeMovieTitlesArabic []activeMovieTitlesArabic
	var result Result
	var finalResult []Result
	var urlset Urlset
	url := os.Getenv("SITE_URL") + "/ar/movie/"
	if err := db.Debug().Raw("select c.content_key, cpi.arabic_title from content c join content_primary_info cpi on cpi.id = c.primary_info_id join about_the_content_info atci on atci.id = c.about_the_content_info_id join content_variance cv on cv.content_id = c.id join playback_item pi1 on pi1.id = cv.playback_item_id join content_rights cr on cr.id = pi1.rights_id  where c.status = 1 and c.deleted_by_user_id is null and (pi1.scheduling_date_time <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and (c.content_tier) = 1").Find(&activeMovieTitlesArabic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	urlset.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	urlset.SchemaLocation = "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	urlset.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	for _, value := range activeMovieTitlesArabic {
		result.SchemaLocation = url + value.ContentKey + "/" + value.ArabicTitle
		result.Changefrequency = "monthly"
		result.Priority = "0.5"
		finalResult = append(finalResult, result)
	}
	urlset.URL = finalResult
	xmlout, _ := xml.Marshal(urlset)
	_ = xml.Unmarshal([]byte(xmlout), &urlset)
	buf, _ := xml.MarshalIndent(urlset, "", "\t")
	filename := "Movie_ar_wyk2.xml"
	err1 := ioutil.WriteFile(filename, buf, 0644)
	if err1 != nil {
		log.Println(err1)
	}
	UploadFileToS3(buf, filename)
	c.XML(http.StatusOK, urlset)
}

//  availableMoviesINEnglish
// GET /api/sitemap/sitemap_Movie-ar.xml
// @Summary Get All movie details IN English-XML file
// @Description Get All movie details IN English-XML file
// @Tags sitemap
// @Accept xml
// @Produce  xml
// @Success 200 {array} object c.xml
// @Router /api/sitemap/sitemap_Movie-en.xml [GET]
//Get all movie titles in english-xml file
func (hs *HandlerService) GetAllmoviedetailsINEnglishXMLfile(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var activeMovieTitlesEnglish []activeMovieTitlesEnglish
	var result Result
	var finalResult []Result
	var urlset Urlset
	url := os.Getenv("SITE_URL") + "/en/movie/"
	if err := db.Debug().Raw("select c.content_key, cpi.transliterated_title from content c join content_primary_info cpi on cpi.id = c.primary_info_id join about_the_content_info atci on atci.id = c.about_the_content_info_id join content_variance cv on cv.content_id = c.id join playback_item pi1 on pi1.id = cv.playback_item_id join content_rights cr on cr.id = pi1.rights_id  where c.status = 1 and c.deleted_by_user_id is null and (pi1.scheduling_date_time <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and  (c.content_tier) = 1").Find(&activeMovieTitlesEnglish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	urlset.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	urlset.SchemaLocation = "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	urlset.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	for _, value := range activeMovieTitlesEnglish {
		result.SchemaLocation = url + value.ContentKey + "/" + value.TransliteratedTitle
		result.Changefrequency = "monthly"
		result.Priority = "0.5"
		finalResult = append(finalResult, result)
	}
	urlset.URL = finalResult
	xmlout, _ := xml.Marshal(urlset)
	_ = xml.Unmarshal([]byte(xmlout), &urlset)
	buf, _ := xml.MarshalIndent(urlset, "", "\t")
	filename := "Movie_en_wyk2.xml"
	err1 := ioutil.WriteFile(filename, buf, 0644)
	if err1 != nil {
		log.Println(err1)
	}
	UploadFileToS3(buf, filename)
	c.XML(http.StatusOK, urlset)
}

//  availableseriesINArabic
// GET /api/sitemap/sitemap_series-ar.xml
// @Summary Get All series details IN Arabic-XML file
// @Description Get All series details IN Arabic-XML file
// @Tags sitemap
// @Accept xml
// @Produce  xml
// @Success 200 {array} object c.xml
// @Router /api/sitemap/sitemap_series-ar.xml [GET]
//Get all series titles in arabic-xml file
func (hs *HandlerService) GetAllseriesdetailsINArabicXMLfile(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var activeSeriesTitlesArabic []activeSeriesTitlesArabic
	var result Result
	var finalResult []Result
	var urlset Urlset
	url := os.Getenv("SITE_URL") + "/ar/series/"
	if err := db.Debug().Raw("select distinct c.content_key, cpi.arabic_title from content c join content_primary_info cpi on cpi.id = c.primary_info_id join season s on s.content_id = c.id join episode e on e.season_id = s.id join playback_item pi1 on pi1.id = e.playback_item_id join content_rights cr on cr.id = s.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where (pi1.scheduling_date_time <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and c.status = 1 and c.deleted_by_user_id is null and s.status = 1 and s.deleted_by_user_id is null and e.status = 1 and e.deleted_by_user_id is null and (c.content_tier) = 2").Find(&activeSeriesTitlesArabic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	urlset.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	urlset.SchemaLocation = "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	urlset.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	for _, value := range activeSeriesTitlesArabic {
		result.SchemaLocation = url + value.ContentKey + "/" + value.ArabicTitle
		result.Changefrequency = "monthly"
		result.Priority = "0.5"
		finalResult = append(finalResult, result)
	}
	urlset.URL = finalResult
	xmlout, _ := xml.Marshal(urlset)
	_ = xml.Unmarshal([]byte(xmlout), &urlset)
	buf, _ := xml.MarshalIndent(urlset, "", "\t")
	filename := "series_ar_wyk2.xml"
	err1 := ioutil.WriteFile(filename, buf, 0644)
	if err1 != nil {
		log.Println(err1)
	}
	UploadFileToS3(buf, filename)
	c.XML(http.StatusOK, urlset)
}

//  availableseriesINEnglish
// GET /api/sitemap/sitemap_series-en.xml
// @Summary Get All series details IN English-XML file
// @Description Get All series details IN English-XML file
// @Tags sitemap
// @Accept xml
// @Produce  xml
// @Success 200 {array} object c.xml
// @Router /api/sitemap/sitemap_series-en.xml [GET]
//Get all series titles in english-xml file
func (hs *HandlerService) GetAllseriesdetailsINEnglishXMLfile(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var activeSeriesTitlesEnglish []activeSeriesTitlesEnglish
	var result Result
	var finalResult []Result
	var urlset Urlset
	url := os.Getenv("SITE_URL") + "/en/series/"
	if err := db.Debug().Raw("select distinct c.content_key, cpi.transliterated_title from content c join content_primary_info cpi on cpi.id = c.primary_info_id join season s on s.content_id = c.id join episode e on e.season_id = s.id join playback_item pi1 on pi1.id = e.playback_item_id join content_rights cr on cr.id = s.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where (pi1.scheduling_date_time <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and c.status = 1 and c.deleted_by_user_id is null and s.status = 1 and s.deleted_by_user_id is null and e.status = 1 and e.deleted_by_user_id is null and (c.content_tier) = 2").Find(&activeSeriesTitlesEnglish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	urlset.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	urlset.SchemaLocation = "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	urlset.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	for _, value := range activeSeriesTitlesEnglish {
		result.SchemaLocation = url + value.ContentKey + "/" + value.TransliteratedTitle
		result.Changefrequency = "monthly"
		result.Priority = "0.5"
		finalResult = append(finalResult, result)
	}
	urlset.URL = finalResult
	xmlout, _ := xml.Marshal(urlset)
	_ = xml.Unmarshal([]byte(xmlout), &urlset)
	buf, _ := xml.MarshalIndent(urlset, "", "\t")
	filename := "series_en_wyk2.xml"
	err1 := ioutil.WriteFile(filename, buf, 0644)
	if err1 != nil {
		log.Println(err1)
	}
	UploadFileToS3(buf, filename)
	c.XML(http.StatusOK, urlset)
}

//  EpisodedetailsINEnglish
// GET /api/sitemap/sitemap_shows-en.xml
// @Summary Get All series details IN English-XML file
// @Description Get All series details IN English-XML file
// @Tags sitemap
// @Accept xml
// @Produce  xml
// @Success 200 {array} object c.xml
// @Router /api/sitemap/sitemap_shows-en.xml [GET]
//Get All Episode Titles In English-xml file
func (hs *HandlerService) GetAllEpisodedetailsINEnglishXMLfile(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var activeEpisodeDetailsEnglish []activeEpisodeDetailsEnglish
	var result EpisodeUrl
	var finalResult []EpisodeUrl
	var urlset Episode

	url := os.Getenv("SITE_URL") + "/en/player/episode/"
	if err := db.Debug().Raw("select distinct e.episode_key, cpi.transliterated_title from season s join episode e on e.season_id = s.id join content_primary_info cpi on cpi.id = e.primary_info_id join playback_item pi1 on pi1.id = e.playback_item_id join content_rights cr on cr.id = s.rights_id where e.status = 1 and e.deleted_by_user_id is null and s.deleted_by_user_id is null and (pi1.scheduling_date_time <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) order by e.episode_key asc").Find(&activeEpisodeDetailsEnglish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	urlset.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	urlset.SchemaLocation = "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	urlset.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	for _, value := range activeEpisodeDetailsEnglish {
		result.Location = url + value.EpisodeKey + "/" + value.TransliteratedTitle
		result.Changefreqency = "monthly"
		result.Priority = "0.5"
		finalResult = append(finalResult, result)
	}
	urlset.URL = finalResult
	xmlout, _ := xml.Marshal(urlset)
	_ = xml.Unmarshal([]byte(xmlout), &urlset)
	buf, _ := xml.MarshalIndent(urlset, "", "\t")
	filename := "shows_en_wyk2.xml"
	err1 := ioutil.WriteFile(filename, buf, 0644)
	if err1 != nil {
		log.Println(err1)
	}
	UploadFileToS3(buf, filename)
	c.XML(http.StatusOK, urlset)
}

//  EpisodedetailsINArabic
// GET /api/sitemap/sitemap_shows-ar.xml
// @Summary Get All series details IN Arabic-XML file
// @Description Get All series details IN Arabic-XML file
// @Tags sitemap
// @Accept xml
// @Produce  xml
// @Success 200 {array} object c.xml
// @Router /api/sitemap/sitemap_shows-ar.xml [GET]
//Get All Episode Titles In Arabic-xml file
func (hs *HandlerService) GetAllEpisodedetailsINArabicXMLfile(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var activeEpisodeDetailsEnglish []activeEpisodeDetailsEnglish
	var result EpisodeUrl
	var finalResult []EpisodeUrl
	var urlset Episode

	url := os.Getenv("SITE_URL") + "/ar/player/episode/"
	if err := db.Debug().Raw("select e.episode_key, cpi.arabic_title from season s join episode e on e.season_id = s.id join content_primary_info cpi on cpi.id = e.primary_info_id join playback_item pi1 on pi1.id = e.playback_item_id join content_rights cr on cr.id = s.rights_id where e.status = 1 and e.deleted_by_user_id is null and s.deleted_by_user_id is null and (pi1.scheduling_date_time <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) order by e.episode_key asc").Find(&activeEpisodeDetailsEnglish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	urlset.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	urlset.SchemaLocation = "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	urlset.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	for _, value := range activeEpisodeDetailsEnglish {
		result.Location = url + value.EpisodeKey + "/" + value.ArabicTitle
		result.Changefreqency = "monthly"
		result.Priority = "0.5"
		finalResult = append(finalResult, result)
	}
	urlset.URL = finalResult
	xmlout, _ := xml.Marshal(urlset)
	_ = xml.Unmarshal([]byte(xmlout), &urlset)
	buf, _ := xml.MarshalIndent(urlset, "", "\t")
	filename := "shows_ar_wyk2.xml"
	err1 := ioutil.WriteFile(filename, buf, 0644)
	if err1 != nil {
		log.Println(err1)
	}
	UploadFileToS3(buf, filename)
	c.XML(http.StatusOK, urlset)
}

func UploadFileToS3(buffer []byte, filename string) (string, error) {
	s, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYOGUWMUMGAQLMW3U",                     // id
			"Jb0NV2eHwXAJg6UADb5vs3BgAyuUsvhgREi/hWRj", // secret
			""), // token can be left blank for now
	})

	tempFileName := "sitemap-qa/" + filename
	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("S3_SITEMAPBUCKET")),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		StorageClass:         aws.String("STANDARD"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		fmt.Printf("Unable to upload %q, %v", tempFileName, err)
	} else {
		fmt.Printf("Successfully uploaded %q", tempFileName)
	}
	return tempFileName, err

}
