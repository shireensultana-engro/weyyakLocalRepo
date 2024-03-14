package admin

import (

	// u "pdfGenerator"

	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"user/common"
	l "user/logger"
	"user/register"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/thanhpk/randstr"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	qrg := r.Group("/api")
	qrg.Use(common.ValidateToken())
	qrg.PUT("/admins", hs.CreateAdminDetails)
	qrg.GET("/admins", hs.GetAllAdminDetails)
	qrg.POST("/admin/:id", hs.UpdateAdminDetails)
	qrg.DELETE("/admins/:id", hs.DeleteAdminDetails)
	r.PUT("api/admin/reset_password_emails", hs.ResetAdminPassword)
	r.PUT("api/admin/password", hs.SetAdminPassword)
	r.POST("api/admin/password", hs.SetAdminPassword)

	/*Error code Exception URL*/
	//create admin
	/*qrg.GET("/admins", hs.CreateAdminDetails)
	qrg.POST("/admins", hs.CreateAdminDetails)
	qrg.DELETE("/admins", hs.CreateAdminDetails)
	//get all admins
	qrg.POST("/admins", hs.GetAllAdminDetails)
	qrg.DELETE("/admins", hs.GetAllAdminDetails)
	qrg.PUT("/admins", hs.GetAllAdminDetails)
	//Update admin details//
	qrg.GET("/admin/:id", hs.UpdateAdminDetails)
	qrg.PUT("/admin/:id", hs.UpdateAdminDetails)
	qrg.DELETE("/admin/:id", hs.UpdateAdminDetails)
	//Delete admin details//
	qrg.GET("/admins/:id", hs.DeleteAdminDetails)
	qrg.POST("/admins/:id", hs.DeleteAdminDetails)
	qrg.PUT("/admins/:id", hs.DeleteAdminDetails)*/
}

// CreateAdminDetails -  create admin details
// PUT /api/admins
// @Summary Create admin details
// @Description Create admin details
// @Tags Admin
// @Accept  json
// @Security Authorization
// @Produce  json
// @Success 200
// @Router /api/admins [put]
func (hs *HandlerService) CreateAdminDetails(c *gin.Context) {
	// /*405(Request-method)*/
	// if c.Request.Method != http.MethodPut {
	// 	l.JSON(c, http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
	// 	return
	// }
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var request AdminDetailsRequest
	c.ShouldBindJSON(&request)
	/*Input-Validations*/
	var emailcheck Emailcheck
	var emailError Email
	err := db.Table("user").Where("email = ? and is_deleted='false'", request.Email)
	err.Scan(&emailcheck)
	if len(emailcheck.Email) > 0 {
		emailError = Email{"error_user_email_already_exists", "Specified email already exists."}
	}
	if !common.RegEmail(request.Email) && request.Email != "" {
		emailError = Email{"error_user_email_invalid", "Email is invalid."}
	}
	if request.Email == "" {
		emailError = Email{"error_user_email_required", "Email is required."}
	}
	if emailError.Code != "" {
		invalid := Invalid{Email: &emailError}
		finalErrorResponse := FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
		l.JSON(c, http.StatusBadRequest, finalErrorResponse)
		return
	}
	/*End for Validations*/
	user := register.User{Email: request.Email, FirstName: request.FirstName, LastName: request.LastName, IsBackOfficeUser: true, LanguageId: 1, IsAdult: true, LastActivityAt: time.Now(), RegisteredAt: time.Now(), ModifiedAt: time.Now(), EmailConfirmed: true, Version: 2, UserName: request.Email, RoleId: "90f15b92-97fd-e611-814f-0af7afba4acb"}
	if admincreate := db.Create(&user).Error; admincreate != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	// var token string
	//ToDo: Token to be generated and url to be changed
	token := common.EncodeToString(6)
	templateData := struct {
		PasswordChangeUrl string
	}{
		PasswordChangeUrl: os.Getenv("BOPASSWORDCHANGEURL") + "email=" + request.Email + "&resetPasswordToken=" + token,
	}
	templatePath := "BackOfficeUserSetupPasswordBody.html"
	con, _ := ioutil.ReadFile(templatePath)
	content := string(con)
	fmt.Println(string(content))
	// if err := r.ParseTemplate(templatePath, templateData); err == nil {
	fmt.Println(templateData, templatePath, "rrrrrrrrrrrrrrrrrrr")
	content = strings.Replace(content, "{{.PasswordChangeUrl}}", templateData.PasswordChangeUrl, 1)
	message := template.HTML(content)
	error := common.SendMail(request.Email, string(message), "Subject: Welcome to Weyyak!")
	if error != nil {
		fmt.Println("Email has not sent- ", error)
	}
	return
}

