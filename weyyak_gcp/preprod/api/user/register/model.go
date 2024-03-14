package register

import (
	"time"

	_ "github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User - struct for DB binding
type User struct {
	Id                   string    `json:"id" gorm:"primary_key"`
	Country              int       `json:"countryId"`
	FirstName            string    `json:"firstName"`
	IsBackOfficeUser     bool      `json:"isBackOfficeUser"`
	LanguageId           int       `json:"languageId"`
	LastName             string    `json:"lastName"`
	NewslettersEnabled   bool      `json:"newslettersEnabled"`
	PromotionsEnabled    bool      `json:"promotionsEnabled"`
	LastActivityAt       time.Time `json:"lastActivityAt"`
	RegistrationSource   int       `json:"registrationSource"`
	RegisteredAt         time.Time `json:"registeredAt"`
	Email                string    `json:"email" gorm:"unique"`
	EmailConfirmed       bool      `json:"emailConfirmed"`
	PasswordHash         string    `json:"passwordHash" binding:"required"`
	SecurityStamp        string    `json:"securityStamp"`
	UserName             string    `json:"userName"`
	IsDeleted            bool      `json:"isDeleted"`
	PhoneNumber          string    `json:"phoneNumber" gorm:"unique"`
	PhoneNumberConfirmed bool      `json:"phoneNumberConfirmed"`
	CallingCode          string    `json:"callingCode"`
	NationalNumber       string    `json:"nationalNumber"`
	CountryName          string    `json:"countryName"`
	SearchableText       string    `json:"searchableText"`
	Paycmsstatus         bool      `json:"payCMSStatus"`
	IsAdult              bool      `json:"isAdult"`
	PrivacyPolicy        bool      `json:"privacyPolicy"`
	IsRecommend          bool      `json:"isRecommend"`
	UserLead             string    `json:"userLead"`
	Performance          bool      `json:"performance"`
	GoogleAnalytics      bool      `json:"googleAnalytics"`
	Firebase             bool      `json:"firebase"`
	AppFlyer             bool      `json:"appFlyer"`
	Advertising          bool      `json:"advertising"`
	Aique                bool      `json:"aique"`
	GoogleAds            bool      `json:"googleAds"`
	FacebookAds          bool      `json:"facebookAds"`
	IsGdprAccepted       bool      `json:"isGdprAccepted"`
	SaltStored           string    `json:"-"`
	Version              int       `json:"-"`
	RoleId               string    `json:"-"`
}

// RequestRegisterUserUsingEmail - struct for DB binding
type RequestRegisterUserUsingEmail struct {
	Password      string `json:"password"`
	LanguageId    int    `json:"languageId"`
	PrivacyPolicy bool   `json:"PrivacyPolicy"`
	IsAdult       bool   `json:"isAdult"`
	IsRecommend   bool   `json:"IsRecommend"`
	Email         string `json:"email"`
	CountryName   string `json:"countryName"`
	Alpha2code    string `json:"Alpha2code"`
	Source        string `json:"source"` /* sync usecase */
	UserId        string `json:"userId"` /* sync usecase */

}

// RequestRegisterUserUsingSMS - struct for DB binding
type RequestRegisterUserUsingSMS struct {
	Password           string `json:"password"`
	LanguageId         int    `json:"languageId"`
	PrivacyPolicy      bool   `json:"PrivacyPolicy"`
	IsAdult            bool   `json:"isAdult"`
	IsRecommend        bool   `json:"IsRecommend"`
	PhoneNumber        string `json:"phonenumber"`
	Silentregistration bool   `json:"silentregistration"`
	Source             string `json:"source"` /* sync usecase */
	UserId             string `json:"userId"` /* sync usecase */

}

type SilentUser struct {
	Id          string `json:"id"`
	UserName    string `json:"userName"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

// HashResponse - struct for mapping hashresponse
type HashResponse struct {
	Password     string `json:"password"`
	PasswordHash string `json:"passwordHash"`
}

// User - struct for DB binding
type UpdateUser struct {
	FirstName            *string `json:"firstName"`
	LastName             *string `json:"lastName"`
	Email                string  `json:"email" gorm:"unique"`
	NewslettersEnabled   *bool   `json:"newslettersEnabled"`
	PromotionsEnabled    *bool   `json:"promotionsEnabled"`
	Country              int     `json:"countryId,omitempty"`
	CountryName          string  `json:"countryName"`
	IsAdult              bool    `json:"isAdult"`
	UserLead             string  `json:"userLead"`
	PrivacyPolicy        bool    `json:"privacyPolicy"`
	IsRecommend          bool    `json:"isRecommend"`
	Performance          bool    `json:"performance"`
	GoogleAnalytics      bool    `json:"googleAnalytics"`
	Firebase             bool    `json:"firebase"`
	AppFlyer             bool    `json:"appFlyer"`
	Advertising          bool    `json:"advertising"`
	Aique                bool    `json:"aique"`
	GoogleAds            bool    `json:"googleAds"`
	FacebookAds          bool    `json:"facebookAds"`
	IsGdprAccepted       bool    `json:"isGdprAccepted"`
	CleverTap            bool    `json:"cleverTap"`
	LanguageId           int     `json:"languageId"`
	RegistrationSource   int     `json:"registrationSource"`
	PhoneNumber          string  `json:"phoneNumber" gorm:"unique"`
	EmailConfirmed       bool    `json:"emailConfirmed"`
	PhoneNumberConfirmed bool    `json:"phoneNumberConfirmed"`
}

//UserDevicesResponse - struct for mapping userDevicesResponse

type UserDevicesResponse struct {
	Name     string `json:"name"`
	Platform int    `json:"platform"`
	Id       string `json:"id" gorm:"primary_key"`
}

type userProfileResponse struct {
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	Email                string `json:"email"`
	NewslettersEnabled   bool   `json:"newslettersEnabled"`
	PromotionsEnabled    bool   `json:"promotionsEnabled"`
	Country              int    `json:"countryId"`
	CountryName          string `json:"countryName"`
	IsAdult              bool   `json:"isAdult"`
	UserLead             string `json:"userLead"`
	PrivacyPolicy        bool   `json:"privacyPolicy"`
	IsRecommend          bool   `json:"isRecommend"`
	Performance          bool   `json:"performance"`
	GoogleAnalytics      bool   `json:"googleAnalytics"`
	Firebase             bool   `json:"firebase"`
	AppFlyer             bool   `json:"appFlyer"`
	Advertising          bool   `json:"advertising"`
	Aique                bool   `json:"aique"`
	GoogleAds            bool   `json:"googleAds"`
	FacebookAds          bool   `json:"facebookAds"`
	IsGdprAccepted       bool   `json:"isGdprAccepted"`
	CleverTap            bool   `json:"cleverTap"`
	LanguageId           int    `json:"languageId"`
	LanguageName         string `json:"languageName"`
	RegistrationSource   int    `json:"registrationSource"`
	PhoneNumber          string `json:"phoneNumber"`
	EmailConfirmed       bool   `json:"emailConfirmed"`
	PhoneNumberConfirmed bool   `json:"phoneNumberConfirmed"`
	VerificationStatus   bool   `json:"verificationStatus"`
}

type LanguageDetails struct {
	EnglishName string
	ArabicName  string
}

type ApplicationSetting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type UserDevicesLimitCount struct {
	UserDevicesLimit int `json:"userDevicesLimit"`
}

type LoginRequest struct {
	GrantType      string `json:"grant_type" validate:"required"`
	Username       string `json:"username" validate:"required"`
	Password       string `json:"Password" validate:"required"`
	DeviceId       string `json:"DeviceId" validate:"required"`
	DevicePlatform string `json:"DevicePlatform" validate:"required"`
	DeviceName     string `json:"DeviceName" validate:"required"`
}

type ConfirmEmail struct {
	ConfirmationToken string `json:"confirmationToken"`
	DateTimeToken     string `json:"dateTimeToken"`
	Source            string `json:"source"` /* sync usecase */
	UserId            string `json:"userId"` /* sync usecase */
}

type PaycmsStatus struct {
	Userid string `json:"userid"`
}

type PhoneNumber struct {
	PhoneNumber          string `json:"phoneNumber"`
	PhoneNumberConfirmed bool   `json:"phone_number_confirmed"`
}

type Emailcheck struct {
	Email string `json:"email"`
}

// Reset password with email
type ResendEmail struct {
	Email      string `json:"email"`
	Alpha2code string `json:"Alpha2code"`
}
type ValidateEmail struct {
	Email      string `json:"email"`
	LanguageId int    `json:"languageId"`
	ID         string `json:"id"`
}
type CollectEmailOtp struct {
	Message string    `json:"message"`
	SentOn  time.Time `json:"sentOn"`
}

// ERROR RESPONSE
type ErrorCode struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type PasswordError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ConfirmationToken struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type DateTimeToken struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type phoneNumberError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type RequestType struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type DevicePlatform struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ResetPasswordToken struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type PhoneNumberError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type OtpValidator struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ValidatePassword struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type NewPasswordValidate struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type LanguageId struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

type Invalid struct {
	Email               *ErrorCode           `json:"email,omitempty"`
	ConfirmationToken   *ConfirmationToken   `json:"confirmationToken,omitempty"`
	DateTimeToken       *DateTimeToken       `json:"dateTimeToken,omitempty"`
	Password            *PasswordError       `json:"password,omitempty"`
	PhoneNumber         *phoneNumberError    `json:"phoneNumber,omitempty"`
	RequestType         *RequestType         `json:"requestType,omitempty"`
	DevicePlatform      *DevicePlatform      `json:"devicePlatform,omitempty"`
	ResetPasswordToken  *ResetPasswordToken  `json:"resetPasswordToken,omitempty"`
	PhoneNumberError    *PhoneNumberError    `json:"phoneNumberError,omitempty"`
	OtpValidator        *OtpValidator        `json:"otp,omitempty"`
	ValidatePassword    *ValidatePassword    `json:"oldpassword,omitempty"`
	NewPasswordValidate *NewPasswordValidate `json:"Password,omitempty"`
	LanguageId          *LanguageId          `json:"languageId,omitempty"`
}
type FinalErrorResponse struct {
	Error       string  `json:"error"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	RequestId   string  `json:"requestId"`
	Invalid     Invalid `json:"invalid,omitempty"`
}
type FinalResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
	Code        string `json:"code"`
	RequestId   string `json:"requestId"`
}

