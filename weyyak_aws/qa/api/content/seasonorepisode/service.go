package seasonorepisode

import (
	"content/common"
	l "content/logger"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/thanhpk/randstr"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	qrg := r.Group("/api")
	qrg.Use(common.ValidateToken())
	qrg.POST("/seasons/:id/status/:status", hs.UpdateSeason)
	qrg.POST("/episodes/:id/status/:status", hs.UpdateEpisode)
	qrg.POST("/contentvariances/:id/status/:status", hs.UpdateContentVariance)
	qrg.POST("/seasons/:id/digitalrightstype/:status", hs.ChangeDigitalRightsDuringSeasonUpdation)
	qrg.DELETE("/seasons/:id", hs.DeleteSeason)
	qrg.DELETE("/episodes/:id", hs.DeleteEpisode)
	qrg.DELETE("/contentvariances/:id", hs.DeleteContentVariance)
}

// For  Change Digital Rights During Season Updation -Change Digital Rights During Season Updation
// POST /seasons/{id}/digitalrightstype/{status}
// @Summary Show Change Digital Rights During Season Updation
// @Description post Change Digital Rights During Season Updation
// @Tags Season
// @Accept  json
// @Produce  json
// @security Authorization
// @Param id path string true "Id"
// @Param status path string true "status"
// @Success 200 {array} object c.JSON
// @Router /seasons/{id}/digitalrightstype/{status} [post]
func (hs *HandlerService) ChangeDigitalRightsDuringSeasonUpdation(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var errorFlag bool
	errorFlag = false
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var invalidError InvalidError
	var count int
	type ForRights struct {
		RightsId string `json:"rights_id"`
	}
	var forrights ForRights
	db.Debug().Table("season").Select("rights_id").Where("id=? and deleted_by_user_id is null", c.Param("id")).Find(&forrights)
	db.Debug().Table("season").Where("id=? and deleted_by_user_id is null", c.Param("id")).Count(&count)
	if count > 0 {
		db.Debug().Table("content_rights").Where("id=? ", forrights.RightsId).Update(UpdateRights{DigitalRightsType: c.Param("status")})
		res := map[string]string{
			"id": c.Param("id"),
		}
		l.JSON(c, http.StatusOK, gin.H{"data": res})
		return
	} else {
		errorFlag = true
		invalidError = InvalidError{"error_content_not_found", "The specified condition was not met for 'Id'."}
	}
	var invalid Invalid
	invalid = Invalid{Id: &invalidError}
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"not_found", "Not found.", "", randstr.String(32), invalid}
	if errorFlag {
		l.JSON(c, http.StatusBadRequest, finalErrorResponse)
		return
	}
}

// For  Update Season -Update Season
// POST /api/seasons/{id}/status/{status}
// @Summary Show Update Season
// @Description post Update Season
// @Tags Season
// @Accept  json
// @Produce  json
// @security Authorization
// @Param id path string true "Id"
// @Param status path string true "status"
// @Success 200 {array} object c.JSON
// @Router /api/seasons/{id}/status/{status} [post]
func (hs *HandlerService) UpdateSeason(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var errorFlag bool
	errorFlag = false
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var invalidError InvalidError
	var count int
	var seasondetails common.EpisodeDetails
	db.Debug().Table("season").Where("id=? and deleted_by_user_id is null", c.Param("id")).Count(&count)
	if count > 0 {
		db.Debug().Table("season").Where("id=? and deleted_by_user_id is null", c.Param("id")).Update(UpdateStatus{Status: c.Param("status"), ModifiedAt: time.Now()})
		//	db.Table("episode").Where("season_id=? and deleted_by_user_id is null", c.Param("id")).Update(UpdateStatus{Status: c.Param("status"), ModifiedAt: time.Now()})
		res := map[string]string{
			"id": c.Param("id"),
		}
		db.Debug().Raw("select c.content_key,c.content_type,s.content_id from content c join season s on s.content_id = c.id where s.id=?", c.Param("id")).Find(&seasondetails)
		contentkeyconverted := strconv.Itoa(seasondetails.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(seasondetails.ContentType, c)
		go common.ContentSynching(seasondetails.ContentId, c)
		common.ClearRedisKeyKeys(c, "Menus_Slider_*")
		l.JSON(c, http.StatusOK, gin.H{"data": res})
		return
	} else {
		errorFlag = true
		invalidError = InvalidError{"error_content_not_found", "The specified condition was not met for 'Id'."}
	}
	var invalid Invalid
	invalid = Invalid{Id: &invalidError}
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"not_found", "Not found.", "", randstr.String(32), invalid}
	if errorFlag {
		l.JSON(c, http.StatusBadRequest, finalErrorResponse)
		return
	}
}

