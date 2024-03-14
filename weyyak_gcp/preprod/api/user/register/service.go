package register

import (
	"encoding/base64"
	//	"encoding/json"
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"

	// u "pdfGenerator"
	"strconv"
	"time"
	"user/common"

	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	// "github.com/aws/aws-sdk-go/service/sns"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/gomail.v2"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/nyaruka/phonenumbers"
	"github.com/thanhpk/randstr"
	"github.com/xuri/excelize/v2"

	"github.com/dkoston/xlsx2csv/xlsx"
	// "encoding/csv"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	sqrg := r.Group("/users")
	sqrg1 := r.Group("/v1/users")
	sqrg2 := r.Group("/v1/user")
	sqrg.POST("/register_email", hs.RegisterUserUsingEmail)
	sqrg1.POST("/register_email", hs.RegisterUserUsingEmail)
	sqrg.POST("/register_sms", hs.RegisterUserUsingSMS)
	sqrg1.POST("/register_sms", hs.RegisterUserUsingSMS)
	rqrg := r.Group("/user")
	rqrg.POST("/registration_confirmation", hs.RegistrationConfirmation)
	sqrg2.POST("/registration_confirmation", hs.RegistrationConfirmation)
	rqrg.PUT("/reset_password_emails", hs.ResetPasswordEmails)
	sqrg2.PUT("/reset_password_emails", hs.ResetPasswordEmails)
	rqrg.POST("/password", hs.ResetPasswordWithEmail)
	sqrg2.POST("/password", hs.ResetPasswordWithEmail)
	qrg := r.Group("/v1/users")
	qrg.Use(common.ValidateToken())
	qrg.POST("/self", hs.UpdateUserProfile)
	qrg.GET("/self", hs.GetUserDetails)
	dqrg := r.Group("/v1")
	dqrg.Use(common.ValidateToken())
	dqrg.GET("/devices/", hs.GetUserDevices)
	dqrg.GET("/devices", hs.GetUserDevices)
	doqrg := r.Group("/api")
	doqrg.Use(common.ValidateToken())
	doqrg.GET("/devices/limit", hs.GetUserDevicesLimit)
	doqrg.POST("/devices/limit", hs.UpdateUserDevicesLimitCount)
	dqrg.GET("/devices/limit", hs.GetUserDevicesLimit)
	dqrg.POST("/devices/limit", hs.UpdateUserDevicesLimitCount)

	sqrg.POST("/self/paycms_status", hs.UpdatePaycmsStatus)
	sqrg1.POST("/self/paycms_status", hs.UpdatePaycmsStatus)
	sqrg.POST("/resend_email", hs.ResendEmailVerification)
	sqrg1.POST("/resend_email", hs.ResendEmailVerification)
	sqrg.POST("/send_otp", hs.SendOtp)
	sqrg.POST("/verify_otp", hs.VerifyOtp)
	sqrg.POST("/password_otp", hs.ForgotPasswordOtp)
	mqrg := r.Group("/users")
	mqrg.Use(common.ValidateToken())
	mqrg.POST("/self/phone_number", hs.PhonenumberChangeOtp)
	mqrg.POST("/self/password", hs.ChangePassword)
	dqrg.POST("/logout", hs.Logout)
	oqrg := r.Group("/oauth2")
	oqrg.POST("/device/code", hs.GeneratePairingCode)
	uoqrg := r.Group("/oauth2")
	uoqrg.Use(common.ValidateToken())
	uoqrg.POST("/device/auth", hs.VerifyPairingCode)

	//twitter related APIs
	r.GET("/:lang/usertoken", hs.TwitterUserToken)
	r.GET("/:lang/getAccessToken", hs.GetTwitterAccessToken)

	bqrg := r.Group("/api/users")
	// boqrg := r.Group("/api/user")
	// boqrg.Use(common.ValidateToken())
	bqrg.Use(common.ValidateToken())
	bqrg.GET("/export", hs.ExportUserDetails)
	bqrg.POST("/:id", hs.UpdateUserDetailsByUserid)
	bqrg.GET("/filters", hs.UserFilterslist)
	bqrg.GET("/:id/viewactivities", hs.UserViewActivitybyFilters)
	bqrg.GET("", hs.UsersListandSearchbyFilterswithPagination)
	// r.POST("/oauth2/token", hs.Login)

	/* Exception Urls */
	//Register with Email
	sqrg.PUT("/register_email", hs.RegisterUserUsingEmail)
	sqrg.DELETE("/register_email", hs.RegisterUserUsingEmail)
	sqrg.GET("/register_email", hs.RegisterUserUsingEmail)
	sqrg2.PUT("/register_email", hs.RegisterUserUsingEmail)
	sqrg2.DELETE("/register_email", hs.RegisterUserUsingEmail)
	sqrg2.GET("/register_email", hs.RegisterUserUsingEmail)
	//Register with SMS
	sqrg.PUT("/register_sms", hs.RegisterUserUsingSMS)
	sqrg.DELETE("/register_sms", hs.RegisterUserUsingSMS)
	sqrg.GET("/register_sms", hs.RegisterUserUsingSMS)
	sqrg1.PUT("/register_sms", hs.RegisterUserUsingSMS)
	sqrg1.DELETE("/register_sms", hs.RegisterUserUsingSMS)
	sqrg1.GET("/register_sms", hs.RegisterUserUsingSMS)
	//Registration Confirmation
	rqrg.PUT("/registration_confirmation", hs.RegistrationConfirmation)
	rqrg.GET("/registration_confirmation", hs.RegistrationConfirmation)
	rqrg.DELETE("/registration_confirmation", hs.RegistrationConfirmation)
	sqrg2.PUT("/registration_confirmation", hs.RegistrationConfirmation)
	sqrg2.GET("/registration_confirmation", hs.RegistrationConfirmation)
	sqrg2.DELETE("/registration_confirmation", hs.RegistrationConfirmation)
	//Resend Email for Registration
	sqrg.PUT("/resend_email", hs.ResendEmailVerification)
	sqrg.GET("/resend_email", hs.ResendEmailVerification)
	sqrg.DELETE("/resend_email", hs.ResendEmailVerification)
	sqrg1.PUT("/resend_email", hs.ResendEmailVerification)
	sqrg1.GET("/resend_email", hs.ResendEmailVerification)
	sqrg1.DELETE("/resend_email", hs.ResendEmailVerification)
	// Reset password request with email
	rqrg.POST("/reset_password_emails", hs.ResetPasswordEmails)
	rqrg.GET("/reset_password_emails", hs.ResetPasswordEmails)
	rqrg.DELETE("/reset_password_emails", hs.ResetPasswordEmails)
	sqrg2.POST("/reset_password_emails", hs.ResetPasswordEmails)
	sqrg2.GET("/reset_password_emails", hs.ResetPasswordEmails)
	sqrg2.DELETE("/reset_password_emails", hs.ResetPasswordEmails)
	//Change password using reset password email
	rqrg.PUT("/password", hs.ResetPasswordWithEmail)
	rqrg.GET("/password", hs.ResetPasswordWithEmail)
	rqrg.DELETE("/password", hs.ResetPasswordWithEmail)
	sqrg2.PUT("/password", hs.ResetPasswordWithEmail)
	sqrg2.GET("/password", hs.ResetPasswordWithEmail)
	sqrg2.DELETE("/password", hs.ResetPasswordWithEmail)
	//Check paycms status
	sqrg.GET("/self/paycms_status", hs.UpdatePaycmsStatus)
	sqrg.PUT("/self/paycms_status", hs.UpdatePaycmsStatus)
	sqrg.DELETE("/self/paycms_status", hs.UpdatePaycmsStatus)
	sqrg1.GET("/self/paycms_status", hs.UpdatePaycmsStatus)
	sqrg1.PUT("/self/paycms_status", hs.UpdatePaycmsStatus)
	sqrg1.DELETE("/self/paycms_status", hs.UpdatePaycmsStatus)
	dqrg.POST("/devices/", hs.GetUserDevices)
	dqrg.PUT("/devices/", hs.GetUserDevices)
	dqrg.DELETE("/devices/:deviceid", hs.DeleteUserDevices)

	brg := r.Group("/api/users")
	brg.Use(common.ValidateToken())
	brg.GET("/:id/ratings", hs.GetUserRatingsDetailsWithSearchText)
	r.Use(common.ValidateToken())
	r.GET("api/viewactivities/:id/watchingissues", hs.GetUserWAtchingIssues)

}

// RegisterUserUsingEmail -  Creates a new user using email id
// POST /users/register_email
// @Summary Creates a user using email id
// @Description Creates a user using email id
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body RequestRegisterUserUsingEmail true "Raw JSON string"
// @Success 200 {array} User
// @Router /users/register_email [post]
func (hs *HandlerService) RegisterUserUsingEmail(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var templatePath, confirmationURL string
	var registeruseremail RequestRegisterUserUsingEmail
	var emailcheck Emailcheck
	var errorFlag bool
	errorFlag = false
	if err := c.ShouldBindJSON(&registeruseremail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	var emailError ErrorCode
	err := db.Table("user").Where("email = ?", registeruseremail.Email)
	err.Scan(&emailcheck)
	if len(emailcheck.Email) > 0 {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_already_exists", "Specified email already exists."}
	}
	if !common.RegEmail(registeruseremail.Email) && registeruseremail.Email != "" {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_invalid", "Email is invalid."}
	}
	if registeruseremail.Email == "" {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_required", "Email is required."}
	}
	var passwordError PasswordError

	if len(registeruseremail.Password) < 8 || len(registeruseremail.Password) > 255 && registeruseremail.Password != "" {
		errorFlag = true
		passwordError = PasswordError{"error_user_password_length_invalid", "Password length should be between 8 and 255 characters."}
	}
	if registeruseremail.Password == "" {
		errorFlag = true
		passwordError = PasswordError{"error_user_password_required", "Password is required."}
	}

	var langidError LanguageId
	if !(registeruseremail.LanguageId == 1 || registeruseremail.LanguageId == 2) {
		errorFlag = true
		langidError = LanguageId{"PredicateValidator", "Unexpected value."}
	}
	hashedPassword, saltString := common.HashPassword(registeruseremail.Password)
	countryId := int(common.Countrys(registeruseremail.Alpha2code))
	// TODO: Role id is hard coded today. Need to be modified once the strategy of login handling of backoffice and frontoffice users..
	// In current system they differentiating the user role based on basepath of the respective backoffice and frontoffice
	var user User
	if registeruseremail.Source == ".net" {
		user = User{Id: registeruseremail.UserId, Email: registeruseremail.Email, PasswordHash: hashedPassword, LanguageId: registeruseremail.LanguageId, PrivacyPolicy: registeruseremail.PrivacyPolicy, IsAdult: registeruseremail.IsAdult, IsRecommend: registeruseremail.IsRecommend, Version: 2, SaltStored: saltString, UserName: registeruseremail.Email, RoleId: "91f15b92-97fd-e611-814f-0af7afba4acb", RegistrationSource: 1, CountryName: registeruseremail.CountryName, Country: countryId, RegisteredAt: time.Now()}
	} else {
		user = User{Email: registeruseremail.Email, PasswordHash: hashedPassword, LanguageId: registeruseremail.LanguageId, PrivacyPolicy: registeruseremail.PrivacyPolicy, IsAdult: registeruseremail.IsAdult, IsRecommend: registeruseremail.IsRecommend, Version: 2, SaltStored: saltString, UserName: registeruseremail.Email, RoleId: "91f15b92-97fd-e611-814f-0af7afba4acb", RegistrationSource: 1, CountryName: registeruseremail.CountryName, Country: countryId, RegisteredAt: time.Now()}
	}
	var invalid Invalid
	if emailError.Code != "" {
		invalid = Invalid{Email: &emailError}
	}
	if langidError.Code != "" {
		invalid = Invalid{LanguageId: &langidError}
	}
	if passwordError.Code != "" {
		invalid = Invalid{Password: &passwordError}
	}
	if emailError.Code != "" && langidError.Code != "" {
		invalid = Invalid{Email: &emailError, LanguageId: &langidError}
	}
	if emailError.Code != "" && passwordError.Code != "" {
		invalid = Invalid{Email: &emailError, Password: &passwordError}
	}
	if passwordError.Code != "" && langidError.Code != "" {
		invalid = Invalid{LanguageId: &langidError, Password: &passwordError}
	}
	if emailError.Code != "" && langidError.Code != "" && passwordError.Code != "" {
		invalid = Invalid{Email: &emailError, LanguageId: &langidError, Password: &passwordError}
	}

	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}

	// Stroing data into User table and validation
	db.Create(&user)
	// Sending Email Notification
	ConfirmationToken := base64.StdEncoding.EncodeToString([]byte(user.Id))
	currentTime := time.Now().Local().Add(time.Minute * time.Duration(1440))
	timeToString := currentTime.String()
	DateTimeToken := base64.StdEncoding.EncodeToString([]byte(timeToString))
	if registeruseremail.Source != ".net" {
		if registeruseremail.LanguageId == 1 {
			templatePath = "CreateFrontOfficeUserBodyEN.html"
			confirmationURL = os.Getenv("BASE_URL") + "/en/confirm-email?confirmationToken=" + ConfirmationToken + "&dateTimeToken=" + DateTimeToken
		} else {
			templatePath = "CreateFrontOfficeUserBodyAR.html"
			confirmationURL = os.Getenv("BASE_URL") + "/ar/confirm-email?confirmationToken=" + ConfirmationToken + "&dateTimeToken=" + DateTimeToken
		}
		templateData := struct {
			EmailHeadImageUrl    string
			EmailContentImageUrl string
			CallbackUrl          string
		}{
			EmailHeadImageUrl:    string(os.Getenv("EMAILIMAGEBASEURL")) + string(os.Getenv("EMAILHEADIMAGEFILENAME")),
			EmailContentImageUrl: string(os.Getenv("EMAILIMAGEBASEURL")) + string(os.Getenv("EMAILCONTENTIMAGEFILENAME")),
			CallbackUrl:          confirmationURL,
		}
		con, _ := ioutil.ReadFile(templatePath)
		content := string(con)
		fmt.Println(string(content))
		// if err := r.ParseTemplate(templatePath, templateData); err == nil {
		fmt.Println(templateData, templatePath, "rrrrrrrrrrrrrrrrrrr")
		content = strings.Replace(content, "{{.EmailHeadImageUrl}}", templateData.EmailHeadImageUrl, 1)
		content = strings.Replace(content, "{{.EmailContentImageUrl}}", templateData.EmailContentImageUrl, 1)
		content = strings.Replace(content, "{{.CallbackUrl}}", templateData.CallbackUrl, 1)
		/* otp and mail commented need to uncomment in future
		message := template.HTML(content)
		error := common.SendMail(registeruseremail.Email, string(message), "Subject: Welcome to Weyyak!")
		if error != nil {
			fmt.Println("Email has not sent- ", error)
		}*/
	}
	// }
	// Sending details to PayCMS
	CreateRecordPayCMS(db, user)

	c.JSON(http.StatusOK, gin.H{"message": "Confirmation Token Sent to User.", "status": 1})
}

