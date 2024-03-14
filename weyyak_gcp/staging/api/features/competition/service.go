package competition

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"feature/common"

	"github.com/thanhpk/randstr"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	srg := r.Group("/api/v1")
	srg.POST("/users", hs.PostCompetitionUserDetails)
	srg.GET("/country", hs.GetCountries)
	srg.GET("/agegroup", hs.GetAgeGroup)
	r.GET("/v1/agegroup", hs.GetAgeGroup)

}

// CreateCompetitionUserDetails -  create competition user details
// POST /api/v1/users
// @Summary Create competition user details
// @Description Create competition user details
// @Tags Competition Users
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param body body CompetitionUsers true "Raw JSON string"
// @Success 200
// @Router /api/v1/users [post]
func (hs *HandlerService) PostCompetitionUserDetails(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var request CompetitionUsers
	var templatePath, subjt string
	var emailcheck Emailcheck
	var errorFlag bool
	errorFlag = false
	var emailError ErrorCode
	var invalid Invalid
	var finalErrorResponse FinalErrorResponse
	var templateData struct {
		Name string
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.Table("competition_users").Where("email = ? or mobile= ?", request.Email, request.Mobile)
	err.Scan(&emailcheck)
	if len(emailcheck.Email) > 0 {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_or_mobile_already_exists", "Specified email or mobile already exists."}
	}
	if !common.RegEmail(request.Email) && request.Email != "" {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_invalid", "Email is invalid."}
	}
	if request.Email == "" {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_required", "Email is required."}
	}

	if emailError.Code != "" {
		invalid = Invalid{Email: &emailError}
	}
	finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
	if err := db.Debug().Create(&request).Error; err != nil {
		fmt.Println("errr", err)
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest})
		return
	}
	if strings.ToLower(request.Language) == "en" {
		templatePath = "emailer.html"
		templateData.Name = request.EnglishFullName
		subjt = "You participation in Weyyak Safra wa Soufra"
	} else {
		templatePath = "emailerar.html"
		templateData.Name = request.EnglishFullName
		subjt = "اشتراكك برحلة وياك سَفرة وسُفرة"
	}
	fmt.Println(templatePath, "llllllllll")
	con, _ := ioutil.ReadFile(templatePath)
	content := string(con)
	fmt.Println(string(content), "jjjjjjjj")
	fmt.Println(templateData, templatePath, "rrrrrrrrrrrrrrrrrrr")
	content = strings.Replace(content, "{{.UserName}}", templateData.Name, 1)
	message := template.HTML(content)
	erro := common.AwsSendMail(request.Email, string(message), subjt)
	if erro != "" {
		fmt.Println("Email has not sent - ", erro)
		request.MailSent = "Error in sending mail"
		db.Table("competition_users").Where("id = ?", request.Id).Update("mail_sent", request.MailSent)
	} else {
		request.MailSent = "Successfully email sent to user"
		db.Table("competition_users").Where("id = ?", request.Id).Update("mail_sent", request.MailSent)
		db.Table("competition_users").Where("id = ?", request.Id).Update("email_confirmed", true)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registered Successfully"})
}

// GetCompetitionUserCountry -  get competition user country details
// GET /api/v1/country
// @Summary Get competition user country details
// @Description Get competition user country details
// @Tags Country
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param body body CompetitionUsers true "Raw JSON string"
// @Success 200
// @Router /api/v1/country [get]
func (hs *HandlerService) GetCountries(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var request []Country

	if err := db.Find(&request).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": request})

}

// GetCompetitionUserC -  get competition user age group details
// POST /api/v1/agegroup
// @Summary Get competition user age group details
// @Description Get competition user age group details
// @Tags Age Group
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param body body CompetitionUsers true "Raw JSON string"
// @Success 200
// @Router /api/v1/agegroup [get]
func (hs *HandlerService) GetAgeGroup(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var request []AgeGroup
	if err := db.Find(&request).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": request})

}
