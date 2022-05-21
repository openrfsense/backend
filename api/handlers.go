package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/openrfsense/backend/nats"
	"github.com/openrfsense/common/keystore"
	"github.com/openrfsense/common/stats"
	"github.com/openrfsense/common/types"
)

// @summary      Request an Emitter key
// @description  Returns an [Emitter channel key](https://emitter.io/develop/getting-started/) for a specific channel and access mode.
// @tags         internal
// @security     BasicAuth
// @accept       json
// @produce      plain
// @param        message  body      types.KeyRequest  true  "Channel name and access modes string"
// @success      200      {string}  string            "A valid private key for the requested channel"
// @failure      422      "When the request JSON is malformed"
// @failure      500      "When the kestore cannot retrieve the key"
// @router       /key [post]
func KeyPost(ctx *fiber.Ctx) error {
	keyReq := new(types.KeyRequest)
	if err := ctx.BodyParser(keyReq); err != nil {
		return err
	}

	// TODO: move validation elsewhere
	if keyReq.Access == "" || keyReq.Channel == "" {
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	key, err := keystore.Must(keyReq.Channel, keyReq.Access)
	if err != nil {
		return err
	}

	return ctx.SendString(key)
}

// @summary      List nodes
// @description  Returns a list of all connected nodes by their hardware ID. Will time out in 300ms if any one of the nodes does not respond.
// @tags         nodes
// @security     BasicAuth
// @produce      json
// @success      200  {array}  stats.Stats  "Bare statistics for all the running and connected nodes"
// @failure      500  "When the internal timeout for information retrieval expires"
// @router       /nodes [get]
func ListGet(ctx *fiber.Ctx) error {
	nodes, err := nats.Presence("node.all")
	if err != nil {
		return err
	}

	// Listen for responses on node.get.all (arbitrary)
	statsChan := make(chan stats.Stats)
	sub, err := nats.Conn().Subscribe("node.get.all", func(s *stats.Stats) {
		statsChan <- *s
	})
	defer sub.Unsubscribe()

	// Request stats with timeout
	err = nats.Conn().PublishRequest("node.all", "node.get.all", "")
	if err != nil {
		return err
	}
	err = nats.Conn().FlushTimeout(300 * time.Millisecond)
	if err != nil {
		return err
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
			log.Error("ListGet timed out")
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	return ctx.JSON(statsAll)
}

// @summary      Get stats from a node
// @description  Returns full stats from the node with given hardware ID. Will time out in 300ms the node does not respond.
// @tags         nodes
// @security     BasicAuth
// @produce      json
// @success      200  {object}  stats.Stats  "Full system statistics for the node associated to the given ID"
// @failure      500  "When the internal timeout for information retrieval expires"
// @router       /nodes/{id}/stats [get]
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
