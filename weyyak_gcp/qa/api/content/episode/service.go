package episode

import (
	"bytes"
	"content/common"
	"content/fragments"
	l "content/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"github.com/thanhpk/randstr"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"time"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	srg := r.Group("/api")
	srg.Use(common.ValidateToken())
	srg.GET("/episodes", hs.GetEpisodeBasedOnSeasonId)
	srg.GET("/seasons", hs.GetListofSeasonsFromContentId)
	srg.POST("/episodes/published", hs.CreateOrUpdateEpisodes)
	srg.POST("/episodes/draft", hs.CreateOrUpdateDraftEpisode)
	srg.POST("/episodes/published/:id", hs.CreateOrUpdateEpisodes)
	srg.POST("/episodes/draft/:id", hs.CreateOrUpdateDraftEpisode)
	srg.GET("/seasons/:id", hs.GetSeasonDetailsBySeasonId)
	srg.GET("/contents/onetier/:id", hs.GetOneTierContentDetailsBasedonContentID)

}

// GetEpisodeBasedOnSeasonId -Get Episodes Based on Season ID
// GET /api/episodes
// @Summary Get Episodes Based on Season ID
// @Description Get Episodes Based on Season ID
// @Tags Episode
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param seasonId path string false "seasonId"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Success 200 {array} object c.JSON
// @Router /api/episodes [get]
func (hs *HandlerService) GetEpisodeBasedOnSeasonId(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	db := c.MustGet("DB").(*gorm.DB)
	userdb := c.MustGet("UDB").(*gorm.DB)
	var limit, offset int64
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	language := strings.ToLower(c.Param("lang"))
	if language != "ar" {
		language = "en"
	}
	serverError := common.ServerErrorResponse()
	if limit == 0 {
		limit = 10
	}
	var seasonId string
	var episodelist EpisodeDetailsSummary
	var finalEpisodesResult []FinalEpisodesResult
	finalResult := []EpisodeDetailsSummary{}
	if c.Request.URL.Query()["seasonId"] != nil {
		seasonId = strings.ToLower(c.Request.URL.Query()["seasonId"][0])
	}
	var totalcount int
	db.Debug().Raw("select count(*) from season s join episode e on e.season_id =s.id where s.id='" + seasonId + "' and e.deleted_by_user_id is null").Count(&totalcount)
	if totalcount > 0 {
		if err := db.Debug().Table("season as s").
			Select("s.status as season_status,e.id ,case when e.primary_info_id is null then false else true end as is_primary,s.content_id,e.episode_key,e.season_id ,e.status,case when s.status=1  then true else false end as status_can_be_changed,e.status as sub_status ,e.number,e.synopsis_english ,e.synopsis_arabic ,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.intro_start ,cpi.outro_start,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date,ds.name as sub_status_name,pi2.video_content_id,pi2.scheduling_date_time,s.created_by_user_id").
			Joins("left join episode e on e.season_id =s.id").
			Joins("left join content_primary_info cpi on cpi.id =e.primary_info_id").
			Joins("left join content_translation ct on ct.id =s.translation_id").
			Joins("left join content_rights cr on cr.id =s.rights_id").
			Joins("left join display_status ds on ds.id =e.status").
			Joins("left join playback_item pi2 on pi2 .id =e.playback_item_id").
			Where("s.id = ? and e.deleted_by_user_id is null", seasonId).Order("e.number asc").
			Limit(limit).Offset(offset).
			Find(&finalEpisodesResult).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		userId := finalEpisodesResult[0].CreatedByUserId
		var userInfo UserInfo
		userdb.Debug().Table("user").Select("user_name").Where("id=?", userId).Scan(&userInfo)

		createdBy := userInfo.UserName
		for _, eplist := range finalEpisodesResult {
			episodelist.IsPrimary = eplist.IsPrimary
			episodelist.UserId = "00000000-0000-0000-0000-000000000000"
			episodelist.SecondarySeasonId = "00000000-0000-0000-0000-000000000000"
			episodelist.VarianceIds = eplist.VarianceIds
			episodelist.EpisodeIds = eplist.EpisodeIds
			episodelist.SecondaryEpisodeId = "00000000-0000-0000-0000-000000000000"
			episodelist.ContentId = eplist.ContentId
			episodelist.EpisodeKey = eplist.EpisodeKey
			episodelist.SeasonId = eplist.SeasonId
			if eplist.SeasonStatus == 2 {
				newdate := eplist.DigitalRightsEndDate
				anotherdate := newdate.Format("2006-01-02")
				startsdate := eplist.DigitalRightsStartDate
				newdates := startsdate.Format("2006-01-02")
				if anotherdate < time.Now().Format("2006-01-02") {
					status := "Digital Rights Exceeded"
					episodelist.Status = 2
					episodelist.SubStatusName = status
					episodelist.StatusCanBeChanged = false
				} else if newdates > time.Now().Format("2006-01-02") {
					status := "Unpublished"
					episodelist.Status = 2
					episodelist.SubStatusName = status
					episodelist.StatusCanBeChanged = false
				} else {
					if episodelist.Status == 2 {
						status := "Unpublished"
						episodelist.SubStatusName = status
					} else if episodelist.Status == 1 {
						status := "Published"
						episodelist.SubStatusName = status
					} else if episodelist.Status == 3 {
						status := "Draft"
						episodelist.SubStatusName = status
					}
					episodelist.Status = 2
					episodelist.StatusCanBeChanged = false
				}
			} else if eplist.Status == 1 {
				newdate := eplist.DigitalRightsEndDate
				anotherdate := newdate.Format("2006-01-02")
				startsdate := eplist.DigitalRightsStartDate
				newdates := startsdate.Format("2006-01-02")
				if anotherdate < time.Now().Format("2006-01-02") {
					status := "Digital Rights Exceeded"
					episodelist.Status = 2
					episodelist.SubStatusName = status
					episodelist.StatusCanBeChanged = false
				} else if newdates > time.Now().Format("2006-01-02") {
					status := "Unpublished"
					episodelist.Status = 2
					episodelist.SubStatusName = status
					episodelist.StatusCanBeChanged = false
				} else {
					status := "Published"
					episodelist.SubStatusName = status
					episodelist.Status = 1
					episodelist.StatusCanBeChanged = true
				}
			} else if eplist.Status == 3 {
				status := "Draft"
				episodelist.SubStatusName = status
				episodelist.Status = 2
				episodelist.StatusCanBeChanged = false
			} else if eplist.Status == 2 {
				status := "Unpublished"
				episodelist.Status = 2
				episodelist.SubStatusName = status
				episodelist.StatusCanBeChanged = true
			}

			// if episodelist.DigitalRightsStartDate >= time.Now() || eplist.DigitalRightsEndDate <= time.Now() {
			// 	episodelist.StatusCanBeChanged = false
			// 	episodelist.Status = 2
			// 	subStatus := "Digital Rights Exceeded"
			// 	episodelist.SubStatusName = subStatus
			// }
			episodelist.SubStatus = eplist.Status
			// episodelist.SubStatusName = eplist.SubStatusName
			episodelist.DigitalRightsType = eplist.DigitalRightsType
			episodelist.DigitalRightsStartDate = eplist.DigitalRightsStartDate
			episodelist.DigitalRightsEndDate = eplist.DigitalRightsEndDate
			episodelist.CreatedBy = createdBy
			episodelist.PrimaryInfo.Number = eplist.Number
			episodelist.PrimaryInfo.VideoContentId = eplist.VideoContentId
			episodelist.PrimaryInfo.SynopsisEnglish = eplist.SynopsisEnglish
			episodelist.PrimaryInfo.SynopsisArabic = eplist.SynopsisArabic
			episodelist.PrimaryInfo.OriginalTitle = eplist.OriginalTitle
			episodelist.PrimaryInfo.AlternativeTitle = eplist.AlternativeTitle
			episodelist.PrimaryInfo.ArabicTitle = eplist.ArabicTitle
			episodelist.PrimaryInfo.TransliteratedTitle = eplist.TransliteratedTitle
			episodelist.PrimaryInfo.Notes = eplist.Notes
			episodelist.PrimaryInfo.IntroStart = eplist.IntroStart
			episodelist.PrimaryInfo.OutroStart = eplist.OutroStart
			episodelist.Cast = eplist.Cast
			episodelist.Music = eplist.Music
			episodelist.TagInfo = eplist.TagInfo
			episodelist.NonTextualData = eplist.NonTextualData
			episodelist.Translation.LanguageType = LanguageOriginTypes(eplist.LanguageType)
			episodelist.Translation.DubbingLanguage = eplist.DubbingLanguage
			episodelist.Translation.DubbingDialectId = eplist.DubbingDialectId
			episodelist.Translation.SubtitlingLanguage = eplist.SubtitlingLanguage
			episodelist.SchedulingDateTime = eplist.SchedulingDateTime
			episodelist.PublishingPlatforms = eplist.PublishingPlatforms
			episodelist.SeoDetails = eplist.SeoDetails
			episodelist.Id = eplist.Id
			finalResult = append(finalResult, episodelist)

		}
	}
	pages := map[string]int{
		"size":   totalcount,
		"offset": int(offset),
		"limit":  int(limit),
	}
	l.JSON(c, http.StatusOK, gin.H{"pagination": pages, "data": finalResult})
}

// GetListofSeasonsFromContentId -Get List of Seasons from content id
// GET /api/seasons
// @Summary Get List of Seasons from content id
// @Description Get List of Seasons from content id
// @Tags Season
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param ContentId path string false "Content Id"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Success 200 {array} object c.JSON
// @Router /api/seasons [get]
func (hs *HandlerService) GetListofSeasonsFromContentId(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	db := c.MustGet("DB").(*gorm.DB)
	userdb := c.MustGet("UDB").(*gorm.DB)
	var limit, offset int64
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	language := strings.ToLower(c.Param("lang"))
	if language != "ar" {
		language = "en"
	}
	serverError := common.ServerErrorResponse()
	if limit == 0 {
		limit = 10
	}
	var contentId string
	var seasonlist SeasonDetailsSummary
	var finalSeasonResult []FinalSeasonResult
	finalResult := []SeasonDetailsSummary{}
	if c.Request.URL.Query()["contentId"] != nil {
		contentId = strings.ToLower(c.Request.URL.Query()["contentId"][0])
	}
	var count int
	db.Debug().Raw("select count(*) from content c join season s on s.content_id =c.id  where c.id='" + contentId + "' and s.deleted_by_user_id is null").Count(&count)
	if count > 0 {
		if err := db.Debug().Table("content c").Select("c.status as content_status,s.rights_id,c.created_by_user_id,s.content_id ,s.id ,s.season_key ,s.status,s.modified_at ,s.number as season_number,ds.name as sub_status_name,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title,cpi.notes ,cpi.intro_start ,cpi.outro_start,ct.language_type ,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date ,crp.subscription_plan_id,atci.intro_duration,atci.intro_start as about_intro_start,atci.outro_duration ,atci.outro_start as about_outro_start").
			Joins("left join season s on s.content_id =c.id ").
			Joins("left join display_status ds on ds.id =s.status").
			Joins("left join content_primary_info cpi on cpi.id =s.primary_info_id").
			Joins("left join content_translation ct on ct.id=s.translation_id").
			Joins("left join content_rights cr on cr.id =s.rights_id").
			Joins("left join content_rights_plan crp on crp.rights_id =s.rights_id").
			Joins("left join about_the_content_info atci on atci.id =s.about_the_content_info_id").
			Where("c.id = ? and s.deleted_by_user_id is null", contentId).
			Limit(limit).Offset(offset).
			Find(&finalSeasonResult).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var userInfo UserInfo
		userdb.Debug().Table("user").Select("user_name").Where("id=?", finalSeasonResult[0].CreatedByUserId).Find(&userInfo)

		createdBy := userInfo.UserName
		for _, SeasonResult := range finalSeasonResult {
			seasonlist.ContentId = SeasonResult.ContentId
			seasonlist.SeasonKey = SeasonResult.SeasonKey
			seasonlist.Status = SeasonResult.Status
			seasonlist.StatusCanBeChanged = SeasonResult.StatusCanBeChanged
			seasonlist.SubStatusName = SeasonResult.SubStatusName
			if SeasonResult.ContentStatus == 2 {
				if SeasonResult.Status == 2 {
					status := "Unpublished"
					seasonlist.SubStatusName = status
				} else if SeasonResult.Status == 1 {
					status := "Published"
					seasonlist.SubStatusName = status
				} else if SeasonResult.Status == 3 {
					status := "Draft"
					seasonlist.SubStatusName = status
				}
				//	status := "Unpublished"
				seasonlist.Status = 2
				//	sDetails.SubStatusName = &status
				seasonlist.StatusCanBeChanged = false
			} else if SeasonResult.Status == 1 {
				newdate := SeasonResult.DigitalRightsEndDate
				anotherdate := newdate.Format("2006-01-02")
				startsdate := SeasonResult.DigitalRightsStartDate
				newdates := startsdate.Format("2006-01-02")
				if anotherdate < time.Now().Format("2006-01-02") {
					status := "Digital Rights Exceeded"
					seasonlist.Status = 2
					seasonlist.SubStatusName = status
					seasonlist.StatusCanBeChanged = false
				} else if newdates > time.Now().Format("2006-01-02") {
					status := "Unpublished"
					seasonlist.Status = 2
					seasonlist.SubStatusName = status
					seasonlist.StatusCanBeChanged = false
				} else {
					status := "Published"
					seasonlist.Status = SeasonResult.Status
					seasonlist.SubStatusName = status
					seasonlist.StatusCanBeChanged = true
				}
			} else if SeasonResult.Status == 3 {
				status := "Draft"
				seasonlist.SubStatusName = status
				seasonlist.Status = 2
				seasonlist.StatusCanBeChanged = false
			} else if SeasonResult.Status == 2 {
				status := "Unpublished"
				seasonlist.Status = 2
				seasonlist.SubStatusName = status
				seasonlist.StatusCanBeChanged = true
			}

			seasonlist.ModifiedAt = SeasonResult.ModifiedAt
			seasonlist.PrimaryInfo.SeasonNumber = SeasonResult.SeasonNumber
			seasonlist.PrimaryInfo.OriginalTitle = SeasonResult.OriginalTitle
			seasonlist.PrimaryInfo.AlternativeTitle = SeasonResult.AlternativeTitle
			seasonlist.PrimaryInfo.ArabicTitle = SeasonResult.ArabicTitle
			seasonlist.PrimaryInfo.TransliteratedTitle = SeasonResult.TransliteratedTitle
			seasonlist.PrimaryInfo.Notes = SeasonResult.Notes
			seasonlist.PrimaryInfo.IntroStart = SeasonResult.IntroStart
			seasonlist.PrimaryInfo.OutroStart = SeasonResult.OutroStart
			seasonlist.Cast = SeasonResult.Cast
			seasonlist.Music = SeasonResult.Music
			seasonlist.TagInfo = SeasonResult.TagInfo
			seasonlist.SeasonGenres = SeasonResult.SeasonGenres
			seasonlist.AboutTheContent = SeasonResult.AboutTheContent
			seasonlist.Translation.LanguageType = LanguageOriginTypes(SeasonResult.LanguageType)
			seasonlist.Translation.DubbingLanguage = SeasonResult.DubbingLanguage
			seasonlist.Translation.DubbingDialectId = SeasonResult.DubbingDialectId
			seasonlist.Translation.SubtitlingLanguage = SeasonResult.SubtitlingLanguage
			seasonlist.Episodes = SeasonResult.Episodes
			seasonlist.NonTextualData = SeasonResult.NonTextualData
			seasonlist.Rights.DigitalRightsType = SeasonResult.DigitalRightsType
			seasonlist.Rights.DigitalRightsStartDate = SeasonResult.DigitalRightsStartDate
			seasonlist.Rights.DigitalRightsEndDate = SeasonResult.DigitalRightsEndDate
			//Fetch Regions
			var digitalRightsRegions []DigitalRightsRegions
			db.Debug().Table("content_rights_country").Select("country_id").Where("content_rights_id=?", SeasonResult.RightsId).Scan(&digitalRightsRegions)
			var cIds []int
			for _, idarr := range digitalRightsRegions {
				cIds = append(cIds, idarr.CountryId)
			}
			seasonlist.Rights.DigitalRightsRegions = cIds
			if SeasonResult.SubscriptionPlans < 1 {
				buffer := make([]int, 0)
				seasonlist.Rights.SubscriptionPlans = buffer
			}
			seasonlist.CreatedBy = createdBy
			seasonlist.IntroDuration = SeasonResult.IntroDuration
			seasonlist.IntroStart = *SeasonResult.AboutIntroStart
			seasonlist.OutroDuration = SeasonResult.OutroDuration
			seasonlist.OutroStart = *SeasonResult.AboutOutroStart
			seasonlist.Products = SeasonResult.Products
			seasonlist.SeoDetails = SeasonResult.SeoDetails
			seasonlist.VarianceTrailers = SeasonResult.VarianceTrailers
			seasonlist.Id = SeasonResult.Id
			finalResult = append(finalResult, seasonlist)
		}
	}
	pages := map[string]int{
		"size":   count,
		"offset": int(offset),
		"limit":  int(limit),
	}
	l.JSON(c, http.StatusOK, gin.H{"pagination": pages, "data": finalResult})
}