// For  Update  Episode status -Update  Episode status
// POST /api/episodes/{id}/status/{status}
// @Summary Show Update Episode status
// @Description post Update Episode status
// @Tags  episode
// @Accept  json
// @Produce  json
// @security Authorization
// @Param id path string true "Id"
// @Param status path string true "status"
// @Success 200 {array} object c.JSON
// @Router /api/episodes/{id}/status/{status} [post]
func (hs *HandlerService) UpdateEpisode(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var errorFlag bool
	errorFlag = false
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var invalidError InvalidError
	var count int
	var episodedetails common.EpisodeDetails
	db.Debug().Table("episode").Where("id=? and deleted_by_user_id is null", c.Param("id")).Count(&count)
	if count > 0 {
		db.Debug().Table("episode").Where("id=? and deleted_by_user_id is null", c.Param("id")).Update(UpdateStatus{Status: c.Param("status"), ModifiedAt: time.Now()})
		res := map[string]string{
			"id": c.Param("id"),
		}
		db.Debug().Raw("select c.content_key,c.content_type,s.content_id from content c join season s on s.content_id = c.id join episode e on e.season_id  = s.id where e.id =?", c.Param("id")).Find(&episodedetails)
		contentkeyconverted := strconv.Itoa(episodedetails.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(episodedetails.ContentType, c)
		fmt.Println("cont id", episodedetails.ContentId)
		go common.ContentSynching(episodedetails.ContentId, c)
		common.ClearRedisKeyKeys(c, "Menus_Slider_*")
		l.JSON(c, http.StatusOK, gin.H{"data": res})
		return
	} else {
		errorFlag = true
		invalidError = InvalidError{"error_content_not_found", "The specified condition was not met for 'Id'."}
	}
	var invalid Invalid
	invalid = Invalid{Id: &invalidError}
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"not_found", "Not found.", "", randstr.String(32), invalid}
	if errorFlag {
		l.JSON(c, http.StatusBadRequest, finalErrorResponse)
		return
	}
}

// For  Update contentvariance status-Update contentvariance status
// POST /api/contentvariances/{id}/status/{status}
// @Summary Show contentvariance status
// @Description post Update  contentvariance status
// @Tags  contentvariance
// @Accept  json
// @Produce  json
// @security Authorization
// @Param id path string true "Id"
// @Param status path string true "status"
// @Success 200 {array} object c.JSON
// @Router /api/contentvariances/{id}/status/{status} [post]
func (hs *HandlerService) UpdateContentVariance(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var errorFlag bool
	errorFlag = false
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var invalidError InvalidError
	var count int
	var contenttype common.ContentType
	var contentid common.ContentID
	db.Debug().Table("content_variance").Where("id=? and deleted_by_user_id is null", c.Param("id")).Count(&count)
	if count > 0 {
		db.Debug().Table("content_variance").Where("id=? and deleted_by_user_id is null", c.Param("id")).Update(UpdateStatus{Status: c.Param("status"), ModifiedAt: time.Now()})
		db.Debug().Raw("select content_id from content_variance where id=?", c.Param("id")).Find(&contentid)
		db.Debug().Raw("select content_type from content where id=?", contentid.ContentId).Find(&contenttype)
		go common.CreateRedisKeyForContentType(contenttype.ContentType, c)
		go common.ContentSynching(contentid.ContentId, c)
		common.ClearRedisKeyKeys(c, "Menus_Slider_*")
	} else {
		errorFlag = true
		invalidError = InvalidError{"error_content_not_found", "The specified condition was not met for 'Id'."}
	}
	var invalid Invalid
	invalid = Invalid{Id: &invalidError}
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"not_found", "Not found.", "", randstr.String(32), invalid}
	if errorFlag {
		l.JSON(c, http.StatusBadRequest, finalErrorResponse)
		return
	}
}

