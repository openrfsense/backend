package stream

import (
	"sync"
	"time"

	"github.com/reugn/go-streams"
)

type DiscardingWindow struct {
	slidingInterval    time.Duration
	lastTimestamp      int64
	timestampExtractor func(any) int64
	in                 chan any
	out                chan any
	done               chan struct{}
	sync.Mutex
}

// Verify DiscardingWindow satisfies the Flow interface.
var _ streams.Flow = (*DiscardingWindow)(nil)

// NewDiscardingWindow returns a new processing time based DiscardingWindow.
// Processing time refers to the system time of the machine that is executing the respective operation.
//
// size is the Duration of generated windows.
// slide is the sliding interval of generated windows.
func NewDiscardingWindow(slide time.Duration) *DiscardingWindow {
	return NewDiscardingWindowWithTSExtractor(slide, nil)
}

// NewDiscardingWindowWithTSExtractor returns a new event time based DiscardingWindow.
// Event time is the time that each individual event occurred on its producing device.
// Gives correct results on out-of-order events, late events, or on replays of data.
//
// slide is the sliding interval of generated windows.
// timestampExtractor is the record timestamp (in nanoseconds) extractor.
func NewDiscardingWindowWithTSExtractor(slide time.Duration, timestampExtractor func(any) int64) *DiscardingWindow {
	window := &DiscardingWindow{
		slidingInterval:    slide,
		timestampExtractor: timestampExtractor,
		in:                 make(chan any),
		out:                make(chan any),
	}

	go window.receive()
	return window
}

// Via streams data through the given flow
func (dw *DiscardingWindow) Via(flow streams.Flow) streams.Flow {
	go dw.transmit(flow)
	return flow
}

// To streams data to the given sink
func (dw *DiscardingWindow) To(sink streams.Sink) {
	dw.transmit(sink)
}

// Out returns an output channel for sending data
func (dw *DiscardingWindow) Out() <-chan any {
	return dw.out
}

// In returns an input channel for receiving data
func (dw *DiscardingWindow) In() chan<- any {
	return dw.in
}

// transmit submits newly created windows to the next Inlet.
func (dw *DiscardingWindow) transmit(inlet streams.Inlet) {
	for elem := range dw.Out() {
		inlet.In() <- elem
	}
	close(inlet.In())
}

// timestamp extracts the timestamp from a record if the timestampExtractor is set.
// Returns system clock time otherwise.
func (dw *DiscardingWindow) timestamp(elem any) int64 {
	if dw.timestampExtractor == nil {
		return time.Now().UTC().UnixNano()
	}
	return dw.timestampExtractor(elem)
}

func (dw *DiscardingWindow) receive() {
	for elem := range dw.in {
		elemTs := dw.timestamp(elem)
		if (elemTs - dw.lastTimestamp) > int64(dw.slidingInterval) {
			dw.lastTimestamp = elemTs
			dw.out <- elem
		}
	}
	close(dw.done)
	close(dw.out)
}
