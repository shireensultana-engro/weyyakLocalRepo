package faq

type Faq struct {
	Id          string `json:"id" gorm:"primary_key" swaggerignore:"true"`
	Question    string `json:"question"`
	Description string `json:"description"`
}
