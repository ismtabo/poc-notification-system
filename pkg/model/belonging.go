package model

import "time"

type Belonging struct {
	UserID    ID         `json:"userId"`
	GroupID   ID         `json:"groupId"`
	CreatedAt *time.Time `json:"createdAt"`
}
