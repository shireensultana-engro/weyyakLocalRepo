package page

import (
	"bytes"
	"context"
	"fmt"
	"frontend_config/common"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	qrg.DELETE("/pages/:id", hs.DeletePageDetailsById)
	qrg.POST("/pages/:id/available", hs.PageAvailabilityUpdate)
	qrg.GET("/pages", hs.GetAllPagesListBySearchText)
	qrg.GET("/pages/ordered", hs.PageOrder)
	qrg.POST("/pages/ordered", hs.UpdatePageOrder)
	qrg.GET("/pages/summary", hs.GetPagesListBySearchText)
	qrg.POST("/pages", hs.CreateOrUpdatePageDetails)
	qrg.POST("/pages/:id", hs.CreateOrUpdatePageDetails)
	qrg.GET("/pages/:id", hs.PageDetailsbyPageid)
	qrg.GET("/pages/:id/region", hs.GetallPageRegionsByPageId)
}

// Get page order details -  fetches page order details
// GET /pages/ordered
// @Summary Show page order details
// @Description get page order details
// @Tags Frontend
// @Security Authorization
// @Accept  json
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /pages/ordered [get]
func (hs *HandlerService) PageOrder(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	serverError := common.ServerErrorResponse()
	db := c.MustGet("DB").(*gorm.DB)
	var publishPlatformDetails []PublishPlatformDetails
	var platformdetails []PlatformDetails
	if err := db.Debug().Table("publish_platform").Select("id").Find(&publishPlatformDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	data := make(map[int]interface{})
	final := make(map[string]interface{})
	for _, value := range publishPlatformDetails {
		fields := "id,english_title as englishTitle,ptp.page_order_number as pageOrderNumber,case when page_type = 1 then true else false end as isHome"
		if err := db.Debug().Table("page p").Select(fields).Joins("join page_target_platform ptp on p.id=ptp.page_id").
			Where("ptp.target_platform=? and p.deleted_by_user_id is null", value.Id).
			Order("ptp.target_platform,ptp.page_order_number asc").Find(&platformdetails).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		data[value.Id] = platformdetails
	}
	final["publishingPlatformOrderedDetails"] = data
	c.JSON(http.StatusOK, gin.H{"data": final})
	return
}

// Get page details by pageid -  fetches page details by pageid
// GET /api/pages/{id}
// @Summary Show page details by pageid
// @Description get page details by pageid
// @Tags Page
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param id path string true "Id"
// @Success 200 {array} object c.JSON
// @Router /api/pages/{id} [get]
func (hs *HandlerService) PageDetailsbyPageid(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var pagedetails PageDetails
	var sliderid []Sliderid
	var sliderarr Sliderdetails
	var sliderdefaultarr Sliderdetails
	res := make(map[string]interface{})
	var errorresponse = common.ServerErrorResponse()

	/*sigle query to fetch all page related things TODO some changes */
	fields := "p.english_title,p.arabic_title,p.page_order_number,p.english_page_friendly_url,p.arabic_page_friendly_url,p.english_meta_title,p.arabic_meta_title,p.english_meta_description,p.arabic_meta_description,json_agg(distinct ptp.target_platform)::varchar as platforms,string_agg(distinct concat(ptp.target_platform::varchar,':',ptp.page_order_number::int), ',') as platform_order ,json_agg(distinct pc.country_id)::varchar as regions,case when p.page_type = 1 then true else false end as isHome, p.is_disabled,case when p.has_mobile_menu = true then 'https://z5content-uat.s3.amazonaws.com/" + c.Param("id") + "/mobile-menu' else '' end as mobileMenu,case when p.has_menu_poster_image = true then 'https://z5content-uat.s3.amazonaws.com/" + c.Param("id") + "/menu-poster-image' else '' end as menuPosterImage ,case when p.has_mobile_menu_poster_image = true then 'https://z5content-uat.s3.amazonaws.com/" + c.Param("id") + "/mobile-menu-poster-image' else '' end as mobileMenuPosterImage"
	joins := "left join page_target_platform ptp on ptp.page_id =p.id left join page_country pc on pc.page_id =p.id"
	group := "p.english_title,p.arabic_title,p.page_order_number,p.english_page_friendly_url,p.arabic_page_friendly_url,p.english_meta_title,p.arabic_meta_title,p.english_meta_description,p.page_type,p.is_disabled,p.arabic_meta_description,p.has_mobile_menu,p.has_menu_poster_image ,p.has_mobile_menu_poster_image"
	if pagedetailsdata := db.Debug().Table("page p").Select(fields).Joins(joins).Where("id=? and deleted_by_user_id is null", c.Param("id")).Group(group).Find(&pagedetails).Error; pagedetailsdata != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	/*platforms */
	if pagedetails.Platforms == "[null]" {
		buffer := make([]int, 0)
		res["publishingPlatforms"] = buffer
	} else {
		Platforms, _ := common.JsonStringToIntSliceOrMap(pagedetails.Platforms)
		res["publishingPlatforms"] = Platforms
	}
	/* Regions */
	if pagedetails.Regions == "[null]" {
		buffer := make([]int, 0)
		res["regions"] = buffer
	} else {
		pageRegions, _ := common.JsonStringToIntSliceOrMap(pagedetails.Regions)
		res["regions"] = pageRegions
	}
	var dataqqq []string
	data11 := map[string]int{}
	if pagedetails.PlatformOrder == ":" {
		data11 = map[string]int{}
	} else {
		dataqqq = strings.Split(pagedetails.PlatformOrder, ",")
		for _, aa := range dataqqq {
			da := strings.Split(aa, ":")
			key := da[0]
			value, _ := strconv.Atoi(da[1])
			data11[key] = value
		}
	}
	res["englishTitle"] = pagedetails.EnglishTitle
	res["arabicTitle"] = pagedetails.ArabicTitle
	res["pageOrderNumber"] = pagedetails.PageOrderNumber
	res["englishPageFriendlyUrl"] = pagedetails.EnglishPageFriendlyUrl
	res["arabicPageFriendlyUrl"] = pagedetails.ArabicPageFriendlyUrl
	res["englishMetaTitle"] = pagedetails.EnglishMetaTitle
	res["arabicMetaTitle"] = pagedetails.ArabicMetaTitle
	res["englishMetaDescription"] = pagedetails.EnglishMetaDescription
	res["arabicMetaDescription"] = pagedetails.ArabicMetaDescription
	res["isHome"] = pagedetails.Ishome
	res["isDisabled"] = pagedetails.IsDisabled

	var playlist, multiPlaylist []Playlistdetails
	var singleplaylist Playlistdetails
	if playlistdata := db.Debug().Table("page_playlist pp").Select("p.english_title,p.arabic_title,p.is_disabled,p.id,json_agg(distinct plp.target_platform)::varchar as platforms").Where("p.deleted_by_user_id is null and pp.page_id =?", c.Param("id")).
		Joins("join playlist p on p.id = pp.playlist_id left join play_list_platform plp on plp.play_list_id = p.id").
		Group("p.english_title,p.arabic_title,p.is_disabled,p.id,pp.order").
		Order("pp.order").
		Find(&playlist).Error; playlistdata != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	for _, playlistInfo := range playlist {
		singleplaylist.EnglishTitle = playlistInfo.EnglishTitle
		singleplaylist.ArabicTitle = playlistInfo.ArabicTitle
		singleplaylist.IsDisabled = playlistInfo.IsDisabled
		//platforms, _ := common.JsonStringToIntSliceOrMap(playlistInfo.Platforms) /*no need to fetch platforms here if incase of any use with this platforms make sure unblock this comment */
		singleplaylist.PublishingPlatforms = nil
		singleplaylist.Id = playlistInfo.Id
		multiPlaylist = append(multiPlaylist, singleplaylist)
	}
	/* Fetch Playlist Details */
	res["playlists"] = multiPlaylist
	if pagesliderdata := db.Debug().Table("page_slider").Select("slider_id,page_slider.order").Where("page_id=?", c.Param("id")).Find(&sliderid).Error; pagesliderdata != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	var sliderfinal []Sliderdetails
	for _, sliderange := range sliderid {
		fmt.Println(sliderange.SliderId)
		if sliderange.Order == 0 {
			rows := db.Debug().Table("slider").Select("name,is_disabled,id").Where("id=? and deleted_by_user_id is null", sliderange.SliderId).Find(&sliderdefaultarr)
			fmt.Println("lllllllllll", sliderdefaultarr)
			count := int(rows.RowsAffected)
			if count < 1 {
				sliderdefaultarr.Id = ""
				sliderdefaultarr.IsDisabled = nil
				sliderdefaultarr.Name = ""
			}
		} else {
			var count int
			db.Debug().Table("slider").Select("name,is_disabled,id").Where("id=? and deleted_by_user_id is null ", sliderange.SliderId).Find(&sliderarr).Count(&count)
			if sliderarr.Name != "" && sliderarr.Id != "" {
				sliderfinal = append(sliderfinal, sliderarr)
			}
		}
	}
	res["sliders"] = sliderfinal
	res["defaultSlider"] = sliderdefaultarr //sliderdefaultfinal //reason of change expecting in get page details of direct object not array of object
	res["publishingPlatformsOrder"] = data11
	/* images */
	textfields := make(map[string]string)
	textfields["mobileMenu"] = pagedetails.Mobilemenu
	textfields["menuPosterImage"] = pagedetails.Menuposterimage
	textfields["mobileMenuPosterImage"] = pagedetails.Mobilemenuposterimage

	res["nonTextualData"] = textfields
	res["id"] = c.Param("id")
	c.JSON(http.StatusOK, gin.H{"data": res})

}

// DeletePageDetailsById - Delete page details by page id
// DELETE /api/pages/:id
// @Summary Delete page
// @Description delete Delete Page
// @Tags Page
// @Accept json
// @Security Authorization
// @Produce json
// @Param id path string true "Id"
// @Success 200 {array} object c.JSON
// @Router /api/pages/{id} [delete]
func (hs *HandlerService) DeletePageDetailsById(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	userid := c.MustGet("userid")
	serverError := common.ServerErrorResponse()
	notFound := common.NotFoundErrorResponse()
	var totalCount int
	if err := db.Debug().Table("page").Where("id=? and deleted_by_user_id is null", c.Param("id")).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if totalCount < 1 {
		c.JSON(http.StatusNotFound, notFound)
		return
	} else {
		if err := db.Debug().Table("page").Where("id=?", c.Param("id")).Updates(map[string]interface{}{"deleted_by_user_id": userid, "modified_at": time.Now()}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
	}
}

// PageAvailabilityUpdate - Disable or enable based on page id
// POST /api/pages/:id/available
// @Summary Show disable or enable page
// @Description post disable or enable page
// @Tags Page
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param id path string true "Id"
// @Param body body PageAvailability true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /api/pages/{id}/available [post]
func (hs *HandlerService) PageAvailabilityUpdate(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var pageUpdate PageAvailability
	if err := c.ShouldBindJSON(&pageUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": http.StatusBadRequest})
		return
	}
	pageIdUrl := c.Param("id")
	if pageIdUrl != pageUpdate.Id {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide valid page Id.", "Status": http.StatusBadRequest})
		return
	}
	/*checking that the page exist in database-Table or not*/
	var count int
	if pageDetails := db.Debug().Table("page").Select("id").Where("id=? and deleted_by_user_id is null", c.Param("id")).Count(&count).Error; pageDetails != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	if !(count < 1) {
		// if err := db.Debug().Table("page").Where("id=?", c.Param("id")).Update("is_disabled", pageUpdate.IsDisabled).Error; err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		// 	return
		// }
		if err := db.Debug().Exec("UPDATE page set is_disabled = ? ,modified_at = ? where id = ?", pageUpdate.IsDisabled, time.Now(), c.Param("id")).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
			return
		}
		/* update dirty count in page_sync table */
		type PageDetails struct {
			Id      string
			PageKey int
		}
		var pageDetails PageDetails
		db.Debug().Raw("select p.id,p.page_key from page p where p.id=?", c.Param("id")).Find(&pageDetails)
		go common.PageSynching(pageDetails.Id, pageDetails.PageKey, c)
		c.JSON(http.StatusOK, gin.H{})
	} else {
		id := Id{"error_page_not_found", "The specified condition was not met for 'Id'."}
		invalid := Invalid{Id: &id}
		finalErrorResponse := FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
		c.JSON(http.StatusBadRequest, finalErrorResponse)
	}
}

// GetPagesListBySearchText -Get Pages List By Search Text
// GET /pages/summary
// @Summary Get Pages List By Search Text
// @Description Get Pages List By Search Text
// @Tags Page
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param searchText path string false "Search Text"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Success 200 {array} object c.JSON
// @Router /pages/summary [get]
func (hs *HandlerService) GetPagesListBySearchText(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var limit, offset int64

	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	if limit == 0 {
		limit = 10
	}
	var searchText string
	var pagelist, pagelists, totalCount []PagelistSummary

	if c.Request.URL.Query()["searchText"] != nil {
		searchText = strings.ToLower(c.Request.URL.Query()["searchText"][0])
	}
	rawquery := "select p.id,p.english_title,p.arabic_title,p.is_disabled,case when rtrim(substring( string_agg( p.page_type ::varchar, ',' order by p.page_type desc ) from '(?:[^,]+,){1,1}'), ',') = '1'  then 1 else 0 end  as is_home from page p where p.deleted_by_user_id is null"
	if searchText != "" {
		rawquery += " and ( lower(p.english_title) like '%" + searchText + "%' OR  lower(p.arabic_title) like '%" + searchText + "%' )"
	}
	rawquery += "  group by p.id order by p.created_at desc  "
	if data := db.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&pagelist).Error; data != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	var platformDetails []PlanNames
	for _, pltid := range pagelist {
		rawpidlist := "select ptp.target_platform from page p join page_target_platform ptp on p.id = ptp.page_id where p.id = '" + pltid.ID + "' "
		if dataId := db.Debug().Raw(rawpidlist).Scan(&platformDetails).Error; dataId != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
			return
		}
		var Ids []int
		for _, idarr := range platformDetails {
			Ids = append(Ids, idarr.TargetPlatform)
		}
		pltid.PublishingPlatforms = Ids
		pagelists = append(pagelists, pltid)
	}
	if errCount := db.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	pages := map[string]int{
		"size":   len(totalCount),
		"offset": int(offset),
		"limit":  int(limit),
	}
	c.JSON(http.StatusOK, gin.H{"pagination": pages, "data": pagelists})
}

