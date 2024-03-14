package menu

import (
	"context"
	"encoding/json"
	"fmt"
	"frontend_config/common"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	// "github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const numKeys = 1000000
const routines = 16
const portion = numKeys / routines

var uuids [numKeys]string
var wg sync.WaitGroup

// var pool *redis.Pool

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	qrg := r.Group("/api")
	qrg.Use(common.ValidateToken())
	qrg.GET("/ratings/filters", hs.Menu)
	qrg.GET("/publishingplatforms", hs.GetPublishingPlatforms)
	r.GET("/:lang/menu/:pageKey", hs.HomePage)
	r.GET("/:lang/getmenu", hs.TvHomePage)
	r.GET("/api/pagesync", hs.UpdateAllPageCache)
	r.GET("/api/playlistsync", hs.UdateDirtyCountPlaylistRelatedPages)
	r.GET("/api/slidersync", hs.UpdateDirtycountSliderRelatedPages)

	// r.GET("/redisCache/:pageKey", hs.CreateRedisCacheForHomePage)

	//Error code exception URL
	/*filters*/
	qrg.PUT("/ratings/filters", hs.Menu)
	qrg.POST("/ratings/filters", hs.Menu)
	qrg.DELETE("/ratings/filters", hs.Menu)
	/*PublishingPlatforms*/
	qrg.PUT("/publishingplatforms", hs.GetPublishingPlatforms)
	qrg.POST("/publishingplatforms", hs.GetPublishingPlatforms)
	qrg.DELETE("/publishingplatforms", hs.GetPublishingPlatforms)
}

