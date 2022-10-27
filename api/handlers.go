package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/openrfsense/backend/nats"
	"github.com/openrfsense/common/id"
	"github.com/openrfsense/common/stats"
	"github.com/openrfsense/common/types"
)

// List nodes
//
// @summary     List nodes
// @description Returns a list of all connected nodes by their hardware ID. Will time out in 300ms if any one of the nodes does not respond.
// @tags        administration
// @security    BasicAuth
// @produce     json
// @success     200 {array} stats.Stats "Bare statistics for all the running and connected nodes"
// @failure     500 "When the internal timeout for information retrieval expires"
// @router      /nodes [get]
func ListGet(ctx *fiber.Ctx) error {
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
// @param       id path string true "Node hardware ID"
// @produce     json
// @success     200 {object} stats.Stats "Full system statistics for the node associated to the given ID"
// @failure     500 "When the internal timeout for information retrieval expires"
// @router      /nodes/{id}/stats [get]
func NodeStatsGet(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
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

// Starts a measurement on a list of nodes and returns an aggregated spectrum measurement
//
// @summary     Get an aggregated spectrum measurement from a list of nodes
// @description Sends an aggregated measurement request to the nodes specified in `sensors` and returns a list of `stats.Stats` objects for all sensors taking part in the campaign. Will time out in `300ms` if any sensor does not respond.
// @tags        measurement
// @security    BasicAuth
// @param       id body types.AggregatedMeasurementRequest true "Measurement request object"
// @produce     json
// @success     200 {array} stats.Stats "Bare statistics for all nodes in the measurement campaign. Will always include sensor status information."
// @failure     500 "When the internal timeout for information retrieval expires"
// @router      /aggregated [post]
func NodeAggregatedPost(ctx *fiber.Ctx) error {
	amr := types.AggregatedMeasurementRequest{}
	err := ctx.BodyParser(&amr)
	if err != nil {
		return err
	}

	err = amr.Validate()
	if err != nil {
		return err
	}

	// Overwrite campaign ID
	amr.CampaignId = id.Generate(9)

	statsAll, err := nats.Ping[stats.Stats]("node.all", "node.get.all", nats.PingConfig{
		Message: amr,
		HowMany: len(amr.Sensors),
		Timeout: 100 * time.Millisecond,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(statsAll)
}

// Starts a measurement on a node and returns the raw spectrum measurement
//
// @summary     Get a raw spectrum measurement from a list of nodes
// @description Sends a raw measurement request to the nodes specified in `sensors` and returns a list of `stats.Stats` objects for all sensors taking part in the campaign. Will time out in `300ms` if any sensor does not respond.
// @tags        measurement
// @security    BasicAuth
// @param       id body types.RawMeasurementRequest true "Measurement request object"
// @produce     json
// @success     200 {array} stats.Stats "Bare statistics for all nodes in the measurement campaign. Will always include sensor status information."
// @failure     500 "When the internal timeout for information retrieval expires"
// @router      /raw [post]
func NodeRawPost(ctx *fiber.Ctx) error {
	rmr := types.RawMeasurementRequest{}
	err := ctx.BodyParser(&rmr)
	if err != nil {
		return err
	}

	err = rmr.Validate()
	if err != nil {
		return err
	}

	// Overwrite campaign ID
	rmr.CampaignId = id.Generate(9)

	statsAll, err := nats.Ping[stats.Stats]("node.all", "node.get.all", nats.PingConfig{
		Message: rmr,
		HowMany: len(rmr.Sensors),
		Timeout: 100 * time.Millisecond,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(statsAll)
}
