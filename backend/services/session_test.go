package services_test

import (
	"backend/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSession(t *testing.T) {
	// Login to get session id
	res, _, _, err := getLogin()
	assert.NoError(t, err)

	// Find all sessions for user:
	sessions, count, err := services.FindSession(&services.FindSessionBody{
		Offset: 0,
		Count:  100,
	}, res)

	// Check that all ok
	assert.NoError(t, err)
	assert.Greater(t, count, 0)
	assert.Equal(t, len(sessions), count)
}

func TestRevokeSession(t *testing.T) {
	// Login to get session id
	res, _, _, err := getLogin()
	assert.NoError(t, err)

	// Try to revoke
	rev, err := services.RevokeSession(&services.RevokeSessionBody{
		SessionId: res.SessionID,
	}, res)

	// Check that all ok
	assert.NoError(t, err)
	assert.Equal(t, rev.Id, res.SessionID)
	assert.Equal(t, rev.Revoked, true)
}
