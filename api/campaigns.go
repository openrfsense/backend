package api

import (
	"context"
	"strings"

	"github.com/openrfsense/backend/database"
	"github.com/openrfsense/backend/database/models"

	"github.com/gofiber/fiber/v2"
)

// List campaigns
//
// @summary     List campaigns
// @description Returns a list of campaigns that were successfully started. Will return all campaigns unless either of the query parameters is set.
// @tags        data
// @security    BasicAuth
// @param       sensors    path string false "Matches campigns which contain ALL these sensors as a comma-separated list."
// @param       campaignId path string false "Matches a single campaign by its unique ID."
// @produce     json
// @success     200 {array} models.Campaign "All recorded campaigns which match the given parameters"
// @failure     500 "Generally a database error"
// @router      /campaigns [get]
func CampaignsGet(ctx *fiber.Ctx) error {
	campaignId := ctx.Query("campaignId")
	sensors := ctx.Query("sensors")

	builder := database.Instance().Select("*").From("campaigns")
	if len(campaignId) > 0 {
		builder = builder.Where("campaign_id = ?", campaignId)
	}

	if len(sensors) > 0 {
		builder = builder.Where("sensors @> ?", strings.Split(sensors, ","))
	}

	sql, args, _ := builder.ToSql()
	campaigns, err := database.Multiple[models.Campaign](
		context.Background(),
		sql,
		args...,
	)
	if len(campaigns) == 0 {
		return ctx.JSON(campaigns)
	}
	if err != nil {
		return err
	}

	return ctx.JSON(campaigns)
}
