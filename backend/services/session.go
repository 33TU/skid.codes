package services

import (
	"backend/claims"
	"backend/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

const (
	findSessionTimeout   = time.Second * 5
	revokeSessionTimeout = time.Second * 5
)

var (
	ErrSessionNotFound = fiber.NewError(404, "session not found")
)

type FindSessionBody struct {
	Offset      int  `json:"offset" validate:"min=0"`
	Count       int  `json:"count" validate:"required,min=1,max=100"`
	ShowRevoked bool `json:"showRevoked"`
}

type RevokeSessionBody struct {
	SessionID int `json:"sid" validate:"required"`
}

type FindSessionResult struct {
	ID      int       `json:"id"`
	Country string    `json:"country"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Revoked bool      `json:"revoked"`
	Count   int       `json:"-"`
}

type RevokeSessionResult struct {
	ID      int       `json:"id"`
	Country string    `json:"country"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Revoked bool      `json:"revoked"`
}

// FindSession finds session based on parameters.
func FindSession(body *FindSessionBody, session *claims.AuthClaims) (res []*FindSessionResult, count int, err error) {
	res, err = database.Select[FindSessionResult](
		findSessionTimeout,
		database.QuerySessionFind,
		session.UserID, body.ShowRevoked, body.Offset, body.Count,
	)

	// Check for error
	if err != nil {
		return
	}

	// Get total count
	count = 0
	if len(res) != 0 {
		count = res[0].Count
	}

	// Return result
	return
}

// RevokeSession revokes sessions from refreshing JWT-token.
func RevokeSession(body *RevokeSessionBody, session *claims.AuthClaims) (res *RevokeSessionResult, cookie *fiber.Cookie, err error) {
	res, err = database.SelectOne[RevokeSessionResult](
		revokeSessionTimeout,
		database.QuerySessionRevoke,
		session.UserID, body.SessionID,
	)

	if err == pgx.ErrNoRows {
		err = ErrSessionNotFound
		return
	} else if err != nil {
		return
	}

	// Set expired cookie for refresh_token
	cookie = &fiber.Cookie{
		Name:     "refresh_token",
		Value:    "expired",
		Expires:  time.Time{},
		SameSite: "Strict",
		Path:     "/api/auth/refresh",
		HTTPOnly: true,
		Secure:   true,
	}

	// Return result
	return
}