// GetAllPagesListBySearchText -Get all pages list and pages by search text
// GET /api/pages
// @Summary Get all pages list and pages by search text
// @Description Get all pages list and pages by search text
// @Tags Page
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param searchText path string false "Search Text"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Success 200 {array} object c.JSON
// @Router /api/pages [get]
func (hs *HandlerService) GetAllPagesListBySearchText(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
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
	var pagelist, totalCount []PageDetailsSummary
	pagelists := []PageDetailsSummary{}
	if c.Request.URL.Query()["searchText"] != nil {
		searchText = strings.ToLower(c.Request.URL.Query()["searchText"][0])
	}
	rawquery := "select p.page_order_number,p.id,p.english_title,p.arabic_title,p.is_disabled,case when rtrim(substring( string_agg( p.page_type ::varchar, ',' order by p.page_type desc ) from '(?:[^,]+,){1,1}'), ',') = '1'  then 1 else 0 end  as isome_removed, case when page_type='1' then 1 else 0 end as is_home from page p where p.deleted_by_user_id is null"
	if searchText != "" {
		rawquery += " and ( lower(p.english_title) like '%" + searchText + "%' OR  lower(p.arabic_title) like '%" + searchText + "%' )"
	}
	rawquery += "  group by p.id order by is_home desc, p.modified_at desc  "
	if data := db.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&pagelist).Error; data != nil {
		c.JSON(http.StatusInternalServerError, FinalErrorResponse{Error: "server_error", Description: "Server error.", Code: "error_server_error", RequestId: randstr.String(32)})
		return
	}
	var platformDetails []PlanNames
	for _, pltid := range pagelist {
		//collect Target Platforms
		rawpidlist := "select ptp.target_platform from page p join page_target_platform ptp on p.id = ptp.page_id where p.id = '" + pltid.ID + "' "
		if dataId := db.Debug().Raw(rawpidlist).Scan(&platformDetails).Error; dataId != nil {
			c.JSON(http.StatusInternalServerError, FinalErrorResponse{Error: "server_error", Description: "Server error.", Code: "error_server_error", RequestId: randstr.String(32)})
			return
		}
		var Ids []int
		for _, idarr := range platformDetails {
			Ids = append(Ids, idarr.TargetPlatform)
		}
		if len(Ids) < 1 {
			pltid.HasPublishingPlatforms = true
			buffer := make([]int, 0)
			pltid.PublishingPlatforms = buffer
		} else {
			pltid.HasPublishingPlatforms = true
			pltid.PublishingPlatforms = Ids
		}
		//collect Regions
		var details PageRegionDetails
		rawRegionlist := "select case when (select count(*) as countrys_count from page p join page_country pc on pc.page_id = p.id join country c on c.id = pc.country_id where p.id = '" + pltid.ID + "') = 248 then 'All Regions' else (select string_agg(c.english_name,',') from page p join page_country pc on pc.page_id = p.id join country c on c.id = pc.country_id where p.id = '" + pltid.ID + "') end as details"
		if regionresult := db.Debug().Raw(rawRegionlist).Scan(&details).Error; regionresult != nil {
			c.JSON(http.StatusInternalServerError, FinalErrorResponse{Error: "server_error", Description: "Server error.", Code: "error_server_error", RequestId: randstr.String(32)})
			return
		}
		if details.Details == "All Regions" {
			pltid.Region = "All Regions"
			pltid.HasMoreRegions = false
		} else {
			var regions []string
			array := strings.Split(details.Details, ",")
			for _, val := range array {
				regions = append(regions, val)
			}
			if len(regions) >= 5 {
				output := regions[:5]
				pltid.Region = strings.Join(output, ",")
				pltid.HasMoreRegions = true
			} else {
				pltid.Region = details.Details
				pltid.HasMoreRegions = false
			}
		}
		//check page_id having playlist&slider items or not
		var hasplaylist HasPlaylist
		rawPlaylist := "select distinct case when pp.page_id is null then false else true end as has_playlist_id,case when ps.page_id is null then false else true end as has_slider_id from page p left join page_playlist pp on pp.page_id =p.id left join page_slider ps on ps.page_id =p.id where p.id ='" + pltid.ID + "'"
		db.Debug().Raw(rawPlaylist).Scan(&hasplaylist)
		pltid.HasPlaylists = hasplaylist.HasPlaylistId
		pltid.HasSliders = hasplaylist.HasSliderId
		pagelists = append(pagelists, pltid)

	}
	if errCount := db.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
		c.JSON(http.StatusInternalServerError, FinalErrorResponse{Error: "server_error", Description: "Server error.", Code: "error_server_error", RequestId: randstr.String(32)})
		return
	}
	pages := map[string]int{
		"size":   len(totalCount),
		"offset": int(offset),
		"limit":  int(limit),
	}
	c.JSON(http.StatusOK, gin.H{"pagination": pages, "data": pagelists})
}

