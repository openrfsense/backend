package nats

import (
	"github.com/nats-io/nats.go"
)

func startClient(url string, token string, encoder string) (*nats.EncodedConn, error) {
	c, err := nats.Connect(
		natsServer.ClientURL(),
		nats.Token(token),
	)
	if err != nil {
		return nil, err
	}

	conn, err := nats.NewEncodedConn(c, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	_, err = conn.Subscribe("node.all.error", func(data *map[string]interface{}) {
		d := *data
		log.Debugf("received error: %#v", d)
	})
	if err != nil {
		return nil, err
	}

	_, err = conn.Subscribe("node.all.output", func(data *map[string]interface{}) {
		d := *data
		log.Debugf("received output: %#v", d)
	})
	if err != nil {
		return nil, err
	}

	return conn, nil
}
