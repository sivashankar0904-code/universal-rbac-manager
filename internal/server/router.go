package server

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host string
	Port string
}

type Server struct {
	cfg     *Config
	log     *slog.Logger
	engine  *gin.Engine
	handler *Handler
}

func New(cfg *Config, log *slog.Logger, handler *Handler) *Server {
	engine := gin.New()
	// No reverse proxy in front of this yet — disable trusting any
	// X-Forwarded-For header rather than defaulting to "trust everyone",
	// which lets a client spoof its own IP. Revisit if/when a real proxy
	// is added in front of URM (would need its actual IP/CIDR listed here).
	if err := engine.SetTrustedProxies(nil); err != nil {
		log.Error("failed to set trusted proxies", "err", err)
	}

	s := &Server{
		cfg:     cfg,
		log:     log,
		engine:  engine,
		handler: handler,
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
		s.log.Error("panic recovered",
			"err", recovered,
			"stack", string(debug.Stack()),
		)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (s *Server) Start() error {
	return s.engine.Run(s.cfg.Host + ":" + s.cfg.Port)
}
