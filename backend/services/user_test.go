package services_test

import (
	"backend/services"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	res, err := services.GetUser(username)
	assert.NoError(t, err)
	assert.Equal(t, res.Username, username)
}

func TestFindUser(t *testing.T) {
	res, count, err := services.FindUser(&services.FindUserRequest{
		Username: "admin",
		Offset:   0,
		Count:    1,
	})

	assert.NoError(t, err)
	assert.Greater(t, count, 0)
	assert.Equal(t, res[0].Username, username)
}

func TestUpdateUser(t *testing.T) {
	// Login to get auth
	res, _, _, err := getLogin()
	assert.NoError(t, err)

	// Update user
	email := fmt.Sprint("admin", rand.Intn(123), "@skid.codes")
	ures, err := services.UpdateUser(&services.UpdateUserRequest{
		Email: &email,
	}, res)

	// Check that all ok
	assert.NoError(t, err)
	assert.Equal(t, ures.Email, email)
}

func TestCreateUser(t *testing.T) {
	// Try recreating same account
	res, err := services.CreateUser(&services.CreateUserRequest{
		Username: username,
		Email:    "admin@skid.tools",
		Password: "dev12345",
	})

	// This should fail
	assert.Error(t, err)
	assert.Nil(t, res)
}
