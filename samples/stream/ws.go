package stream

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/websocket/v2"
	"github.com/reugn/go-streams"
	"github.com/reugn/go-streams/flow"
)

// Message represents a WebSocket message container.
// Message types are defined in RFC 6455, section 11.8.
type Message struct {
	MsgType int
	Payload []byte
}

// WebSocketSource represents a Fiber WebSocket source connector.
type WebSocketSource struct {
	ctx        context.Context
	connection *websocket.Conn
	out        chan any
}

// NewWebSocketSource creates and returns a new WebSocketSource using an existing connection.
func NewWebSocketSource(ctx context.Context, connection *websocket.Conn) *WebSocketSource {
	source := &WebSocketSource{
		ctx:        ctx,
		connection: connection,
		out:        make(chan any),
	}

	go source.init()
	return source
}

// init starts the main loop
func (wsock *WebSocketSource) init() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-sigchan:
			break loop

		case <-wsock.ctx.Done():
			break loop

		default:
			t, msg, err := wsock.connection.ReadMessage()
			if err != nil {
				log.Printf("Error reading WebSocket message: %s", err)
				continue
			}
			wsock.out <- Message{
				MsgType: t,
				Payload: msg,
			}

			// exit loop on CloseMessage
			if t == websocket.CloseMessage {
				break loop
			}
		}
	}

	log.Print("Closing WebSocketSource connection")
	close(wsock.out)
	wsock.connection.Close()
}

// Via streams data through the given flow
func (wsock *WebSocketSource) Via(_flow streams.Flow) streams.Flow {
	flow.DoStream(wsock, _flow)
	return _flow
}

// Out returns an output channel for sending data
func (wsock *WebSocketSource) Out() <-chan any {
	return wsock.out
}

// WebSocketSink represents a Fiber WebSocket sink connector.
type WebSocketSink struct {
	ctx        context.Context
	connection *websocket.Conn
	in         chan any
}

// NewWebSocketSink creates and returns a new WebSocketSink using an existing connection.
func NewWebSocketSink(ctx context.Context, connection *websocket.Conn) *WebSocketSink {
	sink := &WebSocketSink{
		ctx:        ctx,
		connection: connection,
		in:         make(chan any),
	}

	go sink.init()
	return sink
}

// init starts the main loop
func (wsock *WebSocketSink) init() {
	for msg := range wsock.in {
		var err error
		switch m := msg.(type) {
		case Message:
			err = wsock.connection.WriteMessage(m.MsgType, m.Payload)

		case string:
			err = wsock.connection.WriteMessage(websocket.TextMessage, []byte(m))

		case []byte:
			err = wsock.connection.WriteMessage(websocket.BinaryMessage, m)

		default:
			log.Printf("WebSocketSink Unsupported message type %v", m)
		}

		if err != nil {
			log.Printf("Error processing WebSocket message: %s", err)
		}
	}

	log.Print("Closing WebSocketSink connection")
	wsock.connection.Close()
}

// In returns an input channel for receiving data
func (wsock *WebSocketSink) In() chan<- any {
	return wsock.in
}
