package api

import (
	"github.com/openrfsense/backend/database/models"

	"github.com/gofiber/fiber/v2"
)

// Get samples
//
// @summary     Get samples
// @description Returns a list of all the samples recorded during a campaign by a specific sensors partaking in said campaign.
// @tags        data
// @security    BasicAuth
// @param       sensorId   path string true  "Sensor which the samples belong to"
// @param       campaignId path string true  "Campaign which the samples belong to"
// @param       from       path string false "Samples returned will have been received strictly later than this date (must be in ISO 8601/RFC 3339)"
// @param       to         path string false "Samples returned will have been received strictly before this date (must be in ISO 8601/RFC 3339)"
// @produce     json
// @success     200 {array} models.Sample "All samples which respect the given conditions"
// @failure     500 "Generally a database error"
// @router      /samples [get]
func SamplesGet(ctx *fiber.Ctx) error {
	// sensorId := ctx.Query("sensorId")
	// campaignId := ctx.Query("campaignId")

	// if len(sensorId) == 0 || len(campaignId) == 0 {
	// 	return fiber.NewError(fiber.StatusBadRequest, "sensorId and campaignId must be defined")
	// }

	// var err error
	// var from time.Time
	// if fromStr := ctx.Query("from"); len(fromStr) > 0 {
	// 	from, err = time.Parse(time.RFC3339, fromStr)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// var to time.Time
	// if toStr := ctx.Query("to"); len(toStr) > 0 {
	// 	to, err = time.Parse(time.RFC3339, toStr)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	samples := []models.Sample{}
	// samples, err := samples.RetrieveSamples(campaignId, sensorId, from, to)
	// if err != nil {
	// 	return err
	// }
	// if len(samples) == 0 {
	// 	return ctx.JSON(samples)
	// }

	return ctx.JSON(samples)
}