// RegisterUserUsingSMS -  Creates a new user using SMS
// POST /users/register_sms
// @Summary Creates a user using sms
// @Description Creates a user using sms
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body RequestRegisterUserUsingSMS true "Raw JSON string"
// @Success 200 {array} User
// @Router /users/register_sms [post]
func (hs *HandlerService) RegisterUserUsingSMS(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var registerusersms RequestRegisterUserUsingSMS
	var phonenumbercheck PhoneNumber
	var errorFlag bool
	errorFlag = false
	if err := c.ShouldBindJSON(&registerusersms); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	// Split PhoneNumber & collect(callingCode,NationalNumber) from phoneNumber
	num, err := phonenumbers.Parse(registerusersms.PhoneNumber, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	regionNumber := phonenumbers.GetRegionCodeForNumber(num)
	countryCode := phonenumbers.GetCountryCodeForRegion(regionNumber)
	callingCode := fmt.Sprint("+", countryCode)
	nationalNumber := strings.Split(registerusersms.PhoneNumber, callingCode)
	if registerusersms.Silentregistration == true {
		var silentUser SilentUser
		db.Table("user").Select("id,user_name,phone_number").Where("phone_number = ?", registerusersms.PhoneNumber).Find(&silentUser)
		if silentUser.PhoneNumber != "" {
			c.JSON(http.StatusOK, gin.H{"status": 2, "username": silentUser.UserName, "PhoneNumber": silentUser.PhoneNumber, "id": silentUser.Id})
			return
		} else {
			password := randstr.String(8)
			hashedPassword, saltString := common.HashPassword(password)
			user := User{PhoneNumber: registerusersms.PhoneNumber, CallingCode: callingCode, NationalNumber: nationalNumber[1], PhoneNumberConfirmed: true, RegisteredAt: time.Now(), PasswordHash: hashedPassword, LanguageId: registerusersms.LanguageId, PrivacyPolicy: registerusersms.PrivacyPolicy, IsAdult: registerusersms.IsAdult, IsRecommend: registerusersms.IsRecommend, Version: 2, SaltStored: saltString, UserName: registerusersms.PhoneNumber, RoleId: "91f15b92-97fd-e611-814f-0af7afba4acb", RegistrationSource: 4}
			if newSilentUser := db.Create(&user).Error; newSilentUser != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": newSilentUser.Error(), "status": http.StatusInternalServerError})
				return
			}
			if collectSilentUser := db.Table("user").Select("id,user_name,phone_number").Where("phone_number = ?", user.PhoneNumber).Find(&silentUser).Error; collectSilentUser != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": collectSilentUser.Error(), "status": http.StatusInternalServerError})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": 1, "username": silentUser.UserName, "PhoneNumber": silentUser.PhoneNumber, "password": password, "id": silentUser.Id})
			return
		}
	}
	// Error_codes & Validations
	var phoneError phoneNumberError
	// Validations user mobile number should be unique
	if err := db.Table("user").Where("phone_number = ?", registerusersms.PhoneNumber).First(&phonenumbercheck).Error; err == nil {
		errorFlag = true
		phoneError = phoneNumberError{"error_phone_number_registered", "Phone Number Already Exists"}
	}
	if !common.RegMobile(registerusersms.PhoneNumber) && registerusersms.PhoneNumber != "" {
		errorFlag = true
		phoneError = phoneNumberError{"error_phone_number_invalid", "Invalid Phone number"}
	}
	if registerusersms.PhoneNumber == "" {
		errorFlag = true
		phoneError = phoneNumberError{"NotEmptyValidator", "'Phone Number' should not be empty."}
	}
	var passwordError PasswordError
	if len(registerusersms.Password) < 8 || len(registerusersms.Password) > 255 && registerusersms.Password != "" {
		errorFlag = true
		passwordError = PasswordError{"error_user_password_length_invalid", "Password length should be between 8 and 255 characters."}
	}
	if registerusersms.Password == "" {
		errorFlag = true
		passwordError = PasswordError{"error_user_password_required", "Password is required."}
	}
	var langidError LanguageId
	if !(registerusersms.LanguageId == 1 || registerusersms.LanguageId == 2) {
		errorFlag = true
		langidError = LanguageId{"PredicateValidator", "Unexpected value."}
	}

	var invalid Invalid
	if phoneError.Code != "" {
		invalid = Invalid{PhoneNumber: &phoneError}
	}
	if langidError.Code != "" {
		invalid = Invalid{LanguageId: &langidError}
	}
	if passwordError.Code != "" {
		invalid = Invalid{Password: &passwordError}
	}
	if phoneError.Code != "" && langidError.Code != "" {
		invalid = Invalid{PhoneNumber: &phoneError, LanguageId: &langidError}
	}
	if phoneError.Code != "" && passwordError.Code != "" {
		invalid = Invalid{PhoneNumber: &phoneError, Password: &passwordError}
	}
	if passwordError.Code != "" && langidError.Code != "" {
		invalid = Invalid{LanguageId: &langidError, Password: &passwordError}
	}
	if phoneError.Code != "" && langidError.Code != "" && passwordError.Code != "" {
		invalid = Invalid{PhoneNumber: &phoneError, LanguageId: &langidError, Password: &passwordError}
	}
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
	// End of Error_codes & Validations

	hashedPassword, saltString := common.HashPassword(registerusersms.Password)
	// TODO: Role id is hard coded today. Need to be modified once the strategy of login handling of backoffice and frontoffice users..
	// In current system they differentiating the user role based on basepath of the respective backoffice and frontoffice
	type Country struct {
		EnglishName string `json:"english_name"`
		Id          int    `json:"id"`
	}
	var country Country
	fdb := c.MustGet("FDB").(*gorm.DB)
	fdb.Table("country").Select("english_name,id").Where("calling_code=?", callingCode).Find(&country)
	var user User
	if registerusersms.Source == ".net" {
		user = User{Id: registerusersms.UserId, PhoneNumber: registerusersms.PhoneNumber, CallingCode: callingCode, NationalNumber: nationalNumber[1], RegisteredAt: time.Now(), PasswordHash: hashedPassword, LanguageId: registerusersms.LanguageId, PrivacyPolicy: registerusersms.PrivacyPolicy, IsAdult: registerusersms.IsAdult, IsRecommend: registerusersms.IsRecommend, Version: 2, SaltStored: saltString, UserName: registerusersms.PhoneNumber, RoleId: "91f15b92-97fd-e611-814f-0af7afba4acb", RegistrationSource: 4, CountryName: country.EnglishName, Country: country.Id}
	} else {
		user = User{PhoneNumber: registerusersms.PhoneNumber, CallingCode: callingCode, NationalNumber: nationalNumber[1], RegisteredAt: time.Now(), PasswordHash: hashedPassword, LanguageId: registerusersms.LanguageId, PrivacyPolicy: registerusersms.PrivacyPolicy, IsAdult: registerusersms.IsAdult, IsRecommend: registerusersms.IsRecommend, Version: 2, SaltStored: saltString, UserName: registerusersms.PhoneNumber, RoleId: "91f15b92-97fd-e611-814f-0af7afba4acb", RegistrationSource: 4, CountryName: country.EnglishName, Country: country.Id}
	}
	db.Create(&user)

	if registerusersms.Source != ".net" {
		// sending otp to mobille number using AWS SNS Service
		otp := common.EncodeToString(4)
		var language string
		if user.LanguageId == 1 {
			language = `Your OTP to verify your Weyyak account is :` + string(otp)
		} else {
			language = `استخدم الرمز التعريفي ` + string(otp) + `لتفعيل حسابك في وياك`
		}
		// AwsRegion := "us-east-1"
		otpRecord := OtpRecord{Phone: user.PhoneNumber, Message: otp, SentOn: time.Now(), Number: 1}
		// Access_key := "AKIAYOGUWMUMEEQD6CPW"
		// Secret_Key := "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"
		if final := db.Table("otp_record").Create(&otpRecord).Error; final != nil {
			fmt.Println(final.Error(), "While inserting record into otp_record table")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "حدث خطأ ما", "code": "error_server_error", "requestId": randstr.String(32)})
			return
		}
		body := map[string]interface{}{
			"phonenumber": registerusersms.PhoneNumber,
			"Message":     language,
			"CallingCode": callingCode,
		}
		langdetails1, err := common.PostCurlCall("POST", "https://api-backoffice-production.weyyak.com/users/send_otp", body)
		fmt.Println("registerrrrrrrrrrrrrrrrrrrrrrr res", string(langdetails1))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
			return
		}
		// sess, res := session.NewSession(&aws.Config{
		// 	Region:      aws.String(AwsRegion),
		// 	Credentials: credentials.NewStaticCredentials(Access_key, Secret_Key, ""),
		// })
		// if res != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "حدث خطأ ما", "code": "error_server_error", "requestId": randstr.String(32)})
		// 	return
		// }
		// svc := sns.New(sess)
		// params := &sns.PublishInput{
		// 	Message: aws.String(language),
		// 	/* otp and mail commented need to uncomment in future */
		// 	PhoneNumber: aws.String(user.PhoneNumber),
		// }
		// _, sample := svc.Publish(params)
		// if sample != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "حدث خطأ ما", "code": "error_server_error", "requestId": randstr.String(32)})
		// 	return
		// }
	}
	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Otp sent to user"})
	// Sending details to PayCMS
	CreateRecordPayCMS(db, user)
}

func CreateRecordPayCMS(db *gorm.DB, user User) error {
	currentTime := time.Now()
	PostData := make(map[string]string)
	PostData["id"] = user.Id
	PostData["firstName"] = user.FirstName
	PostData["lastName"] = user.LastName
	PostData["email"] = user.Email
	PostData["phoneNumber"] = user.PhoneNumber
	PostData["registeredAt"] = currentTime.Format("2006-01-02 15:04:05")
	PostData["countryName"] = user.CountryName
	payCMSResponse, err := common.PostCurlCall("POST", "https://zpmsapi-preprod.engro.in/payment/registration/insert", PostData)
	fmt.Println("payCMSResponse", string(payCMSResponse))
	// var user User
	if err == nil {
		db.Model(&user).Where("id = ?", user.Id).Update("paycmsstatus", true)
		return nil
	}
	return err
}

// UpdateUserProfile -  Update User Profile
// POST /v1/users/self
// @Summary Update User Profile
// @Description Update User Profile by user id
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Param body body register.UpdateUser true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /v1/users/self [POST]
func (hs *HandlerService) UpdateUserProfile(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	//lgdb := c.MustGet("FDB").(*gorm.DB)
	var user User
	userid := c.MustGet("userid")
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request.", "Status": http.StatusUnauthorized})
		return
	}
	if err := db.Where("id=?", userid).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Record does not exist.Please provide valid User Id.", "Status": http.StatusBadRequest})
		return
	}
	var input UpdateUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	fmt.Println(input.Country, "?????????????")
	if err := db.Table("user").Where("id=?", userid).Update(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}
	if err := db.Table("user").Where("id=?", userid).Update(map[string]interface{}{"google_ads": input.GoogleAds, "google_analytics": input.GoogleAnalytics, "performance": input.Performance, "is_gdpr_accepted": input.IsGdprAccepted, "clever_tap": input.CleverTap, "facebook_ads": input.FacebookAds, "aique": input.Aique, "advertising": input.Advertising, "app_flyer": input.AppFlyer, "firebase": input.Firebase}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
	// } else {
	// 	if c.Request.FormValue("countryId") != "" {
	// 		db.Table("user").Where("id=?", userid).Update("country", c.Request.FormValue("countryId"))
	// 	}
	// 	if input.CountryName != "" {
	// 		fmt.Println("jj")
	// 		var countryId CountryDetails
	// 		lgdb.Table("country").Select("id").Where("english_name=?", input.CountryName).Scan(&countryId)
	// 		db.Table("user").Where("id=?", userid).Update("country", countryId.Id)
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{})
	// }

}

// UserDevices -  Getting User Logged In Devices
// GET /v1/devices/
// @Summary Get User Logged In Devices
// @Description Get User Logged In Devices
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} UserDevicesResponse
// @Router /v1/devices/ [get]
func (hs *HandlerService) GetUserDevices(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var limit, offset, current_page int64

	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["page"] != nil {
		current_page, _ = strconv.ParseInt(c.Request.URL.Query()["page"][0], 10, 64)
	}
	if limit == 0 {
		limit = 50
	}
	offset = current_page * limit
	userId := c.MustGet("userid") //common.USERID
	var userDevices, totalCount []UserDevicesResponse
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var user User
	if finalresult := db.Raw(`SELECT * FROM "user" WHERE id=?`, userId).Scan(&user).Error; finalresult != nil {
		c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}

	if data := db.Raw("SELECT device.name,device.device_id as id,device.platform FROM user_device INNER JOIN device ON user_device.device_id = device.device_id where user_device.user_id=? and token is not null and token != '' ", userId).Limit(limit).Offset(offset).Scan(&userDevices).Error; data != nil {
		if user.LanguageId == 1 {
			c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
			return
		} else {
			c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "حدث خطأ ما", "error_server_error", randstr.String(32)})
			return
		}
	}
	if errCount := db.Raw("SELECT device.name,device.device_id as id,device.platform FROM user_device INNER JOIN device ON user_device.device_id = device.device_id where user_device.user_id=? and token is not null and token != '' ", userId).Scan(&totalCount).Error; errCount != nil {
		if user.LanguageId == 1 {
			c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
			return
		} else {
			c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "حدث خطأ ما", "error_server_error", randstr.String(32)})
			return
		}
	}
	commits := map[string]int{
		"size":   len(totalCount),
		"offset": int(offset),
		"limit":  int(limit),
	}

	c.JSON(http.StatusOK, gin.H{"pagination": commits, "data": userDevices})
}

// UserProfileDetails -  Getting User Profile Details
// GET /v1/users/self
// @Summary Get User Profile Details
// @Description Get User Profile Details
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} userProfileResponse
// @Router /v1/users/self [get]
func (hs *HandlerService) GetUserDetails(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request.", "Status": http.StatusUnauthorized})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	fdb := c.MustGet("FDB").(*gorm.DB)
	userId := c.MustGet("userid") //common.USERID
	var user userProfileResponse
	if data := db.Raw(`SELECT * FROM "user" WHERE id=?`, userId).Scan(&user).Error; data != nil {
		c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}
	var langName LanguageDetails
	if lang := fdb.Raw(`SELECT * FROM "language" WHERE id=?`, user.LanguageId).Scan(&langName).Error; lang != nil {
		if user.LanguageId == 1 {
			c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
			return
		} else {
			c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "حدث خطأ ما", "error_server_error", randstr.String(32)})
			return
		}
	}
	user.LanguageName = langName.EnglishName
	if user.EmailConfirmed == true || user.PhoneNumberConfirmed == true {
		user.VerificationStatus = true
	} else {
		user.VerificationStatus = false
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GetUserDevicesLimit -  Getting User Devices Limit
// GET /devices/limit
// @Summary Get User Devices Limit
// @Description Get User Devices Limit
// @Tags User
// @security Authorization
// @Accept  json
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /devices/limit [get]
func (hs *HandlerService) GetUserDevicesLimit(c *gin.Context) {
	db := c.MustGet("FDB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var errorresponse = common.ServerErrorResponse()
	var limit ApplicationSetting
	var devicelim int64
	field := "value"
	whereparam := "UserDevicesLimit"
	if data := db.Table("application_setting").Select(field).Where("name=?", whereparam).Scan(&limit).Error; data != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	result := limit.Value
	devicelim, _ = strconv.ParseInt(result, 10, 64)
	c.JSON(http.StatusOK, gin.H{"data": devicelim})
}

// UpdateUserDevicesLimitCount -  Updating User Devices Limit Count
// POST /devices/limit
// @Summary Updating User Devices Limit Count
// @Description Updating User Devices Limit Count
// @Tags User
// @security Authorization
// @Accept  json
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /devices/limit [post]
func (hs *HandlerService) UpdateUserDevicesLimitCount(c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	var errorresponse = common.ServerErrorResponse()
	var input UserDevicesLimitCount
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	// Check Post data is json formated or Not
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	validinput := input.UserDevicesLimit
	if validinput > 0 {
		var applicationSetting ApplicationSetting
		applicationSetting.Value = strconv.Itoa(input.UserDevicesLimit)

		if err := fdb.Table("application_setting").Where("name='UserDevicesLimit'").Updates(&applicationSetting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Record Updated Successfully.", "Status": http.StatusOK})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide a valid input limit number.", "Status": http.StatusBadRequest})
	}

}

// RegistrationConfirmation -  update user Status
// POST /user/registration_confirmation
// @Summary Updating User information by conformation
// @Description Updating information by conformation
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body ConfirmEmail true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /user/registration_confirmation [post]
func (hs *HandlerService) RegistrationConfirmation(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	currentTime := time.Now()
	var confirmEmail ConfirmEmail
	var user User
	if err := c.ShouldBindJSON(&confirmEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}

	/* Sync usecase starts here */
	if confirmEmail.Source == ".net" {		
		results := db.Model(&user).Where("id = ?", confirmEmail.UserId).Update("emailConfirmed", true)
		if results.RowsAffected > 0 {
			c.JSON(http.StatusOK, gin.H{"message": "Your Email has been verified.", "Status": 1})
			return
		}

		/* Sync usecase end here */
	} else {
		userId, _ := base64.StdEncoding.DecodeString(confirmEmail.ConfirmationToken)
		Timetoken, _ := base64.StdEncoding.DecodeString(confirmEmail.DateTimeToken)
		tokentime := string(Timetoken)
		layout := "2006-01-02 15:04:05 -0700 MST"
		validTokenTime, err := time.Parse(layout, tokentime)
		if err != nil {
			fmt.Println(err)
		}
		diff := currentTime.Sub(validTokenTime)
		var errorFlag bool
		errorFlag = false
		var confirmToken ConfirmationToken
		var datetimetoken DateTimeToken
		var invalid Invalid
		var count int
		if !(diff < 0) {
			errorFlag = true
			datetimetoken = DateTimeToken{"error_expired_confirmation_token", "Confirmation token is expired."}
			confirmToken = ConfirmationToken{Code: "error_invalid_confirmation_token", Description: "Confirmation Token is Invalid."}
		}
		db.Table("user").Where("id=?", userId).Count(&count)		
		if count < 1 {
			errorFlag = true
			confirmToken = ConfirmationToken{Code: "error_invalid_confirmation_token", Description: "Confirmation Token is Invalid."}

		}
		if confirmEmail.ConfirmationToken == "" {
			errorFlag = true
			confirmToken = ConfirmationToken{Code: "NotEmptyValidator", Description: "Confirmation Token' should not be empty."}
		}
		if !common.ValidTime(tokentime) {
			errorFlag = true
			datetimetoken = DateTimeToken{Code: "error_dateTime_token_invalid", Description: "DateTime Token is Invalid."}
		}
		if confirmEmail.DateTimeToken == "" {
			errorFlag = true
			datetimetoken = DateTimeToken{"NotEmptyValidator", "Date Time Token' should not be empty."}
		}
		if confirmToken.Code != "" {
			invalid = Invalid{ConfirmationToken: &confirmToken}
		}
		if datetimetoken.Code != "" {
			invalid = Invalid{DateTimeToken: &datetimetoken}
		}
		if confirmToken.Code != "" && datetimetoken.Code != "" {
			invalid = Invalid{ConfirmationToken: &confirmToken, DateTimeToken: &datetimetoken}
		}
		var finalErrorResponse FinalErrorResponse
		//invalid := Invalid{ConfirmationToken: &confirmToken, DateTimeToken: &datetimetoken}
		finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
		if errorFlag {
			c.JSON(http.StatusBadRequest, finalErrorResponse)
			return
		}
		if err := db.Model(&user).Where("id = ?", userId).Update("emailConfirmed", true).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Your Email has been verified.", "Status": 1})
		return
	}

}

// UpdatePaycmsStatus -  update paycms status in user table
// POST /users/self/paycms_status
// @Summary Updating Paycms status in user table
// @Description Updating Paycms status in user table
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body PaycmsStatus true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /users/self/paycms_status [post]
func (hs *HandlerService) UpdatePaycmsStatus(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var paycmsStatus PaycmsStatus
	var user User
	if err := c.ShouldBindJSON(&paycmsStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}
	validinput := paycmsStatus.Userid
	if err := db.Model(&user).Where("id = ?", validinput).Update("paycmsstatus", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Your Paycms status has been updated successfully.", "status": 1})
	return
}

// ResendEmailVerification - Resend an email for Registration
// POST /users/resend_email
// @Summary Send an Email for verification
// @Description  Send an Email for verification
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body ResendEmail true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /users/resend_email [post]
func (hs *HandlerService) ResendEmailVerification(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var resendEmail ResendEmail
	var validateEmail ValidateEmail
	var templatePath, confirmationURL string
	var errorFlag bool
	errorFlag = false
	var finalErrorResponse FinalErrorResponse
	var emailError ErrorCode
	if err := c.ShouldBindJSON(&resendEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}
	if err := db.Table("user").Where("email = ?", resendEmail.Email).First(&validateEmail).Error; err != nil {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_unregistered", "User could not be found with this Email"}
	}
	if !common.RegEmail(resendEmail.Email) {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_invalid", "Email is invalid"}
	}
	if resendEmail.Email == "" {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_required", "Email is required."}
	}
	invalid := Invalid{Email: &emailError}
	finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
	ConfirmationToken := base64.StdEncoding.EncodeToString([]byte(validateEmail.ID))
	currentTime := time.Now().Local().Add(time.Minute * time.Duration(1440))
	timeToString := currentTime.String()
	DateTimeToken := base64.StdEncoding.EncodeToString([]byte(timeToString))

	if validateEmail.LanguageId == 1 {
		templatePath = "CreateFrontOfficeUserBodyEN.html"
		confirmationURL = os.Getenv("BASE_URL") + "/en/confirm-email?confirmationToken=" + ConfirmationToken + "&dateTimeToken=" + DateTimeToken
	} else {
		templatePath = "CreateFrontOfficeUserBodyAR.html"
		confirmationURL = os.Getenv("BASE_URL") + "/ar/confirm-email?confirmationToken=" + ConfirmationToken + "&dateTimeToken=" + DateTimeToken
	}
	templateData := struct {
		EmailHeadImageUrl    string
		EmailContentImageUrl string
		CallbackUrl          string
	}{
		EmailHeadImageUrl:    string(os.Getenv("EMAILIMAGEBASEURL")) + string(os.Getenv("EMAILHEADIMAGEFILENAME")),
		EmailContentImageUrl: string(os.Getenv("EMAILIMAGEBASEURL")) + string(os.Getenv("EMAILCONTENTIMAGEFILENAME")),
		CallbackUrl:          confirmationURL,
	}
	con, _ := ioutil.ReadFile(templatePath)
	content := string(con)
	content = strings.Replace(content, "{{.EmailHeadImageUrl}}", templateData.EmailHeadImageUrl, 1)
	content = strings.Replace(content, "{{.EmailContentImageUrl}}", templateData.EmailContentImageUrl, 1)
	content = strings.Replace(content, "{{.CallbackUrl}}", templateData.CallbackUrl, 1)
	/* otp and mail commented need to uncomment in future
	message := template.HTML(content)
	error := common.SendMail(validateEmail.Email, string(message), "Subject: Welcome to Weyyak!")
	if error != nil {
		fmt.Println("Email has not sent- ", error)
	}
	*/
	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Confirmation Token Sent to User."})
}

// ResetPasswordEmails - Reset password with email
// PUT /user/reset_password_emails
// @Summary Send an Email for password reset
// @Description  Send an Email for password reset
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body ResendEmail true "Raw JSON string"
// @Success 202 {array} object c.JSON
// @Router /user/reset_password_emails [put]
func (hs *HandlerService) ResetPasswordEmails(c *gin.Context) {
	if c.Request.Method != http.MethodPut {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var resendEmail ResendEmail
	var validateEmail ValidateEmail
	var templatePath, confirmationURL string
	PasswordToken := common.EncodeToString(6)
	var collectEmailOtp CollectEmailOtp
	if err := c.ShouldBindJSON(&resendEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}
	var errorFlag bool
	errorFlag = false
	var emailError ErrorCode
	details := EmailOtpRecord{Phone: resendEmail.Email, Message: PasswordToken, SentOn: time.Now()}
	if !common.RegEmail(resendEmail.Email) {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_invalid", "Email is invalid"}
	}
	if resendEmail.Email == "" {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_required", "Email is required."}
	}
	if err := db.Table("user").Where("email = ?", resendEmail.Email).First(&validateEmail).Error; err != nil {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_unregistered", "User could not be found with this Email"}
	}
	invalid := Invalid{Email: &emailError}
	finalErrorResponse := FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
	db.Table("otp_record").Select("message").Where("phone = ?", validateEmail.Email).Find(&collectEmailOtp)
	if len(collectEmailOtp.Message) == 0 {
		if final := db.Table("otp_record").Create(&details).Error; final != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": final.Error(), "status": http.StatusInternalServerError})
			return
		}
	} else {
		if result := db.Table("otp_record").Where("phone=(?)", details.Phone).Update(EmailOtpRecord{Message: PasswordToken, SentOn: time.Now()}).Error; result != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error(), "status": http.StatusInternalServerError})
			return
		}
	}
	if resendEmail.Alpha2code == "EG" {
		if validateEmail.LanguageId == 1 {
			templatePath = "FrontOfficeUserResetPasswordBodyEN.html"
			confirmationURL = os.Getenv("EGYPTBASE_URL") + "/en/reset-password?email=" + validateEmail.Email + "&resetPasswordToken=" + PasswordToken
		} else {
			templatePath = "FrontOfficeUserResetPasswordBodyAR.html"
			confirmationURL = os.Getenv("EGYPTBASE_URL") + "/ar/reset-password?email=" + validateEmail.Email + "&resetPasswordToken=" + PasswordToken
		}
	}
	if validateEmail.LanguageId == 1 {
		templatePath = "FrontOfficeUserResetPasswordBodyEN.html"
		confirmationURL = os.Getenv("BASE_URL") + "/en/reset-password?email=" + validateEmail.Email + "&resetPasswordToken=" + PasswordToken
	} else {
		templatePath = "FrontOfficeUserResetPasswordBodyAR.html"
		confirmationURL = os.Getenv("BASE_URL") + "/ar/reset-password?email=" + validateEmail.Email + "&resetPasswordToken=" + PasswordToken
	}

	templateData := struct {
		EmailHeadImageUrl    string
		EmailContentImageUrl string
		CallbackUrl          string
	}{
		EmailHeadImageUrl:    string(os.Getenv("EMAILIMAGEBASEURL")) + string(os.Getenv("EMAILHEADIMAGEFILENAME")),
		EmailContentImageUrl: string(os.Getenv("EMAILIMAGEBASEURL")) + string(os.Getenv("EMAILCONTENTIMAGEFILENAME")),
		CallbackUrl:          confirmationURL,
	}
	con, _ := ioutil.ReadFile(templatePath)
	content := string(con)
	content = strings.Replace(content, "{{@Model.EmailHeadImageUrl}}", templateData.EmailHeadImageUrl, 1)
	content = strings.Replace(content, "{{@Model.EmailContentImageUrl}}", templateData.EmailContentImageUrl, 1)
	content = strings.Replace(content, "{{@Model.ResetPasswordUri}}", templateData.CallbackUrl, 1)
	/* otp and mail commented need to uncomment in future
	message := template.HTML(content)
	error := common.SendMail(validateEmail.Email, string(message), "Subject: Reset Password")
	if error != nil {
		fmt.Println("Email has not sent- ", error)
	}*/

	c.JSON(http.StatusAccepted, gin.H{})
}