// End of ERROR RESPONSE

// PlaylistedContent -- struct for db binding
type PlaylistedContent struct {
	UserId    string    `json:"user_id"`
	ContentId string    `json:"content_id"`
	AddedAt   time.Time `json:"added_at"`
}

// PlaylistedContentRequest -- Request for playlisted content
type PlaylistedContentRequest struct {
	Id          int           `json:"id"`
	Title       string        `json:"title"`
	ContentType string        `json:"contentType"`
	Genres      []interface{} `json:"genres"`
}

// Content -- struct for db binding
type Content struct {
	Id                     string    `json:"id"`
	AverageRating          float64   `json:"average_rating"`
	AverageRatingUpdatedAt time.Time `json:"average_rating_updated_at"`
	ContentKey             int       `json:"content_key"`
	ContentType            string    `json:"content_type"`
	PrimaryInfoId          string    `json:"primary_info_id"`
	AboutTheContentInfoId  string    `json:"about_the_content_info_id"`
	CastId                 string    `json:"cast_id"`
	MusicId                string    `json:"music_id"`
	TagInfoId              string    `json:"tag_info_id"`
}
type OtpRecord struct {
	Phone   string `json:"Phoneumber"`
	Message string `json:"message"`
	SentOn  time.Time
	Number  int `json:"number"` // for count
}

