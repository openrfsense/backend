package collector

import (
	"net"
	"testing"

	"github.com/linkedin/goavro/v2"
)

var avroPacket = map[string]interface{}{
	"sensorId":   "sensor",
	"campaignId": "campaign",
	"sampleType": "PSD",
	"config": map[string]interface{}{
		"hoppingStrategy":           nil,
		"antennaGain":               nil,
		"frontendGain":              nil,
		"samplingRate":              nil,
		"centerFreq":                0,
		"frequencyCorrectionFactor": nil,
		"antennaId":                 nil,
		"rfSync":                    nil,
		"systemSync":                nil,
		"sigStrengthCalibration":    nil,
		"iqBalanceCalibration":      nil,
		"estNoiseFloor":             nil,
		"extraConf":                 nil,
	},
	"time": map[string]interface{}{
		"seconds":      0,
		"microseconds": 0,
	},
	"data": []float64{
		0.1,
	},
}

func TestHandleRequest(t *testing.T) {
	codecBytes, err := schemasFs.ReadFile("sample.avsc")
	if err != nil {
		t.Fatal(err)
	}
	codec, err = goavro.NewCodec(string(codecBytes))
	if err != nil {
		t.Fatal(err)
	}

	bin, err := codec.BinaryFromNative(nil, map[string]interface{}(avroPacket))
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.Dial("tcp", ":2222")
	if err != nil {
		t.Error("Could not connect to server (wrong port?): ", err)
	}
	defer conn.Close()

	if _, err := conn.Write(bin); err != nil {
		t.Error("could not write payload to TCP server:", err)
	}
}
