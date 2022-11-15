package models

import (
	"time"

	"github.com/lib/pq"
)

// Type Campaign represents a measurement campaign which has been successfully launched
// and stored in the database.
type Campaign struct {
	// The textual, random ID for the campaign
	CampaignId string `json:"campaignId" db:"campaign_id"`

	// The list of sensor partaking in the campaign
	Sensors pq.StringArray `json:"sensors"`

	// The type of measurements requested
	Type string `json:"type"`

	// The time at which the campaign is supposed to start
	Begin time.Time `json:"begin"`

	// The time at which the campaign will end
	End time.Time `json:"end"`

	// Database-specific data
	ID        uint      `json:"-"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
