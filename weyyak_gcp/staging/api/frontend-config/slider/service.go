package slider

import (
	"net/http"
	"strconv"
	"strings"

	//"strings"
	"context"
	"frontend_config/common"
	"frontend_config/fragments"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"github.com/thanhpk/randstr"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	qrg := r.Group("/api")
	qrg.Use(common.ValidateToken())
	qrg.POST("/sliders/:id/available", hs.SliderAvailabilityUpdate)
	qrg.GET("/slider/types", hs.GetSliderTypes)
	qrg.DELETE("/sliders/:id", hs.DeleteSliderList)
	qrg.GET("/sliders", hs.GetSliderList)
	//qrg.POST("/sliders", hs.CreateAndUpdateSlider)
	//qrg.POST("/sliders/:id", hs.CreateAndUpdateSlider)
	qrg.POST("/sliders/:id", hs.UpdateSlider)
	qrg.POST("/sliders", hs.CreateSlider)
	qrg.GET("/sliders/:id", hs.GetSliderDetailsById)
	qrg.GET("/slider/previewlayouts", hs.GetSliderPreviewlayouts)
	qrg.GET("/sliders/summary", hs.GetSlidersBasedOnSearchInPage)
	qrg.GET("/sliders/notifications", hs.GetSliderNotifications)
	qrg.GET("/sliders/:id/region", hs.GetallRegionsBasedOnSliderId)

	/*Error code Exception URL*/
	/*Get Slider Types*/
	qrg.POST("/slider/types", hs.GetSliderTypes)
	qrg.PUT("/slider/types", hs.GetSliderTypes)
	qrg.DELETE("/slider/types", hs.GetSliderTypes)

}

// For slider update -disables or enables slider
// GET /api/sliders/:id/available
// @Summary Show disable or enable slider
// @Description post disable or enable slider
// @Tags Slider
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param id path string true "Id"
// @Param body body SliderUpdate true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /api/sliders/{id}/available [post]
func (hs *HandlerService) SliderAvailabilityUpdate(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var sliderAvailability SliderAvailability
	var errorresponse = common.ServerErrorResponse()
	if err := c.ShouldBindJSON(&sliderAvailability); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": http.StatusBadRequest})
		return
	}
	sliderurlId := c.Param("id")
	sliderId := sliderAvailability.Id
	if sliderurlId != sliderId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide valid slider Id.", "Status": http.StatusBadRequest})
		return
	}
	var totalcount int
	if err := db.Debug().Table("slider").Where("id=?", c.Param("id")).Count(&totalcount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	if totalcount < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "error_slider_not_found", "message": "The specified condition was not met for 'Id'.", "Status": http.StatusBadRequest})
		return
	} else {
		if err := db.Debug().Table("slider").Where("id=?", c.Param("id")).Update(map[string]interface{}{"is_disabled": sliderAvailability.IsDisabled, "modified_at": time.Now()}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		/* update dirty count when slider get updated */
		go common.SliderSynching(c.Param("id"), c)
		c.JSON(http.StatusOK, gin.H{})
	}
}