// ResetPasswordWithEmail - Reset Password with Email&string
// POST /user/password
// @Summary Send an Email for password reset
// @Description  Send an Email for password reset
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body VerifyEmail true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /user/password [post]
func (hs *HandlerService) ResetPasswordWithEmail(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var verify VerifyEmail
	var collectOtp GetOtpDetails
	var errorFlag bool
	errorFlag = false

	if data := c.ShouldBindJSON(&verify); data != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": data.Error(), "status": http.StatusBadRequest})
		return
	}
	currentTime := time.Now()
	if result := db.Table("otp_record").Select("message,sent_on").Where("phone=(?)", verify.Email).Find(&collectOtp).Error; result != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Reset request is no longer valid. Please try again.", "code": "error_reset_password_failed", "requestId": randstr.String(32)})
		return
	}
	TimeDiffrence := currentTime.Sub(collectOtp.SentOn)
	TimeInMinutes := int(TimeDiffrence.Minutes())
	hashedPassword, saltStored := common.HashPassword(verify.Password)
	/* Sync usecase starts here */
	if verify.Source == ".net" {
		if userresult := db.Table("user").Where("email=?", verify.Email).Update(UpdatePassword{PasswordHash: hashedPassword, SaltStored: saltStored, Version: 2}).Error; userresult != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Reset request is no longer valid. Please try again.", "code": "error_reset_password_failed", "requestId": randstr.String(32)})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": "password changed successfully"})
		return
	}
	/* Sync usecase end here */
	/* Error_codes && Validations */
	var emailError ErrorCode
	if !common.RegEmail(verify.Email) || len(verify.Email) < 0 {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_invalid", "Email is invalid"}
	}
	if verify.Email == "" {
		errorFlag = true
		emailError = ErrorCode{"error_user_email_required", "Email is required."}
	}
	var resetPasswordTokenError ResetPasswordToken
	if verify.ResetPasswordToken == "" {
		errorFlag = true
		resetPasswordTokenError = ResetPasswordToken{"NotEmptyValidator", "'Reset Password Token' should not be empty."}
	}
	var passwordError PasswordError
	if len(verify.Password) < 8 || len(verify.Password) > 255 && verify.Password != "" {
		errorFlag = true
		passwordError = PasswordError{"error_user_password_length_invalid", "Password length should be between 8 and 255 characters."}
	}
	if verify.Password == "" {
		errorFlag = true
		passwordError = PasswordError{"error_user_password_required", "Password is required."}
	}

	var invalid Invalid
	if emailError.Code != "" {
		invalid = Invalid{Email: &emailError}
	}
	if resetPasswordTokenError.Code != "" {
		invalid = Invalid{ResetPasswordToken: &resetPasswordTokenError}
	}
	if passwordError.Code != "" {
		invalid = Invalid{Password: &passwordError}
	}
	if emailError.Code != "" && resetPasswordTokenError.Code != "" {
		invalid = Invalid{Email: &emailError, ResetPasswordToken: &resetPasswordTokenError}
	}
	if emailError.Code != "" && passwordError.Code != "" {
		invalid = Invalid{Email: &emailError, Password: &passwordError}
	}
	if passwordError.Code != "" && resetPasswordTokenError.Code != "" {
		invalid = Invalid{ResetPasswordToken: &resetPasswordTokenError, Password: &passwordError}
	}
	if emailError.Code != "" && resetPasswordTokenError.Code != "" && passwordError.Code != "" {
		invalid = Invalid{Email: &emailError, ResetPasswordToken: &resetPasswordTokenError, Password: &passwordError}
	}
	//End of Error_codes && Validations
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
	if collectOtp.Message == verify.ResetPasswordToken && TimeInMinutes < 30 {
		if userresult := db.Table("user").Where("email=?", verify.Email).Update(UpdatePassword{PasswordHash: hashedPassword, SaltStored: saltStored, Version: 2}).Error; userresult != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Reset request is no longer valid. Please try again.", "code": "error_reset_password_failed", "requestId": randstr.String(32)})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": "password changed successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Reset request is no longer valid. Please try again.", "code": "error_reset_password_failed", "requestId": randstr.String(32)})
	}
}

// Logout -  User Logout
// POST /v1/logout
// @Summary User Logout
// @Description User logout by token id
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body LogoutToken true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /v1/logout [post]
func (hs *HandlerService) Logout(c *gin.Context) {
	/* Authorization */
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	deviceId := c.MustGet("device_id")
	userId := c.MustGet("userid")
	var logouttoken LogoutToken
	var usertoken UserToken
	if err := c.ShouldBindJSON(&logouttoken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}
	/*Update last_activity in user table */
	if userErr := db.Table("user").Where("id=?", c.MustGet("userid")).Update("last_activity_at", time.Now()).Error; userErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	/* Remove device from user_device */
	var userDevice UserDevice
	userDevice.Token = ""
	if deviceErr := db.Table("user_device").Where("device_id=? and user_id=?", deviceId, userId).Delete(&userDevice).Error; deviceErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	/*Remove user related token info */
	refreshtoken := logouttoken.RefreshToken
	var oauthTokens Oauth2Tokens
	if oauthErr := db.Table("oauth2_tokens").Where("refresh=?", refreshtoken).Delete(&oauthTokens).Error; oauthErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	/* Remove token from user token table */
	if tokenErr := db.Table("user_token").Where("token=?", refreshtoken).Delete(&usertoken).Error; tokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}

}

// GeneratePairingCode - Pairing Code generates
// POST /oauth2/device/code
// @Summary Pairing Code Generates
// @Description Pairing Code Generates
// @Tags Device
// @Accept  json
// @Produce  json
// @Param body body RequestPairingCode true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /oauth2/device/code [post]
func (hs *HandlerService) GeneratePairingCode(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var request RequestPairingCode
	var platform Platform
	var response ResponsePairingCode
	fmt.Println(request, "reeeeeeeeeeee")
	// if err := c.ShouldBindJSON(&request); err != nil {
	// 	fmt.Println("error in code.....", err.Error())
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
	// 	return
	// }
	request.DeviceId = c.Request.FormValue("deviceId")
	request.DeviceName = c.Request.FormValue("deviceName")
	request.DevicePlatform = c.Request.FormValue("devicePlatform")
	var errorFlag bool
	errorFlag = false
	var devicePlatform DevicePlatform
	if request.DevicePlatform != "smart_tv" && request.DevicePlatform != "apple_tv" {
		errorFlag = true
		devicePlatform = DevicePlatform{"error_invalid_value", "Invalid value."}
	}
	if request.DevicePlatform == "" {
		errorFlag = true
		devicePlatform = DevicePlatform{"NotEmptyValidator", "'Device Platform' should not be empty."}
	}
	invalid := Invalid{DevicePlatform: &devicePlatform}
	finalErrorResponse := FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
	// generate device and user code
	userCode := randstr.String(4)
	userCode = strings.ToUpper((userCode))

	deviceCode := randstr.String(45)
	deviceCode = strings.ToUpper(deviceCode)
	var count int
	expireTime := time.Now().Local().Add(time.Minute * time.Duration(30))
	pairingCode := PairingCode{DeviceId: request.DeviceId, DeviceCode: deviceCode, UserCode: userCode, CreatedAt: time.Now(), ExpiresAt: expireTime}
	if idresults := db.Table("pairing_code").Select("device_id").Where("device_id=?", request.DeviceId).Count(&count).Error; idresults != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": idresults.Error(), "status": http.StatusInternalServerError})
		return
	}
	if platformError := db.Raw("select platform_id,name from public.platform where name = ?", request.DevicePlatform).Scan(&platform).Error; platformError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": platformError.Error(), "status": http.StatusInternalServerError})
		return
	}
	device := Device{DeviceId: request.DeviceId, Name: request.DeviceName, Platform: strconv.Itoa(platform.PlatformId), CreatedAt: time.Now()}
	if count == 0 {
		// create a record in pairing code table
		if pairingCodeError := db.Create(&pairingCode).Error; pairingCodeError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": pairingCodeError.Error(), "status": http.StatusInternalServerError})
			return
		}
		// create a record in device table
		if deviceResult := db.Create(&device).Error; deviceResult != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": deviceResult.Error(), "status": http.StatusInternalServerError})
			return
		}
	} else {
		// If DeviceId already exist in table then update a record in pairing code table
		if result := db.Table("pairing_code").Where("device_id=(?)", request.DeviceId).Update(PairingCode{DeviceCode: deviceCode, UserCode: userCode, CreatedAt: time.Now(), ExpiresAt: expireTime}).Error; result != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error(), "status": http.StatusInternalServerError})
			return
		}
		// If DeviceId already exist in table then update a record in device table
		if result := db.Table("device").Where("device_id=(?)", request.DeviceId).Update(Device{Name: request.DeviceName, Platform: strconv.Itoa(platform.PlatformId), CreatedAt: time.Now()}).Error; result != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error(), "status": http.StatusInternalServerError})
			return
		}
	}
	response.DeviceCode = deviceCode
	response.UserCode = userCode
	response.VerificationUri = "https://webuat.weyyak.com/tv-pair"
	response.ExpiresIn = 1800
	response.Interval = 5
	c.JSON(http.StatusOK, response)
}

// VerifyPairingCode - Verify pairing code
// POST /oauth2/device/auth
// @Summary Verify pairing code
// @Description Verify pairing code
// @Tags Device
// @Accept  json
// @Produce  json
// @security Authorization
// @Param body body VerifyPairingCode true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /oauth2/device/auth [post]
func (hs *HandlerService) VerifyPairingCode(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var verifyPairCode VerifyPairingCode
	userId := c.MustGet("userid")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}

	if err := c.ShouldBindJSON(&verifyPairCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	var paircode PairingCode
	var newdevice Device
	var deviceid []DeviceIds
	deviceresult := db.Debug().Table("user_device").Select("device_id").Where("user_id=? and token is not null", userId)
	deviceresult.Scan(&deviceid)
	var devicecount []string
	for _, value := range deviceid {
		devicecount = append(devicecount, value.DeviceId)
	}
	if len(devicecount) > 5 {
		c.JSON(http.StatusBadRequest, FinalResponse{"invalid_grant", "Maximum number of allowed devices was reached.", "error_device_limit_reached", randstr.String(32)})
		return
	}
	// Verify the usercode and subscription date across the table before expiration
	row := db.Raw("select * from pairing_code pc where lower(user_code) =?", strings.ToLower(verifyPairCode.UserCode.UserCode))
	row.Scan(&paircode)
	a := time.Now()
	differnce := paircode.ExpiresAt.Sub(a)
	d := int(differnce.Seconds())
	fmt.Println("Date - ", d)
	if d <= 1800 && d > 0 {
		if paircode.DeviceId != "" {
			userDevice := UserDevice{DeviceId: paircode.DeviceId, UserId: userId.(string)}
			// userdevice record to be created with device and user details
			db.Create(&userDevice)
			// update subscription date in pairing code table
			db.Table("pairing_code").Where("device_id=?", paircode.DeviceId).Update("subscription_date", verifyPairCode.UserCode.SubscriptionDate)
			db.Table("pairing_code").Where("device_id=?", paircode.DeviceId).Update("user_id", userId.(string))
			// update user lead with device name in user table if createdat and subscription date are equal
			newdevicedata := db.Table("device").Select("name").Where("device_id=?", paircode.DeviceId)
			newdevicedata.Scan(&newdevice)
			newdate := paircode.CreatedAt
			new := newdate.Format("2006-01-02")
			anotherdate := verifyPairCode.UserCode.SubscriptionDate
			if new == strings.Replace(anotherdate, "/", "-", 3) {
				db.Table("user").Where("id=?", userId).Update(User{UserLead: newdevice.Name})
			}
			c.JSON(http.StatusOK, gin.H{"message": "Verified pairing code", "Status": http.StatusOK})
		}
	} else {
		c.JSON(http.StatusBadRequest, FinalResponse{"error_pairing_code_invalid", "Pairing code is invalid or expired.", "error_pairing_code_invalid", randstr.String(32)})
		return
	}

}