// GetRatingfilters -  fetches filters list
// GET /api/ratings/filters
// @Summary Show rating filters
// @Description get rating filters
// @Tags Frontend
// @Security Authorization
// @Accept  json
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /api/ratings/filters [get]
func (hs *HandlerService) Menu(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	var newmenu []NewMenu
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var errorresponse = common.ServerErrorResponse()
	if err := db.Debug().Table("menu").Select("distinct(url)").Find(&newmenu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	newarr := []string{}
	contenttype := make(map[string][]string)
	for _, element := range newmenu {
		newarr = append(newarr, strings.Title(element.Url))
	}
	contenttype["contentTypes"] = newarr
	c.JSON(http.StatusOK, contenttype)
}

// GetPublishingPlatforms -  Get all Publishing Platforms
// GET /api/publishingplatforms
// @Summary Show Publishing Platforms
// @Description get Publishing Platforms
// @Tags Frontend
// @Accept  json
// @Security Authorization
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /api/publishingplatforms [get]
func (hs *HandlerService) GetPublishingPlatforms(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var publishingplatforms []Publishingplatforms
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var errorresponse = common.ServerErrorResponse()
	if err := db.Debug().Table("publish_platform").Select("id, platform as name").Order("id asc").Find(&publishingplatforms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": publishingplatforms})
}

type RedisCacheRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

/* Page Fragment */
// func (hs *HandlerService) HomePage(c *gin.Context) {
// 	fdb := c.MustGet("FDB").(*gorm.DB)
// 	db := c.MustGet("DB").(*gorm.DB)
// 	cdb := c.MustGet("CDB").(*gorm.DB)
// 	pageKey := c.Param("pageKey")
// 	language := c.Param("lang")
// 	// country := c.Param("country")
// 	// platform := c.Param("platform")

// 	var country, platform int
// 	platform = 0
// 	// flag for validating slider data
// 	var IsCheck bool
// 	type CountryDetails struct {
// 		Id int
// 	}
// 	type PlatformDetails struct {
// 		Id int
// 	}
// 	type PageType struct {
// 		PageType int
// 	}
// 	var countryDetails CountryDetails
// 	var platformDetails PlatformDetails
// 	fmt.Println(c.Request.URL.Query()["country"], "country")
// 	if c.Request.URL.Query()["country"] != nil {
// 		db.Debug().Raw("select id from country where alpha2code=?", c.Request.URL.Query()["country"][0]).Find(&countryDetails)
// 		country = countryDetails.Id
// 		if country == 0 {
// 			country = 784
// 		}
// 		// country = strings.ToUpper(c.Request.URL.Query()["country"][0])
// 	} else {
// 		country = 784
// 	}
// 	fmt.Println(c.Request.URL.Query()["platform"], "platform")
// 	if c.Request.URL.Query()["platform"] != nil {
// 		db.Debug().Raw("select id from publish_platform where platform=?", c.Request.URL.Query()["platform"][0]).Find(&countryDetails)
// 		platform = platformDetails.Id
// 	} else {
// 		platform = 0
// 	}
// 	var Result MenuPage
// 	var menuPage MenuPageDetails
// 	var imageryDetails ImageryDetails
// 	var playlistIds, contentIds []string
// 	var featured FeaturedResponse
// 	var featuredPlaylistsResponse []FeaturedPlaylistsResponse
// 	var featuredPlaylist FeaturedPlaylists
// 	var featuredPlaylists []FeaturedPlaylists
// 	var menuPlaylsts []MenuPlaylists
// 	var menuPlaylst MenuPlaylists

// 	key := c.Request.Host + pageKey + language + strconv.Itoa(country) + strconv.Itoa(platform)
// 	key = strings.ReplaceAll(key, "/", "_")
// 	key = strings.ReplaceAll(key, "?", "_")
// 	url := os.Getenv("REDIS_CACHE_URL") + "/" + key
// 	fmt.Println("Redis url -", url, "???????????")
// 	response, err := common.GetCurlCall(url)
// 	if err != nil {
// 		fmt.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
// 		return
// 	}
// 	type RedisCacheResponse struct {
// 		Value string `json:"value"`
// 		Error string `json:"error"`
// 	}
// 	var RedisResponse RedisCacheResponse
// 	json.Unmarshal(response, &RedisResponse)
// 	if RedisResponse.Value != "" {
// 		if err := json.Unmarshal([]byte(RedisResponse.Value), &Result); err != nil {
// 			fmt.Println(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
// 			return
// 		}
// 	} else {
// 		if menupageError := db.Debug().Raw("select p.id, p.page_key,p.has_menu_poster_image,case when 'en' = lower(?) then english_page_friendly_url when 'ar' = lower(?) then arabic_page_friendly_url end as friendly_url ,	case when 'en' = lower(?) then english_meta_description when 'ar' = lower(?) then arabic_meta_description end as seo_description,case when 'en' = lower(?) then p.english_title when 'ar' = lower(?) then p.arabic_title end as title , case when page_type = 1 then 'Home' when page_type = 8 then 'Settings' when page_type = 16 then 'Favourites' when jsonb_agg(distinct ps.slider_id) = '[null]' then 'VOD' else 'Home' end as page_type from page p left join page_slider ps on ps.page_id = p.id left join slider s on s.id = ps.slider_id left join page_country pc on pc.page_id = p.id left join country c on c.id = pc.country_id where p.page_key =? and p.is_disabled = false and p.deleted_by_user_id is null group by p.id,s.slider_key ,p.page_key", language, language, language, language, language, language, pageKey).Find(&menuPage).Error; menupageError != nil {
// 			fmt.Println(menupageError)
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": menupageError})
// 		}
// 		Result.ID = menuPage.PageKey
// 		Result.Type = menuPage.PageType
// 		fmt.Println("page type is", menuPage.PageType)
// 		Result.SeoDescription = menuPage.SeoDescription
// 		Result.Title = menuPage.Title
// 		Result.FriendlyUrl = menuPage.FriendlyUrl
// 		/* Featured Details */
// 		// here slider scheduling has been handled as per .net
// 		var FeaturedResult Featured
// 		var pagetype PageType
// 		db.Debug().Raw("select page_type from page where page_key = ?", pageKey).Find(&pagetype)
// 		fmt.Println("page type", pagetype.PageType)
// 		if pagetype.PageType == 1 {
// 			// Here checking for general slider if not true go get deafult slider
// 			res := db.Debug().Raw("select jsonb_build_object('id', slider_key, 'type', st.name, 'playlists', jsonb_agg(distinct concat(red_area_playlist_id, ',' , green_area_playlist_id, ',' ,black_area_playlist_id))) as featured from page p join page_slider ps on ps.page_id = p.id join slider s on s.id = ps.slider_id join slider_types st on st.id = s.type where p.page_key = ? and p.is_disabled = false and s.deleted_by_user_id is null and s.is_disabled = false and (s.scheduling_start_date  <= now() or s.scheduling_start_date is null ) and (s.scheduling_end_date > now() or s.scheduling_end_date is null) and p.deleted_by_user_id is null and ps.order = 1 group by s.slider_key,st.name", pageKey).Find(&FeaturedResult).RowsAffected
// 			if res == 0 {
// 				db.Debug().Raw("select jsonb_build_object('id', slider_key, 'type', st.name, 'playlists', jsonb_agg(distinct concat(red_area_playlist_id, ',' , green_area_playlist_id, ',' ,black_area_playlist_id))) as featured from page p join page_slider ps on ps.page_id = p.id join slider s on s.id = ps.slider_id join slider_types st on st.id = s.type where p.page_key = ? and p.is_disabled = false and s.deleted_by_user_id is null and s.is_disabled = false and p.deleted_by_user_id is null and ps.order = 0 group  by  s.slider_key,st.name", pageKey).Find(&FeaturedResult)
// 			}
// 		} else {
// 			db.Debug().Raw("select jsonb_build_object('id', slider_key, 'type', st.name, 'playlists', jsonb_agg(distinct concat(red_area_playlist_id, ',' , green_area_playlist_id, ',' ,black_area_playlist_id))) as featured from page p join page_slider ps on ps.page_id = p.id join slider s on s.id = ps.slider_id join slider_types st on st.id = s.type where p.page_key = ? and p.is_disabled = false and s.deleted_by_user_id is null and s.is_disabled = false and (s.scheduling_start_date  <= now() or s.scheduling_start_date is null ) and (s.scheduling_end_date > now() or s.scheduling_end_date is null) and p.deleted_by_user_id is null group  by s.slider_key,st.name", pageKey).Find(&FeaturedResult)
// 		}
// 		// (p.scheduling_start_date  <= now() or p.scheduling_start_date is null ) and (p.scheduling_end_date > now() or p.scheduling_end_date is null)
// 		data, _ := json.Marshal(FeaturedResult.Featured)
// 		json.Unmarshal(data, &featured)
// 		fmt.Println(string(data))
// 		fmt.Println("FeaturedResult", FeaturedResult.Featured)
// 		if featured.ID == 0 && featured.Type == "" {
// 			Result.Featured = nil
// 		} else {
// 			var stringsnew []string
// 			for _, ids := range featured.Playlists {
// 				stringsnew = strings.Split(ids, ",")
// 			}
// 			// var contentType []string
// 			var redcheck, blackcheck, greencheck bool
// 			for keyv, val := range stringsnew {
// 				var ctype string
// 				if keyv == 0 {
// 					ctype = "red_playlist"
// 				}
// 				if keyv == 1 {
// 					ctype = "green_playlist"
// 				}
// 				if keyv == 2 {
// 					ctype = "black_playlist"
// 				}
// 				//contentType = append(contentType, ctype)

// 				/* playlist response */
// 				var playlistResponse FeaturedPlaylistsResponse
// 				db.Debug().Raw("select p.playlist_key as id,'"+ctype+"' as playlist_type,case when 'en' = lower(?) then english_title when 'ar' = lower(?) then arabic_title end as title,jsonb_agg(coalesce(case when pi2.multi_tier_content_id is not null then pi2.multi_tier_content_id end, case when pi2.one_tier_content_id is not null then pi2.one_tier_content_id end, pic.content_id) order by pi2.order)as content from playlist p left join playlist_item pi2 on pi2.playlist_id = p.id left join playlist_item_content pic on pic.playlist_item_id = pi2.id left join page_playlist pp on pp.playlist_id = p.id and pp.page_id = ? where p.id =? and p.is_disabled='false' and (p.scheduling_start_date  <= now() or p.scheduling_start_date is null ) and (p.scheduling_end_date > now() or p.scheduling_end_date is null) group by p.id,pp.order order by pp.order", language, language, menuPage.ID, val).Find(&playlistResponse)
// 				// var Ids string
// 				// cdb.Debug().Raw("select c.id from content c where c.id in (?)and c.status = 1 and c.deleted_by_user_id is null", playlistResponse.Content).Find(&Ids)
// 				if playlistResponse.ID != 0 {
// 					featuredPlaylistsResponse = append(featuredPlaylistsResponse, playlistResponse)
// 				}
// 				if ctype == "red_playlist" {
// 					if playlistResponse.ID == 0 {
// 						redcheck = true
// 					}
// 				} else if ctype == "green_playlist" {
// 					if playlistResponse.ID == 0 {
// 						greencheck = true
// 					}
// 				} else if ctype == "black_playlist" {
// 					if playlistResponse.ID == 0 {
// 						blackcheck = true
// 					}
// 				}
// 			}
// 			if redcheck == true && greencheck == true && blackcheck == true {
// 				/* here for making slider data null if no playlist has rights or active state */
// 				IsCheck = true
// 			}
// 			for _, playlists := range featuredPlaylistsResponse {
// 				featuredPlaylist.ID = playlists.ID
// 				featuredPlaylist.PlaylistType = playlists.PlaylistType
// 				featuredPlaylist.Title = playlists.Title
// 				contentId, _ := json.Marshal(playlists.Content)
// 				json.Unmarshal(contentId, &contentIds)
// 				var content1 Cont
// 				res := strconv.Itoa(country)
// 				if featured.Type == "Layout A – Smart TV" {
// 					if playlists.PlaylistType == "red_playlist" {
// 						var countcheck int
// 						for _, val := range contentIds {
// 							var forcount int
// 							fdb.Debug().Raw("select count(*) from content_fragment cf where cf.content_id = ? and country like '%"+res+"%' ", val).Count(&forcount)
// 							countcheck = countcheck + forcount
// 						}
// 						if countcheck < 14 {
// 							IsCheck = true
// 						}
// 					} else if playlists.PlaylistType == "green_playlist" {
// 						var countcheck int
// 						for _, val := range contentIds {
// 							var forcount int
// 							fdb.Debug().Raw("select count(*) from content_fragment cf where cf.content_id = ? and country like '%"+res+"%' ", val).Count(&forcount)
// 							countcheck = countcheck + forcount
// 							if countcheck < 1 {
// 								IsCheck = true
// 							}
// 						}
// 					} else if playlists.PlaylistType == "black_playlist" {
// 						var countcheck int
// 						for _, val := range contentIds {
// 							var forcount int
// 							fdb.Debug().Raw("select count(*) from content_fragment cf where cf.content_id = ? and country like '%"+res+"%' ", val).Count(&forcount)
// 							countcheck = countcheck + forcount
// 							if countcheck < 1 {
// 								IsCheck = true
// 							}
// 						}
// 					}
// 				} else if featured.Type == "Layout B - STV / Website / Apple TV" {
// 					// scenario may need to be changed based on requirment
// 					if playlists.PlaylistType == "red_playlist" {
// 						var countcheck int
// 						for _, val := range contentIds {
// 							var forcount int
// 							fdb.Debug().Raw("select count(*) from content_fragment cf where cf.content_id = ? and country like '%"+res+"%' ", val).Count(&forcount)
// 							countcheck = countcheck + forcount
// 						}
// 						// need to be changed as per how much needed here it has been set to 14 as redplaylist must have 7 contents and double we get from fragments
// 						if countcheck < 14 {
// 							IsCheck = true
// 						}
// 					} else if playlists.PlaylistType == "black_playlist" {
// 						var countcheck int
// 						for _, val := range contentIds {
// 							var forcount int
// 							fdb.Debug().Raw("select count(*) from content_fragment cf where cf.content_id = ? and country like '%"+res+"%' ", val).Count(&forcount)
// 							countcheck = countcheck + forcount
// 							if countcheck < 1 {
// 								IsCheck = true
// 							}
// 						}
// 					}
// 				} else if featured.Type == "layout C - STV - Website - Apple TV" {
// 					//if playlists.PlaylistType == "black_playlist" {
// 					var countcheck int
// 					for _, val := range contentIds {
// 						var forcount int
// 						fdb.Debug().Raw("select count(*) from content_fragment cf where cf.content_id = ? and country like '%"+res+"%' ", val).Count(&forcount)
// 						countcheck = countcheck + forcount
// 					}
// 					if countcheck < 1 {
// 						IsCheck = true
// 					}
// 					//}
// 				}
// 				// var finalcontents []string
// 				// for _, value := range contentIds {
// 				// 	var rescount int
// 				// 	fdb.Debug().Raw("select count(*) from content_fragment cf where cf.content_id = ? and country like '%"+res+"%' ", value).Count(&rescount)
// 				// 	if rescount > 0 {
// 				// 		finalcontents = append(finalcontents, value)
// 				// 	}
// 				// }
// 				var finalcontents []string
// 				var fetchvalidids ContentIds
// 				cdb.Debug().Raw("select jsonb_agg(id order by array_position(array[?], id::text)) as ids from (select distinct c.id from content c left join content_variance cv on cv.content_id =c.id left join season s on s.content_id = c.id left join variance_trailer vt on vt.content_variance_id = s.id or vt.content_variance_id = cv.id left join episode e on e.season_id = s.id join playback_item pi2 on pi2.id = cv.playback_item_id or pi2.id = e.playback_item_id join content_rights_country crc on crc.content_rights_id = s.rights_id or crc.content_rights_id = pi2.rights_id where c.status = 1 and (cv.status =1 or cv.status is null) and (s.status = 1 or s.status is null) and ((e.status =1 or e.status is null) or (vt.id is not null)) and crc.country_id = ? and c.id in (?)) as foo", contentIds, country, contentIds).Find(&fetchvalidids)
// 				validcontentidswithorder, _ := json.Marshal(fetchvalidids.Ids)
// 				json.Unmarshal(validcontentidswithorder, &finalcontents)
// 				fmt.Println(len(finalcontents), "length of final contents")
// 				if contentError := fdb.Debug().Raw("select jsonb_agg(details order by array_position(array[?], cf.content_id::text)) as cnt from content_fragment cf where cf.content_id in (?) and language in (?)  and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(country,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(platform,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and rights_start_date <= now() and rights_end_date >= now() and language = ? and country like '%"+res+"%' ", finalcontents, finalcontents, language, finalcontents, country, country, finalcontents, platform, platform, language).Find(&content1).Error; contentError != nil {
// 					fmt.Println(contentError)
// 					c.JSON(http.StatusInternalServerError, gin.H{"message": contentError})

// 				}
// 				featuredPlaylist.Content = content1.Cnt
// 				featuredPlaylists = append(featuredPlaylists, featuredPlaylist)
// 			}

// 			Result.Featured = &FeaturedDetails{featured.ID, featured.Type, featuredPlaylists}
// 		}
// 		if menuPage.HasMenuPosterImage == true {
// 			imageryDetails.MobileMenu = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/mobile-menu" + "?c=" + strconv.Itoa(rand.Intn(200000))
// 			imageryDetails.MobilePosterImage = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
// 			imageryDetails.MobileMenuPosterImage = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/mobile-menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
// 		} else {
// 			imageryDetails.MobileMenu = ""
// 			imageryDetails.MobilePosterImage = ""
// 			imageryDetails.MobileMenuPosterImage = ""
// 		}
// 		if imageryDetails.MobileMenu == "" || imageryDetails.MobilePosterImage == "" || imageryDetails.MobileMenuPosterImage == "" {
// 			Result.Imagery = nil
// 		} else {
// 			Result.Imagery = &imageryDetails
// 		}
// 		if IsCheck {
// 			Result.Featured = nil
// 		}
// 		if Result.Featured == nil && menuPage.Title != "Home" && menuPage.Title != "الصفحة الرئيسية" {
// 			Result.Type = "VOD"
// 		}
// 		// Page Playlists
// 		var plays Plays
// 		db.Debug().Raw("select jsonb_agg(p4.id order by pp2.order)::varchar as playlists from page p3 join page_playlist pp2 on pp2.page_id = p3.id join playlist p4 on p4.id = pp2.playlist_id where p4.deleted_by_user_id is null and p4.is_disabled = false and p3.page_key = ?", pageKey).Find(&plays)
// 		var parts []string
// 		if err := json.Unmarshal([]byte(plays.Playlists), &parts); err != nil {
// 			fmt.Println(err)
// 		}
// 		for _, ids := range parts {
// 			playlistIds = append(playlistIds, string(ids))
// 		}

// 		if playListError := db.Debug().Raw("select p.playlist_key as id,case when 'en'=lower('"+language+"') then p.english_title when 'ar'=lower('"+language+"') then p.arabic_title end as playlist_type,jsonb_agg(coalesce(case when pi2.multi_tier_content_id is not null then pi2.multi_tier_content_id end, case when pi2.one_tier_content_id is not null then pi2.one_tier_content_id end,case when pi2.season_id is not null then pi2.season_id end,pic.content_id) order by pi2.order)as content from playlist p join playlist_item pi2 on pi2.playlist_id = p.id join playlist_item_content pic on pic.playlist_item_id = pi2.id join page_playlist pp on pp.playlist_id = p.id and pp.page_id =? where p.id in (?) and p.is_disabled='false' and (p.scheduling_start_date  <= now() or p.scheduling_start_date is null ) and (p.scheduling_end_date > now() or p.scheduling_end_date is null) group by p.id,pp.order order by pp.order", menuPage.ID, playlistIds).Find(&featuredPlaylistsResponse).Error; playListError != nil {
// 			fmt.Println(playListError)
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": playListError})
// 			return
// 		}
// 		menuply := make([]MenuPlaylists, 0)
// 		if len(featuredPlaylistsResponse) <= 0 {
// 			menuPlaylsts = menuply
// 			fmt.Println(len(featuredPlaylistsResponse))
// 		}
// 		for _, playlists := range featuredPlaylistsResponse {
// 			menuPlaylst.ID = playlists.ID
// 			menuPlaylst.Title = playlists.PlaylistType
// 			contentId, _ := json.Marshal(playlists.Content)
// 			json.Unmarshal(contentId, &contentIds)
// 			var totalCount int
// 			cdb.Debug().Raw("select count(id) as not_valid_content from content where id in (?) and status!=2", contentIds).Count(&totalCount)
// 			// if totalCount > 0 {
// 			// 	var content2 Cont
// 			// 	res1 := strconv.Itoa(country)
// 			// 	if contentError := fdb.Debug().Raw("select jsonb_agg(details) as cnt from content_fragment cf where cf.content_id in (?) and language in (?)  and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(country,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(platform,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and rights_start_date <= now() and rights_end_date >= now() and language = ? and country like '%"+res1+"%' ", contentIds, language, contentIds, country, country, contentIds, platform, platform, language).Find(&content2).Error; contentError != nil {
// 			// 		fmt.Println(contentError)
// 			// 		c.JSON(http.StatusInternalServerError, gin.H{"message": contentError})
// 			// 		return
// 			// 	}
// 			// for _, val := range contentIds {
// 			// 	//check weather season amd MTC both exists
// 			// 	cdb.Debug().Raw()
// 			// }
// 			//here query changed to contents in playlist as per order
// 			var finalcontentids []string
// 			var contentidsfororder []string
// 			fmt.Println("printing  contentids", contentIds)
// 			for _, val := range contentIds {
// 				var mltcount int
// 				var moviecount int
// 				var noncontentidcount int
// 				var varianceidforonetire seasondetails
// 				var conid seasondetail
// 				var id seasondetails
// 				res2 := strconv.Itoa(country)
// 				//here fetching content ids for seasons added to playlist
// 				cdb.Debug().Raw("select count(id) from content where id = ?", val).Count(&noncontentidcount)
// 				if noncontentidcount <= 0 {
// 					// here fetching content id from season id
// 					var seasoncheck int
// 					cdb.Debug().Raw("select count(id) from season s where id = ? ", val).Count(&seasoncheck)
// 					if seasoncheck > 0 {
// 						cdb.Debug().Raw("select s.content_id  from season s where id  = ? ", val).Find(&conid)
// 						fmt.Println("printing contetn id here", conid.ContentID)
// 						if conid.ContentID != "" {
// 							contentidsfororder = append(contentidsfororder, conid.ContentID)
// 						}
// 					} else {
// 						contentidsfororder = append(contentidsfororder, val)
// 					}
// 				} else {
// 					contentidsfororder = append(contentidsfororder, val)
// 				}
// 				cdb.Debug().Raw("select count(id) from content where (content_type = 'Series' or content_type = 'Program')and status =1 and id = ? ", val).Count(&mltcount)
// 				cdb.Debug().Raw("select count(id) from content where (content_type = 'Movie' or content_type = 'LiveTV' or content_type  = 'Play') and status=1 and id = ?", val).Count(&moviecount)
// 				if mltcount > 0 {
// 					//fetching content season 1 id as per .net validation
// 					cdb.Debug().Raw("select s.id from content c join season s on s.content_id = c.id where c.id = ? and s.status = 1 and s.deleted_by_user_id is null and c.status =1 order by s.number asc fetch first row only", val).Find(&id)
// 					if id.ID != "" {
// 						finalcontentids = append(finalcontentids, id.ID)
// 					}
// 				} else if moviecount > 0 {
// 					// fetching variance id for one tire content
// 					cdb.Debug().Raw("select cv.id from  content c left join content_variance cv on cv.content_id =c.id left join playback_item pi2 on pi2.id = cv.playback_item_id left join content_rights_country crc on crc.content_rights_id = pi2.rights_id where c.id = ? and crc.country_id = ? and cv.status = 1 and c.status = 1", val, res2).Find(&varianceidforonetire)
// 					if varianceidforonetire.ID != "" {
// 						finalcontentids = append(finalcontentids, varianceidforonetire.ID)
// 					}
// 					// } else if mltcount <= 0 && moviecount <= 0 {
// 					// 	var validcontentcount int
// 					// 	cdb.Debug().Raw("select count(id) from content c where c.status = 1 and c.id = ?", val).Count(&validcontentcount)
// 					// 	if validcontentcount > 0 {
// 					// 		finalcontentids = append(finalcontentids, val)
// 					// 	}
// 				}
// 			}
// 			if totalCount > 0 {
// 				var content2 Cont
// 				res1 := strconv.Itoa(country)
// 				if contentError := fdb.Debug().Raw("select jsonb_agg(details order by array_position(array[?], cf.content_id::text)) as cnt from content_fragment cf where (cf.content_id in (?) or  cf.content_variance_id in (?))and language in (?)  and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(country,'[]','{}')) from content_fragment cf where (cf.content_id in (?) or  cf.content_variance_id in (?))) as ss )aa where cid= ? )= ? and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(platform,'[]','{}')) from content_fragment cf where (cf.content_id in (?) or cf.content_variance_id in (?))) as ss )aa where cid= ? )= ? and rights_start_date <= now() and rights_end_date >= now() and language = ? and country like '%"+res1+"%' ", contentidsfororder, finalcontentids, finalcontentids, language, finalcontentids, finalcontentids, country, country, finalcontentids, finalcontentids, platform, platform, language).Find(&content2).Error; contentError != nil {
// 					fmt.Println(contentError)
// 					c.JSON(http.StatusInternalServerError, gin.H{"message": contentError})
// 					return
// 				}
// 				menuPlaylst.Content = content2.Cnt
// 				menuPlaylsts = append(menuPlaylsts, menuPlaylst)
// 			} else {
// 				continue
// 			}
// 		}

// 		Result.Playlists = menuPlaylsts
// 		jsonData, _ := json.Marshal(Result)
// 		var request RedisCacheRequest
// 		url := os.Getenv("REDIS_CACHE_URL")
// 		request.Key = key //pageKey + language + strconv.Itoa(country) + strconv.Itoa(platform)
// 		request.Value = string(jsonData)
// 		_, err := common.PostCurlCall("POST", url, request)
// 		if err != nil {
// 			fmt.Println(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": Result})
// }

// changed to new if need old api use above api if needed new uncomment below api
func (hs *HandlerService) HomePage(c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	db := c.MustGet("DB").(*gorm.DB)
	cdb := c.MustGet("CDB").(*gorm.DB)
	pageKey := c.Param("pageKey")
	language := c.Param("lang")

	var country, platform int
	platform = 0
	var IsCheck bool
	type PageType struct {
		PageType int
	}
	if c.Request.URL.Query()["country"] != nil {
		countryAlphacode := c.Request.URL.Query()["country"][0]
		country = int(common.Countrys(countryAlphacode))
		if country == 0 {
			country = 784
		}
	} else {
		country = 784
	}
	if c.Request.URL.Query()["platform"] != nil {
		platforminput := c.Request.URL.Query()["platform"][0]
		platform = int(common.PublishingPlatforms(platforminput))
	} else {
		platform = 0
	}
	var Result MenuPage
	var menuPage MenuPageDetails
	var imageryDetails ImageryDetails
	var playlistIds, contentIds []string
	var featured FeaturedResponse
	var featuredPlaylistsResponse []FeaturedPlaylistsResponse
	var featuredPlaylist FeaturedPlaylists
	var featuredPlaylists []FeaturedPlaylists
	var menuPlaylsts []MenuPlaylists
	var menuPlaylst MenuPlaylists
	key := c.Request.Host + pageKey + language + strconv.Itoa(country) + strconv.Itoa(platform)
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "?", "_")
	url := os.Getenv("REDIS_CACHE_URL") + "/" + key
	fmt.Println("Redis url -", url, "???????????")
	response, err := common.GetCurlCall(url)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	type RedisCacheResponse struct {
		Value string `json:"value"`
		Error string `json:"error"`
	}
	var RedisResponse RedisCacheResponse
	json.Unmarshal(response, &RedisResponse)
	if RedisResponse.Value != "" {
		if err := json.Unmarshal([]byte(RedisResponse.Value), &Result); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	} else {
		if menupageError := db.Debug().Raw("select p.id, p.page_key,p.has_menu_poster_image,case when 'en' = lower(?) then english_page_friendly_url when 'ar' = lower(?) then arabic_page_friendly_url end as friendly_url ,	case when 'en' = lower(?) then english_meta_description when 'ar' = lower(?) then arabic_meta_description end as seo_description,case when 'en' = lower(?) then p.english_title when 'ar' = lower(?) then p.arabic_title end as title , case when page_type = 1 then 'Home' when page_type = 8 then 'Settings' when page_type = 16 then 'Favourites' when jsonb_agg(distinct ps.slider_id) = '[null]' then 'VOD' else 'Home' end as page_type from page p left join page_slider ps on ps.page_id = p.id left join slider s on s.id = ps.slider_id left join page_country pc on pc.page_id = p.id left join country c on c.id = pc.country_id where p.page_key =? and p.is_disabled = false and p.deleted_by_user_id is null group by p.id,s.slider_key ,p.page_key", language, language, language, language, language, language, pageKey).Find(&menuPage).Error; menupageError != nil {
			fmt.Println(menupageError)
			c.JSON(http.StatusInternalServerError, gin.H{"message": menupageError})
		}
		Result.ID = menuPage.PageKey
		Result.Type = menuPage.PageType
		fmt.Println("page type is", menuPage.PageType)
		Result.SeoDescription = menuPage.SeoDescription
		Result.Title = menuPage.Title
		Result.FriendlyUrl = menuPage.FriendlyUrl
		/* Featured Details */
		// here slider scheduling has been handled as per .net
		var FeaturedResult Featured
		var pagetype PageType
		db.Debug().Raw("select page_type from page where page_key = ?", pageKey).Find(&pagetype)
		fmt.Println("page type", pagetype.PageType)
		if pagetype.PageType == 1 {
			// Here checking for general slider if not true go get deafult slider
			res := db.Debug().Raw("select jsonb_build_object('id', slider_key, 'type', st.name, 'playlists', jsonb_agg(distinct concat(red_area_playlist_id, ',' , green_area_playlist_id, ',' ,black_area_playlist_id))) as featured from page p join page_slider ps on ps.page_id = p.id join slider s on s.id = ps.slider_id join slider_types st on st.id = s.type where p.page_key = ? and p.is_disabled = false and s.deleted_by_user_id is null and s.is_disabled = false and (s.scheduling_start_date  <= now() or s.scheduling_start_date is null ) and (s.scheduling_end_date > now() or s.scheduling_end_date is null) and p.deleted_by_user_id is null and ps.order = 1 group by s.slider_key,st.name", pageKey).Find(&FeaturedResult).RowsAffected
			if res == 0 {
				db.Debug().Raw("select jsonb_build_object('id', slider_key, 'type', st.name, 'playlists', jsonb_agg(distinct concat(red_area_playlist_id, ',' , green_area_playlist_id, ',' ,black_area_playlist_id))) as featured from page p join page_slider ps on ps.page_id = p.id join slider s on s.id = ps.slider_id join slider_types st on st.id = s.type where p.page_key = ? and p.is_disabled = false and s.deleted_by_user_id is null and s.is_disabled = false and p.deleted_by_user_id is null and ps.order = 0 group  by  s.slider_key,st.name", pageKey).Find(&FeaturedResult)
			}
		} else {
			db.Debug().Raw("select jsonb_build_object('id', slider_key, 'type', st.name, 'playlists', jsonb_agg(distinct concat(red_area_playlist_id, ',' , green_area_playlist_id, ',' ,black_area_playlist_id))) as featured from page p join page_slider ps on ps.page_id = p.id join slider s on s.id = ps.slider_id join slider_types st on st.id = s.type where p.page_key = ? and p.is_disabled = false and s.deleted_by_user_id is null and s.is_disabled = false and (s.scheduling_start_date  <= now() or s.scheduling_start_date is null ) and (s.scheduling_end_date > now() or s.scheduling_end_date is null) and p.deleted_by_user_id is null group  by s.slider_key,st.name", pageKey).Find(&FeaturedResult)
		}
		data, _ := json.Marshal(FeaturedResult.Featured)
		json.Unmarshal(data, &featured)
		fmt.Println(string(data))
		fmt.Println("FeaturedResult", FeaturedResult.Featured)
		if featured.ID == 0 && featured.Type == "" {
			Result.Featured = nil
		} else {
			var stringsnew []string
			for _, ids := range featured.Playlists {
				stringsnew = strings.Split(ids, ",")
			}
			var redcheck, blackcheck, greencheck bool
			for keyv, val := range stringsnew {
				var ctype string
				if keyv == 0 {
					ctype = "red_playlist"
				}
				if keyv == 1 {
					ctype = "green_playlist"
				}
				if keyv == 2 {
					ctype = "black_playlist"
				}
				/* playlist response */
				var playlistResponse FeaturedPlaylistsResponse
				db.Debug().Raw("select p.playlist_key as id,'"+ctype+"' as playlist_type,case when 'en' = lower(?) then english_title when 'ar' = lower(?) then arabic_title end as title,jsonb_agg(coalesce(case when pi2.multi_tier_content_id is not null then pi2.multi_tier_content_id end, case when pi2.one_tier_content_id is not null then pi2.one_tier_content_id end, pic.content_id) order by pi2.order)as content from playlist p left join playlist_item pi2 on pi2.playlist_id = p.id left join playlist_item_content pic on pic.playlist_item_id = pi2.id left join page_playlist pp on pp.playlist_id = p.id and pp.page_id = ? where p.id =? and p.is_disabled='false' and (p.scheduling_start_date  <= now() or p.scheduling_start_date is null ) and (p.scheduling_end_date > now() or p.scheduling_end_date is null) group by p.id,pp.order order by pp.order", language, language, menuPage.ID, val).Find(&playlistResponse)
				if playlistResponse.ID != 0 {
					featuredPlaylistsResponse = append(featuredPlaylistsResponse, playlistResponse)
				}
				if ctype == "red_playlist" {
					if playlistResponse.ID == 0 {
						redcheck = true
					}
				} else if ctype == "green_playlist" {
					if playlistResponse.ID == 0 {
						greencheck = true
					}
				} else if ctype == "black_playlist" {
					if playlistResponse.ID == 0 {
						blackcheck = true
					}
				}
			}
			if redcheck == true && greencheck == true && blackcheck == true {
				/* here for making slider data null if no playlist has rights or active state */
				IsCheck = true
			}
			for _, playlists := range featuredPlaylistsResponse {
				featuredPlaylist.ID = playlists.ID
				featuredPlaylist.PlaylistType = playlists.PlaylistType
				featuredPlaylist.Title = playlists.Title
				contentId, _ := json.Marshal(playlists.Content)
				json.Unmarshal(contentId, &contentIds)
				var content1 Cont
				res := strconv.Itoa(country)

				var finalcontents []string
				var fetchvalidids ContentIds
				cdb.Debug().Raw("select jsonb_agg(id order by array_position(array[?], id::text)) as ids from (select distinct c.id from content c left join content_variance cv on cv.content_id =c.id left join season s on s.content_id = c.id left join variance_trailer vt on vt.content_variance_id = s.id or vt.content_variance_id = cv.id left join episode e on e.season_id = s.id join playback_item pi2 on pi2.id = cv.playback_item_id or pi2.id = e.playback_item_id join content_rights_country crc on crc.content_rights_id = s.rights_id or crc.content_rights_id = pi2.rights_id where c.status = 1 and (cv.status =1 or cv.status is null) and (s.status = 1 or s.status is null) and ((e.status =1 or e.status is null) or (vt.id is not null)) and crc.country_id = ? and c.id in (?)) as foo", contentIds, country, contentIds).Find(&fetchvalidids)
				validcontentidswithorder, _ := json.Marshal(fetchvalidids.Ids)
				json.Unmarshal(validcontentidswithorder, &finalcontents)
				fmt.Println(len(finalcontents), "lenght of final contents")
				if featured.Type == "Layout A – Smart TV" {
					if playlists.PlaylistType == "red_playlist" {
						if len(finalcontents) < 7 {
							IsCheck = true
						}
					} else if playlists.PlaylistType == "green_playlist" {
						if len(finalcontents) < 1 {
							IsCheck = true
						}
					} else if playlists.PlaylistType == "black_playlist" {
						if len(finalcontents) < 1 {
							IsCheck = true
						}
					}
				} else if featured.Type == "Layout B - STV / Website / Apple TV" {
					// scenario may need to be changed based on requirment
					if playlists.PlaylistType == "red_playlist" {
						if len(finalcontents) < 7 {
							IsCheck = true
						}
					} else if playlists.PlaylistType == "black_playlist" {
						if len(finalcontents) < 1 {
							IsCheck = true
						}
					}
				} else if featured.Type == "layout C - STV - Website - Apple TV" {
					//if playlists.PlaylistType == "black_playlist" {
					if len(finalcontents) < 1 {
						IsCheck = true
					}
					//}
				}
				if contentError := fdb.Debug().Raw("select jsonb_agg(details order by array_position(array[?], cf.content_id::text)) as cnt from content_fragment cf where cf.content_id in (?) and language in (?)  and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(country,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(platform,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and rights_start_date <= now() and rights_end_date >= now() and language = ? and country like '%"+res+"%' ", finalcontents, contentIds, language, contentIds, country, country, finalcontents, platform, platform, language).Find(&content1).Error; contentError != nil {
					fmt.Println(contentError)
					c.JSON(http.StatusInternalServerError, gin.H{"message": contentError})

				}

				var SliderPlaylistVerification VerifyPlaylist
				SliderPlaylistResponseVerifiyBytes, _ := json.Marshal(content1.Cnt)
				json.Unmarshal(SliderPlaylistResponseVerifiyBytes, &SliderPlaylistVerification)
				featuredPlaylist.Content = content1.Cnt
				if len(SliderPlaylistVerification) > 0 {
					featuredPlaylists = append(featuredPlaylists, featuredPlaylist)
				}
			}
			Result.Featured = &FeaturedDetails{featured.ID, featured.Type, featuredPlaylists}
			if len(featuredPlaylists) <= 0 {
				IsCheck = true
			}
		}
		if menuPage.HasMenuPosterImage == true {
			imageryDetails.MobileMenu = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/mobile-menu" + "?c=" + strconv.Itoa(rand.Intn(200000))
			imageryDetails.MobilePosterImage = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
			imageryDetails.MobileMenuPosterImage = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/mobile-menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
		} else {
			imageryDetails.MobileMenu = ""
			imageryDetails.MobilePosterImage = ""
			imageryDetails.MobileMenuPosterImage = ""
		}
		if imageryDetails.MobileMenu == "" || imageryDetails.MobilePosterImage == "" || imageryDetails.MobileMenuPosterImage == "" {
			Result.Imagery = nil
		} else {
			Result.Imagery = &imageryDetails
		}
		if IsCheck {
			Result.Featured = nil
		}
		if Result.Featured == nil && menuPage.Title != "Home" && menuPage.Title != "الصفحة الرئيسية" {
			Result.Type = "VOD"
		}
		// Page Playlists
		var plays Plays
		db.Debug().Raw("select jsonb_agg(p4.id order by pp2.order)::varchar as playlists from page p3 join page_playlist pp2 on pp2.page_id = p3.id join playlist p4 on p4.id = pp2.playlist_id where p4.deleted_by_user_id is null and p4.is_disabled = false and p3.page_key = ?", pageKey).Find(&plays)
		var parts []string
		if err := json.Unmarshal([]byte(plays.Playlists), &parts); err != nil {
			fmt.Println(err)
		}
		for _, ids := range parts {
			playlistIds = append(playlistIds, string(ids))
		}

		if playListError := db.Debug().Raw("select p.playlist_key as id,case when 'en'=lower('"+language+"') then p.english_title when 'ar'=lower('"+language+"') then p.arabic_title end as playlist_type,jsonb_agg(coalesce(case when pi2.multi_tier_content_id is not null then pi2.multi_tier_content_id end, case when pi2.one_tier_content_id is not null then pi2.one_tier_content_id end,case when pi2.season_id is not null then pi2.season_id end,pic.content_id) order by pi2.order)as content from playlist p join playlist_item pi2 on pi2.playlist_id = p.id join playlist_item_content pic on pic.playlist_item_id = pi2.id join page_playlist pp on pp.playlist_id = p.id and pp.page_id =? where p.id in (?) and p.is_disabled='false' and (p.scheduling_start_date  <= now() or p.scheduling_start_date is null ) and (p.scheduling_end_date > now() or p.scheduling_end_date is null) group by p.id,pp.order order by pp.order", menuPage.ID, playlistIds).Find(&featuredPlaylistsResponse).Error; playListError != nil {
			fmt.Println(playListError)
			c.JSON(http.StatusInternalServerError, gin.H{"message": playListError})
			return
		}
		menuply := make([]MenuPlaylists, 0)
		if len(featuredPlaylistsResponse) <= 0 {
			menuPlaylsts = menuply
			fmt.Println(len(featuredPlaylistsResponse))
		}
		for _, playlists := range featuredPlaylistsResponse {
			menuPlaylst.ID = playlists.ID
			menuPlaylst.Title = playlists.PlaylistType
			contentId, _ := json.Marshal(playlists.Content)
			json.Unmarshal(contentId, &contentIds)
			var finalcontentids []string
			var contentidsfororder []string
			var makestring string
			for i, val := range contentIds {
				j := i + 1
				value := strconv.Itoa(j)
				samplestring := " when c.id='" + val + "' or s.id ='" + val + "' then " + value
				makestring = makestring + samplestring
			}

			var filteredids []FilteredIds
			cdb.Raw("select case when s.id is not null then s.id when cv.id is not null then cv.id else c.id end as filteredids,c.id from content c left join season s on s.content_id  =  c.id left join content_variance cv on cv.content_id = c.id left join variance_trailer vt on vt.content_variance_id = s.id or vt.content_variance_id = cv.id left join episode e on e.season_id = s.id where (c.id in (?) or s.id in (?)) and c.status = 1 and (cv.status = 1 or cv.status is null) and (s.status = 1 or s.status is null) and ((e.status = 1 or e.status is null) or (vt.id is not null)) group by c.id,s.id,cv.id order by case "+makestring+" end", contentIds, contentIds).Find(&filteredids)
			for _, val := range filteredids {
				finalcontentids = append(finalcontentids, val.Filteredids)
				contentidsfororder = append(contentidsfororder, val.Id)
			}
			if len(finalcontentids) > 0 {
				var content2 Cont
				res1 := strconv.Itoa(country)
				if contentError := fdb.Debug().Raw("select jsonb_agg(details order by array_position(array[?], cf.content_id::text)) as cnt from content_fragment cf where (cf.content_id in (?) or  cf.content_variance_id in (?))and language in (?)  and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(country,'[]','{}')) from content_fragment cf where (cf.content_id in (?) or  cf.content_variance_id in (?))) as ss )aa where cid= ? )= ? and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(platform,'[]','{}')) from content_fragment cf where (cf.content_id in (?) or cf.content_variance_id in (?))) as ss )aa where cid= ? )= ? and rights_start_date <= now() and rights_end_date >= now() and language = ? and country like '%"+res1+"%' ", contentidsfororder, finalcontentids, finalcontentids, language, finalcontentids, finalcontentids, country, country, finalcontentids, finalcontentids, platform, platform, language).Find(&content2).Error; contentError != nil {
					fmt.Println(contentError)
					c.JSON(http.StatusInternalServerError, gin.H{"message": contentError})
					return
				}
				var PlaylistVerification VerifyPlaylist
				contentId, _ := json.Marshal(content2.Cnt)
				json.Unmarshal(contentId, &PlaylistVerification)
				menuPlaylst.Content = content2.Cnt
				if len(PlaylistVerification) > 0 {
					menuPlaylsts = append(menuPlaylsts, menuPlaylst)
				}
			} else {
				continue
			}
		}
		if len(menuPlaylsts) <= 0 {
			menuPlaylsts = menuply
		}
		Result.Playlists = menuPlaylsts
		jsonData, _ := json.Marshal(Result)
		var request RedisCacheRequest
		url := os.Getenv("REDIS_CACHE_URL")
		request.Key = key //pageKey + language + strconv.Itoa(country) + strconv.Itoa(platform)
		request.Value = string(jsonData)
		_, err := common.PostCurlCall("POST", url, request)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": Result})
}

type PageDetails struct {
	Region   string `json:"region"`
	Platform string `json:"platform"`
}

// func (hs *HandlerService) CreateRedisCacheForHomePage(pageKey int,c *gin.Context) {
func PrepareMenuPageCache(pageKey int, c *gin.Context) int {
	db := c.MustGet("DB").(*gorm.DB)
	fdb := c.MustGet("FDB").(*gorm.DB)
	//pageKey := c.Param("pageKey")
	var pageDetails PageDetails
	var Result MenuPage
	var imageryDetails ImageryDetails
	var featured FeaturedResponse
	var featuredPlaylistsResponse []FeaturedPlaylistsResponse
	var featuredPlaylist FeaturedPlaylists
	var featuredPlaylists []FeaturedPlaylists
	var menuPlaylsts []MenuPlaylists
	var menuPlaylst MenuPlaylists
	var stringsnew, parts, contentIds []string
	var content1 Cont
	// var content []Content
	//var request RedisCacheRequest
	languages := [2]string{"en", "ar"}
	fmt.Println(languages)

	/* Fetching Regions & Platforms Reguarding PageKey */
	if fetcherror := db.Debug().Raw("select STRING_AGG (distinct c.alpha2code, ',') as  region,STRING_AGG (distinct pp2.platform, ',') as platform from page p join page_country pc on pc.page_id = p.id join page_playlist pp on pp.page_id = p.id join play_list_platform plp on plp.play_list_id = pp.playlist_id join country c on c.id = pc.country_id join publish_platform pp2 on pp2.id = plp.target_platform where p.page_key = ?", pageKey).Find(&pageDetails).Error; fetcherror != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
		return 400
	}
	/* Preparing Redis Cache for individual country */
	for _, language := range languages {
		var menuPage MenuPageDetails
		if Error := db.Debug().Raw("select p.id,p.page_key,p.has_menu_poster_image, case when 'en'= '"+language+"' then english_page_friendly_url when 'ar'='"+language+"' then arabic_page_friendly_url end as friendly_url ,case when 'en'= '"+language+"' then english_meta_description when 'ar'='"+language+"' then arabic_meta_description end as seo_description,case when 'en'= '"+language+"' then p.english_title  when 'ar'='"+language+"' then p.arabic_title end as title , jsonb_build_object('id', slider_key, 'type', st.name,'playlists', jsonb_agg(distinct concat(red_area_playlist_id, ',' , green_area_playlist_id, ',' , black_area_playlist_id))) as featured,case when page_type = 1 then 'Home' when page_type = 8 then 'Settings' when page_type = 16 then 'Favourites' when page_type = 0 then 'None' when jsonb_agg(distinct ps.slider_id) = null then 'VOD' end as page_type, jsonb_agg(p2.id)::varchar as playlists from page p join page_slider ps on ps.page_id = p.id join slider s on s.id = ps.slider_id join slider_types st on st.id = s.type join page_playlist pp on pp.page_id = p.id join playlist p2 on p2.id = pp.playlist_id join page_country pc on pc.page_id =p.id join country c on c.id =pc.country_id  where s.deleted_by_user_id is null and s.is_disabled = false and s.scheduling_start_date <= NOW() and s.scheduling_end_date >= NOW() and p.page_key = ?  and p.is_disabled = false  and p.deleted_by_user_id is null group by p.id,s.slider_key ,p.page_key,st.name", pageKey).Find(&menuPage).Error; Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server-Error"})
			return 400
		}
		Result.ID = menuPage.PageKey
		Result.Type = menuPage.PageType
		Result.SeoDescription = menuPage.SeoDescription
		Result.Title = menuPage.Title
		Result.FriendlyUrl = menuPage.FriendlyUrl
		/* Featured Details */
		data, _ := json.Marshal(menuPage.Featured)
		json.Unmarshal(data, &featured)
		Result.Featured.ID = featured.ID
		Result.Featured.Type = featured.Type
		var playlistIds []string
		for _, ids := range featured.Playlists {
			stringsnew = strings.Split(ids, ",")
			playlistIds = stringsnew
		}
		/* playlist response */
		if playListError := db.Debug().Raw("select p.playlist_key as id,p.playlist_type,jsonb_agg(coalesce(pi2.multi_tier_content_id, pi2.one_tier_content_id))as content from playlist p join playlist_item pi2 on pi2.playlist_id = p.id where p.id in (?) group by p.id", playlistIds).Find(&featuredPlaylistsResponse).Error; playListError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server-Error"})
			return 400
		}
		for _, playlists := range featuredPlaylistsResponse {
			featuredPlaylist.ID = playlists.ID
			featuredPlaylist.PlaylistType = playlists.PlaylistType
			contentId, _ := json.Marshal(playlists.Content)
			json.Unmarshal(contentId, &contentIds)

			if contentError := fdb.Debug().Raw("select jsonb_agg(details)::text as cnt from content_fragment cf where cf.content_id in (?) and language in ('"+language+"')", contentIds).Find(&content1).Error; contentError != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "server-Error"})
				return 400
			}
			// json.Unmarshal([]byte(content1.Cnt), &content)
			featuredPlaylist.Content = content1.Cnt
			featuredPlaylists = append(featuredPlaylists, featuredPlaylist)
		}
		Result.Featured.Playlists = featuredPlaylists
		if menuPage.HasMenuPosterImage == true {
			imageryDetails.MobileMenu = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/mobile-menu" + "?c=" + strconv.Itoa(rand.Intn(200000))
			imageryDetails.MobilePosterImage = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
			imageryDetails.MobileMenuPosterImage = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/mobile-menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
		} else {
			imageryDetails.MobileMenu = ""
			imageryDetails.MobilePosterImage = ""
			imageryDetails.MobileMenuPosterImage = ""
		}
		Result.Imagery = &imageryDetails
		if err := json.Unmarshal([]byte(menuPage.Playlists), &parts); err != nil {
			fmt.Println(err)
		}
		var playlists []string
		for _, ids := range parts {
			stringsnew = strings.Split(ids, ",")
			playlists = stringsnew
		}
		if playListError := db.Debug().Raw("select p.playlist_key as id,case when 'en'=lower('"+language+"') then p.english_title when 'ar'=lower('"+language+"') then p.arabic_title end as playlist_type,jsonb_agg(coalesce(pi2.multi_tier_content_id, pi2.one_tier_content_id))as content from playlist p join playlist_item pi2 on pi2.playlist_id = p.id where p.id in (?) group by p.id", playlists).Find(&featuredPlaylistsResponse).Error; playListError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server-Error"})
			return 400
		}
		for _, playlists := range featuredPlaylistsResponse {
			menuPlaylst.ID = playlists.ID
			menuPlaylst.Title = playlists.PlaylistType
			contentId, _ := json.Marshal(playlists.Content)
			json.Unmarshal(contentId, &contentIds)
			if contentError := fdb.Debug().Raw("select jsonb_agg(details)::text as cnt from content_fragment cf where cf.content_id in (?) and language in (?)", contentIds, language).Find(&content1).Error; contentError != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "server-Error"})
				return 400
			}
			// json.Unmarshal([]byte(content1.Cnt), &content)
			menuPlaylst.Content = content1.Cnt
			menuPlaylsts = append(menuPlaylsts, menuPlaylst)
		}
		Result.Playlists = menuPlaylsts
		var country, platform []string
		country = strings.Split(pageDetails.Region, ",")
		platform = strings.Split(pageDetails.Platform, ",")
		fmt.Println(time.Now(), "start time redis")
		var request RedisCacheRequest
		rdb := redis.NewClient(&redis.Options{
			// Addr: "10.33.82.25:6379",
			Addr:     "172.31.32.59:6379",
			Password: "",
			DB:       0,
		})

		pipe := rdb.TxPipeline()
		for _, countryId := range country {
			for _, platformId := range platform {
				jsonData, _ := json.Marshal(Result)
				platform := strings.Join(strings.Fields(strings.ToLower(platformId)), "")
				key := strconv.Itoa(pageKey) + language + countryId + platform
				request.Key = key
				request.Value = string(jsonData)
				errSet := pipe.Set(context.Background(), key, request.Value, 0).Err()
				if errSet != nil {
					fmt.Println(errSet)
					return 400
				}
			}
		}
		fmt.Println(time.Now(), "End time rdis")
		pipe.Exec(context.Background())
	}
	return 200
}

