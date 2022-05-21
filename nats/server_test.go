package nats

import (
	"testing"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

const token = "super_secure_token"

func createClientConnSubscribeAndPublish(t *testing.T, s *server.Server, subjects ...string) *nats.Conn {
	nc, err := nats.Connect(
		s.ClientURL(),
		nats.Token(token),
	)
	if err != nil {
		t.Fatalf("Error creating client: %v to: %s\n", err, s.ClientURL())
	}

	ch := make(chan bool)
	inbox := nats.NewInbox()
	sub, err := nc.Subscribe(inbox, func(m *nats.Msg) { ch <- true })
	if err != nil {
		t.Fatalf("Error subscribing to `%s`: %v\n", inbox, err)
	}
	nc.Publish(inbox, []byte("Hello"))

	for _, s := range subjects {
		nc.Publish(s, []byte("Hello"))
	}

	if nc.LastError() != nil {
		t.Fail()
	}

	// Wait for message
	<-ch
	sub.Unsubscribe()
	close(ch)
	nc.Flush()
	return nc
}

func TestServer(t *testing.T) {
	s, err := startServer(&server.Options{
		Host:          "",
		Port:          0,
		JetStream:     false,
		Authorization: token,
	})
	if err != nil {
		t.Fail()
	}

	c := createClientConnSubscribeAndPublish(t, s, "node.all")
	t.Cleanup(c.Close)
}
