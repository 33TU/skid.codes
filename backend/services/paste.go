package services

import (
	"backend/claims"
	"backend/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

const (
	findPasteTimeout   = time.Second * 5
	fetchPasteTimeout  = time.Second * 5
	createPasteTimeout = time.Second * 5
	deletePasteTimeout = time.Second * 5
	updatePasteTimeout = time.Second * 5
)

var (
	ErrPasteNotFound        = fiber.NewError(404, "paste not found")
	ErrPasteCreate          = fiber.NewError(400, "failed to create paste")
	ErrPasteUpdateNoContent = fiber.NewError(400, "password can not be set without content")
)

type FetchPasteRequest struct {
	ID       string  `json:"id" validate:"required"`
	Password *string `json:"password"`
}

type FindPasteRequest struct {
	UserID   *int    `json:"uid"`
	Username *string `json:"username"`
	Language *string `json:"language"`
	Title    *string `json:"title"`
	Content  *string `json:"content"`

	// Only works with authenticated user
	Private  *bool `json:"private"`
	Unlisted *bool `json:"unlisted"`
	Password *bool `json:"password"`

	// Created and end
	CreatedBegin *time.Time `json:"createdBegin"`
	CreatedEnd   *time.Time `json:"createdEnd"`

	Offset int `json:"offset" validate:"min=0"`
	Count  int `json:"count" validate:"required,min=1,max=100"`
}

type UpdatePasteRequest struct {
	ID       string  `json:"id" validate:"required"`
	Language *string `json:"language"`
	Title    *string `json:"title"`
	Private  *bool   `json:"private"`
	Unlisted *bool   `json:"unlisted"`
	Password *string `json:"password,omitempty"`
	Content  *string `json:"content,omitempty"`
}

type DeletePasteRequest struct {
	ID string `json:"id" validate:"required"`
}

type CreatePasteRequest struct {
	Language string  `json:"language" validate:"required"`
	Content  string  `json:"content" validate:"required"`
	Title    *string `json:"title"`
	Password *string `json:"password"`
	Private  bool    `json:"private"`
	Unlisted bool    `json:"unlisted"`
}

type FetchPasteResponse struct {
	ID       string    `json:"id"`
	UID      int       `json:"uid"`
	Title    *string   `json:"title"`
	Private  bool      `json:"private"`
	Unlisted bool      `json:"unlisted"`
	Created  time.Time `json:"created"`
	Password bool      `json:"password"`
	Content  *string   `json:"content"`

	Language struct {
		Name string   `json:"name"`
		Mode string   `json:"mode"`
		Mime string   `json:"mime"`
		Ext  []string `json:"ext"`
	} `json:"language"`
}

type FindPasteResponse struct {
	ID       string    `json:"id"`
	UID      int       `json:"uid"`
	Count    int       `json:"-"`
	Title    *string   `json:"title"`
	Content  *string   `json:"content"`
	Created  time.Time `json:"created"`
	Private  bool      `json:"private"`
	Language struct {
		Ext  []string `json:"ext"`
		Mime string   `json:"mime"`
		Mode string   `json:"mode"`
		Name string   `json:"name"`
	} `json:"language"`
	Password bool   `json:"password"`
	Unlisted bool   `json:"unlisted"`
	Username string `json:"username"`
}

type UpdatePasteResponse struct {
	ID string `json:"id"`
}

type DeletePasteResponse struct {
	ID string `json:"id"`
}

type CreatePasteResponse struct {
	ID string `json:"id"`
}

// FetchPaste fetches paste based on parameters. Caller is user id.
func FetchPaste(body *FetchPasteRequest, caller *int) (res *FetchPasteResponse, err error) {
	res, err = database.SelectOne[FetchPasteResponse](
		fetchPasteTimeout,
		database.QueryPasteFetch,
		body.ID, caller, body.Password,
	)

	if err == pgx.ErrNoRows {
		err = ErrPasteNotFound
	}

	return
}

// FindPaste finds paste based on parameters. Caller is user id.
func FindPaste(body *FindPasteRequest, caller *int) (res []*FindPasteResponse, count int, err error) {
	res, err = database.Select[FindPasteResponse](
		findPasteTimeout,
		database.QueryPasteFind,
		body.UserID, body.Username, body.Language, body.Title,
		body.Password, body.Content,
		body.CreatedBegin, body.CreatedEnd,
		body.Private, body.Unlisted, caller,
		body.Offset, body.Count,
	)

	if err == pgx.ErrNoRows {
		err = ErrPasteNotFound
	}

	// Get total count
	count = 0
	if len(res) != 0 {
		count = res[0].Count
	}

	return
}

// UpdatePaste updates paste's details.
func UpdatePaste(body *UpdatePasteRequest, session *claims.AuthClaims) (res *UpdatePasteResponse, err error) {
	if body.Password != nil && body.Content == nil {
		err = ErrPasteUpdateNoContent
		return
	}

	res, err = database.SelectOne[UpdatePasteResponse](
		updatePasteTimeout,
		database.QueryPasteUpdate,
		body.ID, session.UserID, body.Language, body.Title, body.Private, body.Unlisted, body.Password, body.Content,
	)

	if err == pgx.ErrNoRows {
		err = ErrPasteNotFound
	}

	return
}

// DeletePaste deletes paste.
func DeletePaste(body *DeletePasteRequest, session *claims.AuthClaims) (res *DeletePasteResponse, err error) {
	res, err = database.SelectOne[DeletePasteResponse](
		deletePasteTimeout,
		database.QueryPasteDelete,
		session.UserID, body.ID,
	)

	if err == pgx.ErrNoRows {
		err = ErrPasteNotFound
	}

	return
}

// CreatePaste creates new paste.
func CreatePaste(body *CreatePasteRequest, session *claims.AuthClaims) (res *CreatePasteResponse, err error) {
	res, err = database.SelectOne[CreatePasteResponse](
		createPasteTimeout,
		database.QueryPasteCreate,
		session.UserID, body.Language, body.Title, body.Password, body.Private, body.Unlisted, body.Content,
	)

	if err == pgx.ErrNoRows {
		err = ErrPasteCreate
	}

	return
}
