package services

import (
	"backend/claims"
	"backend/database"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	findSessionTimeout   = time.Second * 5
	revokeSessionTimeout = time.Second * 5
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
func RevokeSession(body *RevokeSessionBody, session *claims.AuthClaims) (res *RevokeSessionResult, err error) {
	res, err = database.SelectOne[RevokeSessionResult](
		revokeSessionTimeout,
		database.QuerySessionRevoke,
		session.UserID, body.SessionID,
	)

	if err == pgx.ErrNoRows {
		err = ErrNotFound
		return
	} else if err != nil {
		return
	}

	// Return result
	return
}
