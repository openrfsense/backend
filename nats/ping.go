package nats

import (
	"context"
	"fmt"
	"time"
)

type PingConfig struct {
	HowMany int
	Timeout time.Duration
	Message interface{}
}

var defaultConfig = PingConfig{
	HowMany: 0,
	Message: "",
	Timeout: 300 * time.Millisecond,
}

func Ping[T any](subject string, reply string, config ...PingConfig) ([]T, error) {
	cfg := configDefault(config...)

	// If not specified, try to count all subscribers to the given subject
	if cfg.HowMany == 0 {
		nodes, err := Presence(subject)
		if err != nil {
			return nil, err
		}

		// Early return if there are no subcribers to the subject
		if nodes == 0 {
			return []T{}, nil
		}

		cfg.HowMany = nodes
	}

	// Asynchronously collect messages in channel as soon as they are received
	collectorChan := make(chan T)
	sub, err := natsConn.Subscribe(reply, func(t *T) {
		collectorChan <- *t
	})
	defer func() {
		err = sub.Unsubscribe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		return nil, err
	}

	// Send the message but register a flush with the given timeout
	err = natsConn.PublishRequest(subject, reply, cfg.Message)
	if err != nil {
		return nil, err
	}
	err = natsConn.FlushTimeout(cfg.Timeout)
	if err != nil {
		return nil, err
	}

	// Collect received repsonses with timeout
	c, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()
	collector := make([]T, 0, cfg.HowMany)
	// Terminate if either one of these conditions is true:
	// - The number of responses equals the given expected quantity
	// - The timeout expires
	for i := 0; i < cfg.HowMany; i++ {
		select {
		case p := <-collectorChan:
			collector = append(collector, p)
		case <-c.Done():
			log.Debugf("Ping timeout for request %#v on %s -> %s", cfg, subject, reply)
			err = fmt.Errorf("Ping timed out after %v for request on '%s'", cfg.Timeout, subject)
		}
	}

	return collector, err
}

// Helper for setting default values of PingConfig
func configDefault(config ...PingConfig) PingConfig {
	if len(config) == 0 {
		return defaultConfig
	}

	cfg := config[0]
	if cfg.Message == nil {
		cfg.Message = defaultConfig.Message
	}

	if int(cfg.Timeout) <= 0 {
		cfg.Timeout = defaultConfig.Timeout
	}

	return cfg
}
