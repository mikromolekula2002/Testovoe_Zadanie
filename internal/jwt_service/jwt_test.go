package jwt_service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAccessToken(t *testing.T) {
	jwtTest := InitJWT()

	token, err := jwtTest.GenerateAccessToken("testUser", "192.896.2.4", []byte("secretKeyTest"), 15)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}
