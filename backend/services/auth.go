package services

import (
	"backend/claims"
	"backend/database"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ip2location/ip2location-go/v9"
	"github.com/jackc/pgx/v5"
)

const (
	loginQueryTimeout   = time.Second * 10
	refreshQueryTimeout = time.Second * 10
)

var (
	ErrLoginInvalidCredentials = errors.New("missing credentials")
	ErrRefreshNoPermission     = errors.New("no permission to refresh this token")
	ErrRefreshTokenInvalid     = errors.New("invalid refresh token")
	ErrRefreshTokenExpired     = errors.New("expired refresh token")
)

type LoginBody struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=1,max=32"`
	Email    string `json:"email,omitempty" validate:"omitempty,email,min=6,max=255"`
	Password string `json:"password" validate:"required"`
}

type RefreshBody struct {
	Token string `json:"token" validate:"required"`
}

// LoginHandler logins and creates session into db returning JWT token.
func Login(body *LoginBody, lookup *ip2location.IP2Locationrecord, ip string) (res *claims.AuthClaims, auth string, cookie *fiber.Cookie, err error) {
	if body.Email == "" && body.Username == "" {
		err = ErrLoginInvalidCredentials
		return
	}

	// Use email instead of username
	useEmail := body.Email != ""
	country := lookup.Country_short

	// Query loginResponse from db
	res, err = database.SelectOne[claims.AuthClaims](
		loginQueryTimeout,                                               // Timeout
		database.QueryAuthLogin,                                         // Query
		useEmail, body.Username, body.Email, body.Password, ip, country, // Params
	)

	// Check if no results or err
	if err == pgx.ErrNoRows {
		err = ErrNotFound
		return
	} else if err != nil {
		return
	}

	// Create auth and refesh token
	var refreshExpires time.Time
	var refresh string
	auth, refresh, refreshExpires, err = claims.CreateTokens(res, res.Refresh())
	if err != nil {
		return
	}

	// Store refresh token as cookie
	cookie = &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Expires:  refreshExpires,
		SameSite: "Strict",
		Path:     "/api/auth/refresh",
		HTTPOnly: true,
	}

	// Return response
	return
}

// Refresh refreshes the JWT token.
func Refresh(session *claims.RefreshClaims, ip string) (res *claims.AuthClaims, auth string, cookie *fiber.Cookie, err error) {
	// Query loginResponse from db
	res, err = database.SelectOne[claims.AuthClaims](
		refreshQueryTimeout,       // Timeout
		database.QueryAuthRefresh, // Query
		session.SessionID, ip,     // Params
	)

	// Check if no results or err
	if err == pgx.ErrNoRows {
		err = ErrNotFound
		return
	} else if err != nil {
		return
	}

	// Update auth and refesh tokens
	var refreshExpires time.Time
	var refresh string
	auth, refresh, refreshExpires, err = claims.CreateTokens(res, res.Refresh())
	if err != nil {
		return
	}

	// Store refresh token as cookie
	cookie = &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Expires:  refreshExpires,
		SameSite: "Strict",
		Path:     "/api/auth/refresh",
		HTTPOnly: true,
	}

	// Return response
	return
}
