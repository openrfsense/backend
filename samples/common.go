package samples

import (
	"embed"

	"github.com/hamba/avro/v2"
	"github.com/openrfsense/common/logging"
)

var log = logging.New().
	WithPrefix("samples").
	WithFlags(logging.FlagsDevelopment).
	WithLevel(logging.DebugLevel)

//go:embed sample.avsc
var schemasFs embed.FS

var DefaultSchema avro.Schema

func init() {
	// Initialize schema
	schemaBytes, err := schemasFs.ReadFile("sample.avsc")
	if err != nil {
		log.Fatal(err)
	}
	DefaultSchema = avro.MustParse(string(schemaBytes))
}

func makePrefix(campaignId string, sensorId string) []byte {
	return append(
		append([]byte(campaignId), byte('_')),
		[]byte(sensorId)...,
	)
}