type PageKeys struct {
	PageKey int
	PageId  string
}
type PlaylistIds struct {
	playlistId string
}
type SliderIds struct {
	SliderId string
}
type Counts struct {
	Slidercount   int
	Playlistcount int
	IsDisabled    int
}

var ctx = context.Background()

/*create PageCache for all pending pages */
func (hs *HandlerService) UpdateAllPageCache(c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	fcdb := c.MustGet("DB").(*gorm.DB)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_RDS"),
		Password: "",
		DB:       0,
	})
	var pageKeys []PageKeys
	fdb.Debug().Table("page_sync").Select("page_key").Where("dirty_count >?", 0).Find(&pageKeys)
	var pageSync common.PageSync
	if len(pageKeys) > 0 {
		for _, keys := range pageKeys {
			/*sidemenu delete when page update */
			search := "msapiuat-frontend.z5.com_v1*"
			fmt.Println(search, "??????????????")
			if len(os.Args) > 1 {
				search = os.Args[1]
			}
			iter2 := redisClient.Keys(ctx, search)
			for _, Rkeysidemenu := range iter2.Val() {
				fmt.Println(Rkeysidemenu, ">>>>>>>>>>>>>>.")
				redisClient.Del(ctx, Rkeysidemenu)
				fmt.Println("Redis key " + Rkeysidemenu + "is deleted")
			}
			var counts Counts
			fcdb.Debug().Table("page p").Select("SUM(CASE WHEN ps.slider_id IS NOT NULL THEN 1 ELSE 0 END) AS slidercount,SUM(CASE WHEN pp.playlist_id IS NOT NULL THEN 1 ELSE 0 END) AS playlistcount,SUM(case when p.is_disabled=true then 1 else 0 end) as is_disabled").Joins("left join page_slider ps on ps.page_id = p.id left join page_playlist pp on pp.page_id = p.id").Where("p.page_key =?", keys.PageKey).Find(&counts)
			if (counts.Slidercount > 0 || counts.Playlistcount > 0) || counts.IsDisabled == 0 {
				/* Delete Redis keys Related page_key */
				searchPattern := c.Request.Host + strconv.Itoa(keys.PageKey) + "*"
				if len(os.Args) > 1 {
					searchPattern = os.Args[1]
				}
				fmt.Println(searchPattern, "??????????????")
				iter := redisClient.Keys(ctx, searchPattern)
				for _, Rkey := range iter.Val() {
					fmt.Println(Rkey, ">>>>>>>>>>>>>>>>>")
					redisClient.Del(ctx, Rkey)
					fmt.Println("Redis key " + Rkey + "is deleted")
				}
				redisClient.Close()
				/*prepare Radis cache for page */
				// go MenuPageCache(keys.PageKey, c)
				/* delete page_key from pagesynctable if redis cache cleared */
				fdb.Debug().Where("page_key =?", keys.PageKey).Delete(&pageSync)
				c.JSON(http.StatusOK, "job success")
			} else {
				c.JSON(http.StatusOK, "The page (psge_key "+strconv.Itoa(keys.PageKey)+") has no sliders or playlists. so, its not mandatory to prepare cache for this page.")
				// return
			}
		}
		for _, val := range pageKeys {
			pagekey := strconv.Itoa(val.PageKey)
			go common.ClearRedisKeyForPages(pagekey, c)
		}
		c.JSON(http.StatusOK, "job success")
		return
	}
}