// Update page order details -  Update page order details
// POST /api/pages/ordered
// @Summary  Update page order details
// @Description post  Update page order details
// @Tags Page
// @Accept  json
// @Security Authorization
// @Produce  json
// @Security Authorization
// @Param body body FinalResponse true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /api/pages/ordered [post]
func (hs *HandlerService) UpdatePageOrder(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var errorresponse = common.ServerErrorResponse()
	var response FinalResponse
	ctx := context.Background()
	tx := db.BeginTx(ctx, nil)
	if errresult := c.ShouldBindJSON(&response); errresult != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errresult.Error(), "status": http.StatusBadRequest})
		return
	}
	var pagedetailsresponse []PageTargetPlatform
	// if deleteresult := tx.Debug().Table("page_target_platform").Delete(&pagedetailsresponse).Error; deleteresult != nil {
	// 	tx.Rollback()
	// 	c.JSON(http.StatusInternalServerError, errorresponse)
	// 	return
	// }
	var pagedata0 []interface{}
	for _, k0 := range response.PublishingPlatformOrderedDetails.TargetPlatform0 {
		if deleteresult := tx.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", k0.PageId, 0).Delete(&pagedetailsresponse).Error; deleteresult != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if k0.EnglishTitle == "Home" && k0.PageOrderNumber != 0 {
			pageresponse := PageTargetPlatform{PageId: k0.PageId, TargetPlatform: 0, PageOrderNumber: 0}
			pagedata0 = append(pagedata0, pageresponse)
		} else if k0.EnglishTitle != "Home" && k0.PageOrderNumber == 0 {
			pageresponse := PageTargetPlatform{PageId: k0.PageId, TargetPlatform: 0, PageOrderNumber: 1}
			pagedata0 = append(pagedata0, pageresponse)
		} else {
			pageresponse := PageTargetPlatform{PageId: k0.PageId, TargetPlatform: 0, PageOrderNumber: k0.PageOrderNumber}
			pagedata0 = append(pagedata0, pageresponse)
		}
	}
	if err := gormbulk.BulkInsert(tx, pagedata0, common.BULK_INSERT_LIMIT); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	var pagedata1 []interface{}
	for _, k1 := range response.PublishingPlatformOrderedDetails.TargetPlatform1 {
		if deleteresult := tx.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", k1.PageId, 1).Delete(&pagedetailsresponse).Error; deleteresult != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if k1.EnglishTitle == "Home" && k1.PageOrderNumber != 0 {
			pageresponse := PageTargetPlatform{PageId: k1.PageId, TargetPlatform: 1, PageOrderNumber: 0}
			pagedata1 = append(pagedata1, pageresponse)
		} else if k1.EnglishTitle != "Home" && k1.PageOrderNumber == 0 {
			pageresponse := PageTargetPlatform{PageId: k1.PageId, TargetPlatform: 1, PageOrderNumber: 1}
			pagedata1 = append(pagedata1, pageresponse)
		} else {
			pageresponse := PageTargetPlatform{PageId: k1.PageId, TargetPlatform: 1, PageOrderNumber: k1.PageOrderNumber}
			pagedata1 = append(pagedata1, pageresponse)
		}
	}
	if err := gormbulk.BulkInsert(tx, pagedata1, common.BULK_INSERT_LIMIT); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	var pagedata2 []interface{}
	for _, k2 := range response.PublishingPlatformOrderedDetails.TargetPlatform2 {
		if deleteresult := tx.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", k2.PageId, 2).Delete(&pagedetailsresponse).Error; deleteresult != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if k2.EnglishTitle == "Home" && k2.PageOrderNumber != 0 {
			pageresponse := PageTargetPlatform{PageId: k2.PageId, TargetPlatform: 2, PageOrderNumber: 0}
			pagedata2 = append(pagedata2, pageresponse)
		} else if k2.EnglishTitle != "Home" && k2.PageOrderNumber == 0 {
			pageresponse := PageTargetPlatform{PageId: k2.PageId, TargetPlatform: 2, PageOrderNumber: 1}
			pagedata2 = append(pagedata2, pageresponse)
		} else {
			pageresponse := PageTargetPlatform{PageId: k2.PageId, TargetPlatform: 2, PageOrderNumber: k2.PageOrderNumber}
			pagedata2 = append(pagedata2, pageresponse)
		}
	}
	if err := gormbulk.BulkInsert(tx, pagedata2, common.BULK_INSERT_LIMIT); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	var pagedata3 []interface{}
	for _, k3 := range response.PublishingPlatformOrderedDetails.TargetPlatform3 {
		if deleteresult := tx.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", k3.PageId, 3).Delete(&pagedetailsresponse).Error; deleteresult != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if k3.EnglishTitle == "Home" && k3.PageOrderNumber != 0 {
			pageresponse := PageTargetPlatform{PageId: k3.PageId, TargetPlatform: 3, PageOrderNumber: 0}
			pagedata3 = append(pagedata3, pageresponse)
		} else if k3.EnglishTitle != "Home" && k3.PageOrderNumber == 0 {
			pageresponse := PageTargetPlatform{PageId: k3.PageId, TargetPlatform: 3, PageOrderNumber: 1}
			pagedata3 = append(pagedata3, pageresponse)
		} else {
			pageresponse := PageTargetPlatform{PageId: k3.PageId, TargetPlatform: 3, PageOrderNumber: k3.PageOrderNumber}
			pagedata3 = append(pagedata3, pageresponse)
		}
	}
	if err := gormbulk.BulkInsert(tx, pagedata3, common.BULK_INSERT_LIMIT); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	var pagedata4 []interface{}
	for _, k4 := range response.PublishingPlatformOrderedDetails.TargetPlatform4 {
		if deleteresult := tx.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", k4.PageId, 4).Delete(&pagedetailsresponse).Error; deleteresult != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if k4.EnglishTitle == "Home" && k4.PageOrderNumber != 0 {
			pageresponse := PageTargetPlatform{PageId: k4.PageId, TargetPlatform: 4, PageOrderNumber: 0}
			pagedata4 = append(pagedata4, pageresponse)
		} else if k4.EnglishTitle != "Home" && k4.PageOrderNumber == 0 {
			pageresponse := PageTargetPlatform{PageId: k4.PageId, TargetPlatform: 4, PageOrderNumber: 1}
			pagedata4 = append(pagedata4, pageresponse)
		} else {
			pageresponse := PageTargetPlatform{PageId: k4.PageId, TargetPlatform: 4, PageOrderNumber: k4.PageOrderNumber}
			pagedata4 = append(pagedata4, pageresponse)
		}
	}
	if err := gormbulk.BulkInsert(tx, pagedata4, common.BULK_INSERT_LIMIT); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	var pagedata5 []interface{}
	for _, k5 := range response.PublishingPlatformOrderedDetails.TargetPlatform5 {
		if deleteresult := tx.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", k5.PageId, 5).Delete(&pagedetailsresponse).Error; deleteresult != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if k5.EnglishTitle == "Home" && k5.PageOrderNumber != 0 {
			pageresponse := PageTargetPlatform{PageId: k5.PageId, TargetPlatform: 5, PageOrderNumber: 0}
			pagedata5 = append(pagedata5, pageresponse)
		} else if k5.EnglishTitle != "Home" && k5.PageOrderNumber == 0 {
			pageresponse := PageTargetPlatform{PageId: k5.PageId, TargetPlatform: 5, PageOrderNumber: 1}
			pagedata5 = append(pagedata5, pageresponse)
		} else {
			pageresponse := PageTargetPlatform{PageId: k5.PageId, TargetPlatform: 5, PageOrderNumber: k5.PageOrderNumber}
			pagedata5 = append(pagedata5, pageresponse)
		}
	}
	if err := gormbulk.BulkInsert(tx, pagedata5, common.BULK_INSERT_LIMIT); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	var pagedata6 []interface{}
	for _, k6 := range response.PublishingPlatformOrderedDetails.TargetPlatform6 {
		if deleteresult := tx.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", k6.PageId, 6).Delete(&pagedetailsresponse).Error; deleteresult != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if k6.EnglishTitle == "Home" && k6.PageOrderNumber != 0 {
			pageresponse := PageTargetPlatform{PageId: k6.PageId, TargetPlatform: 6, PageOrderNumber: 0}
			pagedata6 = append(pagedata6, pageresponse)
		} else if k6.EnglishTitle != "Home" && k6.PageOrderNumber == 0 {
			pageresponse := PageTargetPlatform{PageId: k6.PageId, TargetPlatform: 6, PageOrderNumber: 1}
			pagedata6 = append(pagedata6, pageresponse)
		} else {
			pageresponse := PageTargetPlatform{PageId: k6.PageId, TargetPlatform: 6, PageOrderNumber: k6.PageOrderNumber}
			pagedata6 = append(pagedata6, pageresponse)
		}
	}
	if err := gormbulk.BulkInsert(tx, pagedata6, common.BULK_INSERT_LIMIT); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	var pagedata7 []interface{}
	for _, k7 := range response.PublishingPlatformOrderedDetails.TargetPlatform7 {
		if deleteresult := tx.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", k7.PageId, 7).Delete(&pagedetailsresponse).Error; deleteresult != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if k7.EnglishTitle == "Home" && k7.PageOrderNumber != 0 {
			pageresponse := PageTargetPlatform{PageId: k7.PageId, TargetPlatform: 7, PageOrderNumber: 0}
			pagedata7 = append(pagedata7, pageresponse)
		} else if k7.EnglishTitle != "Home" && k7.PageOrderNumber == 0 {
			pageresponse := PageTargetPlatform{PageId: k7.PageId, TargetPlatform: 7, PageOrderNumber: 1}
			pagedata7 = append(pagedata7, pageresponse)
		} else {
			pageresponse := PageTargetPlatform{PageId: k7.PageId, TargetPlatform: 7, PageOrderNumber: k7.PageOrderNumber}
			pagedata7 = append(pagedata7, pageresponse)
		}
	}
	if err := gormbulk.BulkInsert(tx, pagedata7, common.BULK_INSERT_LIMIT); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	var pagedata9 []interface{}
	for _, k9 := range response.PublishingPlatformOrderedDetails.TargetPlatform9 {
		if deleteresult := tx.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", k9.PageId, 9).Delete(&pagedetailsresponse).Error; deleteresult != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if k9.EnglishTitle == "Home" && k9.PageOrderNumber != 0 {
			pageresponse := PageTargetPlatform{PageId: k9.PageId, TargetPlatform: 9, PageOrderNumber: 0}
			pagedata9 = append(pagedata9, pageresponse)
		} else if k9.EnglishTitle != "Home" && k9.PageOrderNumber == 0 {
			pageresponse := PageTargetPlatform{PageId: k9.PageId, TargetPlatform: 9, PageOrderNumber: 1}
			pagedata9 = append(pagedata9, pageresponse)
		} else {
			pageresponse := PageTargetPlatform{PageId: k9.PageId, TargetPlatform: 9, PageOrderNumber: k9.PageOrderNumber}
			pagedata9 = append(pagedata9, pageresponse)
		}
	}
	if err := gormbulk.BulkInsert(tx, pagedata9, common.BULK_INSERT_LIMIT); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	var pagedata10 []interface{}
	for _, k10 := range response.PublishingPlatformOrderedDetails.TargetPlatform10 {
		if deleteresult := tx.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", k10.PageId, 10).Delete(&pagedetailsresponse).Error; deleteresult != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if k10.EnglishTitle == "Home" && k10.PageOrderNumber != 0 {
			pageresponse := PageTargetPlatform{PageId: k10.PageId, TargetPlatform: 10, PageOrderNumber: 0}
			pagedata10 = append(pagedata10, pageresponse)
		} else if k10.EnglishTitle != "Home" && k10.PageOrderNumber == 0 {
			pageresponse := PageTargetPlatform{PageId: k10.PageId, TargetPlatform: 10, PageOrderNumber: 1}
			pagedata10 = append(pagedata10, pageresponse)
		} else {
			pageresponse := PageTargetPlatform{PageId: k10.PageId, TargetPlatform: 10, PageOrderNumber: k10.PageOrderNumber}
			pagedata10 = append(pagedata10, pageresponse)
		}
	}
	if err := gormbulk.BulkInsert(tx, pagedata10, common.BULK_INSERT_LIMIT); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
}