// Sending otp to User -   Sending otp to User
// POST /users/send_otp
// @Summary  Sending otp to User
// @Description  Sending otp to User
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body Number true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /users/send_otp [post]
func (hs *HandlerService) SendOtp(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var phn Number
	var language string
	var mes Message
	fmt.Println("langiage",language)
	// var newnumber string
	var finalErrorResponse FinalErrorResponse
	var invalid Invalid
	otp := common.EncodeToString(4)
	// AwsRegion := os.Getenv("AWS_SNS_REGION")
	var errorFlag bool
	errorFlag = false
	if number := c.ShouldBindJSON(&phn); number != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": number.Error(), "status": http.StatusBadRequest})
		return
	}
	var phonenumbercheck PhoneNumber
	var phoneError phoneNumberError

	db.Table("public.user").Where("phone_number = ?", phn.Phone).Find(&phonenumbercheck)
	fmt.Println("################", phonenumbercheck)
	fmt.Println(phonenumbercheck.PhoneNumberConfirmed, "confiremd")
	var OTPuser Users
	db.Table("public.user").Where("phone_number = ?", phn.Phone).Find(&OTPuser)
	

	if phn.RequestType == "nm" || phn.RequestType == "up" {
		
		if phonenumbercheck.PhoneNumber != "" && phonenumbercheck.PhoneNumberConfirmed {
			phoneError = phoneNumberError{"error_phone_number_registered", "Phone Number Already Exists"}
			invalid = Invalid{PhoneNumber: &phoneError}
			finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
			c.JSON(http.StatusBadRequest, finalErrorResponse)
			return
		}
	} /*else if phn.RequestType == "up" {
		if phonenumbercheck.PhoneNumber != "" {
			phoneError = phoneNumberError{"error_phone_number_registered", "Phone Number Already Exists"}
			invalid = Invalid{PhoneNumber: &phoneError}
			finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
			c.JSON(http.StatusBadRequest, finalErrorResponse)
			return
		}
	}*/
	// if phn.RequestType == "up" {
	// 	c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error.", "error_server_error", randstr.String(32)})
	// 	return
	// }
	if len(phn.Phone) == 0 {
		errorFlag = true
		phoneError = phoneNumberError{"NotEmptyValidator", "'Phone Number' should not be empty."}
	}
	var requesttype RequestType
	fmt.Println(phn.RequestType, ";;;;;;;;;;;")

	if !(phn.RequestType == "nm" || phn.RequestType == "fp" || phn.RequestType == "up") {
		errorFlag = true
		requesttype = RequestType{"error_request_type_invalid", "Request Type Invalid. Must be 'nm' or 'fp' or 'up'"}
		fmt.Println(requesttype, "............")
	}
	if phn.RequestType == "" {
		errorFlag = true
		requesttype = RequestType{"NotEmptyValidator", "'Request Type' should not be empty."}
	}
	var userlangdetails UserLangDetails
	// var countryCodeDetails CountryCodeDetails
	var phnConfirmed PhnConfirmed

	if phn.RequestType == "fp" {
		if strings.Contains(phn.Phone, "+") {
			langdetails := db.Table("user").Select("id,language_id,phone_number").Where("phone_number=(?)", phn.Phone)
			langdetails.Scan(&userlangdetails)
			
		} else {
			// countrydetails := db.Table("user").Select("calling_code").Where("user_name like '%+"+phn.Phone+"%' or national_number=?", phn.Phone)
			// countrydetails.Scan(&countryCodeDetails)
			// newnumber = countryCodeDetails.CallingCode + phn.Phone
			// fmt.Println(newnumber, "///////////")
			languagedetails := db.Table("user").Select("id,language_id,phone_number").Where("USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ?", phn.Phone, phn.Phone, phn.Phone, ("+" + phn.Phone), ("+" + phn.Phone), ("+" + phn.Phone), strings.TrimLeft(phn.Phone, "0"), strings.TrimLeft(phn.Phone, "0"), strings.TrimLeft(phn.Phone, "0"))
			languagedetails.Scan(&userlangdetails)
		}

	} else {
		db.Table("user").Select("phone_number_confirmed").Where("USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ?", phn.Phone, phn.Phone, phn.Phone, ("+" + phn.Phone), ("+" + phn.Phone), ("+" + phn.Phone), strings.TrimLeft(phn.Phone, "0"), strings.TrimLeft(phn.Phone, "0"), strings.TrimLeft(phn.Phone, "0")).Find(&phnConfirmed)
		if phnConfirmed.PhoneNumberConfirmed == false {
			langdetails := db.Table("user").Select("id,language_id,phone_number").Where("USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ?", phn.Phone, phn.Phone, phn.Phone, ("+" + phn.Phone), ("+" + phn.Phone), ("+" + phn.Phone), strings.TrimLeft(phn.Phone, "0"), strings.TrimLeft(phn.Phone, "0"), strings.TrimLeft(phn.Phone, "0"))
			langdetails.Scan(&userlangdetails)
		} else {
			errorFlag = true
			phoneError = phoneNumberError{"error_phone_number_verified", "Phone Number Is Already Verified"}
			fmt.Println(phoneError, "////////")
		}
	}

	if len(userlangdetails.Id) == 0 && phnConfirmed.PhoneNumberConfirmed == false && phn.RequestType == "nm" {
		errorFlag = true
		phoneError = phoneNumberError{"error_phone_number_unregistered", "Phone Number Not Found"}
	}
	if len(userlangdetails.Id) == 0 && phn.RequestType == "fp" {
		errorFlag = true
		phoneError = phoneNumberError{"error_phone_number_unregistered", "Phone Number Not Found"}
	}
	if phoneError.Code != "" {
		invalid = Invalid{PhoneNumber: &phoneError}
	}
	if requesttype.Code != "" {
		invalid = Invalid{RequestType: &requesttype}
	}
	if requesttype.Code != "" && phoneError.Code != "" {
		invalid = Invalid{RequestType: &requesttype, PhoneNumber: &phoneError}
	}
	finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
	if phn.RequestType == "up" || phn.RequestType == "nm" {
		if userlangdetails.LanguageId == 1 {
			language = `Your OTP to verify your Weyyak account is :` + string(otp)
		} else {
			language = `استخدم الرمز التعريفي ` + string(otp) + `لتفعيل حسابك في وياك`
		}
	} else {
		if userlangdetails.LanguageId == 1 {
			language = `Forget your password? your OTP to reset your weyyak password is :` + string(otp)
		} else {
			language = string(otp) + `نسيت كلمة السر؟ الرمز التعريفي لتجديد كلمة السر هو`
		}
	}
	user := UserDetails{Phone: phn.Phone, Message: otp, SentOn: time.Now(), Number: 1}
	body := map[string]interface{}{
		"phonenumber": phn.Phone,
		"Message":     language,
		"CallingCode": OTPuser.CallingCode,
	}
	langdetails1, err := common.PostCurlCall("POST", "https://api-backoffice-production.weyyak.com/users/send_otp", body)
	fmt.Println("lang2@@@@@@@@@@@@@@@@", string(langdetails1))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	// Access_key := os.Getenv("ACCESS_SECRET")
	// Secret_Key := os.Getenv("REFRESH_SECRET")
	// fmt.Println("region............", AwsRegion)
	// sess, res := session.NewSession(&aws.Config{
	// 	Region:      aws.String(AwsRegion),
	// 	Credentials: credentials.NewStaticCredentials(Access_key, Secret_Key, ""),
	// })
	// if res != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": res.Error(), "status": http.StatusInternalServerError})
	// 	return
	// }

	// svc := sns.New(sess)
	var phonenumberfinal string
	if phn.RequestType == "fp" {
		if strings.Contains(phn.Phone, "+") {
			phonenumberfinal = phn.Phone
		} else {
			phonenumberfinal = "+" + string(phn.Phone) //newnumber
		}
	} else {
		phonenumberfinal = phn.Phone
	}
	fmt.Println(phonenumberfinal, userlangdetails.PhoneNumber, "kkkkkkkk")
	// params := &sns.PublishInput{
	// 	Message: aws.String(language),
	// 	/* otp and mail commented need to uncomment in future */
	// 	PhoneNumber: aws.String(userlangdetails.PhoneNumber),
	// }
	// // params for grant type up
	// var paramsfp *sns.PublishInput
	// if phn.RequestType == "up" {
	// 	paramsfp = &sns.PublishInput{
	// 		Message: aws.String(language),
	// 		/* otp and mail commented need to uncomment in future */
	// 		PhoneNumber: aws.String(phonenumberfinal),
	// 	}
	// }
	var forcount UserDetails
	data := db.Table("otp_record").Select("message,phone").Where("phone=?", userlangdetails.PhoneNumber)
	fmt.Println("@@@@@@@@@@@@@@@@", userlangdetails.PhoneNumber)
	data.Scan(&mes)
	if len(mes.Message) == 0 || mes.Phone == "" {
		if final := db.Table("otp_record").Create(&user).Error; final != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": final.Error(), "status": http.StatusInternalServerError})
			return
		}
	} else {
		db.Table("otp_record").Select("number").Where("phone=(?)", userlangdetails.PhoneNumber).Find(&forcount)
		if forcount.Number <= 6 {
			if result := db.Table("otp_record").Where("phone=(?)", userlangdetails.PhoneNumber).Update(UserDetails{Message: otp, SentOn: time.Now(), Number: forcount.Number + 1}).Error; result != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error(), "status": http.StatusInternalServerError})
				return
			}
		}
	}
	// if err := db.Table("otp_record").Select("number").Where("phone=(?)", phn.Phone).Find(&forcount).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
	// 	return
	// }
	if forcount.Number < 6 {
		if phn.RequestType == "up" {
			langdetails1, err := common.PostCurlCall("POST", "https://api-backoffice-production.weyyak.com/users/send_otp", body)
			fmt.Println("lang2@@@@@@@@@@@@@@@@", string(langdetails1))
				if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
				return
			
	}
			
		} else {
			langdetails1, err := common.PostCurlCall("POST", "https://api-backoffice-production.weyyak.com/users/send_otp", body)
			fmt.Println("lang2@@@@@@@@@@@@@@@@", string(langdetails1))
				if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
				return
			
	}
		}
	} else {
		type SentTime struct {
			SentOn time.Time `json:"sent_on"`
		}
		var senttime SentTime
		db.Table("otp_record").Select("sent_on").Where("phone=(?)", userlangdetails.PhoneNumber).Find(&senttime)
		fmt.Println("sent time", senttime)
		lastSent := senttime.SentOn.Format("2006-01-02")
		currentTime := time.Now().Format("2006-01-02")
		//senttime.SentOn.Before(time.Now())
		if lastSent < currentTime {
			fmt.Println("inside the time loop")
			langdetails1, err := common.PostCurlCall("POST", "https://api-backoffice-production.weyyak.com/users/send_otp", body)
			fmt.Println("lang2@@@@@@@@@@@@@@@@", string(langdetails1))
				if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
				return
			
	}
			if result := db.Table("otp_record").Where("phone=(?)", userlangdetails.PhoneNumber).Update(UserDetails{Message: otp, SentOn: time.Now(), Number: 1}).Error; result != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error(), "status": http.StatusInternalServerError})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, FinalResponse{Error: "not_found", Description: "You have exceeded the maximum daily limit. Please try again after 24 hours", Code: "", RequestId: randstr.String(32)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "OTP Sent to User"})
}

// Verify registration by sending otp  -  Verify registration by sending otp
// POST /users/verify_otp
// @Summary  Sending otp to User
// @Description  Sending otp to User
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body Valid true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /users/verify_otp [post]
func (hs *HandlerService) VerifyOtp(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var verify Valid
	var details Details
	var errorFlag bool
	errorFlag = false
	var otpValidator OtpValidator
	if data := c.ShouldBindJSON(&verify); data != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": data.Error(), "status": http.StatusBadRequest})
		return
	}
	/* Sync usecase starts here */
	if verify.Source == ".net" {
		if confirmError := db.Debug().Table("user").Where("phone_number=(?)", verify.PhoneNumber).Update("phone_number_confirmed", true).Error; confirmError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": confirmError.Error(), "Status": http.StatusInternalServerError})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Otp Verified"})
		return
	}
	/* Sync usecase end here */

	var phoneError PhoneNumberError
	if len(verify.PhoneNumber) == 0 {
		errorFlag = true
		phoneError = PhoneNumberError{"NotEmptyValidator", "'Phone Number' should not be empty."}
	}
	if result := db.Table("otp_record").Select("message,sent_on").Where("phone=(?)", verify.PhoneNumber).Find(&details).Error; result != nil {
		if len(verify.PhoneNumber) != 0 {
			c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
			return
		}
	}
	a := time.Now()
	b := a.Sub(details.SentOn)
	d := int(b.Minutes())

	if d > 30 {
		errorFlag = true
		//	c.JSON(http.StatusBadRequest, gin.H{"message": "otp time expired."})
		otpValidator = OtpValidator{"error_otp_expired", "OTP Expired"}
	} else {
		if details.Message == verify.Message && d <= 30 {

			db.Debug().Table("user").Where("phone_number=(?)", verify.PhoneNumber).Update("phone_number_confirmed", true)
			c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Otp Verified"})
			return
		} else if len(verify.Message) == 0 {
			errorFlag = true
			otpValidator = OtpValidator{"NotEmptyValidator", "'OTP' should not be empty"}
		} else {
			errorFlag = true
			otpValidator = OtpValidator{"error_otp_invalid", "OTP Not Matched"}
		}
	}
	var invalid Invalid
	if otpValidator.Code != "" {
		invalid = Invalid{OtpValidator: &otpValidator}
	}
	if phoneError.Code != "" {
		invalid = Invalid{PhoneNumberError: &phoneError}
	}
	if otpValidator.Code != "" && phoneError.Code != "" {
		invalid = Invalid{OtpValidator: &otpValidator, PhoneNumberError: &phoneError}
	}
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
}

// Change Password by verifying otp  -  Change Password by verifying otp
// POST /users/password_otp
// @Summary  Change Password by verifying otp
// @Description  Change Password by verifying otp
// @Tags User
// @Accept  json
// @Produce  json
// @Param body body VerifyDetails true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /users/password_otp [post]
func (hs *HandlerService) ForgotPasswordOtp(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var verify VerifyDetails
	var details Details
	var errorFlag bool
	errorFlag = false
	var otpValidator OtpValidator
	if data := c.ShouldBindJSON(&verify); data != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": data.Error(), "status": http.StatusBadRequest})
		return
	}
	hashedPassword, saltStored := common.HashPassword(verify.Password)

	/* Sync usecase starts here */
	if verify.Source == ".net" {
		if userresult := db.Table("user").Where("national_number=(?)", verify.PhoneNumber).Update(UpdatePassword{PasswordHash: hashedPassword, SaltStored: saltStored, Version: 2}).Error; userresult != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": userresult.Error(), "status": http.StatusInternalServerError})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Your password has been changed successfully"})
		return
	}
	/* Sync usecase end here */

	var phoneError PhoneNumberError
	if len(verify.PhoneNumber) == 0 {
		errorFlag = true
		phoneError = PhoneNumberError{"NotEmptyValidator", "'Phone Number' should not be empty."}
	}
	var passswordvalidate NewPasswordValidate
	if len(verify.Password) < 8 {
		errorFlag = true
		passswordvalidate = NewPasswordValidate{"error_user_password_length_invalid", "Password length should be between 8 and 255 characters"}
	}
	if len(verify.Password) == 0 {
		errorFlag = true
		passswordvalidate = NewPasswordValidate{"error_user_password_required", "Password is required"}
	}
	if len(verify.Message) == 0 {
		errorFlag = true
		otpValidator = OtpValidator{"NotEmptyValidator", "'OTP' should not be empty"}
	}
	if result := db.Table("otp_record").Select("message,sent_on").Where("phone=(?)", verify.PhoneNumber).Find(&details).Error; result != nil {
		if len(verify.PhoneNumber) != 0 {
			c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
			return
		}
	}

	var invalid Invalid
	if passswordvalidate.Code != "" {
		invalid = Invalid{NewPasswordValidate: &passswordvalidate}
	}
	if phoneError.Code != "" {
		invalid = Invalid{PhoneNumberError: &phoneError}
	}
	if otpValidator.Code != "" {
		invalid = Invalid{OtpValidator: &otpValidator}
	}
	if phoneError.Code != "" && passswordvalidate.Code != "" {
		invalid = Invalid{PhoneNumberError: &phoneError, NewPasswordValidate: &passswordvalidate}
	}
	if otpValidator.Code != "" && passswordvalidate.Code != "" {
		invalid = Invalid{OtpValidator: &otpValidator, NewPasswordValidate: &passswordvalidate}
	}
	if otpValidator.Code != "" && phoneError.Code != "" {
		invalid = Invalid{OtpValidator: &otpValidator, PhoneNumberError: &phoneError}
	}

	var finalErrorResponse FinalErrorResponse

	finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}

	a := time.Now()
	b := a.Sub(details.SentOn)
	d := int(b.Minutes())

	if d > 30 {
		errorFlag = true
		otpValidator = OtpValidator{"error_otp_expired", "OTP Expired"}
		if passswordvalidate.Code != "" {
			invalid = Invalid{NewPasswordValidate: &passswordvalidate}
		}
		if phoneError.Code != "" {
			invalid = Invalid{PhoneNumberError: &phoneError}
		}
		if otpValidator.Code != "" {
			invalid = Invalid{OtpValidator: &otpValidator}
		}
		if phoneError.Code != "" && passswordvalidate.Code != "" {
			invalid = Invalid{PhoneNumberError: &phoneError, NewPasswordValidate: &passswordvalidate}
		}
		if otpValidator.Code != "" && passswordvalidate.Code != "" {
			invalid = Invalid{OtpValidator: &otpValidator, NewPasswordValidate: &passswordvalidate}
		}
		if otpValidator.Code != "" && phoneError.Code != "" {
			invalid = Invalid{OtpValidator: &otpValidator, PhoneNumberError: &phoneError}
		}
		if otpValidator.Code != "" && phoneError.Code != "" && passswordvalidate.Code != "" {
			invalid = Invalid{OtpValidator: &otpValidator, PhoneNumberError: &phoneError, NewPasswordValidate: &passswordvalidate}
		}

		finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
		if errorFlag {
			c.JSON(http.StatusBadRequest, finalErrorResponse)
			return
		}
	} else {
		if details.Message == verify.Message && d < 30 {

			if userresult := db.Table("user").Where("national_number=(?)", verify.PhoneNumber).Update(UpdatePassword{PasswordHash: hashedPassword, SaltStored: saltStored, Version: 2}).Error; userresult != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": userresult.Error(), "status": http.StatusInternalServerError})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Your password has been changed successfully"})
		} else {
			errorFlag = true
			otpValidator = OtpValidator{"error_otp_invalid", "OTP Not Matched"}
			if passswordvalidate.Code != "" {
				invalid = Invalid{NewPasswordValidate: &passswordvalidate}
			}
			if phoneError.Code != "" {
				invalid = Invalid{PhoneNumberError: &phoneError}
			}
			if otpValidator.Code != "" {
				invalid = Invalid{OtpValidator: &otpValidator}
			}
			if phoneError.Code != "" && passswordvalidate.Code != "" {
				invalid = Invalid{PhoneNumberError: &phoneError, NewPasswordValidate: &passswordvalidate}
			}
			if otpValidator.Code != "" && passswordvalidate.Code != "" {
				invalid = Invalid{OtpValidator: &otpValidator, NewPasswordValidate: &passswordvalidate}
			}
			if otpValidator.Code != "" && phoneError.Code != "" {
				invalid = Invalid{OtpValidator: &otpValidator, PhoneNumberError: &phoneError}
			}
			if otpValidator.Code != "" && phoneError.Code != "" && passswordvalidate.Code != "" {
				invalid = Invalid{OtpValidator: &otpValidator, PhoneNumberError: &phoneError, NewPasswordValidate: &passswordvalidate}
			}
			finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
			if errorFlag {
				c.JSON(http.StatusBadRequest, finalErrorResponse)
				return
			}
		}
	}
}

// Change PhoneNumber by verifying otp  -  Change PhoneNumber by verifying otp
// POST /users/self/phone_number
// @Summary  Change PhoneNumber by verifying otp
// @Description  Change PhoneNumber by verifying otp
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Param body body Verify true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /users/self/phone_number [post]
func (hs *HandlerService) PhonenumberChangeOtp(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var verify Verify
	var details Details
	var errorFlag bool
	errorFlag = false
	var otpValidator OtpValidator
	userId := c.MustGet("userid")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	if data := c.ShouldBindJSON(&verify); data != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": data.Error(), "status": http.StatusBadRequest})
		return
	}
	var phoneError PhoneNumberError
	if len(verify.PhoneNumber) == 0 {
		errorFlag = true
		phoneError = PhoneNumberError{"NotEmptyValidator", "'Phone Number' should not be empty."}
	}
	if len(verify.Message) == 0 {
		errorFlag = true
		otpValidator = OtpValidator{"NotEmptyValidator", "'OTP' should not be empty"}
	}
	num, err := phonenumbers.Parse(verify.PhoneNumber, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	regionNumber := phonenumbers.GetRegionCodeForNumber(num)
	countryCode := phonenumbers.GetCountryCodeForRegion(regionNumber)
	callingCode := fmt.Sprint("+", countryCode)
	nationalnumber := strings.Split(verify.PhoneNumber, callingCode)
	if result := db.Table("otp_record").Select("message,sent_on").Where("phone=(?)", verify.PhoneNumber).Find(&details).Error; result != nil {
		c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}
	/* Sync usecase starts here */
	if verify.Source == ".net" {
		if userresult := db.Table("user").Where("id=(?)", userId).Update(ChangePhoneNumber{PhoneNumber: verify.PhoneNumber, CallingCode: callingCode, NationalNumber: nationalnumber[1], PhoneNumberConfirmed: true}).Error; userresult != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": userresult.Error(), "status": http.StatusInternalServerError})

		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Your phone number has been updated successfully."})
		return
	}
	/* Sync usecase end here */
	a := time.Now()
	b := a.Sub(details.SentOn)
	d := int(b.Minutes())
	if d > 30 {
		errorFlag = true
		otpValidator = OtpValidator{"error_otp_expired", "OTP Expired"}
	} else {
		if details.Message == verify.Message && d < 30 {
			if userresult := db.Table("user").Where("id=(?)", userId).Update(ChangePhoneNumber{PhoneNumber: verify.PhoneNumber, CallingCode: callingCode, NationalNumber: nationalnumber[1], PhoneNumberConfirmed: true}).Error; userresult != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": userresult.Error(), "status": http.StatusInternalServerError})

			}
			c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Your phone number has been updated successfully."})

		} else {
			errorFlag = true
			otpValidator = OtpValidator{"error_otp_invalid", "OTP Not Matched"}
		}
	}

	var invalid Invalid
	if otpValidator.Code != "" {
		invalid = Invalid{OtpValidator: &otpValidator}
	}
	if phoneError.Code != "" {
		invalid = Invalid{PhoneNumberError: &phoneError}
	}
	if otpValidator.Code != "" && phoneError.Code != "" {
		invalid = Invalid{OtpValidator: &otpValidator, PhoneNumberError: &phoneError}
	}
	var finalErrorResponse FinalErrorResponse

	finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
}

