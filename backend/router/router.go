package router

import (
	h "backend/handler"
	m "backend/middleware"
	s "backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	// Internal fiber.App
	app *fiber.App
)

// init setup middlewares and routes.
func init() {
	app = fiber.New()

	// Middlewares
	app.Use(cors.New())
	app.Use(m.Lookup())

	// Api
	api := app.Group("/api", logger.New())

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", m.Validate[s.LoginBody](), h.LoginHandler)
	auth.Use(m.Refresh()) // Protected paths
	auth.Get("/refresh", h.RefreshHandler)

	// User
	user := api.Group("/user")
	user.Get("/:username", h.GetUserHandler)
	user.Post("/create", m.Validate[s.CreateUserBody](), h.CreateUserHandler)
	user.Post("/find", m.Validate[s.FindUserBody](), h.FindUserHandler)
	user.Use(m.Auth()) // Protected paths
	user.Post("/update", m.Validate[s.UpdateUserBody](), h.UpdateUserHandler)

	// Paste
	paste := api.Group("/paste")
	paste.Post("/fetch", m.Validate[s.FetchPasteBody](), h.FetchPasteHandler)
	paste.Post("/find", m.Validate[s.FindPasteBody](), h.FindPasteHandler)
	paste.Use(m.Auth()) // Protected paths
	paste.Post("/ufetch", m.Validate[s.FetchPasteBody](), h.FetchUserPasteHandler)
	paste.Post("/ufind", m.Validate[s.FindPasteBody](), h.FindPasteHandler)
	paste.Post("/create", m.Validate[s.CreatePasteBody](), h.CreatePasteHandler)
	paste.Post("/update", m.Validate[s.UpdatePasteBody](), h.UpdatePasteHandler)
	paste.Post("/delete", m.Validate[s.DeletePasteBody](), h.DeletePasteHandler)

	// Session
	session := api.Group("/session")
	session.Use(m.Auth()) // Protected paths
	session.Post("/find", m.Validate[s.FindSessionBody](), h.FindSessionHandler)
	session.Post("/revoke", m.Validate[s.RevokeSessionBody](), h.RevokeSessionHandler)
}

// App returns the internal *fiber.App.
func App() *fiber.App {
	return app
}