// Create Or Update Page Details -  Create Or Update Page Details
// POST /api/pages/{id}
// @Summary  Create Or Update Page Details
// @Description post  Create Or Update Page Details
// @Tags Page
// @Accept  json
// @Produce  json
// @security Authorization
// @Param id path string true "Id"
// @Param body body PageRequest true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /api/pages/{id} [post]
func (hs *HandlerService) CreateOrUpdatePageDetails(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	fmt.Println("Entering into create or update page")
	db := c.MustGet("DB").(*gorm.DB)
	ctx := context.Background()
	tx := db.BeginTx(ctx, nil)
	errorresponse := common.ServerErrorResponse()
	var pageRequest PageRequest
	var errorFlag bool
	errorFlag = false
	c.ShouldBindJSON(&pageRequest)
	var englisharabicnames []Page
	db.Debug().Table("page").Select("english_title,arabic_title,id").Where("english_title=? or arabic_title=?", pageRequest.EnglishTitle, pageRequest.ArabicTitle).Find(&englisharabicnames)
	var englishTitleError EnglishTitleError
	var arabicTitleError ArabicTitleError
	for _, englisharabicname := range englisharabicnames {
		if c.Param("id") == "" {
			if englisharabicname.EnglishTitle == pageRequest.EnglishTitle {
				errorFlag = true
				englishTitleError = EnglishTitleError{"error_page_english_title_not_unique", "Page with specified 'English Title' of " + pageRequest.EnglishTitle + " already exists."}
			}
			if englisharabicname.ArabicTitle == pageRequest.ArabicTitle {
				errorFlag = true
				arabicTitleError = ArabicTitleError{"error_page_arabic_title_not_unique", "Page with specified 'Arabic Title' of " + pageRequest.ArabicTitle + " already exists."}
			}
		} else if c.Param("id") != "" {
			if englisharabicname.EnglishTitle == pageRequest.EnglishTitle && c.Param("id") != englisharabicname.Id {
				errorFlag = true
				englishTitleError = EnglishTitleError{"error_page_english_title_not_unique", "Page with specified 'English Title' of " + pageRequest.EnglishTitle + " already exists."}
			}
			if englisharabicname.ArabicTitle == pageRequest.ArabicTitle && c.Param("id") != englisharabicname.Id {
				errorFlag = true
				arabicTitleError = ArabicTitleError{"error_page_arabic_title_not_unique", "Page with specified 'Arabic Title' of  " + pageRequest.ArabicTitle + " already exists."}
			}
		}
	}
	if pageRequest.EnglishTitle == "" {
		errorFlag = true
		englishTitleError = EnglishTitleError{"NotEmptyValidator", "'English Title' should not be empty."}
	}
	if pageRequest.ArabicTitle == "" {
		errorFlag = true
		arabicTitleError = ArabicTitleError{"NotEmptyValidator", "'Arabic Title' should not be empty."}
	}
	var englishPageFriendlyError EnglishPageFriendlyError
	if pageRequest.EnglishPageFriendlyUrl == "" {
		errorFlag = true
		englishPageFriendlyError = EnglishPageFriendlyError{"NotEmptyValidator", "'English Page Friendly Url' should not be empty."}
	}
	var arabicPageFriendlyError ArabicPageFriendlyError
	if pageRequest.ArabicPageFriendlyUrl == "" {
		errorFlag = true
		arabicPageFriendlyError = ArabicPageFriendlyError{"NotEmptyValidator", "'Arabic Page Friendly Url' should not be empty."}
	}
	var englishMetaTitleError EnglishMetaTitleError
	if pageRequest.EnglishMetaTitle == "" {
		errorFlag = true
		englishMetaTitleError = EnglishMetaTitleError{"NotEmptyValidator", "'English Meta Title' should not be empty."}
	}
	var arabicMetaTitleError ArabicMetaTitleError
	if pageRequest.ArabicMetaTitle == "" {
		errorFlag = true
		arabicMetaTitleError = ArabicMetaTitleError{"NotEmptyValidator", "'Arabic Meta Title' should not be empty."}
	}
	var englishMetaDescriptionError EnglishMetaDescriptionError
	if pageRequest.EnglishMetaDescription == "" {
		errorFlag = true
		englishMetaDescriptionError = EnglishMetaDescriptionError{"NotEmptyValidator", "'English Meta Description' should not be empty."}
	}
	var arabicMetaDescriptionError ArabicMetaDescriptionError
	if pageRequest.ArabicMetaDescription == "" {
		errorFlag = true
		arabicMetaDescriptionError = ArabicMetaDescriptionError{"NotEmptyValidator", "'Arabic Meta Description' should not be empty."}
	}
	fmt.Println("Create or update page %s : validation completed", pageRequest.EnglishTitle)
	var sliderIds SlidersIds
	if len(pageRequest.SlidersIds) > 1 {
		type Names struct {
			Name string
		}
		var names Names
		db.Debug().Table("slider s").Select("string_agg(s.name,',')::varchar as name ").Where("s.id in (?)", pageRequest.SlidersIds).Find(&names)
		split := strings.Split(names.Name, ",")
		errorFlag = true
		sliderIds = SlidersIds{"error_page_sliders_should_not_have_page_platforms_intersections", "Sliders should not have the page platforms intersections Sliders with intersections: " + strings.Join(split, ",")}
		fmt.Println("Create or update page %s : Sliders %s", pageRequest.EnglishTitle, names)
	}
	// var pageOrderNUmberError PageOrderNUmberError
	// if pageRequest.PageOrderNumber == 0 {
	// 	errorFlag = true
	// 	pageOrderNUmberError = PageOrderNUmberError{"NotEmptyValidator", "'Page Order Number' should not be empty."}
	// }
	var publishingPlatformsError PublishingPlatformsError
	// fmt.Println("Create or update page %s : Sliders %s",pageRequest.EnglishTitle,names)
	// fmt.Println((pageRequest.PublishingPlatforms), "..........")
	if len(pageRequest.PublishingPlatforms) == 0 {
		errorFlag = true
		publishingPlatformsError = PublishingPlatformsError{"NotEmptyValidator", "'Publishing Platforms' should not be empty."}
	}
	var regionssError RegionssError
	if len(pageRequest.Regions) == 0 {
		errorFlag = true
		regionssError = RegionssError{"NotEmptyValidator", "'Regions' should not be empty."}
	}
	var invalid Invalid
	if englishTitleError.Code != "" {
		invalid.EnglishTitleError = englishTitleError
	}
	if arabicTitleError.Code != "" {
		invalid.ArabicTitleError = arabicTitleError
	}
	if englishPageFriendlyError.Code != "" {
		invalid.EnglishPageFriendlyError = englishPageFriendlyError
	}
	if arabicPageFriendlyError.Code != "" {
		invalid.ArabicPageFriendlyError = arabicPageFriendlyError
	}
	if englishMetaTitleError.Code != "" {
		invalid.EnglishMetaTitleError = englishMetaTitleError
	}
	if arabicMetaTitleError.Code != "" {
		invalid.ArabicMetaTitleError = arabicMetaTitleError
	}
	if englishMetaDescriptionError.Code != "" {
		invalid.EnglishMetaDescriptionError = englishMetaDescriptionError
	}
	if arabicMetaDescriptionError.Code != "" {
		invalid.ArabicMetaDescriptionError = arabicMetaDescriptionError
	}
	// if pageOrderNUmberError.Code != "" {
	// 	invalid.PageOrderNUmberError = pageOrderNUmberError
	// }
	if publishingPlatformsError.Code != "" {
		invalid.PublishingPlatformsError = publishingPlatformsError
	}
	if regionssError.Code != "" {
		invalid.RegionssError = regionssError
	}
	if sliderIds.Code != "" {
		invalid.SlidersIds = sliderIds
	}
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}

	var page Page
	// var pageplaylist PagePlaylist
	if c.Param("id") != "" {
		tx.Debug().Table("page").Select("id").Where("id=?", c.Param("id")).Find(&page)
	}
	var pagetype int
	if pageRequest.IsHome {
		pagetype = 1
	} else {
		pagetype = 0
	}
	var menuposterimage bool
	if pageRequest.NonTextualData.MenuPosterImage != "" {
		menuposterimage = true
	} else {
		menuposterimage = false
	}
	var mobilemenuposterimage bool
	if pageRequest.NonTextualData.MobileMenuPosterImage != "" {
		mobilemenuposterimage = true
	} else {
		mobilemenuposterimage = false
	}
	var mobilemenu bool
	if pageRequest.NonTextualData.MobileMenu != "" {
		mobilemenu = true
	} else {
		mobilemenu = false
	}
	// create page
	if c.Param("id") == "" {
		var pagekey Page
		if err := tx.Debug().Table("page").Select("max(page_key)+1 as page_key,max(third_party_page_key) as third_party_page_key").Find(&pagekey).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		fmt.Println("Page key value........", pagekey)
		// add id for creating old pages with .net
		// Getting count of playlists which are not hidden
		// Getting count of slider which are not hidden
		var Activeslidercount int
		var Activeplaylistcount int
		for _, data := range pageRequest.SlidersIds {
			db.Debug().Raw("select count(*) from slider s where s.deleted_by_user_id is null and s.is_disabled = ? and id = ?", false, data).Count(&Activeslidercount)
		}
		for _, data := range pageRequest.PlaylistsIds {
			db.Debug().Raw("select count(*) from playlist p where p.deleted_by_user_id is null and p.is_disabled = ? and id = ?", false, data).Count(&Activeplaylistcount)
		}
		var IsCheck bool
		fmt.Println("ACTIVE PALYLIST COUNT", Activeplaylistcount)
		fmt.Println("ACTIVE SLIDER COUNT", Activeslidercount)
		if Activeplaylistcount > 0 || Activeslidercount > 0 {
			IsCheck = false
		} else {
			IsCheck = true
		}
		// for sync remove ,removed this Id: pageRequest.PageId, on below
		// for sync remove ,replaced this PageKey: pagekey.PageKey in place of pageRequest.PageKey,
		pageres := Page{EnglishTitle: pageRequest.EnglishTitle, ArabicTitle: pageRequest.ArabicTitle, PageOrderNumber: 1, EnglishPageFriendlyUrl: pageRequest.EnglishPageFriendlyUrl, ArabicPageFriendlyUrl: pageRequest.ArabicPageFriendlyUrl, EnglishMetaTitle: pageRequest.EnglishMetaTitle, ArabicMetaTitle: pageRequest.ArabicMetaTitle, EnglishMetaDescription: pageRequest.EnglishMetaDescription, ArabicMetaDescription: pageRequest.ArabicMetaDescription, IsDisabled: IsCheck, PageKey: pagekey.PageKey, PageType: pagetype, CreatedAt: time.Now(), ModifiedAt: time.Now(), HasMenuPosterImage: menuposterimage, HasMobileMenuPosterImage: mobilemenuposterimage, HasMobileMenu: mobilemenu, ThirdPartyPageKey: pagekey.ThirdPartyPageKey + 1}
		if pageresult := tx.Debug().Table("page").Create(&pageres).Error; pageresult != nil {
			tx.Rollback()
			fmt.Println("Create or update page %s : Error in creating page query", pageRequest.EnglishTitle, pageresult)
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		var pageplayarr []interface{}
		order := 0
		for _, data := range pageRequest.PlaylistsIds {
			order = order + 1
			pageplayres := PagePlaylist{PageId: pageres.Id, PlaylistId: data, Order: order}
			pageplayarr = append(pageplayarr, pageplayres)
		}
		fmt.Println("before bulk insertion")
		if pageplaylistres := gormbulk.BulkInsert(tx, pageplayarr, common.BULK_INSERT_LIMIT); pageplaylistres != nil {
			tx.Rollback()
			fmt.Println("Create or update page %s : Error in page playlists", pageRequest.EnglishTitle, pageplaylistres)
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		// var publisharr []interface{}
		// for _, value := range pageRequest.PublishingPlatforms {
		// 	publishres := PageTargetPlatform{PageId: pageres.Id, TargetPlatform: value, PageOrderNumber: pageRequest.PageOrderNumber}
		// 	publisharr = append(publisharr, publishres)
		// }
		// if publishresult := gormbulk.BulkInsert(tx, publisharr, common.BULK_INSERT_LIMIT); publishresult != nil {
		// 	tx.Rollback()
		// 	c.JSON(http.StatusInternalServerError, errorresponse)
		// 	return
		// }
		/*page-order(page target_platform) */
		for _, req := range pageRequest.PublishingPlatforms {
			fmt.Println("middle of the for loop11", req)
			// if updateError := db.Table("page_target_platform").Where("target_platform=? and page_order_number!=0 ", req).Create("page_order_number=page_order_number+1").Error; updateError != nil {
			// 	fmt.Println(updateError)
			// 	return
			// }
			// 	TODO
			if updateError := tx.Debug().Table("page_target_platform").Where("target_platform=? and page_order_number!=0 ", req).Update(map[string]interface{}{"page_order_number": gorm.Expr("page_order_number + ?", 1)}).Error; updateError != nil {
				fmt.Println(updateError)
				return
			}
			// fmt.Println("here printing pageorder number", pageRequest.PageOrderNumber)
			// if updateError := db.Debug().Table("page_target_platform").Where("target_platform=? and page_order_number!=0 and page_id=?", req, pageRequest.PageId).Update("page_order_number", pageRequest.PageOrderNumber).Error; updateError != nil {
			// 	fmt.Println(updateError)
			// 	return
			// }
			fmt.Println("middle of the for loop")
			//db.Debug().Raw("UPDATE page_target_platform	SET page_order_number=page_order_number+1 WHERE page_order_number != 0  AND target_platform=?", req)
			var newRecord PageTargetPlatform
			//newRecord.PageId = pageres.Id
			//newRecord.PageOrderNumber = pageRequest.PageOrderNumber
			// here we create page order number to 1 as it is creation
			// for sync remove newRecord.PageId = pageRequest.PageId
			newRecord.PageId = pageres.Id
			newRecord.PageOrderNumber = 1
			newRecord.TargetPlatform = req
			if insertError := db.Debug().Create(&newRecord).Error; insertError != nil {
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
		}
		var pagecountry []interface{}
		for _, result := range pageRequest.Regions {
			pagecountryres := PageCountry{PageId: pageres.Id, CountryId: result}
			pagecountry = append(pagecountry, pagecountryres)
		}
		fmt.Println("before bulk insertion page platform")
		gormbulk.BulkInsert(tx, pagecountry, common.BULK_INSERT_LIMIT)
		fmt.Println("before bulk insertion page platform after")
		for _, pagesliderrange := range pageRequest.SlidersIds {
			pagesliderres := PageSlider{PageId: pageres.Id, SliderId: pagesliderrange, Order: 1}
			if res := tx.Debug().Table("page_slider").Create(&pagesliderres).Error; res != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
		}
		if len(pageRequest.DefaultSliderId) != 0 {
			for _, pagesliderrange := range pageRequest.SlidersIds {
				pagesliderres := PageSlider{PageId: pageres.Id, SliderId: pagesliderrange, Order: 0}
				if response := tx.Debug().Table("page_slider").Create(&pagesliderres).Error; response != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, errorresponse)
					return
				}
			}
		}
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		/* This Fragment code written by tarun
		response := make(chan fragments.FragmentUpdate)
		go fragments.UpdatePageFragment(page.Id, c, response, pageRequest.PlaylistsIds, pageRequest.SlidersIds, pageRequest.DefaultSliderId)
		outPut := <-response
		if outPut.Err != nil {
			errorresponse.Description = outPut.Err.Error()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}*/
		fmt.Println("before file updload and sync operations")
		/*Upload Images In S3 with Respect to PageId */
		go PageFileUPload(pageRequest.NonTextualData, pageres.Id)
		/* update dirty count in page_sync table */
		go common.PageSynching(pageres.Id, pageres.PageKey, c)
		fmt.Println("before bulk insertion after page sync")
		finPageId := map[string]string{"id": pageres.Id}

		common.ClearRedisKeyFollowKeys(c, "Menus_Slider_*")

		c.JSON(http.StatusOK, gin.H{"data": finPageId})
		return

	} else {
		// here updation
		type PageKey struct {
			PageKey int `json:"page_key"`
		}
		var pagekey PageKey
		db.Debug().Table("page").Where("id=?", c.Param("id")).Find(&pagekey)
		// Getting count of playlists which are not hidden
		// Getting count of slider which are not hidden
		var Activeslidercount int
		var Activeplaylistcount int
		for _, data := range pageRequest.SlidersIds {
			db.Debug().Raw("select count(*) from slider s where s.deleted_by_user_id is null and s.is_disabled = ? and id = ?", false, data).Count(&Activeslidercount)
		}
		for _, data := range pageRequest.PlaylistsIds {
			db.Debug().Raw("select count(*) from playlist p where p.deleted_by_user_id is null and p.is_disabled = ? and id = ?", false, data).Count(&Activeplaylistcount)
		}
		var IsCheck bool
		fmt.Println("ACTIVE PALYLIST COUNT IN UPDATION", Activeplaylistcount)
		fmt.Println("ACTIVE SLIDER COUNT IN UPDATION", Activeslidercount)
		if Activeplaylistcount > 0 || Activeslidercount > 0 {
			IsCheck = false
		} else {
			IsCheck = true
		}
		pageres := Page{EnglishTitle: pageRequest.EnglishTitle, ArabicTitle: pageRequest.ArabicTitle, PageOrderNumber: pageRequest.PageOrderNumber, PageType: pagetype, EnglishPageFriendlyUrl: pageRequest.EnglishPageFriendlyUrl, ArabicPageFriendlyUrl: pageRequest.ArabicPageFriendlyUrl, EnglishMetaTitle: pageRequest.EnglishMetaTitle, ArabicMetaTitle: pageRequest.ArabicMetaTitle, EnglishMetaDescription: pageRequest.EnglishMetaDescription, ArabicMetaDescription: pageRequest.ArabicMetaDescription, IsDisabled: IsCheck, ModifiedAt: time.Now(), HasMenuPosterImage: menuposterimage, HasMobileMenuPosterImage: mobilemenuposterimage, HasMobileMenu: mobilemenu}
		if res := db.Debug().Table("page").Where("id=? and page_key = ?", c.Param("id"), pagekey.PageKey).Update(pageres).Error; res != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		var pageplay PagePlaylist
		tx.Debug().Table("page_playlist").Where("page_id=?", c.Param("id")).Delete(&pageplay)
		if len(pageRequest.PlaylistsIds) != 0 {
			for i, data := range pageRequest.PlaylistsIds {
				pageplaylistres := PagePlaylist{PageId: c.Param("id"), PlaylistId: data, Order: i + 1}
				if res := tx.Table("page_playlist").Create(&pageplaylistres).Error; res != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, errorresponse)
					return
				}
			}
		}
		/* PageOrder for new platform */
		/*var newarr []int
		var fetchPagePlatforms []PageTargetPlatform
		for _, platformId := range pageRequest.PublishingPlatforms {
			newarr = append(newarr, platformId)
		}
		var pagetargetplatform PageTargetPlatform
		if pagefinal := db.Debug().Table("page_target_platform").Where("page_id=? and target_platform not in (?)", c.Param("id"), newarr).Delete(&pagetargetplatform).Error; pagefinal != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		var newarrays []int
		db.Debug().Table("page_target_platform").Select("*").Where("page_id=?", c.Param("id")).Find(&fetchPagePlatforms)
		var exist bool
		for _, req := range pageRequest.PublishingPlatforms {
			exist = false
			for _, order := range fetchPagePlatforms {
				if req == order.TargetPlatform {
					exist = true
					break
				}
			}
			if !(exist) {
				//TODO
				// if updateError := db.Table("page_target_platform").Where("target_platform=? page_order_number!=0", req).Update("page_order_number=page_order_number+1").Error; updateError != nil {
				// 	fmt.Println(updateError)
				// 	return
				// }
				var newRecord PageTargetPlatform
				newRecord.PageId = c.Param("id")
				newRecord.PageOrderNumber = pageRequest.PageOrderNumber
				newRecord.TargetPlatform = req
				if insertError := db.Debug().Update(&newRecord).Error; insertError != nil {
					c.JSON(http.StatusInternalServerError, errorresponse)
					return
				}
				// db.Table("page_target_platform").Where("page_order_number!=0 and target_platform=?", req).Update("page_order_number=page_order_number+1")
				// db.Table("page_target_platform").Where("page_id=? and target_platform=?",c.Param("id"),req).Update("")
				// if new := db.Raw("Update page_target_platform set page_order_number=? where page_id=? target_platform in (?)", 1, c.Param("id"), newarrays).Error; new != nil {
				// 	c.JSON(http.StatusInternalServerError, new)
				// 	return
				// }
				newarrays = append(newarrays, req)
			}
		}*/
		//fmt.Println(newarrays, ".....................")
		// if updateError := tx.Table("page_target_platform").Where("page_order_number!=0 and target_platform in (?)", newarrays).Update("page_order_number=page_order_number+1").Error; updateError != nil {
		// 	fmt.Println(updateError, ";;;;;;;;")
		// 	c.JSON(http.StatusInternalServerError, updateError)
		// 	return
		// }
		// if new := tx.Raw("Update page_target_platform set page_order_number=? where page_id=? target_platform in (?)", 1, c.Param("id"), newarrays).Error; new != nil {
		// 	c.JSON(http.StatusInternalServerError, new)
		// 	return
		// }

		// var pagetargetplatform PageTargetPlatform
		// if pagefinal := tx.Table("page_target_platform").Where("page_id=?", c.Param("id")).Delete(&pagetargetplatform).Error; pagefinal != nil {
		// 	tx.Rollback()
		// 	c.JSON(http.StatusInternalServerError, errorresponse)
		// 	return
		// }
		// var pagetargerarr []interface{}
		// for _, value := range pageRequest.PublishingPlatforms {
		// 	publishres := PageTargetPlatform{PageId: c.Param("id"), TargetPlatform: value, PageOrderNumber: pageRequest.PageOrderNumber}
		// 	pagetargerarr = append(pagetargerarr, publishres)
		// }
		// if err := gormbulk.BulkInsert(tx, pagetargerarr, common.BULK_INSERT_LIMIT); err != nil {
		// 	tx.Rollback()
		// 	c.JSON(http.StatusInternalServerError, errorresponse)
		// 	return
		// }
		// fetching previous platforms here
		var previousplatforms []PageTargetPlatform
		if err := db.Debug().Table("page_target_platform").Where("page_id = ? ", c.Param("id")).Find(&previousplatforms).Error; err != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		var temp []int
		for _, val := range previousplatforms {
			temp = append(temp, val.TargetPlatform)
		}
		fmt.Println("previous platforrms", temp)
		fmt.Println("current platforms", pageRequest.PublishingPlatforms)
		var notneededplatforms []int
		temp = append(temp, pageRequest.PublishingPlatforms...)
		fmt.Println("temp after append", temp)
		// notneededplatforms = unique(temp)
		// final:=make(map[int]int)
		final := DupCount(temp)
		for item, count := range final {
			if count == 1 {
				notneededplatforms = append(notneededplatforms, item)
			}
		}
		fmt.Println("not needed platforms", notneededplatforms)
		// here deleting the not needed platforms
		for _, delvalue := range notneededplatforms {
			var pageordernum PageTargetPlatform
			var existingcount int
			db.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", c.Param("id"), delvalue).Find(&pageordernum).Count(&existingcount)
			fmt.Println("existing count", existingcount)
			if existingcount > 0 {
				var ptp PageTargetPlatform
				if err := db.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", c.Param("id"), delvalue).Delete(&ptp).Error; err != nil {
					c.JSON(http.StatusInternalServerError, errorresponse)
					return
				}
				// here decrementing the page order by 1 for all the pages whose page order number is greater than current page order number
				if updateError := db.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? and page_order_number > ?  ", c.Param("id"), delvalue, pageordernum.PageOrderNumber).Update(map[string]interface{}{"page_order_number": gorm.Expr("page_order_number - ?", 1)}).Error; updateError != nil {
					fmt.Println(updateError)
					return
				}
			}
		}
		for _, value := range pageRequest.PublishingPlatforms {
			var forcount int
			db.Debug().Table("page_target_platform").Where("page_id = ? and target_platform = ? ", c.Param("id"), value).Count(&forcount)
			if forcount < 1 {
				var newRecord PageTargetPlatform
				newRecord.PageId = c.Param("id")
				newRecord.PageOrderNumber = 1
				newRecord.TargetPlatform = value
				if insertError := db.Debug().Create(&newRecord).Error; insertError != nil {
					c.JSON(http.StatusInternalServerError, errorresponse)
					return
				}

			}
		}
		var pagecountry PageCountry
		if err := tx.Debug().Table("page_country").Where("page_id=?", c.Param("id")).Delete(&pagecountry).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			fmt.Println(errorresponse)
			return
		}
		var pagecountryfinal []interface{}
		for _, result := range pageRequest.Regions {
			pagecountryres := PageCountry{PageId: c.Param("id"), CountryId: result}
			pagecountryfinal = append(pagecountryfinal, pagecountryres)
		}
		if err := gormbulk.BulkInsert(tx, pagecountryfinal, common.BULK_INSERT_LIMIT); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
		}
		// tx.Debug().Table("page_slider ps").Where("ps.page_id=? and ps.order=?", c.Param("id"), 1).Delete(&PageSlider{})
		// here where condition is removed because it is not satifiying based on implementations we need to change in future if required
		tx.Debug().Table("page_slider ps").Where("ps.page_id=?", c.Param("id")).Delete(&PageSlider{})
		if len(pageRequest.SlidersIds) != 0 {
			for _, pagesliderrange := range pageRequest.SlidersIds {
				pagesliderres := PageSlider{PageId: c.Param("id"), SliderId: pagesliderrange, Order: 1}
				if res := tx.Debug().Table("page_slider").Create(&pagesliderres).Error; res != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, errorresponse)
					return
				}
			}
		}
		/*if len(pageRequest.DefaultSliderId) != 0 {
			for _, pagesliderrange := range pageRequest.SlidersIds {
				pagesliderres := PageSlider{SliderId: pagesliderrange, Order: 0}
				if res := tx.Debug().Table("page_slider").Where("page_id=?", c.Param("id")).Update(pagesliderres).Error; res != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, errorresponse)
					return
				}
			}
		}*/
		if pageRequest.DefaultSliderId != "" {
			pagedefaultsliderres := PageSlider{PageId: c.Param("id"), SliderId: pageRequest.DefaultSliderId, Order: 0}
			if res := tx.Debug().Table("page_slider").Create(&pagedefaultsliderres).Error; res != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
		}
		/*Upload Images In S3 with Respect to PageId */
		go PageFileUPload(pageRequest.NonTextualData, c.Param("id"))
		/* update dirty count in page_sync table */
		go common.PageSynching(c.Param("id"), pagekey.PageKey, c)

		finPageId := map[string]string{"id": c.Param("id")}

		common.ClearRedisKeyFollowKeys(c, "Menus_Slider_*")

		c.JSON(http.StatusOK, gin.H{"data": finPageId})
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	/* This Fragment code written by tarun
	response := make(chan fragments.FragmentUpdate)
	go fragments.UpdatePageFragment(page.Id, c, response, pageRequest.PlaylistsIds, pageRequest.SlidersIds, pageRequest.DefaultSliderId)
	outPut := <-response
	if outPut.Err != nil {
		errorresponse.Description = outPut.Err.Error()
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}*/
}
func DupCount(list []int) map[int]int {

	duplicate_frequency := make(map[int]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		_, exist := duplicate_frequency[item]

		if exist {
			duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
		} else {
			duplicate_frequency[item] = 1 // else start counting from 1
		}
	}
	return duplicate_frequency
}

