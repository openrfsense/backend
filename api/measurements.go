package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/openrfsense/backend/database"
	"github.com/openrfsense/backend/nats"
	"github.com/openrfsense/common/id"
	"github.com/openrfsense/common/stats"
	"github.com/openrfsense/common/types"
)

// Starts a measurement on a list of nodes and returns an aggregated spectrum measurement
//
// @summary     Get an aggregated spectrum measurement from a list of nodes
// @description Sends an aggregated measurement request to the nodes specified in `sensors` and returns a list of `stats.Stats` objects for all sensors taking part in the campaign. Will time out in `300ms` if any sensor does not respond.
// @tags        measurement
// @security    BasicAuth
// @accept      json
// @param       id body types.AggregatedMeasurementRequest true "Measurement request object"
// @produce     json
// @success     201 {array}  stats.Stats "Bare statistics for all nodes in the measurement campaign. Will always include sensor status information."
// @header      201 {string} Location    "Location of the new campaign object."
// @failure     500 "When the internal timeout for information retrieval expires"
// @router      /aggregated [post]
func AggregatedPost(ctx *fiber.Ctx) error {
	amr := types.AggregatedMeasurementRequest{}
	err := ctx.BodyParser(&amr)
	if err != nil {
		return err
	}

	err = amr.Validate()
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return err
	}

	// Overwrite campaign ID
	amr.CampaignId = id.Generate(9)

	statsAll, err := nats.Ping[stats.Stats]("node.all.aggregated", "node.get.all.aggregated", nats.PingConfig{
		Message: amr,
		HowMany: len(amr.Sensors),
	})
	if err != nil {
		return err
	}

	err = database.Do(
		ctx.Context(),
		`insert into campaigns ("campaign_id", "sensors", "type", "begin", "end") values ($1, $2, 'PSD', $3, $4)`,
		amr.CampaignId,
		amr.Sensors,
		amr.Begin,
		amr.End,
	)
	if err != nil {
		return err
	}

	ctx.Set("Location", "/campaigns?campaignId="+amr.CampaignId)
	return ctx.Status(fiber.StatusCreated).JSON(statsAll)
}

// Starts a measurement on a node and returns the raw spectrum measurement
//
// @summary     Get a raw spectrum measurement from a list of nodes
// @description Sends a raw measurement request to the nodes specified in `sensors` and returns a list of `stats.Stats` objects for all sensors taking part in the campaign. Will time out in `300ms` if any sensor does not respond.
// @tags        measurement
// @security    BasicAuth
// @accept      json
// @param       id body types.RawMeasurementRequest true "Measurement request object"
// @produce     json
// @success     201 {array}  stats.Stats "Bare statistics for all nodes in the measurement campaign. Will always include sensor status information."
// @header      201 {string} Location    "Location of the new campaign object."
// @failure     500 "When the internal timeout for information retrieval expires"
// @router      /raw [post]
func RawPost(ctx *fiber.Ctx) error {
	rmr := types.RawMeasurementRequest{}
	err := ctx.BodyParser(&rmr)
	if err != nil {
		return err
	}

	err = rmr.Validate()
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return err
	}

	// Overwrite campaign ID
	rmr.CampaignId = id.Generate(9)

	statsAll, err := nats.Ping[stats.Stats]("node.all.raw", "node.get.all.raw", nats.PingConfig{
		Message: rmr,
		HowMany: len(rmr.Sensors),
	})
	if err != nil {
		return err
	}

	err = database.Do(
		ctx.Context(),
		`insert into campaigns ("campaign_id", "sensors", "type", "begin", "end") values ($1, $2, 'IQ', $3, $4)`,
		rmr.CampaignId,
		rmr.Sensors,
		rmr.Begin,
		rmr.End,
	)
	if err != nil {
		return err
	}

	ctx.Set("Location", "/campaigns?campaignId="+rmr.CampaignId)
	return ctx.Status(fiber.StatusCreated).JSON(statsAll)
}