// For Delete Season by season id -Delete Season by season id
// DELETE /api/seasons/{id}
// @Summary Delete Season by season id
// @Description Delete Season by season id
// @Tags Season
// @Accept  json
// @Produce  json
// @security Authorization
// @Param id path string true "Id"
// @Success 200 {array} object c.JSON
// @Router /api/{value}/{id} [delete]
func (hs *HandlerService) DeleteSeason(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var errorFlag bool
	errorFlag = false
	var count int
	var seasondetails common.EpisodeDetails
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var userId = c.MustGet("userid")
	var invalidError InvalidError
	fmt.Println(count, "count is")
	db.Debug().Table("season").Where("id=? and deleted_by_user_id is null", c.Param("id")).Count(&count)
	if count > 0 {
		db.Debug().Table("season").Where("id=? and deleted_by_user_id is null", c.Param("id")).Update(UpdateDetails{DeletedByUserId: userId.(string), ModifiedAt: time.Now()})
		db.Debug().Raw("select c.content_key,c.content_type from content c join season s on s.content_id = c.id where s.id=?", c.Param("id")).Find(&seasondetails)
		contentkeyconverted := strconv.Itoa(seasondetails.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(seasondetails.ContentType, c)
	} else {
		errorFlag = true
		invalidError = InvalidError{"error_content_not_found", "The specified condition was not met for 'Id'."}
	}
	var invalid Invalid
	invalid = Invalid{Id: &invalidError}
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"not_found", "Not found.", "", randstr.String(32), invalid}
	if errorFlag {
		l.JSON(c, http.StatusBadRequest, finalErrorResponse)
		return
	}
}

// For Delete episode details by episode id-Delete episode details by episode id
// DELETE /api/episodes/{id}
// @Summary Delete episode details by episode id
// @Description Delete episode details by episode id
// @Tags  episode
// @Accept  json
// @Produce  json
// @security Authorization
// @Param id path string true "Id"
// @Success 200 {array} object c.JSON
// @Router /api/episodes/{id} [delete]
func (hs *HandlerService) DeleteEpisode(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var errorFlag bool
	errorFlag = false
	var count int
	var episodedetails common.EpisodeDetails
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var userId = c.MustGet("userid")
	var invalidError InvalidError
	db.Debug().Table("episode").Where("id=? and deleted_by_user_id is null", c.Param("id")).Count(&count)
	if count > 0 {
		db.Debug().Table("episode").Where("id=? and deleted_by_user_id is null", c.Param("id")).Update(UpdateDetails{DeletedByUserId: userId.(string), ModifiedAt: time.Now()})
		db.Debug().Raw("select c.content_key,c.content_type from content c join season s on s.content_id = c.id join episode e on e.season_id  = s.id where e.id =?", c.Param("id")).Find(&episodedetails)
		contentkeyconverted := strconv.Itoa(episodedetails.ContentKey)
		go common.CreateRedisKeyForContent(contentkeyconverted, c)
		go common.CreateRedisKeyForContentType(episodedetails.ContentType, c)
	} else {
		errorFlag = true
		invalidError = InvalidError{"error_content_not_found", "The specified condition was not met for 'Id'."}
	}
	var invalid Invalid
	invalid = Invalid{Id: &invalidError}
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"not_found", "Not found.", "", randstr.String(32), invalid}
	if errorFlag {
		l.JSON(c, http.StatusBadRequest, finalErrorResponse)
		return
	}
}

// For Delete content variances details by contentvariance id-Delete content variances details by contentvariance id
// DELETE /api/contentvariances/{id}
// @Summary Delete content variances details by contentvariance id
// @Description Delete content variances details by contentvariance id
// @Tags contentvariance
// @Accept  json
// @Produce  json
// @security Authorization
// @Param id path string true "Id"
// @Success 200 {array} object c.JSON
// @Router /api/contentvariances/{id} [delete]
func (hs *HandlerService) DeleteContentVariance(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var errorFlag bool
	errorFlag = false
	var count int
	var contenttype common.ContentType
	var contentid common.ContentID
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var userId = c.MustGet("userid")
	var invalidError InvalidError
	db.Debug().Table("content_variance").Where("id=? and deleted_by_user_id is null", c.Param("id")).Count(&count)
	if count > 0 {
		db.Debug().Table("content_variance").Where("id=? and deleted_by_user_id is null", c.Param("id")).Update(UpdateDetails{DeletedByUserId: userId.(string), ModifiedAt: time.Now()})
		db.Debug().Raw("select content_id from content_variance where id=?", c.Param("id")).Find(&contentid)
		db.Debug().Raw("select content_type from content where id=?", contentid.ContentId).Find(&contenttype)
		go common.CreateRedisKeyForContentType(contenttype.ContentType, c)
	} else {
		errorFlag = true
		invalidError = InvalidError{"error_content_not_found", "The specified condition was not met for 'Id'."}

	}
	var invalid Invalid
	invalid = Invalid{Id: &invalidError}
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"not_found", "Not found.", "", randstr.String(32), invalid}
	if errorFlag {
		l.JSON(c, http.StatusBadRequest, finalErrorResponse)
		return
	}
}
