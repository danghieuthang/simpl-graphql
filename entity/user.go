package entity

import (
	"time"
)

type Base struct {
	CreatedAt      time.Time `json:"createdAt"`
	LastModifiedAt time.Time `json:"lastModifiedAt"`
	IsDeleted      bool      `json:"isDeleted"`
}

type User struct {
	Id             int    `gorm:"primaryKey" json:"id"`
	Name           string `json:"name"`
	Email          string `form:"unique" json:"email"`
	Password       string
	RoleId         *int
	Role           Role      `gorm:"constraint:OnDelete:SET NULL;`
	CreatedAt      time.Time `json:"createdAt"`
	LastModifiedAt time.Time `json:"lastModifiedAt"`
	IsDeleted      bool      `json:"isDeleted"`
}
