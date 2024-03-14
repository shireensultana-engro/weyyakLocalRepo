package admin

import "time"

type AdminDetailsRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserRole  string `json:"userRole"`
	Password  string `json:"password"`
}

type PaginationResult struct {
	Size   int   `json:"size"`
	Offset int   `json:"offset"`
	Limit  int64 `json:"limit"`
}

type UserKeys struct {
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Email        string    `json:"email"`
	UserRole     string    `json:"userRole"`
	AllowDelete  bool      `json:"allowDelete"`
	RegisteredAt time.Time `json:"registeredAt"`
	Id           string    `json:"id"`
}

// ERROR-CODES
type Emailcheck struct {
	Email string `json:"email"`
}
type Email struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type Invalid struct {
	Email *Email `json:"email,omitempty"`
}
type FinalErrorResponse struct {
	Error       string  `json:"error"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	RequestId   string  `json:"requestId"`
	Invalid     Invalid `json:"invalid,omitempty"`
}

type ResetPasswordAdminRequest struct {
	Email string `json:"email"`
}

type RequestSetPasswordParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
type EmailOtpRecord struct {
	Phone   string    `json:"phone"`
	Message string    `json:"message"`
	SentOn  time.Time `json:"sentOn"`
}