/* Moniter playlist and slider*/
func (hs *HandlerService) UdateDirtyCountPlaylistRelatedPages(c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	db := c.MustGet("DB").(*gorm.DB)

	var playlistId []common.PlaylistSync
	var playlistSync common.PlaylistSync
	if countError := fdb.Debug().Table("playlist_sync").Select("playlist_id,dirty_count").Where("dirty_count >?", 0).Find(&playlistId).Error; countError != nil {
		fmt.Println(countError)
	}
	fmt.Println(playlistId)
	for _, playIds := range playlistId {
		var pageKeys []PageKeys
		if playlistError := db.Debug().Table("playlist p").Select("distinct p2.page_key,p2.id as page_id").
			Joins("left join page_playlist pp on pp.playlist_id =p.id").
			Joins("left join slider s on s.red_area_playlist_id = p.id or s.green_area_playlist_id = p.id or s.black_area_playlist_id = p.id").
			Joins("left join page_slider ps on ps.slider_id = s.id").
			Joins("join page p2 on p2.id =pp.page_id or p2.id = ps.page_id").
			Where("p.id=?", playIds.PlaylistId).Find(&pageKeys).Error; playlistError != nil {
			fmt.Println(playlistError)
			return
		}
		fmt.Println(pageKeys)
		if len(pageKeys) > 0 {
			for _, Keys := range pageKeys {
				updateResult := common.PageSynching(Keys.PageId, Keys.PageKey, c)
				if updateResult == "success" {
					if deleteError := fdb.Where("playlist_id =?", playIds.PlaylistId).Delete(&playlistSync).Error; deleteError != nil {
						fmt.Println(deleteError)
						return
					}
				}
			}
		} else {
			if deleteError := fdb.Where("playlist_id =?", playIds.PlaylistId).Delete(&playlistSync).Error; deleteError != nil {
				fmt.Println(deleteError)
				return
			}
		}
	}
}

