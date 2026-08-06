package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host string
	Port string
}

type Server struct {
	cfg    *Config
	log    *slog.Logger
	engine *gin.Engine
}

func New(cfg *Config, log *slog.Logger) *Server {
	engine := gin.New()

	s := &Server{
		cfg:    cfg,
		log:    log,
		engine: engine,
	}

	engine.Use(s.requestLogger())
	engine.Use(gin.CustomRecoveryWithWriter(nil, s.recovery()))

	s.registerRoutes()
	return s
}

func (s *Server) requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		s.log.Info("request",
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"latency", time.Since(start),
		)
	}
}

func (s *Server) recovery() gin.RecoveryFunc {
	return func(c *gin.Context, recovered any) {
		s.log.Error("panic recovered", "err", recovered)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (s *Server) Start() error {
	return s.engine.Run(s.cfg.Host + ":" + s.cfg.Port)
}
