package server

import "github.com/gin-gonic/gin"

type Config struct {
	Host string
	Port string
}

type Server struct {
	cfg    *Config
	engine *gin.Engine
}

func New(cfg *Config) *Server {
	s := &Server{
		cfg:    cfg,
		engine: gin.Default(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) Start() error {
	return s.engine.Run(s.cfg.Host + ":" + s.cfg.Port)
}
