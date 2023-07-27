package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateLoginResponse(t *testing.T) {
	id := RandomNumString(20)
	username := RandomString(10)
	loginResponse, token, err := GenerateLoginResponse(id, username)
	assert.NoError(t, err)
	assert.Equal(t, loginResponse.Username, username)
	assert.Equal(t, loginResponse.AccessToken, *token.AccessToken.Token)
	assert.Equal(t, loginResponse.RefreshToken, *token.RefreshToken.Token)
	assert.Equal(t, id, token.AccessToken.UserId)
	assert.Equal(t, username, token.AccessToken.Username)
	assert.Equal(t, id, token.RefreshToken.UserId)
	assert.Equal(t, username, token.RefreshToken.Username)
	assert.NotZero(t, token.AccessToken.ExpiresAt)
	assert.NotZero(t, token, token.RefreshToken.ExpiresAt)
	assert.Equal(t, token.AccessToken.TokenID, token.RefreshToken.TokenID)
}
