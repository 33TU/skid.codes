package router

import (
	"backend/config"
	e "backend/errors"
	h "backend/handler"
	m "backend/middleware"
	s "backend/services"
	"log"
	"net"
	"os"
	"time"

	json "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/sys/unix"
)

// Config
var (
	// HTTP addr and network
	addr       string
	network    string
	trustProxy string

	// Internal fiber.App
	app *fiber.App
)

// Limiters
var (
	loginLimiter       = m.IPLimiter(10, time.Minute, false)
	refreshLimiter     = m.IPLimiter(10, time.Minute, false)
	userCreateLimiter  = m.IPLimiter(1, time.Minute*5, true)
	pasteCreateLimiter = m.AuthLimiter(10, time.Minute, false)
	pasteFetchLimiter  = m.IPLimiter(30, time.Minute, false)
	pasteFindLimiter   = m.IPLimiter(30, time.Minute, false)
	userLimiter        = m.AuthLimiter(100, time.Minute, false)
)

func initConfig() {
	var ok bool

	addr, ok = config.Get("HTTP_ADDR")
	if !ok {
		log.Fatalln("HTTP_ADDR env not found.")
	}

	network, ok = config.Get("HTTP_NETWORK")
	if !ok {
		log.Fatalln("HTTP_NETWORK env not found.")
	}

	trustProxy, ok = config.Get("HTTP_TRUST_PROXY")
	if !ok {
		log.Fatalln("HTTP_TRUST_PROXY env not found.")
	}
}

// init setup middlewares and routes.
func init() {
	// Init config
	initConfig()

	// Internal fiber app
	app = fiber.New(fiber.Config{
		Network:                 network,
		EnableTrustedProxyCheck: trustProxy == "true",
		ErrorHandler:            e.SendError,

		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Middlewares
	app.Use(cors.New())
	app.Use(m.Lookup())

	// Api
	api := app.Group("/api", logger.New())

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", loginLimiter, m.Validate[s.LoginRequest](), h.LoginHandler)
	auth.Get("/refresh", m.Refresh(), refreshLimiter, h.RefreshHandler)
	auth.Get("/logout", m.Auth(), userLimiter, h.LogoutHandler)

	// User
	user := api.Group("/user")
	user.Get("/:username", h.GetUserHandler)
	user.Post("/create", userCreateLimiter, m.Validate[s.CreateUserRequest](), h.CreateUserHandler)
	user.Post("/find", m.Validate[s.FindUserRequest](), h.FindUserHandler)
	user.Use(m.Auth(), userLimiter) // Protected paths
	user.Post("/update", m.Validate[s.UpdateUserRequest](), h.UpdateUserHandler)

	// Paste
	paste := api.Group("/paste")
	paste.Post("/fetch", pasteFetchLimiter, m.Validate[s.FetchPasteRequest](), h.FetchPasteHandler)
	paste.Post("/find", pasteFindLimiter, m.Validate[s.FindPasteRequest](), h.FindPasteHandler)
	paste.Use(m.Auth(), userLimiter) // Protected paths
	paste.Post("/ufetch", m.Validate[s.FetchPasteRequest](), h.FetchUserPasteHandler)
	paste.Post("/ufind", m.Validate[s.FindPasteRequest](), h.FindPasteHandler)
	paste.Post("/create", pasteCreateLimiter, m.Validate[s.CreatePasteRequest](), h.CreatePasteHandler)
	paste.Post("/update", m.Validate[s.UpdatePasteRequest](), h.UpdatePasteHandler)
	paste.Post("/delete", m.Validate[s.DeletePasteRequest](), h.DeletePasteHandler)

	// Session
	session := api.Group("/session")
	session.Use(m.Auth(), userLimiter) // Protected paths
	session.Post("/find", m.Validate[s.FindSessionRequest](), h.FindSessionHandler)
	session.Post("/revoke", m.Validate[s.RevokeSessionRequest](), h.RevokeSessionHandler)
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
