package content

import (
	"encoding/json"
	"fmt"
	"frontend_service/common"
	l "frontend_service/logger"
	"frontend_service/menu"
	"math/rand"

	// "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"

	// "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/thanhpk/randstr"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes

	qrg := r.Group("/v1")
	qrg.Use(common.ValidateToken())
	qrg.GET("/:lang/contents/plans", hs.GetContentsByPlanId)
	qrg.GET("/contents/plans", hs.GetContentsByPlanId)
	// implemented for logs purpose
	lqrg := qrg.Group("/:lang")
	lqrg.Use(gin.LoggerWithFormatter(common.Logs))
	lqrg.GET("/series", hs.GetSeasonsByContent)
	qrg.GET("/:lang/contents/moviedetails", hs.GetMovieDetailsByContent)
	qrg.GET("/:lang/contents/moviedetailsresponse", hs.PrepareMovieDetailsByContent)
	qrg.GET("/contents/:ckeyctype", hs.GetContentRating)
	qrg.GET("/:lang/episode/:episode_key", hs.GetEpisodeDetailsByEpisodeKey)
	qrg.GET("/:lang/contents/contentType", hs.GetContentTypeDetails)
	qrg.GET("/:lang/contents/contentTrailer", hs.GetContentTrailerDetails)
	qrg.GET("/:lang/mediaobject/:ids", hs.GetMediaObjectDetails)
	qrg.GET("/:lang/related", hs.GetRelatedContents)
	qrg.GET("/contents/playlist", hs.GetUserPlaylists)
	qrg.GET("/contents/rated", hs.GetRatedContents)
	qrg.GET("/contents/resumable", hs.GetResumbleContents)
	qrg.GET("/contents/watching", hs.GetWatchingContents)
	qrg.POST("/contents/watching", hs.AddViewActivity)
	qrg.POST("/contents/watching/user", hs.AddViewActivity)
	qrg.GET("/:lang/search", hs.GetsearchContent)
	qrg.POST("/contents/playlist", hs.AddUserPlaylist)
	qrg.POST("/contents/playlist/user", hs.AddUserPlaylist)
	qrg.POST("/contents/rated", hs.AddRatingForContentByUser)
	qrg.POST("/contents/rated/user", hs.AddRatingForContentByUser)
	qrg.POST("/contents/watching/:ckeyctype/issues", hs.ReportContentIssue)
	qrg.DELETE("/contents/playlist/:ckeyctype", hs.RemoveContentsUserPlaylist)
	qrg.GET("/:lang/searchbyGenre", hs.GetSearchbyGenre)
	qrg.GET("/:lang/searchbyCast", hs.GetSearchbyCast)
	qrg.DELETE("/contents/watching/:ckeyctype", hs.RemoveWatchingHistory)
	qrg.DELETE("/contents/rated/:ckeyctype", hs.RemoveRatedContent)

	//For Flutter

	qrg.GET("/flutter/continuewatching", hs.GetResumbleContentsFlutterContinuewatching)

	//Error code exception URL
	qrg.POST("/:lang/searchbyGenre", hs.GetSearchbyGenre)
	qrg.PUT("/:lang/searchbyGenre", hs.GetSearchbyGenre)
	qrg.DELETE("/:lang/searchbyGenre", hs.GetSearchbyGenre)
	qrg.GET("/:lang/playlist/:playlistkey", hs.GetPlaylistDetails)
	qrg.GET("/contentFragments", hs.contentFragments)

}

