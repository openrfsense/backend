// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "OpenRFSense",
            "url": "https://github.com/openrfsense/backend/issues"
        },
        "license": {
            "name": "AGPLv3",
            "url": "https://spdx.org/licenses/AGPL-3.0-or-later.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/aggregated": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Sends an aggregated measurement request to the nodes specified in ` + "`" + `sensors` + "`" + ` and returns a list of ` + "`" + `stats.Stats` + "`" + ` objects for all sensors taking part in the campaign. Will time out in ` + "`" + `300ms` + "`" + ` if any sensor does not respond.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "measurement"
                ],
                "summary": "Get an aggregated spectrum measurement from a list of nodes",
                "parameters": [
                    {
                        "description": "Measurement request object",
                        "name": "id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.AggregatedMeasurementRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Bare statistics for all nodes in the measurement campaign. Will always include sensor status information.",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/stats.Stats"
                            }
                        }
                    },
                    "500": {
                        "description": "When the internal timeout for information retrieval expires"
                    }
                }
            }
        },
        "/campaigns": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Returns a list of all recorded campaigns (that were successfully started).",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data"
                ],
                "summary": "List campaigns",
                "responses": {
                    "200": {
                        "description": "All recorded campaigns",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/database.Campaign"
                            }
                        }
                    },
                    "500": {
                        "description": "Generally a database error"
                    }
                }
            }
        },
        "/campaigns/{campaign_id}": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Returns the campaign object corresponding to the given unique ID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data"
                ],
                "summary": "Get a single campaign object",
                "responses": {
                    "200": {
                        "description": "The campaign with the given ID",
                        "schema": {
                            "$ref": "#/definitions/database.Campaign"
                        }
                    },
                    "500": {
                        "description": "Generally a database error"
                    }
                }
            }
        },
        "/campaigns/{campaign_id}/samples": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Returns a list of all the samples recorded during a campaign by the sensors partakin in said campaign.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data"
                ],
                "summary": "Get all samples recorded during a specific campaign",
                "responses": {
                    "200": {
                        "description": "All samples received during the campaign",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/database.Sample"
                            }
                        }
                    },
                    "500": {
                        "description": "Generally a database error"
                    }
                }
            }
        },
        "/nodes": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Returns a list of all connected nodes by their hardware ID. Will time out in 300ms if any one of the nodes does not respond.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "administration"
                ],
                "summary": "List nodes",
                "responses": {
                    "200": {
                        "description": "Bare statistics for all the running and connected nodes",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/stats.Stats"
                            }
                        }
                    },
                    "500": {
                        "description": "When the internal timeout for information retrieval expires"
                    }
                }
            }
        },
        "/nodes/{sensor_id}/samples": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Returns all samples received by the backend from the sensor with the given ID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data"
                ],
                "summary": "Get all samples received from a specific node",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Node hardware ID",
                        "name": "sensor_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of samples received by the given sensor",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/database.Sample"
                            }
                        }
                    },
                    "500": {
                        "description": "Generally a database error"
                    }
                }
            }
        },
        "/nodes/{sensor_id}/stats": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Returns full stats from the node with given hardware ID. Will time out in ` + "`" + `300ms` + "`" + ` if the node does not respond.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "administration"
                ],
                "summary": "Get stats from a node",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Node hardware ID",
                        "name": "sensor_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Full system statistics for the node associated to the given ID",
                        "schema": {
                            "$ref": "#/definitions/stats.Stats"
                        }
                    },
                    "500": {
                        "description": "When the internal timeout for information retrieval expires"
                    }
                }
            }
        },
        "/raw": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Sends a raw measurement request to the nodes specified in ` + "`" + `sensors` + "`" + ` and returns a list of ` + "`" + `stats.Stats` + "`" + ` objects for all sensors taking part in the campaign. Will time out in ` + "`" + `300ms` + "`" + ` if any sensor does not respond.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "measurement"
                ],
                "summary": "Get a raw spectrum measurement from a list of nodes",
                "parameters": [
                    {
                        "description": "Measurement request object",
                        "name": "id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.RawMeasurementRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Bare statistics for all nodes in the measurement campaign. Will always include sensor status information.",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/stats.Stats"
                            }
                        }
                    },
                    "500": {
                        "description": "When the internal timeout for information retrieval expires"
                    }
                }
            }
        }
    },
    "definitions": {
        "database.Campaign": {
            "type": "object",
            "properties": {
                "begin": {
                    "description": "The time at which the campaign is supposed to start",
                    "type": "string"
                },
                "campaignId": {
                    "description": "The textual, random ID for the campaign",
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "end": {
                    "description": "The time at which the campaign will end",
                    "type": "string"
                },
                "sensors": {
                    "description": "The list of sensor partaking in the campaign",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "type": {
                    "description": "The type of measurements requested",
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "database.Sample": {
            "type": "object",
            "properties": {
                "campaignId": {
                    "description": "Unique identifier for the campaign this sample belongs to",
                    "type": "string"
                },
                "config": {
                    "description": "Sensor configuration for the recorded data set",
                    "$ref": "#/definitions/database.SampleConfig"
                },
                "createdAt": {
                    "type": "string"
                },
                "data": {
                    "description": "Actual measurement data. Unit depends on measurement type",
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "deletedAt": {
                    "type": "string"
                },
                "lossRate": {
                    "description": "Sample loss rate in percent",
                    "type": "number"
                },
                "obfuscation": {
                    "description": "Method used to obfuscate IQ spectrum data",
                    "type": "string"
                },
                "sampleType": {
                    "description": "Sample type string (IQ, PSD, DEC)",
                    "type": "string"
                },
                "sensorId": {
                    "description": "The unique hardware id of the sensor",
                    "type": "string"
                },
                "time": {
                    "description": "Sample timestamp with microseconds precision",
                    "$ref": "#/definitions/database.SampleTime"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "database.SampleConfig": {
            "type": "object",
            "properties": {
                "antennaGain": {
                    "description": "Antenna gain in dBi",
                    "type": "number"
                },
                "antennaId": {
                    "description": "Identifier for the antenna being used if device has multiple antennas",
                    "type": "string"
                },
                "centerFreq": {
                    "description": "Center frequency in Hz to which the RF front-end was tuned to while recording the associated spectrum data",
                    "type": "integer"
                },
                "estNoiseFloor": {
                    "description": "Estimated noise floor in dB",
                    "type": "number"
                },
                "extraConf": {
                    "description": "Extra configuration for arbitrary data",
                    "type": "object",
                    "additionalProperties": true
                },
                "frequencyCorrectionFactor": {
                    "description": "Correction factor for center frequency in Hz. The correction is already applied to the center frequency (0.0 for no correction)",
                    "type": "number"
                },
                "frontendGain": {
                    "description": "RF front-end gain in dB (-1 for automatic gain control)",
                    "type": "number"
                },
                "hoppingStrategy": {
                    "description": "Hopping strategy  used to overcome the bandwidth limitations of the RF front-end (0:Sequential, 1:Random, 2:Similarity)",
                    "type": "integer"
                },
                "iqBalanceCalibration": {
                    "description": "True if IQ samples are balanced",
                    "type": "boolean"
                },
                "rfSync": {
                    "description": "Time synchronization of the radio frontend (0: none, 1: GPS, 2: Reference Clock, 5: Other)",
                    "type": "string"
                },
                "samplingRate": {
                    "description": "Sensor's sampling rate in samples per second",
                    "type": "integer"
                },
                "sigStrengthCalibration": {
                    "description": "True if signal strength is calibrated",
                    "type": "boolean"
                },
                "systemSync": {
                    "description": "Time synchronization of the system (0: none, 1: GPS, 2: Reference Clock, 3: NTP, 4: OpenSky, 5: Other)",
                    "type": "string"
                }
            }
        },
        "database.SampleTime": {
            "type": "object",
            "properties": {
                "microseconds": {
                    "description": "Microseconds extension for the UNIX time stamp",
                    "type": "integer"
                },
                "seconds": {
                    "description": "Number of seconds since the UNIX epoch start on January 1st, 1970 at UTC",
                    "type": "integer"
                }
            }
        },
        "stats.Stats": {
            "type": "object",
            "properties": {
                "hostname": {
                    "description": "Hostname of the system",
                    "type": "string"
                },
                "id": {
                    "description": "A unique identifier for the node (a hardware-bound ID is recommended)",
                    "type": "string"
                },
                "model": {
                    "description": "The model/vendor of the system's hardware, useful for identification",
                    "type": "string"
                },
                "providers": {
                    "description": "Extra, more in-depth information about the system as dynamically returned by providers.",
                    "type": "object",
                    "additionalProperties": true
                },
                "uptime": {
                    "description": "Uptime of the system",
                    "type": "integer"
                }
            }
        },
        "types.AggregatedMeasurementRequest": {
            "type": "object",
            "properties": {
                "begin": {
                    "description": "Start time in ISO 8601",
                    "type": "string"
                },
                "campaignId": {
                    "description": "Campaign ID. For internal use only, will be ignored if not null",
                    "type": "string"
                },
                "end": {
                    "description": "End time in ISO 8601",
                    "type": "string"
                },
                "freqMax": {
                    "description": "Upper bound for frequency in Hz",
                    "type": "integer"
                },
                "freqMin": {
                    "description": "Lower bound for frequency in Hz",
                    "type": "integer"
                },
                "freqRes": {
                    "description": "Frequency resolution in Hz",
                    "type": "integer"
                },
                "sensors": {
                    "description": "List of sensor hardware IDs to run the measurement campaign on",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "timeRes": {
                    "description": "Time resolution in seconds",
                    "type": "integer"
                }
            }
        },
        "types.RawMeasurementRequest": {
            "type": "object",
            "properties": {
                "begin": {
                    "description": "Start time in ISO 8601",
                    "type": "string"
                },
                "campaignId": {
                    "description": "Campaign ID. For internal use only, will be ignored if not null",
                    "type": "string"
                },
                "end": {
                    "description": "End time in ISO 8601",
                    "type": "string"
                },
                "sensors": {
                    "description": "List of sensor hardware IDs to run the measurement campaign on",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "OpenRFSense backend API",
	Description:      "OpenRFSense backend API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
