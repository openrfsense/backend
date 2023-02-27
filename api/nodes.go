package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/openrfsense/backend/nats"
	"github.com/openrfsense/common/stats"
)

// List nodes
//
// @summary     List nodes
// @description Returns a list of all connected nodes by their hardware ID. Will time out in `300ms` if any one of the nodes does not respond.
// @tags        administration
// @security    BasicAuth
// @produce     json
// @success     200 {array} stats.Stats "Bare statistics for all the running and connected nodes"
// @failure     500 "When the internal timeout for information retrieval expires"
// @router      /nodes [get]
func NodesGet(ctx *fiber.Ctx) error {
	statsAll, err := nats.Ping[stats.Stats]("node.all", "node.get.all", nats.PingConfig{
		Timeout: 100 * time.Millisecond,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(statsAll)
}

// Get stats from a node
//
// @summary     Get stats from a node
// @description Returns full stats from the node with given hardware ID. Will time out in `300ms` if the node does not respond.
// @tags        administration
// @security    BasicAuth
// @param       sensor_id path string true "Node hardware ID"
// @produce     json
// @success     200 {object} stats.Stats "Full system statistics for the node associated to the given ID"
// @failure     500 "When the internal timeout for information retrieval expires"
// @router      /nodes/{sensor_id} [get]
func NodeGet(ctx *fiber.Ctx) error {
	id := ctx.Params("sensor_id")
	if id == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	stat := stats.Stats{}
	channel := fmt.Sprintf("node.%s.stats", strings.Trim(id, "."))
	err := nats.Conn().Request(channel, "", &stat, 300*time.Millisecond)
	if err != nil {
		return err
	}

	return ctx.JSON(stat)
}
