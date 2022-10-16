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
                    "nodes",
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
        "/nodes/{id}/aggregated": {
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
                    "nodes",
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
        "/nodes/{id}/stats": {
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
                    "nodes",
                    "administration"
                ],
                "summary": "Get stats from a node",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Node hardware ID",
                        "name": "id",
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
        }
    },
    "definitions": {
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
                    "description": "Start time in milliseconds since epoch (Unix time)",
                    "type": "integer"
                },
                "end": {
                    "description": "End time in milliseconds since epoch (Unix time)",
                    "type": "integer"
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
