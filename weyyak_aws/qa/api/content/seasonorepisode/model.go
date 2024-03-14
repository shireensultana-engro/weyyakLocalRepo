package seasonorepisode

import (
	"time"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Details struct {
	Id string `json:"id"`
}

type UpdateDetails struct {
	DeletedByUserId string    `json:"deleted_by_user_id"`
	ModifiedAt      time.Time `json:"modified_at"`
}
type UpdateStatus struct {
	Status     string    `json:"status"`
	ModifiedAt time.Time `json:"modified_at"`
}
type UpdateRights struct {
	DigitalRightsType string `json:"digital_rights_type"`
}

// Error codes
type FinalErrorResponse struct {
	Error       string  `json:"error"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	RequestId   string  `json:"requestId"`
	Invalid     Invalid `json:"invalid,omitempty"`
}
type Invalid struct {
	Id *InvalidError `json:"id"`
}

type InvalidError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
