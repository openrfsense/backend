package mqtt

import (
	"fmt"
	"log"
	"net/url"
	"time"

	emitter "github.com/emitter-io/go/v2"
	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/keystore"
)

var Client *emitter.Client

var DefaultTTL = 600 * time.Second

// TODO: make a better init procedure and/or move to openrfsense-common
func InitClient() error {
	brokerHost := fmt.Sprintf("%s:%d", config.Must[string]("mqtt.host"), config.Must[int]("mqtt.port"))
	brokerUrl := url.URL{
		Scheme: config.Get[string]("mqtt.protocol"),
		Host:   brokerHost,
	}
	Client = emitter.NewClient(
		emitter.WithBrokers(brokerUrl.String()),
		emitter.WithAutoReconnect(true),
		emitter.WithConnectTimeout(10*time.Second),
		emitter.WithKeepAlive(10*time.Second),
		emitter.WithMaxReconnectInterval(2*time.Minute),
	)

	Client.OnMessage(func(_ *emitter.Client, msg emitter.Message) {
		fmt.Printf("[emitter] -> [B] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})

	err := Client.Connect()
	if err != nil {
		return err
	}

	Client.OnConnect(func(_ *emitter.Client) {
		log.Println("Regained connection to MQTT broker")
	})
	// Client.OnDisconnect(nil)

	return nil
}

// Custom keystore.Retriever which fetches channel keys from the Emitter broker
// using the main secret key
func NewBrokerRetriever() keystore.Retriever {
	secret := config.Get[string]("mqtt.secret")

	return func(channel string, access string) (string, error) {
		log.Printf("asking broker for channel '%s' and access '%s'", channel, access)
		key, err := Client.GenerateKey(secret, channel, access, int(DefaultTTL/time.Second))
		if err != nil {
			return "", err
		}
		return key, nil
	}
}

// Disconnect will end the connection with the server, but not before waiting the specified
// time for existing work to be completed.
func Disconnect(waitTime time.Duration) {
	Client.Disconnect(waitTime)
}

// Wrapper around emitter.Presence with automatic key management.
// Presence sends a presence request to the broker.
func Presence(channel string, status, changes bool) error {
	key, err := keystore.Must(channel, "p")
	if err != nil {
		return err
	}
	return Client.Presence(key, channel, status, changes)
}

// Wrapper around emitter.Publish with automatic key management.
// Publish will publish a message with the specified QoS and content to the specified topic.
// Returns a token to track delivery of the message to the broker
func Publish(channel string, payload interface{}, options ...emitter.Option) error {
	key, err := keystore.Must(channel, "w")
	if err != nil {
		return err
	}
	return Client.Publish(key, channel, payload, options...)
}

// Wrapper around emitter.Subscribe with automatic key management.
// Subscribe starts a new subscription. Provide a MessageHandler to be executed when a
// message is published on the topic provided.
func Subscribe(channel string, optionalHandler emitter.MessageHandler, options ...emitter.Option) error {
	key, err := keystore.Must(channel, "r")
	if err != nil {
		return err
	}
	return Client.Subscribe(key, channel, optionalHandler, options...)
}
