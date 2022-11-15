package collector

import (
	"net"
	"testing"

	"github.com/openrfsense/backend/database/models"

	"github.com/hamba/avro/v2"
)

var avroPacket = models.Sample{
	SensorId:   "sensor",
	CampaignId: "campaign",
	SampleType: "PSD",
	SampleConfig: models.SampleConfig{
		CenterFreq: 0,
	},
	SampleTime: models.SampleTime{
		Seconds:      0,
		Microseconds: 0,
	},
	Data: []float32{0.1},
}

func TestHandleRequest(t *testing.T) {
	schemaBytes, err := schemasFs.ReadFile("sample.avsc")
	if err != nil {
		t.Fatal(err)
	}

	schema := avro.MustParse(string(schemaBytes))
	bin, err := avro.Marshal(schema, avroPacket)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.Dial("tcp", ":2022")
	if err != nil {
		t.Error("Could not connect to server (wrong port?): ", err)
	}
	defer conn.Close()

	if _, err := conn.Write(bin); err != nil {
		t.Error("could not write payload to TCP server:", err)
	}

	select {
	case err = <-errors:
		if err != nil {
			t.Fatal(err)
		}
	default:
	}
}
