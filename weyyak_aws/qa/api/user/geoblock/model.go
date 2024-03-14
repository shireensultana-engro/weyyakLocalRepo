package geoblock

//HandlerService ... export handlebar
type HandlerService struct{}

//NewsLetter ... newsletter request details
type NewsLetter struct {
	Email string `json:"email" gorm:"email"`
}
