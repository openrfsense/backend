package samples

import (
	"context"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/hamba/avro/v2"
	"github.com/knadh/koanf"
	"github.com/openrfsense/backend/database/models"
	"github.com/openrfsense/backend/samples/stream"
	"github.com/reugn/go-streams/extension"
)

type waterfallMessage struct {
	Center    int64     `json:"center"`
	Framerate int       `json:"framerate"`
	Gain      float32   `json:"gain"`
	Span      int       `json:"span"`
	S         []float32 `json:"s"`
}

func StartWebsocket(ctx context.Context, config *koanf.Koanf, router *fiber.App) {
	handler := makeHandler(ctx, config.MustString("backend.storage"), DefaultSchema)

	router.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	router.Get("/ws/:sensor_id/:campaign_id", websocket.New(handler))
}

func makeHandler(ctx context.Context, badgerPath string, schema avro.Schema) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		campaignId := c.Params("campaign_id")
		sensorId := c.Params("sensor_id")
		defer c.Close()

		conf := waterfallMessage{}
		err := c.ReadJSON(&conf)
		if err != nil {
			log.Error(err)
			return
		}

		source, err := stream.NewBadgerSource(ctx, badger.DefaultOptions(badgerPath), makePrefix(campaignId, sensorId))
		if err != nil {
			log.Error(err)
			return
		}

		// tick := time.Second / time.Duration(conf.Framerate)
		// dw := stream.NewDiscardingWindowWithTSExtractor(tick, sampleTSExtractor)

		// // toWaterfallMessage := flow.NewMap(func(sample models.Sample) waterfallMessage {
		// // 	return waterfallMessage{
		// // 		Center:    sample.SampleConfig.CenterFreq,
		// // 		Framerate: conf.Framerate,
		// // 		Gain:      *sample.SampleConfig.FrontendGain,
		// // 		S:         sample.Data,
		// // 	}
		// // }, 1)

		// // sink := stream.NewWebSocketSink(ctx, c)
		sink := extension.NewStdoutSink()

		go source.
			Via(stream.NewAvroUnmarshal[models.Sample](schema)).
			// Via(dw).
			// Via(flow.NewThrottler(1, tick, 1, flow.Backpressure)).
			// Via(toWaterfallMessage).
			// Via(stream.NewJsonMarshal[waterfallMessage]()).
			To(sink)

		for {
			err := c.ReadJSON(&conf)
			if err != nil {
				log.Error(err)
				break
			}
			log.Debugf("%#v", conf)
			err = c.WriteJSON(conf)
			if err != nil {
				log.Error(err)
				break
			}
		}
	}
}

func sampleTSExtractor(i interface{}) int64 {
	v := i.(models.Sample)
	return v.SampleTime.Seconds*int64(time.Second) +
		int64(v.SampleTime.Microseconds)*int64(time.Microsecond)
}