type UserDevice struct {
	Token    string `json:"token"`
	UserId   string `json:"user_id"`
	DeviceId string `json:"device_id"`
}

type PairingCode struct {
	DeviceId         string    `json:"device_id"`
	DeviceCode       string    `json:"device_code"`
	UserCode         string    `json:"user_code"`
	CreatedAt        time.Time `json:"created_at"`
	ExpiresAt        time.Time `json:"expires_at"`
	UserId           string    `json:"user_id"`
	SubscriptionDate time.Time `json:"subscription_date"`
}

type DeviceIds struct {
	DeviceId string `json:"device_id"`
}

type Number struct {
	Phone       string `json:"Phonenumber"`
	RequestType string `json:"requestType"`
}
type UserDetails struct {
	Phone       string `json:"Phoneumber"`
	Message     string `json:"message"`
	SentOn      time.Time
	Number      int    `json:"number"` // for count
	CallingCode string `json:"callingCode"`
}

type Message struct {
	Message string `json:"message"`
	Phone   string `json:"phone"`
}

type PhnConfirmed struct {
	PhoneNumberConfirmed bool `json:"phone_number_confirmed"`
}
type CountryCodeDetails struct {
	CallingCode string `json:"calling_code"`
}

type Verify struct {
	Message     string `json:"otp"`
	PhoneNumber string `json:"phoneNumber"`
	Source      string `json:"source"` /* sync usecase */
}
type VerifyDetails struct {
	Message     string `json:"otp"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
	Source      string `json:"source"` /* sync usecase */

}

type Valid struct {
	Message     string `json:"otp"`
	PhoneNumber string `json:"phoneNumber"`
	Source      string `json:"source"` /* sync usecase */

}
type Details struct {
	Message string `json:"message"`
	SentOn  time.Time
}

type Password struct {
	PasswordHash string `json:"oldPassword"`
	NewPassword  string `json:"Password"`
}
type PasswordDetails struct {
	PasswordHash string `json:"password"`
	Version      int    `json:"version"`
	SaltStored   string `json:"saltstored"`
}
type UpdatePassword struct {
	PasswordHash string `json:"passwordhash"`
	SaltStored   string `json:"saltstored"`
	Version      int    `json:"version"`
}
type ChangePhoneNumber struct {
	PhoneNumber          string `json:"phonenumber"`
	CallingCode          string `json:"callingcode"`
	NationalNumber       string `json:"nationalnumber"`
	PhoneNumberConfirmed bool   `json:"phone_number_confirmed"`
}

type UserLangDetails struct {
	Id          string `json:"id"`
	LanguageId  int    `json:"lanugageid"`
	PhoneNumber string `json:"phoneNumber"`
}

type RequestPairingCode struct {
	DeviceId       string `json:"deviceId"`
	DeviceName     string `json:"deviceName"`
	DevicePlatform string `json:"devicePlatform"`
}

type Device struct {
	DeviceId  string    `json:"deviceId"`
	Name      string    `json:"deviceName"`
	Platform  string    `json:"devicePlatform"`
	CreatedAt time.Time `json:"created_at"`
}

type Platform struct {
	PlatformId int    `json:"platform_id"`
	Name       string `json:"name"`
}

type ResponsePairingCode struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type VerifyPairingCode struct {
	UserCode UsercodeRequest `json:"user_code"`
}

type UsercodeRequest struct {
	UserCode         string `json:"user_code" binding:"required"`
	SubscriptionDate string `json:"subscriptiondate"`
}

type LogoutToken struct {
	RefreshToken string `json:"refresh_token"`
}

type UserToken struct {
	Token string `json:"token"`
}

type UserID struct {
	UserId string `json:"userid"`
}

type UserFinal struct {
	IsDeleted bool `json:"isdeleted"`
}

// Verify Email
type VerifyEmail struct {
	Email              string `json:"email"`
	ResetPasswordToken string `json:"resetPasswordToken"`
	Password           string `json:"password"`
	Source             string `json:"source"` /* sync usecase */
}
type EmailOtpRecord struct {
	Phone   string    `json:"phone"`
	Message string    `json:"message"`
	SentOn  time.Time `json:"sentOn"`
}
type GetOtpDetails struct {
	Message string    `json:"message"`
	SentOn  time.Time `json:"sent_on"`
}

// End of verif Email

// Update user details by user id

type UpdateUserDetails struct {
	Country    int    `json:"countryId"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	LanguageId int    `json:"languageId"`
}

