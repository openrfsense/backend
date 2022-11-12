package collector

import (
	"net"
	"testing"

	"github.com/hamba/avro/v2"
	"github.com/openrfsense/backend/database"
)

var avroPacket = database.Sample{
	SensorID:   "sensor",
	CampaignId: "campaign",
	SampleType: "PSD",
	Config: database.SampleConfig{
		CenterFreq: 0,
	},
	Time: database.SampleTime{
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
}
