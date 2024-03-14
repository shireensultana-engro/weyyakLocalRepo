package geoblock

//HandlerService ... export handlebar
type HandlerService struct{}

//NewsLetter ... newsletter request details
type GeoBlock struct {
	Email string `json:"email" gorm:"email"`
}
