package entities

import "gorm.io/gorm"

type Group struct {
	Id          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	OwnerId     string         `json:"owner_id"`
	DeleteAt    gorm.DeletedAt `gorm:"index"`
}
