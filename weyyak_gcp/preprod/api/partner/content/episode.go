package content

import (
	"fmt"
	common "masterdata/common"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetEpisodeDetailsByEpisodeId - Get Episode Details By Episode Id
// GET /v1/episode/:contentId
// @Description Get Episode Details By Episode Id
// @Tags MultiTier
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param contentId path string true "Content Id"
// @Param Country query string false "Country code of the user"
// @Success 200  object EpisodeResult
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /v1/episode/{contentId} [get]
func (hs *HandlerService) GetEpisodeDetailsByEpisodeId(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	// udb := c.MustGet("UDB").(*gorm.DB)
	// _ = udb

	UserId := c.MustGet("userid")
	if UserId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization", "status": http.StatusUnauthorized})
		return
	}
	serverError := common.ServerErrorResponse()
	notFound := common.NotFoundErrorResponse()
	episode_key, _ := strconv.Atoi(c.Param("contentId"))
	var count int
	var CountryResult int32
	if err := db.Table("episode").Where("third_party_episode_key=?", episode_key).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if count < 1 {
		c.JSON(http.StatusNotFound, notFound)
		return
	}
	type ContentId struct {
		ContentId string
	}
	var watchnowsupplier ContentId
	if UserId == os.Getenv("WATCH_NOW") {
		if watchnowResult := db.Debug().Table("content as c").Select("c.id as content_id").Joins("left join season s on s.content_id = c.id").Joins("left join episode e on e.season_id = s.id").Where("c.watch_now_supplier = 'true' and e.third_party_episode_key =?", episode_key).Find(&watchnowsupplier).Error; watchnowResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	}
	var finalContentResult FinalSeasonResult
	var contentResult EpisodeResult
	if episodeResult := db.Debug().Table("episode as e").Select(`
				e.number as episode_number ,
				s.content_id ,
				e.third_party_episode_key as episode_key ,
				pi2.duration, pi2.video_content_id,
				e.synopsis_english , e.synopsis_arabic, e.season_id, e.has_poster_image ,
				e.has_dubbing_script ,e.has_subtitling_script ,
				e.english_meta_title ,e.arabic_meta_title ,e.english_meta_description ,
				e.arabic_meta_description ,e.created_at ,e.modified_at ,
				e.id ,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,
				cpi.notes ,cpi.intro_start ,cpi.outro_start,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,
				ct.subtitling_language ,e.cast_id ,e.music_id ,e.tag_info_id ,s.rights_id
		`).
		Joins("left join season s on e.season_id =s.id").
		Joins("left join content_rights cr2 on cr2.id = s.rights_id").
		Joins("left join playback_item pi2 on pi2.id =e.playback_item_id").
		Joins("left join content_primary_info cpi on cpi.id =e.primary_info_id").
		Joins("left join content_translation ct on ct.id =pi2.translation_id").
		Joins("join content_rights cr on cr.id = pi2.rights_id").
		Joins("join about_the_content_info atci on  atci.Id = s.about_the_content_info_id").
		Where(`e.deleted_by_user_id is null and e.status = 1 and 
		( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and 
		(cr2.digital_rights_start_date <= NOW() or cr2.digital_rights_start_date is null) and 
		(cr2.digital_rights_end_date >= NOW() or cr2.digital_rights_end_date is null) and 
		atci.supplier !='Others' and
		e.third_party_episode_key =?`, episode_key).Find(&finalContentResult).Error; episodeResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var country string
	if c.Request.URL.Query()["Country"] != nil {
		country = c.Request.URL.Query()["Country"][0]
	}

	// if country == "" || country == "all" || country == "All" {
	// 	CountryResult = 1
	// } else if country != "" {
	CountryResult = common.Countrys(country)
	// }

	contentResult.EpisodeNumber = finalContentResult.EpisodeNumber
	contentResult.ContentKey = finalContentResult.ContentId
	contentResult.EpisodeKey = finalContentResult.EpisodeKey
	contentResult.Length = finalContentResult.Duration
	contentResult.VideoContentUrl = os.Getenv("CONTENT_URL") + finalContentResult.VideoContentId
	contentResult.SynopsisEnglish = finalContentResult.SynopsisEnglish
	contentResult.SynopsisArabic = finalContentResult.SynopsisArabic
	contentResult.SeasonId = finalContentResult.SeasonId
	/*Fetch Digital_right_Regions*/
	var digitalRightsRegions []DigitalRightsRegions
	if countryError := db.Debug().Table("content_rights_country").Select("country_id").Where("content_rights_id=?", finalContentResult.RightsId).Scan(&digitalRightsRegions).Error; countryError != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	// var RegionRights []int
	// for _, idarr := range digitalRightsRegions {
	// 	RegionRights = append(RegionRights, idarr.CountryId)
	// }
	/*for digital rights*/

	var digitalRights []int

	fmt.Println("CountryResult-------->", CountryResult)

	if CountryResult == 0 {

		if country != "" {
			if country == "all" || country == "All" {

				if common.CountryCount() == len(digitalRightsRegions) {
					digitalRights = nil
				} else {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}

			}
		} else if country == "" {
			for _, region := range digitalRightsRegions {
				digitalRights = append(digitalRights, region.CountryId)
			}
		}

	} else if CountryResult == 1 {
		for _, region := range digitalRightsRegions {
			digitalRights = append(digitalRights, region.CountryId)
		}

	} else {
		var isCheck bool = false
		for _, region := range digitalRightsRegions {
			if CountryResult == int32(region.CountryId) {
				isCheck = true
			}
		}

		if !isCheck {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	}

	contentResult.DigitalRightsRegions = digitalRights

	contentResult.PrimaryInfo.Number = finalContentResult.EpisodeNumber
	contentResult.PrimaryInfo.VideoContentId = finalContentResult.VideoContentId
	contentResult.PrimaryInfo.SynopsisEnglish = finalContentResult.SynopsisEnglish
	contentResult.PrimaryInfo.SynopsisArabic = finalContentResult.SynopsisArabic
	contentResult.PrimaryInfo.OriginalTitle = finalContentResult.OriginalTitle
	contentResult.PrimaryInfo.AlternativeTitle = finalContentResult.AlternativeTitle
	contentResult.PrimaryInfo.ArabicTitle = finalContentResult.ArabicTitle
	contentResult.PrimaryInfo.TransliteratedTitle = finalContentResult.TransliteratedTitle
	contentResult.PrimaryInfo.Notes = finalContentResult.Notes
	/* Fetch content_cast*/
	var contentCast Cast
	if castResult := db.Table("content_cast cc").Select("cc.main_actor_id,cc.main_actress_id,actor.english_name as main_actor_english,actor.arabic_name as main_actor_arabic,actress.english_name as main_actress_english,actress.arabic_name as main_actress_arabic").
		Joins("left join actor actor on actor.id =cc.main_actor_id").
		Joins("left join actor actress on actress.id =cc.main_actress_id").
		Where("cc.id=?", finalContentResult.CastId).Scan(&contentCast).Error; castResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	contentResult.Cast.CastId = finalContentResult.CastId
	contentResult.Cast.MainActorId = contentCast.MainActorId
	contentResult.Cast.MainActressId = contentCast.MainActressId
	contentResult.Cast.MainActorEnglish = contentCast.MainActorEnglish
	contentResult.Cast.MainActorArabic = contentCast.MainActorArabic
	contentResult.Cast.MainActressEnglish = contentCast.MainActressEnglish
	contentResult.Cast.MainActressArabic = contentCast.MainActressArabic
	var contentActor []ContentActor
	if actorResult := db.Table("content_actor ca").Select("a.english_name as actor_english,a.arabic_name as actor_arabic,a.id as actor_id,w.id as writer_id,w.english_name as writer_english,w.arabic_name as writer_arabic,d.id as director_id,d.english_name as director_english,d.arabic_name as director_arabic").
		Joins("left join actor a on a.id =ca.actor_id").
		Joins("left join content_writer cw on cw.cast_id =ca.cast_id").
		Joins("left join writer w on w.id =cw.writer_id").
		Joins("left join content_director cd on cd.cast_id =ca.cast_id").
		Joins("left join director d on d.id =cd.director_id").
		Where("ca.cast_id=?", finalContentResult.CastId).Scan(&contentActor).Error; actorResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var actorEnglish, actorId, actorArabic, writerId, writerEnglish, writerArabic, directorId, directorEnglish, directorArabic []string
	for _, actorIds := range contentActor {
		actorId = append(actorId, actorIds.ActorId)
		actorEnglish = append(actorEnglish, actorIds.ActorEnglish)
		actorArabic = append(actorArabic, actorIds.ActorArabic)
		writerId = append(writerId, actorIds.WriterId)
		writerEnglish = append(writerEnglish, actorIds.WriterEnglish)
		writerArabic = append(writerArabic, actorIds.WriterArabic)
		directorId = append(directorId, actorIds.DirectorId)
		directorEnglish = append(directorEnglish, actorIds.DirectorEnglish)
		directorArabic = append(directorArabic, actorIds.DirectorArabic)
	}
	contentResult.Cast.ActorIds = common.RemoveDuplicateValues(actorId)
	contentResult.Cast.ActorEnglish = common.RemoveDuplicateValues(actorEnglish)
	contentResult.Cast.ActorArabic = common.RemoveDuplicateValues(actorArabic)
	contentResult.Cast.WriterId = common.RemoveDuplicateValues(writerId)
	contentResult.Cast.WriterEnglish = common.RemoveDuplicateValues(writerEnglish)
	contentResult.Cast.WriterArabic = common.RemoveDuplicateValues(writerArabic)
	contentResult.Cast.DirectorIds = common.RemoveDuplicateValues(directorId)
	contentResult.Cast.DirectorEnglish = common.RemoveDuplicateValues(directorEnglish)
	contentResult.Cast.DirectorArabic = common.RemoveDuplicateValues(directorArabic)

	/* Fetch content_music*/
	var contentMusic []ContentMusic
	if actorResult := db.Debug().Table("content_singer cs").Select("s.id as singer_ids,s.english_name as singers_english,s.arabic_name as singers_arabic,mc.id as music_composer_ids,mc.english_name as music_composers_english ,mc.arabic_name as music_composers_arabic,sw.id as song_writer_ids,sw.english_name as song_writers_english,sw.arabic_name as song_writers_arabic").
		Joins("left join singer s on s.id =cs.singer_id").
		Joins("left join content_music_composer cmc on cmc.music_id =cs.music_id").
		Joins("left join music_composer mc on mc.id =cmc.music_composer_id").
		Joins("left join content_song_writer csw on csw.music_id =cs.music_id ").
		Joins("left join song_writer sw on sw.id =csw.song_writer_id").
		Where(" cs.music_id=?", finalContentResult.MusicId).Scan(&contentMusic).Error; actorResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	var singerEnglish, singerId, singerArabic, composerId, composerEnglish, SongWriterId, composerArabic, SongWriterEnglish, SongWriterArabic []string
	for _, musicIds := range contentMusic {
		singerId = append(singerId, musicIds.SingerIds)
		singerEnglish = append(singerEnglish, musicIds.SingersEnglish)
		singerArabic = append(singerArabic, musicIds.SingersArabic)
		composerId = append(composerId, musicIds.MusicComposerIds)
		composerEnglish = append(composerEnglish, musicIds.MusicComposersEnglish)
		composerArabic = append(composerArabic, musicIds.MusicComposersArabic)
		SongWriterId = append(SongWriterId, musicIds.SongWriterIds)
		SongWriterEnglish = append(SongWriterEnglish, musicIds.SongWritersEnglish)
		SongWriterArabic = append(SongWriterArabic, musicIds.SongWritersArabic)
	}
	contentResult.Music.MusicId = finalContentResult.MusicId
	contentResult.Music.SingerIds = common.RemoveDuplicateValues(singerId)
	contentResult.Music.SingersEnglish = common.RemoveDuplicateValues(singerEnglish)
	contentResult.Music.SingersArabic = common.RemoveDuplicateValues(singerArabic)
	contentResult.Music.MusicComposerIds = common.RemoveDuplicateValues(composerId)
	contentResult.Music.MusicComposersEnglish = common.RemoveDuplicateValues(composerEnglish)
	contentResult.Music.MusicComposersArabic = common.RemoveDuplicateValues(composerArabic)
	contentResult.Music.SongWriterIds = common.RemoveDuplicateValues(SongWriterId)
	contentResult.Music.SongWritersEnglish = common.RemoveDuplicateValues(SongWriterEnglish)
	contentResult.Music.SongWritersArabic = common.RemoveDuplicateValues(SongWriterArabic)
	/*Fetch tag_info*/
	var contentTagInfo []ContentTag
	db.Table("content_tag ct").Select("tdt.name").
		Joins("left join textual_data_tag tdt on tdt.id =ct.textual_data_tag_id").
		Where("ct.tag_info_id=?", finalContentResult.TagInfoId).Scan(&contentTagInfo)
	var tagInfo []string
	for _, tagInfoIds := range contentTagInfo {
		tagInfo = append(tagInfo, tagInfoIds.Name)
	}
	contentResult.TagInfo.Tags = tagInfo
	if len(tagInfo) < 1 {
		buffer := make([]string, 0)
		contentResult.TagInfo.Tags = buffer
	}
	/*Non_textual Data*/
	if finalContentResult.HasPosterImage {
		contentResult.NonTextualDataEpisode.PosterImage = os.Getenv("IMAGE_URL_GO") + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/" + finalContentResult.Id + os.Getenv("POSTER_IMAGE")
	} else {
		contentResult.NonTextualDataEpisode.PosterImage = ""
	}
	if finalContentResult.HasDubbingScript {
		contentResult.NonTextualDataEpisode.DubbingScript = os.Getenv("IMAGE_URL_GO") + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/" + finalContentResult.Id + os.Getenv("DUBBLING_SCRIPT")
	} else {
		contentResult.NonTextualDataEpisode.DubbingScript = ""
	}
	if finalContentResult.HasSubtitlingScript {
		contentResult.NonTextualDataEpisode.SubtitlingScript = os.Getenv("IMAGE_URL_GO") + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/" + finalContentResult.Id + "/subtitling-script"
	} else {
		contentResult.NonTextualDataEpisode.SubtitlingScript = ""
	}
	/*Translation data*/
	contentResult.Translation.LanguageType = common.LanguageOriginTypes(finalContentResult.LanguageType)
	contentResult.Translation.DubbingLanguage = finalContentResult.DubbingLanguage
	contentResult.Translation.DubbingDialectId = finalContentResult.DubbingDialectId
	contentResult.Translation.DubbingDialectName = common.DialectIdname(finalContentResult.DubbingDialectId, "en")
	contentResult.Translation.SubtitlingLanguage = finalContentResult.SubtitlingLanguage
	/*SeoDetails*/
	contentResult.SeoDetails.EnglishMetaTitle = finalContentResult.EnglishMetaTitle
	contentResult.SeoDetails.ArabicMetaTitle = finalContentResult.ArabicMetaTitle
	contentResult.SeoDetails.EnglishMetaDescription = finalContentResult.EnglishMetaDescription
	contentResult.SeoDetails.ArabicMetaDescription = finalContentResult.ArabicMetaDescription

	contentResult.CreatedAt = finalContentResult.CreatedAt
	contentResult.ModifiedAt = finalContentResult.ModifiedAt
	contentResult.Id = finalContentResult.Id
	// if CountryResult != 0 {
	c.JSON(http.StatusOK, gin.H{"data": contentResult})
	// }
}