// user filters
type CountryDetailsRespone struct {
	EnglishName string `json:"name"`
	Id          int    `json:"id"`
}

type DeviceResponse struct {
	Id       int    `json:"id"`
	Platform string `json:"name"`
}

type Status struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type UserManagementFilterResponse struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type TotalResponse struct {
	CountryDetailsRespone   []CountryDetailsRespone        `json:"countries"`
	DeviceResponse          []DeviceResponse               `json:"devicePlatforms"`
	UserstatusResponse      []Status                       `json:"userStatuses"`
	PaycmsConfirmedResponse []Status                       `json:"paycmsConfirmed"`
	PhoneNumberConfirmeds   []Status                       `json:"phoneNumberConfirmeds"`
	VerificationStatuses    []Status                       `json:"verificationStatuses"`
	EmailConfirmeds         []Status                       `json:"emailConfirmeds"`
	NewsLetters             []Status                       `json:"newsLetters"`
	PromotionEnabled        []Status                       `json:"promotionEnabled"`
	SourceTypes             []Status                       `json:"sourceTypes"`
	UserManagementFilter    []UserManagementFilterResponse `json:"userManagementFilter"`
	UserLeads               []Status                       `json:"userLeads"`
}

// user view activity by filters
type ViewActivityRespone struct {
	Id                string    `json:"viewActivityId,omitempty"`
	ViewedAt          time.Time `json:"viewedAt,omitempty"`
	LastWatchPosition int       `json:"lastWatchPositionSeconds,omitempty"`
	IsHidden          bool      `json:"isHidden,omitempty"`
	ContentId         string    `json:"content_id,omitempty"`
	DeviceId          string    `json:"device_id,omitempty"`
	PlaybackItemId    string    `json:"play_back_item_id,omitempty"`
}

