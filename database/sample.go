package database

import (
	"time"

	"github.com/lib/pq"
)

// Type Sample describes a single sample received by a sensor
type Sample struct {
	// The unique hardware id of the sensor
	SensorID string `json:"sensorId" avro:"sensorId"`

	// Unique identifier for the campaign this sample belongs to
	CampaignId string `json:"campaignId" avro:"campaignId"`

	// Sample type string (IQ, PSD, DEC)
	SampleType string `json:"sampleType" avro:"sampleType"`

	// Sample timestamp with microseconds precision
	Time SampleTime `json:"time" avro:"time" gorm:"embedded;embeddedPrefix:time_"`

	// Sensor configuration for the recorded data set
	Config SampleConfig `json:"config" avro:"config" gorm:"embedded;embeddedPrefix:config_"`

	// Sample loss rate in percent
	// LossRate *float64 `json:"lossRate,omitempty" avro:"lossRate"`

	// Method used to obfuscate IQ spectrum data
	// Obfuscation *string `json:"obfuscation,omitempty" avro:"obfuscation"`

	// Actual measurement data. Unit depends on measurement type
	Data pq.Float32Array `json:"data" avro:"data" gorm:"type:numeric[]"`

	// Database-oriented data
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

// Type SampleTime contains timing information for a single sample
type SampleTime struct {
	// Number of seconds since the UNIX epoch start on January 1st, 1970 at UTC
	Seconds int64 `json:"seconds" avro:"seconds"`

	// Microseconds extension for the UNIX time stamp
	Microseconds int32 `json:"microseconds" avro:"microseconds"`
}

// Type SampleConfig describes sensor configuration in use when collection a sample
type SampleConfig struct {
	// Antenna gain in dBi
	AntennaGain *float32 `json:"antennaGain,omitempty" avro:"antennaGain"`

	// Identifier for the antenna being used if device has multiple antennas
	AntennaID *string `json:"antennaId,omitempty" avro:"antennaId"`

	// Center frequency in Hz to which the RF front-end was tuned to while recording the associated spectrum data
	CenterFreq int64 `json:"centerFreq" avro:"centerFreq"`

	// Estimated noise floor in dB
	EstNoiseFloor *float32 `json:"estNoiseFloor,omitempty" avro:"estNoiseFloor"`

	// Correction factor for center frequency in Hz. The correction is already applied to the center frequency (0.0 for no correction)
	FrequencyCorrectionFactor *float32 `json:"frequencyCorrectionFactor,omitempty" avro:"frequencyCorrectionFactor"`

	// RF front-end gain in dB (-1 for automatic gain control)
	FrontendGain *float32 `json:"frontendGain,omitempty" avro:"frontendGain"`

	// Hopping strategy  used to overcome the bandwidth limitations of the RF front-end (0:Sequential, 1:Random, 2:Similarity)
	HoppingStrategy *int32 `json:"hoppingStrategy,omitempty" avro:"hoppingStrategy"`

	// True if IQ samples are balanced
	IqBalanceCalibration *bool `json:"iqBalanceCalibration,omitempty" avro:"iqBalanceCalibration"`

	// Time synchronization of the radio frontend (0: none, 1: GPS, 2: Reference Clock, 5: Other)
	RFSync *string `json:"rfSync,omitempty" avro:"rfSync"`

	// Sensor's sampling rate in samples per second
	SamplingRate *int32 `json:"samplingRate,omitempty" avro:"samplingRate"`

	// True if signal strength is calibrated
	SigStrengthCalibration *bool `json:"sigStrengthCalibration,omitempty" avro:"sigStrengthCalibration"`

	// Time synchronization of the system (0: none, 1: GPS, 2: Reference Clock, 3: NTP, 4: OpenSky, 5: Other)
	SystemSync *string `json:"systemSync,omitempty" avro:"systemSync"`

	// Extra configuration for arbitrary data
	ExtraConf map[string]interface{} `json:"extraConf,omitempty" avro:"extraConf" gorm:"serializer:json"`
}
