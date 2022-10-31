package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/openrfsense/backend/nats"
	"github.com/openrfsense/common/stats"
)

func fetchAllSensorStats() ([]stats.Stats, error) {
	statsAll, err := nats.Ping[stats.Stats]("node.all", "node.get.all", nats.PingConfig{
		Timeout: 100 * time.Millisecond,
	})
	if err != nil {
		return nil, err
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
