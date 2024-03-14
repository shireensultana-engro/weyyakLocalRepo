package geoblock

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
	l "user/logger"
	"user/register"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Bootstrap ... router export
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	r.POST("/geoblock", hs.PostNewsletters)
}

var userTable string = "user"

// PostNewsletters ... post geoblock email
// RegisterUserUsingEmail -  Creates a new user using email id
// POST /geoblock
// @Summary post geoblock email
// @Description post geoblock email
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body NewsLetter true "Raw JSON string"
// @Success 200
// @Router /geoblock [post]
func (hs *HandlerService) PostNewsletters(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	var (
		input NewsLetter
	)

	if err := c.ShouldBindJSON(&input); err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}

	if input.Email == "" {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": "Email is required", "Status": http.StatusBadRequest})
		return
	}

	if !validateEmail(input.Email) {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": "Invalid email", "Status": http.StatusBadRequest})
		return
	}

	var count int64
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(userTable).Where("email = ?", input.Email).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("email '%s' is already email already exists", input.Email)
		}

		if err := tx.Table(userTable).Create(&register.User{
			Email:              input.Email,
			UserName:           input.Email,
			RoleId:             "91f15b92-97fd-e611-814f-0af7afba4acb",
			RegistrationSource: 6,
			RegisteredAt:       time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("failed to create newsletter entry: %v", err)
		}

		return nil
	})

	if err != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	l.JSON(c, http.StatusOK, gin.H{"status": 1, "message": "Email updated successfully"})
}

func validateEmail(email string) bool {
	// regular expression for email validation
	// this regex is just a simple example and may not cover all edge cases
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
