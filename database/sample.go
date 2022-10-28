package database

import (
	"time"

	"github.com/lib/pq"
)

// Type Sample describes a single sample received by a sensor
type Sample struct {
	// The unique hardware id of the sensor
	SensorID string `json:"sensorId"`

	// Unique identifier for the campaign this sample belongs to
	CampaignId string `json:"campaignId"`

	// Sample type string (IQ, PSD, DEC)
	SampleType string `json:"sampleType"`

	// Sample timestamp with microseconds precision
	Time SampleTime `json:"time" gorm:"embedded;embeddedPrefix:time_"`

	// Sensor configuration for the recorded data set
	Config SampleConfig `json:"config" gorm:"embedded;embeddedPrefix:config_"`

	// Sample loss rate in percent
	LossRate *float64 `json:"lossRate,omitempty"`

	// Method used to obfuscate IQ spectrum data
	Obfuscation *string `json:"obfuscation,omitempty"`

	// Actual measurement data. Unit depends on measurement type
	Data pq.Float64Array `json:"data" gorm:"type:numeric[]"`

	// Database-oriented data
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

// Type SampleTime contains timing information for a single sample
type SampleTime struct {
	// Number of seconds since the UNIX epoch start on January 1st, 1970 at UTC
	Seconds int64 `json:"seconds"`

	// Microseconds extension for the UNIX time stamp
	Microseconds int64 `json:"microseconds"`
}

// Type SampleConfig describes sensor configuration in use when collection a sample
type SampleConfig struct {
	// Antenna gain in dBi
	AntennaGain *float64 `json:"antennaGain,omitempty"`

	// Identifier for the antenna being used if device has multiple antennas
	AntennaID *string `json:"antennaId,omitempty"`

	// Center frequency in Hz to which the RF front-end was tuned to while recording the associated spectrum data
	CenterFreq int64 `json:"centerFreq"`

	// Estimated noise floor in dB
	EstNoiseFloor *float64 `json:"estNoiseFloor,omitempty"`

	// Correction factor for center frequency in Hz. The correction is already applied to the center frequency (0.0 for no correction)
	FrequencyCorrectionFactor *float64 `json:"frequencyCorrectionFactor,omitempty"`

	// RF front-end gain in dB (-1 for automatic gain control)
	FrontendGain *float64 `json:"frontendGain,omitempty"`

	// Hopping strategy  used to overcome the bandwidth limitations of the RF front-end (0:Sequential, 1:Random, 2:Similarity)
	HoppingStrategy *int64 `json:"hoppingStrategy,omitempty"`

	// True if IQ samples are balanced
	IqBalanceCalibration *bool `json:"iqBalanceCalibration,omitempty"`

	// Time synchronization of the radio frontend (0: none, 1: GPS, 2: Reference Clock, 5: Other)
	RFSync *string `json:"rfSync,omitempty"`

	// Sensor's sampling rate in samples per second
	SamplingRate *int64 `json:"samplingRate,omitempty"`

	// True if signal strength is calibrated
	SigStrengthCalibration *bool `json:"sigStrengthCalibration,omitempty"`

	// Time synchronization of the system (0: none, 1: GPS, 2: Reference Clock, 3: NTP, 4: OpenSky, 5: Other)
	SystemSync *string `json:"systemSync,omitempty"`

	// Extra configuration for arbitrary data
	ExtraConf map[string]interface{} `json:"extraConf,omitempty" gorm:"serializer:json"`
}