// Change Password by confirming old password  -  Change Password by confirming old password
// POST /users/self/Password
// @Summary  Change Password by confirming old password
// @Description  Change Password by confirming old password
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Param body body Password true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /users/self/password [post]
func (hs *HandlerService) ChangePassword(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var pwd Password
	var pwddetails PasswordDetails
	var errorFlag bool
	errorFlag = false
	userId := c.MustGet("userid")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}

	if res := c.ShouldBindJSON(&pwd); res != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": res.Error(), "status": http.StatusBadRequest})
		return
	}
	var pwdvalidate ValidatePassword
	if len(pwd.PasswordHash) == 0 {
		errorFlag = true
		pwdvalidate = ValidatePassword{"NotEmptyValidator", "'Old Password' should not be empty."}
	}
	var newpasswordvalidate NewPasswordValidate
	fmt.Println(len(pwd.NewPassword), ".................")
	if len(pwd.NewPassword) < 8 {
		errorFlag = true
		newpasswordvalidate = NewPasswordValidate{"error_user_password_length_invalid", "Password length should be between 8 and 255 characters"}
	}
	if len(pwd.NewPassword) == 0 {
		errorFlag = true
		newpasswordvalidate = NewPasswordValidate{"error_user_password_required", "Password is required."}
	}
	hashedPassword, saltstored := common.HashPassword(pwd.NewPassword)
	if pwdresult := db.Table("user").Select("password_hash,version,salt_stored").Where("id=(?)", userId).Find(&pwddetails).Error; pwdresult != nil {
		c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}
	var invalid Invalid
	if pwdvalidate.Code != "" {
		invalid = Invalid{ValidatePassword: &pwdvalidate}
	}
	if newpasswordvalidate.Code != "" {
		invalid = Invalid{NewPasswordValidate: &newpasswordvalidate}
	}
	if pwdvalidate.Code != "" && newpasswordvalidate.Code != "" {
		invalid = Invalid{ValidatePassword: &pwdvalidate, NewPasswordValidate: &newpasswordvalidate}
	}
	var finalErrorResponse FinalErrorResponse
	finalErrorResponse = FinalErrorResponse{"invalid_request", "Validation failed.", "error_validation_failed", randstr.String(32), invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}

	decryptpassword := common.VerifyHashPassword(pwddetails.PasswordHash, pwd.PasswordHash, pwddetails.Version, pwddetails.SaltStored)
	if decryptpassword {
		if finalres := db.Table("user").Where("id=(?)", userId).Update(UpdatePassword{PasswordHash: hashedPassword, SaltStored: saltstored, Version: 2}).Error; finalres != nil {
			c.JSON(http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Your password has been updated successfully."})
		return
	} else {
		c.JSON(http.StatusBadRequest, FinalResponse{"invalid_grant", "The username or password is incorrect", "error_user_invalid_credentials", randstr.String(32)})
		return
	}
}

// TwitterUserToken -  Getting User token for twitter
// GET /:lang/usertoken
// @Summary Get User token twitter
// @Description Get User token twitter
// @Tags User
// @Accept  json
// @Produce  json
// @Param lang path string true "lang"
// @Param callback query string false "callback url"
// @Success 200 {array} object c.JSON
// @Router /{lang}/usertoken [get]
func (hs *HandlerService) TwitterUserToken(c *gin.Context) {
	config := oauth1.Config{
		ConsumerKey:    "z4CAfGU3ToZlln6v440wRoA4x",
		ConsumerSecret: "Oip8cfbwSDykJ2OdfY5iWMjbpxTGK4t1UqHSf5SVCLQuXgUGmp",
		CallbackURL:    os.Getenv("BASE_URL") + "/en/twitter-token",
		Endpoint:       twitter.AuthorizeEndpoint,
	}
	requestToken, requestSecret, _ := config.RequestToken()
	fmt.Println(requestToken, requestSecret)
	db := c.MustGet("DB").(*gorm.DB)
	// url := "https://api.twitter.com/oauth/request_token"
	// method := "POST"
	// //os.Getenv("BASE_URL")+"/en/twitter-token"
	// type TwitterRequest struct {
	// 	OauthCallback string `json:"oauth_callback"`
	// }
	// var request TwitterRequest
	// request.OauthCallback = c.Request.URL.Query()["callback"][0]
	// payloadBytes, _ := json.Marshal(request)
	// payload := bytes.NewReader(payloadBytes)
	// client := &http.Client{}
	// req, err := http.NewRequest(method, url, payload)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("Authorization", "OAuth oauth_consumer_key=\""+oauth_consumer_key+"\",oauth_signature_method=\"HMAC-SHA1\",oauth_timestamp=\""+strconv.FormatInt(time.Now().UTC().Unix(), 10)+"\",oauth_nonce=\"JKvmCxoZmki\",oauth_version=\"1.0\",oauth_signature=\"Bz8%2B%2BUtJEvbM1RE3I%2FyGV7DWgNQ%3D\"")
	// req.Header.Add("Authorization", "OAuth oauth_consumer_key=\"z4CAfGU3ToZlln6v440wRoA4x\",oauth_consumer_secret=\"Oip8cfbwSDykJ2OdfY5iWMjbpxTGK4t1UqHSf5SVCLQuXgUGmp\"")
	// req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Cookie", "personalization_id=\"v1_9figFO/m0qwhU+KQBIs1AQ==\"; guest_id=v1%3A163058097006789858; _twitter_sess=BAh7CSIKZmxhc2hJQzonQWN0aW9uQ29udHJvbGxlcjo6Rmxhc2g6OkZsYXNo%250ASGFzaHsABjoKQHVzZWR7ADoPY3JlYXRlZF9hdGwrCCf4N6Z7AToMY3NyZl9p%250AZCIlNDQ2OTkzODBhYTE1NWFlZmFkYWJhYzdlNWU5MWQ1ZDk6B2lkIiVkMmVk%250AMDk5YTVmM2NhN2U5NWIyZDI0YTEzZjgxM2Y4NQ%253D%253D--07cea9b69cd61caa98a1b52a811970653b2445e1")
	// res, err := client.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer res.Body.Close()

	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(body))
	// respSplit := strings.Split(string(body), "&")
	// oauthToken := strings.Replace(respSplit[0], "oauth_token=", "", 1)
	// oauthTokenSecret := strings.Replace(respSplit[1], "oauth_token_secret=", "", 1)
	db.Exec("INSERT INTO public.twitter_request_token(oauth_token, oauth_token_secret, created_at)VALUES(?, ?, ?);", requestToken, requestSecret, time.Now())
	// // db.Raw("INSERT INTO public.twitter_request_token(oauth_token, oauth_token_secret, created_at)VALUES(?, ?, ?);", oauthToken, oauthTokenSecret, time.Now())
	type Response struct {
		OauthToken       string `json:"oauth_token"`
		OauthTokenSecret string `json:"oauth_token_secret"`
	}
	var resp Response
	resp.OauthToken = requestToken
	resp.OauthTokenSecret = requestSecret
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// GetTwitterAccessToken -  Getting User access token for twitter
// GET /:lang/getAccessToken
// @Summary Get User access token twitter
// @Description Get User access token twitter
// @Tags User
// @Accept  json
// @Produce  json
// @Param lang path string true "lang"
// @Param oauth_token query string false "oauth_token"
// @Param oauth_verifier query string false "oauth_verifier"
// @Success 200 {array} object c.JSON
// @Router /{lang}/getAccessToken [get]
func (hs *HandlerService) GetTwitterAccessToken(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	config := oauth1.Config{
		ConsumerKey:    "z4CAfGU3ToZlln6v440wRoA4x",
		ConsumerSecret: "Oip8cfbwSDykJ2OdfY5iWMjbpxTGK4t1UqHSf5SVCLQuXgUGmp",
		CallbackURL:    os.Getenv("BASE_URL") + "/en/twitter-token",
		Endpoint:       twitter.AuthorizeEndpoint,
	}
	type TwitterSecret struct {
		OauthTokenSecret string
	}
	var twitterSecret TwitterSecret
	db.Table("public.twitter_request_token").Select("oauth_token_secret").Where("oauth_token=?", c.Request.URL.Query()["oauth_token"][0]).Scan(&twitterSecret)
	accessToken, accessSecret, _ := config.AccessToken(c.Request.URL.Query()["oauth_token"][0], twitterSecret.OauthTokenSecret, c.Request.URL.Query()["oauth_verifier"][0])
	// requestToken, requestSecret, _ := config.RequestToken()
	// fmt.Println(requestToken, requestSecret)
	// url := "https://api.twitter.com/oauth/access_token"
	// method := "POST"
	// //os.Getenv("BASE_URL")+"/en/twitter-token"
	// type TwitterRequest struct {
	// 	OauthToken    string `json:"oauth_token"`
	// 	OauthVerifier string `json:"oauth_verifier"`
	// }
	// var request TwitterRequest
	// request.OauthToken = c.Request.URL.Query()["oauth_token"][0]
	// request.OauthVerifier = c.Request.URL.Query()["oauth_verifier"][0]
	// payloadBytes, _ := json.Marshal(request)
	// payload := bytes.NewReader(payloadBytes)
	// client := &http.Client{}
	// req, err := http.NewRequest(method, url, payload)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// req.Header.Add("Authorization", "OAuth oauth_consumer_key=\"z4CAfGU3ToZlln6v440wRoA4x\",oauth_signature_method=\"HMAC-SHA1\",oauth_timestamp=\"1630589958\",oauth_nonce=\"sw9763Rirln\",oauth_version=\"1.0\",oauth_signature=\"fB0fx1esGe7WSaih0uLey6%2B3Ofo%3D\",oauth_token=\""+request.OauthToken+"\"")
	// req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Cookie", "personalization_id=\"v1_9figFO/m0qwhU+KQBIs1AQ==\"; guest_id=v1%3A163058097006789858; _twitter_sess=BAh7CSIKZmxhc2hJQzonQWN0aW9uQ29udHJvbGxlcjo6Rmxhc2g6OkZsYXNo%250ASGFzaHsABjoKQHVzZWR7ADoPY3JlYXRlZF9hdGwrCCf4N6Z7AToMY3NyZl9p%250AZCIlNDQ2OTkzODBhYTE1NWFlZmFkYWJhYzdlNWU5MWQ1ZDk6B2lkIiVkMmVk%250AMDk5YTVmM2NhN2U5NWIyZDI0YTEzZjgxM2Y4NQ%253D%253D--07cea9b69cd61caa98a1b52a811970653b2445e1")
	// res, err := client.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer res.Body.Close()

	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(body))
	// respSplit := strings.Split(string(body), "&")
	// oauthToken := strings.Replace(respSplit[0], "oauth_token=", "", 1)
	// oauthTokenSecret := strings.Replace(respSplit[1], "oauth_token_secret=", "", 1)
	db.Exec("delete from public.twitter_request_token where oauth_token=?;", c.Request.URL.Query()["oauth_token"][0])
	db.Raw("INSERT INTO public.twitter_request_token(oauth_token, oauth_token_secret, created_at)VALUES(?, ?, ?);", accessToken, accessSecret, time.Now())
	// type Response struct {
	// 	OauthToken       string `json:"oauth_token"`
	// 	OauthTokenSecret string `json:"oauth_token_secret"`
	// }
	// var resp Response
	// resp.OauthToken = oauthToken
	// resp.OauthTokenSecret = oauthTokenSecret
	// token := oauth1.NewToken(accessToken, accessSecret)
	c.JSON(http.StatusOK, gin.H{"data": string("oauth_token=" + accessToken + "&oauth_token_secret=" + accessSecret)})
}

// post UpdateUserDetailsByUserid -  UpdateUserDetailsByUserid
// GET /api/users/{id}
// @Summary Show UpdateUserDetailsByUserid
// @Description post UpdateUserDetailsByUserid
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "Id"
// @Param body body UpdateUserDetails true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /api/users/{id} [post]
func (hs *HandlerService) UpdateUserDetailsByUserid(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	fdb := c.MustGet("FDB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var errorresponse = common.ServerErrorResponse()
	var request UpdateUserDetails
	if result := c.ShouldBindJSON(&request); result != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": result.Error(), "status": http.StatusBadRequest})
	}
	type Country struct {
		EnglishName string `json:"english_name"`
	}
	var country Country
	fdb.Table("country").Select("english_name").Where("id=?", request.Country).Find(&country)
	// languageid, _ := strconv.Atoi(request.LanguageId)
	// countryid, _ := strconv.Atoi(request.Country)
	updateuser := User{FirstName: request.FirstName, LastName: request.LastName, LanguageId: request.LanguageId, Country: request.Country, CountryName: country.EnglishName}
	if updateresult := db.Table("user").Where("id=?", c.Param("id")).Update(updateuser).Error; updateresult != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

}

// get UserFilterslist -  UserFilterslist
// GET /api/users/filters
// @Summary Show UserFilterslist
// @Description post UserFilterslist
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} object c.JSON
// @Router /api/users/filters [get]
func (hs *HandlerService) UserFilterslist(c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var response []CountryDetailsRespone
	var totalresponse TotalResponse
	var errorresponse = common.ServerErrorResponse()
	if result := fdb.Table("country").Select("english_name,id").Find(&response).Error; result != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	var deviceresponse []DeviceResponse
	deviceresponse = append(deviceresponse, DeviceResponse{1, "iOS"})
	deviceresponse = append(deviceresponse, DeviceResponse{2, "Android"})
	deviceresponse = append(deviceresponse, DeviceResponse{3, "Apple TV"})
	deviceresponse = append(deviceresponse, DeviceResponse{4, "Smart TV"})
	deviceresponse = append(deviceresponse, DeviceResponse{5, "Roku"})
	deviceresponse = append(deviceresponse, DeviceResponse{6, "Xbox One"})
	deviceresponse = append(deviceresponse, DeviceResponse{7, "PlayStation"})
	deviceresponse = append(deviceresponse, DeviceResponse{9, "Android TV"})
	deviceresponse = append(deviceresponse, DeviceResponse{10, "Amazon Fire TV"})
	// if result := fdb.Table("publish_platform").Select("id,platform").Find(&deviceresponse).Error; result != nil {
	// 	c.JSON(http.StatusInternalServerError, errorresponse)
	// 	return
	// }

	response = append(response, CountryDetailsRespone{"No country", -2})
	deviceresponse = append(deviceresponse, DeviceResponse{-1, "No Active Devices"})
	status1 := Status{0, "Active"}
	status2 := Status{1, "Inactive"}
	var status []Status
	status = append(status, status1)
	status = append(status, status2)
	var sourcetypes []Status
	sourcetypes1 := Status{1, "Reg-Email"}
	sourcetypes2 := Status{2, "SM-Twitter"}
	sourcetypes3 := Status{3, "SM-Facebook"}
	sourcetypes4 := Status{4, "Reg-Mobile"}
	sourcetypes5 := Status{5, "SM-Apple"}

	sourcetypes = append(sourcetypes, sourcetypes1)
	sourcetypes = append(sourcetypes, sourcetypes2)
	sourcetypes = append(sourcetypes, sourcetypes3)
	sourcetypes = append(sourcetypes, sourcetypes4)
	sourcetypes = append(sourcetypes, sourcetypes5)

	var userResponse []UserManagementFilterResponse
	final := UserManagementFilterResponse{"0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z"}
	userResponse = append(userResponse, final)

	var userLeads []Status
	userleads1 := Status{1, "Foxxum"}
	userleads2 := Status{2, "Vidaa"}
	userleads3 := Status{3, "SamsungTV"}
	userleads4 := Status{4, "Apple TV"}
	userleads6 := Status{6, "Roku"}
	userleads7 := Status{7, "AndroidTV"}
	userleads8 := Status{8, "WebOS"}

	userLeads = append(userLeads, userleads1)
	userLeads = append(userLeads, userleads2)
	userLeads = append(userLeads, userleads3)
	userLeads = append(userLeads, userleads4)
	userLeads = append(userLeads, userleads6)
	userLeads = append(userLeads, userleads7)
	userLeads = append(userLeads, userleads8)

	totalresponse.CountryDetailsRespone = response
	totalresponse.DeviceResponse = deviceresponse
	totalresponse.UserstatusResponse = status
	totalresponse.PaycmsConfirmedResponse = status
	totalresponse.PhoneNumberConfirmeds = status
	totalresponse.VerificationStatuses = status
	totalresponse.EmailConfirmeds = status
	totalresponse.NewsLetters = status
	totalresponse.PromotionEnabled = status
	totalresponse.SourceTypes = sourcetypes
	totalresponse.UserManagementFilter = userResponse
	totalresponse.UserLeads = userLeads
	c.JSON(http.StatusOK, totalresponse)
}

// get User view activity by filters -  User view activity by filters
// GET /api/users/:id/viewactivities
// @Summary Show User view activity by filters
// @Description get User view activity by filters
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Param searchText path string false "Search Text"
// @Param contentType path string false "content Type"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Param id path string true "Id"
// @Success 200 {array} object c.JSON
// @Router /api/users/{id}/viewactivities [get]
func (hs *HandlerService) UserViewActivitybyFilters(c *gin.Context) {
	db := c.MustGet("CDB").(*gorm.DB)
	udb := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var errorresponse = common.ServerErrorResponse()
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
	if c.Request.URL.Query()["searchText"] != nil {
		searchText = c.Request.URL.Query()["searchText"][0]
	}
	var contentType string
	if c.Request.URL.Query()["contentType"] != nil {
		contentType = c.Request.URL.Query()["contentType"][0]
	}
	var rawquery string
	var contentDetails, totalCount []ContentDetails
	var where string
	where += " where va.user_id='" + c.Param("id") + "' "
	if contentType == "Movie" {
		where += " and c.content_type::text='" + contentType + "'"
	} else if contentType == "Episode" {
		where += " and c.content_type='Series' "
	}
	if searchText != "" {

		where += " and lower(cpc.transliterated_title) like('%" + strings.ToLower(searchText) + "%')"
	}
	rawquery += "select  distinct va.id, c.content_type,cpc.transliterated_title,va.viewed_at,va.last_watch_position,va.is_hidden,va.device_id ,string_agg(g.english_name,',') as english_name ,pi2.duration  from view_activity va	left join public.content c on va.content_id = c.id   join content_primary_info cpc on c.primary_info_id = cpc.id join playback_item pi2 on pi2.id = va.playback_item_id	join content_genre cg on cg.content_id = c.id   join genre g on g.id=cg.genre_id " + where + "  group by va.id,c.content_type,cpc.transliterated_title,pi2.duration,va.viewed_at ,va.last_watch_position ,va.is_hidden ,va.device_id"

	if res := db.Raw(rawquery).Limit(limit).Offset(offset).Find(&contentDetails).Error; res != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	if res := db.Raw(rawquery).Scan(&totalCount).Error; res != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	finalres := []FinalUserRespone{}
	var finalUserRespone FinalUserRespone
	var newplatform string
	for _, contentrange := range contentDetails {
		finalUserRespone.ViewActivityId = contentrange.Id
		finalUserRespone.ViewedAt = contentrange.ViewedAt
		finalUserRespone.Title = contentrange.TransliteratedTitle
		//new merge based on validation said
		if contentrange.ContentType == "Series" {
			finalUserRespone.ContentType = "Episode"
		} else {
			finalUserRespone.ContentType = contentrange.ContentType
		}
		finalUserRespone.LastWatchPosition = contentrange.LastWatchPosition
		finalUserRespone.DurationSeconds = contentrange.Duration

		var device DeviceName
		udb.Table("device").Select("platform").Where("device_id =(?)", contentrange.DeviceId).Find(&device)

		var count int
		udb.Table("watching_issue").Where("view_activity_id =(?)", contentrange.Id).Count(&count)

		newplatform = common.DeviceNames(device.Platform)
		if newplatform == "web" {
			finalUserRespone.ViewedOnPlatformName = "Website"
		} else {
			finalUserRespone.ViewedOnPlatformName = newplatform
		}
		if count == 0 {
			finalUserRespone.HasWatchingIssues = false
		} else {
			finalUserRespone.HasWatchingIssues = true
		}
		fmt.Println(finalUserRespone.HasWatchingIssues, ";;;;;;;;;;;;;;;;")
		finalUserRespone.IsHidden = contentrange.IsHidden
		var new []string
		new = strings.Split(contentrange.EnglishName, ",")
		finalUserRespone.Genres = new
		finalres = append(finalres, finalUserRespone)
	}

	pages := map[string]int{
		"size":   len(totalCount),
		"offset": int(offset),
		"limit":  int(limit),
	}
	c.JSON(http.StatusOK, gin.H{"pagination": pages, "data": finalres})
}

