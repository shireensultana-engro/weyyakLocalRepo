package menu

import (
	"context"
	"encoding/json"
	"fmt"
	_ "fmt"
	"frontend_service/common"
	"frontend_service/pagination"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	r.POST("page/response", hs.CreatePageResponse)
	r.POST("slider/response", hs.CreateSliderResponse)
	r.POST("playlist/response", hs.CreatePlaylistResponse)
	route := r.Group("/v1")
	route.Use(common.ValidateToken())
	route.GET("/:lang/menu/:pagekey", hs.GetMenuDetails)
	r.GET("/v1/:lang/menu", hs.GetSideMenuPageDetails)
	route.GET("/:lang/contenttype", hs.GetTopMenuDetails)
}
func (hs *HandlerService) CreatePageResponse(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	ctx := context.Background()
	tx := db.BeginTx(ctx, nil)
	var request PageDataSyncRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var pageFragment PageFragment
	var pageFragments []interface{}
	pageId := request.PageId
	pageKey := request.PageKey
	country := request.Country
	pageResponse := make(map[string]interface{})
	pageResponse["en_response"] = request.PageDetails.En
	pageResponse["ar_response"] = request.PageDetails.Ar
	pagedata, _ := json.Marshal(pageResponse)
	var pr FragmentResponse
	json.Unmarshal(pagedata, &pr)
	for _, details := range request.PageOrder {
		language := "en"
		for i := 0; i < 2; i++ {
			if i == 1 {
				language = "ar"
				pageFragment.Details = pr.ArResponseData
			} else {
				pageFragment.Details = pr.EnResponseData
			}
			pageFragment.PageId = pageId
			pageFragment.PageOrder = details.PageOrderNumber
			pageFragment.PageKey = pageKey
			pageFragment.Country = country
			pageFragment.Platform = common.DeviceNames(details.TargetPlarform)
			pageFragment.Language = language
			pageFragments = append(pageFragments, pageFragment)
		}
	}
	if err := tx.Where("page_id=? and country=?", pageId, country).Delete(&pageFragment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	err = gormbulk.BulkInsert(tx, pageFragments, common.BULK_INSERT_LIMIT)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	if request.SliderDetails != nil {
		for _, slider := range request.SliderDetails {
			sliderId := slider.SliderId
			sliderResponse := make(map[string]interface{})
			sliderResponse["en_response"] = slider.Details.En
			sliderResponse["ar_response"] = slider.Details.Ar
			sliderdata, _ := json.Marshal(sliderResponse)
			var sr FragmentResponse
			json.Unmarshal(sliderdata, &sr)
			var sliderFragment SliderFragment
			var sliderFragments []interface{}
			for _, details := range request.PageOrder {
				language := "en"
				for i := 0; i < 2; i++ {
					if i == 1 {
						language = "ar"
						sliderFragment.Details = sr.ArResponseData
					} else {
						sliderFragment.Details = sr.EnResponseData
					}
					sliderFragment.PageId = pageId
					sliderFragment.PageKey = pageKey
					sliderFragment.SliderId = sliderId
					sliderFragment.Country = country
					sliderFragment.Platform = common.DeviceNames(details.TargetPlarform)
					sliderFragment.Language = language
					sliderFragments = append(sliderFragments, sliderFragment)
				}
			}
			if err := tx.Where("page_id=? and country=? and slider_id=?", pageId, country, sliderId).Delete(&sliderFragment).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
				return
			}
			err = gormbulk.BulkInsert(tx, sliderFragments, common.BULK_INSERT_LIMIT)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
				return
			}
		}
	}

	if request.PlaylistDetails != nil {
		for _, playlist := range request.PlaylistDetails {
			var playlistFragment PlaylistFragment
			var playlistFragments []interface{}
			playlistId := playlist.PlaylistId
			playlistResponse := make(map[string]interface{})
			playlistResponse["en_response"] = playlist.Details.En
			playlistResponse["ar_response"] = playlist.Details.Ar
			playlistdata, _ := json.Marshal(playlistResponse)
			var pr FragmentResponse
			json.Unmarshal(playlistdata, &pr)
			for _, details := range request.PageOrder {
				language := "en"
				for i := 0; i < 2; i++ {
					if i == 1 {
						language = "ar"
						playlistFragment.Details = pr.ArResponseData
					} else {
						playlistFragment.Details = pr.EnResponseData
					}
					playlistFragment.PageId = pageId
					playlistFragment.PageKey = pageKey
					playlistFragment.PlaylistId = playlistId
					playlistFragment.Country = country
					playlistFragment.Platform = common.DeviceNames(details.TargetPlarform)
					playlistFragment.Language = language
					playlistFragments = append(playlistFragments, playlistFragment)
				}
			}
			if err := tx.Where("page_id=? and country=? and playlist_id=?", pageId, country, playlistId).Delete(&playlistFragment).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
				return
			}
			err = gormbulk.BulkInsert(tx, playlistFragments, common.BULK_INSERT_LIMIT)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
				return
			}
		}
	}
	err = tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Success", "Status": http.StatusOK})
	return
}
func (hs *HandlerService) CreateSliderResponse(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	ctx := context.Background()
	tx := db.BeginTx(ctx, nil)
	var request SliderDataSyncRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sliderId := request.SliderId
	country := request.Country
	var sliderFragments []interface{}
	var sliderFragment SliderFragment
	sliderResponse := make(map[string]interface{})
	sliderResponse["en_response"] = request.SliderDetails.En
	sliderResponse["ar_response"] = request.SliderDetails.Ar
	sliderdata, _ := json.Marshal(sliderResponse)
	var sr FragmentResponse
	json.Unmarshal(sliderdata, &sr)
	for _, platform := range request.PublishingPlatforms {
		for _, details := range request.SliderAvailablePages {
			language := "en"
			for i := 0; i < 2; i++ {
				if i == 1 {
					language = "ar"
					sliderFragment.Details = sr.ArResponseData
				} else {
					sliderFragment.Details = sr.EnResponseData
				}
				sliderFragment.PageId = details.PageId
				sliderFragment.PageKey = details.PageKey
				sliderFragment.SliderId = sliderId
				sliderFragment.Country = country
				sliderFragment.Platform = common.DeviceNames(platform)
				sliderFragment.Language = language
				sliderFragments = append(sliderFragments, sliderFragment)
			}
		}
	}
	if err := tx.Table("slider_fragment").Where("country=? and slider_id=?", country, sliderId).Delete(&sliderFragment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	err = gormbulk.BulkInsert(tx, sliderFragments, common.BULK_INSERT_LIMIT)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	err = tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Success", "Status": http.StatusOK})
	return
}
func (hs *HandlerService) CreatePlaylistResponse(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	ctx := context.Background()
	tx := db.BeginTx(ctx, nil)
	var request PlaylistDataSyncRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	playlistId := request.PlaylistId
	country := request.Country
	var playlistFragments []interface{}
	var playlistFragment PlaylistFragment
	playlistResponse := make(map[string]interface{})
	playlistResponse["en_response"] = request.PlaylistDetails.En
	playlistResponse["ar_response"] = request.PlaylistDetails.Ar
	playlistdata, _ := json.Marshal(playlistResponse)
	var sr FragmentResponse
	json.Unmarshal(playlistdata, &sr)
	for _, platform := range request.PublishingPlatforms {
		for _, details := range request.PlaylistAvailablePages {
			language := "en"
			for i := 0; i < 2; i++ {
				if i == 1 {
					language = "ar"
					playlistFragment.Details = sr.ArResponseData
				} else {
					playlistFragment.Details = sr.EnResponseData
				}
				playlistFragment.PageId = details.PageId
				playlistFragment.PageKey = details.PageKey
				playlistFragment.PlaylistId = playlistId
				playlistFragment.Country = country
				playlistFragment.Platform = common.DeviceNames(platform)
				playlistFragment.Language = language
				playlistFragments = append(playlistFragments, playlistFragment)
			}
		}
	}
	if err := tx.Where("country=? and playlist_id=?", country, playlistId).Delete(&playlistFragment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	err = gormbulk.BulkInsert(tx, playlistFragments, common.BULK_INSERT_LIMIT)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}
	err = tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Success", "Status": http.StatusOK})
	return
}

// @Tags Menu
// @Summary Get Menu Details By Page Id
// @Description Get Menu Details By Page Id
// @Accept  json
// @Produce json
// @Param pagekey path string true "Page Key"
// @Param lang path string true "Language Code"
// @Param country query string false "Country"
// @Success 200 {object} object	GetMenu
// @Failure 400 {object} object err
// @Failure 404 {object} object err
// @Router /v1/{lang}/menu/{pagekey} [get]
func (hs *HandlerService) GetMenuDetails(c *gin.Context) {
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	serverError := common.ServerErrorResponse(language)
	key := c.Request.Host + c.Request.URL.String()
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "?", "_")
	// var country, platform string
	// if c.Request.URL.Query()["country"] != nil {
	// 	country = strings.ToUpper(c.Request.URL.Query()["country"][0])
	// }
	// if c.Request.URL.Query()["platform"] != nil {
	// 	platform = strings.ToLower(c.Request.URL.Query()["platform"][0])
	// }
	// key := c.Param("pagekey") + c.Param("lang") + country + strings.Join(strings.Fields(platform), "")
	
	//countryRestriction
	var CountryCode string
	if c.Request.URL.Query()["country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
	}

	S3url, err := common.GetCurlCall(os.Getenv("S3_URLFORCONFIG"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	var result interface{}
	err = json.Unmarshal(S3url, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	var mapResult map[string]interface{}
	err1 := json.Unmarshal(S3url, &mapResult)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	// Check if "restrictedCountries" key exists and is of type []interface{}
	restrictedCountriesInterface, ok := mapResult["restrictedCountries"].([]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, "Invalid restrictedCountries data")
		return
	}
	// Convert the interface slice to a string slice
	var countryList []string
	for _, v := range restrictedCountriesInterface {
		if str, ok := v.(string); ok {
			countryList = append(countryList, str)
		}
	}
    isCountryRestricted := false
	fmt.Println("countryList" , countryList)
    for _, country := range countryList {
        if country == CountryCode {
            isCountryRestricted = true
            break
        }
    }
	if isCountryRestricted{
		// c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden,"message":"country restriction" })
		c.JSON(http.StatusForbidden, gin.H{"Message": "Invalid restrictedCountries data"})
		return
	}
	//countryRestriction
	
	url := os.Getenv("REDIS_CACHE_URL") + "/" + key
	fmt.Println("Redis url -", url)
	response, err := common.GetCurlCall(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var menuPage MenuPage
	type RedisCacheResponse struct {
		Value string `json:"value"`
		Error string `json:"error"`
	}
	var RedisResponse RedisCacheResponse
	json.Unmarshal(response, &RedisResponse)
	if RedisResponse.Value != "" {
		if err := json.Unmarshal([]byte(RedisResponse.Value), &menuPage); err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	} else {
		db := c.MustGet("FCDB").(*gorm.DB)
		fdb := c.MustGet("DB").(*gorm.DB)
		var blackPlaylistCount, redPlaylistCount, greenPlaylistCount int
		var featured *FeaturedDetails
		var imageryDetails ImageryDetails
		var featuredDetails FeaturedDetails
		var featuredPlaylist FeaturedPlaylists
		featuredDetails.Playlists = nil
		var countryCode string
		if c.Request.URL.Query()["country"] != nil {
			countryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
		}
		if len(countryCode) != 2 {
			countryCode = "AE"
		}
		platformName := "web"
		PageKey := c.Param("pagekey")
		//page details
		type Page struct {
			Id                     string
			PageKey                int
			EnglishPageFriendlyUrl string
			ArabicPageFriendlyUrl  string
			EnglishMetaDescription string
			ArabicMetaDescription  string
			EnglishTitle           string
			ArabicTitle            string
			PageType               int
			HasMobileMenu          bool
		}
		var details Page
		if err := db.Where("page_key=?", PageKey).Find(&details).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		type PageIds struct {
			Id string `json:"id"`
		}
		var pageids []PageIds
		var ids []string
		if err := db.Table("page p").Select("p.id").Joins("inner join page_slider ps on ps.page_id=p.id inner join slider s on s.id = ps.slider_id").Where("s.deleted_by_user_id  is null and s.is_disabled =false and p.is_disabled=false and p.deleted_by_user_id is null  and s.scheduling_start_date <=NOW() and s.scheduling_end_date >=NOW() and p.id=?", details.Id).Find(&pageids).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		if pageids != nil {
			for _, pageid := range pageids {
				ids = append(ids, pageid.Id)
			}
		}
		menuPage.ID = details.PageKey
		menuPage.FriendlyUrl = details.EnglishPageFriendlyUrl
		menuPage.SeoDescription = details.EnglishMetaDescription
		menuPage.Title = details.EnglishTitle
		if language != "en" {
			menuPage.FriendlyUrl = details.ArabicPageFriendlyUrl
			menuPage.SeoDescription = details.ArabicMetaDescription
			menuPage.Title = details.ArabicTitle
		}
		menuPage.Type = common.PageTypes(details.PageType)
		menuPage.Featured = featured
		if details.PageType != 16 && details.PageType != 8 {
			exists := common.FindString(ids, details.Id)
			if (details.PageType == 0 && exists == true) || details.PageType == 1 {
				menuPage.Type = "Home"
			} else {
				menuPage.Type = "VOD"
			}
		}
		if details.HasMobileMenu == true {
			imageryDetails.MobileMenu = os.Getenv("IMAGES") + strings.ToLower(details.Id) + "/mobile-menu" + "?c=" + strconv.Itoa(rand.Intn(200000))
			imageryDetails.MobilePosterImage = os.Getenv("IMAGES") + strings.ToLower(details.Id) + "/menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
			imageryDetails.MobileMenuPosterImage = os.Getenv("IMAGES") + strings.ToLower(details.Id) + "/mobile-menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
		} else {
			imageryDetails.MobileMenu = ""
			imageryDetails.MobilePosterImage = ""
			imageryDetails.MobileMenuPosterImage = ""
		}

		menuPage.Imagery = imageryDetails
		//page slider details
		var slider Slider
		if err := db.Select("s.*").Table("slider s").Joins("inner join page_slider ps on ps.slider_id=s.id").Where("s.deleted_by_user_id  is null and s.is_disabled =false and ps.page_id=? and (s.scheduling_start_date <=NOW() or ps.order =0) and (s.scheduling_end_date >=NOW()  or ps.order =0)", details.Id).Limit(1).Order("ps.order desc").Find(&slider).Error; err != nil && err.Error() != "record not found" {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		if slider.SliderKey != 0 {
			var featuredPlaylists []FeaturedPlaylists
			featuredDetails.ID = int64(slider.SliderKey)
			featuredDetails.Type = common.SliderTypes(slider.Type)
			if slider.BlackAreaPlaylistId != "" || slider.RedAreaPlaylistId != "" || slider.GreenAreaPlaylistId != "" {
				playlists, _ := SliderPlaylists(slider.BlackAreaPlaylistId, slider.RedAreaPlaylistId, slider.GreenAreaPlaylistId, c)
				for _, playlist := range playlists {
					featuredPlaylist.ID = int32(playlist.PlaylistKey)
					featuredPlaylist.PlaylistType = playlist.PlaylistType
					contentIds, err := PlaylistItemContents(playlist.ID, c)
					if err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
					var Ids []string
					for _, content := range contentIds {
						Ids = append(Ids, content.ContentId)
					}
					var playlistContent PlaylistContent
					var playlistContents []PlaylistContent
					type ContentFragmentDetails struct {
						Details string `json:"details"`
					}
					var contentFragmentDetails []ContentFragmentDetails
					if err := fdb.Table("content_fragment").Select("details::text as details").Where("content_id in(?) and country=? and language=? and platform=?", Ids, countryCode, language, platformName).Find(&contentFragmentDetails).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
					for _, content := range contentFragmentDetails {
						if err := json.Unmarshal([]byte(content.Details), &playlistContent); err != nil {
							c.JSON(http.StatusInternalServerError, serverError)
							return
						}
						if playlistContent.ID != 0 && playlistContent.ID > 0 {
							playlistContents = append(playlistContents, playlistContent)
						}
					}
					if playlistContents == nil {
						continue
					}
					if playlist.PlaylistType == "black_playlist" {
						blackPlaylistCount = len(playlistContents)
					} else if playlist.PlaylistType == "red_playlist" {
						redPlaylistCount = len(playlistContents)
					} else {
						greenPlaylistCount = len(playlistContents)
					}
					var contents []PlaylistContent
					for _, id := range contentIds {
						for _, content := range playlistContents {
							if id.ContentId == content.ContentId {
								contents = append(contents, content)
							}
						}
					}
					featuredPlaylist.Content = contents
					featuredPlaylists = append(featuredPlaylists, featuredPlaylist)
				}
			}
			if details.PageType == 1 && blackPlaylistCount >= common.BlackPlaylistCount && redPlaylistCount == common.RedPlaylistCount && greenPlaylistCount >= common.GreenPlaylistCount {
				featuredDetails.Playlists = featuredPlaylists
			} else if details.PageType != 1 {
				featuredDetails.Playlists = featuredPlaylists
			}
			if len(featuredDetails.Playlists) > 0 {
				menuPage.Featured = &featuredDetails
			}
		}
		//page palylist details
		var playlists []Playlist
		fields := ",english_title"
		if language != "en" {
			fields = ",arabic_title"
		}
		if err := db.Select("p.id"+fields+",p.scheduling_start_date,p.scheduling_end_date,p.deleted_by_user_id,p.is_disabled,p.created_at,p.playlist_key,p.modified_at,p.playlist_type").Table("page_playlist pp").Joins("join playlist p on p.id =pp.playlist_id").Where("p.is_disabled =false and p.deleted_by_user_id is null and pp.page_id =? and (p.scheduling_start_date <=now() or p.scheduling_start_date is null) and (p.scheduling_end_date >=now() or p.scheduling_end_date is null)", details.Id).Order("pp.order asc").Find(&playlists).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var pagePlaylists []MenuPlaylists
		var wg sync.WaitGroup
		for _, playlist := range playlists {
			pagePlaylist := MenuPlaylists{}
			pagePlaylist.ID = int32(playlist.PlaylistKey)
			pagePlaylist.PlaylistType = playlist.PlaylistType

			pagePlaylist.Content = []PlaylistContent{}
			pagePlaylist.PageContent = []PageContent{}
			pagePlaylist.Title = playlist.EnglishTitle
			if language != "en" {
				pagePlaylist.Title = playlist.ArabicTitle
			}
			if playlist.PlaylistType == "pagecontent" {
				pagePlaylist.Title = "null"
				playlistPages := make(chan []PageContent)
				wg.Add(1)
				go PlaylistPages(playlist.ID, language, playlistPages, c)
				defer wg.Done()
				if playlistPages != nil {
					pagePlaylist.PageContent = <-playlistPages
				}
			} else if playlist.PlaylistType == "content" {
				contentIds, err := PlaylistItemContents(playlist.ID, c)
				if err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
				var Ids []string
				for _, content := range contentIds {
					Ids = append(Ids, content.ContentId)
				}
				var playlistContents []PlaylistContent
				type ContentFragmentDetails struct {
					Details string `json:"details"`
				}
				var contentFragmentDetails []ContentFragmentDetails
				if err := fdb.Table("content_fragment").Select("details::text as details").Where("content_id in(?) and country=? and language=? and platform=?", Ids, countryCode, language, platformName).Find(&contentFragmentDetails).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
				for _, content := range contentFragmentDetails {
					var playlistContent PlaylistContent
					if err := json.Unmarshal([]byte(content.Details), &playlistContent); err != nil {
						c.JSON(http.StatusInternalServerError, err.Error())
						return
					}
					if playlistContent.ID != 0 {
						playlistContents = append(playlistContents, playlistContent)
					}
				}
				if playlistContents == nil {
					continue
				}
				var contents []PlaylistContent
				for _, id := range contentIds {
					for _, content := range playlistContents {
						if id.ContentId == content.ContentId {
							contents = append(contents, content)
						}
					}
				}
				pagePlaylist.Content = contents
			}
			pagePlaylists = append(pagePlaylists, pagePlaylist)
		}
		menuPage.Playlists = pagePlaylists
		jsonData, _ := json.Marshal(menuPage)
		type RedisCacheRequest struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		var request RedisCacheRequest
		url := os.Getenv("REDIS_CACHE_URL")
		request.Key = key
		request.Value = string(jsonData)
		_, err := common.PostCurlCall("POST", url, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	}
	fmt.Println("menupage..", menuPage)
	c.JSON(http.StatusOK, gin.H{"data": menuPage})
	return
}

// func (hs *HandlerService) GetMenuDetails(c *gin.Context) {
// 	db := c.MustGet("DB").(*gorm.DB)
// 	var menu GetMenu
// 	var page MenuPageDetails
// 	var slider *FeaturedDetails
// 	var playlists []MenuPlaylists
// 	var CountryCode string
// 	if c.Request.URL.Query()["country"] != nil {
// 		CountryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
// 	}
// 	if len(CountryCode) != 2 {
// 		CountryCode = "AE"
// 	}
// 	PageKey := c.Param("pagekey")
// 	language := strings.ToLower(c.Param("lang"))
// 	if language != "en" {
// 		language = "ar"
// 	}
// 	serverError := common.ServerErrorResponse(language)
// 	type PageFragmentDetails struct {
// 		Page string `json:"page"`
// 	}
// 	where := "country='" + CountryCode + "' and page_key =" + PageKey + " and platform ='web' and language ='" + language + "'"
// 	var pageFragmentDetails PageFragmentDetails
// 	if err := db.Table("page_fragment").Select("details::text as page").Where(where).Find(&pageFragmentDetails).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, serverError)
// 		return
// 	}
// 	if err := json.Unmarshal([]byte(pageFragmentDetails.Page), &page); err != nil {
// 		c.JSON(http.StatusInternalServerError, serverError)
// 		return

// 	}
// 	type SliderFragmentDetails struct {
// 		Slider string `json:"slider"`
// 	}
// 	var sliderFragmentDetails SliderFragmentDetails
// 	if err := db.Table("slider_fragment").Select("details::text as slider").Where(where).Find(&sliderFragmentDetails).Error; err == nil {
// 		if err := json.Unmarshal([]byte(sliderFragmentDetails.Slider), &slider); err != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return

// 		}
// 		if slider.Playlists != nil {
// 			if slider.Playlists[0].Content != nil {
// 				page.Featured = slider
// 			}
// 		}
// 	}

// 	type PlaylistFragmentDetails struct {
// 		Playlists string `json:"playlists"`
// 	}
// 	var playlistDetails []PlaylistFragmentDetails
// 	if err := db.Table("playlist_fragment").Select("details::text as playlists").Where("country=? and page_key =? and platform ='web' and language =?", CountryCode, PageKey, language).Find(&playlistDetails).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, serverError)
// 		return
// 	}
// 	for _, details := range playlistDetails {
// 		var playlist MenuPlaylists
// 		if err := json.Unmarshal([]byte(details.Playlists), &playlist); err != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		if playlist.Content != nil || playlist.PageContent != nil {
// 			playlists = append(playlists, playlist)
// 		}
// 	}
// 	page.Playlists = playlists
// 	menu.Data = page
// 	c.JSON(http.StatusOK, menu)
// 	return
// }

// GetSideMenuPageDetails -  Get side  menu details
// GET /v1/:lang/menu
// @Summary Get side menu pages details
// @Description Get side menu pages details
// @Tags Menu
// @Accept  json
// @Produce  json
// @Param lang path string true "Language Code"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Param device query string true "Device"
// @Param country query string false "Country"
// @Param cascade query string false "Cascade"
// @Success 200 {array} object c.JSON
// @Router /v1/{lang}/menu [get]
func (hs *HandlerService) GetSideMenuPageDetails(c *gin.Context) {
	db := c.MustGet("FCDB").(*gorm.DB)
	fdb := c.MustGet("DB").(*gorm.DB)
	//	notFound := common.NotFoundErrorResponse()
	var limit, offset, current_page int64
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["page"] != nil {
		current_page, _ = strconv.ParseInt(c.Request.URL.Query()["page"][0], 10, 64)
	}
	if limit == 0 {
		limit, _ = strconv.ParseInt(os.Getenv("DEFAULT_PAGE_SIZE"), 10, 64)
	}
	offset = current_page * limit

	var CountryCode, Casecade string
	if c.Request.URL.Query()["device"] == nil || c.Request.URL.Query()["device"][0] == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}
	if c.Request.URL.Query()["country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
	}
	if c.Request.URL.Query()["cascade"] != nil {
		Casecade = c.Request.URL.Query()["cascade"][0]
	}
	DeviceName := strings.ToLower(c.Request.URL.Query()["device"][0])
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}

	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	serverError := common.ServerErrorResponse(language)
	key := c.Request.Host + c.Request.URL.String()
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "?", "_")
	url := os.Getenv("REDIS_CACHE_URL") + "/" + key
	fmt.Println("Redis url -", url)
	response, err := common.GetCurlCall(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	type RedisCacheResponse struct {
		Value string `json:"value"`
		Error string `json:"error"`
	}
	var RedisResponse RedisCacheResponse
	json.Unmarshal(response, &RedisResponse)
	if RedisResponse.Value != "" {
		var menu SideMenuDetails
		if err := json.Unmarshal([]byte(RedisResponse.Value), &menu); err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		c.JSON(http.StatusOK, menu)
		return
	} else {
		var menuPageDetails []MenuPage
		var menu SideMenuDetails
		//page details
		type Page struct {
			Id                     string
			PageKey                int
			EnglishPageFriendlyUrl string
			ArabicPageFriendlyUrl  string
			EnglishMetaDescription string
			ArabicMetaDescription  string
			EnglishTitle           string
			ArabicTitle            string
			PageType               int
			HasMobileMenu          bool
		}
		var pageDetails []Page
		var myplaylistpagedetials Page
		rows := db.Table("page p").Select("p.*").Joins("inner join page_target_platform ptp on ptp.page_id=p.id").Where("p.is_disabled=false and p.deleted_by_user_id is null and ptp.target_platform=(select id from publish_platform pp where lower(platform)=lower(?))", DeviceName).Order("ptp.page_order_number asc").Find(&pageDetails).RowsAffected
		changedlimit := limit - 1
		if err := db.Table("page p").Select("p.*").Joins("inner join page_target_platform ptp on ptp.page_id=p.id").Where("p.is_disabled=false and p.page_type != 16 and p.deleted_by_user_id is null and ptp.target_platform=(select id from publish_platform pp where lower(platform)=lower(?))", DeviceName).Order("ptp.page_order_number asc").Limit(changedlimit).Offset(offset).Find(&pageDetails).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		// for fetching my playlist at last
		if err := db.Table("page p").Select("p.*").Joins("inner join page_target_platform ptp on ptp.page_id=p.id").Where("p.is_disabled=false and p.page_type = 16 and p.deleted_by_user_id is null and ptp.target_platform=(select id from publish_platform pp where lower(platform)=lower(?))", DeviceName).Order("ptp.page_order_number asc").Find(&myplaylistpagedetials).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		type PageIds struct {
			Id string `json:"id"`
		}
		var pageids []PageIds
		var ids []string
		if err := db.Table("page p").Select("p.id").Joins("inner join page_slider ps on ps.page_id=p.id inner join slider s on s.id = ps.slider_id").Where("s.deleted_by_user_id  is null and s.is_disabled =false and p.is_disabled=false and p.deleted_by_user_id is null  and s.scheduling_start_date <=NOW() and s.scheduling_end_date >=NOW()").Find(&pageids).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		if pageids != nil {
			for _, pageid := range pageids {
				ids = append(ids, pageid.Id)
			}
		}
		for _, details := range pageDetails {
			var SliderCount int
			var PlaylistCount int
			db.Table("page_slider").Select("slider_id").Where("page_id= ?", details.Id).Count(&SliderCount)
			db.Table("page_playlist").Select("playlist_id").Where("page_id= ?", details.Id).Count(&PlaylistCount)
			fmt.Println("printing details before", details.EnglishTitle)
			if SliderCount > 0 || PlaylistCount > 0 || details.EnglishTitle == "My Playlist" || details.ArabicTitle == "قائمتي" {
				fmt.Println("PRINTING DETAILS HERE", details.EnglishTitle)
				var menuPage MenuPage
				var blackPlaylistCount, redPlaylistCount, greenPlaylistCount int
				var featured *FeaturedDetails
				var imageryDetails ImageryDetails
				var featuredDetails FeaturedDetails
				var featuredPlaylist FeaturedPlaylists
				featuredDetails.Playlists = nil
				var countryCode string
				if c.Request.URL.Query()["country"] != nil {
					countryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
				}
				if len(countryCode) != 2 {
					countryCode = "AE"
				}
				platformName := "web"
				menuPage.ID = details.PageKey
				menuPage.FriendlyUrl = details.EnglishPageFriendlyUrl
				menuPage.SeoDescription = details.EnglishMetaDescription
				menuPage.Title = details.EnglishTitle
				if language != "en" {
					menuPage.FriendlyUrl = details.ArabicPageFriendlyUrl
					menuPage.SeoDescription = details.ArabicMetaDescription
					menuPage.Title = details.ArabicTitle
				}
				menuPage.Type = common.PageTypes(details.PageType)
				menuPage.Featured = featured
				if details.PageType != 16 && details.PageType != 8 {
					exists := common.FindString(ids, details.Id)
					if (details.PageType == 0 && exists == true) || details.PageType == 1 {
						menuPage.Type = "Home"
					} else {
						menuPage.Type = "VOD"
					}
				}
				if details.HasMobileMenu == true {
					imageryDetails.MobileMenu = os.Getenv("IMAGES") + strings.ToLower(details.Id) + "/mobile-menu" + "?c=" + strconv.Itoa(rand.Intn(200000))
					imageryDetails.MobilePosterImage = os.Getenv("IMAGES") + strings.ToLower(details.Id) + "/menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
					imageryDetails.MobileMenuPosterImage = os.Getenv("IMAGES") + strings.ToLower(details.Id) + "/mobile-menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
				} else {
					imageryDetails.MobileMenu = ""
					imageryDetails.MobilePosterImage = ""
					imageryDetails.MobileMenuPosterImage = ""
				}
				menuPage.Imagery = imageryDetails
				if Casecade == "2" && DeviceName == "smarttv" {
					//page slider details
					var slider Slider
					if err := db.Select("s.*").Table("slider s").Joins("inner join page_slider ps on ps.slider_id=s.id").Where("s.deleted_by_user_id  is null and s.is_disabled =false and ps.page_id=? and (s.scheduling_start_date <=NOW() or ps.order =0) and (s.scheduling_end_date >=NOW()  or ps.order =0)", details.Id).Limit(1).Order("ps.order desc").Find(&slider).Error; err != nil && err.Error() != "record not found" {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
					if slider.SliderKey != 0 {
						var featuredPlaylists []FeaturedPlaylists
						featuredDetails.ID = int64(slider.SliderKey)
						featuredDetails.Type = common.SliderTypes(slider.Type)
						if slider.BlackAreaPlaylistId != "" || slider.RedAreaPlaylistId != "" || slider.GreenAreaPlaylistId != "" {
							playlists, _ := SliderPlaylists(slider.BlackAreaPlaylistId, slider.RedAreaPlaylistId, slider.GreenAreaPlaylistId, c)
							for _, playlist := range playlists {
								featuredPlaylist.ID = int32(playlist.PlaylistKey)
								featuredPlaylist.PlaylistType = playlist.PlaylistType
								contentIds, err := PlaylistItemContents(playlist.ID, c)
								if err != nil {
									c.JSON(http.StatusInternalServerError, serverError)
									return
								}
								var Ids []string
								for _, content := range contentIds {
									Ids = append(Ids, content.ContentId)
								}
								var playlistContent PlaylistContent
								var playlistContents []PlaylistContent
								type ContentFragmentDetails struct {
									Details string `json:"details"`
								}
								var contentFragmentDetails []ContentFragmentDetails
								if err := fdb.Table("content_fragment").Select("details::text as details").Where("content_id in(?) and country=? and language=? and platform=?", Ids, countryCode, language, platformName).Find(&contentFragmentDetails).Error; err != nil {
									c.JSON(http.StatusInternalServerError, serverError)
									return
								}
								for _, content := range contentFragmentDetails {
									if err := json.Unmarshal([]byte(content.Details), &playlistContent); err != nil {
										c.JSON(http.StatusInternalServerError, serverError)
										return
									}
									if playlistContent.ID != 0 && playlistContent.ID > 0 {
										playlistContents = append(playlistContents, playlistContent)
									}
								}
								if playlistContents == nil {
									continue
								}
								if playlist.PlaylistType == "black_playlist" {
									blackPlaylistCount = len(playlistContents)
								} else if playlist.PlaylistType == "red_playlist" {
									redPlaylistCount = len(playlistContents)
								} else {
									greenPlaylistCount = len(playlistContents)
								}
								var contents []PlaylistContent
								for _, id := range contentIds {
									for _, content := range playlistContents {
										if id.ContentId == content.ContentId {
											contents = append(contents, content)
										}
									}
								}
								featuredPlaylist.Content = contents
								featuredPlaylists = append(featuredPlaylists, featuredPlaylist)
							}
						}
						if details.PageType == 1 && blackPlaylistCount >= common.BlackPlaylistCount && redPlaylistCount == common.RedPlaylistCount && greenPlaylistCount >= common.GreenPlaylistCount {
							featuredDetails.Playlists = featuredPlaylists
						} else if details.PageType != 1 {
							featuredDetails.Playlists = featuredPlaylists
						}
						if len(featuredDetails.Playlists) > 0 {
							menuPage.Featured = &featuredDetails
						}
					}
					//page palylist details
					var playlists []Playlist
					fields := ",english_title"
					if language != "en" {
						fields = ",arabic_title"
					}
					if err := db.Select("p.id"+fields+",p.scheduling_start_date,p.scheduling_end_date,p.deleted_by_user_id,p.is_disabled,p.created_at,p.playlist_key,p.modified_at,p.playlist_type").Table("page_playlist pp").Joins("join playlist p on p.id =pp.playlist_id").Where("p.is_disabled =false and p.deleted_by_user_id is null and pp.page_id =? and (p.scheduling_start_date <=now() or p.scheduling_start_date is null) and (p.scheduling_end_date >=now() or p.scheduling_end_date is null)", details.Id).Order("pp.order asc").Find(&playlists).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
					var pagePlaylists []MenuPlaylists
					var wg sync.WaitGroup
					for _, playlist := range playlists {
						pagePlaylist := MenuPlaylists{}
						pagePlaylist.ID = int32(playlist.PlaylistKey)
						pagePlaylist.PlaylistType = playlist.PlaylistType

						pagePlaylist.Content = []PlaylistContent{}
						pagePlaylist.PageContent = []PageContent{}
						pagePlaylist.Title = playlist.EnglishTitle
						if language != "en" {
							pagePlaylist.Title = playlist.ArabicTitle
						}
						if playlist.PlaylistType == "pagecontent" {
							pagePlaylist.Title = "null"
							playlistPages := make(chan []PageContent)
							wg.Add(1)
							go PlaylistPages(playlist.ID, language, playlistPages, c)
							defer wg.Done()
							if playlistPages != nil {
								pagePlaylist.PageContent = <-playlistPages
							}
						} else if playlist.PlaylistType == "content" {
							contentIds, err := PlaylistItemContents(playlist.ID, c)
							if err != nil {
								c.JSON(http.StatusInternalServerError, serverError)
								return
							}
							var Ids []string
							for _, content := range contentIds {
								Ids = append(Ids, content.ContentId)
							}
							var playlistContents []PlaylistContent
							type ContentFragmentDetails struct {
								Details string `json:"details"`
							}
							var contentFragmentDetails []ContentFragmentDetails
							if err := fdb.Table("content_fragment").Select("details::text as details").Where("content_id in(?) and country=? and language=? and platform=?", Ids, countryCode, language, platformName).Find(&contentFragmentDetails).Error; err != nil {
								c.JSON(http.StatusInternalServerError, serverError)
								return
							}
							for _, content := range contentFragmentDetails {
								var playlistContent PlaylistContent
								if err := json.Unmarshal([]byte(content.Details), &playlistContent); err != nil {
									c.JSON(http.StatusInternalServerError, err.Error())
									return
								}
								if playlistContent.ID != 0 {
									playlistContents = append(playlistContents, playlistContent)
								}
							}
							if playlistContents == nil {
								continue
							}
							var contents []PlaylistContent
							for _, id := range contentIds {
								for _, content := range playlistContents {
									if id.ContentId == content.ContentId {
										contents = append(contents, content)
									}
								}
							}
							pagePlaylist.Content = contents
						}
						pagePlaylists = append(pagePlaylists, pagePlaylist)
					}
					menuPage.Playlists = pagePlaylists
				}
				menuPageDetails = append(menuPageDetails, menuPage)
			}
		}
		// sending my playlist details here seperately
		var myplaylistpage MenuPage
		myplaylistpage.FriendlyUrl = myplaylistpagedetials.EnglishPageFriendlyUrl
		myplaylistpage.SeoDescription = myplaylistpagedetials.EnglishMetaDescription
		myplaylistpage.Title = myplaylistpagedetials.EnglishTitle
		if language != "en" {
			myplaylistpage.FriendlyUrl = myplaylistpagedetials.ArabicPageFriendlyUrl
			myplaylistpage.SeoDescription = myplaylistpagedetials.ArabicMetaDescription
			myplaylistpage.Title = myplaylistpagedetials.ArabicTitle
		}
		myplaylistpage.ID = myplaylistpagedetials.PageKey
		if myplaylistpagedetials.HasMobileMenu == true {
			myplaylistpage.Imagery.MobileMenu = os.Getenv("IMAGES") + strings.ToLower(myplaylistpagedetials.Id) + "/mobile-menu" + "?c=" + strconv.Itoa(rand.Intn(200000))
			myplaylistpage.Imagery.MobilePosterImage = os.Getenv("IMAGES") + strings.ToLower(myplaylistpagedetials.Id) + "/menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
			myplaylistpage.Imagery.MobileMenuPosterImage = os.Getenv("IMAGES") + strings.ToLower(myplaylistpagedetials.Id) + "/mobile-menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
		} else {
			myplaylistpage.Imagery.MobileMenu = ""
			myplaylistpage.Imagery.MobilePosterImage = ""
			myplaylistpage.Imagery.MobileMenuPosterImage = ""
		}
		myplaylistpage.SeoDescription = ""
		myplaylistpage.Title = myplaylistpagedetials.EnglishTitle
		if language != "en" {
			myplaylistpage.Title = myplaylistpagedetials.ArabicTitle
		}
		myplaylistpage.Type = common.PageTypes(myplaylistpagedetials.PageType)
		fmt.Println("APPEND LAST", myplaylistpage.Title, myplaylistpage.ID)
		menuPageDetails = append(menuPageDetails, myplaylistpage)
		lastPage := rows / limit
		menu.Total = rows
		menu.PerPage = limit
		menu.CurrentPage = current_page
		menu.LastPage = lastPage
		Host := ""
		if c.Request.Host == "localhost:3003" {
			Host = "http://" + c.Request.Host
		} else {
			Host = os.Getenv("BASE_URL")
		}
		if current_page < lastPage {
			menu.NextPageUrl = Host + "v1/" + language + "/menu?device=" + DeviceName + "&cascade=" + Casecade + "&limit=" + strconv.FormatInt(limit, 10) + "&page=" + strconv.FormatInt(current_page+1, 10)
		} else {
			menu.NextPageUrl = ""
		}
		if current_page-1 > 0 {
			menu.PrevPageUrl = Host + "v1/" + language + "/menu?device=" + DeviceName + "&cascade=" + Casecade + "&limit=" + strconv.FormatInt(limit, 10) + "&page=" + strconv.FormatInt(current_page-1, 10)
		} else {
			menu.PrevPageUrl = ""
		}
		menu.From = offset
		menu.To = offset + limit
		menu.Data = menuPageDetails
		jsonData, _ := json.Marshal(menu)
		type RedisCacheRequest struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		var request RedisCacheRequest
		url := os.Getenv("REDIS_CACHE_URL")
		request.Key = key
		request.Value = string(jsonData)
		_, err := common.PostCurlCall("POST", url, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		c.JSON(http.StatusOK, menu)
		return
	}
}

// func (hs *HandlerService) GetSideMenuPageDetails(c *gin.Context) {
// 	db := c.MustGet("DB").(*gorm.DB)
// 	//	notFound := common.NotFoundErrorResponse()
// 	var limit, offset, current_page int64
// 	if c.Request.URL.Query()["limit"] != nil {
// 		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
// 	}
// 	if c.Request.URL.Query()["page"] != nil {
// 		current_page, _ = strconv.ParseInt(c.Request.URL.Query()["page"][0], 10, 64)
// 	}
// 	if limit == 0 {
// 		limit, _ = strconv.ParseInt(os.Getenv("DEFAULT_PAGE_SIZE"), 10, 64)
// 	}
// 	offset = current_page * limit

// 	var CountryCode, Casecade string
// 	if c.Request.URL.Query()["device"] == nil || c.Request.URL.Query()["device"][0] == "" {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
// 		return
// 	}
// 	if c.Request.URL.Query()["country"] != nil {
// 		CountryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
// 	}
// 	if c.Request.URL.Query()["cascade"] != nil {
// 		Casecade = c.Request.URL.Query()["cascade"][0]
// 	}
// 	DeviceName := strings.ToLower(c.Request.URL.Query()["device"][0])
// 	var menu SideMenuDetails
// 	var page MenuPageDetails
// 	var pageDetails []MenuPageDetails
// 	// var slider *FeaturedDetails
// 	var pageFragment []GetPageFragment
// 	if len(CountryCode) != 2 {
// 		CountryCode = "AE"
// 	}

// 	language := strings.ToLower(c.Param("lang"))
// 	if language != "en" {
// 		language = "ar"
// 	}

// 	serverError := common.ServerErrorResponse(language)
// 	type PageIds struct {
// 		Id string `json:"id"`
// 	}
// 	rows := db.Table("page_fragment").Where("country=? and platform =? and language =?", CountryCode, DeviceName, language).Order("page_order").Find(&pageFragment).RowsAffected
// 	if err := db.Table("page_fragment").Where("country=? and platform =? and language =?", CountryCode, DeviceName, language).Order("page_order").Limit(limit).Offset(offset).Find(&pageFragment).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, serverError)
// 		return
// 	}
// 	for _, pageData := range pageFragment {
// 		var slider *FeaturedDetails
// 		if err := json.Unmarshal([]byte(pageData.Details), &page); err != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return

// 		}
// 		if Casecade == "2" && DeviceName == "smarttv" {
// 			type SliderFragmentDetails struct {
// 				Slider string `json:"slider"`
// 			}
// 			page.Featured = nil
// 			var sliderFragmentDetails SliderFragmentDetails
// 			if err := db.Table("slider_fragment").Select("details::text as slider").Where("country=? and page_key =? and platform ='smarttv' and language =?", CountryCode, pageData.PageKey, language).Find(&sliderFragmentDetails).Error; err == nil {
// 				if err := json.Unmarshal([]byte(sliderFragmentDetails.Slider), &slider); err != nil {
// 					c.JSON(http.StatusInternalServerError, serverError)
// 					return

// 				}
// 				if slider.Playlists != nil {
// 					if slider.Playlists[0].Content != nil {
// 						page.Featured = slider
// 					}
// 				}
// 			}

// 			type PlaylistFragmentDetails struct {
// 				Playlists string `json:"playlists"`
// 			}
// 			var playlistDetails []PlaylistFragmentDetails
// 			if err := db.Table("playlist_fragment").Select("details::text as playlists").Where("country=? and page_key =? and platform ='smarttv' and language =?", CountryCode, pageData.PageKey, language).Find(&playlistDetails).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 			var playlists []MenuPlaylists
// 			for _, details := range playlistDetails {
// 				var playlist MenuPlaylists
// 				if err := json.Unmarshal([]byte(details.Playlists), &playlist); err != nil {
// 					c.JSON(http.StatusInternalServerError, serverError)
// 					return
// 				}
// 				if playlist.Content != nil || playlist.PageContent != nil {
// 					playlists = append(playlists, playlist)
// 				}
// 			}
// 			page.Playlists = playlists
// 		}
// 		pageDetails = append(pageDetails, page)
// 	}

// 	lastPage := rows / limit
// 	menu.Total = rows
// 	menu.PerPage = limit
// 	menu.CurrentPage = current_page
// 	menu.LastPage = lastPage
// 	if current_page < lastPage {
// 		menu.NextPageUrl = os.Getenv("BASE_URL") + "v1/" + language + "/menu?device=" + DeviceName + "&cascade=" + Casecade + "&limit=" + strconv.FormatInt(limit, 10) + "&page=" + strconv.FormatInt(current_page+1, 10)
// 	} else {
// 		menu.NextPageUrl = ""
// 	}
// 	if current_page-1 > 0 {
// 		menu.PrevPageUrl = os.Getenv("BASE_URL") + "v1/" + language + "/menu?device=" + DeviceName + "&cascade=" + Casecade + "&limit=" + strconv.FormatInt(limit, 10) + "&page=" + strconv.FormatInt(current_page-1, 10)
// 	} else {
// 		menu.PrevPageUrl = ""
// 	}
// 	menu.From = offset
// 	menu.To = offset + limit
// 	menu.Data = pageDetails
// 	c.JSON(http.StatusOK, menu)
// 	return
// }

// GetTopMenuList -  fetches topmenu list
// GET /v1/:lang/contenttype
// @Summary Show a list of topmenu's
// @Description get list of all topmenu list
// @Tags Menu
// @Accept  json
// @Produce  json
// @Param lang query string true "Language"
// @Param device query string true "Device"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Param page query string false "Page"
// @Success 200 {array} MenuDetails
// @Router /v1/{lang}/contenttype [get]
func (hs *HandlerService) GetTopMenuDetails(c *gin.Context) {
	db := c.MustGet("FCDB").(*gorm.DB)
	
	//countryRestriction
	var CountryCode string
	if c.Request.URL.Query()["country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
	}
	fmt.Println("-=-=-=-=-=-=" , os.Getenv("S3_URLFORCONFIG"))
	S3url, err := common.GetCurlCall(os.Getenv("S3_URLFORCONFIG"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	var result interface{}
	err = json.Unmarshal(S3url, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	var mapResult map[string]interface{}
	err1 := json.Unmarshal(S3url, &mapResult)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	// Check if "restrictedCountries" key exists and is of type []interface{}
	restrictedCountriesInterface, ok := mapResult["restrictedCountries"].([]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, "Invalid restrictedCountries data")
		return
	}
	// Convert the interface slice to a string slice
	var countryList []string
	for _, v := range restrictedCountriesInterface {
		if str, ok := v.(string); ok {
			countryList = append(countryList, str)
		}
	}
    isCountryRestricted := false
	fmt.Println("countryList" , countryList)
    for _, country := range countryList {
        if country == CountryCode {
            isCountryRestricted = true
            break
        }
    }
	if isCountryRestricted{
		// c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden,"message":"country restriction" })
		c.JSON(http.StatusForbidden, gin.H{"Message": "Invalid restrictedCountries data"})
		return
	}
	//countryRestriction
	// notFound := common.NotFoundErrorResponse()
	if c.Request.URL.Query()["device"] == nil || c.Request.URL.Query()["device"][0] == "" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	DeviceName := strings.ToLower(c.Request.URL.Query()["device"][0])
	langCode := c.Param("lang")
	serverError := common.ServerErrorResponse(langCode)
	var menu []MenuDetails
	var finalmenu []MenuDetails
	fields := "device,menu_type as menuType,slider_key as sliderKey,url,menu.order"
	if langCode == "en" {
		fields += ",menu_english_name as title"
	} else {
		fields += ",menu_arabic_name as title"
	}
	var limit, offset, current_page int64
	// var CountryCode string
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["page"] != nil {
		current_page, _ = strconv.ParseInt(c.Request.URL.Query()["page"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	if c.Request.URL.Query()["country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}
	fmt.Println("limit =", limit)
	if limit == 0 {
		limit, _ = strconv.ParseInt(os.Getenv("DEFAULT_PAGE_SIZE"), 10, 64)
	}
	sort := "menu.order"
	orderby := ""
	where := "device='" + DeviceName + "' and is_published=true"
	queryParams := "device=" + DeviceName
	if c.Request.URL.Query()["limit"] != nil {
		paginationvar := pagination.GeneratePaginationRequest(c, limit, offset, current_page, sort, orderby)
		menuDetails := pagination.PaginationServices(c, paginationvar, fields, where, queryParams, "menu", &menu)
		c.JSON(http.StatusOK, menuDetails.Details)
		return
	} else {
		var seriesresult, moviesresult, programsresult, livetvresult, playsresult []AllAvailableSeasons
		seriesUrl := os.Getenv("CONTENT_TYPE_URL") + "series&Country=" + CountryCode + os.Getenv("CONTENT_TYPE_URL_PAGINATION")
		moviesUrl := os.Getenv("CONTENT_TYPE_URL") + "movie&Country=" + CountryCode + os.Getenv("CONTENT_TYPE_URL_PAGINATION")
		ProgramsUrl := os.Getenv("CONTENT_TYPE_URL") + "program&Country=" + CountryCode + os.Getenv("CONTENT_TYPE_URL_PAGINATION")
		livetvUrl := os.Getenv("CONTENT_TYPE_URL") + "livetv&Country=" + CountryCode + os.Getenv("CONTENT_TYPE_URL_PAGINATION")
		playsUrl := os.Getenv("CONTENT_TYPE_URL") + "play&Country=" + CountryCode + os.Getenv("CONTENT_TYPE_URL_PAGINATION")
		seriesres, _ := common.GetCurlCall(seriesUrl)
		if err := json.Unmarshal(seriesres, &seriesresult); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		moviesres, _ := common.GetCurlCall(moviesUrl)
		if err := json.Unmarshal(moviesres, &moviesresult); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		programsres, _ := common.GetCurlCall(ProgramsUrl)
		if err := json.Unmarshal(programsres, &programsresult); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		livetvres, _ := common.GetCurlCall(livetvUrl)
		if err := json.Unmarshal(livetvres, &livetvresult); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		playsres, _ := common.GetCurlCall(playsUrl)
		if err := json.Unmarshal(playsres, &playsresult); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		if data := db.Table("menu").Select(fields).Where(where).Order("menu.order").Find(&menu).Error; data != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		for _, val := range menu {
			if val.Url == "series" {
				if len(seriesresult) > 0 {
					finalmenu = append(finalmenu, val)
				}
			} else if val.Url == "movie" {
				if len(moviesresult) > 0 {
					finalmenu = append(finalmenu, val)
				}
			} else if val.Url == "program" {
				if len(programsresult) > 0 {
					finalmenu = append(finalmenu, val)
				}
			} else if val.Url == "livetv" {
				if len(livetvresult) > 0 {
					finalmenu = append(finalmenu, val)
				}
			} else if val.Url == "play" {
				if len(playsresult) > 0 {
					finalmenu = append(finalmenu, val)
				}
			}
		}
		// here made nil to overrode current query fetched result
		menu = nil
		menu = finalmenu
	}

	c.JSON(http.StatusOK, menu)
	return
}
func SliderPlaylists(BlackAreaPlaylistId string, RedAreaPlaylistId string, GreenAreaPlaylistId string, c *gin.Context) ([]Playlist, error) {
	db := c.MustGet("FCDB").(*gorm.DB)
	var playlists []Playlist
	playlist := []interface{}{}
	fields := "id,playlist_key,case"
	if BlackAreaPlaylistId != "" && RedAreaPlaylistId != "" && GreenAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 'black_playlist'  when id='" + RedAreaPlaylistId + "' then 'red_playlist' when id='" + GreenAreaPlaylistId + "' then 'green_playlist'"
		playlist = []interface{}{BlackAreaPlaylistId, RedAreaPlaylistId, GreenAreaPlaylistId}
	} else if BlackAreaPlaylistId != "" && RedAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 'black_playlist'  when id='" + RedAreaPlaylistId + "' then 'red_playlist'"
		playlist = []interface{}{BlackAreaPlaylistId, RedAreaPlaylistId}
	} else if BlackAreaPlaylistId != "" && GreenAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 'black_playlist' when id='" + GreenAreaPlaylistId + "' then 'green_playlist'"
		playlist = []interface{}{BlackAreaPlaylistId, GreenAreaPlaylistId}
	} else if BlackAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 'black_playlist'"
		playlist = []interface{}{BlackAreaPlaylistId}
	}
	fields += " end as playlist_type,"
	fields += " case"
	if BlackAreaPlaylistId != "" && RedAreaPlaylistId != "" && GreenAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 1  when id='" + RedAreaPlaylistId + "' then 2 when id='" + GreenAreaPlaylistId + "' then 3"
	} else if BlackAreaPlaylistId != "" && RedAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 1  when id='" + RedAreaPlaylistId + "' then 2"
	} else if BlackAreaPlaylistId != "" && GreenAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 1 when id='" + GreenAreaPlaylistId + "' then 3"
	} else if BlackAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 1"
	}
	fields += " else 0 end as playlist_order"
	if err := db.Select(fields).Where("id in(?) and (scheduling_start_date <=now() or scheduling_start_date is null) and (scheduling_end_date >=now() or scheduling_end_date is null)", playlist).Order("playlist_type desc").Find(&playlists).Error; err != nil {
		return nil, err
	}
	return playlists, nil
}
func PlaylistItemContents(playlistId string, c *gin.Context) ([]PlaylistContentIds, error) {
	db := c.MustGet("FCDB").(*gorm.DB)
	var contentIds []PlaylistContentIds
	if err := db.Table("playlist_item_content pic").Select("pic.content_id").Joins("inner join playlist_item pi2 on pi2.id=pic.playlist_item_id inner join playlist p on p.id=pi2.playlist_id").Where("p.id =?", playlistId).Order("pi2.order asc").Find(&contentIds).Error; err != nil {
		return nil, err
	}
	return contentIds, nil
}
func PlaylistPages(playlistId, language string, pageContents chan []PageContent, c *gin.Context) {
	db := c.MustGet("FCDB").(*gorm.DB)
	var pageContent []PageContent
	var finalResult []PageContent
	fields := "p.id as key,p.page_key as id,p.page_type::text as type"
	if language == "en" {
		fields += ",p.english_title as title,p.english_page_friendly_url as friendly_url,p.english_meta_description as seo_description"
	} else {
		fields += ",p.arabic_title as title,p.arabic_page_friendly_url as friendly_url,p.arabic_meta_description as seo_description"
	}
	if err := db.Table("playlist_item pi2").Select(fields).Joins("join page p on p.id=pi2.group_by_page_id").Where("playlist_id=? and p.is_disabled =false and p.deleted_by_user_id is null", playlistId).Order("pi2.order").Find(&pageContent).Error; err != nil {
		return
	}
	for _, page := range pageContent {
		var imageryDetails ImageryDetails
		imageryDetails.MobileMenu = os.Getenv("IMAGES") + strings.ToLower(page.Key) + "/mobile-menu" + "?c=" + strconv.Itoa(rand.Intn(200000))
		imageryDetails.MobilePosterImage = os.Getenv("IMAGES") + strings.ToLower(page.Key) + "/menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
		imageryDetails.MobileMenuPosterImage = os.Getenv("IMAGES") + strings.ToLower(page.Key) + "/mobile-menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
		page.Imagery = imageryDetails
		finalResult = append(finalResult, page)
	}
	pageContents <- finalResult
	return
}
