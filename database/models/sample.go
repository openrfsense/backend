package models

import (
	"time"

	"github.com/lib/pq"
)

// Type Sample describes a single sample received by a sensor
type Sample struct {
	// The unique hardware id of the sensor
	SensorId string `json:"sensorId" avro:"sensorId" db:"sensor_id"`

	// Unique identifier for the campaign this sample belongs to
	CampaignId string `json:"campaignId" avro:"campaignId" db:"campaign_id"`

	// Sample type string (IQ, PSD, DEC)
	SampleType string `json:"sampleType" avro:"sampleType" db:"sample_type"`

	// Sample timestamp with microseconds precision
	SampleTime SampleTime `json:"time" avro:"time" db:"embedded"`

	// Sensor configuration for the recorded data set
	SampleConfig SampleConfig `json:"config" avro:"config" db:"embedded"`

	// Sample loss rate in percent
	// LossRate *float64 `json:"lossRate,omitempty" avro:"lossRate"`

	// Method used to obfuscate IQ spectrum data
	// Obfuscation *string `json:"obfuscation,omitempty" avro:"obfuscation"`

	// Actual measurement data. Unit depends on measurement type
	Data pq.Float32Array `json:"data" avro:"data"`

	// Database-oriented data
	Id        uint      `json:"-"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// Type SampleTime contains timing information for a single sample
type SampleTime struct {
	// Number of seconds since the UNIX epoch start on January 1st, 1970 at UTC
	Seconds int64 `json:"seconds" avro:"seconds" db:"time_seconds"`

	// Microseconds extension for the UNIX time stamp
	Microseconds int32 `json:"microseconds" avro:"microseconds" db:"time_microseconds"`
}

// Type SampleConfig describes sensor configuration in use when collection a sample
type SampleConfig struct {
	// Antenna gain in dBi
	AntennaGain *float32 `json:"antennaGain,omitempty" avro:"antennaGain" db:"config_antenna_gain"`

	// Identifier for the antenna being used if device has multiple antennas
	AntennaId *int `json:"antennaId,omitempty" avro:"antennaId" db:"config_antenna_id"`

	// Center frequency in Hz to which the RF front-end was tuned to while recording the associated spectrum data
	CenterFreq int64 `json:"centerFreq" avro:"centerFreq" db:"config_center_freq"`

	// Estimated noise floor in dB
	EstNoiseFloor *float32 `json:"estNoiseFloor,omitempty" avro:"estNoiseFloor" db:"config_est_noise_floor"`

	// Correction factor for center frequency in Hz. The correction is already applied to the center frequency (0.0 for no correction)
	FrequencyCorrectionFactor *float32 `json:"frequencyCorrectionFactor,omitempty" avro:"frequencyCorrectionFactor" db:"config_frequency_correction_factor"`

	// RF front-end gain in dB (-1 for automatic gain control)
	FrontendGain *float32 `json:"frontendGain,omitempty" avro:"frontendGain" db:"config_frontend_gain"`

	// Hopping strategy  used to overcome the bandwidth limitations of the RF front-end (0:Sequential, 1:Random, 2:Similarity)
	HoppingStrategy *int `json:"hoppingStrategy,omitempty" avro:"hoppingStrategy" db:"config_hopping_strategy"`

	// True if IQ samples are balanced
	IqBalanceCalibration *bool `json:"iqBalanceCalibration,omitempty" avro:"iqBalanceCalibration" db:"config_iq_balance_calibration"`

	// Time synchronization of the radio frontend (0: none, 1: GPS, 2: Reference Clock, 5: Other)
	RfSync *int32 `json:"rfSync,omitempty" avro:"rfSync" db:"config_rf_sync"`

	// Sensor's sampling rate in samples per second
	SamplingRate *int `json:"samplingRate,omitempty" avro:"samplingRate" db:"config_sampling_rate"`

	// True if signal strength is calibrated
	SigStrengthCalibration *bool `json:"sigStrengthCalibration,omitempty" avro:"sigStrengthCalibration" db:"config_sig_strength_calibration"`

	// Time synchronization of the system (0: none, 1: GPS, 2: Reference Clock, 3: NTP, 4: OpenSky, 5: Other)
	SystemSync *string `json:"systemSync,omitempty" avro:"systemSync" db:"config_system_sync"`

	// Extra configuration for arbitrary data
	ExtraConf map[string]interface{} `json:"extraConf,omitempty" avro:"extraConf" db:"config_extra_conf"`
}