// get Users List and Search by Filters with Pagination -  Users List and Search by Filters with Pagination
// GET /api/users
// @Summary ShowUsers List and Search by Filters with Pagination
// @Description get Users List and Search by Filters with Pagination
// @Tags User
// @Accept  json
// @Produce  json
// @Param userStatus path string false "userStatus"
// @Param RegistrationSourceType path string false "RegistrationSourceType"
// @Param NewsLetter path string false "NewsLetter"
// @Param PromotionsEnabled path string false "PromotionsEnabled"
// @Param VerificationStatus path string false "VerificationStatus"
// @Param UserLead path string false "UserLead"
// @Param activeDevicePlatform path string false "activeDevicePlatform"
// @Param countryId path string false "countryId"
// @Param searchText path string false "searchText"
// @Param StartDate path string false "StartDate"
// @Param EndDate path string false "EndDate"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Success 200 {array} object c.JSON
// @Router /api/users [get]
func (hs *HandlerService) UsersListandSearchbyFilterswithPagination(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	//	var errorresponse = common.ServerErrorResponse()
	db := c.MustGet("DB").(*gorm.DB)
	//	fdb := c.MustGet("FDB").(*gorm.DB)
	var limits, offsets string
	if c.Request.URL.Query()["limit"] != nil {
		limits = c.Request.URL.Query()["limit"][0]
	}
	if c.Request.URL.Query()["offset"] != nil {
		offsets = c.Request.URL.Query()["offset"][0]
	}
	if limits == "0" {
		limits = "10"
	}
	var userStatus string
	if c.Request.URL.Query()["userStatus"] != nil {
		userStatus = c.Request.URL.Query()["userStatus"][0]
	}
	var RegistrationSourceType string
	if c.Request.URL.Query()["RegistrationSourceType"] != nil {
		RegistrationSourceType = c.Request.URL.Query()["RegistrationSourceType"][0]
	}
	var NewsLetter string
	if c.Request.URL.Query()["NewsLetter"] != nil {
		NewsLetter = c.Request.URL.Query()["NewsLetter"][0]
	}
	var PromotionsEnabled string
	if c.Request.URL.Query()["PromotionsEnabled"] != nil {
		PromotionsEnabled = c.Request.URL.Query()["PromotionsEnabled"][0]
	}
	var VerificationStatus string
	if c.Request.URL.Query()["VerificationStatus"] != nil {
		VerificationStatus = c.Request.URL.Query()["VerificationStatus"][0]
	}
	var UserLead string
	if c.Request.URL.Query()["UserLead"] != nil {
		UserLead = c.Request.URL.Query()["UserLead"][0]
	}
	var activeDevicePlatform string
	if c.Request.URL.Query()["activeDevicePlatform"] != nil {
		activeDevicePlatform = c.Request.URL.Query()["activeDevicePlatform"][0]
	}
	var countryId string
	if c.Request.URL.Query()["countryId"] != nil {
		countryId = c.Request.URL.Query()["countryId"][0]
	}
	var searchText string
	if c.Request.URL.Query()["searchText"] != nil {
		searchText = c.Request.URL.Query()["searchText"][0]
	}
	var StartDate string
	if c.Request.URL.Query()["StartDate"] != nil {
		StartDate = c.Request.URL.Query()["StartDate"][0]
	}
	var EndDate string
	if c.Request.URL.Query()["EndDate"] != nil {
		EndDate = c.Request.URL.Query()["EndDate"][0]
	}

	user := []Users{}
	var filtersFinalResponse FiltersFinalResponse
	var role Role
	db.Table("role").Select("id,name").Where("name='User'").Find(&role)
	var rawquery string
	isdeleted := " is_deleted=false and role_id='" + role.Id + "' "
	rawquery += isdeleted
	if RegistrationSourceType != "" {
		rawquery += " and registration_source='" + RegistrationSourceType + "' "
	}
	if NewsLetter != "" {
		rawquery += "and newsletters_enabled='" + NewsLetter + "' "
	}
	if PromotionsEnabled != "" {
		rawquery += " and promotions_enabled='" + PromotionsEnabled + "' "
	}
	if UserLead != "" {
		rawquery += "and  lower(user_lead)=lower('" + UserLead + "') "
	}
	if searchText != "" {
		rawquery += " and (lower(searchable_text) like  lower('%" + searchText + "%') or lower(first_name) like  lower('%" + searchText + "%') or lower(last_name) like  lower('%" + searchText + "%') or lower(email) like  lower('%" + searchText + "%')) "
	}
	if StartDate != "" && EndDate != "" {
		rawquery += "and registered_at BETWEEN '" + StartDate + " 00:00:00" + "' and '" + EndDate + " 23:59:00 ' "
	}
	if countryId != "" {
		if countryId == "-2" {
			rawquery += "and country IS NULL "
		} else {
			rawquery += "and country=" + countryId + " "
		}
	}
	if userStatus == "0" {
		rawquery += " and EXTRACT(EPOCH FROM (now() - last_activity_at))<1800 "
	} else if userStatus == "1" {
		rawquery += " and (EXTRACT(EPOCH FROM (now() - last_activity_at))>1800 or last_activity_at is null )  "
	}
	if activeDevicePlatform == "-1" {
		rawquery += ""
	} else if activeDevicePlatform != "" && activeDevicePlatform != "-1" {
		rawquery += " and platform=" + activeDevicePlatform + "  "
	}
	if VerificationStatus == "true" {
		rawquery += "and (phone_number_confirmed='" + VerificationStatus + "'  or email_confirmed=' " + VerificationStatus + "') "
	} else if VerificationStatus == "false" {
		rawquery += "and phone_number_confirmed='" + VerificationStatus + "'  and email_confirmed=' " + VerificationStatus + "' "
	}
	var join string
	join += "left join user_device ud on ud.user_id = u.id "
	var rawquery1 string
	var countquery1 string
	if activeDevicePlatform != "-1" {
		rawquery1 = "	select  u.id as id,u.country as country,u.first_name as first_name,u.language_id as language_id,u.last_name as last_name,	u.newsletters_enabled as newsletters_enabled,	u.promotions_enabled as promotions_enabled,u.last_activity_at as last_activity_at,	u.registration_source as registration_source,	u.registered_at as registered_at,	u.email as email,u.email_confirmed as email_confirmed,	u.password_hash as password_hash,	u.user_name as user_name,u.is_deleted as is_deleted,u.phone_number as phone_number,u.phone_number_confirmed as phone_number_confirmed,u.calling_code as calling_code,u.national_number as national_number,	u.country_name as country_name,	u.searchable_text as searchable_text,u.is_adult as is_adult,u.privacy_policy as privacy_policy,u.is_recommend as is_recommend,u.user_lead as user_lead,u.performance as performance,u.google_analytics as google_analytices,	u.firebase as firebase,	u.app_flyer as app_flyer,	u.advertising as advertising,	u.aique as aique,	u.google_ads as google_ads,	u.facebook_ads as facebook_ads,	u.is_gdpr_accepted as is_gdpr_accepted,string_agg(device2.platform::varchar,',') as platform	from (select * from device d left join user_device ud	on ud.device_id = d.device_id where ud.token is not null  ) device2	left join public.user u	on device2.user_id = u.id where " + rawquery + " group by u.id order by registered_at desc"
		countquery1 = "select  u.id as id from (select * from device d left join user_device ud	on ud.device_id = d.device_id where ud.token is not null  ) device2	left join public.user u	on device2.user_id = u.id where " + rawquery + " group by u.id order by registered_at desc "
	} else {
		rawquery1 = "	select  u.id as id,u.country as country,u.first_name as first_name,u.language_id as language_id,u.last_name as last_name,	u.newsletters_enabled as newsletters_enabled,	u.promotions_enabled as promotions_enabled,u.last_activity_at as last_activity_at,	u.registration_source as registration_source,	u.registered_at as registered_at,	u.email as email,u.email_confirmed as email_confirmed,	u.password_hash as password_hash,	u.user_name as user_name,u.is_deleted as is_deleted,u.phone_number as phone_number,u.phone_number_confirmed as phone_number_confirmed,u.calling_code as calling_code,u.national_number as national_number,	u.country_name as country_name,	u.searchable_text as searchable_text,u.is_adult as is_adult,u.privacy_policy as privacy_policy,u.is_recommend as is_recommend,u.user_lead as user_lead,u.performance as performance,u.google_analytics as google_analytices,	u.firebase as firebase,	u.app_flyer as app_flyer,	u.advertising as advertising,	u.aique as aique,	u.google_ads as google_ads,	u.facebook_ads as facebook_ads,	u.is_gdpr_accepted as is_gdpr_accepted	from (select * from device d left join user_device ud	on ud.device_id = d.device_id where ud.token is null  ) device2	left join public.user u	on device2.user_id = u.id where " + rawquery + " group by u.id "
		countquery1 = "select  u.id as id from (select * from device d left join user_device ud	on ud.device_id = d.device_id where ud.token is null  ) device2	left join public.user u	on device2.user_id = u.id where " + rawquery + " group by u.id "
	}
	final := []FiltersFinalResponse{}
	//var totalCount []Users
	type Counting struct {
		ID string `json:"id"`
	}
	var dataforsize []Counting
	//var totalCount int
	var totalCount1 int
	if activeDevicePlatform != "" {
		db.Raw(rawquery1).Limit(limits).Offset(offsets).Find(&user)
		db.Raw(countquery1).Find(&dataforsize)
		for _, data := range user {
			filtersFinalResponse.RegistrationSourceName = common.RegistrationSource(data.RegistrationSource)
			/* diff checking  */
			diffffff := time.Now().Sub(data.LastActivityAt)
			d := int(diffffff.Seconds())
			if userStatus == "" {
				if d < 1800 {
					filtersFinalResponse.Status = 0
					filtersFinalResponse.StatusName = "Active"
				} else {
					filtersFinalResponse.Status = 1
					filtersFinalResponse.StatusName = "Inactive"
				}
			}
			if userStatus == "0" {
				filtersFinalResponse.Status = 0
				filtersFinalResponse.StatusName = "Active"
			} else if userStatus == "1" {
				filtersFinalResponse.Status = 1
				filtersFinalResponse.StatusName = "Inactive"
			}
			filtersFinalResponse.Id = string(data.ID)
			filtersFinalResponse.RegisteredAt = data.RegisteredAt
			filtersFinalResponse.FirstName = data.FirstName
			filtersFinalResponse.LastName = data.LastName
			filtersFinalResponse.Email = data.Email
			filtersFinalResponse.RegistrationSource = data.RegistrationSource
			filtersFinalResponse.NewslettersEnabled = data.NewslettersEnabled
			filtersFinalResponse.PromotionsEnabled = data.PromotionsEnabled
			//	var countryidDetails CountryDetails
			filtersFinalResponse.Country = data.Country
			type Country struct {
				EnglishName string `json:"english_name"`
			}
			//	var country Country
			//	fdb.Table("country").Select("english_name").Where("id=?", data.Country).Find(&country)
			filtersFinalResponse.CountryName = common.CountryName(data.Country)
			filtersFinalResponse.IsAdult = data.IsAdult
			filtersFinalResponse.UserLead = data.UserLead
			filtersFinalResponse.PrivacyPolicy = data.PrivacyPolicy
			filtersFinalResponse.IsRecommend = data.IsRecommend
			filtersFinalResponse.Performance = data.Performance
			filtersFinalResponse.GoogleAnalytics = data.GoogleAnalytics
			filtersFinalResponse.Firebase = data.Firebase
			filtersFinalResponse.AppFlyer = data.AppFlyer
			filtersFinalResponse.Advertising = data.Advertising
			filtersFinalResponse.Aique = data.Aique
			filtersFinalResponse.GoogleAds = data.GoogleAds
			filtersFinalResponse.FacebookAds = data.FacebookAds
			filtersFinalResponse.IsGdprAccepted = data.IsGdprAccepted
			filtersFinalResponse.LanguageId = data.LanguageId
			filtersFinalResponse.PhoneNumber = data.PhoneNumber
			filtersFinalResponse.EmailConfirmed = data.EmailConfirmed
			filtersFinalResponse.PhoneNumberConfirmed = data.PhoneNumberConfirmed
			if data.NewslettersEnabled {
				filtersFinalResponse.NewslettersEnabledDisplayName = "Enabled"
			} else {
				filtersFinalResponse.NewslettersEnabledDisplayName = "Disabled"
			}
			if data.PromotionsEnabled {
				filtersFinalResponse.PromotionsEnabledDisplayName = "Enabled"
			} else {
				filtersFinalResponse.PromotionsEnabledDisplayName = "Disabled"
			}
			if data.LanguageId == 1 {
				filtersFinalResponse.LanguageName = "English"
			} else {
				filtersFinalResponse.LanguageName = "Arabic"
			}
			if data.PhoneNumberConfirmed || data.EmailConfirmed {
				filtersFinalResponse.VerificationStatus = true
			} else {
				filtersFinalResponse.VerificationStatus = false
			}

			if data.PhoneNumberConfirmed || data.EmailConfirmed {
				filtersFinalResponse.VerificationEnabledDisplayName = "Verified"
			} else if !data.PhoneNumberConfirmed && !data.EmailConfirmed {
				filtersFinalResponse.VerificationEnabledDisplayName = "Non-Verified"
			}

			var newstring string
			newstring = ""
			//	var newarr []string
			var anotherarray []string
			if activeDevicePlatform == "-1" {
				filtersFinalResponse.ActiveDevicePlatformNames = newstring
				filtersFinalResponse.NumberOfActiveDevices = 0
			} else if data.Platform != "" {
				newarr := strings.Split(data.Platform, ",")
				for _, token := range newarr {
					anotherarray = append(anotherarray, common.DeviceName(token))
				}
				new := common.DupCount(anotherarray)
				for platform, number := range new {
					if len(newstring) == 0 && number > 1 {
						newstring = newstring + strconv.Itoa(number) + "x" + platform
					} else if len(newstring) == 0 && number == 1 {
						newstring = newstring + platform
					} else if len(newstring) != 0 && number > 1 {
						newstring = newstring + "," + strconv.Itoa(number) + "x" + platform
					} else if len(newstring) != 0 && number == 1 {
						newstring = newstring + "," + platform
					}
				}
				filtersFinalResponse.ActiveDevicePlatformNames = newstring
				filtersFinalResponse.NumberOfActiveDevices = len(anotherarray)
			} else {
				filtersFinalResponse.ActiveDevicePlatformNames = newstring
				filtersFinalResponse.NumberOfActiveDevices = 0
			}
			final = append(final, filtersFinalResponse)
		}
	} else {
		rawquery2 := "select * from (select * from public.user u2 where" + rawquery + "order by u2.registered_at desc limit " + limits + " offset " + offsets + " ) as user_details	left join (select distinct (user_id) from user_device ud where ud.token is not null) active_users2 on user_details.id = active_users2.user_id order by registered_at desc"
		db.Raw(rawquery2).Find(&user)

		db.Raw("select count(*) from (select * from public.user u2 where" + rawquery + " ) as user_details	left join (select distinct (user_id) from user_device ud where ud.token is not null) active_users2 on user_details.id = active_users2.user_id").Count(&totalCount1)

		for _, data := range user {
			type PlatformDetails struct {
				Platform string `json:"platform"`
			}
			var platformdetails PlatformDetails
			filtersFinalResponse.RegistrationSourceName = common.RegistrationSource(data.RegistrationSource)
			rawquery3 := "select string_agg(device.platform::varchar,',') as platform from (select ud.user_id,ud.device_id from user_device ud where ud.token is not null and ud.user_id='" + data.ID + "' ) as active_user2 left join device on active_user2.device_id = device.device_id"
			db.Raw(rawquery3).Find(&platformdetails)
			/* diff checking  */
			diffffff := time.Now().Sub(data.LastActivityAt)

			d := int(diffffff.Seconds())
			if userStatus == "" {
				if d < 1800 {
					filtersFinalResponse.Status = 0
					filtersFinalResponse.StatusName = "Active"
				} else {
					filtersFinalResponse.Status = 1
					filtersFinalResponse.StatusName = "Inactive"
				}
			}
			if userStatus == "0" {
				filtersFinalResponse.Status = 0
				filtersFinalResponse.StatusName = "Active"
			} else if userStatus == "1" {
				filtersFinalResponse.Status = 1
				filtersFinalResponse.StatusName = "Inactive"
			}
			filtersFinalResponse.Id = string(data.ID)
			filtersFinalResponse.RegisteredAt = data.RegisteredAt
			filtersFinalResponse.FirstName = data.FirstName
			filtersFinalResponse.LastName = data.LastName
			filtersFinalResponse.Email = data.Email
			filtersFinalResponse.RegistrationSource = data.RegistrationSource
			filtersFinalResponse.NewslettersEnabled = data.NewslettersEnabled
			filtersFinalResponse.PromotionsEnabled = data.PromotionsEnabled
			//	var countryidDetails CountryDetails
			filtersFinalResponse.Country = data.Country
			type Country struct {
				EnglishName string `json:"english_name"`
			}
			//	var country Country
			//	fdb.Table("country").Select("english_name").Where("id=?", data.Country).Find(&country)
			filtersFinalResponse.CountryName = common.CountryName(data.Country)
			filtersFinalResponse.IsAdult = data.IsAdult
			filtersFinalResponse.UserLead = data.UserLead
			filtersFinalResponse.PrivacyPolicy = data.PrivacyPolicy
			filtersFinalResponse.IsRecommend = data.IsRecommend
			filtersFinalResponse.Performance = data.Performance
			filtersFinalResponse.GoogleAnalytics = data.GoogleAnalytics
			filtersFinalResponse.Firebase = data.Firebase
			filtersFinalResponse.AppFlyer = data.AppFlyer
			filtersFinalResponse.Advertising = data.Advertising
			filtersFinalResponse.Aique = data.Aique
			filtersFinalResponse.GoogleAds = data.GoogleAds
			filtersFinalResponse.FacebookAds = data.FacebookAds
			filtersFinalResponse.IsGdprAccepted = data.IsGdprAccepted
			filtersFinalResponse.LanguageId = data.LanguageId
			filtersFinalResponse.PhoneNumber = data.PhoneNumber
			filtersFinalResponse.EmailConfirmed = data.EmailConfirmed
			filtersFinalResponse.PhoneNumberConfirmed = data.PhoneNumberConfirmed
			if data.NewslettersEnabled {
				filtersFinalResponse.NewslettersEnabledDisplayName = "Enabled"
			} else {
				filtersFinalResponse.NewslettersEnabledDisplayName = "Disabled"
			}
			if data.PromotionsEnabled {
				filtersFinalResponse.PromotionsEnabledDisplayName = "Enabled"
			} else {
				filtersFinalResponse.PromotionsEnabledDisplayName = "Disabled"
			}
			if data.LanguageId == 1 {
				filtersFinalResponse.LanguageName = "English"
			} else {
				filtersFinalResponse.LanguageName = "Arabic"
			}
			if data.PhoneNumberConfirmed || data.EmailConfirmed {
				filtersFinalResponse.VerificationStatus = true
			} else {
				filtersFinalResponse.VerificationStatus = false
			}

			if data.PhoneNumberConfirmed || data.EmailConfirmed {
				filtersFinalResponse.VerificationEnabledDisplayName = "Verified"
			} else if !data.PhoneNumberConfirmed && !data.EmailConfirmed {
				filtersFinalResponse.VerificationEnabledDisplayName = "Non-Verified"
			}

			var newstring string
			newstring = ""

			// var newarr []string
			// newarr=strings.Split(platformdetails.platform,",")
			var anotherarray []string
			fmt.Println(platformdetails.Platform, ";;;;;;;;;;;")
			//		if platformdetails.platform != "" {

			newarr := strings.Split(platformdetails.Platform, ",")
			fmt.Println(len(newarr), newarr, "}}}}}}}}}}}}}}}}}}}}}}}}}}}]]")

			//		db.Table("user_device ud").Select("d.Platform").Joins("join device d on d.device_id=ud.device_id").Where("ud.token IN(?)", newarr).Find(&userdevice)
			if platformdetails.Platform != "" {
				if len(newarr) > 0 {
					for _, token := range newarr {
						anotherarray = append(anotherarray, common.DeviceName(token))

					}
					fmt.Println(newarr, "+++++++++++")
					fmt.Println(anotherarray, "------------------------")
					vb := len(anotherarray)
					fmt.Println(vb, "???????????????")
					new := common.DupCount(anotherarray)
					fmt.Println(len(new), "???????????????")

					for platform, number := range new {

						if len(newstring) == 0 && number > 1 {
							newstring = newstring + strconv.Itoa(number) + "x" + platform
						} else if len(newstring) == 0 && number == 1 {
							newstring = newstring + platform
						} else if len(newstring) != 0 && number > 1 {
							newstring = newstring + "," + strconv.Itoa(number) + "x" + platform
						} else if len(newstring) != 0 && number == 1 {
							newstring = newstring + "," + platform
						}

					}
					filtersFinalResponse.ActiveDevicePlatformNames = newstring
					filtersFinalResponse.NumberOfActiveDevices = len(anotherarray)
				}
			} else {
				filtersFinalResponse.ActiveDevicePlatformNames = newstring
				filtersFinalResponse.NumberOfActiveDevices = 0
			}

			final = append(final, filtersFinalResponse)
		}
	}
	offset1, _ := strconv.Atoi(offsets)
	limit1, _ := strconv.Atoi(limits)
	var totalCount3 int
	if activeDevicePlatform != "" {
		totalCount3 = len(dataforsize)
	} else {
		totalCount3 = totalCount1
	}
	pages := map[string]int{
		"size":   totalCount3,
		"offset": offset1,
		"limit":  limit1,
	}

	c.JSON(http.StatusOK, gin.H{"pagination": pages, "data": final})
}

