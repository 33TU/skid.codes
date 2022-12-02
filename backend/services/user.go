package services

import (
	"backend/claims"
	"backend/database"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	getUserTimeout    = time.Second * 5
	updateUserTimeout = time.Second * 5
	createUserTimeout = time.Second * 10
	findUserTimeout   = time.Second * 15
)

type CreateUserBody struct {
	Username string `json:"username" validate:"required,min=1,max=32"`
	Email    string `json:"email" validate:"required,email,min=6,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

type FindUserBody struct {
	Username string `json:"username" validate:"required,min=1,max=32"`
	Offset   int    `json:"offset" validate:"min=0"`
	Count    int    `json:"count" validate:"required,min=1,max=100"`
}

type UpdateUserBody struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=1,max=32"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email,min=6,max=255"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=8"`
}

type GetUserResult struct {
	Id       int        `json:"id"`
	Username string     `json:"username"`
	Role     string     `json:"role"`
	Pastes   int        `json:"pastes"`
	Online   *time.Time `json:"online"`
}

type FindUserResult struct {
	Id       int        `json:"id"`
	Username string     `json:"username"`
	Role     string     `json:"role"`
	Pastes   int        `json:"pastes"`
	Online   *time.Time `json:"online"`
	Count    int        `json:"-"`
}

type UpdateUserResult struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type CreateUserResult struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// GetUser gets user based on parameters.
func GetUser(username string) (*GetUserResult, error) {
	return database.SelectOne[GetUserResult](
		getUserTimeout,
		database.QueryUserGet,
		username,
	)
}

// FindPaste finds user based on parameters.
func FindUser(body *FindUserBody) (res []*FindUserResult, count int, err error) {
	res, err = database.Select[FindUserResult](
		findUserTimeout,
		database.QueryUserFind,
		body.Username, body.Offset, body.Count,
	)

	// Return if err
	if err != nil {
		return
	}

	// Get total count
	count = 0
	if len(res) != 0 {
		count = res[0].Count
	}

	// Return
	return
}

// UpdatePaste updates user's details.
func UpdateUser(body *UpdateUserBody, session *claims.AuthClaims) (res *UpdateUserResult, err error) {
	res, err = database.SelectOne[UpdateUserResult](
		updateUserTimeout,
		database.QueryUserUpdate,
		session.UserID, body.Username, body.Email, body.Password,
	)

	if err == pgx.ErrNoRows {
		err = ErrNotFound
	}

	return
}

// CreatePaste creates new user.
func CreateUser(body *CreateUserBody) (res *CreateUserResult, err error) {
	res, err = database.SelectOne[CreateUserResult](
		createUserTimeout,
		database.QueryUserCreate,
		body.Username, body.Email, body.Password,
	)

	if err == pgx.ErrNoRows {
		err = ErrNotFound
	}

	return
}