func (hs *HandlerService) UpdateDirtycountSliderRelatedPages(c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	db := c.MustGet("DB").(*gorm.DB)
	var pageKeys []PageKeys
	var sliderIds []SliderIds
	var sliderSync common.SliderSync
	if countError := fdb.Debug().Table("slider_sync").Select("slider_id").Where("dirty_count >?", 0).Find(&sliderIds).Error; countError != nil {
		fmt.Println(countError)
	}
	for _, sliIds := range sliderIds {
		if fetchError := db.Debug().Table("page_slider ps").Select("p.id as page_id ,p.page_key ").
			Joins("join page p on p.id=ps.page_id ").
			Where("ps.slider_id =?", sliIds.SliderId).Find(&pageKeys).Error; fetchError != nil {
			fmt.Println(fetchError)
			return
		}
		if len(pageKeys) > 0 {
			for _, pkeys := range pageKeys {
				updateResult := common.PageSynching(pkeys.PageId, pkeys.PageKey, c)
				if updateResult == "success" {
					if deleteError := fdb.Debug().Where("slider_id =?", sliIds.SliderId).Delete(&sliderSync).Error; deleteError != nil {
						fmt.Println(deleteError)
						return
					}
				}
			}
		} else {
			if deleteError := fdb.Debug().Where("slider_id =?", sliIds.SliderId).Delete(&sliderSync).Error; deleteError != nil {
				fmt.Println(deleteError)
				return
			}
		}
	}
}

