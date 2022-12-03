package services

import (
	"backend/claims"
	"backend/database"
	"errors"
	"time"

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
	ErrPasteNotFound        = errors.New("paste not found")
	ErrPasteCreate          = errors.New("failed to create paste")
	ErrPasteUpdateNoContent = errors.New("password can not be set without content")
)

type FetchPasteBody struct {
	ID       string  `json:"id" validate:"required"`
	Password *string `json:"password"`
}

type FindPasteBody struct {
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

type UpdatePasteBody struct {
	ID       string  `json:"id" validate:"required"`
	Language *string `json:"language"`
	Title    *string `json:"title"`
	Private  *bool   `json:"private"`
	Unlisted *bool   `json:"unlisted"`
	Password *string `json:"password,omitempty"`
	Content  *string `json:"content,omitempty"`
}

type DeletePasteBody struct {
	ID string `json:"id" validate:"required"`
}

type CreatePasteBody struct {
	Language string  `json:"language" validate:"required"`
	Content  string  `json:"content" validate:"required"`
	Title    *string `json:"title"`
	Password *string `json:"password"`
	Private  bool    `json:"private"`
	Unlisted bool    `json:"unlisted"`
}

type FetchPasteResult struct {
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

type FindPasteResult struct {
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

type UpdatePasteResult struct {
	ID string `json:"id"`
}

type DeletePasteResult struct {
	ID string `json:"id"`
}

type CreatePasteResult struct {
	ID string `json:"id"`
}

// FetchPaste fetches paste based on parameters. Caller is user id.
func FetchPaste(body *FetchPasteBody, caller *int) (res *FetchPasteResult, err error) {
	res, err = database.SelectOne[FetchPasteResult](
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
func FindPaste(body *FindPasteBody, caller *int) (res []*FindPasteResult, count int, err error) {
	res, err = database.Select[FindPasteResult](
		findPasteTimeout,
		database.QueryPasteFind,
		body.UserID, body.Username, body.Language, body.Title, body.Password, body.Content,
		body.CreatedBegin, body.CreatedEnd,
		body.Private, body.Unlisted,
		caller, // Caller user_id
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
func UpdatePaste(body *UpdatePasteBody, session *claims.AuthClaims) (res *UpdatePasteResult, err error) {
	if body.Password != nil && body.Content == nil {
		err = ErrPasteUpdateNoContent
		return
	}

	res, err = database.SelectOne[UpdatePasteResult](
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
func DeletePaste(body *DeletePasteBody, session *claims.AuthClaims) (res *DeletePasteResult, err error) {
	res, err = database.SelectOne[DeletePasteResult](
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
func CreatePaste(body *CreatePasteBody, session *claims.AuthClaims) (res *CreatePasteResult, err error) {
	res, err = database.SelectOne[CreatePasteResult](
		createPasteTimeout,
		database.QueryPasteCreate,
		session.UserID, body.Language, body.Title, body.Password, body.Private, body.Unlisted, body.Content,
	)

	if err == pgx.ErrNoRows {
		err = ErrPasteCreate
	}

	return
}
