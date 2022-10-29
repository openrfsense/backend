package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/openrfsense/backend/database"
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
// @param       sensor_id   path string true "Node hardware ID"
// @produce     json
// @success     200 {object} stats.Stats "Full system statistics for the node associated to the given ID"
// @failure     500 "When the internal timeout for information retrieval expires"
// @router      /nodes/{sensor_id}/stats [get]
func NodeStatsGet(ctx *fiber.Ctx) error {
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

// Get all samples received from a specific node
//
// @summary     Get all samples received from a specific node
// @description Returns all samples received by the backend from the sensor with the given ID.
// @tags        data
// @security    BasicAuth
// @param       sensor_id path string true "Node hardware ID"
// @produce     json
// @success     200 {array} database.Sample "List of samples received by the given sensor"
// @failure     500 "Generally a database error"
// @router      /nodes/{sensor_id}/samples [get]
func NodeSamplesGet(ctx *fiber.Ctx) error {
	id := ctx.Params("sensor_id")
	if id == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	samples := []database.Sample{}
	err := database.Instance().
		Model(&database.Sample{}).
		Where("sensor_id = ?", id).
		Find(&samples).
		Error
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(samples)
}

// Get all campaigns a specific node took part in
//
// @summary     Get all campaigns a specific node took part in
// @description Returns all campaigns where the given sensor was requested to take part in.
// @tags        data
// @security    BasicAuth
// @param       sensor_id path string true "Node hardware ID"
// @produce     json
// @success     200 {array} database.Campaign "List of campaign the sensor took part in"
// @failure     500 "Generally a database error"
// @router      /nodes/{sensor_id}/campaigns [get]
func NodeCampaignsGet(ctx *fiber.Ctx) error {
	id := ctx.Params("sensor_id")
	if id == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	campaigns := []database.Campaign{}
	err := database.Instance().
		Model(&database.Campaign{}).
		Where("? = any (sensors)", id).
		Find(&campaigns).
		Error
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(campaigns)
}

// Get all samples received from a specific sensor and belonging to a specific campaign
//
// @summary     Get all samples received from a specific sensor and belonging to a specific campaign
// @description Returns all samples received by the given sensor and belonging to the given campaign.
// @tags        data
// @security    BasicAuth
// @param       sensor_id path string true "Node hardware ID"
// @param       campaign_id path string true "Campaign ID"
// @produce     json
// @success     200 {array} database.Sample "List of samples received by the given sensor during the given campaign"
// @failure     500 "Generally a database error"
// @router      /nodes/{sensor_id}/campaigns/{campaign_id} [get]
func NodeCampaignSamplesGet(ctx *fiber.Ctx) error {
	sensorId := ctx.Params("sensor_id")
	if sensorId == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	campaignId := ctx.Params("campaign_id")
	if campaignId == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	samples := []database.Sample{}
	err := database.Instance().
		Model(&database.Sample{}).
		Where("sensor_id = ? and campaign_id = ?", sensorId, campaignId).
		Find(&samples).
		Error
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(samples)
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

	campaign := &database.Campaign{
		CampaignId: amr.CampaignId,
		Sensors:    amr.Sensors,
		Type:       "PSD",
		Begin:      amr.Begin,
		End:        amr.End,
	}
	err = database.Instance().Create(campaign).Error
	if err != nil {
		log.Error(err)
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

	campaign := &database.Campaign{
		CampaignId: rmr.CampaignId,
		Sensors:    rmr.Sensors,
		Type:       "IQ",
		Begin:      rmr.Begin,
		End:        rmr.End,
	}
	err = database.Instance().Create(campaign).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(statsAll)
}

// List campaign
//
// @summary     List campaigns
// @description Returns a list of all recorded campaigns (that were successfully started).
// @tags        data
// @security    BasicAuth
// @produce     json
// @success     200 {array} database.Campaign "All recorded campaigns"
// @failure     500 "Generally a database error"
// @router      /campaigns [get]
func CampaignsGet(ctx *fiber.Ctx) error {
	campaigns := []database.Campaign{}
	err := database.Instance().Find(&campaigns).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(campaigns)
}

// Get a single campaign object
//
// @summary     Get a single campaign object
// @description Returns the campaign object corresponding to the given unique ID.
// @tags        data
// @security    BasicAuth
// @produce     json
// @success     200 {object} database.Campaign "The campaign with the given ID"
// @failure     500 "Generally a database error"
// @router      /campaigns/{campaign_id} [get]
func CampaignGet(ctx *fiber.Ctx) error {
	id := ctx.Params("campaign_id")
	if id == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	campaign := database.Campaign{}
	err := database.Instance().
		First(&campaign, "campaign_id = ?", id).
		Error
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(campaign)
}

// Get all samples recorded during a specific campaign
//
// @summary     Get all samples recorded during a specific campaign
// @description Returns a list of all the samples recorded during a campaign by the sensors partaking in said campaign.
// @tags        data
// @security    BasicAuth
// @produce     json
// @success     200 {array} database.Sample "All samples received during the campaign"
// @failure     500 "Generally a database error"
// @router      /campaigns/{campaign_id}/samples [get]
func CampaignSamplesGet(ctx *fiber.Ctx) error {
	id := ctx.Params("campaign_id")
	if id == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	campaign := database.Campaign{}
	err := database.Instance().First(&campaign, "campaign_id = ?", id).Error
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	samples := []database.Sample{}
	err = database.Instance().
		Where("campaign_id = ?", campaign.CampaignId).
		Find(&samples).
		Error
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(samples)
}
