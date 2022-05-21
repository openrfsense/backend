package nats

import (
	"errors"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	natsgo "github.com/nats-io/nats.go"

	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/logging"
)

var (
	natsConn   *natsgo.EncodedConn
	natsServer *server.Server

	log = logging.New().
		WithPrefix("nats").
		WithLevel(logging.DebugLevel).
		WithFlags(logging.FlagsDevelopment)
)

var ErrNotReady = errors.New("server still isn't ready for connections")

// Start the embedded NATS server. If options are passed as parameters, they will override the internal
// options (the common config module is used).
func Start(options ...server.Options) error {
	token := config.Must[string]("nats.token")
	opts := server.Options{
		Host:          config.GetOrDefault("backend.host", ""),
		Port:          config.GetOrDefault("nats.port", 0),
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
	c, err := natsgo.Connect(
		natsServer.ClientURL(),
		nats.Token(token),
	)
	if err != nil {
		return err
	}

	natsConn, err = natsgo.NewEncodedConn(c, natsgo.JSON_ENCODER)
	return err
}

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
func Conn() *natsgo.EncodedConn {
	return natsConn
}

// Drain and close the internal NATS connection.
func Disconnect() {
	natsConn.Drain()
	natsConn.Close()
}