// GetContentsByPlanId -  Get contents by plan
// GET /v1/:lang/contents/plans
// @Summary Get contents by plan
// @Description Get contents based on subscripti on plan id
// @Tags Content
// @Accept  json
// @Produce  json
// @Param lang path string true "Language Code"
// @Param Country query string false "Country Code"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Param plan query string false "Plan Id"
// @Success 200
// @Router /v1/{lang}/contents/plans [get]
func (hs *HandlerService) GetContentsByPlanId(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	notFoundError := common.NotFoundErrorResponse()
	serverError := common.ServerErrorResponse(language)
	var limit, offset, current_page int64
	var platform string
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		current_page, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	if limit == 0 {
		limit, _ = strconv.ParseInt(os.Getenv("DEFAULT_PAGE_SIZE"), 10, 64)
	}
	offset = current_page

	var CountryCode, PlanId string
	if c.Request.URL.Query()["Country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["Country"][0])
	}
	if c.Request.URL.Query()["plan"] != nil {
		PlanId = strings.ToUpper(c.Request.URL.Query()["plan"][0])
	}
	if c.Request.URL.Query()["platform"] != nil {
		platform = strings.ToUpper(c.Request.URL.Query()["platform"][0])
	}
	if platform == "" {
		platform = "0"
	}
	if PlanId == "" {
		PlanId = "3"
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}
	var Fragment string
	if c.Request.URL.Query()["fragment"] != nil {
		Fragment = c.Request.URL.Query()["fragment"][0]
	}
	country := int(common.Countrys(CountryCode))
	countryId := strconv.Itoa(int(country))
	var spcds []ContentIdDetails
	var contents []PlaylistContent
	query := common.ContentsByPlansQuery(countryId, PlanId)
	rows := db.Raw(query).Find(&spcds).RowsAffected
	if rows == 0 {
		l.JSON(c, http.StatusNotFound, notFoundError)
		return
	}
	if err := db.Raw(query).Limit(limit).Offset(offset).Find(&spcds).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	if Fragment == "" {
		ContentDetails := make(chan PlaylistContent)
		for _, spcd := range spcds {
			if spcd.ContentTier == 1 {
				go OneTierContentDetails(spcd.Id, language, country, c, ContentDetails)
			} else {
				go MultiTierContentDetails(spcd.Id, language, country, c, ContentDetails, 1, platform)
			}
			contents = append(contents, <-ContentDetails)
		}
	} else if Fragment != "" {
		var ids []string
		for _, spcd := range spcds {
			ids = append(ids, spcd.Id)
		}
		fdb := c.MustGet("FDB").(*gorm.DB)
		var contentFragment []ContentFragmentDetails
		if err := fdb.Table("content_fragment").Select("details::text").Where("content_id in(?) and language='en' and (scheduling_date_time <=NOW() or scheduling_date_time is null) and (digital_rights_start_date <=NOW() or digital_rights_start_date is null) and (digital_rights_end_date >=NOW() or digital_rights_end_date is null) and country=?", ids, CountryCode).Find(&contentFragment).Error; err == nil {
			for _, content := range contentFragment {
				var playlistContent PlaylistContent
				if err := json.Unmarshal([]byte(content.Details), &playlistContent); err != nil {
					l.JSON(c, http.StatusInternalServerError, serverError)
					return
				}
				contents = append(contents, playlistContent)
			}
		}
	}

	var pagination PaginationResult
	pagination.Size = rows
	pagination.Offset = offset
	pagination.Limit = limit
	l.JSON(c, http.StatusOK, gin.H{"pagination": pagination, "data": contents})
	return
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
		/*Based on frontend Developer below Requirements change from lower to original response
		1.0 API response they are sending the content_type as "LiveTV" ,but in 2.0 API response they are sending the content_type value as "livetv" , so it is not recognising the content as LiveTv, Please update once
		*/
		// playlistContent.ContentType = strings.ToLower(onetierContentResult.ContentType)
		if onetierContentResult.ContentType == "LiveTV" {
			playlistContent.ContentType = onetierContentResult.ContentType
		} else {
			playlistContent.ContentType = strings.ToLower(onetierContentResult.ContentType)
		}
		//playlistContent.ContentType = onetierContentResult.ContentType
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
			if err := db.Table("content_variance cv").Select(fields).Joins(join).Where(where, details.Id).Find(&SubsPlans).Error; err != nil {
				fmt.Println("err", err)
				return
			}
			plans := make([]int, 0)
			plansName := make([]string, 0)
			rights := make([]int, 0)
			for _, plan := range SubsPlans {
				type SubscriptionPlan struct {
					Id   string `json:"id"`
					Name string `json:"series_id"`
				}

				var planNameSubs SubscriptionPlan
				if err := db.Debug().Table("subscription_plan").Where("id = ?", plan.SubscriptionPlanId).Find(&planNameSubs).Error; err != nil {
					return
				}

				plansName = append(plansName, planNameSubs.Name)
				plans = append(plans, plan.SubscriptionPlanId)
			}
			movie.ID = 0
			movie.Title = details.Title
			movie.Geoblock = false
			movie.DigitalRightType = details.DigitalRightsType
			movie.DigitalRightsRegions = rights
			movie.SubscriptiontPlans = plans
			movie.SubscriptionPlansName = plansName
			movie.InsertedAt = details.InsertedAt
			if onetierContentResult.ContentType != "LiveTV" && onetierContentResult.ContentType != "livetv" && onetierContentResult.ContentType != "liveTV" {
				movie.IntroStart = details.IntroStart
				movie.IntroDuration = details.IntroDuration
			}
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
			err = db.Raw("select "+actorFields+" from actor where id =?", actorIds.MainActorId).Find(&mainActor).Error
			playlistContent.MainActor = mainActor.Name
			err = db.Raw("select "+actorFields+" from actor where id=?", actorIds.MainActressId).Find(&mainActress).Error
			if playlistContent.MainActor != mainActress.Name {
				playlistContent.MainActress = mainActress.Name
			}
			if err != nil {
				fmt.Println(err)
			}
		} else if actorIds.MainActorId != "" {
			err = db.Raw("select "+actorFields+" from actor where id=?", actorIds.MainActorId).Find(&mainActor).Error
			playlistContent.MainActor = mainActor.Name
			if err != nil {
				fmt.Println(err)
			}
		} else if actorIds.MainActressId != "" {
			err = db.Raw("select "+actorFields+" from actor where id=?", actorIds.MainActressId).Find(&mainActress).Error
			if playlistContent.MainActor != mainActress.Name {
				playlistContent.MainActress = mainActress.Name
			}
			if err != nil {
				fmt.Println(err)
			}
		}
		playlistContent.Cast = actorsList
		// for index, element := range playlistContent.Cast {
		// 	if element == playlistContent.MainActor {
		// 		playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
		// 	} else if element == playlistContent.MainActress {
		// 		// playlistContent.Cast = RemoveIndex(playlistContent.Cast, index)
		// 	}
		// }

		for index, element := range playlistContent.Cast {
			if element == playlistContent.MainActor {
				playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
			}
		}

		for index, element := range playlistContent.Cast {
			if element == playlistContent.MainActress {
				playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
			}
		}

		var genreNames []Names
		fields, join, where, groupBy = common.ContentGenresQuery(language)
		if err := db.Debug().Select(fields).Table("genre g").Joins(join).Where(where, contentId).Group(groupBy).Order("cg.order").Find(&genreNames).Error; err != nil {
			fmt.Println("err", err)
			return
		}
		for _, genre := range genreNames {
			genres = append(genres, genre.Name)
		}
		playlistContent.Genres = genres
		var tagNames []Names
		fields, join, where, groupBy = common.ContentTagsQuery()
		if err := db.Select(fields).Table("textual_data_tag tdt").Joins(join).Where(where, contentId).Group(groupBy).Order("tdt.name").Find(&tagNames).Error; err != nil {
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
func MultiTierContentDetails(contentId string, language string, country int, c *gin.Context, ContentDetails chan PlaylistContent, seasonLimit int, platform string) {
	fmt.Println("****inside go routine MultiTierContentDetails")
	db := c.MustGet("CDB").(*gorm.DB)
	var playlistContent PlaylistContent
	var onetierContentResult OnetierContentResult
	// var contentImageryDetails ContentImageryDetails
	fields, join, where, groupBy := common.MultitierContentQuery(contentId, language)
	fmt.Println("***inside go routine MultiTierContentDetails query")
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
			fmt.Println("***inside go routine MultiTierContentDetails SeasonDetailsQuery query")
			if err := db.Debug().Select(fields).Table("season s").Joins(join).Where(where, contentId, country).Group(groupBy).Limit(1).Order("s.number,s.created_at asc").Find(&seasonDetails).Error; err != nil {
				fmt.Println("err", err)
				return
			}
		} else {
			fmt.Println("***inside go routine MultiTierContentDetails SeasonDetailsQuery query 2")
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
				fmt.Println("***inside go routine MultiTierContentDetails ContentPlansQuery query")
				if err := db.Debug().Table("content_rights_plan crp").Select(fields).Joins(join).Where(where, details.ID).Find(&SubsPlans).Error; err != nil {
					fmt.Println("err", err)
					return
				}
				plans := make([]int, 0)
				plansName := make([]string, 0)
				// rights := make([]int, 0)
				for _, plan := range SubsPlans {
					type SubscriptionPlan struct {
						Id   string `json:"id"`
						Name string `json:"series_id"`
					}
					var planNameSubs SubscriptionPlan
					if err := db.Debug().Table("subscription_plan").Where("id = ?", plan.SubscriptionPlanId).Find(&planNameSubs).Error; err != nil {
						return
					}
					plansName = append(plansName, planNameSubs.Name)
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
				season.SubscriptionPlansName = plansName
				/* Skip Intro Implementation */
				season.IntroStart = details.IntroStart
				season.IntroDuration = details.IntroDuration
				season.OutroStart = details.OutroStart
				season.OutroDuration = details.OutroDuration
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

					fields, join, where, groupBy := common.SeasonEpisodesQuery(language, platform)
					fmt.Println("***inside go routine MultiTierContentDetails  SeasonEpisodesQuery query")
					if err := db.Debug().Select(fields).Table("episode e").
						Joins(join).Where(where, details.ID).
						Group(groupBy).Order("e.number asc").Find(&seasonEpisodes).Error; err == nil {
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

							if c.MustGet("AuthorizationRequired") == 1 {
								episode.LastWatchPosition = 0
							} else if c.MustGet("AuthorizationRequired") == 0 {
								UserId := c.MustGet("userid")
								var episodeviewActivity ViewActivity

								// episode view activity details based on last activity
								if errVA := db.Debug().Table("view_activity").Where("user_id=? and playback_item_id=?", UserId, episode.PlaybackItemId).Find(&episodeviewActivity).Error; errVA == nil {
									episode.LastWatchPosition = episodeviewActivity.LastWatchPosition
								}

								// if errVA := db.Debug().Table("view_activity_history").Where("user_id=? and content_key=?", UserId, episode.Id).Order("viewed_at desc").Limit(1).Find(&episodeviewActivity).Error; errVA == nil {
								// 	episode.LastWatchPosition = episodeviewActivity.LastWatchPosition
								// }

							}

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
			if playlistContent.MainActor != mainActress.Name {
				playlistContent.MainActress = mainActress.Name
			}
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
			if playlistContent.MainActor != mainActress.Name {
				playlistContent.MainActress = mainActress.Name
			}
			if err != nil {
				fmt.Println(err)
			}
		}
		playlistContent.Cast = actorsList

		for index, element := range playlistContent.Cast {
			if element == playlistContent.MainActor {
				playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
			}
		}

		for index, element := range playlistContent.Cast {
			if element == playlistContent.MainActress {
				playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
			}
		}

		var genreNames []Names
		fields, join, where, groupBy = common.ContentGenresQuery(language)
		// fields, join, where, groupBy = common.SeriesGenresQuery(language)
		if err := db.Debug().Select(fields).Table("genre g").Joins(join).Where(where, contentId).Group(groupBy).Order("cg.order").Find(&genreNames).Error; err != nil {
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
	fmt.Println("***end go routine MultiTierContentDetails query")
	return
}
func MultiTierContentDetailsWithoutEpisode(contentId string, language string, country int, c *gin.Context, ContentDetails chan PlaylistContentForWithoutEpisode, seasonLimit int) {
	fmt.Println("inside go routinte MultiTierContentDetailsWithoutEpisode")
	db := c.MustGet("CDB").(*gorm.DB)
	var playlistContent PlaylistContentForWithoutEpisode
	var onetierContentResult OnetierContentResult
	// var contentImageryDetails ContentImageryDetails
	fields, join, where, groupBy := common.MultitierContentQueryWithoutEpisode(contentId, language)
	fmt.Println("inside go routinte MultiTierContentDetailsWithoutEpisode and query MultitierContentQueryWithoutEpisode")
	if row := db.Debug().Table("content c").Select(fields).Joins(join).Where(where, contentId, country).Group(groupBy).Order("s.number").Limit(1).Find(&onetierContentResult).RowsAffected; row != 0 {
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
		var season PlaylistContentSeasonsForWithoutEpisode
		var seasons []PlaylistContentSeasonsForWithoutEpisode
		fields, join, where, groupBy := common.SeasonDetailsQueryWithoutEpisode(language)
		fmt.Println("inside go routinte MultiTierContentDetailsWithoutEpisode and query SeasonDetailsQueryWithoutEpisode")
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
				fmt.Println("inside go routinte MultiTierContentDetailsWithoutEpisode and query ContentPlansQuery")
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
				/* Skip Intro Implementation */
				season.IntroStart = details.IntroStart
				season.IntroDuration = details.IntroDuration
				season.OutroStart = details.OutroStart
				season.OutroDuration = details.OutroDuration
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
					//var seasonEpisodes []SeasonEpisodes
					//var episodesResponse []SeasonEpisodes
					/*fields, join, where, groupBy := common.SeasonEpisodesQuery(language)
					if err := db.Select(fields).Table("episode e").Joins(join).Where(where, details.ID).Group(groupBy).Order("e.number asc").Find(&seasonEpisodes).Error; err == nil {
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
					}*/
				}
				// var episodesResponse []SeasonEpisodes
				season.Episodes = make([]SeasonEpisodes, 0)
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
		fmt.Println("inside go routinte MultiTierContentDetailsWithoutEpisode and query ContentActorsQuery")
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
		fmt.Println("inside go routinte MultiTierContentDetailsWithoutEpisode and query Actors Query")
		if actorIds.MainActorId != "" && actorIds.MainActressId != "" {
			err = db.Debug().Raw("select "+actorFields+" from actor where id =?", actorIds.MainActorId).Find(&mainActor).Error
			playlistContent.MainActor = mainActor.Name
			err = db.Debug().Raw("select "+actorFields+" from actor where id=?", actorIds.MainActressId).Find(&mainActress).Error
			if playlistContent.MainActor != mainActress.Name {
				playlistContent.MainActress = mainActress.Name
			}
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
			if playlistContent.MainActor != mainActress.Name {
				playlistContent.MainActress = mainActress.Name
			}
			if err != nil {
				fmt.Println(err)
			}
		}
		playlistContent.Cast = actorsList
		// for index, element := range playlistContent.Cast {
		// 	if element == playlistContent.MainActor {
		// 		playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
		// 	} else if element == playlistContent.MainActress {
		// 		playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
		// 	}
		// }

		for index, element := range playlistContent.Cast {
			if element == playlistContent.MainActor {
				playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
			}
		}

		for index, element := range playlistContent.Cast {
			if element == playlistContent.MainActress {
				playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
			}
		}

		var genreNames []Names
		fmt.Println("inside go routinte MultiTierContentDetailsWithoutEpisode and query ContentGenresQuery")
		fields, join, where, groupBy = common.ContentGenresQuery(language)
		if err := db.Debug().Select(fields).Table("genre g").Joins(join).Where(where, contentId).Group(groupBy).Order("cg.order").Find(&genreNames).Error; err != nil {
			fmt.Println("err", err)
			return
		}
		for _, genre := range genreNames {
			genres = append(genres, genre.Name)
		}
		playlistContent.Genres = genres
		var tagNames []Names
		fmt.Println("inside go routinte MultiTierContentDetailsWithoutEpisode and query ContentTagsQuery")
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
func MultiTierContentDetailsWithoutEpisodeForSearch(contentId string, language string, country int, c *gin.Context, ContentDetails chan PlaylistContent, seasonLimit int) {
	db := c.MustGet("CDB").(*gorm.DB)
	var playlistContent PlaylistContent
	var onetierContentResult OnetierContentResult
	// var contentImageryDetails ContentImageryDetails
	fields, join, where, groupBy := common.MultitierContentQueryWithoutEpisode(contentId, language)
	if row := db.Table("content c").Select(fields).Joins(join).Where(where, contentId, country).Group(groupBy).Order("s.number").Limit(1).Find(&onetierContentResult).RowsAffected; row != 0 {
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
		fields, join, where, groupBy := common.SeasonDetailsQueryWithoutEpisode(language)
		if seasonLimit == 1 {
			if err := db.Select(fields).Table("season s").Joins(join).Where(where, contentId, country).Group(groupBy).Limit(1).Order("s.number,s.created_at asc").Find(&seasonDetails).Error; err != nil {
				fmt.Println("err", err)
				return
			}
		} else {
			if err := db.Select(fields).Table("season s").Joins(join).Where(where, contentId, country).Group(groupBy).Order("s.number,s.created_at asc").Find(&seasonDetails).Error; err != nil {
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
				if err := db.Table("content_rights_plan crp").Select(fields).Joins(join).Where(where, details.ID).Find(&SubsPlans).Error; err != nil {
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
				/* Skip Intro Implementation */
				season.IntroStart = details.IntroStart
				season.IntroDuration = details.IntroDuration
				season.OutroStart = details.OutroStart
				season.OutroDuration = details.OutroDuration
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
					//var seasonEpisodes []SeasonEpisodes
					//var episodesResponse []SeasonEpisodes
					/*fields, join, where, groupBy := common.SeasonEpisodesQuery(language)
					if err := db.Select(fields).Table("episode e").Joins(join).Where(where, details.ID).Group(groupBy).Order("e.number asc").Find(&seasonEpisodes).Error; err == nil {
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
					}*/
				}
				// var episodesResponse []SeasonEpisodes
				season.Episodes = make([]SeasonEpisodes, 0)
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
		if err := db.Select(fields).Table("content_cast cc").Joins(join).Where(where, castId).Group(groupBy).Find(&actorIds).Error; err == nil {
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
			err = db.Raw("select "+actorFields+" from actor where id =?", actorIds.MainActorId).Find(&mainActor).Error
			playlistContent.MainActor = mainActor.Name
			err = db.Raw("select "+actorFields+" from actor where id=?", actorIds.MainActressId).Find(&mainActress).Error
			if playlistContent.MainActor != mainActress.Name {
				playlistContent.MainActress = mainActress.Name
			}
			if err != nil {
				fmt.Println(err)
			}
		} else if actorIds.MainActorId != "" {
			err = db.Raw("select "+actorFields+" from actor where id=?", actorIds.MainActorId).Find(&mainActor).Error
			playlistContent.MainActor = mainActor.Name
			if err != nil {
				fmt.Println(err)
			}
		} else if actorIds.MainActressId != "" {
			err = db.Raw("select "+actorFields+" from actor where id=?", actorIds.MainActressId).Find(&mainActress).Error
			if playlistContent.MainActor != mainActress.Name {
				playlistContent.MainActress = mainActress.Name
			}
			if err != nil {
				fmt.Println(err)
			}
		}
		playlistContent.Cast = actorsList
		// for index, element := range playlistContent.Cast {
		// 	if element == playlistContent.MainActor {
		// 		playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
		// 	} else if element == playlistContent.MainActress {
		// 		playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
		// 	}
		// }

		for index, element := range playlistContent.Cast {
			if element == playlistContent.MainActor {
				playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
			}
		}

		for index, element := range playlistContent.Cast {
			if element == playlistContent.MainActress {
				playlistContent.Cast = common.RemoveIndex(playlistContent.Cast, index)
			}
		}

		var genreNames []Names
		fields, join, where, groupBy = common.ContentGenresQuery(language)
		if err := db.Select(fields).Table("genre g").Joins(join).Where(where, contentId).Group(groupBy).Order("cg.order").Find(&genreNames).Error; err != nil {
			fmt.Println("err", err)
			return
		}
		for _, genre := range genreNames {
			genres = append(genres, genre.Name)
		}
		playlistContent.Genres = genres
		var tagNames []Names
		fields, join, where, groupBy = common.ContentTagsQuery()
		if err := db.Select(fields).Table("textual_data_tag tdt").Joins(join).Where(where, contentId).Group(groupBy).Order("tdt.name").Find(&tagNames).Error; err != nil {
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

// GetSeasonsByContent -  Get seasons by content
// GET /v1/:lang/series
// @Summary Get seasons by content
// @Description Get seasons based on content key
// @Tags Content
// @Accept  json
// @Produce  json
// @Param platform query string false "Platform"
// @Param country query string false "Country Code"
// @Param contentkey query string true "Content Key"
// @Param lang path string true "Language Code"
// @Success 200
// @Router /v1/{lang}/series [get]
func (hs *HandlerService) GetSeasonsByContent(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	fmt.Println("inside series api before redis")
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	notFoundError := common.NotFoundErrorResponse()
	serverError := common.ServerErrorResponse(language)
	if c.Request.URL.Query()["contentkey"] == nil || c.Request.URL.Query()["contentkey"][0] == "" {
		serverError.Description = "Please Provide Content Key."
		l.JSON(c, http.StatusBadRequest, serverError)
		return
	}
	ContentKey := c.Request.URL.Query()["contentkey"][0]
	var CountryCode, platform string
	if c.Request.URL.Query()["country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
	}
	if c.Request.URL.Query()["platform"] != nil {
		platform = strings.ToUpper(c.Request.URL.Query()["platform"][0])
	}
	if platform == "" {
		platform = "0"
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}
	country := int(common.Countrys(CountryCode))
	var content Content
	result := PlaylistContent{}
	/* get Result From Redis */
	key := os.Getenv("REDIS_CONTENT_KEY") + ContentKey + language + strconv.Itoa(country)
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "?", "_")
	url := os.Getenv("REDIS_CACHE_URL") + "/" + key
	response, err := common.GetCurlCall(url)
	if err != nil {
		fmt.Println(err)
		l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	type RedisCacheResponse struct {
		Value string `json:"value"`
		Error string `json:"error"`
	}
	var RedisResponse RedisCacheResponse
	json.Unmarshal(response, &RedisResponse)
	// if RedisResponse.Value != "" {
	if "" != "" {
		fmt.Println("inside series api getting from redis")
		if err := json.Unmarshal([]byte(RedisResponse.Value), &result); err != nil {
			fmt.Println(err)
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	} else {
		fmt.Println("inside series api getting from DB")
		if err := db.Debug().Where("content_key=? and content_tier=2", ContentKey).Find(&content).Error; err != nil {
			l.JSON(c, http.StatusBadRequest, serverError)
			return
		}

		ContentDetails := make(chan PlaylistContent)
		fmt.Println("inside series api getting from DB entering go routine")
		go MultiTierContentDetails(content.Id, language, country, c, ContentDetails, 0, platform)
		details := <-ContentDetails
		if details.ID == 0 {
			ContentDetails := make(chan PlaylistContentForWithoutEpisode)
			fmt.Println("inside series api getting entering go routine MultiTierContentDetailsWithoutEpisode")
			go MultiTierContentDetailsWithoutEpisode(content.Id, language, country, c, ContentDetails, 0)
			details1 := <-ContentDetails
			if details1.ID == 0 && details.ID == 0 {
				l.JSON(c, http.StatusNotFound, notFoundError)
				return
			} else {
				/* Create Redis Key for content Type */
				jsonData, _ := json.Marshal(details1)
				var request RedisCacheRequest
				url := os.Getenv("REDIS_CACHE_URL")
				request.Key = key //os.Getenv("REDIS_CONTENT_KEY") + ContentKey + language + strconv.Itoa(country)
				request.Value = string(jsonData)
				_, err := common.PostCurlCall("POST", url, request)
				if err != nil {
					fmt.Println(err)
					l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
					return
				}
				l.JSON(c, http.StatusOK, gin.H{"data": details1})
				return
			}
		}
		/* Create Redis Key for content Type */
		jsonData, _ := json.Marshal(details)
		var request RedisCacheRequest
		url := os.Getenv("REDIS_CACHE_URL")
		request.Key = key //pageKey + language + strconv.Itoa(country) + strconv.Itoa(platform)
		request.Value = string(jsonData)
		_, err := common.PostCurlCall("POST", url, request)
		if err != nil {
			fmt.Println(err)
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		l.JSON(c, http.StatusOK, gin.H{"data": details})
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": result})
	return
}

/* without redis
func (hs *HandlerService) GetSeasonsByContent(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	notFoundError := common.NotFoundErrorResponse()
	serverError := common.ServerErrorResponse(language)
	if c.Request.URL.Query()["contentkey"] == nil || c.Request.URL.Query()["contentkey"][0] == "" {
		serverError.Description = "Please Provide Content Key."
		l.JSON(c, http.StatusBadRequest, serverError)
		return
	}
	ContentKey := c.Request.URL.Query()["contentkey"][0]
	var CountryCode string
	if c.Request.URL.Query()["country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}
	country := int(common.Countrys(CountryCode))
	var content Content
	if err := db.Where("content_key=? and content_tier=2", ContentKey).Find(&content).Error; err != nil {
		l.JSON(c, http.StatusBadRequest, serverError)
		return
	}
	ContentDetails := make(chan PlaylistContent)
	go MultiTierContentDetails(content.Id, language, country, c, ContentDetails, 0)
	details := <-ContentDetails
	if details.ID == 0 {
		ContentDetails := make(chan PlaylistContentForWithoutEpisode)
		go MultiTierContentDetailsWithoutEpisode(content.Id, language, country, c, ContentDetails, 0)
		details1 := <-ContentDetails
		if details1.ID == 0 && details.ID == 0 {
			l.JSON(c, http.StatusNotFound, notFoundError)
			return
		} else {
			l.JSON(c, http.StatusOK, gin.H{"data": details1})
			return
		}
	}
	l.JSON(c, http.StatusOK, gin.H{"data": details})
	return
}*/

// GetMovieDetailsByContent -  Get Movie details by content
// GET /v1/:lang/contents/moviedetails
// @Summary Get movie details by content
// @Description Get movie details based on content key
// @Tags Content
// @Accept  json
// @Produce  json
// @Param Country query string false "Country Code"
// @Param contentkey query string true "Content Key"
// @Param lang path string true "Language Code"
// @Success 200
// @Router /v1/{lang}/contents/moviedetails [get]
func (hs *HandlerService) GetMovieDetailsByContent(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	//notFoundError := common.NotFoundErrorResponse()
	serverError := common.ServerErrorResponse(language)
	if c.Request.URL.Query()["contentkey"] == nil || c.Request.URL.Query()["contentkey"][0] == "" {
		serverError.Description = "Please Provide Content Key."
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	ContentKey := c.Request.URL.Query()["contentkey"][0]
	var CountryCode string
	if c.Request.URL.Query()["Country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["Country"][0])
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}

	country := int(common.Countrys(CountryCode))
	var content Content
	result := PlaylistContent{}
	/* get Result From Redis */
	key := os.Getenv("REDIS_CONTENT_KEY") + ContentKey + language + strconv.Itoa(country)
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "?", "_")
	url := os.Getenv("REDIS_CACHE_URL") + "/" + key
	response, err := common.GetCurlCall(url)
	if err != nil {
		fmt.Println(err)
		l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	type RedisCacheResponse struct {
		Value string `json:"value"`
		Error string `json:"error"`
	}
	var RedisResponse RedisCacheResponse
	json.Unmarshal(response, &RedisResponse)
	// TODO:
	// if "" != "" {
	if RedisResponse.Value != "" {
		if err := json.Unmarshal([]byte(RedisResponse.Value), &result); err != nil {
			fmt.Println(err)
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	} else {
		if err := db.Where("content_key=?", ContentKey).Find(&content).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		ContentDetails := make(chan PlaylistContent)
		go OneTierContentDetails(content.Id, language, country, c, ContentDetails)
		details := <-ContentDetails
		fmt.Println("details not found")
		if details.ID == 0 {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* Create Redis Key for content Type */
		jsonData, _ := json.Marshal(details)
		var request RedisCacheRequest
		url := os.Getenv("REDIS_CACHE_URL")
		request.Key = key //os.Getenv("REDIS_CONTENT_KEY") + ContentKey + language + strconv.Itoa(country)
		request.Value = string(jsonData)
		_, err := common.PostCurlCall("POST", url, request)
		if err != nil {
			fmt.Println(err)
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		l.JSON(c, http.StatusOK, gin.H{"data": details})
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": result})
	return
}

/*with out redis
func (hs *HandlerService) GetMovieDetailsByContent(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	//notFoundError := common.NotFoundErrorResponse()
	serverError := common.ServerErrorResponse(language)
	if c.Request.URL.Query()["contentkey"] == nil || c.Request.URL.Query()["contentkey"][0] == "" {
		serverError.Description = "Please Provide Content Key."
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	ContentKey := c.Request.URL.Query()["contentkey"][0]
	var CountryCode string
	if c.Request.URL.Query()["Country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["Country"][0])
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}

	country := int(common.Countrys(CountryCode))
	var content Content
	if err := db.Where("content_key=? and content_tier=1", ContentKey).Find(&content).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	ContentDetails := make(chan PlaylistContent)
	go OneTierContentDetails(content.Id, language, country, c, ContentDetails)
	details := <-ContentDetails
	fmt.Println("details not found")
	if details.ID == 0 {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": details})
	return
}*/

// GetContentRating -  Get content rating details
// GET /v1/contents/:ckeyctype
// @Summary Get content rating details
// @Description Get content rating details based on content key
// @Tags Content
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param ckeyctype path string true "ContentKey,ContentType"
// @Success 200
// @Router /v1/contents/{ckeyctype} [get]
func (hs *HandlerService) GetContentRating(c *gin.Context) {
	// db := c.MustGet("CDB").(*gorm.DB)
	UserId := c.MustGet("userid")

	if c.Param("ckeyctype") == "" {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": "Please Provide Content Key and Content Type.", "status": http.StatusBadRequest})
		return
	}
	request := strings.Split(c.Param("ckeyctype"), ",")
	var ContentKey, ContentType string
	if len(request) < 2 {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": "Please Provide valid Content Key and Content Type.", "status": http.StatusBadRequest})
		return
	} else {
		ContentKey = request[0]
		ContentType = request[1]
	}
	fmt.Println("sssss", request[1])
	if ContentType == "undefined" {
		ContentType = "episode"
	}
	if strings.ToLower(ContentType) == "episode" {
		details, err := GetEpisodeRatingDetails(ContentKey, UserId.(string), 1, 0, c)
		if err != nil {
			l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
			return
		}
		l.JSON(c, http.StatusOK, gin.H{"data": details})
	} else {
		details, err := GetContentRatingDetails(ContentKey, UserId.(string), 1, 0, c)
		if err != nil {
			l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
			return
		}
		l.JSON(c, http.StatusOK, gin.H{"data": details})
	}

	return
}

// GetEpisodeDetailsByEpisodeKey -  Fetch details of episode based on episode key
// GET /v1/:lang/episode/:episode_key
// @Summary Get episode details based on episode key
// @Description Get episode details based on episode key
// @Tags Content
// @Accept  json
// @Produce  json
// @Param episode_key path string true "Episode Key"
// @Param country query string false "Country Code"
// @Success 200
// @Router /v1/{lang}/episode/{episode_key} [get]
func (hs *HandlerService) GetEpisodeDetailsByEpisodeKey(c *gin.Context) {
	var response EpisodeResponse
	var genreResponse []GenreNameResponse
	var tagResponse []TagResponse
	var subscriptionplans []SubscriptionPlan
	if c.Request.URL.Query()["country"] != nil {
		response.Geoblock = false
	} else {
		response.Geoblock = true
	}
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	serverError := common.ServerErrorResponse(language)
	var fields string
	if language == "en" {
		fields = ",cpi.transliterated_title as title, cpi.arabic_title as translated_title"
	} else {
		fields = ",cpi.arabic_title as title, cpi.transliterated_title as translated_title"
	}
	episodeKey := c.Param("episode_key")
	db := c.MustGet("CDB").(*gorm.DB)
	if err := db.Table("episode as e").
		Select("e.episode_key as id,c.content_key as series_id,s.number as season_number,season_key as season_id,e.number as episode_number,video_content_id as video_id,pi.duration as length,e.created_at,cpi.intro_start, cpi.outro_start, cr.digital_rights_type,c.id as contentid,s.id as seasonid, e.id as episodeid, e.tag_info_id as tagid,s.rights_id as rightsid"+fields+"").
		Joins("left join season s on s.id=e.season_id").
		Joins("left join content_primary_info cpi on cpi.id=e.primary_info_id").
		Joins("left join content c on c.id=s.content_id").
		Joins("left join playback_item pi on pi.id=e.playback_item_id").
		Joins("left join content_rights cr on cr.id= pi.rights_id").
		Where("e.episode_key = ?", episodeKey).
		Find(&response).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	gFields := "g.english_name as genre"
	if language == "ar" {
		gFields = "g.arabic_name as genre"
	}
	if err := db.Table("content_genre cg").
		Select(gFields).
		Joins("left join genre g on g.id=cg.genre_id").
		Where("content_id = ?", response.Contentid).
		Find(&genreResponse).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	for _, s := range genreResponse {
		response.Genres = append(response.Genres, s.Genre)
	}

	if err := db.Table("content_tag ct").
		Select("tdt.name as tags").
		Joins("left join textual_data_tag tdt on tdt.id=ct.textual_data_tag_id").
		Where("ct.tag_info_id = ?", response.Tagid).
		Find(&tagResponse).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	for _, t := range tagResponse {
		response.Tags = append(response.Tags, t.Tags)
	}

	if err := db.Table("content_rights cr").
		Select("sp.name as subscription_plan").
		Joins("left join content_rights_plan crp on crp.rights_id=cr.id").
		Joins("left join subscription_plan sp on sp.id::int=crp.subscription_plan_id::int").
		Where("cr.id = ?", response.Rightsid).
		Find(&subscriptionplans).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	for _, i := range subscriptionplans {
		response.SubscriptionPlans = append(response.SubscriptionPlans, i.Name)
	}

	//Forming of imagery dict
	ImageryDetails := make(chan ContentImageryDetails)
	go EpisodeImagery(response.Contentid, response.Seasonid, response.Episodeid, ImageryDetails)
	response.Imagery = <-ImageryDetails
	l.JSON(c, http.StatusOK, gin.H{"data": response})
}

type RedisCacheRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetContentTypeDetails -  Get contents by type
// GET /v1/:lang/contents/contentType
// @Summary Get contents by type
// @Description Get contents details based on content type
// @Tags Content
// @Accept  json
// @Produce  json
// @Param lang path string true "Language Code"
// @Param Country query string false "Country Code"
// @Param RowCountPerPage query string false "Row Count Per Page"
// @Param pageNo query string false "Page No"
// @Param IsPaging query string false "Is Paging"
// @Param OrderBy query string false "Order By"
// @Param contentType query string true "Content Type"
// @Param platform query string true "Platform"
// @Success 200
// @Router /v1/{lang}/contents/contentType [get]
func (hs *HandlerService) GetContentTypeDetails(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	if c.Request.URL.Query()["contentType"] == nil || c.Request.URL.Query()["contentType"][0] == "" {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": "Please provide content type."})
		return
	}
	contentType := strings.ToLower(c.Request.URL.Query()["contentType"][0])
	var limit, offset, current_page int64
	var platform string
	if c.Request.URL.Query()["RowCountPerPage"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["RowCountPerPage"][0], 10, 64)
	}
	if c.Request.URL.Query()["pageNo"] != nil {
		current_page, _ = strconv.ParseInt(c.Request.URL.Query()["pageNo"][0], 10, 64)
	}
	if c.Request.URL.Query()["platform"] != nil {
		platform = c.Request.URL.Query()["platform"][0]
	}
	if platform == "" {
		platform = "0"
	}
	if limit == 0 {
		limit, _ = strconv.ParseInt(os.Getenv("DEFAULT_PAGE_SIZE"), 10, 64)
	}
	if current_page == 0 || current_page == 1 {
		offset = 0
	} else {
		offset = (current_page * limit) - limit
	}

	var CountryCode, OrderBy string
	var IsPaging int64
	if c.Request.URL.Query()["Country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["Country"][0])
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}
	if c.Request.URL.Query()["IsPaging"] != nil {
		IsPaging, _ = strconv.ParseInt(c.Request.URL.Query()["IsPaging"][0], 10, 64)
	}
	if c.Request.URL.Query()["OrderBy"] != nil {
		OrderBy = strings.ToUpper(c.Request.URL.Query()["OrderBy"][0])
	}
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	country := int(common.Countrys(CountryCode))
	var contents []AllAvailableSeasons
	result := []AllAvailableSeasons{}

	/* get Result From Redis */
	key := os.Getenv("REDIS_CONTENT_KEY") + contentType + language + strconv.Itoa(country)
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "?", "_")
	url := os.Getenv("REDIS_CACHE_URL") + "/" + key
	response, err := common.GetCurlCall(url)
	// if err != nil {
	// 	fmt.Println(err)
	// 	l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
	// 	return
	// }
	type RedisCacheResponse struct {
		Value string `json:"value"`
		Error string `json:"error"`
	}
	var RedisResponse RedisCacheResponse
	json.Unmarshal(response, &RedisResponse)
	// to make redis disabled -- To Do Lter we need to enhance on this
	//RedisResponse.Value = ""

	fmt.Println("key-------->", key)
	if RedisResponse.Value != "" {
		if err := json.Unmarshal([]byte(RedisResponse.Value), &result); err != nil {
			fmt.Println(err)
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	} else {
		if contentType == "series" || contentType == "program" {
			//fields, join, where, groupBy := common.GetSeriesQuery(language)
			if IsPaging == 1 {
				// if err := db.Select(fields).Table("content_primary_info cpi").Joins(join).Where(where, country, contentType).Group(groupBy).Order("c.created_at " + OrderBy + "").Limit(limit).Offset(offset).Find(&contents).Error; err != nil {
				// 	l.JSON(c, http.StatusInternalServerError, gin.H{"error": err})
				// 	return
				// }
				if err := db.Raw("select c.id as content_id,min(s.id::text) as season_id,c.content_key as id,c.content_tier,min(pi1.video_content_id) as video_id,replace(lower(cpi.transliterated_title), ' ', '_') as friendly_url,cpi.transliterated_title as title,c.created_at as time from content_primary_info cpi join content c on c.primary_info_id = cpi.id join season s on s.content_id = c.id join episode e on e.season_id = s.id join playback_item pi1 on pi1.id = e.playback_item_id join content_rights cr on cr.id = s.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where(( pi1.scheduling_date_time <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and c.status = 1 and c.deleted_by_user_id is null and s.status = 1 and s.deleted_by_user_id is null and (e.status = 1 or e.status is null) and e.deleted_by_user_id is null and crc.country_id = ? and lower(c.content_type) = ? and ? in (select pitp.target_platform from playback_item_target_platform pitp where pitp.playback_item_id = e.playback_item_id)) group by c.id,c.content_key,cpi.transliterated_title,cpi.arabic_title,c.created_at union select c.id as content_id,min(s.id::text) as season_id,c.content_key as id,c.content_tier,min(vt.video_trailer_id) as video_id,replace(lower(cpi.transliterated_title), ' ', '_') as friendly_url,cpi.transliterated_title as title,c.created_at as time from content_primary_info cpi left join content c on c.primary_info_id = cpi.id left join season s on s.content_id = c.id left join episode e on e.season_id = s.id  left join variance_trailer vt on vt.season_id = s.id left join content_rights cr on cr.id = s.rights_id left join content_rights_country crc on crc.content_rights_id = cr.id where((cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null)and c.status = 1 and c.deleted_by_user_id is null and s.status = 1 and s.deleted_by_user_id is null and vt.id is not null and crc.country_id = ? and lower(c.content_type) = ? and ? in (select pitp.target_platform from playback_item_target_platform pitp where pitp.playback_item_id = e.playback_item_id) and e.id is null) group by c.id,c.content_key,cpi.transliterated_title,cpi.arabic_title,c.created_at order by time "+OrderBy+" ", country, contentType, platform, country, contentType, platform).Limit(limit).Offset(offset).Find(&contents).Error; err != nil {
					l.JSON(c, http.StatusInternalServerError, gin.H{"error": err})
					return
				}

			} else {
				/*if err := db.Select(fields).Table("content_primary_info cpi").Joins(join).Where(where, country, contentType).Group(groupBy).Order("c.created_at " + OrderBy + "").Find(&contents).Error; err != nil {
					l.JSON(c, http.StatusInternalServerError, gin.H{"error": err})
					return
				}*/
				// if err := db.Raw("select c.id as content_id,min(s.id::text) as season_id,c.content_key as id,c.content_tier,min(pi1.video_content_id) as video_id,replace(lower(cpi.transliterated_title), ' ', '_') as friendly_url,cpi.transliterated_title as title,c.created_at as time from content_primary_info cpi join content c on c.primary_info_id = cpi.id join season s on s.content_id = c.id join episode e on e.season_id = s.id join playback_item pi1 on pi1.id = e.playback_item_id join content_rights cr on cr.id = s.rights_id join content_rights_country crc on crc.content_rights_id = cr.id left join variance_trailer vt on vt.season_id = s.id where(( pi1.scheduling_date_time <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and c.status = 1 and c.deleted_by_user_id is null and s.status = 1 and s.deleted_by_user_id is null and (e.status = 1 or e.status is null) and crc.country_id = ? and lower(c.content_type) = ? and ( ( ? in ( select pitp2.target_platform from playback_item_target_platform pitp2 where playback_item_id = pi1.id) and e.id is not null ) or  (e.id is null and vt.id is not null) ) ) group by c.id,c.content_key,cpi.transliterated_title,cpi.arabic_title,c.created_at union select c.id as content_id,min(s.id::text) as season_id,c.content_key as id,c.content_tier,min(vt.video_trailer_id) as video_id,replace(lower(cpi.transliterated_title), ' ', '_') as friendly_url,cpi.transliterated_title as title,c.created_at as time from content_primary_info cpi left join content c on c.primary_info_id = cpi.id left join season s on s.content_id = c.id left join episode e on e.season_id = s.id  left join variance_trailer vt on vt.season_id = s.id left join content_rights cr on cr.id = s.rights_id left join content_rights_country crc on crc.content_rights_id = cr.id where((cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null)and c.status = 1 and c.deleted_by_user_id is null and s.status = 1 and s.deleted_by_user_id is null and vt.id is not null and crc.country_id = ? and lower(c.content_type) = ? and ? in (select pitp.target_platform from playback_item_target_platform pitp where pitp.playback_item_id = e.playback_item_id) and e.id is null or (e.id is null and vt.id is not null)) group by c.id,c.content_key,cpi.transliterated_title,cpi.arabic_title,c.created_at order by time "+OrderBy+" ", country, contentType, platform, country, contentType, platform).Find(&contents).Error; err != nil {
				if err := db.Raw("select c.id as content_id,min(s.id::text) as season_id,c.content_key as id,c.content_tier,min(pi1.video_content_id) as video_id,replace(lower(cpi.transliterated_title), ' ', '_') as friendly_url,cpi.transliterated_title as title,c.created_at as time from content_primary_info cpi join content c on c.primary_info_id = cpi.id join season s on s.content_id = c.id join episode e on e.season_id = s.id join playback_item pi1 on pi1.id = e.playback_item_id join content_rights cr on cr.id = s.rights_id join content_rights_country crc on crc.content_rights_id = cr.id left join variance_trailer vt on vt.season_id = s.id where(( pi1.scheduling_date_time <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and c.status = 1 and c.deleted_by_user_id is null and s.status = 1 and s.deleted_by_user_id is null and (e.status = 1 or e.status is null) and crc.country_id = ? and lower(c.content_type) = ?  ) group by c.id,c.content_key,cpi.transliterated_title,cpi.arabic_title,c.created_at union select c.id as content_id,min(s.id::text) as season_id,c.content_key as id,c.content_tier,min(vt.video_trailer_id) as video_id,replace(lower(cpi.transliterated_title), ' ', '_') as friendly_url,cpi.transliterated_title as title,c.created_at as time from content_primary_info cpi left join content c on c.primary_info_id = cpi.id left join season s on s.content_id = c.id left join episode e on e.season_id = s.id  left join variance_trailer vt on vt.season_id = s.id left join content_rights cr on cr.id = s.rights_id left join content_rights_country crc on crc.content_rights_id = cr.id where((cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null)and c.status = 1 and c.deleted_by_user_id is null and s.status = 1 and s.deleted_by_user_id is null and vt.id is not null and crc.country_id = ? and lower(c.content_type) = ? and ? in (select pitp.target_platform from playback_item_target_platform pitp where pitp.playback_item_id = e.playback_item_id) and e.id is null or (e.id is null and vt.id is not null)) and c.deleted_by_user_id is null group by c.id,c.content_key,cpi.transliterated_title,cpi.arabic_title,c.created_at order by time "+OrderBy+" ", country, contentType, country, contentType, platform).Find(&contents).Error; err != nil {
					// if err := db.Raw(`
					// 		select
					// 			c.id as content_id,
					// 			min(s.id::text) as season_id,
					// 			c.content_key as id,
					// 			c.content_tier,
					// 			min(pi1.video_content_id) as video_id,
					// 			replace(lower(cpi.transliterated_title), ' ', '_') as friendly_url,
					// 			cpi.transliterated_title as title,
					// 			c.created_at as time
					// 		from
					// 			content_primary_info cpi
					// 		join content c on
					// 			c.primary_info_id = cpi.id
					// 		join season s on
					// 			s.content_id = c.id
					// 		join episode e on
					// 			e.season_id = s.id
					// 		join playback_item pi1 on
					// 			pi1.id = e.playback_item_id
					// 		join content_rights cr on
					// 			cr.id = s.rights_id
					// 		join content_rights_country crc on
					// 			crc.content_rights_id = cr.id
					// 		left join variance_trailer vt on
					// 			vt.season_id = s.id
					// 		where
					// 			(( pi1.scheduling_date_time <= NOW()
					// 				or pi1.scheduling_date_time is null)
					// 			and (cr.digital_rights_start_date <= NOW()
					// 				or cr.digital_rights_start_date is null)
					// 			and (cr.digital_rights_end_date >= NOW()
					// 				or cr.digital_rights_end_date is null)
					// 			and c.status = 1
					// 			and c.deleted_by_user_id is null
					// 			and s.status = 1
					// 			and s.deleted_by_user_id is null
					// 			and (e.status = 1
					// 				or e.status is null)
					// 			and crc.country_id = ?
					// 			and lower(c.content_type) = ?
					// 				and ( ( ? in (
					// 				select
					// 					pitp2.target_platform
					// 				from
					// 					playback_item_target_platform pitp2
					// 				where
					// 					playback_item_id = pi1.id)
					// 				and e.id is not null )
					// 				or (e.id is null
					// 					and vt.id is not null) ) )
					// 		group by
					// 			c.id,
					// 			c.content_key,
					// 			cpi.transliterated_title,
					// 			cpi.arabic_title,
					// 			c.created_at
					// 		union
					// 		select
					// 			c.id as content_id,
					// 			min(s.id::text) as season_id,
					// 			c.content_key as id,
					// 			c.content_tier,
					// 			min(vt.video_trailer_id) as video_id,
					// 			replace(lower(cpi.transliterated_title), ' ', '_') as friendly_url,
					// 			cpi.transliterated_title as title,
					// 			c.created_at as time
					// 		from
					// 			content_primary_info cpi
					// 		left join content c on
					// 			c.primary_info_id = cpi.id
					// 		left join season s on
					// 			s.content_id = c.id
					// 		left join episode e on
					// 			e.season_id = s.id
					// 		left join variance_trailer vt on
					// 			vt.season_id = s.id
					// 		left join content_rights cr on
					// 			cr.id = s.rights_id
					// 		left join content_rights_country crc on
					// 			crc.content_rights_id = cr.id
					// 		left join playback_item pi1 on
					// 			pi1.id = e.playback_item_id
					// 		where
					// 			((cr.digital_rights_start_date <= NOW()
					// 				or cr.digital_rights_start_date is null)
					// 			and (cr.digital_rights_end_date >= NOW()
					// 				or cr.digital_rights_end_date is null)
					// 			and c.status = 1
					// 			and c.deleted_by_user_id is null
					// 			and s.status = 1
					// 			and s.deleted_by_user_id is null
					// 			and vt.id is not null
					// 			and crc.country_id = ?
					// 			and lower(c.content_type) = ?

					// 			and ( ( ? in (
					// 				select
					// 					pitp2.target_platform
					// 				from
					// 					playback_item_target_platform pitp2
					// 				where
					// 					playback_item_id = pi1.id)
					// 				and e.id is not null )
					// 				or (e.id is null
					// 					and vt.id is not null) ) )

					// 		group by
					// 			c.id,
					// 			c.content_key,
					// 			cpi.transliterated_title,
					// 			cpi.arabic_title,
					// 			c.created_at
					// 		order by
					// 			time `+OrderBy+` `, country, contentType, platform, country, contentType, platform).Find(&contents).Error; err != nil {
					l.JSON(c, http.StatusInternalServerError, gin.H{"error": err})
					return
				}
			}

		} else {
			fields, join, where, groupBy := common.GetMoviesQuery(language)
			if IsPaging == 1 {
				if err := db.Debug().Select(fields).Table("content c").Joins(join).Where(where, country, contentType, platform).Group(groupBy).Order("c.created_at  " + OrderBy + "").Limit(limit).Offset(offset).Find(&contents).Error; err != nil {
					l.JSON(c, http.StatusInternalServerError, gin.H{"error": err})
					return
				}
			} else {
				if err := db.Debug().Select(fields).Table("content c").Joins(join).Where(where, country, contentType, platform).Group(groupBy).Order("c.created_at  " + OrderBy + "").Find(&contents).Error; err != nil {
					l.JSON(c, http.StatusInternalServerError, gin.H{"error": err})
					return
				}
			}
		}
		for _, content := range contents {
			ImageryDetails := make(chan ContentImageryDetails)
			if content.ContentTier == 1 {
				go OnetierImagery(content.ContentId, content.SeasonId, ImageryDetails)
			} else {
				go MultitierImagery(content.ContentId, content.SeasonId, ImageryDetails)
			}
			// fmt.Println(content.SeasonId, "season id is")
			content.Imagery = <-ImageryDetails
			result = append(result, content)
		}
		/* Create Redis Key for content Type */
		jsonData, _ := json.Marshal(result)
		var request RedisCacheRequest
		url = os.Getenv("REDIS_CACHE_URL")
		request.Key = key //pageKey + language + strconv.Itoa(country) + strconv.Itoa(platform)
		request.Value = string(jsonData)
		_, err = common.PostCurlCall("POST", url, request)
		if err != nil {
			fmt.Println(err)
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	}
	l.JSON(c, http.StatusOK, result)
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

// GetContentTrailerDetails -  Get content trailers
// GET /v1/:lang/contents/contentTrailer
// @Summary Get content trailer details
// @Description Get content trailers based on content key
// @Tags Content
// @Accept  json
// @Produce  json
// @Param Country	 query string false "Country Code"
// @Param contenttype query string true "Content Type"
// @Param contentkey query int true "Content Key"
// @Param lang path string true "Language Code"
// @Success 200
// @Router /v1/{lang}/contents/contentTrailer [get]
func (hs *HandlerService) GetContentTrailerDetails(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	notFoundError := common.NotFoundErrorResponse()
	if c.Request.URL.Query()["contentkey"] == nil || c.Request.URL.Query()["contentkey"][0] == "" {
		notFoundError.Description = "Invalid Content key"
		l.JSON(c, http.StatusBadRequest, notFoundError)
		return
	}
	if c.Request.URL.Query()["contenttype"] == nil || c.Request.URL.Query()["contenttype"][0] == "" {
		notFoundError.Description = "Invalid Content type"
		l.JSON(c, http.StatusBadRequest, notFoundError)
		return
	}
	ContentKey := c.Request.URL.Query()["contentkey"][0]
	ContentType := strings.ToLower(c.Request.URL.Query()["contenttype"][0])
	var CountryCode string
	if c.Request.URL.Query()["Country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["Country"][0])
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	serverError := common.ServerErrorResponse(language)
	country := int(common.Countrys(CountryCode))
	var trailers []ContentTrailers
	result := []ContentTrailers{}
	if ContentType == "movie" || ContentType == "play" || ContentType == "LiveTV" || ContentType == "livetv" {
		fmt.Println("CONTENT TYPE ", ContentType)
		fields, join, where := common.GetMovieTrailerQuery(language)
		if err := db.Select(fields).Table("variance_trailer vt").Joins(join).Where(where, ContentKey, country).Order("vt.order").Find(&trailers).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
	} else {
		fields, join, where := common.GetSeasonTrailerQuery(language)
		if err := db.Select(fields).Table("variance_trailer vt").Joins(join).Where(where, ContentKey, country).Order("vt.order").Find(&trailers).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
	}
	for _, trailer := range trailers {
		var imagery TrailerImagery
		imagery.TrailerPosterImage = os.Getenv("IMAGES") + trailer.ContentId + "/" + trailer.VarianceId + "/" + trailer.Id + "/trailer-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
		trailer.Imagery = imagery
		result = append(result, trailer)
	}
	l.JSON(c, http.StatusOK, result)
	return
}

// GetMediaObjectDetails -  Get mediaobject details
// GET /v1/:lang/mediaobject/:ids
// @Summary Get content rating details
// @Description Get content rating details based on content key
// @Tags Content
// @Accept  json
// @Produce  json
// @Param lang path string true "Language Code"
// @Param ids path string true "Mediaobject Ids"
// @Success 200
// @Router /v1/{lang}/mediaobject/{ids} [get]
func (hs *HandlerService) GetMediaObjectDetails(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	serverError := common.ServerErrorResponse(language)
	var contentIds, episodeIds []string
	for _, data := range strings.Split(c.Param("ids"), ",") {
		details := strings.Split(data, ".")
		if len(details) == 2 {
			if details[1] == "Episode" {
				episodeIds = append(episodeIds, details[0])
			}
			contentIds = append(contentIds, details[0])
		}
	}
	var results []MediaObjectDetails
	query := common.MeadiaObjectQuery(language)
	if err := db.Raw(query, contentIds, contentIds).Find(&results).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	if episodeIds != nil {
		var episodeResults []MediaObjectDetails
		query := common.EpisodeMeadiaObjectQuery(language)
		if err := db.Raw(query, episodeIds).Find(&episodeResults).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		for _, episode := range episodeResults {
			results = append(results, episode)
		}
	}

	var mediaObjects []MediaObjectDetails
	for _, id := range contentIds {
		conentKey, _ := strconv.Atoi(id)
		for _, details := range results {
			if conentKey == details.Id {
				rating, _ := strconv.Atoi(details.AgeRating)
				details.AgeRating = common.AgeRatings(rating, language)
				details.Geoblock = false
				ImageryDetails := make(chan ContentImageryDetails)
				fmt.Println("details", details.ContentType)
				if details.ContentType == "movie" || details.ContentType == "LiveTV" || details.ContentType == "play" || details.ContentType == "livetv" {
					go OnetierImagery(details.ContentId, details.VarianceId, ImageryDetails)
				} else if details.ContentType == "series" || details.ContentType == "program" {
					go MultitierImagery(details.ContentId, details.VarianceId, ImageryDetails)
				} else {
					go EpisodeImagery(details.ContentId, details.VarianceId, details.EpisodeId, ImageryDetails)
				}
				type BaseContentType struct {
					ContentType string `json:"content_type"`
				}

				var baseContentType BaseContentType

				if err := db.Raw("select content_type from content where id=?", details.ContentId).Find(&baseContentType).Error; err != nil {
					l.JSON(c, http.StatusInternalServerError, serverError)
					return
				}

				details.Imagery = <-ImageryDetails
				details.BaseContentType = strings.ToLower(baseContentType.ContentType)
				mediaObjects = append(mediaObjects, details)
			}
		}
	}

	l.JSON(c, http.StatusOK, gin.H{"data": mediaObjects})
	return
}

// GetRelatedContents -  Get related contents
// GET /v1/:lang/related
// @Summary Get related contents
// @Description Get related contents based on content key
// @Tags Content
// @Accept  json
// @Produce  json
// @Param id query string true "Content Key"
// @Param country query string false "Country Code"
// @Param size query string false "Limit"
// @Param q query string false "Query string"
// @Param lang path string true "Language Code"
// @Success 200
// @Router /v1/{lang}/related [get]
func (hs *HandlerService) GetRelatedContents(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	notFoundError := common.NotFoundErrorResponse()
	serverError := common.ServerErrorResponse(language)
	if c.Request.URL.Query()["id"] == nil || c.Request.URL.Query()["id"][0] == "" {
		notFoundError.Description = "ContentKey should not be empty"
		l.JSON(c, http.StatusNotFound, notFoundError)
		return
	}
	ContentKey := c.Request.URL.Query()["id"][0]
	var CountryCode string
	if c.Request.URL.Query()["country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}

	country := int(common.Countrys(CountryCode))
	var genre1, genre2, sgenre1, sgenre2, originalLang string
	var relatedDetails []RelatedContentGenres
	query := common.GetRelatedContentGenreQuery()
	if err := db.Raw(query, ContentKey, ContentKey).Limit(2).Find(&relatedDetails).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	for i, details := range relatedDetails {
		if i == 0 {
			genre1 = details.GenreId
			sgenre1 = details.SubgenreId
			originalLang = details.OriginalLanguage
		} else {
			genre2 = details.GenreId
			sgenre2 = details.SubgenreId
			originalLang = details.OriginalLanguage
		}
	}
	var contentsIds []ContentDetails
	var finalResult []ContentDetails
	var finalids []int
	if finalids == nil {
		finalids = append(finalids, 1)
	}
	if genre1 != "" && genre2 != "" && sgenre1 != "" && sgenre2 != "" {
		query := common.RelatedContentsQuery(1, genre1, genre2, sgenre1, sgenre2, originalLang, country, "20", language, ContentKey)
		if err := db.Raw(query).Find(&contentsIds).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		for _, details := range contentsIds {
			finalids = append(finalids, details.ContentKey)
		}
	}
	finalResult = contentsIds
	if len(finalResult) < 20 && genre1 != "" && genre2 != "" && sgenre1 != "" {
		limit := strconv.Itoa(20 - len(finalResult))
		var contentsIds []ContentDetails
		query := common.RelatedContentsQuery(2, genre1, genre2, sgenre1, sgenre2, originalLang, country, limit, language, ContentKey)
		if err := db.Raw(query, finalids, finalids).Find(&contentsIds).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		for _, details := range contentsIds {
			finalResult = append(finalResult, details)
			finalids = append(finalids, details.ContentKey)
		}
	}
	if len(finalResult) < 20 && genre1 != "" && sgenre1 != "" {
		limit := strconv.Itoa(20 - len(finalResult))
		var contentsIds []ContentDetails
		query := common.RelatedContentsQuery(3, genre1, genre2, sgenre1, sgenre2, originalLang, country, limit, language, ContentKey)
		if err := db.Raw(query, finalids, finalids).Find(&contentsIds).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		for _, details := range contentsIds {
			finalResult = append(finalResult, details)
			finalids = append(finalids, details.ContentKey)
		}
	}
	if len(finalResult) < 20 && originalLang != "" {
		limit := strconv.Itoa(20 - len(finalResult))
		var contentsIds []ContentDetails
		query := common.RelatedContentsQuery(4, genre1, genre2, sgenre1, sgenre2, originalLang, country, limit, language, ContentKey)
		if err := db.Raw(query, finalids, finalids).Find(&contentsIds).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		for _, details := range contentsIds {
			finalResult = append(finalResult, details)
		}
	}
	// var contents []PlaylistContent
	// for _, result := range finalResult {
	// 	ContentDetails := make(chan PlaylistContent)
	// 	if result.ContentTier == 1 {
	// 		go OneTierContentDetails(result.ContentId, language, country, c, ContentDetails)
	// 	} else {
	// 		go MultiTierContentDetails(result.ContentId, language, country, c, ContentDetails, 1)
	// 	}
	// 	contents = append(contents, <-ContentDetails)
	// }
	var contents []ContentDetails
	for _, details := range finalResult {
		friendlyUrl := strings.ToLower(details.FriendlyUrl)
		details.FriendlyUrl = strings.Replace(friendlyUrl, " ", "-", -1)
		details.Geoblock = false
		ImageryDetails := make(chan ContentImageryDetails)
		if details.ContentTier == 1 {
			go OnetierImagery(details.ContentId, details.SeasonOrVarienceId, ImageryDetails)
		} else {
			go MultitierImagery(details.ContentId, details.SeasonOrVarienceId, ImageryDetails)
		}
		details.Imagery = <-ImageryDetails
		contents = append(contents, details)
	}
	l.JSON(c, http.StatusOK, gin.H{"data": contents})
	return
}

// GetUserPlaylists -  Get user playlist
// GET /v1/contents/playlist
// @Summary Get all user playlists
// @Description Get all user playlists by user id
// @Tags User
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param limit query string false "Limit"
// @Success 200
// @Router /v1/contents/playlist [get]
func (hs *HandlerService) GetUserPlaylists(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	UserId := c.MustGet("userid")
	if c.MustGet("AuthorizationRequired") == 1 {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request.", "status": http.StatusUnauthorized})
		return
	}

	if UserId == "" {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Invalid Authorization", "status": http.StatusUnauthorized})
		return
	}
	limit := 100
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.Atoi(c.Request.URL.Query()["limit"][0])
	}
	if limit > 100 {
		limit = 100
	}
	response := []ContentRatingDetails{}
	var contentKeys []ContentKeys
	var pagination PaginationResult
	var totalCount int
	if err := db.Table("content c").Select("c.content_key::text").Joins("join playlisted_content pc on pc.content_id =c.id").Where("pc.user_id=?", UserId).Order("pc.added_at desc").Count(&totalCount).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "حدث خطأ ما", "code": "error_server_error", "requestId": UserId.(string)})
		return
	}
	if err := db.Table("content c").Select("c.content_key::text").Joins("join playlisted_content pc on pc.content_id =c.id").Where("pc.user_id=?", UserId).Order("pc.added_at desc").Limit(limit).Find(&contentKeys).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "حدث خطأ ما", "code": "error_server_error", "requestId": UserId.(string)})
		return
	}
	for _, cKey := range contentKeys {
		details, err := GetContentRatingDetails(cKey.ContentKey, UserId.(string), 0, 0, c)
		if err != nil {
			l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "حدث خطأ ما", "code": "error_server_error", "requestId": UserId.(string)})
			return
		}
		if details.Content.Id != 0 {
			response = append(response, details)
		}
	}
	pagination.Size = int64(totalCount)
	pagination.Offset = 0
	pagination.Limit = int64(limit)
	l.JSON(c, http.StatusOK, gin.H{"pagination": pagination, "data": response})
	return
}

// GetRatedContents -  Get Rated Contents
// GET /v1/contents/rated
// @Summary Get all rated contents
// @Description Get all user rated contents by user id
// @Tags User
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param limit query string false "Limit"
// @Success 200
// @Router /v1/contents/rated [get]
func (hs *HandlerService) GetRatedContents(c *gin.Context) {
	//db := c.MustGet("CDB").(*gorm.DB)
	dbro := c.MustGet("CDBRO").(*gorm.DB)
	serverError := common.ServerErrorResponse("en")
	if c.MustGet("AuthorizationRequired") == 1 {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	UserId := c.MustGet("userid")
	limit := 10
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.Atoi(c.Request.URL.Query()["limit"][0])
	}
	if limit > 10 {
		limit = 10
	}
	response := []ContentRatingDetails{}
	var contentKeys []ContentKeys
	var pagination PaginationResult
	var totalCount int
	if err := dbro.Table("content c").Select("c.content_key::text").Joins("join rated_content rc on rc.content_id =c.id").Where("rc.user_id=? and rc.is_hidden=false", UserId).Order("rc.rated_at desc").Count(&totalCount).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	if err := dbro.Table("content c").Select("c.content_key::text").Joins("join rated_content rc on rc.content_id =c.id").Where("rc.user_id=? and rc.is_hidden=false", UserId).Order("rc.rated_at desc").Limit(limit).Find(&contentKeys).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	for _, cKey := range contentKeys {
		details, err := GetContentRatingDetails(cKey.ContentKey, UserId.(string), 0, 0, c)
		if err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		if details.Content.Id != 0 {
			response = append(response, details)
		}
	}
	pagination.Size = int64(totalCount)
	pagination.Offset = 0
	pagination.Limit = int64(limit)
	l.JSON(c, http.StatusOK, gin.H{"pagination": pagination, "data": response})
	return
}

// GetWatchingContents -  Get Watching Contents
// GET /v1/contents/watching
// @Summary Get all watching contents
// @Description Get all watching contents by user id
// @Tags User
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param limit query string false "Limit"
// @Success 200
// @Router /v1/contents/watching [get]
func (hs *HandlerService) GetWatchingContents(c *gin.Context) {
	dbro := c.MustGet("CDBRO").(*gorm.DB)
	serverError := common.ServerErrorResponse("en")
	if c.MustGet("AuthorizationRequired") == 1 {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	UserId := c.MustGet("userid")
	// Hard coded limit to 100 as per .net production
	limit := 100
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.Atoi(c.Request.URL.Query()["limit"][0])
	}
	response := []ContentRatingDetails{}
	var contentKeys []ContentKeys
	var pagination PaginationResult
	var totalCount int
	if err := dbro.Table("content c").Select("c.content_key::text").Joins("join view_activity va on va.content_id =c.id and va.user_id=?", UserId).Where("va.user_id=? and va.is_hidden=?", UserId, false).Order("va.viewed_at desc").Count(&totalCount).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	if err := dbro.Table("content c").Select("c.content_key::text").Joins("join view_activity va on va.content_id =c.id and va.user_id=?", UserId).Where("va.user_id=? and va.is_hidden=?", UserId, false).Order("va.viewed_at desc").Limit(limit).Find(&contentKeys).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	for _, cKey := range contentKeys {
		details, err := GetContentRatingDetails(cKey.ContentKey, UserId.(string), 1, 1, c)
		if err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		if details.Content.Id != 0 && details.UserData.ViewActivity != nil {
			response = append(response, details)
		}
	}
	pagination.Size = int64(totalCount)
	pagination.Offset = 0
	pagination.Limit = int64(limit)
	l.JSON(c, http.StatusOK, gin.H{"pagination": pagination, "data": response})
	return
}

// GetResumbleContents -  Get Resumble Contents
// GET /v1/contents/resumable
// @Summary Get all resumble contents
// @Description Get all continue watching contents by user id
// @Tags User
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 200
// @Router /v1/contents/resumable [get]
func (hs *HandlerService) GetResumbleContents(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	serverError := common.ServerErrorResponse("en")
	if c.MustGet("AuthorizationRequired") == 1 {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	UserId := c.MustGet("userid")
	limit := 100
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.Atoi(c.Request.URL.Query()["limit"][0])
	}
	if limit > 100 {
		limit = 100
	}
	response := []ResumbleContentRatingDetails{}
	var contentRating []ContentRating
	var contentKeys []ContentKeys
	var pagination PaginationResult
	// Here Duration condition is handled because some contents have duration zero,This must be changed in coming stages as per .net prod scenarios
	// and pi2.duration != 0
	if err := db.Debug().Table("content c").Select("c.id").Joins("join view_activity va on va.content_id =c.id join playback_item pi2 on pi2.id=va.playback_item_id ").Where("va.user_id=? and va.last_watch_position*100/pi2.duration <95 and va.last_watch_position != 0 and pi2.duration != 0", UserId).Order("va.viewed_at desc").Limit(limit).Find(&contentKeys).Error; err != nil {
		serverError.Description = common.EN_SERVER_ERROR_DESCRIPTION
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var ids []string
	for _, cKey := range contentKeys {
		ids = append(ids, cKey.Id)
	}
	query := common.ResumbleContentsQuery()
	if err := db.Raw(query, UserId, UserId).Order("viewed_at desc").Find(&contentRating).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var playlistedContent []PlaylistedContent
	db.Where("user_id=? and content_id in(?)", UserId, ids).Find(&playlistedContent)
	var ratedContent []RatedContent
	db.Where("user_id=? and content_id in(?) and is_hidden=false", UserId, ids).Find(&ratedContent)
	genres := make([]string, 0)
	for _, details := range contentRating {
		var resumableContent ResumbleContentRatingDetails
		var userRating ResumbleUserRating
		var ratingViewActivity ResumbleRatingViewActivity
		details.Genres = genres
		resumableContent.Content = details
		resumableContent.UserData = &userRating
		ratingViewActivity.ViewedAt = details.ViewedAt
		ratingViewActivity.ResumeWatchPosition = details.LastWatchPosition
		userRating.ViewActivity = ratingViewActivity
		for _, rcontent := range ratedContent {
			if rcontent.ContentId == details.ContentId {
				userRating.Rating = &rcontent.Rating
				userRating.RatedAt = &rcontent.RatedAt
				userRating.IsTailored = false
				break
			}
		}
		for _, pcontent := range playlistedContent {
			if pcontent.ContentId == details.ContentId {
				userRating.AddedToPlaylistAt = &pcontent.AddedAt
				break
			}
		}
		resumableContent.UserData = &userRating
		response = append(response, resumableContent)
	}
	pagination.Size = int64(len(response))
	pagination.Offset = 0
	pagination.Limit = int64(limit)
	l.JSON(c, http.StatusOK, gin.H{"pagination": pagination, "data": response})
	return
}

// GetResumbleContents -  Get Resumble Contents
// GET /v1/contents/resumable
// @Summary Get all resumble contents
// @Description Get all continue watching contents by user id
// @Tags User
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 200
// @Router /flutter/continuewatching [get]
func (hs *HandlerService) GetResumbleContentsFlutterContinuewatching(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	serverError := common.ServerErrorResponse("en")
	if c.MustGet("AuthorizationRequired") == 1 {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	UserId := c.MustGet("userid")
	limit := 100
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.Atoi(c.Request.URL.Query()["limit"][0])
	}
	if limit > 100 {
		limit = 100
	}
	response := []ResumbleContentRatingDetails{}

	finalContent := []ContentRatingFlutter{}

	var finalResponse ResumbleContentRatingDetailsForFlutter

	var contentRating []ContentRatingQueryFlutter
	var contentKeys []ContentKeys
	var pagination PaginationResult
	// Here Duration condition is handled because some contents have duration zero,This must be changed in coming stages as per .net prod scenarios
	// and pi2.duration != 0
	if err := db.Debug().Table("content c").Select("c.id").Joins("join view_activity va on va.content_id =c.id join playback_item pi2 on pi2.id=va.playback_item_id ").Where("va.user_id=? and va.last_watch_position*100/pi2.duration <95 and va.last_watch_position != 0 and pi2.duration != 0", UserId).Order("va.viewed_at desc").Limit(limit).Find(&contentKeys).Error; err != nil {
		serverError.Description = common.EN_SERVER_ERROR_DESCRIPTION
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}

	// var results []MediaObjectDetails
	// query := common.MeadiaObjectQuery(language)
	// if err := db.Raw(query, contentIds, contentIds).Find(&results).Error; err != nil {
	// 	l.JSON(c, http.StatusInternalServerError, serverError)
	// 	return
	// }
	// if episodeIds != nil {
	// 	var episodeResults []MediaObjectDetails
	// 	query := common.EpisodeMeadiaObjectQuery(language)
	// 	if err := db.Raw(query, episodeIds).Find(&episodeResults).Error; err != nil {
	// 		l.JSON(c, http.StatusInternalServerError, serverError)
	// 		return
	// 	}
	// 	for _, episode := range episodeResults {
	// 		results = append(results, episode)
	// 	}
	// }

	var ids []string
	for _, cKey := range contentKeys {
		ids = append(ids, cKey.Id)
	}

	query := common.ResumbleContentsQueryForFlutter()
	if err := db.Raw(query, UserId, UserId).Order("viewed_at desc").Find(&contentRating).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}

	var playlistedContent []PlaylistedContent
	db.Where("user_id=? and content_id in(?)", UserId, ids).Find(&playlistedContent)
	var ratedContent []RatedContent
	db.Where("user_id=? and content_id in(?) and is_hidden=false", UserId, ids).Find(&ratedContent)
	genres := make([]string, 0)

	for _, details := range contentRating {
		var resumableContent ResumbleContentRatingDetails
		var userRating ResumbleUserRating
		var ratingViewActivity ResumbleRatingViewActivity

		var contentImagesFlutter []ContentImagesFlutter

		details.Genres = genres
		// resumableContent.Content = details

		imageID := details.ContentId

		if details.ContentType == "Movie" {
			imageID = details.ContentId
		} else {
			imageID = details.ContentId + "/" + details.MultiTierContentId
		}

		//posterImage
		contentImagesFlutter = append(contentImagesFlutter, ContentImagesFlutter{
			ImageCategory: "posterImage",
			ImageUrl: []string{
				fmt.Sprintf("%s%s/poster-image", os.Getenv("IMAGES"), imageID),
			},
		})

		//detailsPageBackground
		contentImagesFlutter = append(contentImagesFlutter, ContentImagesFlutter{
			ImageCategory: "detailsPageBackground",
			ImageUrl: []string{
				fmt.Sprintf("%s%s/details-background", os.Getenv("IMAGES"), imageID),
			},
		})

		//mobileDetailsPageBackground
		contentImagesFlutter = append(contentImagesFlutter, ContentImagesFlutter{
			ImageCategory: "mobileDetailsPageBackground",
			ImageUrl: []string{
				fmt.Sprintf("%s%s/mobile-details-background", os.Getenv("IMAGES"), imageID),
			},
		})

		var playlistedContentfor PlaylistedContent
		Isaddtoplaylist := false
		var BaseContentType string

		if rows := db.Table("playlisted_content pc").Where("user_id=? and content_id = (?)", UserId, details.ContentId).Find(&playlistedContentfor).RowsAffected; rows != 0 {
			Isaddtoplaylist = true
		}

		switch details.ContentType {
		case "Episode":
			BaseContentType = "Series"
		case "Movie":
			BaseContentType = "Movie"
		}

		finalContent = append(finalContent, ContentRatingFlutter{
			ContentType:        details.ContentType,
			BaseContentType:    BaseContentType,
			ID:                 details.ContentId,
			Key:                details.Id,
			Images:             contentImagesFlutter,
			MultiTierContentId: details.MultiTierContentId,
			Name:               details.Title,
			OneTierContentId:   details.OneTierContentId,
			Isaddtoplaylist:    Isaddtoplaylist,
		})

		resumableContent.UserData = &userRating
		ratingViewActivity.ViewedAt = details.ViewedAt
		ratingViewActivity.ResumeWatchPosition = details.LastWatchPosition
		userRating.ViewActivity = ratingViewActivity
		for _, rcontent := range ratedContent {
			if rcontent.ContentId == details.ContentId {
				userRating.Rating = &rcontent.Rating
				userRating.RatedAt = &rcontent.RatedAt
				userRating.IsTailored = false
				break
			}
		}
		for _, pcontent := range playlistedContent {
			if pcontent.ContentId == details.ContentId {
				userRating.AddedToPlaylistAt = &pcontent.AddedAt
				break
			}
		}
		resumableContent.UserData = &userRating
		response = append(response, resumableContent)
	}
	pagination.Size = int64(len(response))
	pagination.Offset = 0
	pagination.Limit = int64(limit)

	lang := "en"
	if c.Request.URL.Query()["lang"] != nil {
		lang = c.Request.URL.Query()["lang"][0]
	}

	finalResponse.Language = common.OriginalLanguage(lang)
	finalResponse.Content = finalContent
	finalResponse.Type = "watching"
	finalResponse.Title = "continue watching"

	l.JSON(c, http.StatusOK, gin.H{"data": finalResponse, "status": "redisResponse-main"})
	return
}

func GetEpisodeRatingDetails(ContentKey string, UserId string, episode, vActivity int, c *gin.Context) (ContentRatingDetails, error) {
	dbro := c.MustGet("CDBRO").(*gorm.DB)
	var details ContentRatingDetails
	var result ContentRatingQueryEpisode
	query := common.ContentRatingQueryForEpisode(ContentKey)
	if err := dbro.Raw(query).Find(&result).Error; err == nil {
		genres := make([]string, 0)
		var contentRating ContentRating
		var viewActivity ViewActivityDetails
		var activityDetails ViewActivityDetails
		contentRating.AverageRating = result.AverageRating
		contentRating.Title = result.TransliteratedTitle
		contentRating.DigitalRightsType = result.DigitalRightsType
		if result.ContentTier == 1 {
			contentRating.Id = result.ContentKey
			contentRating.ContentType = result.ContentType
			contentRating.Duration = result.Length
		} else {
			contentRating.Id = result.ContentKey
			contentRating.ContentType = result.ContentType
			contentRating.Duration = result.Length
			if episode == 1 && UserId != "" {
				fields, join, where := common.ViewActivityDetailsEpisodeQuery()
				if err := dbro.Table("view_activity va").Select(fields).Joins(join).Where(where, result.PlaybackItemId, UserId).Order("va.viewed_at desc").Limit(1).Find(&viewActivity).Error; err == nil {
					contentRating.Id = viewActivity.EpisodeKey
					contentRating.ContentType = viewActivity.ContentType
					contentRating.Duration = viewActivity.Duration
				}
			}
		}
		contentRating.Genres = genres
		var userRating UserRating
		userRating.ViewActivity = nil
		// if vActivity == 1 {
		var ratingViewActivity RatingViewActivity

		// if err := dbro.Table("view_activity_history").Select("viewed_at,last_watch_position").Where("(content_type_name = 'Episode' or content_type_name = 'episode') and content_key = ? and user_id =?", ContentKey, UserId).Order("viewed_at desc").Limit(1).Find(&activityDetails).Error; err == nil {
		// 	ratingViewActivity.ViewedAt = activityDetails.ViewedAt
		// 	ratingViewActivity.ResumeWatchPosition = &activityDetails.LastWatchPosition
		// 	userRating.ViewActivity = &ratingViewActivity
		// }

		if err := dbro.Table("view_activity").Select("viewed_at,last_watch_position").Where("playback_item_id =? and user_id =? and is_hidden=false", result.PlaybackItemId, UserId).Order("viewed_at desc").Limit(1).Find(&activityDetails).Error; err == nil {
			ratingViewActivity.ViewedAt = activityDetails.ViewedAt
			ratingViewActivity.ResumeWatchPosition = &activityDetails.LastWatchPosition
			userRating.ViewActivity = &ratingViewActivity
		}

		// }
		if UserId != "" {
			details.UserData = &userRating
			var playlistedContent PlaylistedContent
			if err := dbro.Where("user_id=? and content_id=?", UserId, result.Id).Find(&playlistedContent).Error; err == nil {
				userRating.AddedToPlaylistAt = &playlistedContent.AddedAt
			}
			var ratedContent RatedContent
			if err := dbro.Where("user_id=? and content_id=? and is_hidden=false", UserId, result.Id).Find(&ratedContent).Error; err == nil {
				userRating.Rating = &ratedContent.Rating
				userRating.RatedAt = &ratedContent.RatedAt
				userRating.IsTailored = false
				details.UserData = &userRating
			}
		}
		details.Content = contentRating
	}
	return details, nil
}

func GetContentRatingDetails(ContentKey string, UserId string, episode, vActivity int, c *gin.Context) (ContentRatingDetails, error) {
	dbro := c.MustGet("CDBRO").(*gorm.DB)
	var details ContentRatingDetails
	var result ContentRatingQuery
	query := common.ContentRatingQuery(ContentKey)
	if err := dbro.Raw(query).Find(&result).Error; err == nil {
		genres := make([]string, 0)
		var contentRating ContentRating
		var viewActivity ViewActivityDetails
		var activityDetails ViewActivityDetails
		contentRating.AverageRating = result.AverageRating
		contentRating.Title = result.TransliteratedTitle
		contentRating.DigitalRightsType = result.DigitalRightsType
		if result.ContentTier == 1 {
			contentRating.Id = result.ContentKey
			contentRating.ContentType = result.ContentType
			contentRating.Duration = result.Length
		} else {
			contentRating.Id = result.ContentKey
			contentRating.ContentType = result.ContentType
			contentRating.Duration = result.Length
			if episode == 1 && UserId != "" {
				fields, join, where := common.ViewActivityDetailsQuery()
				if err := dbro.Table("view_activity va").Select(fields).Joins(join).Where(where, result.Id, UserId).Order("va.viewed_at desc").Limit(1).Find(&viewActivity).Error; err == nil {
					contentRating.Id = viewActivity.EpisodeKey
					contentRating.ContentType = viewActivity.ContentType
					contentRating.Duration = viewActivity.Duration
				}
			}
		}
		contentRating.Genres = genres
		var userRating UserRating
		userRating.ViewActivity = nil
		// if vActivity == 1 {
		var ratingViewActivity RatingViewActivity
		if err := dbro.Table("view_activity").Select("viewed_at,last_watch_position").Where("content_id =? and user_id =? and is_hidden=false", result.Id, UserId).Order("viewed_at desc").Limit(1).Find(&activityDetails).Error; err == nil {
			ratingViewActivity.ViewedAt = activityDetails.ViewedAt
			ratingViewActivity.ResumeWatchPosition = &activityDetails.LastWatchPosition
			userRating.ViewActivity = &ratingViewActivity
		}
		// }
		if UserId != "" {
			details.UserData = &userRating
			var playlistedContent PlaylistedContent
			if err := dbro.Where("user_id=? and content_id=?", UserId, result.Id).Find(&playlistedContent).Error; err == nil {
				userRating.AddedToPlaylistAt = &playlistedContent.AddedAt
			}
			var ratedContent RatedContent
			if err := dbro.Where("user_id=? and content_id=? and is_hidden=false", UserId, result.Id).Find(&ratedContent).Error; err == nil {
				userRating.Rating = &ratedContent.Rating
				userRating.RatedAt = &ratedContent.RatedAt
				userRating.IsTailored = false
				details.UserData = &userRating
			}
		}
		details.Content = contentRating
	}
	return details, nil
}

// AddViewActivity -  Add watching contents into history
// POST /v1/contents/watching
// @Summary Add watching contents into history
// @Description Add watching contents into history by user id
// @Tags User
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param body body AddViewActivityRequest true "Raw JSON string"
// @Success 200
// @Router /v1/contents/watching [post]
func (hs *HandlerService) AddViewActivity(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	dbro := c.MustGet("CDBRO").(*gorm.DB)
	serverError := common.ServerErrorResponse("en")
	badRequestError := common.BadRequestErrorResponse()
	var invalidRequestError common.InvalidRequestError
	var invalidRequest common.InvalidRequest
	if c.MustGet("AuthorizationRequired") == 1 {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	UserId := c.MustGet("userid")
	deviceId := c.MustGet("device_id")
	var request AddViewActivityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		badRequestError.Description = err.Error()
		badRequestError.Invalid = &invalidRequest
		l.JSON(c, http.StatusBadRequest, badRequestError)
		return
	}
	if request.Content.ContentType != "LiveTV" {
		if request.LastWatchPosition > request.Content.Duration {
			invalidRequestError.Code = "lessthanorequal_error"
			invalidRequestError.Description = "'Last Watch Position' must be less than or equal to 'Duration'."
			invalidRequest.LastWatchPosition = &invalidRequestError
			badRequestError.Invalid = &invalidRequest
			l.JSON(c, http.StatusBadRequest, badRequestError)
			return
		}
		var cDetails ContentIdDetails
		query := common.GetWatchingContentDetailsQuery(request.Content.Id)
		if err := dbro.Raw(query).Find(&cDetails).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var input, update ViewActivity
		var historyInput ViewActivityHistory
		var viewActivity ViewActivity
		update.IsHidden = "true"
		if err := db.Table("view_activity").Where("content_id=? and user_id=? and is_hidden=false", cDetails.Id, UserId).Update(&update).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		if rows := dbro.Table("view_activity").Where("playback_item_id=? and content_id=? and user_id=?", cDetails.PlaybackItemId, cDetails.Id, UserId).Find(&viewActivity).RowsAffected; rows != 0 {
			if err := db.Table("view_activity va").Where("id=?", viewActivity.Id).Update(map[string]interface{}{
				"viewed_at":           time.Now(),
				"last_watch_position": request.LastWatchPosition,
				"watch_session_id":    request.WatchSessionId,
				"is_hidden":           false,
				"device_id":           deviceId.(string),
			}).Error; err != nil {
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}

			// db.Debug().Raw("UPDATE view_activity SET is_hidden = 'false'  WHERE (content_id=? and user_id=?)", cDetails.Id, UserId)
		} else {
			input.ContentId = cDetails.Id
			input.UserId = UserId.(string)
			input.DeviceId = deviceId.(string)
			input.LastWatchPosition = request.LastWatchPosition
			input.ViewedAt = time.Now()
			input.WatchSessionId = request.WatchSessionId
			input.IsHidden = "false"
			input.PlaybackItemId = cDetails.PlaybackItemId
			if err := db.Create(&input).Error; err != nil {
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}
		}
		historyInput.ContentTypeName = request.Content.ContentType
		historyInput.ContentKey = request.Content.Id
		historyInput.UserId = UserId.(string)
		historyInput.DeviceId = deviceId.(string)
		historyInput.LastWatchPosition = request.LastWatchPosition
		historyInput.ViewedAt = time.Now()
		historyInput.WatchSessionId = request.WatchSessionId
		// historyJson, err := json.Marshal(historyInput)
		// if err != nil {
		// 	fmt.Println("Cannot encode to JSON ", err)
		// }
		// fmt.Println(string(historyJson))
		if err := db.Create(&historyInput).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		//TODO: Episode watching history updated
		if strings.ToLower(request.Content.ContentType) == "episode" {

			var (
				episodesIds []int
				episode     EpisodeDetailsSummary
				season      SeasonDetailsSummary
				seasonIds   []string
			)

			if err := db.Table("episode").Where("episode_key = ?", request.Content.Id).Find(&episode).Error; err != nil {
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}

			if err := db.Table("season").Where("id = ?", episode.SeasonId).Find(&season).Error; err != nil {
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}

			if err := db.Table("season").Select("id").Where("content_id = ?", season.ContentId).Pluck("id", &seasonIds).Error; err != nil {
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}

			if err := db.Table("episode").Select("episode_key").Where("season_id in (?) AND episode_key != ?", seasonIds, request.Content.Id).Pluck("episode_key", &episodesIds).Error; err != nil {
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}

			var episodesIdstr []string
			for _, ep := range episodesIds {
				episodesIdstr = append(episodesIdstr, strconv.Itoa(ep))
			}

			episode_keys := strings.Join(episodesIdstr, ",")

			query := common.GetWatchingContentDetailsQueryStr(episode_keys, UserId.(string))

			var viewActivity []string

			if err := db.Raw(query).Pluck("id", &viewActivity).Error; err != nil {
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}

			if err := db.Table("view_activity va").Where("id in (?)", viewActivity).Update(map[string]interface{}{
				"last_watch_position": 0,
			}).Error; err != nil {
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}

		}
		//TODO: Episode watching history updated

	}
	l.JSON(c, http.StatusOK, gin.H{"message": "History Added successfully."})
	return
}

// searchContent
// GetsearchContent- Get content by search
// GET /v1/:lang/search
// @Summary Get movie details by content
// @Description Get movie details based on content key
// @Tags Content
// @Accept json
// @Produce json
// @Param country query string true "Country Code"
// @Param q query string true "Q string"
// @Param lang path string true "Language Code"
// @Success 200
// @Router /v1/{lang}/search [get]
func (hs *HandlerService) GetsearchContent(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	serverError := common.ServerErrorResponse(language)
	if c.Request.URL.Query()["q"] == nil || c.Request.URL.Query()["q"][0] == "" {
		serverError.Description = "Please Provide valid search text."
		l.JSON(c, http.StatusBadRequest, serverError)
		return
	}

	querystringwithExtras := strings.ToLower(c.Request.URL.Query()["q"][0])
	querystring := strings.Replace(querystringwithExtras, "'", "''", -1)

	var CountryCode string
	if c.Request.URL.Query()["Country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["Country"][0])
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}
	country := int(common.Countrys(CountryCode))
	countrystr := strconv.Itoa(country)

	var Platform string

	if c.Request.URL.Query()["platform"] != nil {
		Platform = strings.ToUpper(c.Request.URL.Query()["platform"][0])
	}

	if Platform == "" {
		Platform = "0"
	}

	// not to delete
	/*var searchcontents []SearchContent
	contents := []PlaylistContent{}
	var finalResult []SearchContent
	query := "select c.id, c.content_type,c.content_key,c.content_tier from content_primary_info cpi join content c on cpi.id = c.primary_info_id where c.deleted_by_user_id is null and c.status =1 and  (lower(cpi.transliterated_title) like '%" + querystring + "%' or lower(cpi.arabic_title) like '%" + querystring + "%')"

	if err := db.Raw(query).Find(&searchcontents).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	finalResult = searchcontents
	if len(finalResult) < 5 {
		searchcontents = nil
		aquery := common.SearchContentByCastQuery(querystring)
		if err := db.Raw(aquery).Find(&searchcontents).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		for _, actor := range searchcontents {
			finalResult = append(finalResult, actor)
		}
	}

	ContentDetails := make(chan PlaylistContent)

	for _, searchcontent := range finalResult {
		if searchcontent.ContentTier == 1 {
			go OneTierContentDetails(searchcontent.Id, language, country, c, ContentDetails)

		} else {
			go MultiTierContentDetails(searchcontent.Id, language, country, c, ContentDetails, 1)
		}
		details := <-ContentDetails
		if details.ID != 0 {
			contents = append(contents, details)
		}
		if searchcontent.ContentTier != 1 {
			if details.ID <= 0 {
				go MultiTierContentDetailsWithoutEpisodeForSearch(searchcontent.Id, language, country, c, ContentDetails, 1)
				details1 := <-ContentDetails
				if details1.ID != 0 {
					contents = append(contents, details1)
				}
			}
		}

	}
	l.JSON(c, http.StatusOK, gin.H{"data": contents})
	return*/
	var searchcontents []SearchDetails
	result := []Search{}
	var temp Search
	// searchquery := "((select distinct on (c.content_key)c.content_key,c.id ,c.content_type ,cr.digital_rights_type ,c.content_tier ,c.created_at ,pi1.video_content_id ,cpi.transliterated_title ,cpi.arabic_title,c.has_poster_image,cv.id as variance_id from content c join content_primary_info cpi on cpi.id = c.primary_info_id join about_the_content_info atci on atci.id = c.about_the_content_info_id join content_variance cv on cv.content_id = c.id join playback_item pi1 on pi1.id = cv.playback_item_id join content_rights cr on cr.id = pi1.rights_id join content_translation ct on ct.id = pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id = cr.id full outer join content_rights_plan crp on crp.rights_id = cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time <= now() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= now() or cr.digital_rights_end_date is null) and crc.country_id = " + countrystr + " and (cpi.original_title ilike '%" + querystring + "%' or cpi.transliterated_title ilike '%" + querystring + "%' or cpi.alternative_title ilike '%" + querystring + "%' or cpi.arabic_title ilike '%" + querystring + "%') and " + Platform + " in (select pitp2.target_platform from playback_item_target_platform pitp2 where playback_item_id = pi1.id) group by c.content_key,c.id,c.content_type ,cr.digital_rights_type ,c.content_tier ,c.created_at ,cpi.transliterated_title ,cpi.arabic_title ,pi1.video_content_id,cv.id) union(select distinct on (c.content_key)c.content_key ,c.id ,c.content_type ,cr.digital_rights_type ,c.content_tier ,c.created_at ,min(pi1.video_content_id),cpi.transliterated_title ,cpi.arabic_title,s.has_poster_image,s.id as variance_id from content c join season s on s.content_id = c.id left  join episode e on e.season_id = s.id left join content_primary_info cpi on cpi.id = c.primary_info_id left join about_the_content_info atci on atci.id = s.about_the_content_info_id left join playback_item pi1 on pi1.id = e.playback_item_id left join content_rights cr on cr.id = s.rights_id left join content_translation ct on ct.id = pi1.translation_id full outer  join content_rights_country crc on crc.content_rights_id = cr.id full outer  join content_rights_plan crp on crp.rights_id = cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time  <= now() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date  is null) and (cr.digital_rights_end_date  >= now() or cr.digital_rights_end_date  is null) and s.status = 1 and s.deleted_by_user_id  is null and (e.status = 1 or e.status is null) and e.deleted_by_user_id  is null and crc.country_id  = " + countrystr + " and (cpi.original_title  ilike '%" + querystring + "%' or cpi.transliterated_title  ilike '%" + querystring + "%' or cpi.alternative_title  ilike '%" + querystring + "%' or cpi.arabic_title  ilike '%" + querystring + "%') and " + Platform + " in (select pitp2.target_platform from playback_item_target_platform pitp2 where playback_item_id = pi1.id) group by c.content_key,c.id,c.content_type,cr.digital_rights_type,c.content_tier,c.created_at,cpi.transliterated_title ,cpi.arabic_title,s.has_poster_image,s.id ))order by created_at desc"

	// searchquery := ` (( select distinct on (c.content_key)c.content_key, c.id , c.content_type , cr.digital_rights_type , c.content_tier , c.created_at , pi1.video_content_id , cpi.transliterated_title , cpi.arabic_title, c.has_poster_image, cv.id as variance_id from content c join content_primary_info cpi on cpi.id = c.primary_info_id join about_the_content_info atci on atci.id = c.about_the_content_info_id join content_variance cv on cv.content_id = c.id join playback_item pi1 on pi1.id = cv.playback_item_id join content_rights cr on cr.id = pi1.rights_id join content_translation ct on ct.id = pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id = cr.id full outer join content_rights_plan crp on crp.rights_id = cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time <= now() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= now() or cr.digital_rights_end_date is null) and crc.country_id = ` + countrystr + ` and (cpi.original_title ilike '%` + querystring + `%' or cpi.transliterated_title ilike '%` + querystring + `%' or cpi.alternative_title ilike '%` + querystring + `%' or cpi.arabic_title ilike '%` + querystring + `%') and ` + Platform + ` in ( select pitp2.target_platform from playback_item_target_platform pitp2 where playback_item_id = pi1.id) group by c.content_key, c.id, c.content_type , cr.digital_rights_type , c.content_tier , c.created_at , cpi.transliterated_title , cpi.arabic_title , pi1.video_content_id, cv.id) union( select distinct on (c.content_key)c.content_key , c.id , c.content_type , cr.digital_rights_type , c.content_tier , c.created_at , min(pi1.video_content_id), cpi.transliterated_title , cpi.arabic_title, s.has_poster_image, s.id as variance_id from content c join season s on s.content_id = c.id left join episode e on e.season_id = s.id left join content_primary_info cpi on cpi.id = c.primary_info_id left join about_the_content_info atci on atci.id = s.about_the_content_info_id left join playback_item pi1 on pi1.id = e.playback_item_id left join content_rights cr on cr.id = s.rights_id left join content_translation ct on ct.id = pi1.translation_id left join variance_trailer vt on vt.season_id = s.id full outer join content_rights_country crc on crc.content_rights_id = cr.id full outer join content_rights_plan crp on crp.rights_id = cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time <= now() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= now() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null and (e.status = 1 or e.status is null) and ( e.deleted_by_user_id is null or vt.id is not null ) and crc.country_id = ` + countrystr + ` and (cpi.original_title ilike '%` + querystring + `%' or cpi.transliterated_title ilike '%` + querystring + `%' or cpi.alternative_title ilike '%` + querystring + `%' or cpi.arabic_title ilike '%` + querystring + `%') and ( (` + Platform + ` in (select pitp2.target_platform from  playback_item_target_platform pitp2 where playback_item_id = pi1.id) and e.id is not null ) or (e.id is null and vt.id is not null) ) group by c.content_key, c.id, c.content_type, cr.digital_rights_type, c.content_tier, c.created_at, cpi.transliterated_title , cpi.arabic_title, s.has_poster_image, s.id )) order by created_at desc `
	searchquery := ` (( select distinct on (c.content_key)c.content_key, c.id , c.content_type , cr.digital_rights_type , c.content_tier , c.created_at , pi1.video_content_id , cpi.transliterated_title , cpi.arabic_title, c.has_poster_image, cv.id as variance_id from content c join content_primary_info cpi on cpi.id = c.primary_info_id join about_the_content_info atci on atci.id = c.about_the_content_info_id join content_variance cv on cv.content_id = c.id join playback_item pi1 on pi1.id = cv.playback_item_id join content_rights cr on cr.id = pi1.rights_id join content_translation ct on ct.id = pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id = cr.id full outer join content_rights_plan crp on crp.rights_id = cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time <= now() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= now() or cr.digital_rights_end_date is null) and crc.country_id = ` + countrystr + ` and (cpi.original_title ilike '%` + querystring + `%' or cpi.transliterated_title ilike '%` + querystring + `%' or cpi.alternative_title ilike '%` + querystring + `%' or cpi.arabic_title ilike '%` + querystring + `%') and ` + Platform + ` in ( select pitp2.target_platform from playback_item_target_platform pitp2 where playback_item_id = pi1.id) group by c.content_key, c.id, c.content_type , cr.digital_rights_type , c.content_tier , c.created_at , cpi.transliterated_title , cpi.arabic_title , pi1.video_content_id, cv.id) union( select distinct on (c.content_key)c.content_key , c.id , c.content_type , cr.digital_rights_type , c.content_tier , c.created_at , min(pi1.video_content_id), cpi.transliterated_title , cpi.arabic_title, s.has_poster_image, s.id as variance_id from content c join season s on s.content_id = c.id left join episode e on e.season_id = s.id left join content_primary_info cpi on cpi.id = c.primary_info_id left join about_the_content_info atci on atci.id = s.about_the_content_info_id left join playback_item pi1 on pi1.id = e.playback_item_id left join content_rights cr on cr.id = s.rights_id left join content_translation ct on ct.id = pi1.translation_id left join variance_trailer vt on vt.season_id = s.id full outer join content_rights_country crc on crc.content_rights_id = cr.id full outer join content_rights_plan crp on crp.rights_id = cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time <= now() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= now() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null and (e.status = 1 or e.status is null) and ( e.deleted_by_user_id is null or vt.id is not null ) and crc.country_id = ` + countrystr + ` and (cpi.original_title ilike '%` + querystring + `%' or cpi.transliterated_title ilike '%` + querystring + `%' or cpi.alternative_title ilike '%` + querystring + `%' or cpi.arabic_title ilike '%` + querystring + `%') ` + ` group by c.content_key, c.id, c.content_type, cr.digital_rights_type, c.content_tier, c.created_at, cpi.transliterated_title , cpi.arabic_title, s.has_poster_image, s.id )) order by created_at desc `

	if err := db.Debug().Limit(100).Raw(searchquery).Find(&searchcontents).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}

	var MultitierIds []string
	for _, val := range searchcontents {
		if val.ContentTier == 2 {
			MultitierIds = append(MultitierIds, val.Id)
		}
	}
	var ValidMultitierIds []string
	var multitiercheck []MultiTierCheck
	db.Debug().Raw("select count(e.id) as episode_count,count(vt.id) as trailer_count,c.id  from content c join season s on s.content_id = c.id left join episode e on e.season_id = s.id left join variance_trailer vt on vt.season_id = s.id where c.id in (?) group by c.id", MultitierIds).Find(&multitiercheck)
	for _, val := range multitiercheck {
		if val.EpisodeCount > 0 || val.TrailerCount > 0 {
			ValidMultitierIds = append(ValidMultitierIds, val.Id)
		}
	}
	for _, val := range searchcontents {
		temp.ID = val.ContentKey
		temp.VideoID = val.VideoContentId
		friendlyurl := strings.ReplaceAll(val.TransliteratedTitle, " ", "-")
		temp.FriendlyURL = friendlyurl
		if language == "en" {
			temp.Title = val.TransliteratedTitle
		} else {
			temp.Title = val.ArabicTitle
		}
		var Imagery ContentImageryDetails
		if val.ContentTier == 1 {
			ImageryDetails := make(chan ContentImageryDetails)
			if val.HasPosterImage {
				go OnetierImagery(val.Id, val.VarianceId, ImageryDetails)
				Imagery = <-ImageryDetails
			} else {
				Imagery.Thumbnail = ""
				Imagery.Backdrop = ""
				Imagery.MobileImg = ""
				Imagery.FeaturedImg = ""
				Imagery.Banner = ""
			}
		} else if val.ContentTier == 2 {
			//Season Imaginery Details
			ImageryDetails := make(chan ContentImageryDetails)
			if val.HasPosterImage {
				go MultitierImagery(val.Id, val.VarianceId, ImageryDetails)
				Imagery = <-ImageryDetails
			} else {
				Imagery.Thumbnail = ""
				Imagery.Backdrop = ""
				Imagery.MobileImg = ""
				Imagery.FeaturedImg = ""
				Imagery.Banner = ""
			}
		}
		temp.Imagery.Thumbnail = Imagery.Thumbnail
		temp.Imagery.Backdrop = Imagery.Backdrop
		temp.Imagery.MobileImg = Imagery.MobileImg
		temp.Imagery.Banner = Imagery.Banner
		temp.ContentType = strings.ToLower(val.ContentType)
		temp.DigitalRighttype = val.DigitalRightsType
		temp.Geoblock = false
		if val.ContentTier == 1 {
			result = append(result, temp)
		} else if val.ContentTier == 2 {
			if stringInSlice(val.Id, ValidMultitierIds) {
				result = append(result, temp)
			}
		}
	}
	l.JSON(c, http.StatusOK, gin.H{"data": result})

}
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// AddUserPlaylist -  Add user playlist
// POST /v1/contents/playlist
// @Summary Add contents into user playlist
// @Description Add contents into user playlist
// @Tags User
// @security Authorization
// @Accept  json
// @Produce  json
// @Param body body PlaylistedContentRequest true "Raw JSON string"
// @Success 200
// @Router /v1/contents/playlist [post]
func (hs *HandlerService) AddUserPlaylist(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	userid := c.MustGet("userid")
	if c.MustGet("AuthorizationRequired") == 1 {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request.", "status": http.StatusUnauthorized})
		return
	}
	var input PlaylistedContent
	var request PlaylistedContentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Invalid request", "code": "invalid_request", "requestId": userid.(string)})
		return
	}
	var content Content
	if err := db.Where("content_key=?", request.Id).Find(&content).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "حدث خطأ ما", "code": "error_server_error", "requestId": userid.(string)})
		return
	}
	input.UserId = userid.(string)
	input.ContentId = content.Id
	input.AddedAt = time.Now()
	if err := db.Create(&input).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"status": http.StatusOK})
}

// Add Rating for content by user  -  Add Rating for content by user
// POST /v1/contents/rated
// @Summary   Add Rating for content by user
// @Description   Add Rating for content by user
// @Tags User
// @security Authorization
// @Accept  json
// @Produce  json
// @Param body body Addcontent true "Raw JSON string"
// @Success 200
// @Router /v1/contents/rated [post]
func (hs *HandlerService) AddRatingForContentByUser(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	dbro := c.MustGet("CDBRO").(*gorm.DB)
	serverError := common.ServerErrorResponse("en")
	if c.MustGet("AuthorizationRequired") == 1 {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	userId := c.MustGet("userid")
	Deviceid := c.MustGet("device_id")
	if userId == "" {
		serverError.Description = "Invalid Authorization"
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	badRequestError := common.BadRequestErrorResponse()
	var invalidRequestError common.InvalidRequestError
	var invalidRequest common.InvalidRequest
	var addcontent Addcontent
	var contentid ContentId
	if res := c.ShouldBindJSON(&addcontent); res != nil {
		l.JSON(c, http.StatusBadRequest, badRequestError)
		return
	}
	if addcontent.Rating <= 0 {
		invalidRequestError.Code = "NotEmptyValidator"
		invalidRequestError.Description = "'Rating' should not be empty."
		invalidRequest.Rating = &invalidRequestError
		badRequestError.Invalid = &invalidRequest
		l.JSON(c, http.StatusBadRequest, badRequestError)
		return
	}
	if addcontent.ContentRequest.Id < 1 {
		invalidRequestError.Code = "GreaterThanValidator"
		invalidRequestError.Description = "'Id' must be greater than '0'."
		invalidRequest.ContentId = &invalidRequestError
		badRequestError.Invalid = &invalidRequest
		l.JSON(c, http.StatusBadRequest, badRequestError)
		return
	}
	if addcontent.ContentRequest.ContentType == "" {
		invalidRequestError.Code = "NotEmptyValidator"
		invalidRequestError.Description = "'Content Type' should not be empty."
		invalidRequest.ContentType = &invalidRequestError
		badRequestError.Invalid = &invalidRequest
		l.JSON(c, http.StatusBadRequest, badRequestError)
		return
	}

	// var episodeid EpisodeDetailsSummary
	// var seasonid SeasonDetailsSummary

	/*Based on content-key fetch content-id */
	// if idresults := dbro.Table("episode").Where("episode_key=(?)", addcontent.ContentRequest.Id).Find(&episodeid).Error; idresults != nil {
	// 	l.JSON(c, http.StatusBadRequest, badRequestError)
	// 	return
	// }

	// if idresults := dbro.Table("season").Where("id=(?)", episodeid.SeasonId).Find(&seasonid).Error; idresults != nil {
	// 	l.JSON(c, http.StatusBadRequest, badRequestError)
	// 	return
	// }

	if idresults := dbro.Table("content c").Select("c.id as id").Where("content_key=(?)", addcontent.ContentRequest.Id).Find(&contentid).Error; idresults != nil {
		l.JSON(c, http.StatusBadRequest, badRequestError)
		return
	}

	/*collect rated_content details based on user_id & content_id*/
	var count int
	if idresults := dbro.Table("rated_content").Select("rating").Where("content_id=? and user_id=?", contentid.Id, userId).Count(&count).Error; idresults != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	ratedContent := CreateRatedContent{Rating: addcontent.Rating, RatedAt: time.Now(), ContentId: contentid.Id, UserID: userId.(string), DeviceId: Deviceid.(string), IsHidden: false}
	if count < 1 {
		/* If User not rated for Content yet then Add Rating in rated_content table by user */
		if newRecord := db.Table("rated_content").Create(&ratedContent).Error; newRecord != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
	} else {
		/* If user already Rated for content then updated rating */
		if updateRecord := db.Table("rated_content").Where("content_id=? and user_id=?", contentid.Id, userId).Update(map[string]interface{}{"rated_at": time.Now(), "rating": ratedContent.Rating, "is_hidden": false, "device_id": Deviceid.(string)}).Error; updateRecord != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
	}
	/* Add record in rated_content_history */
	addToHistory := AddRatingToHistory{Rating: addcontent.Rating, RatedAt: time.Now(), ContentId: contentid.Id, UserID: userId.(string), DeviceId: Deviceid.(string)}
	// historyJson, err := json.Marshal(addToHistory)
	// if err != nil {
	// 	fmt.Println("Cannot encode to JSON ", err)
	// }
	// var userEvent UserEvent
	// userEvent.Timestamp = time.Now()
	// userEvent.UserID = userId.(string)
	// userEvent.EventType = "rating_history_activity"
	// userEvent.Details = string(historyJson)
	// url := os.Getenv("USER_ACTIVITY_URL")
	// response, err := common.PostCurlCall("POST", url, userEvent)
	// if err != nil {
	// 	l.JSON(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// fmt.Println(response)
	if newHistory := db.Table("rated_content_history").Create(&addToHistory).Error; newHistory != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"message": newHistory.Error(), "status": http.StatusInternalServerError})
		return
	}
	/* Find most-recent average_rating based on content_id */
	var average float64
	row := dbro.Raw("select avg(rating) from rated_content where content_id =?", contentid.Id)
	row.Count(&average)
	/* update average rating on content table*/
	avgContent := ResponseContent{AverageRating: average, AverageRatingUpdatedAt: time.Now()}
	if updatecontent := db.Table("content").Where("content_key=(?)", addcontent.ContentRequest.Id).Update(ResponseContent{AverageRating: avgContent.AverageRating, AverageRatingUpdatedAt: avgContent.AverageRatingUpdatedAt}).Error; updatecontent != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
}

// Report content issue by user -  Report content issue by user
// POST /v1/contents/watching/:ckeyctype/issues
// @Summary  Report content issue by user
// @Description  Report content issue by user
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Param ckeyctype path string true "ContentKey,ContentType"
// @Param body body ViewActivityDetailRequest true "Raw JSON string"
// @Success 200
// @Router /v1/contents/watching/{ckeyctype}/issues [post]
func (hs *HandlerService) ReportContentIssue(c *gin.Context) {
	cdb := c.MustGet("CDB").(*gorm.DB)
	db := c.MustGet("UDB").(*gorm.DB)
	serverError := common.ServerErrorResponse("en")
	var viewactivitydetails ViewActivityDetailRequest
	//	var contentid ContentId
	var viewid []ViewDetails
	//	var watchdetails WatchDetails
	contentKey := c.Param("ckeyctype")
	userId := c.MustGet("userid")
	if c.MustGet("AuthorizationRequired") == 1 {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	resultvar := strings.Split(contentKey, ",")
	if len(resultvar) != 2 {
		serverError.Description = "Please provide contentkey and contenttype"
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	if res := c.ShouldBindJSON(&viewactivitydetails); res != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var viewActivityDetailss ViewActivityDetailss
	viewActivityDetailss.Description = viewactivitydetails.Description
	viewActivityDetailss.IsCommunication = viewactivitydetails.IsCommunication
	viewActivityDetailss.IsSound = viewactivitydetails.IsSound
	viewActivityDetailss.IsTranslation = viewactivitydetails.IsTranslation
	viewActivityDetailss.IsVideo = viewactivitydetails.IsVideo
	viewActivityDetailss.ReportedAt = time.Now()
	badRequestError := common.BadRequestErrorResponse()
	var invalidRequestError common.InvalidRequestError
	var invalidRequest common.InvalidRequest
	if viewactivitydetails.Description == "" {
		invalidRequestError.Code = "error_watching_issue_description_required"
		invalidRequestError.Description = "Description is required."
		invalidRequest.LastWatchPosition = &invalidRequestError
		badRequestError.Invalid = &invalidRequest
		l.JSON(c, http.StatusBadRequest, badRequestError)
		return
	}
	if len(viewactivitydetails.Description) > 255 {
		invalidRequestError.Code = "error_watching_issue_description_length_invalid"
		invalidRequestError.Description = "Length should not exceed 255 characters."
		invalidRequest.LastWatchPosition = &invalidRequestError
		badRequestError.Invalid = &invalidRequest
		l.JSON(c, http.StatusBadRequest, badRequestError)
		return
	}
	if resultvar[1] == "movie" {
		if contentresult := cdb.Table("content c").Select("va.id").Joins("left join content_variance cv on cv.content_id=c.id").Joins("left join view_activity va on va.playback_item_id=cv.playback_item_id").Where("c.content_key=(?) and va.user_id=(?)", resultvar[0], userId).Find(&viewid).Error; contentresult != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
	} else {
		if episoderesult := cdb.Table("episode e").Select("va.id").Joins("left join view_activity va on va.playback_item_id=e.playback_item_id").Where("e.episode_key=(?) and va.user_id=(?)", resultvar[0], userId).Find(&viewid).Error; episoderesult != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
	}
	for _, value := range viewid {
		var count int
		db.Table("watching_issue").Select("view_activity_id").Where("view_activity_id =(?)", value.Id).Count(&count)
		if count == 0 {
			viewActivityDetailss.ViewActivityId = value.Id
			if finalwatchresult := db.Table("watching_issue").Create(&viewActivityDetailss).Error; finalwatchresult != nil {
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}
		} else {
			if finalresult := db.Table("watching_issue").Where("view_activity_id=(?)", value.Id).Update(&viewActivityDetailss).Error; finalresult != nil {
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}
		}
	}
}

// Remove contents to user playlist -  Remove contents to user playlist
// DELETE /v1/contents/playlist/:ckeyctype
// @Summary Remove contents to user playlist
// @Description  Remove contents to user playlist
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Param ckeyctype path string true "ContentKey,ContentType"
// @Success 200
// @Router /v1/contents/playlist/{ckeyctype} [delete]
func (hs *HandlerService) RemoveContentsUserPlaylist(c *gin.Context) {
	cdb := c.MustGet("CDB").(*gorm.DB)
	contentkey := c.Param("ckeyctype")
	res := strings.Split(contentkey, ",")
	var contentdetails ContentId
	var playlistedContent PlaylistedContent
	userId := c.MustGet("userid")
	fmt.Println(userId, "............")
	if len(res) != 2 {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": "please provide contentkey and contenttype"})
	}
	ckeys, _ := strconv.Atoi(res[0])
	if result := cdb.Table("content c").Select("id").
		Joins("join playlisted_content pc on pc.content_id =c.id").
		Where("content_key=? and lower(content_type)=lower(?) and pc.user_id=?", ckeys, res[1], userId.(string)).Find(&contentdetails).Error; result != nil {
		l.JSON(c, http.StatusNotFound, gin.H{"error": "not_found", "description": "Not found", "code": "", "requestId": randstr.String(32)})
		return
	}
	//if userresult := cdb.Raw("DELETE FROM playlisted_content WHERE content_id=? and user_id=?", contentdetails.Id, userId.(string) ).Error; userresult != nil {
	if userresult := cdb.Where("content_id=? and user_id=?", contentdetails.Id, userId.(string)).Delete(&playlistedContent).Error; userresult != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "حدث خطأ ما", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
}

// GetSearchbyGenre - Get content by genre
// GET /v1/:lang/searchbyGenre
// @Summary Get content details by genre
// @Description Get content details based on by genre
// @Tags Content
// @Accept json
// @Produce json
// @Param country query string true "Country Code"
// @Param q query string true "q string"
// @Param lang path string true "Language Code"
// @Success 200
// @Router /v1/{lang}/searchbyGenre [get]
func (hs *HandlerService) GetSearchbyGenre(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		l.JSON(c, http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("CDB").(*gorm.DB)
	if c.Request.URL.Query()["q"] == nil || c.Request.URL.Query()["q"][0] == "" {
		l.JSON(c, http.StatusBadRequest, gin.H{})
		return
	}
	querystring := strings.TrimSpace(strings.ToLower(c.Request.URL.Query()["q"][0]))

	var CountryCode, language_translated_title string
	if c.Request.URL.Query()["country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	var country string
	country = strconv.Itoa(int(common.Countrys(CountryCode)))

	var finalResult []ContentDetails
	var contents []ContentDetails

	if strings.ToLower(c.Param("lang")) == "en" {
		language_translated_title = ", cpi.transliterated_title as translated_title"
	} else {
		language_translated_title = ", cpi.arabic_title as translated_title"
	}

	querydata := "select  g.id as genre_id,c.content_key, c.content_key as id, c.id as content_id, c.content_tier,  c.created_at as inserted_at, c.modified_at, lower(c.content_type) as content_type, cpi.transliterated_title as title  " + language_translated_title + " , cpi.transliterated_title as friendly_url,  min(pi1.video_content_id) as video_id ,cv.id as season_or_varience_id  from genre g  join content_genre cg on g.id = cg.genre_id join content c  on cg.content_id = c.id  join content_primary_info cpi on c.primary_info_id = cpi.id  join content_variance cv on cv.content_id =c.id  join playback_item pi1 on pi1.id =cv.playback_item_id  join content_rights cr on cr.id=pi1.rights_id  join content_translation ct on ct.id=pi1.translation_id  full outer join content_rights_country crc on crc.content_rights_id =cr.id  full outer join content_rights_plan crp on crp.rights_id =cr.id where  c.deleted_by_user_id is null   and c.status =1    and c.deleted_by_user_id is null  and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null)   and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null)  and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null)"

	if querystring != "" {
		querydata += "  and crc.country_id='" + country + "'  and (  lower(g.english_name) like '%" + querystring + "%'  or lower(g.arabic_name)  like '%" + querystring + "%'  )"
	}
	querydata += "  group by g.id,c.content_key,c.id,cpi.transliterated_title,cpi.arabic_title,cv.id union  select g.id as genre_id,c.content_key, c.content_key as id, c.id as content_id, c.content_tier,  c.created_at as inserted_at, c.modified_at, lower(c.content_type) as content_type, cpi.transliterated_title as title   " + language_translated_title + " ,  cpi.transliterated_title as friendly_url,  min(pi1.video_content_id) as video_id ,s.id as season_or_varience_id from genre g  join content_genre cg on g.id = cg.genre_id join content c  on cg.content_id = c.id  join content_primary_info cpi on c.primary_info_id = cpi.id  join season s on c.id = s.content_id  join episode e on s.id = e.season_id  join playback_item pi1 on e.playback_item_id =pi1.id  join content_rights cr on cr.id=pi1.rights_id  join content_translation ct on ct.id=pi1.translation_id  full outer join content_rights_country crc on crc.content_rights_id =cr.id  full outer join content_rights_plan crp on crp.rights_id =cr.id where  c.deleted_by_user_id is null   and c.status =1    and c.deleted_by_user_id is null  and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null)   and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null)  and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null)  and crc.country_id='" + country + "'  "

	if querystring != "" {
		querydata += "  and (  lower(g.english_name) like '%" + querystring + "%'  or lower(g.arabic_name)  like '%" + querystring + "%'  )  "
	}
	querydata += " group by g.id,c.content_key,c.id,cpi.transliterated_title,cpi.arabic_title ,s.id"

	if err := db.Raw(querydata).Find(&finalResult).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}

	for _, details := range finalResult {
		friendlyUrl := strings.ToLower(details.FriendlyUrl)
		details.FriendlyUrl = strings.Replace(friendlyUrl, " ", "-", -1)
		details.Geoblock = false
		ImageryDetails := make(chan ContentImageryDetails)
		if details.ContentTier == 1 {
			go OnetierImagery(details.ContentId, details.SeasonOrVarienceId, ImageryDetails)
		} else {
			go MultitierImagery(details.ContentId, details.SeasonOrVarienceId, ImageryDetails)
		}
		details.Imagery = <-ImageryDetails
		contents = append(contents, details)
	}
	l.JSON(c, http.StatusOK, gin.H{"data": contents})
	return
}

// GetSearchbyCast - Get content search by cast
// GET /v1/:lang/searchbyCast
// @Summary Get movie details by content
// @Description Get movie details based on search by cast
// @Tags Content
// @Accept json
// @Produce json
// @Param country query string true "Country Code"
// @Param q query string true "q string"
// @Param lang path string true "Language Code"
// @Success 200
// @Router /v1/{lang}/searchbyCast [get]
func (hs *HandlerService) GetSearchbyCast(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	serverError := common.ServerErrorResponse(language)
	if c.Request.URL.Query()["q"] == nil || c.Request.URL.Query()["q"][0] == "" {
		serverError.Description = "Please Provide search string to search string."
		l.JSON(c, http.StatusBadRequest, serverError)
		return
	}
	querystring := (strings.ToLower(c.Request.URL.Query()["q"][0]))

	var CountryCode string
	if c.Request.URL.Query()["country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}
	var country string
	country = strconv.Itoa(int(common.Countrys(CountryCode)))

	var finalResult []ContentDetails
	var contents []ContentDetails
	var language_translated_title string
	if language == "en" {
		language_translated_title = ", cpi.transliterated_title as translated_title"
	} else {
		language_translated_title = ", cpi.arabic_title as translated_title"
	}

	var stringtext string
	stringtext = " and ( lower(a.english_name)   ='" + querystring + "' or lower(a.arabic_name) ='" + querystring + "'"
	// if len(strings.Split(querystring, " ")) > 1 {
	// 	for _, stringArray := range strings.Split(querystring, " ") {
	// 		stringtext += " or ( lower(a.english_name)  like '%" + stringArray + "%' or lower(a.arabic_name) like  '%" + stringArray + "%' ) "
	// 	}
	// }
	stringtext += " ) "
	querydata := "select distinct on (c.content_key)c.content_key,c.id as content_id,c.content_tier,c.created_at as inserted_at,c.modified_at,cpi.transliterated_title as title " + language_translated_title + " ,cpi.transliterated_title as friendly_url,lower(c.content_type) as content_type,cr.digital_rights_type,min(pi1.video_content_id) as video_id,cv.id as season_or_varience_id from actor a join content_actor ca on a.id = ca.actor_id join content c on c.cast_id = ca.cast_id left join content_primary_info cpi on cpi.id = c.primary_info_id left join content_variance cv on cv.content_id = c.id left join playback_item pi1 on pi1.id = cv. playback_item_id left join content_rights cr on cr.id = pi1.rights_id left join content_translation ct on ct.id = pi1.translation_id left join content_rights_country crc on crc.content_rights_id = cr.id left join content_rights_plan crp on crp.rights_id = cr.id where c.deleted_by_user_id is null and c.status = 1 and ( pi1.scheduling_date_time <= NOW()	 or pi1.scheduling_date_time is null) and  (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)  and(cr.digital_rights_end_date >= NOW()	 or cr.digital_rights_end_date is null) and crc.country_id = '" + country + "' "
	querydata += stringtext
	querydata += " group by a.id,c.id,cpi.transliterated_title,cpi.arabic_title,cr.digital_rights_type,cv.id  union  select distinct on (c.content_key)c.content_key, c.id as content_id,c.content_tier,c.created_at as inserted_at,c.modified_at, cpi.transliterated_title as title  " + language_translated_title + "  ,  cpi.transliterated_title as friendly_url,lower(c.content_type) as content_type,cr.digital_rights_type,min(pi1.video_content_id) as video_id,cv.id as season_or_varience_id from actor a left join content_cast cc on (a.id = cc.main_actor_id or a.id =cc.main_actress_id) left join content c on cc.id = c.cast_id left join content_primary_info cpi on cpi.id=c.primary_info_id  left join content_variance cv on cv.content_id =c.id left join playback_item pi1 on pi1.id =cv.playback_item_id  left join content_rights cr on cr.id=pi1.rights_id left join content_translation ct on ct.id=pi1.translation_id  full outer join content_rights_country crc on crc.content_rights_id =cr.id  full outer join content_rights_plan crp on crp.rights_id =cr.id  where 	c.deleted_by_user_id is null 	and c.status =1	 and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) 	and  (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null)  and  (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null)  and  crc.country_id='" + country + "' "
	querydata += stringtext
	querydata += " group by a.id,c.id,cpi.transliterated_title,cpi.arabic_title,cr.digital_rights_type,cv.id  union select distinct on (c.content_key)c.content_key, c.id as content_id,c.content_tier,c.created_at as inserted_at,c.modified_at, cpi.transliterated_title as title " + language_translated_title + " ,  cpi.transliterated_title as friendly_url,lower(c.content_type) as content_type, cr.digital_rights_type, min(pi1.video_content_id) as video_id,s.id as season_or_varience_id  from actor a  left join content_cast cc on (a.id = cc.main_actor_id or a.id =cc.main_actress_id)   left join content c on cc.id = c.cast_id  left join content_primary_info cpi on cpi.id=c.primary_info_id left join about_the_content_info atci on atci.id=c.about_the_content_info_id  left join season s on c.id = s.content_id left join episode e on s.id = e.season_id  left join playback_item pi1 on e.playback_item_id =pi1.id  left join content_rights cr on cr.id=pi1.rights_id  left join content_translation ct on ct.id=pi1.translation_id  full outer join content_rights_country crc on crc.content_rights_id =cr.id  full outer join content_rights_plan crp on crp.rights_id =cr.id  where 	c.deleted_by_user_id is null 	and c.status =1	 and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) 	 and  (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null)  and  (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null)  and  crc.country_id='" + country + "' "
	querydata += stringtext
	querydata += " group by a.id,c.id,cpi.transliterated_title,cpi.arabic_title,cr.digital_rights_type,s.id  union select distinct on (c.content_key)c.content_key, s.content_id as content_id,c.content_tier,c.created_at as inserted_at,c.modified_at, cpi.transliterated_title as title" + language_translated_title + " , cpi.transliterated_title as friendly_url,lower(c.content_type) as content_type, cr.digital_rights_type, min(pi1.video_content_id) as video_id,s.id as season_or_varience_id from actor a left join content_cast cc on (a.id = cc.main_actor_id or a.id =cc.main_actress_id) left join season s on s.cast_id = cc.id left join episode e on s.id = e.season_id left join content c on c.id = s.content_id  left join content_primary_info cpi on cpi.id=c.primary_info_id left join about_the_content_info atci on atci.id=c.about_the_content_info_id left join playback_item pi1 on e.playback_item_id =pi1.id left join content_rights cr on cr.id=pi1.rights_id  left join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id  full outer join content_rights_plan crp on crp.rights_id =cr.id where s.deleted_by_user_id is null and e.status=1 and s.status =1 and c.status =1	 and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and  (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null)  and  (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and  crc.country_id='" + country + "' "
	querydata += stringtext
	querydata += "group by a.id,c.id,cpi.transliterated_title,cpi.arabic_title,cr.digital_rights_type,s.content_id,s.id union select distinct on (c.content_key)c.content_key, s.content_id as content_id,c.content_tier,c.created_at as inserted_at,c.modified_at, cpi.transliterated_title as title,cpi.transliterated_title as translated_title, cpi.transliterated_title as friendly_url,lower(c.content_type) as content_type, cr.digital_rights_type, min(pi1.video_content_id) as video_id,s.id as season_or_varience_id from actor a left join content_actor ca on ca.actor_id = a.id left join content_cast cc on cc.id=ca.cast_id  left join season s on s.cast_id = cc.id left join episode e on s.id = e.season_id left join content c on c.id = s.content_id  left join content_primary_info cpi on cpi.id=c.primary_info_id left join about_the_content_info atci on atci.id=c.about_the_content_info_id left join playback_item pi1 on e.playback_item_id =pi1.id left join content_rights cr on cr.id=pi1.rights_id left join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id full outer join content_rights_plan crp on crp.rights_id =cr.id where s.deleted_by_user_id is null 	and s.status =1 and c.status =1 and e.status=1 and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and  (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and  crc.country_id='" + country + "' "
	querydata += stringtext
	querydata += "group by a.id,c.id,cpi.transliterated_title,cpi.arabic_title,cr.digital_rights_type,s.content_id,s.id order by inserted_at DESC"
	// fmt.Println(querydata)
	if err := db.Raw(querydata).Find(&finalResult).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	for _, details := range finalResult {
		friendlyUrl := strings.ToLower(details.FriendlyUrl)
		details.FriendlyUrl = strings.Replace(friendlyUrl, " ", "-", -1)
		details.Geoblock = false
		ImageryDetails := make(chan ContentImageryDetails)
		if details.ContentTier == 1 {
			go OnetierImagery(details.ContentId, details.SeasonOrVarienceId, ImageryDetails)
		} else {
			go MultitierImagery(details.ContentId, details.SeasonOrVarienceId, ImageryDetails)
		}
		details.Imagery = <-ImageryDetails
		contents = append(contents, details)
	}
	l.JSON(c, http.StatusOK, gin.H{"data": contents})

	return
}

// RemoveWatchingHistory -  Remove watch contents to user history
// DELETE /v1/contents/watching/:ckeyctype
// @Summary Remove watching contents to user watch history
// @Description  Remove watching contents to user watch history
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Param ckeyctype path string true "ContentKey,ContentType"
// @Success 200
// @Router /v1/contents/watching/{ckeyctype} [delete]
func (hs *HandlerService) RemoveWatchingHistory(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	serverError := common.ServerErrorResponse("en")
	contentkey := c.Param("ckeyctype")
	res := strings.Split(contentkey, ",")
	if c.MustGet("AuthorizationRequired") == 1 {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	UserId := c.MustGet("userid")
	if len(res) != 2 {
		serverError.Description = "Please provide contentkey and contenttype"
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	contentKey, _ := strconv.Atoi(res[0])
	var cDetails ContentIdDetails
	query := common.GetWatchingContentDetailsQuery(contentKey)
	if err := db.Raw(query).Find(&cDetails).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var update ViewActivity
	update.IsHidden = "true"
	update.ViewedAt = time.Now()
	if err := db.Model(&update).Where("playback_item_id=? and content_id=? and user_id=?", cDetails.PlaybackItemId, cDetails.Id, UserId).Update(map[string]interface{}{
		"viewed_at":           time.Now(),
		"is_hidden":           false,
		"last_watch_position": 0,
		"watch_session_id":    update.WatchSessionId,
	}).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}

	if res[1] == "episode" || res[1] == "Episode" || res[1] == "EPISODE" {
		var historyInput ViewActivityHistory
		historyInput.ContentTypeName = "episode"
		historyInput.ContentKey = contentKey
		historyInput.UserId = UserId.(string)
		// historyInput.DeviceId = deviceId.(string)
		historyInput.LastWatchPosition = 0
		historyInput.ViewedAt = time.Now()
		if err := db.Create(&historyInput).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

	}

}

// RemoveRatedContent -  Remove Rating by user
// DELETE /v1/contents/rated/:ckeyctype
// @Summary Remove Rating by user
// @Description  Remove Rating by user
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Param ckeyctype path string true "ContentKey,ContentType"
// @Success 200
// @Router /v1/contents/rated/{ckeyctype} [delete]
func (hs *HandlerService) RemoveRatedContent(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("userid") == "" {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("CDB").(*gorm.DB)
	serverError := common.ServerErrorResponse("en")
	notFound := common.NotFoundErrorResponse()
	contentkey := c.Param("ckeyctype")
	res := strings.Split(contentkey, ",")
	UserId := c.MustGet("userid")
	if len(res) != 2 {
		serverError.Description = "Please provide contentkey and contenttype"
		l.JSON(c, http.StatusBadRequest, serverError)
		return
	}
	var ratedContentId RatedContentByUser
	contentKey, _ := strconv.Atoi(res[0])
	if ratingresult := db.Table("content c").Select("rc.id,rc.content_id").
		Joins("left join rated_content rc on rc.content_id =c.id").
		Where("c.content_key =? and rc.user_id =? and rc.is_hidden = ?", contentKey, UserId, false).Find(&ratedContentId).Error; ratingresult != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	if ratedContentId.Id == "" {
		l.JSON(c, http.StatusNotFound, notFound)
		return
	}
	var ratedContent RatedContent
	ratedContent.IsHidden = true
	fmt.Println(c.MustGet("userid"), ratedContentId.Id, ">>>>>>>>>>>>>>>>>>>>>>>")
	if err := db.Model(&ratedContent).Where("id=? and user_id=?", ratedContentId.Id, c.MustGet("userid")).Update(&ratedContent).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{})
	return
}

// @Tags Content
// @Summary Get Playlist Details By Playlist Id
// @Description Get Playlist Details By Playlist Id
// @Accept  json
// @Produce json
// @Param playlistkey path string true "Playlist Key"
// @Param lang path string true "Language Code"
// @Param country query string false "Country"
// @Param cascade query string false "Cascade"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /v1/{lang}/playlist/{playlistkey} [get]
func (hs *HandlerService) GetPlaylistDetails(c *gin.Context) {
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	serverError := common.ServerErrorResponse(language)
	db := c.MustGet("FCDB").(*gorm.DB)
	fdb := c.MustGet("DB").(*gorm.DB)
	var response PlaylistDetailsResponse
	var contents []PlaylistContent
	var countryCode string
	if c.Request.URL.Query()["country"] != nil {
		countryCode = strings.ToUpper(c.Request.URL.Query()["country"][0])
	}
	if len(countryCode) != 2 {
		countryCode = "AE"
	}
	// platformName := "web"
	PlaylistKey := c.Param("playlistkey")
	//page palylist details
	var playlist Playlist
	if err := db.Select("id,english_title,arabic_title,scheduling_start_date,scheduling_end_date,deleted_by_user_id,is_disabled,created_at,playlist_key,modified_at,playlist_type").Where("is_disabled =false and deleted_by_user_id is null and playlist_key =? and (scheduling_start_date <=now() or scheduling_start_date is null) and (scheduling_end_date >=now() or scheduling_end_date is null)", PlaylistKey).Find(&playlist).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var wg sync.WaitGroup
	pagePlaylist := MenuPlaylistsDetails{}
	pagePlaylist.ID = playlist.PlaylistKey
	pagePlaylist.PlaylistType = playlist.PlaylistType

	// pagePlaylist.Content = []PlaylistContent{}
	pagePlaylist.PageContent = []menu.PageContent{}
	pagePlaylist.Title = playlist.EnglishTitle
	if language != "en" {
		pagePlaylist.Title = playlist.ArabicTitle
	}
	if playlist.PlaylistType == "pagecontent" {
		pagePlaylist.Title = "null"
		playlistPages := make(chan []menu.PageContent)
		wg.Add(1)
		go menu.PlaylistPages(playlist.ID, language, playlistPages, c)
		defer wg.Done()
		if playlistPages != nil {
			pagePlaylist.PageContent = <-playlistPages
		}
	} else if playlist.PlaylistType == "content" {
		contentIds, err := menu.PlaylistItemContentsJ(playlist.ID, c)
		fmt.Println(contentIds, "lll")
		contentId, _ := json.Marshal(contentIds)
		var contentIdss []string
		json.Unmarshal(contentId, &contentIdss)
		if err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		// var Ids []string
		// for _, content := range contentIds {
		// 	Ids = append(Ids, content.ContentId)
		// }
		// var playlistContents []PlaylistContent
		// type ContentFragmentDetails struct {
		// 	Details string `json:"details"`
		// }
		// var contentFragmentDetails []ContentFragmentDetails
		// if err := fdb.Debug().Raw("select jsonb_agg(details) as cnt from content_fragment cf where cf.content_id in (?) and language in (?)  and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(country,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(platform,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and rights_start_date <= now() and rights_end_date >= now()", contentIds, language, contentIds, countryId, countryId, contentIds, platformId, platformId).Find(&content1).Error; err != nil {
		// if err := fdb.Table("content_fragment").Select("details::text as details").Where("content_id in(?) and country=? and language=? and platform=?", Ids, countryCode, language, platformName).Find(&contentFragmentDetails).Error; err != nil {
		// 	l.JSON(c, http.StatusInternalServerError, serverError)
		// 	return
		// }

		var content1 Cont
		if contentError := fdb.Debug().Raw("select jsonb_agg(details) as cnt from content_fragment cf where cf.content_id in (?) and language in (?)  and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(country,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and (select cid from (select distinct unnest(translate::text[]) as cid from (select (translate(platform,'[]','{}')) from content_fragment cf where cf.content_id in (?)) as ss )aa where cid= ? )= ? and rights_start_date <= now() and rights_end_date >= now()", contentIdss, language, contentIdss, "356", "356", contentIdss, "1", "1").Find(&content1).Error; contentError != nil {
			fmt.Println("Fetch content Errors", contentError)
		}
		pagePlaylist.Content = content1.Cnt
		// for _, content := range contentFragmentDetails {
		// 	var playlistContent PlaylistContent
		// 	if err := json.Unmarshal([]byte(content.Details), &playlistContent); err != nil {
		// 		l.JSON(c, http.StatusInternalServerError, err.Error())
		// 		return
		// 	}
		// 	if playlistContent.ID != 0 {
		// 		playlistContents = append(playlistContents, playlistContent)
		// 	}
		// }
		// for _, id := range contentIds {
		// 	for _, content := range playlistContents {
		// 		if id.ContentId == content.ContentId {
		// 			contents = append(contents, content)
		// 		}
		// 	}
		// }
		// pagePlaylist.Content = contents
	}
	response.Total = len(contents)
	response.PerPage = 50
	response.CurrentPage = 1
	response.LastPage = 1
	response.NextPageUrl = nil
	response.PrevPageUrl = nil
	response.From = 1
	response.To = len(contents)
	response.Data = pagePlaylist
	l.JSON(c, http.StatusOK, response)
	return
}

/* for getting response and preparing redis key for movie */
func (hs *HandlerService) PrepareMovieDetailsByContent(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "en" {
		language = "ar"
	}
	//notFoundError := common.NotFoundErrorResponse()
	serverError := common.ServerErrorResponse(language)
	if c.Request.URL.Query()["contentkey"] == nil || c.Request.URL.Query()["contentkey"][0] == "" {
		serverError.Description = "Please Provide Content Key."
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	ContentKey := c.Request.URL.Query()["contentkey"][0]
	var CountryCode string
	if c.Request.URL.Query()["Country"] != nil {
		CountryCode = strings.ToUpper(c.Request.URL.Query()["Country"][0])
	}
	if len(CountryCode) != 2 {
		CountryCode = "AE"
	}
	country := int(common.Countrys(CountryCode))
	var content Content
	if err := db.Where("content_key=? and content_tier=1", ContentKey).Find(&content).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	ContentDetails := make(chan PlaylistContent)
	go OneTierContentDetails(content.Id, language, country, c, ContentDetails)
	details := <-ContentDetails
	fmt.Println("details not found")
	if details.ID == 0 {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": details})
	return
}

func (hs *HandlerService) contentFragments(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var contentFragment []ContentFragment1
	rows := db.Table("content_fragment").Scan(&contentFragment)
	fmt.Println("-=-=-", rows.RowsAffected)
	var resultMap map[string]interface{}

	for _, r := range contentFragment {
		// fmt.Println("-=-=-=-" , r.Details)
		// Unmarshal the JSON data into the map
		detailsBytes, err := r.Details.MarshalJSON()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if detailsBytes == nil {
			fmt.Println("rrrr", r.Id)
			continue
		}
		fmt.Println("rrrr", r.Id)
		// Unmarshal the JSON data into the map
		err = json.Unmarshal(detailsBytes, &resultMap) // Note the use of "&" to get the address of resultMap
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if resultMap == nil {
			fmt.Println("rrrr", r.Id)
			fmt.Println("Error: Unmarshaling resulted in a nil map.")
			continue
		}
		var ok bool
		bannerImg, ok := resultMap["imagery"].(map[string]interface{})["banner"].(string)
		if ok {
			removeC := strings.Split(bannerImg, "?c")
			newBanner := removeC[0]
			resultMap["imagery"].(map[string]interface{})["banner"] = newBanner
			fmt.Println("bbb", newBanner)
		}
		var ok1 bool
		backdropImg, ok1 := resultMap["imagery"].(map[string]interface{})["backdrop"].(string)
		if ok1 {
			removeC := strings.Split(backdropImg, "?c")
			newbackdropImg := removeC[0]
			resultMap["imagery"].(map[string]interface{})["backdrop"] = newbackdropImg
			fmt.Println("newbackdropImg", newbackdropImg)
		}
		var ok2 bool
		thumbnailImg, ok2 := resultMap["imagery"].(map[string]interface{})["thumbnail"].(string)
		if ok2 {
			removeC := strings.Split(thumbnailImg, "?c")
			newthumbnailImg := removeC[0]
			resultMap["imagery"].(map[string]interface{})["thumbnail"] = newthumbnailImg
			fmt.Println("bbb", newthumbnailImg)
		}
		var ok3 bool
		mobile_img, ok3 := resultMap["imagery"].(map[string]interface{})["mobile_img"].(string)
		if ok3 {
			removeC := strings.Split(mobile_img, "?c")
			newmobile_img := removeC[0]
			resultMap["imagery"].(map[string]interface{})["mobile_img"] = newmobile_img
			fmt.Println("bbb", newmobile_img)
		}
		var ok4 bool
		featured_img, ok4 := resultMap["imagery"].(map[string]interface{})["featured_img"].(string)
		if ok4 {
			removeC := strings.Split(featured_img, "?c")
			newfeatured_img := removeC[0]
			fmt.Println("bbb", newfeatured_img)
			resultMap["imagery"].(map[string]interface{})["featured_img"] = newfeatured_img
		}
		updatedDetailsBytes, err := json.Marshal(resultMap)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// Update the r.Details field with the new JSON data
		r.Details = postgres.Jsonb{RawMessage: updatedDetailsBytes}
		rows := db.Table("content_fragment").Where("id = ?", r.Id).Update(&r)
		fmt.Println("=-=-=-", rows.RowsAffected)
	}
}
