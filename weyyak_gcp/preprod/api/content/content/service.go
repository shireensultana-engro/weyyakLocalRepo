package content

import (
	"bytes"
	"content/common"
	"content/episode"
	"content/fragments"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"github.com/thanhpk/randstr"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const DateTimeFormat string = "2006-01-02T15:04:05.000Z"

type ContentService struct{}

// All the services should be protected by auth token
func (hs *ContentService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	qrg := r.Group("/api")
	qrg.Use(common.ValidateToken())
	//	qrg.POST("/contents/onetier/published", hs.CreateOnetierContent)
	qrg.GET("/contents/multitier/:result", hs.GetMultitierContentDetails)
	qrg.GET("/episodes/:id", hs.GetEpisodeDetailsBYepisodeId)
	qrg.POST("/seasons/published", hs.CreateSeason)
	qrg.POST("/seasons/published/:id", hs.CreateSeason)
	qrg.POST("/seasons/draft", hs.DraftSeason)
	qrg.POST("/seasons/draft/:id", hs.DraftSeason)
	qrg.GET("/contents", hs.GetAllContentDetails)
	r.GET("/api/contentsync", hs.PageSyncWithContentId)
	//	qrg.POST("/seasonvarince/published", hs.CreateOrUpdateSeasonVariance)
	qrg.POST("/nontextual/content/poster-image", hs.UploadMenuPosterImageGcp)
	qrg.POST("/nontextual/content/details-background", hs.DetailsBackgroundImageGcp)
	qrg.POST("/nontextual/content/mobile-details-background", hs.MobileDetailsBackgroundImageGcp)
	qrg.POST("/nontextual/content/overlay-poster-image", hs.OverlayPosterImageGcp)
	qrg.POST("/nontextual/content/menu-poster-image", hs.MenuPosterImageGcp)
	qrg.POST("/nontextual/content/mobile-menu", hs.MobileMenuGcp)
	qrg.POST("/nontextual/content/mobile-menu-poster-image", hs.MobileMenuPosterImageGcp)
	qrg.POST("/nontextual/content/trailer-poster-image", hs.TrailerPosterImageGcp)

}

// CreateOnetierContent -  Create onetier content
// POST /onetier/published
// @Summary Create onetier content details
// @Description Create onetier content details
// @Tags Content
// @Accept  json
// @Produce  json
// @Param body body content.OnetierContentRequest true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /api/contents/onetier/published [post]
func (hs *ContentService) CreateOnetierContent(c *gin.Context) {
	var request OnetierContentRequest
	var content Content
	var primaryInfo ContentPrimaryInfo
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("DB").(*gorm.DB)
	MusicId := ""
	TaginfoId := ""
	primaryInfo.OriginalTitle = request.TextualData.PrimaryInfo.OriginalTitle
	primaryInfo.AlternativeTitle = request.TextualData.PrimaryInfo.AlternativeTitle
	primaryInfo.ArabicTitle = request.TextualData.PrimaryInfo.ArabicTitle
	primaryInfo.TransliteratedTitle = request.TextualData.PrimaryInfo.TransliteratedTitle
	primaryInfo.Notes = request.TextualData.PrimaryInfo.Notes
	if err := db.Debug().Create(&primaryInfo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}
	fmt.Println("9")
	var aboutTheContent AboutTheContentInfo
	aboutTheContent.OriginalLanguage = request.TextualData.AboutTheContent.OriginalLanguage
	aboutTheContent.Supplier = request.TextualData.AboutTheContent.Supplier
	aboutTheContent.AcquisitionDepartment = request.TextualData.AboutTheContent.AcquisitionDepartment
	aboutTheContent.EnglishSynopsis = request.TextualData.AboutTheContent.EnglishSynopsis
	aboutTheContent.ArabicSynopsis = request.TextualData.AboutTheContent.ArabicSynopsis
	aboutTheContent.ProductionYear = request.TextualData.AboutTheContent.ProductionYear
	aboutTheContent.ProductionHouse = request.TextualData.AboutTheContent.ProductionHouse
	aboutTheContent.AgeGroup = request.TextualData.AboutTheContent.AgeGroup
	if err := db.Debug().Create(&aboutTheContent).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}
	fmt.Println("10")
	var productionCountry ProductionCountry
	if request.TextualData.AboutTheContent.ProductionCountries != nil && len(request.TextualData.AboutTheContent.ProductionCountries) > 0 {
		for _, country := range request.TextualData.AboutTheContent.ProductionCountries {
			productionCountry.AboutTheContentInfoId = aboutTheContent.Id
			productionCountry.CountryId = country
			if err := db.Debug().Create(&productionCountry).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
				return
			}
		}
	}
	fmt.Println("11")
	var contentKey Content
	db.Debug().Select("content_key").Order("content_key desc").Limit(1).Find(&contentKey)

	content.ContentKey = contentKey.ContentKey + 1
	content.ContentType = request.TextualData.PrimaryInfo.ContentType
	content.Status = 1
	content.ContentTier = 1
	content.PrimaryInfoId = primaryInfo.ID
	content.AboutTheContentInfoId = aboutTheContent.Id
	content.CastId = ""
	content.CreatedByUserId = "00000000-0000-0000-0000-000000000000"
	content.MusicId = MusicId
	content.TagInfoId = TaginfoId
	content.CreatedAt = time.Now()
	content.EnglishMetaTitle = request.TextualData.SeoDetails.EnglishMetaTitle
	content.ArabicMetaTitle = request.TextualData.SeoDetails.ArabicMetaTitle
	content.EnglishMetaDescription = request.TextualData.SeoDetails.EnglishMetaDescription
	content.ArabicMetaDescription = request.TextualData.SeoDetails.ArabicMetaDescription
	if err := db.Debug().Create(&content).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil})
		return
	} else {
		var contentVariance ContentVariance
		if request.TextualData.ContentVariances != nil {
			for j, variance := range request.TextualData.ContentVariances {
				var contentRights ContentRights
				contentRights.DigitalRightsType = variance.DigitalRightsType
				//	DRSDate, _ := time.Parse(DateTimeFormat, variance.DigitalRightsStartDate)
				contentRights.DigitalRightsStartDate = variance.DigitalRightsStartDate
				//	DREDate, _ := time.Parse(DateTimeFormat, variance.DigitalRightsEndDate)
				contentRights.DigitalRightsEndDate = variance.DigitalRightsEndDate
				if err := db.Debug().Create(&contentRights).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
					return
				}
				fmt.Println("13")
				if variance.DigitalRightsRegions != nil && len(variance.DigitalRightsRegions) > 0 {
					var contentRightsCountry ContentRightsCountry
					var contentRightsCountrys []interface{}
					for _, country := range variance.DigitalRightsRegions {
						contentRightsCountry.ContentRightsId = contentRights.Id
						contentRightsCountry.CountryId = country
						contentRightsCountrys = append(contentRightsCountrys, contentRightsCountry)
					}
					err = gormbulk.BulkInsert(db, contentRightsCountrys, common.BulkInsertLimit)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
						return
					}
				}
				fmt.Println("14")
				var contentRightsPlan ContentRightsPlan
				var contentRightsPlans []interface{}
				if variance.SubscriptionPlans != nil && len(variance.SubscriptionPlans) > 0 {
					for _, plan := range variance.SubscriptionPlans {
						contentRightsPlan.RightsId = contentRights.Id
						contentRightsPlan.SubscriptionPlanId = plan
						contentRightsPlans = append(contentRightsPlans, contentRightsPlan)
					}
					err = gormbulk.BulkInsert(db, contentRightsPlans, common.BulkInsertLimit)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
						return
					}
				}
				fmt.Println("15")
				var rightsProduct RightsProduct
				var rightsProducts []interface{}
				if variance.Products != nil && len(variance.Products) > 0 {
					for _, product := range variance.Products {
						rightsProduct.RightsId = contentRights.Id
						rightsProduct.ProductName = product
						rightsProducts = append(rightsProducts, rightsProduct)
					}
					err = gormbulk.BulkInsert(db, rightsProducts, common.BulkInsertLimit)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
						return
					}
				}
				fmt.Println("16")
				var contentTranslationRequest ContentTranslationRequest
				contentTranslationRequest.LanguageType = common.LanguageOriginTypes(variance.LanguageType)
				contentTranslationRequest.DubbingLanguage = &variance.DubbingLanguage
				contentTranslationRequest.DubbingDialectId = &variance.DubbingDialectId
				contentTranslationRequest.SubtitlingLanguage = &variance.SubtitlingLanguage
				// message,StatusCode,duration:=GetVideoDuration(variance.VideoContentId)
				fmt.Println("18")
				var playbackItem PlaybackItem
				playbackItem.VideoContentId = variance.VideoContentId
				SDT, _ := time.Parse(DateTimeFormat, variance.SchedulingDateTime)
				playbackItem.SchedulingDateTime = SDT
				// playbackItem.CreatedByUserId = variance.CreatedBy
				playbackItem.CreatedByUserId = "00000000-0000-0000-0000-000000000000"
				playbackItem.TranslationId = ""
				playbackItem.RightsId = contentRights.Id
				// playbackItem.Duration = <-duration
				playbackItem.Duration = 0
				if err := db.Debug().Create(&playbackItem).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
					return
				}
				fmt.Println("19")
				var PITPlatform PlaybackItemTargetPlatform
				var PITPlatforms []interface{}
				if variance.PublishingPlatforms != nil && len(variance.PublishingPlatforms) > 0 {
					for _, platform := range variance.PublishingPlatforms {
						PITPlatform.RightsId = contentRights.Id
						PITPlatform.TargetPlatform = platform
						PITPlatform.PlaybackItemId = playbackItem.Id
						PITPlatforms = append(PITPlatforms, PITPlatform)
					}
					err = gormbulk.BulkInsert(db, PITPlatforms, common.BulkInsertLimit)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
						return
					}
				}
				fmt.Println("20")
				contentVariance.ContentId = content.Id
				contentVariance.PlaybackItemId = playbackItem.Id
				contentVariance.Order = j
				contentVariance.HasDubbingScript = false
				contentVariance.HasSubtitlingScript = false
				contentVariance.Status = variance.Status
				contentVariance.HasAllRights = variance.CountryCheck
				contentVariance.IntroStart = 0
				contentVariance.IntroDuration = 0
				if err := db.Debug().Create(&contentVariance).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
					return
				}
				fmt.Println("21")
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "Content Created Successfully.", "status": http.StatusOK})
		return
	}
}