/*create dirty count with slider*/
func UpdateDirtycountSliderIdRelatedPages(sliderId string, c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	db := c.MustGet("DB").(*gorm.DB)
	var pageKeys []PageKeys
	var sliderSync common.SliderSync

	if fetchError := db.Debug().Table("page_slider ps").Select("p.id as page_id ,p.page_key ").
		Joins("join page p on p.id=ps.page_id ").
		Where("ps.slider_id =?", sliderId).Find(&pageKeys).Error; fetchError != nil {
		fmt.Println(fetchError)
		return
	}
	for _, Keys := range pageKeys {
		updateResult := common.PageSynching(Keys.PageId, Keys.PageKey, c)
		if updateResult == "success" {
			if deleteError := fdb.Debug().Where("slider_sync =?", sliderId).Delete(&sliderSync).Error; deleteError != nil {
				fmt.Println(deleteError)
				return
			}
		}
	}
}

// func prepareUUids(pageKey string, c *gin.Context) {
// 	db := c.MustGet("DB").(*gorm.DB)
// 	var pageDetails PageDetails
// 	if fetcherror := db.Raw("select STRING_AGG (distinct c.alpha2code, ',') as  region,STRING_AGG (distinct pp2.platform, ',') as platform from page p join page_country pc on pc.page_id = p.id join page_playlist pp on pp.page_id = p.id join play_list_platform plp on plp.play_list_id = pp.playlist_id join country c on c.id = pc.country_id join publish_platform pp2 on pp2.id = plp.target_platform where p.page_key = ?", pageKey).Find(&pageDetails).Error; fetcherror != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
// 		return
// 	}
// 	var country, platform []string
// 	languages := [2]string{"en", "ar"}
// 	for _, countryId := range country {
// 		for _, platformId := range platform {
// 			for _, language := range languages {
// 				for i := 0; i < numKeys; i++ {
// 					platform := strings.Join(strings.Fields(strings.ToLower(platformId)), "")
// 					key := pageKey + language + countryId + platform
// 					uuids[i] = key
// 				}
// 			}
// 		}
// 	}
// }

