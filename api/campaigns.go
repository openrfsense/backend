package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/openrfsense/backend/database"
)

// List campaigns
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