// GetMultitierContentDetails -  Get multitier content details
// GET /api/contents/multitier/:result
// @Summary Get multitier content details
// @Description Get multitier content details by content id
// @Tags Content
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param id path string true "Content Id"
// @Success 200 {array} object c.JSON
// @Router /api/contents/multitier/{result} [get]
func (hs *ContentService) GetMultitierContentDetails(c *gin.Context) {
	AuthorizationRequired := c.MustGet("AuthorizationRequired")
	userid := c.MustGet("userid")
	if AuthorizationRequired == 1 || userid == "" || c.MustGet("is_back_office_user") == "false" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	serverError := common.ServerErrorResponse()
	notFound := common.NotFoundErrorResponse()
	db := c.MustGet("DB").(*gorm.DB)
	var content Content
	if err := db.Debug().Where("id=? and deleted_by_user_id is null", c.Param("result")).Find(&content).Error; err != nil {
		c.JSON(http.StatusNotFound, notFound)
		return
	}
	var multitierContentDetails MultitierContentDetails
	var primaryInfo PrimaryInfo
	var contentGenre ContentGenres
	var contentGenres []ContentGenres
	var contentSeason ContentSeasons
	contentSeasons := []ContentSeasons{}
	// var nonTextualData NonTextualData
	// var rights Rights
	var seoDetails SeoDetails
	//query result variables
	var multitierDetails MultitierContentQueryDetails
	var genres []ContentGeneresQueryDetails
	var seasonDetails []ContentSeasonsQueryDetails
	var episodeDetails []SeasonEpisodesQueryDetails
	if err := db.Debug().Table("content c").Select("c.content_key,null::numeric as duration,c.status,c.content_type,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title,cpi.notes,cpi.intro_start,cpi.outro_start,c.english_meta_title,c.arabic_meta_title,c.english_meta_description,c.arabic_meta_description,c.id").Joins("join content_primary_info cpi on cpi.id = c.primary_info_id").Where("c.id=?", c.Param("result")).Find(&multitierDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	multitierContentDetails.ContentKey = multitierDetails.ContentKey
	multitierContentDetails.Duration = multitierDetails.Duration
	multitierContentDetails.Status = multitierDetails.Status
	primaryInfo.ContentType = multitierDetails.ContentType
	primaryInfo.OriginalTitle = multitierDetails.OriginalTitle
	primaryInfo.AlternativeTitle = multitierDetails.AlternativeTitle
	primaryInfo.ArabicTitle = multitierDetails.ArabicTitle
	primaryInfo.TransliteratedTitle = multitierDetails.TransliteratedTitle
	primaryInfo.Notes = multitierDetails.Notes
	primaryInfo.IntroStart = nil //multitierDetails.IntroStart
	primaryInfo.OutroStart = nil //multitierDetails.OutroStart
	multitierContentDetails.PrimaryInfo = primaryInfo
	if err := db.Debug().Table("content_genre cg").Select("cg.genre_id,json_agg(cs.subgenre_id  order by cs.order)::varchar as subgenres_id,json_agg(cs.order  order by cs.order)::varchar as sub_genre_order,cg.id").Joins("join content_subgenre cs on cs.content_genre_id = cg.id join content c on c.id = cg.content_id").Where("cg.content_id = ?", c.Param("result")).Group("cg.genre_id,cg.id").Order("cg.order").Find(&genres).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	for _, genre := range genres {
		contentGenre.GenreId = genre.GenreId
		subGenres, err := JsonStringToStringSliceOrMap(genre.SubgenresId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		contentGenre.SubgenresId = subGenres
		contentGenre.Id = genre.Id
		contentGenres = append(contentGenres, contentGenre)
	}
	multitierContentDetails.ContentGenres = contentGenres
	if err := db.Debug().Table("season s").Select("s.content_id,s.season_key,s.status,s.modified_at,s.number,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title,cpi.notes,cpi.intro_start,cpi.outro_start,atci.original_language,atci.supplier,atci.acquisition_department,atci.english_synopsis,atci.arabic_synopsis,atci.production_year,atci.production_house,atci.age_group,atci.intro_duration,atci.intro_start as atci_intro_start,atci.outro_duration,atci.outro_start as atci_intro_start,ct.language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date,json_agg(crc.country_id)::varchar AS digital_rights_regions,s.created_by_user_id,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,s.arabic_meta_description,s.id,s.rights_id,s.about_the_content_info_id").Joins("join content c on c.id=s.content_id join content_primary_info cpi on cpi.id =s.primary_info_id join about_the_content_info atci on atci.id=s.about_the_content_info_id join content_translation ct on ct.id =s.translation_id join content_rights cr on cr.id = s.rights_id join content_rights_country crc on crc.content_rights_id = cr.id").Where("s.content_id =? and s.deleted_by_user_id is null", c.Param("result")).Group("s.content_id,s.season_key,s.status,s.modified_at,s.number,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title,cpi.notes,cpi.intro_start,cpi.outro_start,atci.original_language,atci.supplier,atci.acquisition_department,atci.english_synopsis,atci.arabic_synopsis,atci.production_year,atci.production_house,atci.age_group,atci.intro_duration,atci.intro_start,atci.outro_duration,atci.outro_start,ct.language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date,s.created_by_user_id,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,s.arabic_meta_description,s.id").Order("s.number").Find(&seasonDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	for _, season := range seasonDetails {
		var cast Cast
		var music Music
		var tagInfo TagInfo
		var aboutTheContent AboutTheContent
		var contentTranslation ContentTranslation
		seasonEpisodes := []SeasonEpisodes{}
		var rightsDetails Rights
		var castDetails CastQueryDetails
		var musicDetails MusicQueryDetails
		var tagsDetails TagInfoQueryDetails
		actorResult, writerResult, directorResult := make([]string, 0), make([]string, 0), make([]string, 0)
		singerResult, musicComposerResult, songWriterResult := make([]string, 0), make([]string, 0), make([]string, 0)
		contentSeason.ContentId = season.ContentId
		contentSeason.SeasonKey = season.SeasonKey
		now := time.Now().Unix()
		sDate := season.DigitalRightsStartDate.Unix()
		eDate := season.DigitalRightsEndDate.Unix()
		contentSeason.SubStatusName = nil
		contentSeason.Status = season.Status
		if season.Status == 1 || season.Status == 2 {
			contentSeason.StatusCanBeChanged = true
			contentSeason.Status = season.Status
		} else if season.Status == 3 {
			contentSeason.StatusCanBeChanged = false
			contentSeason.Status = 3
			subStatus := "Draft"
			contentSeason.SubStatusName = &subStatus
		}
		if sDate > now || eDate < now {
			contentSeason.StatusCanBeChanged = false
			contentSeason.Status = 2
			subStatus := "Digital Rights Exceeded"
			contentSeason.SubStatusName = &subStatus
		}
		contentSeason.ModifiedAt = season.ModifiedAt
		//primary info details
		var primaryInfo PrimaryInfo
		primaryInfo.SeasonNumber = season.Number
		primaryInfo.OriginalTitle = season.OriginalTitle
		primaryInfo.AlternativeTitle = season.AlternativeTitle
		primaryInfo.ArabicTitle = season.ArabicTitle
		primaryInfo.TransliteratedTitle = season.TransliteratedTitle
		primaryInfo.Notes = season.Notes
		primaryInfo.IntroStart = season.IntroStart
		primaryInfo.OutroStart = season.OutroStart
		contentSeason.PrimaryInfo = primaryInfo
		if err := db.Debug().Table("content_cast cc").Select("cc.main_actor_id,cc.main_actress_id,json_agg(ca.actor_id)::varchar as actors,json_agg(cw.writer_id)::varchar as writers,json_agg(cd.director_id)::varchar as directors").Joins("join season s on s.cast_id=cc.id full outer join content_actor ca on ca.cast_id =cc.id full outer join content_writer cw on cw.cast_id =cc.id full outer join content_director cd on cd.cast_id =cc.id").Where("s.id=?", season.Id).Group("cc.main_actor_id,cc.main_actress_id").Find(&castDetails).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		cast.MainActorId = castDetails.MainActorId
		cast.MainActressId = castDetails.MainActressId
		if castDetails.Actors != "" {
			actors, _ := JsonStringToStringSliceOrMap(castDetails.Actors)
			actorResult = RemoveDuplicateStringValues(actors)
			if actorResult[0] == "" {
				actorResult = make([]string, 0)
			}
		}
		cast.Actors = actorResult
		if castDetails.Writers != "" {
			writers, _ := JsonStringToStringSliceOrMap(castDetails.Writers)
			writerResult = RemoveDuplicateStringValues(writers)
			if writerResult[0] == "" {
				writerResult = make([]string, 0)
			}
		}
		cast.Writers = writerResult
		if castDetails.Directors != "" {
			directors, _ := JsonStringToStringSliceOrMap(castDetails.Directors)
			directorResult = RemoveDuplicateStringValues(directors)
			if directorResult[0] == "" {
				directorResult = make([]string, 0)
			}
		}
		cast.Directors = directorResult
		contentSeason.Cast = cast
		if err := db.Debug().Table("season s").Select("json_agg(cs.singer_id)::varchar as singers,json_agg(cmc.music_composer_id)::varchar as music_composers,json_agg(csw.song_writer_id)::varchar as song_writers").Joins("full outer join content_singer cs on cs.music_id =s.music_id full outer join content_music_composer cmc on cmc.music_id =s.music_id full outer join content_song_writer csw on csw.music_id =s.music_id").Where("s.id=?", season.Id).Find(&musicDetails).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		if musicDetails.Singers != "" {
			singers, _ := JsonStringToStringSliceOrMap(musicDetails.Singers)
			singerResult = RemoveDuplicateStringValues(singers)
			if singerResult[0] == "" {
				singerResult = make([]string, 0)
			}
		}
		music.Singers = singerResult
		if musicDetails.MusicComposers != "" {
			musicComposers, _ := JsonStringToStringSliceOrMap(musicDetails.MusicComposers)
			musicComposerResult = RemoveDuplicateStringValues(musicComposers)
			if musicComposerResult[0] == "" {
				musicComposerResult = make([]string, 0)
			}
		}
		music.MusicComposers = musicComposerResult
		if musicDetails.SongWriters != "" {
			songWriters, _ := JsonStringToStringSliceOrMap(musicDetails.SongWriters)
			songWriterResult = RemoveDuplicateStringValues(songWriters)
			if songWriterResult[0] == "" {
				songWriterResult = make([]string, 0)
			}
		}
		music.SongWriters = songWriterResult
		contentSeason.Music = music
		if err := db.Debug().Table("content_tag ct").Select("json_agg(ct.textual_data_tag_id order by ct.textual_data_tag_id desc)::varchar as tags").Joins("join season s on s.tag_info_id =ct.tag_info_id").Where("s.id=?", season.Id).Find(&tagsDetails).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		tagsResult := make([]string, 0)
		if tagsDetails.Tags != "" {
			tags, _ := JsonStringToStringSliceOrMap(tagsDetails.Tags)
			tagsResult = RemoveDuplicateStringValues(tags)
			if tagsResult[0] == "" {
				tagsResult = make([]string, 0)
			}
		}
		tagInfo.Tags = tagsResult
		contentSeason.TagInfo = tagInfo
		genres, contentGenres = nil, nil
		if err := db.Debug().Table("season_genre sg").Select("sg.genre_id,json_agg(ss.subgenre_id order by ss.order)::varchar as subgenres_id,json_agg(ss.order order by ss.order)::varchar as sub_genre_order,sg.id").Joins("join season_subgenre ss on ss.season_genre_id = sg.id join season c on c.id = sg.season_id").Where("sg.season_id = ?", season.Id).Group("sg.genre_id,sg.id").Order("sg.order").Find(&genres).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		for _, genre := range genres {
			contentGenre.GenreId = genre.GenreId
			subGenres, err := JsonStringToStringSliceOrMap(genre.SubgenresId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			contentGenre.SubgenresId = subGenres
			contentGenre.Id = genre.Id
			contentGenres = append(contentGenres, contentGenre)
		}
		contentSeason.SeasonGenres = contentGenres
		aboutTheContent.OriginalLanguage = season.OriginalLanguage
		aboutTheContent.Supplier = season.Supplier
		aboutTheContent.AcquisitionDepartment = season.AcquisitionDepartment
		aboutTheContent.EnglishSynopsis = season.EnglishSynopsis
		aboutTheContent.ArabicSynopsis = season.ArabicSynopsis
		aboutTheContent.ProductionYear = season.ProductionYear
		aboutTheContent.ProductionHouse = season.ProductionHouse
		aboutTheContent.AgeGroup = season.AgeGroup
		aboutTheContent.IntroDuration = *&season.IntroDuration
		aboutTheContent.IntroStart = *&season.AtciIntroStart
		aboutTheContent.OutroDuration = season.OutroDuration
		aboutTheContent.OutroStart = season.AtciOutroStart

		productResult := make([]int, 0)
		if season.ProductionCountries != "" {
			products, _ := JsonStringToIntSliceOrMap(season.ProductionCountries)
			productResult = RemoveDuplicateValues(products)
			if productResult[0] == 0 {
				productResult = make([]int, 0)
			}
		}
		type ProductionCountry struct {
			CountryId int `json:"subscription_plan_id"`
		}
		var country []ProductionCountry
		if err := db.Debug().Table("production_country").Where("about_the_content_info_id=?", season.AboutTheContentInfoId).Order("country_id asc").Find(&country).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		newarr := make([]int, 0)
		for _, plans := range country {
			newarr = append(newarr, plans.CountryId)
		}
		aboutTheContent.ProductionCountries = newarr
		contentSeason.AboutTheContent = aboutTheContent
		contentSeason.Products = productResult
		contentTranslation.LanguageType = common.ContentLanguageOriginTypesName(season.LanguageType)
		contentTranslation.DubbingLanguage = season.DubbingLanguage
		contentTranslation.DubbingDialectId = season.DubbingDialectId
		contentTranslation.SubtitlingLanguage = season.SubtitlingLanguage
		contentSeason.Translation = contentTranslation
		if err := db.Debug().Table("episode e").Select("e.episode_key,e.season_id,e.status,cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date,e.created_at,e.number,pi2.video_content_id,atci.english_synopsis,atci.arabic_synopsis,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title ,cpi.notes,cpi.intro_start,cpi.outro_start,ct.language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,pi2.scheduling_date_time,json_agg(distinct pitp.target_platform)::varchar AS publishing_platforms,e.english_meta_title,e.arabic_meta_title,e.english_meta_description,e.arabic_meta_description,e.id,cc.main_actor_id,cc.main_actress_id,json_agg(ca.actor_id)::varchar as actors,json_agg(cw.writer_id)::varchar as writers,json_agg(cd.director_id)::varchar as directors,json_agg(cs.singer_id)::varchar as singers,json_agg(cmc.music_composer_id)::varchar as music_composers,json_agg(csw.song_writer_id)::varchar as song_writers,json_agg(ct1.textual_data_tag_id)::varchar as tags").Joins("join season s on s.id =e.season_id join content_primary_info cpi on cpi.id =e.primary_info_id join about_the_content_info atci on atci.id =s.about_the_content_info_id join content_translation ct on ct.id =s.translation_id join playback_item pi2 on pi2.id =e.playback_item_id join content_rights cr on cr.id =pi2.rights_id full outer join playback_item_target_platform pitp on pitp.playback_item_id = pi2.id full outer join content_cast cc on cc.id = e.cast_id full outer join content_actor ca on ca.cast_id =cc.id full outer join content_writer cw on cw.cast_id =cc.id full outer join content_director cd on cd.cast_id =cc.id full outer join content_singer cs on cs.music_id =e.music_id full outer join content_music_composer cmc on cmc.music_id =e.music_id full outer join content_song_writer csw on csw.music_id =e.music_id full outer join content_tag ct1 on ct1.tag_info_id =e.tag_info_id").Where("e.season_id =? and e.deleted_by_user_id is null", season.Id).Group("e.episode_key,e.season_id,e.status,cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date,e.created_at,e.number,pi2.video_content_id,atci.english_synopsis,atci.arabic_synopsis,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title,cpi.notes,cpi.intro_start,cpi.outro_start,ct.language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,pi2.scheduling_date_time,e.english_meta_title,e.arabic_meta_title,e.english_meta_description,e.arabic_meta_description,e.id,cc.main_actor_id,cc.main_actress_id").Find(&episodeDetails).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		for _, episode := range episodeDetails {
			var seasonEpisode SeasonEpisodes
			var cast Cast
			var music Music
			var tagInfo TagInfo
			var contentTranslation ContentTranslation
			// var castDetails CastQueryDetails
			// var musicDetails MusicQueryDetails
			// var tagsDetails TagInfoQueryDetails
			actorResult, writerResult, directorResult := make([]string, 0), make([]string, 0), make([]string, 0)
			singerResult, musicComposerResult, songWriterResult := make([]string, 0), make([]string, 0), make([]string, 0)
			seasonEpisode.IsPrimary = true
			seasonEpisode.UserId = "00000000-0000-0000-0000-000000000000"
			seasonEpisode.SecondarySeasonId = "00000000-0000-0000-0000-000000000000"
			seasonEpisode.VarianceIds = nil
			seasonEpisode.EpisodeIds = nil
			seasonEpisode.SecondaryEpisodeId = "00000000-0000-0000-0000-000000000000"
			seasonEpisode.ContentId = "00000000-0000-0000-0000-000000000000"
			seasonEpisode.EpisodeKey = episode.EpisodeKey
			seasonEpisode.SeasonId = episode.SeasonId
			seasonEpisode.Status = episode.Status
			now := time.Now().Unix()
			sDate := episode.DigitalRightsStartDate.Unix()
			eDate := episode.DigitalRightsEndDate.Unix()
			seasonEpisode.SubStatusName = nil
			seasonEpisode.Status = episode.Status
			seasonEpisode.SubStatus = episode.Status
			if episode.Status == 1 || episode.Status == 2 {
				seasonEpisode.StatusCanBeChanged = true
			} else if episode.Status == 3 {
				seasonEpisode.StatusCanBeChanged = false
				seasonEpisode.Status = 2
				seasonEpisode.SubStatus = episode.Status
				subStatus := "Draft"
				seasonEpisode.SubStatusName = &subStatus
			}
			if sDate > now || eDate < now {
				seasonEpisode.StatusCanBeChanged = false
				seasonEpisode.Status = 2
				seasonEpisode.SubStatus = episode.Status
				subStatus := "Digital Rights Exceeded"
				seasonEpisode.SubStatusName = &subStatus
			}
			seasonEpisode.DigitalRightsType = nil      // episode.DigitalRightsType
			seasonEpisode.DigitalRightsStartDate = nil // episode.DigitalRightsStartDate
			seasonEpisode.DigitalRightsEndDate = nil   //episode.DigitalRightsEndDate
			seasonEpisode.CreatedBy = nil
			var primaryInfo PrimaryInfo
			primaryInfo.Number = episode.Number
			primaryInfo.VideoContentId = episode.VideoContentId
			primaryInfo.SynopsisEnglish = episode.EnglishSynopsis
			primaryInfo.SynopsisArabic = episode.ArabicSynopsis
			primaryInfo.OriginalTitle = episode.OriginalTitle
			primaryInfo.AlternativeTitle = episode.AlternativeTitle
			primaryInfo.ArabicTitle = episode.ArabicTitle
			primaryInfo.TransliteratedTitle = episode.TransliteratedTitle
			primaryInfo.Notes = episode.Notes
			primaryInfo.IntroStart = episode.IntroStart
			primaryInfo.OutroStart = episode.OutroStart
			seasonEpisode.PrimaryInfo = primaryInfo
			// if err:=db.Table("content_cast cc").Select("cc.main_actor_id,cc.main_actress_id,json_agg(ca.actor_id)::varchar as actors,json_agg(cw.writer_id)::varchar as writers,json_agg(cd.director_id)::varchar as directors").Joins("join episode e on e.cast_id=cc.id full outer join content_actor ca on ca.cast_id =cc.id full outer join content_writer cw on cw.cast_id =cc.id full outer join content_director cd on cd.cast_id =cc.id").Where("e.season_id=?",season.Id).Group("cc.main_actor_id,cc.main_actress_id").Find(&castDetails).Error;err!=nil{
			// 	c.JSON(http.StatusInternalServerError,serverError)
			// 	return
			// }
			cast.MainActorId = episode.MainActorId
			cast.MainActressId = episode.MainActressId
			if episode.Actors != "" {
				actors, _ := JsonStringToStringSliceOrMap(episode.Actors)
				actorResult = RemoveDuplicateStringValues(actors)
				if actorResult[0] == "" {
					actorResult = make([]string, 0)
				}
			}
			cast.Actors = actorResult
			if episode.Writers != "" {
				writers, _ := JsonStringToStringSliceOrMap(episode.Writers)
				writerResult = RemoveDuplicateStringValues(writers)
				if writerResult[0] == "" {
					writerResult = make([]string, 0)
				}
			}
			cast.Writers = writerResult
			if episode.Directors != "" {
				directors, _ := JsonStringToStringSliceOrMap(episode.Directors)
				directorResult = RemoveDuplicateStringValues(directors)
				if directorResult[0] == "" {
					directorResult = make([]string, 0)
				}
			}
			cast.Directors = directorResult
			seasonEpisode.Cast = cast
			// if err:=db.Table("episode e").Select("json_agg(cs.singer_id)::varchar as singers,json_agg(cmc.music_composer_id)::varchar as music_composers,json_agg(csw.song_writer_id)::varchar as song_writers").Joins("full outer join content_singer cs on cs.music_id =e.music_id full outer join content_music_composer cmc on cmc.music_id =e.music_id full outer join content_song_writer csw on csw.music_id =e.music_id").Where("e.season_id=?",season.Id).Find(&musicDetails).Error;err!=nil{
			// 	c.JSON(http.StatusInternalServerError,serverError)
			// 	return
			// }
			if episode.Singers != "" {
				singers, _ := JsonStringToStringSliceOrMap(episode.Singers)
				singerResult = RemoveDuplicateStringValues(singers)
				if singerResult[0] == "" {
					singerResult = make([]string, 0)
				}
			}
			music.Singers = singerResult
			if episode.MusicComposers != "" {
				musicComposers, _ := JsonStringToStringSliceOrMap(episode.MusicComposers)
				musicComposerResult = RemoveDuplicateStringValues(musicComposers)
				if musicComposerResult[0] == "" {
					musicComposerResult = make([]string, 0)
				}
			}
			music.MusicComposers = musicComposerResult
			if episode.SongWriters != "" {
				songWriters, _ := JsonStringToStringSliceOrMap(episode.SongWriters)
				songWriterResult = RemoveDuplicateStringValues(songWriters)
				if songWriterResult[0] == "" {
					songWriterResult = make([]string, 0)
				}
			}
			music.SongWriters = songWriterResult
			seasonEpisode.Music = music
			// if err:=db.Table("content_tag ct").Select("json_agg(ct.textual_data_tag_id)::varchar as tags").Joins("join episode e on e.tag_info_id =ct.tag_info_id").Where("e.season_id=?",season.Id).Find(&tagsDetails).Error;err!=nil{
			// 	c.JSON(http.StatusInternalServerError,serverError)
			// 	return
			// }
			tagsResult := make([]string, 0)
			if episode.Tags != "" {
				tags, _ := JsonStringToStringSliceOrMap(episode.Tags)
				tagsResult = RemoveDuplicateStringValues(tags)
				if tagsResult[0] == "" {
					tagsResult = make([]string, 0)
				}
			}
			tagInfo.Tags = tagsResult
			seasonEpisode.TagInfo = tagInfo

			/*if episode.HasPosterImage {
				seasonEpisode.NonTextualData.PosterImage = os.Getenv("IMAGERY_URL") + episode.ContentId + "/" + episode.SeasonId + "/" + episode.Id + "/poster-image"
			}
			if episode.HasDubbingScript {
				seasonEpisode.NonTextualData.DubbingScript = os.Getenv("IMAGERY_URL") + episode.ContentId + "/" + episode.SeasonId + "/" + episode.Id + "/dubbing-script"
			}
			if episode.HasSubtitlingScript {
				seasonEpisode.NonTextualData.SubtitlingScript = os.Getenv("IMAGERY_URL") + episode.ContentId + "/" + episode.SeasonId + "/" + episode.Id + "/subtitling-script"
			} */
			seasonEpisode.NonTextualData = nil
			contentTranslation.LanguageType = common.ContentLanguageOriginTypesName(season.LanguageType)
			contentTranslation.DubbingLanguage = season.DubbingLanguage
			contentTranslation.DubbingDialectId = season.DubbingDialectId
			contentTranslation.SubtitlingLanguage = season.SubtitlingLanguage
			seasonEpisode.Translation = contentTranslation
			seasonEpisode.SchedulingDateTime = episode.SchedulingDateTime
			/* Episode Platforms */
			var platformResult []int
			if episode.PublishingPlatforms == "[null]" {
				buffer := make([]int, 0)
				platformResult = buffer
			} else {
				platforms, _ := JsonStringToIntSliceOrMap(episode.PublishingPlatforms)
				platformResult = platforms
			}
			// if episode.PublishingPlatforms != "" {
			// 	platforms, _ := JsonStringToIntSliceOrMap(episode.PublishingPlatforms)
			// 	platformResult = RemoveDuplicateValues(platforms)
			// 	if platformResult[0] == 0 {
			// 		platformResult = make([]int, 0)
			// 	}
			// }
			seasonEpisode.PublishingPlatforms = platformResult
			seasonEpisode.SeoDetails = nil
			// seasonEpisode.SeoDetails.ArabicMetaTitle = episode.ArabicMetaTitle
			// seasonEpisode.SeoDetails.EnglishMetaTitle = episode.EnglishMetaTitle
			// seasonEpisode.SeoDetails.ArabicMetaDescription = episode.ArabicMetaDescription
			// seasonEpisode.SeoDetails.EnglishMetaDescription = episode.EnglishMetaDescription
			seasonEpisode.Id = episode.Id
			seasonEpisodes = append(seasonEpisodes, seasonEpisode)
		}
		contentSeason.Episodes = seasonEpisodes
		contentSeason.NonTextualData = nil
		rightsDetails.DigitalRightsType = season.DigitalRightsType
		rightsDetails.DigitalRightsStartDate = season.DigitalRightsStartDate
		rightsDetails.DigitalRightsEndDate = season.DigitalRightsEndDate
		regionsResult := make([]int, 0)
		if season.DigitalRightsRegions != "" {
			regions, _ := JsonStringToIntSliceOrMap(season.DigitalRightsRegions)
			regionsResult = RemoveDuplicateValues(regions)
			if regionsResult[0] == 0 {
				regionsResult = make([]int, 0)
			}
		}
		rightsDetails.DigitalRightsRegions = regionsResult
		// if season.SubscriptionPlans != "" {
		// 	plans, _ := JsonStringToIntSliceOrMap(season.SubscriptionPlans)
		// 	plansResult = RemoveDuplicateValues(plans)
		// 	if plansResult[0] == 0 {
		// 		plansResult = make([]int, 0)
		// 	}
		// }
		type ContentRightsPlan struct {
			SubscriptionPlanId int `json:"subscription_plan_id"`
		}
		var plan []ContentRightsPlan
		if err := db.Debug().Table("content_rights_plan").Select("subscription_plan_id").Where("rights_id=?", season.RightsId).Find(&plan).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var newplan []int
		for _, plans := range plan {
			newplan = append(newplan, plans.SubscriptionPlanId)
		}
		rightsDetails.SubscriptionPlans = newplan
		contentSeason.Rights = rightsDetails
		contentSeason.CreatedBy = nil //season.CreatedByUserId
		contentSeason.IntroDuration = "00:00:00"
		contentSeason.IntroStart = "00:00:00"
		contentSeason.OutroDuration = "00:00:00"
		contentSeason.OutroStart = "00:00:00"
		products := make([]int, 0)
		contentSeason.Products = products
		contentSeason.SeoDetails = nil
		contentSeason.VarianceTrailers = nil
		contentSeason.Id = season.Id
		contentSeasons = append(contentSeasons, contentSeason)
	}

	multitierContentDetails.ContentSeasons = contentSeasons
	seoDetails.EnglishMetaTitle = multitierDetails.EnglishMetaTitle
	seoDetails.ArabicMetaTitle = multitierDetails.ArabicMetaTitle
	seoDetails.EnglishMetaDescription = multitierDetails.EnglishMetaDescription
	seoDetails.ArabicMetaDescription = multitierDetails.ArabicMetaDescription
	multitierContentDetails.SeoDetails = &seoDetails
	multitierContentDetails.Id = multitierDetails.Id
	c.JSON(http.StatusOK, gin.H{"data": multitierContentDetails})
	return
}
func JsonStringToStringSliceOrMap(data string) ([]string, error) {
	output := make([]string, 1000)
	err := json.Unmarshal([]byte(data), &output)
	if err != nil {
		return nil, err
	}
	//sort.Strings(output)
	return output, nil
}
func JsonStringToIntSliceOrMap(data string) ([]int, error) {
	output := make([]int, 1000)
	err := json.Unmarshal([]byte(data), &output)
	if err != nil {
		return nil, err
	}
	sort.Ints(output)
	return output, nil
}
func RemoveDuplicateValues(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
func RemoveDuplicateStringValues(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// CreateSeason -  Create Season details
// POST /seasons/published/:id
// @Summary Create season details
// @Description Create season details by content id
// @Tags Content
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param body body content.CreateSeasonRequest true "Raw JSON string"
// @Param id path string false "Id"
// @Success 200 {array} object c.JSON
// @Router /api/seasons/published/{id} [post]
func (hs *ContentService) CreateSeason(c *gin.Context) {
	AuthorizationRequired := c.MustGet("AuthorizationRequired")
	userid := c.MustGet("userid")
	if AuthorizationRequired == 1 || userid == "" || c.MustGet("is_back_office_user") == "false" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	serverError := common.ServerErrorResponse()
	notFound := common.NotFoundErrorResponse()
	var request CreateSeasonRequestValidation
	var season Season
	var errorFlag bool
	errorFlag = false
	c.ShouldBindJSON(&request)
	db := c.MustGet("DB").(*gorm.DB)
	ctx := context.Background()
	tx := db.Debug().BeginTx(ctx, nil)
	var contentiderror common.Contentiderror
	if request.ContentId == nil {
		contentiderror = common.Contentiderror{Code: "error_multitiercontent_not_found", Description: "The specified condition was not met for 'Content Id'."}
	}
	var primaryInfoError common.PrimaryInfoError
	if request.PrimaryInfo == nil {
		errorFlag = true
		primaryInfoError = common.PrimaryInfoError{Code: "NotEmptyValidator", Description: "'Primary Info' should not be empty."}

	}

	var rightserror common.RigthsError
	if request.Rights == nil {
		errorFlag = true
		rightserror = common.RigthsError{Code: "NotEmptyValidator", Description: "'Rights' should not be empty."}
	}

	var productserror common.ProductsError
	if request.Products == nil {
		errorFlag = true
		productserror = common.ProductsError{Code: "error_multitiercontent_season_products_is_null", Description: "'Products' must not be empty."}
	}
	var translationerror common.TranslationError
	if request.Translation == nil {
		errorFlag = true
		translationerror = common.TranslationError{Code: "NotEmptyValidator", Description: "'Translation' should not be empty."}
	}
	var casterror common.CastError
	if request.Cast == nil {
		errorFlag = true
		casterror = common.CastError{Code: "NotEmptyValidator", Description: "'Cast' should not be empty."}
	}
	var musicError common.MusicError
	if request.Music == nil {
		errorFlag = true
		musicError = common.MusicError{Code: "NotEmptyValidator", Description: "'Music' should not be empty."}
	}
	var taginfoError common.TaginfoError
	if request.TagInfo == nil {
		errorFlag = true
		taginfoError = common.TaginfoError{Code: "NotEmptyValidator", Description: "'Tag Info' should not be empty."}
	}

	var nontextualerrror common.NonTextualDataError
	if request.NonTextualData == nil {
		errorFlag = true
		nontextualerrror = common.NonTextualDataError{Code: "NotEmptyValidator", Description: "'Non Textual Data' must not be empty."}
	}
	var genrerror common.GenresError
	type ContentGenre struct {
		Id string `json:"id"`
	}

	var contentgenre []ContentGenre
	tx.Debug().Table("content_genre").Select("id").Where("content_id=?", request.ContentId).Find(&contentgenre)
	var newgenre []string
	for _, new := range contentgenre {
		newgenre = append(newgenre, new.Id)
	}
	if len(newgenre) < 2 || len(newgenre) > 6 {
		errorFlag = true
		genrerror = common.GenresError{Code: "error_multitiercontentGener_not_found", Description: "From (2) to (6) genre(s) are required in 'Content Genres' field in Series"}
	}
	// var seasonserror common.SeasonGenresError
	// if len(req.ContentGenres) == 0 {
	// 	errorFlag = true
	// 	seasonserror = common.SeasonGenresError{"NotEmptyValidator", "'SeasonGenres' should not be empty."}
	// }
	var invalid common.Invalidsepisode
	if contentiderror.Code != "" {
		invalid.Contentiderror = contentiderror
	}
	if primaryInfoError.Code != "" {
		fmt.Println(primaryInfoError.Code)
		invalid.PrimaryInfoError = primaryInfoError
	}

	if rightserror.Code != "" {
		invalid.RightsError = rightserror
	}

	if casterror.Code != "" {
		invalid.CastError = casterror
	}
	if musicError.Code != "" {
		invalid.MusicError = musicError
	}
	if taginfoError.Code != "" {
		invalid.TaginfoError = taginfoError
	}

	if nontextualerrror.Code != "" {
		invalid.NonTextualDataError = nontextualerrror
	}
	// if seasonserror.Code != "" {
	// 	invalid.SeasonGenresError = seasonserror
	// }
	if productserror.Code != "" {
		invalid.ProductsError = productserror
	}
	if translationerror.Code != "" {
		invalid.TranslationError = translationerror
	}
	if genrerror.Code != "" {
		invalid.GenresError = genrerror
	}
	var finalErrorResponse common.FinalErrorResponseepisode
	finalErrorResponse = common.FinalErrorResponseepisode{Error: "invalid_request", Description: "Validation failed.", Code: "error_validation_failed", RequestId: randstr.String(32), Invalid: invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
	//	ContentType := "series"
	var seasondetails common.EpisodeDetails
	if c.Param("id") == "" && request.SeasonId == "" {

		var contentCast ContentCast
		contentCast.MainActorId = request.Cast.MainActorId
		contentCast.MainActressId = request.Cast.MainActressId
		if err := tx.Create(&contentCast).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		castId := contentCast.Id
		/*create music-id for insert singer,music-composer,songwtriter info*/
		type ContentMusic struct {
			Id string `json:"id"`
		}
		var contentMusic ContentMusic
		if err := tx.Debug().Create(&contentMusic).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		MusicId := contentMusic.Id
		/*create content-tags-info for episode*/
		type ContentTagInfo struct {
			Id string `json:"id"`
		}
		var contentTagInfo ContentTagInfo
		if err := tx.Debug().Create(&contentTagInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		TaginfoId := contentTagInfo.Id
		var primaryInfo ContentPrimaryInfo
		primaryInfo.OriginalTitle = request.PrimaryInfo.OriginalTitle
		primaryInfo.AlternativeTitle = request.PrimaryInfo.AlternativeTitle
		primaryInfo.ArabicTitle = request.PrimaryInfo.ArabicTitle
		primaryInfo.TransliteratedTitle = request.PrimaryInfo.TransliteratedTitle
		primaryInfo.Notes = request.PrimaryInfo.Notes
		primaryInfo.IntroStart = request.IntroStart
		primaryInfo.OutroStart = ""
		if err := db.Debug().Create(&primaryInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var aboutTheContent AboutTheContentInfo
		aboutTheContent.OriginalLanguage = request.AboutTheContent.OriginalLanguage
		aboutTheContent.Supplier = request.AboutTheContent.Supplier
		aboutTheContent.AcquisitionDepartment = request.AboutTheContent.AcquisitionDepartment
		aboutTheContent.EnglishSynopsis = request.AboutTheContent.EnglishSynopsis
		aboutTheContent.ArabicSynopsis = request.AboutTheContent.ArabicSynopsis
		aboutTheContent.ProductionYear = request.AboutTheContent.ProductionYear
		aboutTheContent.ProductionHouse = request.AboutTheContent.ProductionHouse
		aboutTheContent.AgeGroup = request.AboutTheContent.AgeGroup
		aboutTheContent.IntroDuration = request.AboutTheContent.IntroDuration
		aboutTheContent.OutroDuration = request.AboutTheContent.OutroDuration
		aboutTheContent.IntroStart = request.AboutTheContent.IntroStart
		aboutTheContent.OutroStart = request.AboutTheContent.OutroStart
		if err := db.Debug().Create(&aboutTheContent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var productionCountry ProductionCountry
		if request.AboutTheContent.ProductionCountries != nil && len(request.AboutTheContent.ProductionCountries) > 0 {
			for _, country := range request.AboutTheContent.ProductionCountries {
				productionCountry.AboutTheContentInfoId = aboutTheContent.Id
				productionCountry.CountryId = country
				if err := tx.Debug().Create(&productionCountry).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		var contentTranslation ContentTranslationRequest
		contentTranslation.LanguageType = common.ContentLanguageOriginTypes(request.Translation.LanguageType)
		contentTranslation.DubbingLanguage = request.Translation.DubbingLanguage
		contentTranslation.DubbingDialectId = request.Translation.DubbingDialectId
		contentTranslation.SubtitlingLanguage = request.Translation.SubtitlingLanguage
		if err := tx.Debug().Table("content_translation").Create(&contentTranslation).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var contentRights ContentRights
		contentRights.DigitalRightsType = request.Rights.DigitalRightsType
		//	DRSDate, _ := time.Parse(DateTimeFormat, request.Rights.DigitalRightsStartDate)
		contentRights.DigitalRightsStartDate = request.Rights.DigitalRightsStartDate
		// DREDate, _ := time.Parse(DateTimeFormat, request.Rights.DigitalRightsEndDate)
		contentRights.DigitalRightsEndDate = request.Rights.DigitalRightsEndDate
		if err := tx.Debug().Create(&contentRights).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		RightsId := contentRights.Id
		var seasonKey Season
		tx.Debug().Select("season_key,number").Order("season_key desc").Limit(1).Find(&seasonKey)

		// season.SeasonKey = seasonKey.SeasonKey + 1
		season.SeasonKey = request.SeasonKey
		fmt.Println(request.ContentId, "request content id is")
		// for creating old seaesons take season id and user id from request body
		season.Id = request.SecondarySeasonId
		//	season.CreatedByUserId = request.CreatedByUserId

		season.ContentId = *request.ContentId
		season.Status = 1
		season.PrimaryInfoId = primaryInfo.ID
		season.AboutTheContentInfoId = aboutTheContent.Id
		season.Number = request.PrimaryInfo.SeasonNumber
		season.TranslationId = contentTranslation.Id
		season.CastId = castId
		season.CreatedByUserId = userid.(string)
		season.MusicId = MusicId
		season.TagInfoId = TaginfoId
		season.RightsId = RightsId
		season.CreatedAt = time.Now()
		season.ModifiedAt = time.Now()
		season.EnglishMetaTitle = request.SeoDetails.EnglishMetaTitle
		season.ArabicMetaTitle = request.SeoDetails.ArabicMetaTitle
		season.EnglishMetaDescription = request.SeoDetails.EnglishMetaDescription
		season.ArabicMetaDescription = request.SeoDetails.ArabicMetaDescription
		if request.NonTextualData.PosterImage != "" && request.NonTextualData.OverlayPosterImage != "" && request.NonTextualData.DetailsBackground != "" && request.NonTextualData.MobileDetailsBackground != "" {
			season.HasPosterImage = "true"
			season.HasOverlayPosterImage = "true"
			season.HasDetailsBackground = "true"
			season.HasMobileDetailsBackground = "true"
		} else {
			season.HasPosterImage = "false"
			season.HasOverlayPosterImage = "false"
			season.HasDetailsBackground = "false"
			season.HasMobileDetailsBackground = "false"
		}
		if len(request.Rights.DigitalRightsRegionsint) == 241 {
			season.HasAllRights = true
		}
		if err := tx.Debug().Create(&season).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		if request.SeasonGenres != nil {
			for i, genre := range request.SeasonGenres {
				var seasonGenre SeasonGenre
				seasonGenre.GenreId = genre.GenreId
				seasonGenre.SeasonId = season.Id
				seasonGenre.Order = i + 1
				if err := tx.Debug().Create(&seasonGenre).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}

				for j, subgenre := range genre.SubgenresId {
					var seasonSubgenre SeasonSubgenre
					seasonSubgenre.SeasonGenreId = seasonGenre.Id
					seasonSubgenre.SubgenreId = subgenre
					seasonSubgenre.Order = j + 1
					if err := tx.Debug().Create(&seasonSubgenre).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}
		}

		if len(*request.VarianceTrailers) != 0 {
			for i, trailerrange := range *request.VarianceTrailers {
				if trailerrange.VideoTrailerId != "" {
					_, _, duration := common.GetVideoDuration(trailerrange.VideoTrailerId)
					if duration == 0 {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "InValid Content TrailerId", "description": "Please provide valid Video TrailerId", "code": "", "requestId": randstr.String(32)})
						return
					}
					var varianceTrailer VarianceTrailer
					if trailerrange.TrailerPosterImage != "" {
						varianceTrailer.HasTrailerPosterImage = true
					} else {
						varianceTrailer.HasTrailerPosterImage = false
					}
					// for synchronization add traielr Id
					if len(request.VarianceTrailerIds) > 0 {
						varianceTrailer.Id = request.VarianceTrailerIds[i]
					}
					varianceTrailer.EnglishTitle = trailerrange.EnglishTitle
					varianceTrailer.ArabicTitle = trailerrange.ArabicTitle
					varianceTrailer.VideoTrailerId = trailerrange.VideoTrailerId
					varianceTrailer.SeasonId = season.Id
					varianceTrailer.Order = i + 1
					varianceTrailer.Duration = duration
					if err := tx.Debug().Table("variance_trailer").Create(&varianceTrailer).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
					go SeasonTrailerImageUploadGcp(*request.ContentId, season.Id, varianceTrailer.Id, trailerrange.TrailerPosterImage)
				}
			}
		}

		if request.Cast.Actors != nil {
			for _, actor := range request.Cast.Actors {
				var contentActor ContentActor
				contentActor.CastId = castId
				contentActor.ActorId = actor
				if err := tx.Debug().Create(&contentActor).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Cast.Writers != nil {
			for _, writer := range request.Cast.Writers {
				var contentWriter ContentWriter
				contentWriter.CastId = castId
				contentWriter.WriterId = writer
				if err := tx.Debug().Create(&contentWriter).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Cast.Actors != nil {
			for _, director := range request.Cast.Directors {
				var contentDirector ContentDirector
				contentDirector.CastId = castId
				contentDirector.DirectorId = director
				if err := tx.Debug().Create(&contentDirector).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Music.Singers != nil {
			for _, singer := range request.Music.Singers {
				var contentSingers ContentSinger
				contentSingers.MusicId = MusicId
				contentSingers.SingerId = singer
				if err := tx.Debug().Create(&contentSingers).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Music.MusicComposers != nil {
			for _, musicComposer := range request.Music.MusicComposers {
				var contentMusicComposer ContentMusicComposer
				contentMusicComposer.MusicId = MusicId
				contentMusicComposer.MusicComposerId = musicComposer
				if err := tx.Debug().Create(&contentMusicComposer).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Music.SongWriters != nil {
			for _, songWriter := range request.Music.SongWriters {
				var contentSongWriter ContentSongWriter
				contentSongWriter.MusicId = MusicId
				contentSongWriter.SongWriterId = songWriter
				if err := tx.Debug().Create(&contentSongWriter).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.TagInfo.Tags != nil {
			for _, tag := range request.TagInfo.Tags {
				var contentTag ContentTag
				contentTag.TagInfoId = TaginfoId
				contentTag.TextualDataTagId = tag
				if err := tx.Debug().Create(&contentTag).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		var contentRightsCountry ContentRightsCountry
		var contentRightsCountrys []interface{}
		for _, country := range request.Rights.DigitalRightsRegionsint {
			fmt.Println(country)
			contentRightsCountry.ContentRightsId = RightsId
			contentRightsCountry.CountryId = country
			contentRightsCountrys = append(contentRightsCountrys, contentRightsCountry)
		}
		if err := gormbulk.BulkInsert(tx, contentRightsCountrys, common.BulkInsertLimit); err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		if request.Rights.SubscriptionPlans != nil && len(request.Rights.SubscriptionPlans) > 0 {
			for _, plan := range request.Rights.SubscriptionPlans {
				var contentRightsPlan ContentRightsPlan
				contentRightsPlan.RightsId = RightsId
				contentRightsPlan.SubscriptionPlanId = plan
				if err := tx.Debug().Create(&contentRightsPlan).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		var rightsProduct RightsProduct
		if request.Products != nil && len(*request.Products) > 0 {
			for _, product := range *request.Products {
				rightsProduct.RightsId = RightsId
				rightsProduct.ProductName = product
				if err := tx.Debug().Create(&rightsProduct).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		if err := tx.Debug().Table("content").Where("id=?", season.ContentId).Update("modified_at", time.Now()).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		/*commit changes*/
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		res := map[string]string{
			"id": season.Id,
		}
		seasonId := season.Id
		/* Image Upload */
		go SeasonFileUPloadGcp(request, seasonId, season.ContentId)
		/* Fragment Creation */
		go fragments.CreateContentFragment(season.ContentId, c)
		/* update dirty count in content_sync table */
		go common.ContentSynching(season.ContentId, c)
		/* update dirty count in page_sync with contentId relation*/
		//go common.PageSyncWithContentId(season.ContentId, c)
		/* Prepare Redis Cache for all contents */
		db.Debug().Raw("select content_key,content_type from content where id=?", season.ContentId).Find(&seasondetails)
		/* Prepare Redis Cache for single content*/
		contentkeyconverted := strconv.Itoa(seasondetails.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(seasondetails.ContentType, c)
		c.JSON(http.StatusOK, gin.H{"data": res})
		return
	} else if request.SeasonId != "" && c.Param("id") == "" {
		var countryid []ContentRightsCountry
		tx.Debug().Table("season s").Select("crc.country_id").Joins("join content_rights_country crc on crc.content_rights_id=s.rights_id").Where("s.content_id=? and s.number=?  and s.deleted_by_user_id is null", request.ContentId, request.PrimaryInfo.SeasonNumber).Find(&countryid)
		var countryflag bool
		countryflag = false
		for _, data := range countryid {
			for _, value := range request.Rights.DigitalRightsRegionsint {
				if data.CountryId == value {
					countryflag = true
					break
				}
			}
		}
		if countryflag {
			c.JSON(http.StatusBadRequest, common.ServerError{Error: "countries exists", Description: "Selected countries for this variant are not allowed.", Code: "", RequestId: randstr.String(32)})
			return
		}
		var contentCast ContentCast
		contentCast.MainActorId = request.Cast.MainActorId
		contentCast.MainActressId = request.Cast.MainActressId
		if err := tx.Debug().Create(&contentCast).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		castId := contentCast.Id
		/*create music-id for insert singer,music-composer,songwtriter info*/
		type ContentMusic struct {
			Id string `json:"id"`
		}
		var contentMusic ContentMusic
		if err := tx.Debug().Create(&contentMusic).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		MusicId := contentMusic.Id
		/*create content-tags-info for episode*/
		type ContentTagInfo struct {
			Id string `json:"id"`
		}
		var contentTagInfo ContentTagInfo
		if err := tx.Debug().Create(&contentTagInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		TaginfoId := contentTagInfo.Id
		var primaryInfo ContentPrimaryInfo
		primaryInfo.OriginalTitle = request.PrimaryInfo.OriginalTitle
		primaryInfo.AlternativeTitle = request.PrimaryInfo.AlternativeTitle
		primaryInfo.ArabicTitle = request.PrimaryInfo.ArabicTitle
		primaryInfo.TransliteratedTitle = request.PrimaryInfo.TransliteratedTitle
		primaryInfo.Notes = request.PrimaryInfo.Notes
		primaryInfo.IntroStart = request.IntroStart
		primaryInfo.OutroStart = ""
		if err := tx.Debug().Create(&primaryInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var aboutTheContent AboutTheContentInfo
		aboutTheContent.OriginalLanguage = request.AboutTheContent.OriginalLanguage
		aboutTheContent.Supplier = request.AboutTheContent.Supplier
		aboutTheContent.AcquisitionDepartment = request.AboutTheContent.AcquisitionDepartment
		aboutTheContent.EnglishSynopsis = request.AboutTheContent.EnglishSynopsis
		aboutTheContent.ArabicSynopsis = request.AboutTheContent.ArabicSynopsis
		aboutTheContent.ProductionYear = request.AboutTheContent.ProductionYear
		aboutTheContent.ProductionHouse = request.AboutTheContent.ProductionHouse
		aboutTheContent.AgeGroup = request.AboutTheContent.AgeGroup
		aboutTheContent.IntroDuration = request.AboutTheContent.IntroDuration
		aboutTheContent.OutroDuration = request.AboutTheContent.OutroDuration
		aboutTheContent.IntroStart = request.AboutTheContent.IntroStart
		aboutTheContent.OutroStart = request.AboutTheContent.OutroStart
		if err := tx.Debug().Create(&aboutTheContent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var productionCountry ProductionCountry
		if request.AboutTheContent.ProductionCountries != nil && len(request.AboutTheContent.ProductionCountries) > 0 {
			for _, country := range request.AboutTheContent.ProductionCountries {
				productionCountry.AboutTheContentInfoId = aboutTheContent.Id
				productionCountry.CountryId = country
				if err := tx.Debug().Create(&productionCountry).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		var contentTranslation ContentTranslationRequest
		contentTranslation.LanguageType = common.ContentLanguageOriginTypes(request.Translation.LanguageType)
		contentTranslation.DubbingLanguage = request.Translation.DubbingLanguage
		contentTranslation.DubbingDialectId = request.Translation.DubbingDialectId
		contentTranslation.SubtitlingLanguage = request.Translation.SubtitlingLanguage
		if err := tx.Debug().Table("content_translation").Create(&contentTranslation).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var contentRights ContentRights
		contentRights.DigitalRightsType = request.Rights.DigitalRightsType
		//	DRSDate, _ := time.Parse(DateTimeFormat, request.Rights.DigitalRightsStartDate)
		contentRights.DigitalRightsStartDate = request.Rights.DigitalRightsStartDate
		//	DREDate, _ := time.Parse(DateTimeFormat, request.Rights.DigitalRightsEndDate)
		contentRights.DigitalRightsEndDate = request.Rights.DigitalRightsEndDate
		if err := tx.Debug().Create(&contentRights).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		RightsId := contentRights.Id
		var seasonKey Season
		tx.Debug().Select("season_key,number").Order("season_key desc").Limit(1).Find(&seasonKey)
		season.SeasonKey = request.SeasonKey
		// season.SeasonKey = seasonKey.SeasonKey + 1
		fmt.Println(request.ContentId, "request content id is")
		// for creating old seaesons take season id and user id from request body
		season.Id = request.SecondarySeasonId
		//	season.CreatedByUserId = request.CreatedByUserId

		season.ContentId = *request.ContentId
		season.Status = 1
		season.PrimaryInfoId = primaryInfo.ID
		season.AboutTheContentInfoId = aboutTheContent.Id
		season.Number = request.PrimaryInfo.SeasonNumber
		season.TranslationId = contentTranslation.Id
		season.CastId = castId
		season.CreatedByUserId = userid.(string)
		season.MusicId = MusicId
		season.TagInfoId = TaginfoId
		season.RightsId = RightsId
		season.CreatedAt = time.Now()
		season.ModifiedAt = time.Now()
		season.EnglishMetaTitle = request.SeoDetails.EnglishMetaTitle
		season.ArabicMetaTitle = request.SeoDetails.ArabicMetaTitle
		season.EnglishMetaDescription = request.SeoDetails.EnglishMetaDescription
		season.ArabicMetaDescription = request.SeoDetails.ArabicMetaDescription
		/*if request.NonTextualData.PosterImage != "" && request.NonTextualData.OverlayPosterImage != "" && request.NonTextualData.DubbingScript != "" && request.NonTextualData.MobileDetailsBackground != "" {
			season.HasPosterImage = "true"
			season.HasOverlayPosterImage = "true"
			season.HasDetailsBackground = "true"
			season.HasMobileDetailsBackground = "true"
		} else {
			season.HasPosterImage = "false"
			season.HasOverlayPosterImage = "false"
			season.HasDetailsBackground = "false"
			season.HasMobileDetailsBackground = "false"
		}*/
		// changed present need to check in future
		if request.NonTextualData.PosterImage == "" && request.NonTextualData.OverlayPosterImage == "" && request.NonTextualData.DubbingScript == "" && request.NonTextualData.MobileDetailsBackground == "" {
			season.HasPosterImage = "false"
			season.HasOverlayPosterImage = "false"
			season.HasDetailsBackground = "false"
			season.HasMobileDetailsBackground = "false"
		} else {
			season.HasPosterImage = "true"
			season.HasOverlayPosterImage = "true"
			season.HasDetailsBackground = "true"
			season.HasMobileDetailsBackground = "true"
		}
		if len(request.Rights.DigitalRightsRegionsint) == 241 {
			season.HasAllRights = true
		}

		if err := tx.Debug().Create(&season).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		if request.SeasonGenres != nil {
			for i, genre := range request.SeasonGenres {
				var seasonGenre SeasonGenre
				seasonGenre.GenreId = genre.GenreId
				seasonGenre.SeasonId = season.Id
				seasonGenre.Order = i + 1
				if err := tx.Debug().Create(&seasonGenre).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}

				for j, subgenre := range genre.SubgenresId {
					var seasonSubgenre SeasonSubgenre
					seasonSubgenre.SeasonGenreId = seasonGenre.Id
					seasonSubgenre.SubgenreId = subgenre
					seasonSubgenre.Order = j + 1
					if err := tx.Debug().Create(&seasonSubgenre).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}
		}

		if len(*request.VarianceTrailers) != 0 {
			for i, trailerrange := range *request.VarianceTrailers {
				if trailerrange.VideoTrailerId != "" {
					_, _, duration := common.GetVideoDuration(trailerrange.VideoTrailerId)
					if duration == 0 {
						serverError = common.ServerError{Error: "InValid Content TrailerId", Description: "Please provide valid Video TrailerId", Code: "", RequestId: randstr.String(32)}
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
					var varianceTrailer VarianceTrailer
					if trailerrange.TrailerPosterImage != "" {
						varianceTrailer.HasTrailerPosterImage = true
					} else {
						varianceTrailer.HasTrailerPosterImage = false
					}
					// for sync add trailer id
					if len(request.VarianceTrailerIds) > 0 {
						varianceTrailer.Id = request.VarianceTrailerIds[i]
					}
					varianceTrailer.EnglishTitle = trailerrange.EnglishTitle
					varianceTrailer.ArabicTitle = trailerrange.ArabicTitle
					varianceTrailer.VideoTrailerId = trailerrange.VideoTrailerId
					varianceTrailer.SeasonId = season.Id
					varianceTrailer.Order = i + 1
					varianceTrailer.Duration = duration
					if err := tx.Debug().Table("variance_trailer").Create(&varianceTrailer).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
					go SeasonVarianceTrailerImageUploadGcp(*request.ContentId, request.SeasonId, trailerrange.Id, *request.ContentId, season.Id, varianceTrailer.Id, trailerrange.TrailerPosterImage)
				}
			}
		}

		if request.Cast.Actors != nil {
			for _, actor := range request.Cast.Actors {
				var contentActor ContentActor
				contentActor.CastId = castId
				contentActor.ActorId = actor
				if err := tx.Debug().Create(&contentActor).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Cast.Writers != nil {
			for _, writer := range request.Cast.Writers {
				var contentWriter ContentWriter
				contentWriter.CastId = castId
				contentWriter.WriterId = writer
				if err := tx.Debug().Create(&contentWriter).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Cast.Actors != nil {
			for _, director := range request.Cast.Directors {
				var contentDirector ContentDirector
				contentDirector.CastId = castId
				contentDirector.DirectorId = director
				if err := tx.Debug().Create(&contentDirector).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Music.Singers != nil {
			for _, singer := range request.Music.Singers {
				var contentSingers ContentSinger
				contentSingers.MusicId = MusicId
				contentSingers.SingerId = singer
				if err := tx.Debug().Create(&contentSingers).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Music.MusicComposers != nil {
			for _, musicComposer := range request.Music.MusicComposers {
				var contentMusicComposer ContentMusicComposer
				contentMusicComposer.MusicId = MusicId
				contentMusicComposer.MusicComposerId = musicComposer
				if err := tx.Debug().Create(&contentMusicComposer).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Music.SongWriters != nil {
			for _, songWriter := range request.Music.SongWriters {
				var contentSongWriter ContentSongWriter
				contentSongWriter.MusicId = MusicId
				contentSongWriter.SongWriterId = songWriter
				if err := tx.Debug().Create(&contentSongWriter).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.TagInfo.Tags != nil {
			for _, tag := range request.TagInfo.Tags {
				var contentTag ContentTag
				contentTag.TagInfoId = TaginfoId
				contentTag.TextualDataTagId = tag
				if err := tx.Debug().Create(&contentTag).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		var contentRightsCountry ContentRightsCountry
		var contentRightsCountrys []interface{}
		for _, country := range request.Rights.DigitalRightsRegionsint {
			fmt.Println(country)
			contentRightsCountry.ContentRightsId = RightsId
			contentRightsCountry.CountryId = country
			contentRightsCountrys = append(contentRightsCountrys, contentRightsCountry)
		}

		if err := gormbulk.BulkInsert(tx.Debug(), contentRightsCountrys, common.BULK_INSERT_LIMIT); err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		if request.Rights.SubscriptionPlans != nil && len(request.Rights.SubscriptionPlans) > 0 {
			for _, plan := range request.Rights.SubscriptionPlans {
				var contentRightsPlan ContentRightsPlan
				contentRightsPlan.RightsId = RightsId
				contentRightsPlan.SubscriptionPlanId = plan
				if err := tx.Debug().Create(&contentRightsPlan).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		var rightsProduct RightsProduct
		if request.Products != nil && len(*request.Products) > 0 {
			for _, product := range *request.Products {
				rightsProduct.RightsId = RightsId
				rightsProduct.ProductName = product
				if err := tx.Debug().Create(&rightsProduct).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		if err := tx.Debug().Table("content").Where("id=?", season.ContentId).Update("modified_at", time.Now()).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		//	 notFound := common.NotFoundErrorResponse()

		var episode []Episode
		// fetch episode
		if err := tx.Debug().Raw("select distinct on (number) number ,* from episode where season_id=? ", request.SeasonId).Find(&episode).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		// episode details
		for _, data := range episode {
			// fetching primary info details
			var primaryinodetails ContentPrimaryInfo
			fmt.Println(data.PrimaryInfoId, "primary info id is")
			if err := tx.Debug().Raw("select * from content_primary_info where id =?", data.PrimaryInfoId).Find(&primaryinodetails).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			// creating primary info
			var primarydata ContentPrimaryInfo
			primarydata = ContentPrimaryInfo{OriginalTitle: primaryinodetails.OriginalTitle, AlternativeTitle: primaryinodetails.AlternativeTitle, ArabicTitle: primaryinodetails.AlternativeTitle, TransliteratedTitle: primaryinodetails.TransliteratedTitle, Notes: primaryinodetails.Notes, IntroStart: primaryinodetails.IntroStart, OutroStart: primaryinodetails.OutroStart}
			if episodeprimaryinfo := tx.Debug().Table("content_primary_info").Create(&primarydata).Error; episodeprimaryinfo != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			// feltching playback item details
			var playbackitem PlaybackItem
			if err := tx.Debug().Raw("select * from playback_item where id =?", data.PlaybackItemId).Find(&playbackitem).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			//	fetching contentrights
			var contentrights ContentRights
			if err := tx.Debug().Raw("select * from content_rights where id =?", playbackitem.RightsId).Find(&contentrights).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			// creating content rights

			contentrightsdata := ContentRights{DigitalRightsType: contentrights.DigitalRightsType,
				DigitalRightsStartDate: contentrights.DigitalRightsEndDate, DigitalRightsEndDate: contentrights.DigitalRightsEndDate}
			//contentrightsdetails = append(contentrightsdetails, contentrightsdata)
			if err := tx.Debug().Table("content_rights").Create(&contentrightsdata).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			//	fetching content translation details
			var contenttranslation ContentTranslationData
			if err := tx.Debug().Raw("select * from content_translation where id =?", playbackitem.TranslationId).Find(&contenttranslation).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			var contenttransdata ContentTranslationData
			// creating content  translation details
			contenttransdata.LanguageType = contenttranslation.LanguageType
			contenttransdata.DubbingLanguage = contenttranslation.DubbingLanguage
			contenttransdata.DubbingDialectId = contenttranslation.DubbingDialectId
			contenttransdata.SubtitlingLanguage = contenttranslation.SubtitlingLanguage
			//	contenttarnsfinal = append(contenttarnsfinal, contentdata)
			if err := tx.Debug().Table("content_translation").Create(&contenttransdata).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			//create playback item
			playbackdata := PlaybackItem{VideoContentId: playbackitem.VideoContentId, SchedulingDateTime: playbackitem.SchedulingDateTime, TranslationId: contenttransdata.Id, RightsId: contentrightsdata.Id, Duration: playbackitem.Duration}
			if createPlaybackItem := tx.Debug().Table("playback_item").Create(&playbackdata).Error; createPlaybackItem != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			// feltching platform details
			var newplatform []interface{}
			var TargetPlatformvalues []PlaybackItemTargetPlatform
			tx.Table("playback_item_target_platform").Select("target_platform").Where("playback_item_id=?", playbackitem.Id).Find(&TargetPlatformvalues)

			for _, value := range TargetPlatformvalues {
				playbackdetails := PlaybackItemTargetPlatform{PlaybackItemId: playbackdata.Id, TargetPlatform: value.TargetPlatform, RightsId: contentrightsdata.Id}
				newplatform = append(newplatform, playbackdetails)
			}

			// inserting publishing platforms
			if err := gormbulk.BulkInsert(tx.Debug(), newplatform, common.BULK_INSERT_LIMIT); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			// fetching content cast details
			var contentcast ContentCast
			if err := tx.Debug().Raw("select * from content_cast where id=?", data.CastId).Find(&contentcast).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			// inserting content cast details
			contentcastdata := ContentCast{MainActorId: contentcast.MainActorId, MainActressId: contentcast.MainActressId}
			if err := tx.Debug().Table("content_cast").Create(&contentcastdata).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			// fetching content actor details
			var contentactor ContentActor
			tx.Debug().Raw("select * from content_actor where cast_id=?", data.CastId).Find(&contentactor)
			// inserting content actor details
			if contentactor.ActorId != "" {
				contentactordata := ContentActor{CastId: contentcastdata.Id, ActorId: contentactor.ActorId}
				if err := tx.Debug().Table("content_actor").Create(&contentactordata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			// fetching content writer details
			fmt.Println(data.CastId, "cast id is")
			var contentwriter ContentWriter
			tx.Debug().Raw("select * from content_writer where cast_id=?", data.CastId).Find(&contentwriter)

			// inserting content writer details
			if contentwriter.WriterId != "" {
				contentwriterdata := ContentWriter{CastId: contentcastdata.Id, WriterId: contentwriter.WriterId}
				if err := tx.Debug().Table("content_writer").Create(&contentwriterdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			// fetching content director details
			var contentdirector ContentDirector
			tx.Debug().Raw("select * from content_director where cast_id=?", data.CastId).Find(&contentdirector)
			// inserting content director details
			fmt.Println(contentdirector.DirectorId, "director id is")
			if contentdirector.DirectorId != "" {
				contentdirectorrdata := ContentDirector{CastId: contentcastdata.Id, DirectorId: contentdirector.DirectorId}
				if err := tx.Debug().Table("content_director").Create(&contentdirectorrdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			// inserting contetn music details
			var contentmusic ContentMusic
			if err := tx.Debug().Table("content_music").Create(&contentmusic).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			fmt.Println(contentmusic.Id, "content music")
			// fetching content music composer details
			var contentMusicComposer ContentMusicComposer
			tx.Debug().Raw("select * from content_music_composer where music_id=?", data.MusicId).Find(&contentMusicComposer)

			// inserting content music composer  details
			if contentMusicComposer.MusicComposerId != "" {
				contentmusicdata := ContentMusicComposer{MusicComposerId: contentMusicComposer.MusicComposerId, MusicId: contentmusic.Id}
				if err := tx.Debug().Table("content_music_composer").Create(&contentmusicdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			// fetching content music singer details
			var contentsinger ContentSinger
			tx.Debug().Raw("select * from content_singer where music_id=?", data.MusicId).Find(&contentsinger)

			// inserting content music singer  details
			if contentsinger.SingerId != "" {
				contentsingerdata := ContentSinger{SingerId: contentsinger.SingerId, MusicId: contentmusic.Id}
				if err := tx.Debug().Table("content_singer").Create(&contentsingerdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}

			// fetching content music song writer details
			var contentsongwriter ContentSongWriter
			tx.Debug().Raw("select * from content_song_writer where music_id=?", data.MusicId).Find(&contentsongwriter)

			// inserting content music song writer  details
			if contentsongwriter.SongWriterId != "" {
				contentsongdata := ContentSongWriter{SongWriterId: contentsongwriter.SongWriterId, MusicId: contentmusic.Id}
				if err := tx.Debug().Table("content_song_writer").Create(&contentsongdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			// inserting tag info details
			var contenttaginfo ContentTagInfo
			if err := tx.Debug().Table("content_tag_info").Create(&contenttaginfo).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			// fetching content _tag details
			var contenttag ContentTag
			tx.Debug().Raw("select * from content_tag where tag_info_id=?", data.TagInfoId).Find(&contenttag)
			var textualdatatag TextualDataTag
			if contenttag.TextualDataTagId != "" {
				// fetching textual data tag details
				tx.Debug().Raw("select * from textual_data_tag where id=?", contenttag.TextualDataTagId).Find(&textualdatatag)
			}

			// inserting textual data tag details
			textualdata := TextualDataTag{Name: textualdatatag.Name}
			if textualdatatag.Name != "" {
				if err := tx.Debug().Table("textual_data_tag").Create(&textualdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}

			// inserting content_tag details
			if textualdata.Id != "" {
				contenttagdata := ContentTag{TagInfoId: contenttaginfo.Id, TextualDataTagId: textualdata.Id}
				if err := tx.Debug().Table("content_tag").Create(&contenttagdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			var episodekey Episode
			tx.Debug().Table("episode").Select("max(episode_key) as episode_key").Find(&episodekey)
			//		var episodefinal []interface{}
			var episodeDetails Episode
			episodeDetails.Number = data.Number
			episodeDetails.SeasonId = season.Id
			episodeDetails.PrimaryInfoId = primarydata.ID
			episodeDetails.PlaybackItemId = playbackdata.Id
			episodeDetails.Status = 1
			episodeDetails.SynopsisEnglish = data.SynopsisArabic
			episodeDetails.SynopsisArabic = data.SynopsisArabic
			episodeDetails.CastId = contentcastdata.Id
			episodeDetails.MusicId = contentMusic.Id
			episodeDetails.TagInfoId = contenttaginfo.Id
			episodeDetails.EpisodeKey = episodekey.EpisodeKey + 1
			episodeDetails.CreatedAt = time.Now()
			episodeDetails.ModifiedAt = time.Now()
			episodeDetails.EnglishMetaTitle = data.EnglishMetaTitle
			episodeDetails.ArabicMetaTitle = data.ArabicMetaTitle
			episodeDetails.EnglishMetaDescription = data.EnglishMetaDescription
			episodeDetails.ArabicMetaDescription = data.ArabicMetaDescription
			episodeDetails.HasPosterImage = data.HasPosterImage
			episodeDetails.HasSubtitlingScript = data.HasSubtitlingScript
			episodeDetails.HasDubbingScript = data.HasDubbingScript
			var image []Images
			image = append(image, Images{Imagename: "poster-image", HasImage: data.HasPosterImage})
			image = append(image, Images{Imagename: "dubbing-script", HasImage: data.HasDubbingScript})
			image = append(image, Images{Imagename: "subtitling-script", HasImage: data.HasSubtitlingScript})
			if err := tx.Debug().Table("episode").Create(&episodeDetails).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				fmt.Println(serverError, "server error")
				return
			}
			go SeasonVarianceEpisodeImageUploadGcp(*request.ContentId, season.Id, episodeDetails.Id, image, *request.ContentId, request.SeasonId, data.Id)
		}
		/*commit changes*/
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		seasonId := season.Id
		SeasonVarianceFileUPloadGcp(request.SeasonId, seasonId, season.ContentId)
		/* Fragment Creation */
		go fragments.CreateContentFragment(season.ContentId, c)
		/* Prepare Redis Cache for all contents */
		db.Debug().Raw("select content_key,content_type from content where id=?", season.ContentId).Find(&seasondetails)
		/* Prepare Redis Cache for single content*/
		contentkeyconverted := strconv.Itoa(seasondetails.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(seasondetails.ContentType, c)
		c.JSON(http.StatusOK, gin.H{"message": "Season Created Successfully.", "status": http.StatusOK, "id": season.Id})
		return
	} else {
		if err := tx.Debug().Where("id=?", c.Param("id")).Find(&season).Error; err != nil {
			c.JSON(http.StatusNotFound, notFound)
			return
		}
		result, err := UpdateSeasonDetails(request, c, season, 1)
		fmt.Println(err, ".....................")
		if result == "" {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		res := map[string]string{
			"id": c.Param("id"),
		}
		seasonId := c.Param("id")
		/* Image Upload */
		go SeasonFileUPloadGcp(request, seasonId, season.ContentId)
		/* Fragment Creation */
		go fragments.CreateContentFragment(season.ContentId, c)
		/* update dirty count in content_sync table */
		go common.ContentSynching(season.ContentId, c)
		/* update dirty count in page_sync with contentId relation*/
		//go common.PageSyncWithContentId(season.ContentId, c)
		/* Prepare Redis Cache for all contents */
		db.Debug().Raw("select content_key,content_type from content where id=?", season.ContentId).Find(&seasondetails)
		fmt.Println("KEY", seasondetails.ContentType, seasondetails.ContentKey)
		/* Prepare Redis Cache for single content*/
		contentkeyconverted := strconv.Itoa(seasondetails.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(seasondetails.ContentType, c)
		c.JSON(http.StatusOK, gin.H{"data": res})
		return
	}

}

// CreateSeason -  Draft Season details
// POST /api/seasons/draft/:id
// @Summary Draft season details
// @Description Draft season details by content id
// @Tags Content
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param body body content.CreateSeasonRequest true "Raw JSON string"
// @Param id path string false "Id"
// @Success 200 {array} object c.JSON
// @Router /api/seasons/draft/{id} [post]
func (hs *ContentService) DraftSeason(c *gin.Context) {
	AuthorizationRequired := c.MustGet("AuthorizationRequired")
	userid := c.MustGet("userid")
	if AuthorizationRequired == 1 || userid == "" || c.MustGet("is_back_office_user") == "false" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	serverError := common.ServerErrorResponse()
	notFound := common.NotFoundErrorResponse()

	var request CreateSeasonRequestValidation
	var season Season
	var errorFlag bool
	errorFlag = false
	c.ShouldBindJSON(&request)
	db := c.MustGet("DB").(*gorm.DB)
	ctx := context.Background()
	tx := db.Debug().BeginTx(ctx, nil)
	var contentiderror common.Contentiderror
	if request.ContentId == nil {
		contentiderror = common.Contentiderror{Code: "error_multitiercontent_not_found", Description: "The specified condition was not met for 'Content Id'."}
	}
	var primaryInfoError common.PrimaryInfoError
	if request.PrimaryInfo == nil {
		errorFlag = true
		primaryInfoError = common.PrimaryInfoError{Code: "NotEmptyValidator", Description: "'Primary Info' should not be empty."}
	}

	var rightserror common.RigthsError
	if request.Rights == nil {
		errorFlag = true
		rightserror = common.RigthsError{Code: "NotEmptyValidator", Description: "'Rights' should not be empty."}
	}
	// var seasonserror common.SeasonGenresError
	// fmt.Println(req.ContentGenres, "lllllllllllll")
	// if len(req.ContentGenres) <1{
	// 	errorFlag = true
	// 	fmt.Println(errorFlag, ";;;;;;;;;;;;")
	// 	seasonserror = common.SeasonGenresError{"NotEmptyValidator", "'SeasonGenres' should not be empty."}
	// }
	var productserror common.ProductsError
	if request.Products == nil {
		errorFlag = true
		productserror = common.ProductsError{Code: "error_multitiercontent_season_products_is_null", Description: "'Products' must not be empty."}
	}
	var translationerror common.TranslationError
	if request.Translation == nil {
		errorFlag = true
		translationerror = common.TranslationError{Code: "NotEmptyValidator", Description: "'Translation' should not be empty."}
	}
	var casterror common.CastError
	if request.Cast == nil {
		errorFlag = true
		casterror = common.CastError{Code: "NotEmptyValidator", Description: "'Cast' should not be empty."}
	}
	var musicError common.MusicError
	if request.Music == nil {
		errorFlag = true
		musicError = common.MusicError{Code: "NotEmptyValidator", Description: "'Music' should not be empty."}
	}
	var taginfoError common.TaginfoError
	if request.TagInfo == nil {
		errorFlag = true
		taginfoError = common.TaginfoError{Code: "NotEmptyValidator", Description: "'Tag Info' should not be empty."}
	}

	var nontextualerrror common.NonTextualDataError
	if request.NonTextualData == nil {
		errorFlag = true
		nontextualerrror = common.NonTextualDataError{Code: "NotEmptyValidator", Description: "'Non Textual Data' must not be empty."}
	}

	var invalid common.Invalidsepisode
	if contentiderror.Code != "" {
		invalid.Contentiderror = contentiderror
	}
	if primaryInfoError.Code != "" {
		fmt.Println(primaryInfoError.Code)
		invalid.PrimaryInfoError = primaryInfoError
	}
	var genrerror common.GenresError
	type ContentGenre struct {
		Id string `json:"id"`
	}

	var contentgenre []ContentGenre
	tx.Debug().Table("content_genre").Select("id").Where("content_id=?", request.ContentId).Find(&contentgenre)
	var newgenre []string
	for _, new := range contentgenre {
		newgenre = append(newgenre, new.Id)
	}
	if len(newgenre) < 2 || len(newgenre) > 6 {
		errorFlag = true
		genrerror = common.GenresError{Code: "error_multitiercontentGener_not_found", Description: "From (2) to (6) genre(s) are required in 'Content Genres' field in Series"}
	}

	if rightserror.Code != "" {
		invalid.RightsError = rightserror
	}

	if casterror.Code != "" {
		invalid.CastError = casterror
	}
	if musicError.Code != "" {
		invalid.MusicError = musicError
	}
	if taginfoError.Code != "" {
		invalid.TaginfoError = taginfoError
	}

	if nontextualerrror.Code != "" {
		invalid.NonTextualDataError = nontextualerrror
	}
	// if seasonserror.Code != "" {
	// 	invalid.SeasonGenresError = seasonserror
	// }
	if productserror.Code != "" {
		invalid.ProductsError = productserror
	}
	if translationerror.Code != "" {
		invalid.TranslationError = translationerror
	}
	if genrerror.Code != "" {
		invalid.GenresError = genrerror
	}
	var finalErrorResponse common.FinalErrorResponseepisode
	finalErrorResponse = common.FinalErrorResponseepisode{Error: "invalid_request", Description: "Validation failed.", Code: "error_validation_failed", RequestId: randstr.String(32), Invalid: invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}

	c.ShouldBindJSON(&request)
	var seasondetails common.EpisodeDetails
	if c.Param("id") == "" && request.SeasonId == "" {
		var contentCast ContentCast
		contentCast.MainActorId = request.Cast.MainActorId
		contentCast.MainActressId = request.Cast.MainActressId
		if err := tx.Create(&contentCast).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		castId := contentCast.Id
		/*create music-id for insert singer,music-composer,songwtriter info*/
		type ContentMusic struct {
			Id string `json:"id"`
		}
		var contentMusic ContentMusic
		if err := tx.Debug().Create(&contentMusic).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		MusicId := contentMusic.Id
		/*create content-tags-info for episode*/
		type ContentTagInfo struct {
			Id string `json:"id"`
		}
		var contentTagInfo ContentTagInfo
		if err := tx.Debug().Create(&contentTagInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		TaginfoId := contentTagInfo.Id
		var primaryInfo ContentPrimaryInfo
		primaryInfo.OriginalTitle = request.PrimaryInfo.OriginalTitle
		primaryInfo.AlternativeTitle = request.PrimaryInfo.AlternativeTitle
		primaryInfo.ArabicTitle = request.PrimaryInfo.ArabicTitle
		primaryInfo.TransliteratedTitle = request.PrimaryInfo.TransliteratedTitle
		primaryInfo.Notes = request.PrimaryInfo.Notes
		primaryInfo.IntroStart = request.IntroStart
		primaryInfo.OutroStart = ""
		if err := tx.Debug().Create(&primaryInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var aboutTheContent AboutTheContentInfo
		aboutTheContent.OriginalLanguage = request.AboutTheContent.OriginalLanguage
		aboutTheContent.Supplier = request.AboutTheContent.Supplier
		aboutTheContent.AcquisitionDepartment = request.AboutTheContent.AcquisitionDepartment
		aboutTheContent.EnglishSynopsis = request.AboutTheContent.EnglishSynopsis
		aboutTheContent.ArabicSynopsis = request.AboutTheContent.ArabicSynopsis
		aboutTheContent.ProductionYear = request.AboutTheContent.ProductionYear
		aboutTheContent.ProductionHouse = request.AboutTheContent.ProductionHouse
		aboutTheContent.AgeGroup = request.AboutTheContent.AgeGroup
		aboutTheContent.IntroDuration = request.AboutTheContent.IntroDuration
		aboutTheContent.OutroDuration = request.AboutTheContent.OutroDuration
		aboutTheContent.IntroStart = request.AboutTheContent.IntroStart
		aboutTheContent.OutroStart = request.AboutTheContent.OutroStart
		if err := tx.Debug().Create(&aboutTheContent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			fmt.Println(serverError, "about error is")
			return
		}
		var productionCountry ProductionCountry
		if request.AboutTheContent.ProductionCountries != nil && len(request.AboutTheContent.ProductionCountries) > 0 {
			for _, country := range request.AboutTheContent.ProductionCountries {
				productionCountry.AboutTheContentInfoId = aboutTheContent.Id
				productionCountry.CountryId = country
				if err := tx.Debug().Create(&productionCountry).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					fmt.Println(serverError, "server error is")
					return
				}
			}
		}
		var contentTranslation ContentTranslationRequest
		contentTranslation.LanguageType = common.ContentLanguageOriginTypes(request.Translation.LanguageType)
		contentTranslation.DubbingLanguage = request.Translation.DubbingLanguage
		contentTranslation.DubbingDialectId = request.Translation.DubbingDialectId
		contentTranslation.SubtitlingLanguage = request.Translation.SubtitlingLanguage
		if err := tx.Debug().Table("content_translation").Create(&contentTranslation).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			fmt.Println(serverError, "trans error")
			return
		}
		var contentRights ContentRights
		contentRights.DigitalRightsType = request.Rights.DigitalRightsType
		//	DRSDate, _ := time.Parse(DateTimeFormat, request.Rights.DigitalRightsStartDate)
		contentRights.DigitalRightsStartDate = request.Rights.DigitalRightsStartDate
		//	DREDate, _ := time.Parse(DateTimeFormat, request.Rights.DigitalRightsEndDate)
		contentRights.DigitalRightsEndDate = request.Rights.DigitalRightsEndDate
		if err := db.Debug().Create(&contentRights).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		RightsId := contentRights.Id
		var seasonKey Season
		tx.Debug().Select("season_key,number").Order("season_key desc").Limit(1).Find(&seasonKey)
		// for creating old seaesons take season id and user id from request body
		season.Id = request.SecondarySeasonId
		//	season.CreatedByUserId = request.CreatedByUserId
		season.SeasonKey = request.SeasonKey
		// season.SeasonKey = seasonKey.SeasonKey + 1
		season.ContentId = *request.ContentId
		season.Status = 3
		season.PrimaryInfoId = primaryInfo.ID
		season.AboutTheContentInfoId = aboutTheContent.Id
		season.Number = request.PrimaryInfo.SeasonNumber
		season.TranslationId = contentTranslation.Id
		season.CastId = castId
		season.CreatedByUserId = userid.(string)
		season.MusicId = MusicId
		season.TagInfoId = TaginfoId
		season.RightsId = RightsId
		season.CreatedAt = time.Now()
		season.ModifiedAt = time.Now()
		season.EnglishMetaTitle = request.SeoDetails.EnglishMetaTitle
		season.ArabicMetaTitle = request.SeoDetails.ArabicMetaTitle
		season.EnglishMetaDescription = request.SeoDetails.EnglishMetaDescription
		season.ArabicMetaDescription = request.SeoDetails.ArabicMetaDescription
		if request.NonTextualData.PosterImage != "" && request.NonTextualData.OverlayPosterImage != "" && request.NonTextualData.DetailsBackground != "" && request.NonTextualData.MobileDetailsBackground != "" {
			season.HasPosterImage = "true"
			season.HasOverlayPosterImage = "true"
			season.HasDetailsBackground = "true"
			season.HasMobileDetailsBackground = "true"
		} else {
			season.HasPosterImage = "false"
			season.HasOverlayPosterImage = "false"
			season.HasDetailsBackground = "false"
			season.HasMobileDetailsBackground = "false"
		}
		if len(request.Rights.DigitalRightsRegionsint) == 241 {
			season.HasAllRights = true
		}

		if err := tx.Debug().Create(&season).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		} else {
			if request.SeasonGenres != nil {
				for i, genre := range request.SeasonGenres {
					var seasonGenre SeasonGenre
					seasonGenre.GenreId = genre.GenreId
					seasonGenre.SeasonId = season.Id
					seasonGenre.Order = i + 1
					if err := tx.Debug().Create(&seasonGenre).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						fmt.Println("seasongnere error", serverError)
						return
					}

					for j, subgenre := range genre.SubgenresId {
						var seasonSubgenre SeasonSubgenre
						seasonSubgenre.SeasonGenreId = seasonGenre.Id
						seasonSubgenre.SubgenreId = subgenre
						seasonSubgenre.Order = j + 1
						if err := tx.Debug().Create(&seasonSubgenre).Error; err != nil {
							c.JSON(http.StatusInternalServerError, serverError)
							return
						}
					}
				}
			}
			if len(*request.VarianceTrailers) != 0 {

				for i, trailerrange := range *request.VarianceTrailers {
					if trailerrange.VideoTrailerId != "" {
						_, _, duration := common.GetVideoDuration(trailerrange.VideoTrailerId)
						if duration == 0 {
							c.JSON(http.StatusInternalServerError, gin.H{"error": "InValid Content TrailerId", "description": "Please provide valid Video TrailerId", "code": "", "requestId": randstr.String(32)})
							return
						}
						var varianceTrailer VarianceTrailer
						if trailerrange.TrailerPosterImage != "" {
							varianceTrailer.HasTrailerPosterImage = true
						} else {
							varianceTrailer.HasTrailerPosterImage = false
						}
						// for sync add trailer id
						if len(request.VarianceTrailerIds) > 0 {
							varianceTrailer.Id = request.VarianceTrailerIds[i]
						}
						varianceTrailer.EnglishTitle = trailerrange.EnglishTitle
						varianceTrailer.ArabicTitle = trailerrange.ArabicTitle
						varianceTrailer.VideoTrailerId = trailerrange.VideoTrailerId
						varianceTrailer.SeasonId = season.Id
						varianceTrailer.Order = i + 1
						varianceTrailer.Duration = duration
						if err := tx.Debug().Table("variance_trailer").Create(&varianceTrailer).Error; err != nil {
							c.JSON(http.StatusInternalServerError, serverError)
							return
						}
						go SeasonTrailerImageUploadGcp(*request.ContentId, season.Id, varianceTrailer.Id, trailerrange.TrailerPosterImage)
					}
				}
			}

			// fmt.Println("1")
			if request.Cast.Actors != nil {
				for _, actor := range request.Cast.Actors {
					var contentActor ContentActor
					contentActor.CastId = castId
					contentActor.ActorId = actor
					if err := tx.Debug().Create(&contentActor).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}
			// fmt.Println("2")
			if request.Cast.Writers != nil {
				for _, writer := range request.Cast.Writers {
					var contentWriter ContentWriter
					contentWriter.CastId = castId
					contentWriter.WriterId = writer
					if err := tx.Debug().Create(&contentWriter).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}
			// fmt.Println("3")
			if request.Cast.Actors != nil {
				for _, director := range request.Cast.Directors {
					var contentDirector ContentDirector
					contentDirector.CastId = castId
					contentDirector.DirectorId = director
					if err := tx.Debug().Create(&contentDirector).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}
			// fmt.Println("4")
			if request.Music.Singers != nil {
				for _, singer := range request.Music.Singers {
					var contentSingers ContentSinger
					contentSingers.MusicId = MusicId
					contentSingers.SingerId = singer
					if err := tx.Debug().Create(&contentSingers).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}
			// fmt.Println("5")
			if request.Music.MusicComposers != nil {
				for _, musicComposer := range request.Music.MusicComposers {
					var contentMusicComposer ContentMusicComposer
					contentMusicComposer.MusicId = MusicId
					contentMusicComposer.MusicComposerId = musicComposer
					if err := tx.Debug().Create(&contentMusicComposer).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}
			// fmt.Println("6")
			if request.Music.SongWriters != nil {
				for _, songWriter := range request.Music.SongWriters {
					var contentSongWriter ContentSongWriter
					contentSongWriter.MusicId = MusicId
					contentSongWriter.SongWriterId = songWriter
					if err := tx.Debug().Create(&contentSongWriter).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}
			// fmt.Println("7")
			if request.TagInfo.Tags != nil {
				for _, tag := range request.TagInfo.Tags {
					var contentTag ContentTag
					contentTag.TagInfoId = TaginfoId
					contentTag.TextualDataTagId = tag
					if err := tx.Debug().Create(&contentTag).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						fmt.Println(serverError, "content tag error is")
						return
					}
				}
			}
			// fmt.Println("8")

			var contentRightsCountry ContentRightsCountry
			var contentRightsCountrys []interface{}

			for _, country := range request.Rights.DigitalRightsRegionsint {
				contentRightsCountry.ContentRightsId = RightsId
				contentRightsCountry.CountryId = country
				contentRightsCountrys = append(contentRightsCountrys, contentRightsCountry)
			}
			err = gormbulk.BulkInsert(tx.Debug(), contentRightsCountrys, common.BulkInsertLimit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			if request.Rights.SubscriptionPlans != nil && len(request.Rights.SubscriptionPlans) > 0 {
				for _, plan := range request.Rights.SubscriptionPlans {
					var contentRightsPlan ContentRightsPlan
					contentRightsPlan.RightsId = RightsId
					contentRightsPlan.SubscriptionPlanId = plan
					if err := tx.Debug().Create(&contentRightsPlan).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}
			var rightsProduct RightsProduct
			if request.Products != nil && len(*request.Products) > 0 {
				for _, product := range *request.Products {
					rightsProduct.RightsId = RightsId
					rightsProduct.ProductName = product
					if err := tx.Debug().Create(&rightsProduct).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}
		}
		if err := tx.Debug().Table("content").Where("id=?", season.ContentId).Update("modified_at", time.Now()).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		/*commit changes*/
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		res := map[string]string{
			"id": season.Id,
		}
		seasonId := season.Id
		SeasonFileUPloadGcp(request, seasonId, season.ContentId)
		/* Prepare Redis Cache for all contents */
		db.Debug().Raw("select content_key,content_type from content where id=?", season.ContentId).Find(&seasondetails)
		/* Prepare Redis Cache for single content*/
		contentkeyconverted := strconv.Itoa(seasondetails.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(seasondetails.ContentType, c)

		c.JSON(http.StatusOK, gin.H{"data": res})
		return

	} else if request.SeasonId != "" && c.Param("id") == "" {

		var countryid []ContentRightsCountry
		tx.Debug().Table("season s").Select("crc.country_id").Joins("join content_rights_country crc on crc.content_rights_id=s.rights_id").Where("s.content_id=? and s.number=? and s.deleted_by_user_id is null", request.ContentId, request.PrimaryInfo.SeasonNumber).Find(&countryid)
		var countryflag bool
		countryflag = false
		for _, data := range countryid {
			for _, value := range request.Rights.DigitalRightsRegionsint {
				if data.CountryId == value {
					countryflag = true
					break
				}
			}
		}
		if countryflag {
			c.JSON(http.StatusBadRequest, common.ServerError{Error: "countries exists", Description: "Selected countries for this variant are not allowed.", Code: "", RequestId: randstr.String(32)})
			return
		}
		var contentCast ContentCast
		contentCast.MainActorId = request.Cast.MainActorId
		contentCast.MainActressId = request.Cast.MainActressId
		if err := tx.Debug().Create(&contentCast).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		castId := contentCast.Id
		/*create music-id for insert singer,music-composer,songwtriter info*/
		type ContentMusic struct {
			Id string `json:"id"`
		}
		var contentMusic ContentMusic
		if err := tx.Debug().Create(&contentMusic).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		MusicId := contentMusic.Id
		/*create content-tags-info for episode*/
		type ContentTagInfo struct {
			Id string `json:"id"`
		}
		var contentTagInfo ContentTagInfo
		if err := tx.Debug().Create(&contentTagInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		TaginfoId := contentTagInfo.Id
		var primaryInfo ContentPrimaryInfo
		primaryInfo.OriginalTitle = request.PrimaryInfo.OriginalTitle
		primaryInfo.AlternativeTitle = request.PrimaryInfo.AlternativeTitle
		primaryInfo.ArabicTitle = request.PrimaryInfo.ArabicTitle
		primaryInfo.TransliteratedTitle = request.PrimaryInfo.TransliteratedTitle
		primaryInfo.Notes = request.PrimaryInfo.Notes
		primaryInfo.IntroStart = request.IntroStart
		primaryInfo.OutroStart = ""
		if err := tx.Debug().Create(&primaryInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var aboutTheContent AboutTheContentInfo
		aboutTheContent.OriginalLanguage = request.AboutTheContent.OriginalLanguage
		aboutTheContent.Supplier = request.AboutTheContent.Supplier
		aboutTheContent.AcquisitionDepartment = request.AboutTheContent.AcquisitionDepartment
		aboutTheContent.EnglishSynopsis = request.AboutTheContent.EnglishSynopsis
		aboutTheContent.ArabicSynopsis = request.AboutTheContent.ArabicSynopsis
		aboutTheContent.ProductionYear = request.AboutTheContent.ProductionYear
		aboutTheContent.ProductionHouse = request.AboutTheContent.ProductionHouse
		aboutTheContent.AgeGroup = request.AboutTheContent.AgeGroup
		aboutTheContent.IntroDuration = request.AboutTheContent.IntroDuration
		aboutTheContent.OutroDuration = request.AboutTheContent.OutroDuration
		aboutTheContent.IntroStart = request.AboutTheContent.IntroStart
		aboutTheContent.OutroStart = request.AboutTheContent.OutroStart
		if err := tx.Debug().Create(&aboutTheContent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var productionCountry ProductionCountry
		if request.AboutTheContent.ProductionCountries != nil && len(request.AboutTheContent.ProductionCountries) > 0 {
			for _, country := range request.AboutTheContent.ProductionCountries {
				productionCountry.AboutTheContentInfoId = aboutTheContent.Id
				productionCountry.CountryId = country
				if err := tx.Debug().Create(&productionCountry).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		var contentTranslation ContentTranslationRequest
		contentTranslation.LanguageType = common.ContentLanguageOriginTypes(request.Translation.LanguageType)
		contentTranslation.DubbingLanguage = request.Translation.DubbingLanguage
		contentTranslation.DubbingDialectId = request.Translation.DubbingDialectId
		contentTranslation.SubtitlingLanguage = request.Translation.SubtitlingLanguage
		if err := tx.Debug().Table("content_translation").Create(&contentTranslation).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var contentRights ContentRights
		contentRights.DigitalRightsType = request.Rights.DigitalRightsType
		//	DRSDate, _ := time.Parse(DateTimeFormat, request.Rights.DigitalRightsStartDate)
		contentRights.DigitalRightsStartDate = request.Rights.DigitalRightsStartDate
		//	DREDate, _ := time.Parse(DateTimeFormat, request.Rights.DigitalRightsEndDate)
		contentRights.DigitalRightsEndDate = request.Rights.DigitalRightsEndDate
		if err := tx.Debug().Create(&contentRights).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		RightsId := contentRights.Id
		var seasonKey Season
		tx.Debug().Select("season_key,number").Order("season_key desc").Limit(1).Find(&seasonKey)
		// for creating old seaesons take season id and user id from request body
		season.Id = request.SecondarySeasonId
		//	season.CreatedByUserId = request.CreatedByUserId
		season.SeasonKey = request.SeasonKey
		// season.SeasonKey = seasonKey.SeasonKey + 1
		fmt.Println(request.ContentId, "request content id is")
		season.ContentId = *request.ContentId
		season.Status = 3
		season.PrimaryInfoId = primaryInfo.ID
		season.AboutTheContentInfoId = aboutTheContent.Id
		season.Number = request.PrimaryInfo.SeasonNumber
		season.TranslationId = contentTranslation.Id
		season.CastId = castId
		season.CreatedByUserId = userid.(string)
		season.MusicId = MusicId
		season.TagInfoId = TaginfoId
		season.RightsId = RightsId
		season.CreatedAt = time.Now()
		season.ModifiedAt = time.Now()
		season.EnglishMetaTitle = request.SeoDetails.EnglishMetaTitle
		season.ArabicMetaTitle = request.SeoDetails.ArabicMetaTitle
		season.EnglishMetaDescription = request.SeoDetails.EnglishMetaDescription
		season.ArabicMetaDescription = request.SeoDetails.ArabicMetaDescription
		if request.NonTextualData.PosterImage != "" && request.NonTextualData.OverlayPosterImage != "" && request.NonTextualData.DetailsBackground != "" && request.NonTextualData.MobileDetailsBackground != "" {
			season.HasPosterImage = "true"
			season.HasOverlayPosterImage = "true"
			season.HasDetailsBackground = "true"
			season.HasMobileDetailsBackground = "true"
		} else {
			season.HasPosterImage = "false"
			season.HasOverlayPosterImage = "false"
			season.HasDetailsBackground = "false"
			season.HasMobileDetailsBackground = "false"
		}
		if len(request.Rights.DigitalRightsRegionsint) == 241 {
			season.HasAllRights = true
		}
		if err := tx.Debug().Create(&season).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		if request.SeasonGenres != nil {
			for i, genre := range request.SeasonGenres {
				var seasonGenre SeasonGenre
				seasonGenre.GenreId = genre.GenreId
				seasonGenre.SeasonId = season.Id
				seasonGenre.Order = i + 1
				if err := tx.Debug().Create(&seasonGenre).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}

				for j, subgenre := range genre.SubgenresId {
					var seasonSubgenre SeasonSubgenre
					seasonSubgenre.SeasonGenreId = seasonGenre.Id
					seasonSubgenre.SubgenreId = subgenre
					seasonSubgenre.Order = j + 1
					if err := tx.Debug().Create(&seasonSubgenre).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}
		}

		if len(*request.VarianceTrailers) != 0 {
			for i, trailerrange := range *request.VarianceTrailers {
				if trailerrange.VideoTrailerId != "" {
					_, _, duration := common.GetVideoDuration(trailerrange.VideoTrailerId)
					if duration == 0 {
						serverError = common.ServerError{Error: "InValid Content TrailerId", Description: "Please provide valid Video TrailerId", Code: "", RequestId: randstr.String(32)}
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
					var varianceTrailer VarianceTrailer
					if trailerrange.TrailerPosterImage != "" {
						varianceTrailer.HasTrailerPosterImage = true
					} else {
						varianceTrailer.HasTrailerPosterImage = false
					}
					// for sync add trailer id
					if len(request.VarianceTrailerIds) > 0 {
						varianceTrailer.Id = request.VarianceTrailerIds[i]
					}
					varianceTrailer.EnglishTitle = trailerrange.EnglishTitle
					varianceTrailer.ArabicTitle = trailerrange.ArabicTitle
					varianceTrailer.VideoTrailerId = trailerrange.VideoTrailerId
					varianceTrailer.SeasonId = season.Id
					varianceTrailer.Order = i + 1
					varianceTrailer.Duration = duration
					if err := tx.Debug().Table("variance_trailer").Create(&varianceTrailer).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
					go SeasonVarianceTrailerImageUploadGcp(*request.ContentId, request.SeasonId, trailerrange.Id, *request.ContentId, season.Id, varianceTrailer.Id, trailerrange.TrailerPosterImage)
				}
			}
		}

		if request.Cast.Actors != nil {
			for _, actor := range request.Cast.Actors {
				var contentActor ContentActor
				contentActor.CastId = castId
				contentActor.ActorId = actor
				if err := tx.Debug().Create(&contentActor).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Cast.Writers != nil {
			for _, writer := range request.Cast.Writers {
				var contentWriter ContentWriter
				contentWriter.CastId = castId
				contentWriter.WriterId = writer
				if err := tx.Debug().Create(&contentWriter).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Cast.Actors != nil {
			for _, director := range request.Cast.Directors {
				var contentDirector ContentDirector
				contentDirector.CastId = castId
				contentDirector.DirectorId = director
				if err := tx.Debug().Create(&contentDirector).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Music.Singers != nil {
			for _, singer := range request.Music.Singers {
				var contentSingers ContentSinger
				contentSingers.MusicId = MusicId
				contentSingers.SingerId = singer
				if err := tx.Debug().Create(&contentSingers).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Music.MusicComposers != nil {
			for _, musicComposer := range request.Music.MusicComposers {
				var contentMusicComposer ContentMusicComposer
				contentMusicComposer.MusicId = MusicId
				contentMusicComposer.MusicComposerId = musicComposer
				if err := tx.Debug().Create(&contentMusicComposer).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.Music.SongWriters != nil {
			for _, songWriter := range request.Music.SongWriters {
				var contentSongWriter ContentSongWriter
				contentSongWriter.MusicId = MusicId
				contentSongWriter.SongWriterId = songWriter
				if err := tx.Debug().Create(&contentSongWriter).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}

		if request.TagInfo.Tags != nil {
			for _, tag := range request.TagInfo.Tags {
				var contentTag ContentTag
				contentTag.TagInfoId = TaginfoId
				contentTag.TextualDataTagId = tag
				if err := tx.Debug().Create(&contentTag).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		var contentRightsCountry ContentRightsCountry
		var contentRightsCountrys []interface{}
		for _, country := range request.Rights.DigitalRightsRegionsint {
			fmt.Println(country)
			contentRightsCountry.ContentRightsId = RightsId
			contentRightsCountry.CountryId = country
			contentRightsCountrys = append(contentRightsCountrys, contentRightsCountry)
		}

		if err := gormbulk.BulkInsert(tx.Debug(), contentRightsCountrys, common.BULK_INSERT_LIMIT); err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		if request.Rights.SubscriptionPlans != nil && len(request.Rights.SubscriptionPlans) > 0 {
			for _, plan := range request.Rights.SubscriptionPlans {
				var contentRightsPlan ContentRightsPlan
				contentRightsPlan.RightsId = RightsId
				contentRightsPlan.SubscriptionPlanId = plan
				if err := tx.Debug().Create(&contentRightsPlan).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		var rightsProduct RightsProduct
		if request.Products != nil && len(*request.Products) > 0 {
			for _, product := range *request.Products {
				rightsProduct.RightsId = RightsId
				rightsProduct.ProductName = product
				if err := tx.Debug().Create(&rightsProduct).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
		}
		if err := tx.Debug().Table("content").Where("id=?", season.ContentId).Update("modified_at", time.Now()).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		//	 notFound := common.NotFoundErrorResponse()

		var episode []Episode
		// fetch episode
		if err := tx.Debug().Raw("select distinct on (number) number ,* from episode where season_id=? ", request.SeasonId).Find(&episode).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		// episode details
		for _, data := range episode {
			// fetching primary info details
			var primaryinodetails ContentPrimaryInfo
			fmt.Println(data.PrimaryInfoId, "primary info id is")
			if err := tx.Debug().Raw("select * from content_primary_info where id =?", data.PrimaryInfoId).Find(&primaryinodetails).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			// creating primary info
			var primarydata ContentPrimaryInfo
			primarydata = ContentPrimaryInfo{OriginalTitle: primaryinodetails.OriginalTitle, AlternativeTitle: primaryinodetails.AlternativeTitle, ArabicTitle: primaryinodetails.AlternativeTitle, TransliteratedTitle: primaryinodetails.TransliteratedTitle, Notes: primaryinodetails.Notes, IntroStart: primaryinodetails.IntroStart, OutroStart: primaryinodetails.OutroStart}
			if episodeprimaryinfo := tx.Debug().Table("content_primary_info").Create(&primarydata).Error; episodeprimaryinfo != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			// feltching playback item details
			var playbackitem PlaybackItem
			if err := tx.Debug().Raw("select * from playback_item where id =?", data.PlaybackItemId).Find(&playbackitem).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			//	fetching contentrights
			var contentrights ContentRights
			if err := tx.Debug().Raw("select * from content_rights where id =?", playbackitem.RightsId).Find(&contentrights).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			// creating content rights

			contentrightsdata := ContentRights{DigitalRightsType: contentrights.DigitalRightsType,
				DigitalRightsStartDate: contentrights.DigitalRightsEndDate, DigitalRightsEndDate: contentrights.DigitalRightsEndDate}
			//contentrightsdetails = append(contentrightsdetails, contentrightsdata)
			if err := tx.Debug().Table("content_rights").Create(&contentrightsdata).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			//	fetching content translation details
			var contenttranslation ContentTranslationData
			if err := tx.Debug().Raw("select * from content_translation where id =?", playbackitem.TranslationId).Find(&contenttranslation).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			var contenttransdata ContentTranslationData
			// creating content  translation details
			contenttransdata.LanguageType = contenttranslation.LanguageType
			contenttransdata.DubbingLanguage = contenttranslation.DubbingLanguage
			contenttransdata.DubbingDialectId = contenttranslation.DubbingDialectId
			contenttransdata.SubtitlingLanguage = contenttranslation.SubtitlingLanguage
			//	contenttarnsfinal = append(contenttarnsfinal, contentdata)
			if err := tx.Debug().Table("content_translation").Create(&contenttransdata).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			//create playback item
			playbackdata := PlaybackItem{VideoContentId: playbackitem.VideoContentId, SchedulingDateTime: playbackitem.SchedulingDateTime, TranslationId: contenttransdata.Id, RightsId: contentrightsdata.Id, Duration: playbackitem.Duration}
			if createPlaybackItem := tx.Debug().Table("playback_item").Create(&playbackdata).Error; createPlaybackItem != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			// feltching platform details
			var newplatform []interface{}
			var TargetPlatformvalues []PlaybackItemTargetPlatform
			tx.Debug().Table("playback_item_target_platform").Select("target_platform").Where("playback_item_id=?", playbackitem.Id).Find(&TargetPlatformvalues)

			for _, value := range TargetPlatformvalues {
				playbackdetails := PlaybackItemTargetPlatform{PlaybackItemId: playbackdata.Id, TargetPlatform: value.TargetPlatform, RightsId: contentrightsdata.Id}
				newplatform = append(newplatform, playbackdetails)
			}

			// inserting publishing platforms
			if err := gormbulk.BulkInsert(tx, newplatform, common.BULK_INSERT_LIMIT); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			// fetching content cast details
			var contentcast ContentCast
			if err := tx.Debug().Raw("select * from content_cast where id=?", data.CastId).Find(&contentcast).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			// inserting content cast details
			contentcastdata := ContentCast{MainActorId: contentcast.MainActorId, MainActressId: contentcast.MainActressId}
			if err := tx.Debug().Table("content_cast").Create(&contentcastdata).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			// fetching content actor details
			var contentactor ContentActor
			tx.Debug().Raw("select * from content_actor where cast_id=?", data.CastId).Find(&contentactor)
			// inserting content actor details
			if contentactor.ActorId != "" {
				contentactordata := ContentActor{CastId: contentcastdata.Id, ActorId: contentactor.ActorId}
				if err := tx.Debug().Table("content_actor").Create(&contentactordata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			// fetching content writer details
			var contentwriter ContentWriter
			tx.Debug().Raw("select * from content_writer where cast_id=?", data.CastId).Find(&contentwriter)

			// inserting content writer details
			if contentwriter.WriterId != "" {
				contentwriterdata := ContentWriter{CastId: contentcastdata.Id, WriterId: contentwriter.WriterId}
				if err := tx.Debug().Table("content_writer").Create(&contentwriterdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			// fetching content director details
			var contentdirector ContentDirector
			tx.Debug().Raw("select * from content_director where cast_id=?", data.CastId).Find(&contentdirector)
			// inserting content director details
			if contentdirector.DirectorId != "" {
				contentdirectorrdata := ContentDirector{CastId: contentcastdata.Id, DirectorId: contentdirector.DirectorId}
				if err := tx.Debug().Table("content_director").Create(&contentdirectorrdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			// inserting contetn music details
			var contentmusic ContentMusic
			if err := tx.Debug().Table("content_music").Create(&contentmusic).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			fmt.Println(contentmusic.Id, "content music")
			// fetching content music composer details
			var contentMusicComposer ContentMusicComposer
			tx.Debug().Raw("select * from content_music_composer where music_id=?", data.MusicId).Find(&contentMusicComposer)

			// inserting content music composer  details
			if contentMusicComposer.MusicComposerId != "" {
				contentmusicdata := ContentMusicComposer{MusicComposerId: contentMusicComposer.MusicComposerId, MusicId: contentmusic.Id}
				if err := tx.Debug().Table("content_music_composer").Create(&contentmusicdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			// fetching content music singer details
			var contentsinger ContentSinger
			tx.Debug().Raw("select * from content_singer where music_id=?", data.MusicId).Find(&contentsinger)

			// inserting content music singer  details
			if contentsinger.SingerId != "" {
				contentsingerdata := ContentSinger{SingerId: contentsinger.SingerId, MusicId: contentMusic.Id}
				if err := tx.Debug().Table("content_singer").Create(&contentsingerdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}

			// fetching content music song writer details
			var contentsongwriter ContentSongWriter
			tx.Debug().Raw("select * from content_song_writer where music_id=?", data.MusicId).Find(&contentsongwriter)

			// inserting content music song writer  details
			if contentsongwriter.SongWriterId != "" {
				contentsongdata := ContentSongWriter{SongWriterId: contentsongwriter.SongWriterId, MusicId: contentMusic.Id}
				if err := tx.Debug().Table("content_song_writer").Create(&contentsongdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			// inserting tag info details
			var contenttaginfo ContentTagInfo
			if err := tx.Debug().Table("content_tag_info").Create(&contenttaginfo).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			// fetching content _tag details
			var contenttag ContentTag
			tx.Debug().Raw("select * from content_tag where tag_info_id=?", data.TagInfoId).Find(&contenttag)
			var textualdatatag TextualDataTag
			if contenttag.TextualDataTagId != "" {
				// fetching textual data tag details
				tx.Debug().Raw("select * from textual_data_tag where id=?", contenttag.TextualDataTagId).Find(&textualdatatag)
			}

			// inserting textual data tag details
			textualdata := TextualDataTag{Name: textualdatatag.Name}
			if textualdatatag.Name != "" {
				if err := tx.Debug().Table("textual_data_tag").Create(&textualdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}

			// inserting content_tag details
			if textualdata.Id != "" {
				contenttagdata := ContentTag{TagInfoId: contenttaginfo.Id, TextualDataTagId: textualdata.Id}
				if err := tx.Debug().Table("content_tag").Create(&contenttagdata).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
			}
			var episodekey Episode
			tx.Debug().Table("episode").Select("max(episode_key) as episode_key").Find(&episodekey)
			//	var episodefinal []interface{}
			var episodeDetails Episode
			episodeDetails.Number = data.Number
			episodeDetails.SeasonId = season.Id
			episodeDetails.PrimaryInfoId = primarydata.ID
			episodeDetails.PlaybackItemId = playbackdata.Id
			episodeDetails.Status = 3
			episodeDetails.SynopsisEnglish = data.SynopsisArabic
			episodeDetails.SynopsisArabic = data.SynopsisArabic
			episodeDetails.CastId = contentcastdata.Id
			episodeDetails.MusicId = contentMusic.Id
			episodeDetails.TagInfoId = contenttaginfo.Id
			episodeDetails.EpisodeKey = episodekey.EpisodeKey + 1
			episodeDetails.CreatedAt = time.Now()
			episodeDetails.ModifiedAt = time.Now()
			episodeDetails.EnglishMetaTitle = data.EnglishMetaTitle
			episodeDetails.ArabicMetaTitle = data.ArabicMetaTitle
			episodeDetails.EnglishMetaDescription = data.EnglishMetaDescription
			episodeDetails.ArabicMetaDescription = data.ArabicMetaDescription
			episodeDetails.HasPosterImage = data.HasPosterImage
			episodeDetails.HasSubtitlingScript = data.HasSubtitlingScript
			episodeDetails.HasDubbingScript = data.HasDubbingScript
			//	episodefinal = append(episodefinal, episodeDetails)

			var image []Images

			image = append(image, Images{Imagename: "poster-image", HasImage: data.HasPosterImage})
			image = append(image, Images{Imagename: "dubbing-script", HasImage: data.HasDubbingScript})
			image = append(image, Images{Imagename: "subtitling-script", HasImage: data.HasSubtitlingScript})

			if err := tx.Debug().Table("episode").Create(&episodeDetails).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				fmt.Println(serverError, "server error")
				return
			}
			go SeasonVarianceEpisodeImageUploadGcp(*request.ContentId, season.Id, episodeDetails.Id, image, *request.ContentId, request.SeasonId, data.Id)
		}
		/*commit changes*/
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		seasonId := season.Id
		SeasonVarianceFileUPloadGcp(request.SeasonId, seasonId, season.ContentId)
		/* Prepare Redis Cache for all contents */
		db.Debug().Raw("select content_key,content_type from content where id=?", season.ContentId).Find(&seasondetails)
		/* Prepare Redis Cache for single content*/
		contentkeyconverted := strconv.Itoa(seasondetails.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(seasondetails.ContentType, c)
		c.JSON(http.StatusOK, gin.H{"message": "Season Created Successfully.", "status": http.StatusOK, "id": season.Id})
		return
	} else {
		if err := tx.Debug().Where("id=?", c.Param("id")).Find(&season).Error; err != nil {
			c.JSON(http.StatusNotFound, notFound)
			return
		}
		result, _ := UpdateSeasonDetails(request, c, season, 3)
		if result == "" {
			c.JSON(http.StatusNotFound, serverError)
			return
		}
		res := map[string]string{
			"id": c.Param("id"),
		}
		seasonId := c.Param("id")
		SeasonFileUPloadGcp(request, seasonId, season.ContentId)
		/* Prepare Redis Cache for all contents */
		db.Debug().Raw("select content_key,content_type from content where id=?", season.ContentId).Find(&seasondetails)
		/* Prepare Redis Cache for single content*/
		fmt.Println("KEY", seasondetails.ContentType, seasondetails.ContentKey)
		contentkeyconverted := strconv.Itoa(seasondetails.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(seasondetails.ContentType, c)
		c.JSON(http.StatusOK, gin.H{"data": res})
		return
	}
}
func UpdateSeasonDetails(request CreateSeasonRequestValidation, c *gin.Context, seasonDetails Season, status int) (string, common.ServerError) {

	db := c.MustGet("DB").(*gorm.DB)
	ctx := context.Background()
	tx := db.Debug().BeginTx(ctx, nil)
	var countryid []ContentRightsCountry
	tx.Debug().Table("season s").Select("crc.country_id").Joins("join content_rights_country crc on crc.content_rights_id=s.rights_id").Where("s.content_id=? and s.number=? and s.id!=(?) and s.deleted_by_user_id is null ", request.ContentId, request.PrimaryInfo.SeasonNumber, c.Param("id")).Find(&countryid)
	var countryflag bool
	countryflag = false
	for _, data := range countryid {
		for _, value := range request.Rights.DigitalRightsRegionsint {
			if data.CountryId == value {
				countryflag = true
				break
			}
		}
	}

	error11 := common.ServerError{Error: "countries exists", Description: "Selected countries for this variant are not allowed.", Code: "", RequestId: randstr.String(32)}
	if countryflag {
		return "", error11
	}
	userid := c.MustGet("userid")
	var season Season
	serverError := common.ServerErrorResponse()
	var contentCast ContentCast
	contentCast.MainActorId = request.Cast.MainActorId
	contentCast.MainActressId = request.Cast.MainActressId
	if err := tx.Debug().Model(&contentCast).Where("id=?", seasonDetails.CastId).Update(&contentCast).Error; err != nil {
		return "", serverError
	}
	castId := seasonDetails.CastId
	MusicId := seasonDetails.MusicId
	TaginfoId := seasonDetails.TagInfoId
	var primaryInfo ContentPrimaryInfo
	primaryInfo.OriginalTitle = request.PrimaryInfo.OriginalTitle
	primaryInfo.AlternativeTitle = request.PrimaryInfo.AlternativeTitle
	primaryInfo.ArabicTitle = request.PrimaryInfo.ArabicTitle
	primaryInfo.TransliteratedTitle = request.PrimaryInfo.TransliteratedTitle
	primaryInfo.Notes = request.PrimaryInfo.Notes
	primaryInfo.IntroStart = request.IntroStart
	primaryInfo.OutroStart = ""
	if err := tx.Debug().Model(&primaryInfo).Where("id=?", seasonDetails.PrimaryInfoId).Update(&primaryInfo).Error; err != nil {
		return "", serverError
	}
	var aboutTheContent AboutTheContentInfo
	aboutTheContent.OriginalLanguage = request.AboutTheContent.OriginalLanguage
	aboutTheContent.Supplier = request.AboutTheContent.Supplier
	aboutTheContent.AcquisitionDepartment = request.AboutTheContent.AcquisitionDepartment
	aboutTheContent.EnglishSynopsis = request.AboutTheContent.EnglishSynopsis
	aboutTheContent.ArabicSynopsis = request.AboutTheContent.ArabicSynopsis
	aboutTheContent.ProductionYear = request.AboutTheContent.ProductionYear
	aboutTheContent.ProductionHouse = request.AboutTheContent.ProductionHouse
	aboutTheContent.AgeGroup = request.AboutTheContent.AgeGroup
	aboutTheContent.IntroDuration = request.AboutTheContent.IntroDuration
	aboutTheContent.OutroDuration = request.AboutTheContent.OutroDuration
	aboutTheContent.IntroStart = request.AboutTheContent.IntroStart
	aboutTheContent.OutroStart = request.AboutTheContent.OutroStart
	if err := tx.Debug().Model(&aboutTheContent).Where("id=?", seasonDetails.AboutTheContentInfoId).Update(&aboutTheContent).Error; err != nil {
		return "", serverError
	}
	var productionCountry ProductionCountry
	db.Debug().Where("about_the_content_info_id=?", seasonDetails.AboutTheContentInfoId).Delete(&productionCountry)
	if request.AboutTheContent.ProductionCountries != nil && len(request.AboutTheContent.ProductionCountries) > 0 {
		for _, country := range request.AboutTheContent.ProductionCountries {
			var productionCountry ProductionCountry
			productionCountry.AboutTheContentInfoId = seasonDetails.AboutTheContentInfoId
			productionCountry.CountryId = country
			if err := db.Debug().Create(&productionCountry).Error; err != nil {
				return "", serverError
			}
		}
	}
	var contentTranslation ContentTranslationRequest
	contentTranslation.LanguageType = common.ContentLanguageOriginTypes(request.Translation.LanguageType)
	contentTranslation.DubbingLanguage = request.Translation.DubbingLanguage
	contentTranslation.DubbingDialectId = request.Translation.DubbingDialectId
	contentTranslation.SubtitlingLanguage = request.Translation.SubtitlingLanguage
	if err := db.Debug().Table("content_translation").Where("id=?", seasonDetails.TranslationId).Update(&contentTranslation).Error; err != nil {
		return "", serverError
	}
	var contentRights ContentRights
	contentRights.DigitalRightsType = request.Rights.DigitalRightsType
	//	DRSDate, _ := time.Parse(DateTimeFormat, request.Rights.DigitalRightsStartDate)
	contentRights.DigitalRightsStartDate = request.Rights.DigitalRightsStartDate
	//	DREDate, _ := time.Parse(DateTimeFormat, request.Rights.DigitalRightsEndDate)
	contentRights.DigitalRightsEndDate = request.Rights.DigitalRightsEndDate
	if err := tx.Debug().Model(&contentRights).Where("id=?", seasonDetails.RightsId).Update(&contentRights).Error; err != nil {
		return "", serverError
	}
	RightsId := seasonDetails.RightsId
	season.SeasonKey = seasonDetails.SeasonKey
	season.ContentId = *request.ContentId
	season.Status = status
	season.PrimaryInfoId = seasonDetails.PrimaryInfoId
	season.AboutTheContentInfoId = seasonDetails.AboutTheContentInfoId
	season.Number = seasonDetails.Number
	season.TranslationId = seasonDetails.TranslationId
	season.CastId = castId
	season.CreatedByUserId = userid.(string)
	season.MusicId = MusicId
	season.TagInfoId = TaginfoId
	season.RightsId = RightsId
	//	season.DeletedByUserId = "00000000-0000-0000-0000-000000000000"
	season.ModifiedAt = time.Now()
	season.EnglishMetaTitle = request.SeoDetails.EnglishMetaTitle
	season.ArabicMetaTitle = request.SeoDetails.ArabicMetaTitle
	season.EnglishMetaDescription = request.SeoDetails.EnglishMetaDescription
	season.ArabicMetaDescription = request.SeoDetails.ArabicMetaDescription
	if request.NonTextualData.PosterImage != "" && request.NonTextualData.OverlayPosterImage != "" && request.NonTextualData.DetailsBackground != "" && request.NonTextualData.MobileDetailsBackground != "" {
		season.HasPosterImage = "true"
		season.HasOverlayPosterImage = "true"
		season.HasDetailsBackground = "true"
		season.HasMobileDetailsBackground = "true"
	} else {
		season.HasPosterImage = "false"
		season.HasOverlayPosterImage = "false"
		season.HasDetailsBackground = "false"
		season.HasMobileDetailsBackground = "false"
	}

	if len(request.Rights.DigitalRightsRegionsint) == 241 {
		season.HasAllRights = true
	}
	if err := tx.Debug().Model(&season).Where("id=?", seasonDetails.Id).Update(&season).Error; err != nil {
		return "", serverError
	}
	var genreDeatils []SeasonGenre
	db.Debug().Where("season_id=?", seasonDetails.Id).Find(&genreDeatils)
	var ids []string
	for _, details := range genreDeatils {
		ids = append(ids, details.GenreId)
	}
	var seasonGenre SeasonGenre
	var seasonSubgenre SeasonSubgenre
	tx.Debug().Where("season_id=?", seasonDetails.Id).Delete(&seasonGenre)
	tx.Debug().Where("season_genre_id in(?)", ids).Delete(&seasonSubgenre)
	if request.SeasonGenres != nil {
		for i, genre := range request.SeasonGenres {
			var seasonGenre SeasonGenre
			seasonGenre.GenreId = genre.GenreId
			seasonGenre.SeasonId = c.Param("id")
			seasonGenre.Order = i
			if err := tx.Debug().Create(&seasonGenre).Error; err != nil {
				return "", serverError
			}

			for j, subgenre := range genre.SubgenresId {
				var seasonSubgenre SeasonSubgenre
				seasonSubgenre.SeasonGenreId = seasonGenre.Id
				seasonSubgenre.SubgenreId = subgenre
				seasonSubgenre.Order = j
				if err := tx.Debug().Create(&seasonSubgenre).Error; err != nil {
					return "", serverError
				}
			}
		}
	}
	var varincetrailers []VarianceTrailer
	tx.Debug().Table("variance_trailer").Select("id").Where("season_id=?", c.Param("id")).Find(&varincetrailers)

	var newarr []string
	var exists bool
	for _, trailers := range varincetrailers {
		exists = false
		for _, newtratrailers := range *request.VarianceTrailers {
			if trailers.Id == newtratrailers.Id {
				exists = true
				break
			}
		}
		if !exists {
			newarr = append(newarr, trailers.Id)
		}
	}

	if len(newarr) != 0 {
		var varincetrailer VarianceTrailer
		if err := tx.Debug().Table("variance_trailer").Where("id in(?)", newarr).Delete(&varincetrailer).Error; err != nil {
			return "", serverError
		}
	}
	var orders int
	orders = 0
	if len(*request.VarianceTrailers) != 0 {
		for i, trailerrange := range *request.VarianceTrailers {
			if trailerrange.Id == "" {
				if trailerrange.VideoTrailerId != "" {
					_, _, duration := common.GetVideoDuration(trailerrange.VideoTrailerId)
					if duration == 0 {
						serverError = common.ServerError{Error: "InValid Content TrailerId", Description: "Please provide valid Video TrailerId", Code: "", RequestId: randstr.String(32)}
						return "", serverError
					}
					var varianceTrailer VarianceTrailer
					if trailerrange.TrailerPosterImage != "" {
						varianceTrailer.HasTrailerPosterImage = true
					} else {
						varianceTrailer.HasTrailerPosterImage = false
					}
					// for sync add trailer id
					fmt.Println(orders, "llllllllllllllll")
					if len(request.VarianceTrailerIds) > 0 {
						varianceTrailer.Id = request.VarianceTrailerIds[orders]
					}
					varianceTrailer.EnglishTitle = trailerrange.EnglishTitle
					varianceTrailer.ArabicTitle = trailerrange.ArabicTitle
					varianceTrailer.VideoTrailerId = trailerrange.VideoTrailerId
					varianceTrailer.SeasonId = seasonDetails.Id
					varianceTrailer.Order = i + 1
					varianceTrailer.Duration = duration
					if err := tx.Debug().Table("variance_trailer").Create(&varianceTrailer).Error; err != nil {
						return "", serverError
					}
					orders = orders + 1
					go SeasonTrailerImageUploadGcp(*request.ContentId, c.Param("id"), varianceTrailer.Id, trailerrange.TrailerPosterImage)
				}
			} else {
				var trailers VarianceTrailer
				if trailerrange.VideoTrailerId == "" {
					if err := tx.Debug().Table("variance_trailer").Where("id = ?", trailerrange.Id).Delete(&VarianceTrailer{}).Error; err != nil {
						return "", serverError
					}
				} else {
					_, _, duration := common.GetVideoDuration(trailerrange.VideoTrailerId)
					if duration == 0 {
						serverError = common.ServerError{Error: "InValid Content TrailerId", Description: "Please provide valid Video TrailerId", Code: "", RequestId: randstr.String(32)}
						return "", serverError
					}
					trailers.EnglishTitle = trailerrange.EnglishTitle
					trailers.ArabicTitle = trailerrange.ArabicTitle
					trailers.VideoTrailerId = trailerrange.VideoTrailerId
					trailers.Duration = duration
					if err := tx.Debug().Table("variance_trailer").Where("id=?", trailerrange.Id).Update(&trailers).Error; err != nil {
						return "", serverError
					}
					go SeasonTrailerImageUploadGcp(*request.ContentId, c.Param("id"), trailerrange.Id, trailerrange.TrailerPosterImage)
				}
			}
		}
	}

	var contentActor ContentActor
	tx.Debug().Where("cast_id=?", castId).Delete(&contentActor)
	if request.Cast.Actors != nil {
		for _, actor := range request.Cast.Actors {
			var contentActor ContentActor
			contentActor.CastId = castId
			contentActor.ActorId = actor
			if err := tx.Debug().Create(&contentActor).Error; err != nil {
				return "", serverError
			}
		}
	}
	var contentWriter ContentWriter
	tx.Debug().Where("cast_id=?", castId).Delete(&contentWriter)
	if request.Cast.Writers != nil {
		for _, writer := range request.Cast.Writers {
			var contentWriter ContentWriter
			contentWriter.CastId = castId
			contentWriter.WriterId = writer
			if err := tx.Debug().Create(&contentWriter).Error; err != nil {
				return "", serverError
			}
		}
	}
	var contentDirector ContentDirector
	tx.Debug().Where("cast_id=?", castId).Delete(&contentDirector)
	if request.Cast.Directors != nil {
		for _, director := range request.Cast.Directors {
			var contentDirector ContentDirector
			contentDirector.CastId = castId
			contentDirector.DirectorId = director
			if err := tx.Debug().Create(&contentDirector).Error; err != nil {
				return "", serverError
			}
		}
	}
	var contentSingers ContentSinger
	tx.Debug().Where("music_id=?", MusicId).Delete(&contentSingers)
	if request.Music.Singers != nil {
		for _, singer := range request.Music.Singers {
			var contentSingers ContentSinger
			contentSingers.MusicId = MusicId
			contentSingers.SingerId = singer
			if err := tx.Debug().Create(&contentSingers).Error; err != nil {
				return "", serverError
			}
		}
	}
	var contentMusicComposer ContentMusicComposer
	tx.Debug().Where("music_id=?", MusicId).Delete(&contentMusicComposer)
	if request.Music.MusicComposers != nil {
		for _, musicComposer := range request.Music.MusicComposers {
			var contentMusicComposer ContentMusicComposer
			contentMusicComposer.MusicId = MusicId
			contentMusicComposer.MusicComposerId = musicComposer
			if err := tx.Debug().Create(&contentMusicComposer).Error; err != nil {
				return "", serverError
			}
		}
	}
	var contentSongWriter ContentSongWriter
	tx.Debug().Where("music_id=?", MusicId).Delete(&contentSongWriter)
	if request.Music.SongWriters != nil {
		for _, songWriter := range request.Music.SongWriters {
			var contentSongWriter ContentSongWriter
			contentSongWriter.MusicId = MusicId
			contentSongWriter.SongWriterId = songWriter
			if err := tx.Debug().Create(&contentSongWriter).Error; err != nil {
				return "", serverError
			}
		}
	}
	var contentTag ContentTag
	tx.Debug().Where("tag_info_id=?", TaginfoId).Delete(&contentTag)
	if request.TagInfo.Tags != nil {
		for _, tag := range request.TagInfo.Tags {
			var contentTag ContentTag
			contentTag.TagInfoId = TaginfoId
			contentTag.TextualDataTagId = tag
			if err := tx.Debug().Create(&contentTag).Error; err != nil {
				return "", serverError
			}
		}
	}
	var contentRightsCountry ContentRightsCountry
	var contentRightsCountrys []interface{}
	var contentsdelete ContentRightsCountry
	if request.Rights.DigitalRightsRegionsint != nil {
		if err := tx.Debug().Table("content_rights_country").Where("content_rights_id=?", RightsId).Delete(&contentsdelete).Error; err != nil {
			return "", serverError
		}
	}
	if request.Rights.DigitalRightsRegionsint != nil {
		for _, country := range request.Rights.DigitalRightsRegionsint {
			contentRightsCountry.ContentRightsId = RightsId
			contentRightsCountry.CountryId = country
			contentRightsCountrys = append(contentRightsCountrys, contentRightsCountry)
		}
	}
	if err := gormbulk.BulkInsert(tx.Debug(), contentRightsCountrys, common.BULK_INSERT_LIMIT); err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return "", serverError
	}

	var contentRightsPlans ContentRightsPlan
	db.Debug().Where("rights_id=?", RightsId).Delete(&contentRightsPlans)
	if request.Rights.SubscriptionPlans != nil && len(request.Rights.SubscriptionPlans) > 0 {
		for _, plan := range request.Rights.SubscriptionPlans {
			var contentRightsPlan ContentRightsPlan
			contentRightsPlan.RightsId = RightsId
			contentRightsPlan.SubscriptionPlanId = plan
			if err := tx.Debug().Create(&contentRightsPlan).Error; err != nil {
				return "", serverError
			}
		}
	}
	var rightsProducts RightsProduct
	var rightsProduct RightsProduct
	tx.Debug().Where("rights_id=?", RightsId).Delete(&rightsProducts)
	if request.Products != nil && len(*request.Products) > 0 {
		for _, product := range *request.Products {
			rightsProduct.RightsId = RightsId
			rightsProduct.ProductName = product
			if err := tx.Debug().Create(&rightsProduct).Error; err != nil {
				return "", serverError
			}
		}
	}
	if err := tx.Debug().Table("content").Where("id=?", *request.ContentId).Update("modified_at", time.Now()).Error; err != nil {
		return "", serverError
	}
	fmt.Println(status, "status is")
	/*commit changes*/
	if err := tx.Commit().Error; err != nil {
		return "", serverError

	}
	if status == 1 {
		return "Season updated successfully.", serverError
	} else {
		return "Season drafted successfully.", serverError
	}

}

// GetContents -  Get all contents
// GET /api/contents
// @Summary Get all contents
// @Description Get all contents by filters
// @Tags Content
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Param contentType query string false "Content Type"
// @Param searchText query string false "Search Text"
// @Success 200 {array} object c.JSON
// @Router /api/contents [get]
func (hs *ContentService) GetAllContentDetails(c *gin.Context) {
	AuthorizationRequired := c.MustGet("AuthorizationRequired")
	userid := c.MustGet("userid")
	if AuthorizationRequired == 1 || userid == "" || c.MustGet("is_back_office_user") == "false" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	serverError := common.ServerErrorResponse()
	db := c.MustGet("DB").(*gorm.DB)
	udb := c.MustGet("UDB").(*gorm.DB)
	var limit, offset int64
	var searchText, contentType, Id string
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	if limit == 0 {
		limit = 50
	} else if limit > 100 {
		limit = 100
	}
	if c.Request.URL.Query()["contentType"] != nil {
		contentType = c.Request.URL.Query()["contentType"][0]
	}
	if c.Request.URL.Query()["searchText"] != nil {
		searchText = c.Request.URL.Query()["searchText"][0]
	}
	if c.Request.URL.Query()["id"] != nil {
		Id = c.Request.URL.Query()["id"][0]
	}
	//query result variables
	response := []GetAllContentDetails{}
	var contentDetails []GetAllContentDetails
	where := "c.deleted_by_user_id is null"
	if contentType != "" && contentType != "All" {
		where += " and c.content_type = '" + contentType + "'"
	}
	if Id != "" {
		where += " and c.id = '" + Id + "'"
	}
	if searchText != "" {
		where += " and (cpi.transliterated_title ilike '%" + searchText + "%' or cpi.arabic_title ilike '%" + searchText + "%')"
	}
	var totalCount int
	if err := db.Debug().Table("content c").Select("c.content_tier as type,c.content_key,case when c.status = 1 then 'Published' when c.status = 2 then 'Unpublished' when c.status = 3 then 'Draft' end as sub_status_name,case when c.status = 3 then 2 else c.status end as status,case when c.status = 3 then false else true end as status_can_be_changed,cpi.transliterated_title,c.id,c.created_by_user_id").Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").Where(where).Order("c.modified_at desc").Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if err := db.Debug().Table("content c").Select("c.content_tier as type,c.content_key,c.status,cpi.transliterated_title,c.id,c.created_by_user_id").Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").Where(where).Order("c.modified_at desc").Limit(limit).Offset(offset).Find(&contentDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	for _, content := range contentDetails {
		var userDetails UserDetails
		if err := udb.Debug().Table("user").Where("id=?", content.CreatedByUserId).Find(&userDetails).Error; err != nil && err.Error() != "record not found" {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		if content.Status == 1 {
			status := "Published"
			content.SubStatusName = status
			content.StatusCanBeChanged = true
		} else if content.Status == 2 {
			status := "Unpublished"
			content.SubStatusName = status
			content.StatusCanBeChanged = true
		} else if content.Status == 3 {
			status := "Draft"
			content.SubStatusName = status
			content.Status = 2
			content.StatusCanBeChanged = false
		}
		content.CreatedBy = userDetails.UserName
		variences := []ContentVariancesDetails{}
		content.ContentVariances = variences
		seasons := []ContentSeasonsDetails{}
		content.ContentSeasons = seasons
		if content.Type == 1 {
			var varienceResponse []ContentVariancesDetails
			var varianceDetails []ContentVariancesDetails
			if err := db.Debug().Table("content c").Select("ct.language_type,ct.dubbing_language,cv.status,ct.dubbing_dialect_id,ct.subtitling_language,cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date,json_agg(crc.country_id)::varchar as digital_rights,cv.has_all_rights as country_check,cv.id,pi2.scheduling_date_time").Joins("join content_variance cv on cv.content_id = c.id join playback_item pi2 on pi2.id = cv.playback_item_id join content_rights cr on cr.id = pi2.rights_id join content_rights_country crc on crc.content_rights_id = cr.id join content_translation ct on ct.id =pi2.translation_id").Where("c.id=? and cv.deleted_by_user_id is null", content.Id).Order("c.modified_at desc").Group("cv.status,cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date,cv.id,ct.language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,c.modified_at,pi2.scheduling_date_time").Find(&varianceDetails).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			for _, details := range varianceDetails {
				if content.Status == 2 {
					// content object related main content details, details obj related to variance details
					//	status := "Unpublished"
					if details.Status == 2 {
						details.SubStatusName = "Unpublished"
					} else if details.Status == 1 {
						// details.SubStatusName = "Published"
						// newdate := details.DigitalRightsEndDate
						// anotherdate := newdate.Format("2006-01-02")
						// startsdate := details.DigitalRightsStartDate
						// newdates := startsdate.Format("2006-01-02")
						var anotherdate string
						var newdates string
						if details.DigitalRightsEndDate != nil {
							newdate := details.DigitalRightsEndDate
							anotherdate = newdate.Format("2006-01-02")
						}
						if details.DigitalRightsEndDate != nil {
							startsdate := details.DigitalRightsStartDate
							newdates = startsdate.Format("2006-01-02")
						}
						if anotherdate != "" && newdates != "" {
							if anotherdate <= time.Now().Format("2006-01-02") {
								status := "Digital Rights Exceeded"
								details.Status = 2
								details.SubStatusName = status
								details.StatusCanBeChanged = false
							} else if newdates > time.Now().Format("2006-01-02") {
								status := "Unpublished"
								details.Status = 2
								details.SubStatusName = status
								details.StatusCanBeChanged = false
							} else {
								status := "Published"
								details.SubStatusName = status
								details.StatusCanBeChanged = true
							}
						} else if anotherdate == "" && newdates == "" {
							status := "Published"
							details.SubStatusName = status
							details.StatusCanBeChanged = true
						}
					} else if details.Status == 3 {
						details.SubStatusName = "Draft"
					}
					details.Status = 2
					details.StatusCanBeChanged = false
				} else if content.Status == 3 {
					status := "Draft"
					details.Status = 3
					details.SubStatusName = status
					details.StatusCanBeChanged = false
				} else if details.Status == 1 {
					fmt.Println("PPPPPPPPPPPPPPPPP")
					var anotherdate string
					var newdates string
					if details.DigitalRightsEndDate != nil {
						newdate := details.DigitalRightsEndDate
						anotherdate = newdate.Format("2006-01-02")
					}
					if details.DigitalRightsEndDate != nil {
						startsdate := details.DigitalRightsStartDate
						newdates = startsdate.Format("2006-01-02")
					}
					var scheduledate string
					if details.SchedulingDateTime != nil {
						//scheduledate = details.SchedulingDateTime.UTC().Format("2006-01-02 11:33:41")
						scheduledate = details.SchedulingDateTime.String()
					}
					if anotherdate != "" && newdates != "" {
						if anotherdate <= time.Now().Format("2006-01-02") {
							status := "Digital Rights Exceeded"
							details.Status = 2
							details.SubStatusName = status
							details.StatusCanBeChanged = false
						} else if newdates > time.Now().Format("2006-01-02") {
							status := "Scheduled"
							details.Status = 2
							details.SubStatusName = status
							details.StatusCanBeChanged = false
						} else if scheduledate > time.Now().String() && details.SchedulingDateTime != nil {
							//	status := "Unpublished"
							details.Status = 2
							details.SubStatusName = "Scheduled"
							details.StatusCanBeChanged = false
						} else {
							status := "Published"
							details.SubStatusName = status
							details.StatusCanBeChanged = true
							details.Status = 1
						}
					} else if anotherdate == "" && newdates == "" {
						status := "Published"
						details.SubStatusName = status
						details.StatusCanBeChanged = true
						details.Status = 1
					}
				} else if details.Status == 3 {
					status := "Draft"
					details.SubStatusName = status
					details.Status = 2
					details.StatusCanBeChanged = false
				} else if details.Status == 2 {
					details.Status = 2
					status := "Unpublished"
					details.SubStatusName = status
					details.StatusCanBeChanged = true
				}
				details.VideoContentId = nil
				details.OverlayPosterImage = nil
				details.DubbingScript = nil
				details.SubtitlingScript = nil
				regionsResult := make([]int, 0)
				if details.DigitalRights != "" {
					regions, _ := JsonStringToIntSliceOrMap(details.DigitalRights)
					regionsResult = RemoveDuplicateValues(regions)
					if regionsResult[0] == 0 {
						regionsResult = make([]int, 0)
					}
				}
				if details.SchedulingDateTime != nil {
					fmt.Println(details.SchedulingDateTime, "llllllllllllllllll", details.Status)
					// if details.Status == 1 {
					//scheduledate = details.SchedulingDateTime.UTC().Format("2006-01-02 11:33:41")
					scheduledatetime := details.SchedulingDateTime.String()
					if scheduledatetime > time.Now().String() && (details.Status == 1 || details.Status == 2) {
						//fmt.Println("helloooooooooooo")
						status := "Scheduled"
						details.SubStatusName = status
					} else if details.SchedulingDateTime.Format("2006-01-02 11:33:41") < time.Now().Format("2006-01-02 11:33:41") && details.Status == 1 {
						substatusname := "Published"
						details.SubStatusName = substatusname

						//else {
						//details.SubStatusName = "Draft"
					} else if details.SchedulingDateTime.Format("2006-01-02 11:33:41") < time.Now().Format("2006-01-02 11:33:41") && details.Status == 2 {
						//fmt.Println("inside")
						substatusname := "Unpublished"
						details.SubStatusName = substatusname
					}
				}
				details.DigitalRightsRegions = regionsResult
				details.CreatedBy = userDetails.UserName
				details.PublishingPlatforms = nil
				if details.SchedulingDateTime != nil {
					scheduletime := details.SchedulingDateTime
					details.SchedulingDateTime = scheduletime
				} else {
					details.SchedulingDateTime = nil
				}
				details.Products = nil
				details.SubscriptionPlans = nil
				details.IntroDuration = nil
				details.IntroStart = nil
				details.VarianceTrailers = nil
				varienceResponse = append(varienceResponse, details)
			}
			content.ContentVariances = varienceResponse
		} else {
			var seasonDetails []ContentSeasonsDetails
			if err := db.Debug().Table("content c").Select("s.content_id,s.season_key,s.modified_at,s.status,s.number as season_number,cpi.transliterated_title,atci.original_language,ct.language_type,ct.dubbing_language,ct.dubbing_dialect_id,cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date,json_agg(crc.country_id)::varchar as digital_rights,s.id").Joins("join season s on s.content_id = c.id join content_primary_info cpi on cpi.id =s.primary_info_id join about_the_content_info atci on atci.id =s.about_the_content_info_id join content_translation ct on ct.id = s.translation_id join content_rights cr on cr.id = s.rights_id join content_rights_country crc on crc.content_rights_id = cr.id").Where("c.id=? and s.deleted_by_user_id is null", content.Id).Order("s.number asc").Group("s.status,cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date,s.id,cpi.transliterated_title,s.content_id,s.season_key,atci.original_language,ct.language_type,ct.dubbing_language,ct.dubbing_dialect_id,c.modified_at").Find(&seasonDetails).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			seasonResponse := []ContentSeasonsDetails{}
			for _, sDetails := range seasonDetails {
				episodeResponse := []SeasonEpisodesDetails{}
				var episodeDetails []SeasonEpisodesDetails
				var userDetails UserDetails
				if err := udb.Debug().Table("user").Where("id=?", content.CreatedByUserId).Find(&userDetails).Error; err != nil && err.Error() != "record not found" {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
				sDetails.CreatedBy = &userDetails.UserName
				var primaryDetails PrimaryInfoDetails
				var aboutTheContent AboutTheContentDetails
				var translation ContentTranslationDetails
				var rights Rights
				if content.Status == 2 {
					if sDetails.Status == 2 {
						status := "Unpublished"
						sDetails.SubStatusName = &status
					} else if sDetails.Status == 1 {
						// status := "Published"
						// sDetails.SubStatusName = &status
						newdate := sDetails.DigitalRightsEndDate
						anotherdate := newdate.Format("2006-01-02")
						startsdate := sDetails.DigitalRightsStartDate
						newdates := startsdate.Format("2006-01-02")
						if anotherdate <= time.Now().Format("2006-01-02") {
							status := "Digital Rights Exceeded"
							sDetails.Status = 2
							sDetails.SubStatusName = &status
							sDetails.StatusCanBeChanged = false
						} else if newdates > time.Now().Format("2006-01-02") {
							status := "Scheduled"
							sDetails.Status = 2
							sDetails.SubStatusName = &status
							sDetails.StatusCanBeChanged = false
						} else {
							status := "Published"
							sDetails.Status = 1
							sDetails.SubStatusName = &status
							sDetails.StatusCanBeChanged = true
						}
					} else if sDetails.Status == 3 {
						status := "Draft"
						sDetails.SubStatusName = &status
					}
					//	status := "Unpublished"
					sDetails.Status = 2
					//	sDetails.SubStatusName = &status
					sDetails.StatusCanBeChanged = false
				} else if sDetails.Status == 1 {
					newdate := sDetails.DigitalRightsEndDate
					endDate := newdate.Format("2006-01-02")
					startsdate := sDetails.DigitalRightsStartDate
					startDate := startsdate.Format("2006-01-02")
					if endDate <= time.Now().Format("2006-01-02") {
						status := "Digital Rights Exceeded"
						sDetails.Status = 2
						sDetails.SubStatusName = &status
						sDetails.StatusCanBeChanged = false
					} else if startDate > time.Now().Format("2006-01-02") {
						status := "Scheduled"
						sDetails.Status = 2
						sDetails.SubStatusName = &status
						sDetails.StatusCanBeChanged = false
					} else {
						status := "Published"
						sDetails.Status = 1
						sDetails.SubStatusName = &status
						sDetails.StatusCanBeChanged = true
					}
				} else if sDetails.Status == 3 {
					status := "Draft"
					sDetails.SubStatusName = &status
					sDetails.Status = 2
					sDetails.StatusCanBeChanged = false
				} else if sDetails.Status == 2 {
					status := "Unpublished"
					sDetails.Status = 2
					sDetails.SubStatusName = &status
					sDetails.StatusCanBeChanged = true
				}
				primaryDetails.SeasonNumber = sDetails.SeasonNumber
				primaryDetails.OriginalTitle = nil
				primaryDetails.AlternativeTitle = nil
				primaryDetails.ArabicTitle = nil
				primaryDetails.TransliteratedTitle = sDetails.TransliteratedTitle
				primaryDetails.Notes = nil
				primaryDetails.IntroStart = nil
				primaryDetails.OutroStart = nil
				sDetails.PrimaryInfo = primaryDetails
				sDetails.Cast = nil
				sDetails.Music = nil
				sDetails.TagInfo = nil
				sDetails.SeasonGenres = nil
				aboutTheContent.OriginalLanguage = sDetails.OriginalLanguage
				aboutTheContent.Supplier = nil
				aboutTheContent.AcquisitionDepartment = nil
				aboutTheContent.ArabicSynopsis = nil
				aboutTheContent.ProductionYear = nil
				aboutTheContent.ProductionHouse = nil
				aboutTheContent.AgeGroup = nil
				aboutTheContent.IntroDuration = nil
				aboutTheContent.IntroStart = nil
				aboutTheContent.OutroDuration = nil
				aboutTheContent.OutroStart = nil
				aboutTheContent.ProductionCountries = nil
				sDetails.AboutTheContent = aboutTheContent
				translation.LanguageType = common.ContentLanguageOriginTypesName(sDetails.LanguageType)
				translation.DubbingLanguage = &sDetails.DubbingLanguage
				translation.DubbingDialectId = &sDetails.DubbingDialectId
				translation.SubtitlingLanguage = nil
				sDetails.Translation = translation
				if err := db.Debug().Table("episode e").Select("cr.digital_rights_type,e.status,cr.digital_rights_start_date,cr.digital_rights_end_date,e.number,cpi.transliterated_title,e.id,e.episode_key,pi2.scheduling_date_time").Joins("join content_primary_info cpi on cpi.id=e.primary_info_id join playback_item pi2 on pi2.id =e.playback_item_id join content_rights cr on cr.id =pi2.rights_id").Where("e.season_id=? and e.deleted_by_user_id is null", sDetails.Id).Order("e.number asc").Find(&episodeDetails).Error; err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
				// episodeResponse
				for _, eDetails := range episodeDetails {
					var primaryDetails PrimaryInfoDetails
					eDetails.IsPrimary = true
					eDetails.UserId = "00000000-0000-0000-0000-000000000000"
					eDetails.SecondarySeasonId = "00000000-0000-0000-0000-000000000000"
					eDetails.VarianceIds = nil
					eDetails.EpisodeIds = nil
					eDetails.SecondaryEpisodeId = "00000000-0000-0000-0000-000000000000"
					eDetails.ContentId = "00000000-0000-0000-0000-000000000000"
					eDetails.SeasonId = "00000000-0000-0000-0000-000000000000"
					eDetails.SubStatus = eDetails.Status
					if sDetails.Status == 2 {
						newdate := sDetails.DigitalRightsEndDate
						anotherdate := newdate.Format("2006-01-02")
						startsdate := sDetails.DigitalRightsStartDate
						newdates := startsdate.Format("2006-01-02")
						if anotherdate <= time.Now().Format("2006-01-02") && eDetails.Status == 1 {
							status := "Digital Rights Exceeded"
							eDetails.Status = 2
							eDetails.SubStatusName = &status
							eDetails.StatusCanBeChanged = false
						} else if newdates > time.Now().Format("2006-01-02") && eDetails.Status == 1 {
							status := "Scheduled"
							eDetails.Status = 2
							eDetails.SubStatusName = &status
							eDetails.StatusCanBeChanged = false
						} else if newdates > time.Now().Format("2006-01-02") && eDetails.Status == 1 {
							status := "Draft"
							eDetails.Status = 2
							eDetails.SubStatusName = &status
							eDetails.StatusCanBeChanged = false
						} else {
							if eDetails.Status == 2 {
								status := "Unpublished"
								eDetails.SubStatusName = &status
							}
							if eDetails.Status == 3 {
								status := "Draft"
								eDetails.SubStatusName = &status
							}
							eDetails.Status = 2
							eDetails.StatusCanBeChanged = false
						}
					} else if eDetails.Status == 1 {
						newdate := sDetails.DigitalRightsEndDate
						endDate := newdate.Format("2006-01-02")
						startsdate := sDetails.DigitalRightsStartDate
						startDate := startsdate.Format("2006-01-02")
						var scheduledate string
						if eDetails.SchedulingDateTime != nil {
							//scheduledate = eDetails.SchedulingDateTime.UTC().Format("2006-01-02 11:00:00")
							scheduledate = eDetails.SchedulingDateTime.String()
						}
						if endDate <= time.Now().Format("2006-01-02") {
							status := "Digital Rights Exceeded"
							eDetails.Status = 2
							eDetails.SubStatusName = &status
							eDetails.StatusCanBeChanged = false
						} else if startDate > time.Now().Format("2006-01-02") {
							status := "Scheduled"
							eDetails.Status = 2
							eDetails.SubStatusName = &status
							eDetails.StatusCanBeChanged = false
						} else if scheduledate > time.Now().String() {
							status := "Scheduled"
							eDetails.Status = 2
							eDetails.SubStatusName = &status
							eDetails.StatusCanBeChanged = false
						} else {
							status := "Published"
							eDetails.Status = 1
							eDetails.SubStatusName = &status
							eDetails.StatusCanBeChanged = true
						}
					} else if eDetails.Status == 3 {
						status := "Draft"
						eDetails.SubStatusName = &status
						eDetails.Status = 2
						eDetails.StatusCanBeChanged = false
					} else if eDetails.Status == 2 {
						status := "Unpublished"
						eDetails.Status = 2
						eDetails.SubStatusName = &status
						eDetails.StatusCanBeChanged = true
					}
					if eDetails.SchedulingDateTime != nil {
						if eDetails.SchedulingDateTime.Format("2006-01-02 11:33:41") > time.Now().Format("2006-01-02 11:33:41") && eDetails.Status == 1 {
							substatusname := "Scheduled"
							eDetails.SubStatusName = &substatusname
						} else if eDetails.SchedulingDateTime.Format("2006-01-02 11:33:41") < time.Now().Format("2006-01-02 11:33:41") && eDetails.Status == 1 {
							substatusname := "Published"
							eDetails.SubStatusName = &substatusname
						}
					}
					eDetails.CreatedBy = &userDetails.UserName
					primaryDetails.Number = eDetails.Number
					primaryDetails.VideoContentId = nil
					primaryDetails.SynopsisEnglish = nil
					primaryDetails.SynopsisArabic = nil
					primaryDetails.OriginalTitle = nil
					primaryDetails.AlternativeTitle = nil
					primaryDetails.ArabicTitle = nil
					primaryDetails.TransliteratedTitle = eDetails.TransliteratedTitle
					primaryDetails.Notes = nil
					primaryDetails.IntroStart = nil
					primaryDetails.OutroStart = nil
					eDetails.PrimaryInfo = primaryDetails
					eDetails.Cast = nil
					eDetails.Music = nil
					eDetails.TagInfo = nil
					eDetails.NonTextualData = nil
					if eDetails.SchedulingDateTime != nil {
						scheduledate := eDetails.SchedulingDateTime
						eDetails.SchedulingDateTime = scheduledate
					} else {
						eDetails.SchedulingDateTime = nil
					}
					eDetails.PublishingPlatforms = nil
					eDetails.SeoDetails = nil
					episodeResponse = append(episodeResponse, eDetails)
				}
				sDetails.Episodes = episodeResponse
				sDetails.NonTextualData = nil
				rights.DigitalRightsType = sDetails.DigitalRightsType
				rights.DigitalRightsStartDate = sDetails.DigitalRightsStartDate
				rights.DigitalRightsEndDate = sDetails.DigitalRightsEndDate
				regionsResult := make([]int, 0)
				if sDetails.DigitalRights != "" {
					regions, _ := JsonStringToIntSliceOrMap(sDetails.DigitalRights)
					regionsResult = RemoveDuplicateValues(regions)
					if regionsResult[0] == 0 {
						regionsResult = make([]int, 0)
					}
				}
				rights.DigitalRightsRegions = regionsResult
				rights.SubscriptionPlans = nil
				sDetails.Rights = rights
				sDetails.IntroDuration = "00:00:00"
				sDetails.IntroStart = "00:00:00"
				sDetails.OutroDuration = "00:00:00"
				sDetails.OutroStart = "00:00:00"
				sDetails.Products = nil
				sDetails.SeoDetails = nil
				sDetails.VarianceTrailers = nil
				seasonResponse = append(seasonResponse, sDetails)
			}
			content.ContentSeasons = seasonResponse
			// if len(seasonResponse) < 1 {
			// 	// buffer := make(map[string]string) // make(map[string]string[]interface{}
			// 	//	var new []buffer
			// 	var newarr []ContentSeasonsDetails
			// 	content.ContentSeasons = buffer
			// }
		}
		response = append(response, content)
	}
	type Response struct {
		Pagination Pagination             `json:"pagination"`
		Data       []GetAllContentDetails `json:"data"`
	}
	var pagination Pagination
	pagination.Size = totalCount
	pagination.Offset = int(offset)
	pagination.Limit = int(limit)
	var finalResult Response
	finalResult.Pagination = pagination
	finalResult.Data = response
	c.JSON(http.StatusOK, finalResult)
	return
}

func (hs *ContentService) UploadMenuPosterImage(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYOGUWMUMGAQLMW3U",                     // id
			"Jb0NV2eHwXAJg6UADb5vs3BgAyuUsvhgREi/hWRj", // secret
			""), // token can be left blank for now
	})
	if err != nil {
		fmt.Println("from s3 session", err)
		fmt.Println("Could not upload file -- session")
	}
	fileName, errr := UploadFileToS3(s, file, fileHeader, "poster-image_")
	if errr != nil {
		fmt.Println("from s3 upload", errr)
		fmt.Println("Could not upload file")
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) DetailsBackgroundImage(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYOGUWMUMGAQLMW3U",                     // id
			"Jb0NV2eHwXAJg6UADb5vs3BgAyuUsvhgREi/hWRj", // secret
			""), // token can be left blank for now
	})
	if err != nil {
		fmt.Println("Could not upload file")
	}
	fileName, err := UploadFileToS3(s, file, fileHeader, "details-background_")
	if err != nil {
		fmt.Println("Could not upload file")
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) MobileDetailsBackgroundImage(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYOGUWMUMGAQLMW3U",                     // id
			"Jb0NV2eHwXAJg6UADb5vs3BgAyuUsvhgREi/hWRj", // secret
			""), // token can be left blank for now
	})
	if err != nil {
		fmt.Println("Could not upload file")
	}
	fileName, err := UploadFileToS3(s, file, fileHeader, "mobile-details-background_")
	if err != nil {
		fmt.Println("Could not upload file")
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) OverlayPosterImage(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYOGUWMUMGAQLMW3U",                     // id
			"Jb0NV2eHwXAJg6UADb5vs3BgAyuUsvhgREi/hWRj", // secret
			""), // token can be left blank for now
	})
	if err != nil {
		fmt.Println("Could not upload file")
	}
	fileName, err := UploadFileToS3(s, file, fileHeader, "overlay-poster-image_")
	if err != nil {
		fmt.Println("Could not upload file")
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) MenuPosterImage(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYOGUWMUMGAQLMW3U",                     // id
			"Jb0NV2eHwXAJg6UADb5vs3BgAyuUsvhgREi/hWRj", // secret
			""), // token can be left blank for now
	})
	fileName, uploaderr := UploadFileToS3(s, file, fileHeader, "menu-poster-image_")
	if uploaderr != nil {
		fmt.Println("Could not upload file")
		return
	}
	fmt.Println("Image uploaded successfully", fileName)
	fmt.Println(fileName)
	filetrim := strings.Split(fileName, "/")
	fmt.Println(filetrim[1])
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) MobileMenu(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	// create an AWS session which can be
	// reused if we're uploading many files
	s, sessionerr := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYOGUWMUMGAQLMW3U",                     // id
			"Jb0NV2eHwXAJg6UADb5vs3BgAyuUsvhgREi/hWRj", // secret
			""), // token can be left blank for now
	})
	if sessionerr != nil {
		fmt.Println("Session Error", sessionerr)
		return
	}
	fileName, err := UploadFileToS3(s, file, fileHeader, "mobile-menu_")
	if err != nil {
		fmt.Println("Could not upload file", err)
		return
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) MobileMenuPosterImage(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYOGUWMUMGAQLMW3U",                     // id
			"Jb0NV2eHwXAJg6UADb5vs3BgAyuUsvhgREi/hWRj", // secret
			""), // token can be left blank for now
	})
	if err != nil {
		fmt.Println("Could not upload file")
	}
	fileName, err := UploadFileToS3(s, file, fileHeader, "mobile-menu-poster-image_")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) TrailerPosterImage(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYOGUWMUMGAQLMW3U",                     // id
			"Jb0NV2eHwXAJg6UADb5vs3BgAyuUsvhgREi/hWRj", // secret
			""), // token can be left blank for now
	})
	if err != nil {
		fmt.Println("Could not upload file")
	}
	fileName, err := UploadFileToS3(s, file, fileHeader, "trailer-poster-image_")
	if err != nil {
		fmt.Println("Could not upload file")
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

// UploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func UploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader, imagetype string) (string, error) {
	// get the file size and read
	// the file content into a buffer
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)
	tempFileName := "temp/" + imagetype + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)
	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("S3_BUCKET")),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		StorageClass:         aws.String("STANDARD"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		fmt.Printf("Unable to upload %q, %v", tempFileName, err)
	}
	fmt.Printf("Successfully uploaded %q", tempFileName)
	return tempFileName, err
}

/*Uploade image Based on season Id*/
func SeasonFileUPload(request CreateSeasonRequestValidation, seasonId string, contentId string) {
	bucketName := os.Getenv("S3_BUCKET")
	var newarr []string
	newarr = append(newarr, request.NonTextualData.PosterImage)
	newarr = append(newarr, request.NonTextualData.OverlayPosterImage)
	newarr = append(newarr, request.NonTextualData.DetailsBackground)
	newarr = append(newarr, request.NonTextualData.MobileDetailsBackground)
	for k := 0; k < len(newarr); k++ {
		item := newarr[k]
		filetrim := strings.Split(item, "_")
		Destination := contentId + "/" + seasonId + "/" + filetrim[0]
		source := bucketName + "/" + "temp/" + item
		s, err := session.NewSession(&aws.Config{
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
		errorr := SizeUploadFileToS3(s, filetrim[0], contentId, seasonId)
		if errorr != nil {
			fmt.Println("error in uploading size upload", errorr)
		}
		fmt.Println("Success!")
	}
}

// SizeUploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func SizeUploadFileToS3(s *session.Session, fileName string, contentId string, seasonId string) error {
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
		s3file := sizeValue[i] + contentId + "/" + seasonId + "/" + fileName
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
			fmt.Printf("Unable to upload %q, %v", fileName, er)
		}
		fmt.Printf("Successfully uploaded %q", fileName)
	}
	return er
}

/*Get Episode Details By Episode Id */
func (hs *ContentService) GetEpisodeDetailsBYepisodeId(c *gin.Context) {
	AuthorizationRequired := c.MustGet("AuthorizationRequired")
	userid := c.MustGet("userid")
	if AuthorizationRequired == 1 || userid == "" || c.MustGet("is_back_office_user") == "false" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	//serverError := common.ServerErrorResponse()
	notfound := common.NotFoundErrorResponse()
	db := c.MustGet("DB").(*gorm.DB)
	Id := c.Param("id")
	var episodeDetails SeasonEpisodesQueryDetails
	if err := db.Debug().Table("episode e").Select("e.episode_key,e.season_id,e.status,cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date,e.created_at,e.number,pi2.video_content_id,e.synopsis_english as english_synopsis,e.synopsis_arabic as arabic_synopsis,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title ,cpi.notes,cpi.intro_start,cpi.outro_start,ct.language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,pi2.scheduling_date_time,json_agg(pitp.target_platform)::varchar AS publishing_platforms,e.english_meta_title,e.arabic_meta_title,e.english_meta_description,e.arabic_meta_description,e.id,cc.main_actor_id,cc.main_actress_id,json_agg(ca.actor_id)::varchar as actors,json_agg(cw.writer_id)::varchar as writers,json_agg(cd.director_id)::varchar as directors,json_agg(cs.singer_id)::varchar as singers,json_agg(cmc.music_composer_id)::varchar as music_composers,json_agg(csw.song_writer_id)::varchar as song_writers,json_agg(ct1.textual_data_tag_id)::varchar as tags ,s.content_id,e.has_poster_image,e.has_dubbing_script,e.has_subtitling_script").Joins("join season s on s.id =e.season_id join content_primary_info cpi on cpi.id =e.primary_info_id join about_the_content_info atci on atci.id =s.about_the_content_info_id join content_translation ct on ct.id =s.translation_id join playback_item pi2 on pi2.id =e.playback_item_id join content_rights cr on cr.id =pi2.rights_id full outer join playback_item_target_platform pitp on pitp.playback_item_id = pi2.id full outer join content_cast cc on cc.id = e.cast_id full outer join content_actor ca on ca.cast_id =cc.id full outer join content_writer cw on cw.cast_id =cc.id full outer join content_director cd on cd.cast_id =cc.id full outer join content_singer cs on cs.music_id =e.music_id full outer join content_music_composer cmc on cmc.music_id =e.music_id full outer join content_song_writer csw on csw.music_id =e.music_id full outer join content_tag ct1 on ct1.tag_info_id =e.tag_info_id").Where("e.id =? and e.deleted_by_user_id is null", Id).Group("e.episode_key,e.season_id,e.status,cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date,e.created_at,e.number,pi2.video_content_id,atci.english_synopsis,atci.arabic_synopsis,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title,cpi.notes,cpi.intro_start,cpi.outro_start,ct.language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,pi2.scheduling_date_time,e.english_meta_title,e.arabic_meta_title,e.english_meta_description,e.arabic_meta_description,e.id,cc.main_actor_id,cc.main_actress_id,s.content_id").Find(&episodeDetails).Error; err != nil {
		c.JSON(http.StatusNotFound, notfound)
		return
	}
	var seasonEpisode SeasonEpisode
	var cast Cast
	var music Music
	var tagInfo TagInfo
	var contentTranslation ContentTranslation
	actorResult, writerResult, directorResult := make([]string, 0), make([]string, 0), make([]string, 0)
	singerResult, musicComposerResult, songWriterResult := make([]string, 0), make([]string, 0), make([]string, 0)
	seasonEpisode.IsPrimary = true
	seasonEpisode.UserId = "00000000-0000-0000-0000-000000000000"
	seasonEpisode.SecondarySeasonId = "00000000-0000-0000-0000-000000000000"
	seasonEpisode.VarianceIds = nil
	seasonEpisode.EpisodeIds = nil
	seasonEpisode.SecondaryEpisodeId = "00000000-0000-0000-0000-000000000000"
	seasonEpisode.ContentId = episodeDetails.ContentId
	seasonEpisode.EpisodeKey = episodeDetails.EpisodeKey
	seasonEpisode.SeasonId = episodeDetails.SeasonId
	seasonEpisode.Status = episodeDetails.Status
	now := time.Now().Unix()
	sDate := episodeDetails.DigitalRightsStartDate.Unix()
	eDate := episodeDetails.DigitalRightsEndDate.Unix()
	seasonEpisode.SubStatusName = nil
	seasonEpisode.Status = episodeDetails.Status
	seasonEpisode.SubStatus = episodeDetails.Status
	if episodeDetails.Status == 1 || episodeDetails.Status == 2 {
		seasonEpisode.StatusCanBeChanged = true
	} else if episodeDetails.Status == 3 {
		seasonEpisode.StatusCanBeChanged = false
		seasonEpisode.Status = 2
		seasonEpisode.SubStatus = episodeDetails.Status
		subStatus := "Draft"
		seasonEpisode.SubStatusName = &subStatus
	}
	if sDate > now || eDate < now {
		seasonEpisode.StatusCanBeChanged = false
		seasonEpisode.Status = 2
		seasonEpisode.SubStatus = episodeDetails.Status
		subStatus := "Digital Rights Exceeded"
		seasonEpisode.SubStatusName = &subStatus
	}
	seasonEpisode.DigitalRightsType = episodeDetails.DigitalRightsType
	seasonEpisode.DigitalRightsStartDate = episodeDetails.DigitalRightsStartDate
	seasonEpisode.DigitalRightsEndDate = episodeDetails.DigitalRightsEndDate
	seasonEpisode.CreatedBy = nil
	var primaryInfo PrimaryInfo
	primaryInfo.Number = episodeDetails.Number
	primaryInfo.VideoContentId = episodeDetails.VideoContentId
	primaryInfo.SynopsisEnglish = episodeDetails.EnglishSynopsis
	primaryInfo.SynopsisArabic = episodeDetails.ArabicSynopsis
	primaryInfo.OriginalTitle = episodeDetails.OriginalTitle
	primaryInfo.AlternativeTitle = episodeDetails.AlternativeTitle
	primaryInfo.ArabicTitle = episodeDetails.ArabicTitle
	primaryInfo.TransliteratedTitle = episodeDetails.TransliteratedTitle
	primaryInfo.Notes = episodeDetails.Notes
	primaryInfo.IntroStart = episodeDetails.IntroStart
	primaryInfo.OutroStart = episodeDetails.OutroStart
	seasonEpisode.PrimaryInfo = primaryInfo
	cast.MainActorId = episodeDetails.MainActorId
	cast.MainActressId = episodeDetails.MainActressId
	if episodeDetails.Actors != "" {
		actors, _ := JsonStringToStringSliceOrMap(episodeDetails.Actors)
		actorResult = RemoveDuplicateStringValues(actors)
		if actorResult[0] == "" {
			actorResult = make([]string, 0)
		}
	}
	cast.Actors = actorResult
	if episodeDetails.Writers != "" {
		writers, _ := JsonStringToStringSliceOrMap(episodeDetails.Writers)
		writerResult = RemoveDuplicateStringValues(writers)
		if writerResult[0] == "" {
			writerResult = make([]string, 0)
		}
	}
	cast.Writers = writerResult
	if episodeDetails.Directors != "" {
		directors, _ := JsonStringToStringSliceOrMap(episodeDetails.Directors)
		directorResult = RemoveDuplicateStringValues(directors)
		if directorResult[0] == "" {
			directorResult = make([]string, 0)
		}
	}
	cast.Directors = directorResult
	seasonEpisode.Cast = cast
	if episodeDetails.Singers != "" {
		singers, _ := JsonStringToStringSliceOrMap(episodeDetails.Singers)
		singerResult = RemoveDuplicateStringValues(singers)
		if singerResult[0] == "" {
			singerResult = make([]string, 0)
		}
	}
	music.Singers = singerResult
	if episodeDetails.MusicComposers != "" {
		musicComposers, _ := JsonStringToStringSliceOrMap(episodeDetails.MusicComposers)
		musicComposerResult = RemoveDuplicateStringValues(musicComposers)
		if musicComposerResult[0] == "" {
			musicComposerResult = make([]string, 0)
		}
	}
	music.MusicComposers = musicComposerResult
	if episodeDetails.SongWriters != "" {
		songWriters, _ := JsonStringToStringSliceOrMap(episodeDetails.SongWriters)
		songWriterResult = RemoveDuplicateStringValues(songWriters)
		if songWriterResult[0] == "" {
			songWriterResult = make([]string, 0)
		}
	}
	music.SongWriters = songWriterResult
	seasonEpisode.Music = music
	tagsResult := make([]string, 0)
	if episodeDetails.Tags != "" {
		tags, _ := JsonStringToStringSliceOrMap(episodeDetails.Tags)
		tagsResult = RemoveDuplicateStringValues(tags)
		if tagsResult[0] == "" {
			tagsResult = make([]string, 0)
		}
	}
	tagInfo.Tags = tagsResult
	seasonEpisode.TagInfo = tagInfo
	if episodeDetails.HasPosterImage {
		seasonEpisode.NonTextualData.PosterImage = os.Getenv("IMAGERY_URL") + episodeDetails.ContentId + "/" + episodeDetails.SeasonId + "/" + episodeDetails.Id + "/poster-image"
	}
	if episodeDetails.HasDubbingScript {
		seasonEpisode.NonTextualData.DubbingScript = os.Getenv("IMAGERY_URL") + episodeDetails.ContentId + "/" + episodeDetails.SeasonId + "/" + episodeDetails.Id + "/dubbing-script"
	}
	if episodeDetails.HasSubtitlingScript {
		seasonEpisode.NonTextualData.SubtitlingScript = os.Getenv("IMAGERY_URL") + episodeDetails.ContentId + "/" + episodeDetails.SeasonId + "/" + episodeDetails.Id + "/subtitling-script"
	}
	contentTranslation.LanguageType = common.ContentLanguageOriginTypesName(episodeDetails.LanguageType)
	contentTranslation.DubbingLanguage = episodeDetails.DubbingLanguage
	contentTranslation.DubbingDialectId = episodeDetails.DubbingDialectId
	contentTranslation.SubtitlingLanguage = episodeDetails.SubtitlingLanguage
	seasonEpisode.Translation = contentTranslation
	seasonEpisode.SchedulingDateTime = episodeDetails.SchedulingDateTime
	platformResult := make([]int, 0)
	if episodeDetails.PublishingPlatforms != "" {
		platforms, _ := JsonStringToIntSliceOrMap(episodeDetails.PublishingPlatforms)
		platformResult = RemoveDuplicateValues(platforms)
	}
	seasonEpisode.PublishingPlatforms = platformResult
	seasonEpisode.SeoDetails.ArabicMetaTitle = episodeDetails.ArabicMetaTitle
	seasonEpisode.SeoDetails.EnglishMetaTitle = episodeDetails.EnglishMetaTitle
	seasonEpisode.SeoDetails.ArabicMetaDescription = episodeDetails.ArabicMetaDescription
	seasonEpisode.SeoDetails.EnglishMetaDescription = episodeDetails.EnglishMetaDescription
	seasonEpisode.Id = episodeDetails.Id
	c.JSON(http.StatusOK, gin.H{"data": seasonEpisode})
	return
}

/*Uploade image in S3 bucket  Based on variance and trailer Id*/
func SeasonVarianceTrailerImageUpload(oldcontentId string, oldseasonid string, oldTrailerId string, contentId string, seasonId string, TrailerId string, trailerPosterImage string) {
	bucketName := os.Getenv("S3_BUCKET")
	var newarr []string

	// for _, value := range Variances.VarianceTrailer {
	newarr = append(newarr, trailerPosterImage)
	for k := 0; k < len(newarr); k++ {
		item := newarr[k]
		if strings.Contains(item, "_") {
			filetrim := strings.Split(item, "_")
			Destination := contentId + "/" + seasonId + "/" + TrailerId + "/" + filetrim[0]
			source := bucketName + "/" + "temp/" + item
			s, err := session.NewSession(&aws.Config{
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
			fmt.Println(result, "reseult......")
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
			errorr := SizeUploadFileToS3(s, filetrim[0], contentId, seasonId)
			if errorr != nil {
				fmt.Println("error in uploading size upload", errorr)
			}
			fmt.Println("Success!")
		} else {
			Destination := contentId + "/" + seasonId + "/" + TrailerId + "/" + "trailer-poster-image"
			source := bucketName + "/" + oldcontentId + "/" + oldseasonid + "/" + oldTrailerId + "/" + "trailer-poster-image"
			s, err := session.NewSession(&aws.Config{
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
			fmt.Println(result, "reseult......")
		}
	}
	// }

}

/*Uploade image in S3 bucket  Based on variance and trailer Id*/
func SeasonTrailerImageUpload(contentId string, seasonId string, TrailerId string, trailerPosterImage string) {
	bucketName := os.Getenv("S3_BUCKET")
	var newarr []string

	// for _, value := range Variances.VarianceTrailer {
	newarr = append(newarr, trailerPosterImage)
	for k := 0; k < len(newarr); k++ {
		item := newarr[k]
		if strings.Contains(item, "_") {
			filetrim := strings.Split(item, "_")
			Destination := contentId + "/" + seasonId + "/" + TrailerId + "/" + filetrim[0]
			source := bucketName + "/" + "temp/" + item
			s, err := session.NewSession(&aws.Config{
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
			fmt.Println(result, "reseult......")
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
			errorr := SizeUploadFileToS3(s, filetrim[0], contentId, seasonId)
			if errorr != nil {
				fmt.Println("error in uploading size upload", errorr)
			}
			fmt.Println("Success!")
		}
	}
	// }

}

type PageSync struct {
	PageId     string `json:"pageId"`
	DirtyCount int    `json:"dirtyCount,omitempty"`
	PageKey    int    `json:"pageKey"`
}

type ContentIds struct {
	ContentId string
}
type ContentSync struct {
	ContentId  string `json:"contentId"`
	DirtyCount int    `json:"dirtyCount"`
}
type Contents struct {
	Id          string `json:"Id"`
	ContentTier int    `json:"contentTier"`
	ContentKey  int    `json:"contentKey"`
}

/* sync tables updation */
func (hs *ContentService) PageSyncWithContentId(c *gin.Context) {
	fcdb := c.MustGet("FCDB").(*gorm.DB)
	fdb := c.MustGet("FDB").(*gorm.DB)
	db := c.MustGet("DB").(*gorm.DB)
	fmt.Println("AAAAAAAAAAAAAAAAA")
	var pageSync []PageSync
	var contentIds []ContentIds
	/* checking for scheduled contents */
	var contentdata []Contents
	db.Debug().Raw("select c.id,content_tier,c.content_key  from content c join content_variance cv on cv.content_id = c.id join playback_item pi2 on pi2.id  = cv.playback_item_id where pi2.scheduling_date_time between timezone('UTC', now()) - interval '5 minute' and timezone('UTC', now()) union select c.id,content_tier,c.content_key from content c join season s on s.content_id = c.id join episode e on e.season_id = s.id join playback_item pi2 on pi2.id =e.playback_item_id where pi2.scheduling_date_time between timezone('UTC', now()) - interval '5 minute' and timezone('UTC', now())").Find(&contentdata)
	for _, val := range contentdata {
		if val.ContentTier == 1 {
			go common.CreateRedisKeyForContentTypeOTC(c)
		} else if val.ContentTier == 2 {
			go common.CreateRedisKeyForContentTypeMTC(c)
		}
		go common.ContentSynching(val.Id, c)
		contentkeyconverted := strconv.Itoa(val.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
	}
	if fetchResult := fdb.Debug().Table("content_sync").Select(" distinct content_id").Where("dirty_count > ?", 0).Find(&contentIds).Error; fetchResult != nil {
		fmt.Println(fetchResult)
		return
	}
	fmt.Println(contentIds)
	for _, ids := range contentIds {
		/*Fetch page details with content Id through playlist and slider */
		contentcount := fcdb.Debug().Table("playlist_item pi2").Select("p2.page_key,p2.id as page_id").Joins("join playlist p on p.id = pi2.playlist_id").Joins("join page_playlist pp on pp.playlist_id = p.id").Joins("join page p2 on p2.id = pp.page_id").Joins("join slider s on s.black_area_playlist_id = p.id or s.red_area_playlist_id =p.id or s.green_area_playlist_id =p.id").Where("pi2.one_tier_content_id = ? or pi2.multi_tier_content_id = ?", ids.ContentId, ids.ContentId).Find(&pageSync)
		totalcount := int(contentcount.RowsAffected)
		if totalcount < 1 {
			/*Fetch page details with content Id through playlist */
			if playError := fcdb.Debug().Table("playlist_item pi2").Select("p2.page_key,p2.id as page_id").
				Joins("join playlist p on p.id = pi2.playlist_id").
				Joins("join page_playlist pp on pp.playlist_id = p.id").
				Joins("join page p2 on p2.id = pp.page_id").
				Where("pi2.one_tier_content_id =? or pi2.multi_tier_content_id =?", ids.ContentId, ids.ContentId).
				Find(&pageSync).Error; playError != nil {
				fmt.Println(playError)
				return
			}
		}
		for _, pkeys := range pageSync {
			var pageSyncDetails PageSync
			rows := fdb.Debug().Table("page_sync ps").Select("ps.page_id,ps.page_key,ps.dirty_count").Where("ps.page_id=?", pkeys.PageId).Find(&pageSyncDetails)
			totalcount := int(rows.RowsAffected)
			var insertPageSync PageSync
			if totalcount > 0 {
				insertPageSync.PageId = pageSyncDetails.PageId
				insertPageSync.DirtyCount = pageSyncDetails.DirtyCount + 1
				insertPageSync.PageKey = pageSyncDetails.PageKey
				if updateError := fdb.Debug().Table("page_sync").Where("page_id=?", pkeys.PageId).Update(&insertPageSync).Error; updateError != nil {
					fmt.Println(updateError)
					return
				}
			} else {
				insertPageSync.PageId = pkeys.PageId
				insertPageSync.DirtyCount = 1
				insertPageSync.PageKey = pkeys.PageKey
				if updateError := fdb.Debug().Table("page_sync").Create(&insertPageSync).Error; updateError != nil {
					fmt.Println(updateError)
					return
				}
			}
		}
		var content Content
		contentrow := db.Debug().Table("content").Where("id=? and status = 1 ", ids.ContentId).Find(&content)
		if contentrow.RowsAffected == 0 {
			var deletecontentSync ContentSync
			if deleteError := fdb.Debug().Table("content_fragment").Where("content_id=?", ids.ContentId).Delete(&deletecontentSync).Error; deleteError != nil {
				fmt.Println(deleteError)
				return
			}
		}
		var deletecontentSyncFC ContentSync
		if deleteContentError := fdb.Debug().Table("content_sync").Where("content_id=?", ids.ContentId).Delete(&deletecontentSyncFC).Error; deleteContentError != nil {
			fmt.Println(deleteContentError)
			return
		}
	}
}

/*Uploade image in S3 bucket  Based on variance and trailer Id*/
func SeasonVarianceEpisodeImageUpload(contentId string, seasonId string, EpisodeId string, nontextualImages []Images, OldcontentId string, OldseasonId string, OldEpisodeId string) {
	bucketName := os.Getenv("S3_BUCKET")
	for _, nontextrange := range nontextualImages {
		if nontextrange.HasImage {
			source := bucketName + "/" + OldcontentId + "/" + OldseasonId + "/" + OldEpisodeId + "/" + nontextrange.Imagename
			Destination := contentId + "/" + seasonId + "/" + EpisodeId + "/" + nontextrange.Imagename

			s, err := session.NewSession(&aws.Config{
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
			fmt.Println(result, "reseult......")
			url := "https://" + bucketName + ".s3.ap-south-1.amazonaws.com/" + Destination
			// don't worry about errors
			response, e := http.Get(url)
			if e != nil {
				log.Fatal(e)
			}
			defer response.Body.Close()

			//open a file for writing
			file, err := os.Create(nontextrange.Imagename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			// Use io.Copy to just dump the response body to the file. This supports huge files
			_, err = io.Copy(file, response.Body)
			if err != nil {
				log.Fatal(err)
			}
			errorr := episode.SizeUploadFileToS3(s, nontextrange.Imagename, contentId, seasonId, EpisodeId)
			if errorr != nil {
				fmt.Println("error in uploading size upload", errorr)
			}
			fmt.Println("Success!")
		}
	}

}

/*Uploade image Based on season Id*/
func SeasonVarianceFileUPload(oldseasonid string, seasonId string, contentId string) {
	bucketName := os.Getenv("S3_BUCKET")
	var newarr []string
	newarr = append(newarr, "poster-image")
	newarr = append(newarr, "overlay-poster-image")
	newarr = append(newarr, "details-background")
	newarr = append(newarr, "mobile-details-background")
	for k := 0; k < len(newarr); k++ {
		item := newarr[k]
		//	filetrim := strings.Split(item, "_")
		Destination := contentId + "/" + seasonId + "/" + item
		source := bucketName + "/" + contentId + "/" + oldseasonid + "/" + item
		s, err := session.NewSession(&aws.Config{
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
		file, err := os.Create(item)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal(err)
		}
		errorr := SizeUploadFileToS3(s, item, contentId, seasonId)
		if errorr != nil {
			fmt.Println("error in uploading size upload", errorr)
		}
		fmt.Println("Success!")
	}
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
		"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC1nukPNB8e8iXM\nMhHr1iuvSGwienHPd4avprk0yUIXAlBZzwvbcK9i8V5yfpZzz6RUQcwshPOs5k9r\n3hBMy7zTGiZeyh1tPHdSumn3c4o7vL90RivGKff0VFvbPk4GdcuUFrEOJEH5gMS3\nYBTJtKKhOxKK3hqG/e0WQVjPJidfZwIKFthq9+z+d/4GMpldJAI3CPRpw9l+xzoC\n+vOueZ0aUCaSMvLKqgVzsKp5+YxGAEZbdxeYPMZGffJlZVedBwFNnyELBL8uKmVi\nmnCABwMjCTRwL3bPSgJ9mHLm2FiIK3heJ6Tg5HFjjIHIrxcdbVG57lKoXOt2wKed\n23l24T9tAgMBAAECggEAAvM8+unWbG6qjzmvLPtn1kzLpXEoEEEd8ssxMqJIqCOM\nLHCGOubJnZXZ4evNMbH3BcjHirUcWTvluUW2Rh4GiA/KIdEKIdoXL1bzORTMvG7d\nhoOI/69agNtAgwIp6ZTO+K24QODQnBrNPtccJ7cXaamqoFI4XgsHc7Q2jsfNC2bp\nAaIi4ZZHLhhQf94KFmfqOsVMhX7nmaBjaVZrpIfSM+5g0ESKbYbaLgdg/yVwpbdQ\nrLjolOvZw5r8e8ZVstdfU/GwihHuNsbgTbU511IeUYd+YmxoCZ1fkJq0Xf0uJ6Cy\nz1byOXfFOps8RurZhR1hkUknfeBaTBGVrujlHcoV6wKBgQDgBPKgOlzEWh/I5Epz\nviJZKa3TTJS2kinIGDYjbiJNp5NucQZfkJu0xBn0vztwmFnIdYnIOV7kiWuWulVM\nzjC3KoSBiC+GVGjABukuU/dlcWpbSttRuKg90gJ/gOtF4FuYLWhZkJGM2iCndanv\nkFmylCMoq6aiPnC73VGX2mfvpwKBgQDPjHTAzU1RMWiymF9yO25RZ6jzcyB6hPXf\n2NG2YJ0luM41pZMx5DlRFi8ky6YGK2gFwNsyBBXhRh5AdciGD96FwbqtLEv6eKCn\nC1BxZceYAdA/P6Aa6h/4Wv7J67THKEYbwzGgPYE5Js+jmJXG3tKhtPqnlWa7Zx6H\nLS5uNM7aywKBgF4z1m9oe3AaUflhfqlzT/BcpXsQXgz0I9u/yqxVeNlc2ZN8tehj\n4AZA3IVeETnE5yRzwM/QyEWkP/jvPEWDA1tS5sutoAaF4lK11UKlDoi7C7V+IgIY\ne68ba++AH++PbBTvK01WjM5FP6wLv70832tH/gzxOa5KQY/OfqwzrLdLAoGBAMoa\nAoLQJ+7dRw9KEv9AYf9BCqLtw32qxWYRUrzePYhC+gIBVmEp1KpiCMwyxluRnvyj\nPI7qrYes6L5qMzZgc5YZ/LauwNmI5x9ihBW4P3CEq407XqN2wmTr7tke/e1FCWf1\nXfiki5XkdiLe7VI3HjI68i2H7P6lvnNxCppkL92bAoGAERQ76Setyiz6VOiy+AJ6\nRVcQg2PDpoP0woWj7AOstCwbP91AT1h0Wq/aXRm1lk10Yvq0zm8RNzAUkPrfqLwx\n1pC5SQPrut0h+RZaQPRzUJERxfzMzej/WStGh51E9gRyFdSKfL/iOQ6ZT8PysEwy\ntcDR1sRoS24TmmtlgwyP91o=\n-----END PRIVATE KEY-----\n",
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

func (hs *ContentService) UploadMenuPosterImageGcp(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()

	fileName, err := gcpUpload(file, fileHeader, "poster-image_")
	if err != nil {
		fmt.Println("from gcp upload", err)
		fmt.Println("Could not upload file")
		c.JSON(http.StatusBadGateway, gin.H{"data": "Could not upload file"})
		return
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) DetailsBackgroundImageGcp(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()

	fileName, err := gcpUpload(file, fileHeader, "details-background_")
	if err != nil {
		fmt.Println("from gcp upload", err)
		fmt.Println("Could not upload file")
		c.JSON(http.StatusBadGateway, gin.H{"data": "Could not upload file"})
		return
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")

	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) MobileDetailsBackgroundImageGcp(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()

	fileName, err := gcpUpload(file, fileHeader, "mobile-details-background_")
	if err != nil {
		fmt.Println("from gcp upload", err)
		fmt.Println("Could not upload file")
		c.JSON(http.StatusBadGateway, gin.H{"data": "Could not upload file"})
		return
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) OverlayPosterImageGcp(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()

	fileName, err := gcpUpload(file, fileHeader, "overlay-poster-image_")
	if err != nil {
		fmt.Println("from gcp upload", err)
		fmt.Println("Could not upload file")
		c.JSON(http.StatusBadGateway, gin.H{"data": "Could not upload file"})
		return
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) MenuPosterImageGcp(c *gin.Context) {
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	fileName, err := gcpUpload(file, fileHeader, "menu-poster-image_")
	if err != nil {
		fmt.Println("from gcp upload", err)
		fmt.Println("Could not upload file")
		c.JSON(http.StatusBadGateway, gin.H{"data": "Could not upload file"})
		return
	}
	fmt.Println("Image uploaded successfully", fileName)
	fmt.Println(fileName)
	filetrim := strings.Split(fileName, "/")
	fmt.Println(filetrim[1])
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) MobileMenuGcp(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	fileName, err := gcpUpload(file, fileHeader, "mobile-menu_")
	if err != nil {
		fmt.Println("from gcp upload", err)
		fmt.Println("Could not upload file")
		c.JSON(http.StatusBadGateway, gin.H{"data": "Could not upload file"})
		return
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) MobileMenuPosterImageGcp(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	fileName, err := gcpUpload(file, fileHeader, "mobile-menu-poster-image_")
	if err != nil {
		fmt.Println("from gcp upload", err)
		fmt.Println("Could not upload file")
		c.JSON(http.StatusBadGateway, gin.H{"data": "Could not upload file"})
		return
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func (hs *ContentService) TrailerPosterImageGcp(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	fileName, err := gcpUpload(file, fileHeader, "trailer-poster-image_")
	if err != nil {
		fmt.Println("from gcp upload", err)
		fmt.Println("Could not upload file")
		c.JSON(http.StatusBadGateway, gin.H{"data": "Could not upload file"})
		return
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	c.JSON(http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

func gcpUpload(file multipart.File, fileHeader *multipart.FileHeader, imagetype string) (string, error) {
	ctx := context.Background()
	fileContent, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	tempFileName := "temp/" + imagetype + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)
	contentType := mime.TypeByExtension(filepath.Ext(tempFileName))

	client, gcperr := getGCPClient()
	if gcperr != nil {
		fmt.Println("from gcp Connection", err)
		return "err", err
	}
	bucketName := os.Getenv("BUCKET_NAME")

	objectName := tempFileName
	obj := client.Bucket(bucketName).Object(objectName)
	wc := obj.NewWriter(ctx)

	_, err = io.Copy(wc, fileContent)
	if err != nil {
		fmt.Println("Error copying file content:", err)
		return "", err
	}
	wc.ContentType = contentType
	if err := wc.Close(); err != nil {
		return "", err
	}
	fileLocation := tempFileName
	// var wr io.Writer
	// signedObjectURL, err := generateV4GetObjectSignedURL(wr, bucketName, objectName)
	// _ = fileLocation

	return fileLocation, nil

}

// generateV4GetObjectSignedURL generates object signed URL with GET method.
func generateV4GetObjectSignedURL(w io.Writer, bucket, object string) (string, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	fmt.Println("google cloud stroage")
	// ctx := context.Background()
	// client, err := storage.NewClient(ctx)
	// if err != nil {
	// 	fmt.Println("ERROR===>",err)
	// 		return "", fmt.Errorf("storage.NewClient: %w", err)
	// }
	// defer client.Close()
	client, err := getGCPClient()
	if err != nil {
		fmt.Println("from gcp Connection", err)
		return "err", err
	}
	// Signing a URL requires credentials autos.Getenv("CLIENT_EMAIL")horized to sign a URL. You can pass
	// these in through SignedURLOptions with one of the following options:
	//    a. a Google service account private key, obtainable from the Google Developers Console
	//    b. a Google Access ID with iam.serviceAccounts.signBlob permissions
	//    c. a SignBytes function implementing custom signing.
	// In this example, none of these options are used, which means the SignedURL
	// function attempts to use the same authentication that was used to instantiate
	// the Storage client. This authentication must include a private key or have
	// iam.serviceAccounts.signBlob permissions.

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute),
		GoogleAccessID: os.Getenv("CLIENT_EMAIL"),
		PrivateKey:     []byte(os.Getenv("PRIVATE_KEY")),
	}

	u, err := client.Bucket(bucket).SignedURL(object, opts)
	fmt.Println("ERROR ===>", u, err)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %w", bucket, err)
	}

	fmt.Fprintln(w, "Generated GET signed URL:")
	fmt.Fprintf(w, "%q\n", u)
	fmt.Fprintln(w, "You can use this URL with any user agent, for example:")
	fmt.Fprintf(w, "curl %q\n", u)
	return u, nil
}

func SizeUploadFileToGCP(client *storage.Client, fileName string, contentId string, seasonId string, fileUrl string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	sizeValue := []string{
		"100x100/",
		"150x150/",
		"200x200/",
		"250x250/",
		"270x270/",
		"300x300/",
		"420x420/",
		"450x450/",
		"570x570/",
		"400x200/",
		"600x600/",
		"620x620/",
		"800x384/",
		"800x800/",
		"811x811/",
		"900x900/",
		"1125x240/",
		"1920x1080/",
		"1125x540/",
	}
	ctx := context.Background()
	var er error
	for i := 0; i < len(sizeValue); i++ {
		filetrim := strings.Split(sizeValue[i], "/")
		filetri := strings.Split(filetrim[0], "x")
		width := filetri[0]
		height := filetri[1]
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
			fmt.Println("err1" , err )
		}
		defer file.Close()

		// Get file size and read content into a buffer
		fileInfo, _ := file.Stat()
		var size = fileInfo.Size()
		buffer := make([]byte, size)
		file.Read(buffer)

		gcsFile := sizeValue[i] + contentId + "/" + seasonId + "/" + fileName
		wc := client.Bucket(os.Getenv("BUCKET_NAME")).Object(gcsFile).NewWriter(ctx)
		wc.ContentType = http.DetectContentType(buffer)
		wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
		if _, err := wc.Write(buffer); err != nil {
			return fmt.Errorf("unable to upload %s: %v", fileName, err)
		}
		if err := wc.Close(); err != nil {
			return fmt.Errorf("unable to close writer for %s: %v", fileName, err)
		}

		fmt.Printf("Successfully uploaded %q\n", fileName)
	}

	return er
}

func SizeUploadFileToGcpEpi(ctx context.Context, client *storage.Client, fileName string, episodeId string, contentId string, seasonId string, fileUrl string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	sizeValue := []string{
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
		filetrim := strings.Split(sizeValue[i], "/")
		filetri := strings.Split(filetrim[0], "x")
		width := filetri[0]
		height := filetri[1]
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
			fmt.Println(err)
		}
		defer file.Close()

		// Get file size and read content into a buffer
		fileInfo, _ := file.Stat()
		var size = fileInfo.Size()
		buffer := make([]byte, size)
		file.Read(buffer)

		gcsFile := sizeValue[i] + contentId + "/" + seasonId + "/" + episodeId + "/" + fileName
		wc := client.Bucket(os.Getenv("BUCKET_NAME")).Object(gcsFile).NewWriter(ctx)
		wc.ContentType = http.DetectContentType(buffer)
		wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
		if _, err := wc.Write(buffer); err != nil {
			return fmt.Errorf("unable to upload %s: %v", fileName, err)
		}
		if err := wc.Close(); err != nil {
			return fmt.Errorf("unable to close writer for %s: %v", fileName, err)
		}

		fmt.Printf("Successfully uploaded %q\n", fileName)
	}

	return er
}

func SeasonTrailerImageUploadGcp(contentId string, seasonId string, TrailerId string, trailerPosterImage string) {
	bucketName := os.Getenv("BUCKET_NAME")
	var newarr []string

	newarr = append(newarr, trailerPosterImage)
	for k := 0; k < len(newarr); k++ {
		item := newarr[k]
		if strings.Contains(item, "_") {
			filetrim := strings.Split(item, "_")
			Destination := contentId + "/" + seasonId + "/" + TrailerId + "/" + filetrim[0]
			source := "temp/" + item

			ctx := context.Background()
			client, gcperr := getGCPClient()
			if gcperr != nil {
				fmt.Println("from gcp Connection", gcperr)

			}

			// Copy the object from one directory to another in GCS
			srcObject := client.Bucket(bucketName).Object(source)
			attrs, err := srcObject.Attrs(ctx)
			_ = attrs
			if err != nil {
				// Handle the case where the source object doesn't exist
				fmt.Printf("Source object does not exist: %v\n", err)

				// Modify the source path if needed
				filetrims := strings.Split(item, "/")
				source = contentId + "/" + seasonId + "/" + TrailerId + "/" + filetrims[len(filetrims)-1]
				Destination = contentId + "/" + seasonId + "/" + TrailerId + "/" + filetrims[len(filetrims)-1]
				srcObject = client.Bucket(bucketName).Object(source)
				filetrim[0] = filetrims[len(filetrims)-1]

				// Retry the Attrs call
				attrs, err = srcObject.Attrs(ctx)
				if err != nil {
					// Handle the case where the modified source object also doesn't exist
					fmt.Printf("Modified source object does not exist: %v\n", err)
					continue
				}
			}
			dstObject := client.Bucket(bucketName).Object(Destination)
			if _, err := dstObject.CopierFrom(srcObject).Run(ctx); err != nil {
				fmt.Printf("Error copying object: %v", err)
				continue
			}

			url := "https://storage.googleapis.com/" + bucketName + "/" + Destination
			// Don't worry about errors
			response, e := http.Get(url)
			if e != nil {
				fmt.Println(e)
			}
			defer response.Body.Close()

			// Open a file for writing
			file, err := os.Create(filetrim[0])
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()

			// Use io.Copy to just dump the response body to the file
			_, err = io.Copy(file, response.Body)
			if err != nil {
				fmt.Println(err)
			}

			errorr := SizeUploadFileToGCP(client, filetrim[0], contentId, seasonId, url)
			if errorr != nil {
				fmt.Println("error in uploading size upload", errorr)
			}
			os.Remove(filetrim[0])
			fmt.Println("Success!")
		}
	}
	// }
}

func SeasonFileUPloadGcp(request CreateSeasonRequestValidation, seasonId string, contentId string) {
	bucketName := os.Getenv("BUCKET_NAME")
	var newarr []string
	newarr = append(newarr, request.NonTextualData.PosterImage)
	newarr = append(newarr, request.NonTextualData.OverlayPosterImage)
	newarr = append(newarr, request.NonTextualData.DetailsBackground)
	newarr = append(newarr, request.NonTextualData.MobileDetailsBackground)
	// newarr = append(newarr, request.NonTextualData.SeasonLogo)

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
		Destination := contentId + "/" + seasonId + "/" + filetrim[0]
		source := "temp/" + item

		// Copy the object from one directory to another in GCS
		srcObject := client.Bucket(bucketName).Object(source)
		attrs, err := srcObject.Attrs(ctx)
		_ = attrs
		if err != nil {
			// Handle the case where the source object doesn't exist
			fmt.Printf("Source object does not exist: %v\n", err)

			// Modify the source path if needed
			filetrims := strings.Split(item, "/")
			source = contentId + "/" + seasonId + "/" + filetrims[len(filetrims)-1]
			Destination = contentId + "/" + seasonId + "/" + filetrims[len(filetrims)-1]
			srcObject = client.Bucket(bucketName).Object(source)
			filetrim[0] = filetrims[len(filetrims)-1]

			// Retry the Attrs call
			attrs, err = srcObject.Attrs(ctx)
			if err != nil {
				// Handle the case where the modified source object also doesn't exist
				fmt.Printf("Modified source object does not exist: %v\n", err)
				continue
			}
		}
		dstObject := client.Bucket(bucketName).Object(Destination)
		if _, err := dstObject.CopierFrom(srcObject).Run(ctx); err != nil {
			fmt.Printf("Error copying object: %v", err)
			continue
		}

		url := "https://storage.googleapis.com/" + bucketName + "/" + Destination
		// Don't worry about errors
		response, e := http.Get(url)
		if e != nil {
			fmt.Println(e)
		}
		defer response.Body.Close()

		// Open a file for writing
		file, err := os.Create(filetrim[0])
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		// Use io.Copy to just dump the response body to the file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Println(err)
		}

		errorr := SizeUploadFileToGCP(client, filetrim[0], contentId, seasonId, url)
		if errorr != nil {
			fmt.Println("error in uploading size upload", errorr)
		}
		os.Remove(filetrim[0])
		fmt.Println("Success!")
	}
}

func SeasonVarianceEpisodeImageUploadGcp(contentId string, seasonId string, EpisodeId string, nontextualImages []Images, OldcontentId string, OldseasonId string, OldEpisodeId string) {
	bucketName := os.Getenv("BUCKET_NAME")
	ctx := context.Background()
	client, gcperr := getGCPClient()
	if gcperr != nil {
		fmt.Println("from gcp Connection", gcperr)
	}
	defer client.Close()

	for _, nontextrange := range nontextualImages {
		if nontextrange.HasImage {
			source := "temp/" + OldcontentId + "/" + OldseasonId + "/" + OldEpisodeId + "/" + nontextrange.Imagename
			Destination := contentId + "/" + seasonId + "/" + EpisodeId + "/" + nontextrange.Imagename

			// Copy the object from one directory to another in GCS
			srcObject := client.Bucket(bucketName).Object(source)
			attrs, err := srcObject.Attrs(ctx)
			_ = attrs
			if err != nil {
				// Handle the case where the source object doesn't exist
				fmt.Printf("Source object does not exist: %v\n", err)

				// Modify the source path if needed
				// filetrims := strings.Split(item, "/")
				source = contentId + "/" + seasonId + "/" + EpisodeId + "/" + nontextrange.Imagename
				Destination = contentId + "/" + seasonId + "/" + EpisodeId + "/" + nontextrange.Imagename
				srcObject = client.Bucket(bucketName).Object(source)
				// filetrim[0]=filetrims[len(filetrims)-1]

				// Retry the Attrs call
				attrs, err = srcObject.Attrs(ctx)
				if err != nil {
					// Handle the case where the modified source object also doesn't exist
					fmt.Printf("Modified source object does not exist: %v\n", err)
					continue
				}
			}
			dstObject := client.Bucket(bucketName).Object(Destination)
			if _, err := dstObject.CopierFrom(srcObject).Run(ctx); err != nil {
				fmt.Printf("Error copying object: %v", err)
				continue
			}

			url := "https://storage.googleapis.com/" + bucketName + "/" + Destination
			// Don't worry about errors
			response, e := http.Get(url)
			if e != nil {
				fmt.Println(e)
			}
			defer response.Body.Close()

			// Open a file for writing
			file, err := os.Create(nontextrange.Imagename)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()

			// Use io.Copy to just dump the response body to the file
			_, err = io.Copy(file, response.Body)
			if err != nil {
				fmt.Println(err)
			}

			errorr := SizeUploadFileToGcpEpi(ctx, client, nontextrange.Imagename, contentId, seasonId, EpisodeId, url)
			if errorr != nil {
				fmt.Println("error in uploading size upload", errorr)
			}
			os.Remove(nontextrange.Imagename)
			fmt.Println("Success!")
		}
	}
}

func SeasonVarianceFileUPloadGcp(oldseasonid string, seasonId string, contentId string) {
	bucketName := os.Getenv("BUCKET_NAME")
	var newarr []string
	newarr = append(newarr, "poster-image")
	newarr = append(newarr, "overlay-poster-image")
	newarr = append(newarr, "details-background")
	newarr = append(newarr, "mobile-details-background")
	ctx := context.Background()
	client, gcperr := getGCPClient()
	if gcperr != nil {
		fmt.Println("from gcp Connection", gcperr)
	}
	defer client.Close()
	for k := 0; k < len(newarr); k++ {
		item := newarr[k]
		Destination := contentId + "/" + seasonId + "/" + item
		source := contentId + "/" + oldseasonid + "/" + item

		srcObject := client.Bucket(bucketName).Object(source)
		attrs, err := srcObject.Attrs(ctx)
		_ = attrs
		if err != nil {
			// Handle the case where the source object doesn't exist
			fmt.Printf("Source object does not exist: %v\n", err)

			// Modify the source path if needed
			filetrims := strings.Split(item, "/")
			source = contentId + "/" + seasonId + "/" + filetrims[len(filetrims)-1]
			Destination = contentId + "/" + seasonId + "/" + filetrims[len(filetrims)-1]
			srcObject = client.Bucket(bucketName).Object(source)
			item = filetrims[len(filetrims)-1]

			// Retry the Attrs call
			attrs, err = srcObject.Attrs(ctx)
			if err != nil {
				// Handle the case where the modified source object also doesn't exist
				fmt.Printf("Modified source object does not exist: %v\n", err)
				continue
			}
		}
		dstObject := client.Bucket(bucketName).Object(Destination)
		if _, err := dstObject.CopierFrom(srcObject).Run(ctx); err != nil {
			fmt.Printf("Error copying object: %v", err)
			continue
		}

		url := "https://storage.googleapis.com/" + bucketName + "/" + Destination
		response, e := http.Get(url)
		if e != nil {
			fmt.Println(e)
		}
		defer response.Body.Close()

		file, err := os.Create(item)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Println(err)
		}

		errorr := SizeUploadFileToGCP(client, item, contentId, seasonId, url)
		if errorr != nil {
			fmt.Println("error in uploading size upload", errorr)
		}
		os.Remove(item)
		fmt.Println("Success!")
	}
}

func SeasonVarianceTrailerImageUploadGcp(oldcontentId string, oldseasonid string, oldTrailerId string, contentId string, seasonId string, TrailerId string, trailerPosterImage string) {
	bucketName := os.Getenv("BUCKET_NAME")
	var newarr []string

	newarr = append(newarr, trailerPosterImage)
	for k := 0; k < len(newarr); k++ {
		item := newarr[k]
		if strings.Contains(item, "_") {
			filetrim := strings.Split(item, "_")
			Destination := contentId + "/" + seasonId + "/" + TrailerId + "/" + filetrim[0]
			source := "temp/" + item

			ctx := context.Background()
			client, err := storage.NewClient(ctx)
			if err != nil {
				fmt.Println("Failed to create GCS client: ", err)
			}
			defer client.Close()

			// Copy the object from one directory to another in GCS
			srcObject := client.Bucket(bucketName).Object(source)
			dstObject := client.Bucket(bucketName).Object(Destination)
			if _, err := dstObject.CopierFrom(srcObject).Run(ctx); err != nil {
				fmt.Printf("Error copying object: %v", err)
				continue
			}

			url := "https://storage.googleapis.com/" + bucketName + "/" + Destination
			// Don't worry about errors
			response, e := http.Get(url)
			if e != nil {
				fmt.Println(e)
			}
			defer response.Body.Close()

			// Open a file for writing
			file, err := os.Create(filetrim[0])
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()

			// Use io.Copy to just dump the response body to the file
			_, err = io.Copy(file, response.Body)
			if err != nil {
				fmt.Println(err)
			}

			errorr := SizeUploadFileToGCP(client, filetrim[0], contentId, seasonId, url)
			if errorr != nil {
				fmt.Println("error in uploading size upload", errorr)
			}
			os.Remove(filetrim[0])
			fmt.Println("Success!")
		} else {
			Destination := contentId + "/" + seasonId + "/" + TrailerId + "/" + "trailer-poster-image"
			source := oldcontentId + "/" + oldseasonid + "/" + oldTrailerId + "/" + "trailer-poster-image"

			ctx := context.Background()
			client, err := storage.NewClient(ctx)
			if err != nil {
				fmt.Println("Failed to create GCS client: ", err)
			}
			defer client.Close()

			// Copy the object from one directory to another in GCS
			srcObject := client.Bucket(bucketName).Object(source)
			dstObject := client.Bucket(bucketName).Object(Destination)
			if _, err := dstObject.CopierFrom(srcObject).Run(ctx); err != nil {
				fmt.Printf("Error copying object: %v", err)
				continue
			}

			fmt.Println("Success!")
		}
	}
}
