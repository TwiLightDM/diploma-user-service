package entities

import "gorm.io/gorm"

type GroupMember struct {
	Id        string         `json:"id"`
	UserId    string         `json:"user_id"`
	GroupId   string         `json:"group_id"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
