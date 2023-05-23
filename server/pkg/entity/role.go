package entity

import "time"

type Role struct {
	Id             int       `gorm:"primaryKey" json:"id"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"createdAt"`
	LastModifiedAt time.Time `json:"lastModifiedAt"`
	IsDeleted      bool      `json:"isDeleted"`
}
