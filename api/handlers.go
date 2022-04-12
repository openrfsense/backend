package api

import (
	"github.com/gofiber/fiber/v2"

	"github.com/openrfsense/backend/mqtt"
	"github.com/openrfsense/common/keystore"
	"github.com/openrfsense/common/types"
)

// @summary      Request an Emitter key
// @description  Returns an [Emitter channel key](https://emitter.io/develop/getting-started/) for a specific channel and access mode
// @security     BasicAuth
// @accept       json
// @produce      plain
// @param        message  body      types.KeyRequest  true  "Channel name and access modes string"
// @success      200      {string}  string            "A valid private key for the requested channel"
// @failure      422
// @failure      500
// @router       /key [post]
func KeyPost(ctx *fiber.Ctx) error {
	keyReq := new(types.KeyRequest)
	if err := ctx.BodyParser(keyReq); err != nil {
		return err
	}

	key, err := keystore.Must(keyReq.Channel, keyReq.Access)
	if err != nil {
		return err
	}

	return ctx.SendString(key)
}

func ListGet(ctx *fiber.Ctx) error {
	// TODO: implement this + swagger doc
	err := mqtt.Presence("sensors/all/", true, false)
	if err != nil {
		return err
	}

	// err = mqtt.Publish("control/", "es_sensor")
	// if err != nil {
	// 	return err
	// }

	return ctx.SendString("")
}
