package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"urm/internal/apperr"
	"urm/internal/model"
)

func (h *Handler) CreateOrganisation(c *gin.Context) {
	var o model.Organisation
	if err := c.ShouldBindJSON(&o); err != nil {
		h.respondError(c, apperr.InputBody(err.Error()))
		return
	}

	created, err := h.serviceSvc.OrgSvc.Create(c.Request.Context(), o)
	if err != nil {
		h.respondError(c, err)
		return
	}

	h.respondJSON(c, http.StatusCreated, created)
}

func (h *Handler) GetOrganisation(c *gin.Context) {
	o, err := h.serviceSvc.OrgSvc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.respondError(c, err)
		return
	}

	h.respondJSON(c, http.StatusOK, o)
}

func (h *Handler) ListOrganisations(c *gin.Context) {
	orgs, err := h.serviceSvc.OrgSvc.List(c.Request.Context())
	if err != nil {
		h.respondError(c, err)
		return
	}

	h.respondJSON(c, http.StatusOK, orgs)
}

func (h *Handler) UpdateOrganisation(c *gin.Context) {
	var o model.Organisation
	if err := c.ShouldBindJSON(&o); err != nil {
		h.respondError(c, apperr.InputBody(err.Error()))
		return
	}

	updated, err := h.serviceSvc.OrgSvc.Update(c.Request.Context(), c.Param("id"), o)
	if err != nil {
		h.respondError(c, err)
		return
	}

	h.respondJSON(c, http.StatusOK, updated)
}

func (h *Handler) DeleteOrganisation(c *gin.Context) {
	if err := h.serviceSvc.OrgSvc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) OnboardOrganisation(c *gin.Context) {
	o, err := h.serviceSvc.OrgSvc.Onboard(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.respondError(c, err)
		return
	}

	h.respondJSON(c, http.StatusOK, o)
}