// GetSliderTypes -  fetches all Slider types
// GET /api/slider/types
// @Summary Show a list of all slidertypes
// @Description get list of all slider types
// @Tags Slider
// @Accept  json
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /api/slider/types [get]
func (hs *HandlerService) GetSliderTypes(c *gin.Context) {
	/*405(Request-method) Handling*/
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var SliderType []sliderTypes
	var errorresponse = common.ServerErrorResponse()
	fields := "id, name"
	if data := db.Debug().Table("slider_types").Select(fields).Scan(&SliderType).Error; data != nil {
		c.JSON(http.StatusBadRequest, errorresponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": SliderType})
}

// GetSliderList - Get Slider List
// GET /sliders/
// @Summary Get Slider List
// @Description Get Slider List
// @Tags slider
// @Accept  json
// @Produce  json
// @Param searchText path string false "Search Text"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Success 200 {array} object c.JSON
// @Router /sliders/ [get]
func (hs *HandlerService) GetSliderList(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var errorresponse = common.ServerErrorResponse()
	db := c.MustGet("DB").(*gorm.DB)
	var limit, offset int64

	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	if limit == 0 {
		limit = 50
	}
	var searchText string
	if c.Request.URL.Query()["searchText"] != nil {
		searchText = c.Request.URL.Query()["searchText"][0]
	}
	// rawquery := "select s.id,s.name,s.is_disabled,s.scheduling_end_date,string_agg(distinct p.english_title::text, ',')::varchar as found_in,string_agg(distinct case when ps.order = 0 then 'true' else 'false' end,',')::bool as is_default_for_any_page,case when count(distinct pp.id) = 10 then 'All Platforms' else string_agg(distinct pp.platform :: varchar, ',') end as available_on,case when count(distinct c.english_name)= 248 then 'All Regions' else string_agg(distinct c.english_name ::varchar, ',' order by c.english_name asc ) end as region,case when count(distinct sc.country_id)<5 or count(distinct sc.country_id)= 248 then 'false' else 'true' end as has_more_regions from slider s left join page_slider ps on ps.slider_id =s.id left join page p on p.id =ps.page_id join slider_country sc on s.id = sc.slider_id join country c on sc.country_id = c.id and sc.slider_id = s.id join slider_target_platform stp on s.id = stp.slider_id join publish_platform pp on stp.target_platform = pp.id where s.deleted_by_user_id is null and s.name != ''"
	rawquery := `
			select
				s.id,
				s.name,
				s.is_disabled,
				s.scheduling_end_date,
				string_agg(distinct p.english_title::text,
				',')::varchar as found_in,
				string_agg(distinct case
					when ps.order = 0 then 'true'
					else 'false'
				end,
				',')::boolean as is_default_for_any_page,
				case
					when count(distinct pp.id) = 10 then 'All Platforms'
					else string_agg(distinct pp.platform :: varchar,
					',')
				end as available_on,
				case when count(distinct c.english_name) = 248 then 'All Regions' 
					else (SELECT string_agg(c.english_name::varchar, ',' ORDER BY c.english_name)
							FROM slider_country sc
							JOIN country c ON sc.country_id = c.id AND sc.slider_id = s.id
						)
				end as region,
				
				case
					when count(distinct sc.country_id)<5
					or count(distinct sc.country_id)= 248 then 'false'
					else 'true'
				end as has_more_regions
			from
				slider s
			left join page_slider ps on
				ps.slider_id = s.id
			left join page p on
				p.id = ps.page_id
			join slider_country sc on
				s.id = sc.slider_id
			join country c on
				sc.country_id = c.id
				and sc.slider_id = s.id
			join slider_target_platform stp on
				s.id = stp.slider_id
			join publish_platform pp on
				stp.target_platform = pp.id
			where
				s.deleted_by_user_id is null
				and s.name != ''
		`

	if searchText != "" {
		rawquery += " and s.name ilike '%" + searchText + "%'"
	}
	rawquery += "  group by s.id order by s.modified_at DESC "

	var sliderlist, finalsliderlist, totalCount []SliderList

	if data := db.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sliderlist).Error; data != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	if errCount := db.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	for _, sliderResponse := range sliderlist {
		fmt.Println("@@@", sliderResponse.Region)
		var regionsarray []string
		if sliderResponse.Region != "All Regions" {
			regions := strings.Split(sliderResponse.Region, ",")
			for _, val := range regions {
				regionsarray = append(regionsarray, val)
			}
			if len(regionsarray) >= 5 {
				regionarray := regionsarray[:5]
				sliderResponse.Region = strings.Join(regionarray, ",")
				// sliderResponse.HasMoreRegions = true
			} else {
				sliderResponse.Region = strings.Join(regionsarray, ",")
				// sliderResponse.HasMoreRegions = false
			}
		} else {
			sliderResponse.Region = "All Regions"
			// sliderResponse.HasMoreRegions = false
		}
		finalsliderlist = append(finalsliderlist, sliderResponse)
	}
	pages := map[string]int{
		"size":   len(totalCount),
		"offset": int(offset),
		"limit":  int(limit),
	}
	c.JSON(http.StatusOK, gin.H{"pagination": pages, "data": finalsliderlist})
}

// DeleteSliderList - delete slider from slider table
// DELETE /api/slider/:id
// @Summary Delete slider
// @Description delete Delete slider
// @Tags Slider
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 200 {array} object c.JSON
// @Router /api/slider/{id} [delete]
func (hs *HandlerService) DeleteSliderList(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var Slider SliderList
	userId := c.MustGet("userid")
	serverError := common.ServerErrorResponse()
	notFound := common.NotFoundErrorResponse()
	if err := db.Debug().Table("slider").Where("id=?", c.Param("id")).Find(&Slider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if len(Slider.ID) > 0 {
		if err := db.Debug().Table("slider").Where("id=?", c.Param("id")).Update(UpdateDetails{DeletedByUserId: userId.(string), ModifiedAt: time.Now()}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	} else {
		c.JSON(http.StatusNotFound, notFound)
	}
}

// GetSliderDetailsById - Get slider details by slider id
// GET /api/sliders/:id
// @Summary Get  slider details
// @Description Get Slider details by slider id
// @Tags Slider
// @Accept json
// @Security Authorization
// @Produce json
// @Param id path string true "Slider Id"
// @Success 200 {array} object c.JSON
// @Router /api/sliders/{id} [get]
func (hs *HandlerService) GetSliderDetailsById(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	serverError := common.ServerErrorResponse()
	notFound := common.NotFoundErrorResponse()
	db := c.MustGet("DB").(*gorm.DB)
	var Slider SliderDetails
	var platforms []SliderTargetPlatform
	var regions []SliderCountry
	var pages []SliderPageDetails
	if err := db.Debug().Table("slider").Select("id,name,scheduling_start_date,scheduling_end_date,type,black_area_playlist_id,green_area_playlist_id,red_area_playlist_id").Where("id=?", c.Param("id")).Find(&Slider).Error; err != nil {
		c.JSON(http.StatusNotFound, notFound)
		return
	}
	if err := db.Debug().Table("slider_country").Select("country_id").Where("slider_id=?", c.Param("id")).Find(&regions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	for _, region := range regions {
		Slider.Regions = append(Slider.Regions, region.CountryId)
	}
	if err := db.Debug().Table("slider_target_platform").Select("target_platform").Where("slider_id=?", c.Param("id")).Find(&platforms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	for _, platform := range platforms {
		Slider.PublishingPlatforms = append(Slider.PublishingPlatforms, platform.TargetPlatform)
	}
	if Slider.BlackAreaPlaylistId != "" && Slider.BlackAreaPlaylistId != "00000000-0000-0000-0000-000000000000" {
		var playlist SliderPlaylistDetails
		if err := db.Debug().Table("playlist").Select("id,english_title,arabic_title,is_disabled").Where("id=?", Slider.BlackAreaPlaylistId).Find(&playlist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		Slider.BlackAreaPlaylist = &playlist
	}
	if Slider.GreenAreaPlaylistId != "" && Slider.GreenAreaPlaylistId != "00000000-0000-0000-0000-000000000000" {
		var playlist SliderPlaylistDetails
		if err := db.Debug().Table("playlist").Select("id,english_title,arabic_title,is_disabled").Where("id=?", Slider.GreenAreaPlaylistId).Find(&playlist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		Slider.GreenAreaPlaylist = &playlist
	}
	if Slider.RedAreaPlaylistId != "" && Slider.RedAreaPlaylistId != "00000000-0000-0000-0000-000000000000" {
		var playlist SliderPlaylistDetails
		if err := db.Debug().Table("playlist").Select("id,english_title,arabic_title,is_disabled").Where("id=?", Slider.RedAreaPlaylistId).Find(&playlist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		Slider.RedAreaPlaylist = &playlist
	}
	if err := db.Debug().Table("page_slider ps").Select("p.id,p.english_title,p.arabic_title,p.is_disabled,case when ps.order=0 then true else false end as is_default,case when p.page_type=1 then true else false end as is_home,json_agg(distinct ptp.target_platform)::varchar as platforms").
		Joins("left join page p on p.id = ps.page_id").
		Joins("left join page_target_platform ptp on ptp.page_id=p.id").
		Where("ps.slider_id =?", Slider.Id).
		Group("p.id,ps.order").Find(&pages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var siglePage SliderPageDetails
	var allPages []SliderPageDetails
	for _, pageinfo := range pages {
		siglePage.Id = pageinfo.Id
		siglePage.EnglishTitle = pageinfo.EnglishTitle
		siglePage.ArabicTitle = pageinfo.ArabicTitle
		siglePage.IsDefault = pageinfo.IsDefault
		siglePage.IsDisabled = pageinfo.IsDisabled
		siglePage.IsHome = pageinfo.IsHome
		/* Platforms */
		if pageinfo.Platforms == "[null]" {
			buffer := make([]int, 0)
			siglePage.PublishingPlatforms = buffer
		} else {
			platforms, _ := common.JsonStringToIntSliceOrMap(pageinfo.Platforms)
			siglePage.PublishingPlatforms = platforms
		}
		allPages = append(allPages, siglePage)
	}
	Slider.Pages = allPages
	if len(allPages) < 1 {
		buffer := make([]SliderPageDetails, 0)
		Slider.Pages = buffer
	}
	c.JSON(http.StatusOK, gin.H{"data": Slider})
	return
}

// UpdateSlider - update slider
// GET /api/sliders/:id
// @Summary Update slider details
// @Description Update slider details
// @Tags Slider
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param id path string false "Id"
// @Param body body CreateUpdateSliderRequest true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /api/sliders/{id} [post]
func (hs *HandlerService) UpdateSlider(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	serverError := common.ServerErrorResponse()
	db := c.MustGet("DB").(*gorm.DB)
	ctx := context.Background()
	tx := db.BeginTx(ctx, nil)
	var request CreateUpdateSliderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var errorFlag bool
	errorFlag = false

	if request.Regions == nil || request.PublishingPlatforms == nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	/* Name Error */
	var nameUniqueCheck int
	var nameError common.NameError
	db.Debug().Select("slider").Where("name=? and id!=? ", request.Name, c.Param("id")).Count(&nameUniqueCheck)
	if nameUniqueCheck > 0 {
		errorFlag = true
		nameError = common.NameError{Code: "error_slider_name_not_unique", Description: "Slider with specified 'Name' of " + request.Name + " already exists."}
	}
	// if !common.NumberRegex(request.Name) {
	// 	errorFlag = true
	// 	nameError = common.NameError{Code: "error_slider_Name_invalid", Description: "At least one alphanumeric character is required(" + request.Name + ")."}
	// }
	/* Black-Playlist validations */
	var blackerror common.BlackAreaPlaylistError
	// if request.BlackAreaPlaylistId != "" {
	// 	var count int
	// 	db.Table("playlist p").Where("p.id =? and p.scheduling_start_date <= ? and p.scheduling_end_date >=?", request.BlackAreaPlaylistId, request.SchedulingStartDate, request.SchedulingEndDate).Count(&count)
	// 	if count < 1 {
	// 		errorFlag = true
	// 		blackerror = common.BlackAreaPlaylistError{Code: "error_slider_balckareaplaylist_scheduling_invalid", Description: "Black Area playlist should match to slider scheduling."}
	// 	}
	// }
	// if request.BlackAreaPlaylistId == "" && request.Type == 3 {
	// 	errorFlag = true
	// 	blackerror = common.BlackAreaPlaylistError{Code: "error_slider_blackareaplaylist_not_empty", Description: "Black Area playlist should be assigned."}
	// }
	/* Green-Playlist Validatopns */
	var greenerror common.GreenAreaPlaylistError
	// if request.GreenAreaPlaylistId != "" {
	// 	var greencount int
	// 	db.Table("playlist p").Where("p.id =? and p.scheduling_start_date <= ? and p.scheduling_end_date >=?", request.GreenAreaPlaylistId, request.SchedulingStartDate, request.SchedulingEndDate).Count(&greencount)
	// 	if greencount < 1 {
	// 		errorFlag = true
	// 		greenerror = common.GreenAreaPlaylistError{Code: "error_slider_greenareaplaylist_scheduling_invalid", Description: "Green Area playlist should match to slider scheduling."}
	// 	}
	// }
	// if request.GreenAreaPlaylistId == "" && request.Type == 1 {
	// 	errorFlag = true
	// 	greenerror = common.GreenAreaPlaylistError{Code: "error_slider_greenareaplaylist_not_empty", Description: "Green Area playlist should be assigned."}
	// }
	/* Red-Playlist Validations */
	var rederror common.RedAreaPlaylistError
	// if request.RedAreaPlaylistId != "" {
	// 	var redcount int
	// 	db.Table("playlist p").Where("p.id =? and p.scheduling_start_date <= ? and p.scheduling_end_date >=?", request.RedAreaPlaylistId, request.SchedulingStartDate, request.SchedulingEndDate).Count(&redcount)
	// 	if redcount > 0 {
	// 		var totalCount int
	// 		db.Raw("select count(pi2.id) from playlist_item pi2 left join playlist_item_content pic on pic.playlist_item_id = pi2.id where pi2.playlist_id =?", request.RedAreaPlaylistId).Count(&totalCount)
	// 		if totalCount < 7 {
	// 			errorFlag = true
	// 			rederror = common.RedAreaPlaylistError{Code: "error_slider_redareaplaylist_invalid", Description: "Red Area should contain not less than 7 published Movies/Seasons/Series with matching scheduling."}
	// 		}
	// 	} else {
	// 		fmt.Println(redcount)
	// 		errorFlag = true
	// 		rederror = common.RedAreaPlaylistError{Code: "error_slider_redareaplaylist_invalid", Description: "Red Area playlist should match to slider scheduling."}
	// 	}
	// }

	// if request.RedAreaPlaylistId == "" && request.Type == 2 {
	// 	errorFlag = true
	// 	rederror = common.RedAreaPlaylistError{Code: "error_slider_redareaplaylist_not_empty", Description: "Red Area playlist should be assigned."}
	// }
	/* End of playlist validations */
	var pageValidation common.PagesIds
	var pagetype common.PageType
	var pageNames []string
	type PageName struct {
		EnglishTitle string `json:"english_title"`
	}
	var pageName PageName
	for _, pageid := range request.PagesIds {
		// check if home page go to default slider scenario
		db.Debug().Table("page p").Select("p.page_type").Where("p.id = ?", pageid).Find(&pagetype)
		if pagetype.PageType != 1 {
			db.Debug().Table("page_slider ps").Select("p.english_title").Joins("join page p on p.id =ps.page_id").Where("ps.page_id=? and ps.slider_id !=?", pageid, c.Param("id")).Find(&pageName)
			if pageName.EnglishTitle != "" {
				pageNames = append(pageNames, pageName.EnglishTitle)
			}
		} else if pagetype.PageType == 1 {
			db.Debug().Table("page_slider ps").Select("p.english_title").Joins("join page p on p.id =ps.page_id").Where("ps.page_id=? and ps.slider_id !=? and ps.order != 1", pageid, c.Param("id")).Find(&pageName)
			if pageName.EnglishTitle != "" {
				pageNames = append(pageNames, pageName.EnglishTitle)
			}
		}
	}
	if len(pageNames) > 0 {
		fmt.Println(len(pageNames), pageNames)
		errorFlag = true
		pageValidation = common.PagesIds{Code: "error_slider_should_not_conflict_with_other_pages_sliders_platforms", Description: "Slider has platforms intersections for the next sliders and pages:" + strings.Join(pageNames, " ,")}
	}

	var invalid common.Invalidsslider
	if nameError.Code != "" {
		invalid.NameError = nameError
	}
	if blackerror.Code != "" {
		invalid.BlackAreaPlaylistError = blackerror
	}
	if rederror.Code != "" {
		invalid.RedAreaPlaylistError = rederror
	}
	if greenerror.Code != "" {
		invalid.RedAreaPlaylistError = greenerror
	}
	if pageValidation.Code != "" {
		invalid.PagesIds = pageValidation
	}
	var finalErrorResponse common.FinalErrorResponseslider
	finalErrorResponse = common.FinalErrorResponseslider{Error: "invalid_request", Description: "Validation failed.", Code: "error_validation_failed", RequestId: randstr.String(32), Invalid: invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
	type PageSliderOrder struct {
		Order int `json:"order"`
	}
	var platform SliderTargetPlatform
	var platforms []interface{}
	var region SliderCountry
	var regions []interface{}
	var pageSlider PageSlider
	var pageSliders []interface{}
	sliderId := c.Param("id")
	var slider Slider
	slider.Name = request.Name
	//res, _ := strconv.Atoi(request.Type)
	slider.Type = request.Type
	slider.DeletedByUserId = nil
	slider.IsDisabled = false
	//	layout := "2006-01-02T15:04:05.000Z"
	//	sdate, _ := time.Parse(layout, request.SchedulingStartDate)
	//	edate, _ := time.Parse(layout, request.SchedulingEndDate)
	slider.SchedulingStartDate = request.SchedulingStartDate
	slider.SchedulingEndDate = request.SchedulingEndDate
	if request.GreenAreaPlaylistId != "" {
		slider.GreenAreaPlaylistId = request.GreenAreaPlaylistId
	} else {
		slider.GreenAreaPlaylistId = "00000000-0000-0000-0000-000000000000"
	}

	if request.BlackAreaPlaylistId != "" {
		slider.BlackAreaPlaylistId = request.BlackAreaPlaylistId
	} else {
		slider.BlackAreaPlaylistId = "00000000-0000-0000-0000-000000000000"
	}
	if request.RedAreaPlaylistId != "" {
		slider.RedAreaPlaylistId = request.RedAreaPlaylistId
	} else {
		slider.RedAreaPlaylistId = "00000000-0000-0000-0000-000000000000"
	}
	if sliderId != "" {
		fmt.Println("---------------------------------update---------------------------------")
		//name Validation while create new records Start
		// var slidernames Slider
		// db.Table("slider").Select("name,id").Where("name=?", request.Name).Find(&slidernames)

		//var name common.NameError
		// var isStringAlphabetic = regexp.MustCompile(`^[a-zA-Z0-9_ ]*$`).MatchString
		// if !isStringAlphabetic(request.Name) {
		// 	errorFlag = true
		// 	name = common.NameError{Code: "At least one alphanumeric character is required" + request.Name, Description: "At least one alphanumeric character is required" + request.Name}
		// } else
		// if slidernames.Name != "" {
		//	if slidernames.Id != c.Param("id") {
		//		errorFlag = true
		//		name = common.NameError{Code: "error_slider_name_not_unique", Description: "Slider with specified 'Name' of " + request.Name + " already exists."}
		//	}
		// }

		//if name.Code != "" {
		//	invalid.NameError = name
		//}
		finalErrorResponse = common.FinalErrorResponseslider{Error: "invalid_request", Description: "Validation failed.", Code: "error_validation_failed", RequestId: randstr.String(32), Invalid: invalid}
		if errorFlag {
			c.JSON(http.StatusBadRequest, finalErrorResponse)
			return
		}
		//name Validation while create new records END
		slider.ModifiedAt = time.Now()
		if err := tx.Debug().Table("slider").Where("id=?", c.Param("id")).Update(&slider).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		tx.Debug().Where("slider_id=?", sliderId).Delete(&platform)
		if request.PublishingPlatforms != nil {
			for _, platformId := range *request.PublishingPlatforms {
				platform.SliderId = sliderId
				platform.TargetPlatform = platformId
				platforms = append(platforms, platform)
			}
			err := gormbulk.BulkInsert(tx, platforms, common.BULK_INSERT_LIMIT)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		}
		tx.Debug().Where("slider_id=?", sliderId).Delete(&region)
		if request.Regions != nil {
			for _, regionId := range *request.Regions {
				region.SliderId = sliderId
				region.CountryId = regionId
				regions = append(regions, region)
			}
			err := gormbulk.BulkInsert(tx, regions, common.BULK_INSERT_LIMIT)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		}
		tx.Debug().Where("slider_id=? ", sliderId).Delete(&pageSlider)
		if len(request.PagesIds) != 0 {
			for _, pageId := range request.PagesIds {
				var pageSliderOrdercheck PageSliderOrder
				rows := db.Debug().Table("page_slider ps").Select("ps.order as order").Where("page_id=? and slider_id=?", pageId, sliderId).Find(&pageSliderOrdercheck)
				pageslidercount := int(rows.RowsAffected)
				fmt.Println("page slider count ", pageslidercount)
				if pageslidercount > 0 {
					// if update keep order same
					pageSlider.SliderId = sliderId
					pageSlider.PageId = pageId
					pageSlider.Order = pageSliderOrdercheck.Order
					pageSliders = append(pageSliders, pageSlider)
				} else {
					// if create increment order by 1
					var pageSliderOrder PageSliderOrder
					if err := db.Debug().Table("page_slider ps").Select("max(ps.order) as order").Where("page_id=?", pageId).Find(&pageSliderOrder).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
					pageSlider.SliderId = sliderId
					pageSlider.PageId = pageId
					pageSlider.Order = pageSliderOrder.Order + 1
					pageSliders = append(pageSliders, pageSlider)
				}
			}
			err := gormbulk.BulkInsert(tx, pageSliders, common.BULK_INSERT_LIMIT)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		}
		err := tx.Commit().Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		if len(request.PagesIds) != 0 {
			response := make(chan fragments.FragmentUpdate)
			go fragments.CreateSliderResponse(sliderId, "", c, response, 0, 0, 0)
			outPut := <-response
			if outPut.Err != nil {
				serverError.Description = outPut.Err.Error()
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		} else {
			fdb := c.MustGet("FDB").(*gorm.DB)
			var sliderFragment fragments.SliderFragment
			fdb.Debug().Where("slider_id=?", sliderId).Delete(&sliderFragment)
		}
		/* update dirty count in slider_sync table */
		go common.SliderSynching(c.Param("id"), c)
		res := map[string]string{
			"id": c.Param("id"),
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
		return
	}
}

// CreateSlider - Create slider
// GET /api/sliders
// @Summary Create slider details
// @Description Create slider details
// @Tags Slider
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param id path string false "Id"
// @Param body body CreateUpdateSliderRequestCreate true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /api/sliders [post]
func (hs *HandlerService) CreateSlider(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	serverError := common.ServerErrorResponse()
	db := c.MustGet("DB").(*gorm.DB)
	ctx := context.Background()
	tx := db.BeginTx(ctx, nil)
	var request CreateUpdateSliderRequestCreate
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println("Request Body Binding Error in Create Slider: ", err)
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	fmt.Println("Request Body Binding Error in Create Slider: ", request)
	var errorFlag bool
	errorFlag = false
	if request.Regions == nil || request.PublishingPlatforms == nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	/* Name Error */
	var nameUniqueCheck int
	var nameError common.NameError
	db.Debug().Select("slider").Where("name=?", request.Name).Count(&nameUniqueCheck)
	if nameUniqueCheck > 0 {
		errorFlag = true
		nameError = common.NameError{Code: "error_slider_name_not_unique", Description: "Slider with specified 'Name' of " + request.Name + " already exists."}
	}
	// if !common.NumberRegex(request.Name) {
	// 	errorFlag = true
	// 	nameError = common.NameError{Code: "error_slider_Name_invalid", Description: "At least one alphanumeric character is required(" + request.Name + ")."}
	// }
	/* Black Playlist Error */
	var blackerror common.BlackAreaPlaylistError
	// if request.BlackAreaPlaylistId != "" {
	// 	var count int
	// 	db.Table("playlist p").Where("p.id =? and p.scheduling_start_date <= ? and p.scheduling_end_date >=?", request.BlackAreaPlaylistId, request.SchedulingStartDate, request.SchedulingEndDate).Count(&count)
	// 	if count < 1 {
	// 		errorFlag = true
	// 		blackerror = common.BlackAreaPlaylistError{Code: "error_slider_balckareaplaylist_scheduling_invalid", Description: "Black Area playlist should match to slider scheduling."}
	// 	}
	// }
	// if request.BlackAreaPlaylistId == "" && request.Type == 3 {
	// 	errorFlag = true
	// 	blackerror = common.BlackAreaPlaylistError{Code: "error_slider_blackareaplaylist_not_empty", Description: "Black Area playlist should be assigned."}
	// }
	/* Green Playlist Error */
	var greenerror common.GreenAreaPlaylistError
	// if request.GreenAreaPlaylistId != "" {
	// 	var greencount int
	// 	db.Table("playlist p").Where("p.id =? and p.scheduling_start_date <= ? and p.scheduling_end_date >=?", request.GreenAreaPlaylistId, request.SchedulingStartDate, request.SchedulingEndDate).Count(&greencount)
	// 	if greencount < 1 {
	// 		errorFlag = true
	// 		greenerror = common.GreenAreaPlaylistError{Code: "error_slider_greenareaplaylist_scheduling_invalid", Description: "Green Area playlist should match to slider scheduling."}
	// 	}
	// }
	// if request.GreenAreaPlaylistId == "" && request.Type == 1 {
	// 	errorFlag = true
	// 	greenerror = common.GreenAreaPlaylistError{Code: "error_slider_greenareaplaylist_not_empty", Description: "Green Area playlist should be assigned."}
	// }
	/* Red Playlist Error */
	var rederror common.RedAreaPlaylistError
	// if request.RedAreaPlaylistId != "" {
	// 	var redcount int
	// 	db.Table("playlist p").Where("p.id =? and p.scheduling_start_date <= ? and p.scheduling_end_date >=?", request.RedAreaPlaylistId, request.SchedulingStartDate, request.SchedulingEndDate).Count(&redcount)
	// 	if redcount > 0 {
	// 		var totalCount int
	// 		db.Raw("select count(pi2.id) from playlist_item pi2 left join playlist_item_content pic on pic.playlist_item_id = pi2.id where pi2.playlist_id =?", request.RedAreaPlaylistId).Count(&totalCount)
	// 		if totalCount < 7 {
	// 			errorFlag = true
	// 			rederror = common.RedAreaPlaylistError{Code: "error_slider_redareaplaylist_invalid", Description: "Red Area should contain not less than 7 published Movies/Seasons/Series with matching scheduling."}
	// 		}
	// 	} else {
	// 		fmt.Println(redcount)
	// 		errorFlag = true
	// 		rederror = common.RedAreaPlaylistError{Code: "error_slider_redareaplaylist_invalid", Description: "Red Area playlist should match to slider scheduling."}
	// 	}
	// }
	// if request.RedAreaPlaylistId == "" && request.Type == 2 {
	// 	errorFlag = true
	// 	rederror = common.RedAreaPlaylistError{Code: "error_slider_redareaplaylist_not_empty", Description: "Red Area playlist should be assigned."}
	// }
	var pageValidation common.PagesIds
	var pageNames []string
	type PageName struct {
		EnglishTitle string `json:"english_title"`
	}
	var pageName PageName
	for _, pageid := range request.PagesIds {
		db.Debug().Table("page_slider ps").Select("p.english_title").Joins("join page p on p.id =ps.page_id").Where("ps.page_id=?", pageid).Find(&pageName)
		if pageName.EnglishTitle != "" {
			pageNames = append(pageNames, pageName.EnglishTitle)
		}
	}
	if len(pageNames) > 0 {
		fmt.Println(len(pageNames), pageNames)
		errorFlag = true
		pageValidation = common.PagesIds{Code: "error_slider_should_not_conflict_with_other_pages_sliders_platforms", Description: "Slider has platforms intersections for the next sliders and pages:" + strings.Join(pageNames, " ,")}
	}
	var invalid common.Invalidsslider
	if nameError.Code != "" {
		invalid.NameError = nameError
	}
	if blackerror.Code != "" {
		invalid.BlackAreaPlaylistError = blackerror
	}
	if rederror.Code != "" {
		invalid.RedAreaPlaylistError = rederror
	}
	if greenerror.Code != "" {
		invalid.GreenAreaPlaylistError = greenerror
	}
	if pageValidation.Code != "" {
		invalid.PagesIds = pageValidation
	}
	var finalErrorResponse common.FinalErrorResponseslider
	finalErrorResponse = common.FinalErrorResponseslider{Error: "invalid_request", Description: "Validation failed.", Code: "error_validation_failed", RequestId: randstr.String(32), Invalid: invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
	type PageSliderOrder struct {
		Order int `json:"order"`
	}
	var platform SliderTargetPlatform
	var platforms []interface{}
	var region SliderCountry
	var regions []interface{}
	var pageSlider PageSlider
	var pageSliders []interface{}
	sliderId := c.Param("id")
	var slider Slider
	slider.Name = request.Name
	restype := request.Type
	slider.Type = restype
	slider.DeletedByUserId = nil
	slider.IsDisabled = false
	//	layout := "2006-01-02T15:04:05.000Z"
	//	sdate, _ := time.Parse(layout, request.SchedulingStartDate)
	//	edate, _ := time.Parse(layout, request.SchedulingEndDate)
	slider.SchedulingStartDate = request.SchedulingStartDate
	slider.SchedulingEndDate = request.SchedulingEndDate

	if request.GreenAreaPlaylistId != "" {
		slider.GreenAreaPlaylistId = request.GreenAreaPlaylistId
	} else {
		slider.GreenAreaPlaylistId = "00000000-0000-0000-0000-000000000000"
	}

	if request.BlackAreaPlaylistId != "" {
		slider.BlackAreaPlaylistId = request.BlackAreaPlaylistId
	} else {
		slider.BlackAreaPlaylistId = "00000000-0000-0000-0000-000000000000"
	}
	if request.RedAreaPlaylistId != "" {
		slider.RedAreaPlaylistId = request.RedAreaPlaylistId
	} else {
		slider.RedAreaPlaylistId = "00000000-0000-0000-0000-000000000000"
	}

	fmt.Println("---------------------------------Insert---------------------------------")
	// //name Validation while create new records Start
	// var slidername Slider
	// db.Table("slider").Select("Name").Where("name=?", request.Name).Find(&slidername)
	// if request.Regions == nil || request.PublishingPlatforms == nil {
	// 	c.JSON(http.StatusInternalServerError, serverError)
	// 	return
	// }
	// var name common.NameError
	// var isStringAlphabetic = regexp.MustCompile(`[0-9]`).MatchString
	// if !isStringAlphabetic(request.Name) {
	// 	errorFlag = true
	// 	name = common.NameError{Code: "At least one alphanumeric character is required" + request.Name, Description: "At least one alphanumeric character is required" + request.Name}
	// } else if slidername.Name != "" {
	// 	errorFlag = true
	// 	name = common.NameError{Code: "error_slider_name_not_unique", Description: "Slider with specified 'Name' of " + request.Name + " already exists."}
	// }

	// if name.Code != "" {
	// 	invalid.NameError = name
	// }
	// finalErrorResponse = common.FinalErrorResponseslider{Error: "invalid_request", Description: "Validation failed.", Code: "error_validation_failed", RequestId: randstr.String(32), Invalid: invalid}
	// if errorFlag {
	// 	c.JSON(http.StatusBadRequest, finalErrorResponse)
	// 	return
	// }
	// //name Validation while create new records END
	type SliderKey struct {
		Key int `json:"key"`
	}
	var sliderKey SliderKey
	if err := db.Debug().Table("slider").Select("max(slider_key)+1 as key").Find(&sliderKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	slider.CreatedAt = time.Now()
	slider.ModifiedAt = time.Now()
	slider.SliderKey = request.SliderKey
	slider.Id = request.SliderId // slider id for creating old sliders with .net
	if err := tx.Debug().Create(&slider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if request.PublishingPlatforms != nil {
		for _, platformId := range *request.PublishingPlatforms {
			platform.SliderId = slider.Id
			platform.TargetPlatform = platformId
			platforms = append(platforms, platform)
		}
		err := gormbulk.BulkInsert(tx, platforms, common.BULK_INSERT_LIMIT)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	}
	if request.Regions != nil {
		for _, regionId := range *request.Regions {
			region.SliderId = slider.Id
			//region.CountryId = regionId
			region.CountryId = regionId //regionId

			regions = append(regions, region)
		}
		err := gormbulk.BulkInsert(tx, regions, common.BULK_INSERT_LIMIT)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	}
	if len(request.PagesIds) != 0 {
		for _, pageId := range request.PagesIds {
			var pageSliderOrder PageSliderOrder
			if err := db.Debug().Table("page_slider ps").Select("max(ps.order) as order").Where("page_id=?", pageId).Find(&pageSliderOrder).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			pageSlider.SliderId = slider.Id
			pageSlider.PageId = pageId
			pageSlider.Order = pageSliderOrder.Order + 1
			pageSliders = append(pageSliders, pageSlider)
		}
		err := gormbulk.BulkInsert(tx, pageSliders, common.BULK_INSERT_LIMIT)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	}
	err := tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if len(request.PagesIds) != 0 {
		response := make(chan fragments.FragmentUpdate)
		go fragments.CreateSliderResponse(sliderId, "", c, response, 0, 0, 0)
		outPut := <-response
		if outPut.Err != nil {
			serverError.Description = outPut.Err.Error()
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	}
	/* update dirty count in slider_sync table */
	go common.SliderSynching(slider.Id, c)
	result := map[string]string{
		"id": slider.Id,
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
	return
}

// GetSliderPreviewlayouts -  Get Slider Preview Layouts
// GET /api/slider/previewlayouts
// @Summary Show Slider Preview Layouts
// @Description Get Slider Preview Layouts
// @Tags Slider
// @Security Authorization
// @Accept  json
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /api/slider/previewlayouts [get]
func (hs *HandlerService) GetSliderPreviewlayouts(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var result []PreviewLayouts
	var newarr []PreviewLayouts
	if data := db.Debug().Table("preview_layout").Select("slider_type , platform,id, preview_image_key").Find(&result).Error; data != nil {
		c.JSON(http.StatusBadRequest, common.ServerErrorResponse())
		return
	}
	for _, reg := range result {
		reg.PreviewImageKey = "https://s3.ap-south-1.amazonaws.com/z5backofficecontent/slider-preview-layouts/" + reg.PreviewImageKey
		newarr = append(newarr, reg)
	}
	c.JSON(http.StatusOK, gin.H{"data": newarr})
}

// GetSliderList - Get Sliders based on search in page
// GET /api/sliders/summary
// @Summary Get Sliders based on search in page
// @Description Get Sliders based on search in page
// @Tags slider
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param searchText path string false "SearchText"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Success 200 {array} object c.JSON
// @Router /api/sliders/summary [get]
func (hs *HandlerService) GetSlidersBasedOnSearchInPage(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	serverError := common.ServerErrorResponse()
	var limit, offset int64
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["page"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	if limit == 0 {
		limit = 10
	}
	var searchText string
	if c.Request.URL.Query()["searchText"] != nil {
		searchText = c.Request.URL.Query()["searchText"][0]
	}
	rawquery := "select id,name,is_disabled from slider where deleted_by_user_id is null"
	if searchText != "" {
		rawquery += " and name like '%" + strings.Title(searchText) + "%'"
	}
	rawquery += " group by id"
	var sliderlist, totalCount []PageSummary

	if data := db.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sliderlist).Error; data != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if errCount := db.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	pages := map[string]int{
		"size":   len(totalCount),
		"offset": int(offset),
		"limit":  int(limit),
	}
	c.JSON(http.StatusOK, gin.H{"pagination": pages, "data": sliderlist})
}

// GetSliderNotifications - Get Slider Notifications
// GET /api/sliders/notifications
// @Summary Get Slider Notifications
// @Description Get Slider Notifications details
// @Tags Slider
// @Accept  json
// @Security Authorization
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /api/sliders/notifications [get]
func (hs *HandlerService) GetSliderNotifications(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	serverError := common.ServerErrorResponse()
	db := c.MustGet("DB").(*gorm.DB)
	var sliders []Slider
	if data := db.Debug().Where("is_disabled =false and deleted_by_user_id is null  and scheduling_start_date <=now() and scheduling_end_date >now()").Find(&sliders).Error; data != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var Notifications []SliderNotificationError
	for _, slider := range sliders {
		startDate := slider.SchedulingStartDate //.Format("2006-01-02") for sync
		endDate := slider.SchedulingEndDate     //.Format("2006-01-02") for sync
		if slider.BlackAreaPlaylistId != "" {
			var playlist SliderPlaylistDetails
			if err := db.Debug().Table("playlist").Where("id=?", slider.BlackAreaPlaylistId).Find(&playlist).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			countError := GetPlaylistContentsErrors(slider.BlackAreaPlaylistId, startDate, endDate, slider.Name, common.BLACK_AREA_PLAYLIST_CONUT, c, playlist)
			if countError.Message != "" {
				Notifications = append(Notifications, countError)
			}
			if countError.Message == "" {
				scheduleError := CheckPlaylistSheduleDate(slider.BlackAreaPlaylistId, startDate, endDate, slider.Name, c, playlist)
				if scheduleError.Message != "" {
					Notifications = append(Notifications, scheduleError)
				}
				disableError := CheckPlaylistDisabledOrNot(slider.BlackAreaPlaylistId, slider.Name, c, playlist)
				if disableError.Message != "" {
					Notifications = append(Notifications, disableError)
				}
			}
		}
		if slider.RedAreaPlaylistId != "" {
			var playlist SliderPlaylistDetails
			if err := db.Debug().Table("playlist").Where("id=?", slider.RedAreaPlaylistId).Find(&playlist).Error; err != nil {
				/* c.JSON(http.StatusInternalServerError, serverError)
				return */
			}
			countError := GetPlaylistContentsErrors(slider.RedAreaPlaylistId, startDate, endDate, slider.Name, common.RED_AREA_PLAYLIST_CONUT, c, playlist)
			if countError.Message != "" {
				Notifications = append(Notifications, countError)
			}
			if countError.Message == "" {
				scheduleError := CheckPlaylistSheduleDate(slider.RedAreaPlaylistId, startDate, endDate, slider.Name, c, playlist)
				if scheduleError.Message != "" {
					Notifications = append(Notifications, scheduleError)
				}
				disableError := CheckPlaylistDisabledOrNot(slider.RedAreaPlaylistId, slider.Name, c, playlist)
				if disableError.Message != "" {
					Notifications = append(Notifications, disableError)
				}
			}
		}
		if slider.GreenAreaPlaylistId != "" {
			var playlist SliderPlaylistDetails
			if err := db.Debug().Table("playlist").Where("id=?", slider.GreenAreaPlaylistId).Find(&playlist).Error; err != nil {
				/* c.JSON(http.StatusInternalServerError, serverError)
				return */
			}
			countError := GetPlaylistContentsErrors(slider.GreenAreaPlaylistId, startDate, endDate, slider.Name, common.GREEN_AREA_PLAYLIST_CONUT, c, playlist)
			if countError.Message != "" {
				Notifications = append(Notifications, countError)
			}
			if countError.Message == "" {
				scheduleError := CheckPlaylistSheduleDate(slider.GreenAreaPlaylistId, startDate, endDate, slider.Name, c, playlist)
				if scheduleError.Message != "" {
					Notifications = append(Notifications, scheduleError)
				}
				disableError := CheckPlaylistDisabledOrNot(slider.GreenAreaPlaylistId, slider.Name, c, playlist)
				if disableError.Message != "" {
					Notifications = append(Notifications, disableError)
				}
			}
		}
	}
	c.JSON(http.StatusOK, Notifications)
	return
}
func GetPlaylistContentsErrors(playlistId, startDate, endDate, sliderName string, count int, c *gin.Context, playlistDetails SliderPlaylistDetails) SliderNotificationError {
	ecount := strconv.Itoa(count)
	db := c.MustGet("DB").(*gorm.DB)
	cdb := c.MustGet("CDB").(*gorm.DB)
	var playlistContents []PlaylistContents
	var contents []PlaylistContentsCount
	var errorMessage SliderNotificationError
	var contentIds []string
	if rows := db.Debug().Table("playlist_item pi1").Select("pic.content_id,p.english_title as playlist_name").Joins("join playlist_item_content pic on pic.playlist_item_id = pi1.id join playlist p on p.id=pi1.playlist_id").Where("pi1.playlist_id =? and p.is_disabled=false", playlistId).Find(&playlistContents).RowsAffected; rows == 0 {
		return errorMessage
	}
	for _, playlist := range playlistContents {
		contentIds = append(contentIds, playlist.ContentId)
	}
	query := "SELECT DISTINCT(c.content_key) as id,c.content_type,cpi.transliterated_title  FROM content c join season s on s.content_id =c.id join episode e on e.season_id =s.id join content_primary_info cpi on cpi.id=c.primary_info_id join about_the_content_info atci on atci.id=s.about_the_content_info_id join playback_item pi1 on pi1.id =e.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id WHERE (c.id in(?) and c.status = 1 and c.deleted_by_user_id is null and (cr.digital_rights_start_date <='" + endDate + "' or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= '" + startDate + "' or cr.digital_rights_end_date is null) and s.status =1 and s.deleted_by_user_id is null and e.status =1 and e.deleted_by_user_id is null) GROUP BY c.content_key,c.content_type,cpi.transliterated_title UNION SELECT c.content_key as id,c.content_type,cpi.transliterated_title FROM content c join content_primary_info cpi on cpi.id=c.primary_info_id join about_the_content_info atci on atci.id=c.about_the_content_info_id join content_variance cv on cv.content_id =c.id join playback_item pi1 on pi1.id =cv.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id WHERE (c.id in(?) and c.status = 1 and c.deleted_by_user_id is null and (pi1.scheduling_date_time <='" + startDate + "' or pi1.scheduling_date_time is null) and (pi1.scheduling_date_time < '" + endDate + "' or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <='" + endDate + "' or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= '" + startDate + "' or cr.digital_rights_end_date is null)) GROUP BY c.content_key,c.content_type,cpi.transliterated_title"
	if rows := cdb.Debug().Raw(query, contentIds, contentIds).Find(&contents).RowsAffected; int(rows) < count {
		fmt.Println("rows", rows)
		fmt.Println("count", count)
		acount := strconv.Itoa(int(rows))
		replacer := strings.NewReplacer("#stitle#", sliderName, "#ptitle#", playlistDetails.EnglishTitle, "#ecount#", ecount, "#acount#", acount)
		output := replacer.Replace(common.SNP_ITEM_COUNT_MISMATCH_MESSAGE)
		errorMessage.Message = output
		errorMessage.Code = common.SNP_ITEM_COUNT_MISMATCH_CODE
		errorMessage.Type = 2
		return errorMessage
	}
	return errorMessage
}
func CheckPlaylistSheduleDate(playlistId, startDate, endDate, sliderName string, c *gin.Context, playlistDetails SliderPlaylistDetails) SliderNotificationError {
	db := c.MustGet("DB").(*gorm.DB)
	var playlist SliderPlaylistDetails
	var errorMessage SliderNotificationError
	if rows := db.Debug().Table("playlist").Where("id =? and (scheduling_start_date <= ? or scheduling_start_date is null) and (scheduling_end_date >= ? or scheduling_end_date is null)", playlistId, startDate, endDate).Find(&playlist).RowsAffected; rows == 0 {
		pstartDate := playlistDetails.SchedulingStartDate.Format("2006-01-02")
		pendDate := playlistDetails.SchedulingEndDate.Format("2006-01-02")
		replacer := strings.NewReplacer("#stitle#", sliderName, "#sssdate#", startDate, "#ssedate#", endDate, "#ptitle#", playlistDetails.EnglishTitle, "#pssdate#", pstartDate, "#psedate#", pendDate)
		output := replacer.Replace(common.SNP_SCHEDULING_MISMATCH_MESSAGE)
		errorMessage.Message = output
		errorMessage.Code = common.SNP_SCHEDULING_MISMATCH_CODE
		errorMessage.Type = 2
		return errorMessage
	}
	return errorMessage
}
func CheckPlaylistDisabledOrNot(playlistId, sliderName string, c *gin.Context, playlistDetails SliderPlaylistDetails) SliderNotificationError {
	db := c.MustGet("DB").(*gorm.DB)
	var playlist SliderPlaylistDetails
	var errorMessage SliderNotificationError
	if rows := db.Debug().Table("playlist").Where("id =? and is_disabled=false", playlistId).Find(&playlist).RowsAffected; rows == 0 {
		replacer := strings.NewReplacer("#stitle#", sliderName, "#ptitle#", playlistDetails.EnglishTitle)
		output := replacer.Replace(common.SNP_UNPUBLISHED_MESSAGE)
		errorMessage.Message = output
		errorMessage.Code = common.SNP_UNPUBLISHED_CODE
		errorMessage.Type = 2
		return errorMessage
	}
	return errorMessage
}

// GetallRegionsBasedOnSliderId -  Get all Regions Based On Slider Id
// GET /api/sliders/:id/region
// @Summary   Get all Regions Based On Slider Id
// @Description   Get all Regions Based On Slider Id
// @Tags Slider
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param id path string true "id"
// @Success 200 {array} object c.JSON
// @Router /api/sliders/{id}/region [get]
func (hs *HandlerService) GetallRegionsBasedOnSliderId(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	id := c.Param("id")
	var errorresponse = common.ServerErrorResponse()
	var sliderRegions []Country
	/*Fetch slider regions*/
	if resultError := db.Debug().Table("slider_country sc").Select("c.english_name as name,c.id").
		Joins("left join country c  on c.id =sc.country_id").
		Where("sc.slider_id=?", id).Find(&sliderRegions).Error; resultError != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
	}
	c.JSON(http.StatusOK, gin.H{"data": sliderRegions})
}
