package services_test

import (
	"backend/claims"
	"backend/services"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/ip2location/ip2location-go/v9"
	"github.com/stretchr/testify/assert"
)

// NOTE: these tests assume db have user called admin
// and it's credentials must be admin:dev12345

const (
	username = "admin"
	password = "dev12345"
)

// getLogin gets login res for user.
func getLogin() (*claims.AuthClaims, string, *fiber.Cookie, error) {
	ipGeo := &ip2location.IP2Locationrecord{}

	return services.Login(&services.LoginBody{
		Username: username,
		Password: password,
		Email:    "",
	}, ipGeo, "127.0.0.1")
}

func TestLogin(t *testing.T) {
	res, _, _, err := getLogin()

	assert.NoError(t, err)
	assert.Equal(t, res.Username, username)
}

func TestRefresh(t *testing.T) {
	// Login
	res, _, _, err := getLogin()
	assert.NoError(t, err)

	// Refresh
	res, _, _, err = services.Refresh(&claims.RefreshClaims{
		SessionID: res.SessionID,
	}, "127.0.0.1")

	assert.NoError(t, err)
	assert.Equal(t, res.Username, username)
}