// GetallPageRegionsByPageId -  Get pageRegions Based on pageId
// GET /api/pages/:id/region
// @Summary Get pageRegions Based on pageId
// @Description  Get pageRegions Based on pageId
// @Tags Page
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param id path string true "id"
// @Success 200 {array} object c.JSON
// @Router /api/pages/{id}/region [get]
func (hs *HandlerService) GetallPageRegionsByPageId(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	id := c.Param("id")
	var errorresponse = common.ServerErrorResponse()
	var pageRegions []Country
	/*Fetch page regions*/
	if resultError := db.Debug().Table("page_country pc").Select("c.english_name as name,c.id").
		Joins("left join country c on pc.country_id =c.id").
		Where("pc.page_id=?", id).Find(&pageRegions).Error; resultError != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
	}
	c.JSON(http.StatusOK, gin.H{"data": pageRegions})
}

/*Uploade image Based on Page Id*/
func PageFileUPload(request NonTextualData, pageId string) {
	bucketName := os.Getenv("S3_BUCKET")
	var newarr []string
	newarr = append(newarr, request.MenuPosterImage)
	newarr = append(newarr, request.MobileMenuPosterImage)
	newarr = append(newarr, request.MobileMenu)
	for k := 0; k < len(newarr); k++ {
		item := newarr[k]
		filetrim := strings.Split(item, "_")
		Destination := pageId + "/" + filetrim[0]
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
		errorr := SizeUploadFileToS3(s, filetrim[0], pageId)
		if errorr != nil {
			fmt.Println("error in uploading size upload", errorr)
		}
		fmt.Println("Success!")
	}
}

// SizeUploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func SizeUploadFileToS3(s *session.Session, fileName string, contentId string) error {
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
		s3file := sizeValue[i] + contentId + "/" + fileName
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
			fmt.Println("Unable to upload", er)
		}
		fmt.Printf("Successfully uploaded %q", fileName)
	}
	return er
}
