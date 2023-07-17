package core

import (
	"time"
)

type Item struct {
	Id          int       `json:"id,omitempty"`
	CampaignId  int       `json:"campaignId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,default:-,omitempty"`
	Priority    int       `json:"priority,omitempty"`
	Removed     bool      `json:"removed,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}