// GetAllAdminDetails -  Get user details
// GET /api/admins
// @Summary Get all user details
// @Description Get all user details
// @Tags User
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Param searchText query string false "Search Text"
// @Success 200
// @Router /api/admins [get]
func (hs *HandlerService) GetAllAdminDetails(c *gin.Context) {
	// /*405(Request-method)*/
	// if c.Request.Method != http.MethodGet {
	// 	l.JSON(c, http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
	// 	return
	// }
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
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
	var searchText string
	if c.Request.URL.Query()["searchText"] != nil {
		searchText = c.Request.URL.Query()["searchText"][0]
	}
	where := "is_back_office_user ='true' and email_confirmed ='true' and is_deleted ='false' "
	if searchText != "" {
		where += " and (first_name ilike '%" + searchText + "%' or last_name ilike '%" + searchText + "%' or email like '%" + searchText + "%')"
	}
	if limit == 0 {
		limit = 50
	}
	// response := []UserDetails{}
	var userKeys []UserKeys
	var pagination PaginationResult
	var totalCount int
	userId := c.MustGet("userid")
	if err := db.Table("public.user").Select("first_name ,last_name ,user_name as email,'admin' as userRole,registered_at ,id").
		Where(where).Count(&totalCount).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	if err := db.Table("public.user").Select("first_name ,last_name ,user_name as email,'Admin' as user_role,registered_at ,id,case when id=? then false else true end as allow_delete", userId).
		Where(where).Order("registered_at desc").
		Limit(limit).Offset(offset).Find(&userKeys).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	pagination.Size = int(totalCount)
	pagination.Offset = int(offset)
	pagination.Limit = limit
	l.JSON(c, http.StatusOK, gin.H{"pagination": pagination, "data": userKeys})
	return
}

