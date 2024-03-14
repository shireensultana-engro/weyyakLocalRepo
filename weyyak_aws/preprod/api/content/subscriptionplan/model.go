package subscriptionplan

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type SubscriptionPlan struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
