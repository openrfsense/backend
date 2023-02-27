package samples

import (
	"context"
	"time"

	"github.com/openrfsense/backend/database/models"
)

// TODO: force pagination?
func RetrieveSamples(ctx context.Context, campaignId string, sensorId string, from time.Time, to time.Time) ([]models.Sample, error) {
	samples := []models.Sample{}
	compareFrom := !from.IsZero()
	compareTo := !to.IsZero()

	var fromSampleTime, toSampleTime models.SampleTime
	if compareFrom {
		fromSampleTime.Seconds = from.Unix()
		fromSampleTime.Microseconds = int32(from.UnixMicro() - (from.Unix() * 10e5))
	}
	if compareTo {
		toSampleTime.Seconds = to.Unix()
		toSampleTime.Microseconds = int32(to.UnixMicro() - (to.Unix() * 10e5))
	}

	// source, err := stream.NewBadgerSource(ctx, badger.DefaultOptions(""), makePrefix(campaignId, sensorId))
	// if err != nil {
	// 	return nil, err
	// }

	return samples, nil
}
