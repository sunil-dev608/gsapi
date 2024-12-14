package server

import (
	"context"

	"gsapi/config"
	"gsapi/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
	cfg  *config.Config
}

func New(cfg *config.Config) *Server {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	packsHandler := handlers.NewPacksHandler(cfg.PackSizes)

	// Setup routes
	v1 := e.Group("/api/v1")
	v1.POST("/packs-for-items", packsHandler.GetPacksForItems)

	return &Server{
		echo: e,
		cfg:  cfg,
	}
}

func (s *Server) Start() error {
	return s.echo.Start(":8080")
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
