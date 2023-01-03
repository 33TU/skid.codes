package services

import (
	"backend/claims"
	"backend/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
)

const (
	getUserTimeout    = time.Second * 5
	updateUserTimeout = time.Second * 5
	createUserTimeout = time.Second * 10
	findUserTimeout   = time.Second * 15
)

var (
	ErrUserNotFound     = fiber.NewError(404, "user not found")
	ErrUserAlreadyExist = fiber.NewError(400, "user or email aready exist")
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=1,max=32"`
	Email    string `json:"email" validate:"required,email,min=6,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

type FindUserRequest struct {
	Username string `json:"username" validate:"required,min=1,max=32"`
	Offset   int    `json:"offset" validate:"min=0"`
	Count    int    `json:"count" validate:"required,min=1,max=100"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=1,max=32"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email,min=6,max=255"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=8"`
}

type GetUserResponse struct {
	ID       int        `json:"id"`
	Username string     `json:"username"`
	Role     string     `json:"role"`
	Pastes   int        `json:"pastes"`
	Online   *time.Time `json:"online"`
}

type FindUserResponse struct {
	ID       int        `json:"id"`
	Username string     `json:"username"`
	Role     string     `json:"role"`
	Pastes   int        `json:"pastes"`
	Online   *time.Time `json:"online"`
	Count    int        `json:"-"`
}

type UpdateUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type CreateUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// GetUser gets user based on parameters.
func GetUser(username string) (res *GetUserResponse, err error) {
	res, err = database.SelectOne[GetUserResponse](
		getUserTimeout,
		database.QueryUserGet,
		username,
	)

	if err == pgx.ErrNoRows {
		err = ErrUserNotFound
	}

	return
}

// FindPaste finds user based on parameters.
func FindUser(body *FindUserRequest) (res []*FindUserResponse, count int, err error) {
	res, err = database.Select[FindUserResponse](
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
func UpdateUser(body *UpdateUserRequest, session *claims.AuthClaims) (res *UpdateUserResponse, err error) {
	res, err = database.SelectOne[UpdateUserResponse](
		updateUserTimeout,
		database.QueryUserUpdate,
		session.UserID, body.Username, body.Email, body.Password,
	)

	if err == pgx.ErrNoRows {
		err = ErrUserNotFound
	}

	return
}

// CreatePaste creates new user.
func CreateUser(body *CreateUserRequest) (res *CreateUserResponse, err error) {
	res, err = database.SelectOne[CreateUserResponse](
		createUserTimeout,
		database.QueryUserCreate,
		body.Username, body.Email, body.Password,
	)

	if err == pgx.ErrNoRows {
		err = ErrUserNotFound
	}

	if database.ErrorCodeEquals(err, pgerrcode.UniqueViolation) {
		err = ErrUserAlreadyExist
	}

	return
}