type FinalUserRespone struct {
	ViewActivityId       string    `json:"viewActivityId"`
	ViewedAt             time.Time `json:"viewedAt"`
	Title                string    `json:"title"`
	ContentType          string    `json:"contentType"`
	LastWatchPosition    int       `json:"lastWatchPositionSeconds"`
	IsHidden             bool      `json:"isHidden"`
	DeviceId             string    `json:"device_id"`
	Genres               []string  `json:"genres"`
	ViewedOnPlatformName string    `json:"viewedOnPlatformName"`
	DurationSeconds      int       `json:"durationSeconds"`
	HasWatchingIssues    bool      `json:"hasWatchingIssues"`
}
type ContentDetails struct {
	Id                  string    `json:"viewActivityId,omitempty"`
	ViewedAt            time.Time `json:"viewedAt,omitempty"`
	TransliteratedTitle string    `json:"transliterated_title,omitempty"`
	ContentType         string    `json:"contentType,omitempty"`
	LastWatchPosition   int       `json:"lastWatchPositionSeconds,omitempty"`
	DeviceId            string    `json:"device_id,omitempty"`
	IsHidden            bool      `json:"isHidden"`
	EnglishName         string    `json:"name"`
	Duration            int       `json:"durationSeconds"`
}
type DeviceName struct {
	Platform int `json:"platform"`
}

//  Get Users List and search by filters with pagination

