package nats

import (
	"errors"
	"time"

	"github.com/knadh/koanf"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"

	"github.com/openrfsense/common/logging"
)

var (
	natsConn   *nats.EncodedConn
	natsServer *server.Server

	log = logging.New().
		WithPrefix("nats").
		WithLevel(logging.DebugLevel).
		WithFlags(logging.FlagsDevelopment)
)

var ErrNotReady = errors.New("server still isn't ready for connections")

// Start the embedded NATS server. If options are passed as parameters, they will override the internal
// options (the common config module is used).
func Start(config *koanf.Koanf, options ...server.Options) error {
	token := config.MustString("nats.token")
	opts := server.Options{
		Host:          config.String("backend.host"),
		Port:          config.MustInt("nats.port"),
		JetStream:     false,
		Authorization: token,
		Debug:         true,
	}
	if len(options) > 0 {
		opts = options[0]
	}

	var err error
	natsServer, err = startServer(&opts)
	if err != nil {
		return err
	}

	// Create a default client connected to the embedded server
	natsConn, err = startClient(natsServer.ClientURL(), token, nats.JSON_ENCODER)

	return err
}

// Returns the number of subscribers on a given subject.
func Presence(subject string) (int, error) {
	subsz, err := natsServer.Subsz(&server.SubszOptions{
		Subscriptions: true,
		Test:          subject,
	})
	if err != nil {
		return 0, err
	}

	return len(subsz.Subs), nil
}

// Returns the default internal client.
func Conn() *nats.EncodedConn {
	return natsConn
}

// Drains and closes the connection to the embedded NATS server.
func Disconnect() {
	_ = natsConn.Drain()
	natsConn.Close()
}

// Creates server object and starts it in a goroutine.
func startServer(opts *server.Options) (*server.Server, error) {
	s, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}

	go func() {
		s.Start()
		s.WaitForShutdown()
	}()

	if !s.ReadyForConnections(100 * time.Millisecond) {
		return nil, ErrNotReady
	}

	return s, nil
}
