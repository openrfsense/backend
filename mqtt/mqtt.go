package mqtt

import (
	"fmt"
	"net/url"
	"time"

	emitter "github.com/emitter-io/go/v2"

	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/keystore"
	"github.com/openrfsense/common/logging"
)

var (
	Client     *emitter.Client
	DefaultTTL = 600 * time.Second

	log = logging.New(
		logging.WithPrefix("mqtt"),
		logging.WithLevel(logging.DebugLevel),
		logging.WithFlags(logging.FlagsDevelopment),
	)
)

// TODO: make a better init procedure and/or move to openrfsense-common
func Init() {
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
		count := 0
		ticker := time.NewTicker(30 * time.Second)
		for !Client.IsConnected() && count < 5 {
			log.Warn("Could not connect to MQTT broker, trying again")
			<-ticker.C
			Client.Connect()
			count++
		}
		ticker.Stop()
	}

	Client.OnConnect(func(_ *emitter.Client) {
		log.Info("Connected to MQTT broker")
	})
}

// Custom keystore.Retriever which fetches channel keys from the Emitter broker
// using the main secret key
func NewBrokerRetriever() keystore.Retriever {
	secret := config.Get[string]("mqtt.secret")

	return func(channel string, access string) (string, error) {
		log.Debugf("asking broker for channel '%s' and access '%s'", channel, access)
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
