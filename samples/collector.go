package samples

import (
	"context"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/hamba/avro/v2"
	"github.com/knadh/koanf"
	"github.com/openrfsense/backend/database/models"
	"github.com/openrfsense/backend/samples/stream"
	"github.com/reugn/go-streams/flow"
)

func StartCollector(ctx context.Context, config *koanf.Koanf) error {
	// Define a channel-based TCP listener
	addr := fmt.Sprintf(":%d", config.MustInt("collector.port"))
	source, err := stream.NewTCPSource(ctx, addr)
	if err != nil {
		return err
	}

	// Define a BadgerDB sink for data with a custom bucket key extractor
	opts := badger.DefaultOptions(config.MustString("backend.storage"))
	opts.SyncWrites = true
	sink, err := stream.NewBadgerSinkWithPrefixExtractor(opts, extractPrefix(DefaultSchema))
	if err != nil {
		return err
	}

	go source.
		Via(flow.NewPassThrough()).
		To(sink)

	return nil
}

// func (c *Collector) Start() {
// 	c.source.
// 		Via(flow.NewPassThrough()).
// 		To(c.sink)
// }

func extractPrefix(schema avro.Schema) stream.BadgerPrefixExtractor {
	return func(b []byte) []byte {
		var s models.Sample
		err := avro.Unmarshal(schema, b, &s)
		if err != nil {
			log.Error(err)
			return nil
		}
		return makePrefix(s.CampaignId, s.SensorId)
	}
}
