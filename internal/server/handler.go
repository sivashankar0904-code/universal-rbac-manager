package server

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"urm/internal/apperr"
	"urm/internal/service"
)

// Handler holds business-layer dependencies and the endpoint methods.
// Server stays limited to HTTP transport/lifecycle concerns.
type Handler struct {
	log        *slog.Logger
	serviceSvc *service.Service
}

func NewHandler(
	log *slog.Logger,
	serviceSvc *service.Service,

) *Handler {
	return &Handler{
		log:        log,
		serviceSvc: serviceSvc,
	}
}

func (h *Handler) respondJSON(c *gin.Context, status int, payload any) {
	c.JSON(status, gin.H{"data": payload, "status": status})
}

func (h *Handler) respondError(c *gin.Context, err error) {
	var ae *apperr.Error
	if errors.As(err, &ae) {
		if ae.Kind() == apperr.KindInternal {
			h.log.Error("internal error", "err", ae)
		} else {
			h.log.Info("request failed", "err", ae)
		}

		status := ae.HTTPStatusCode()
		body := gin.H{"data": gin.H{}, "status": status, "message": ae.Message(), "code": string(ae.Kind())}
		if fields := ae.Fields(); len(fields) > 0 {
			body["fields"] = fields
		}
		c.JSON(status, body)
		return
	}

	h.log.Error("unhandled error", "err", err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"data": gin.H{}, "status": http.StatusInternalServerError, "message": "internal error",
	})
}
