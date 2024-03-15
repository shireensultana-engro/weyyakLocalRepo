package register

import (
	"time"

	_ "github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User - struct for DB binding
type User struct {
	Id                   string    `json:"id" gorm:"type:uuid;primaryKey"`
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
	NickName             string    `json:"nickName"`
	Email                string    `json:"email" gorm:"unique"`
	EmailConfirmed       bool      `json:"emailConfirmed"`
	PasswordHash         string    `json:"passwordHash" binding:"required"`
	SecurityStamp        string    `json:"securityStamp"`
	UserName             string    `json:"userName"`
	IsDeleted            bool      `json:"isDeleted"`
	PhoneNumber          string    `json:"phoneNumber" gorm:"unique"`
	PhoneNumberConfirmed bool      `json:"phoneNumberConfirmed"`
	DeleteInitiatesAt    time.Time `gorm:"default:null" json:"delete_initiates_at"`
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
	ModifiedAt           time.Time `json:"modified_at"`
	AppleUserId          bool      `json:"appleDeviceId"`
	AppleSignIn          bool      `json:"appleSignIn"`
	CleverTap            bool      `json:"cleverTap"`
	RegistrationPlatform string    `json:"registrationPlatform"`
}

// RequestRegisterUserUsingEmail - struct for DB binding
type RequestRegisterUserUsingEmail struct {
	Password             string `json:"password"`
	LanguageId           int    `json:"languageId"`
	PrivacyPolicy        bool   `json:"PrivacyPolicy"`
	IsAdult              bool   `json:"isAdult"`
	IsRecommend          bool   `json:"IsRecommend"`
	Email                string `json:"email"`
	CountryName          string `json:"countryName"`
	Alpha2code           string `json:"Alpha2code"`
	Source               string `json:"source"`                             /* sync usecase */
	UserId               string `json:"userId" gorm:"type:uuid;primaryKey"` /* sync usecase */
	RegistrationPlatform string `json:"registrationPlatform"`
}

// RequestRegisterUserUsingSMS - struct for DB binding
type RequestRegisterUserUsingSMS struct {
	Password             string `json:"password"`
	LanguageId           int    `json:"languageId"`
	PrivacyPolicy        bool   `json:"PrivacyPolicy"`
	IsAdult              bool   `json:"isAdult"`
	IsRecommend          bool   `json:"IsRecommend"`
	PhoneNumber          string `json:"phonenumber"`
	Silentregistration   bool   `json:"silentregistration"`
	Source               string `json:"source"` /* sync usecase */
	UserId               string `json:"userId"` /* sync usecase */
	RecaptchaToken       string `json:"recaptcha_token"`
	DeviceId             string `json:"DeviceId"`
	RegistrationPlatform string `json:"registrationPlatform"`
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
	FirstName            *string     `json:"firstName"`
	LastName             *string     `json:"lastName"`
	NickName             string      `json:"nickName"`
	Email                string      `json:"email" gorm:"unique"`
	NewslettersEnabled   interface{} `json:"newslettersEnabled"`
	PromotionsEnabled    interface{} `json:"promotionsEnabled"`
	Country              int         `json:"countryId,omitempty"`
	CountryName          string      `json:"countryName"`
	IsAdult              interface{} `json:"isAdult"`
	UserLead             string      `json:"userLead"`
	PrivacyPolicy        interface{} `json:"privacyPolicy"`
	IsRecommend          interface{} `json:"isRecommend"`
	Performance          interface{} `json:"performance"`
	GoogleAnalytics      interface{} `json:"googleAnalytics"`
	Firebase             interface{} `json:"firebase"`
	AppFlyer             interface{} `json:"appFlyer"`
	Advertising          interface{} `json:"advertising"`
	Aique                interface{} `json:"aique"`
	GoogleAds            interface{} `json:"googleAds"`
	FacebookAds          interface{} `json:"facebookAds"`
	IsGdprAccepted       interface{} `json:"isGdprAccepted"`
	CleverTap            interface{} `json:"cleverTap"`
	LanguageId           int         `json:"languageId"`
	RegistrationSource   int         `json:"registrationSource"`
	PhoneNumber          string      `json:"phoneNumber" gorm:"unique"`
	EmailConfirmed       bool        `json:"emailConfirmed"`
	PhoneNumberConfirmed bool        `json:"phoneNumberConfirmed"`
	AppleUserId          bool        `json:"appleDeviceId"`
	AppleSignIn          bool        `json:"appleSignIn"`
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
	NickName             string `json:"nickName"`
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
	CallingCode          string `json:"callingCode"`
	PhoneNumber          string `json:"phoneNumber"`
	NationalNumber       string `json:"nationalNumber"`
	EmailConfirmed       bool   `json:"emailConfirmed"`
	PhoneNumberConfirmed bool   `json:"phoneNumberConfirmed"`
	VerificationStatus   bool   `json:"verificationStatus"`
	AppleUserId          bool   `json:"appleDeviceId"`
	AppleSignIn          bool   `json:"appleSignIn"`
	UserCount            int    `json:"user_count"`
	RegistrationPlatform string `json:"registrationPlatform"`
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
	PhoneNumber          string    `json:"phoneNumber"`
	PhoneNumberConfirmed bool      `json:"phone_number_confirmed"`
	DeleteInitiatesAt    time.Time `json:"delete_initiates_at"`
}

type Emailcheck struct {
	Email             string    `json:"email"`
	DeleteInitiatesAt time.Time `json:"delete_initiates_at"`
	EmailConfirmed    bool      `json:"email_confirmed"`
	UserId            string    `json:"userId"`
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
	Id                    string `json:"id"`
	FirstName             string `json:"firstName"`
	LastName              string `json:"lastName"`
	Status                string `json:"status"`
	Country               string `json:"country"`
	RegisteredAt          string `json:"registeredAt"`
	Email                 string `json:"email" gorm:"unique"`
	PhoneNumber           string `json:"phone_number" gorm:"unique"`
	TailoredGenres        string `json:"tailored_genres"`
	ActiveDevices         string `json:"active_devices_name"`
	NumberOfActiveDevices int    `json:"active_device"`
	Language              string `json:"language"`
	NewslettersEnabled    string `json:"newsletters_enabled"`
	PromotionsEnabled     string `json:"promotions_enabled"`
	RegistrationSource    string `json:"registration_source"`
	UserLead              string `json:"user_lead"`
	VerificationStatus    string `json:"verification_status"`
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
	NickName             string    `json:"nickName"`
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

type TwitterUser struct {
	ContributorsEnabled bool   `json:"contributors_enabled"`
	CreatedAt           string `json:"created_at"`
	DefaultProfile      bool   `json:"default_profile"`
	DefaultProfileImage bool   `json:"default_profile_image"`
	Description         string `json:"description"`
	Email               string `json:"email"`
	Entities            struct {
		Description struct {
			Urls []interface{} `json:"urls"`
		} `json:"description"`
	} `json:"entities"`
	FavouritesCount                int64         `json:"favourites_count"`
	FollowRequestSent              bool          `json:"follow_request_sent"`
	FollowersCount                 int64         `json:"followers_count"`
	Following                      bool          `json:"following"`
	FriendsCount                   int64         `json:"friends_count"`
	GeoEnabled                     bool          `json:"geo_enabled"`
	HasExtendedProfile             bool          `json:"has_extended_profile"`
	ID                             float64       `json:"id"`
	IDStr                          string        `json:"id_str"`
	IsTranslationEnabled           bool          `json:"is_translation_enabled"`
	IsTranslator                   bool          `json:"is_translator"`
	Lang                           string        `json:"lang"`
	ListedCount                    int64         `json:"listed_count"`
	Location                       string        `json:"location"`
	Name                           string        `json:"name"`
	NeedsPhoneVerification         bool          `json:"needs_phone_verification"`
	Notifications                  bool          `json:"notifications"`
	ProfileBackgroundColor         string        `json:"profile_background_color"`
	ProfileBackgroundImageURL      interface{}   `json:"profile_background_image_url"`
	ProfileBackgroundImageURLHttps interface{}   `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool          `json:"profile_background_tile"`
	ProfileImageURL                string        `json:"profile_image_url"`
	ProfileImageURLHttps           string        `json:"profile_image_url_https"`
	ProfileLinkColor               string        `json:"profile_link_color"`
	ProfileSidebarBorderColor      string        `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string        `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string        `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool          `json:"profile_use_background_image"`
	Protected                      bool          `json:"protected"`
	ScreenName                     string        `json:"screen_name"`
	StatusesCount                  int64         `json:"statuses_count"`
	Suspended                      bool          `json:"suspended"`
	TimeZone                       string        `json:"time_zone"`
	TranslatorType                 string        `json:"translator_type"`
	URL                            interface{}   `json:"url"`
	UTCOffset                      interface{}   `json:"utc_offset"`
	Verified                       bool          `json:"verified"`
	WithheldInCountries            []interface{} `json:"withheld_in_countries"`
}

type DeleteDeviceUserLogout struct {
	DeviceID string `json:"deviceid"`
}

type DeviceToken struct {
	UserId   string `json:"userid"`
	DeviceId string `json:"deviceid"`
	Token    string `json:"token"`
}

type Addcontent struct {
	ContentRequest ContentRequest `json:"content"`
	Rating         int            `json:"rating"`
}

type ContentRequest struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	ContentType string   `json:"contentType"`
	Duration    int      `json:"duration"`
	Genres      []string `json:"genres"`
}

type ContentId struct {
	Id string `json:"id"`
}

type CreateRatedContent struct {
	Rating    int       `json:"rating"`
	RatedAt   time.Time `json:"rated_at"`
	ContentId string    `json:"content_id"`
	UserID    string    `json:"user_id"`
	IsHidden  bool      `json:"is_hidden"`
	DeviceId  string    `json:"device_id"`
}

type AddRatingToHistory struct {
	Rating    int       `json:"rating"`
	RatedAt   time.Time `json:"rated_at"`
	ContentId string    `json:"content_id"`
	UserID    string    `json:"user_id"`
	DeviceId  string    `json:"device_id"`
}

type ResponseContent struct {
	AverageRating          float64 `json:"average_rating"`
	AverageRatingUpdatedAt time.Time
}

type EpisodeDetailsSummary struct {
	IsPrimary              bool   `json:"isPrimary"`
	UserId                 string `json:"userId"`
	SecondarySeasonId      string `json:"secondarySeasonId" `
	VarianceIds            []int  `json:"varianceIds"`
	EpisodeIds             []int  `json:"episodeIds"`
	SecondaryEpisodeId     string `json:"secondaryEpisodeId"`
	ContentId              string `json:"contentId"`
	EpisodeKey             int    `json:"episodeKey"`
	SeasonId               string `json:"seasonId"`
	Status                 int    `json:"status"`
	StatusCanBeChanged     bool   `json:"statusCanBeChanged"`
	SubStatus              int    `json:"subStatus"`
	SubStatusName          string `json:"subStatusName"`
	DigitalRightsType      int    `json:"digitalRightsType"`
	DigitalRightsStartDate string `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   string `json:"digitalRightsEndDate"`
	CreatedBy              string `json:"createdBy"`
	PrimaryInfo            PrimaryInfo
	Cast                   []int `json:"cast"`
	Music                  []int `json:"music"`
	TagInfo                []int `json:"tagInfo"`
	NonTextualData         []int `json:"nonTextualData"`
	Translation            Translation
	SchedulingDateTime     []int  `json:"schedulingDateTime"`
	PublishingPlatforms    []int  `json:"publishingPlatforms"`
	SeoDetails             []int  `json:"seoDetails"`
	Id                     string `json:"id"`
}

type PrimaryInfo struct {
	Number              int     `json:"number ,omitempty"`
	VideoContentId      string  `json:"videoContentId ,omitempty"`
	SynopsisEnglish     string  `json:"synopsisEnglish ,omitempty"`
	SynopsisArabic      string  `json:"synopsisArabic ,omitempty"`
	SeasonNumber        int     `json:"seasonNumber ,omitempty"`
	OriginalTitle       string  `json:"originalTitle"`
	AlternativeTitle    string  `json:"alternativeTitle"`
	ArabicTitle         string  `json:"arabicTitle"`
	TransliteratedTitle string  `json:"transliteratedTitle"`
	Notes               string  `json:"notes"`
	IntroStart          *string `json:"introStart"`
	OutroStart          *string `json:"outroStart"`
}
type Translation struct {
	LanguageType       string  `json:"languageType"`
	DubbingLanguage    *string `json:"dubbingLanguage"`
	DubbingDialectId   *int    `json:"dubbingDialectId"`
	SubtitlingLanguage *string `json:"subtitlingLanguage"`
}

// Get Episode Details Based in contentId
type SeasonDetailsSummary struct {
	ContentId          string    `json:"contentId"`
	SeasonKey          int       `json:"seasonKey"`
	Status             int       `json:"status"`
	StatusCanBeChanged bool      `json:"statusCanBeChanged"`
	SubStatusName      string    `json:"subStatusName"`
	ModifiedAt         time.Time `json:"modifiedAt"`
	PrimaryInfo        PrimaryInfo
	Cast               *string `json:"cast"`
	Music              *string `json:"music"`
	TagInfo            *string `json:"tagInfo"`
	SeasonGenres       *string `json:"seasonGenres"`
	AboutTheContent    *string `json:"aboutTheContent"`
	Translation        Translation
	Episodes           *string `json:"episodes"`
	NonTextualData     *string `json:"nonTextualData"`
	Rights             Rights  `json:"rights"`
	CreatedBy          string  `json:"createdBy"`
	IntroDuration      string  `json:"introDuration"`
	IntroStart         *string `json:"introStart"`
	OutroDuration      string  `json:"outroDuration"`
	OutroStart         *string `json:"outroStart"`
	Products           *string `json:"products"`
	SeoDetails         *string `json:"seoDetails"`
	VarianceTrailers   *string `json:"varianceTrailers"`
	Id                 string  `json:"id"`
}
type Rights struct {
	DigitalRightsType      int    `json:"digitalRightsType"`
	DigitalRightsStartDate string `json:"digitalRightsStartDate"`
	DigitalRightsEndDate   string `json:"digitalRightsEndDate"`
	DigitalRightsRegions   []int  `json:"digitalRightsRegions"`
	SubscriptionPlans      []int  `json:"subscriptionPlans"`
}

type ViewActivity struct {
	Id                string    `json:"id" gorm:"primary_id"`
	ContentId         string    `json:"content_id"`
	UserId            string    `json:"user_id"`
	DeviceId          string    `json:"device_id"`
	LastWatchPosition int       `json:"last_watch_position"`
	ViewedAt          time.Time `json:"viewed_at"`
	WatchSessionId    string    `json:"watch_session_id"`
	IsHidden          string    `json:"is_hidden"`
	PlaybackItemId    string    `json:"playback_item_id"`
}

type ViewActivityHistory struct {
	Id                string    `json:"id,omitempty" gorm:"primary_id"`
	ContentTypeName   string    `json:"content_type_name"`
	ContentKey        int       `json:"content_key"`
	UserId            string    `json:"user_id"`
	DeviceId          string    `json:"device_id"`
	LastWatchPosition int       `json:"last_watch_position"`
	ViewedAt          time.Time `json:"viewed_at"`
	WatchSessionId    string    `json:"watch_session_id"`
}

type ContentIdDetails struct {
	Id             string    `json:"id"`
	ContentKey     int       `json:"content_key"`
	ContentTier    int       `json:"content_tier"`
	CreatedAt      time.Time `json:"created_at"`
	PlaybackItemId string    `json:"playback_item_id,omitempty"`
}

type AddViewActivityRequest struct {
	Content           ViewActivityRequestContent `json:"content"`
	LastWatchPosition int                        `json:"lastWatchPosition"`
	WatchSessionId    string                     `json:"watchSessionId"`
}
type ViewActivityRequestContent struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	ContentType string   `json:"contentType"`
	Duration    int      `json:"duration"`
	Genres      []string `json:"genres"`
}
type UpdateUserCookies struct {
	FirstName            *string     `json:"firstName"`
	LastName             *string     `json:"lastName"`
	NickName             string      `json:"nickName"`
	Email                string      `json:"email" gorm:"unique"`
	NewslettersEnabled   interface{} `json:"newslettersEnabled"`
	PromotionsEnabled    interface{} `json:"promotionsEnabled"`
	Country              int         `json:"countryId,omitempty"`
	CountryName          string      `json:"countryName"`
	IsAdult              interface{} `json:"isAdult"`
	UserLead             string      `json:"userLead"`
	PrivacyPolicy        interface{} `json:"privacyPolicy"`
	IsRecommend          interface{} `json:"isRecommend"`
	Performance          interface{} `json:"performance"`
	GoogleAnalytics      interface{} `json:"googleAnalytics"`
	Firebase             interface{} `json:"firebase"`
	AppFlyer             interface{} `json:"appFlyer"`
	Advertising          interface{} `json:"advertising"`
	Aique                interface{} `json:"aique"`
	GoogleAds            interface{} `json:"googleAds"`
	FacebookAds          interface{} `json:"facebookAds"`
	IsGdprAccepted       interface{} `json:"isGdprAccepted"`
	CleverTap            interface{} `json:"cleverTap"`
	LanguageId           int         `json:"languageId"`
	RegistrationSource   int         `json:"registrationSource"`
	PhoneNumber          string      `json:"phoneNumber" gorm:"unique"`
	EmailConfirmed       bool        `json:"emailConfirmed"`
	PhoneNumberConfirmed bool        `json:"phoneNumberConfirmed"`
	UserId               string      `json:"user_id"`
	Platform             int         `json:"platform"`
	CreatedAt            time.Time   `json:"createdAt"`
	LastActivityAt       time.Time   `json:"last_activity_at"`
	AppleUserId          bool        `json:"appleDeviceId"`
	AppleSignIn          bool        `json:"appleSignIn"`
}
type CookieUserProfileResponse struct {
	FirstName            string    `json:"firstName"`
	LastName             string    `json:"lastName"`
	NickName             string    `json:"nickName"`
	Email                string    `json:"email"`
	NewslettersEnabled   bool      `json:"newslettersEnabled"`
	PromotionsEnabled    bool      `json:"promotionsEnabled"`
	Country              int       `json:"countryId"`
	CountryName          string    `json:"countryName"`
	IsAdult              bool      `json:"isAdult"`
	UserLead             string    `json:"userLead"`
	PrivacyPolicy        bool      `json:"privacyPolicy"`
	IsRecommend          bool      `json:"isRecommend"`
	Performance          bool      `json:"performance"`
	GoogleAnalytics      bool      `json:"googleAnalytics"`
	Firebase             bool      `json:"firebase"`
	AppFlyer             bool      `json:"appFlyer"`
	Advertising          bool      `json:"advertising"`
	Aique                bool      `json:"aique"`
	GoogleAds            bool      `json:"googleAds"`
	FacebookAds          bool      `json:"facebookAds"`
	IsGdprAccepted       bool      `json:"isGdprAccepted"`
	CleverTap            bool      `json:"cleverTap"`
	LanguageId           int       `json:"languageId"`
	LanguageName         string    `json:"languageName"`
	RegistrationSource   int       `json:"registrationSource"`
	CallingCode          string    `json:"callingCode"`
	PhoneNumber          string    `json:"phoneNumber"`
	NationalNumber       string    `json:"nationalNumber"`
	EmailConfirmed       bool      `json:"emailConfirmed"`
	PhoneNumberConfirmed bool      `json:"phoneNumberConfirmed"`
	VerificationStatus   bool      `json:"verificationStatus"`
	UserID               string    `json:"user_id"`
	Platform             int       `json:"platform"`
	LastActivityAt       time.Time `json:"last_activity_at"`
	AppleUserId          bool      `json:"appleDeviceId"`
	AppleSignIn          bool      `json:"appleSignIn"`
	UserCount            int       `json:"user_count"`
	RegistrationPlatform string    `json:"registrationPlatform"`
}
