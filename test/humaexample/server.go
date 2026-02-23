package main

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Server wraps the Fiber app and the Huma API instance.
type Server struct {
	app *fiber.App
	api huma.API
}

// newServer creates a Server, registers all routes, and returns it ready to listen.
func newServer(db *gorm.DB) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	humaConfig := huma.DefaultConfig("Go PageToken Example API", "1.0.0")
	api := humafiber.New(app, humaConfig)

	s := &Server{app: app, api: api}
	registerRoutes(s.api, db)
	return s
}

// Listen starts the HTTP server on the given address (e.g. ":8080").
func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

// Shutdown gracefully drains in-flight requests.
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
