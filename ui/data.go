package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openrfsense/backend/nats"
	"github.com/openrfsense/common/stats"
)

func fetchAllSensorStats() ([]stats.Stats, error) {
	nodes, err := nats.Presence("node.all")
	if err != nil {
		return nil, err
	}

	// Listen for responses on node.get.all (arbitrary)
	statsChan := make(chan stats.Stats)
	sub, err := nats.Conn().Subscribe("node.get.all", func(s *stats.Stats) {
		statsChan <- *s
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

	// Request stats with timeout
	err = nats.Conn().PublishRequest("node.all", "node.get.all", "")
	if err != nil {
		return nil, err
	}
	err = nats.Conn().FlushTimeout(300 * time.Millisecond)
	if err != nil {
		return nil, err
	}

	// Collect received stats with timeout
	c, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	statsAll := make([]stats.Stats, 0, nodes)
	for i := 0; i < nodes; i++ {
		select {
		case p := <-statsChan:
			statsAll = append(statsAll, p)
		case <-c.Done():
			log.Error("FetchStats timed out")
		}
	}

	return statsAll, nil
}

func fetchSensorStats(id string) (stats.Stats, error) {
	stat := stats.Stats{}
	channel := fmt.Sprintf("node.%s.stats", strings.Trim(id, "."))
	err := nats.Conn().Request(channel, "", &stat, 300*time.Millisecond)
	if err != nil {
		return stats.Stats{}, err
	}

	return stat, nil
}
