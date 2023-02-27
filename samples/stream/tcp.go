package stream

import (
	"bufio"
	"context"
	"encoding/binary"
	"net"
	"time"

	"github.com/reugn/go-streams"
)

// TCPSource represents an inbound network socket connector.
type TCPSource struct {
	ctx      context.Context
	listener net.Listener
	out      chan any
}

// NewTCPSource returns a new instance of TCPSource.
func NewTCPSource(ctx context.Context, address string) (*TCPSource, error) {
	var err error
	var listener net.Listener
	out := make(chan any)

	listener, err = net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	source := &TCPSource{
		ctx:      ctx,
		listener: listener,
		out:      out,
	}
	go source.acceptConnections()
	go source.listenCtx()

	return source, nil
}

func (ts *TCPSource) listenCtx() {
	<-ts.ctx.Done()

	close(ts.out)

	if ts.listener != nil {
		ts.listener.Close()
	}
}

// acceptConnections accepts new TCP connectiots.
func (ts *TCPSource) acceptConnections() {
	for {
		// accept a new connection
		conn, err := ts.listener.Accept()
		if err != nil {
			return
		}

		// handle the new connection
		go handleConnection(conn, ts.out)
	}
}

// handleConnection handles new connectiots.
func handleConnection(conn net.Conn, out chan<- any) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	dataSize := make([]byte, 4)
	_, err := reader.Read(dataSize)
	if err != nil {
		return
	}

	bufferBytes := make([]byte, binary.BigEndian.Uint32(dataSize))
	_, err = reader.Read(bufferBytes)
	if err == nil && len(bufferBytes) > 0 {
		out <- bufferBytes
	}
}

// Streams data through the given flow
func (ts *TCPSource) Via(_flow streams.Flow) streams.Flow {
	go doStream(ts, _flow)
	return _flow
}

// Returns an output channel for sending data
func (ts *TCPSource) Out() <-chan any {
	return ts.out
}

// TCPSink represents an outbound network socket connector.
type TCPSink struct {
	conn net.Conn
	in   chan any
}

// NewTCPSink returns a new instance of TCPSink with an optional timeout (defaults to 10 seconds).
func NewTCPSink(address string, timeout ...time.Duration) (*TCPSink, error) {
	var err error
	var conn net.Conn

	tOut := 10 * time.Second
	if len(timeout) > 0 {
		tOut = timeout[0]
	}

	conn, err = net.DialTimeout("tcp", address, tOut)
	if err != nil {
		return nil, err
	}

	sink := &TCPSink{
		conn: conn,
		in:   make(chan any),
	}

	go sink.init()
	return sink, nil
}

// init starts the main loop
func (ts *TCPSink) init() {
	writer := bufio.NewWriter(ts.conn)

	for msg := range ts.in {
		switch m := msg.(type) {
		case string:
			_, err := writer.WriteString(m)
			if err == nil {
				writer.Flush()
			}
		case []byte:
			_, err := writer.Write(m)
			if err == nil {
				writer.Flush()
			}
		}
	}

	ts.conn.Close()
}

// In returns an input channel for receiving data
func (ts *TCPSink) In() chan<- any {
	return ts.in
}