type FiltersFinalResponse struct {
	Id                             string    `json:"id"`
	Status                         int       `json:"status"`
	StatusName                     string    `json:"statusName"`
	TailoredGenres                 string    `json:"tailoredGenres"`
	RegisteredAt                   time.Time `json:"registeredAt"`
	ActiveDevicePlatformNames      string    `json:"activeDevicePlatformNames"`
	NumberOfActiveDevices          int       `json:"numberOfActiveDevices"`
	RegistrationSource             int       `json:"registrationSource"`
	RegistrationSourceName         string    `json:"registrationSourceName"`
	PromotionsEnabledDisplayName   string    `json:"promotionsEnabledDisplayName"`
	NewslettersEnabledDisplayName  string    `json:"newslettersEnabledDisplayName"`
	VerificationEnabledDisplayName string    `json:"verificationEnabledDisplayName"`
	FirstName                      string    `json:"firstName"`
	LastName                       string    `json:"lastName"`
	Email                          string    `json:"email"`
	NewslettersEnabled             bool      `json:"newslettersEnabled"`
	PromotionsEnabled              bool      `json:"promotionsEnabled"`
	Country                        int       `json:"countryId"`
	CountryName                    string    `json:"countryName"`
	IsAdult                        bool      `json:"isAdult"`
	UserLead                       string    `json:"userLead"`
	PrivacyPolicy                  bool      `json:"privacyPolicy"`
	IsRecommend                    bool      `json:"isRecommend"`
	Performance                    bool      `json:"performance"`
	GoogleAnalytics                bool      `json:"googleAnalytics"`
	Firebase                       bool      `json:"firebase"`
	AppFlyer                       bool      `json:"appFlyer"`
	Advertising                    bool      `json:"advertising"`
	Aique                          bool      `json:"aique"`
	GoogleAds                      bool      `json:"googleAds"`
	FacebookAds                    bool      `json:"facebookAds"`
	IsGdprAccepted                 bool      `json:"isGdprAccepted"`
	LanguageId                     int       `json:"languageId"`
	LanguageName                   string    `json:"languageName"`
	PhoneNumber                    string    `json:"phoneNumber"`
	EmailConfirmed                 bool      `json:"emailConfirmed"`
	PhoneNumberConfirmed           bool      `json:"phoneNumberConfirmed"`
	VerificationStatus             bool      `json:"verificationStatus"`
}

type RegistrationSource struct {
	Id         string `json:"id,omitempty"`
	SourceName string `json:"sourncename,omitempty"`
}
type PlatformValues struct {
	Platform int `json:"platform,omitempty"`
}
type Platformdetails struct {
	Platform string `json:"platform"`
}
type Role struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type userExportResponse struct {
	FirstName             string    `json:"firstName"`
	LastName              string    `json:"lastName"`
	Status                string    `json:"status"`
	Country               string    `json:"countryId"`
	RegisteredAt          time.Time `json:"registeredAt"`
	Email                 string    `json:"email" gorm:"unique"`
	PhoneNumber           string    `json:"phoneNumber" gorm:"unique"`
	TailoredGenres        string    `json:"tailoredGenres"`
	ActiveDevices         string    `json:"activeDevices"`
	NumberOfActiveDevices int       `json:"numberOfActiveDevices"`
	LanguageId            int       `json:"languageId"`
	NewslettersEnabled    bool      `json:"newslettersEnabled"`
	PromotionsEnabled     bool      `json:"promotionsEnabled"`
	RegistrationSource    int       `json:"registrationSource"`
	UserLead              string    `json:"userLead"`
	Verified              string    `json:"verified"`
}

type AdminDetails struct {
	Email string
}

type CountryDetails struct {
	Id int
}

