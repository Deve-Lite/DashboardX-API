package handler

import (
	"net/http"
	"time"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application"
	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	ae "github.com/Deve-Lite/DashboardX-API/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler interface {
	Broadcast(ctx *gin.Context)
}

type eventHandler struct {
	c  *config.Config
	es application.EventService
}

func NewEventHandler(c *config.Config, es application.EventService) EventHandler {
	return &eventHandler{c, es}
}

// EventHandler godoc
//
//	@Summary		Subscribe a client (user) to the events bus
//	@Description	The event bus allows you to receive events along a specific user,
//					necessary to maintain synchronization between multiple clients (use of several devices at the same time)
//	@Tags			Events
//	@Accept			json
//	@Produce		text/event-stream
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/events [get]
func (h *eventHandler) Broadcast(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		return
	}

	eventChannels, channelID := h.es.Subscribe(ctx, userID)

	var cause string

	defer func() {
		h.es.Unsubscribe(ctx, userID, channelID, cause)
	}()

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	age := int(time.Duration(365 * 24 * float32(time.Hour)).Milliseconds())

	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(string(enum.EventChannelCookie), channelID.String(), age, "/", h.c.Server.Domain, true, true)
	ctx.Writer.Flush()

	for {
		select {
		case <-ctx.Request.Context().Done():
			cause = "user disconnected"
			return
		case evt := <-eventChannels.Channels[channelID.String()]:
			ctx.SSEvent(string(evt.Data.Action), evt.String())
			ctx.Writer.Flush()

			if evt.Data.Action == enum.ChannelClosedAction {
				cause = "user logged out"
				return
			}
		}
	}
}

func (h *eventHandler) getUserID(ctx *gin.Context) (uuid.UUID, error) {
	userID, err := uuid.Parse(ctx.MustGet("UserID").(string))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.HTTPError{Message: "could not retrieve user info"})
		return uuid.Nil, err
	}

	return userID, nil
}
