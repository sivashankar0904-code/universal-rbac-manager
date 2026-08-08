package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerRoutes() {
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{}, "status": http.StatusOK})
	})

	s.engine.POST("/organisations", s.handler.CreateOrganisation)
	s.engine.GET("/organisations", s.handler.ListOrganisations)
	s.engine.GET("/organisations/:id", s.handler.GetOrganisation)
	s.engine.PUT("/organisations/:id", s.handler.UpdateOrganisation)
	s.engine.DELETE("/organisations/:id", s.handler.DeleteOrganisation)
	s.engine.POST("/organisations/:id/onboard", s.handler.OnboardOrganisation)
}
