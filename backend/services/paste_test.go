package services_test

import (
	"backend/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchPaste(t *testing.T) {
	// Login to get auth
	res, _, _, err := getLogin()
	assert.NoError(t, err)

	// Params
	content := "int main() {}"
	title := "A C program"
	password := "password123"
	private := true
	unlisted := true

	// Create paste
	cres, err := services.CreatePaste(&services.CreatePasteRequest{
		Language: "C",
		Content:  content,
		Title:    &title,
		Password: &password,
		Private:  private,
		Unlisted: unlisted,
	}, res)

	assert.NoError(t, err)

	// Fetch the paste
	fres, err := services.FetchPaste(&services.FetchPasteRequest{
		ID:       cres.ID,
		Password: &password,
	}, &res.UserID)

	assert.NoError(t, err)
	assert.Equal(t, *fres.Content, content)
	assert.Equal(t, *fres.Title, title)
	assert.Equal(t, fres.Private, private)
	assert.Equal(t, fres.Unlisted, unlisted)
}

func TestFindPaste(t *testing.T) {
	language := "C"

	// Find a paste
	res, count, err := services.FindPaste(&services.FindPasteRequest{
		Language: &language,
		Offset:   0,
		Count:    1,
	}, nil)

	// Assert
	assert.NoError(t, err)
	assert.Greater(t, count, 0)
	assert.Equal(t, res[0].Language.Name, language)
}

func TestUpdatePaste(t *testing.T) {
	// Login to get auth
	res, _, _, err := getLogin()
	assert.NoError(t, err)

	// Find one paste
	fres, count, err := services.FindPaste(&services.FindPasteRequest{
		UserID: &res.UserID,
		Offset: 0,
		Count:  1,
	}, &res.UserID)

	assert.NoError(t, err)
	assert.Greater(t, count, 0)

	// Update paste
	newTitle := "updated"
	content := "int main() {}"
	password := "password123"

	ures, err := services.UpdatePaste(&services.UpdatePasteRequest{
		ID:       fres[0].ID,
		Title:    &newTitle,
		Content:  &content,
		Password: &password,
	}, res)

	assert.NoError(t, err)
	assert.Equal(t, ures.ID, fres[0].ID)
}

func TestCreatePaste(t *testing.T) {
	// Login to get auth
	res, _, _, err := getLogin()
	assert.NoError(t, err)

	// Params
	content := "TestCreatePaste"
	title := "TestCreatePaste"
	password := "TestCreatePaste"
	private := false
	unlisted := false

	// Create paste
	cres, err := services.CreatePaste(&services.CreatePasteRequest{
		Language: "C",
		Content:  content,
		Title:    &title,
		Password: &password,
		Private:  private,
		Unlisted: unlisted,
	}, res)

	assert.NoError(t, err)
	assert.Greater(t, len(cres.ID), 0)
}

func TestDeletePaste(t *testing.T) {
	// Login to get auth
	res, _, _, err := getLogin()
	assert.NoError(t, err)

	// Find one paste
	fres, count, err := services.FindPaste(&services.FindPasteRequest{
		UserID: &res.UserID,
		Offset: 0,
		Count:  1,
	}, &res.UserID)

	assert.NoError(t, err)
	assert.Greater(t, count, 0)

	// Delete the paste
	dres, err := services.DeletePaste(&services.DeletePasteRequest{
		ID: fres[0].ID,
	}, res)

	assert.NoError(t, err)
	assert.Equal(t, dres.ID, fres[0].ID)
}