// DeleteAdminDetails - Deleting admin details
// DELETE /api/admins/{id}
// @Summary Delete user details
// @Description Delete user details
// @Tags User
// @Security Authorization
// @Accept  json
// @Produce  json
// @Param id path string true "Id"
// @Success 200
// @Router /api/admins/{id} [delete]
func (hs *HandlerService) DeleteAdminDetails(c *gin.Context) {
	// /*405(Request-method) Handling*/
	// if c.Request.Method != http.MethodDelete {
	// 	l.JSON(c, http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
	// 	return
	// }
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	/*Validating for right admin to delete admin Detail*/
	UserId := c.MustGet("userid")
	if UserId == c.Param("id") {
		l.JSON(c, http.StatusBadRequest, gin.H{"Message": " Not allowed to delete himself."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var totalCount int
	if adminresult := db.Table("public.user").Select("id").Where("id=? and is_back_office_user=?", c.Param("id"), true).Count(&totalCount).Error; adminresult != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	if totalCount > 0 {
		if deleteadmin := db.Table("public.user").Where("id = ?", c.Param("id")).Update("is_deleted", true).Error; deleteadmin != nil {
			l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
			return
		}
		l.JSON(c, http.StatusOK, gin.H{})
	} else {
		l.JSON(c, http.StatusNotFound, gin.H{"error": "not_found", "description": "Not found.", "code": "", "requestId": randstr.String(32)})
		return
	}

}

// UpdateAdminDetails -  update admin details
// POST /api/admins
// @Summary Update admin details
// @Description Update admin details
// @Tags Admin
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param id path string true "Id"
// @Success 200
// @Router /api/admins [post]
func (hs *HandlerService) UpdateAdminDetails(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	// /*405(Request-method)*/
	// if c.Request.Method != http.MethodPost {
	// 	l.JSON(c, http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
	// 	return
	// }
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var request AdminDetailsRequest
	c.ShouldBindJSON(&request)
	/*Input-Validations*/
	var emailError Email
	if !common.RegEmail(request.Email) && request.Email != "" {
		emailError = Email{"error_user_email_invalid", "Email is invalid."}
	}
	if request.Email == "" {
		emailError = Email{"error_user_email_required", "Email is required."}
	}
	if emailError.Code != "" {
		invalid := Invalid{Email: &emailError}
		finalErrorResponse := FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
		l.JSON(c, http.StatusBadRequest, finalErrorResponse)
		return
	}
	passwordHash, saltStored := common.HashPassword(request.Password)
	db.Table("public.user").Where("id = ?", c.Param("id")).Updates(map[string]interface{}{"email": request.Email, "first_name": request.FirstName, "last_name": request.LastName, "password_hash": passwordHash, "salt_stored": saltStored, "version": 2, "modified_at": time.Now(), "last_activity_at": time.Now()})
	return
}

// ResetAdminPassword -  Admin password reset
// PUT /api/admin/reset_password_emails
// @Summary Admin password reset
// @Description Admin password reset
// @Tags Admin
// @Accept  json
// @Produce  json
// @Success 200
// @Router /api/admin/reset_password_emails [put]
func (hs *HandlerService) ResetAdminPassword(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var request ResetPasswordAdminRequest
	c.ShouldBindJSON(&request)
	/*Input-Validations*/
	var emailError Email
	var count int
	db.Table("user").Where("email = ?", request.Email).Count(&count)
	if count < 1 {
		emailError = Email{"error_user_email_does_not_exists", "Specified email does not exists."}
		l.JSON(c, http.StatusBadRequest, emailError)
		return
	}
	/*End for Validations*/
	token := base64.StdEncoding.EncodeToString([]byte(request.Email))
	details := EmailOtpRecord{Phone: request.Email, Message: token, SentOn: time.Now()}
	if final := db.Table("otp_record").Create(&details).Error; final != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"message": final.Error(), "status": http.StatusInternalServerError})
		return
	}
	templateData := struct {
		PasswordChangeUrl string
	}{
		PasswordChangeUrl: "https://wyk2boqa.weyyak.com/reset-password?" + "email=" + request.Email + "&resetPasswordToken=" + token,
	}
	templatePath := "BackOfficeUserSetupPasswordBody.html"
	con, _ := ioutil.ReadFile(templatePath)
	content := string(con)
	content = strings.Replace(content, "{{.PasswordChangeUrl}}", templateData.PasswordChangeUrl, 1)
	message := template.HTML(content)
	error := common.SendMail(request.Email, string(message), "Subject: Welcome to Weyyak!")
	if error != nil {
		fmt.Println("Email has not sent- ", error)
	}
	l.JSON(c, http.StatusAccepted, gin.H{})
	return
}

// SetAdminPassword -  Admin password set
// PUT /api/admin/password
// @Summary Admin password set
// @Description Admin password set
// @Tags Admin
// @Accept  json
// @Produce  json
// @Success 200
// @Router /api/admin/password [put]
func (hs *HandlerService) SetAdminPassword(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var request RequestSetPasswordParameters
	c.ShouldBindJSON(&request)
	/*Input-Validations*/
	hashedPassword, saltStored := common.HashPassword(request.Password)
	db.Table("user").Where("email=?", request.Email).Update("password_hash", hashedPassword)
	db.Table("user").Where("email=?", request.Email).Update("salt_stored", saltStored)
	l.JSON(c, http.StatusOK, gin.H{"status": 1, "message": "password changed successfully"})
}
