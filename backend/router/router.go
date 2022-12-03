package router

import (
	"backend/config"
	h "backend/handler"
	m "backend/middleware"
	s "backend/services"
	"log"
	"net"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/sys/unix"
)

var (
	// HTTP addr and network
	addr    string
	network string

	// Internal fiber.App
	app *fiber.App
)

// init setup middlewares and routes.
func init() {
	var ok bool

	addr, ok = config.Get("HTTP_ADDR")
	if !ok {
		log.Fatalln("HTTP_ADDR env not found.")
	}

	network, ok = config.Get("HTTP_NETWORK")
	if !ok {
		log.Fatalln("HTTP_NETWORK env not found.")
	}

	// Internal fiber app
	app = fiber.New(fiber.Config{
		Network: network,
	})

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

// Start starts listening for connections.
func Start() {
	if network == "unix" {
		// Remove old unix socket
		os.Remove(addr)

		// Set unix mask for permissions so we can share the unix socket.
		old := unix.Umask(0)
		listener, err := net.Listen("unix", addr)
		unix.Umask(old)

		// Check if error
		if err != nil {
			log.Fatalln(err)
		}

		// Start listening
		log.Fatalln(app.Listener(listener))
	} else {
		// Start with tcp network
		log.Fatalln(app.Listen(addr))
	}
}