// userslist
// User - struct for DB binding
type Users struct {
	ID                   string    `json:"id" gorm:"primary_key"`
	Country              int       `json:"countryId"`
	Status               int       `json:"status"`
	FirstName            string    `json:"firstName"`
	IsBackOfficeUser     bool      `json:"isBackOfficeUser"`
	LanguageId           int       `json:"languageId"`
	LastName             string    `json:"lastName"`
	NewslettersEnabled   bool      `json:"newslettersEnabled"`
	PromotionsEnabled    bool      `json:"promotionsEnabled"`
	LastActivityAt       time.Time `json:"lastActivityAt"`
	RegistrationSource   int       `json:"registrationSource"`
	RegisteredAt         time.Time `json:"registeredAt"`
	Email                string    `json:"email" gorm:"unique"`
	EmailConfirmed       bool      `json:"emailConfirmed"`
	PasswordHash         string    `json:"passwordHash" binding:"required"`
	SecurityStamp        string    `json:"securityStamp"`
	UserName             string    `json:"userName"`
	IsDeleted            bool      `json:"isDeleted"`
	PhoneNumber          string    `json:"phoneNumber" gorm:"unique"`
	PhoneNumberConfirmed bool      `json:"phoneNumberConfirmed"`
	CallingCode          string    `json:"callingCode"`
	NationalNumber       string    `json:"nationalNumber"`
	CountryName          string    `json:"countryName"`
	SearchableText       string    `json:"searchableText"`
	Paycmsstatus         bool      `json:"payCMSStatus"`
	IsAdult              bool      `json:"isAdult"`
	PrivacyPolicy        bool      `json:"privacyPolicy"`
	IsRecommend          bool      `json:"isRecommend"`
	UserLead             string    `json:"userLead"`
	Performance          bool      `json:"performance"`
	GoogleAnalytics      bool      `json:"googleAnalytics"`
	Firebase             bool      `json:"firebase"`
	AppFlyer             bool      `json:"appFlyer"`
	Advertising          bool      `json:"advertising"`
	Aique                bool      `json:"aique"`
	GoogleAds            bool      `json:"googleAds"`
	FacebookAds          bool      `json:"facebookAds"`
	IsGdprAccepted       bool      `json:"isGdprAccepted"`
	SaltStored           string    `json:"-"`
	Version              int       `json:"-"`
	RoleId               string    `json:"-"`
	SourceName           string    `json:"registrationSourceName"`
	Platform             string    `json:"platform"`
	Token                string    `json:"token"`
}

/*Get User Ratings details with search text*/
type RatingRecords struct {
	Id          string    `json:"id"`
	RatedAt     time.Time `json:"ratedAt"`
	Title       string    `json:"title"`
	ContentType string    `json:"contentType"`
	DeviceId    string    `json:"deviceId"`
	Rating      float64   `json:"rating"`
	IsHidden    bool      `json:"isHidden"`
}

type RatingByUser struct {
	RatedAt             time.Time `json:"ratedAt"`
	Title               string    `json:"title"`
	ContentType         string    `json:"contentType"`
	Genres              []string  `json:"genres"`
	RatedOnPlatformName string    `json:"ratedOnPlatformName"`
	Rating              float64   `json:"rating"`
	IsHidden            bool      `json:"isHidden"`
}
type Genname struct {
	EnglishName string `json:"englishName"`
}
type PlatformName struct {
	RatedOnPlatformName int `json:"ratedOnPlatformName"`
}

// watching issue
type WatchingIssue struct {
	ReportedAt      time.Time `json:"reportedAt"`
	IsVideo         bool      `json:"isVideo"`
	IsSound         bool      `json:"isSound"`
	IsTranslation   bool      `json:"isTranslation"`
	IsCommunication bool      `json:"isCommunication"`
	Description     string    `json:"description"`
}

/*oauth2_tokens */
type Oauth2Tokens struct {
	Id        int            `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	ExpiresAt time.Time      `json:"expiresAt"`
	Code      string         `json:"code"`
	Access    string         `json:"access"`
	Refresh   string         `json:"refresh"`
	Data      postgres.Jsonb `json:"data"`
}