// DeleteUserDevices -  delete User Logged In Devices
// DELETE /v1/devices/:deviceid
// @Summary Get User Logged In Devices
// @Description Get User Logged In Devices
// @Tags User
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} UserDevicesResponse
// @Router /v1/devices/{deviceid} [delete]
func (hs *HandlerService) DeleteUserDevices(c *gin.Context) {
	if c.Request.Method != http.MethodDelete {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	userId := c.MustGet("userid") //common.USERID
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	//var userdevice UserDevice
	if err := db.Table("user_device").Where("device_id=? and user_id=?", c.Param("deviceid"), userId).Update("token", gorm.Expr("NULL")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
}

// serDetails -  Export user details
// GET /api/users/export
// @Summary Updating Paycms status in user table
// @Description Updating Paycms status in user table
// @Tags User
// @Accept  json
// @Produce  json
// @Param type query string false "File type"
// @Param start_date query string false "Start Date"
// @Param end_date query string false "End Date"
// @Success 200 {array} object c.JSON
// @Router /api/users/export [GET]
func (hs *HandlerService) ExportUserDetails(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var limits, offsets string
	if c.Request.URL.Query()["limit"] != nil {
		limits = c.Request.URL.Query()["limit"][0]
	}
	if c.Request.URL.Query()["offset"] != nil {
		offsets = c.Request.URL.Query()["offset"][0]
	}
	if limits == "0" {
		limits = "10"
	}
	var userStatus string
	if c.Request.URL.Query()["userStatus"] != nil {
		userStatus = c.Request.URL.Query()["userStatus"][0]
	}
	var RegistrationSourceType string
	if c.Request.URL.Query()["RegistrationSourceType"] != nil {
		RegistrationSourceType = c.Request.URL.Query()["RegistrationSourceType"][0]
	}
	var NewsLetter string
	if c.Request.URL.Query()["NewsLetter"] != nil {
		NewsLetter = c.Request.URL.Query()["NewsLetter"][0]
	}
	var PromotionsEnabled string
	if c.Request.URL.Query()["PromotionsEnabled"] != nil {
		PromotionsEnabled = c.Request.URL.Query()["PromotionsEnabled"][0]
	}
	var VerificationStatus string
	if c.Request.URL.Query()["VerificationStatus"] != nil {
		VerificationStatus = c.Request.URL.Query()["VerificationStatus"][0]
	}
	var UserLead string
	if c.Request.URL.Query()["UserLead"] != nil {
		UserLead = c.Request.URL.Query()["UserLead"][0]
	}
	var activeDevicePlatform string
	if c.Request.URL.Query()["activeDevicePlatform"] != nil {
		activeDevicePlatform = c.Request.URL.Query()["activeDevicePlatform"][0]
	}
	var countryId string
	if c.Request.URL.Query()["countryId"] != nil {
		countryId = c.Request.URL.Query()["countryId"][0]
	}
	var searchText string
	if c.Request.URL.Query()["searchText"] != nil {
		searchText = c.Request.URL.Query()["searchText"][0]
	}
	var StartDate string
	if c.Request.URL.Query()["StartDate"] != nil {
		StartDate = c.Request.URL.Query()["StartDate"][0]
	}
	var EndDate string
	if c.Request.URL.Query()["EndDate"] != nil {
		EndDate = c.Request.URL.Query()["EndDate"][0]
	}

	user := []Users{}
	var filtersFinalResponse FiltersFinalResponse
	var role Role
	db.Table("role").Select("id,name").Where("name='User'").Find(&role)
	var rawquery string
	isdeleted := " is_deleted=false and role_id='" + role.Id + "' "
	rawquery += isdeleted
	if RegistrationSourceType != "" {
		rawquery += " and registration_source='" + RegistrationSourceType + "' "
	}
	if NewsLetter != "" {
		rawquery += "and newsletters_enabled='" + NewsLetter + "' "
	}
	if PromotionsEnabled != "" {
		rawquery += " and promotions_enabled='" + PromotionsEnabled + "' "
	}
	if UserLead != "" {
		rawquery += "and  lower(user_lead)=lower('" + UserLead + "') "
	}
	if searchText != "" {
		rawquery += " and (lower(searchable_text) like  lower('%" + searchText + "%') or lower(first_name) like  lower('%" + searchText + "%') or lower(last_name) like  lower('%" + searchText + "%') or lower(email) like  lower('%" + searchText + "%')) "
	}
	if StartDate != "" && EndDate != "" {
		rawquery += "and registered_at BETWEEN '" + StartDate + " 00:00:00" + "' and '" + EndDate + " 23:59:00 ' "
	}
	if countryId != "" {
		if countryId == "-2" {
			rawquery += "and country IS NULL "
		} else {
			rawquery += "and country=" + countryId + " "
		}
	}
	if userStatus == "0" {
		rawquery += " and EXTRACT(EPOCH FROM (now() - last_activity_at))<1800 "
	} else if userStatus == "1" {
		rawquery += " and (EXTRACT(EPOCH FROM (now() - last_activity_at))>1800 or last_activity_at is null )  "
	}
	if activeDevicePlatform == "-1" {
		rawquery += ""
	} else if activeDevicePlatform != "" && activeDevicePlatform != "-1" {
		rawquery += " and platform=" + activeDevicePlatform + "  "
	}
	if VerificationStatus == "true" {
		rawquery += "and (phone_number_confirmed='" + VerificationStatus + "'  or email_confirmed=' " + VerificationStatus + "') "
	} else if VerificationStatus == "false" {
		rawquery += "and phone_number_confirmed='" + VerificationStatus + "'  and email_confirmed=' " + VerificationStatus + "' "
	}
	var join string
	join += "left join user_device ud on ud.user_id = u.id "
	var rawquery1 string
	var countquery1 string
	if activeDevicePlatform != "-1" {
		rawquery1 = "	select  u.id as id,u.country as country,u.first_name as first_name,u.language_id as language_id,u.last_name as last_name,	u.newsletters_enabled as newsletters_enabled,	u.promotions_enabled as promotions_enabled,u.last_activity_at as last_activity_at,	u.registration_source as registration_source,	u.registered_at as registered_at,	u.email as email,u.email_confirmed as email_confirmed,	u.password_hash as password_hash,	u.user_name as user_name,u.is_deleted as is_deleted,u.phone_number as phone_number,u.phone_number_confirmed as phone_number_confirmed,u.calling_code as calling_code,u.national_number as national_number,	u.country_name as country_name,	u.searchable_text as searchable_text,u.is_adult as is_adult,u.privacy_policy as privacy_policy,u.is_recommend as is_recommend,u.user_lead as user_lead,u.performance as performance,u.google_analytics as google_analytices,	u.firebase as firebase,	u.app_flyer as app_flyer,	u.advertising as advertising,	u.aique as aique,	u.google_ads as google_ads,	u.facebook_ads as facebook_ads,	u.is_gdpr_accepted as is_gdpr_accepted,string_agg(device2.platform::varchar,',') as platform	from (select * from device d left join user_device ud	on ud.device_id = d.device_id where ud.token is not null  ) device2	left join public.user u	on device2.user_id = u.id where " + rawquery + " group by u.id order by registered_at desc"
		countquery1 = "select  u.id as id from (select * from device d left join user_device ud	on ud.device_id = d.device_id where ud.token is not null  ) device2	left join public.user u	on device2.user_id = u.id where " + rawquery + " group by u.id order by registered_at desc "
	} else {
		rawquery1 = "	select  u.id as id,u.country as country,u.first_name as first_name,u.language_id as language_id,u.last_name as last_name,	u.newsletters_enabled as newsletters_enabled,	u.promotions_enabled as promotions_enabled,u.last_activity_at as last_activity_at,	u.registration_source as registration_source,	u.registered_at as registered_at,	u.email as email,u.email_confirmed as email_confirmed,	u.password_hash as password_hash,	u.user_name as user_name,u.is_deleted as is_deleted,u.phone_number as phone_number,u.phone_number_confirmed as phone_number_confirmed,u.calling_code as calling_code,u.national_number as national_number,	u.country_name as country_name,	u.searchable_text as searchable_text,u.is_adult as is_adult,u.privacy_policy as privacy_policy,u.is_recommend as is_recommend,u.user_lead as user_lead,u.performance as performance,u.google_analytics as google_analytices,	u.firebase as firebase,	u.app_flyer as app_flyer,	u.advertising as advertising,	u.aique as aique,	u.google_ads as google_ads,	u.facebook_ads as facebook_ads,	u.is_gdpr_accepted as is_gdpr_accepted	from (select * from device d left join user_device ud	on ud.device_id = d.device_id where ud.token is null  ) device2	left join public.user u	on device2.user_id = u.id where " + rawquery + " group by u.id "
		countquery1 = "select  u.id as id from (select * from device d left join user_device ud	on ud.device_id = d.device_id where ud.token is null  ) device2	left join public.user u	on device2.user_id = u.id where " + rawquery + " group by u.id "
	}
	final := []FiltersFinalResponse{}
	//var totalCount []Users
	type Counting struct {
		ID string `json:"id"`
	}
	var dataforsize []Counting
	//var totalCount int
	var totalCount1 int
	if activeDevicePlatform != "" {
		db.Raw(rawquery1).Limit(limits).Offset(offsets).Find(&user)
		db.Raw(countquery1).Find(&dataforsize)
		for _, data := range user {
			filtersFinalResponse.RegistrationSourceName = common.RegistrationSource(data.RegistrationSource)
			/* diff checking  */
			diffffff := time.Now().Sub(data.LastActivityAt)
			d := int(diffffff.Seconds())
			if userStatus == "" {
				if d < 1800 {
					filtersFinalResponse.Status = 0
					filtersFinalResponse.StatusName = "Active"
				} else {
					filtersFinalResponse.Status = 1
					filtersFinalResponse.StatusName = "Inactive"
				}
			}
			if userStatus == "0" {
				filtersFinalResponse.Status = 0
				filtersFinalResponse.StatusName = "Active"
			} else if userStatus == "1" {
				filtersFinalResponse.Status = 1
				filtersFinalResponse.StatusName = "Inactive"
			}
			filtersFinalResponse.Id = string(data.ID)
			filtersFinalResponse.RegisteredAt = data.RegisteredAt
			filtersFinalResponse.FirstName = data.FirstName
			filtersFinalResponse.LastName = data.LastName
			filtersFinalResponse.Email = data.Email
			filtersFinalResponse.RegistrationSource = data.RegistrationSource
			filtersFinalResponse.NewslettersEnabled = data.NewslettersEnabled
			filtersFinalResponse.PromotionsEnabled = data.PromotionsEnabled
			//	var countryidDetails CountryDetails
			filtersFinalResponse.Country = data.Country
			type Country struct {
				EnglishName string `json:"english_name"`
			}
			//	var country Country
			//	fdb.Table("country").Select("english_name").Where("id=?", data.Country).Find(&country)
			filtersFinalResponse.CountryName = common.CountryName(data.Country)
			filtersFinalResponse.IsAdult = data.IsAdult
			filtersFinalResponse.UserLead = data.UserLead
			filtersFinalResponse.PrivacyPolicy = data.PrivacyPolicy
			filtersFinalResponse.IsRecommend = data.IsRecommend
			filtersFinalResponse.Performance = data.Performance
			filtersFinalResponse.GoogleAnalytics = data.GoogleAnalytics
			filtersFinalResponse.Firebase = data.Firebase
			filtersFinalResponse.AppFlyer = data.AppFlyer
			filtersFinalResponse.Advertising = data.Advertising
			filtersFinalResponse.Aique = data.Aique
			filtersFinalResponse.GoogleAds = data.GoogleAds
			filtersFinalResponse.FacebookAds = data.FacebookAds
			filtersFinalResponse.IsGdprAccepted = data.IsGdprAccepted
			filtersFinalResponse.LanguageId = data.LanguageId
			filtersFinalResponse.PhoneNumber = data.PhoneNumber
			filtersFinalResponse.EmailConfirmed = data.EmailConfirmed
			filtersFinalResponse.PhoneNumberConfirmed = data.PhoneNumberConfirmed
			if data.NewslettersEnabled {
				filtersFinalResponse.NewslettersEnabledDisplayName = "Enabled"
			} else {
				filtersFinalResponse.NewslettersEnabledDisplayName = "Disabled"
			}
			if data.PromotionsEnabled {
				filtersFinalResponse.PromotionsEnabledDisplayName = "Enabled"
			} else {
				filtersFinalResponse.PromotionsEnabledDisplayName = "Disabled"
			}
			if data.LanguageId == 1 {
				filtersFinalResponse.LanguageName = "English"
			} else {
				filtersFinalResponse.LanguageName = "Arabic"
			}
			if data.PhoneNumberConfirmed || data.EmailConfirmed {
				filtersFinalResponse.VerificationStatus = true
			} else {
				filtersFinalResponse.VerificationStatus = false
			}

			if data.PhoneNumberConfirmed || data.EmailConfirmed {
				filtersFinalResponse.VerificationEnabledDisplayName = "Verified"
			} else if !data.PhoneNumberConfirmed && !data.EmailConfirmed {
				filtersFinalResponse.VerificationEnabledDisplayName = "Non-Verified"
			}

			var newstring string
			newstring = ""
			//	var newarr []string
			var anotherarray []string
			if activeDevicePlatform == "-1" {
				filtersFinalResponse.ActiveDevicePlatformNames = newstring
				filtersFinalResponse.NumberOfActiveDevices = 0
			} else if data.Platform != "" {
				newarr := strings.Split(data.Platform, ",")
				for _, token := range newarr {
					anotherarray = append(anotherarray, common.DeviceName(token))
				}
				new := common.DupCount(anotherarray)
				for platform, number := range new {
					if len(newstring) == 0 && number > 1 {
						newstring = newstring + strconv.Itoa(number) + "x" + platform
					} else if len(newstring) == 0 && number == 1 {
						newstring = newstring + platform
					} else if len(newstring) != 0 && number > 1 {
						newstring = newstring + "," + strconv.Itoa(number) + "x" + platform
					} else if len(newstring) != 0 && number == 1 {
						newstring = newstring + "," + platform
					}
				}
				filtersFinalResponse.ActiveDevicePlatformNames = newstring
				filtersFinalResponse.NumberOfActiveDevices = len(anotherarray)
			} else {
				filtersFinalResponse.ActiveDevicePlatformNames = newstring
				filtersFinalResponse.NumberOfActiveDevices = 0
			}
			final = append(final, filtersFinalResponse)
		}
	} else {
		rawquery2 := "select * from (select * from public.user u2 where" + rawquery + "order by u2.registered_at desc ) as user_details	left join (select distinct (user_id) from user_device ud where ud.token is not null) active_users2 on user_details.id = active_users2.user_id order by registered_at desc"
		db.Raw(rawquery2).Find(&user)

		db.Raw("select count(*) from (select * from public.user u2 where" + rawquery + " ) as user_details	left join (select distinct (user_id) from user_device ud where ud.token is not null) active_users2 on user_details.id = active_users2.user_id").Count(&totalCount1)

		for _, data := range user {
			type PlatformDetails struct {
				Platform string `json:"platform"`
			}
			var platformdetails PlatformDetails
			filtersFinalResponse.RegistrationSourceName = common.RegistrationSource(data.RegistrationSource)
			rawquery3 := "select string_agg(device.platform::varchar,',') as platform from (select ud.user_id,ud.device_id from user_device ud where ud.token is not null and ud.user_id='" + data.ID + "' ) as active_user2 left join device on active_user2.device_id = device.device_id"
			db.Raw(rawquery3).Find(&platformdetails)
			/* diff checking  */
			diffffff := time.Now().Sub(data.LastActivityAt)

			d := int(diffffff.Seconds())
			if userStatus == "" {
				if d < 1800 {
					filtersFinalResponse.Status = 0
					filtersFinalResponse.StatusName = "Active"
				} else {
					filtersFinalResponse.Status = 1
					filtersFinalResponse.StatusName = "Inactive"
				}
			}
			if userStatus == "0" {
				filtersFinalResponse.Status = 0
				filtersFinalResponse.StatusName = "Active"
			} else if userStatus == "1" {
				filtersFinalResponse.Status = 1
				filtersFinalResponse.StatusName = "Inactive"
			}
			filtersFinalResponse.Id = string(data.ID)
			filtersFinalResponse.RegisteredAt = data.RegisteredAt
			filtersFinalResponse.FirstName = data.FirstName
			filtersFinalResponse.LastName = data.LastName
			filtersFinalResponse.Email = data.Email
			filtersFinalResponse.RegistrationSource = data.RegistrationSource
			filtersFinalResponse.NewslettersEnabled = data.NewslettersEnabled
			filtersFinalResponse.PromotionsEnabled = data.PromotionsEnabled
			//	var countryidDetails CountryDetails
			filtersFinalResponse.Country = data.Country
			type Country struct {
				EnglishName string `json:"english_name"`
			}
			//	var country Country
			//	fdb.Table("country").Select("english_name").Where("id=?", data.Country).Find(&country)
			filtersFinalResponse.CountryName = common.CountryName(data.Country)
			filtersFinalResponse.IsAdult = data.IsAdult
			filtersFinalResponse.UserLead = data.UserLead
			filtersFinalResponse.PrivacyPolicy = data.PrivacyPolicy
			filtersFinalResponse.IsRecommend = data.IsRecommend
			filtersFinalResponse.Performance = data.Performance
			filtersFinalResponse.GoogleAnalytics = data.GoogleAnalytics
			filtersFinalResponse.Firebase = data.Firebase
			filtersFinalResponse.AppFlyer = data.AppFlyer
			filtersFinalResponse.Advertising = data.Advertising
			filtersFinalResponse.Aique = data.Aique
			filtersFinalResponse.GoogleAds = data.GoogleAds
			filtersFinalResponse.FacebookAds = data.FacebookAds
			filtersFinalResponse.IsGdprAccepted = data.IsGdprAccepted
			filtersFinalResponse.LanguageId = data.LanguageId
			filtersFinalResponse.PhoneNumber = data.PhoneNumber
			filtersFinalResponse.EmailConfirmed = data.EmailConfirmed
			filtersFinalResponse.PhoneNumberConfirmed = data.PhoneNumberConfirmed
			if data.NewslettersEnabled {
				filtersFinalResponse.NewslettersEnabledDisplayName = "Enabled"
			} else {
				filtersFinalResponse.NewslettersEnabledDisplayName = "Disabled"
			}
			if data.PromotionsEnabled {
				filtersFinalResponse.PromotionsEnabledDisplayName = "Enabled"
			} else {
				filtersFinalResponse.PromotionsEnabledDisplayName = "Disabled"
			}
			if data.LanguageId == 1 {
				filtersFinalResponse.LanguageName = "English"
			} else {
				filtersFinalResponse.LanguageName = "Arabic"
			}
			if data.PhoneNumberConfirmed || data.EmailConfirmed {
				filtersFinalResponse.VerificationStatus = true
			} else {
				filtersFinalResponse.VerificationStatus = false
			}

			if data.PhoneNumberConfirmed || data.EmailConfirmed {
				filtersFinalResponse.VerificationEnabledDisplayName = "Verified"
			} else if !data.PhoneNumberConfirmed && !data.EmailConfirmed {
				filtersFinalResponse.VerificationEnabledDisplayName = "Non-Verified"
			}

			var newstring string
			newstring = ""

			// var newarr []string
			// newarr=strings.Split(platformdetails.platform,",")
			var anotherarray []string
			fmt.Println(platformdetails.Platform, ";;;;;;;;;;;")
			//		if platformdetails.platform != "" {

			newarr := strings.Split(platformdetails.Platform, ",")
			fmt.Println(len(newarr), newarr, "}}}}}}}}}}}}}}}}}}}}}}}}}}}]]")

			//		db.Table("user_device ud").Select("d.Platform").Joins("join device d on d.device_id=ud.device_id").Where("ud.token IN(?)", newarr).Find(&userdevice)
			if platformdetails.Platform != "" {
				if len(newarr) > 0 {
					for _, token := range newarr {
						anotherarray = append(anotherarray, common.DeviceName(token))

					}
					fmt.Println(newarr, "+++++++++++")
					fmt.Println(anotherarray, "------------------------")
					vb := len(anotherarray)
					fmt.Println(vb, "???????????????")
					new := common.DupCount(anotherarray)
					fmt.Println(len(new), "???????????????")

					for platform, number := range new {

						if len(newstring) == 0 && number > 1 {
							newstring = newstring + strconv.Itoa(number) + "x" + platform
						} else if len(newstring) == 0 && number == 1 {
							newstring = newstring + platform
						} else if len(newstring) != 0 && number > 1 {
							newstring = newstring + "," + strconv.Itoa(number) + "x" + platform
						} else if len(newstring) != 0 && number == 1 {
							newstring = newstring + "," + platform
						}

					}
					filtersFinalResponse.ActiveDevicePlatformNames = newstring
					filtersFinalResponse.NumberOfActiveDevices = len(anotherarray)
				}
			} else {
				filtersFinalResponse.ActiveDevicePlatformNames = newstring
				filtersFinalResponse.NumberOfActiveDevices = 0
			}

			final = append(final, filtersFinalResponse)
		}
	}
	fileType := c.Request.URL.Query()["type"][0]
	userId := c.MustGet("userid") //common.USERID
	var admin AdminDetails
	db.Raw(`SELECT email FROM "user" WHERE id=?`, userId).Scan(&admin)
	f := excelize.NewFile()
	var headerStyle int
	var er error
	// Create a new sheet.
	index := f.NewSheet("Sheet1")
	if headerStyle, er = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "FFFFFF", Bold: false, Size: 14, Family: "Arial"},
		Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"FFA500"}, Pattern: 1},
	}); er != nil {
		fmt.Println(er, "error er")
		return
	}
	f.AutoFilter("Sheet1", "A1", "P16", "")
	f.SetColWidth("Sheet1", "A", "P", 20)
	f.SetCellStyle("Sheet1", "A1", "P1", headerStyle)
	// Set value of a cell.
	// f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "A1", "First Name")
	f.SetCellValue("Sheet1", "B1", "Last Name")
	f.SetCellValue("Sheet1", "C1", "Status")
	f.SetCellValue("Sheet1", "D1", "Country")
	f.SetCellValue("Sheet1", "E1", "Registration Date")
	f.SetCellValue("Sheet1", "F1", "Email Address")
	f.SetCellValue("Sheet1", "G1", "PhoneNumber")
	f.SetCellValue("Sheet1", "H1", "Tailored Genres")
	f.SetCellValue("Sheet1", "I1", "Active Devices")
	f.SetCellValue("Sheet1", "J1", "Number of Active Devices")
	f.SetCellValue("Sheet1", "K1", "Language")
	f.SetCellValue("Sheet1", "L1", "Newsletter")
	f.SetCellValue("Sheet1", "M1", "Promotions")
	f.SetCellValue("Sheet1", "N1", "Source")
	f.SetCellValue("Sheet1", "O1", "Lead Device")
	f.SetCellValue("Sheet1", "P1", "Verification Status")

	fmt.Println("data from db.........", user)

	for i := 0; i < len(final); i++ {
		// fmt.Println("A" + strconv.Itoa(i+2))
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), final[i].FirstName)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), final[i].LastName)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), final[i].StatusName)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), final[i].CountryName)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), final[i].RegisteredAt)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), final[i].Email)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i+2), final[i].PhoneNumber)
		f.SetCellValue("Sheet1", "H"+strconv.Itoa(i+2), final[i].TailoredGenres)
		f.SetCellValue("Sheet1", "I"+strconv.Itoa(i+2), final[i].ActiveDevicePlatformNames)
		f.SetCellValue("Sheet1", "J"+strconv.Itoa(i+2), final[i].NumberOfActiveDevices)
		f.SetCellValue("Sheet1", "K"+strconv.Itoa(i+2), final[i].LanguageName)
		f.SetCellValue("Sheet1", "L"+strconv.Itoa(i+2), final[i].NewslettersEnabledDisplayName)
		f.SetCellValue("Sheet1", "M"+strconv.Itoa(i+2), final[i].PromotionsEnabledDisplayName)
		f.SetCellValue("Sheet1", "N"+strconv.Itoa(i+2), final[i].RegistrationSourceName)
		f.SetCellValue("Sheet1", "O"+strconv.Itoa(i+2), final[i].UserLead)
		f.SetCellValue("Sheet1", "P"+strconv.Itoa(i+2), final[i].VerificationEnabledDisplayName)
	}
	var filename string
	filename = "user_details.xlsx"
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	myfile, e := os.Create(filename)
	if e != nil {
		fmt.Println(e)
	}
	myfile.Close()
	p, _ := os.Getwd()
	fmt.Println(p, "path")
	// Save spreadsheet by the given path.
	if err := f.SaveAs(filename); err != nil {
		fmt.Println(err)
	}
	if fileType == "csv" || fileType == "xlsx" {
		filename_out := "user_details.csv"
		csvOpts := func(cw *csv.Writer) { cw.Comma = ([]rune(","))[0] }
		// Open our XSLX file
		file, err := xlsx.New(filename)
		if err != nil {
			fmt.Println("error in opening file", err)
		}
		mydir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		outFile, err := xlsx.GetOutFile(filename_out, mydir)
		//fmt.Println(outFile.Name, "filenaaaaaaaaaa")
		if err != nil {
			fmt.Println("error in getting path of the output file", err)
		}
		err = file.GenerateCSVFromSheet(outFile, 0, csvOpts)
		if err != nil {
			fmt.Println("error in generating csv from sheet", err)
		}
		err = outFile.Close()
		if err != nil {
			fmt.Println("error in closing file", err)
		}
		// data, err := ioutil.ReadFile("user_details.csv")
		// if err != nil {
		// 	fmt.Println("ERROR IN READING THE FILE", err)
		// }
		// input, errr := buildEmailInput("support@weyyak.com", admin.Email, "Users Report", "Users Report", data)
		// fmt.Println(errr, "error from sending email through ses")
		msg := gomail.NewMessage()
		msg.SetHeader("From", "support@weyyak.com")
		msg.SetHeader("To", admin.Email)
		msg.SetHeader("Subject", "Users Report")
		//	msg.SetBody("Message", "Please donot reply")
		msg.Attach("user_details.xlsx")
		var emailRaw bytes.Buffer
		msg.WriteTo(&emailRaw)
		message := ses.RawMessage{Data: emailRaw.Bytes()}
		//session, err := session.NewSession()
		//svc := ses.New(session, &aws.Config{Region: aws.String(os.Getenv("SES_REGION"))})
		source := aws.String("support@weyyak.com")
		destinations := []*string{aws.String(admin.Email)}
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(os.Getenv("SES_REGION")),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("SES_ID"),     // id
				os.Getenv("SES_SECRET"), // secret
				""),                     // token can be left blank for now
		})
		fmt.Println(sess)
		svc := ses.New(sess)
		//  message := ses.RawMessage{Data: []byte(` From: xxx <xxx@xxx.com>\\nTo: xxx  <xxx@xxx.com>\\nSubject: Test email (contains an attachment)\\nMIME-Version: 1.0\\nContent-type: Multipart/Mixed; boundary=\"NextPart\"\\n\\n--NextPart\\nContent-Type: text/plain\\n\\nThis is the message body.\\n\\n--NextPart\\nContent-Type: text/plain;\\nContent-Disposition: attachment; filename=\"sample.txt\"\\n\\nThis is the text in the attachment.\\n\\n--NextPart--" `)}
		input := ses.SendRawEmailInput{Source: source, Destinations: destinations, RawMessage: &message}
		output, err := svc.SendRawEmail(&input)
		if err != nil {
			fmt.Println("ERROR WHILE SENDING A MAIL")
		}
		fmt.Println(output)
		// message := ses.RawMessage{Data: emailRaw.Bytes()}
		// sess, _ := session.NewSession(&aws.Config{
		// 	Region: aws.String(os.Getenv("SES_REGION")),
		// 	Credentials: credentials.NewStaticCredentials(
		// 		os.Getenv("SES_ID"),     // id
		// 		os.Getenv("SES_SECRET"), // secret
		// 		""),                     // token can be left blank for now
		// })
		// fmt.Println(sess)
		// svc := ses.New(sess)
		// result, err := svc.SendRawEmail(message)
		// if err != nil {
		// 	fmt.Println("error in sending raw email", err)
		// }
		// fmt.Println(result)
	}
	// sender := common.New()
	// m := common.NewMessage("Users Report", "Users Report *This is an automatically generated email, please do not reply*")
	// m.To = []string{admin.Email}
	// // m.CC = []string{"copy1@gmail.com", "copy2@gmail.com"}
	// // m.BCC = []string{"bc@gmail.com"}
	// m.AttachFile(filename)
	// fmt.Println(sender.Send(m))
	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "User Report Sent to User Email Succesfully."})
}

