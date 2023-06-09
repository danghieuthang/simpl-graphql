// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type AuthenType struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PageArgs struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type PageInfo struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type Role struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type RolePage struct {
	Data     []*Role   `json:"data"`
	PageInfo *PageInfo `json:"pageInfo"`
}

type User struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	CreatedAt      time.Time  `json:"createdAt"`
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`
	CreatedBy      *string    `json:"createdBy,omitempty"`
	UpdatedBy      *string    `json:"updatedBy,omitempty"`
	Role           *Role      `json:"role,omitempty"`
}

type UserInput struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPage struct {
	Data     []*User   `json:"data"`
	PageInfo *PageInfo `json:"pageInfo"`
}
