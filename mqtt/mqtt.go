package mqtt

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	emitter "github.com/emitter-io/go/v2"

	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/keystore"
	"github.com/openrfsense/common/logging"
)

// A default time-to-live for keys, useful for keystore initialization and
// calls to emitter.Client.GenerateKey.
var DefaultTTL = 600 * time.Second

var (
	client *emitter.Client

	log = logging.New().
		WithPrefix("mqtt").
		WithLevel(logging.DebugLevel).
		WithFlags(logging.FlagsDevelopment)
)

// Type Payload represents a payload which can be sent over MQTT.
// Since emitter rejects anything which isn't a string or byte array,
// using generics ensures the handlers always return a payload which
// can be sent. In general, message sending has to be ensured.
type Payload interface {
	~string | ~[]byte
}

// Starts the internal MQTT client (and connects to the broker).
func Init() {
	brokerHost := fmt.Sprintf(
		"%s:%d",
		config.Must[string]("mqtt.host"),
		config.GetWeakInt("mqtt.port"),
	)
	brokerUrl := url.URL{
		Scheme: config.Get[string]("mqtt.protocol"),
		Host:   brokerHost,
	}

	client = emitter.NewClient(
		emitter.WithBrokers(brokerUrl.String()),
		emitter.WithAutoReconnect(true),
		emitter.WithConnectTimeout(10*time.Second),
		emitter.WithKeepAlive(10*time.Second),
		emitter.WithMaxReconnectInterval(2*time.Minute),
	)

	err := client.Connect()
	if err != nil {
		count := 0
		ticker := time.NewTicker(30 * time.Second)
		for !client.IsConnected() && count < 5 {
			log.Warn("Could not connect to MQTT broker, trying again")
			<-ticker.C
			client.Connect()
			count++
		}
		ticker.Stop()
	}

	client.OnConnect(func(_ *emitter.Client) {
		log.Info("Connected to MQTT broker")
	})
}

// Sends a GET request (empty payload) to node/get/channel/ and handles the
// response at node/channel/ with the given handler.
func Get(channel string, handler emitter.MessageHandler) error {
	channelReq := fmt.Sprintf("node/get/%s/", strings.Trim(channel, "/"))
	channelResp := fmt.Sprintf("node/%s/", strings.Trim(channel, "/"))
	keyReq, err := keystore.Must(channelReq, "w")
	if err != nil {
		return err
	}

	keyResp, err := keystore.Must(channelResp, "r")
	if err != nil {
		return err
	}

	err = client.Subscribe(keyResp, channelResp, func(c *emitter.Client, m emitter.Message) {
		handler(c, m)
		c.Unsubscribe(keyResp, channelResp)
	})
	if err != nil {
		return err
	}

	return client.Publish(keyReq, channelReq, []byte{}, emitter.WithoutEcho())
}

// Sends a POST request (with a payload) to node/get/channel/ and handles the
// response at node/channel/ with the given handler.
func Post[P Payload](channel string, payload P, handler emitter.MessageHandler) error {
	channelReq := fmt.Sprintf("node/post/%s/", strings.Trim(channel, "/"))
	channelResp := fmt.Sprintf("node/%s/", strings.Trim(channel, "/"))
	keyReq, err := keystore.Must(channelReq, "w")
	if err != nil {
		return err
	}

	keyResp, err := keystore.Must(channelResp, "r")
	if err != nil {
		return err
	}

	err = client.Subscribe(keyResp, channelResp, func(c *emitter.Client, m emitter.Message) {
		handler(c, m)
		c.Unsubscribe(keyResp, channelResp)
	})
	if err != nil {
		return err
	}

	return client.Publish(keyReq, channelReq, payload, emitter.WithoutEcho())
}

// Custom keystore.Retriever which fetches channel keys from the Emitter broker
// using the main secret key.
func NewBrokerRetriever() keystore.Retriever {
	secret := config.Get[string]("mqtt.secret")

	return func(channel string, access string) (string, error) {
		log.Debugf("asking broker for channel '%s' and access '%s'", channel, access)
		key, err := client.GenerateKey(secret, channel, access, int(DefaultTTL/time.Second))
		if err != nil {
			return "", err
		}
		return key, nil
	}
}
