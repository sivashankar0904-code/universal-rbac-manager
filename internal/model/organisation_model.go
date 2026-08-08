package model

import (
	"encoding/json"
	"time"
)

type Organisation struct {
	ID             string          `json:"id"`
	Key            string          `json:"key"`
	Name           string          `json:"name"`
	Description    *string         `json:"description,omitempty"`
	ParentID       *string         `json:"parent_id,omitempty"`
	PrimaryContact *Contact        `json:"primaryContact,omitempty"`
	PrimaryAddress *Address        `json:"primaryAddress,omitempty"`
	Contacts       []Contact       `json:"contacts,omitempty"`
	Addresses      []Address       `json:"addresses,omitempty"`
	IsActive       bool            `json:"is_active"`
	Meta           json.RawMessage `json:"meta,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