// GetUserRatingsDetailsWithSearchText -Get User Ratings details with search text
// GET /users/{id}/ratings
// @Summary Get User Ratings details with search text
// @Description Get User Ratings details with search text
// @Tags User
// @Accept  json
// @Security Authorization
// @Param searchText path string false "SearchText"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Success 200 {array} object c.JSON
// @Router /users/{id}/ratings [get]
func (hs *HandlerService) GetUserRatingsDetailsWithSearchText(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var errorresponse = common.ServerErrorResponse()
	db := c.MustGet("CDB").(*gorm.DB)
	userdb := c.MustGet("DB").(*gorm.DB)

	var limit, offset int64

	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if limit == 0 {
		limit = 10
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	paramUserId := strings.ToLower(c.Param("id"))

	var searchText string
	var totalcount int
	var ratingrecords []RatingRecords
	var ratingByUser RatingByUser
	ratingByUsers := []RatingByUser{}

	rawquery := "select c.id, rc.rated_at,cpi.transliterated_title as title, c.content_type,rc.device_id,rc.rating, rc.is_hidden from content c join content_primary_info cpi on cpi.id=c.primary_info_id join rated_content rc on rc.content_id =c.id where rc.user_id ='" + paramUserId + "' "
	if c.Request.URL.Query()["searchText"] != nil {
		searchText = strings.TrimSpace(strings.ToLower(c.Request.URL.Query()["searchText"][0]))
		rawquery += " and lower(cpi.transliterated_title) like '%" + searchText + "%' "
	}
	if c.Request.URL.Query()["contentType"] != nil {
		var contentType1 string
		contentType1 = strings.TrimSpace(c.Request.URL.Query()["contentType"][0])
		if contentType1 == "Livetv" {
			contentType1 = "LiveTV"
		} else if contentType1 == "Episode" {
			contentType1 = "Series"
		}
		fmt.Println(contentType1, ";;;;;;;;;;;")
		rawquery += " and c.content_type = '" + contentType1 + "' "
	}

	if err := db.Raw(rawquery).Scan(&ratingrecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	totalcount = len(ratingrecords)

	if err := db.Raw(rawquery).Limit(limit).Offset(offset).Find(&ratingrecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	for _, cKey := range ratingrecords {
		ratingByUser.Title = cKey.Title
		ratingByUser.ContentType = cKey.ContentType
		ratingByUser.RatedAt = cKey.RatedAt

		var genname []Genname
		querygenre := "select g.english_name from content c left join content_genre cg on c.id = cg.content_id left join genre g on cg.genre_id = g.id where c.id='" + cKey.Id + "' "
		if err := db.Raw(querygenre).Find(&genname).Error; err != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		var gstring []string
		for _, idarr := range genname {
			gstring = append(gstring, idarr.EnglishName)
		}
		ratingByUser.Genres = gstring

		var platformName PlatformName
		queryplatform := "select d.platform as rated_on_platform_name from device d where  d.device_id ='" + cKey.DeviceId + "' "
		userdb.Raw(queryplatform).Limit(1).Find(&platformName)
		// 	c.JSON(http.StatusInternalServerError, errorresponse)
		// 	return
		// }
		ratingByUser.RatedOnPlatformName = common.DeviceNames(platformName.RatedOnPlatformName)
		ratingByUser.Rating = cKey.Rating
		ratingByUser.IsHidden = cKey.IsHidden
		ratingByUsers = append(ratingByUsers, ratingByUser)
	}
	pagination := map[string]int{
		"size":   totalcount,
		"offset": int(offset),
		"limit":  int(limit),
	}
	c.JSON(http.StatusOK, gin.H{"pagination": pagination, "data": ratingByUsers})
	return
}

func (hs *HandlerService) GetUserWAtchingIssues(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var errorresponse = common.ServerErrorResponse()
	userdb := c.MustGet("DB").(*gorm.DB)
	var watchingissues []WatchingIssue
	if err := userdb.Table("watching_issue").Select("reported_at,is_video,is_sound,is_translation,is_communication,description").Where("view_activity_id=?", c.Param("id")).Find(&watchingissues).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": watchingissues})
}

// func buildEmailInput(source, destination, subject, message string,
// 	csvFile []byte) (*ses.SendRawEmailInput, error) {

// 	buf := new(bytes.Buffer)
// 	writer := multipart.NewWriter(buf)

// 	// email main header:
// 	h := make(textproto.MIMEHeader)
// 	h.Set("From", source)
// 	h.Set("To", destination)
// 	h.Set("Return-Path", source)
// 	h.Set("Subject", subject)
// 	h.Set("Content-Language", "en-US")
// 	h.Set("Content-Type", "multipart/mixed; boundary=\""+writer.Boundary()+"\"")
// 	h.Set("MIME-Version", "1.0")
// 	//	records := readCsvFile("/home/vivekk/go-qa/weyyak-ms-go-api/user/user_details.csv")
// 	//fmt.Println(records, "jjjjjjjjjjj")
// 	_, err := writer.CreatePart(h)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// body:
// 	h = make(textproto.MIMEHeader)
// 	h.Set("Content-Transfer-Encoding", "7bit")
// 	h.Set("Content-Type", "text/plain; charset=us-ascii")
// 	part, err := writer.CreatePart(h)
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, err = part.Write([]byte(message))
// 	if err != nil {
// 		return nil, err
// 	}

// 	// file attachment:
// 	fn := "user_details.csv"
// 	h = make(textproto.MIMEHeader)
// 	h.Set("Content-Disposition", "attachment; filename="+fn)
// 	h.Set("Content-Type", "text/csv; x-unix-mode=0644; name=\""+fn+"\"")
// 	h.Set("Content-Transfer-Encoding", "7bit")

// 	part, err = writer.CreatePart(h)
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, err = part.Write(csvFile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = writer.Close()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Strip boundary line before header (doesn't work with it present)
// 	s := buf.String()
// 	if strings.Count(s, "\n") < 2 {
// 		return nil, fmt.Errorf("invalid e-mail content")
// 	}
// 	s = strings.SplitN(s, "\n", 2)[1]

// 	raw := ses.RawMessage{
// 		Data: []byte(s),
// 	}
// 	input := &ses.SendRawEmailInput{
// 		Destinations: []*string{aws.String(destination)},
// 		Source:       aws.String(source),
// 		RawMessage:   &raw,
// 	}

// 	return input, nil
// }

// func readCsvFile(filePath string) [][]string {
// 	f, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println("Unable to read input file "+filePath, err)
// 	}
// 	defer f.Close()

// 	csvReader := csv.NewReader(f)
// 	records, err := csvReader.ReadAll()
// 	if err != nil {
// 		fmt.Println("Unable to parse file as CSV for "+filePath, err)
// 	}

// 	return records
// }