// func createPool() {
// 	pool = &redis.Pool{
// 		MaxIdle:     16,
// 		MaxActive:   16,
// 		IdleTimeout: 3600 * time.Second,
// 		Dial: func() (redis.Conn, error) {
// 			c, err := net.Dial("tcp", "172.31.32.59:6379")
// 			if err != nil {
// 				log.Fatal(err)
// 				return nil, err
// 			}
// 			return redis.NewConn(c, 10*time.Second, 10*time.Second), nil
// 		},
// 		TestOnBorrow: func(c redis.Conn, t time.Time) error {
// 			_, err := c.Do("PING")
// 			return err
// 		},
// 	}
// }

// func massImport(Result interface{}) {
// 	wg.Add(routines)
// 	for i := 0; i < routines; i++ {
// 		go importRoutine(i, pool.Get(), Result)
// 	}

// 	wg.Wait()
// }

// func importRoutine(t int, client redis.Conn, Result interface{}) {
// 	defer client.Flush()
// 	defer wg.Done()
// 	for i := t * portion; i < (t+1)*portion; i++ {
// 		client.Send("SET", uuids[i], Result)
// 	}
// }

// func closePool() {
// 	pool.Close()
// }

/* HomePage-TV(Fetch All Pages) */
func (hs *HandlerService) TvHomePage(c *gin.Context) {
	language := c.Param("lang")
	var country, device string
	if c.Request.URL.Query()["Country"] != nil {
		country = c.Request.URL.Query()["Country"][0]
	} else {
		country = "AE"
	}

	DeviceName := strings.ToLower(c.Request.URL.Query()["device"][0])

	device = strings.ReplaceAll(DeviceName, " ", "")

	if c.Request.URL.Query()["device"] == nil || device == "" {
		device = "web"
	}

	var menuPage SideMenuDetails
	var finalResult []Data
	/* prepare Redis Cache */
	key := c.Request.Host + language + country + device
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "?", "_")
	url := os.Getenv("REDIS_CACHE_URL") + "/" + key

	response, err := common.GetCurlCall(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
		return
	}
	var RedisResponse RedisCacheResponse
	json.Unmarshal(response, &RedisResponse)
	// if "" != "" {
	if RedisResponse.Value != "" {
		if err := json.Unmarshal([]byte(RedisResponse.Value), &menuPage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
			return
		}
	} else {
		/* calling side menu api to collect all page_keys */
		//https://msapifo-uat.weyyak.z5.com/v1/
		//url := "https://ynk2yz6oak.execute-api.ap-south-1.amazonaws.com/weyyak-fo-ms-api-qa/v1/" + language + "/menu?device=" + device
		url := os.Getenv("HOME_PAGE_FOR_TV") + language + "/menu?device=" + device
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
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		json.Unmarshal(body, &menuPage)

		fmt.Println("menuPage.Data------------------>", menuPage.Data)

		var myPlaylist Data

		seen := make(map[string][]Data)

		for _, singlePage := range menuPage.Data {

			fmt.Println("Name--------------->", singlePage.Title)

			/* Calling menu api to collect single page data(playlists,content..etc) */
			// url := "https://ynk2yz6oak.execute-api.ap-south-1.amazonaws.com/weyyak-fo-ms-api-qa/v1/" + language + "/menu/" + strconv.Itoa(singlePage.ID) + "?cascade=2&country=" + country
			url := os.Getenv("HOME_PAGE_FOR_TV") + language + "/menu/" + strconv.Itoa(singlePage.ID) + "?cascade=2&country=" + country
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

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			var data Result
			json.Unmarshal(body, &data)
			if data.Data.Title == "My Playlist" || data.Data.Title == "قائمتي" {
				data.Data.Type = "Favourites"
				data.Data.Imagery = nil
				data.Data.Featured = nil
			}

			if data.Data.ID == 66 {
				myPlaylist = data.Data
			}

			key := data.Data.FriendlyURL
			if _, ok := seen[key]; ok {
				fmt.Println("Duplicate found")
			} else {
				if data.Data.ID != 74 && data.Data.ID != 66 {
					seen[key] = []Data{data.Data}
					finalResult = append(finalResult, data.Data)
				}
			}

		}

		finalResult = append(finalResult, myPlaylist)
		// for fetching staticpages settings,search,exit

		fdb := c.MustGet("DB").(*gorm.DB)
		var StaticPagesResult []Page
		if StaticPagesResultError := fdb.Raw("select * from page p where page_key = 67 or page_key = 89 or page_key = 74 order by english_title desc").Find(&StaticPagesResult).Error; StaticPagesResultError != nil {
			fmt.Println(StaticPagesResultError)
			return
		}
		for _, v := range StaticPagesResult {
			var data Result
			emptyarr := make(Playlists, 0)
			if language == "en" {
				data.Data.ID = v.PageKey
				data.Data.FriendlyURL = v.EnglishPageFriendlyUrl
				data.Data.SeoDescription = v.EnglishMetaDescription
				data.Data.Title = v.EnglishTitle
			} else {
				data.Data.ID = v.PageKey
				data.Data.FriendlyURL = v.ArabicPageFriendlyUrl
				data.Data.SeoDescription = v.ArabicMetaDescription
				data.Data.Title = v.ArabicTitle
			}
			if v.PageKey == 67 {
				data.Data.Type = "Settings"
			} else {
				data.Data.Type = "VOD"
			}
			data.Data.Playlists = emptyarr
			data.Data.Imagery = nil
			data.Data.Featured = nil
			finalResult = append(finalResult, data.Data)
		}

		menuPage.Data = finalResult

		/* prepare redis key pair to serve next-time from redis cache */
		jsonData, _ := json.Marshal(menuPage)
		var request RedisCacheRequest
		RedisUrl := os.Getenv("REDIS_CACHE_URL")
		request.Key = key
		request.Value = string(jsonData)
		_, Curlerr := common.PostCurlCall("POST", RedisUrl, request)
		if Curlerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
			return
		}
	}
	c.JSON(http.StatusOK, menuPage)
	return
}

