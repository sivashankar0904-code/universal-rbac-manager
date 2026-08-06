package server

import "github.com/gin-gonic/gin"

func (s *Server) registerRoutes() {
	s.engine.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
