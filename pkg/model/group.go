package model

import "time"

type Group struct {
	ID        ID         `json:"groupId"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
