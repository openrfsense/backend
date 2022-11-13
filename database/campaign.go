package database

import (
	"time"

	"github.com/lib/pq"
)

// Type Campaign represents a measurement campaign which has been successfully launched
// and stored in the database.
type Campaign struct {
	// The textual, random ID for the campaign
	CampaignId string `json:"campaignId" gorm:"not null"`

	// The list of sensor partaking in the campaign
	Sensors pq.StringArray `json:"sensors" gorm:"not null;type:text[]"`

	// The type of measurements requested
	Type string `json:"type" gorm:"not null"`

	// The time at which the campaign is supposed to start
	Begin time.Time `json:"begin" gorm:"not null;type:timestamp with time zone"`

	// The time at which the campaign will end
	End time.Time `json:"end" gorm:"not null;type:timestamp with time zone"`

	// Database-oriented data
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
