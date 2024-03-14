package agegroup

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type AgeGroup struct{
	EnglishName string `json:"englishName"`
	ArabicName string `json:"arabicName"`
	Id string `json:"id"`
} 