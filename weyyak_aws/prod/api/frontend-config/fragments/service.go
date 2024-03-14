package fragments

import (
	// "context"
	"encoding/json"
	"fmt"
	"frontend_config/common"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

func UpdatePageFragment(pageId string, c *gin.Context, response chan FragmentUpdate, playlists, sliders []string, defaultSlider string) {
	var fragmentResponse FragmentUpdate
	fragmentResponse.Response = ""
	fragmentResponse.Err = nil
	response <- fragmentResponse
	time.Sleep(1 * time.Second)
	var pageFragment PageFragment
	var menuPage MenuPage
	var pageDetails []MenuPageDetails
	var imageryDetails ImageryDetails
	var featured *FeaturedDetails
	db := c.MustGet("DB").(*gorm.DB)
	fdb := c.MustGet("FDB").(*gorm.DB)
	var countryDetails []Country
	var platformDetails []PublishPlatform
	if err := db.Find(&countryDetails).Error; err != nil {
		fragmentResponse.Response = ""
		fragmentResponse.Err = err
		response <- fragmentResponse
	}
	if err := db.Find(&platformDetails).Error; err != nil {
		fragmentResponse.Response = ""
		fragmentResponse.Err = err
		response <- fragmentResponse
	}
	for _, country := range countryDetails {
		countryCode := common.CountryNames(country.Id)
		for _, platform := range platformDetails {
			platformName := common.DeviceNames(platform.Id)
			var pageFragments []PageFragment
			fdb.Debug().Where("country=? and platform=? and page_id=?", countryCode, platformName, pageId).Delete(&pageFragment)
			for i := 1; i <= 2; i++ {
				var language string
				fields := "id,page_key,page_type as type,ptp.page_order_number,has_mobile_menu,has_menu_poster_image,has_mobile_menu_poster_image"
				if i == 1 {
					language = "en"
					fields += ", english_title as title,english_page_friendly_url as friendly_url,english_meta_description as seo_description"
				} else {
					language = "ar"
					fields += ", arabic_title as title,arabic_page_friendly_url as friendly_url,arabic_meta_description as seo_description"
				}
				groupby := "id,page_key,english_title,arabic_title,english_page_friendly_url,arabic_page_friendly_url,english_meta_description,arabic_meta_description,page_type,ptp.page_order_number"
				if err := db.Debug().Table("page").Select(fields).Joins("inner join page_target_platform ptp on ptp.page_id=page.id inner join page_country pc on pc.page_id=page.id").Where("page.is_disabled=? and page.deleted_by_user_id is null and pc.country_id=? and ptp.target_platform=? and page.id=?", false, country.Id, platform.Id, pageId).Group(groupby).Find(&pageDetails).Error; err != nil {
					fragmentResponse.Response = ""
					fragmentResponse.Err = err
					response <- fragmentResponse
				}
				type PageIds struct {
					Id string `json:"id"`
				}
				var pageids []PageIds
				var ids []string
				if err := db.Debug().Table("page p").Select("p.id").Joins("inner join page_slider ps on ps.page_id=p.id inner join slider s on s.id = ps.slider_id").Where("s.deleted_by_user_id  is null and s.is_disabled =false and p.is_disabled=false and p.deleted_by_user_id is null  and s.scheduling_start_date <=NOW() and s.scheduling_end_date >=NOW() and ps.page_id=?", pageId).Find(&pageids).Error; err != nil {
					fragmentResponse.Response = ""
					fragmentResponse.Err = err
					response <- fragmentResponse
				}
				if pageids != nil {
					for _, pageid := range pageids {
						ids = append(ids, pageid.Id)
					}
				}
				for _, p := range pageDetails {
					menuPage.ID = p.PageKey
					menuPage.FriendlyUrl = p.FriendlyUrl
					menuPage.SeoDescription = p.SeoDescription
					menuPage.Title = p.Title
					menuPage.Type = common.PageTypes(p.Type)
					menuPage.Featured = featured
					if p.Type != 16 && p.Type != 8 {
						exists := common.FindString(ids, p.ID)
						if (p.Type == 0 && exists == true) || p.Type == 1 {
							menuPage.Type = "Home"
						} else {
							menuPage.Type = "VOD"
						}
					}
					if p.HasMobileMenu == true {
						imageryDetails.MobileMenu = os.Getenv("IMAGES") + strings.ToLower(p.ID) + "/mobile-menu" + "?c=" + strconv.Itoa(rand.Intn(200000))
						imageryDetails.MobilePosterImage = os.Getenv("IMAGES") + strings.ToLower(p.ID) + "/menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
						imageryDetails.MobileMenuPosterImage = os.Getenv("IMAGES") + strings.ToLower(p.ID) + "/mobile-menu-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
					} else {
						imageryDetails.MobileMenu = ""
						imageryDetails.MobilePosterImage = ""
						imageryDetails.MobileMenuPosterImage = ""
					}

					menuPage.Imagery = imageryDetails
					Response := make(map[string]interface{})
					Response["response_data"] = menuPage
					data, _ := json.Marshal(Response)
					type PageResponse struct {
						ResponseData postgres.Jsonb `json:"response_data"`
					}
					var pr PageResponse
					json.Unmarshal(data, &pr)
					pageFragment.PageId = p.ID
					pageFragment.PageKey = p.PageKey
					pageFragment.Country = common.CountryNames(country.Id)
					pageFragment.Platform = common.DeviceNames(platform.Id)
					pageFragment.Language = language
					pageFragment.Details = pr.ResponseData
					pageFragment.PageOrder = p.PageOrderNumber
					pageFragments = append(pageFragments, pageFragment)
				}
			}
			fmt.Printf("%v - %v ", country.Alpha2code, len(pageFragments))
			var InserPageFragments []interface{}
			for _, page := range pageFragments {
				InserPageFragments = append(InserPageFragments, page)
			}
			err := gormbulk.BulkInsert(fdb, InserPageFragments, 3000)
			if err != nil {
				fragmentResponse.Response = ""
				fragmentResponse.Err = err
				response <- fragmentResponse
			}
			response := make(chan FragmentUpdate)
			if len(playlists) > 0 {
				pWhere := ""
				for i, playlist := range playlists {
					order := strconv.Itoa(i + 1)
					if i == 0 {
						pWhere += "pp.page_id='" + pageId + "' and ((pp.playlist_id='" + playlist + "' and pp.order=" + order + ")"
					} else {
						pWhere += " or (pp.playlist_id='" + playlist + "' and pp.order=" + order + ")"
					}
					pWhere += ")"
				}
				var pagePlaylist PagePlaylist
				if rows := db.Debug().Table("page_playlist pp").Where(pWhere).Find(&pagePlaylist).RowsAffected; int(rows) != len(playlists) {
					for _, playlist := range playlists {
						go UpdatePlaylistFragment(playlist, pageId, c, response, country.Id, platform.Id, 1)
					}
				}
			} else {
				var playlistFrgment PlaylistFragment
				fdb.Debug().Where("page_id=? and country=? and platform=?", pageId, countryCode, platformName).Delete(&playlistFrgment)
			}
			if len(sliders) > 0 {
				sWhere := ""
				for i, slider := range sliders {
					order := strconv.Itoa(i + 1)
					if i == 0 {
						//sWhere += "ps.page_id='" + pageId + "' and ((ps.slider_id='" + slider + "' and ps.order=" + order + ")"  // old code for bkp having extra "("
						sWhere += "ps.page_id='" + pageId + "' and (ps.slider_id='" + slider + "' and ps.order=" + order + ")"
					} else {
						sWhere += " or (ps.slider_id='" + slider + "' and ps.order=" + order + ")"
					}
					//sWhere += ")" // old code extra close ")"
				}
				var pageSlider PageSlider
				if rows := db.Debug().Table("page_slider ps").Where(sWhere).Find(&pageSlider).RowsAffected; int(rows) != len(sliders) {
					for _, slider := range sliders {
						go CreateSliderResponse(slider, pageId, c, response, country.Id, platform.Id, 1)
					}
				}
			} else if defaultSlider != "" {
				var sliderFrgment SliderFragment
				fdb.Debug().Where("page_id=? and country=? and platform=? and slider_id!=?", pageId, countryCode, platformName, defaultSlider).Delete(&sliderFrgment)
			} else {
				var sliderFrgment SliderFragment
				fdb.Debug().Where("page_id=? and country=? and platform=?", pageId, countryCode, platformName).Delete(&sliderFrgment)

			}
			if defaultSlider != "" {
				var pageSlider PageSlider
				if rows := db.Debug().Where("page_id=? and slider_id=? and order=0").Find(&pageSlider).RowsAffected; rows == 0 {
					go CreateSliderResponse(defaultSlider, pageId, c, response, country.Id, platform.Id, 1)
				}
			} else {
				var sliderFrgment SliderFragment
				fdb.Debug().Where("page_id=? and country=? and platform=?", pageId, countryCode, platformName).Delete(&sliderFrgment)
			}
		}
	}
	fragmentResponse.Response = "Page Fragment Updated Successfully."
	fragmentResponse.Err = nil
	response <- fragmentResponse
}
func UpdatePlaylistFragment(playlistId, pageId string, c *gin.Context, response chan FragmentUpdate, countryId, platformId, ignoreAsync int) {
	var fragmentResponse FragmentUpdate
	if ignoreAsync == 0 {
		fragmentResponse.Response = ""
		fragmentResponse.Err = nil
		response <- fragmentResponse
		time.Sleep(1 * time.Second)
	}
	//playlistId = "4adcd871-8e9f-4ba1-82fe-e15c3a468b4f"

	//db := c.MustGet("CDB").(*gorm.DB)
	fcdb := c.MustGet("DB").(*gorm.DB)
	// cdb := c.MustGet("CDB").(*gorm.DB)
	type PlaylistPageIds struct {
		PageId  string `json:"page_id"`
		PageKey int    `json:"page_key"`
	}
	var pageIds []PlaylistPageIds
	// var ids []string
	var where string
	if pageId != "" {
		where = "pp.playlist_id='" + playlistId + "' and pp.page_id='" + pageId + "'"
	} else {
		where = "pp.playlist_id='" + playlistId + "'"
	}
	if err := fcdb.Debug().Table("page_playlist pp").Select("distinct(pp.page_id),p.page_key").Joins("join page p on p.id =pp.page_id").Where(where).Find(&pageIds).Error; err != nil {
		fragmentResponse.Response = ""
		fragmentResponse.Err = err
		response <- fragmentResponse
	}
	var playlist Playlist
	if err := fcdb.Debug().Select("id,english_title,arabic_title,scheduling_start_date,scheduling_end_date,deleted_by_user_id,is_disabled,created_at,playlist_key,modified_at,playlist_type").Table("playlist").Where("id=? and is_disabled =false and deleted_by_user_id is null and (scheduling_start_date <=now() or scheduling_start_date is null) and (scheduling_end_date >=now() or scheduling_end_date is null)", playlistId).Find(&playlist).Error; err != nil {
		fragmentResponse.Response = ""
		fragmentResponse.Err = err
		response <- fragmentResponse
	}
	var countryDetails []Country
	var platformDetails []PublishPlatform
	if countryId != 0 {
		if err := fcdb.Debug().Where("id=?", countryId).Find(&countryDetails).Error; err != nil {
			fragmentResponse.Response = ""
			fragmentResponse.Err = err
			response <- fragmentResponse
		}
	} else {
		if err := fcdb.Debug().Find(&countryDetails).Error; err != nil {
			fragmentResponse.Response = ""
			fragmentResponse.Err = err
			response <- fragmentResponse
		}
	}

	if platformId != 0 {
		if err := fcdb.Debug().Where("id=?", platformId).Find(&platformDetails).Error; err != nil {
			fragmentResponse.Response = ""
			fragmentResponse.Err = err
			response <- fragmentResponse
		}
	} else {
		if err := fcdb.Debug().Find(&platformDetails).Error; err != nil {
			fragmentResponse.Response = ""
			fragmentResponse.Err = err
			response <- fragmentResponse
		}
	}
	contentIds, err := PlaylistItemContents(playlist.ID, c)
	if err != nil {
		fragmentResponse.Response = ""
		fragmentResponse.Err = err
		response <- fragmentResponse
	}
	for _, country := range countryDetails {
		countryCode := common.CountryNames(country.Id)
		for _, platform := range platformDetails {
			platformName := common.DeviceNames(platform.Id)
			var playlistFragment []PlaylistFragment
			language := "en"
			for i := 1; i <= 2; i++ {
				if i == 2 {
					language = "ar"
				}
				var pagePlaylist MenuPlaylists

				for _, details := range pageIds {
					var wg sync.WaitGroup
					pagePlaylist.ID = playlist.PlaylistKey
					pagePlaylist.PlaylistType = playlist.PlaylistType

					pagePlaylist.Content = []PlaylistContent{}
					pagePlaylist.PageContent = []PageContent{}
					if playlist.PlaylistType == "pagecontent" {
						pagePlaylist.Title = nil
						playlistPages := make(chan []PageContent)
						wg.Add(1)
						go PlaylistPages(playlist.ID, language, playlistPages, c)
						defer wg.Done()
						if playlistPages != nil {
							pagePlaylist.PageContent = <-playlistPages
						}
					} else if playlist.PlaylistType == "content" {
						pagePlaylist.Title = &playlist.EnglishTitle
						if language != "en" {
							pagePlaylist.Title = &playlist.ArabicTitle
						}
						playlistContent := make(chan []PlaylistContent)
						wg.Add(1)
						go PlaylistContentDetails(contentIds, language, country.Id, c, playlistContent)
						defer wg.Done()
						contentDetails := <-playlistContent
						if contentDetails == nil {
							continue
						}
						var contents []PlaylistContent
						for _, id := range contentIds {
							for _, content := range contentDetails {
								if id.ContentId == content.ContentId {
									contents = append(contents, content)
								}
							}
						}
						pagePlaylist.Content = contents
					}

					playlistFragmentdetails := make(chan PlaylistFragment)
					wg.Add(1)
					go PlaylistFragmentDetails(pagePlaylist, details.PageKey, details.PageId, countryCode, platformName, language, playlist.ID, playlistFragmentdetails)
					defer wg.Done()
					playlistFragment = append(playlistFragment, <-playlistFragmentdetails)
				}
			}
			var pFragment PlaylistFragment
			fdb := c.MustGet("FDB").(*gorm.DB)
			fdb.Debug().Table("playlist_fragment").Where("country=? and playlist_id=? and platform=?", countryCode, playlistId, platformName).Delete(&pFragment)
			// db.Where("country=? and playlist_id=?", countryCode, playlistId).Delete(&pFragment)
			fmt.Printf("%v - %v ", platformName, len(playlistFragment))
			var InsertPlaylistFragments []interface{}
			for _, playlist := range playlistFragment {
				InsertPlaylistFragments = append(InsertPlaylistFragments, playlist)
			}
			err := gormbulk.BulkInsert(fdb, InsertPlaylistFragments, 3000)
			if err != nil {
				fragmentResponse.Response = ""
				fragmentResponse.Err = err
				response <- fragmentResponse
			}
		}
	}
	fragmentResponse.Response = "Playlist Fragment Updated Successfully."
	fragmentResponse.Err = nil
	response <- fragmentResponse
}
func PlaylistContentDetails(contentIds []PlaylistContentIds, language string, country int, c *gin.Context, playlistContentDetails chan []PlaylistContent) {
	var request PlaylistContentRequest
	var Ids []string
	for _, content := range contentIds {
		Ids = append(Ids, content.ContentId)
	}
	request.Ids = Ids
	request.Language = language
	request.Country = country
	go GetPlaylistContents(request, c, playlistContentDetails)
	return
}
func GetPlaylistContents(request PlaylistContentRequest, c *gin.Context, playlistContentDetails chan []PlaylistContent) {
	db := c.MustGet("CDB").(*gorm.DB)
	var contents []Content
	var playlistContent PlaylistContent
	var playlistContents []PlaylistContent
	if err := db.Debug().Where("id in(?) and status=? and deleted_by_user_id is null", request.Ids, 1).Find(&contents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, content := range contents {
		ContentDetails := make(chan PlaylistContent)
		if content.ContentTier == 1 {
			go OneTierContentDetails(content.Id, request.Language, request.Country, c, ContentDetails)
			playlistContent = <-ContentDetails
			playlistContent.ContentId = content.Id
		} else {
			go MultiTierContentDetails(content.Id, request.Language, request.Country, c, ContentDetails, 1)
			playlistContent = <-ContentDetails
			playlistContent.ContentId = content.Id
		}

		if playlistContent.ID != 0 {
			playlistContents = append(playlistContents, playlistContent)
		}
	}
	playlistContentDetails <- playlistContents
}

func PlaylistPages(playlistId, language string, pageContents chan []PageContent, c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var pageContent []PageContent
	var finalResult []PageContent
	fields := "p.id as key,p.page_key as id,p.page_type::text as type"
	if language == "en" {
		fields += ",p.english_title as title,p.english_page_friendly_url as friendly_url,p.english_meta_description as seo_description"
	} else {
		fields += ",p.arabic_title as title,p.arabic_page_friendly_url as friendly_url,p.arabic_meta_description as seo_description"
	}
	if err := db.Debug().Table("playlist_item pi2").Select(fields).Joins("join page p on p.id=pi2.group_by_page_id").Where("playlist_id=? and p.is_disabled =false and p.deleted_by_user_id is null", playlistId).Order("pi2.order").Find(&pageContent).Error; err != nil {
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

func PlaylistFragmentDetails(details MenuPlaylists, pageKey int, pageId, countryCode, platformName, languageCode, playlistId string, playlistFragmentDetails chan PlaylistFragment) {
	var playlistFragment PlaylistFragment
	playlistResponse := make(map[string]interface{})
	playlistResponse["playlist_response"] = details
	playlistdata, _ := json.Marshal(playlistResponse)
	type PlaylistResponse struct {
		ResponseData postgres.Jsonb `json:"playlist_response"`
	}
	var playlist PlaylistResponse
	json.Unmarshal(playlistdata, &playlist)
	playlistFragment.PageId = pageId
	playlistFragment.PageKey = pageKey
	playlistFragment.Country = countryCode
	playlistFragment.Platform = platformName
	playlistFragment.Details = playlist.ResponseData
	playlistFragment.Language = languageCode
	playlistFragment.PlaylistId = playlistId
	playlistFragmentDetails <- playlistFragment
	return
}

func PlaylistItemContents(playlistId string, c *gin.Context) ([]PlaylistContentIds, error) {
	db := c.MustGet("DB").(*gorm.DB)
	var contentIds []PlaylistContentIds
	if err := db.Debug().Table("playlist_item_content pic").Select("pic.content_id").Joins("inner join playlist_item pi2 on pi2.id=pic.playlist_item_id inner join playlist p on p.id=pi2.playlist_id").Where("p.id =?", playlistId).Order("pi2.order asc").Find(&contentIds).Error; err != nil {
		return nil, err
	}
	return contentIds, nil
}

func OneTierContentDetails(contentId string, language string, country int, c *gin.Context, ContentDetails chan PlaylistContent) {
	db := c.MustGet("CDB").(*gorm.DB)
	var playlistContent PlaylistContent
	var onetierContentResult OnetierContentResult
	// var contentImageryDetails ContentImageryDetails
	fields, join, where, groupBy := common.OnetierContentQuery(contentId, language)
	if row := db.Debug().Table("content c").Select(fields).Joins(join).Where(where, contentId, country).Group(groupBy).Find(&onetierContentResult).RowsAffected; row != 0 {
		playlistContent.ID = onetierContentResult.ID
		playlistContent.AgeRating = common.AgeRatings(onetierContentResult.AgeRating, language)
		playlistContent.VideoId = onetierContentResult.VideoId
		friendlyUrl := strings.ToLower(onetierContentResult.FriendlyUrl)
		playlistContent.FriendlyUrl = strings.Replace(friendlyUrl, " ", "-", -1)
		playlistContent.ContentType = strings.ToLower(onetierContentResult.ContentType)
		playlistContent.Synopsis = onetierContentResult.Synopsis
		playlistContent.ProductionYear = onetierContentResult.ProductionYear
		playlistContent.Length = onetierContentResult.Length
		playlistContent.Title = onetierContentResult.Title
		if onetierContentResult.SeoDescription != "" {
			playlistContent.SeoDescription = onetierContentResult.SeoDescription
		} else {
			playlistContent.SeoDescription = onetierContentResult.Synopsis
		}
		playlistContent.TranslatedTitle = onetierContentResult.TranslatedTitle
		if onetierContentResult.SeoTitle != "" {
			playlistContent.SeoTitle = onetierContentResult.SeoTitle
		} else {
			playlistContent.SeoTitle = onetierContentResult.Title
		}

		playlistContent.InsertedAt = onetierContentResult.InsertedAt
		playlistContent.ModifiedAt = onetierContentResult.ModifiedAt
		playlistContent.Geoblock = false
		castId := onetierContentResult.CastId
		//movie details
		var movieDetails []ContentMovieDetails
		var movie PlaylistMovie
		var movies []PlaylistMovie
		movieFields, join, where, groupBy := common.MovieDetailsQuery(language)
		if err := db.Debug().Table("content c").Select(movieFields).Joins(join).Where(where, contentId, country).Group(groupBy).Order("cv.order").Find(&movieDetails).Error; err != nil {
			fmt.Println("err", err)
			return
		}
		for _, details := range movieDetails {
			var SubsPlans []ContentSubsPlans
			fields, join, where := common.ContentPlansQuery(1)
			if err := db.Debug().Table("content_variance cv").Select(fields).Joins(join).Where(where, details.Id).Find(&SubsPlans).Error; err != nil {
				fmt.Println("err", err)
				return
			}
			plans := make([]int, 0)
			rights := make([]int, 0)
			for _, plan := range SubsPlans {
				plans = append(plans, plan.SubscriptionPlanId)
			}
			movie.ID = 0
			movie.Title = details.Title
			movie.Geoblock = false
			movie.DigitalRightType = details.DigitalRightsType
			movie.DigitalRightsRegions = rights
			movie.SubscriptiontPlans = plans
			movie.InsertedAt = details.InsertedAt
			movies = append(movies, movie)
		}
		playlistContent.Movies = movies
		//Imaginery Details
		var Imagery ContentImageryDetails
		ImageryDetails := make(chan ContentImageryDetails)
		if onetierContentResult.HasPosterImage == true {
			go OnetierImagery(contentId, onetierContentResult.ContentVersionId, ImageryDetails)
			Imagery = <-ImageryDetails
		} else {
			Imagery.Thumbnail = ""
			Imagery.Backdrop = ""
			Imagery.MobileImg = ""
			Imagery.FeaturedImg = ""
			Imagery.Banner = ""
		}
		playlistContent.Imagery = Imagery

		actors := make([]string, 0)
		genres := make([]string, 0)
		tags := make([]string, 0)
		var actorIds ActorIds
		fields, join, where, groupBy := common.ContentActorsQuery(language)
		if err := db.Debug().Select(fields).Table("content_cast cc").Joins(join).Where(where, castId).Group(groupBy).Find(&actorIds).Error; err == nil {
			if actorIds.Actors != "" {
				actors = strings.Split(actorIds.Actors, ",")
			}
		}
		keys := make(map[string]bool)
		actorsList := make([]string, 0)
		for _, entry := range actors {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				actorsList = append(actorsList, entry)
			}
		}
		var actorFields string
		if language == "en" {
			actorFields += "english_name as name"
		} else {
			actorFields += "arabic_name as name"
		}
		type Names struct {
			Name string `json:"name"`
		}
		var mainActor, mainActress Names
		var err error
		if actorIds.MainActorId != "" && actorIds.MainActressId != "" {
			err = db.Debug().Raw("select "+actorFields+" from actor where id =?", actorIds.MainActorId).Find(&mainActor).Error
			playlistContent.MainActor = mainActor.Name
			err = db.Raw("select "+actorFields+" from actor where id=?", actorIds.MainActressId).Find(&mainActress).Error
			playlistContent.MainActress = mainActress.Name
			if err != nil {
				fmt.Println(err)
			}
		} else if actorIds.MainActorId != "" {
			err = db.Debug().Raw("select "+actorFields+" from actor where id=?", actorIds.MainActorId).Find(&mainActor).Error
			playlistContent.MainActor = mainActor.Name
			if err != nil {
				fmt.Println(err)
			}
		} else if actorIds.MainActressId != "" {
			err = db.Debug().Raw("select "+actorFields+" from actor where id=?", actorIds.MainActressId).Find(&mainActress).Error
			playlistContent.MainActress = mainActress.Name
			if err != nil {
				fmt.Println(err)
			}
		}
		playlistContent.Cast = actorsList
		var genreNames []Names
		fields, join, where, groupBy = common.ContentGenresQuery(language)
		if err := db.Debug().Select(fields).Table("genre g").Joins(join).Where(where, contentId).Group(groupBy).Order("g.id").Find(&genreNames).Error; err != nil {
			fmt.Println("err", err)
			return
		}
		for _, genre := range genreNames {
			genres = append(genres, genre.Name)
		}
		playlistContent.Genres = genres
		var tagNames []Names
		fields, join, where, groupBy = common.ContentTagsQuery()
		if err := db.Debug().Select(fields).Table("textual_data_tag tdt").Joins(join).Where(where, contentId).Group(groupBy).Order("tdt.name").Find(&tagNames).Error; err != nil {
			fmt.Println("err", err)
			return
		}
		for _, tag := range tagNames {
			tags = append(tags, tag.Name)
		}
		playlistContent.Tags = tags
	}
	playlistContent.SchedulingDateTime = onetierContentResult.SchedulingDateTime
	playlistContent.DigitalRightsEndDate = onetierContentResult.DigitalRightsEndDate
	playlistContent.DigitalRightsStartDate = onetierContentResult.DigitalRightsStartDate
	ContentDetails <- playlistContent
	return
}

func MultiTierContentDetails(contentId string, language string, country int, c *gin.Context, ContentDetails chan PlaylistContent, seasonLimit int) {
	db := c.MustGet("CDB").(*gorm.DB)
	var playlistContent PlaylistContent
	var onetierContentResult OnetierContentResult
	// var contentImageryDetails ContentImageryDetails
	fields, join, where, groupBy := common.MultitierContentQuery(contentId, language)
	if row := db.Debug().Table("content c").Select(fields).Joins(join).Where(where, contentId, country).Group(groupBy).Order("e.number,s.number").Limit(1).Find(&onetierContentResult).RowsAffected; row != 0 {
		playlistContent.ID = onetierContentResult.ID
		playlistContent.AgeRating = common.AgeRatings(onetierContentResult.AgeRating, language)
		playlistContent.VideoId = onetierContentResult.VideoId
		friendlyUrl := strings.ToLower(onetierContentResult.FriendlyUrl)
		playlistContent.FriendlyUrl = strings.Replace(friendlyUrl, " ", "-", -1)
		playlistContent.ContentType = strings.ToLower(onetierContentResult.ContentType)
		playlistContent.Synopsis = onetierContentResult.Synopsis
		playlistContent.ProductionYear = onetierContentResult.ProductionYear
		playlistContent.Length = onetierContentResult.Length
		playlistContent.Title = onetierContentResult.Title
		if onetierContentResult.SeoDescription != "" {
			playlistContent.SeoDescription = onetierContentResult.SeoDescription
		} else {
			playlistContent.SeoDescription = onetierContentResult.Synopsis
		}
		playlistContent.TranslatedTitle = onetierContentResult.TranslatedTitle
		if onetierContentResult.SeoTitle != "" {
			playlistContent.SeoTitle = onetierContentResult.SeoTitle
		} else {
			playlistContent.SeoTitle = onetierContentResult.Title
		}
		//playlistContent.SeoDescription = onetierContentResult.SeoDescription
		//playlistContent.TranslatedTitle = onetierContentResult.TranslatedTitle
		//playlistContent.SeoTitle = onetierContentResult.SeoTitle
		playlistContent.InsertedAt = onetierContentResult.InsertedAt
		playlistContent.ModifiedAt = onetierContentResult.ModifiedAt
		playlistContent.Geoblock = false
		castId := onetierContentResult.CastId
		//season details
		var seasonDetails []ContentSeasonDetails
		var season PlaylistContentSeasons
		var seasons []PlaylistContentSeasons
		fields, join, where, groupBy := common.SeasonDetailsQuery(language)
		if seasonLimit == 1 {
			if err := db.Debug().Select(fields).Table("season s").Joins(join).Where(where, contentId, country).Group(groupBy).Limit(1).Order("s.number,s.created_at asc").Find(&seasonDetails).Error; err != nil {
				fmt.Println("err", err)
				return
			}
		} else {
			if err := db.Debug().Select(fields).Table("season s").Joins(join).Where(where, contentId, country).Group(groupBy).Order("s.number,s.created_at asc").Find(&seasonDetails).Error; err != nil {
				fmt.Println("err", err)
				return
			}
		}
		var seasonId string
		if seasonDetails != nil {
			for i, details := range seasonDetails {
				if i == 0 {
					seasonId = details.ID
				}
				var SubsPlans []ContentSubsPlans
				fields, join, where := common.ContentPlansQuery(2)
				if err := db.Debug().Table("content_rights_plan crp").Select(fields).Joins(join).Where(where, details.ID).Find(&SubsPlans).Error; err != nil {
					fmt.Println("err", err)
					return
				}
				plans := make([]int, 0)
				// rights := make([]int, 0)
				for _, plan := range SubsPlans {
					plans = append(plans, plan.SubscriptionPlanId)
				}
				if details.LanguageType == 2 {
					season.Dubbed = true
				} else {
					season.Dubbed = false
				}
				season.ID = details.SeasonKey
				season.Title = details.Title
				season.Geoblock = false
				season.DigitalRightType = details.DigitalRightsType
				season.SeasonNumber = details.SeasonNumber
				// season.SeoDescription = details.SeoDescription
				// season.SeoTitle = details.SeoTitle
				if details.SeoDescription != "" {
					season.SeoDescription = details.SeoDescription
				} else {
					season.SeoDescription = details.Synopsis
				}
				if details.SeoTitle != "" {
					season.SeoTitle = details.SeoTitle
				} else {
					season.SeoTitle = details.Title
				}
				season.DigitalRightsRegions = nil
				season.SubscriptiontPlans = plans
				if seasonLimit == 0 {
					//Season Imaginery Details
					var Imagery ContentImageryDetails
					ImageryDetails := make(chan ContentImageryDetails)
					if details.HasPosterImage == true {
						go MultitierImagery(contentId, details.ID, ImageryDetails)
						Imagery = <-ImageryDetails
					} else {
						Imagery.Thumbnail = ""
						Imagery.Backdrop = ""
						Imagery.MobileImg = ""
						Imagery.FeaturedImg = ""
						Imagery.Banner = ""
					}
					season.Imagery = &Imagery
					var seasonEpisodes []SeasonEpisodes
					var episodesResponse []SeasonEpisodes
					fields, join, where, groupBy := common.SeasonEpisodesQuery(language)
					if err := db.Debug().Select(fields).Table("episode e").Joins(join).Where(where, details.ID).Group(groupBy).Order("e.number asc").Find(&seasonEpisodes).Error; err == nil {
						for _, episode := range seasonEpisodes {
							//episode Imaginery Details
							var Imagery ContentImageryDetails
							ImageryDetails := make(chan ContentImageryDetails)
							if episode.HasPosterImage == true {
								go EpisodeImagery(contentId, details.ID, episode.EpisodeId, ImageryDetails)
								Imagery = <-ImageryDetails
							} else {
								Imagery.Thumbnail = ""
								Imagery.Backdrop = ""
								Imagery.MobileImg = ""
								Imagery.FeaturedImg = ""
								Imagery.Banner = ""
							}
							episode.Imagery = Imagery
							episodesResponse = append(episodesResponse, episode)
						}
						season.Episodes = episodesResponse
					}
				}
				seasons = append(seasons, season)
			}
		}
		playlistContent.Seasons = seasons
		//Imaginery Details
		var Imagery ContentImageryDetails
		ImageryDetails := make(chan ContentImageryDetails)
		if onetierContentResult.HasPosterImage == true {
			go MultitierImagery(contentId, seasonId, ImageryDetails)
			Imagery = <-ImageryDetails
		} else {
			Imagery.Thumbnail = ""
			Imagery.Backdrop = ""
			Imagery.MobileImg = ""
			Imagery.FeaturedImg = ""
			Imagery.Banner = ""
		}
		playlistContent.Imagery = Imagery
		actors := make([]string, 0)
		genres := make([]string, 0)
		tags := make([]string, 0)
		var actorIds ActorIds

		fields, join, where, groupBy = common.ContentActorsQuery(language)
		if err := db.Debug().Select(fields).Table("content_cast cc").Joins(join).Where(where, castId).Group(groupBy).Find(&actorIds).Error; err == nil {
			if actorIds.Actors != "" {
				actors = strings.Split(actorIds.Actors, ",")
			}
		}
		keys := make(map[string]bool)
		actorsList := make([]string, 0)
		for _, entry := range actors {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				actorsList = append(actorsList, entry)
			}
		}
		var actorFields string
		if language == "en" {
			actorFields += "english_name as name"
		} else {
			actorFields += "arabic_name as name"
		}
		type Names struct {
			Name string `json:"name"`
		}
		var mainActor, mainActress Names
		var err error
		if actorIds.MainActorId != "" && actorIds.MainActressId != "" {
			err = db.Debug().Raw("select "+actorFields+" from actor where id =?", actorIds.MainActorId).Find(&mainActor).Error
			playlistContent.MainActor = mainActor.Name
			err = db.Debug().Raw("select "+actorFields+" from actor where id=?", actorIds.MainActressId).Find(&mainActress).Error
			playlistContent.MainActress = mainActress.Name
			if err != nil {
				fmt.Println(err)
			}
		} else if actorIds.MainActorId != "" {
			err = db.Debug().Raw("select "+actorFields+" from actor where id=?", actorIds.MainActorId).Find(&mainActor).Error
			playlistContent.MainActor = mainActor.Name
			if err != nil {
				fmt.Println(err)
			}
		} else if actorIds.MainActressId != "" {
			err = db.Debug().Raw("select "+actorFields+" from actor where id=?", actorIds.MainActressId).Find(&mainActress).Error
			playlistContent.MainActress = mainActress.Name
			if err != nil {
				fmt.Println(err)
			}
		}
		playlistContent.Cast = actorsList
		var genreNames []Names
		fields, join, where, groupBy = common.ContentGenresQuery(language)
		if err := db.Debug().Select(fields).Table("genre g").Joins(join).Where(where, contentId).Group(groupBy).Order("g.id").Find(&genreNames).Error; err != nil {
			fmt.Println("err", err)
			return
		}
		for _, genre := range genreNames {
			genres = append(genres, genre.Name)
		}
		playlistContent.Genres = genres
		var tagNames []Names
		fields, join, where, groupBy = common.ContentTagsQuery()
		if err := db.Debug().Select(fields).Table("textual_data_tag tdt").Joins(join).Where(where, contentId).Group(groupBy).Order("tdt.name").Find(&tagNames).Error; err != nil {
			fmt.Println("err", err)
			return
		}
		for _, tag := range tagNames {
			tags = append(tags, tag.Name)
		}
		playlistContent.Tags = tags
	}
	ContentDetails <- playlistContent
	return
}

func OnetierImagery(contentId, varianceId string, ImageryDetails chan ContentImageryDetails) {
	var Imagery ContentImageryDetails
	Imagery.Thumbnail = os.Getenv("IMAGES") + contentId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.Backdrop = os.Getenv("IMAGES") + contentId + "/details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.MobileImg = os.Getenv("IMAGES") + contentId + "/mobile-details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.FeaturedImg = os.Getenv("IMAGES") + contentId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.Banner = os.Getenv("IMAGES") + contentId + "/" + varianceId + "/overlay-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	ImageryDetails <- Imagery
}
func MultitierImagery(contentId, seasonId string, ImageryDetails chan ContentImageryDetails) {
	var Imagery ContentImageryDetails
	Imagery.Thumbnail = os.Getenv("IMAGES") + contentId + "/" + seasonId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.Backdrop = os.Getenv("IMAGES") + contentId + "/" + seasonId + "/details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.MobileImg = os.Getenv("IMAGES") + contentId + "/" + seasonId + "/mobile-details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.FeaturedImg = os.Getenv("IMAGES") + contentId + "/" + seasonId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.Banner = os.Getenv("IMAGES") + contentId + "/" + seasonId + "/overlay-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	ImageryDetails <- Imagery
}
func EpisodeImagery(contentId, seasonId, episodeId string, ImageryDetails chan ContentImageryDetails) {
	var Imagery ContentImageryDetails
	Imagery.Thumbnail = os.Getenv("IMAGES") + contentId + "/" + seasonId + "/" + episodeId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.Backdrop = os.Getenv("IMAGES") + contentId + "/" + seasonId + "/details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.MobileImg = os.Getenv("IMAGES") + contentId + "/" + seasonId + "/mobile-details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.FeaturedImg = os.Getenv("IMAGES") + contentId + "/" + seasonId + "/" + episodeId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.Banner = os.Getenv("IMAGES") + contentId + "/" + seasonId + "/overlay-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	ImageryDetails <- Imagery
}
func CreateSliderResponse(sliderId, pageId string, c *gin.Context, response chan FragmentUpdate, countryId, platformId, ignoreAsync int) {
	var fragmentResponse FragmentUpdate
	if ignoreAsync == 0 {
		fragmentResponse.Response = ""
		fragmentResponse.Err = nil
		response <- fragmentResponse
		time.Sleep(1 * time.Second)
	}
	db := c.MustGet("DB").(*gorm.DB)
	fdb := c.MustGet("FDB").(*gorm.DB)
	var pagedetails []Page
	var featuredDetails FeaturedDetails
	var featuredPlaylist FeaturedPlaylists
	var countryDetails []Country
	var platformDetails []PublishPlatform
	if countryId != 0 {
		if err := db.Debug().Where("id=?", countryId).Find(&countryDetails).Error; err != nil {
			fragmentResponse.Response = ""
			fragmentResponse.Err = err
			response <- fragmentResponse
		}
	} else {
		if err := db.Debug().Find(&countryDetails).Error; err != nil {
			fragmentResponse.Response = ""
			fragmentResponse.Err = err
			response <- fragmentResponse
		}
	}

	if platformId != 0 {
		if err := db.Debug().Where("id=?", platformId).Find(&platformDetails).Error; err != nil {
			fragmentResponse.Response = ""
			fragmentResponse.Err = err
			response <- fragmentResponse
		}
	} else {
		if err := db.Debug().Find(&platformDetails).Error; err != nil {
			fragmentResponse.Response = ""
			fragmentResponse.Err = err
			response <- fragmentResponse
		}
	}

	var where string
	if pageId != "" {
		where = "s.deleted_by_user_id  is null and s.is_disabled =false and p.is_disabled=false and p.deleted_by_user_id is null and ps.slider_id='" + sliderId + "' and p.id='" + pageId + "'"
	} else {
		where = "s.deleted_by_user_id  is null and s.is_disabled =false and p.is_disabled=false and p.deleted_by_user_id is null and ps.slider_id='" + sliderId + "'"
	}
	if err := db.Debug().Table("page p").Select("distinct(p.id),p.page_key,p.page_type").Joins("inner join page_slider ps on ps.page_id=p.id inner join slider s on s.id = ps.slider_id ").Where(where).Find(&pagedetails).Error; err != nil {
		fragmentResponse.Response = ""
		fragmentResponse.Err = err
		response <- fragmentResponse
	}
	var slider Slider
	if err := db.Debug().Select("*").Where("deleted_by_user_id  is null and is_disabled =false and scheduling_start_date <=NOW() and scheduling_end_date >=NOW() and id=?", sliderId).Find(&slider).Error; err != nil && err.Error() != "record not found" {
		fragmentResponse.Response = ""
		fragmentResponse.Err = err
		response <- fragmentResponse
	}
	for _, country := range countryDetails {
		countryCode := common.CountryNames(country.Id)
		for _, platform := range platformDetails {
			platformName := common.DeviceNames(platform.Id)
			var sliderFragment []SliderFragment
			var sFragment SliderFragment
			fdb.Debug().Where("country=? and slider_id=? and platform=?", countryCode, sliderId, platformName).Delete(&sFragment)
			for _, details := range pagedetails {
				var blackPlaylistCount, redPlaylistCount, greenPlaylistCount int
				featuredDetails.Playlists = nil
				language := "en"
				for i := 1; i <= 2; i++ {
					if i == 2 {
						language = "ar"
					}
					if slider.SliderKey == 0 {
						break
					}
					var featuredPlaylists []FeaturedPlaylists
					featuredDetails.ID = slider.SliderKey
					featuredDetails.Type = common.SliderTypes(slider.Type)
					if slider.BlackAreaPlaylistId != "" || slider.RedAreaPlaylistId != "" || slider.GreenAreaPlaylistId != "" {
						playlists, _ := SliderPlaylists(slider.BlackAreaPlaylistId, slider.RedAreaPlaylistId, slider.GreenAreaPlaylistId, c)
						for _, playlist := range playlists {
							featuredPlaylist.ID = playlist.PlaylistKey
							featuredPlaylist.PlaylistType = playlist.PlaylistType
							contentIds, err := PlaylistItemContents(playlist.ID, c)
							if err != nil {
								fragmentResponse.Response = ""
								fragmentResponse.Err = err
								response <- fragmentResponse
							}
							playlistContent := make(chan []PlaylistContent)
							go PlaylistContentDetails(contentIds, language, country.Id, c, playlistContent)
							contentDetails := <-playlistContent
							if playlist.PlaylistType == "black_playlist" {
								blackPlaylistCount = len(contentDetails)
							} else if playlist.PlaylistType == "red_playlist" {
								redPlaylistCount = len(contentDetails)
							} else {
								greenPlaylistCount = len(contentDetails)
							}
							if contentDetails == nil {
								continue
							}
							var contents []PlaylistContent
							for _, id := range contentIds {
								for _, content := range contentDetails {
									if id.ContentId == content.ContentId {
										contents = append(contents, content)
									}
								}
							}
							featuredPlaylist.Content = contents
							featuredPlaylists = append(featuredPlaylists, featuredPlaylist)
						}
					}
					if details.PageType == 1 && blackPlaylistCount >= common.BLACK_AREA_PLAYLIST_CONUT && redPlaylistCount == common.RED_AREA_PLAYLIST_CONUT && greenPlaylistCount >= common.GREEN_AREA_PLAYLIST_CONUT {
						featuredDetails.Playlists = featuredPlaylists
					} else if details.PageType != 1 {
						featuredDetails.Playlists = featuredPlaylists
					}
					sliderFragmentdetails := make(chan SliderFragment)
					go SliderFragmentDetails(featuredDetails, details.PageKey, details.Id, countryCode, platformName, language, slider.Id, sliderFragmentdetails)
					sliderFragment = append(sliderFragment, <-sliderFragmentdetails)
				}
			}
			fmt.Printf("%v - %v ", country.Alpha2code, len(sliderFragment))
			fdb := c.MustGet("FDB").(*gorm.DB)
			var InserSliderFragments []interface{}
			for _, slider := range sliderFragment {
				InserSliderFragments = append(InserSliderFragments, slider)
			}
			err := gormbulk.BulkInsert(fdb, InserSliderFragments, common.BULK_INSERT_LIMIT)
			if err != nil {
				fragmentResponse.Response = ""
				fragmentResponse.Err = err
				response <- fragmentResponse
			}

		}
	}
	fragmentResponse.Response = "successfully updated."
	fragmentResponse.Err = nil
	response <- fragmentResponse
}
func SliderPlaylists(BlackAreaPlaylistId string, RedAreaPlaylistId string, GreenAreaPlaylistId string, c *gin.Context) ([]Playlist, error) {
	db := c.MustGet("DB").(*gorm.DB)
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
	if err := db.Debug().Select(fields).Where("id in(?) and (scheduling_start_date <=now() or scheduling_start_date is null) and (scheduling_end_date >=now() or scheduling_end_date is null)", playlist).Order("playlist_type desc").Find(&playlists).Error; err != nil {
		return nil, err
	}
	return playlists, nil
}
func SliderFragmentDetails(details FeaturedDetails, pageKey int, pageId, countryCode, platformName, languageCode, sliderId string, sliderFragmentDetails chan SliderFragment) {
	var sliderFragment SliderFragment
	sliderResponse := make(map[string]interface{})
	sliderResponse["slider_response"] = details
	sliderdata, _ := json.Marshal(sliderResponse)
	type SliderResponse struct {
		ResponseData postgres.Jsonb `json:"slider_response"`
	}
	var slider SliderResponse
	json.Unmarshal(sliderdata, &slider)
	sliderFragment.PageId = pageId
	sliderFragment.PageKey = pageKey
	sliderFragment.Country = countryCode
	sliderFragment.Platform = platformName
	sliderFragment.Details = slider.ResponseData
	sliderFragment.Language = languageCode
	sliderFragment.SliderId = sliderId
	sliderFragmentDetails <- sliderFragment
}