// CreateOrUpdateEpisodes - Create or Update episodes
// POST /api/episodes/published/:id
// @Summary Create or Update episodes
// @Description Create or Update episodes
// @Tags Page
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param id path string true "Id"
// @Param body body Episodes true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /api/episodes/published/{id} [post]
func (hs *HandlerService) CreateOrUpdateEpisodes(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }

	var episodes Episodes
	// for redis
	var contentkey common.ContentKey
	var contenttype common.ContentType

	// decoder := json.NewDecoder(c.Request.Body)
	// decoder.Decode(&episodes)
	userId := c.MustGet("userid")
	db := c.MustGet("DB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "ar" {
		language = "en"
	}
	serverError := common.ServerErrorResponse()
	c.ShouldBindJSON(&episodes)

	if episodes.PrimaryInfo == nil || episodes.Rights == nil || episodes.Cast == nil || episodes.Music == nil || episodes.TagInfo == nil || episodes.NonTextualData == nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	// for _, cvdata := range request.Rights.SubscriptionPlans {
	if int32(episodes.Rights.DigitalRightsType) != common.ContentRightsTypes("Avod") {
		if len(episodes.Rights.SubscriptionPlans) > 1 {
			l.JSON(c, http.StatusBadRequest, common.ServerError{
				Error:       "Multiple Subscription Plans not allowed",
				Description: "If multiple subscription plans have been assigned to the season content id " + *&episodes.SeasonId + ", only one plan can be selected",
				Code:        "",
				RequestId:   randstr.String(32),
			})
			return
		}
	}
	var duplicatesEpisodeCheck int

	if checkEpisodeDuplicate := db.Debug().Table("episode").Where("season_id = ? and number = ? and deleted_by_user_id is null", episodes.SeasonId, episodes.PrimaryInfo.Number).Count(&duplicatesEpisodeCheck).Error; checkEpisodeDuplicate != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}

	_, _, duration := common.GetVideoDuration(episodes.PrimaryInfo.VideoContentId)
	if duration == 0 {
		l.JSON(c, 400, gin.H{
			"error":       "Invalid Content ContentId",
			"description": episodes.PrimaryInfo.VideoContentId + " Content Id is wrong, Please provide valid Video ContentId",
			"code":        "",
			"requestId":   randstr.String(32),
		})
		return
	}

	if *episodes.Cast.MainActorId == "" {
		l.JSON(c, http.StatusBadRequest, common.FinalErrorResponseepisode{
			Error:       "invalid_request",
			Description: "Main Actor is required field",
			Code:        "error_validation_failed",
			RequestId:   randstr.String(32)})
		return
	}

	if *episodes.Cast.MainActressId == "" {
		l.JSON(c, http.StatusBadRequest, common.FinalErrorResponseepisode{
			Error:       "invalid_request",
			Description: "Main Actress is required field",
			Code:        "error_validation_failed",
			RequestId:   randstr.String(32)})
		return
	}

	// if duplicatesEpisodeCheck > 0 {
	// 	serverError.Description = "Episode Already Created"
	// 	l.JSON(c, http.StatusInternalServerError, serverError)
	// 	return
	// }

	// var errorFlag bool
	// errorFlag = false
	// var primaryInfoError common.PrimaryInfoError

	// if req.PrimaryInfo == nil {
	// 	errorFlag = true
	// 	primaryInfoError = common.PrimaryInfoError{"NotEmptyValidator", "'Primary Info' should not be empty."}
	// 	fmt.Println(primaryInfoError, ",,,,,,,,,,,,")
	// }

	// var rightserror common.RigthsError
	// if req.Rights == nil {
	// 	errorFlag = true
	// 	rightserror = common.RigthsError{"NotEmptyValidator", "Rights' should not be empty."}
	// }

	// var casterror common.CastError
	// if req.Cast == nil {
	// 	errorFlag = true
	// 	casterror = common.CastError{"NotEmptyValidator", "'Cast' should not be empty."}
	// }
	// var musicError common.MusicError
	// if req.Music == nil {
	// 	errorFlag = true
	// 	musicError = common.MusicError{"NotEmptyValidator", "'Music' should not be empty."}
	// }
	// var taginfoError common.TaginfoError
	// if req.TagInfo == nil {
	// 	errorFlag = true
	// 	taginfoError = common.TaginfoError{"NotEmptyValidator", "'Tag Info' should not be empty."}
	// }

	// var nontextualerrror common.NonTextualDataError
	// if req.NonTextualData == nil {
	// 	errorFlag = true
	// 	nontextualerrror = common.NonTextualDataError{"NotEmptyValidator", "'Non Textual Data' must not be empty."}
	// }

	// var invalid common.Invalidsepisode
	// if primaryInfoError.Code != "" {
	// 	fmt.Println(primaryInfoError.Code)
	// 	invalid.PrimaryInfoError = primaryInfoError
	// 	fmt.Println(primaryInfoError, "................................")
	// }

	// if rightserror.Code != "" {
	// 	invalid.RightsError = rightserror
	// }

	// if casterror.Code != "" {
	// 	invalid.CastError = casterror
	// }
	// if musicError.Code != "" {
	// 	invalid.MusicError = musicError
	// }
	// if taginfoError.Code != "" {
	// 	invalid.TaginfoError = taginfoError
	// }

	// if nontextualerrror.Code != "" {
	// 	invalid.NonTextualDataError = nontextualerrror
	// }
	// var finalErrorResponse common.FinalErrorResponseepisode
	// finalErrorResponse = common.FinalErrorResponseepisode{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	// if errorFlag {
	// 	l.JSON(c, http.StatusBadRequest, finalErrorResponse)
	// 	return
	// }

	ctx := context.Background()
	tx := db.Debug().BeginTx(ctx, nil)
	var HasPosterImage, HasDubbingScript, HasSubtitlingScript bool
	HasPosterImage = false
	HasDubbingScript = false
	HasSubtitlingScript = false
	fmt.Println(episodes.NonTextualData.PosterImage, "poster")
	if episodes.NonTextualData.PosterImage != "" {
		HasPosterImage = true
	} else {
		HasPosterImage = false
	}
	if episodes.NonTextualData.DubbingScript != "" {
		HasDubbingScript = true
	} else {
		HasDubbingScript = false
	}
	if episodes.NonTextualData.SubtitlingScript != "" {
		HasSubtitlingScript = true
	} else {

		HasSubtitlingScript = false
	}

	/*Gathering Information for either create or update */
	/*content_primary_info*/
	createPrimaryInfo := CreatePrimaryInfo{OriginalTitle: episodes.PrimaryInfo.OriginalTitle, AlternativeTitle: episodes.PrimaryInfo.AlternativeTitle, ArabicTitle: episodes.PrimaryInfo.ArabicTitle, TransliteratedTitle: episodes.PrimaryInfo.TransliteratedTitle, Notes: episodes.PrimaryInfo.Notes, IntroStart: episodes.PrimaryInfo.IntroStart, OutroStart: episodes.PrimaryInfo.OutroStart}
	/*content_rights*/
	contentRights := ContentRights{DigitalRightsType: episodes.Rights.DigitalRightsType, DigitalRightsStartDate: episodes.Rights.DigitalRightsStartDate, DigitalRightsEndDate: episodes.Rights.DigitalRightsEndDate}
	/*content_cast*/
	insertCast := InsertCast{MainActorId: *episodes.Cast.MainActorId, MainActressId: *episodes.Cast.MainActressId}
	/*fetch translation id to insert in play_back_item */
	type Translation struct {
		TranslationId string `json:"translationId"`
	}
	var translation Translation
	if translationError := db.Table("season").Select("translation_id").Where("id=?", episodes.SeasonId).Find(&translation).Error; translationError != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	/*check if episode not exists in table then create new episode*/
	if c.Param("id") == "" {
		if episodeprimaryinfo := tx.Debug().Table("content_primary_info").Create(&createPrimaryInfo).Error; episodeprimaryinfo != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(createPrimaryInfo.Id)
		/*create content-rights information*/
		if createRights := tx.Debug().Table("content_rights").Create(&contentRights).Error; createRights != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(contentRights.Id)
		/*create content-rights-country information*/
		var createcontentRightscountry []interface{}
		for _, countries := range episodes.Rights.DigitalRightsRegions {
			contentRightsCountry := ContentRightsCountry{ContentRightsId: contentRights.Id, CountryId: countries}
			createcontentRightscountry = append(createcontentRightscountry, contentRightsCountry)
		}
		if err := gormbulk.BulkInsert(tx, createcontentRightscountry, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* create cast_id for insert actor,writer,director*/
		if createCast := tx.Debug().Table("content_cast").Create(&insertCast).Error; createCast != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(insertCast.Id)
		/* create actor for episode*/
		var createActor []interface{}
		for _, actors := range episodes.Cast.Actors {
			contentActor := ContentActor{CastId: insertCast.Id, ActorId: actors}
			createActor = append(createActor, contentActor)
		}
		if err := gormbulk.BulkInsert(tx, createActor, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create writer for episode*/
		var createWriter []interface{}
		for _, writers := range episodes.Cast.Writers {
			contentWriter := ContentWriter{CastId: insertCast.Id, WriterId: writers}
			createWriter = append(createWriter, contentWriter)
		}
		if err := gormbulk.BulkInsert(tx, createWriter, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* create director for episode*/
		var createDirector []interface{}
		for _, directors := range episodes.Cast.Directors {
			contentDirector := ContentDirector{CastId: insertCast.Id, DirectorId: directors}
			createDirector = append(createDirector, contentDirector)
		}
		if err := gormbulk.BulkInsert(tx, createDirector, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create music-id for insert singer,music-composer,songwtriter info*/
		type ContentMusic struct {
			Id string `json:"id"`
		}
		var contentMusic ContentMusic
		if createContentMusic := tx.Debug().Table("content_music").Create(&contentMusic).Error; createContentMusic != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(contentMusic.Id)
		/*create singers for episode*/
		var createSinger []interface{}
		for _, singers := range episodes.Music.Singers {
			contentSinger := ContentSinger{MusicId: contentMusic.Id, SingerId: singers}
			createSinger = append(createSinger, contentSinger)
		}
		if err := gormbulk.BulkInsert(tx, createSinger, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create music-composeer for episode*/
		var createMusicComposer []interface{}
		for _, musicComposers := range episodes.Music.MusicComposers {
			contentMusicComposer := ContentMusicComposer{MusicId: contentMusic.Id, MusicComposerId: musicComposers}
			createMusicComposer = append(createMusicComposer, contentMusicComposer)
		}
		if err := gormbulk.BulkInsert(tx, createMusicComposer, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create song-writer for episode*/
		var createSongWriters []interface{}
		for _, songWriters := range episodes.Music.SongWriters {
			ContentSongWriter := ContentSongWriter{MusicId: contentMusic.Id, SongWriterId: songWriters}
			createSongWriters = append(createSongWriters, ContentSongWriter)
		}
		if err := gormbulk.BulkInsert(tx, createSongWriters, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create content-tags-info for episode*/
		type ContentTagInfo struct {
			Id string `json:"id"`
		}
		var contentTagInfo ContentTagInfo
		if createContentTagInfo := tx.Debug().Table("content_tag_info").Create(&contentTagInfo).Error; createContentTagInfo != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(contentTagInfo.Id)
		/*create content-tags for episode*/
		var createContentTags []interface{}
		for _, contentTags := range episodes.TagInfo.Tags {
			contentTag := ContentTag{TagInfoId: contentTagInfo.Id, TextualDataTagId: contentTags}
			createContentTags = append(createContentTags, contentTag)
		}
		if err := gormbulk.BulkInsert(tx, createContentTags, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create playbackItem for episode*/
		// take created by userid from request body for creating old contents else take user id from generated token
		_, _, duration := common.GetVideoDuration(episodes.PrimaryInfo.VideoContentId)
		playbackItem := PlaybackItem{VideoContentId: episodes.PrimaryInfo.VideoContentId, SchedulingDateTime: episodes.SchedulingDateTime, Duration: duration, CreatedByUserId: userId.(string), TranslationId: translation.TranslationId, RightsId: contentRights.Id}
		if createPlaybackItem := tx.Debug().Table("playback_item").Create(&playbackItem).Error; createPlaybackItem != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(playbackItem.Id)
		/*create playback_item_target_platforms*/
		var createplaybackTargetPlatforms []interface{}
		for _, platforms := range episodes.PublishingPlatforms {
			playbackItemTargetPlatform := PlaybackItemTargetPlatform{PlaybackItemId: playbackItem.Id, TargetPlatform: platforms, RightsId: contentRights.Id}
			createplaybackTargetPlatforms = append(createplaybackTargetPlatforms, playbackItemTargetPlatform)
		}
		if err := gormbulk.BulkInsert(tx, createplaybackTargetPlatforms, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		// Create a content_rights_plan for episode level
		fmt.Println("-------Creating content_rights_plan ------")
		var contentRightsPlan ContentRightsPlan
		if len(episodes.Rights.SubscriptionPlans) > 0 {
			for _, contentplanrange := range episodes.Rights.SubscriptionPlans {
				contentRightsPlan = ContentRightsPlan{RightsId: contentRights.Id, SubscriptionPlanId: contentplanrange}
				if res := tx.Debug().Table("content_rights_plan").Create(&contentRightsPlan).Error; res != nil {
					tx.Rollback()
					l.JSON(c, http.StatusInternalServerError, nil)
					return
				}
			}
		}
		fmt.Println("-------Creating content_rights_plan ------")
		/*Create An Episode*/
		var episodeKey FetchEpisodeDetails
		if contentkeyresult := tx.Debug().Table("episode").Select("max(episode_key) as episode_key,max(third_party_episode_key) as third_party_episode_key").Find(&episodeKey).Error; contentkeyresult != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": contentkeyresult.Error(), "status": http.StatusInternalServerError})
		}
		createEpisode := CreateEpisode{
			// for removing sync below line is commented
			//	Id:                  episodes.SecondaryEpisodeId, // id for creating old episodes using .net
			SeasonId:             episodes.SeasonId,
			Number:               episodes.PrimaryInfo.Number,
			PrimaryInfoId:        createPrimaryInfo.Id,
			PlaybackItemId:       playbackItem.Id,
			Status:               1,
			SynopsisEnglish:      episodes.PrimaryInfo.SynopsisEnglish,
			SynopsisArabic:       episodes.PrimaryInfo.SynopsisArabic,
			CastId:               insertCast.Id,
			MusicId:              contentMusic.Id,
			TagInfoId:            contentTagInfo.Id,
			HasPosterImage:       HasPosterImage,
			HasDubbingScript:     HasDubbingScript,
			HasSubtitlingScript:  HasSubtitlingScript,
			EpisodeKey:           episodeKey.EpisodeKey + 1,
			ThirdPartyEpisodeKey: episodeKey.ThirdPartyEpisodeKey + 1, // for creating thirdparty episode key
			// EpisodeKey:             episodes.EpisodeKey, // for removing sync commented this line and uncommented above line
			CreatedAt:              time.Now(),
			ModifiedAt:             time.Now(),
			EnglishMetaTitle:       episodes.SeoDetails.EnglishMetaTitle,
			ArabicMetaTitle:        episodes.SeoDetails.ArabicMetaTitle,
			EnglishMetaDescription: episodes.SeoDetails.EnglishMetaDescription,
			ArabicMetaDescription:  episodes.SeoDetails.ArabicMetaDescription,
		}

		if err := tx.Debug().Table("episode").Create(&createEpisode).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var seasonDetails SeasonResult
		if seasonresult := tx.Debug().Table("season").Select("content_id").Where("id=?", episodes.SeasonId).Find(&seasonDetails).Error; seasonresult != nil {
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": seasonresult.Error(), "status": http.StatusInternalServerError})
		}
		if err := tx.Debug().Table("content").Where("id=?", seasonDetails.ContentId).Update("modified_at", time.Now()).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* update dirty count in content_sync table */
		go common.ContentSynching(seasonDetails.ContentId, c)
		/* upload images to S3 bucket based on episode Id*/
		go EpisodeFileUploadGcp(episodes, createEpisode.Id, seasonDetails.ContentId, episodes.SeasonId)
		/* content Fragment Creation */
		go fragments.CreateContentFragment(seasonDetails.ContentId, c)
		db.Debug().Raw("select content_key from content where id=?", seasonDetails.ContentId).Find(&contentkey)
		db.Debug().Raw("select content_type from content where id=?", seasonDetails.ContentId).Find(&contenttype)
		/* Prepare Redis Cache for single content*/
		contentkeyconverted := strconv.Itoa(contentkey.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(contenttype.ContentType, c)
		go common.RedisFlush(c)
		go common.GetMenu()
		go common.Pagekey()
		go common.Contenttype()
		finEpiId := map[string]string{"id": createEpisode.Id}
		l.JSON(c, http.StatusOK, gin.H{"data": finEpiId})

	} else {
		var checkEpisode int
		if epiError := db.Debug().Table("episode").Select("id").Where("id=?", c.Param("id")).Count(&checkEpisode).Error; epiError != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		if checkEpisode < 1 {
			l.JSON(c, http.StatusBadRequest, gin.H{"message": "error_season_not_found", "status": http.StatusBadRequest})
			return
		}
		/*check if episode exists in table then update episode with episodeKey*/
		var fetchEpisodeResults FetchEpisodeDetails
		if episodeResult := db.Debug().Table("episode e").Where("e.id=?", c.Param("id")).Select("e.id,e.playback_item_id,e.primary_info_id,e.cast_id,e.music_id,e.tag_info_id,pi.rights_id").
			Joins("left join playback_item pi on pi.id =e.playback_item_id").
			Scan(&fetchEpisodeResults).Error; episodeResult != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(fetchEpisodeResults)
		/*update content_primary_info*/
		// if updatePrimaryInfo := tx.Debug().Table("content_primary_info").Where("id=?", fetchEpisodeResults.PrimaryInfoId).Update(&createPrimaryInfo).Error; updatePrimaryInfo != nil {
		// 	tx.Rollback()
		// 	l.JSON(c, http.StatusInternalServerError, serverError)
		// 	return
		// }

		if primaryinfoupdate := tx.Debug().Table("content_primary_info").Where("id=?", fetchEpisodeResults.PrimaryInfoId).Update(map[string]interface{}{
			"alternative_title":    createPrimaryInfo.AlternativeTitle,
			"arabic_title":         createPrimaryInfo.ArabicTitle,
			"intro_start":          createPrimaryInfo.IntroStart,
			"notes":                createPrimaryInfo.Notes,
			"original_title":       createPrimaryInfo.OriginalTitle,
			"outro_start":          createPrimaryInfo.OutroStart,
			"transliterated_title": createPrimaryInfo.TransliteratedTitle,
		}).Error; primaryinfoupdate != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}

		/* update content_rights*/
		// if updateContentRights := tx.Debug().Table("content_rights").Where("id=?", fetchEpisodeResults.RightsId).Update(&contentRights).Error; updateContentRights != nil {
		// 	tx.Rollback()
		// 	l.JSON(c, http.StatusInternalServerError, serverError)
		// 	return
		// }
		if updateContentRights := tx.Debug().Select("Id", "DigitalRightsType", "DigitalRightsStartDate", "DigitalRightsEndDate", contentRights.Id, contentRights.DigitalRightsType, contentRights.DigitalRightsStartDate, contentRights.DigitalRightsEndDate).Create(&contentRights).Error; updateContentRights != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* update content_cast*/
		if updateContentcast := tx.Debug().Table("content_cast").Where("id=?", fetchEpisodeResults.CastId).Update(&insertCast).Error; updateContentcast != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update playback_item*/
		_, _, duration := common.GetVideoDuration(episodes.PrimaryInfo.VideoContentId)
		playbackItem := PlaybackItem{VideoContentId: episodes.PrimaryInfo.VideoContentId, Duration: duration, SchedulingDateTime: episodes.SchedulingDateTime, CreatedByUserId: userId.(string), TranslationId: translation.TranslationId, RightsId: fetchEpisodeResults.RightsId}
		//new implementation
		if episodes.SchedulingDateTime == nil {
			if res := tx.Debug().Table("playback_item").Select("scheduling_date_time").Where("id=?", fetchEpisodeResults.PlaybackItemId).Updates(map[string]interface{}{"scheduling_date_time": gorm.Expr("NULL")}).Error; res != nil {
				tx.Rollback()
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}
			if updateContentcast := tx.Debug().Table("playback_item").Where("id=?", fetchEpisodeResults.PlaybackItemId).Update(&playbackItem).Error; updateContentcast != nil {
				tx.Rollback()
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}
		} else {
			if updateContentcast := tx.Debug().Table("playback_item").Where("id=?", fetchEpisodeResults.PlaybackItemId).Update(&playbackItem).Error; updateContentcast != nil {
				tx.Rollback()
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}
		}
		if updateContentcast := tx.Debug().Exec("update playback_item pi2 set rights_id = ? where id = ?", contentRights.Id, fetchEpisodeResults.PlaybackItemId).Error; updateContentcast != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update content_rights_country*/
		var contentRightsCountry ContentRightsCountry
		// don't delete content rights here while deleting content rights,season rights is deleted due to season and episode has same rights earlier
		// generate new rights id for ep dont disturb season rights_id
		/*if err := tx.Debug().Table("content_rights_country").Where("content_rights_id=?", fetchEpisodeResults.RightsId).Delete(&contentRightsCountry).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}*/
		var createcontentRightscountry []interface{}
		for _, countries := range episodes.Rights.DigitalRightsRegions {
			contentRightsCountry = ContentRightsCountry{ContentRightsId: contentRights.Id, CountryId: countries}
			createcontentRightscountry = append(createcontentRightscountry, contentRightsCountry)
		}
		if err := gormbulk.BulkInsert(tx, createcontentRightscountry, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		// contentRightsPlan for episodes
		fmt.Println("-------episodes-----contentRightsPlan-----")
		var contentRightsPlans ContentRightsPlan
		db.Debug().Where("rights_id=?", contentRights.Id).Delete(&contentRightsPlans)
		if episodes.Rights.SubscriptionPlans != nil && len(episodes.Rights.SubscriptionPlans) > 0 {
			for _, plan := range episodes.Rights.SubscriptionPlans {
				var contentRightsPlan ContentRightsPlan
				contentRightsPlan.RightsId = contentRights.Id
				contentRightsPlan.SubscriptionPlanId = plan
				if err := tx.Debug().Create(&contentRightsPlan).Error; err != nil {
					l.JSON(c, http.StatusInternalServerError, serverError)
					// return "", serverError, 0
				}
			}
		}
		fmt.Println("-------episodes-----contentRightsPlan-----")
		/*update content_actor*/
		var contentActor ContentActor
		if err := tx.Debug().Table("content_actor").Where("cast_id=?", fetchEpisodeResults.CastId).Delete(&contentActor).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createActor []interface{}
		for _, actors := range episodes.Cast.Actors {
			contentActor := ContentActor{CastId: fetchEpisodeResults.CastId, ActorId: actors}
			createActor = append(createActor, contentActor)
		}
		if err := gormbulk.BulkInsert(tx, createActor, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update content_writer*/
		var contentWriter ContentWriter
		if err := tx.Debug().Table("content_writer").Where("cast_id=?", fetchEpisodeResults.CastId).Delete(&contentWriter).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createWriter []interface{}
		for _, writers := range episodes.Cast.Writers {
			contentWriter := ContentWriter{CastId: fetchEpisodeResults.CastId, WriterId: writers}
			createWriter = append(createWriter, contentWriter)
		}
		if err := gormbulk.BulkInsert(tx, createWriter, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update content_director*/
		var contentDirector ContentDirector
		if err := tx.Debug().Table("content_director").Where("cast_id=?", fetchEpisodeResults.CastId).Delete(&contentDirector).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createDirector []interface{}
		for _, directors := range episodes.Cast.Directors {
			contentDirector := ContentDirector{CastId: fetchEpisodeResults.CastId, DirectorId: directors}
			createDirector = append(createDirector, contentDirector)
		}
		if err := gormbulk.BulkInsert(tx, createDirector, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}

		/*update content_singers*/
		var contentSinger ContentSinger
		if err := tx.Debug().Table("content_singer").Where("music_id=?", fetchEpisodeResults.MusicId).Delete(&contentSinger).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createSinger []interface{}
		for _, singers := range episodes.Music.Singers {
			contentSinger := ContentSinger{MusicId: fetchEpisodeResults.MusicId, SingerId: singers}
			createSinger = append(createSinger, contentSinger)
		}
		if err := gormbulk.BulkInsert(tx, createSinger, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* update content_music_composer */
		var contentMusicComposer ContentMusicComposer
		if err := tx.Debug().Table("content_music_composer").Where("music_id=?", fetchEpisodeResults.MusicId).Delete(&contentMusicComposer).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createMusicComposer []interface{}
		for _, musicComposers := range episodes.Music.MusicComposers {
			contentMusicComposer := ContentMusicComposer{MusicId: fetchEpisodeResults.MusicId, MusicComposerId: musicComposers}
			createMusicComposer = append(createMusicComposer, contentMusicComposer)
		}
		if err := gormbulk.BulkInsert(tx, createMusicComposer, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* update content_song_writers */
		var contentSongWriter ContentSongWriter
		if err := tx.Debug().Table("content_song_writer").Where("music_id=?", fetchEpisodeResults.MusicId).Delete(&contentSongWriter).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createSongWriters []interface{}
		for _, songWriters := range episodes.Music.SongWriters {
			ContentSongWriter := ContentSongWriter{MusicId: fetchEpisodeResults.MusicId, SongWriterId: songWriters}
			createSongWriters = append(createSongWriters, ContentSongWriter)
		}
		if err := gormbulk.BulkInsert(tx, createSongWriters, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update content-tags*/
		var contentTags ContentTag
		if err := tx.Debug().Table("content_tag").Where("tag_info_id=?", fetchEpisodeResults.TagInfoId).Delete(&contentTags).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createContentTags []interface{}
		for _, contentTags := range episodes.TagInfo.Tags {
			contentTag := ContentTag{TagInfoId: fetchEpisodeResults.TagInfoId, TextualDataTagId: contentTags}
			createContentTags = append(createContentTags, contentTag)
		}
		if err := gormbulk.BulkInsert(tx, createContentTags, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update playback_item_target_platforms*/
		var playbackItemTargetPlatform PlaybackItemTargetPlatform
		if err := tx.Debug().Table("playback_item_target_platform").Where("playback_item_id=?", fetchEpisodeResults.PlaybackItemId).Delete(&playbackItemTargetPlatform).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createplaybackTargetPlatforms []interface{}
		for _, platforms := range episodes.PublishingPlatforms {
			playbackItemTargetPlatform := PlaybackItemTargetPlatform{PlaybackItemId: fetchEpisodeResults.PlaybackItemId, TargetPlatform: platforms, RightsId: fetchEpisodeResults.RightsId}
			createplaybackTargetPlatforms = append(createplaybackTargetPlatforms, playbackItemTargetPlatform)
		}
		if err := gormbulk.BulkInsert(tx, createplaybackTargetPlatforms, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update episode*/
		Episode := CreateEpisode{
			SeasonId:               episodes.SeasonId,
			Number:                 episodes.PrimaryInfo.Number,
			SynopsisEnglish:        episodes.PrimaryInfo.SynopsisEnglish,
			SynopsisArabic:         episodes.PrimaryInfo.SynopsisArabic,
			HasPosterImage:         HasPosterImage,
			Status:                 1,
			HasDubbingScript:       HasDubbingScript,
			HasSubtitlingScript:    HasSubtitlingScript,
			ModifiedAt:             time.Now(),
			EnglishMetaTitle:       episodes.SeoDetails.EnglishMetaTitle,
			ArabicMetaTitle:        episodes.SeoDetails.ArabicMetaTitle,
			EnglishMetaDescription: episodes.SeoDetails.EnglishMetaDescription,
			ArabicMetaDescription:  episodes.SeoDetails.ArabicMetaDescription,
		}
		if updateEpisode := tx.Debug().Table("episode").Where("id=?", fetchEpisodeResults.Id).Update(&Episode).Error; updateEpisode != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		res := map[string]string{
			"id": c.Param("id"),
		}
		var seasonDetails SeasonResult
		if seasonresult := tx.Debug().Table("season").Select("content_id").Where("id=?", episodes.SeasonId).Find(&seasonDetails).Error; seasonresult != nil {
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": seasonresult.Error(), "status": http.StatusInternalServerError})
		}
		if err := tx.Debug().Table("content").Where("id=?", seasonDetails.ContentId).Update("modified_at", time.Now()).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* upload images to S3 bucket based on episode Id*/
		go EpisodeFileUploadGcp(episodes, c.Param("id"), seasonDetails.ContentId, episodes.SeasonId)
		/* update dirty count if content in content_sync table */
		go common.ContentSynching(seasonDetails.ContentId, c)
		/* content Fragment Creation */
		go fragments.CreateContentFragment(seasonDetails.ContentId, c)
		db.Debug().Raw("select content_key from content where id=?", seasonDetails.ContentId).Find(&contentkey)
		db.Debug().Raw("select content_type from content where id=?", seasonDetails.ContentId).Find(&contenttype)
		/* Prepare Redis Cache for single content*/
		contentkeyconverted := strconv.Itoa(contentkey.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(contenttype.ContentType, c)
		go common.RedisFlush(c)
		go common.GetMenu()
		go common.Pagekey()
		go common.Contenttype()

		l.JSON(c, http.StatusOK, gin.H{"data": res})
	}
	/*commit changes*/
	if err := tx.Commit().Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}

}

// CreateOrUpdateDraftEpisode - Create or Update draft episode
// POST /api/episodes/draft/:id
// @Summary Create or Update draft episode
// @DescriptionCreate or Update draft episode
// @Tags Page
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param id path string true "Id"
// @Param body body Episodes true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /api/episodes/draft/{id} [post]
func (hs *HandlerService) CreateOrUpdateDraftEpisode(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	userId := c.MustGet("userid")
	var episodes Episodes
	decoder := json.NewDecoder(c.Request.Body)
	decoder.Decode(&episodes)
	db := c.MustGet("DB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "ar" {
		language = "en"
	}
	serverError := common.ServerErrorResponse()
	c.ShouldBindJSON(&episodes)

	if episodes.PrimaryInfo == nil || episodes.Rights == nil || episodes.Cast == nil || episodes.Music == nil || episodes.TagInfo == nil || episodes.NonTextualData == nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}

	// _, _, duration := common.GetVideoDuration(episodes.PrimaryInfo.VideoContentId)
	// if duration == 0 {
	// 	l.JSON(c, 400, gin.H{
	// 		"error":       "Invalid Content ContentId",
	// 		"description": episodes.PrimaryInfo.VideoContentId + " Content Id is wrong, Please provide valid Video ContentId",
	// 		"code":        "",
	// 		"requestId":   randstr.String(32),
	// 	})
	// 	return
	// }

	if *episodes.Cast.MainActorId == "" {
		l.JSON(c, http.StatusBadRequest, common.FinalErrorResponseepisode{
			Error:       "invalid_request",
			Description: "Main Actor is required field",
			Code:        "error_validation_failed",
			RequestId:   randstr.String(32)})
		return
	}

	if *episodes.Cast.MainActressId == "" {
		l.JSON(c, http.StatusBadRequest, common.FinalErrorResponseepisode{
			Error:       "invalid_request",
			Description: "Main Actress is required field",
			Code:        "error_validation_failed",
			RequestId:   randstr.String(32)})
		return
	}

	ctx := context.Background()
	tx := db.Debug().BeginTx(ctx, nil)
	var HasPosterImage, HasDubbingScript, HasSubtitlingScript bool
	HasPosterImage = false
	HasDubbingScript = false
	HasSubtitlingScript = false
	if episodes.NonTextualData.PosterImage != "" {
		HasPosterImage = true
	} else {
		HasPosterImage = false
	}
	if episodes.NonTextualData.DubbingScript != "" {
		HasDubbingScript = true
	} else {
		HasDubbingScript = false
	}
	if episodes.NonTextualData.SubtitlingScript != "" {
		HasSubtitlingScript = true
	} else {

		HasSubtitlingScript = false
	}

	/*Gathering Information for either create or update */
	/*content_primary_info*/
	createPrimaryInfo := CreatePrimaryInfo{OriginalTitle: episodes.PrimaryInfo.OriginalTitle, AlternativeTitle: episodes.PrimaryInfo.AlternativeTitle, ArabicTitle: episodes.PrimaryInfo.ArabicTitle, TransliteratedTitle: episodes.PrimaryInfo.TransliteratedTitle, Notes: episodes.PrimaryInfo.Notes, IntroStart: episodes.PrimaryInfo.IntroStart, OutroStart: episodes.PrimaryInfo.OutroStart}
	/*content_rights*/
	contentRights := ContentRights{DigitalRightsType: episodes.Rights.DigitalRightsType, DigitalRightsStartDate: episodes.Rights.DigitalRightsStartDate, DigitalRightsEndDate: episodes.Rights.DigitalRightsEndDate}
	/*content_cast*/
	insertCast := InsertCast{MainActorId: *episodes.Cast.MainActorId, MainActressId: *episodes.Cast.MainActressId}
	/*fetch translation id to insert in play_back_item */
	type Translation struct {
		TranslationId string `json:"translationId"`
	}
	var translation Translation
	if translationError := db.Debug().Table("season").Select("translation_id").Where("id=?", episodes.SeasonId).Find(&translation).Error; translationError != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}

	/*check if episode not exists in table then create new episode*/
	if c.Param("id") == "" {
		if episodeprimaryinfo := tx.Debug().Table("content_primary_info").Create(&createPrimaryInfo).Error; episodeprimaryinfo != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(createPrimaryInfo.Id)
		/*create content-rights information*/
		if createRights := tx.Debug().Table("content_rights").Create(&contentRights).Error; createRights != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(contentRights.Id)
		/*create content-rights-country information*/
		var createcontentRightscountry []interface{}
		for _, countries := range episodes.Rights.DigitalRightsRegions {
			contentRightsCountry := ContentRightsCountry{ContentRightsId: contentRights.Id, CountryId: countries}
			createcontentRightscountry = append(createcontentRightscountry, contentRightsCountry)
		}
		if err := gormbulk.BulkInsert(tx, createcontentRightscountry, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* create cast_id for insert actor,writer,director*/
		if createCast := tx.Debug().Table("content_cast").Create(&insertCast).Error; createCast != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(insertCast.Id)
		/* create actor for episode*/
		var createActor []interface{}
		for _, actors := range episodes.Cast.Actors {
			contentActor := ContentActor{CastId: insertCast.Id, ActorId: actors}
			createActor = append(createActor, contentActor)
		}
		if err := gormbulk.BulkInsert(tx, createActor, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create writer for episode*/
		var createWriter []interface{}
		for _, writers := range episodes.Cast.Writers {
			contentWriter := ContentWriter{CastId: insertCast.Id, WriterId: writers}
			createWriter = append(createWriter, contentWriter)
		}
		if err := gormbulk.BulkInsert(tx, createWriter, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* create director for episode*/
		var createDirector []interface{}
		for _, directors := range episodes.Cast.Directors {
			contentDirector := ContentDirector{CastId: insertCast.Id, DirectorId: directors}
			createDirector = append(createDirector, contentDirector)
		}
		if err := gormbulk.BulkInsert(tx, createDirector, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create music-id for insert singer,music-composer,songwtriter info*/
		type ContentMusic struct {
			Id string `json:"id"`
		}
		var contentMusic ContentMusic
		if createContentMusic := tx.Debug().Table("content_music").Create(&contentMusic).Error; createContentMusic != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(contentMusic.Id)
		/*create singers for episode*/
		var createSinger []interface{}
		for _, singers := range episodes.Music.Singers {
			contentSinger := ContentSinger{MusicId: contentMusic.Id, SingerId: singers}
			createSinger = append(createSinger, contentSinger)
		}
		if err := gormbulk.BulkInsert(tx, createSinger, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create music-composeer for episode*/
		var createMusicComposer []interface{}
		for _, musicComposers := range episodes.Music.MusicComposers {
			contentMusicComposer := ContentMusicComposer{MusicId: contentMusic.Id, MusicComposerId: musicComposers}
			createMusicComposer = append(createMusicComposer, contentMusicComposer)
		}
		if err := gormbulk.BulkInsert(tx, createMusicComposer, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create song-writer for episode*/
		var createSongWriters []interface{}
		for _, songWriters := range episodes.Music.SongWriters {
			ContentSongWriter := ContentSongWriter{MusicId: contentMusic.Id, SongWriterId: songWriters}
			createSongWriters = append(createSongWriters, ContentSongWriter)
		}
		if err := gormbulk.BulkInsert(tx, createSongWriters, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create content-tags-info for episode*/
		type ContentTagInfo struct {
			Id string `json:"id"`
		}
		var contentTagInfo ContentTagInfo
		if createContentTagInfo := tx.Debug().Table("content_tag_info").Create(&contentTagInfo).Error; createContentTagInfo != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(contentTagInfo.Id)
		/*create content-tags for episode*/
		var createContentTags []interface{}
		for _, contentTags := range episodes.TagInfo.Tags {
			contentTag := ContentTag{TagInfoId: contentTagInfo.Id, TextualDataTagId: contentTags}
			createContentTags = append(createContentTags, contentTag)
		}
		if err := gormbulk.BulkInsert(tx, createContentTags, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*create playbackItem for episode*/
		// take created by userid from request body for creating old contents else take user id from generated token
		_, _, duration := common.GetVideoDuration(episodes.PrimaryInfo.VideoContentId)
		playbackItem := PlaybackItem{VideoContentId: episodes.PrimaryInfo.VideoContentId, Duration: duration, SchedulingDateTime: episodes.SchedulingDateTime, CreatedByUserId: userId.(string), TranslationId: translation.TranslationId, RightsId: contentRights.Id}
		if createPlaybackItem := tx.Debug().Table("playback_item").Create(&playbackItem).Error; createPlaybackItem != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(playbackItem.Id)
		/*create playback_item_target_platforms*/
		var createplaybackTargetPlatforms []interface{}
		for _, platforms := range episodes.PublishingPlatforms {
			playbackItemTargetPlatform := PlaybackItemTargetPlatform{PlaybackItemId: playbackItem.Id, TargetPlatform: platforms, RightsId: contentRights.Id}
			createplaybackTargetPlatforms = append(createplaybackTargetPlatforms, playbackItemTargetPlatform)
		}
		if err := gormbulk.BulkInsert(tx, createplaybackTargetPlatforms, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		// Create a content_rights_plan for episode level
		fmt.Println("-------Creating content_rights_plan ------")
		var contentRightsPlan ContentRightsPlan
		if len(episodes.Rights.SubscriptionPlans) > 0 {
			for _, contentplanrange := range episodes.Rights.SubscriptionPlans {
				contentRightsPlan = ContentRightsPlan{RightsId: contentRights.Id, SubscriptionPlanId: contentplanrange}
				if res := tx.Debug().Table("content_rights_plan").Create(&contentRightsPlan).Error; res != nil {
					tx.Rollback()
					l.JSON(c, http.StatusInternalServerError, nil)
					return
				}
			}
		}
		fmt.Println("-------Creating content_rights_plan ------")
		/*Create An Episode*/
		var episodeKey FetchEpisodeDetails
		if contentkeyresult := tx.Debug().Table("episode").Select("max(episode_key) as episode_key,max(third_party_episode_key) as third_party_episode_key").Find(&episodeKey).Error; contentkeyresult != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": contentkeyresult.Error(), "status": http.StatusInternalServerError})
		}
		createEpisode := CreateEpisode{
			// for removing sync below line is commented
			//	Id:                  episodes.SecondaryEpisodeId, // id for creating old episodes using .net
			SeasonId:             episodes.SeasonId,
			Number:               episodes.PrimaryInfo.Number,
			PrimaryInfoId:        createPrimaryInfo.Id,
			PlaybackItemId:       playbackItem.Id,
			Status:               3,
			SynopsisEnglish:      episodes.PrimaryInfo.SynopsisEnglish,
			SynopsisArabic:       episodes.PrimaryInfo.SynopsisArabic,
			CastId:               insertCast.Id,
			MusicId:              contentMusic.Id,
			TagInfoId:            contentTagInfo.Id,
			HasPosterImage:       HasPosterImage,
			HasDubbingScript:     HasDubbingScript,
			HasSubtitlingScript:  HasSubtitlingScript,
			EpisodeKey:           episodeKey.EpisodeKey + 1,
			ThirdPartyEpisodeKey: episodeKey.ThirdPartyEpisodeKey + 1,
			//	EpisodeKey:             episodes.EpisodeKey, // for removing sync commented this line and uncommented above line
			CreatedAt:              time.Now(),
			ModifiedAt:             time.Now(),
			EnglishMetaTitle:       episodes.SeoDetails.EnglishMetaTitle,
			ArabicMetaTitle:        episodes.SeoDetails.ArabicMetaTitle,
			EnglishMetaDescription: episodes.SeoDetails.EnglishMetaDescription,
			ArabicMetaDescription:  episodes.SeoDetails.ArabicMetaDescription,
		}

		if err := tx.Debug().Table("episode").Create(&createEpisode).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var seasonDetails SeasonResult
		if seasonresult := tx.Debug().Table("season").Select("content_id").Where("id=?", episodes.SeasonId).Find(&seasonDetails).Error; seasonresult != nil {
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": seasonresult.Error(), "status": http.StatusInternalServerError})
		}
		if err := tx.Debug().Table("content").Where("id=?", seasonDetails.ContentId).Update("modified_at", time.Now()).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* upload images to S3 bucket based on episode Id*/
		go EpisodeFileUploadGcp(episodes, createEpisode.Id, seasonDetails.ContentId, episodes.SeasonId)
		res := map[string]string{
			"id": createEpisode.Id,
		}
		fmt.Println(createEpisode.Id)
		go common.RedisFlush(c)
		go common.GetMenu()
		go common.Pagekey()
		go common.Contenttype()
		l.JSON(c, http.StatusOK, gin.H{"data": res})

	} else {
		var checkEpisode int
		if epiError := db.Debug().Table("episode").Select("id").Where("id=?", c.Param("id")).Count(&checkEpisode).Error; epiError != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		if checkEpisode < 1 {
			l.JSON(c, http.StatusBadRequest, gin.H{"message": "error_season_not_found", "status": http.StatusBadRequest})
			return
		}
		/*check if episode exists in table then update episode with episodeKey*/
		var fetchEpisodeResults FetchEpisodeDetails
		if episodeResult := db.Debug().Table("episode e").Where("e.id=?", c.Param("id")).Select("e.id,e.playback_item_id,e.primary_info_id,e.cast_id,e.music_id,e.tag_info_id,pi.rights_id").
			Joins("left join playback_item pi on pi.id =e.playback_item_id").
			Scan(&fetchEpisodeResults).Error; episodeResult != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		fmt.Println(fetchEpisodeResults)
		/*update content_primary_info*/
		if updatePrimaryInfo := tx.Debug().Table("content_primary_info").Where("id=?", fetchEpisodeResults.PrimaryInfoId).Update(&createPrimaryInfo).Error; updatePrimaryInfo != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* update content_rights*/
		if updateContentRights := tx.Debug().Table("content_rights").Where("id=?", fetchEpisodeResults.RightsId).Update(&contentRights).Error; updateContentRights != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* update content_cast*/
		if updateContentcast := tx.Debug().Table("content_cast").Where("id=?", fetchEpisodeResults.CastId).Update(&insertCast).Error; updateContentcast != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update playback_item*/
		_, _, duration := common.GetVideoDuration(episodes.PrimaryInfo.VideoContentId)
		playbackItem := PlaybackItem{VideoContentId: episodes.PrimaryInfo.VideoContentId, Duration: duration, SchedulingDateTime: episodes.SchedulingDateTime, CreatedByUserId: userId.(string), TranslationId: translation.TranslationId, RightsId: fetchEpisodeResults.RightsId}
		//new implementation
		if episodes.SchedulingDateTime == nil {
			if res := tx.Debug().Table("playback_item").Select("scheduling_date_time").Where("id=?", fetchEpisodeResults.PlaybackItemId).Updates(map[string]interface{}{"scheduling_date_time": gorm.Expr("NULL")}).Error; res != nil {
				tx.Rollback()
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}
			if updateContentcast := tx.Debug().Table("playback_item").Where("id=?", fetchEpisodeResults.PlaybackItemId).Update(&playbackItem).Error; updateContentcast != nil {
				tx.Rollback()
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}
		} else {
			if updateContentcast := tx.Debug().Table("playback_item").Where("id=?", fetchEpisodeResults.PlaybackItemId).Update(&playbackItem).Error; updateContentcast != nil {
				tx.Rollback()
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}
		}
		/*update content_rights_country*/
		var contentRightsCountry ContentRightsCountry
		if err := tx.Debug().Table("content_rights_country").Where("content_rights_id=?", fetchEpisodeResults.RightsId).Delete(&contentRightsCountry).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createcontentRightscountry []interface{}
		for _, countries := range episodes.Rights.DigitalRightsRegions {
			contentRightsCountry = ContentRightsCountry{ContentRightsId: fetchEpisodeResults.RightsId, CountryId: countries}
			createcontentRightscountry = append(createcontentRightscountry, contentRightsCountry)
		}
		if err := gormbulk.BulkInsert(tx, createcontentRightscountry, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		// contentRightsPlan for episodes
		fmt.Println("-------episodes-----contentRightsPlan-----")
		var contentRightsPlans ContentRightsPlan
		db.Debug().Where("rights_id=?", contentRights.Id).Delete(&contentRightsPlans)
		if episodes.Rights.SubscriptionPlans != nil && len(episodes.Rights.SubscriptionPlans) > 0 {
			for _, plan := range episodes.Rights.SubscriptionPlans {
				var contentRightsPlan ContentRightsPlan
				contentRightsPlan.RightsId = contentRights.Id
				contentRightsPlan.SubscriptionPlanId = plan
				if err := tx.Debug().Create(&contentRightsPlan).Error; err != nil {
					l.JSON(c, http.StatusInternalServerError, serverError)
					// return "", serverError, 0
				}
			}
		}
		fmt.Println("-------episodes-----contentRightsPlan-----")
		/*update content_actor*/
		var contentActor ContentActor
		if err := tx.Debug().Table("content_actor").Where("cast_id=?", fetchEpisodeResults.CastId).Delete(&contentActor).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createActor []interface{}
		for _, actors := range episodes.Cast.Actors {
			contentActor := ContentActor{CastId: fetchEpisodeResults.CastId, ActorId: actors}
			createActor = append(createActor, contentActor)
		}
		if err := gormbulk.BulkInsert(tx, createActor, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update content_writer*/
		var contentWriter ContentWriter
		if err := tx.Debug().Table("content_writer").Where("cast_id=?", fetchEpisodeResults.CastId).Delete(&contentWriter).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createWriter []interface{}
		for _, writers := range episodes.Cast.Writers {
			contentWriter := ContentWriter{CastId: fetchEpisodeResults.CastId, WriterId: writers}
			createWriter = append(createWriter, contentWriter)
		}
		if err := gormbulk.BulkInsert(tx, createWriter, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update content_director*/
		var contentDirector ContentDirector
		if err := tx.Debug().Table("content_director").Where("cast_id=?", fetchEpisodeResults.CastId).Delete(&contentDirector).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createDirector []interface{}
		for _, directors := range episodes.Cast.Directors {
			contentDirector := ContentDirector{CastId: fetchEpisodeResults.CastId, DirectorId: directors}
			createDirector = append(createDirector, contentDirector)
		}
		if err := gormbulk.BulkInsert(tx, createDirector, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}

		/*update content_singers*/
		var contentSinger ContentSinger
		if err := tx.Debug().Table("content_singer").Where("music_id=?", fetchEpisodeResults.MusicId).Delete(&contentSinger).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createSinger []interface{}
		for _, singers := range episodes.Music.Singers {
			contentSinger := ContentSinger{MusicId: fetchEpisodeResults.MusicId, SingerId: singers}
			createSinger = append(createSinger, contentSinger)
		}
		if err := gormbulk.BulkInsert(tx, createSinger, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* update content_music_composer */
		var contentMusicComposer ContentMusicComposer
		if err := tx.Debug().Table("content_music_composer").Where("music_id=?", fetchEpisodeResults.MusicId).Delete(&contentMusicComposer).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createMusicComposer []interface{}
		for _, musicComposers := range episodes.Music.MusicComposers {
			contentMusicComposer := ContentMusicComposer{MusicId: fetchEpisodeResults.MusicId, MusicComposerId: musicComposers}
			createMusicComposer = append(createMusicComposer, contentMusicComposer)
		}
		if err := gormbulk.BulkInsert(tx, createMusicComposer, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* update content_song_writers */
		var contentSongWriter ContentSongWriter
		if err := tx.Debug().Table("content_song_writer").Where("music_id=?", fetchEpisodeResults.MusicId).Delete(&contentSongWriter).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createSongWriters []interface{}
		for _, songWriters := range episodes.Music.SongWriters {
			ContentSongWriter := ContentSongWriter{MusicId: fetchEpisodeResults.MusicId, SongWriterId: songWriters}
			createSongWriters = append(createSongWriters, ContentSongWriter)
		}
		if err := gormbulk.BulkInsert(tx, createSongWriters, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update content-tags*/
		var contentTags ContentTag
		if err := tx.Debug().Table("content_tag").Where("tag_info_id=?", fetchEpisodeResults.TagInfoId).Delete(&contentTags).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createContentTags []interface{}
		for _, contentTags := range episodes.TagInfo.Tags {
			contentTag := ContentTag{TagInfoId: fetchEpisodeResults.TagInfoId, TextualDataTagId: contentTags}
			createContentTags = append(createContentTags, contentTag)
		}
		if err := gormbulk.BulkInsert(tx, createContentTags, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update playback_item_target_platforms*/
		var playbackItemTargetPlatform PlaybackItemTargetPlatform
		if err := tx.Debug().Table("playback_item_target_platform").Where("playback_item_id=?", fetchEpisodeResults.PlaybackItemId).Delete(&playbackItemTargetPlatform).Error; err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var createplaybackTargetPlatforms []interface{}
		for _, platforms := range episodes.PublishingPlatforms {
			playbackItemTargetPlatform := PlaybackItemTargetPlatform{PlaybackItemId: fetchEpisodeResults.PlaybackItemId, TargetPlatform: platforms, RightsId: fetchEpisodeResults.RightsId}
			createplaybackTargetPlatforms = append(createplaybackTargetPlatforms, playbackItemTargetPlatform)
		}
		if err := gormbulk.BulkInsert(tx, createplaybackTargetPlatforms, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/*update episode*/
		Episode := CreateEpisode{
			SeasonId:               episodes.SeasonId,
			Number:                 episodes.PrimaryInfo.Number,
			SynopsisEnglish:        episodes.PrimaryInfo.SynopsisEnglish,
			SynopsisArabic:         episodes.PrimaryInfo.SynopsisArabic,
			HasPosterImage:         HasPosterImage,
			HasDubbingScript:       HasDubbingScript,
			HasSubtitlingScript:    HasSubtitlingScript,
			ModifiedAt:             time.Now(),
			Status:                 3,
			EnglishMetaTitle:       episodes.SeoDetails.EnglishMetaTitle,
			ArabicMetaTitle:        episodes.SeoDetails.ArabicMetaTitle,
			EnglishMetaDescription: episodes.SeoDetails.EnglishMetaDescription,
			ArabicMetaDescription:  episodes.SeoDetails.ArabicMetaDescription,
		}
		if updateEpisode := tx.Debug().Table("episode").Where("id=?", fetchEpisodeResults.Id).Update(&Episode).Error; updateEpisode != nil {
			tx.Rollback()
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var seasonDetails SeasonResult
		if seasonresult := tx.Debug().Table("season").Select("content_id").Where("id=?", episodes.SeasonId).Find(&seasonDetails).Error; seasonresult != nil {
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": seasonresult.Error(), "status": http.StatusInternalServerError})
		}
		if err := tx.Debug().Table("content").Where("id=?", seasonDetails.ContentId).Update("modified_at", time.Now()).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		/* upload images to S3 bucket based on episode Id*/
		go EpisodeFileUploadGcp(episodes, c.Param("id"), seasonDetails.ContentId, episodes.SeasonId)
		res := map[string]string{
			"id": c.Param("id"),
		}
		/*create redis keys */
		var episodedetails common.EpisodeDetails
		db.Debug().Raw("select c.content_key,c.content_type from content c join season s on s.content_id = c.id join episode e on e.season_id  = s.id where e.id =?", c.Param("id")).Find(&episodedetails)
		contentkeyconverted := strconv.Itoa(episodedetails.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(episodedetails.ContentType, c)
		go common.RedisFlush(c)
		go common.GetMenu()
		go common.Pagekey()
		go common.Contenttype()
		l.JSON(c, http.StatusOK, gin.H{"data": res})
	}
	/*commit changes*/
	if err := tx.Commit().Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}

}

// GetSeasonDetailsBySeasonId - Get Season details by season id
// GET /api/seasons/:id
// @Summary Get Season details by season id
// @Description Get Season details by season id
// @Tags Page
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param id path string true "Id"
// @Success 200 {array} object c.JSON
// @Router/api/seasons/{id} [get]
func (hs *HandlerService) GetSeasonDetailsBySeasonId(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	db := c.MustGet("DB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "ar" {
		language = "en"
	}
	var seasonResult SeasonResult
	serverError := common.ServerErrorResponse()
	var finalSeasonResult FinalSeasonResult
	seasonId := c.Param("id")
	if err := db.Debug().Table("season s").Select("s.id,s.content_id,s.season_key ,s.status ,case when s.status = 1 then true else false end as status_can_be_changed,s.modified_at,s.number as season_number,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.intro_start ,cpi.outro_start, cc.main_actor_id ,cc.main_actress_id ,atci.original_language ,atci.supplier ,atci.acquisition_department ,atci.english_synopsis ,atci.arabic_synopsis ,atci.production_year ,atci.production_house ,atci.age_group ,atci.intro_duration as about_intro_duration,atci.intro_start as about_intro_start,atci.outro_duration as about_outro_duration,atci.outro_start as about_outro_start,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date ,s.created_by_user_id as created_by,s.english_meta_title ,s.arabic_meta_title ,s.english_meta_description ,s.arabic_meta_description,s.cast_id,s.music_id,s.rights_id,s.tag_info_id,s.about_the_content_info_id").
		Joins("left join content_primary_info cpi on cpi.id =s.primary_info_id").
		Joins("left join content_cast cc on cc.id =s.cast_id").
		Joins("left join about_the_content_info atci on atci.id =s.about_the_content_info_id").
		Joins("left join content_translation ct on ct.id =s.translation_id").
		Joins("left join content_rights cr on cr.id=s.rights_id").
		Where("s.id = ? and s.deleted_by_user_id is null", seasonId).
		Find(&finalSeasonResult).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		fmt.Println(serverError, "..............")
		return
	}
	seasonResult.ContentId = finalSeasonResult.ContentId
	seasonResult.SeasonKey = finalSeasonResult.SeasonKey
	seasonResult.Status = finalSeasonResult.Status
	seasonResult.StatusCanBeChanged = finalSeasonResult.StatusCanBeChanged
	seasonResult.ModifiedAt = finalSeasonResult.ModifiedAt
	seasonResult.PrimaryInfo.SeasonNumber = finalSeasonResult.SeasonNumber
	seasonResult.PrimaryInfo.OriginalTitle = finalSeasonResult.OriginalTitle
	seasonResult.PrimaryInfo.AlternativeTitle = finalSeasonResult.AlternativeTitle
	seasonResult.PrimaryInfo.ArabicTitle = finalSeasonResult.ArabicTitle
	seasonResult.PrimaryInfo.TransliteratedTitle = finalSeasonResult.TransliteratedTitle
	seasonResult.PrimaryInfo.Notes = finalSeasonResult.Notes
	seasonResult.PrimaryInfo.IntroStart = finalSeasonResult.IntroStart
	seasonResult.PrimaryInfo.OutroStart = finalSeasonResult.OutroStart
	seasonResult.Cast.MainActorId = finalSeasonResult.MainActorId
	seasonResult.Cast.MainActressId = finalSeasonResult.MainActressId
	/* Fetch content_actors*/
	var contentActor []ContentActor
	if actorResult := db.Debug().Table("content_actor").Select("actor_id").Where("cast_id=?", finalSeasonResult.CastId).Scan(&contentActor).Error; actorResult != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var actors []string
	for _, actorIds := range contentActor {
		actors = append(actors, actorIds.ActorId)
	}
	seasonResult.Cast.Actors = actors
	if len(actors) < 1 {
		buffer := make([]string, 0)
		seasonResult.Cast.Actors = buffer
	}
	/* Fetch content_writers*/
	var contentWriters []ContentWriter
	if writerResult := db.Debug().Table("content_writer").Select("writer_id").Where("cast_id=?", finalSeasonResult.CastId).Scan(&contentWriters).Error; writerResult != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var writers []string
	for _, writerIds := range contentWriters {
		writers = append(writers, writerIds.WriterId)
	}
	seasonResult.Cast.Writers = writers
	if len(writers) < 1 {
		buffer := make([]string, 0)
		seasonResult.Cast.Writers = buffer
	}
	/* Fetch content_directors*/
	var contentDirectors []ContentDirector
	if directorResult := db.Debug().Table("content_director").Select("director_id").Where("cast_id=?", finalSeasonResult.CastId).Scan(&contentDirectors).Error; directorResult != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var directors []string
	for _, directorIds := range contentDirectors {
		directors = append(directors, directorIds.DirectorId)
	}
	seasonResult.Cast.Directors = directors
	if len(directors) < 1 {
		buffer := make([]string, 0)
		seasonResult.Cast.Directors = buffer
	}
	/* Fetch content_singers*/
	var contentSingers []ContentSinger
	if singerResult := db.Debug().Table("content_singer").Select("singer_id").Where("music_id=?", finalSeasonResult.MusicId).Scan(&contentSingers).Error; singerResult != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var singers []string
	for _, singerIds := range contentSingers {
		singers = append(singers, singerIds.SingerId)
	}
	seasonResult.Music.Singers = singers
	if len(singers) < 1 {
		buffer := make([]string, 0)
		seasonResult.Music.Singers = buffer
	}
	/* Fetch content_music_composers*/
	var contentMusicComposers []ContentMusicComposer
	db.Debug().Table("content_music_composer").Select("music_composer_id").Where("music_id=?", finalSeasonResult.MusicId).Scan(&contentMusicComposers)
	var musicComposers []string
	for _, composerIds := range contentMusicComposers {
		musicComposers = append(musicComposers, composerIds.MusicComposerId)
	}
	seasonResult.Music.MusicComposers = musicComposers
	if len(musicComposers) < 1 {
		buffer := make([]string, 0)
		seasonResult.Music.MusicComposers = buffer
	}
	/* Fetch content_writers*/
	var contentSongWriters []ContentSongWriter
	db.Debug().Table("content_song_writer").Select("song_writer_id").Where("music_id=?", finalSeasonResult.MusicId).Scan(&contentSongWriters)
	var songWriters []string
	for _, songWritersIds := range contentSongWriters {
		songWriters = append(songWriters, songWritersIds.SongWriterId)
	}
	seasonResult.Music.SongWriters = songWriters
	if len(songWriters) < 1 {
		buffer := make([]string, 0)
		seasonResult.Music.SongWriters = buffer
	}
	/*Fetch tag_info*/
	var contentTagInfo []ContentTag
	db.Debug().Table("content_tag").Select("textual_data_tag_id").Where("tag_info_id=?", finalSeasonResult.TagInfoId).Scan(&contentTagInfo)
	var tagInfo []string
	for _, tagInfoIds := range contentTagInfo {
		tagInfo = append(tagInfo, tagInfoIds.TextualDataTagId)
	}
	seasonResult.TagInfo.Tags = tagInfo
	if len(tagInfo) < 1 {
		buffer := make([]string, 0)
		seasonResult.TagInfo.Tags = buffer
	}
	/*Fetch Season_geners*/
	//var seasonGenres []SeasonGenres
	var finalSeasonGenre []interface{}
	var newSeasonGenres NewSeasonGenres
	// if genreResult := db.Debug().Table("season_genre").Select("id,genre_id").Where("season_id=?", seasonId).Scan(&seasonGenres).Error; genreResult != nil {
	// 	l.JSON(c, http.StatusInternalServerError, serverError)
	// 	fmt.Println(".................................")
	// 	return
	// }
	// for _, tagInfoIds := range seasonGenres {
	// 	var seasonSubgenre []SeasonSubgenre
	// 	if subgenreResult := db.Debug().Table("season_subgenre").Select("subgenre_id").Where("season_genre_id=?", tagInfoIds.Id).Scan(&seasonSubgenre).Error; subgenreResult != nil {
	// 		l.JSON(c, http.StatusInternalServerError, serverError)
	// 		fmt.Println("///////////////////")
	// 		return
	// 	}
	// 	var Subgenre []string
	// 	for _, data := range seasonSubgenre {
	// 		Subgenre = append(Subgenre, data.SubgenreId)
	// 		newSeasonGenres.GenreId = tagInfoIds.GenreId
	// 		newSeasonGenres.Id = tagInfoIds.Id
	// 		newSeasonGenres.SubgenresId = Subgenre
	// 		finalSeasonGenre = append(finalSeasonGenre, newSeasonGenres)
	// 	}
	// }
	var genres []ContentGeneresQueryDetails
	//var contentGenres []ContentGenres
	genres = nil
	if err := db.Debug().Table("season_genre sg").Select("sg.genre_id,json_agg(ss.subgenre_id order by ss.order)::varchar as subgenres_id,json_agg(ss.order order by ss.order)::varchar as sub_genre_order,sg.id").Joins("join season_subgenre ss on ss.season_genre_id = sg.id join season c on c.id = sg.season_id").Where("sg.season_id = ?", seasonId).Group("sg.genre_id,sg.id").Order("sg.order").Find(&genres).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	for _, genre := range genres {
		newSeasonGenres.GenreId = genre.GenreId
		subGenres, err := JsonStringToStringSliceOrMap(genre.SubgenresId)
		if err != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		newSeasonGenres.SubgenresId = subGenres
		newSeasonGenres.Id = genre.Id
		finalSeasonGenre = append(finalSeasonGenre, newSeasonGenres)
	}

	seasonResult.SeasonGenres = finalSeasonGenre
	seasonResult.AboutTheContent.OriginalLanguage = finalSeasonResult.OriginalLanguage
	seasonResult.AboutTheContent.Supplier = finalSeasonResult.Supplier
	seasonResult.AboutTheContent.AcquisitionDepartment = finalSeasonResult.AcquisitionDepartment
	seasonResult.AboutTheContent.EnglishSynopsis = finalSeasonResult.EnglishSynopsis
	seasonResult.AboutTheContent.ArabicSynopsis = finalSeasonResult.ArabicSynopsis
	seasonResult.AboutTheContent.ProductionYear = finalSeasonResult.ProductionYear
	seasonResult.AboutTheContent.ProductionHouse = finalSeasonResult.ProductionHouse
	seasonResult.AboutTheContent.AgeGroup = finalSeasonResult.AgeGroup
	seasonResult.AboutTheContent.IntroDuration = finalSeasonResult.AboutIntroDuration
	seasonResult.AboutTheContent.IntroStart = finalSeasonResult.AboutIntroStart
	seasonResult.AboutTheContent.OutroDuration = finalSeasonResult.AboutOutroDuration
	seasonResult.AboutTheContent.OutroStart = finalSeasonResult.AboutOutroStart
	/*Fetch Production_country*/
	var productionCountry []ProductionCountry
	if productionCountryResult := db.Debug().Table("production_country").Select("country_id").Where("about_the_content_info_id=?", finalSeasonResult.AboutTheContentInfoId).Order("country_id").Scan(&productionCountry).Error; productionCountryResult != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var countries []int
	for _, prcountries := range productionCountry {
		countries = append(countries, prcountries.CountryId)
	}
	seasonResult.AboutTheContent.ProductionCountries = countries

	// if len(tagInfo) < 1 {
	// 	buffer := make([]int, 0)
	// 	seasonResult.AboutTheContent.ProductionCountries = buffer
	// }

	seasonResult.Translation.LanguageType = LanguageOriginTypes(finalSeasonResult.LanguageType)
	seasonResult.Translation.DubbingLanguage = finalSeasonResult.DubbingLanguage
	seasonResult.Translation.DubbingDialectId = finalSeasonResult.DubbingDialectId
	seasonResult.Translation.SubtitlingLanguage = finalSeasonResult.SubtitlingLanguage
	/*calling GetEpisodes function to fetch all episodes which belongs to seasonId*/
	var episodecount int
	db.Debug().Raw("select count(*) from season s join episode e on e.season_id =s.id where s.id='" + seasonId + "' and e.deleted_by_user_id is null").Count(&episodecount)
	if episodecount > 0 {
		seasonResult.Episodes = GetEpisodes(c, seasonId)
	} else {
		seasonResult.Episodes = nil
	}
	seasonResult.NonTextualData.PosterImage = os.Getenv("IMAGERY_URL") + finalSeasonResult.ContentId + "/" + finalSeasonResult.Id + "/poster-image"
	seasonResult.NonTextualData.OverlayPosterImage = os.Getenv("IMAGERY_URL") + finalSeasonResult.ContentId + "/" + finalSeasonResult.Id + "/overlay-poster-image"
	seasonResult.NonTextualData.DetailsBackground = os.Getenv("IMAGERY_URL") + finalSeasonResult.ContentId + "/" + finalSeasonResult.Id + "/details-background"
	seasonResult.NonTextualData.MobileDetailsBackground = os.Getenv("IMAGERY_URL") + finalSeasonResult.ContentId + "/" + finalSeasonResult.Id + "/mobile-details-background"
	seasonResult.NonTextualData.SeasonLogo = os.Getenv("IMAGERY_URL") + finalSeasonResult.ContentId + "/" + finalSeasonResult.Id + "/season-logo"
	seasonResult.Rights.DigitalRightsType = finalSeasonResult.DigitalRightsType
	seasonResult.Rights.DigitalRightsStartDate = finalSeasonResult.DigitalRightsStartDate
	seasonResult.Rights.DigitalRightsEndDate = finalSeasonResult.DigitalRightsEndDate
	/*Fetch Digital_right_Regions*/
	var digitalRightsRegions []DigitalRightsRegions
	if countryError := db.Debug().Table("content_rights_country").Select("country_id").Where("content_rights_id=?", finalSeasonResult.RightsId).Scan(&digitalRightsRegions).Error; countryError != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var RegionRights []int
	for _, idarr := range digitalRightsRegions {
		RegionRights = append(RegionRights, idarr.CountryId)
	}
	seasonResult.Rights.DigitalRightsRegions = RegionRights
	if len(RegionRights) < 1 {
		buffer := make([]int, 0)
		seasonResult.Rights.DigitalRightsRegions = buffer
	}
	/*Fetch Subscriptions*/
	var subscriptionPlans []SubscriptionPlans
	if subscriptionError := db.Debug().Table("content_rights_plan").Select("subscription_plan_id").Where("rights_id=?", finalSeasonResult.RightsId).Scan(&subscriptionPlans).Error; subscriptionError != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var subscriptions []int
	for _, idarr := range subscriptionPlans {
		subscriptions = append(subscriptions, idarr.SubscriptionPlanId)
	}
	seasonResult.Rights.SubscriptionPlans = subscriptions
	if len(subscriptions) < 1 {
		buffer := make([]int, 0)
		seasonResult.Rights.SubscriptionPlans = buffer
	}

	seasonResult.CreatedBy = finalSeasonResult.CreatedBy
	seasonResult.IntroDuration = "00:00:00"
	seasonResult.IntroStart = "00:00:00"
	seasonResult.OutroDuration = "00:00:00"
	seasonResult.OutroStart = "00:00:00"
	var rightProduct []RightProduct
	if productResult := db.Debug().Table("rights_product").Select("product_name").Where("rights_id=?", finalSeasonResult.RightsId).Order("product_name").Scan(&rightProduct).Error; productResult != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var products []int
	for _, productNumbers := range rightProduct {
		products = append(products, productNumbers.ProductName)
	}
	seasonResult.Products = products
	if len(products) < 1 {
		buffer := make([]int, 0)
		seasonResult.Products = buffer
	}
	seasonResult.SeoDetails.EnglishMetaTitle = finalSeasonResult.EnglishMetaTitle
	seasonResult.SeoDetails.ArabicMetaTitle = finalSeasonResult.ArabicMetaTitle
	seasonResult.SeoDetails.EnglishMetaDescription = finalSeasonResult.EnglishMetaDescription
	seasonResult.SeoDetails.ArabicMetaDescription = finalSeasonResult.ArabicMetaDescription
	/* Fetch Variance_Trailers*/
	var varianceTrailers []VarianceTrailers
	if varianceTrailersError := db.Debug().Raw("select * from variance_trailer vt where vt.season_id=? order by vt.order asc ", seasonId).Find(&varianceTrailers).Error; varianceTrailersError != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	var variance VarianceTrailers
	var finalVarianceResult []interface{}
	for _, varianceData := range varianceTrailers {
		variance.Order = varianceData.Order
		variance.VideoTrailerId = varianceData.VideoTrailerId
		variance.EnglishTitle = varianceData.EnglishTitle
		variance.ArabicTitle = varianceData.ArabicTitle
		variance.Duration = varianceData.Duration
		variance.HasTrailerPosterImage = varianceData.HasTrailerPosterImage
		//const url = "https://z5content-qa.s3.amazonaws.com/"
		variance.TrailerposterImage = os.Getenv("IMAGERY_URL") + seasonResult.ContentId + "/" + seasonId + "/" + varianceData.Id + "/trailer-poster-image"
		variance.Id = varianceData.Id
		finalVarianceResult = append(finalVarianceResult, variance)
	}
	seasonResult.VarianceTrailers = finalVarianceResult
	seasonResult.Id = finalSeasonResult.Id
	l.JSON(c, http.StatusOK, gin.H{"data": seasonResult})
	return
}

func LanguageOriginTypes(originType int) string {
	OriginTypesArray := map[int]string{1: "Original", 2: "Dubbed", 3: "Subtitled"}
	return OriginTypesArray[originType]
}

/* Get all Episodes by seasonId */
func GetEpisodes(c *gin.Context, seasonid string) []interface{} {
	db := c.MustGet("DB").(*gorm.DB)
	userdb := c.MustGet("UDB").(*gorm.DB)
	// serverError := common.ServerErrorResponse()
	var seasonId string
	var episodelist EpisodeDetailsByseasonId
	var finalEpisodesResult []FinalEpisodesResult
	var finalResult []interface{}
	seasonId = seasonid
	var totalcount int
	db.Debug().Raw("select count(*) from season s join episode e on e.season_id =s.id where s.id='" + seasonId + "' and e.deleted_by_user_id is null").Count(&totalcount)
	db.Debug().Table("season as s").Select("e.id ,case when e.primary_info_id is null then false else true end as is_primary,s.content_id,e.episode_key,e.season_id ,s.status,case when s.status=1  then true else false end as status_can_be_changed,e.status as sub_status ,e.number,e.synopsis_english ,e.synopsis_arabic ,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.intro_start ,cpi.outro_start,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date,ds.name as sub_status_name,pi2.video_content_id,s.created_by_user_id,e.cast_id,e.music_id,pi2.rights_id,e.tag_info_id,s.about_the_content_info_id,e.playback_item_id, cc.main_actor_id ,cc.main_actress_id").
		Joins("left join episode e on e.season_id =s.id").
		Joins("left join content_primary_info cpi on cpi.id =e.primary_info_id").
		Joins("left join content_cast cc on cc.id =e.cast_id").
		Joins("left join content_translation ct on ct.id =s.translation_id").
		Joins("left join playback_item pi2 on pi2 .id =e.playback_item_id").
		Joins("left join content_rights cr on cr.id =pi2.rights_id").
		Joins("left join display_status ds on ds.id =e.status").
		Where("s.id = ? and e.deleted_by_user_id is null", seasonId).
		Find(&finalEpisodesResult)

	var userInfo UserInfo
	userdb.Debug().Table("user").Select("user_name").Where("id=?", finalEpisodesResult[0].CreatedByUserId).Scan(&userInfo)
	createdBy := userInfo.UserName
	for _, eplist := range finalEpisodesResult {
		episodelist.IsPrimary = eplist.IsPrimary
		episodelist.UserId = "00000000-0000-0000-0000-000000000000"
		episodelist.SecondarySeasonId = "00000000-0000-0000-0000-000000000000"
		episodelist.VarianceIds = eplist.VarianceIds
		episodelist.EpisodeIds = eplist.EpisodeIds
		episodelist.SecondaryEpisodeId = "00000000-0000-0000-0000-000000000000"
		episodelist.ContentId = eplist.ContentId
		episodelist.EpisodeKey = eplist.EpisodeKey
		episodelist.SeasonId = eplist.SeasonId
		episodelist.Status = eplist.Status
		episodelist.StatusCanBeChanged = eplist.StatusCanBeChanged
		episodelist.SubStatus = eplist.SubStatus
		episodelist.SubStatusName = eplist.SubStatusName
		episodelist.DigitalRightsType = eplist.DigitalRightsType
		episodelist.DigitalRightsStartDate = eplist.DigitalRightsStartDate
		episodelist.DigitalRightsEndDate = eplist.DigitalRightsEndDate
		episodelist.CreatedBy = createdBy
		episodelist.PrimaryInfo.Number = eplist.Number
		episodelist.PrimaryInfo.VideoContentId = eplist.VideoContentId
		episodelist.PrimaryInfo.SynopsisEnglish = eplist.SynopsisEnglish
		episodelist.PrimaryInfo.SynopsisArabic = eplist.SynopsisArabic
		episodelist.PrimaryInfo.OriginalTitle = eplist.OriginalTitle
		episodelist.PrimaryInfo.AlternativeTitle = eplist.AlternativeTitle
		episodelist.PrimaryInfo.ArabicTitle = eplist.ArabicTitle
		episodelist.PrimaryInfo.TransliteratedTitle = eplist.TransliteratedTitle
		episodelist.PrimaryInfo.Notes = eplist.Notes
		episodelist.PrimaryInfo.IntroStart = eplist.IntroStart
		episodelist.PrimaryInfo.OutroStart = eplist.OutroStart
		episodelist.Cast.MainActorId = eplist.MainActorId
		episodelist.Cast.MainActressId = eplist.MainActressId
		/* Fetch content_actors*/
		var contentActor []ContentActor
		db.Debug().Table("content_actor").Select("actor_id").Where("cast_id=?", eplist.CastId).Scan(&contentActor)
		var actors []string
		for _, actorIds := range contentActor {
			actors = append(actors, actorIds.ActorId)
		}
		episodelist.Cast.Actors = actors
		if len(actors) < 1 {
			buffer := make([]string, 0)
			episodelist.Cast.Actors = buffer
		}
		/* Fetch content_writers*/
		var contentWriters []ContentWriter
		db.Debug().Table("content_writer").Select("writer_id").Where("cast_id=?", eplist.CastId).Scan(&contentWriters)

		var writers []string
		for _, writerIds := range contentWriters {
			writers = append(writers, writerIds.WriterId)
		}
		episodelist.Cast.Writers = writers
		if len(writers) < 1 {
			buffer := make([]string, 0)
			episodelist.Cast.Writers = buffer
		}
		/* Fetch content_directors*/
		var contentDirectors []ContentDirector
		db.Debug().Table("content_director").Select("director_id").Where("cast_id=?", eplist.CastId).Scan(&contentDirectors)

		var directors []string
		for _, directorIds := range contentDirectors {
			directors = append(directors, directorIds.DirectorId)
		}
		episodelist.Cast.Directors = directors
		if len(directors) < 1 {
			buffer := make([]string, 0)
			episodelist.Cast.Directors = buffer
		}
		/* Fetch content_singers*/
		var contentSingers []ContentSinger
		db.Debug().Table("content_singer").Select("singer_id").Where("music_id=?", eplist.MusicId).Scan(&contentSingers)

		var singers []string
		for _, singerIds := range contentSingers {
			singers = append(singers, singerIds.SingerId)
		}
		episodelist.Music.Singers = singers
		if len(singers) < 1 {
			buffer := make([]string, 0)
			episodelist.Music.Singers = buffer
		}
		/* Fetch content_music_composers*/
		var contentMusicComposers []ContentMusicComposer
		db.Debug().Table("content_music_composer").Select("music_composer_id").Where("music_id=?", eplist.MusicId).Scan(&contentMusicComposers)
		var musicComposers []string
		for _, composerIds := range contentMusicComposers {
			musicComposers = append(musicComposers, composerIds.MusicComposerId)
		}
		episodelist.Music.MusicComposers = musicComposers
		if len(musicComposers) < 1 {
			buffer := make([]string, 0)
			episodelist.Music.MusicComposers = buffer
		}
		/* Fetch content_writers*/
		var contentSongWriters []ContentSongWriter
		db.Debug().Table("content_song_writer").Select("song_writer_id").Where("music_id=?", eplist.MusicId).Scan(&contentSongWriters)

		var songWriters []string
		for _, songWritersIds := range contentSongWriters {
			songWriters = append(songWriters, songWritersIds.SongWriterId)
		}
		episodelist.Music.SongWriters = songWriters
		if len(songWriters) < 1 {
			buffer := make([]string, 0)
			episodelist.Music.SongWriters = buffer
		}
		/*Fetch tag_info*/
		var contentTagInfo []ContentTag
		db.Debug().Table("content_tag").Select("textual_data_tag_id").Where("tag_info_id=?", eplist.TagInfoId).Scan(&contentTagInfo)
		var tagInfo []string
		for _, tagInfoIds := range contentTagInfo {
			tagInfo = append(tagInfo, tagInfoIds.TextualDataTagId)
		}
		episodelist.TagInfo.Tags = tagInfo
		if len(tagInfo) < 1 {
			buffer := make([]string, 0)
			episodelist.TagInfo.Tags = buffer
		}
		episodelist.NonTextualData = eplist.NonTextualData
		episodelist.Translation.LanguageType = LanguageOriginTypes(eplist.LanguageType)
		episodelist.Translation.DubbingLanguage = eplist.DubbingLanguage
		episodelist.Translation.DubbingDialectId = eplist.DubbingDialectId
		episodelist.Translation.SubtitlingLanguage = eplist.SubtitlingLanguage
		episodelist.SchedulingDateTime = eplist.SchedulingDateTime
		//Fetch Publishing_platforms
		var targetPlatforms []PlaybackItemTargetPlatform
		db.Debug().Table("playback_item_target_platform").Select("target_platform").Where("playback_item_id=?", eplist.PlaybackItemId).Scan(&targetPlatforms)
		var platforms []int
		for _, idarr := range targetPlatforms {
			platforms = append(platforms, idarr.TargetPlatform)
		}
		episodelist.PublishingPlatforms = platforms
		if len(platforms) < 1 {
			buffer := make([]int, 0)
			episodelist.PublishingPlatforms = buffer
		}
		episodelist.SeoDetails = eplist.SeoDetails
		episodelist.Id = eplist.Id
		finalResult = append(finalResult, episodelist)
	}
	return finalResult
}

// GetOneTierContentDetailsBasedonContentID - Get One Tier Content Details Based on Content ID
// GET /api/contents/onetier/:id
// @Summary Get One Tier Content Details Based on Content ID
// @Description Get One Tier Content Details Based on Content ID
// @Tags Page
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param id path string true "Id"
// @Success 200 {array} object c.JSON
// @Router /api/contents/onetier/{id} [get]
func (hs *HandlerService) GetOneTierContentDetailsBasedonContentID(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	db := c.MustGet("DB").(*gorm.DB)
	language := strings.ToLower(c.Param("lang"))
	if language != "ar" {
		language = "en"
	}
	var contentResult OnetireContent
	serverError := common.ServerErrorResponse()
	var finalContentResults []FinalSeasonResultNew
	contentId := c.Param("id")
	if err := db.Debug().Debug().Table("content c").Select("c.id ,c.content_key , c.status,c.content_type ,c.english_meta_title ,c.arabic_meta_title ,c.english_meta_description ,c.arabic_meta_description ,c.has_poster_image ,c.has_details_background ,c.has_mobile_details_background ,c.cast_id,c.music_id,c.tag_info_id,c.about_the_content_info_id,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.intro_start ,cpi.outro_start ,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language,pi2.id as playback_item_id,pi2.video_content_id ,pi2.rights_id,pi2.scheduling_date_time,cv.id as variance_id,cv.status as variance_status,cv.has_overlay_poster_image ,cv.has_dubbing_script ,cv.has_subtitling_script ,cv.intro_duration as variance_intro_duration ,cv.intro_start as variance_intro_start ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date ,atci.original_language ,atci.supplier ,atci.acquisition_department ,atci.english_synopsis ,atci.arabic_synopsis ,atci.production_year ,atci.production_house ,atci.age_group ,atci.intro_duration as about_intro_duration,atci.intro_start as about_intro_start,atci.outro_duration as about_outro_duration,atci.outro_start as about_outro_start,cc.main_actor_id ,cc.main_actress_id").
		Joins("left join content_primary_info cpi on cpi.id =c.primary_info_id").
		Joins("left join content_variance cv on cv.content_id =c.id").
		Joins("left join playback_item pi2 on pi2.id =cv.playback_item_id").
		Joins("left join content_translation ct on ct.id =pi2.translation_id").
		Joins("left join content_rights cr on cr.id =pi2.rights_id").
		Joins("left join content_cast cc on cc.id =c.cast_id").
		Joins("left join about_the_content_info atci on atci.id=c.about_the_content_info_id").
		Where("c.id = ? and c.content_tier=? and c.deleted_by_user_id is null and cv.deleted_by_user_id is null", contentId, 1).
		Order("cv.order asc").Find(&finalContentResults).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		fmt.Println(serverError, "...")
		return
	}
	/*content-Data*/
	for _, finalContentResult := range finalContentResults {
		contentResult.ContentKey = finalContentResult.ContentKey
		contentResult.Duration = finalContentResult.Duration
		contentResult.Status = finalContentResult.Status
		/*Textual-Data*/
		contentResult.TextualData.PrimaryInfo.ContentType = finalContentResult.ContentType
		contentResult.TextualData.PrimaryInfo.OriginalTitle = finalContentResult.OriginalTitle
		contentResult.TextualData.PrimaryInfo.AlternativeTitle = finalContentResult.AlternativeTitle
		contentResult.TextualData.PrimaryInfo.ArabicTitle = finalContentResult.ArabicTitle
		contentResult.TextualData.PrimaryInfo.TransliteratedTitle = finalContentResult.TransliteratedTitle
		contentResult.TextualData.PrimaryInfo.Notes = finalContentResult.Notes
		contentResult.TextualData.PrimaryInfo.IntroStart = finalContentResult.IntroStart
		contentResult.TextualData.PrimaryInfo.OutroStart = finalContentResult.OutroStart
		/*Fetch content_geners*/
		// var contentGenres []SeasonGenres
		var contentGenres []ContentGenre
		var finalContentGenre []interface{}
		var newContentGenres NewSeasonGenres
		if genreResult := db.Debug().Table("content_genre").Select("id,genre_id").Where("content_id=?", c.Param("id")).Scan(&contentGenres).Error; genreResult != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		for _, tagInfoIds := range contentGenres {
			// var contentSubgenre []SeasonSubgenre
			var contentSubgenre []ContentSubgenre
			if subgenreResult := db.Debug().Table("content_subgenre as cs").Select("subgenre_id").Where("content_genre_id=?", tagInfoIds.Id).Order("cs.order").Scan(&contentSubgenre).Error; subgenreResult != nil {
				l.JSON(c, http.StatusInternalServerError, serverError)
				return
			}
			var Subgenre []string
			for _, data := range contentSubgenre {
				Subgenre = append(Subgenre, data.SubgenreId)

			}
			newContentGenres.GenreId = tagInfoIds.GenreId
			newContentGenres.Id = tagInfoIds.Id
			newContentGenres.SubgenresId = Subgenre
			finalContentGenre = append(finalContentGenre, newContentGenres)
		}

		contentResult.TextualData.ContentGenres = finalContentGenre
	}
	/*content_variance*/

	ContentVariance := []ContentVariances{}
	for _, finalContentResult := range finalContentResults {
		var contentVariances ContentVariances
		contentVariances.Status = finalContentResult.VarianceStatus
		contentVariances.VideoContentId = finalContentResult.VideoContentId
		// contentVariances.YoutubeVideoId = finalContentResult.YoutubeVideoId
		//	contentVariances.VideoContentId = finalContentResult.VideoContentId
		contentVariances.LanguageType = common.ContentLanguageOriginTypesName(finalContentResult.LanguageType)
		if finalContentResult.HasOverlayPosterImage {
			contentVariances.OverlayPosterImage = os.Getenv("IMAGERY_URL") + finalContentResult.Id + "/" + finalContentResult.VarianceId + "/overlay-poster-image"
		}
		if finalContentResult.HasDubbingScript {
			contentVariances.DubbingScript = os.Getenv("IMAGERY_URL") + finalContentResult.Id + "/" + finalContentResult.VarianceId + "/dubbing-script"
		}
		if finalContentResult.HasSubtitlingScript {
			contentVariances.SubtitlingScript = os.Getenv("IMAGERY_URL") + finalContentResult.Id + "/" + finalContentResult.VarianceId + "/subtitling-script"
		}
		contentVariances.DubbingLanguage = finalContentResult.DubbingLanguage
		contentVariances.DubbingDialectId = finalContentResult.DubbingDialectId
		contentVariances.SubtitlingLanguage = finalContentResult.SubtitlingLanguage
		contentVariances.DigitalRightsType = finalContentResult.DigitalRightsType
		contentVariances.DigitalRightsStartDate = finalContentResult.DigitalRightsStartDate
		contentVariances.DigitalRightsEndDate = finalContentResult.DigitalRightsEndDate
		/*Fetch Digital_right_Regions*/
		var digitalRightsRegions []DigitalRightsRegions
		if countryError := db.Debug().Table("content_rights_country").Select("country_id").Where("content_rights_id=?", finalContentResult.RightsId).Scan(&digitalRightsRegions).Error; countryError != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var RegionRights []int
		for _, idarr := range digitalRightsRegions {
			RegionRights = append(RegionRights, idarr.CountryId)
		}
		contentVariances.DigitalRightsRegions = RegionRights
		if len(RegionRights) < 1 {
			buffer := make([]int, 0)
			contentVariances.DigitalRightsRegions = buffer
		}
		contentVariances.SchedulingDateTime = finalContentResult.SchedulingDateTime
		contentVariances.CreatedBy = finalContentResult.CreatedByUserId
		//Fetch Publishing_platforms
		var targetPlatforms []PlaybackItemTargetPlatform
		if Errorplatforms := db.Debug().Table("playback_item_target_platform").Select("target_platform").Where("playback_item_id=?", finalContentResult.PlaybackItemId).Scan(&targetPlatforms).Error; Errorplatforms != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var platforms []int
		for _, idarr := range targetPlatforms {
			platforms = append(platforms, idarr.TargetPlatform)
		}
		contentVariances.PublishingPlatforms = platforms
		if len(platforms) < 1 {
			buffer := make([]int, 0)
			contentVariances.PublishingPlatforms = buffer
		}
		/*fetch product data*/
		var rightProduct []RightProduct
		if productResult := db.Debug().Table("rights_product").Select("product_name").Where("rights_id=?", finalContentResult.RightsId).Order("product_name asc").Scan(&rightProduct).Error; productResult != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var products []int
		for _, productNumbers := range rightProduct {
			products = append(products, productNumbers.ProductName)
		}
		contentVariances.Products = products
		if len(products) < 1 {
			buffer := make([]int, 0)
			contentVariances.Products = buffer
		}
		/*Fetch Subscriptions*/
		var subscriptionPlans []SubscriptionPlans
		if subscriptionError := db.Debug().Table("content_rights_plan").Select("subscription_plan_id").Where("rights_id=?", finalContentResult.RightsId).Scan(&subscriptionPlans).Error; subscriptionError != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var subscriptions []int
		for _, idarr := range subscriptionPlans {
			subscriptions = append(subscriptions, idarr.SubscriptionPlanId)
		}
		contentVariances.SubscriptionPlans = subscriptions
		if len(subscriptions) < 1 {
			buffer := make([]int, 0)
			contentVariances.SubscriptionPlans = buffer
		}
		contentVariances.IntroDuration = finalContentResult.VarianceIntroDuration
		contentVariances.IntroStart = finalContentResult.VarianceIntroStart
		/* Fetch Variance_Trailers*/
		var varianceTrailers []VarianceTrailers
		if varianceTrailersError := db.Debug().Raw("select * from variance_trailer vt where content_variance_id=? order by vt.order asc", finalContentResult.VarianceId).Find(&varianceTrailers).Error; varianceTrailersError != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var variance VarianceTrailers
		finalVarianceResult := []VarianceTrailers{}
		for _, varianceData := range varianceTrailers {
			variance.Order = varianceData.Order
			variance.VideoTrailerId = varianceData.VideoTrailerId
			variance.EnglishTitle = varianceData.EnglishTitle
			variance.ArabicTitle = varianceData.ArabicTitle
			variance.Duration = varianceData.Duration
			variance.HasTrailerPosterImage = varianceData.HasTrailerPosterImage
			var url = os.Getenv("IMAGERY_URL")
			variance.TrailerposterImage = url + finalContentResult.Id + "/" + finalContentResult.VarianceId + "/" + varianceData.Id + "/trailer-poster-image"
			variance.Id = varianceData.Id
			finalVarianceResult = append(finalVarianceResult, variance)
		}
		contentVariances.VarianceTrailers = finalVarianceResult
		contentVariances.Id = finalContentResult.VarianceId
		ContentVariance = append(ContentVariance, contentVariances)
	}
	contentResult.TextualData.ContentVariances = ContentVariance
	// fetching content cast details
	for _, finalContentResult := range finalContentResults {
		contentResult.TextualData.Cast.MainActorId = finalContentResult.MainActorId
		contentResult.TextualData.Cast.MainActressId = finalContentResult.MainActressId
		/* Fetch content_actors*/
		var contentActor []ContentActor
		if actorResult := db.Debug().Table("content_actor").Select("actor_id").Where("cast_id=?", finalContentResult.CastId).Scan(&contentActor).Error; actorResult != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var actors []string
		for _, actorIds := range contentActor {
			actors = append(actors, actorIds.ActorId)
		}
		contentResult.TextualData.Cast.Actors = actors
		if len(actors) < 1 {
			buffer := make([]string, 0)
			contentResult.TextualData.Cast.Actors = buffer
		}
		/* Fetch content_writers*/
		var contentWriters []ContentWriter
		if writerResult := db.Debug().Table("content_writer").Select("writer_id").Where("cast_id=?", finalContentResult.CastId).Scan(&contentWriters).Error; writerResult != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var writers []string
		for _, writerIds := range contentWriters {
			writers = append(writers, writerIds.WriterId)
		}
		contentResult.TextualData.Cast.Writers = writers
		if len(writers) < 1 {
			buffer := make([]string, 0)
			contentResult.TextualData.Cast.Writers = buffer
		}
		/* Fetch content_directors*/
		var contentDirectors []ContentDirector
		if directorResult := db.Debug().Table("content_director").Select("director_id").Where("cast_id=?", finalContentResult.CastId).Scan(&contentDirectors).Error; directorResult != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var directors []string
		for _, directorIds := range contentDirectors {
			directors = append(directors, directorIds.DirectorId)
		}
		contentResult.TextualData.Cast.Directors = directors
		if len(directors) < 1 {
			buffer := make([]string, 0)
			contentResult.TextualData.Cast.Directors = buffer
		}
		/* Fetch content_singers*/
		var contentSingers []ContentSinger
		if singerResult := db.Debug().Table("content_singer").Select("singer_id").Where("music_id=?", finalContentResult.MusicId).Scan(&contentSingers).Error; singerResult != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var singers []string
		for _, singerIds := range contentSingers {
			singers = append(singers, singerIds.SingerId)
		}
		contentResult.TextualData.Music.Singers = singers
		if len(singers) < 1 {
			buffer := make([]string, 0)
			contentResult.TextualData.Music.Singers = buffer
		}
		/* Fetch content_music_composers*/
		var contentMusicComposers []ContentMusicComposer
		db.Debug().Table("content_music_composer").Select("music_composer_id").Where("music_id=?", finalContentResult.MusicId).Scan(&contentMusicComposers)
		var musicComposers []string
		for _, composerIds := range contentMusicComposers {
			musicComposers = append(musicComposers, composerIds.MusicComposerId)
		}
		contentResult.TextualData.Music.MusicComposers = musicComposers
		if len(musicComposers) < 1 {
			buffer := make([]string, 0)
			contentResult.TextualData.Music.MusicComposers = buffer
		}
		/* Fetch content_writers*/
		var contentSongWriters []ContentSongWriter
		db.Debug().Table("content_song_writer").Select("song_writer_id").Where("music_id=?", finalContentResult.MusicId).Scan(&contentSongWriters)
		var songWriters []string
		for _, songWritersIds := range contentSongWriters {
			songWriters = append(songWriters, songWritersIds.SongWriterId)
		}
		contentResult.TextualData.Music.SongWriters = songWriters
		if len(songWriters) < 1 {
			buffer := make([]string, 0)
			contentResult.TextualData.Music.SongWriters = buffer
		}
		/*Fetch tag_info*/
		var contentTagInfo []ContentTag
		db.Debug().Table("content_tag").Select("textual_data_tag_id").Where("tag_info_id=?", finalContentResult.TagInfoId).Scan(&contentTagInfo)
		var tagInfo []string
		for _, tagInfoIds := range contentTagInfo {
			tagInfo = append(tagInfo, tagInfoIds.TextualDataTagId)
		}
		contentResult.TextualData.TagInfo.Tags = tagInfo
		if len(tagInfo) < 1 {
			buffer := make([]string, 0)
			contentResult.TextualData.TagInfo.Tags = buffer
		}
		contentResult.TextualData.AboutTheContent.OriginalLanguage = finalContentResult.OriginalLanguage
		contentResult.TextualData.AboutTheContent.Supplier = finalContentResult.Supplier
		contentResult.TextualData.AboutTheContent.AcquisitionDepartment = finalContentResult.AcquisitionDepartment
		contentResult.TextualData.AboutTheContent.EnglishSynopsis = finalContentResult.EnglishSynopsis
		contentResult.TextualData.AboutTheContent.ArabicSynopsis = finalContentResult.ArabicSynopsis
		contentResult.TextualData.AboutTheContent.ProductionYear = finalContentResult.ProductionYear
		contentResult.TextualData.AboutTheContent.ProductionHouse = finalContentResult.ProductionHouse
		contentResult.TextualData.AboutTheContent.AgeGroup = finalContentResult.AgeGroup
		contentResult.TextualData.AboutTheContent.IntroDuration = finalContentResult.AboutIntroDuration
		contentResult.TextualData.AboutTheContent.IntroStart = finalContentResult.AboutIntroStart
		contentResult.TextualData.AboutTheContent.OutroDuration = finalContentResult.AboutOutroDuration
		contentResult.TextualData.AboutTheContent.OutroStart = finalContentResult.AboutOutroStart
		/*Fetch Production_country*/
		var productionCountry []ProductionCountry
		if productionCountryResult := db.Debug().Table("production_country").Select("country_id").Where("about_the_content_info_id=?", finalContentResult.AboutTheContentInfoId).Order("country_id asc").Scan(&productionCountry).Error; productionCountryResult != nil {
			l.JSON(c, http.StatusInternalServerError, serverError)
			return
		}
		var countries []int
		for _, prcountries := range productionCountry {
			countries = append(countries, prcountries.CountryId)
		}
		contentResult.TextualData.AboutTheContent.ProductionCountries = countries
		fmt.Println(contentResult.TextualData.AboutTheContent.ProductionCountries, "coutnr.........", countries)
		// if len(tagInfo) < 1 {
		// 	buffer := make([]int, 0)
		// 	contentResult.TextualData.AboutTheContent.ProductionCountries = buffer
		// }
		/*SeoDetails*/
		contentResult.TextualData.SeoDetails.EnglishMetaTitle = finalContentResult.EnglishMetaTitle
		contentResult.TextualData.SeoDetails.ArabicMetaTitle = finalContentResult.ArabicMetaTitle
		contentResult.TextualData.SeoDetails.EnglishMetaDescription = finalContentResult.EnglishMetaDescription
		contentResult.TextualData.SeoDetails.ArabicMetaDescription = finalContentResult.ArabicMetaDescription
		/*non_textual Data*/
		if finalContentResult.HasPosterImage {
			contentResult.NonTextualData.PosterImage = os.Getenv("IMAGERY_URL") + finalContentResult.Id + "/poster-image"
		}
		if finalContentResult.HasDetailsBackground {
			contentResult.NonTextualData.DetailsBackground = os.Getenv("IMAGERY_URL") + finalContentResult.Id + "/details-background"
		}
		if finalContentResult.HasDetailsBackground {
			contentResult.NonTextualData.MobileDetailsBackground = os.Getenv("IMAGERY_URL") + finalContentResult.Id + "/mobile-details-background"
		}
		contentResult.Id = finalContentResult.Id

		fmt.Println(finalContentResult.EnglishMetaTitle)
	}

	l.JSON(c, http.StatusOK, gin.H{"data": contentResult})

}

/*Uploade image Based on season Id*/
func EpisodeFileUPload(request Episodes, episodeId string, contentId string, seasonId string) {
	bucketName := os.Getenv("S3_BUCKET")
	var newarr []string
	newarr = append(newarr, request.NonTextualData.PosterImage)
	newarr = append(newarr, request.NonTextualData.DubbingScript)
	newarr = append(newarr, request.NonTextualData.SubtitlingScript)
	newarr = append(newarr, request.NonTextualData.MobileDetailsBackground)
	for k := 0; k < len(newarr); k++ {
		item := newarr[k]
		filetrim := strings.Split(item, "_")
		Destination := contentId + "/" + seasonId + "/" + episodeId + "/" + filetrim[0]
		source := bucketName + "/" + "temp/" + item
		s, _ := session.NewSession(&aws.Config{
			Region: aws.String(os.Getenv("S3_REGION")),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("S3_ID"),     // id
				os.Getenv("S3_SECRET"), // secret
				""),                    // token can be left blank for now
		})
		/* Copy object from one directory to another*/
		svc := s3.New(s)
		input := &s3.CopyObjectInput{
			Bucket:     aws.String(bucketName),
			CopySource: aws.String(source),
			Key:        aws.String(Destination),
		}
		result, err := svc.CopyObject(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeObjectNotInActiveTierError:
					fmt.Println(s3.ErrCodeObjectNotInActiveTierError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Println(result)
		url := "https://" + bucketName + ".s3.ap-south-1.amazonaws.com/" + Destination
		// don't worry about errors
		response, e := http.Get(url)
		if e != nil {
			log.Fatal(e)
		}
		defer response.Body.Close()

		//open a file for writing
		file, err := os.Create(filetrim[0])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal(err)
		}
		errorr := SizeUploadFileToS3(s, filetrim[0], episodeId, contentId, seasonId)
		if errorr != nil {
			fmt.Println("error in uploading size upload", errorr)
		}
		fmt.Println("Success!")
	}
}

// SizeUploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func SizeUploadFileToS3(s *session.Session, fileName string, episodeId string, contentId string, seasonId string) error {
	// open the file for use
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// get the file size and read
	// the file content into a buffer
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	sizeValue := [17]string{
		"100x100/",
		"150x150/",
		"200x200/",
		"250x250/",
		"270x270/",
		"300x300/",
		"420x420/",
		"450x450/",
		"570x570/",
		"600x600/",
		"620x620/",
		"800x384/",
		"800x800/",
		"811x811/",
		"900x900/",
		"2048x670/",
		"1125x240/",
	}
	var er error
	for i := 0; i < len(sizeValue); i++ {
		s3file := sizeValue[i] + contentId + "/" + seasonId + "/" + episodeId + "/" + fileName
		_, er = s3.New(s).PutObject(&s3.PutObjectInput{
			Bucket:               aws.String(os.Getenv("S3_BUCKET")),
			Key:                  aws.String(s3file),
			ACL:                  aws.String("public-read"),
			Body:                 bytes.NewReader(buffer),
			ContentLength:        aws.Int64(size),
			ContentType:          aws.String(http.DetectContentType(buffer)),
			ContentDisposition:   aws.String("attachment"),
			StorageClass:         aws.String("STANDARD"),
			ServerSideEncryption: aws.String("AES256"),
		})
		if er != nil {
			fmt.Println("Unable to upload ", fileName, er)
		}
		fmt.Printf("Successfully uploaded %q", fileName)
	}
	os.Remove(fileName)
	return er
}

func EpisodeFileUploadGcp(request Episodes, episodeId string, contentId string, seasonId string) {
	bucketName := os.Getenv("BUCKET_NAME")
	var newarr []string
	newarr = append(newarr, request.NonTextualData.PosterImage)
	newarr = append(newarr, request.NonTextualData.DubbingScript)
	newarr = append(newarr, request.NonTextualData.SubtitlingScript)
	newarr = append(newarr, request.NonTextualData.MobileDetailsBackground)

	ctx := context.Background()

	client, gcperr := getGCPClient()
	if gcperr != nil {
		fmt.Println("from gcp Connection", gcperr)
		// return gcperr
	}
	defer client.Close()

	for k := 0; k < len(newarr); k++ {
		item := newarr[k]
		filetrim := strings.Split(item, "_")
		Destination := contentId + "/" + seasonId + "/" + episodeId + "/" + filetrim[0]
		source := "temp/" + item // Assuming temp is a local directory.

		// Copy object from one bucket to another.
		src := client.Bucket(bucketName).Object(source)
		attrs, err := src.Attrs(ctx)
		_ = attrs
		if err != nil {
			// Handle the case where the source object doesn't exist
			fmt.Printf("Source object does not exist: %v\n", err)

			// Modify the source path if needed
			filetrims := strings.Split(item, "/")
			source = contentId + "/" + seasonId + "/" + episodeId + "/" + filetrims[len(filetrims)-1]
			Destination = contentId + "/" + seasonId + "/" + episodeId + "/" + filetrims[len(filetrims)-1]
			src = client.Bucket(bucketName).Object(source)
			filetrim[0] = filetrims[len(filetrims)-1]

			// Retry the Attrs call
			attrs, err = src.Attrs(ctx)
			if err != nil {
				// Handle the case where the modified source object also doesn't exist
				fmt.Printf("Modified source object does not exist: %v\n", err)
				filetrims := strings.Split(item, "_")
				source = contentId + "/" + seasonId + "/" + episodeId + "/" + filetrims[0]
				Destination = contentId + "/" + seasonId + "/" + episodeId + "/" + filetrims[0]
				src = client.Bucket(bucketName).Object(source)
				filetrim[0] = filetrims[0]
				attrs, err := src.Attrs(ctx)
				_ = attrs
				if err != nil {
					fmt.Println("err-------No image", err)
					fmt.Println("err-------", filetrims[0], source)
					fmt.Println("destination-------", Destination, "-=-=-", source)
				}
			}
		}
		dst := client.Bucket(bucketName).Object(Destination)
		if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
			fmt.Println("CopyObject failed: ", err)
		}

		// Generate the public URL for the uploaded file.
		// url := os.Getenv("IMAGERY_URL") + bucketName + "/" + Destination
		url := os.Getenv("IMAGERY_URL") + Destination

		// Don't worry about errors.
		response, e := http.Get(url)
		if e != nil {
			fmt.Println(e)
		}
		defer response.Body.Close()
		if e == nil {
			// Open a file for writing.
			file, err := os.Create(filetrim[0])
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()
			if err == nil {
				// Use io.Copy to dump the response body to the file.
				_, err = io.Copy(file, response.Body)
				if err != nil {
					fmt.Println(err)
				}

				go SizeUploadFileToGcp(client, filetrim[0], episodeId, contentId, seasonId, url)
				// if err != nil {
				// 	fmt.Println("error in uploading size upload", err)
				// }
			}
			// os.Remove(filetrim[0])
		}

		fmt.Println("Success!")
	}
}

func SizeUploadFileToGcp(client *storage.Client, fileName string, episodeId string, contentId string, seasonId string, fileUrl string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	sizeValue := [17]string{
		"100x100/",
		"150x150/",
		"200x200/",
		"250x250/",
		"270x270/",
		"300x300/",
		"420x420/",
		"450x450/",
		"570x570/",
		"600x600/",
		"620x620/",
		"800x384/",
		"800x800/",
		"811x811/",
		"900x900/",
		"2048x670/",
		"1125x240/",
	}
	ctx := context.Background()
	bucketName := os.Getenv("BUCKET_NAME")

	for i := 0; i < len(sizeValue); i++ {
		filetrim := strings.Split(sizeValue[i], "/")
		filetri := strings.Split(filetrim[0], "x")
		width := filetri[0]
		height := filetri[1]
		if i == 0 {
			width = "200"
			height = "200"
		}
		// https://msapiuat-image.z5.com/crop?width=200&height=200&url=https://weyyak-content-qa.engro.in/2b7d164d-eddd-4b6d-9d9c-84df62ccf01b/28e91598-4b43-40aa-89d6-daadb31ef82b/poster-image
		// Get the URL of the uploaded object in GCS
		url := os.Getenv("RESIZE_IMAGE_URL") + "width=" + width + "&height=" + height + "&url=" + fileUrl
		fmt.Println("urlurlurlurl", url)
		method := "GET"

		client1 := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
			// return
		}
		res, err := client1.Do(req)
		if err != nil {
			fmt.Println(err)
			// return
		}
		defer res.Body.Close()
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		_, err = io.Copy(file, res.Body)
		if err != nil {
			fmt.Println(err)
		}

		file, err = os.Open(fileName)
		if err != nil {
			fmt.Println("err1", err)
		}
		defer file.Close()

		// Get file size and read content into a buffer
		fileInfo, _ := file.Stat()
		var size = fileInfo.Size()
		buffer := make([]byte, size)
		file.Read(buffer)

		s3file := sizeValue[i] + contentId + "/" + seasonId + "/" + episodeId + "/" + fileName
		fmt.Println("s3files3file", s3file)
		wc := client.Bucket(bucketName).Object(s3file).NewWriter(ctx)
		wc.ContentType = http.DetectContentType(buffer)
		wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

		_, err = wc.Write(buffer)
		if err != nil {
			return fmt.Errorf("Unable to upload %s: %v", fileName, err)
		}

		if err := wc.Close(); err != nil {
			return fmt.Errorf("Unable to close writer for %s: %v", fileName, err)
		}

		fmt.Printf("Successfully uploaded %q\n", fileName)
	}

	return nil
}

func getGCPClient() (*storage.Client, error) {
	data := map[string]interface{}{
		// "client_id":       "764086051850-6qr4p6gpi6hn506pt8ejuq83di341hur.apps.googleusercontent.com",
		// "client_secret":    "d-FL95Q19q7MQmFpd7hHD0Ty",
		// "quota_project_id": "engro-project-392708",
		// "refresh_token":    "1//0gCu2SwEAITTxCgYIARAAGBASNwF-L9IrXoW2jiRehyvfOj0yt3jnt5FXmYdlmkXXNIDjKzt5O1a3USJtclNE6sMSlr_W_Mw4xes",
		// "type":             "authorized_user",

		"type":                        os.Getenv("TYPE"),
		"project_id":                  os.Getenv("PROJECT_ID"),
		"private_key_id":              os.Getenv("PRIVATE_KEY_ID"),
		"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCrrOzb27NS3vj2\n/mspxBm/22giPXnCND0DYKLRcQuzsYv0elnbqmcGSJfCqs9C59PGJzRS6RV+ie/N\nhgflwUx3I3zMtd1fgVsjZFYYwmIDsaShkaK0B3eYWGXqYw8swG67qctk3GubNYNl\n04p6qIsUZgTpmH8Jjvn7abwntDvfxYbK9QVC2rjXzzA0pABfrgo3CVZJiXO8t4Ly\na8P5jkNY8Uze6Us+L0XxYL8T08vD09Cqde2kEEMBM4R6okdtgj3Vp8G287Dj0OUe\nMBGoj6YgsH32ZOXxxeQp003tMETaJbwQx6HTKr4N0CDwQnQkMbVUN2N8wOd4e0+2\nG+U36cGVAgMBAAECggEAFn+JDZ8TNyRleD5gs46G2VqFoRxxXSlqEuE9NTlyu8/k\nHtv8nrRhirSaFDbnsUWfE/QwqpTv7i9hhTZayUS1zVSR7GSrvZ0UNo/Vq1T+HKx7\n03i52+IGov54DL7X+ZjBFPLsPCxEJd5eI/Vpy9KpYg5PTSsLqv2udmulmYZzOktP\nYeV/qAaV/h/uQa+yTkxz9q0lixganx+ZSiC/3iTLwQLTI+Em8ayjVcIGQ/A9j6X1\nVCOxHBvy3bcIgZe+ZImwoWvko8ryaHWrdCKz7zVgXPZ9aT6B+VW0qqJGsHS0F5m7\nK0EC8fkdMlRufEiw6DChWUmspg7FYNW3fL7boAXemQKBgQDoOzEwr9khlO32ZXSs\nqIKRGNoL5pZPekVHPc/LI6713Vwg4g52xmtT0ZwgjkUB9QF4CimYVGLHytZG6P0G\nSBAdf4JMeeuBkJtmkXnYdJAlbwRNTiHWz409yAJ9hIyPafLZFKYxMLAqj95wnBxc\nMGq9accLaLIUtG8WGfSUrs5fwwKBgQC9PxTrAMl+ewm2O+a86du+BGsiy8fUvQZX\nJ9xayx9ARjEJXv1cgD4z59mQDn6gzBLrDcH+KY2ZSZUmvPof5LkXUlXXplJWh1Qj\nYvpMx2IOdu2OFFfydtyvq/JbXaEMrvUGU3+pvCF1e7Wxf+jlCTZM4yKwg1Ba9FyT\nCUaPlJFbxwKBgC0wv4y622TWh0voSEEE9Ytoq52fPGaw42ROme3svrInZjMb6jag\nu+fupRQMu077L1L9n0R+P06joPjhg8NCKKik1GUvYG2xBxx5eJ1vaVFvfgXRC3Ky\npsh78Egej/+kXVZy1zhBQja2ElIVfstNvKepOst0jxrKVceWO2rnbU9jAoGAdDH6\nNvxpuyXyZZjL6GwyRq5R1bCHRqC09uh7jKewzXcLfrR7HcOD7bzKQYAU0cfbScVN\nui9rSJX8ZSec794woxgjqt/tKEG5MG0CQAgftb/hxd3Jzg6bG6WYje6kBrSZr0Ov\nW9kuNgM6IPznU1FfrL+9OeG2gdIN0R3d3CSdR1sCgYBxpi1DXXBCXVeU1wIb9XBA\nwewiFSAabF/UtiF7CHkGSMN1lMe/R1AFKM8Irrqbbm0jl00BZ5fgVYV/wVaYZtDw\nPZQmGeO3yGi6FanLnBaxE/bKjk+RkaORM8QoaYGghX59TNoFzHNE1rF0w1lMdrlN\nnsFelOtLls3xNrtNNMxHeg==\n-----END PRIVATE KEY-----\n",
		"client_email":                os.Getenv("CLIENT_EMAIL"),
		"client_id":                   os.Getenv("CLIENT_ID"),
		"auth_uri":                    os.Getenv("AUTH_URI"),
		"token_uri":                   os.Getenv("TOKEN_URI"),
		"auth_provider_x509_cert_url": os.Getenv("AUTH_PROVIDER_X509_CERT_URL"),
		"client_x509_cert_url":        os.Getenv("CLIENT_X509_CERT_URL"),
		"universe_domain":             os.Getenv("UNIVERSE_DOMAIN"),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	creds, err := google.CredentialsFromJSON(ctx, jsonData, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		fmt.Println("Error creating credentials:", err)
	}
	client, err := storage.NewClient(ctx, option.WithCredentials(creds))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return client, err
}
