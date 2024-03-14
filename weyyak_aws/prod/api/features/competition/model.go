package competition

type CompetitionUsers struct {
	Id                 int    `json:"id" gorm:"primary_key"`
	EnglishFullName    string `json:"englishFullName"`
	ArabicFullName     string `json:"arabicFullName"`
	Alpha2code         string `json:"alpha2Code" binding:"required"`
	Mobile             string `json:"mobile" binding:"required"`
	Email              string `json:"email" binding:"required"`
	AgeGroup           string `json:"ageGroup" binding:"required"`
	Gender             string `json:"gender" binding:"required"`
	CountryOfResidence string `json:"countryOfResidence" binding:"required"`
	Nationality        string `json:"nationality" binding:"required"`
	AdultConsent       bool   `json:"adultConsent" binding:"required"`
	TravelConsent      bool   `json:"travelConsent" binding:"required"`
	Language           string `json:"language"`
	MailSent           string `json:"mailSent"`
	EmailConfirmed     bool   `json:"emailConfirmed"`
}

type Country struct {
	Id          int    `json:"id"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
	CallingCode string `json:"callingCode"`
	Alpha2code  string `json:"alpha2code"`
}
type AgeGroup struct {
	Id       string `json:"id"`
	AgeGroup string `json:"ageGroup"`
}

type Emailcheck struct {
	Email string `json:"email"`
}

// ERROR RESPONSE
type ErrorCode struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

type Invalid struct {
	Email *ErrorCode `json:"email,omitempty"`
}

type FinalErrorResponse struct {
	Error       string  `json:"error"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	RequestId   string  `json:"requestId"`
	Invalid     Invalid `json:"invalid,omitempty"`
}