func MenuPageCache(pageKey int, c *gin.Context) int {
	db := c.MustGet("DB").(*gorm.DB)
	fdb := c.MustGet("FDB").(*gorm.DB)
	var pageDetails PageDetails
	languages := [2]string{"en", "ar"}
	/* Fetching Regions & Platforms Reguarding PageKey */
	if fetcherror := db.Debug().Raw("select STRING_AGG (distinct pc.country_id::text ,',') as region,STRING_AGG (distinct ptp.target_platform::text, ',') as platform from page p left join page_country pc on pc.page_id =p.id left join page_target_platform ptp on ptp.page_id =p.id where p.page_key = ?", pageKey).Find(&pageDetails).Error; fetcherror != nil {
		fmt.Println("Fetch platforms and regions Errors", fetcherror)
		return 400
	}
	var country, platform []string
	country = strings.Split(pageDetails.Region, ",")
	platform = strings.Split(pageDetails.Platform, ",")
	/* Preparing Redis Cache for individual country */
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "172.31.32.59:6379",
	// 	Password: "",
	// 	DB:       0,
	// })
	for _, language := range languages {
		// pipe := rdb.TxPipeline()
		for _, countryId := range country {
			for _, platformId := range platform {
				var Result MenuPage
				var menuPage MenuPageDetails
				var imageryDetails ImageryDetails
				var playlistIds, contentIds []string
				var featured FeaturedResponse
				var featuredPlaylistsResponse []FeaturedPlaylistsResponse
				var featuredPlaylist FeaturedPlaylists
				var featuredPlaylists []FeaturedPlaylists
				var menuPlaylsts []MenuPlaylists
				var menuPlaylst MenuPlaylists

				key1 := c.Request.Host + strconv.Itoa(pageKey) + language + countryId + platformId
				key1 = strings.ReplaceAll(key1, "/", "_")
				key1 = strings.ReplaceAll(key1, "?", "_")
				//url := os.Getenv("REDIS_CACHE_URL") + "/" + key1
				if menupageError := db.Debug().Raw("select p.id,p.page_key,p.has_menu_poster_image, case when 'en'= lower('"+language+"') then english_page_friendly_url when 'ar'=lower('"+language+"') then arabic_page_friendly_url end as friendly_url ,case when 'en'= lower('"+language+"') then english_meta_description when 'ar'=lower('"+language+"') then arabic_meta_description end as seo_description,case when 'en'= lower('"+language+"') then p.english_title  when 'ar'=lower('"+language+"') then p.arabic_title end as title , jsonb_build_object('id', slider_key, 'type', st.name,'playlists', jsonb_agg(distinct concat(red_area_playlist_id, ',' , green_area_playlist_id, ',' , black_area_playlist_id))) as featured,case when page_type = 1 then 'Home' when page_type = 8 then 'Settings' when page_type = 16 then 'Favourites' when page_type = 0 then 'None' when jsonb_agg(distinct ps.slider_id) = null then 'VOD' end as page_type, (select  jsonb_agg(p4.id order by pp2.order)::varchar as playlists from page p3 join page_playlist pp2 on pp2.page_id = p3.id join playlist p4 on p4.id = pp2.playlist_id where p4.deleted_by_user_id is null and p4.is_disabled = false and p3.page_key ="+strconv.Itoa(pageKey)+") from page p left join page_slider ps on ps.page_id = p.id left join slider s on s.id = ps.slider_id left join slider_types st on st.id = s.type left join page_playlist pp on pp.page_id = p.id left join playlist p2 on p2.id = pp.playlist_id left join page_country pc on pc.page_id =p.id left join country c on c.id =pc.country_id  where s.deleted_by_user_id is null and s.is_disabled = false and p.page_key = ?  and p.is_disabled = false and c.id=? and p.deleted_by_user_id is null group by p.id,s.slider_key ,p.page_key,st.name", pageKey, countryId).Find(&menuPage).Error; menupageError != nil {
					fmt.Println("Fetch Page Errors", menupageError)
				}
				Result.ID = menuPage.PageKey
				Result.Type = menuPage.PageType
				Result.SeoDescription = menuPage.SeoDescription
				Result.Title = menuPage.Title
				Result.FriendlyUrl = menuPage.FriendlyUrl
				/* Featured Details */
				data, _ := json.Marshal(menuPage.Featured)
				json.Unmarshal(data, &featured)
				Result.Featured.ID = featured.ID
				Result.Featured.Type = featured.Type
				var stringsnew []string
				for _, ids := range featured.Playlists {
					stringsnew = strings.Split(ids, ",")
				}
				var contentType []string
				for keyv, val := range stringsnew {
					var ctype string
					if keyv == 0 {
						ctype = "red_playlist"
					}
					if keyv == 1 {
						ctype = "green_playlist"
					}
					if keyv == 2 {
						ctype = "black_playlist"
					}
					contentType = append(contentType, ctype)
					/* playlist response */
					var playlistResponse FeaturedPlaylistsResponse
					db.Debug().Raw("select p.playlist_key as id,'"+ctype+"' as playlist_type,case when 'en' = lower(?) then english_title when 'ar' = lower(?) then arabic_title end as title,jsonb_agg(coalesce(pi2.multi_tier_content_id, pi2.one_tier_content_id,case when pi2.season_id is not null then pic.content_id end))as content from playlist p left join playlist_item pi2 on pi2.playlist_id = p.id left join playlist_item_content pic on pic.playlist_item_id = pi2.id where p.id =? and p.is_disabled='false' group by p.id", language, language, val).Find(&playlistResponse)
					if playlistResponse.ID != 0 {
						featuredPlaylistsResponse = append(featuredPlaylistsResponse, playlistResponse)
					}
				}
				for _, playlists := range featuredPlaylistsResponse {
					featuredPlaylist.ID = playlists.ID
					featuredPlaylist.PlaylistType = playlists.PlaylistType
					featuredPlaylist.Title = playlists.Title
					contentId, _ := json.Marshal(playlists.Content)
					json.Unmarshal(contentId, &contentIds)
					var content1 Cont
					if contentError := fdb.Debug().Raw("select jsonb_agg(details) as cnt from content_fragment cf where cf.content_id in (?) and language in (?)  and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(country,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(platform,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and rights_start_date <= now() and rights_end_date >= now()", contentIds, language, contentIds, countryId, countryId, contentIds, platformId, platformId).Find(&content1).Error; contentError != nil {
						fmt.Println("Fetch content Errors", contentError)
					}
					featuredPlaylist.Content = content1.Cnt
					featuredPlaylists = append(featuredPlaylists, featuredPlaylist)
				}
				Result.Featured.Playlists = featuredPlaylists
				if menuPage.HasMenuPosterImage == true {
					imageryDetails.MobileMenu = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/mobile-menu" + "?c=" + strconv.Itoa(rand.Intn(200000))
					imageryDetails.MobilePosterImage = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
					imageryDetails.MobileMenuPosterImage = os.Getenv("IMAGES") + strings.ToLower(menuPage.ID) + "/mobile-menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
				} else {
					imageryDetails.MobileMenu = ""
					imageryDetails.MobilePosterImage = ""
					imageryDetails.MobileMenuPosterImage = ""
				}
				Result.Imagery = &imageryDetails

				var parts []string
				if err := json.Unmarshal([]byte(menuPage.Playlists), &parts); err != nil {
					fmt.Println(err)
				}
				for _, ids := range parts {
					playlistIds = append(playlistIds, string(ids))
				}
				if playListError := db.Debug().Raw("select p.playlist_key as id,case when 'en'=lower('"+language+"') then p.english_title when 'ar'=lower('"+language+"') then p.arabic_title end as playlist_type,jsonb_agg(coalesce(pi2.multi_tier_content_id, pi2.one_tier_content_id,case when pi2.season_id is not null then pic.content_id end))as content from playlist p join playlist_item pi2 on pi2.playlist_id = p.id join playlist_item_content pic on pic.playlist_item_id = pi2.id where p.id in (?) and p.is_disabled='false' group by p.id", playlistIds).Find(&featuredPlaylistsResponse).Error; playListError != nil {
					fmt.Println("Fetch playlist Errors", playListError)
				}
				for _, playlists := range featuredPlaylistsResponse {
					menuPlaylst.ID = playlists.ID
					menuPlaylst.Title = playlists.PlaylistType
					contentId, _ := json.Marshal(playlists.Content)
					json.Unmarshal(contentId, &contentIds)
					var content2 Cont

					if contentErrors := fdb.Debug().Raw("select jsonb_agg(details) as cnt from content_fragment cf where cf.content_id in (?) and language in (?)  and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(country,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(platform,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and rights_start_date <= now() and rights_end_date >= now()", contentIds, language, contentIds, countryId, countryId, contentIds, platformId, platformId).Find(&content2).Error; contentErrors != nil {
						fmt.Println("Fetch contentError Errors", contentErrors)
					}
					menuPlaylst.Content = content2.Cnt
					menuPlaylsts = append(menuPlaylsts, menuPlaylst)
				}
				Result.Playlists = menuPlaylsts
				jsonData, _ := json.Marshal(Result)
				var request RedisCacheRequest
				Redisurl := os.Getenv("REDIS_CACHE_URL")
				request.Key = key1 //pageKey + language + strconv.Itoa(country) + strconv.Itoa(platform)
				request.Value = string(jsonData)
				_, err := common.PostCurlCall("POST", Redisurl, request)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{"message": err})
				}
				// request.Key = url
				// request.Value = string(jsonData)
				// errSet := pipe.Set(context.Background(), url, request.Value, 0).Err()
				// if errSet != nil {
				// 	fmt.Println(errSet)
				// 	return 400
				// }

			}
		}
		// pipe.Exec(context.Background())
	}
	return 200
}
